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

func loadSpells202(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_202.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.121-123 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells202SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/cure-wounds","name":"Cure Wounds","source_id":"srd-5.2.1:content/spells/cure-wounds","source_page":121,"scope":"complete-source-record","level":1,"school":"abjuration","lists":["bard","cleric","druid","paladin","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":2,"die_sides":8,"effect_ids":["healing"]},"effects":[{"id":"healing","mechanic":"healing.cure-wounds","timing":"immediate","dice_count":2,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["healing.cure-wounds","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","cure":{"target":"creature-touched","healing_dice_count":2,"healing_die_sides":8,"healing_modifier":"spellcasting-ability-modifier","slot_dice_count":2,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/dancing-lights","name":"Dancing Lights","source_id":"srd-5.2.1:content/spells/dancing-lights","source_page":121,"scope":"complete-source-record","level":0,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"created-effect","count":4,"eligibility":"torch-size-lights-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of phosphorus","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"lights","mechanic":"light.dancing-lights","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"move","mechanic":"movement.dancing-lights-bonus-action","timing":"bonus-action-each-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.dancing-lights","movement.dancing-lights-bonus-action","light.dim-10-foot-radius","concentration.source-duration","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","lights":{"source_pages":[121,122],"light_count":4,"light_size":"torch","forms":["torches","lanterns","glowing-orbs"],"combined_form":"vaguely-humanlike","combined_size":"medium","light_level":"dim-light","light_radius_feet":10,"move_action":"bonus-action","move_feet":60,"move_destination":"space-within-range","max_separation_feet":20,"separation_from":"another-light-created-by-spell","vanishes_when":"exceeds-spell-range"}}}`,
		`{"id":"srd-5.2.1:content/spells/darkness","name":"Darkness","source_id":"srd-5.2.1:content/spells/darkness","source_page":122,"scope":"complete-source-record","level":2,"school":"evocation","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"point-in-range-or-unattended-object"},"components":{"verbal":true,"somatic":false,"material":{"description":"bat fur and a piece of coal","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"darkness","mechanic":"area.magical-darkness-sphere-15","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"object-darkness","mechanic":"area.magical-darkness-emanation-15","timing":"for-duration-from-unattended-object","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dispel","mechanic":"dispel.overlapping-light-spell-level-2-or-lower","timing":"on-overlap","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.magical-darkness-sphere-15","area.magical-darkness-emanation-15","dispel.overlapping-light-spell-level-2-or-lower","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","darkness":{"darkness_kind":"magical-darkness","sphere_radius_feet":15,"sphere_origin":"point-in-range","darkvision":"cannot-see-through","nonmagical_light":"cannot-illuminate","object_alternative":"object-not-worn-or-carried","emanation_feet":15,"object_blocker":"covered-by-opaque-item","blocker_examples":["bowl","helm"],"dispels_lights":["bright-light","dim-light"],"dispels_up_to_level":2}}}`,
		`{"id":"srd-5.2.1:content/spells/darkvision","name":"Darkvision","source_id":"srd-5.2.1:content/spells/darkvision","source_page":122,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a dried carrot","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"darkvision","mechanic":"sense.darkvision-150","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.darkvision-150","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","darkvision":{"target":"willing-creature-touched","sense":"darkvision","sense_range_feet":150}}}`,
		`{"id":"srd-5.2.1:content/spells/daylight","name":"Daylight","source_id":"srd-5.2.1:content/spells/daylight","source_page":122,"scope":"complete-source-record","level":3,"school":"evocation","lists":["cleric","druid","paladin","ranger","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"point-in-range-or-unattended-object"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"daylight","mechanic":"area.sunlight-sphere-60","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"object-daylight","mechanic":"area.sunlight-emanation-60","timing":"for-duration-from-unattended-object","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dispel","mechanic":"dispel.overlapping-darkness-spell-level-3-or-lower","timing":"on-overlap","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.sunlight-sphere-60","area.sunlight-emanation-60","dispel.overlapping-darkness-spell-level-3-or-lower"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","daylight":{"light_kind":"sunlight","sphere_radius_feet":60,"sphere_origin":"point-in-range","sphere_light":"bright-light","dim_light_extra_feet":60,"object_alternative":"object-not-worn-or-carried","emanation_feet":60,"object_blocker":"covered-by-opaque-item","blocker_examples":["bowl","helm"],"dispels_darkness":"darkness-from-spell","dispels_up_to_level":3}}}`,
		`{"id":"srd-5.2.1:content/spells/death-ward","name":"Death Ward","source_id":"srd-5.2.1:content/spells/death-ward","source_page":122,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"drop-to-one","mechanic":"protection.death-ward-drop-to-1","timing":"first-time-target-would-drop-to-0-hit-points","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"instant-death","mechanic":"protection.death-ward-negate-instant-death","timing":"effect-that-kills-instantly-without-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["protection.death-ward-drop-to-1","protection.death-ward-negate-instant-death"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","death_ward":{"target":"creature-touched","drop_trigger":"would-drop-to-0-hit-points-first-time","replacement_hit_points":1,"drop_ends_spell":true,"instant_death_trigger":"effect-that-would-kill-instantly-without-dealing-damage","instant_death_result":"effect-negated","instant_death_ends_spell":true}}}`,
		`{"id":"srd-5.2.1:content/spells/delayed-blast-fireball","name":"Delayed Blast Fireball","source_id":"srd-5.2.1:content/spells/delayed-blast-fireball","source_page":122,"scope":"complete-source-record","level":7,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"a ball of bat guano and sulfur","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":7,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.delayed-blast-fireball","timing":"spell-ends-bead-explodes-save-half-on-success","dice_count":12,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"growth","mechanic":"damage.delayed-blast-fireball-growth","timing":"whenever-caster-turn-ends-and-spell-has-not-ended","dice_count":1,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"touch","mechanic":"interaction.delayed-blast-bead-touch","timing":"creature-touches-bead-before-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ignite","mechanic":"ignite.flammable-objects","timing":"when-bead-explodes","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.delayed-blast-fireball","damage.delayed-blast-fireball-growth","interaction.delayed-blast-bead-touch","ignite.flammable-objects","area.sphere-20-radius","save.dexterity-half","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","fireball":{"bead_origin":"beam-of-yellow-light-from-caster","bead_point":"chosen-point-in-range","explosion_trigger":"spell-ends","radius_feet":20,"save":"dexterity","damage_type":"fire","fail_damage":"total-accumulated-damage","success_damage":"half-accumulated-damage","base_dice_count":12,"base_die_sides":6,"growth_dice_count":1,"growth_die_sides":6,"growth_timing":"whenever-caster-turn-ends-and-spell-has-not-ended","touch_save":"dexterity","touch_failure":"spell-ends-and-bead-explodes","touch_success":"creature-can-throw-bead","throw_feet":40,"thrown_end":"enters-creature-space-or-collides-with-solid-object","thrown_result":"spell-ends-and-bead-explodes","ignites":"flammable-objects-not-worn-or-carried-in-explosion","slot_base_dice_count":1,"slot_base_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/demiplane","name":"Demiplane","source_id":"srd-5.2.1:content/spells/demiplane","source_page":122,"scope":"complete-source-record","level":8,"school":"conjuration","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-flat-solid-surface-in-range"},"components":{"verbal":false,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"door","mechanic":"portal.demiplane-door","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"demiplane","mechanic":"space.demiplane-room","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ending","mechanic":"ending.demiplane-door-vanishes","timing":"when-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["portal.demiplane-door","space.demiplane-room","ending.demiplane-door-vanishes","condition.prone-on-shunt"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","demiplane":{"source_pages":[122,123],"surface":"flat-solid-surface-visible-in-range","door_size":"medium","door_appearance":"shadowy","door_operation":"opened-and-closed","room_kind":"empty-room","room_dimension_feet":30,"room_material_choices":["wood","stone"],"spell_end_door":"vanishes","objects_inside":"remain","creatures_inside":"remain-unless-they-opt-to-be-shunted-through-door","shunted_condition":"prone","shunted_landing":"unoccupied-spaces-closest-to-former-door-space","connection_choices":["new-demiplane","demiplane-from-previous-casting","demiplane-by-another-creature"],"other_creature_needs":"know-nature-and-contents"}}}`,
		`{"id":"srd-5.2.1:content/spells/detect-evil-and-good","name":"Detect Evil and Good","source_id":"srd-5.2.1:content/spells/detect-evil-and-good","source_page":123,"scope":"complete-source-record","level":1,"school":"divination","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sense","mechanic":"sense.detect-evil-and-good","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.detect-evil-and-good","sense.hallow-active","blocker.detection-barriers","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","detect_evil":{"sense_radius_feet":30,"detected_types":["aberration","celestial","elemental","fey","fiend","undead"],"sensed_info":"location","hallow_detection":"whether-active-and-where","blockers":[{"material":"stone-dirt-or-wood","thickness":"1-foot"},{"material":"metal","thickness":"1-inch"},{"material":"lead","thickness":"thin-sheet"}]}}}`,
		`{"id":"srd-5.2.1:content/spells/detect-magic","name":"Detect Magic","source_id":"srd-5.2.1:content/spells/detect-magic","source_page":123,"scope":"complete-source-record","level":1,"school":"divination","lists":["bard","cleric","druid","paladin","ranger","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sense","mechanic":"sense.detect-magic","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"aura","mechanic":"sense.detect-magic-aura-school","timing":"magic-action-when-effects-sensed","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.detect-magic","sense.detect-magic-aura-school","blocker.detection-barriers","cast.ritual","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","detect_magic":{"sense_radius_feet":30,"senses":"presence-of-magical-effects","follow_up_action":"magic-action","aura_subjects":"visible-creature-or-object-bearing-the-magic","aura_appearance":"faint-aura","school_learned":"school-of-magic-if-effect-created-by-spell","blockers":[{"material":"stone-dirt-or-wood","thickness":"1-foot"},{"material":"metal","thickness":"1-inch"},{"material":"lead","thickness":"thin-sheet"}]}}}`,
		`{"id":"srd-5.2.1:content/spells/detect-poison-and-disease","name":"Detect Poison and Disease","source_id":"srd-5.2.1:content/spells/detect-poison-and-disease","source_page":123,"scope":"complete-source-record","level":1,"school":"divination","lists":["cleric","druid","paladin","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a yew leaf","cost_cp":0,"consumed":false}},"ritual":true,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sense","mechanic":"sense.detect-poison-and-disease","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.detect-poison-and-disease","blocker.detection-barriers","cast.ritual","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","detect_poison":{"sense_radius_feet":30,"detected":["poisons","poisonous-or-venomous-creatures","magical-contagions"],"sensed_info":"location","kind_sensed":"kind-of-poison-creature-or-contagion","blockers":[{"material":"stone-dirt-or-wood","thickness":"1-foot"},{"material":"metal","thickness":"1-inch"},{"material":"lead","thickness":"thin-sheet"}]}}}`,
		`{"id":"srd-5.2.1:content/spells/detect-thoughts","name":"Detect Thoughts","source_id":"srd-5.2.1:content/spells/detect-thoughts","source_page":123,"scope":"complete-source-record","level":2,"school":"divination","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"1 Copper Piece","cost_cp":1,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sense","mechanic":"sense.thoughts","timing":"sense-thoughts-option-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"read","mechanic":"read.thoughts","timing":"read-thoughts-option-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"probe","mechanic":"probe.thoughts-wisdom-save","timing":"magic-action-on-next-turn-after-reading","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sense.thoughts","read.thoughts","probe.thoughts-wisdom-save","check.arcana-ends-probe","blocker.detection-barriers","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","detect_thoughts":{"activation_action":"magic-action-on-later-turns","options":["sense-thoughts","read-thoughts"],"sense_radius_feet":30,"sense_subjects":"thoughts-of-creatures-that-know-languages-or-are-telepathic","sense_reads_content":false,"sense_blockers":[{"material":"stone-dirt-or-wood","thickness":"1-foot"},{"material":"metal","thickness":"1-inch"},{"material":"lead","thickness":"thin-sheet"}],"read_range_feet":30,"read_targets":["visible-creature","creature-detected-with-sense-thoughts"],"read_learns":"what-is-most-on-targets-mind-right-now","read_no_language_result":"learn-nothing","probe_action":"magic-action-on-next-turn","probe_save":"wisdom","probe_failure":["reasoning","emotions","something-that-looms-large"],"probe_failure_examples":["worry","love","hate"],"probe_success":"spell-ends","target_knows_probed":true,"target_check_action":"action-on-its-turn","target_check":"intelligence-arcana","target_check_dc":"caster-spell-save-dc","target_check_success":"spell-ends","target_check_window":"until-caster-shifts-attention-from-target-mind"}}}`,
	}
	batch := loadSpells202(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Cure Wounds", "Dancing Lights", "Darkness", "Darkvision", "Daylight", "Death Ward", "Delayed Blast Fireball", "Demiplane", "Detect Evil and Good", "Detect Magic", "Detect Poison and Disease", "Detect Thoughts"}
	pages := []int{121, 121, 122, 122, 122, 122, 122, 122, 123, 123, 123, 123}
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
			// Dancing Lights is a cantrip: it has no slot levels at all.
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

func TestSpells202CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells202(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 84 {
		t.Fatalf("expected seven twelve-family batches, got %d", len(all.Records))
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
func TestSpells202MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_202.json")
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

func TestSpells202JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_202.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"sphere_radius_feet": 15`, `"sphere_radius_feet": 15, "sphere_radius_feet": 16`},
		{`"sphere_radius_feet": 15`, `"sphere_radius_feet": 15, "SPHERE_RADIUS_FEET": 16`},
		{`"sphere_radius_feet": 15`, `"sphere_radius_feet": 15, "unknown": true`},
		{`"room_material_choices": [
            "wood",
            "stone"
          ]`, `"room_material_choices": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/cure-wounds"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "healing"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "action"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 0`, `"casting_seconds": 60`},
		{`"feet": 150`, `"feet": 151`},
		{`"cost_cp": 1`, `"cost_cp": 2`},
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

func TestSpells202InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"cure die", 0, func(s *SpellRecord) { s.SourceDetails.Cure.HealingDieSides = 6 }},
		{"cure scaling row", 0, func(s *SpellRecord) { s.Scaling.DiceCount = 1 }},
		{"cure slot levels", 0, func(s *SpellRecord) { s.SlotLevels = []int{2, 3} }},
		{"lights count", 1, func(s *SpellRecord) { s.SourceDetails.Lights.LightCount = 5 }},
		{"lights continuation page", 1, func(s *SpellRecord) { s.SourceDetails.Lights.SourcePages = []int{121} }},
		{"lights cantrip slots", 1, func(s *SpellRecord) { s.SlotLevels = []int{0} }},
		{"darkness sphere", 2, func(s *SpellRecord) { s.SourceDetails.Darkness.SphereRadiusFeet = 20 }},
		{"darkness dispel level", 2, func(s *SpellRecord) { s.SourceDetails.Darkness.DispelsUpToLevel = 3 }},
		{"darkness duration", 2, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"darkvision range", 3, func(s *SpellRecord) { s.SourceDetails.Darkvision.SenseRangeFeet = 60 }},
		{"daylight dispel level", 4, func(s *SpellRecord) { s.SourceDetails.Daylight.DispelsUpToLevel = 2 }},
		{"daylight concentration", 4, func(s *SpellRecord) { *s.Concentration = true }},
		{"death ward hit points", 5, func(s *SpellRecord) { s.SourceDetails.DeathWard.ReplacementHitPoints = 0 }},
		{"death ward instant death", 5, func(s *SpellRecord) { s.SourceDetails.DeathWard.InstantDeathEndsSpell = false }},
		{"fireball base damage", 6, func(s *SpellRecord) { s.SourceDetails.Fireball.BaseDiceCount = 8 }},
		{"fireball effect dice", 6, func(s *SpellRecord) { s.Effects[0].DiceCount = 8 }},
		{"fireball throw", 6, func(s *SpellRecord) { s.SourceDetails.Fireball.ThrowFeet = 30 }},
		{"fireball scaling level", 6, func(s *SpellRecord) { s.Scaling.AboveLevel = 6 }},
		{"demiplane room", 7, func(s *SpellRecord) { s.SourceDetails.Demiplane.RoomDimensionFeet = 20 }},
		{"demiplane material choice", 7, func(s *SpellRecord) { s.SourceDetails.Demiplane.RoomMaterialChoices[1] = "metal" }},
		{"demiplane continuation page", 7, func(s *SpellRecord) { s.SourceDetails.Demiplane.SourcePages = []int{122} }},
		{"demiplane components", 7, func(s *SpellRecord) { *s.Components.Verbal = true }},
		{"evil type", 8, func(s *SpellRecord) { s.SourceDetails.DetectEvil.DetectedTypes[5] = "humanoid" }},
		{"evil blocker", 8, func(s *SpellRecord) { s.SourceDetails.DetectEvil.Blockers[1].Thickness = "1-foot" }},
		{"magic ritual", 9, func(s *SpellRecord) { *s.Ritual = false }},
		{"magic list", 9, func(s *SpellRecord) { s.Lists = s.Lists[1:] }},
		{"poison detected", 10, func(s *SpellRecord) { s.SourceDetails.DetectPoison.Detected[2] = "diseases" }},
		{"poison material", 10, func(s *SpellRecord) { s.Components.Material.Description = "a leaf" }},
		{"thoughts save", 11, func(s *SpellRecord) { s.SourceDetails.DetectThoughts.ProbeSave = "intelligence" }},
		{"thoughts option", 11, func(s *SpellRecord) { s.SourceDetails.DetectThoughts.Options[1] = "unknown" }},
		{"thoughts material cost", 11, func(s *SpellRecord) { *s.Components.Material.CostCP = 100 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.DetectMagic = &SpellDetectMagic{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells202(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells202RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_202.json")
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

func TestSpells202WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_202.json")
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
		`"thickness": "1-foot", "thickness": "1-inch"`,
		`"thickness": "1-foot", "THICKNESS": "1-inch"`,
		`"thickness": "1-foot", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"thickness": "1-foot"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells202PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_202.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"light_count", 1, 4}, {"light_radius_feet", 1, 10}, {"move_feet", 1, 60}, {"max_separation_feet", 1, 20},
		{"sphere_radius_feet", 1, 15}, {"sphere_radius_feet", 2, 60}, {"emanation_feet", 1, 15}, {"emanation_feet", 2, 60},
		{"dispels_up_to_level", 1, 2}, {"dispels_up_to_level", 2, 3}, {"sense_range_feet", 1, 150},
		{"replacement_hit_points", 1, 1}, {"radius_feet", 1, 20}, {"base_dice_count", 1, 12}, {"throw_feet", 1, 40},
		{"room_dimension_feet", 1, 30}, {"sense_radius_feet", 1, 30}, {"sense_radius_feet", 2, 30},
		{"sense_radius_feet", 3, 30}, {"sense_radius_feet", 4, 30}, {"read_range_feet", 1, 30},
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
// inventory assigns to catalogue issue #202.
func TestSpells202MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 202 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells202(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}
