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

func loadSpells207(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_207.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.133-136 value fixtures include every choice, effect, stat block
// row and scaling value, independently of decoder defaults and structural validation.
func TestSpells207SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/fog-cloud","name":"Fog Cloud","source_id":"srd-5.2.1:content/spells/fog-cloud","source_page":133,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-radius","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["fog"]},"effects":[{"id":"fog","mechanic":"area.heavily-obscured-fog-sphere-20","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"disperse","mechanic":"dispel.strong-wind-disperses-fog","timing":"when-strong-wind-disperses","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.heavily-obscured-fog-sphere-20","dispel.strong-wind-disperses-fog","scaling.radius-per-slot","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","fog":{"radius_feet":20,"shape":"sphere","origin":"point-in-range","obscurement":"heavily-obscured","dispersal":"strong-wind","dispersal_example":"gust-of-wind","additional_radius_feet_per_slot":20}}}`,
		`{"id":"srd-5.2.1:content/spells/forbiddance","name":"Forbiddance","source_id":"srd-5.2.1:content/spells/forbiddance","source_page":133,"scope":"complete-source-record","level":6,"school":"abjuration","lists":["cleric"],"casting_time":"10-minutes","range":{"kind":"touch","feet":0},"target":{"shape":"created-effect","count":1,"eligibility":"floor-area-up-to-40000-square-feet"},"components":{"verbal":true,"somatic":true,"material":{"description":"ruby dust worth 1,000+ GP","cost_cp":100000,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"ward","mechanic":"protection.forbiddance-teleport-portal-ward","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"planar","mechanic":"protection.forbiddance-planar-travel-proofing","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.forbiddance-5d10","timing":"chosen-type-creature-first-enters-on-turn-or-ends-turn-in-area","dice_count":5,"die_sides":10,"damage_type":"radiant-or-necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"password","mechanic":"exemption.forbiddance-password","timing":"creature-speaks-password-as-it-enters","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"overlap","mechanic":"limit.forbiddance-no-overlap","timing":"when-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"permanence","mechanic":"duration.forbiddance-daily-30-day-permanence","timing":"cast-every-day-for-30-days-in-same-location","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.forbiddance-teleport-portal-ward","protection.forbiddance-planar-travel-proofing","damage.forbiddance-5d10","exemption.forbiddance-password","limit.forbiddance-no-overlap","duration.forbiddance-daily-30-day-permanence","cast.ten-minutes","cast.ritual","components.priced-material","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-ended","forbiddance":{"area_square_feet":40000,"height_feet":30,"height_from":"above-the-floor","blocked_travel":["teleporting-into-area","using-portals-to-enter-area"],"portal_example":"gate","planar_proofing":["astral-plane","ethereal-plane","feywild","shadowfell","plane-shift"],"damaged_types":["aberrations","celestials","elementals","fey","fiends","undead"],"type_selection":"choose-one-or-more-when-casting","damage_triggers":["enters-area-first-time-on-turn","ends-turn-in-area"],"damage_dice_count":5,"damage_die_sides":10,"damage_type_choices":["radiant","necrotic"],"damage_type_selection":"chosen-when-casting","password":"may-designate-when-casting","password_exemption":"speaks-password-as-it-enters-takes-no-damage","overlap_rule":"area-cannot-overlap-another-forbiddance","daily_casting_days":30,"daily_casting_place":"same-location","daily_casting_result":"lasts-until-dispelled","material_consumed":"on-the-last-casting"}}}`,
		`{"id":"srd-5.2.1:content/spells/forcecage","name":"Forcecage","source_id":"srd-5.2.1:content/spells/forcecage","source_page":134,"scope":"complete-source-record","level":7,"school":"evocation","lists":["bard","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":100},"target":{"shape":"created-effect","count":1,"eligibility":"area-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"ruby dust worth 1,500+ GP, which the spell consumes","cost_cp":150000,"consumed":true}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"prison","mechanic":"area.forcecage-prison-cage-or-box","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"trap","mechanic":"trap.forcecage-creatures-inside-or-pushed-out","timing":"when-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"exit","mechanic":"save.forcecage-teleport-or-planar-exit","timing":"creature-tries-teleportation-or-interplanar-travel-to-leave","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ethereal","mechanic":"block.forcecage-ethereal-travel","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dispel","mechanic":"immunity.forcecage-dispel-magic","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.forcecage-prison-cage-or-box","trap.forcecage-creatures-inside-or-pushed-out","save.forcecage-teleport-or-planar-exit","block.forcecage-ethereal-travel","immunity.forcecage-dispel-magic","components.priced-material","components.material","concentration.source-duration"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","forcecage":{"prison":"immobile-invisible-cube-of-magical-force","forms":["cage","solid-box"],"form_choice":"caster-choice","cage_max_side_feet":20,"cage_bar_diameter":"1/2-inch","cage_bar_spacing":"1/2-inch","box_max_side_feet":10,"box_blocks":["any-matter-passing-through","spells-cast-into-or-out-of-area"],"trapped_when":"completely-inside-cage-area-when-cast","pushed_away":"partially-inside-or-too-large-pushed-away-from-center-until-completely-outside","nonmagical_exit":"cannot-leave-by-nonmagical-means","travel_save":"charisma","travel_triggers":["teleportation","interplanar-travel"],"save_success":"can-use-that-magic-to-exit","save_failure":"does-not-exit-and-wastes-spell-or-effect","ethereal":"cage-extends-into-ethereal-plane-blocking-ethereal-travel","dispel":"cannot-be-dispelled-by-dispel-magic"}}}`,
		`{"id":"srd-5.2.1:content/spells/foresight","name":"Foresight","source_id":"srd-5.2.1:content/spells/foresight","source_page":134,"scope":"complete-source-record","level":9,"school":"divination","lists":["bard","druid","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a hummingbird feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"advantage","mechanic":"buff.foresight-d20-advantage","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"defense","mechanic":"buff.foresight-attackers-disadvantage","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.foresight-recast","timing":"caster-casts-spell-again","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["buff.foresight-d20-advantage","buff.foresight-attackers-disadvantage","ending.foresight-recast","cast.one-minute","components.material"],"slot_levels":[9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended-or-recast","foresight":{"target":"willing-creature-touched","grants":"limited-ability-to-see-into-immediate-future","target_advantage":"d20-tests","others_disadvantage":"attack-rolls-against-target","ends_early":"caster-casts-it-again"}}}`,
		`{"id":"srd-5.2.1:content/spells/freedom-of-movement","name":"Freedom of Movement","source_id":"srd-5.2.1:content/spells/freedom-of-movement","source_page":134,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["bard","cleric","druid","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a leather strap","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["protection"]},"effects":[{"id":"protection","mechanic":"protection.freedom-of-movement-terrain-and-magic","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"swim","mechanic":"movement.freedom-of-movement-swim-speed","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"escape","mechanic":"movement.freedom-of-movement-escape-restraints","timing":"spend-5-feet-of-movement","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.freedom-of-movement-terrain-and-magic","movement.freedom-of-movement-swim-speed","movement.freedom-of-movement-escape-restraints","scaling.targets-per-slot","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","freedom":{"target":"willing-creature-touched","difficult_terrain":"movement-unaffected","blocked_magic":["reduce-speed","paralyzed-condition","restrained-condition"],"swim_speed":"equal-to-speed","escape_cost_feet":5,"escape_result":"automatically-escape-nonmagical-restraints","escape_examples":["manacles","creature-imposing-the-grappled-condition"],"base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/freezing-sphere","name":"Freezing Sphere","source_id":"srd-5.2.1:content/spells/freezing-sphere","source_page":134,"scope":"complete-source-record","level":6,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":300},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"a miniature crystal sphere","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":6,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.freezing-sphere","timing":"failed-save-full-successful-save-half","dice_count":10,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"freeze","mechanic":"effect.freezing-sphere-ice","timing":"globe-strikes-body-of-water","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"break","mechanic":"check.freezing-sphere-break-free","timing":"action-by-trapped-creature","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"globe","mechanic":"object.freezing-sphere-held-globe","timing":"refrain-from-firing-after-casting","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.freezing-sphere","effect.freezing-sphere-ice","check.freezing-sphere-break-free","object.freezing-sphere-held-globe","area.sphere-60-radius","save.constitution-half","scaling.per-slot","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","freezing":{"radius_feet":60,"origin":"globe-streaks-from-caster-to-chosen-point-in-range","save":"constitution","fail_damage":"full","success_damage":"half","damage_dice_count":10,"damage_die_sides":6,"damage_type":"cold","slot_dice_count":1,"slot_die_sides":6,"water_freeze_depth_inches":6,"water_freeze_area_feet":30,"ice_seconds":60,"trapped_creatures":"creatures-swimming-on-surface-of-frozen-water","trapped_condition":"restrained","break_free_action":"action","break_free_check":"strength-athletics","break_free_dc":"caster-spell-save-dc","hold_option":"refrain-from-firing-after-completing-casting","held_globe_size":"sling-bullet","held_globe_temperature":"cool-to-the-touch","held_globe_location":"caster-hand","thrown_by":["caster","creature-given-the-globe"],"throw_range_feet":40,"sling_range":"slings-normal-range","impact_result":"shatters-with-same-effect-as-normal-casting","set_down":"can-set-down-without-shattering","explode_seconds":60,"explode_condition":"globe-has-not-already-shattered"}}}`,
		`{"id":"srd-5.2.1:content/spells/gaseous-form","name":"Gaseous Form","source_id":"srd-5.2.1:content/spells/gaseous-form","source_page":134,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of gauze","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":3,"dice_count":0,"die_sides":0,"effect_ids":["form"]},"effects":[{"id":"form","mechanic":"transformation.gaseous-form","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"defense","mechanic":"defense.gaseous-form-resistances-and-advantage","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"restrictions","mechanic":"restriction.gaseous-form-no-speech-objects-or-attacks","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.gaseous-form-zero-hit-points-or-magic-action","timing":"target-drops-to-0-hit-points-or-takes-magic-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.gaseous-form","defense.gaseous-form-resistances-and-advantage","restriction.gaseous-form-no-speech-objects-or-attacks","ending.gaseous-form-zero-hit-points-or-magic-action","scaling.targets-per-slot","components.material","concentration.source-duration"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","gaseous":{"source_pages":[134,135],"target":"willing-creature-touched","transforms":"target-and-everything-it-is-wearing-and-carrying","form":"misty-cloud","ends_when":["target-drops-to-0-hit-points","target-takes-magic-action-to-end-spell-on-itself"],"movement_only":"fly-speed","fly_speed_feet":10,"hover":true,"enter_occupied":"can-enter-and-occupy-space-of-another-creature","resistances":["bludgeoning","piercing","slashing"],"condition_immunity":"prone","save_advantage":["strength","dexterity","constitution"],"narrow_openings":"can-pass-through-narrow-openings","liquids":"treated-as-solid-surfaces","prohibited":["talk","manipulate-objects","attack","cast-spells"],"carried_objects":"cannot-be-dropped-used-or-interacted-with","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/gate","name":"Gate","source_id":"srd-5.2.1:content/spells/gate","source_page":135,"scope":"complete-source-record","level":9,"school":"conjuration","lists":["cleric","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a diamond worth 5,000+ GP","cost_cp":500000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"portal","mechanic":"portal.gate-planar","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"named","mechanic":"summon.gate-named-creature","timing":"caster-speaks-creature-name-when-casting","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"rulers","mechanic":"restriction.gate-planar-ruler-prevention","timing":"deity-or-planar-ruler-in-presence-or-domain","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["portal.gate-planar","summon.gate-named-creature","restriction.gate-planar-ruler-prevention","components.priced-material","components.material","concentration.source-duration"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","gate":{"portal_shape":"circular-opening","diameter_min_feet":5,"diameter_max_feet":20,"orientation":"any-direction-caster-chooses","destination":"precise-location-on-different-plane","destination_visible":true,"sides":"front-and-back-on-each-plane","travel":"only-by-moving-through-front","arrival":"instantly-transported-to-unoccupied-space-nearest-portal-on-other-plane","planar_rulers":"can-prevent-portals-opening-in-presence-or-domain","named_creature":"may-speak-name-of-specific-creature-when-casting","invalid_names":["pseudonym","title","nickname"],"named_on_other_plane":"portal-opens-next-to-named-creature-and-transports-it-to-nearest-unoccupied-space-on-casters-side","creature_control":"no-special-power-over-creature","creature_behavior":"free-to-act-as-gm-deems-appropriate","creature_may":["leave","attack-caster","help-caster"]}}}`,
		`{"id":"srd-5.2.1:content/spells/geas","name":"Geas","source_id":"srd-5.2.1:content/spells/geas","source_page":135,"scope":"complete-source-record","level":5,"school":"enchantment","lists":["bard","cleric","druid","paladin","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":2592000,"scaling":{"kind":"slot-duration","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["command"]},"effects":[{"id":"command","mechanic":"charm.geas-command","timing":"target-fails-wisdom-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.geas-psychic","timing":"acts-directly-counter-to-command-no-more-than-once-each-day","dice_count":5,"die_sides":10,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"suicidal","mechanic":"ending.geas-suicidal-command","timing":"caster-issues-suicidal-command","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.geas-listed-spells","timing":"remove-curse-greater-restoration-or-wish-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["charm.geas-command","damage.geas-psychic","ending.geas-suicidal-command","ending.geas-listed-spells","cast.one-minute","scaling.geas-slot-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended-suicidal-command-or-listed-spell","geas":{"save":"wisdom","auto_success":"target-cannot-understand-command","condition":"charmed","command_kinds":["carry-out-service","refrain-from-action-or-course-of-activity"],"damage_dice_count":5,"damage_die_sides":10,"damage_type":"psychic","damage_trigger":"acts-directly-counter-to-command","damage_limit":"no-more-than-once-each-day","excluded_command":"activity-that-would-result-in-certain-death","suicidal_result":"spell-ends","ended_by":["remove-curse","greater-restoration","wish"],"slot_durations":[{"slot_level":7,"duration":"days","duration_days":365},{"slot_level":8,"duration":"days","duration_days":365},{"slot_level":9,"duration":"until-ended-by-listed-spell","duration_days":0}]}}}`,
		`{"id":"srd-5.2.1:content/spells/gentle-repose","name":"Gentle Repose","source_id":"srd-5.2.1:content/spells/gentle-repose","source_page":135,"scope":"complete-source-record","level":2,"school":"necromancy","lists":["cleric","paladin","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"corpse-or-remains-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"2 Copper Pieces, which the spell consumes","cost_cp":2,"consumed":true}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":864000,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"protection","mechanic":"protection.gentle-repose-decay-and-undeath","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"time","mechanic":"time.gentle-repose-extends-raising-limit","timing":"days-under-spell-do-not-count-against-time-limit","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.gentle-repose-decay-and-undeath","time.gentle-repose-extends-raising-limit","cast.ritual","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","repose":{"target":"corpse-or-other-remains-touched","protected_from":"decay","cannot_become":"undead","time_limit_effect":"days-under-spell-do-not-count-against-time-limit","time_limit_example":"raise-dead"}}}`,
		`{"id":"srd-5.2.1:content/spells/giant-insect","name":"Giant Insect","source_id":"srd-5.2.1:content/spells/giant-insect","source_page":135,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot-spell-level","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["summon"]},"effects":[{"id":"summon","mechanic":"summon.giant-insect","timing":"unoccupied-space-visible-in-range","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"control","mechanic":"control.giant-insect-verbal-commands","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"jab","mechanic":"damage.giant-insect-poison-jab-piercing","timing":"melee-attack-roll-hit","dice_count":1,"die_sides":6,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"jab-poison","mechanic":"damage.giant-insect-poison-jab-poison","timing":"melee-attack-roll-hit","dice_count":1,"die_sides":4,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"web","mechanic":"damage.giant-insect-web-bolt","timing":"ranged-attack-roll-hit-spider-only","dice_count":1,"die_sides":10,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"spew","mechanic":"condition.giant-insect-venomous-spew-poisoned","timing":"failed-save-centipede-only","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.giant-insect","control.giant-insect-verbal-commands","damage.giant-insect-poison-jab-piercing","damage.giant-insect-poison-jab-poison","damage.giant-insect-web-bolt","condition.giant-insect-venomous-spew-poisoned","creature.giant-insect-stat-block","initiative.shared-after-caster","scaling.stat-block-spell-level","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","insect":{"source_pages":[135,136],"forms":["centipede","spider","wasp"],"form_choice":"chosen-when-casting","form_details":"form-determines-certain-stat-block-details","placement":"unoccupied-space-visible-in-range","stat_block":"giant-insect","disappears":["drops-to-0-hit-points","spell-ends"],"ally":"caster-and-caster-allies","initiative":"shares-caster-initiative-count-takes-turn-immediately-after-caster","commands":"verbal-no-action-required","no_command":"takes-dodge-action-and-uses-movement-to-avoid-danger","spell_level_source":"spell-slot-level","giant_insect":{"name":"Giant Insect","size":"large","creature_type":"beast","alignment":"unaligned","armor_class_base":11,"armor_class_adds":"spell-level","hit_points_base":30,"hit_points_per_spell_level":10,"hit_points_above_level":4,"speed_feet":40,"climb_feet":40,"fly_feet":40,"fly_forms":["wasp"],"abilities":[{"ability":"strength","score":17,"modifier":3,"save":3},{"ability":"dexterity","score":13,"modifier":1,"save":1},{"ability":"constitution","score":15,"modifier":2,"save":2},{"ability":"intelligence","score":4,"modifier":-3,"save":-3},{"ability":"wisdom","score":14,"modifier":2,"save":2},{"ability":"charisma","score":3,"modifier":-4,"save":-4}],"darkvision_feet":60,"passive_perception":12,"languages":"understands-languages-caster-knows","challenge_rating":"none","xp":0,"proficiency_bonus":"equals-caster-proficiency-bonus","spider_climb":"climbs-difficult-surfaces-including-ceilings-without-ability-check","multiattack":"attacks-equal-to-half-spell-level-rounded-down","poison_jab":{"attack_kind":"melee-attack-roll","attack_bonus":"caster-spell-attack-modifier","reach_feet":10,"piercing_dice_count":1,"piercing_die_sides":6,"piercing_bonus":3,"piercing_adds_spell_level":true,"poison_dice_count":1,"poison_die_sides":4},"web_bolt":{"form":"spider","attack_kind":"ranged-attack-roll","attack_bonus":"caster-spell-attack-modifier","range_feet":60,"bludgeoning_dice_count":1,"bludgeoning_die_sides":10,"bludgeoning_bonus":3,"bludgeoning_adds_spell_level":true,"speed_result":"speed-reduced-to-0-until-start-of-insects-next-turn"},"venomous_spew":{"form":"centipede","action":"bonus-action","save":"constitution","save_dc":"caster-spell-save-dc","target_range_feet":10,"target":"one-creature-the-insect-can-see","failure":"poisoned-condition-until-start-of-insects-next-turn"}}}}}`,
		`{"id":"srd-5.2.1:content/spells/glibness","name":"Glibness","source_id":"srd-5.2.1:content/spells/glibness","source_page":136,"scope":"complete-source-record","level":8,"school":"enchantment","lists":["bard","warlock"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"check","mechanic":"check.glibness-charisma-15","timing":"when-caster-makes-charisma-check","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"truth","mechanic":"truth.glibness-magic-indicates-truthful","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["check.glibness-charisma-15","truth.glibness-magic-indicates-truthful"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","glibness":{"check_ability":"charisma","replace_roll_with":15,"replace_choice":"optional","truth_magic":"magic-that-determines-truthfulness-indicates-truthful","truth_condition":"no-matter-what-caster-says"}}}`,
	}
	batch := loadSpells207(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Fog Cloud", "Forbiddance", "Forcecage", "Foresight", "Freedom of Movement", "Freezing Sphere", "Gaseous Form", "Gate", "Geas", "Gentle Repose", "Giant Insect", "Glibness"}
	pages := []int{133, 133, 134, 134, 134, 134, 134, 135, 135, 135, 135, 136}
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
			for level := r.Level; level <= 9 && r.Level > 0; level++ {
				slots = append(slots, level)
			}
			if !slices.Equal(r.SlotLevels, slots) {
				t.Fatal("source slot metadata incomplete")
			}
		})
	}
}

