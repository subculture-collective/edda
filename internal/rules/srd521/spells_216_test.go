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

func loadSpells216(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_216.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.155-157 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells216SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/prismatic-wall","name":"Prismatic Wall","source_id":"srd-5.2.1:content/spells/prismatic-wall","source_page":155,"scope":"complete-source-record","level":9,"school":"abjuration","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"wall-or-globe-centered-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"wall","mechanic":"barrier.prismatic-wall","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"light","mechanic":"light.prismatic-wall","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blindness","mechanic":"condition.prismatic-wall-blinded","timing":"approach-or-start-turn-within-20-feet","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"red","mechanic":"damage.prismatic-wall-red","timing":"passing-layer-1","dice_count":12,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"orange","mechanic":"damage.prismatic-wall-orange","timing":"passing-layer-2","dice_count":12,"die_sides":6,"damage_type":"acid","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"yellow","mechanic":"damage.prismatic-wall-yellow","timing":"passing-layer-3","dice_count":12,"die_sides":6,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"green","mechanic":"damage.prismatic-wall-green","timing":"passing-layer-4","dice_count":12,"die_sides":6,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blue","mechanic":"damage.prismatic-wall-blue","timing":"passing-layer-5","dice_count":12,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"indigo","mechanic":"condition.prismatic-wall-indigo","timing":"passing-layer-6","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"violet","mechanic":"condition.prismatic-wall-violet","timing":"passing-layer-7","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.prismatic-wall","light.prismatic-wall","condition.prismatic-wall-blinded","damage.prismatic-wall-red","damage.prismatic-wall-orange","damage.prismatic-wall-yellow","damage.prismatic-wall-green","damage.prismatic-wall-blue","condition.prismatic-wall-indigo","condition.prismatic-wall-violet","save.constitution-prismatic-wall-blindness","save.dexterity-prismatic-wall-layers","destruction.prismatic-wall-layers"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-all-layers-destroyed","prismatic_wall":{"appearance":"shimmering-multicolored-plane-of-light","wall_length_feet":90,"wall_height_feet":30,"wall_thickness_inches":1,"wall_orientation":"vertical-opaque","globe_diameter_feet":30,"shape_choices":["vertical-wall","globe"],"occupied_space":"spell-ends-instantly-without-effect","bright_light_feet":100,"dim_light_extra_feet":100,"safe_creatures":["caster","creatures-designated-at-casting"],"blindness_distance_feet":20,"blindness_triggers":["creature-that-can-see-wall-moves-within-distance","creature-starts-turn-there"],"blindness_save":"constitution","blindness_seconds":60,"layer_count":7,"traversal":"one-layer-at-a-time-through-all-layers","layer_save":"dexterity","armor_class":10,"destruction_order":"red-to-violet-one-layer-at-a-time","destroyed_layer":"gone-for-duration","antimagic_field":"no-effect-on-wall","dispel_magic":"affects-only-violet-layer","layers":[{"order":1,"color":"red","failure":"damage-half-on-success","dice_count":12,"die_sides":6,"damage_type":"fire","additional_effect":"nonmagical-ranged-attacks-cannot-pass-destroyed-by-cold-damage","destroy_damage":25,"destroy_damage_type":"cold"},{"order":2,"color":"orange","failure":"damage-half-on-success","dice_count":12,"die_sides":6,"damage_type":"acid","additional_effect":"magical-ranged-attacks-cannot-pass-destroyed-by-strong-wind-such-as-gust-of-wind","destroy_damage":0,"destroy_damage_type":"none"},{"order":3,"color":"yellow","failure":"damage-half-on-success","dice_count":12,"die_sides":6,"damage_type":"lightning","additional_effect":"destroyed-by-force-damage","destroy_damage":60,"destroy_damage_type":"force"},{"order":4,"color":"green","failure":"damage-half-on-success","dice_count":12,"die_sides":6,"damage_type":"poison","additional_effect":"destroyed-by-passwall-or-equal-or-greater-level-solid-surface-portal-spell","destroy_damage":0,"destroy_damage_type":"none"},{"order":5,"color":"blue","failure":"damage-half-on-success","dice_count":12,"die_sides":6,"damage_type":"cold","additional_effect":"destroyed-by-fire-damage","destroy_damage":25,"destroy_damage_type":"fire"},{"order":6,"color":"indigo","failure":"restrained-then-three-saves-or-petrified","dice_count":0,"die_sides":0,"damage_type":"none","additional_effect":"spells-cannot-be-cast-through-destroyed-by-daylight-bright-light","destroy_damage":0,"destroy_damage_type":"none"},{"order":7,"color":"violet","failure":"blinded-then-wisdom-save-or-teleported","dice_count":0,"die_sides":0,"damage_type":"none","additional_effect":"destroyed-by-dispel-magic","destroy_damage":0,"destroy_damage_type":"none"}],"indigo_condition":"restrained","indigo_save":"constitution","indigo_save_timing":"end-of-each-of-its-turns","indigo_successes_to_end":3,"indigo_failures_to_petrify":3,"indigo_petrified_until":"freed-by-effect-like-greater-restoration","indigo_tracking":"need-not-be-consecutive-until-three-of-a-kind","violet_condition":"blinded","violet_save":"wisdom","violet_save_timing":"start-of-casters-next-turn","violet_success":"condition-ends","violet_failure":"condition-ends-and-teleports-to-another-plane-gm-choice"}}}`,
		`{"id":"srd-5.2.1:content/spells/private-sanctum","name":"Private Sanctum","source_id":"srd-5.2.1:content/spells/private-sanctum","source_page":156,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["wizard"],"casting_time":"10-minutes","range":{"kind":"feet","feet":120},"target":{"shape":"cube","count":0,"eligibility":"warded-area-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a thin sheet of lead","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"slot-cube-size","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["ward"]},"effects":[{"id":"ward","mechanic":"ward.private-sanctum","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["ward.private-sanctum","duration.private-sanctum-permanent-after-daily-casting","scaling.cube-size-per-slot","cast.ten-minutes","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-ended-or-permanent-until-dispelled","private_sanctum":{"min_cube_side_feet":5,"max_cube_side_feet":100,"property_choice_timing":"when-cast","choice_count":"any","properties":["sound-cannot-pass-barrier","barrier-dark-and-foggy-blocks-vision-including-darkvision","divination-sensors-cannot-appear-inside-or-pass-perimeter","creatures-inside-cannot-be-targeted-by-divination","nothing-can-teleport-in-or-out","planar-travel-blocked-inside"],"permanent_casting_days":365,"permanent_requirement":"same-spot-every-day","permanent_result":"lasts-until-dispelled","additional_cube_feet_per_slot":100}}}`,
		`{"id":"srd-5.2.1:content/spells/produce-flame","name":"Produce Flame","source_id":"srd-5.2.1:content/spells/produce-flame","source_page":156,"scope":"complete-source-record","level":0,"school":"conjuration","lists":["druid"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"light","mechanic":"light.produce-flame","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.produce-flame","timing":"magic-action-hurl-on-hit","dice_count":1,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.produce-flame","damage.produce-flame","attack.ranged-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-cast-again","produce_flame":{"flame_location":"casters-hand","flame_heat":"emits-no-heat-and-ignites-nothing","bright_light_feet":20,"dim_light_extra_feet":20,"recast":"spell-ends-if-cast-again","hurl_action":"magic-action-until-spell-ends","hurl_range_feet":60,"hurl_targets":["creature","object"],"attack":"ranged-spell-attack","upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
		`{"id":"srd-5.2.1:content/spells/programmed-illusion","name":"Programmed Illusion","source_id":"srd-5.2.1:content/spells/programmed-illusion","source_page":156,"scope":"complete-source-record","level":6,"school":"illusion","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"created-effect","count":1,"eligibility":"triggered-illusion-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"jade dust worth 25+ GP","cost_cp":2500,"consumed":false}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"illusion","mechanic":"illusion.programmed-illusion","timing":"when-trigger-occurs","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.programmed-illusion","trigger.programmed-illusion","check.programmed-illusion-investigation","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"dispelled","programmed_illusion":{"phenomena":["object","creature","other-visible-phenomenon"],"dormant_state":"imperceptible-until-triggered","max_cube_feet":30,"script_choices":["illusion-behavior","sounds"],"script_timing":"decided-when-cast","max_performance_minutes":5,"activation":"springs-into-existence-and-performs-as-described","after_performance":"disappears-and-remains-dormant","reset_minutes":10,"trigger_detail":"as-general-or-detailed-as-caster-likes","trigger_basis":["visual-phenomena","audible-phenomena"],"trigger_range_feet":30,"trigger_example":"image-of-caster-warns-others-who-try-to-open-trapped-door","physical_interaction":"reveals-illusion-because-things-pass-through","examine_action":"study","examine_check":"intelligence-investigation","examine_dc":"caster-spell-save-dc","discerned":"sees-through-image-and-its-noise-sounds-hollow"}}}`,
		`{"id":"srd-5.2.1:content/spells/project-image","name":"Project Image","source_id":"srd-5.2.1:content/spells/project-image","source_page":156,"scope":"complete-source-record","level":7,"school":"illusion","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":2640000},"target":{"shape":"created-effect","count":1,"eligibility":"previously-seen-location-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a statuette of yourself worth 5+ GP","cost_cp":500,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"copy","mechanic":"illusion.project-image","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"control","mechanic":"control.project-image-magic-action","timing":"magic-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.project-image","control.project-image-magic-action","sense.project-image-through-illusion","check.project-image-investigation","ending.project-image-damaged","concentration.source-duration","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-illusion-damaged","project_image":{"source_pages":[156,157],"copy_of":"caster","range_miles":500,"placement":"any-previously-seen-location-in-range-regardless-of-obstacles","likeness":"looks-and-sounds-like-caster","tangibility":"intangible","damage_result":"illusion-disappears-and-spell-ends","senses":"see-and-hear-through-illusion-as-if-in-its-space","move_action":"magic-action","move_max_feet":60,"behaviors":["gesture","speak","behave-as-caster-chooses"],"mannerisms":"mimics-casters-mannerisms-perfectly","physical_interaction":"reveals-illusion-because-things-pass-through","examine_action":"study","examine_check":"intelligence-investigation","examine_dc":"caster-spell-save-dc","discerned":"sees-through-image-and-its-noise-sounds-hollow"}}}`,
		`{"id":"srd-5.2.1:content/spells/protection-from-energy","name":"Protection from Energy","source_id":"srd-5.2.1:content/spells/protection-from-energy","source_page":157,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["cleric","druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"resistance","mechanic":"defense.protection-from-energy-resistance","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.protection-from-energy-resistance","concentration.source-duration"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","protection_from_energy":{"benefit":"resistance-to-one-chosen-damage-type","damage_type_choices":["acid","cold","fire","lightning","thunder"],"choose_count":1}}}`,
		`{"id":"srd-5.2.1:content/spells/protection-from-evil-and-good","name":"Protection from Evil and Good","source_id":"srd-5.2.1:content/spells/protection-from-evil-and-good","source_page":157,"scope":"complete-source-record","level":1,"school":"abjuration","lists":["cleric","druid","paladin","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a flask of Holy Water worth 25+ GP, which the spell consumes","cost_cp":2500,"consumed":true}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"protection","mechanic":"protection.protection-from-evil-and-good","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.protection-from-evil-and-good","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","protection_from_evil_and_good":{"protected_against":["aberration","celestial","elemental","fey","fiend","undead"],"attack_penalty":"those-types-have-disadvantage-on-attacks-against-target","immunities_from_those_types":["possession","charmed","frightened"],"existing_effect_benefit":"advantage-on-new-saves-against-relevant-effect"}}}`,
		`{"id":"srd-5.2.1:content/spells/protection-from-poison","name":"Protection from Poison","source_id":"srd-5.2.1:content/spells/protection-from-poison","source_page":157,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["cleric","druid","paladin","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"cure","mechanic":"ending.protection-from-poison-poisoned","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ward","mechanic":"defense.protection-from-poison","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["ending.protection-from-poison-poisoned","defense.protection-from-poison"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","protection_from_poison":{"condition_ended":"poisoned","save_benefit":"advantage-on-saves-to-avoid-or-end-poisoned","resistance_damage_type":"poison"}}}`,
		`{"id":"srd-5.2.1:content/spells/purify-food-and-drink","name":"Purify Food and Drink","source_id":"srd-5.2.1:content/spells/purify-food-and-drink","source_page":157,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["cleric","druid","paladin"],"casting_time":"action","range":{"kind":"feet","feet":10},"target":{"shape":"sphere","count":0,"eligibility":"nonmagical-food-and-drink-in-sphere"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"purify","mechanic":"utility.purify-food-and-drink","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.purify-food-and-drink","cast.ritual"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","purify_food_and_drink":{"removes":["poison","rot"],"eligible":"nonmagical-food-and-drink","radius_feet":5,"center":"point-within-range"}}}`,
		`{"id":"srd-5.2.1:content/spells/raise-dead","name":"Raise Dead","source_id":"srd-5.2.1:content/spells/raise-dead","source_page":157,"scope":"complete-source-record","level":5,"school":"necromancy","lists":["bard","cleric","paladin"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"dead-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 500+ GP, which the spell consumes","cost_cp":50000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"revival","mechanic":"revive.raise-dead","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"penalty","mechanic":"penalty.raise-dead-d20-tests","timing":"until-reduced-by-long-rests","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["revive.raise-dead","penalty.raise-dead-d20-tests","cast.one-hour","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","raise_dead":{"max_days_dead":10,"death_requirement":"not-undead-when-it-died","return_hit_points":1,"poison_benefit":"neutralizes-poisons-affecting-creature-at-death","wounds":"closes-all-mortal-wounds","missing_parts":"not-restored","automatic_failure":"lacks-body-parts-or-organs-integral-for-survival","failure_example":"missing-head","d20_test_penalty":4,"penalty_reduction_per_long_rest":1,"penalty_ends_at":0}}}`,
		`{"id":"srd-5.2.1:content/spells/ray-of-enfeeblement","name":"Ray of Enfeeblement","source_id":"srd-5.2.1:content/spells/ray-of-enfeeblement","source_page":157,"scope":"complete-source-record","level":2,"school":"necromancy","lists":["warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"minor","mechanic":"disadvantage.ray-of-enfeeblement-next-attack","timing":"on-successful-save-until-start-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"weakness","mechanic":"disadvantage.ray-of-enfeeblement-strength","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage-reduction","mechanic":"debuff.ray-of-enfeeblement-damage","timing":"on-failed-save-for-duration","dice_count":1,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["disadvantage.ray-of-enfeeblement-next-attack","disadvantage.ray-of-enfeeblement-strength","debuff.ray-of-enfeeblement-damage","save.constitution-ray-of-enfeeblement","save.ray-of-enfeeblement-repeat","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-successful-save","ray_of_enfeeblement":{"save":"constitution","success_result":"disadvantage-on-next-attack-roll","success_window":"until-start-of-casters-next-turn","failure_result":"disadvantage-on-strength-based-d20-tests-for-duration","damage_roll_subtraction_dice_count":1,"damage_roll_subtraction_die_sides":8,"repeat_save_timing":"end-of-each-of-its-turns","repeat_success":"spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/ray-of-frost","name":"Ray of Frost","source_id":"srd-5.2.1:content/spells/ray-of-frost","source_page":157,"scope":"complete-source-record","level":0,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.ray-of-frost","timing":"on-hit","dice_count":1,"die_sides":8,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"slow","mechanic":"movement.ray-of-frost-speed","timing":"on-hit-until-start-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.ray-of-frost","movement.ray-of-frost-speed","attack.ranged-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","ray_of_frost":{"appearance":"frigid-beam-of-blue-white-light","attack":"ranged-spell-attack","speed_reduction_feet":10,"speed_reduction_ends":"start-of-casters-next-turn","upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
	}
	batch := loadSpells216(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Prismatic Wall", "Private Sanctum", "Produce Flame", "Programmed Illusion", "Project Image", "Protection from Energy", "Protection from Evil and Good", "Protection from Poison", "Purify Food and Drink", "Raise Dead", "Ray of Enfeeblement", "Ray of Frost"}
	pages := []int{155, 156, 156, 156, 156, 157, 157, 157, 157, 157, 157, 157}
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

func TestSpells216CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells216(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 264 {
		t.Fatalf("expected 264 records, got %d", len(all.Records))
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
func TestSpells216MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_216.json")
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

func TestSpells216JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_216.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"wall_length_feet": 90`, `"wall_length_feet": 90, "wall_length_feet": 91`},
		{`"wall_length_feet": 90`, `"wall_length_feet": 90, "WALL_LENGTH_FEET": 91`},
		{`"wall_length_feet": 90`, `"wall_length_feet": 90, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/prismatic-wall"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "ward"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-hour"`, `"casting_time": "8-hours"`},
		{`"casting_time": "10-minutes"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"casting_seconds": 3600`, `"casting_seconds": 60`},
		{`"casting_seconds": 600`, `"casting_seconds": 601`},
		{`"casting_time": "bonus-action"`, `"casting_time": "action"`},
		{`"feet": 60`, `"feet": 61`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
		{`"ritual": true`, `"ritual": false`},
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

func TestSpells216InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"wall length", 0, func(s *SpellRecord) { s.SourceDetails.PrismaticWall.WallLengthFeet = 100 }},
		{"wall layer destroy", 0, func(s *SpellRecord) { s.SourceDetails.PrismaticWall.Layers[2].DestroyDamage = 25 }},
		{"sanctum cube", 1, func(s *SpellRecord) { s.SourceDetails.PrivateSanctum.MaxCubeSideFeet = 200 }},
		{"sanctum scaling kind", 1, func(s *SpellRecord) { s.Scaling.Kind = "slot-radius" }},
		{"flame range", 2, func(s *SpellRecord) { s.SourceDetails.ProduceFlame.HurlRangeFeet = 30 }},
		{"flame casting", 2, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"illusion reset", 3, func(s *SpellRecord) { s.SourceDetails.ProgrammedIllusion.ResetMinutes = 5 }},
		{"illusion cost", 3, func(s *SpellRecord) { *s.Components.Material.CostCP = 5000 }},
		{"image range", 4, func(s *SpellRecord) { s.Range.Feet = 5280 }},
		{"image pages", 4, func(s *SpellRecord) { s.SourceDetails.ProjectImage.SourcePages = []int{156} }},
		{"energy choices", 5, func(s *SpellRecord) { s.SourceDetails.ProtectionFromEnergy.DamageTypeChoices[4] = "poison" }},
		{"energy concentration", 5, func(s *SpellRecord) { *s.Concentration = false }},
		{"evil good types", 6, func(s *SpellRecord) {
			s.SourceDetails.ProtectionFromEvilAndGood.ProtectedAgainst = s.SourceDetails.ProtectionFromEvilAndGood.ProtectedAgainst[:5]
		}},
		{"evil good consumed", 6, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"poison resistance", 7, func(s *SpellRecord) { s.SourceDetails.ProtectionFromPoison.ResistanceDamageType = "acid" }},
		{"poison duration", 7, func(s *SpellRecord) { s.DurationSeconds = 600 }},
		{"purify radius", 8, func(s *SpellRecord) { s.SourceDetails.PurifyFoodAndDrink.RadiusFeet = 10 }},
		{"purify ritual", 8, func(s *SpellRecord) { *s.Ritual = false }},
		{"raise days", 9, func(s *SpellRecord) { s.SourceDetails.RaiseDead.MaxDaysDead = 100 }},
		{"raise penalty", 9, func(s *SpellRecord) { s.SourceDetails.RaiseDead.D20TestPenalty = 2 }},
		{"enfeeblement save", 10, func(s *SpellRecord) { s.SourceDetails.RayOfEnfeeblement.Save = "wisdom" }},
		{"enfeeblement die", 10, func(s *SpellRecord) { s.SourceDetails.RayOfEnfeeblement.DamageRollSubtractionDieSides = 6 }},
		{"frost slow", 11, func(s *SpellRecord) { s.SourceDetails.RayOfFrost.SpeedReductionFeet = 5 }},
		{"frost upgrade", 11, func(s *SpellRecord) { s.SourceDetails.RayOfFrost.Upgrades[2].Level = 16 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.RayOfFrost = &SpellRayOfFrost{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells216(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells216RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_216.json")
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

func TestSpells216WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_216.json")
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
		`"wall_length_feet": 90, "wall_length_feet": 91`,
		`"wall_length_feet": 90, "WALL_LENGTH_FEET": 91`,
		`"wall_length_feet": 90, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"wall_length_feet": 90`), []byte(replacement), 1)
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
func TestSpells216PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_216.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"wall_length_feet", 1, 90},
		{"wall_height_feet", 1, 30},
		{"wall_thickness_inches", 1, 1},
		{"globe_diameter_feet", 1, 30},
		{"bright_light_feet", 1, 100},
		{"dim_light_extra_feet", 1, 100},
		{"blindness_distance_feet", 1, 20},
		{"blindness_seconds", 1, 60},
		{"layer_count", 1, 7},
		{"armor_class", 1, 10},
		{"order", 1, 1},
		{"dice_count", 12, 12},
		{"die_sides", 12, 6},
		{"destroy_damage", 1, 25},
		{"order", 2, 2},
		{"dice_count", 13, 12},
		{"die_sides", 13, 6},
		{"destroy_damage", 2, 0},
		{"order", 3, 3},
		{"dice_count", 14, 12},
		{"die_sides", 14, 6},
		{"destroy_damage", 3, 60},
		{"order", 4, 4},
		{"dice_count", 15, 12},
		{"die_sides", 15, 6},
		{"destroy_damage", 4, 0},
		{"order", 5, 5},
		{"dice_count", 16, 12},
		{"die_sides", 16, 6},
		{"destroy_damage", 5, 25},
		{"order", 6, 6},
		{"dice_count", 17, 0},
		{"die_sides", 17, 0},
		{"destroy_damage", 6, 0},
		{"order", 7, 7},
		{"dice_count", 18, 0},
		{"die_sides", 18, 0},
		{"destroy_damage", 7, 0},
		{"indigo_successes_to_end", 1, 3},
		{"indigo_failures_to_petrify", 1, 3},
		{"min_cube_side_feet", 1, 5},
		{"max_cube_side_feet", 1, 100},
		{"permanent_casting_days", 1, 365},
		{"additional_cube_feet_per_slot", 1, 100},
		{"bright_light_feet", 2, 20},
		{"dim_light_extra_feet", 2, 20},
		{"hurl_range_feet", 1, 60},
		{"level", 4, 5},
		{"dice_count", 24, 2},
		{"die_sides", 24, 8},
		{"level", 5, 11},
		{"dice_count", 25, 3},
		{"die_sides", 25, 8},
		{"level", 6, 17},
		{"dice_count", 26, 4},
		{"die_sides", 26, 8},
		{"max_cube_feet", 1, 30},
		{"max_performance_minutes", 1, 5},
		{"reset_minutes", 1, 10},
		{"trigger_range_feet", 1, 30},
		{"range_miles", 1, 500},
		{"move_max_feet", 1, 60},
		{"choose_count", 1, 1},
		{"radius_feet", 1, 5},
		{"max_days_dead", 1, 10},
		{"return_hit_points", 1, 1},
		{"d20_test_penalty", 1, 4},
		{"penalty_reduction_per_long_rest", 1, 1},
		{"penalty_ends_at", 1, 0},
		{"damage_roll_subtraction_dice_count", 1, 1},
		{"damage_roll_subtraction_die_sides", 1, 8},
		{"speed_reduction_feet", 1, 10},
		{"level", 16, 5},
		{"dice_count", 51, 2},
		{"die_sides", 51, 8},
		{"level", 17, 11},
		{"dice_count", 52, 3},
		{"die_sides", 52, 8},
		{"level", 18, 17},
		{"dice_count", 53, 4},
		{"die_sides", 53, 8},
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
// inventory assigns to catalogue issue #216.
func TestSpells216MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 216 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells216(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells216ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_216.json")
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
	want := "SRD 5.2.1 PDF pp.155-157 adapted into twelve complete offline spell families for issue #216, including page-break continuation text for Project Image (pp.156-157). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells216PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.RayOfFrost = loadSpells216(t).Records[11].SourceDetails.RayOfFrost
	if r.Validate() == nil {
		t.Fatal("batch 216 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells216RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells216(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell216Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell216Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
