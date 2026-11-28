package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func loadSpells224(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF p.175 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells224SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/wish","name":"Wish","source_id":"srd-5.2.1:content/spells/wish","source_page":175,"scope":"complete-source-record","level":9,"school":"conjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"choice","count":1,"eligibility":"one-chosen-wish-use"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"duplicate","mechanic":"utility.wish-duplicate-spell","timing":"chosen-use","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"alternative","mechanic":"utility.wish-alternative-effect","timing":"chosen-use","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"stress","mechanic":"damage.wish-stress","timing":"on-each-later-spell-cast-until-long-rest","dice_count":1,"die_sides":10,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.wish-duplicate-spell","utility.wish-alternative-effect","damage.wish-stress","penalty.wish-stress","ruling.wish-gm"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","wish":{"basic_use":{"duplicates":"any-other-spell","max_spell_level":8,"requirements":"none-including-costly-components","result":"spell-takes-effect"},"alternative_choose_count":1,"object_creation":{"object_count":1,"max_value_cp":2500000,"excluded":"magic-item","max_dimension_feet":300,"placement":"unoccupied-space-caster-can-see-on-ground"},"instant_health":{"includes":"caster","max_other_creatures":20,"requirement":"creatures-caster-can-see","healing":"regain-all-hit-points","ends":"all-effects-listed-in-greater-restoration"},"resistance":{"max_creatures":10,"requirement":"creatures-caster-can-see","damage_type_count":1,"damage_type_chooser":"caster","duration":"permanent"},"spell_immunity":{"max_creatures":10,"requirement":"creatures-caster-can-see","immunity":"single-spell-or-other-magical-effect","duration_seconds":28800},"sudden_learning":{"replaced_feats":1,"replacement":"another-feat-caster-is-eligible-for","old_feat":"loses-all-benefits","new_feat":"gains-benefits","restriction":"cannot-replace-prerequisite-for-other-feats-or-features"},"roll_redo":{"undoes":"single-recent-event","method":"forced-reroll-of-any-die-roll","window_rounds":1,"window_includes":"casters-last-turn","reality":"reshapes-to-accommodate-new-result","examples":["undo-allys-failed-saving-throw","undo-foes-critical-hit"],"reroll_options":["advantage","disadvantage"],"result_choice":"caster-chooses-reroll-or-original-roll"},"reshape_reality":{"scope":"wish-not-included-in-other-effects","statement":"state-wish-to-gm-as-precisely-as-possible","ruling":"gm-has-great-latitude","risk":"greater-wish-greater-likelihood-something-goes-wrong","possible_failures":["spell-simply-fails","desired-effect-achieved-only-in-part","unforeseen-consequence-from-wording"],"wording_example":"wishing-villain-dead-might-propel-caster-forward-to-time-villain-no-longer-alive","legendary_wish":"legendary-magic-item-or-artifact-might-transport-caster-to-current-owner","wide_consequences":"granted-wish-affecting-community-region-or-world-likely-attracts-powerful-foes","god_wish":"gods-divine-servants-might-intervene-to-prevent-or-encourage-particular-wording","multiverse_wish":"wish-that-would-undo-multiverse-fails"},"stress":{"applies":"any-effect-other-than-duplicating-another-spell","damage_trigger":"each-spell-cast-until-long-rest-finished","damage_dice_count":1,"damage_die_sides":10,"damage_type":"necrotic","damage_per":"level-of-that-spell","damage_prevention":"cannot-be-reduced-or-prevented","strength_score":3,"strength_days_dice_count":2,"strength_days_die_sides":4,"light_activity":"each-day-resting-with-nothing-more-than-light-activity","light_activity_recovery_decrease_days":2,"never_cast_again_percent":33}}}}`,
		`{"id":"srd-5.2.1:content/spells/word-of-recall","name":"Word of Recall","source_id":"srd-5.2.1:content/spells/word-of-recall","source_page":175,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":5},"target":{"shape":"self-and-up-to-count","count":5,"eligibility":"willing-creatures-within-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"recall","mechanic":"teleport.word-of-recall","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.word-of-recall","knowledge.word-of-recall-sanctuary"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","word_of_recall":{"max_other_creatures":5,"within_feet":5,"destination":"previously-designated-sanctuary","arrival":"nearest-unoccupied-space-to-designated-spot","unprepared":"spell-has-no-effect","preparation":"cast-spell-at-location-to-designate-sanctuary","sanctuary_examples":["temple"]}}}`,
		`{"id":"srd-5.2.1:content/spells/zone-of-truth","name":"Zone of Truth","source_id":"srd-5.2.1:content/spells/zone-of-truth","source_page":175,"scope":"complete-source-record","level":2,"school":"enchantment","lists":["bard","cleric","paladin"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-truth-zone-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration_seconds":600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"truth","mechanic":"social.zone-of-truth-no-deliberate-lies","timing":"on-failed-save-while-in-radius","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["social.zone-of-truth-no-deliberate-lies","save.charisma-zone-of-truth"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","zone_of_truth":{"radius_feet":15,"save_triggers":["enters-area-first-time-on-turn","starts-turn-there"],"save":"charisma","save_failure":"cannot-speak-deliberate-lie-while-in-radius","caster_knowledge":"knows-whether-creature-succeeds-or-fails","target_awareness":"affected-creature-is-aware-of-spell","target_evasion":"can-avoid-answering-questions-it-would-normally-answer-with-lie","evasion_limit":"may-be-evasive-but-must-be-truthful"}}}`,
	}
	batch := loadSpells224(t)
	if len(batch.Records) != 3 {
		t.Fatal("exact three families required")
	}
	names := []string{"Wish", "Word of Recall", "Zone of Truth"}
	pages := []int{175, 175, 175}
	for i, r := range batch.Records {
		t.Run(names[i], func(t *testing.T) {
			var want SpellRecord
			if err := json.Unmarshal([]byte(expected[i]), &want); err != nil {
				t.Fatal(err)
			}
			if r.Name != names[i] || r.SourcePage != pages[i] || !reflect.DeepEqual(r, want) {
				t.Fatalf("PDF p.%d source values differ", want.SourcePage)
			}
			if r.MechanicsAvailable() {
				t.Fatal("source effect available")
			}
			// Cantrips use no spell slot: they have no slot levels at all.
			slots := []int{}
			for level := r.Level; level <= 9 && r.Level > 0; level++ {
				slots = append(slots, level)
			}
			if !slices.Equal(r.SlotLevels, slots) {
				t.Fatal("source slot metadata incomplete")
			}
		})
	}
}

