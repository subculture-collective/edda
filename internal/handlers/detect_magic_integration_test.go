//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

func TestIntegrationDetectMagicMagicInitiate(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	creation, err := srd521.SupportedBuildForConcept("Human Cleric", "Acolyte")
	must(err)
	stats, err := creation.Stats()
	must(err)
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Detect Magic fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	must(err)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
	var location uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name)VALUES($1,'Vault')RETURNING id`, c.ID).Scan(&location))
	pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: stats, Hp: 10, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
	must(err)
	campaign, actor := dbutil.FromPgtype(c.ID), dbutil.FromPgtype(pc.ID)
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, location, pc.ID)
	must(err)
	facts := func(distance int, barrier string, visible bool, origin, school string, hidden bool) string {
		return fmt.Sprintf(`{"srd521_magic":{"profile":"detect-magic-v1","actor_ids":["%s"],"location_id":"%s","hidden":%t,"distance_feet":%d,"barrier":"%s","visible":%t,"origin":"%s","school":"%s"}}`, actor, location, hidden, distance, barrier, visible, origin, school)
	}
	npc := func(name, properties string) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id,properties)VALUES($1,$2,5,$3,$4)RETURNING id`, c.ID, name, location, properties).Scan(&id))
		return id
	}
	ward := npc("Warded door", facts(20, "none", true, "spell", "abjuration", false))
	known := []uuid.UUID{
		ward,
		npc("Lead casket", facts(10, "lead-thin-sheet", true, "spell", "necromancy", false)),
		npc("Distant idol", facts(40, "none", true, "spell", "evocation", false)),
		npc("Veiled watcher", facts(10, "none", false, "spell", "illusion", false)),
		npc("Secret sigil", facts(5, "none", true, "spell", "conjuration", true)),
	}
	// Never involved in play: not player-known, so never sensed or disclosed.
	npc("Unmet lich", facts(5, "none", true, "spell", "transmutation", false))
	dice := &lifecycleIntegrationDice{}
	newEngine := func() *engine.Engine {
		e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(dice))
		must(err)
		return e
	}
	prepare := func(kind, entry string, slot int) *game.ActionPreview {
		t.Helper()
		c, p, err := game.ReadActionSnapshot(ctx, testPool, campaign, actor, testUserID, false)
		must(err)
		a, err := game.PrepareStoredAction(ctx, testPool, c, p, game.ActionIntent{ActorID: actor, Kind: kind, EntryID: entry, SlotLevel: slot}, false)
		must(err)
		return a.Preview
	}
	requestFor := func(p *game.ActionPreview) engine.PreparedAction {
		t.Helper()
		if !p.Available {
			t.Fatalf("unavailable: %+v", p)
		}
		r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent}
		if p.Intent.Kind == srd521.SpellCast {
			r.TargetID = actor
		}
		return r
	}
	confirm := func(p *game.ActionPreview) (engine.PreparedAction, *engine.TurnResult) {
		t.Helper()
		r := requestFor(p)
		result, err := newEngine().ConfirmAction(ctx, campaign, r)
		must(err)
		return r, result
	}
	read := func() *srd521.CharacterSheet {
		t.Helper()
		p, err := q.GetPlayerCharacterByID(ctx, pc.ID)
		must(err)
		s, err := srd521.Project(p.Stats, 1, int(p.Hp), 10, 0)
		must(err)
		now, err := game.ReadSpellClock(ctx, testPool, campaign)
		must(err)
		must(game.ApplyFaithAC(ctx, testPool, campaign, actor, now, s))
		return s
	}
	featUses := func(s *srd521.CharacterSheet) int { return s.Spellcasting.Choices.FeatUses }
	// One quiet minute creates history; the scene then involves the known bearers.
	confirm(prepare(srd521.WaitMinute, "", 0))
	_, err = testPool.Exec(ctx, `UPDATE session_logs SET npcs_involved=$1::uuid[] WHERE campaign_id=$2`, known, c.ID)
	must(err)

	if p := prepare(srd521.SpellCast, srd521.DetectMagicID, 1); p.Available || !strings.Contains(p.Reason, "slot") {
		t.Fatalf("slot casting available: %+v", p)
	}
	if p := prepare(srd521.DetectMagicStudy, "", 0); p.Available {
		t.Fatal("study available without concentration")
	}
	// Shield of Faith first, so Detect Magic must replace it and its recipient link.
	faith, _ := confirm(prepare(srd521.SpellCast, "shield-of-faith", 1))
	if s := read(); s.Defenses["armor-class"].Value != 13 || s.Resources[0].Current != 1 {
		t.Fatalf("faith=%+v", s)
	}
	cast := prepare(srd521.SpellCast, srd521.DetectMagicID, 0)
	if !cast.Available || len(cast.Costs) != 2 || cast.Costs[0].ResourceID != "magic-initiate" || dice.calls != 0 {
		t.Fatalf("cast preview=%+v", cast)
	}
	castReq, first := confirm(cast)
	s := read()
	if s.Concentration == nil || s.Concentration.Spell != srd521.DetectMagicID || s.Concentration.ID != castReq.ActionID.String() || s.Concentration.ExpiresAt-s.Concentration.StartedAt != 600 || featUses(s) != 0 || s.Resources[0].Current != 1 || s.Defenses["armor-class"].Value != 11 {
		t.Fatalf("after cast=%+v", s)
	}
	var links int
	must(testPool.QueryRow(ctx, `SELECT count(*) FROM concentration_recipients WHERE cast_id=$1`, faith.ActionID).Scan(&links))
	if links != 0 {
		t.Fatal("replaced Shield of Faith recipient link survived")
	}
	var castPayload struct {
		Receipt srd521.DetectMagicReceipt `json:"receipt"`
	}
	must(json.Unmarshal(first.AppliedToolCalls[0].Result, &castPayload))
	if first.AppliedToolCalls[0].Tool != "srd_detect_magic" || !castPayload.Receipt.Presence || castPayload.Receipt.Before == nil || castPayload.Receipt.Before.ID != faith.ActionID.String() || len(castPayload.Receipt.Auras) != 0 {
		t.Fatalf("cast receipt=%+v", castPayload.Receipt)
	}
	for i := 0; i < 2; i++ {
		replay, err := newEngine().ConfirmAction(ctx, campaign, castReq)
		must(err)
		if !sameTurnResult(first, replay) || featUses(read()) != 0 {
			t.Fatal("cast replay changed")
		}
	}
	if p := prepare(srd521.SpellCast, srd521.DetectMagicID, 0); p.Available || !strings.Contains(p.Reason, "spent") {
		t.Fatalf("spent free use available: %+v", p)
	}
	// The TUI preview selects the owner's character, carries the live
	// concentration for status and writes nothing.
	tui, err := newEngine().PreviewDetectMagic(ctx, campaign, game.ActionIntent{Kind: srd521.DetectMagicStudy})
	must(err)
	if !tui.Preview.Available || tui.Preview.Intent.ActorID != actor || tui.Concentration == nil || tui.Concentration.ID != castReq.ActionID.String() {
		t.Fatalf("TUI preview: %+v", tui)
	}
	if _, err := newEngine().PreviewDetectMagic(ctx, campaign, game.ActionIntent{Kind: srd521.SpellCast, EntryID: "shield-of-faith"}); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("other spell previewed: %v", err)
	}
	// Study discloses only the visible, unblocked, in-range, player-known bearer.
	studyReq, study := confirm(prepare(srd521.DetectMagicStudy, "", 0))
	var studyPayload struct {
		Receipt srd521.DetectMagicReceipt `json:"receipt"`
	}
	must(json.Unmarshal(study.AppliedToolCalls[0].Result, &studyPayload))
	if len(studyPayload.Receipt.Auras) != 1 || studyPayload.Receipt.Auras[0].SubjectID != ward.String() || studyPayload.Receipt.Auras[0].School != "abjuration" {
		t.Fatalf("study=%+v", studyPayload.Receipt)
	}
	for _, secret := range []string{"Lead casket", "Distant idol", "Veiled watcher", "Secret sigil", "Unmet lich", "necromancy", "evocation", "illusion", "conjuration", "transmutation"} {
		if strings.Contains(study.Narrative, secret) || strings.Contains(string(study.AppliedToolCalls[0].Result), secret) {
			t.Fatalf("study disclosed %q", secret)
		}
	}
	replay, err := newEngine().ConfirmAction(ctx, campaign, studyReq)
	must(err)
	if !sameTurnResult(study, replay) {
		t.Fatal("study replay changed")
	}
	// A changed scene ruling invalidates an older preview.
	stale := requestFor(prepare(srd521.DetectMagicStudy, "", 0))
	_, err = testPool.Exec(ctx, `UPDATE npcs SET properties=$1 WHERE id=$2`, facts(20, "metal-1-inch", true, "spell", "abjuration", false), ward)
	must(err)
	if _, err := newEngine().ConfirmAction(ctx, campaign, stale); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("stale facts admitted: %v", err)
	}
	// The unseen bearer is still sensed, but no visible aura remains to study.
	_, blocked := confirm(prepare(srd521.DetectMagicStudy, "", 0))
	if !strings.Contains(blocked.Narrative, "No sensed magic is on a visible") || strings.Contains(blocked.Narrative, "Warded door") {
		t.Fatalf("blocked study=%s", blocked.Narrative)
	}
	// Reload keeps the clock-bound effect; ten quiet minutes end it.
	for i := 0; i < 10; i++ {
		confirm(prepare(srd521.WaitMinute, "", 0))
	}
	if s := read(); s.Concentration != nil {
		t.Fatal("Detect Magic outlived ten minutes")
	}
	// Long Rest recovers the free use. A failed history write rolls back the cast.
	confirm(prepare(srd521.LongRest, "", 0))
	if featUses(read()) != 1 {
		t.Fatal("Long Rest did not recover the free use")
	}
	rollback := requestFor(prepare(srd521.SpellCast, srd521.DetectMagicID, 0))
	_, err = testPool.Exec(ctx, `CREATE FUNCTION reject_detect_history() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.llm_response LIKE 'Detect Magic cast%' THEN RAISE EXCEPTION 'history failure'; END IF; RETURN NEW; END; $$; CREATE TRIGGER reject_detect_history BEFORE INSERT ON session_logs FOR EACH ROW EXECUTE FUNCTION reject_detect_history()`)
	must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `DROP TRIGGER IF EXISTS reject_detect_history ON session_logs; DROP FUNCTION IF EXISTS reject_detect_history()`)
	})
	if _, err := newEngine().ConfirmAction(ctx, campaign, rollback); err == nil {
		t.Fatal("history failure ignored")
	}
	if s := read(); featUses(s) != 1 || s.Concentration != nil {
		t.Fatalf("partial cast survived rollback: %+v", s)
	}
	_, err = testPool.Exec(ctx, `DROP TRIGGER reject_detect_history ON session_logs; DROP FUNCTION reject_detect_history()`)
	must(err)
	recovered, err := newEngine().ConfirmAction(ctx, campaign, rollback)
	must(err)
	again, err := newEngine().ConfirmAction(ctx, campaign, rollback)
	must(err)
	if !sameTurnResult(recovered, again) || featUses(read()) != 0 {
		t.Fatal("recovered cast replay changed")
	}
	// Ending concentration needs no action and refunds nothing.
	confirm(prepare(srd521.EndConcentration, "", 0))
	if s := read(); s.Concentration != nil || featUses(s) != 0 {
		t.Fatal("end concentration changed the feat use")
	}
	// Without the recorded feat choice, the free use is unavailable.
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats #- '{srd521,casting,choices}' WHERE id=$1`, pc.ID)
	must(err)
	if p := prepare(srd521.SpellCast, srd521.DetectMagicID, 0); p.Available || !strings.Contains(p.Reason, "No recorded Magic Initiate") {
		t.Fatalf("missing feat choice available: %+v", p)
	}
}
