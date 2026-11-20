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
	"testing"
)

func loadSpells201(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_201.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// The exact family list and pages come from the issue-80 decision inventory,
// not from the catalogue under test.
func TestSpells201MatchesDecisionInventory(t *testing.T) {
	b, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			SourcePage     int    `json:"source_page"`
			CatalogueIssue int    `json:"catalogue_issue"`
		} `json:"records"`
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	want := map[string]int{}
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 201 {
			want[r.ID] = r.SourcePage
		}
	}
	batch := loadSpells201(t)
	got := map[string]int{}
	for _, r := range batch.Records {
		got[r.ID] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) || len(batch.Records) != 12 {
		t.Fatalf("catalogue %v differs from inventory %v", got, want)
	}
}

// Literal PDF pp.118-121 value fixtures include every choice, effect, table row
// and scaling value, independently of decoder defaults and structural validation.
func TestSpells201SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/conjure-woodland-beings","name":"Conjure Woodland Beings","source_id":"srd-5.2.1:content/spells/conjure-woodland-beings","source_page":118,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"visible-creatures-in-emanation"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":4,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.conjure-woodland-beings","timing":"emanation-enters-creature-space-or-creature-enters-or-ends-turn-save-half-once-per-turn","dice_count":5,"die_sides":8,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disengage","mechanic":"action.conjure-woodland-beings-bonus-disengage","timing":"bonus-action-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.conjure-woodland-beings","action.conjure-woodland-beings-bonus-disengage","area.emanation-10","save.wisdom-half-once-per-turn","scaling.per-slot","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","woodland":{"source_pages":[118,119],"emanation_feet":10,"origin":"caster","triggers":["emanation-enters-creature-space","creature-enters-emanation","creature-ends-turn-in-emanation"],"trigger_visibility":"creature-you-can-see","save":"wisdom","save_success":"half-damage","saves_per_turn":1,"bonus_action":"disengage","bonus_action_duration":"spell-duration"}}}`,
		`{"id":"srd-5.2.1:content/spells/contact-other-plane","name":"Contact Other Plane","source_id":"srd-5.2.1:content/spells/contact-other-plane","source_page":119,"scope":"complete-source-record","level":5,"school":"divination","lists":["warlock","wizard"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"questions","mechanic":"divination.contact-other-plane-questions","timing":"successful-save-before-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"backlash","mechanic":"damage.contact-other-plane-backlash","timing":"failed-save","dice_count":6,"die_sides":6,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"incapacitated","mechanic":"condition.contact-other-plane-incapacitated","timing":"failed-save-until-long-rest-or-greater-restoration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["divination.contact-other-plane-questions","damage.contact-other-plane-backlash","condition.contact-other-plane-incapacitated","save.intelligence-dc-15","response.gm-one-word","cast.one-minute"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","other_plane":{"save":"intelligence","save_dc":15,"success_questions_max":5,"question_timing":"before-spell-ends","answer_format":"one-word","answer_examples":["yes","no","maybe","never","irrelevant","unclear"],"unclear_meaning":"entity-does-not-know-answer","misleading_alternative":"short-phrase","failure_damage_dice_count":6,"failure_damage_die_sides":6,"failure_damage_type":"psychic","failure_condition":"incapacitated","failure_ends":["finish-long-rest","greater-restoration-cast-on-you"]}}}`,
		`{"id":"srd-5.2.1:content/spells/contagion","name":"Contagion","source_id":"srd-5.2.1:content/spells/contagion","source_page":119,"scope":"complete-source-record","level":5,"school":"necromancy","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":604800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.contagion","timing":"failed-initial-save","dice_count":11,"die_sides":8,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"poisoned","mechanic":"condition.contagion-poisoned","timing":"failed-initial-save-until-three-saves-succeed-or-seven-days","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disadvantage","mechanic":"save-penalty.contagion-chosen-ability","timing":"while-poisoned","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"save.contagion-three-successes-or-failures","timing":"end-of-each-target-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"persist","mechanic":"save.contagion-resist-ending-poisoned","timing":"effect-would-end-poisoned","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.contagion","condition.contagion-poisoned","save-penalty.contagion-chosen-ability","save.contagion-three-successes-or-failures","save.contagion-resist-ending-poisoned","choice.contagion-ability","save.constitution"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-three-saves-succeeded","contagion":{"save":"constitution","save_success":"unaffected","ability_choice":"one-ability-chosen-on-cast","ability_effect":"disadvantage-on-saves-with-chosen-ability-while-poisoned","repeat_timing":"end-of-each-target-turn","successes_to_end":3,"success_result":"spell-ends-on-target","failures_to_last":3,"failure_result":"spell-lasts-7-days-on-target","duration_days":7,"end_effect_save":"constitution","end_effect_failure":"poisoned-condition-does-not-end"}}}`,
		`{"id":"srd-5.2.1:content/spells/contingency","name":"Contingency","source_id":"srd-5.2.1:content/spells/contingency","source_page":119,"scope":"complete-source-record","level":6,"school":"abjuration","lists":["wizard"],"casting_time":"10-minutes","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a gem-encrusted statuette of yourself worth 1,500+ GP","cost_cp":150000,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":864000,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"store","mechanic":"contingency.store-spell","timing":"cast-with-contingent-spell-deferred-until-trigger","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"trigger","mechanic":"contingency.trigger-effect","timing":"immediately-after-trigger-first-occurs","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["contingency.store-spell","contingency.trigger-effect","contingency.contingent-spell-limits","contingency.single-active","contingency.self-only","contingency.ends-without-component-on-person","cast.ten-minutes","components.priced-material","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-ended-trigger-resolved-or-component-absent","contingency":{"contingent_max_level":5,"contingent_casting_time":"action","contingent_target":"can-target-you","contingent_caster":"spell-you-can-cast","slots_expended":"both-spells","deferred_on_cast":true,"trigger_described":"when-you-cast-the-two-spells","trigger_example_spell":"water-breathing","trigger_example":"engulfed-in-water-or-similar-liquid","takes_effect":"immediately-after-trigger-first-occurs","regardless_of_wish":true,"ends_after_trigger":true,"affects_only":"caster","max_active":1,"recast_ends":"other-contingency-on-you","component_requirement":"statuette-on-your-person"}}}`,
		`{"id":"srd-5.2.1:content/spells/continual-flame","name":"Continual Flame","source_id":"srd-5.2.1:content/spells/continual-flame","source_page":119,"scope":"complete-source-record","level":2,"school":"evocation","lists":["cleric","druid","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"ruby dust worth 50+ GP, which the spell consumes","cost_cp":5000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"flame","mechanic":"light.continual-flame","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.continual-flame","duration.until-dispelled","components.consumed-priced-material","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"dispelled","flame":{"source":"object-you-touch","bright_light_radius_feet":20,"dim_light_additional_feet":20,"appearance":"regular-flame","heat":"none","fuel":"consumes-no-fuel","concealment":"can-be-covered-or-hidden","extinguish":"cannot-be-smothered-or-quenched"}}}`,
		`{"id":"srd-5.2.1:content/spells/control-water","name":"Control Water","source_id":"srd-5.2.1:content/spells/control-water","source_page":119,"scope":"complete-source-record","level":4,"school":"transmutation","lists":["cleric","druid","wizard"],"casting_time":"action","range":{"kind":"feet","feet":300},"target":{"shape":"cube","count":0,"eligibility":"water-in-cube"},"components":{"verbal":true,"somatic":true,"material":{"description":"a mixture of water and dust","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"flood","mechanic":"water.control-water-flood","timing":"chosen-effect-until-spell-ends-or-different-effect","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"part","mechanic":"water.control-water-part","timing":"chosen-effect-until-spell-ends-or-different-effect","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"redirect","mechanic":"water.control-water-redirect","timing":"chosen-effect-until-spell-ends-or-different-effect","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"whirlpool","mechanic":"water.control-water-whirlpool","timing":"chosen-effect-until-spell-ends-or-different-effect","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"whirlpool-damage","mechanic":"damage.control-water-whirlpool","timing":"first-entry-on-turn-or-end-turn-in-whirlpool-save-half-on-success","dice_count":2,"die_sides":8,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["water.control-water-flood","water.control-water-part","water.control-water-redirect","water.control-water-whirlpool","damage.control-water-whirlpool","area.cube-100","action.magic-repeat-or-change-effect","save.strength-half","check.strength-athletics-escape","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","water_control":{"source_pages":[119,120],"cube_max_feet":100,"repeat_action":"magic","repeat_choices":["same-effect","different-effect"],"effect_choices":["flood","part-water","redirect-flow","whirlpool"],"flood":{"standing_water_rise_feet":20,"wave_height_feet":20,"wave_area":"area-in-large-body-of-water","wave_path":"one-side-of-area-to-other-then-crashes","vehicle_size":"huge-or-smaller","vehicle_effect":"carried-to-other-side","capsize_percent":25,"capsize_target":"huge-or-smaller-vehicles-struck-by-wave","elevated_until":"spell-ends-or-different-effect","wave_repeats":"start-of-your-next-turn-while-flood-lasts"},"part":{"result":"trench-across-area","walls":"separated-water-forms-wall-either-side","lasts":"until-spell-ends-or-different-effect","refill":"water-slowly-fills-trench-over-next-round"},"redirect":{"movement":"flowing-water-moves-in-chosen-direction","obstacles":"over-obstacles-up-walls-or-other-unlikely-directions","beyond_area":"resumes-terrain-flow","lasts":"until-spell-ends-or-different-effect"},"whirlpool":{"minimum_area_side_feet":50,"minimum_depth_feet":25,"location":"center-of-area","base_width_feet":5,"top_width_feet":50,"height_feet":25,"pull_radius_feet":25,"pull_subjects":"creatures-in-water","pull_feet":10,"save_triggers":["first-entry-on-turn","ends-turn-in-whirlpool"],"save":"strength","save_success":"half-damage","escape_action":"take-action-to-pull-away","escape_check":"strength-athletics-vs-spell-save-dc","lasts":"until-different-effect-or-spell-ends"}}}}`,
		`{"id":"srd-5.2.1:content/spells/control-weather","name":"Control Weather","source_id":"srd-5.2.1:content/spells/control-weather","source_page":120,"scope":"complete-source-record","level":8,"school":"transmutation","lists":["cleric","druid","wizard"],"casting_time":"10-minutes","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"burning incense","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"weather","mechanic":"weather.control-weather-stage-change","timing":"on-cast-and-after-new-conditions-take-effect","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"delay","mechanic":"time.control-weather-delay","timing":"before-new-conditions-take-effect","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["weather.control-weather-stage-change","time.control-weather-delay","area.radius-5-miles","requirement.outdoors","table.precipitation","table.temperature","table.wind","cast.ten-minutes","concentration.source-duration","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-concentration-ended-or-went-indoors","weather":{"radius_miles":5,"requires_outdoors":true,"ends_indoors":true,"current_conditions_by":"gm","changeable":["precipitation","temperature","wind"],"delay_dice_count":1,"delay_die_sides":4,"delay_multiplier_minutes":10,"rechange_after_effect":true,"end_behavior":"weather-gradually-returns-to-normal","stage_change":1,"stage_directions":["up","down"],"wind_direction_changeable":true,"precipitation":[{"stage":1,"condition":"clear"},{"stage":2,"condition":"light-clouds"},{"stage":3,"condition":"overcast-or-ground-fog"},{"stage":4,"condition":"rain-hail-or-snow"},{"stage":5,"condition":"torrential-rain-driving-hail-or-blizzard"}],"temperature":[{"stage":1,"condition":"heat-wave"},{"stage":2,"condition":"hot"},{"stage":3,"condition":"warm"},{"stage":4,"condition":"cool"},{"stage":5,"condition":"cold"},{"stage":6,"condition":"freezing"}],"wind":[{"stage":1,"condition":"calm"},{"stage":2,"condition":"moderate-wind"},{"stage":3,"condition":"strong-wind"},{"stage":4,"condition":"gale"},{"stage":5,"condition":"storm"}]}}}`,
		`{"id":"srd-5.2.1:content/spells/counterspell","name":"Counterspell","source_id":"srd-5.2.1:content/spells/counterspell","source_page":120,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["sorcerer","warlock","wizard"],"casting_time":"reaction","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":false,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"interrupt","mechanic":"counter.counterspell-interrupt","timing":"failed-constitution-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"refund","mechanic":"slot.counterspell-not-expended","timing":"failed-save-if-cast-with-spell-slot","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["counter.counterspell-interrupt","slot.counterspell-not-expended","cast.reaction-trigger","save.constitution","cost.wasted-on-failed-save"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["creature-within-60-feet-casting-spell-with-verbal-somatic-or-material-components"],"end_timing":"instantaneous","counter":{"trigger_range_feet":60,"trigger_components":["verbal","somatic","material"],"trigger_visibility":"you-see-the-creature","save":"constitution","failure_result":"spell-dissipates-with-no-effect","wasted_costs":["action","bonus-action","reaction"],"slot_refund":"slot-not-expended-if-cast-with-spell-slot"}}}`,
		`{"id":"srd-5.2.1:content/spells/create-food-and-water","name":"Create Food and Water","source_id":"srd-5.2.1:content/spells/create-food-and-water","source_page":120,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"choice","count":1,"eligibility":"ground-or-containers-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"food","mechanic":"create.food-45-pounds","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"water","mechanic":"create.water-30-gallons","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"spoil","mechanic":"time.food-spoils-24-hours","timing":"uneaten-after-24-hours","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["create.food-45-pounds","create.water-30-gallons","time.food-spoils-24-hours"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","provisions":{"food_pounds":45,"water_gallons":30,"placement":["ground","containers-in-range"],"food_quality":"bland-but-nourishing","food_form":"food-of-your-choice","water_quality":"clean","spoil_hours":24,"spoil_condition":"uneaten"}}}`,
		`{"id":"srd-5.2.1:content/spells/create-undead","name":"Create Undead","source_id":"srd-5.2.1:content/spells/create-undead","source_page":121,"scope":"complete-source-record","level":6,"school":"necromancy","lists":["cleric","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":10},"target":{"shape":"up-to-count","count":3,"eligibility":"corpses-of-medium-or-small-humanoids-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"one 150+ GP black onyx stone for each corpse","cost_cp":15000,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-creature-table","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["animate"]},"effects":[{"id":"animate","mechanic":"undead.create-undead-ghoul","timing":"on-cast-night-only","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"command","mechanic":"command.create-undead-mental-bonus-action","timing":"bonus-action-each-turn-within-120-feet","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"control","mechanic":"control.create-undead-24-hours-reassert","timing":"24-hours-then-reassert-by-recasting","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["undead.create-undead-ghoul","command.create-undead-mental-bonus-action","control.create-undead-24-hours-reassert","restriction.night-only","scaling.creature-count-by-slot","components.priced-material","components.material","cast.one-minute"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"instantaneous","undead":{"night_only":true,"corpse_max":3,"corpse_sizes":["medium","small"],"corpse_type":"humanoid","result":"ghoul-under-your-control","stat_block_reference":"monsters-section-stat-block","command_action":"bonus-action","command_range_feet":120,"command_scope":["any-animated-creature","several-at-once-same-command"],"command_decides":["action","movement-on-next-turn","general-command"],"general_command_example":"guard-a-particular-place","no_command_behavior":"dodge-and-move-only-to-avoid-harm","order_persistence":"follows-order-until-task-complete","control_hours":24,"control_ends":"stops-obeying-after-24-hours","reassert":"cast-on-creature-before-24-hour-period-ends","reassert_creatures_max":3,"reassert_instead":"rather-than-animating-new-ones","onyx_per_corpse":true,"onyx_minimum_cost_cp":15000,"slot_options":[{"slot_level":7,"options":[{"creature":"ghoul","count":4}]},{"slot_level":8,"options":[{"creature":"ghoul","count":5},{"creature":"ghast-or-wight","count":2}]},{"slot_level":9,"options":[{"creature":"ghoul","count":6},{"creature":"ghast-or-wight","count":3},{"creature":"mummy","count":2}]}]}}}`,
		`{"id":"srd-5.2.1:content/spells/create-or-destroy-water","name":"Create or Destroy Water","source_id":"srd-5.2.1:content/spells/create-or-destroy-water","source_page":121,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"choice","count":1,"eligibility":"container-or-cube-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a mix of water and sand","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-increment","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["create","destroy"]},"effects":[{"id":"create","mechanic":"water.create-or-destroy-water-create","timing":"chosen-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"destroy","mechanic":"water.create-or-destroy-water-destroy","timing":"chosen-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["water.create-or-destroy-water-create","water.create-or-destroy-water-destroy","choice.create-or-destroy","area.cube-30","scaling.per-slot-increment","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","water_create":{"choices":["create-water","destroy-water"],"create_gallons":10,"create_container":"open-container","create_alternative":"rain-in-cube-extinguishing-exposed-flames","destroy_gallons":10,"destroy_container":"open-container","destroy_alternative":"destroy-fog-in-cube","cube_feet":30,"additional_gallons_per_slot":10,"additional_cube_feet_per_slot":5}}}`,
		`{"id":"srd-5.2.1:content/spells/creation","name":"Creation","source_id":"srd-5.2.1:content/spells/creation","source_page":121,"scope":"complete-source-record","level":5,"school":"illusion","lists":["sorcerer","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a paintbrush","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"special","duration_seconds":0,"scaling":{"kind":"slot-increment","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["object"]},"effects":[{"id":"object","mechanic":"create.creation-object","timing":"duration-by-material-shortest-if-multiple","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"component-failure","mechanic":"restriction.creation-material-component","timing":"used-as-another-spell-material-component","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["create.creation-object","restriction.creation-material-component","duration.special-by-material","table.creation-materials","limit.cube-5-feet","scaling.per-slot-increment","cast.one-minute","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"object-duration-by-material","creation":{"origin":"shadowfell","matter_kinds":["vegetable-matter","mineral-matter"],"vegetable_examples":["soft-goods","rope","wood"],"mineral_examples":["stone","crystal","metal"],"max_cube_feet":5,"form_requirement":"form-and-material-you-have-seen","duration_rule":"depends-on-object-material","multiple_materials":"shortest-duration","material_component_use":"other-spell-fails","materials":[{"material":"vegetable-matter","duration_seconds":86400},{"material":"stone-or-crystal","duration_seconds":43200},{"material":"precious-metals","duration_seconds":3600},{"material":"gems","duration_seconds":600},{"material":"adamantine-or-mithral","duration_seconds":60}],"additional_cube_feet_per_slot":5}}}`,
	}
	batch := loadSpells201(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Conjure Woodland Beings", "Contact Other Plane", "Contagion", "Contingency", "Continual Flame", "Control Water", "Control Weather", "Counterspell", "Create Food and Water", "Create Undead", "Create or Destroy Water", "Creation"}
	pages := []int{118, 119, 119, 119, 119, 119, 120, 120, 120, 121, 121, 121}
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
			slots := []int{}
			for level := r.Level; level <= 9; level++ {
				slots = append(slots, level)
			}
			if !slices.Equal(r.SlotLevels, slots) {
				t.Fatal("source slot metadata incomplete")
			}
		})
	}
}