func TestSpells224CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells224(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t), loadSpells220(t), loadSpells221(t), loadSpells222(t), loadSpells223(t), loadSpells224(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 339 {
		t.Fatalf("expected 339 records, got %d", len(all.Records))
	}
	seen := map[string]bool{}
	for _, r := range all.Records {
		if r.Scope != "complete-source-record" || r.MechanicsAvailable() || seen[r.ID] {
			t.Fatal("source granted availability or repeated")
		}
		seen[r.ID] = true
	}
	if _, err := ComposeSpellSourceCatalogue(legacy, loadSpells196(t), batch); err == nil {
		t.Fatal("duplicate Acid Arrow accepted")
	}
	if _, err := ComposeSpellSourceCatalogue(batch, batch); err == nil {
		t.Fatal("duplicate canonical IDs accepted")
	}
	for _, r := range bounded.Records {
		if !r.MechanicsAvailable() {
			t.Fatal("bounded availability changed")
		}
		r.SourceDetails = batch.Records[0].SourceDetails
		if r.MechanicsAvailable() {
			t.Fatal("bounded fixture accepted source extension")
		}
	}
	for _, r := range batch.Records {
		r.RequiredMechanics = append(r.RequiredMechanics, "cast.cleric-self-slot-1", "unknown.future")
		if r.MechanicsAvailable() {
			t.Fatal("unsupported source available")
		}
	}
}

