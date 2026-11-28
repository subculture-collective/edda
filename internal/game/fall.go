package game

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// DeclaredFall is one server-owned fall ruling in stats.srd521_fall (#147).
// Every field is required. The model cannot write it: the player-stat tool
// edits only existing integer stats.
type DeclaredFall struct {
	ID         string `json:"id"`
	HeightFeet *int   `json:"height_feet"`
	// Landing is "solid" or "water"; any other medium requires adjudication.
	Landing string `json:"landing"`
	// Flying, OntoCreature and FallPrevention exclude falls out of flight,
	// falls onto another creature and magic that stops or slows the fall.
	Flying         *bool `json:"flying"`
	OntoCreature   *bool `json:"onto_creature"`
	FallPrevention *bool `json:"fall_prevention"`
}

// Conditions is the actor's durable condition state outside an encounter,
// stats.srd521_conditions. An active encounter owns it on its player actor.
type Conditions = srd521.Conditions

// ReadConditions returns the actor's saved conditions outside encounters.
// Unknown condition state reads as none here; PrepareStoredAction refuses it
// through srd521.DecodeConditions before any action is available.
func ReadConditions(stats []byte) Conditions {
	c, _ := srd521.DecodeConditions(stats)
	return c
}

// readDeclaredFall decodes stats.srd521_fall strictly: a declaration with any
// fact this slice does not know is unsupported rather than ignored.
func readDeclaredFall(stats []byte) (*DeclaredFall, bool, error) {
	var envelope struct {
		Fall json.RawMessage `json:"srd521_fall"`
	}
	if err := json.Unmarshal(stats, &envelope); err != nil || len(envelope.Fall) == 0 || string(envelope.Fall) == "null" {
		return nil, false, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(envelope.Fall))
	decoder.DisallowUnknownFields()
	var fall DeclaredFall
	if err := decoder.Decode(&fall); err != nil {
		return nil, true, fmt.Errorf("This fall declares facts the supported rules do not cover; it requires adjudication.")
	}
	return &fall, true, nil
}

// PendingFall names the actor's declared fall for a client, or returns nil when
// no declaration has an ID. It copies no fact from a declaration that Validate
// or the strict decoder refuses; the preview gives the reason.
func PendingFall(stats []byte) *srd521.PendingFall {
	var envelope struct {
		Fall struct {
			ID string `json:"id"`
		} `json:"srd521_fall"`
	}
	if err := json.Unmarshal(stats, &envelope); err != nil || envelope.Fall.ID == "" {
		return nil
	}
	pending := &srd521.PendingFall{ID: envelope.Fall.ID}
	if fall, declared, err := readDeclaredFall(stats); declared && err == nil && fall.Validate() == nil {
		pending.HeightFeet, pending.Landing = fall.HeightFeet, fall.Landing
	}
	return pending
}

// Validate requires every fact and refuses the excluded cases.
func (f *DeclaredFall) Validate() error {
	if f.HeightFeet == nil || *f.HeightFeet < 0 {
		return fmt.Errorf("Clarify the exact height fallen in feet.")
	}
	switch f.Landing {
	case srd521.FallSolid, srd521.FallWater:
	case "":
		return fmt.Errorf("Clarify what the faller lands on.")
	default:
		return fmt.Errorf("Landing on %q requires adjudication; only solid ground and water are supported.", f.Landing)
	}
	if f.Flying == nil || f.OntoCreature == nil || f.FallPrevention == nil {
		return fmt.Errorf("Clarify whether the faller was flying, lands on a creature, or has magic that prevents the fall.")
	}
	if *f.Flying || *f.OntoCreature || *f.FallPrevention {
		return fmt.Errorf("Falls out of flight, onto a creature or with magic that prevents falling require adjudication.")
	}
	return nil
}

// fallSkill maps the water-entry choice to its ability. "" declines the
// Reaction; it is the only choice on solid ground.
func fallSkill(landing, skill string) (string, bool) {
	switch {
	case skill == "":
		return "", true
	case landing != srd521.FallWater:
		return "", false
	case skill == "athletics":
		return "strength", true
	case skill == "acrobatics":
		return "dexterity", true
	}
	return "", false
}