func TestSpells207CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells207(t)
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
func TestSpells207MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_207.json")
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

func TestSpells207JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_207.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"radius_feet": 20`, `"radius_feet": 20, "radius_feet": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "RADIUS_FEET": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/fog-cloud"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"kind": "slot-radius"`, `"kind": "slot-duration"`},
		{`"kind": "slot-spell-level"`, `"kind": "slot-targets"`},
		{`"casting_time": "action"`, `"casting_time": "1-hour"`},
		{`"casting_time": "10-minutes"`, `"casting_time": "1-minute"`},
		{`"casting_seconds": 0`, `"casting_seconds": 60`},
		{`"casting_seconds": 600`, `"casting_seconds": 60`},
		{`"feet": 120`, `"feet": 121`},
		{`"cost_cp": 100000`, `"cost_cp": 100001`},
		{`"ritual": true`, `"ritual": false`},
		{`"consumed": true`, `"consumed": false`},
		{`"shape": "created-effect"`, `"shape": "unknown"`},
		{`"eligibility": "corpse-or-remains-touch"`, `"eligibility": "creature-touch"`},
		{`"score": 17`, `"score": 17, "score": 18`},
		{`"score": 17`, `"score": 17, "SCORE": 18`},
		{`"score": 17`, `"score": 17, "unknown": true`},
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

