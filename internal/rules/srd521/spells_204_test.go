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

func loadSpells204(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_204.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.126-128 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells204SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/dominate-person","name":"Dominate Person","source_id":"srd-5.2.1:content/spells/dominate-person","source_page":126,"scope":"complete-source-record","level":5,"school":"enchantment","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-humanoid-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"slot-duration","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.dominate-person-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"save.dominate-person-repeat-on-damage","timing":"whenever-target-takes-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"link","mechanic":"command.dominate-person-telepathic-link","timing":"each-of-your-turns-no-action-required","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.dominate-person-charmed","save.dominate-person-repeat-on-damage","command.dominate-person-telepathic-link","save.wisdom-advantage-if-fighting","scaling.duration-by-slot","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","dominate":{"target":"humanoid-you-can-see-within-range","save":"wisdom","condition":"charmed","save_advantage":"you-or-allies-fighting-target","repeat_save":"whenever-target-takes-damage","repeat_success":"spell-ends-on-target","link_requirement":"same-plane-of-existence","command_action":"no-action-required","command_timing":"on-your-turn","command_examples":["attack-that-creature","move-over-there","fetch-that-object"],"obedience":"does-its-best-to-obey-on-its-turn","no_direction":"acts-and-moves-as-it-likes-focusing-on-protecting-itself","reaction_command":"can-command-target-to-take-reaction-using-your-reaction","slot_durations":[{"min_slot":6,"max_slot":6,"seconds":600},{"min_slot":7,"max_slot":7,"seconds":3600},{"min_slot":8,"max_slot":9,"seconds":28800}]}}}`,
		`{"id":"srd-5.2.1:content/spells/dragons-breath","name":"Dragon’s Breath","source_id":"srd-5.2.1:content/spells/dragons-breath","source_page":126,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"bonus-action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a hot pepper","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"choice","mechanic":"choice.dragons-breath-damage-type","timing":"on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.dragons-breath","timing":"each-exhalation-until-spell-ends","dice_count":3,"die_sides":6,"damage_type":"chosen-type","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["choice.dragons-breath-damage-type","damage.dragons-breath","area.cone-15","save.dexterity-half","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","breath":{"target":"willing-creature-touched","damage_types":["acid","cold","fire","lightning","poison"],"exhale_action":"magic-action","exhale_timing":"until-spell-ends","area_shape":"cone","area_feet":15,"area_targets":"each-creature-in-area","save":"dexterity","damage_dice_count":3,"damage_die_sides":6,"fail_damage":"full-damage-of-chosen-type","success_damage":"half-damage","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/dream","name":"Dream","source_id":"srd-5.2.1:content/spells/dream","source_page":126,"scope":"complete-source-record","level":5,"school":"illusion","lists":["bard","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"special","feet":0},"target":{"shape":"single","count":1,"eligibility":"known-creature-same-plane"},"components":{"verbal":true,"somatic":true,"material":{"description":"a handful of sand","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration_seconds":28800,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"messenger","mechanic":"state.dream-messenger-trance","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"converse","mechanic":"communication.dream-converse","timing":"while-target-asleep","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"terror","mechanic":"save.dream-terror-wisdom","timing":"when-terrifying-message-delivered","dice_count":3,"die_sides":6,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["state.dream-messenger-trance","communication.dream-converse","save.dream-terror-wisdom","range.special","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","dream":{"target":"creature-you-know-on-same-plane","messenger_choices":["you","willing-creature-you-touch"],"messenger_state":"trance","messenger_conditions":["incapacitated"],"messenger_speed":0,"asleep_effect":"messenger-appears-in-targets-dreams-and-converses-while-it-remains-asleep","dream_shaping":"messenger-shapes-dream-environment-landscapes-objects-and-other-images","emerge":"messenger-can-emerge-from-trance-at-any-time-ending-spell","recall":"target-recalls-dream-perfectly-upon-waking","awake_knowledge":"messenger-knows-target-is-awake","awake_options":["end-trance-and-spell","wait-for-target-to-sleep-then-enter-its-dreams"],"terrify_option":"you-can-make-messenger-terrifying","terrify_message_max_words":10,"terrify_save":"wisdom","terrify_failure":["no-benefit-from-rest","psychic-damage-when-it-wakes-up"],"terrify_damage_dice_count":3,"terrify_damage_die_sides":6,"terrify_damage_type":"psychic"}}}`,
		`{"id":"srd-5.2.1:content/spells/druidcraft","name":"Druidcraft","source_id":"srd-5.2.1:content/spells/druidcraft","source_page":126,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sensor","mechanic":"effect.druidcraft-weather-sensor","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"bloom","mechanic":"effect.druidcraft-bloom","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sensory","mechanic":"effect.druidcraft-sensory","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"fire","mechanic":"effect.druidcraft-fire-play","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["effect.druidcraft-weather-sensor","effect.druidcraft-bloom","effect.druidcraft-sensory","effect.druidcraft-fire-play"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","druidcraft":{"options":["weather-sensor","bloom","sensory-effect","fire-play"],"weather_sensor_size":"tiny","weather_sensor_predicts":"weather-at-your-location","weather_sensor_hours":24,"weather_sensor_persists_rounds":1,"weather_sensor_examples":["golden-orb-clear-skies","cloud-rain","falling-snowflakes-snow"],"bloom_targets":["flower-blossoms","seed-pod-opens","leaf-bud-blooms"],"bloom_timing":"instantly","sensory_examples":["falling-leaves","spectral-dancing-fairies","gentle-breeze","sound-of-an-animal","faint-odor-of-skunk"],"sensory_cube_feet":5,"fire_play_targets":["candle","torch","campfire"],"fire_play_actions":["light","snuff-out"]}}}`,
		`{"id":"srd-5.2.1:content/spells/earthquake","name":"Earthquake","source_id":"srd-5.2.1:content/spells/earthquake","source_page":127,"scope":"complete-source-record","level":8,"school":"transmutation","lists":["cleric","druid","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":500},"target":{"shape":"created-effect","count":1,"eligibility":"visible-ground-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a fractured rock","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"tremor","mechanic":"area.earthquake-difficult-terrain","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"shake","mechanic":"save.earthquake-dexterity-prone","timing":"when-cast-and-end-of-each-of-your-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"fissures","mechanic":"area.earthquake-fissures","timing":"end-of-turn-you-cast","dice_count":1,"die_sides":6,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"structures","mechanic":"damage.earthquake-structures","timing":"when-cast-and-end-of-each-of-your-turns","dice_count":0,"die_sides":0,"damage_type":"bludgeoning","bonus":50,"add_casting_modifier":false,"half_initial_only":false},{"id":"collapse","mechanic":"damage.earthquake-collapse","timing":"when-structure-collapses","dice_count":12,"die_sides":6,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.earthquake-difficult-terrain","save.earthquake-dexterity-prone","area.earthquake-fissures","damage.earthquake-structures","damage.earthquake-collapse","concentration.source-duration","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","earthquake":{"point":"point-on-ground-you-can-see-within-range","radius_feet":100,"area_shape":"circle","terrain":"difficult-terrain","save":"dexterity","save_timing":["when-you-cast-this-spell","end-of-each-of-your-turns-for-the-duration"],"save_targets":"each-creature-on-the-ground-in-area","failure_conditions":["prone"],"failure_breaks_concentration":true,"fissure_count_dice_count":1,"fissure_count_die_sides":6,"fissure_timing":"end-of-turn-you-cast","fissure_locations":"you-choose-locations-not-under-structures","fissure_depth_dice_count":1,"fissure_depth_die_sides":10,"fissure_depth_multiplier_feet":10,"fissure_width_feet":10,"fissure_extent":"from-one-edge-of-area-to-another-edge","fissure_fall_trigger":"creature-in-same-space-as-fissure","fissure_save":"dexterity","fissure_failure":"falls-in","fissure_success":"moves-with-fissures-edge-as-it-opens","structure_damage":50,"structure_damage_type":"bludgeoning","structure_contact":"structure-in-contact-with-ground-in-area","structure_timing":["when-you-cast-the-spell","end-of-each-of-your-turns-until-spell-ends"],"structure_collapse":"collapses-when-it-drops-to-0-hit-points","collapse_distance":"within-half-the-collapsing-structures-height","collapse_save":"dexterity","collapse_damage_dice_count":12,"collapse_damage_die_sides":6,"collapse_damage_type":"bludgeoning","collapse_failure_conditions":["prone","buried-in-rubble"],"collapse_escape_action":"action","collapse_escape_check":"strength-athletics","collapse_escape_dc":20,"collapse_success":"half-damage-only"}}}`,
		`{"id":"srd-5.2.1:content/spells/eldritch-blast","name":"Eldritch Blast","source_id":"srd-5.2.1:content/spells/eldritch-blast","source_page":127,"scope":"complete-source-record","level":0,"school":"evocation","lists":["warlock"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"cantrip-beams","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":["beam"]},"effects":[{"id":"beam","mechanic":"damage.eldritch-blast","timing":"per-beam-on-hit","dice_count":1,"die_sides":10,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["attack.ranged-spell","damage.eldritch-blast","scaling.cantrip-beams"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","blast":{"attack":"ranged-spell-attack","targets":"one-creature-or-object-in-range","damage_dice_count":1,"damage_die_sides":10,"damage_type":"force","base_beams":1,"beam_steps":[{"character_level":5,"beams":2},{"character_level":11,"beams":3},{"character_level":17,"beams":4}],"beam_targeting":"same-target-or-different-targets","attack_rolls":"separate-attack-roll-for-each-beam"}}}`,
		`{"id":"srd-5.2.1:content/spells/elementalism","name":"Elementalism","source_id":"srd-5.2.1:content/spells/elementalism","source_page":127,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration_seconds":0,"duration":"instantaneous","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"air","mechanic":"effect.elementalism-beckon-air","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"earth","mechanic":"effect.elementalism-beckon-earth","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"fire","mechanic":"effect.elementalism-beckon-fire","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"water","mechanic":"effect.elementalism-beckon-water","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sculpt","mechanic":"effect.elementalism-sculpt-element","timing":"chosen-option-on-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["effect.elementalism-beckon-air","effect.elementalism-beckon-earth","effect.elementalism-beckon-fire","effect.elementalism-beckon-water","effect.elementalism-sculpt-element"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","elementalism":{"options":["beckon-air","beckon-earth","beckon-fire","beckon-water","sculpt-element"],"beckon_air_cube_feet":5,"beckon_air_effects":["ripple-cloth","stir-dust","rustle-leaves","close-open-doors-and-shutters"],"beckon_air_unaffected":"doors-and-shutters-held-open-by-someone-or-something","beckon_earth_area_feet":5,"beckon_earth_shroud":"thin-shroud-of-dust-or-sand-covers-surfaces","beckon_earth_alternative":"single-word-in-your-handwriting-in-patch-of-dirt-or-sand","beckon_fire_cube_feet":5,"beckon_fire_cloud":"thin-cloud-of-harmless-embers-and-colored-scented-smoke","beckon_fire_choices":["color","scent"],"beckon_fire_embers":"can-light-candles-torches-or-lamps-in-area","beckon_fire_scent_lingers_seconds":60,"beckon_water_cube_feet":5,"beckon_water_mist":"spray-of-cool-mist-lightly-dampens-creatures-and-objects","beckon_water_alternative_cups":1,"beckon_water_alternative_places":["open-container","surface"],"beckon_water_evaporates_seconds":60,"sculpt_materials":["dirt","sand","fire","smoke","mist","water"],"sculpt_cube_feet":1,"sculpt_shape":"crude-shape-such-as-that-of-a-creature","sculpt_seconds":3600}}}`,
		`{"id":"srd-5.2.1:content/spells/enhance-ability","name":"Enhance Ability","source_id":"srd-5.2.1:content/spells/enhance-ability","source_page":127,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["bard","cleric","druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"fur or a feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":3600,"duration":"timed","scaling":{"kind":"slot-targets","above_level":2,"dice_count":0,"die_sides":0,"effect_ids":["advantage"]},"effects":[{"id":"advantage","mechanic":"check.enhance-ability-advantage","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["check.enhance-ability-advantage","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","enhance":{"target":"creature-touched","abilities":["strength","dexterity","intelligence","wisdom","charisma"],"benefit":"advantage-on-ability-checks-using-chosen-ability","base_targets":1,"additional_targets_per_slot":1,"ability_per_target":"can-choose-different-ability-for-each-target"}}}`,
		`{"id":"srd-5.2.1:content/spells/enlarge-reduce","name":"Enlarge/Reduce","source_id":"srd-5.2.1:content/spells/enlarge-reduce","source_page":127,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["bard","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of powdered iron","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"enlarge","mechanic":"size.enlarge","timing":"for-duration-if-chosen","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"reduce","mechanic":"size.reduce","timing":"for-duration-if-chosen","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"resist","mechanic":"save.enlarge-reduce-constitution-unwilling","timing":"when-target-is-unwilling-creature","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["size.enlarge","size.reduce","save.enlarge-reduce-constitution-unwilling","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","enlarge":{"source_pages":[127,128],"targets":"creature-or-object-you-can-see-within-range","object_condition":"neither-worn-nor-carried","unwilling_save":"constitution","unwilling_success":"spell-has-no-effect","carried_items":"everything-worn-and-carried-changes-size-with-creature","dropped_items":"return-to-normal-size-at-once","thrown_items":"thrown-weapon-or-ammunition-returns-to-normal-size-immediately-after-hit-or-miss","enlarge_size_change":"size-increases-one-category","enlarge_size_example":"medium-to-large","enlarge_advantage":["strength-checks","strength-saving-throws"],"enlarge_damage_dice_count":1,"enlarge_damage_die_sides":4,"enlarge_damage_sources":["enlarged-weapons","unarmed-strikes"],"reduce_size_change":"size-decreases-one-category","reduce_size_example":"medium-to-small","reduce_disadvantage":["strength-checks","strength-saving-throws"],"reduce_damage_dice_count":1,"reduce_damage_die_sides":4,"reduce_damage_sources":["reduced-weapons","unarmed-strikes"],"reduce_minimum_damage":1}}}`,
		`{"id":"srd-5.2.1:content/spells/ensnaring-strike","name":"Ensnaring Strike","source_id":"srd-5.2.1:content/spells/ensnaring-strike","source_page":128,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["ranger"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"hit-target"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"restrain","mechanic":"condition.ensnaring-strike-restrained","timing":"on-failed-strength-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.ensnaring-strike","timing":"start-of-each-restrained-turn","dice_count":1,"die_sides":6,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"escape","mechanic":"escape.ensnaring-strike-athletics","timing":"action-by-target-or-creature-within-reach","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.ensnaring-strike-restrained","damage.ensnaring-strike","escape.ensnaring-strike-athletics","cast.bonus-action-after-hit","scaling.per-slot","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["immediately-after-hitting-creature-with-weapon"],"end_timing":"duration-or-concentration-ended","ensnaring":{"trigger":"as-you-hit-the-target","save":"strength","large_advantage":"large-or-larger-creature-has-advantage-on-save","failure_condition":"restrained","failure_duration":"until-spell-ends","success_result":"vines-shrivel-away-and-spell-ends","damage_dice_count":1,"damage_die_sides":6,"damage_type":"piercing","damage_timing":"start-of-each-of-restrained-targets-turns","escape_actors":["target","creature-within-reach-of-target"],"escape_action":"action","escape_check":"strength-athletics","escape_dc":"caster-spell-save-dc","escape_success":"spell-ends","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/entangle","name":"Entangle","source_id":"srd-5.2.1:content/spells/entangle","source_page":128,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"square","count":0,"eligibility":"creatures-in-visible-ground-area"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"area.entangle-difficult-terrain","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"restrain","mechanic":"save.entangle-strength-restrained","timing":"when-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"escape","mechanic":"escape.entangle-athletics","timing":"action-while-restrained","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.entangle-difficult-terrain","save.entangle-strength-restrained","escape.entangle-athletics","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","entangle":{"area_shape":"square","area_feet":20,"area_origin":"ground-within-range","terrain":"difficult-terrain","plants_disappear":"when-spell-ends","targets":"each-creature-other-than-you-in-area-when-cast","save":"strength","failure_condition":"restrained","failure_duration":"until-spell-ends","escape_action":"action","escape_check":"strength-athletics","escape_dc":"caster-spell-save-dc","escape_success":"frees-itself-and-no-longer-restrained-by-the-plants"}}}`,
		`{"id":"srd-5.2.1:content/spells/enthrall","name":"Enthrall","source_id":"srd-5.2.1:content/spells/enthrall","source_page":128,"scope":"complete-source-record","level":2,"school":"enchantment","lists":["bard","warlock"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"selected-creatures","count":0,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration_seconds":60,"duration":"timed","scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"distract","mechanic":"save.enthrall-wisdom-perception-penalty","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["save.enthrall-wisdom-perception-penalty","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","enthrall":{"targets":"creatures-of-your-choice-you-can-see-within-range","save":"wisdom","automatic_success":"creatures-you-or-your-companions-are-fighting","failure_penalty":-10,"penalized_rolls":["wisdom-perception-checks","passive-perception"],"penalty_duration":"until-spell-ends"}}}`,
	}
	batch := loadSpells204(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Dominate Person", "Dragon’s Breath", "Dream", "Druidcraft", "Earthquake", "Eldritch Blast", "Elementalism", "Enhance Ability", "Enlarge/Reduce", "Ensnaring Strike", "Entangle", "Enthrall"}
	pages := []int{126, 126, 126, 126, 127, 127, 127, 127, 127, 128, 128, 128}
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

func TestSpells204CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells204(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 168 {
		t.Fatalf("expected 168 records, got %d", len(all.Records))
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
func TestSpells204MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_204.json")
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

func TestSpells204JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_204.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"area_feet": 15`, `"area_feet": 15, "area_feet": 16`},
		{`"area_feet": 15`, `"area_feet": 15, "AREA_FEET": 16`},
		{`"area_feet": 15`, `"area_feet": 15, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/dominate-person"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "charm"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 600`},
		{`"kind": "special"`, `"kind": "unlimited"`},
		{`"feet": 500`, `"feet": 501`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"failure_penalty": -10`, `"failure_penalty": -5`},
		{`"concentration": true`, `"concentration": false`},
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

func TestSpells204InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"dominate duration table", 0, func(s *SpellRecord) { s.SourceDetails.Dominate.SlotDurations[2].Seconds = 3600 }},
		{"dominate scaling kind", 0, func(s *SpellRecord) { s.Scaling.Kind = "slot-targets" }},
		{"dominate save", 0, func(s *SpellRecord) { s.SourceDetails.Dominate.Save = "charisma" }},
		{"dominate target", 0, func(s *SpellRecord) { s.Target.Eligibility = "visible-creature-in-range" }},
		{"dominate list", 0, func(s *SpellRecord) { s.Lists = append(s.Lists, "cleric") }},
		{"breath damage types", 1, func(s *SpellRecord) { s.SourceDetails.Breath.DamageTypes[4] = "necrotic" }},
		{"breath dice", 1, func(s *SpellRecord) { s.Effects[1].DiceCount = 2 }},
		{"breath scaling", 1, func(s *SpellRecord) { s.Scaling.AboveLevel = 3 }},
		{"breath material", 1, func(s *SpellRecord) { s.Components.Material.Description = "a pepper" }},
		{"dream range", 2, func(s *SpellRecord) { s.Range = SpellRange{"unlimited", 0} }},
		{"dream casting", 2, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"dream words", 2, func(s *SpellRecord) { s.SourceDetails.Dream.TerrifyMessageMaxWords = 20 }},
		{"dream duration", 2, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"druidcraft option", 3, func(s *SpellRecord) { s.SourceDetails.Druidcraft.Options[3] = "unknown" }},
		{"druidcraft hours", 3, func(s *SpellRecord) { s.SourceDetails.Druidcraft.WeatherSensorHours = 12 }},
		{"earthquake radius", 4, func(s *SpellRecord) { s.SourceDetails.Earthquake.RadiusFeet = 60 }},
		{"earthquake structure damage", 4, func(s *SpellRecord) { s.Effects[3].Bonus = 5 }},
		{"earthquake escape dc", 4, func(s *SpellRecord) { s.SourceDetails.Earthquake.CollapseEscapeDC = 15 }},
		{"earthquake concentration", 4, func(s *SpellRecord) { *s.Concentration = false }},
		{"blast beam step", 5, func(s *SpellRecord) { s.SourceDetails.Blast.BeamSteps[2].Beams = 3 }},
		{"blast step count", 5, func(s *SpellRecord) { s.SourceDetails.Blast.BeamSteps = s.SourceDetails.Blast.BeamSteps[:2] }},
		{"blast scaling kind", 5, func(s *SpellRecord) { s.Scaling.Kind = "cantrip-upgrade" }},
		{"blast cantrip slots", 5, func(s *SpellRecord) { s.SlotLevels = []int{1} }},
		{"elementalism option", 6, func(s *SpellRecord) { s.SourceDetails.Elementalism.Options[4] = "sculpt-metal" }},
		{"elementalism sculpt time", 6, func(s *SpellRecord) { s.SourceDetails.Elementalism.SculptSeconds = 60 }},
		{"enhance ability", 7, func(s *SpellRecord) { s.SourceDetails.Enhance.Abilities[2] = "constitution" }},
		{"enhance targets", 7, func(s *SpellRecord) { s.SourceDetails.Enhance.AdditionalTargetsPerSlot = 2 }},
		{"enhance scaling level", 7, func(s *SpellRecord) { s.Scaling.AboveLevel = 3 }},
		{"enlarge continuation page", 8, func(s *SpellRecord) { s.SourceDetails.Enlarge.SourcePages = []int{127} }},
		{"enlarge die", 8, func(s *SpellRecord) { s.SourceDetails.Enlarge.EnlargeDamageDieSides = 6 }},
		{"enlarge minimum damage", 8, func(s *SpellRecord) { s.SourceDetails.Enlarge.ReduceMinimumDamage = 0 }},
		{"ensnaring trigger", 9, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{} }},
		{"ensnaring components", 9, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"ensnaring slot dice", 9, func(s *SpellRecord) { s.SourceDetails.Ensnaring.SlotDieSides = 8 }},
		{"entangle area", 10, func(s *SpellRecord) { s.SourceDetails.Entangle.AreaFeet = 30 }},
		{"entangle range", 10, func(s *SpellRecord) { s.Range.Feet = 60 }},
		{"enthrall penalty", 11, func(s *SpellRecord) { s.SourceDetails.Enthrall.FailurePenalty = 10 }},
		{"enthrall rolls", 11, func(s *SpellRecord) {
			s.SourceDetails.Enthrall.PenalizedRolls = s.SourceDetails.Enthrall.PenalizedRolls[:1]
		}},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Entangle = &SpellEntangle{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells204(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells204RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_204.json")
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

func TestSpells204WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_204.json")
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
		`"seconds": 600, "seconds": 3600`,
		`"seconds": 600, "SECONDS": 3600`,
		`"seconds": 600, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"seconds": 600`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells204PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_204.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"area_feet", 1, 15}, {"weather_sensor_hours", 1, 24}, {"weather_sensor_persists_rounds", 1, 1}, {"sensory_cube_feet", 1, 5},
		{"radius_feet", 1, 100}, {"fissure_width_feet", 1, 10}, {"structure_damage", 1, 50}, {"collapse_damage_dice_count", 1, 12},
		{"collapse_escape_dc", 1, 20}, {"damage_dice_count", 1, 3}, {"damage_die_sides", 1, 6}, {"terrify_message_max_words", 1, 10},
		{"messenger_speed", 1, 0}, {"base_beams", 1, 1}, {"beams", 1, 2}, {"beams", 2, 3}, {"beams", 3, 4},
		{"character_level", 1, 5}, {"character_level", 2, 11}, {"character_level", 3, 17},
		{"beckon_air_cube_feet", 1, 5}, {"beckon_fire_scent_lingers_seconds", 1, 60}, {"beckon_water_evaporates_seconds", 1, 60},
		{"sculpt_cube_feet", 1, 1}, {"sculpt_seconds", 1, 3600}, {"base_targets", 1, 1}, {"additional_targets_per_slot", 1, 1},
		{"reduce_minimum_damage", 1, 1}, {"min_slot", 1, 6}, {"min_slot", 2, 7}, {"min_slot", 3, 8}, {"max_slot", 3, 9},
		{"seconds", 1, 600}, {"seconds", 2, 3600}, {"seconds", 3, 28800},
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
// inventory assigns to catalogue issue #204.
func TestSpells204MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 204 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells204(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells204ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_204.json")
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
	want := "SRD 5.2.1 PDF pp.126-128 adapted into twelve complete offline spell families for issue #204, including page-break continuation text for Enlarge/Reduce (pp.127-128). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	// Entangle starts and ends on p.128, so it is not a page-break continuation.
	if strings.Contains(root.Source.Changes, "Entangle") {
		t.Fatal("Entangle sits wholly on p.128 and is not a page-break continuation")
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells204PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.Enthrall = loadSpells204(t).Records[11].SourceDetails.Enthrall
	if r.Validate() == nil {
		t.Fatal("batch 204 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells204RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells204(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell204Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell204Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
