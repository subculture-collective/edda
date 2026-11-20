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

func loadSpells205(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_205.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.128-131 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
// Literal PDF pp.128-131 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells205SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/etherealness","name":"Etherealness","source_id":"srd-5.2.1:content/spells/etherealness","source_page":128,"scope":"complete-source-record","level":7,"school":"conjuration","lists":["bard","cleric","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"slot-targets","above_level":7,"dice_count":0,"die_sides":0,"effect_ids":["targets"]},"effects":[{"id":"border","mechanic":"plane.border-ethereal","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"interaction","mechanic":"interaction.ethereal-only","timing":"while-on-ethereal-plane","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"return","mechanic":"return.shunt-force-damage","timing":"when-spell-ends-in-occupied-space","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.instant-on-unsupported-plane","timing":"when-cast-on-ethereal-plane-or-plane-that-does-not-border-it","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"targets","mechanic":"targets.willing-creatures-per-slot","timing":"slot-level-above-7","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["plane.border-ethereal","interaction.ethereal-only","return.shunt-force-damage","ending.instant-on-unsupported-plane","targets.willing-creatures-per-slot","scaling.targets-per-slot"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","etherealness":{"region":"border-ethereal-where-it-overlaps-your-current-plane","movement":"any-direction","vertical_extra_feet_per_foot":1,"perceived_plane_appearance":"gray","perceived_plane_limit_feet":60,"affects_and_affected_by":"only-creatures-objects-and-effects-on-the-ethereal-plane","outside_creatures":"cannot-perceive-or-interact-unless-a-feature-allows","return_place":"plane-you-left-in-space-corresponding-to-your-border-ethereal-space","occupied_result":"shunted-to-nearest-unoccupied-space","shunt_damage_type":"force","shunt_damage_per_foot_moved":2,"instant_end_conditions":["cast-while-on-the-ethereal-plane","cast-on-plane-that-does-not-border-the-ethereal-plane"],"instant_end_example":"outer-planes","base_targets":1,"additional_targets_per_slot":3,"slot_target_kind":"willing-creatures-including-you","slot_target_range_feet":10}}}`,
		`{"id":"srd-5.2.1:content/spells/expeditious-retreat","name":"Expeditious Retreat","source_id":"srd-5.2.1:content/spells/expeditious-retreat","source_page":129,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["sorcerer","warlock","wizard"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"dash","mechanic":"action.dash-on-cast","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"action.dash-as-bonus-action","timing":"bonus-action-each-turn-until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["action.dash-on-cast","action.dash-as-bonus-action","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","retreat":{"cast_action":"dash","repeat_action":"dash","repeat_action_type":"bonus-action","repeat_while":"until-spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/eyebite","name":"Eyebite","source_id":"srd-5.2.1:content/spells/eyebite","source_page":129,"scope":"complete-source-record","level":6,"school":"necromancy","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"save","mechanic":"save.eyebite-wisdom","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"asleep","mechanic":"condition.eyebite-unconscious","timing":"failed-save-if-asleep-chosen","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"panicked","mechanic":"condition.eyebite-frightened-dash-away","timing":"failed-save-if-panicked-chosen","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sickened","mechanic":"condition.eyebite-poisoned","timing":"failed-save-if-sickened-chosen","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"retarget","mechanic":"target.eyebite-magic-action-retarget","timing":"each-of-your-turns-until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["save.eyebite-wisdom","condition.eyebite-unconscious","condition.eyebite-frightened-dash-away","condition.eyebite-poisoned","target.eyebite-magic-action-retarget","concentration.source-duration"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","eyebite":{"appearance":"your-eyes-become-an-inky-void","target_range_feet":60,"target_origin":"within-60-feet-of-you","target_visibility":"you-can-see","save":"wisdom","failure":"affected-by-one-chosen-effect-for-the-duration","options":["asleep","panicked","sickened"],"retarget_action":"magic-action","retarget_timing":"each-of-your-turns-until-the-spell-ends","retarget_exclusion":"creature-that-succeeded-on-a-save-against-this-casting","asleep_condition":"unconscious","asleep_wakes_when":["it-takes-any-damage","another-creature-takes-an-action-to-shake-it-awake"],"panicked_condition":"frightened","panicked_turn_action":"take-dash-action-and-move-away-from-you-by-safest-and-shortest-route","panicked_end_distance_feet":60,"panicked_end_when":"moves-to-a-space-at-least-60-feet-away-where-it-cannot-see-you","sickened_condition":"poisoned"}}}`,
		`{"id":"srd-5.2.1:content/spells/fabricate","name":"Fabricate","source_id":"srd-5.2.1:content/spells/fabricate","source_page":129,"scope":"complete-source-record","level":4,"school":"transmutation","lists":["wizard"],"casting_time":"10-minutes","range":{"kind":"feet","feet":120},"target":{"shape":"choice","count":1,"eligibility":"visible-raw-materials-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"convert","mechanic":"craft.fabricate-same-material","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"size","mechanic":"limit.fabricate-object-size","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"restriction","mechanic":"restriction.fabricate-skilled-items-need-tool-proficiency","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["craft.fabricate-same-material","limit.fabricate-object-size","restriction.fabricate-skilled-items-need-tool-proficiency","cast.ten-minutes"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"instantaneous","fabricate":{"input":"raw-materials-you-can-see-within-range","output":"products-of-the-same-material","examples":["wooden-bridge-from-a-clump-of-trees","rope-from-a-patch-of-hemp","clothes-from-flax-or-wool"],"object_size":"large-or-smaller","object_cube_feet":10,"connected_cube_count":8,"connected_cube_feet":5,"quantity":"sufficient-quantity-of-material","mineral_kinds":["metal","stone","other-mineral-substance"],"mineral_object_size":"medium-or-smaller","mineral_cube_feet":5,"quality":"based-on-the-quality-of-the-raw-materials","cannot_create":["creatures","magic-items"],"skilled_items":"items-requiring-a-high-degree-of-skill-such-as-weapons-and-armor","skilled_items_exception":"proficiency-with-the-type-of-artisans-tools-used-to-craft-such-objects"}}}`,
		`{"id":"srd-5.2.1:content/spells/faerie-fire","name":"Faerie Fire","source_id":"srd-5.2.1:content/spells/faerie-fire","source_page":129,"scope":"complete-source-record","level":1,"school":"evocation","lists":["bard","druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"cube","count":0,"eligibility":"objects-and-creatures-in-cube"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"outline","mechanic":"outline.faerie-fire-objects-and-failed-save-creatures","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"light","mechanic":"light.dim-10-foot-radius","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"invisibility","mechanic":"condition.faerie-fire-blocks-invisible","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"advantage","mechanic":"roll.attack-advantage-if-attacker-sees-target","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["outline.faerie-fire-objects-and-failed-save-creatures","light.dim-10-foot-radius","condition.faerie-fire-blocks-invisible","roll.attack-advantage-if-attacker-sees-target","save.dexterity","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","faerie_fire":{"cube_feet":20,"objects_outlined":"all-objects-in-the-cube","creature_outline":"each-creature-in-the-cube-that-fails-the-save","save":"dexterity","outline_colors":["blue","green","violet"],"color_choice":"your-choice","light_kind":"dim-light","light_radius_feet":10,"light_source":"objects-and-affected-creatures","blocked_condition":"invisible","attack_roll_effect":"advantage-if-the-attacker-can-see-it","attack_roll_target":"affected-creature-or-object"}}}`,
		`{"id":"srd-5.2.1:content/spells/faithful-hound","name":"Faithful Hound","source_id":"srd-5.2.1:content/spells/faithful-hound","source_page":129,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a silver whistle","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"hound","mechanic":"summon.faithful-hound","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"alarm","mechanic":"alarm.faithful-hound-bark","timing":"small-or-larger-creature-within-30-feet-without-password","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"bite","mechanic":"damage.faithful-hound-bite","timing":"start-of-each-of-your-turns","dice_count":4,"die_sides":8,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"move","mechanic":"movement.faithful-hound-magic-action","timing":"magic-action-on-your-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.faithful-hound","alarm.faithful-hound-bark","damage.faithful-hound-bite","movement.faithful-hound-magic-action","save.dexterity","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-caster-and-hound-more-than-300-feet-apart","hound":{"form":"phantom-watchdog","placement":"unoccupied-space-you-can-see-within-range","ends_when":["the-duration-ends","you-and-the-hound-are-more-than-300-feet-apart"],"separation_feet":300,"visible_to":"only-you","traits":["intangible","invulnerable"],"bark_trigger":"small-or-larger-creature-comes-within-30-feet-without-first-speaking-the-password","bark_trigger_feet":30,"password":"specified-by-you-when-you-cast-the-spell","bark_result":"hound-starts-barking-loudly","truesight_feet":30,"bite_timing":"start-of-each-of-your-turns","bite_target":"one-enemy-within-5-feet-of-the-hound","bite_reach_feet":5,"bite_save":"dexterity","bite_failure_damage_dice_count":4,"bite_failure_damage_die_sides":8,"bite_failure_damage_type":"force","move_action":"magic-action","move_timing":"your-later-turns","move_feet":30}}}`,
		`{"id":"srd-5.2.1:content/spells/false-life","name":"False Life","source_id":"srd-5.2.1:content/spells/false-life","source_page":129,"scope":"complete-source-record","level":1,"school":"necromancy","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of alcohol","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-bonus","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["temporary-hit-points"]},"effects":[{"id":"temporary-hit-points","mechanic":"temporary-hit-points.false-life","timing":"immediate","dice_count":2,"die_sides":4,"damage_type":"","bonus":4,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["temporary-hit-points.false-life","scaling.bonus-per-slot","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","false_life":{"recipient":"you","temporary_hit_points_dice_count":2,"temporary_hit_points_die_sides":4,"temporary_hit_points_bonus":4,"additional_temporary_hit_points_per_slot":5}}}`,
		`{"id":"srd-5.2.1:content/spells/fear","name":"Fear","source_id":"srd-5.2.1:content/spells/fear","source_page":130,"scope":"complete-source-record","level":3,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"cone","count":0,"eligibility":"creatures-in-cone"},"components":{"verbal":true,"somatic":true,"material":{"description":"a white feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"save","mechanic":"save.fear-wisdom","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"drop","mechanic":"drop.fear-held-items","timing":"failed-save-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"frightened","mechanic":"condition.fear-frightened-dash-away","timing":"failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"recheck","mechanic":"save.fear-end-turn-without-line-of-sight","timing":"end-of-creature-turn-without-line-of-sight","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["save.fear-wisdom","drop.fear-held-items","condition.fear-frightened-dash-away","save.fear-end-turn-without-line-of-sight","area.cone-30","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","fear":{"cone_feet":30,"origin":"you","save":"wisdom","failure_effects":["drop-whatever-it-is-holding","frightened-condition-for-the-duration"],"condition":"frightened","turn_action":"take-dash-action-and-move-away-from-you-by-the-safest-route","turn_action_exception":"nowhere-to-move","end_turn_trigger":"ends-its-turn-in-a-space-without-line-of-sight-to-you","end_turn_save":"wisdom","end_turn_save_success":"spell-ends-on-that-creature"}}}`,
		`{"id":"srd-5.2.1:content/spells/feather-fall","name":"Feather Fall","source_id":"srd-5.2.1:content/spells/feather-fall","source_page":130,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["bard","sorcerer","wizard"],"casting_time":"reaction","range":{"kind":"feet","feet":60},"target":{"shape":"up-to-count","count":5,"eligibility":"falling-creatures-in-range"},"components":{"verbal":true,"somatic":false,"material":{"description":"a small feather or piece of down","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"descent","mechanic":"descent.feather-fall-60-feet-per-round","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"landing","mechanic":"landing.feather-fall-no-damage-ends-for-creature","timing":"creature-lands-before-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["descent.feather-fall-60-feet-per-round","landing.feather-fall-no-damage-ends-for-creature","cast.reaction-trigger","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["you-or-a-creature-you-can-see-within-60-feet-of-you-falls"],"end_timing":"duration-ended-or-creature-lands","feather_fall":{"trigger_range_feet":60,"trigger":"you-or-a-creature-you-can-see-within-range-falls","max_targets":5,"target_state":"falling-creatures-within-range","descent_feet_per_round":60,"descent_until":"spell-ends","landing_result":"no-fall-damage-and-spell-ends-for-that-creature"}}}`,
		`{"id":"srd-5.2.1:content/spells/find-familiar","name":"Find Familiar","source_id":"srd-5.2.1:content/spells/find-familiar","source_page":130,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["wizard"],"casting_time":"1-hour","range":{"kind":"feet","feet":10},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"burning incense worth 10+ GP, which the spell consumes","cost_cp":1000,"consumed":true}},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"familiar","mechanic":"summon.find-familiar","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"telepathy","mechanic":"communication.familiar-telepathy-100-feet","timing":"while-familiar-within-100-feet","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"senses","mechanic":"sense.familiar-eyes-and-ears-bonus-action","timing":"bonus-action-until-start-of-your-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"touch","mechanic":"delivery.familiar-touch-spells","timing":"reaction-when-you-cast-a-touch-spell-with-familiar-within-100-feet","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"combat","mechanic":"combat.familiar-ally-own-initiative","timing":"while-familiar-exists","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disappearance","mechanic":"state.familiar-disappears-at-zero-hit-points","timing":"familiar-drops-to-0-hit-points","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dismissal","mechanic":"action.familiar-dismiss-and-recall","timing":"magic-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"limit","mechanic":"limit.one-familiar-only","timing":"when-cast-while-you-have-a-familiar","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.find-familiar","communication.familiar-telepathy-100-feet","sense.familiar-eyes-and-ears-bonus-action","delivery.familiar-touch-spells","combat.familiar-ally-own-initiative","state.familiar-disappears-at-zero-hit-points","action.familiar-dismiss-and-recall","limit.one-familiar-only","cast.one-hour","cast.ritual","components.consumed-priced-material","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"instantaneous","familiar":{"form_choices":["bat","cat","frog","hawk","lizard","octopus","owl","rat","raven","spider","weasel","other-beast-with-challenge-rating-0"],"appearance":"unoccupied-space-within-range","statistics":"stat-block-of-the-chosen-form","creature_type_choices":["celestial","fey","fiend"],"creature_type_replaces":"beast","independence":"acts-independently-but-obeys-your-commands","telepathy_range_feet":100,"see_through_action":"bonus-action","see_through_duration":"until-the-start-of-your-next-turn","see_through_benefit":"benefits-of-any-special-senses-it-has","touch_delivery":"familiar-delivers-spells-with-a-range-of-touch","touch_delivery_range_feet":100,"touch_delivery_action":"reaction","touch_delivery_timing":"when-you-cast-the-spell","combat_role":"ally-to-you-and-your-allies","combat_initiative":"rolls-its-own-initiative-and-acts-on-its-own-turn","attack_restriction":"cannot-attack-but-can-take-other-actions","zero_hit_points":"disappears","reappears_when":"you-cast-this-spell-again","dismiss_action":"magic-action","dismiss_options":["temporarily-to-a-pocket-dimension","forever"],"recall_action":"magic-action-while-temporarily-dismissed","recall_range_feet":30,"recall_space":"unoccupied-space","leaves_behind":"anything-it-was-wearing-or-carrying","leaves_behind_when":["drops-to-0-hit-points","disappears-into-the-pocket-dimension"],"familiar_limit":"cannot-have-more-than-one-familiar-at-a-time","recast_result":"existing-familiar-adopts-a-new-eligible-form"}}}`,
		`{"id":"srd-5.2.1:content/spells/find-steed","name":"Find Steed","source_id":"srd-5.2.1:content/spells/find-steed","source_page":130,"scope":"complete-source-record","level":2,"school":"conjuration","lists":["paladin"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"steed","mechanic":"summon.find-steed","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"replacement","mechanic":"limit.one-steed-replaced-by-new-one","timing":"when-cast-while-you-have-a-steed","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"combat","mechanic":"combat.steed-ally-shared-initiative-controlled-mount","timing":"while-steed-exists","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disappearance","mechanic":"state.steed-disappears-at-zero-hit-points-or-caster-death","timing":"steed-drops-to-0-hit-points-or-you-die","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"slot-level","mechanic":"stat-block.slot-level-as-spell-level","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.find-steed","limit.one-steed-replaced-by-new-one","combat.steed-ally-shared-initiative-controlled-mount","state.steed-disappears-at-zero-hit-points-or-caster-death","stat-block.otherworldly-steed","stat-block.slot-level-as-spell-level"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","steed":{"summons":"otherworldly-being-that-appears-as-a-loyal-steed","placement":"unoccupied-space-of-your-choice-within-range","stat_block":"otherworldly-steed","replaces_existing":"steed-from-this-spell-is-replaced-by-the-new-one","appearance":"large-rideable-animal-of-your-choice","appearance_examples":["horse","camel","dire-wolf","elk"],"creature_type_choices":["celestial","fey","fiend"],"creature_type_timing":"whenever-you-cast-the-spell","combat_role":"ally-to-you-and-your-allies","initiative":"shares-your-initiative-count","mount_function":"controlled-mount-while-you-ride-it","incapacitated_behavior":"takes-its-turn-immediately-after-yours-and-acts-independently-focusing-on-protecting-you","disappears_when":["drops-to-0-hit-points","you-die"],"leaves_behind":"anything-it-was-wearing-or-carrying","recast_choice":"summon-the-steed-that-disappeared-or-a-different-one","slot_level_use":"spell-slot-level-is-the-spell-level-in-the-stat-block"}}}`,
		`{"id":"srd-5.2.1:content/spells/find-traps","name":"Find Traps","source_id":"srd-5.2.1:content/spells/find-traps","source_page":131,"scope":"complete-source-record","level":2,"school":"divination","lists":["cleric","druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"selected-objects","count":0,"eligibility":"traps-in-line-of-sight-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"presence","mechanic":"sense.find-traps-presence","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"nature","mechanic":"sense.find-traps-general-nature-of-danger","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.find-traps-presence","sense.find-traps-general-nature-of-danger"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","traps":{"senses":"any-trap-within-range-and-line-of-sight","trap_definition":"object-or-mechanism-created-to-cause-damage-or-other-danger","sensed_examples":["alarm-spell","glyph-of-warding-spell","mechanical-pit-trap"],"not_sensed_examples":["natural-weakness-in-the-floor","unstable-ceiling","hidden-sinkhole"],"reveals":"that-a-trap-is-present","does_not_reveal":"trap-location","learns":"general-nature-of-the-danger-posed-by-a-sensed-trap"}}}`,
	}
	batch := loadSpells205(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Etherealness", "Expeditious Retreat", "Eyebite", "Fabricate", "Faerie Fire", "Faithful Hound", "False Life", "Fear", "Feather Fall", "Find Familiar", "Find Steed", "Find Traps"}
	pages := []int{128, 129, 129, 129, 129, 129, 129, 130, 130, 130, 130, 131}
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
			// Every family here is level 1 or higher, so each lists its slot levels.
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

func TestSpells205CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells205(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 108 {
		t.Fatalf("expected nine twelve-family batches, got %d", len(all.Records))
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
func TestSpells205MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_205.json")
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

func TestSpells205JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_205.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"cone_feet": 30`, `"cone_feet": 30, "cone_feet": 31`},
		{`"cone_feet": 30`, `"cone_feet": 30, "CONE_FEET": 31`},
		{`"cone_feet": 30`, `"cone_feet": 30, "unknown": true`},
		{`"outline_colors": [
            "blue",
            "green",
            "violet"
          ]`, `"outline_colors": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/etherealness"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "targets"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [
        7,
        8,
        9
      ],`, ``},
		{`"slot_levels": [
        7,
        8,
        9
      ],`, `"slot_levels": null,`},
		{`"slot_levels": [
        7,
        8,
        9
      ],`, `"slot_levels": [0],`},
		{`"casting_time": "action"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 0`, `"casting_seconds": 60`},
		{`"feet": 120`, `"feet": 121`},
		{`"cost_cp": 1000`, `"cost_cp": 1001`},
		{`"ritual": true`, `"ritual": false`},
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

func TestSpells205InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"etherealness extra targets", 0, func(s *SpellRecord) { s.SourceDetails.Etherealness.AdditionalTargetsPerSlot = 2 }},
		{"etherealness base targets", 0, func(s *SpellRecord) { s.SourceDetails.Etherealness.BaseTargets = 2 }},
		{"etherealness scaling level", 0, func(s *SpellRecord) { s.Scaling.AboveLevel = 6 }},
		{"etherealness shunt damage", 0, func(s *SpellRecord) { s.SourceDetails.Etherealness.ShuntDamageType = "necrotic" }},
		{"etherealness slot levels", 0, func(s *SpellRecord) { s.SlotLevels = []int{6, 7, 8, 9} }},
		{"etherealness concentration", 0, func(s *SpellRecord) { *s.Concentration = true }},
		{"retreat bonus action", 1, func(s *SpellRecord) { s.SourceDetails.Retreat.RepeatActionType = "action" }},
		{"retreat casting time", 1, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"eyebite save", 2, func(s *SpellRecord) { s.SourceDetails.Eyebite.Save = "constitution" }},
		{"eyebite option", 2, func(s *SpellRecord) { s.SourceDetails.Eyebite.Options[2] = "blinded" }},
		{"eyebite end distance", 2, func(s *SpellRecord) { s.SourceDetails.Eyebite.PanickedEndDistanceFeet = 30 }},
		{"fabricate casting seconds", 3, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 60 }},
		{"fabricate mineral cube", 3, func(s *SpellRecord) { s.SourceDetails.Fabricate.MineralCubeFeet = 10 }},
		{"fabricate target", 3, func(s *SpellRecord) { s.Target.Eligibility = "visible-creature-in-range" }},
		{"faerie fire cube", 4, func(s *SpellRecord) { s.SourceDetails.FaerieFire.CubeFeet = 10 }},
		{"faerie fire colors", 4, func(s *SpellRecord) { s.SourceDetails.FaerieFire.OutlineColors[2] = "red" }},
		{"faerie fire components", 4, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"hound damage dice", 5, func(s *SpellRecord) { s.SourceDetails.Hound.BiteFailureDamageDiceCount = 2 }},
		{"hound effect dice", 5, func(s *SpellRecord) { s.Effects[2].DiceCount = 2 }},
		{"hound separation", 5, func(s *SpellRecord) { s.SourceDetails.Hound.SeparationFeet = 30 }},
		{"hound material", 5, func(s *SpellRecord) { s.Components.Material.Description = "a bronze whistle" }},
		{"false life bonus", 6, func(s *SpellRecord) { s.SourceDetails.FalseLife.TemporaryHitPointsBonus = 2 }},
		{"false life scaling kind", 6, func(s *SpellRecord) { s.Scaling.Kind = "slot" }},
		{"false life increment", 6, func(s *SpellRecord) { s.SourceDetails.FalseLife.AdditionalTemporaryHitPointsPerSlot = 0 }},
		{"false life effect bonus", 6, func(s *SpellRecord) { s.Effects[0].Bonus = 5 }},
		{"fear cone", 7, func(s *SpellRecord) { s.SourceDetails.Fear.ConeFeet = 15 }},
		{"fear end-turn save", 7, func(s *SpellRecord) { s.SourceDetails.Fear.EndTurnSave = "charisma" }},
		{"feather fall trigger", 8, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{"hit-by-attack-roll"} }},
		{"feather fall no trigger", 8, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{} }},
		{"feather fall targets", 8, func(s *SpellRecord) { s.Target.Count = 3 }},
		{"feather fall descent", 8, func(s *SpellRecord) { s.SourceDetails.FeatherFall.DescentFeetPerRound = 30 }},
		{"familiar form", 9, func(s *SpellRecord) { s.SourceDetails.Familiar.FormChoices[11] = "other-beast" }},
		{"familiar material cost", 9, func(s *SpellRecord) { *s.Components.Material.CostCP = 100 }},
		{"familiar consumed", 9, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"familiar ritual", 9, func(s *SpellRecord) { *s.Ritual = false }},
		{"familiar casting seconds", 9, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 600 }},
		{"steed scaling", 10, func(s *SpellRecord) { s.Scaling.Kind = "slot-bonus" }},
		{"steed stat block", 10, func(s *SpellRecord) { s.SourceDetails.Steed.StatBlock = "warhorse" }},
		{"steed required", 10, func(s *SpellRecord) { s.RequiredMechanics = s.RequiredMechanics[:5] }},
		{"traps range", 11, func(s *SpellRecord) { s.Range.Feet = 60 }},
		{"traps not sensed", 11, func(s *SpellRecord) { s.SourceDetails.Traps.NotSensedExamples[1] = "pit" }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 11, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Fear = &SpellFear{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells205(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells205RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_205.json")
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

func TestSpells205WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_205.json")
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
		`"damage_type": "force", "damage_type": "fire"`,
		`"damage_type": "force", "DAMAGE_TYPE": "fire"`,
		`"damage_type": "force", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"damage_type": "force"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells205PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_205.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"perceived_plane_limit_feet", 1, 60}, {"vertical_extra_feet_per_foot", 1, 1}, {"shunt_damage_per_foot_moved", 1, 2},
		{"additional_targets_per_slot", 1, 3}, {"slot_target_range_feet", 1, 10}, {"target_range_feet", 1, 60},
		{"panicked_end_distance_feet", 1, 60}, {"object_cube_feet", 1, 10}, {"connected_cube_count", 1, 8}, {"connected_cube_feet", 1, 5},
		{"mineral_cube_feet", 1, 5}, {"cube_feet", 1, 20}, {"light_radius_feet", 1, 10}, {"separation_feet", 1, 300}, {"bark_trigger_feet", 1, 30},
		{"truesight_feet", 1, 30}, {"bite_reach_feet", 1, 5}, {"bite_failure_damage_dice_count", 1, 4}, {"bite_failure_damage_die_sides", 1, 8},
		{"move_feet", 1, 30}, {"temporary_hit_points_dice_count", 1, 2}, {"temporary_hit_points_die_sides", 1, 4},
		{"temporary_hit_points_bonus", 1, 4}, {"additional_temporary_hit_points_per_slot", 1, 5}, {"cone_feet", 1, 30},
		{"trigger_range_feet", 1, 60}, {"max_targets", 1, 5}, {"descent_feet_per_round", 1, 60}, {"telepathy_range_feet", 1, 100},
		{"touch_delivery_range_feet", 1, 100}, {"recall_range_feet", 1, 30}, {"base_targets", 1, 1},
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
// inventory assigns to catalogue issue #205.
func TestSpells205MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 205 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells205(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells205ChangesNoteAndFamilyBoundaries(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_205.json")
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
	want := "SRD 5.2.1 PDF pp.128-131 adapted into twelve complete offline spell families for issue #205. No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	// Each family sits on one page; the Otherworldly Steed stat block and the
	// tail of Find the Path on p.131 belong to other families and are not imported.
	for _, r := range loadSpells205(t).Records {
		if r.SourceDetails.EndTiming == "" || strings.Contains(r.Name, "Otherworldly") || strings.Contains(r.Name, "Find the Path") {
			t.Fatalf("unexpected family %s", r.Name)
		}
	}
	steed := loadSpells205(t).Records[10]
	if !slices.Contains(steed.RequiredMechanics, "stat-block.otherworldly-steed") || steed.Scaling.Kind != "none" {
		t.Fatal("Find Steed must keep its stat block as an unavailable dependency")
	}
}