// Mutate each field, recursively, rather than only the top-level payload.
// The valid input covers explicit false and zero source values.
func TestSpells224MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(original, &root); err != nil {
		t.Fatal(err)
	}
	for index, record := range root["records"].([]any) {
		r := record.(map[string]any)
		var walk func(map[string]any, string)
		reject := func(label string) {
			t.Helper()
			b, err := json.Marshal(root)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
				t.Fatalf("record %d accepted %s", index, label)
			}
		}
		walk = func(fields map[string]any, path string) {
			for key, value := range fields {
				fieldPath := path + "/" + key
				delete(fields, key)
				reject("omitted " + fieldPath)
				fields[key] = nil
				reject("null " + fieldPath)
				fields[key] = value
				switch v := value.(type) {
				case float64:
					fields[key] = -1
					reject("negative " + fieldPath)
					fields[key] = value
				case string:
					for _, bad := range []string{"", "   "} {
						fields[key] = bad
						reject("blank " + fieldPath)
					}
					fields[key] = value
				case map[string]any:
					walk(v, fieldPath)
				case []any:
					if len(v) == 0 {
						continue
					}
					for _, bad := range []any{[]any{}, []any{nil}, []any{""}, []any{"   "}, []any{v[0], nil}, []any{v[0], map[string]any{"mixed": true}}} {
						fields[key] = bad
						reject("malformed list " + fieldPath)
					}
					fields[key] = value
					for j, item := range v {
						if nested, ok := item.(map[string]any); ok {
							walk(nested, fmt.Sprintf("%s/%d", fieldPath, j))
						}
					}
				}
			}
		}
		details := r["source_details"].(map[string]any)
		for key, payload := range details {
			if obj, ok := payload.(map[string]any); ok {
				walk(obj, key)
			}
		}
		for _, key := range []string{"source_details", "effects", "required_mechanics", "slot_levels", "lists"} {
			value := r[key]
			delete(r, key)
			reject("omitted " + key)
			r[key] = value
		}
	}
}

func TestSpells224JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"radius_feet": 15`, `"radius_feet": 15, "radius_feet": 16`},
		{`"radius_feet": 15`, `"radius_feet": 15, "RADIUS_FEET": 16`},
		{`"radius_feet": 15`, `"radius_feet": 15, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/wish"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": []`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "action"`, `"casting_time": "1-minute"`},
		{`"casting_seconds": 0`, `"casting_seconds": 60`},
		{`"feet": 60`, `"feet": 61`},
		{`"material": null`, `"material": {"description": "a coin", "cost_cp": 0, "consumed": false}`},
		{`"concentration": false`, `"concentration": true`},
		{`"ritual": false`, `"ritual": true`},
		{`"verbal": true`, `"verbal": false`},
		{`"somatic": false`, `"somatic": true`},
	} {
		t.Run(pair[1], func(t *testing.T) {
			b := bytes.Replace(original, []byte(pair[0]), []byte(pair[1]), 1)
			if bytes.Equal(b, original) {
				t.Fatal("fixture mutation did not match")
			}
			if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
				t.Fatal("invalid JSON/source graph accepted")
			}
		})
	}
}

