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

func loadSpells214(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_214.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.150-152 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells214SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/mirror-image","name":"Mirror Image","source_id":"srd-5.2.1:content/spells/mirror-image","source_page":150,"scope":"complete-source-record","level":2,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"duplicates","mechanic":"illusion.mirror-image-duplicates","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"redirect","mechanic":"defense.mirror-image-redirect","timing":"when-attack-roll-hits-caster","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.mirror-image-duplicates","defense.mirror-image-redirect"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-all-duplicates-destroyed","mirror_image":{"duplicate_count":3,"location":"casters-space","behavior":["move-with-caster","mimic-casters-actions","shift-position-so-real-image-cannot-be-tracked"],"redirect_trigger":"creature-hits-caster-with-attack-roll","redirect_die_sides":6,"dice_per_remaining_duplicate":1,"redirect_threshold":3,"redirect_result":"duplicate-hit-instead-and-destroyed","duplicate_ignores":"all-other-damage-and-effects","end_condition":"all-duplicates-destroyed","unaffected_creatures":["blinded-condition","blindsight","truesight"]}}}`,
		`{"id":"srd-5.2.1:content/spells/mislead","name":"Mislead","source_id":"srd-5.2.1:content/spells/mislead","source_page":150,"scope":"complete-source-record","level":5,"school":"illusion","lists":["bard","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":false,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"invisible","mechanic":"condition.mislead-invisible","timing":"until-attack-damage-or-spell","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"double","mechanic":"illusion.mislead-double","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"control","mechanic":"control.mislead-double-magic-action","timing":"magic-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.mislead-invisible","illusion.mislead-double","control.mislead-double-magic-action","sense.mislead-through-double","ending.mislead-invisibility","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","mislead":{"invisibility_timing":"same-time-double-appears","double_location":"where-caster-is-standing","double_duration":"lasts-for-duration","invisibility_end_triggers":["caster-makes-attack-roll","caster-deals-damage","caster-casts-spell"],"move_action":"magic-action","move_speed_multiplier":2,"double_behaviors":["gesture","speak","behave-as-caster-chooses"],"double_traits":["intangible","invulnerable"],"senses":"see-and-hear-through-double-as-if-located-there"}}}`,
		`{"id":"srd-5.2.1:content/spells/misty-step","name":"Misty Step","source_id":"srd-5.2.1:content/spells/misty-step","source_page":150,"scope":"complete-source-record","level":2,"school":"conjuration","lists":["sorcerer","warlock","wizard"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"teleport","mechanic":"teleport.misty-step","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.misty-step"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","misty_step":{"appearance":"briefly-surrounded-by-silvery-mist","teleport_max_feet":30,"destination":"unoccupied-space-caster-can-see"}}}`,
		`{"id":"srd-5.2.1:content/spells/modify-memory","name":"Modify Memory","source_id":"srd-5.2.1:content/spells/modify-memory","source_page":150,"scope":"complete-source-record","level":5,"school":"enchantment","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-memory-window","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["memory"]},"effects":[{"id":"charm","mechanic":"condition.modify-memory-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"memory","mechanic":"memory.modify-memory","timing":"while-charmed-takes-hold-when-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.modify-memory-charmed","memory.modify-memory","save.wisdom-advantage-if-fighting","ending.modify-memory-damage-or-spell","scaling.memory-window-by-slot","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-damage-or-other-spell","modify_memory":{"save":"wisdom","save_advantage":"caster-fighting-creature","failure_condition":"charmed","charmed_condition":"incapacitated","charmed_awareness":"unaware-of-surroundings-but-can-hear-caster","interrupt_triggers":["target-takes-any-damage","target-targeted-by-another-spell"],"interrupt_result":"spell-ends-and-no-memories-modified","event_window_hours":24,"event_max_minutes":10,"memory_options":["permanently-eliminate-all-memory-of-event","recall-event-with-perfect-clarity","change-memory-of-event-details","create-memory-of-other-event"],"description_requirements":["caster-must-speak-to-target","target-must-understand-casters-language"],"gaps":"target-mind-fills-gaps-in-description","unfinished_description":"memory-not-altered-if-spell-ends-first","takes_hold":"when-spell-ends","behavior_effect":"modified-memory-does-not-necessarily-affect-behavior","contradictions":["natural-inclinations","alignment","beliefs"],"illogical_memory":"dismissed-as-bad-dream","illogical_example":"enjoyed-swimming-in-acid","gm_discretion":"may-deem-memory-too-nonsensical-to-affect-creature","restored_by":["remove-curse","greater-restoration"],"slot_windows":[{"slot":6,"days":7},{"slot":7,"days":30},{"slot":8,"days":365}],"any_past_time_slot":9}}}`,
		`{"id":"srd-5.2.1:content/spells/moonbeam","name":"Moonbeam","source_id":"srd-5.2.1:content/spells/moonbeam","source_page":150,"scope":"complete-source-record","level":2,"school":"evocation","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"cylinder","count":0,"eligibility":"creatures-in-cylinder-centered-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a moonseed leaf","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"light","mechanic":"light.moonbeam-dim","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.moonbeam","timing":"appear-moves-into-enter-or-end-turn","dice_count":2,"die_sides":10,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"movement","mechanic":"movement.moonbeam-magic-action","timing":"magic-action-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.moonbeam-dim","damage.moonbeam","movement.moonbeam-magic-action","save.constitution-half","transformation.moonbeam-revert-shape-shift","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","moonbeam":{"source_pages":[150,151],"beam":"silvery-beam-of-pale-light","cylinder_radius_feet":5,"cylinder_height_feet":40,"center":"point-within-range","light":"dim-light-fills-cylinder","move_action":"magic-action-on-later-turns","move_max_feet":60,"save":"constitution","damage_dice_count":2,"damage_die_sides":10,"damage_type":"radiant","shape_shifted_failure":"reverts-to-true-form-and-cannot-shape-shift-until-leaving-cylinder","shape_shift_example":"polymorph","save_success":"half-damage-only","save_triggers":["cylinder-appears-with-creature-in-it","area-moves-into-creature-space","creature-enters-area","creature-ends-turn-in-area"],"saves_per_turn":1,"slot_dice_count":1,"slot_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/move-earth","name":"Move Earth","source_id":"srd-5.2.1:content/spells/move-earth","source_page":151,"scope":"complete-source-record","level":6,"school":"transmutation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"square","count":0,"eligibility":"terrain-area-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a miniature shovel","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":7200,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"terrain.move-earth-reshape","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.move-earth-reshape","limit.move-earth-materials-and-extent","concentration.source-duration","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","move_earth":{"max_side_feet":40,"materials":["dirt","sand","clay"],"changes":["raise-or-lower-elevation","create-or-fill-trench","erect-or-flatten-wall","form-pillar"],"extent_limit":"half-areas-largest-dimension","example_side_feet":40,"example_max_change_feet":20,"change_minutes":10,"creature_safety":"creatures-usually-cannot-be-trapped-or-injured-by-movement","new_area_interval_minutes":10,"new_area":"choose-new-area-within-range-while-concentrating","excluded":["natural-stone","stone-construction"],"rocks_and_structures":"shift-to-accommodate-new-terrain","unstable_structures":"might-collapse","plants":"not-directly-affected-and-carried-with-moved-earth"}}}`,
		`{"id":"srd-5.2.1:content/spells/nondetection","name":"Nondetection","source_id":"srd-5.2.1:content/spells/nondetection","source_page":151,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["bard","ranger","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"willing-creature-place-or-object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of diamond dust worth 25+ GP, which the spell consumes","cost_cp":2500,"consumed":true}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"ward","mechanic":"protection.nondetection","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.nondetection","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","nondetection":{"hidden_from":"divination-spells","target_choices":["willing-creature","place","object"],"max_dimension_feet":10,"protections":["cannot-be-targeted-by-divination-spells","cannot-be-perceived-through-magical-scrying-sensors"]}}}`,
		`{"id":"srd-5.2.1:content/spells/pass-without-trace","name":"Pass without Trace","source_id":"srd-5.2.1:content/spells/pass-without-trace","source_page":151,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"caster-and-chosen-creatures-in-aura"},"components":{"verbal":true,"somatic":true,"material":{"description":"ashes from burned mistletoe","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"stealth","mechanic":"bonus.pass-without-trace-stealth","timing":"while-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":10,"add_casting_modifier":false,"half_initial_only":false},{"id":"tracks","mechanic":"tracks.pass-without-trace-none","timing":"while-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["bonus.pass-without-trace-stealth","tracks.pass-without-trace-none","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","pass_without_trace":{"aura":"concealing-aura","emanation_feet":30,"beneficiaries":"caster-and-each-chosen-creature-while-in-aura","check":"dexterity-stealth","check_bonus":10,"tracks":"leave-no-tracks"}}}`,
		`{"id":"srd-5.2.1:content/spells/passwall","name":"Passwall","source_id":"srd-5.2.1:content/spells/passwall","source_page":151,"scope":"complete-source-record","level":5,"school":"transmutation","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"visible-point-on-wood-plaster-or-stone-surface"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of sesame seeds","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"passage","mechanic":"passage.passwall","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ejection","mechanic":"ejection.passwall","timing":"when-opening-disappears","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["passage.passwall","ejection.passwall","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","passwall":{"surfaces":["wooden","plaster","stone"],"surface_examples":["wall","ceiling","floor"],"dimensions_choice":"chosen-by-caster","max_width_feet":5,"max_height_feet":8,"max_depth_feet":20,"stability":"creates-no-instability-in-surrounding-structure","ejection":"creatures-and-objects-safely-ejected-to-unoccupied-space-nearest-cast-surface"}}}`,
		`{"id":"srd-5.2.1:content/spells/phantasmal-force","name":"Phantasmal Force","source_id":"srd-5.2.1:content/spells/phantasmal-force","source_page":151,"scope":"complete-source-record","level":2,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of fleece","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"phantasm","mechanic":"illusion.phantasmal-force","timing":"on-failed-intelligence-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.phantasmal-force","timing":"each-of-casters-turns","dice_count":2,"die_sides":8,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.phantasmal-force","damage.phantasmal-force","save.intelligence-negates","check.phantasmal-force-investigation","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-discerned","phantasmal_force":{"source_pages":[151,152],"save":"intelligence","phenomena":["phantasmal-object","phantasmal-creature","other-phenomenon"],"max_cube_feet":10,"perceivable_by":"target-only","stimuli":["sound","temperature","other-stimuli"],"examine_action":"study","examine_check":"intelligence-investigation","examine_dc":"caster-spell-save-dc","examine_success":"target-realizes-illusion-and-spell-ends","belief":"treats-phantasm-as-real-and-rationalizes-illogical-outcomes","rationalization_example":"survives-fall-through-phantasmal-bridge-and-still-believes-bridge-exists","damage_requirement":"phantasm-represents-dangerous-creature-or-hazard","damage_dice_count":2,"damage_die_sides":8,"damage_type":"psychic","damage_timing":"on-each-of-casters-turns","damage_reach_feet":5,"damage_position":"target-in-phantasm-area-or-within-reach","damage_perception":"perceived-as-type-appropriate-to-illusion"}}}`,
		`{"id":"srd-5.2.1:content/spells/phantasmal-killer","name":"Phantasmal Killer","source_id":"srd-5.2.1:content/spells/phantasmal-killer","source_page":152,"scope":"complete-source-record","level":4,"school":"illusion","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":4,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.phantasmal-killer","timing":"immediate-and-on-failed-repeat-save","dice_count":4,"die_sides":10,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"impairment","mechanic":"disadvantage.phantasmal-killer","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.phantasmal-killer","disadvantage.phantasmal-killer","save.wisdom-half","save.phantasmal-killer-repeat","scaling.per-slot","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-successful-save","phantasmal_killer":{"illusion":"deepest-fears-visible-only-to-target","save":"wisdom","damage_dice_count":4,"damage_die_sides":10,"damage_type":"psychic","failure_penalties":["disadvantage-on-ability-checks","disadvantage-on-attack-rolls"],"save_success":"half-damage-and-spell-ends","repeat_timing":"end-of-each-of-its-turns","repeat_failure":"takes-psychic-damage-again","repeat_success":"spell-ends","slot_dice_count":1,"slot_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/phantom-steed","name":"Phantom Steed","source_id":"srd-5.2.1:content/spells/phantom-steed","source_page":152,"scope":"complete-source-record","level":3,"school":"illusion","lists":["wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-ground-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"steed","mechanic":"summon.phantom-steed","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.phantom-steed","ending.phantom-steed-damage-and-fade","cast.one-minute","cast.ritual"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended-or-steed-damaged","phantom_steed":{"creature_size":"large","creature_nature":"quasi-real-horselike","appearance":"decided-by-caster","equipment":["saddle","bit","bridle"],"equipment_vanishes":"puff-of-smoke-when-carried-too-far-from-steed","equipment_max_distance_feet":10,"rider":"caster-or-creature-caster-chooses","stat_block":"riding-horse","speed_feet":100,"travel_miles_per_hour":13,"fade":"gradually-fades-when-spell-ends","dismount_minutes":1,"early_end":"steed-takes-any-damage"}}}`,
	}
	batch := loadSpells214(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Mirror Image", "Mislead", "Misty Step", "Modify Memory", "Moonbeam", "Move Earth", "Nondetection", "Pass without Trace", "Passwall", "Phantasmal Force", "Phantasmal Killer", "Phantom Steed"}
	pages := []int{150, 150, 150, 150, 150, 151, 151, 151, 151, 151, 152, 152}
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

func TestSpells214CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells214(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 240 {
		t.Fatalf("expected 240 records, got %d", len(all.Records))
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
func TestSpells214MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_214.json")
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

func TestSpells214JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_214.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"duplicate_count": 3`, `"duplicate_count": 3, "duplicate_count": 4`},
		{`"duplicate_count": 3`, `"duplicate_count": 3, "DUPLICATE_COUNT": 4`},
		{`"duplicate_count": 3`, `"duplicate_count": 3, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/mirror-image"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "memory"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 30`, `"feet": 31`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
		{`"ritual": true`, `"ritual": false`},
		{`"consumed": true`, `"consumed": false`},
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

func TestSpells214InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"mirror duplicates", 0, func(s *SpellRecord) { s.SourceDetails.MirrorImage.DuplicateCount = 4 }},
		{"mirror threshold", 0, func(s *SpellRecord) { s.SourceDetails.MirrorImage.RedirectThreshold = 6 }},
		{"mislead verbal", 1, func(s *SpellRecord) { *s.Components.Verbal = true }},
		{"mislead speed", 1, func(s *SpellRecord) { s.SourceDetails.Mislead.MoveSpeedMultiplier = 3 }},
		{"misty distance", 2, func(s *SpellRecord) { s.SourceDetails.MistyStep.TeleportMaxFeet = 60 }},
		{"memory table", 3, func(s *SpellRecord) { s.SourceDetails.ModifyMemory.SlotWindows[2].Days = 366 }},
		{"memory any time slot", 3, func(s *SpellRecord) { s.SourceDetails.ModifyMemory.AnyPastTimeSlot = 8 }},
		{"memory scaling kind", 3, func(s *SpellRecord) { s.Scaling.Kind = "slot-duration" }},
		{"moonbeam pages", 4, func(s *SpellRecord) { s.SourceDetails.Moonbeam.SourcePages = []int{150} }},
		{"moonbeam dice", 4, func(s *SpellRecord) { s.Effects[1].DiceCount = 3 }},
		{"earth side", 5, func(s *SpellRecord) { s.SourceDetails.MoveEarth.MaxSideFeet = 100 }},
		{"earth exclusions", 5, func(s *SpellRecord) { s.SourceDetails.MoveEarth.Excluded = s.SourceDetails.MoveEarth.Excluded[:1] }},
		{"nondetection cost", 6, func(s *SpellRecord) { *s.Components.Material.CostCP = 2000 }},
		{"nondetection targets", 6, func(s *SpellRecord) { s.SourceDetails.Nondetection.TargetChoices[1] = "plane" }},
		{"trace bonus", 7, func(s *SpellRecord) { s.Effects[0].Bonus = 5 }},
		{"trace radius", 7, func(s *SpellRecord) { s.SourceDetails.PassWithoutTrace.EmanationFeet = 10 }},
		{"passwall depth", 8, func(s *SpellRecord) { s.SourceDetails.Passwall.MaxDepthFeet = 30 }},
		{"passwall concentration", 8, func(s *SpellRecord) { *s.Concentration = true }},
		{"force save", 9, func(s *SpellRecord) { s.SourceDetails.PhantasmalForce.Save = "wisdom" }},
		{"force cube", 9, func(s *SpellRecord) { s.SourceDetails.PhantasmalForce.MaxCubeFeet = 20 }},
		{"killer success", 10, func(s *SpellRecord) { s.SourceDetails.PhantasmalKiller.SaveSuccess = "half-damage" }},
		{"killer slot level", 10, func(s *SpellRecord) { s.Scaling.AboveLevel = 3 }},
		{"steed speed", 11, func(s *SpellRecord) { s.SourceDetails.PhantomSteed.SpeedFeet = 60 }},
		{"steed ritual", 11, func(s *SpellRecord) { *s.Ritual = false }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.PhantomSteed = &SpellPhantomSteed{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells214(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells214RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_214.json")
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

func TestSpells214WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_214.json")
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
		`"duplicate_count": 3, "duplicate_count": 4`,
		`"duplicate_count": 3, "DUPLICATE_COUNT": 4`,
		`"duplicate_count": 3, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"duplicate_count": 3`), []byte(replacement), 1)
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
func TestSpells214PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_214.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"duplicate_count", 1, 3},
		{"redirect_die_sides", 1, 6},
		{"dice_per_remaining_duplicate", 1, 1},
		{"redirect_threshold", 1, 3},
		{"move_speed_multiplier", 1, 2},
		{"teleport_max_feet", 1, 30},
		{"event_window_hours", 1, 24},
		{"event_max_minutes", 1, 10},
		{"slot", 1, 6},
		{"slot", 2, 7},
		{"slot", 3, 8},
		{"days", 1, 7},
		{"days", 2, 30},
		{"days", 3, 365},
		{"any_past_time_slot", 1, 9},
		{"cylinder_radius_feet", 1, 5},
		{"cylinder_height_feet", 1, 40},
		{"move_max_feet", 1, 60},
		{"damage_dice_count", 1, 2},
		{"damage_dice_count", 2, 2},
		{"damage_dice_count", 3, 4},
		{"damage_die_sides", 1, 10},
		{"damage_die_sides", 2, 8},
		{"damage_die_sides", 3, 10},
		{"saves_per_turn", 1, 1},
		{"slot_dice_count", 1, 1},
		{"slot_dice_count", 2, 1},
		{"slot_die_sides", 1, 10},
		{"slot_die_sides", 2, 10},
		{"max_side_feet", 1, 40},
		{"example_side_feet", 1, 40},
		{"example_max_change_feet", 1, 20},
		{"change_minutes", 1, 10},
		{"new_area_interval_minutes", 1, 10},
		{"max_dimension_feet", 1, 10},
		{"emanation_feet", 1, 30},
		{"check_bonus", 1, 10},
		{"max_width_feet", 1, 5},
		{"max_height_feet", 1, 8},
		{"max_depth_feet", 1, 20},
		{"max_cube_feet", 1, 10},
		{"damage_reach_feet", 1, 5},
		{"equipment_max_distance_feet", 1, 10},
		{"speed_feet", 1, 100},
		{"travel_miles_per_hour", 1, 13},
		{"dismount_minutes", 1, 1},
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
// inventory assigns to catalogue issue #214.
func TestSpells214MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 214 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells214(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells214ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_214.json")
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
	want := "SRD 5.2.1 PDF pp.150-152 adapted into twelve complete offline spell families for issue #214, including page-break continuation text for Moonbeam (pp.150-151) and Phantasmal Force (pp.151-152). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells214PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.PhantomSteed = loadSpells214(t).Records[11].SourceDetails.PhantomSteed
	if r.Validate() == nil {
		t.Fatal("batch 214 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells214RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells214(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell214Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell214Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
