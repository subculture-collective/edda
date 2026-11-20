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

func loadSpells206(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_206.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.131-133 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells206SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/find-the-path","name":"Find the Path","source_id":"srd-5.2.1:content/spells/find-the-path","source_page":131,"scope":"complete-source-record","level":6,"school":"divination","lists":["bard","cleric","druid"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a set of divination tools—such as cards or runes—worth 100+ GP","cost_cp":10000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"route","mechanic":"sense.find-the-path-route","timing":"duration-on-same-plane","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.find-the-path-route","sense.find-the-path-invalid-destinations","cast.one-minute","concentration.source-duration","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","find_path":{"destination":"location-you-name","familiarity_required":true,"route_sensed":"most-direct-physical-route","failing_destinations":["destination-on-another-plane","moving-destination","unspecific-destination"],"known_while_on_same_plane":["distance","direction"],"plane_requirement":"same-plane-as-destination","path_choice":"know-which-path-is-most-direct-at-each-choice"}}}`,
		`{"id":"srd-5.2.1:content/spells/finger-of-death","name":"Finger of Death","source_id":"srd-5.2.1:content/spells/finger-of-death","source_page":131,"scope":"complete-source-record","level":7,"school":"necromancy","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.finger-of-death","timing":"failed-save-full-damage-success-half","dice_count":7,"die_sides":8,"damage_type":"necrotic","bonus":30,"add_casting_modifier":false,"half_initial_only":false},{"id":"zombie","mechanic":"summon.finger-of-death-zombie","timing":"humanoid-killed-rises-at-start-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.finger-of-death","damage.necrotic","save.constitution-half","summon.finger-of-death-zombie"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","finger_of_death":{"target_sight":"creature-you-can-see-in-range","save":"constitution","damage_type":"necrotic","dice_count":7,"die_sides":8,"flat_bonus":30,"success_damage":"half-damage","raised_creature":"humanoid-killed-by-spell","rises_when":"start-of-caster-next-turn","rises_as":"zombie","zombie_reference":"monsters-section","zombie_obeys":"caster-verbal-orders"}}}`,
		`{"id":"srd-5.2.1:content/spells/fire-bolt","name":"Fire Bolt","source_id":"srd-5.2.1:content/spells/fire-bolt","source_page":132,"scope":"complete-source-record","level":0,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"choice","count":1,"eligibility":"creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.fire-bolt","timing":"ranged-spell-attack-hit","dice_count":1,"die_sides":10,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ignite","mechanic":"ignite.flammable-objects","timing":"flammable-object-hit-not-worn-or-carried","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.fire-bolt","damage.fire","attack.ranged-spell","ignite.flammable-objects","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","fire_bolt":{"attack":"ranged-spell-attack","target_kind":"creature-or-object-in-range","damage_type":"fire","ignites":"flammable-object-hit-that-is-not-worn-or-carried","upgrades":[{"character_level":5,"dice_count":2,"die_sides":10},{"character_level":11,"dice_count":3,"die_sides":10},{"character_level":17,"dice_count":4,"die_sides":10}]}}}`,
		`{"id":"srd-5.2.1:content/spells/fire-shield","name":"Fire Shield","source_id":"srd-5.2.1:content/spells/fire-shield","source_page":132,"scope":"complete-source-record","level":4,"school":"evocation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of phosphorus or a firefly","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"light","mechanic":"light.fire-shield","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"resistance","mechanic":"resistance.fire-shield","timing":"duration-per-chosen-shield","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"warm-retaliation","mechanic":"damage.fire-shield-warm-retaliation","timing":"creature-within-5-feet-hits-you-with-melee-attack-roll","dice_count":2,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"chill-retaliation","mechanic":"damage.fire-shield-chill-retaliation","timing":"creature-within-5-feet-hits-you-with-melee-attack-roll","dice_count":2,"die_sides":8,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.fire-shield","resistance.fire-shield","damage.fire-shield-warm-retaliation","damage.fire-shield-chill-retaliation","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","fire_shield":{"bright_light_radius_feet":10,"dim_light_extra_feet":10,"shield_choices":["warm-shield","chill-shield"],"warm_resistance":"cold","chill_resistance":"fire","retaliation_trigger":"creature-within-5-feet-hits-you-with-melee-attack-roll","retaliation_range_feet":5,"retaliation_dice_count":2,"retaliation_die_sides":8,"warm_damage_type":"fire","chill_damage_type":"cold"}}}`,
		`{"id":"srd-5.2.1:content/spells/fire-storm","name":"Fire Storm","source_id":"srd-5.2.1:content/spells/fire-storm","source_page":132,"scope":"complete-source-record","level":7,"school":"evocation","lists":["cleric","druid","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"cubes","count":10,"eligibility":"creatures-in-contiguous-cubes"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.fire-storm","timing":"failed-save-full-damage-success-half","dice_count":7,"die_sides":10,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ignite","mechanic":"ignite.flammable-objects","timing":"flammable-objects-in-area-not-worn-or-carried","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.fire-storm","damage.fire","area.fire-storm-cubes","save.dexterity-half","ignite.flammable-objects"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","fire_storm":{"max_cubes":10,"cube_size_feet":10,"arrangement":"as-you-like","contiguity":"each-cube-contiguous-with-at-least-one-other-cube","save":"dexterity","damage_type":"fire","dice_count":7,"die_sides":10,"success_damage":"half-damage","ignites":"flammable-objects-in-area-not-worn-or-carried"}}}`,
		`{"id":"srd-5.2.1:content/spells/fireball","name":"Fireball","source_id":"srd-5.2.1:content/spells/fireball","source_page":131,"scope":"complete-source-record","level":3,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"a ball of bat guano and sulfur","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.fireball","timing":"failed-save-full-damage-success-half","dice_count":8,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ignite","mechanic":"ignite.flammable-objects","timing":"flammable-objects-in-area-not-worn-or-carried","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.fireball","damage.fire","area.sphere-20-radius","save.dexterity-half","ignite.flammable-objects","scaling.per-slot","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","fireball_burst":{"origin":"bright-streak-from-caster-to-chosen-point-in-range","center":"point-in-range","radius_feet":20,"save":"dexterity","damage_type":"fire","dice_count":8,"die_sides":6,"success_damage":"half-damage","ignites":"flammable-objects-in-area-not-worn-or-carried","slot_above_level":3,"slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/flame-blade","name":"Flame Blade","source_id":"srd-5.2.1:content/spells/flame-blade","source_page":132,"scope":"complete-source-record","level":2,"school":"evocation","lists":["druid","sorcerer"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a sumac leaf","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"blade","mechanic":"weapon.flame-blade","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.flame-blade","timing":"melee-spell-attack-hit","dice_count":3,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":true,"half_initial_only":false},{"id":"light","mechanic":"light.flame-blade","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["weapon.flame-blade","damage.flame-blade","damage.fire","attack.melee-spell","light.flame-blade","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","flame_blade":{"blade_form":"fiery-blade-in-free-hand","blade_shape":"similar-in-size-and-shape-to-scimitar","release":"disappears-if-you-let-go","reevoke_action":"bonus-action","attack_action":"magic-action","attack":"melee-spell-attack","damage_type":"fire","dice_count":3,"die_sides":6,"damage_modifier":"spellcasting-ability-modifier","bright_light_radius_feet":10,"dim_light_extra_feet":10,"slot_above_level":2,"slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/flame-strike","name":"Flame Strike","source_id":"srd-5.2.1:content/spells/flame-strike","source_page":132,"scope":"complete-source-record","level":5,"school":"evocation","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"cylinder","count":0,"eligibility":"creatures-in-cylinder"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of sulfur","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":6,"effect_ids":["fire","radiant"]},"effects":[{"id":"fire","mechanic":"damage.flame-strike-fire","timing":"failed-save-full-damage-success-half","dice_count":5,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"radiant","mechanic":"damage.flame-strike-radiant","timing":"failed-save-full-damage-success-half","dice_count":5,"die_sides":6,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.flame-strike-fire","damage.flame-strike-radiant","area.cylinder-10-radius-40-high","save.dexterity-half","scaling.per-slot","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","flame_strike":{"origin":"vertical-column-of-fire-from-above","center":"point-in-range","radius_feet":10,"height_feet":40,"save":"dexterity","fire_dice_count":5,"fire_die_sides":6,"radiant_dice_count":5,"radiant_die_sides":6,"success_damage":"half-damage","slot_above_level":5,"slot_dice_count":1,"slot_die_sides":6,"slot_applies_to":["fire","radiant"]}}}`,
		`{"id":"srd-5.2.1:content/spells/flaming-sphere","name":"Flaming Sphere","source_id":"srd-5.2.1:content/spells/flaming-sphere","source_page":132,"scope":"complete-source-record","level":2,"school":"conjuration","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-ground-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a ball of wax","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"sphere","mechanic":"object.flaming-sphere","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.flaming-sphere","timing":"creature-ends-turn-within-5-feet-or-sphere-moves-into-space","dice_count":2,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"move","mechanic":"movement.flaming-sphere-bonus-action","timing":"bonus-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ignite","mechanic":"ignite.flammable-objects","timing":"flammable-objects-not-worn-or-carried-touched-by-sphere","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"light","mechanic":"light.flaming-sphere","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["object.flaming-sphere","damage.flaming-sphere","damage.fire","save.dexterity-half","movement.flaming-sphere-bonus-action","ignite.flammable-objects","light.flaming-sphere","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","flaming_sphere":{"diameter_feet":5,"space":"unoccupied-space-on-the-ground-in-range","damage_trigger_feet":5,"damage_trigger":"creature-ends-its-turn-within-distance-of-sphere","save":"dexterity","damage_type":"fire","dice_count":2,"die_sides":6,"success_damage":"half-damage","move_action":"bonus-action","move_feet":30,"move_mode":"rolling-along-the-ground","move_into_creature":"creature-makes-the-save-and-sphere-stops-moving-for-the-turn","barrier_height_feet":5,"pit_width_feet":10,"ignites":"flammable-objects-not-worn-or-carried-touched-by-sphere","bright_light_radius_feet":20,"dim_light_extra_feet":20,"slot_above_level":2,"slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/flesh-to-stone","name":"Flesh to Stone","source_id":"srd-5.2.1:content/spells/flesh-to-stone","source_page":133,"scope":"complete-source-record","level":6,"school":"transmutation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a cockatrice feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"restrain","mechanic":"condition.flesh-to-stone-restrained","timing":"failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"slow","mechanic":"speed.flesh-to-stone-zero","timing":"successful-save-until-start-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"petrify","mechanic":"condition.flesh-to-stone-petrified","timing":"three-failed-saves","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.flesh-to-stone-restrained","speed.flesh-to-stone-zero","condition.flesh-to-stone-petrified","save.constitution-repeating","construct.automatic-success","concentration.source-duration","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","flesh_to_stone":{"target_sight":"one-creature-you-can-see-in-range","save":"constitution","failure_condition":"restrained","success_effect":"speed-0-until-start-of-caster-next-turn","constructs_automatically_succeed":true,"repeat_save_timing":"end-of-each-restrained-target-turn","successes_to_end":3,"failures_to_petrify":3,"results_consecutive":false,"petrified_condition":"petrified-for-the-duration","full_duration_result":"petrified-until-ended-by-greater-restoration-or-similar-magic"}}}`,
		`{"id":"srd-5.2.1:content/spells/floating-disk","name":"Floating Disk","source_id":"srd-5.2.1:content/spells/floating-disk","source_page":133,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of mercury","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"disk","mechanic":"object.floating-disk","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"follow","mechanic":"movement.floating-disk-follow","timing":"caster-moves-more-than-20-feet-away","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.floating-disk-overload-or-distance","timing":"weight-over-500-pounds-or-caster-more-than-100-feet-away","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["object.floating-disk","movement.floating-disk-follow","ending.floating-disk-overload-or-distance","cast.ritual","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","floating_disk":{"plane":"circular-horizontal-plane-of-force","diameter_feet":3,"thickness_inches":1,"height_above_ground_feet":3,"space":"unoccupied-space-you-can-see-in-range","capacity_pounds":500,"overload_result":"spell-ends-and-everything-on-the-disk-falls","immobile_within_feet":20,"follow_distance_feet":20,"terrain":"uneven-terrain-stairs-and-slopes","elevation_change_limit_feet":10,"spell_ends_beyond_feet":100}}}`,
		`{"id":"srd-5.2.1:content/spells/fly","name":"Fly","source_id":"srd-5.2.1:content/spells/fly","source_page":133,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot-targets","above_level":3,"dice_count":0,"die_sides":0,"effect_ids":["fly"]},"effects":[{"id":"fly","mechanic":"speed.fly-60-hover","timing":"duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["speed.fly-60-hover","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","fly":{"target":"willing-creature-touched","fly_speed_feet":60,"hover":true,"ending":"falls-if-still-aloft-unless-it-can-stop-the-fall","base_targets":1,"additional_targets_per_slot":1}}}`,
	}
	batch := loadSpells206(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Find the Path", "Finger of Death", "Fire Bolt", "Fire Shield", "Fire Storm", "Fireball", "Flame Blade", "Flame Strike", "Flaming Sphere", "Flesh to Stone", "Floating Disk", "Fly"}
	pages := []int{131, 131, 132, 132, 132, 131, 132, 132, 132, 133, 133, 133}
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
			// Fire Bolt is a cantrip: it has no slot levels at all.
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

func TestSpells206CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells206(t)
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
	if _, err := ComposeSpellSourceCatalogue(batch, batch); err == nil {
		t.Fatal("duplicate canonical IDs accepted")
	}
	for _, r := range bounded.Records {
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
func TestSpells206MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_206.json")
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

func TestSpells206RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_206.json")
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

func TestSpells206WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_206.json")
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
		`"character_level": 5, "character_level": 6`,
		`"character_level": 5, "CHARACTER_LEVEL": 6`,
		`"character_level": 5, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"character_level": 5`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

func TestSpells206JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_206.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"radius_feet": 20`, `"radius_feet": 20, "radius_feet": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "RADIUS_FEET": 21`},
		{`"radius_feet": 20`, `"radius_feet": 20, "unknown": true`},
		{`"shield_choices": [
            "warm-shield",
            "chill-shield"
          ]`, `"shield_choices": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/find-the-path"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "fly"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 600`},
		{`"feet": 150`, `"feet": 151`},
		{`"cost_cp": 10000`, `"cost_cp": 100`},
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

func TestSpells206InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"path destination", 0, func(s *SpellRecord) { s.SourceDetails.FindPath.FailingDestinations[1] = "stationary-destination" }},
		{"path casting seconds", 0, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"path duration", 0, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"path material cost", 0, func(s *SpellRecord) { *s.Components.Material.CostCP = 1000 }},
		{"finger bonus", 1, func(s *SpellRecord) { s.SourceDetails.FingerOfDeath.FlatBonus = 20 }},
		{"finger effect dice", 1, func(s *SpellRecord) { s.Effects[0].DieSides = 6 }},
		{"finger zombie", 1, func(s *SpellRecord) { s.SourceDetails.FingerOfDeath.RisesAs = "skeleton" }},
		{"bolt upgrade", 2, func(s *SpellRecord) { s.SourceDetails.FireBolt.Upgrades[1].DiceCount = 4 }},
		{"bolt cantrip scaling", 2, func(s *SpellRecord) { s.Scaling.DieSides = 8 }},
		{"bolt slots", 2, func(s *SpellRecord) { s.SlotLevels = []int{1} }},
		{"shield resistance", 3, func(s *SpellRecord) { s.SourceDetails.FireShield.WarmResistance = "fire" }},
		{"shield retaliation type", 3, func(s *SpellRecord) { s.Effects[3].DamageType = "fire" }},
		{"shield material", 3, func(s *SpellRecord) { s.Components.Material.Description = "a bit of phosphorus" }},
		{"storm cubes", 4, func(s *SpellRecord) { s.SourceDetails.FireStorm.MaxCubes = 8 }},
		{"storm target", 4, func(s *SpellRecord) { s.Target.Count = 8 }},
		{"storm damage", 4, func(s *SpellRecord) { s.Effects[0].DiceCount = 6 }},
		{"fireball radius", 5, func(s *SpellRecord) { s.SourceDetails.FireballBurst.RadiusFeet = 30 }},
		{"fireball scaling level", 5, func(s *SpellRecord) { s.Scaling.AboveLevel = 4 }},
		{"fireball concentration", 5, func(s *SpellRecord) { *s.Concentration = true }},
		{"blade action", 6, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"blade modifier", 6, func(s *SpellRecord) { *s.Effects[1].AddCastingModifier = false }},
		{"blade attack", 6, func(s *SpellRecord) { s.SourceDetails.FlameBlade.Attack = "ranged-spell-attack" }},
		{"strike height", 7, func(s *SpellRecord) { s.SourceDetails.FlameStrike.HeightFeet = 30 }},
		{"strike radiant", 7, func(s *SpellRecord) { s.Effects[1].DiceCount = 4 }},
		{"strike scaling list", 7, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"fire"} }},
		{"sphere move", 8, func(s *SpellRecord) { s.SourceDetails.FlamingSphere.MoveFeet = 20 }},
		{"sphere barrier", 8, func(s *SpellRecord) { s.SourceDetails.FlamingSphere.BarrierHeightFeet = 10 }},
		{"sphere target", 8, func(s *SpellRecord) { s.Target.Eligibility = "creature-in-range" }},
		{"flesh saves", 9, func(s *SpellRecord) { s.SourceDetails.FleshToStone.SuccessesToEnd = 2 }},
		{"flesh constructs", 9, func(s *SpellRecord) { s.SourceDetails.FleshToStone.ConstructsAutomaticallySucceed = false }},
		{"flesh save", 9, func(s *SpellRecord) { s.SourceDetails.FleshToStone.Save = "wisdom" }},
		{"disk capacity", 10, func(s *SpellRecord) { s.SourceDetails.FloatingDisk.CapacityPounds = 400 }},
		{"disk ritual", 10, func(s *SpellRecord) { *s.Ritual = false }},
		{"disk list", 10, func(s *SpellRecord) { s.Lists = []string{"wizard", "sorcerer"} }},
		{"fly speed", 11, func(s *SpellRecord) { s.SourceDetails.Fly.FlySpeedFeet = 30 }},
		{"fly extra targets", 11, func(s *SpellRecord) { s.SourceDetails.Fly.AdditionalTargetsPerSlot = 2 }},
		{"fly scaling kind", 11, func(s *SpellRecord) { s.Scaling.Kind = "none" }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Fly = &SpellFly{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.DetectMagic = &SpellDetectMagic{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
		{"missing details", 5, func(s *SpellRecord) { s.SourceDetails = nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells206(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells206PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_206.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"flat_bonus", 1, 30}, {"max_cubes", 1, 10}, {"cube_size_feet", 1, 10}, {"radius_feet", 1, 20}, {"radius_feet", 2, 10}, {"height_feet", 1, 40},
		{"diameter_feet", 1, 5}, {"diameter_feet", 2, 3}, {"move_feet", 1, 30}, {"barrier_height_feet", 1, 5}, {"pit_width_feet", 1, 10},
		{"capacity_pounds", 1, 500}, {"immobile_within_feet", 1, 20}, {"follow_distance_feet", 1, 20}, {"elevation_change_limit_feet", 1, 10},
		{"spell_ends_beyond_feet", 1, 100}, {"fly_speed_feet", 1, 60}, {"successes_to_end", 1, 3}, {"failures_to_petrify", 1, 3},
		{"retaliation_range_feet", 1, 5}, {"base_targets", 1, 1}, {"additional_targets_per_slot", 1, 1},
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
// inventory assigns to catalogue issue #206.
func TestSpells206MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 206 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells206(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells206ChangesNoteNamesNoContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_206.json")
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
	want := "SRD 5.2.1 PDF pp.131-133 adapted into twelve complete offline spell families for issue #206. Every family sits on a single page, so no page-break continuation text is involved. No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	if strings.Contains(root.Source.Changes, "continuation text for") {
		t.Fatal("no family in this batch crosses a page break")
	}
}