// prepareFall previews one declared fall. EntryID names the declaration, so
// a replaced declaration never matches an open preview. Preparation never rolls.
func prepareFall(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent) (*StoredAction, error) {
	fall, declared, declErr := readDeclaredFall(pc.Stats)
	ability := ""
	if declared && declErr == nil {
		ability, _ = fallSkill(fall.Landing, intent.Skill)
	}
	base := ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "strength"}
	if ability != "" {
		base.Ability, base.Skill = ability, intent.Skill
	}
	p, err := PrepareCharacterAction(ctx, campaign, pc, base)
	if err != nil {
		return nil, err
	}
	result := &StoredAction{Preview: p}
	p.Intent = intent
	if ability == "" {
		p.Modifier = nil
	}
	if !p.Available {
		return result, nil
	}
	p.Available = false
	if !declared || intent.EntryID == "" || (declErr == nil && fall.ID != intent.EntryID) {
		p.Reason = "No declared fall to resolve."
		return result, nil
	}
	if declErr != nil {
		p.Reason = declErr.Error()
		return result, nil
	}
	if err := fall.Validate(); err != nil {
		p.Reason = err.Error()
		return result, nil
	}
	if _, ok := fallSkill(fall.Landing, intent.Skill); !ok || intent.Ability != "" || intent.SlotLevel != 0 {
		p.Reason = "Choose no check, or for water Strength (Athletics) or Dexterity (Acrobatics)."
		return result, nil
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	if encounter != nil && encounter.Status == "active" {
		p.Reason = "Falls during an encounter are not supported."
		return result, nil
	}
	conditions := ReadConditions(pc.Stats)
	encoded, err := json.Marshal([]any{p.ExpectedRevision, fall, conditions, intent})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("fall-v1:%x", sha256.Sum256(encoded))
	height := *fall.HeightFeet
	count := srd521.FallDice(height)
	damage := "no damage (under 10 feet), and you land on your feet"
	if count > 0 {
		damage = fmt.Sprintf("%dd6 Bludgeoning damage, then the Prone condition unless you take no damage", count)
	}
	capped := ""
	if height >= 10*(srd521.FallMaxDice+1) {
		capped = fmt.Sprintf(" (capped at %dd6)", srd521.FallMaxDice)
	}
	p.Options = []string{fmt.Sprintf("Fall %d feet onto %s: %s%s. Current HP %d.", height, map[string]string{srd521.FallSolid: "solid ground", srd521.FallWater: "water"}[fall.Landing], damage, capped, pc.HP)}
	p.SourceIDs = []string{"srd-5.2.1:content/rules-glossary/falling"}
	if count > 0 {
		p.SourceIDs = append(p.SourceIDs, "srd-5.2.1:content/rules-glossary/prone")
	}
	if conditions.Prone {
		p.Options = append(p.Options, "You are already Prone.")
	}
	if fall.Landing == srd521.FallWater {
		if intent.Skill == "" {
			p.Options = append(p.Options, fmt.Sprintf("No Reaction used. You could instead use your Reaction for a DC %d Strength (Athletics) or Dexterity (Acrobatics) check to halve the damage.", srd521.FallWaterDC))
		} else {
			p.Options = append(p.Options, fmt.Sprintf("Reaction: DC %d %s check to hit the water head or feet first; on a success the damage is halved, rounded down.", srd521.FallWaterDC, map[string]string{"athletics": "Strength (Athletics)", "acrobatics": "Dexterity (Acrobatics)"}[intent.Skill]))
			p.Costs = []ActionCost{{ResourceID: "reaction", Amount: 1}}
			p.SourceIDs = append(p.SourceIDs, "srd-5.2.1:content/combat/reactions")
		}
	}
	p.Available = true
	p.Reason = ""
	result.Fall = fall
	return result, nil
}

// prepareProneStand ends a saved Prone condition outside an encounter. No
// movement budget exists there, so standing records no cost.
func prepareProneStand(ctx context.Context, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "strength"})
	if err != nil {
		return nil, err
	}
	result := &StoredAction{Preview: p}
	p.Intent = intent
	p.Modifier = nil
	if !p.Available {
		return result, nil
	}
	p.Available = false
	if intent.EntryID != "" || intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		p.Reason = "Choose a supported encounter action."
		return result, nil
	}
	if !ReadConditions(pc.Stats).Prone {
		p.Reason = "No supported encounter is active."
		return result, nil
	}
	p.Available = true
	p.Reason = ""
	p.ExpectedRevision = "stand-v1:" + p.ExpectedRevision
	p.Options = []string{"Stand up from Prone. Outside an encounter no movement budget is tracked."}
	p.SourceIDs = []string{"srd-5.2.1:content/rules-glossary/prone"}
	return result, nil
}

// HasDeclaredFall reports whether stats.srd521_fall holds any declaration,
// readable or not.
func HasDeclaredFall(stats []byte) bool {
	_, declared, _ := readDeclaredFall(stats)
	return declared
}
