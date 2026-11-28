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

func loadSpells223(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_223.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.172-174 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells223SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/wall-of-fire","name":"Wall of Fire","source_id":"srd-5.2.1:content/spells/wall-of-fire","source_page":172,"scope":"complete-source-record","level":4,"school":"evocation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"creatures-in-wall-of-fire"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of charcoal","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"slot","above_level":4,"dice_count":1,"die_sides":8,"effect_ids":["appear","side"]},"effects":[{"id":"appear","mechanic":"damage.wall-of-fire-appear","timing":"when-wall-appears","dice_count":5,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"side","mechanic":"damage.wall-of-fire-side","timing":"on-end-turn-near-chosen-side-or-inside-or-first-entry","dice_count":5,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.wall-of-fire-appear","damage.wall-of-fire-side","save.dexterity-half","concentration.source-duration","scaling.per-slot","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","wall_of_fire":{"surface":"solid-surface-within-range","max_length_feet":60,"height_feet":20,"thickness_feet":1,"ring_max_diameter_feet":20,"ring_height_feet":20,"ring_thickness_feet":1,"opacity":"opaque","appear_save":"dexterity","appear_dice_count":5,"appear_die_sides":8,"damage_type":"fire","appear_success":"half-damage","damaging_side":"one-side-selected-at-casting","side_dice_count":5,"side_die_sides":8,"side_distance_feet":10,"side_triggers":["ends-turn-within-distance-of-chosen-side","ends-turn-inside-wall","enters-wall-first-time-on-turn"],"other_side":"deals-no-damage","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/wall-of-force","name":"Wall of Force","source_id":"srd-5.2.1:content/spells/wall-of-force","source_page":172,"scope":"complete-source-record","level":5,"school":"evocation","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"force-wall-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a shard of glass","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"barrier","mechanic":"barrier.wall-of-force","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.wall-of-force","immunity.wall-of-force","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-destroyed-by-disintegrate","wall_of_force":{"visibility":"invisible","orientations":["horizontal","vertical","angled"],"support":["free-floating","resting-on-solid-surface"],"shapes":["hemispherical-dome","globe","flat-surface-of-panels"],"max_radius_feet":10,"panel_count":10,"panel_length_feet":10,"panel_width_feet":10,"panel_rule":"each-contiguous-with-another-panel","thickness_quarter_inches":1,"cut_creature":"pushed-to-side-caster-chooses","passage":"nothing-can-physically-pass-through","damage_immunity":"all-damage","dispel_immunity":"cannot-be-dispelled-by-dispel-magic","destroyed_by":"disintegrate-instantly","ethereal":"extends-into-ethereal-plane-and-blocks-ethereal-travel"}}}`,
		`{"id":"srd-5.2.1:content/spells/wall-of-ice","name":"Wall of Ice","source_id":"srd-5.2.1:content/spells/wall-of-ice","source_page":172,"scope":"complete-source-record","level":6,"school":"evocation","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"creatures-in-ice-wall"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of quartz","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":600,"duration":"timed","scaling":{"kind":"slot-effect-dice","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["appear","frigid-air"]},"effects":[{"id":"appear","mechanic":"damage.wall-of-ice-appear","timing":"when-wall-appears-through-creature-space","dice_count":10,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"frigid-air","mechanic":"damage.wall-of-ice-frigid-air","timing":"on-first-movement-through-sheet-on-turn","dice_count":5,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.wall-of-ice-appear","damage.wall-of-ice-frigid-air","barrier.wall-of-ice","save.dexterity-half","save.constitution-half","concentration.source-duration","scaling.per-slot-by-effect","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","wall_of_ice":{"source_pages":[172,173],"surface":"solid-surface-within-range","shapes":["hemispherical-dome","globe","flat-surface-of-panels"],"max_radius_feet":10,"panel_count":10,"panel_square_feet":10,"panel_rule":"each-contiguous-with-another-panel","thickness_feet":1,"cut_creature":"pushed-to-side-caster-chooses","appear_save":"dexterity","appear_dice_count":10,"appear_die_sides":6,"damage_type":"cold","appear_success":"half-damage","armor_class":12,"hit_points_per_section":30,"section_feet":10,"immunities":["cold","poison","psychic"],"vulnerabilities":["fire"],"breach":"section-at-0-hit-points-destroyed-leaving-sheet-of-frigid-air","frigid_air_trigger":"creature-moving-through-sheet-first-time-on-turn","frigid_air_save":"constitution","frigid_air_dice_count":5,"frigid_air_die_sides":6,"frigid_air_success":"half-damage","slot_dice":[{"effect_id":"appear","dice_count":2,"die_sides":6},{"effect_id":"frigid-air","dice_count":1,"die_sides":6}]}}}`,
		`{"id":"srd-5.2.1:content/spells/wall-of-stone","name":"Wall of Stone","source_id":"srd-5.2.1:content/spells/wall-of-stone","source_page":173,"scope":"complete-source-record","level":5,"school":"evocation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"stone-wall-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a cube of granite","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"barrier","mechanic":"barrier.wall-of-stone","timing":"for-duration-or-permanent-after-full-concentration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.wall-of-stone","save.wall-of-stone-escape","duration.wall-of-stone-becomes-permanent","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended-permanent-after-full-concentration","wall_of_stone":{"composition":"nonmagical-solid-stone","panel_count":10,"thickness_inches":6,"panel_length_feet":10,"panel_height_feet":10,"thin_panel_length_feet":10,"thin_panel_height_feet":20,"thin_panel_thickness_inches":3,"panel_rule":"each-contiguous-with-another-panel","cut_creature":"pushed-to-side-caster-chooses","enclosure_trigger":"creature-would-be-surrounded-on-all-sides-by-wall-or-wall-and-solid-surface","enclosure_save":"dexterity","enclosure_success":"can-use-reaction-to-move-up-to-speed-out-of-enclosure","shape":"any-shape-not-in-creature-or-object-space","orientation":"need-not-be-vertical-or-on-firm-foundation","anchor":"must-merge-with-and-be-solidly-supported-by-existing-stone","uses":["bridge-a-chasm","create-a-ramp"],"long_span_over_feet":20,"long_span_panel_size_divisor":2,"long_span_purpose":"create-supports","crude_shaping":["battlements-and-the-like"],"armor_class":15,"hit_points_per_inch_of_thickness":30,"immunities":["poison","psychic"],"breach":"panel-at-0-hit-points-destroyed","collapse":"connected-panels-might-collapse-at-gm-discretion","full_concentration_seconds":600,"full_concentration_result":"permanent-and-cannot-be-dispelled","otherwise":"disappears-when-spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/wall-of-thorns","name":"Wall of Thorns","source_id":"srd-5.2.1:content/spells/wall-of-thorns","source_page":173,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"creatures-in-wall-of-thorns"},"components":{"verbal":true,"somatic":true,"material":{"description":"a handful of thorns","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":600,"duration":"timed","scaling":{"kind":"slot","above_level":6,"dice_count":1,"die_sides":8,"effect_ids":["appear","passage"]},"effects":[{"id":"appear","mechanic":"damage.wall-of-thorns-appear","timing":"when-wall-appears","dice_count":7,"die_sides":8,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"passage","mechanic":"damage.wall-of-thorns-passage","timing":"on-first-entry-or-end-of-turn-in-wall","dice_count":7,"die_sides":8,"damage_type":"slashing","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.wall-of-thorns-appear","damage.wall-of-thorns-passage","movement.wall-of-thorns","save.dexterity-half","concentration.source-duration","scaling.per-slot","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","wall_of_thorns":{"composition":"tangled-brush-with-needle-sharp-thorns","surface":"solid-surface-within-range","max_length_feet":60,"straight_height_feet":10,"straight_thickness_feet":5,"circle_diameter_feet":20,"circle_max_height_feet":20,"circle_thickness_feet":5,"sight":"blocks-line-of-sight","appear_save":"dexterity","appear_dice_count":7,"appear_die_sides":8,"appear_damage_type":"piercing","appear_success":"half-damage","movement_cost_feet_per_foot":4,"passage_triggers":["enters-space-in-wall-first-time-on-turn","ends-turn-there"],"passage_save":"dexterity","passage_dice_count":7,"passage_die_sides":8,"passage_damage_type":"slashing","passage_success":"half-damage","passage_save_limit":"once-per-turn","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/warding-bond","name":"Warding Bond","source_id":"srd-5.2.1:content/spells/warding-bond","source_page":173,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"other-willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pair of platinum rings worth 50+ GP each, which you and the target must wear for the duration","cost_cp":10000,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"armor-class","mechanic":"defense.warding-bond-armor-class","timing":"while-target-within-range-of-caster","dice_count":0,"die_sides":0,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false},{"id":"saving-throws","mechanic":"defense.warding-bond-saving-throws","timing":"while-target-within-range-of-caster","dice_count":0,"die_sides":0,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false},{"id":"resistance","mechanic":"defense.warding-bond-resistance","timing":"while-target-within-range-of-caster","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"shared-damage","mechanic":"damage.warding-bond-shared","timing":"each-time-target-takes-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.warding-bond-armor-class","defense.warding-bond-saving-throws","defense.warding-bond-resistance","damage.warding-bond-shared","ending.warding-bond","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-caster-at-0-hit-points-separation-or-recast-on-either-creature","warding_bond":{"target_requirement":"another-willing-creature","connection":"mystic-connection-until-spell-ends","benefit_range_feet":60,"armor_class_bonus":1,"saving_throw_bonus":1,"resistance":"all-damage","shared_damage":"caster-takes-same-amount-each-time-target-takes-damage","ring_count":2,"ring_minimum_cost_cp":5000,"ring_wearers":["caster","target"],"separation_over_feet":60,"early_end":["caster-drops-to-0-hit-points","caster-and-target-separated-beyond-distance","spell-cast-again-on-either-connected-creature"]}}}`,
		`{"id":"srd-5.2.1:content/spells/water-breathing","name":"Water Breathing","source_id":"srd-5.2.1:content/spells/water-breathing","source_page":173,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":10,"eligibility":"chosen-willing-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a short reed","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration_seconds":86400,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"breathing","mechanic":"utility.water-breathing","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.water-breathing","cast.ritual","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","water_breathing":{"max_creatures":10,"grants":"ability-to-breathe-underwater-until-spell-ends","normal_respiration":"retained"}}}`,
		`{"id":"srd-5.2.1:content/spells/water-walk","name":"Water Walk","source_id":"srd-5.2.1:content/spells/water-walk","source_page":173,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["cleric","druid","ranger","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":10,"eligibility":"chosen-willing-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of cork","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"walk","mechanic":"movement.water-walk","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.water-walk","cast.ritual","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","water_walk":{"source_pages":[173,174],"max_creatures":10,"surfaces":"any-liquid-surface","surface_examples":["water","acid","mud","snow","quicksand","lava"],"as_if":"harmless-solid-ground","molten_lava":"crossing-creatures-can-still-take-damage-from-heat","transition_action":"bonus-action","transitions":["from-liquid-surface-into-liquid","from-liquid-onto-surface"],"falling":"falling-target-passes-through-surface-into-liquid"}}}`,
		`{"id":"srd-5.2.1:content/spells/web","name":"Web","source_id":"srd-5.2.1:content/spells/web","source_page":174,"scope":"complete-source-record","level":2,"school":"conjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"cube","count":0,"eligibility":"webs-in-cube-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of spiderweb","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"terrain.web","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"restrained","mechanic":"condition.web-restrained","timing":"on-first-entry-or-start-of-turn-in-webs","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"fire","mechanic":"damage.web-fire","timing":"on-start-of-turn-in-burning-webs","dice_count":2,"die_sides":4,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.web","condition.web-restrained","damage.web-fire","save.dexterity-negates","check.web-escape","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-unanchored-collapse","web":{"cube_feet":20,"terrain":"difficult-terrain","obscurement":"lightly-obscured","anchors":["between-two-solid-masses","layered-across-floor-wall-or-ceiling"],"solid_mass_examples":["walls","trees"],"unanchored":"web-collapses-and-spell-ends-at-start-of-casters-next-turn","flat_surface_depth_feet":5,"save":"dexterity","save_triggers":["enters-webs-first-time-on-turn","starts-turn-there"],"failure_condition":"restrained","condition_ends":"no-longer-in-webs-or-breaks-free","escape_action":"action","escape_check":"strength-athletics","escape_dc":"caster-spell-save-dc","escape_success":"no-longer-restrained","flammable_cube_feet":5,"burn_away_rounds":1,"fire_dice_count":2,"fire_die_sides":4,"fire_damage_type":"fire","fire_targets":"any-creature-that-starts-turn-in-fire"}}}`,
		`{"id":"srd-5.2.1:content/spells/weird","name":"Weird","source_id":"srd-5.2.1:content/spells/weird","source_page":174,"scope":"complete-source-record","level":9,"school":"illusion","lists":["warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"selected-creatures","count":0,"eligibility":"chosen-creatures-in-sphere-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.weird","timing":"immediate","dice_count":10,"die_sides":10,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"frightened","mechanic":"condition.weird-frightened","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"damage.weird-repeat","timing":"on-failed-end-of-turn-save-while-frightened","dice_count":5,"die_sides":10,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.weird","condition.weird-frightened","damage.weird-repeat","save.wisdom-half","save.weird-repeat","concentration.source-duration"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","weird":{"illusion":"illusory-terrors-in-others-minds","radius_feet":30,"save":"wisdom","damage_dice_count":10,"damage_die_sides":10,"damage_type":"psychic","failure_condition":"frightened","condition_duration":"for-spell-duration","save_success":"half-damage-only","repeat_save":"wisdom","repeat_timing":"end-of-each-of-frightened-targets-turns","repeat_dice_count":5,"repeat_die_sides":10,"repeat_success":"spell-ends-on-that-target"}}}`,
		`{"id":"srd-5.2.1:content/spells/wind-walk","name":"Wind Walk","source_id":"srd-5.2.1:content/spells/wind-walk","source_page":174,"scope":"complete-source-record","level":6,"school":"transmutation","lists":["druid"],"casting_time":"1-minute","range":{"kind":"feet","feet":30},"target":{"shape":"self-and-up-to-count","count":10,"eligibility":"chosen-willing-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a candle","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":28800,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"cloud-form","mechanic":"transformation.wind-walk-cloud-form","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"stunned","mechanic":"condition.wind-walk-reverting-stunned","timing":"while-reverting-to-normal-form","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.wind-walk-cloud-form","condition.wind-walk-reverting-stunned","movement.wind-walk-descent","cast.one-minute","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","wind_walk":{"max_other_creatures":10,"form":"gaseous-appearing-as-wisps-of-cloud","fly_speed_feet":300,"hover":"can-hover","condition_immunities":["prone"],"damage_resistances":["bludgeoning","piercing","slashing"],"allowed_actions":["dash-action","magic-action-to-begin-reverting-to-normal-form"],"revert_seconds":60,"revert_condition":"stunned","return_to_cloud":"magic-action-then-transformation-until-spell-ends","return_seconds":60,"descent_trigger":"in-cloud-form-and-flying-when-effect-ends","descent_feet_per_round":60,"descent_seconds":60,"safe_landing":"descends-until-landing-safely","no_landing":"falls-remaining-distance-if-cannot-land-in-time"}}}`,
		`{"id":"srd-5.2.1:content/spells/wind-wall","name":"Wind Wall","source_id":"srd-5.2.1:content/spells/wind-wall","source_page":174,"scope":"complete-source-record","level":3,"school":"evocation","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"wall","count":0,"eligibility":"creatures-in-wind-wall"},"components":{"verbal":true,"somatic":true,"material":{"description":"a fan and a feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"appear","mechanic":"damage.wind-wall","timing":"when-wall-appears","dice_count":4,"die_sides":8,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"barrier","mechanic":"barrier.wind-wall","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.wind-wall","barrier.wind-wall","save.strength-half","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","wind_wall":{"origin":"rises-from-ground-at-point-in-range","max_length_feet":50,"max_height_feet":15,"thickness_feet":1,"shape":"one-continuous-path-along-ground","save":"strength","damage_dice_count":4,"damage_die_sides":8,"damage_type":"bludgeoning","save_success":"half-damage","kept_at_bay":["fog","smoke","other-gases"],"blocked":["small-or-smaller-flying-creatures","small-or-smaller-flying-objects"],"loose_materials":"loose-lightweight-materials-brought-into-wall-fly-upward","deflected":["arrows","bolts","other-ordinary-projectiles"],"deflection":"launched-at-targets-behind-wall-deflected-upward-and-miss-automatically","unaffected":["boulders-hurled-by-giants","boulders-hurled-by-siege-engines","similar-projectiles"],"gaseous_creatures":"cannot-pass-through"}}}`,
	}
	batch := loadSpells223(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Wall of Fire", "Wall of Force", "Wall of Ice", "Wall of Stone", "Wall of Thorns", "Warding Bond", "Water Breathing", "Water Walk", "Web", "Weird", "Wind Walk", "Wind Wall"}
	pages := []int{172, 172, 172, 173, 173, 173, 173, 173, 174, 174, 174, 174}
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

func TestSpells223CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells223(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t), loadSpells220(t), loadSpells221(t), loadSpells222(t), loadSpells223(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 336 {
		t.Fatalf("expected 336 records, got %d", len(all.Records))
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
func TestSpells223MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_223.json")
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

func TestSpells223JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_223.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"radius_feet": 30`, `"radius_feet": 30, "radius_feet": 31`},
		{`"radius_feet": 30`, `"radius_feet": 30, "RADIUS_FEET": 31`},
		{`"radius_feet": 30`, `"radius_feet": 30, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/wall-of-fire"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "appear",
          "side"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 120`, `"feet": 121`},
		{`"cost_cp": 10000`, `"cost_cp": 5000`},
		{`"concentration": true`, `"concentration": false`},
		{`"ritual": true`, `"ritual": false`},
		{`"consumed": false`, `"consumed": true`},
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

func TestSpells223InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"fire ring", 0, func(s *SpellRecord) { s.SourceDetails.WallOfFire.RingMaxDiameterFeet = 30 }},
		{"fire scaled effects", 0, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"appear"} }},
		{"force panels", 1, func(s *SpellRecord) { s.SourceDetails.WallOfForce.PanelCount = 12 }},
		{"ice slot dice", 2, func(s *SpellRecord) { s.SourceDetails.WallOfIce.SlotDice[1].DiceCount = 2 }},
		{"ice slot rows", 2, func(s *SpellRecord) { s.SourceDetails.WallOfIce.SlotDice = s.SourceDetails.WallOfIce.SlotDice[:1] }},
		{"ice armor class", 2, func(s *SpellRecord) { s.SourceDetails.WallOfIce.ArmorClass = 15 }},
		{"stone hit points", 3, func(s *SpellRecord) { s.SourceDetails.WallOfStone.HitPointsPerInchOfThickness = 15 }},
		{"thorns movement", 4, func(s *SpellRecord) { s.SourceDetails.WallOfThorns.MovementCostFeetPerFoot = 2 }},
		{"bond bonus", 5, func(s *SpellRecord) { s.Effects[0].Bonus = 2 }},
		{"bond rings", 5, func(s *SpellRecord) { *s.Components.Material.CostCP = 5000 }},
		{"breathing creatures", 6, func(s *SpellRecord) { s.Target.Count = 8 }},
		{"walk surfaces", 7, func(s *SpellRecord) {
			s.SourceDetails.WaterWalk.SurfaceExamples = s.SourceDetails.WaterWalk.SurfaceExamples[:5]
		}},
		{"walk wizard list", 7, func(s *SpellRecord) { s.Lists = append(s.Lists, "wizard") }},
		{"web cube", 8, func(s *SpellRecord) { s.SourceDetails.Web.CubeFeet = 15 }},
		{"web fire die", 8, func(s *SpellRecord) { s.Effects[2].DieSides = 6 }},
		{"weird repeat dice", 9, func(s *SpellRecord) { s.SourceDetails.Weird.RepeatDiceCount = 10 }},
		{"wind walk speed", 10, func(s *SpellRecord) { s.SourceDetails.WindWalk.FlySpeedFeet = 30 }},
		{"wind walk casting", 10, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"wind wall height", 11, func(s *SpellRecord) { s.SourceDetails.WindWall.MaxHeightFeet = 20 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 0, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 11, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.WindWall = &SpellWindWall{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells223(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells223RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_223.json")
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

func TestSpells223WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_223.json")
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
		`"radius_feet": 30, "radius_feet": 31`,
		`"radius_feet": 30, "RADIUS_FEET": 31`,
		`"radius_feet": 30, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"radius_feet": 30`), []byte(replacement), 1)
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
func TestSpells223PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_223.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"max_length_feet", 1, 60},
		{"height_feet", 1, 20},
		{"thickness_feet", 1, 1},
		{"ring_max_diameter_feet", 1, 20},
		{"ring_height_feet", 1, 20},
		{"ring_thickness_feet", 1, 1},
		{"appear_dice_count", 1, 5},
		{"appear_die_sides", 1, 8},
		{"side_dice_count", 1, 5},
		{"side_die_sides", 1, 8},
		{"side_distance_feet", 1, 10},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
		{"max_radius_feet", 1, 10},
		{"panel_count", 1, 10},
		{"panel_length_feet", 1, 10},
		{"panel_width_feet", 1, 10},
		{"thickness_quarter_inches", 1, 1},
		{"max_radius_feet", 2, 10},
		{"panel_count", 2, 10},
		{"panel_square_feet", 1, 10},
		{"thickness_feet", 2, 1},
		{"appear_dice_count", 2, 10},
		{"appear_die_sides", 2, 6},
		{"armor_class", 1, 12},
		{"hit_points_per_section", 1, 30},
		{"section_feet", 1, 10},
		{"frigid_air_dice_count", 1, 5},
		{"frigid_air_die_sides", 1, 6},
		{"dice_count", 9, 2},
		{"die_sides", 9, 6},
		{"dice_count", 10, 1},
		{"die_sides", 10, 6},
		{"panel_count", 3, 10},
		{"thickness_inches", 1, 6},
		{"panel_length_feet", 2, 10},
		{"panel_height_feet", 1, 10},
		{"thin_panel_length_feet", 1, 10},
		{"thin_panel_height_feet", 1, 20},
		{"thin_panel_thickness_inches", 1, 3},
		{"long_span_over_feet", 1, 20},
		{"long_span_panel_size_divisor", 1, 2},
		{"armor_class", 2, 15},
		{"hit_points_per_inch_of_thickness", 1, 30},
		{"full_concentration_seconds", 1, 600},
		{"max_length_feet", 2, 60},
		{"straight_height_feet", 1, 10},
		{"straight_thickness_feet", 1, 5},
		{"circle_diameter_feet", 1, 20},
		{"circle_max_height_feet", 1, 20},
		{"circle_thickness_feet", 1, 5},
		{"appear_dice_count", 3, 7},
		{"appear_die_sides", 3, 8},
		{"movement_cost_feet_per_foot", 1, 4},
		{"passage_dice_count", 1, 7},
		{"passage_die_sides", 1, 8},
		{"slot_dice_count", 2, 1},
		{"slot_die_sides", 2, 8},
		{"benefit_range_feet", 1, 60},
		{"armor_class_bonus", 1, 1},
		{"saving_throw_bonus", 1, 1},
		{"ring_count", 1, 2},
		{"ring_minimum_cost_cp", 1, 5000},
		{"separation_over_feet", 1, 60},
		{"max_creatures", 1, 10},
		{"max_creatures", 2, 10},
		{"cube_feet", 1, 20},
		{"flat_surface_depth_feet", 1, 5},
		{"flammable_cube_feet", 1, 5},
		{"burn_away_rounds", 1, 1},
		{"fire_dice_count", 1, 2},
		{"fire_die_sides", 1, 4},
		{"radius_feet", 1, 30},
		{"damage_dice_count", 1, 10},
		{"damage_die_sides", 1, 10},
		{"repeat_dice_count", 1, 5},
		{"repeat_die_sides", 1, 10},
		{"max_other_creatures", 1, 10},
		{"fly_speed_feet", 1, 300},
		{"revert_seconds", 1, 60},
		{"return_seconds", 1, 60},
		{"descent_feet_per_round", 1, 60},
		{"descent_seconds", 1, 60},
		{"max_length_feet", 3, 50},
		{"max_height_feet", 1, 15},
		{"thickness_feet", 3, 1},
		{"damage_dice_count", 2, 4},
		{"damage_die_sides", 2, 8},
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
// inventory assigns to catalogue issue #223.
func TestSpells223MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 223 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells223(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells223ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_223.json")
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
	want := "SRD 5.2.1 PDF pp.172-174 adapted into twelve complete offline spell families for issue #223, including page-break continuation text for Wall of Ice (pp.172-173) and Water Walk (pp.173-174). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells223PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.WindWall = loadSpells223(t).Records[11].SourceDetails.WindWall
	if r.Validate() == nil {
		t.Fatal("batch 223 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells223RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells223(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell223Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell223Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
