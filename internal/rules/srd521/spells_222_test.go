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

func loadSpells222(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_222.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.169-172 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells222SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/tongues","name":"Tongues","source_id":"srd-5.2.1:content/spells/tongues","source_page":169,"scope":"complete-source-record","level":3,"school":"divination","lists":["bard","cleric","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":false,"material":{"description":"a miniature ziggurat","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"understanding","mechanic":"utility.tongues","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.tongues","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","tongues":{"source_pages":[169,170],"understands":"any-spoken-or-signed-language-the-target-hears-or-sees","understood_when":"target-communicates-by-speaking-or-signing","understood_by":"any-creature-that-knows-at-least-one-language","understood_requirement":"creature-can-hear-the-speech-or-see-the-signing"}}}`,
		`{"id":"srd-5.2.1:content/spells/transport-via-plants","name":"Transport via Plants","source_id":"srd-5.2.1:content/spells/transport-via-plants","source_page":170,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":10},"target":{"shape":"single","count":1,"eligibility":"large-or-larger-inanimate-plant-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"link","mechanic":"teleport.transport-via-plants","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.transport-via-plants"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","transport_via_plants":{"minimum_plant_size":"large","plant_kind":"inanimate-plant","destination":"another-plant-at-any-distance-on-same-plane-of-existence","destination_requirement":"seen-or-touched-at-least-once","travelers":"any-creature","movement_cost_feet":5}}}`,
		`{"id":"srd-5.2.1:content/spells/tree-stride","name":"Tree Stride","source_id":"srd-5.2.1:content/spells/tree-stride","source_page":170,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"stride","mechanic":"movement.tree-stride","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.tree-stride","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","tree_stride":{"destination_tree":"same-kind","range_feet":500,"tree_requirement":"both-living-and-at-least-casters-size","enter_movement_feet":5,"awareness":"instantly-knows-all-other-same-kind-trees-in-range","awareness_feet":500,"entry_choices":["pass-into-one-of-those-trees","step-out-of-entered-tree"],"exit_distance_feet":5,"exit_movement_feet":5,"no_movement":"appear-near-entered-tree","no_movement_distance_feet":5,"uses_per_turn":1,"turn_end":"must-end-each-turn-outside-a-tree"}}}`,
		`{"id":"srd-5.2.1:content/spells/true-polymorph","name":"True Polymorph","source_id":"srd-5.2.1:content/spells/true-polymorph","source_page":170,"scope":"complete-source-record","level":9,"school":"transmutation","lists":["bard","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"choice","count":1,"eligibility":"visible-creature-or-nonmagical-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of mercury, a dollop of gum arabic, and a wisp of smoke","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"transform","mechanic":"transformation.true-polymorph","timing":"for-duration-or-until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"temporary-hit-points","mechanic":"healing.true-polymorph-temporary-hit-points","timing":"on-creature-into-creature","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.true-polymorph","healing.true-polymorph-temporary-hit-points","save.wisdom-negates","duration.true-polymorph-becomes-until-dispelled","concentration.source-duration","components.material"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-target-dies-or-destroyed-until-dispelled-after-full-concentration","true_polymorph":{"target_choices":["creature","nonmagical-object"],"creature_results":["different-creature","nonmagical-object"],"object_result":"creature","object_requirement":"neither-worn-nor-carried","ends_early":["target-dies","target-destroyed"],"full_concentration_seconds":3600,"full_concentration_result":"lasts-until-dispelled","unwilling_save":"wisdom","unwilling_save_success":"unaffected","creature_into_creature":{"form_choice":"any-kind-caster-chooses","form_limit":"challenge-rating-no-higher-than-targets-challenge-rating-or-level","statistics":"replaced-by-new-form-stat-block","retained":["hit-points","hit-point-dice","alignment","personality"],"temporary_hit_points":"equal-to-hit-points-of-new-form","temporary_hit_points_end":"vanish-if-any-remain-when-spell-ends","limits":["actions-limited-by-new-form-anatomy","cannot-speak","cannot-cast-spells"],"gear":"melds-into-new-form","gear_use":"cannot-use-or-benefit-from-equipment"},"object_into_creature":{"creature_kind":"any","size_limit":"no-larger-than-objects-size","max_challenge_rating":9,"attitude":"friendly-to-caster-and-allies","turn_order":"immediately-after-casters-turn","commands":"obeys-casters-commands","control_seconds":3600,"after_control":"caster-no-longer-controls-it","after_control_attitude":"might-remain-friendly-depending-on-treatment"},"creature_into_object":{"gear":"transforms-with-everything-worn-and-carried","size_limit":"object-no-larger-than-creatures-size","statistics":"become-those-of-object","memory":"none-of-time-in-form-after-spell-ends-and-creature-returns-to-normal"}}}}`,
		`{"id":"srd-5.2.1:content/spells/true-resurrection","name":"True Resurrection","source_id":"srd-5.2.1:content/spells/true-resurrection","source_page":170,"scope":"complete-source-record","level":9,"school":"necromancy","lists":["cleric","druid"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-dead-no-longer-than-200-years-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"diamonds worth 25,000+ GP, which the spell consumes","cost_cp":2500000,"consumed":true}},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"revive","mechanic":"revive.true-resurrection","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["revive.true-resurrection","cast.one-hour","components.material"],"slot_levels":[9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","true_resurrection":{"source_pages":[170,171],"max_dead_years":200,"excluded_death":"old-age","revived_hit_points":"all-hit-points","restores":["closes-all-wounds","neutralizes-any-poison","cures-all-magical-contagions","lifts-curses-affecting-creature-when-it-died","replaces-damaged-or-missing-organs-and-limbs"],"undead":"restored-to-non-undead-form","new_body":"can-be-provided-if-original-no-longer-exists","new_body_requirement":"caster-speaks-creatures-name","new_body_space":"unoccupied-space-caster-chooses","new_body_distance_feet":10}}}`,
		`{"id":"srd-5.2.1:content/spells/true-seeing","name":"True Seeing","source_id":"srd-5.2.1:content/spells/true-seeing","source_page":171,"scope":"complete-source-record","level":6,"school":"divination","lists":["bard","cleric","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"mushroom powder worth 25+ GP, which the spell consumes","cost_cp":2500,"consumed":true}},"ritual":false,"concentration":false,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"truesight","mechanic":"sense.true-seeing-truesight","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.true-seeing-truesight","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","true_seeing":{"recipient":"willing-creature-touched","sense":"truesight","truesight_feet":120}}}`,
		`{"id":"srd-5.2.1:content/spells/true-strike","name":"True Strike","source_id":"srd-5.2.1:content/spells/true-strike","source_page":171,"scope":"complete-source-record","level":0,"school":"divination","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":false,"somatic":true,"material":{"description":"a weapon with which you have proficiency and that is worth 1+ CP","cost_cp":1,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":6,"effect_ids":["extra"]},"effects":[{"id":"attack","mechanic":"attack.true-strike","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"extra","mechanic":"damage.true-strike-upgrade","timing":"on-hit-from-level-5","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["attack.true-strike","damage.true-strike-upgrade","ability.true-strike-spellcasting","scaling.cantrip-upgrade","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","true_strike":{"weapon":"weapon-used-in-casting","attack_count":1,"ability":"spellcasting-ability-instead-of-strength-or-dexterity","ability_applies_to":["attack-rolls","damage-rolls"],"damage_type_choices":["radiant","weapons-normal-damage-type"],"damage_type_choice_when":"attack-deals-damage","upgrade_damage_type":"radiant","upgrade_applies":"with-either-damage-type-choice","upgrades":[{"level":5,"dice_count":1,"die_sides":6},{"level":11,"dice_count":2,"die_sides":6},{"level":17,"dice_count":3,"die_sides":6}]}}}`,
		`{"id":"srd-5.2.1:content/spells/tsunami","name":"Tsunami","source_id":"srd-5.2.1:content/spells/tsunami","source_page":171,"scope":"complete-source-record","level":8,"school":"conjuration","lists":["druid"],"casting_time":"1-minute","range":{"kind":"feet","feet":5280},"target":{"shape":"wall","count":0,"eligibility":"creatures-in-water-wall-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":36,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"appear","mechanic":"damage.tsunami-appear","timing":"when-wall-appears","dice_count":6,"die_sides":10,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"moving","mechanic":"damage.tsunami-moving-wall","timing":"when-moving-wall-contains-or-enters-space","dice_count":5,"die_sides":10,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.tsunami-appear","damage.tsunami-moving-wall","movement.tsunami-wall","check.tsunami-swim","save.strength-half","save.strength-tsunami","cast.one-minute","concentration.source-duration"],"slot_levels":[8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-concentration-or-wall-height-0","tsunami":{"max_length_feet":300,"max_height_feet":300,"max_thickness_feet":50,"appear_save":"strength","appear_dice_count":6,"appear_die_sides":10,"appear_damage_type":"bludgeoning","appear_success":"half-damage","move_timing":"start-of-each-casters-turn-after-wall-appears","move_feet":50,"move_direction":"away-from-caster","carried":"creatures-in-wall-move-with-it","moving_max_creature_size":"huge","moving_triggers":["inside-wall","space-wall-enters-when-it-moves"],"moving_save":"strength","moving_dice_count":5,"moving_die_sides":10,"moving_damage_type":"bludgeoning","moving_damage_limit":"once-per-round","reduction_timing":"end-of-the-turn","height_reduction_feet":50,"later_damage_reduction_dice_count":1,"later_damage_reduction_die_sides":10,"ends_at_height_feet":0,"swimming":"caught-creature-can-move-by-swimming","swim_check":"strength-athletics","swim_check_dc":"caster-spell-save-dc","swim_check_failure":"cannot-move","leaving_wall":"creature-falls-to-ground"}}}`,
		`{"id":"srd-5.2.1:content/spells/unseen-servant","name":"Unseen Servant","source_id":"srd-5.2.1:content/spells/unseen-servant","source_page":171,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["bard","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-ground-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of string and of wood","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"servant","mechanic":"summon.unseen-servant","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.unseen-servant","command.unseen-servant","cast.ritual","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-servant-at-0-hit-points-or-commanded-too-far","unseen_servant":{"traits":["invisible","mindless","shapeless"],"size":"medium","kind":"force","purpose":"performs-simple-tasks-at-casters-command","armor_class":10,"hit_points":1,"strength":2,"can_attack":false,"command_action":"bonus-action","command_frequency":"once-on-each-of-casters-turns","command_kind":"mental","move_feet":15,"command_includes":"interact-with-an-object","task_limit":"simple-tasks-a-human-could-do","task_examples":["fetching-things","cleaning","mending","folding-clothes","lighting-fires","serving-food","pouring-drinks"],"after_command":"performs-task-until-complete-then-waits-for-next-command","max_distance_feet":60,"early_end":["servant-drops-to-0-hit-points","commanded-task-would-move-it-beyond-max-distance"]}}}`,
		`{"id":"srd-5.2.1:content/spells/vampiric-touch","name":"Vampiric Touch","source_id":"srd-5.2.1:content/spells/vampiric-touch","source_page":171,"scope":"complete-source-record","level":3,"school":"necromancy","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"target-within-reach"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.vampiric-touch","timing":"on-hit-at-casting-and-on-later-magic-actions","dice_count":3,"die_sides":6,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"healing","mechanic":"healing.vampiric-touch","timing":"on-hit-after-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.vampiric-touch","healing.vampiric-touch","attack.melee-spell","action.vampiric-touch-repeat","concentration.source-duration","scaling.per-slot"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","vampiric_touch":{"attack":"melee-spell-attack","attack_target":"one-creature-within-reach","damage_dice_count":3,"damage_die_sides":6,"damage_type":"necrotic","healing":"half-of-necrotic-damage-dealt","healing_divisor":2,"repeat_action":"magic-action","repeat_timing":"on-each-of-casters-turns-until-spell-ends","repeat_target":"same-creature-or-a-different-one","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/vicious-mockery","name":"Vicious Mockery","source_id":"srd-5.2.1:content/spells/vicious-mockery","source_page":171,"scope":"complete-source-record","level":0,"school":"enchantment","lists":["bard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-or-heard-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.vicious-mockery","timing":"on-failed-save","dice_count":1,"die_sides":6,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disadvantage","mechanic":"condition.vicious-mockery-attack-disadvantage","timing":"on-failed-save-next-attack-before-end-of-targets-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.vicious-mockery","condition.vicious-mockery-attack-disadvantage","save.wisdom-negates","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","vicious_mockery":{"source_pages":[171,172],"save":"wisdom","damage_dice_count":1,"damage_die_sides":6,"damage_type":"psychic","disadvantage_on":"next-attack-roll-target-makes","disadvantage_window":"before-end-of-targets-next-turn","upgrades":[{"level":5,"dice_count":2,"die_sides":6},{"level":11,"dice_count":3,"die_sides":6},{"level":17,"dice_count":4,"die_sides":6}]}}}`,
		`{"id":"srd-5.2.1:content/spells/vitriolic-sphere","name":"Vitriolic Sphere","source_id":"srd-5.2.1:content/spells/vitriolic-sphere","source_page":172,"scope":"complete-source-record","level":4,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-acid-sphere-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of bile","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"slot","above_level":4,"dice_count":2,"die_sides":4,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.vitriolic-sphere","timing":"immediate","dice_count":10,"die_sides":4,"damage_type":"acid","bonus":0,"add_casting_modifier":false,"half_initial_only":true},{"id":"delayed","mechanic":"damage.vitriolic-sphere-delayed","timing":"end-of-targets-next-turn-on-failed-save","dice_count":5,"die_sides":4,"damage_type":"acid","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.vitriolic-sphere","damage.vitriolic-sphere-delayed","save.dexterity-half","scaling.per-slot","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","vitriolic_sphere":{"ball_diameter_feet":1,"radius_feet":20,"save":"dexterity","initial_dice_count":10,"initial_die_sides":4,"delayed_dice_count":5,"delayed_die_sides":4,"damage_type":"acid","delayed_timing":"end-of-targets-next-turn","save_success":"half-initial-damage-only","slot_applies_to":"initial-damage","slot_dice_count":2,"slot_die_sides":4}}}`,
	}
	batch := loadSpells222(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Tongues", "Transport via Plants", "Tree Stride", "True Polymorph", "True Resurrection", "True Seeing", "True Strike", "Tsunami", "Unseen Servant", "Vampiric Touch", "Vicious Mockery", "Vitriolic Sphere"}
	pages := []int{169, 170, 170, 170, 170, 171, 171, 171, 171, 171, 171, 172}
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

func TestSpells222CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells222(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t), loadSpells220(t), loadSpells221(t), loadSpells222(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 324 {
		t.Fatalf("expected 324 records, got %d", len(all.Records))
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
func TestSpells222MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_222.json")
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

func TestSpells222JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_222.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"radius_feet": 20`, `"radius_feet": 20, "radius_feet": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "RADIUS_FEET": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/tongues"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "extra"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-hour"`, `"casting_time": "1-minute"`},
		{`"casting_seconds": 3600`, `"casting_seconds": 60`},
		{`"feet": 5280`, `"feet": 5281`},
		{`"cost_cp": 2500000`, `"cost_cp": 2500001`},
		{`"concentration": true`, `"concentration": false`},
		{`"ritual": true`, `"ritual": false`},
		{`"consumed": true`, `"consumed": false`},
		{`"verbal": false`, `"verbal": true`},
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

func TestSpells222InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"tongues pages", 0, func(s *SpellRecord) { s.SourceDetails.Tongues.SourcePages = []int{169} }},
		{"transport cost", 1, func(s *SpellRecord) { s.SourceDetails.TransportViaPlants.MovementCostFeet = 10 }},
		{"tree range", 2, func(s *SpellRecord) { s.SourceDetails.TreeStride.RangeFeet = 50 }},
		{"tree uses", 2, func(s *SpellRecord) { s.SourceDetails.TreeStride.UsesPerTurn = 2 }},
		{"polymorph retained", 3, func(s *SpellRecord) {
			s.SourceDetails.TruePolymorph.CreatureIntoCreature.Retained = s.SourceDetails.TruePolymorph.CreatureIntoCreature.Retained[:3]
		}},
		{"polymorph challenge", 3, func(s *SpellRecord) { s.SourceDetails.TruePolymorph.ObjectIntoCreature.MaxChallengeRating = 10 }},
		{"resurrection years", 4, func(s *SpellRecord) { s.SourceDetails.TrueResurrection.MaxDeadYears = 100 }},
		{"resurrection casting", 4, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"truesight range", 5, func(s *SpellRecord) { s.SourceDetails.TrueSeeing.TruesightFeet = 60 }},
		{"strike upgrade", 6, func(s *SpellRecord) { s.SourceDetails.TrueStrike.Upgrades[2].DiceCount = 4 }},
		{"strike verbal", 6, func(s *SpellRecord) { *s.Components.Verbal = true }},
		{"tsunami moving dice", 7, func(s *SpellRecord) { s.Effects[1].DiceCount = 6 }},
		{"tsunami rounds", 7, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"servant strength", 8, func(s *SpellRecord) { s.SourceDetails.UnseenServant.Strength = 10 }},
		{"servant attack", 8, func(s *SpellRecord) { s.SourceDetails.UnseenServant.CanAttack = true }},
		{"vampiric healing", 9, func(s *SpellRecord) { s.SourceDetails.VampiricTouch.HealingDivisor = 1 }},
		{"vampiric slot die", 9, func(s *SpellRecord) { s.Scaling.DieSides = 8 }},
		{"mockery upgrades", 10, func(s *SpellRecord) {
			s.SourceDetails.ViciousMockery.Upgrades = s.SourceDetails.ViciousMockery.Upgrades[:2]
		}},
		{"mockery save", 10, func(s *SpellRecord) { s.SourceDetails.ViciousMockery.Save = "charisma" }},
		{"sphere delayed dice", 11, func(s *SpellRecord) { s.Effects[1].DiceCount = 10 }},
		{"sphere half initial", 11, func(s *SpellRecord) { *s.Effects[0].HalfInitialOnly = false }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 0, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 11, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.VitriolicSphere = &SpellVitriolicSphere{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells222(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells222RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_222.json")
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

func TestSpells222WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_222.json")
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
		`"radius_feet": 20, "radius_feet": 21`,
		`"radius_feet": 20, "RADIUS_FEET": 21`,
		`"radius_feet": 20, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"radius_feet": 20`), []byte(replacement), 1)
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
func TestSpells222PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_222.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"movement_cost_feet", 1, 5},
		{"range_feet", 1, 500},
		{"enter_movement_feet", 1, 5},
		{"awareness_feet", 1, 500},
		{"exit_distance_feet", 1, 5},
		{"exit_movement_feet", 1, 5},
		{"no_movement_distance_feet", 1, 5},
		{"uses_per_turn", 1, 1},
		{"full_concentration_seconds", 1, 3600},
		{"max_challenge_rating", 1, 9},
		{"control_seconds", 1, 3600},
		{"max_dead_years", 1, 200},
		{"new_body_distance_feet", 1, 10},
		{"truesight_feet", 1, 120},
		{"attack_count", 1, 1},
		{"level", 8, 5},
		{"dice_count", 17, 1},
		{"die_sides", 17, 6},
		{"level", 9, 11},
		{"dice_count", 18, 2},
		{"die_sides", 18, 6},
		{"level", 10, 17},
		{"dice_count", 19, 3},
		{"die_sides", 19, 6},
		{"max_length_feet", 1, 300},
		{"max_height_feet", 1, 300},
		{"max_thickness_feet", 1, 50},
		{"appear_dice_count", 1, 6},
		{"appear_die_sides", 1, 10},
		{"move_feet", 1, 50},
		{"moving_dice_count", 1, 5},
		{"moving_die_sides", 1, 10},
		{"height_reduction_feet", 1, 50},
		{"later_damage_reduction_dice_count", 1, 1},
		{"later_damage_reduction_die_sides", 1, 10},
		{"ends_at_height_feet", 1, 0},
		{"armor_class", 1, 10},
		{"hit_points", 1, 1},
		{"strength", 1, 2},
		{"move_feet", 2, 15},
		{"max_distance_feet", 1, 60},
		{"damage_dice_count", 1, 3},
		{"damage_die_sides", 1, 6},
		{"healing_divisor", 1, 2},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 6},
		{"damage_dice_count", 2, 1},
		{"damage_die_sides", 2, 6},
		{"level", 15, 5},
		{"dice_count", 31, 2},
		{"die_sides", 31, 6},
		{"level", 16, 11},
		{"dice_count", 32, 3},
		{"die_sides", 32, 6},
		{"level", 17, 17},
		{"dice_count", 33, 4},
		{"die_sides", 33, 6},
		{"ball_diameter_feet", 1, 1},
		{"radius_feet", 1, 20},
		{"initial_dice_count", 1, 10},
		{"initial_die_sides", 1, 4},
		{"delayed_dice_count", 1, 5},
		{"delayed_die_sides", 1, 4},
		{"slot_dice_count", 2, 2},
		{"slot_die_sides", 2, 4},
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
// inventory assigns to catalogue issue #222.
func TestSpells222MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 222 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells222(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells222ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_222.json")
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
	want := "SRD 5.2.1 PDF pp.169-172 adapted into twelve complete offline spell families for issue #222, including page-break continuation text for Tongues (pp.169-170), True Resurrection (pp.170-171) and Vicious Mockery (pp.171-172). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells222PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.VitriolicSphere = loadSpells222(t).Records[11].SourceDetails.VitriolicSphere
	if r.Validate() == nil {
		t.Fatal("batch 222 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells222RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells222(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell222Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell222Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