func TestSpells224InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"wish spell level", 0, func(s *SpellRecord) { s.SourceDetails.Wish.BasicUse.MaxSpellLevel = 9 }},
		{"wish object value", 0, func(s *SpellRecord) { s.SourceDetails.Wish.ObjectCreation.MaxValueCP = 250000 }},
		{"wish health creatures", 0, func(s *SpellRecord) { s.SourceDetails.Wish.InstantHealth.MaxOtherCreatures = 10 }},
		{"wish resistance creatures", 0, func(s *SpellRecord) { s.SourceDetails.Wish.Resistance.MaxCreatures = 20 }},
		{"wish immunity hours", 0, func(s *SpellRecord) { s.SourceDetails.Wish.SpellImmunity.DurationSeconds = 3600 }},
		{"wish feat limit", 0, func(s *SpellRecord) { s.SourceDetails.Wish.SuddenLearning.Restriction = "none" }},
		{"wish reroll options", 0, func(s *SpellRecord) {
			s.SourceDetails.Wish.RollRedo.RerollOptions = s.SourceDetails.Wish.RollRedo.RerollOptions[:1]
		}},
		{"wish multiverse", 0, func(s *SpellRecord) { s.SourceDetails.Wish.ReshapeReality.MultiverseWish = "succeeds" }},
		{"wish stress strength", 0, func(s *SpellRecord) { s.SourceDetails.Wish.Stress.StrengthScore = 8 }},
		{"wish never again", 0, func(s *SpellRecord) { s.SourceDetails.Wish.Stress.NeverCastAgainPercent = 50 }},
		{"wish stress die", 0, func(s *SpellRecord) { s.Effects[2].DieSides = 12 }},
		{"recall creatures", 1, func(s *SpellRecord) { s.SourceDetails.WordOfRecall.MaxOtherCreatures = 6 }},
		{"recall range", 1, func(s *SpellRecord) { s.Range.Feet = 10 }},
		{"zone radius", 2, func(s *SpellRecord) { s.SourceDetails.ZoneOfTruth.RadiusFeet = 30 }},
		{"zone save", 2, func(s *SpellRecord) { s.SourceDetails.ZoneOfTruth.Save = "wisdom" }},
		{"required mechanic", 2, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 0, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 2, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.ZoneOfTruth = &SpellZoneOfTruth{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells224(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells224RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(original, &root); err != nil {
		t.Fatal(err)
	}
	for index, item := range root["records"].([]any) {
		r := item.(map[string]any)
		objects := []map[string]any{r, r["range"].(map[string]any), r["target"].(map[string]any), r["scaling"].(map[string]any), r["components"].(map[string]any)}
		for _, effect := range r["effects"].([]any) {
			objects = append(objects, effect.(map[string]any))
		}
		if material, ok := r["components"].(map[string]any)["material"].(map[string]any); ok {
			objects = append(objects, material)
		}
		for _, obj := range objects {
			for key, value := range obj {
				switch value.(type) {
				case float64, bool:
				default:
					continue
				}
				for _, omit := range []bool{true, false} {
					if omit {
						delete(obj, key)
					} else {
						obj[key] = nil
					}
					b, err := json.Marshal(root)
					if err != nil {
						t.Fatal(err)
					}
					if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
						t.Fatalf("record %d accepted missing/null %s", index, key)
					}
					obj[key] = value
				}
			}
		}
	}
}

func TestSpells224WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(original, &root); err != nil {
		t.Fatal(err)
	}
	var walk func(any)
	walk = func(node any) {
		switch node := node.(type) {
		case map[string]any:
			for key, value := range node {
				if list, ok := value.([]any); ok && len(list) > 0 {
					node[key] = "empty-list-marker"
					b, err := json.Marshal(root)
					node[key] = value
					if err != nil {
						t.Fatal(err)
					}
					b = bytes.Replace(b, []byte(`"empty-list-marker"`), []byte(`[ ]`), 1)
					// Empty scaling effect references are valid only for kind=none, which
					// already has an empty list and therefore never reaches this branch.
					if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
						t.Fatalf("accepted whitespace empty %s", key)
					}
				}
				walk(value)
			}
		case []any:
			for _, item := range node {
				walk(item)
			}
		}
	}
	walk(root)
	for _, replacement := range []string{
		`"radius_feet": 15, "radius_feet": 16`,
		`"radius_feet": 15, "RADIUS_FEET": 16`,
		`"radius_feet": 15, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"radius_feet": 15`), []byte(replacement), 1)
		if bytes.Equal(b, original) {
			t.Fatal("fixture mutation did not match")
		}
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells224PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"max_spell_level", 1, 8},
		{"alternative_choose_count", 1, 1},
		{"object_count", 1, 1},
		{"max_value_cp", 1, 2500000},
		{"max_dimension_feet", 1, 300},
		{"max_other_creatures", 1, 20},
		{"max_creatures", 1, 10},
		{"damage_type_count", 1, 1},
		{"max_creatures", 2, 10},
		{"duration_seconds", 2, 28800},
		{"replaced_feats", 1, 1},
		{"window_rounds", 1, 1},
		{"damage_dice_count", 1, 1},
		{"damage_die_sides", 1, 10},
		{"strength_score", 1, 3},
		{"strength_days_dice_count", 1, 2},
		{"strength_days_die_sides", 1, 4},
		{"light_activity_recovery_decrease_days", 1, 2},
		{"never_cast_again_percent", 1, 33},
		{"max_other_creatures", 2, 5},
		{"within_feet", 1, 5},
		{"radius_feet", 1, 15},
	} {
		matches := regexp.MustCompile(`"`+c.key+`": (\d+)`).FindAllSubmatchIndex(original, -1)
		if len(matches) < c.n {
			t.Fatalf("%s occurrence %d missing", c.key, c.n)
		}
		m := matches[c.n-1]
		if string(original[m[2]:m[3]]) != strconv.Itoa(c.want) {
			t.Fatalf("%s #%d is not %d", c.key, c.n, c.want)
		}
		b := append([]byte{}, original[:m[2]]...)
		b = append(b, []byte(strconv.Itoa(c.want+1))...)
		b = append(b, original[m[3]:]...)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatalf("%s #%d accepted %d", c.key, c.n, c.want+1)
		}
	}
}

