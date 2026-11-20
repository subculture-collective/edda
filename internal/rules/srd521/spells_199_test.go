package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"testing"
)

func loadSpells199(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_199.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.114-116 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells199SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/call-lightning","name":"Call Lightning","source_id":"srd-5.2.1:content/spells/call-lightning","source_page":114,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"cylinder","count":0,"eligibility":"creatures-within-5-feet-of-chosen-point-under-cloud"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.call-lightning","timing":"on-cast-and-each-magic-action-save-half-on-success","dice_count":3,"die_sides":10,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"storm-bonus","mechanic":"damage.call-lightning-existing-storm","timing":"cast-outdoors-in-storm","dice_count":1,"die_sides":10,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.call-lightning","damage.call-lightning-existing-storm","area.cylinder-cloud-60-radius-10-tall","area.bolt-5-feet-of-point","save.dexterity-half","action.magic-repeat-bolt","environment.outdoor-storm-control","scaling.per-slot","concentration.source-duration"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","lightning":{"cloud_shape":"cylinder","cloud_height_feet":10,"cloud_radius_feet":60,"cloud_placement":"visible-point-in-range-above-caster","bolt_point":"visible-point-under-cloud","bolt_radius_feet":5,"save":"dexterity","save_success":"half-damage","initial_bolt":"on-cast","repeat_action":"magic","repeat_point_choices":["same-point","different-point"],"storm_condition":"outdoors-in-storm-on-cast","storm_effect":"control-existing-storm-instead-of-creating-new-one","storm_bonus_dice_count":1,"storm_bonus_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/calm-emotions","name":"Calm Emotions","source_id":"srd-5.2.1:content/spells/calm-emotions","source_page":114,"scope":"complete-source-record","level":2,"school":"enchantment","lists":["bard","cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"sphere","count":0,"eligibility":"humanoids-in-sphere"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"calm","mechanic":"condition.calm-emotions-chosen-effect","timing":"failed-save-until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.calm-emotions-chosen-effect","area.sphere-20-radius","save.charisma","target.humanoid-only","choice.per-creature-effect","immunity.charmed-frightened-suppression","attitude.indifferent-until-damage","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","calm":{"radius_feet":20,"center":"chosen-point-in-range","creature_type":"humanoid","save":"charisma","save_success":"unaffected","choice_per_creature":true,"choose_count":1,"choices":[{"choice":"condition-immunity","conditions":["charmed","frightened"],"existing":"suppressed-for-duration","attitude":"not-applicable","subjects":"not-applicable","ends_early":["not-applicable"],"on_spell_end":"immunity-ends"},{"choice":"indifference","conditions":["not-applicable"],"existing":"not-applicable","attitude":"indifferent","subjects":"caster-chosen-creatures-it-is-hostile-toward","ends_early":["target-takes-damage","target-witnesses-allies-taking-damage"],"on_spell_end":"attitude-returns-to-normal"}]}}}`,
		`{"id":"srd-5.2.1:content/spells/chain-lightning","name":"Chain Lightning","source_id":"srd-5.2.1:content/spells/chain-lightning","source_page":114,"scope":"complete-source-record","level":6,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"chain","count":1,"eligibility":"visible-creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"three silver pins","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-count","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["bolt"]},"effects":[{"id":"damage","mechanic":"damage.chain-lightning","timing":"each-target-save-half-on-success","dice_count":10,"die_sides":8,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"bolt","mechanic":"target.chain-lightning-secondary-bolt","timing":"leaps-from-first-target-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.chain-lightning","target.chain-lightning-secondary-bolt","target.creature-or-object","target.one-bolt-per-target","save.dexterity-half","scaling.per-slot-bolt","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","chain":{"first_target":"visible-creature-or-object-in-range","secondary_bolts":3,"secondary_targets":"up-to-count-other-targets-of-choice","secondary_range_feet":30,"secondary_range_origin":"first-target","target_kinds":["creature","object"],"max_bolts_per_target":1,"save":"dexterity","save_success":"half-damage","additional_bolts_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/charm-monster","name":"Charm Monster","source_id":"srd-5.2.1:content/spells/charm-monster","source_page":114,"scope":"complete-source-record","level":4,"school":"enchantment","lists":["bard","druid","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":1,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.charm-monster","timing":"failed-save-until-spell-ends-or-damaged-by-caster-or-allies","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.charm-monster","save.wisdom-advantage-if-fighting","condition.charmed-friendly","condition.ends-on-caster-or-ally-damage","knowledge.target-knows-charmed-after","scaling.targets-per-slot"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-damaged-by-caster-or-allies","monster":{"source_pages":[114,115],"creature_type":"any-creature","save":"wisdom","save_advantage_if":"caster-or-allies-fighting-target","save_success":"unaffected","condition":"charmed","attitude":"friendly-to-caster","ends_early":["caster-damages-target","caster-allies-damage-target"],"after_spell":"target-knows-it-was-charmed-by-caster","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/charm-person","name":"Charm Person","source_id":"srd-5.2.1:content/spells/charm-person","source_page":115,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","druid","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":1,"eligibility":"visible-humanoids-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.charm-person","timing":"failed-save-until-spell-ends-or-damaged-by-caster-or-allies","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.charm-person","target.humanoid-only","save.wisdom-advantage-if-fighting","condition.charmed-friendly","condition.ends-on-caster-or-ally-damage","knowledge.target-knows-charmed-after","scaling.targets-per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-damaged-by-caster-or-allies","person":{"source_pages":[115],"creature_type":"humanoid","save":"wisdom","save_advantage_if":"caster-or-allies-fighting-target","save_success":"unaffected","condition":"charmed","attitude":"friendly-to-caster","ends_early":["caster-damages-target","caster-allies-damage-target"],"after_spell":"target-knows-it-was-charmed-by-caster","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/chill-touch","name":"Chill Touch","source_id":"srd-5.2.1:content/spells/chill-touch","source_page":115,"scope":"complete-source-record","level":0,"school":"necromancy","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"target-within-reach"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.chill-touch","timing":"melee-spell-attack-hit","dice_count":1,"die_sides":10,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"no-healing","mechanic":"restriction.chill-touch-no-hit-point-regain","timing":"on-hit-until-end-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.chill-touch","restriction.chill-touch-no-hit-point-regain","attack.melee-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","chill":{"attack":"melee-spell-attack","reach":"target-within-reach","hit_restriction":"cannot-regain-hit-points","restriction_ends":"end-of-caster-next-turn","upgrades":[{"character_level":5,"dice_count":2,"die_sides":10},{"character_level":11,"dice_count":3,"die_sides":10},{"character_level":17,"dice_count":4,"die_sides":10}]}}}`,
		`{"id":"srd-5.2.1:content/spells/chromatic-orb","name":"Chromatic Orb","source_id":"srd-5.2.1:content/spells/chromatic-orb","source_page":115,"scope":"complete-source-record","level":1,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"single","count":1,"eligibility":"target-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 50+ GP","cost_cp":5000,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.chromatic-orb","timing":"ranged-spell-attack-hit","dice_count":3,"die_sides":8,"damage_type":"chosen-type","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"leap","mechanic":"target.chromatic-orb-leap","timing":"same-number-on-two-or-more-damage-d8s","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.chromatic-orb","target.chromatic-orb-leap","attack.ranged-spell","damage.chosen-type","roll.matching-damage-dice","scaling.per-slot","scaling.leaps-equal-slot-level","target.once-per-casting","components.priced-material","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","orb":{"damage_types":["acid","cold","fire","lightning","poison","thunder"],"attack":"ranged-spell-attack","leap_trigger":"same-number-on-two-or-more-d8s","leap_die_sides":8,"leap_range_feet":30,"leap_target":"different-target-of-choice-within-range-of-target","leap_attack":"new-attack-roll","leap_damage":"new-damage-roll","base_slot_leaps":1,"leap_again_min_slot":2,"max_leaps":"level-of-slot-expended","targeting_limit":"creature-targeted-only-once-per-casting"}}}`,
		`{"id":"srd-5.2.1:content/spells/circle-of-death","name":"Circle of Death","source_id":"srd-5.2.1:content/spells/circle-of-death","source_page":115,"scope":"complete-source-record","level":6,"school":"necromancy","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"the powder of a crushed black pearl worth 500+ GP","cost_cp":50000,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":6,"dice_count":2,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.circle-of-death","timing":"sphere-save-half-on-success","dice_count":8,"die_sides":8,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.circle-of-death","area.sphere-60-radius","save.constitution-half","scaling.per-slot","components.priced-material","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","death":{"radius_feet":60,"center":"chosen-point-in-range","save":"constitution","save_success":"half-damage","appearance":"negative-energy-ripples"}}}`,
		`{"id":"srd-5.2.1:content/spells/clairvoyance","name":"Clairvoyance","source_id":"srd-5.2.1:content/spells/clairvoyance","source_page":115,"scope":"complete-source-record","level":3,"school":"divination","lists":["bard","cleric","sorcerer","wizard"],"casting_time":"10-minutes","range":{"kind":"feet","feet":5280},"target":{"shape":"location","count":1,"eligibility":"familiar-or-obvious-location-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a focus worth 100+ GP, either a jeweled horn for hearing or a glass eye for seeing","cost_cp":10000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sensor","mechanic":"sensor.clairvoyance","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sensor.clairvoyance","cast.ten-minutes","range.one-mile","sense.chosen-seeing-or-hearing","action.bonus-switch-sense","sensor.invisible-intangible-invulnerable","sensor.visible-to-see-invisibility-truesight","components.priced-focus","components.material","concentration.source-duration"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","clairvoyance":{"range_miles":1,"location_choices":["familiar","obvious-unfamiliar"],"familiar_definition":"visited-or-seen-before","obvious_examples":["behind-a-door","around-a-corner","in-a-grove-of-trees"],"sensor_invisible":true,"sensor_intangible":true,"sensor_invulnerable":true,"sensor_moves":false,"sense_choices":["seeing","hearing"],"sense_use":"as-if-in-sensor-space","switch_action":"bonus-action","focus_choices":[{"focus":"jeweled-horn","sense":"hearing"},{"focus":"glass-eye","sense":"seeing"}],"focus_minimum_cost_cp":10000,"sensor_seen_by_examples":["see-invisibility","truesight"],"visible_sensor":"luminous-orb-about-size-of-caster-fist"}}}`,
		`{"id":"srd-5.2.1:content/spells/clone","name":"Clone","source_id":"srd-5.2.1:content/spells/clone","source_page":115,"scope":"complete-source-record","level":8,"school":"necromancy","lists":["wizard"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-or-flesh-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 1,000+ GP, which the spell consumes, and a sealable vessel worth 2,000+ GP that is large enough to hold the creature being cloned","cost_cp":300000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"clone","mechanic":"clone.grow-inert-duplicate","timing":"finishes-growing-after-120-days","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"soul","mechanic":"clone.soul-transfer-on-death","timing":"original-dies-after-clone-finishes-forming","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["clone.grow-inert-duplicate","clone.soul-transfer-on-death","cast.one-hour","target.creature-or-flesh","time.growth-120-days","choice.clone-age","state.inert-while-vessel-undisturbed","soul.free-and-willing","equipment.none-transferred","remains.inert-cannot-be-revived","components.consumed-priced-diamond","components.priced-vessel","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","clone":{"source_pages":[115,116],"target_choices":["creature","at-least-1-cubic-inch-of-its-flesh"],"flesh_minimum_cubic_inches":1,"forms_in":"vessel-used-in-casting","growth_days":120,"age_choices":["same-age","younger"],"inert_state":"inert-and-endures-indefinitely-while-vessel-undisturbed","transfer_trigger":"original-dies-after-clone-finishes-forming","soul_conditions":["free","willing-to-return"],"identical":["physical-form","personality","memories","abilities"],"excluded":["equipment"],"original_remains":"inert-and-cannot-be-revived","diamond_minimum_cost_cp":100000,"diamond_consumed":true,"vessel_minimum_cost_cp":200000,"vessel_consumed":false,"vessel_requirements":["sealable","large-enough-to-hold-creature-being-cloned"]}}}`,
		`{"id":"srd-5.2.1:content/spells/cloudkill","name":"Cloudkill","source_id":"srd-5.2.1:content/spells/cloudkill","source_page":116,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.cloudkill","timing":"save-on-cast-sphere-moves-in-entry-or-turn-end-once-per-turn-half-on-success","dice_count":5,"die_sides":8,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.cloudkill","area.sphere-20-radius","area.heavily-obscured","area.moves-10-feet-from-caster","save.constitution-half-once-per-turn","environment.strong-wind-ends","scaling.per-slot","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-strong-wind-ended","cloudkill":{"radius_feet":20,"center":"point-in-range","appearance":"yellow-green-fog","obscurement":"heavily-obscured","dispersed_by":"strong-wind","dispersal_example":"gust-of-wind","dispersal_effect":"ends-spell","save":"constitution","save_success":"half-damage","save_triggers":["on-cast-in-sphere","sphere-moves-into-space","enters-sphere","ends-turn-in-sphere"],"saves_per_turn":1,"move_feet":10,"move_direction":"away-from-caster","move_timing":"start-of-each-caster-turn"}}}`,
		`{"id":"srd-5.2.1:content/spells/color-spray","name":"Color Spray","source_id":"srd-5.2.1:content/spells/color-spray","source_page":116,"scope":"complete-source-record","level":1,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"cone","count":0,"eligibility":"creatures-in-cone"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of colorful sand","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"blind","mechanic":"condition.color-spray-blinded","timing":"failed-save-until-end-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.color-spray-blinded","area.cone-15","save.constitution","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","spray":{"cone_feet":15,"origin":"caster","save":"constitution","save_success":"unaffected","condition":"blinded","condition_ends":"end-of-caster-next-turn","appearance":"dazzling-array-of-flashing-colorful-light"}}}`,
	}
	batch := loadSpells199(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Call Lightning", "Calm Emotions", "Chain Lightning", "Charm Monster", "Charm Person", "Chill Touch", "Chromatic Orb", "Circle of Death", "Clairvoyance", "Clone", "Cloudkill", "Color Spray"}
	pages := []int{114, 114, 114, 114, 115, 115, 115, 115, 115, 115, 116, 116}
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
			// Chill Touch is a cantrip: it has no slot levels at all.
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

func TestSpells199CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells199(t)
	all, err := ComposeSpellSourceCatalogue(legacy, loadSpells218(t), loadSpells198(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 37 {
		t.Fatal("expected Acid Arrow and three twelve-family batches")
	}
	for _, r := range all.Records {
		if r.Scope != "complete-source-record" || r.MechanicsAvailable() {
			t.Fatal("source granted availability")
		}
	}
	if _, err := ComposeSpellSourceCatalogue(batch, batch); err == nil {
		t.Fatal("duplicate canonical IDs accepted")
	}
	for _, r := range legacy.Records[1:] {
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
func TestSpells199MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_199.json")
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

func TestSpells199JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_199.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"cloud_height_feet": 10`, `"cloud_height_feet": 10, "cloud_height_feet": 11`},
		{`"cloud_height_feet": 10`, `"cloud_height_feet": 10, "CLOUD_HEIGHT_FEET": 11`},
		{`"cloud_height_feet": 10`, `"cloud_height_feet": 10, "unknown": true`},
		{`"repeat_point_choices": [
            "same-point",
            "different-point"
          ]`, `"repeat_point_choices": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/call-lightning"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "damage"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "1-hour"`, `"casting_time": "action"`},
		{`"casting_seconds": 600`, `"casting_seconds": 0`},
		{`"feet": 5280`, `"feet": 5281`},
		{`"cost_cp": 300000`, `"cost_cp": 100000`},
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

func TestSpells199InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"storm bonus", 0, func(s *SpellRecord) { s.SourceDetails.Lightning.StormBonusDiceCount = 0 }},
		{"repeat point choices", 0, func(s *SpellRecord) { s.SourceDetails.Lightning.RepeatPointChoices[0] = "unknown" }},
		{"calm unknown choice", 1, func(s *SpellRecord) { s.SourceDetails.Calm.Choices[1].Choice = "unknown" }},
		{"calm choice count", 1, func(s *SpellRecord) { s.SourceDetails.Calm.ChooseCount = 2 }},
		{"chain bolt count", 2, func(s *SpellRecord) { s.SourceDetails.Chain.SecondaryBolts = 4 }},
		{"chain zero increment", 2, func(s *SpellRecord) { s.SourceDetails.Chain.AdditionalBoltsPerSlot = 0 }},
		{"chain scaling dependency", 2, func(s *SpellRecord) { s.RequiredMechanics[5] = "scaling.per-slot-ray" }},
		{"monster creature type", 3, func(s *SpellRecord) { s.SourceDetails.Monster.CreatureType = "humanoid" }},
		{"monster continuation page", 3, func(s *SpellRecord) { s.SourceDetails.Monster.SourcePages = []int{114} }},
		{"person creature type", 4, func(s *SpellRecord) { s.SourceDetails.Person.CreatureType = "any-creature" }},
		{"charm payload swapped", 4, func(s *SpellRecord) {
			s.SourceDetails.Monster, s.SourceDetails.Person = s.SourceDetails.Person, nil
		}},
		{"cantrip upgrade row", 5, func(s *SpellRecord) { s.SourceDetails.Chill.Upgrades[2].CharacterLevel = 16 }},
		{"cantrip slots", 5, func(s *SpellRecord) { s.SlotLevels = []int{1} }},
		{"cantrip upgrade on levelled spell", 6, func(s *SpellRecord) { s.Scaling.Kind = "cantrip-upgrade" }},
		{"orb damage type", 6, func(s *SpellRecord) { s.SourceDetails.Orb.DamageTypes[5] = "force" }},
		{"orb leap slot", 6, func(s *SpellRecord) { s.SourceDetails.Orb.LeapAgainMinSlot = 1 }},
		{"circle save", 7, func(s *SpellRecord) { s.SourceDetails.Death.Save = "dexterity" }},
		{"clairvoyance sense", 8, func(s *SpellRecord) { s.SourceDetails.Clairvoyance.SenseChoices[0] = "unknown" }},
		{"clairvoyance focus pairing", 8, func(s *SpellRecord) { s.SourceDetails.Clairvoyance.FocusChoices[0].Sense = "seeing" }},
		{"clairvoyance moving sensor", 8, func(s *SpellRecord) { s.SourceDetails.Clairvoyance.SensorMoves = true }},
		{"clone vessel consumed", 9, func(s *SpellRecord) { s.SourceDetails.Clone.VesselConsumed = true }},
		{"clone diamond kept", 9, func(s *SpellRecord) { s.SourceDetails.Clone.DiamondConsumed = false }},
		{"clone continuation page", 9, func(s *SpellRecord) { s.SourceDetails.Clone.SourcePages = []int{115} }},
		{"clone casting seconds", 9, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"cloudkill frequency", 10, func(s *SpellRecord) { s.SourceDetails.Cloudkill.SavesPerTurn = 2 }},
		{"cloudkill trigger", 10, func(s *SpellRecord) { s.SourceDetails.Cloudkill.SaveTriggers[1] = "unknown" }},
		{"spray condition", 11, func(s *SpellRecord) { s.SourceDetails.Spray.Condition = "stunned" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Spray = &SpellSpray{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells199(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells199RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_199.json")
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

func TestSpells199WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_199.json")
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
		`"choice": "condition-immunity", "choice": "unknown"`,
		`"choice": "condition-immunity", "CHOICE": "unknown"`,
		`"choice": "condition-immunity", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"choice": "condition-immunity"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

func TestSpells199RejectContradictoryTargetScaling(t *testing.T) {
	for _, index := range []int{3, 4} {
		for _, invalid := range []string{"zero increment", "mismatched base"} {
			t.Run(fmt.Sprintf("%d/%s", index, invalid), func(t *testing.T) {
				r := loadSpells199(t)
				s := &r.Records[index]
				if d := s.SourceDetails.Monster; d != nil {
					if invalid == "zero increment" {
						d.AdditionalTargetsPerSlot = 0
					} else {
						d.BaseTargets = 99
					}
				} else {
					d := s.SourceDetails.Person
					if invalid == "zero increment" {
						d.AdditionalTargetsPerSlot = 0
					} else {
						d.BaseTargets = 99
					}
				}
				if err := r.Validate(); err == nil {
					t.Error("contradictory target scaling validated")
				}
				b, err := json.Marshal(r)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
					t.Error("contradictory target scaling decoded")
				}
			})
		}
	}
}