func TestSpells201CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells201(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 72 {
		t.Fatalf("expected six twelve-family batches, got %d", len(all.Records))
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
func TestSpells201MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_201.json")
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

func TestSpells201JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_201.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"emanation_feet": 10`, `"emanation_feet": 10, "emanation_feet": 11`},
		{`"emanation_feet": 10`, `"emanation_feet": 10, "EMANATION_FEET": 11`},
		{`"emanation_feet": 10`, `"emanation_feet": 10, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/conjure-woodland-beings"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "damage"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [
        4,`, `"slot_levels": [
        3,`},
		{`"casting_time": "reaction"`, `"casting_time": "action"`},
		{`"casting_time": "1-minute"`, `"casting_time": "action"`},
		{`"casting_seconds": 600`, `"casting_seconds": 0`},
		{`"duration": "special"`, `"duration": "timed"`},
		{`"duration_seconds": 604800`, `"duration_seconds": 604801`},
		{`"feet": 300`, `"feet": 301`},
		{`"cost_cp": 150000`, `"cost_cp": 100000`},
		{`"cost_cp": 15000,`, `"cost_cp": 1500,`},
		{`"shape": "cube"`, `"shape": "cuboid"`},
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

func TestSpells201InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"whirlpool area side not square feet", 5, func(s *SpellRecord) { s.SourceDetails.WaterControl.Whirlpool.MinimumAreaSideFeet = 2500 }},
		{"woodland continuation page", 0, func(s *SpellRecord) { s.SourceDetails.Woodland.SourcePages = []int{118} }},
		{"woodland save", 0, func(s *SpellRecord) { s.SourceDetails.Woodland.Save = "intelligence" }},
		{"woodland trigger", 0, func(s *SpellRecord) { s.SourceDetails.Woodland.Triggers[1] = "unknown" }},
		{"woodland scaling dice", 0, func(s *SpellRecord) { s.Scaling.DiceCount = 0 }},
		{"plane save dc", 1, func(s *SpellRecord) { s.SourceDetails.OtherPlane.SaveDC = 14 }},
		{"plane answer", 1, func(s *SpellRecord) { s.SourceDetails.OtherPlane.AnswerExamples[2] = "unknown" }},
		{"plane ritual", 1, func(s *SpellRecord) { *s.Ritual = false }},
		{"plane failure end", 1, func(s *SpellRecord) { s.SourceDetails.OtherPlane.FailureEnds[1] = "unknown" }},
		{"contagion success threshold", 2, func(s *SpellRecord) { s.SourceDetails.Contagion.SuccessesToEnd = 2 }},
		{"contagion duration days", 2, func(s *SpellRecord) { s.SourceDetails.Contagion.DurationDays = 6 }},
		{"contingency spell level", 3, func(s *SpellRecord) { s.SourceDetails.Contingency.ContingentMaxLevel = 6 }},
		{"contingency single active", 3, func(s *SpellRecord) { s.SourceDetails.Contingency.MaxActive = 2 }},
		{"contingency statuette cost", 3, func(s *SpellRecord) { *s.Components.Material.CostCP = 100000 }},
		{"contingency duration", 3, func(s *SpellRecord) { s.DurationSeconds = 86400 }},
		{"flame light", 4, func(s *SpellRecord) { s.SourceDetails.Flame.DimLightAdditionalFeet = 10 }},
		{"flame ruby dust kept", 4, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"flame concentration", 4, func(s *SpellRecord) { *s.Concentration = true }},
		{"water control continuation page", 5, func(s *SpellRecord) { s.SourceDetails.WaterControl.SourcePages = []int{119} }},
		{"water control effect list", 5, func(s *SpellRecord) { s.SourceDetails.WaterControl.EffectChoices[1] = "unknown" }},
		{"water control pull", 5, func(s *SpellRecord) { s.SourceDetails.WaterControl.Whirlpool.PullFeet = 15 }},
		{"water control capsize", 5, func(s *SpellRecord) { s.SourceDetails.WaterControl.Flood.CapsizePercent = 50 }},
		{"water control range", 5, func(s *SpellRecord) { s.Range.Feet = 100 }},
		{"weather precipitation row", 6, func(s *SpellRecord) { s.SourceDetails.Weather.Precipitation[4].Condition = "storm" }},
		{"weather temperature rows", 6, func(s *SpellRecord) { s.SourceDetails.Weather.Temperature = s.SourceDetails.Weather.Temperature[:5] }},
		{"weather wind stage", 6, func(s *SpellRecord) { s.SourceDetails.Weather.Wind[2].Stage = 4 }},
		{"weather delay", 6, func(s *SpellRecord) { s.SourceDetails.Weather.DelayDieSides = 6 }},
		{"counterspell trigger", 7, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{} }},
		{"counterspell casting time", 7, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"counterspell wasted cost", 7, func(s *SpellRecord) { s.SourceDetails.Counter.WastedCosts[0] = "spell" }},
		{"counterspell range", 7, func(s *SpellRecord) { s.Range.Feet = 30 }},
		{"provisions food", 8, func(s *SpellRecord) { s.SourceDetails.Provisions.FoodPounds = 40 }},
		{"undead night", 9, func(s *SpellRecord) { s.SourceDetails.Undead.NightOnly = false }},
		{"undead slot option", 9, func(s *SpellRecord) { s.SourceDetails.Undead.SlotOptions[2].Options[2].Count = 3 }},
		{"undead reassert", 9, func(s *SpellRecord) { s.SourceDetails.Undead.ReassertCreaturesMax = 4 }},
		{"undead onyx cost", 9, func(s *SpellRecord) { *s.Components.Material.CostCP = 10000 }},
		{"undead scaling kind", 9, func(s *SpellRecord) { s.Scaling.Kind = "slot-increment" }},
		{"undead casting seconds", 9, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"water create increment", 10, func(s *SpellRecord) { s.SourceDetails.WaterCreate.AdditionalGallonsPerSlot = 0 }},
		{"water create cube", 10, func(s *SpellRecord) { s.SourceDetails.WaterCreate.CubeFeet = 20 }},
		{"water create scaling kind", 10, func(s *SpellRecord) { s.Scaling.Kind = "slot-creature-table" }},
		{"creation material row", 11, func(s *SpellRecord) { s.SourceDetails.Creation.Materials[3].DurationSeconds = 300 }},
		{"creation fixed duration", 11, func(s *SpellRecord) { s.Duration, s.DurationSeconds = "timed", 3600 }},
		{"creation special concentration", 11, func(s *SpellRecord) { *s.Concentration = true }},
		{"creation special mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[2] = "duration.unknown" }},
		{"special duration elsewhere", 8, func(s *SpellRecord) { s.Duration = "special" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Provisions = &SpellProvisions{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Spray = &SpellSpray{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells201(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells201RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_201.json")
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

func TestSpells201WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_201.json")
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
		`"material": "vegetable-matter", "material": "unknown"`,
		`"material": "vegetable-matter", "MATERIAL": "unknown"`,
		`"material": "vegetable-matter", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"material": "vegetable-matter"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

func TestSpells201RejectContradictoryScaling(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"water zero cube increment", 10, func(s *SpellRecord) { s.SourceDetails.WaterCreate.AdditionalCubeFeetPerSlot = 0 }},
		{"water wrong gallons", 10, func(s *SpellRecord) { s.SourceDetails.WaterCreate.AdditionalGallonsPerSlot = 5 }},
		{"creation zero increment", 11, func(s *SpellRecord) { s.SourceDetails.Creation.AdditionalCubeFeetPerSlot = 0 }},
		{"creation wrong increment", 11, func(s *SpellRecord) { s.SourceDetails.Creation.AdditionalCubeFeetPerSlot = 10 }},
		{"undead no table", 9, func(s *SpellRecord) { s.SourceDetails.Undead.SlotOptions = nil }},
		{"scaling above level", 9, func(s *SpellRecord) { s.Scaling.AboveLevel = 7 }},
		{"scaling dice", 10, func(s *SpellRecord) { s.Scaling.DiceCount = 1; s.Scaling.DieSides = 4 }},
		{"scaling mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[5] = "scaling.unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells201(t)
			tc.change(&r.Records[tc.index])
			if err := r.Validate(); err == nil {
				t.Error("contradictory scaling validated")
			}
			b, err := json.Marshal(r)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
				t.Error("contradictory scaling decoded")
			}
		})
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells201PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_201.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"emanation_feet", 1, 10},
		{"save_dc", 1, 15},
		{"success_questions_max", 1, 5},
		{"successes_to_end", 1, 3},
		{"duration_days", 1, 7},
		{"contingent_max_level", 1, 5},
		{"max_active", 1, 1},
		{"bright_light_radius_feet", 1, 20},
		{"dim_light_additional_feet", 1, 20},
		{"cube_max_feet", 1, 100},
		{"wave_height_feet", 1, 20},
		{"capsize_percent", 1, 25},
		{"minimum_area_side_feet", 1, 50},
		{"top_width_feet", 1, 50},
		{"pull_feet", 1, 10},
		{"radius_miles", 1, 5},
		{"delay_multiplier_minutes", 1, 10},
		{"trigger_range_feet", 1, 60},
		{"food_pounds", 1, 45},
		{"water_gallons", 1, 30},
		{"spoil_hours", 1, 24},
		{"corpse_max", 1, 3},
		{"command_range_feet", 1, 120},
		{"control_hours", 1, 24},
		{"reassert_creatures_max", 1, 3},
		{"onyx_minimum_cost_cp", 1, 15000},
		{"create_gallons", 1, 10},
		{"cube_feet", 1, 30},
		{"additional_gallons_per_slot", 1, 10},
		{"additional_cube_feet_per_slot", 1, 5},
		{"additional_cube_feet_per_slot", 2, 5},
		{"max_cube_feet", 1, 5},
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

// The shared validator accepts a special duration only on typed Creation details.
func TestSpellSpecialDurationRejectedOutsideCreation(t *testing.T) {
	for name, load := range map[string]func(*testing.T) SpellRecords{"196": loadSpells196, "197": loadSpells197, "198": loadSpells198} {
		r := load(t)
		r.Records[0].Duration, r.Records[0].DurationSeconds = "special", 0
		*r.Records[0].Concentration = false
		if r.Validate() == nil {
			t.Fatalf("#%s record accepted special duration", name)
		}
	}
}