// The batch must be exactly the families and pages the issue-80 decision
// inventory assigns to catalogue issue #224.
func TestSpells224MatchesDecisionInventory(t *testing.T) {
	b, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			SourcePage    int    `json:"source_page"`
			CatalogueItem int    `json:"catalogue_issue"`
		} `json:"records"`
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	want := map[string]int{}
	for _, r := range inventory.Records {
		if r.CatalogueItem == 224 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells224(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 3 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells224ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_224.json")
	if err != nil {
		t.Fatal(err)
	}
	var root struct {
		Source struct {
			Changes string `json:"changes"`
		} `json:"source"`
	}
	if err := json.Unmarshal(b, &root); err != nil {
		t.Fatal(err)
	}
	want := "SRD 5.2.1 PDF p.175 adapted into three complete offline spell families for issue #224. No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells224PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.ZoneOfTruth = loadSpells224(t).Records[2].SourceDetails.ZoneOfTruth
	if r.Validate() == nil {
		t.Fatal("batch 224 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells224RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells224(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell224Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell224Pins) != 3 {
		t.Fatal("three pinned records required")
	}
}

// #224 is the last spell batch, so the composed source catalogue must hold
// exactly the spell families, names and pages in the issue-80 decision inventory.
func TestSpells224WholeCatalogueMatchesDecisionInventory(t *testing.T) {
	b, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			Family     string `json:"family"`
			ID         string `json:"id"`
			Name       string `json:"name"`
			SourcePage int    `json:"source_page"`
		} `json:"records"`
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	want := map[string]int{}
	for _, r := range inventory.Records {
		if r.Family == "spells" {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	bounded := loadSpellRecords(t)
	bounded.Records = bounded.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells218(t), loadSpells219(t), loadSpells220(t), loadSpells221(t), loadSpells222(t), loadSpells223(t), loadSpells224(t))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]int{}
	for _, r := range all.Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 339 || len(all.Records) != len(want) || !reflect.DeepEqual(got, want) {
		t.Fatalf("composed catalogue differs from the decision inventory: got %d, want %d", len(all.Records), len(want))
	}
}
