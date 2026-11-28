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

func loadSpells221(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_221.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.166-169 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells221SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/summon-dragon","name":"Summon Dragon","source_id":"srd-5.2.1:content/spells/summon-dragon","source_page":166,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"an object with the image of a dragon engraved on it worth 500+ GP","cost_cp":50000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-stat-block","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["spirit"]},"effects":[{"id":"spirit","mechanic":"summon.summon-dragon","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.summon-dragon","stat-block.draconic-spirit","scaling.stat-block-by-slot","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-spirit-at-0-hit-points","summon_dragon":{"summons":"dragon-spirit","placement":"unoccupied-space-caster-can-see-within-range","disappears_when":["drops-to-0-hit-points","spell-ends"],"combat_role":"ally-to-caster-and-casters-allies","initiative":"shares-casters-initiative-count","turn_order":"immediately-after-casters-turn","commands":"obeys-verbal-commands-no-action-required","no_command":"takes-dodge-action-and-moves-to-avoid-danger","slot_level_use":"spell-slot-level-is-the-spell-level-in-the-stat-block","stat_block":{"name":"draconic-spirit","size":"large","creature_type":"dragon","alignment":"neutral","armor_class_base":14,"armor_class_adds":"spells-level","hit_points_base":50,"hit_points_per_level":10,"hit_points_above_level":5,"speed_feet":30,"fly_speed_feet":60,"swim_speed_feet":30,"abilities":[{"ability":"strength","score":19,"modifier":4,"save":4},{"ability":"dexterity","score":14,"modifier":2,"save":2},{"ability":"constitution","score":17,"modifier":3,"save":3},{"ability":"intelligence","score":10,"modifier":0,"save":0},{"ability":"wisdom","score":14,"modifier":2,"save":2},{"ability":"charisma","score":14,"modifier":2,"save":2}],"damage_resistances":["acid","cold","fire","lightning","poison"],"condition_immunities":["charmed","frightened","poisoned"],"blindsight_feet":30,"darkvision_feet":60,"passive_perception":12,"languages":["draconic","understands-languages-caster-knows"],"challenge_rating":"none","xp":0,"proficiency_bonus":"equals-casters-proficiency-bonus","shared_resistance":"caster-gains-resistance-to-one-chosen-spirit-resistance-until-spell-ends","multiattack_rend":"attacks-equal-to-half-spell-level-rounded-down","multiattack_also":"uses-breath-weapon","rend_attack_bonus":"casters-spell-attack-modifier","rend_reach_feet":10,"rend_dice_count":1,"rend_die_sides":6,"rend_bonus":4,"rend_bonus_adds":"spells-level","rend_damage_type":"piercing","breath_save":"dexterity","breath_save_dc":"casters-spell-save-dc","breath_cone_feet":30,"breath_dice_count":2,"breath_die_sides":6,"breath_damage_type":"spirit-resistance-type-chosen-at-casting","breath_success":"half-damage"}}}}`,
		`{"id":"srd-5.2.1:content/spells/sunbeam","name":"Sunbeam","source_id":"srd-5.2.1:content/spells/sunbeam","source_page":166,"scope":"complete-source-record","level":6,"school":"evocation","lists":["cleric","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"line","count":0,"eligibility":"creatures-in-self-origin-line"},"components":{"verbal":true,"somatic":true,"material":{"description":"a magnifying glass","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.sunbeam","timing":"at-casting-and-on-later-magic-actions","dice_count":6,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blinded","mechanic":"condition.sunbeam-blinded","timing":"on-failed-save-until-start-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sunlight","mechanic":"light.sunbeam-sunlight","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.sunbeam","condition.sunbeam-blinded","light.sunbeam-sunlight","save.constitution-half","action.sunbeam-new-line","concentration.source-duration","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","sunbeam":{"line_width_feet":5,"line_length_feet":60,"save":"constitution","damage_dice_count":6,"damage_die_sides":8,"damage_type":"radiant","failure_condition":"blinded","condition_ends":"start-of-casters-next-turn","save_success":"half-damage-only","new_line_action":"magic-action","mote_location":"above-caster","bright_light_feet":30,"dim_light_extra_feet":30,"light_kind":"sunlight"}}}`,
		`{"id":"srd-5.2.1:content/spells/sunburst","name":"Sunburst","source_id":"srd-5.2.1:content/spells/sunburst","source_page":167,"scope":"complete-source-record","level":8,"school":"evocation","lists":["cleric","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of sunstone","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.sunburst","timing":"immediate","dice_count":12,"die_sides":6,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blinded","mechanic":"condition.sunburst-blinded","timing":"on-failed-save-up-to-1-minute","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dispel","mechanic":"ending.sunburst-dispels-darkness","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.sunburst","condition.sunburst-blinded","ending.sunburst-dispels-darkness","save.constitution-half","save.sunburst-repeat","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","sunburst":{"radius_feet":60,"save":"constitution","damage_dice_count":12,"damage_die_sides":6,"damage_type":"radiant","failure_condition":"blinded","condition_seconds":60,"save_success":"half-damage-only","repeat_save":"constitution","repeat_timing":"end-of-each-of-its-turns","repeat_success":"ends-effect-on-self","dispels":"darkness-in-area-created-by-any-spell"}}}`,
		`{"id":"srd-5.2.1:content/spells/symbol","name":"Symbol","source_id":"srd-5.2.1:content/spells/symbol","source_page":167,"scope":"complete-source-record","level":7,"school":"abjuration","lists":["bard","cleric","druid","wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"surface-or-closable-object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"powdered diamond worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"glyph","mechanic":"trap.symbol-glyph","timing":"until-triggered","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"death","mechanic":"damage.symbol-death","timing":"on-trigger-death","dice_count":10,"die_sides":10,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"discord","mechanic":"effect.symbol-discord","timing":"on-trigger-discord","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"fear","mechanic":"condition.symbol-frightened","timing":"on-trigger-fear","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"pain","mechanic":"condition.symbol-incapacitated","timing":"on-trigger-pain","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sleep","mechanic":"condition.symbol-unconscious","timing":"on-trigger-sleep","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"stunning","mechanic":"condition.symbol-stunned","timing":"on-trigger-stunning","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["trap.symbol-glyph","damage.symbol-death","effect.symbol-discord","condition.symbol-frightened","condition.symbol-incapacitated","condition.symbol-unconscious","condition.symbol-stunned","trigger.symbol","check.symbol-perception","cast.one-minute","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"dispelled-triggered-or-glyph-broken","symbol":{"placements":["surface-such-as-floor-or-wall","closable-object-such-as-book-or-chest"],"max_glyph_diameter_feet":10,"max_object_move_feet":10,"moved_result":"glyph-broken-and-spell-ends-untriggered","notice_check":"wisdom-perception","notice_dc":"caster-spell-save-dc","cast_choices":["trigger","effect"],"surface_triggers":["touching-or-stepping-on-glyph","removing-covering-object","approaching-within-set-distance"],"object_triggers":["opening-object","seeing-glyph"],"trigger_refinements":["only-certain-creature-types","exempt-conditions-such-as-password"],"refinement_example":"affect-aberrations","activation_radius_feet":60,"activation_light":"dim-light","active_minutes":10,"active_end":"spell-ends","targeted":["creatures-in-sphere-at-activation","creature-entering-first-time-on-turn","creature-ending-turn-there"],"targets_per_turn":1,"effects":[{"effect":"death","save":"constitution","result":"damage-half-on-success","dice_count":10,"die_sides":10,"damage_type":"necrotic","condition_seconds":0},{"effect":"discord","save":"wisdom","result":"argues-without-meaningful-communication-and-disadvantage-on-attacks-and-ability-checks","dice_count":0,"die_sides":0,"damage_type":"none","condition_seconds":60},{"effect":"fear","save":"wisdom","result":"frightened-must-move-away-from-glyph-each-turn-if-able","dice_count":0,"die_sides":0,"damage_type":"none","condition_seconds":60},{"effect":"pain","save":"constitution","result":"incapacitated","dice_count":0,"die_sides":0,"damage_type":"none","condition_seconds":60},{"effect":"sleep","save":"wisdom","result":"unconscious-awakens-on-damage-or-action-to-shake","dice_count":0,"die_sides":0,"damage_type":"none","condition_seconds":600},{"effect":"stunning","save":"wisdom","result":"stunned","dice_count":0,"die_sides":0,"damage_type":"none","condition_seconds":60}],"fear_min_move_feet":30}}}`,
		`{"id":"srd-5.2.1:content/spells/telekinesis","name":"Telekinesis","source_id":"srd-5.2.1:content/spells/telekinesis","source_page":167,"scope":"complete-source-record","level":5,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"visible-creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"creature","mechanic":"control.telekinesis-creature","timing":"at-casting-and-on-later-magic-actions","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"object","mechanic":"control.telekinesis-object","timing":"at-casting-and-on-later-magic-actions","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["control.telekinesis-creature","control.telekinesis-object","save.strength-telekinesis","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","telekinesis":{"source_pages":[167,168],"use_timing":"when-cast-and-as-magic-action-on-later-turns","target_kinds":["creature","object"],"retargeting":"same-or-new-target-any-time-prior-target-no-longer-affected","max_creature_size":"huge","creature_save":"strength","creature_move_feet":30,"move_limit":"any-direction-within-spell-range","creature_condition":"restrained","condition_ends":"end-of-casters-next-turn","lifted":"suspended-if-lifted","fall":"falls-at-end-of-casters-next-turn-unless-used-again-and-save-fails","max_object_size":"huge","unattended_object":"moved-automatically","object_move_feet":30,"carried_object_save":"strength","carried_object_failure":"object-pulled-away-and-moved","fine_control":"can-exert-fine-control-on-objects","fine_control_examples":["manipulating-simple-tool"],"source_text_end":"printed-sentence-ends-after-first-fine-control-example"}}}`,
		`{"id":"srd-5.2.1:content/spells/telepathic-bond","name":"Telepathic Bond","source_id":"srd-5.2.1:content/spells/telepathic-bond","source_page":168,"scope":"complete-source-record","level":5,"school":"divination","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":8,"eligibility":"chosen-willing-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"two eggs","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"bond","mechanic":"utility.telepathic-bond","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.telepathic-bond","cast.ritual","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","telepathic_bond":{"max_creatures":8,"link":"each-creature-linked-to-all-others","excluded":"creatures-that-cannot-communicate-in-any-language","communication":"telepathic-with-or-without-shared-language","distance":"any-distance","planar_limit":"does-not-extend-to-other-planes"}}}`,
		`{"id":"srd-5.2.1:content/spells/teleport","name":"Teleport","source_id":"srd-5.2.1:content/spells/teleport","source_page":168,"scope":"complete-source-record","level":7,"school":"conjuration","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":10},"target":{"shape":"choice","count":1,"eligibility":"willing-creatures-or-single-object-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"transport","mechanic":"teleport.teleport","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"mishap","mechanic":"damage.teleport-mishap","timing":"on-each-mishap-result","dice_count":3,"die_sides":10,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.teleport","damage.teleport-mishap","table.teleport-outcome"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","teleport":{"max_additional_creatures":8,"creature_requirement":"willing-creatures-caster-can-see","object_alternative":"single-visible-object-in-range","max_object_size":"large","object_restriction":"not-held-or-carried-by-unwilling-creature","destination_requirements":["known-to-caster","same-plane-of-existence"],"table_die_sides":100,"table":[{"familiarity":"permanent-circle","outcomes":[{"outcome":"on-target","low":1,"high":100}]},{"familiarity":"linked-object","outcomes":[{"outcome":"on-target","low":1,"high":100}]},{"familiarity":"very-familiar","outcomes":[{"outcome":"mishap","low":1,"high":5},{"outcome":"similar-area","low":6,"high":13},{"outcome":"off-target","low":14,"high":24},{"outcome":"on-target","low":25,"high":100}]},{"familiarity":"seen-casually","outcomes":[{"outcome":"mishap","low":1,"high":33},{"outcome":"similar-area","low":34,"high":43},{"outcome":"off-target","low":44,"high":53},{"outcome":"on-target","low":54,"high":100}]},{"familiarity":"viewed-once-or-described","outcomes":[{"outcome":"mishap","low":1,"high":43},{"outcome":"similar-area","low":44,"high":53},{"outcome":"off-target","low":54,"high":73},{"outcome":"on-target","low":74,"high":100}]},{"familiarity":"false-destination","outcomes":[{"outcome":"mishap","low":1,"high":50},{"outcome":"similar-area","low":51,"high":100}]}],"familiarity_meanings":[{"familiarity":"permanent-circle","meaning":"permanent-teleportation-circle-with-known-sigil-sequence"},{"familiarity":"linked-object","meaning":"possess-object-taken-from-destination-recently"},{"familiarity":"very-familiar","meaning":"visited-often-carefully-studied-or-visible-at-casting"},{"familiarity":"seen-casually","meaning":"seen-more-than-once-but-not-very-familiar"},{"familiarity":"viewed-once-or-described","meaning":"seen-once-possibly-by-magic-or-known-from-description-such-as-map"},{"familiarity":"false-destination","meaning":"place-that-does-not-exist-or-no-longer-exists"}],"linked_object_months":6,"linked_object_example":"book-from-wizards-library","false_destination_examples":["illusion-seen-while-scrying-enemy-sanctum","location-that-no-longer-exists"],"mishap_dice_count":3,"mishap_die_sides":10,"mishap_damage_type":"force","mishap_targets":"each-teleporting-creature-or-target-object","mishap_reroll":"gm-rerolls-and-multiple-mishaps-damage-each-time","similar_area":"closest-visually-or-thematically-similar-place","similar_area_example":"another-persons-laboratory-in-same-city","off_target_distance_dice_count":2,"off_target_distance_die_sides":12,"off_target_distance_unit":"miles","off_target_direction_die_sides":8,"off_target_directions":["east","southeast","south","southwest","west","northwest","north","northeast"],"on_target":"appear-where-intended"}}}`,
		`{"id":"srd-5.2.1:content/spells/teleportation-circle","name":"Teleportation Circle","source_id":"srd-5.2.1:content/spells/teleportation-circle","source_page":168,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":10},"target":{"shape":"created-effect","count":1,"eligibility":"circle-on-ground-in-range"},"components":{"verbal":true,"somatic":false,"material":{"description":"rare inks worth 50+ GP, which the spell consumes","cost_cp":5000,"consumed":true}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":6,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"portal","mechanic":"teleport.teleportation-circle","timing":"until-end-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.teleportation-circle","knowledge.teleportation-circle-sigils","cast.one-minute","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","teleportation_circle":{"source_pages":[168,169],"circle_radius_feet":5,"destination":"known-permanent-circle-on-same-material-plane-of-existence","portal_closes":"end-of-casters-next-turn","arrival_distance_feet":5,"occupied_arrival":"nearest-unoccupied-space","common_locations":["major-temples","guildhalls","other-important-places"],"sigil_sequence":"unique-string-of-runes-in-particular-pattern","initial_destinations":2,"initial_destination_plane":"material-plane","initial_destination_chooser":"gm","learning":"may-learn-more-during-adventures","memorize_minutes":1,"permanent_casting_days":365,"permanent_requirement":"cast-in-same-location-every-day"}}}`,
		`{"id":"srd-5.2.1:content/spells/thaumaturgy","name":"Thaumaturgy","source_id":"srd-5.2.1:content/spells/thaumaturgy","source_page":169,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"choice","count":1,"eligibility":"one-chosen-minor-wonder-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"wonder","mechanic":"utility.thaumaturgy","timing":"chosen-option","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.thaumaturgy","limit.thaumaturgy-active-effects"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"up-to-duration","thaumaturgy":{"options":[{"choice":"altered-eyes","effect":"alter-appearance-of-casters-eyes","lasts":"1-minute"},{"choice":"booming-voice","effect":"voice-booms-louder-with-advantage-on-charisma-intimidation","lasts":"1-minute"},{"choice":"fire-play","effect":"flames-flicker-brighten-dim-or-change-color","lasts":"1-minute"},{"choice":"invisible-hand","effect":"unlocked-door-or-window-flies-open-or-slams-shut","lasts":"instantaneous"},{"choice":"phantom-sound","effect":"sound-from-chosen-point-such-as-thunder-raven-cry-or-whispers","lasts":"instantaneous"},{"choice":"tremors","effect":"harmless-tremors-in-ground","lasts":"1-minute"}],"choose_count":1,"voice_volume_multiplier":3,"max_active_one_minute":3}}}`,
		`{"id":"srd-5.2.1:content/spells/thunderwave","name":"Thunderwave","source_id":"srd-5.2.1:content/spells/thunderwave","source_page":169,"scope":"complete-source-record","level":1,"school":"evocation","lists":["bard","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"cube","count":0,"eligibility":"creatures-and-unsecured-objects-in-self-origin-cube"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.thunderwave","timing":"immediate","dice_count":2,"die_sides":8,"damage_type":"thunder","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"push","mechanic":"movement.thunderwave-push","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"boom","mechanic":"sound.thunderwave-boom","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.thunderwave","movement.thunderwave-push","sound.thunderwave-boom","save.constitution-half","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","thunderwave":{"cube_feet":15,"save":"constitution","damage_dice_count":2,"damage_die_sides":8,"damage_type":"thunder","creature_push_feet":10,"save_success":"half-damage-only","object_requirement":"unsecured-and-entirely-within-cube","object_push_feet":10,"audible_feet":300,"slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/time-stop","name":"Time Stop","source_id":"srd-5.2.1:content/spells/time-stop","source_page":169,"scope":"complete-source-record","level":9,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"turns","mechanic":"time.time-stop","timing":"immediate","dice_count":1,"die_sides":4,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["time.time-stop","ending.time-stop-affects-others-or-leaves-area"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"turns-taken-or-early-end","time_stop":{"turn_dice_count":1,"turn_die_sides":4,"turn_bonus":1,"others":"no-time-passes","turn_use":"actions-and-movement-as-normal","ending_triggers":["action-or-created-effect-affects-another-creature","action-or-created-effect-affects-object-worn-or-carried-by-another","caster-moves-too-far-from-casting-location"],"max_distance_feet":1000}}}`,
		`{"id":"srd-5.2.1:content/spells/tiny-hut","name":"Tiny Hut","source_id":"srd-5.2.1:content/spells/tiny-hut","source_page":169,"scope":"complete-source-record","level":3,"school":"evocation","lists":["bard","wizard"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"stationary-emanation-around-caster"},"components":{"verbal":true,"somatic":true,"material":{"description":"a crystal bead","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"hut","mechanic":"barrier.tiny-hut","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"interior","mechanic":"environment.tiny-hut","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.tiny-hut","environment.tiny-hut","ending.tiny-hut-leave-or-recast","cast.one-minute","cast.ritual","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-caster-leaves-or-cast-again","tiny_hut":{"emanation_feet":10,"mobility":"stationary","casting_failure":"fails-if-emanation-cannot-fully-encapsulate-all-creatures-in-area","free_passage":"creatures-and-objects-inside-at-casting","barred":"all-other-creatures-and-objects","max_blocked_spell_level":3,"blocked_spells":"cannot-be-cast-through-and-effects-cannot-extend-into-it","atmosphere":"comfortable-and-dry-regardless-of-weather","interior_light_choices":["dim-light","darkness"],"interior_light_action":"no-action-required","exterior":"opaque-in-chosen-color","interior_view":"transparent-from-inside","early_end":["caster-leaves-emanation","spell-cast-again"]}}}`,
	}
	batch := loadSpells221(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Summon Dragon", "Sunbeam", "Sunburst", "Symbol", "Telekinesis", "Telepathic Bond", "Teleport", "Teleportation Circle", "Thaumaturgy", "Thunderwave", "Time Stop", "Tiny Hut"}
	pages := []int{166, 166, 167, 167, 167, 168, 168, 168, 169, 169, 169, 169}
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

func TestSpells221CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells221(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t), loadSpells220(t), loadSpells221(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 312 {
		t.Fatalf("expected 312 records, got %d", len(all.Records))
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
func TestSpells221MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_221.json")
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

func TestSpells221JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_221.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"line_width_feet": 5`, `"line_width_feet": 5, "line_width_feet": 6`},
		{`"line_width_feet": 5`, `"line_width_feet": 5, "LINE_WIDTH_FEET": 6`},
		{`"line_width_feet": 5`, `"line_width_feet": 5, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/summon-dragon"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "spirit"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 60`, `"feet": 61`},
		{`"cost_cp": 50000`, `"cost_cp": 50001`},
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

func TestSpells221InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"dragon armor class", 0, func(s *SpellRecord) { s.SourceDetails.SummonDragon.StatBlock.ArmorClassBase = 12 }},
		{"dragon ability save", 0, func(s *SpellRecord) { s.SourceDetails.SummonDragon.StatBlock.Abilities[3].Save = -1 }},
		{"sunbeam length", 1, func(s *SpellRecord) { s.SourceDetails.Sunbeam.LineLengthFeet = 30 }},
		{"sunbeam dice", 1, func(s *SpellRecord) { s.Effects[0].DiceCount = 8 }},
		{"sunburst radius", 2, func(s *SpellRecord) { s.SourceDetails.Sunburst.RadiusFeet = 30 }},
		{"sunburst repeat", 2, func(s *SpellRecord) { s.SourceDetails.Sunburst.RepeatSave = "wisdom" }},
		{"symbol effects", 3, func(s *SpellRecord) { s.SourceDetails.Symbol.Effects = s.SourceDetails.Symbol.Effects[:5] }},
		{"symbol death dice", 3, func(s *SpellRecord) { s.SourceDetails.Symbol.Effects[0].DiceCount = 8 }},
		{"telekinesis size", 4, func(s *SpellRecord) { s.SourceDetails.Telekinesis.MaxCreatureSize = "gargantuan" }},
		{"telekinesis pages", 4, func(s *SpellRecord) { s.SourceDetails.Telekinesis.SourcePages = []int{167} }},
		{"bond creatures", 5, func(s *SpellRecord) { s.SourceDetails.TelepathicBond.MaxCreatures = 9 }},
		{"bond ritual", 5, func(s *SpellRecord) { *s.Ritual = false }},
		{"teleport table", 6, func(s *SpellRecord) { s.SourceDetails.Teleport.Table[2].Outcomes[0].High = 10 }},
		{"teleport somatic", 6, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"circle radius", 7, func(s *SpellRecord) { s.SourceDetails.TeleportationCircle.CircleRadiusFeet = 10 }},
		{"circle duration", 7, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"thaumaturgy active", 8, func(s *SpellRecord) { s.SourceDetails.Thaumaturgy.MaxActiveOneMinute = 2 }},
		{"thaumaturgy options", 8, func(s *SpellRecord) { s.SourceDetails.Thaumaturgy.Options = s.SourceDetails.Thaumaturgy.Options[:5] }},
		{"thunderwave push", 9, func(s *SpellRecord) { s.SourceDetails.Thunderwave.CreaturePushFeet = 15 }},
		{"thunderwave slot die", 9, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"time stop bonus", 10, func(s *SpellRecord) { s.Effects[0].Bonus = 2 }},
		{"time stop distance", 10, func(s *SpellRecord) { s.SourceDetails.TimeStop.MaxDistanceFeet = 500 }},
		{"hut level", 11, func(s *SpellRecord) { s.SourceDetails.TinyHut.MaxBlockedSpellLevel = 4 }},
		{"hut casting", 11, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.TinyHut = &SpellTinyHut{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells221(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells221RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_221.json")
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

func TestSpells221WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_221.json")
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
		`"line_width_feet": 5, "line_width_feet": 6`,
		`"line_width_feet": 5, "LINE_WIDTH_FEET": 6`,
		`"line_width_feet": 5, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"line_width_feet": 5`), []byte(replacement), 1)
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
func TestSpells221PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_221.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"armor_class_base", 1, 14},
		{"hit_points_base", 1, 50},
		{"hit_points_per_level", 1, 10},
		{"hit_points_above_level", 1, 5},
		{"speed_feet", 1, 30},
		{"fly_speed_feet", 1, 60},
		{"swim_speed_feet", 1, 30},
		{"score", 1, 19},
		{"modifier", 1, 4},
		{"save", 1, 4},
		{"score", 2, 14},
		{"modifier", 2, 2},
		{"save", 2, 2},
		{"score", 3, 17},
		{"modifier", 3, 3},
		{"save", 3, 3},
		{"score", 4, 10},
		{"modifier", 4, 0},
		{"save", 4, 0},
		{"score", 5, 14},
		{"modifier", 5, 2},
		{"save", 5, 2},
		{"score", 6, 14},
		{"modifier", 6, 2},
		{"save", 6, 2},
		{"blindsight_feet", 1, 30},
		{"darkvision_feet", 1, 60},
		{"passive_perception", 1, 12},
		{"xp", 1, 0},
		{"rend_reach_feet", 1, 10},
		{"rend_dice_count", 1, 1},
		{"rend_die_sides", 1, 6},
		{"rend_bonus", 1, 4},
		{"breath_cone_feet", 1, 30},
		{"breath_dice_count", 1, 2},
		{"breath_die_sides", 1, 6},
		{"line_width_feet", 1, 5},
		{"line_length_feet", 1, 60},
		{"damage_dice_count", 1, 6},
		{"damage_die_sides", 1, 8},
		{"bright_light_feet", 1, 30},
		{"dim_light_extra_feet", 1, 30},
		{"radius_feet", 1, 60},
		{"damage_dice_count", 2, 12},
		{"damage_die_sides", 2, 6},
		{"condition_seconds", 1, 60},
		{"max_glyph_diameter_feet", 1, 10},
		{"max_object_move_feet", 1, 10},
		{"activation_radius_feet", 1, 60},
		{"active_minutes", 1, 10},
		{"targets_per_turn", 1, 1},
		{"dice_count", 19, 10},
		{"die_sides", 19, 10},
		{"condition_seconds", 2, 0},
		{"dice_count", 20, 0},
		{"die_sides", 20, 0},
		{"condition_seconds", 3, 60},
		{"dice_count", 21, 0},
		{"die_sides", 21, 0},
		{"condition_seconds", 4, 60},
		{"dice_count", 22, 0},
		{"die_sides", 22, 0},
		{"condition_seconds", 5, 60},
		{"dice_count", 23, 0},
		{"die_sides", 23, 0},
		{"condition_seconds", 6, 600},
		{"dice_count", 24, 0},
		{"die_sides", 24, 0},
		{"condition_seconds", 7, 60},
		{"fear_min_move_feet", 1, 30},
		{"creature_move_feet", 1, 30},
		{"object_move_feet", 1, 30},
		{"max_creatures", 1, 8},
		{"max_additional_creatures", 1, 8},
		{"table_die_sides", 1, 100},
		{"low", 1, 1},
		{"high", 1, 100},
		{"low", 2, 1},
		{"high", 2, 100},
		{"low", 3, 1},
		{"high", 3, 5},
		{"low", 4, 6},
		{"high", 4, 13},
		{"low", 5, 14},
		{"high", 5, 24},
		{"low", 6, 25},
		{"high", 6, 100},
		{"low", 7, 1},
		{"high", 7, 33},
		{"low", 8, 34},
		{"high", 8, 43},
		{"low", 9, 44},
		{"high", 9, 53},
		{"low", 10, 54},
		{"high", 10, 100},
		{"low", 11, 1},
		{"high", 11, 43},
		{"low", 12, 44},
		{"high", 12, 53},
		{"low", 13, 54},
		{"high", 13, 73},
		{"low", 14, 74},
		{"high", 14, 100},
		{"low", 15, 1},
		{"high", 15, 50},
		{"low", 16, 51},
		{"high", 16, 100},
		{"linked_object_months", 1, 6},
		{"mishap_dice_count", 1, 3},
		{"mishap_die_sides", 1, 10},
		{"off_target_distance_dice_count", 1, 2},
		{"off_target_distance_die_sides", 1, 12},
		{"off_target_direction_die_sides", 1, 8},
		{"circle_radius_feet", 1, 5},
		{"arrival_distance_feet", 1, 5},
		{"initial_destinations", 1, 2},
		{"memorize_minutes", 1, 1},
		{"permanent_casting_days", 1, 365},
		{"choose_count", 1, 1},
		{"voice_volume_multiplier", 1, 3},
		{"max_active_one_minute", 1, 3},
		{"cube_feet", 1, 15},
		{"damage_dice_count", 3, 2},
		{"damage_die_sides", 3, 8},
		{"creature_push_feet", 1, 10},
		{"object_push_feet", 1, 10},
		{"audible_feet", 1, 300},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
		{"turn_dice_count", 1, 1},
		{"turn_die_sides", 1, 4},
		{"turn_bonus", 1, 1},
		{"max_distance_feet", 1, 1000},
		{"emanation_feet", 1, 10},
		{"max_blocked_spell_level", 1, 3},
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
// inventory assigns to catalogue issue #221.
func TestSpells221MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 221 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells221(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells221ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_221.json")
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
	want := "SRD 5.2.1 PDF pp.166-169 adapted into twelve complete offline spell families for issue #221, including page-break continuation text for Telekinesis (pp.167-168) and Teleportation Circle (pp.168-169). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells221PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.TinyHut = loadSpells221(t).Records[11].SourceDetails.TinyHut
	if r.Validate() == nil {
		t.Fatal("batch 221 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells221RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells221(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell221Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell221Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
