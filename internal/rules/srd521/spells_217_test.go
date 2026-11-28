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

func loadSpells217(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_217.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.158-159 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells217SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/ray-of-sickness","name":"Ray of Sickness","source_id":"srd-5.2.1:content/spells/ray-of-sickness","source_page":158,"scope":"complete-source-record","level":1,"school":"necromancy","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.ray-of-sickness","timing":"on-hit","dice_count":2,"die_sides":8,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"poisoned","mechanic":"condition.ray-of-sickness-poisoned","timing":"on-hit-until-end-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.ray-of-sickness","condition.ray-of-sickness-poisoned","attack.ranged-spell","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","ray_of_sickness":{"appearance":"greenish-ray","attack":"ranged-spell-attack","condition":"poisoned","condition_ends":"end-of-casters-next-turn","damage_dice_count":2,"damage_die_sides":8,"damage_type":"poison","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/regenerate","name":"Regenerate","source_id":"srd-5.2.1:content/spells/regenerate","source_page":158,"scope":"complete-source-record","level":7,"school":"transmutation","lists":["bard","cleric","druid"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a prayer wheel","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"healing","mechanic":"healing.regenerate-initial","timing":"immediate","dice_count":4,"die_sides":8,"damage_type":"","bonus":15,"add_casting_modifier":false,"half_initial_only":false},{"id":"recurring","mechanic":"healing.regenerate-turn","timing":"start-of-each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false},{"id":"regrowth","mechanic":"restoration.regenerate-body-parts","timing":"after-2-minutes","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.regenerate-initial","healing.regenerate-turn","restoration.regenerate-body-parts","cast.one-minute","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","regenerate":{"initial_dice_count":4,"initial_die_sides":8,"initial_bonus":15,"turn_healing":1,"turn_timing":"start-of-each-of-its-turns-for-duration","severed_parts":"any-severed-body-parts-regrow","regrowth_minutes":2}}}`,
		`{"id":"srd-5.2.1:content/spells/reincarnate","name":"Reincarnate","source_id":"srd-5.2.1:content/spells/reincarnate","source_page":158,"scope":"complete-source-record","level":5,"school":"necromancy","lists":["druid"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"dead-humanoid-or-piece-of-one-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"rare oils worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"new-body","mechanic":"revive.reincarnate","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["revive.reincarnate","table.reincarnate-species","cast.one-hour","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","reincarnate":{"target_forms":["dead-humanoid","piece-of-dead-humanoid"],"max_days_dead":10,"result":"forms-new-body-and-calls-soul-into-it","table_die_sides":10,"species":[{"roll":1,"species":"roll-again"},{"roll":2,"species":"dragonborn"},{"roll":3,"species":"dwarf"},{"roll":4,"species":"elf"},{"roll":5,"species":"gnome"},{"roll":6,"species":"goliath"},{"roll":7,"species":"halfling"},{"roll":8,"species":"human"},{"roll":9,"species":"orc"},{"roll":10,"species":"tiefling"}],"gm_alternative":"gm-chooses-another-playable-species","species_choices":"creature-makes-choices-species-description-offers","memory":"recalls-former-life","retained":"capabilities-of-original-form","trait_change":"loses-previous-species-traits-and-gains-new-species-traits"}}}`,
		`{"id":"srd-5.2.1:content/spells/remove-curse","name":"Remove Curse","source_id":"srd-5.2.1:content/spells/remove-curse","source_page":158,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["cleric","paladin","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-or-object-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"remove","mechanic":"ending.remove-curse","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["ending.remove-curse"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","remove_curse":{"curses_ended":"all-curses-affecting-one-creature-or-object","cursed_magic_item":"curse-remains-but-owners-attunement-breaks","attunement_result":"item-can-be-removed-or-discarded"}}}`,
		`{"id":"srd-5.2.1:content/spells/resilient-sphere","name":"Resilient Sphere","source_id":"srd-5.2.1:content/spells/resilient-sphere","source_page":158,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"large-or-smaller-creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a glass sphere","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"enclosure","mechanic":"barrier.resilient-sphere","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.resilient-sphere","save.dexterity-negates","movement.resilient-sphere-roll","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-disintegrate","resilient_sphere":{"max_target_size":"large","target_kinds":["creature","object"],"save":"dexterity","save_applies":"unwilling-creature","passage":"nothing-passes-barrier-in-or-out","blocked":["physical-objects","energy","other-spell-effects"],"breathing":"creature-inside-can-breathe","immunity":"sphere-immune-to-all-damage","inside_protection":"cannot-be-damaged-by-attacks-or-effects-from-outside","inside_limit":"creature-inside-cannot-damage-anything-outside","size":"weightless-and-just-large-enough-for-contents","push_action":"action","roll_speed_divisor":2,"outside_movement":"can-be-picked-up-and-moved-by-other-creatures","destroyed_by":"disintegrate-targeting-globe-without-harming-contents"}}}`,
		`{"id":"srd-5.2.1:content/spells/resistance","name":"Resistance","source_id":"srd-5.2.1:content/spells/resistance","source_page":158,"scope":"complete-source-record","level":0,"school":"abjuration","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"reduction","mechanic":"defense.resistance-damage-reduction","timing":"on-chosen-damage-once-per-turn","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.resistance-damage-reduction","concentration.source-duration"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","resistance":{"damage_type_choices":["acid","bludgeoning","cold","fire","lightning","necrotic","piercing","poison","radiant","slashing","thunder"],"choose_count":1,"reduction_dice_count":1,"reduction_die_sides":4,"uses_per_turn":1}}}`,
		`{"id":"srd-5.2.1:content/spells/resurrection","name":"Resurrection","source_id":"srd-5.2.1:content/spells/resurrection","source_page":158,"scope":"complete-source-record","level":7,"school":"necromancy","lists":["bard","cleric"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"dead-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"revival","mechanic":"revive.resurrection","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"penalty","mechanic":"penalty.resurrection-d20-tests","timing":"until-reduced-by-long-rests","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"exertion","mechanic":"penalty.resurrection-long-dead-caster","timing":"until-caster-finishes-long-rest","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["revive.resurrection","penalty.resurrection-d20-tests","penalty.resurrection-long-dead-caster","cast.one-hour","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","resurrection":{"source_pages":[158,159],"max_years_dead":100,"death_requirements":["did-not-die-of-old-age","not-undead-when-it-died"],"return_hit_points":"all","poison_benefit":"neutralizes-poisons-affecting-creature-at-death","wounds":"closes-all-mortal-wounds","missing_parts":"restored","d20_test_penalty":4,"penalty_reduction_per_long_rest":1,"penalty_ends_at":0,"taxing_days_dead":365,"taxing_penalties":["cannot-cast-spells","disadvantage-on-d20-tests"],"taxing_ends":"caster-finishes-long-rest"}}}`,
		`{"id":"srd-5.2.1:content/spells/reverse-gravity","name":"Reverse Gravity","source_id":"srd-5.2.1:content/spells/reverse-gravity","source_page":159,"scope":"complete-source-record","level":7,"school":"transmutation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":100},"target":{"shape":"cylinder","count":0,"eligibility":"creatures-and-objects-in-reversed-gravity"},"components":{"verbal":true,"somatic":true,"material":{"description":"a lodestone and iron filings","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"reversal","mechanic":"terrain.reverse-gravity","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.reverse-gravity","save.dexterity-reverse-gravity-grab","concentration.source-duration","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","reverse_gravity":{"radius_feet":50,"height_feet":100,"affected":"unanchored-creatures-and-objects","fall":"fall-upward-to-top-of-cylinder","grab_save":"dexterity","grab_result":"grab-reachable-fixed-object-and-avoid-upward-fall","collision":"strike-ceiling-or-anchored-object-as-in-downward-fall","unobstructed_top":"hovers-at-top-for-duration","ending":"affected-objects-and-creatures-fall-downward"}}}`,
		`{"id":"srd-5.2.1:content/spells/revivify","name":"Revivify","source_id":"srd-5.2.1:content/spells/revivify","source_page":159,"scope":"complete-source-record","level":3,"school":"necromancy","lists":["cleric","druid","paladin","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"dead-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 300+ GP, which the spell consumes","cost_cp":30000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"revival","mechanic":"revive.revivify","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["revive.revivify","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","revivify":{"max_seconds_dead":60,"return_hit_points":1,"cannot_revive":"creature-that-died-of-old-age","missing_parts":"not-restored"}}}`,
		`{"id":"srd-5.2.1:content/spells/rope-trick","name":"Rope Trick","source_id":"srd-5.2.1:content/spells/rope-trick","source_page":159,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"rope-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a segment of rope","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"space","mechanic":"space.rope-trick","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["space.rope-trick","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","rope_trick":{"rope_motion":"one-end-rises-until-perpendicular-or-reaches-ceiling","portal_visibility":"invisible","portal_width_feet":3,"portal_length_feet":5,"portal_location":"ropes-upper-end","space_kind":"extradimensional-until-spell-ends","access":"climb-rope","rope_handling":"rope-can-be-pulled-in-or-dropped-out","max_creatures":8,"max_creature_size":"medium","isolation":"attacks-spells-and-effects-cannot-pass-in-or-out","sight":"creatures-inside-see-through-portal","ending":"contents-drop-out"}}}`,
		`{"id":"srd-5.2.1:content/spells/sacred-flame","name":"Sacred Flame","source_id":"srd-5.2.1:content/spells/sacred-flame","source_page":159,"scope":"complete-source-record","level":0,"school":"evocation","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.sacred-flame","timing":"on-failed-dexterity-save","dice_count":1,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.sacred-flame","save.dexterity-negates","cover.sacred-flame-ignored","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","sacred_flame":{"appearance":"flame-like-radiance","save":"dexterity","ignored_cover":["half-cover","three-quarters-cover"],"upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
		`{"id":"srd-5.2.1:content/spells/sanctuary","name":"Sanctuary","source_id":"srd-5.2.1:content/spells/sanctuary","source_page":159,"scope":"complete-source-record","level":1,"school":"abjuration","lists":["cleric"],"casting_time":"bonus-action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a shard of glass from a mirror","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"ward","mechanic":"protection.sanctuary","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.sanctuary","save.wisdom-sanctuary-retarget","ending.sanctuary-hostile-act","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-warded-creature-attacks-casts-or-damages","sanctuary":{"triggers":["targets-warded-creature-with-attack-roll","targets-warded-creature-with-damaging-spell"],"save":"wisdom","failure_choices":["choose-new-target","lose-attack-or-spell"],"area_effects":"no-protection-from-areas-of-effect","ends_when_warded_creature":["makes-attack-roll","casts-spell","deals-damage"]}}}`,
	}
	batch := loadSpells217(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Ray of Sickness", "Regenerate", "Reincarnate", "Remove Curse", "Resilient Sphere", "Resistance", "Resurrection", "Reverse Gravity", "Revivify", "Rope Trick", "Sacred Flame", "Sanctuary"}
	pages := []int{158, 158, 158, 158, 158, 158, 158, 159, 159, 159, 159, 159}
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

func TestSpells217CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells217(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 276 {
		t.Fatalf("expected 276 records, got %d", len(all.Records))
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
func TestSpells217MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_217.json")
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

func TestSpells217JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_217.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"damage_dice_count": 2`, `"damage_dice_count": 2, "damage_dice_count": 3`},
		{`"damage_dice_count": 2`, `"damage_dice_count": 2, "DAMAGE_DICE_COUNT": 3`},
		{`"damage_dice_count": 2`, `"damage_dice_count": 2, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/ray-of-sickness"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "damage"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_time": "1-hour"`, `"casting_time": "8-hours"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"casting_seconds": 3600`, `"casting_seconds": 60`},
		{`"casting_time": "bonus-action"`, `"casting_time": "action"`},
		{`"feet": 60`, `"feet": 61`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
		{`"consumed": true`, `"consumed": false`},
		{`"somatic": true`, `"somatic": false`},
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

func TestSpells217InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"sickness dice", 0, func(s *SpellRecord) { s.Effects[0].DiceCount = 3 }},
		{"sickness slot die", 0, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"regenerate bonus", 1, func(s *SpellRecord) { s.Effects[0].Bonus = 10 }},
		{"regenerate regrowth", 1, func(s *SpellRecord) { s.SourceDetails.Regenerate.RegrowthMinutes = 10 }},
		{"reincarnate table", 2, func(s *SpellRecord) { s.SourceDetails.Reincarnate.Species[0].Species = "human" }},
		{"reincarnate cost", 2, func(s *SpellRecord) { *s.Components.Material.CostCP = 50000 }},
		{"curse item", 3, func(s *SpellRecord) { s.SourceDetails.RemoveCurse.CursedMagicItem = "curse-ends" }},
		{"curse range", 3, func(s *SpellRecord) { s.Range.Kind = "self" }},
		{"sphere save", 4, func(s *SpellRecord) { s.SourceDetails.ResilientSphere.Save = "strength" }},
		{"sphere size", 4, func(s *SpellRecord) { s.SourceDetails.ResilientSphere.MaxTargetSize = "huge" }},
		{"resistance die", 5, func(s *SpellRecord) { s.SourceDetails.Resistance.ReductionDieSides = 6 }},
		{"resistance types", 5, func(s *SpellRecord) {
			s.SourceDetails.Resistance.DamageTypeChoices = s.SourceDetails.Resistance.DamageTypeChoices[:10]
		}},
		{"resurrection years", 6, func(s *SpellRecord) { s.SourceDetails.Resurrection.MaxYearsDead = 200 }},
		{"resurrection pages", 6, func(s *SpellRecord) { s.SourceDetails.Resurrection.SourcePages = []int{158} }},
		{"gravity height", 7, func(s *SpellRecord) { s.SourceDetails.ReverseGravity.HeightFeet = 50 }},
		{"gravity concentration", 7, func(s *SpellRecord) { *s.Concentration = false }},
		{"revivify window", 8, func(s *SpellRecord) { s.SourceDetails.Revivify.MaxSecondsDead = 600 }},
		{"revivify consumed", 8, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"rope capacity", 9, func(s *SpellRecord) { s.SourceDetails.RopeTrick.MaxCreatures = 9 }},
		{"rope duration", 9, func(s *SpellRecord) { s.DurationSeconds = 28800 }},
		{"flame cover", 10, func(s *SpellRecord) {
			s.SourceDetails.SacredFlame.IgnoredCover = s.SourceDetails.SacredFlame.IgnoredCover[:1]
		}},
		{"flame upgrade", 10, func(s *SpellRecord) { s.SourceDetails.SacredFlame.Upgrades[1].DiceCount = 4 }},
		{"sanctuary save", 11, func(s *SpellRecord) { s.SourceDetails.Sanctuary.Save = "charisma" }},
		{"sanctuary casting", 11, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Sanctuary = &SpellSanctuary{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells217(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells217RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_217.json")
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

func TestSpells217WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_217.json")
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
		`"damage_dice_count": 2, "damage_dice_count": 3`,
		`"damage_dice_count": 2, "DAMAGE_DICE_COUNT": 3`,
		`"damage_dice_count": 2, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"damage_dice_count": 2`), []byte(replacement), 1)
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
func TestSpells217PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_217.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"damage_dice_count", 1, 2},
		{"damage_die_sides", 1, 8},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
		{"initial_dice_count", 1, 4},
		{"initial_die_sides", 1, 8},
		{"initial_bonus", 1, 15},
		{"turn_healing", 1, 1},
		{"regrowth_minutes", 1, 2},
		{"max_days_dead", 1, 10},
		{"table_die_sides", 1, 10},
		{"roll", 1, 1},
		{"roll", 2, 2},
		{"roll", 3, 3},
		{"roll", 4, 4},
		{"roll", 5, 5},
		{"roll", 6, 6},
		{"roll", 7, 7},
		{"roll", 8, 8},
		{"roll", 9, 9},
		{"roll", 10, 10},
		{"roll_speed_divisor", 1, 2},
		{"choose_count", 1, 1},
		{"reduction_dice_count", 1, 1},
		{"reduction_die_sides", 1, 4},
		{"uses_per_turn", 1, 1},
		{"max_years_dead", 1, 100},
		{"d20_test_penalty", 1, 4},
		{"penalty_reduction_per_long_rest", 1, 1},
		{"penalty_ends_at", 1, 0},
		{"taxing_days_dead", 1, 365},
		{"radius_feet", 1, 50},
		{"height_feet", 1, 100},
		{"max_seconds_dead", 1, 60},
		{"return_hit_points", 1, 1},
		{"portal_width_feet", 1, 3},
		{"portal_length_feet", 1, 5},
		{"max_creatures", 1, 8},
		{"level", 12, 5},
		{"dice_count", 28, 2},
		{"die_sides", 28, 8},
		{"level", 13, 11},
		{"dice_count", 29, 3},
		{"die_sides", 29, 8},
		{"level", 14, 17},
		{"dice_count", 30, 4},
		{"die_sides", 30, 8},
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
// inventory assigns to catalogue issue #217.
func TestSpells217MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 217 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells217(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells217ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_217.json")
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
	want := "SRD 5.2.1 PDF pp.158-159 adapted into twelve complete offline spell families for issue #217, including page-break continuation text for Resurrection (pp.158-159). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells217PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.Sanctuary = loadSpells217(t).Records[11].SourceDetails.Sanctuary
	if r.Validate() == nil {
		t.Fatal("batch 217 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells217RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells217(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell217Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell217Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