func TestSpells207InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"fog radius", 0, func(s *SpellRecord) { s.SourceDetails.Fog.RadiusFeet = 30 }},
		{"fog slot increment", 0, func(s *SpellRecord) { s.SourceDetails.Fog.AdditionalRadiusFeetPerSlot = 10 }},
		{"fog scaling level", 0, func(s *SpellRecord) { s.Scaling.AboveLevel = 2 }},
		{"fog concentration", 0, func(s *SpellRecord) { *s.Concentration = false }},
		{"fog components", 0, func(s *SpellRecord) { *s.Components.Somatic = false }},
		{"forbiddance area", 1, func(s *SpellRecord) { s.SourceDetails.Forbiddance.AreaSquareFeet = 4000 }},
		{"forbiddance damage type", 1, func(s *SpellRecord) { s.SourceDetails.Forbiddance.DamageTypeChoices[1] = "fire" }},
		{"forbiddance planar", 1, func(s *SpellRecord) {
			s.SourceDetails.Forbiddance.PlanarProofing = s.SourceDetails.Forbiddance.PlanarProofing[1:]
		}},
		{"forbiddance casting", 1, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"forbiddance ritual", 1, func(s *SpellRecord) { *s.Ritual = false }},
		{"forbiddance cost", 1, func(s *SpellRecord) { *s.Components.Material.CostCP = 150000 }},
		{"forbiddance consumed", 1, func(s *SpellRecord) { *s.Components.Material.Consumed = true }},
		{"forbiddance duration", 1, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"forcecage cage size", 2, func(s *SpellRecord) { s.SourceDetails.Forcecage.CageMaxSide = 10 }},
		{"forcecage box size", 2, func(s *SpellRecord) { s.SourceDetails.Forcecage.BoxMaxSide = 20 }},
		{"forcecage save", 2, func(s *SpellRecord) { s.SourceDetails.Forcecage.TravelSave = "wisdom" }},
		{"forcecage consumed", 2, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"forcecage cost", 2, func(s *SpellRecord) { *s.Components.Material.CostCP = 100000 }},
		{"foresight casting", 3, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"foresight advantage", 3, func(s *SpellRecord) { s.SourceDetails.Foresight.TargetAdvantage = "attack-rolls" }},
		{"foresight duration", 3, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"freedom swim", 4, func(s *SpellRecord) { s.SourceDetails.Freedom.SwimSpeed = "half-speed" }},
		{"freedom escape cost", 4, func(s *SpellRecord) { s.SourceDetails.Freedom.EscapeCostFeet = 10 }},
		{"freedom blocked magic", 4, func(s *SpellRecord) { s.SourceDetails.Freedom.BlockedMagic[2] = "grappled-condition" }},
		{"freedom target scaling", 4, func(s *SpellRecord) { s.SourceDetails.Freedom.AdditionalTargetsPerSlot = 2 }},
		{"freedom scaling level", 4, func(s *SpellRecord) { s.Scaling.AboveLevel = 3 }},
		{"freezing damage", 5, func(s *SpellRecord) { s.SourceDetails.Freezing.DamageDiceCount = 8 }},
		{"freezing effect dice", 5, func(s *SpellRecord) { s.Effects[0].DiceCount = 8 }},
		{"freezing scaling level", 5, func(s *SpellRecord) { s.Scaling.AboveLevel = 5 }},
		{"freezing ice depth", 5, func(s *SpellRecord) { s.SourceDetails.Freezing.WaterFreezeDepth = 12 }},
		{"freezing throw range", 5, func(s *SpellRecord) { s.SourceDetails.Freezing.ThrowRangeFeet = 30 }},
		{"freezing thrower", 5, func(s *SpellRecord) { s.SourceDetails.Freezing.ThrownBy[1] = "any-creature" }},
		{"gaseous fly speed", 6, func(s *SpellRecord) { s.SourceDetails.Gaseous.FlySpeedFeet = 30 }},
		{"gaseous continuation page", 6, func(s *SpellRecord) { s.SourceDetails.Gaseous.SourcePages = []int{134} }},
		{"gaseous resistance", 6, func(s *SpellRecord) { s.SourceDetails.Gaseous.Resistances[2] = "fire" }},
		{"gaseous hover", 6, func(s *SpellRecord) { s.SourceDetails.Gaseous.Hover = false }},
		{"gaseous scaling level", 6, func(s *SpellRecord) { s.Scaling.AboveLevel = 4 }},
		{"gate diameter", 7, func(s *SpellRecord) { s.SourceDetails.Gate.DiameterMaxFeet = 30 }},
		{"gate invalid name", 7, func(s *SpellRecord) { s.SourceDetails.Gate.InvalidNames = s.SourceDetails.Gate.InvalidNames[:2] }},
		{"gate cost", 7, func(s *SpellRecord) { *s.Components.Material.CostCP = 100000 }},
		{"gate duration", 7, func(s *SpellRecord) { s.DurationSeconds = 600 }},
		{"geas die", 8, func(s *SpellRecord) { s.SourceDetails.Geas.DamageDieSides = 8 }},
		{"geas slot duration", 8, func(s *SpellRecord) { s.SourceDetails.Geas.SlotDurations[1].DurationDays = 30 }},
		{"geas slot table", 8, func(s *SpellRecord) { s.SourceDetails.Geas.SlotDurations = s.SourceDetails.Geas.SlotDurations[:2] }},
		{"geas ender", 8, func(s *SpellRecord) { s.SourceDetails.Geas.EndedBy[2] = "dispel-magic" }},
		{"geas duration", 8, func(s *SpellRecord) { s.DurationSeconds = 864000 }},
		{"geas casting", 8, func(s *SpellRecord) { s.CastingTime = "10-minutes" }},
		{"geas scaling mechanic", 8, func(s *SpellRecord) { s.RequiredMechanics[5] = "scaling.unknown" }},
		{"repose target", 9, func(s *SpellRecord) { s.SourceDetails.Repose.CannotBecome = "construct" }},
		{"repose consumed", 9, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"repose duration", 9, func(s *SpellRecord) { s.DurationSeconds = 86400 }},
		{"repose ritual", 9, func(s *SpellRecord) { *s.Ritual = false }},
		{"insect continuation page", 10, func(s *SpellRecord) { s.SourceDetails.Insect.SourcePages = []int{135} }},
		{"insect form", 10, func(s *SpellRecord) { s.SourceDetails.Insect.Forms[2] = "beetle" }},
		{"insect armor class", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.ArmorClassBase = 12 }},
		{"insect hit points", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.HitPointsPerLevel = 5 }},
		{"insect ability", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.Abilities[3].Modifier = 3 }},
		{"insect ability save", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.Abilities[5].Save = 4 }},
		{"insect jab", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.PoisonJab.PoisonDieSides = 6 }},
		{"insect web", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.WebBolt.RangeFeet = 30 }},
		{"insect spew", 10, func(s *SpellRecord) { s.SourceDetails.Insect.GiantInsect.VenomousSpew.Save = "dexterity" }},
		{"insect spell level source", 10, func(s *SpellRecord) { s.SourceDetails.Insect.SpellLevelSource = "spell-level" }},
		{"insect scaling kind", 10, func(s *SpellRecord) { s.Scaling.Kind = "none" }},
		{"insect effect dice", 10, func(s *SpellRecord) { s.Effects[4].DieSides = 8 }},
		{"glibness roll", 11, func(s *SpellRecord) { s.SourceDetails.Glibness.ReplaceRollWith = 20 }},
		{"glibness ability", 11, func(s *SpellRecord) { s.SourceDetails.Glibness.CheckAbility = "wisdom" }},
		{"glibness concentration", 11, func(s *SpellRecord) { *s.Concentration = true }},
		{"glibness components", 11, func(s *SpellRecord) { *s.Components.Verbal = false }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 3, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Gate = &SpellGate{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"payload on wrong record", 0, func(s *SpellRecord) { s.SourceDetails.Glibness = &SpellGlibness{}; s.SourceDetails.Fog = nil }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
		{"source details omitted", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells207(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells207RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_207.json")
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

func TestSpells207WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_207.json")
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
		`"score": 17, "score": 18`,
		`"score": 17, "SCORE": 18`,
		`"score": 17, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"score": 17`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells207PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_207.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"radius_feet", 1, 20},
		{"radius_feet", 2, 60},
		{"additional_radius_feet_per_slot", 1, 20},
		{"area_square_feet", 1, 40000},
		{"height_feet", 1, 30},
		{"daily_casting_days", 1, 30},
		{"damage_dice_count", 1, 5},
		{"damage_dice_count", 2, 10},
		{"damage_dice_count", 3, 5},
		{"cage_max_side_feet", 1, 20},
		{"box_max_side_feet", 1, 10},
		{"escape_cost_feet", 1, 5},
		{"base_targets", 1, 1},
		{"base_targets", 2, 1},
		{"additional_targets_per_slot", 1, 1},
		{"additional_targets_per_slot", 2, 1},
		{"water_freeze_depth_inches", 1, 6},
		{"water_freeze_area_feet", 1, 30},
		{"ice_seconds", 1, 60},
		{"throw_range_feet", 1, 40},
		{"explode_seconds", 1, 60},
		{"fly_speed_feet", 1, 10},
		{"diameter_min_feet", 1, 5},
		{"diameter_max_feet", 1, 20},
		{"slot_level", 1, 7},
		{"slot_level", 2, 8},
		{"slot_level", 3, 9},
		{"duration_days", 1, 365},
		{"duration_days", 2, 365},
		{"armor_class_base", 1, 11},
		{"hit_points_base", 1, 30},
		{"hit_points_per_spell_level", 1, 10},
		{"hit_points_above_level", 1, 4},
		{"speed_feet", 1, 40},
		{"climb_feet", 1, 40},
		{"fly_feet", 1, 40},
		{"darkvision_feet", 1, 60},
		{"passive_perception", 1, 12},
		{"reach_feet", 1, 10},
		{"range_feet", 1, 60},
		{"target_range_feet", 1, 10},
		{"replace_roll_with", 1, 15},
		{"score", 1, 17},
		{"score", 2, 13},
		{"score", 3, 15},
		{"score", 4, 4},
		{"score", 5, 14},
		{"score", 6, 3},
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
// inventory assigns to catalogue issue #207.
func TestSpells207MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 207 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells207(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells207ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_207.json")
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
	want := "SRD 5.2.1 PDF pp.133-136 adapted into twelve complete offline spell families for issue #207, including page-break continuation text for Gaseous Form (pp.134-135) and Giant Insect (pp.135-136) and the Giant Insect stat block. No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	// Geas sits wholly on p.135 and Glibness wholly on p.136.
	if strings.Contains(root.Source.Changes, "Geas") || strings.Contains(root.Source.Changes, "Glibness") {
		t.Fatal("only Gaseous Form and Giant Insect cross a page break")
	}
}
