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

func loadSpells209(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_209.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.139-141 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells209SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/hallucinatory-terrain","name":"Hallucinatory Terrain","source_id":"srd-5.2.1:content/spells/hallucinatory-terrain","source_page":139,"scope":"complete-source-record","level":4,"school":"illusion","lists":["bard","druid","warlock","wizard"],"casting_time":"10-minutes","range":{"kind":"feet","feet":300},"target":{"shape":"cube","count":0,"eligibility":"natural-terrain-in-cube"},"components":{"verbal":true,"somatic":true,"material":{"description":"a mushroom","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"illusion","mechanic":"illusion.hallucinatory-terrain","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.hallucinatory-terrain","check.hallucinatory-terrain-investigation","cast.ten-minutes","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-ended","hallucinatory_terrain":{"cube_feet":150,"affected_terrain":"natural-terrain","disguised_senses":["sight","sound","smell"],"disguise":"another-sort-of-natural-terrain","examples":["open-fields-or-road-as-swamp-hill-crevasse-or-other-difficult-or-impassable-terrain","pond-as-grassy-meadow","precipice-as-gentle-slope","rock-strewn-gully-as-wide-smooth-road"],"unchanged":["manufactured-structures","equipment","creatures"],"tactile_characteristics":"unchanged","touch_result":"creatures-entering-likely-to-notice-illusion","disbelieve_condition":"difference-not-obvious-by-touch","disbelieve_action":"study","disbelieve_check":"intelligence-investigation","disbelieve_dc":"caster-spell-save-dc","discerned":"sees-vague-image-superimposed-on-real-terrain"}}}`,
		`{"id":"srd-5.2.1:content/spells/harm","name":"Harm","source_id":"srd-5.2.1:content/spells/harm","source_page":139,"scope":"complete-source-record","level":6,"school":"necromancy","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.harm","timing":"on-failed-constitution-save","dice_count":14,"die_sides":6,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.harm","save.constitution-half","reduction.harm-hit-point-maximum"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","harm":{"target":"creature-visible-in-range","save":"constitution","damage_dice_count":14,"damage_die_sides":6,"damage_type":"necrotic","failure_hit_point_maximum":"reduced-by-necrotic-damage-taken","success":"half-damage-only","hit_point_maximum_floor":1}}}`,
		`{"id":"srd-5.2.1:content/spells/haste","name":"Haste","source_id":"srd-5.2.1:content/spells/haste","source_page":139,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-willing-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a shaving of licorice root","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"speed","mechanic":"speed.haste-doubled","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"armor-class","mechanic":"defense.haste-armor-class","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":2,"add_casting_modifier":false,"half_initial_only":false},{"id":"saves","mechanic":"advantage.haste-dexterity-saves","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"action","mechanic":"action.haste-additional","timing":"each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"lethargy","mechanic":"condition.haste-lethargy","timing":"when-spell-ends-until-end-of-targets-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["speed.haste-doubled","defense.haste-armor-class","advantage.haste-dexterity-saves","action.haste-additional","condition.haste-lethargy","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","haste":{"target":"willing-creature-visible-in-range","speed_multiplier":2,"armor_class_bonus":2,"advantage":"dexterity-saving-throws","additional_actions_per_turn":1,"allowed_actions":["attack-one-attack-only","dash","disengage","hide","utilize"],"end_condition":"incapacitated","end_speed_feet":0,"end_duration":"until-end-of-targets-next-turn","end_description":"wave-of-lethargy"}}}`,
		`{"id":"srd-5.2.1:content/spells/heal","name":"Heal","source_id":"srd-5.2.1:content/spells/heal","source_page":139,"scope":"complete-source-record","level":6,"school":"abjuration","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-bonus","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["healing"]},"effects":[{"id":"healing","mechanic":"healing.heal","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":70,"add_casting_modifier":false,"half_initial_only":false},{"id":"conditions","mechanic":"remove.heal-conditions","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.heal","remove.heal-conditions","scaling.bonus-per-slot"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","heal":{"target":"creature-visible-in-range","healing_hit_points":70,"conditions_ended":["blinded","deafened","poisoned"],"additional_hit_points_per_slot":10}}}`,
		`{"id":"srd-5.2.1:content/spells/healing-word","name":"Healing Word","source_id":"srd-5.2.1:content/spells/healing-word","source_page":139,"scope":"complete-source-record","level":1,"school":"abjuration","lists":["bard","cleric","druid"],"casting_time":"bonus-action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":2,"die_sides":4,"effect_ids":["healing"]},"effects":[{"id":"healing","mechanic":"healing.healing-word","timing":"immediate","dice_count":2,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["healing.healing-word","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","healing_word":{"target":"creature-of-choice-visible-in-range","healing_dice_count":2,"healing_die_sides":4,"healing_adds":"spellcasting-ability-modifier","slot_dice_count":2,"slot_die_sides":4}}}`,
		`{"id":"srd-5.2.1:content/spells/heat-metal","name":"Heat Metal","source_id":"srd-5.2.1:content/spells/heat-metal","source_page":139,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["bard","druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-manufactured-metal-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of iron and a flame","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.heat-metal","timing":"on-cast-and-bonus-action-later-turns","dice_count":2,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"grip","mechanic":"save.heat-metal-drop-object","timing":"when-holder-or-wearer-takes-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.heat-metal","save.heat-metal-drop-object","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","heat_metal":{"source_pages":[139,140],"object":"manufactured-metal-object","object_examples":["metal-weapon","suit-of-heavy-or-medium-metal-armor"],"object_effect":"glows-red-hot","damaged":"any-creature-in-physical-contact","damage_dice_count":2,"damage_die_sides":8,"damage_type":"fire","initial_timing":"when-spell-cast","repeat_action":"bonus-action-on-later-turns","repeat_requirement":"object-within-range","save_trigger":"holding-or-wearing-object-and-takes-its-damage","save":"constitution","save_failure":"drop-object-if-able","kept_penalty":["disadvantage-on-attack-rolls","disadvantage-on-ability-checks"],"kept_penalty_duration":"until-start-of-casters-next-turn","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/hellish-rebuke","name":"Hellish Rebuke","source_id":"srd-5.2.1:content/spells/hellish-rebuke","source_page":140,"scope":"complete-source-record","level":1,"school":"evocation","lists":["warlock"],"casting_time":"reaction","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-that-damaged-caster"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.hellish-rebuke","timing":"immediate","dice_count":2,"die_sides":10,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.hellish-rebuke","reaction.hellish-rebuke-trigger","save.dexterity-half","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["taking-damage-from-visible-creature-within-60-feet"],"end_timing":"instantaneous","hellish_rebuke":{"target":"creature-that-damaged-caster","trigger_distance_feet":60,"trigger_requirement":"damaged-by-creature-caster-can-see","appearance":"momentarily-surrounded-by-green-flames","save":"dexterity","damage_dice_count":2,"damage_die_sides":10,"damage_type":"fire","success":"half-damage","slot_dice_count":1,"slot_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/heroes-feast","name":"Heroes’ Feast","source_id":"srd-5.2.1:content/spells/heroes-feast","source_page":140,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["bard","cleric","druid"],"casting_time":"10-minutes","range":{"kind":"self","feet":0},"target":{"shape":"created-effect","count":1,"eligibility":"surface-in-unoccupied-cube-next-to-caster"},"components":{"verbal":true,"somatic":true,"material":{"description":"a gem-encrusted bowl worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"feast","mechanic":"create.heroes-feast","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"resistance","mechanic":"resistance.heroes-feast-poison","timing":"for-24-hours-after-feast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"immunity","mechanic":"immunity.heroes-feast-frightened-poisoned","timing":"for-24-hours-after-feast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"hit-points","mechanic":"increase.heroes-feast-hit-point-maximum","timing":"for-24-hours-after-feast","dice_count":2,"die_sides":10,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["create.heroes-feast","resistance.heroes-feast-poison","immunity.heroes-feast-frightened-poisoned","increase.heroes-feast-hit-point-maximum","cast.ten-minutes","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"instantaneous","heroes_feast":{"location":"surface-in-unoccupied-cube-next-to-caster","cube_feet":10,"consume_hours":1,"feast_disappears":"at-end-of-consume-time","benefits_begin":"after-consume-hour-is-over","max_creatures":12,"benefit_hours":24,"damage_resistance":"poison","condition_immunities":["frightened","poisoned"],"hit_point_maximum_dice_count":2,"hit_point_maximum_die_sides":10,"hit_points_gained":"same-number-as-maximum-increase"}}}`,
		`{"id":"srd-5.2.1:content/spells/heroism","name":"Heroism","source_id":"srd-5.2.1:content/spells/heroism","source_page":140,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","paladin"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["immunity","temporary-hit-points"]},"effects":[{"id":"immunity","mechanic":"immunity.heroism-frightened","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"temporary-hit-points","mechanic":"temporary-hit-points.heroism","timing":"start-of-each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["immunity.heroism-frightened","temporary-hit-points.heroism","scaling.targets-per-slot","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","heroism":{"target":"willing-creature-touched","condition_immunity":"frightened","temporary_hit_points":"spellcasting-ability-modifier","temporary_hit_points_timing":"start-of-each-of-its-turns","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/hex","name":"Hex","source_id":"srd-5.2.1:content/spells/hex","source_page":140,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["warlock"],"casting_time":"bonus-action","range":{"kind":"feet","feet":90},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"the petrified eye of a newt","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-duration","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["damage","checks"]},"effects":[{"id":"damage","mechanic":"damage.hex-extra","timing":"whenever-caster-hits-target-with-attack-roll","dice_count":1,"die_sides":6,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"checks","mechanic":"disadvantage.hex-ability-checks","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.hex-extra","disadvantage.hex-ability-checks","retarget.hex-bonus-action","scaling.duration-by-slot","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","hex":{"target":"creature-visible-in-range","extra_damage_dice_count":1,"extra_damage_die_sides":6,"damage_type":"necrotic","damage_trigger":"whenever-caster-hits-target-with-attack-roll","ability_choice_timing":"when-spell-cast","ability_choice_count":1,"check_penalty":"disadvantage-on-ability-checks-with-chosen-ability","retarget_condition":"target-drops-to-0-hit-points-before-spell-ends","retarget_action":"bonus-action-on-later-turn","retarget_result":"curse-new-creature","slot_durations":[{"min_slot":2,"max_slot":2,"seconds":14400},{"min_slot":3,"max_slot":4,"seconds":28800},{"min_slot":5,"max_slot":9,"seconds":86400}]}}}`,
		`{"id":"srd-5.2.1:content/spells/hideous-laughter","name":"Hideous Laughter","source_id":"srd-5.2.1:content/spells/hideous-laughter","source_page":140,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a tart and a feather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["conditions","repeat"]},"effects":[{"id":"conditions","mechanic":"condition.hideous-laughter-prone-incapacitated","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"save.hideous-laughter-repeat","timing":"end-of-each-turn-and-each-time-damaged","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.hideous-laughter-prone-incapacitated","save.hideous-laughter-repeat","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-successful-save","hideous_laughter":{"target":"one-creature-of-choice-visible-in-range","save":"wisdom","failure_conditions":["prone","incapacitated"],"laughter":"laughs-uncontrollably-if-capable-of-laughter","prone_restriction":"cannot-end-prone-condition-on-itself","repeat_triggers":["end-of-each-of-its-turns","each-time-it-takes-damage"],"damage_save":"advantage-if-triggered-by-damage","repeat_success":"spell-ends","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/hold-monster","name":"Hold Monster","source_id":"srd-5.2.1:content/spells/hold-monster","source_page":140,"scope":"complete-source-record","level":5,"school":"enchantment","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a straight piece of iron","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["paralyzed","repeat"]},"effects":[{"id":"paralyzed","mechanic":"condition.hold-monster-paralyzed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"save.hold-monster-repeat","timing":"end-of-each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.hold-monster-paralyzed","save.hold-monster-repeat","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","hold_monster":{"source_pages":[140,141],"target":"creature-visible-in-range","save":"wisdom","failure_condition":"paralyzed","repeat_timing":"end-of-each-of-its-turns","repeat_success":"spell-ends-on-target","base_targets":1,"additional_targets_per_slot":1}}}`,
	}
	batch := loadSpells209(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Hallucinatory Terrain", "Harm", "Haste", "Heal", "Healing Word", "Heat Metal", "Hellish Rebuke", "Heroes’ Feast", "Heroism", "Hex", "Hideous Laughter", "Hold Monster"}
	pages := []int{139, 139, 139, 139, 139, 139, 140, 140, 140, 140, 140, 140}
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

func TestSpells209CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells209(t)
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
func TestSpells209MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_209.json")
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

func TestSpells209JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_209.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"cube_feet": 150`, `"cube_feet": 150, "cube_feet": 151`},
		{`"cube_feet": 150`, `"cube_feet": 150, "CUBE_FEET": 151`},
		{`"cube_feet": 150`, `"cube_feet": 150, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/hallucinatory-terrain"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "healing"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "10-minutes"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 600`, `"casting_seconds": 3600`},
		{`"feet": 300`, `"feet": 301`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
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

func TestSpells209InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"terrain size", 0, func(s *SpellRecord) { s.SourceDetails.HallucinatoryTerrain.CubeFeet = 100 }},
		{"terrain casting", 0, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 60 }},
		{"harm dice", 1, func(s *SpellRecord) { s.Effects[0].DiceCount = 12 }},
		{"harm floor", 1, func(s *SpellRecord) { s.SourceDetails.Harm.HitPointMaximumFloor = 0 }},
		{"haste actions", 2, func(s *SpellRecord) { s.SourceDetails.Haste.AllowedActions[4] = "magic" }},
		{"haste concentration", 2, func(s *SpellRecord) { *s.Concentration = false }},
		{"heal amount", 3, func(s *SpellRecord) { s.Effects[0].Bonus = 80 }},
		{"heal slot bonus", 3, func(s *SpellRecord) { s.SourceDetails.Heal.AdditionalHitPointsPerSlot = 0 }},
		{"healing word casting", 4, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"healing word slot die", 4, func(s *SpellRecord) { s.Scaling.DieSides = 8 }},
		{"healing word somatic", 4, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"heat metal pages", 5, func(s *SpellRecord) { s.SourceDetails.HeatMetal.SourcePages = []int{139} }},
		{"heat metal penalty", 5, func(s *SpellRecord) {
			s.SourceDetails.HeatMetal.KeptPenalty = s.SourceDetails.HeatMetal.KeptPenalty[:1]
		}},
		{"rebuke trigger", 6, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{} }},
		{"rebuke save", 6, func(s *SpellRecord) { s.SourceDetails.HellishRebuke.Save = "constitution" }},
		{"feast cost", 7, func(s *SpellRecord) { *s.Components.Material.CostCP = 10000 }},
		{"feast consumption", 7, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"feast creatures", 7, func(s *SpellRecord) { s.SourceDetails.HeroesFeast.MaxCreatures = 13 }},
		{"heroism targets", 8, func(s *SpellRecord) { s.SourceDetails.Heroism.AdditionalTargetsPerSlot = 2 }},
		{"heroism modifier", 8, func(s *SpellRecord) { *s.Effects[1].AddCastingModifier = false }},
		{"hex duration table", 9, func(s *SpellRecord) { s.SourceDetails.Hex.SlotDurations[1].MaxSlot = 5 }},
		{"hex missing row", 9, func(s *SpellRecord) { s.SourceDetails.Hex.SlotDurations = s.SourceDetails.Hex.SlotDurations[:2] }},
		{"laughter conditions", 10, func(s *SpellRecord) { s.SourceDetails.HideousLaughter.FailureConditions[0] = "stunned" }},
		{"hold monster eligibility", 11, func(s *SpellRecord) { s.Target.Eligibility = "visible-humanoid-in-range" }},
		{"hold monster slot base", 11, func(s *SpellRecord) { s.Scaling.AboveLevel = 4 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.HoldMonster = &SpellHoldMonster{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells209(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells209RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_209.json")
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

func TestSpells209WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_209.json")
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
		`"cube_feet": 150, "cube_feet": 151`,
		`"cube_feet": 150, "CUBE_FEET": 151`,
		`"cube_feet": 150, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"cube_feet": 150`), []byte(replacement), 1)
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
func TestSpells209PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_209.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"cube_feet", 1, 150},
		{"cube_feet", 2, 10},
		{"damage_dice_count", 1, 14},
		{"damage_dice_count", 2, 2},
		{"damage_dice_count", 3, 2},
		{"damage_die_sides", 1, 6},
		{"damage_die_sides", 2, 8},
		{"damage_die_sides", 3, 10},
		{"hit_point_maximum_floor", 1, 1},
		{"speed_multiplier", 1, 2},
		{"armor_class_bonus", 1, 2},
		{"additional_actions_per_turn", 1, 1},
		{"end_speed_feet", 1, 0},
		{"healing_hit_points", 1, 70},
		{"additional_hit_points_per_slot", 1, 10},
		{"healing_dice_count", 1, 2},
		{"healing_die_sides", 1, 4},
		{"slot_dice_count", 1, 2},
		{"slot_dice_count", 2, 1},
		{"slot_dice_count", 3, 1},
		{"slot_die_sides", 1, 4},
		{"slot_die_sides", 2, 8},
		{"slot_die_sides", 3, 10},
		{"trigger_distance_feet", 1, 60},
		{"consume_hours", 1, 1},
		{"max_creatures", 1, 12},
		{"benefit_hours", 1, 24},
		{"hit_point_maximum_dice_count", 1, 2},
		{"hit_point_maximum_die_sides", 1, 10},
		{"base_targets", 1, 1},
		{"base_targets", 2, 1},
		{"base_targets", 3, 1},
		{"additional_targets_per_slot", 1, 1},
		{"additional_targets_per_slot", 2, 1},
		{"additional_targets_per_slot", 3, 1},
		{"extra_damage_dice_count", 1, 1},
		{"extra_damage_die_sides", 1, 6},
		{"ability_choice_count", 1, 1},
		{"min_slot", 1, 2},
		{"min_slot", 2, 3},
		{"min_slot", 3, 5},
		{"max_slot", 1, 2},
		{"max_slot", 2, 4},
		{"max_slot", 3, 9},
		{"seconds", 1, 14400},
		{"seconds", 2, 28800},
		{"seconds", 3, 86400},
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
// inventory assigns to catalogue issue #209.
func TestSpells209MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 209 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells209(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells209ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_209.json")
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
	want := "SRD 5.2.1 PDF pp.139-141 adapted into twelve complete offline spell families for issue #209, including page-break continuation text for Heat Metal (pp.139-140) and Hold Monster (pp.140-141). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells209PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.HoldMonster = loadSpells209(t).Records[11].SourceDetails.HoldMonster
	if r.Validate() == nil {
		t.Fatal("batch 209 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells209RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells209(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell209Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell209Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
