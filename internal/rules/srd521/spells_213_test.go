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

func loadSpells213(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_213.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.147-150 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells213SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/mass-heal","name":"Mass Heal","source_id":"srd-5.2.1:content/spells/mass-heal","source_page":147,"scope":"complete-source-record","level":9,"school":"abjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"selected-creatures","count":0,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"healing","mechanic":"healing.mass-heal-pool","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"conditions","mechanic":"remove.mass-heal-conditions","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.mass-heal-pool","remove.mass-heal-conditions"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","mass_heal":{"hit_point_pool":700,"division":"divided-as-caster-chooses","targets":"any-number-of-creatures-visible-in-range","conditions_removed":["blinded","deafened","poisoned"],"condition_targets":"creatures-healed-by-spell"}}}`,
		`{"id":"srd-5.2.1:content/spells/mass-healing-word","name":"Mass Healing Word","source_id":"srd-5.2.1:content/spells/mass-healing-word","source_page":148,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["bard","cleric"],"casting_time":"bonus-action","range":{"kind":"feet","feet":60},"target":{"shape":"up-to-count","count":6,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":4,"effect_ids":["healing"]},"effects":[{"id":"healing","mechanic":"healing.mass-healing-word","timing":"immediate","dice_count":2,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["healing.mass-healing-word","scaling.per-slot"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","mass_healing_word":{"max_targets":6,"targets":"creatures-of-casters-choice-visible-in-range","healing_dice_count":2,"healing_die_sides":4,"healing_adds":"spellcasting-ability-modifier","slot_dice_count":1,"slot_die_sides":4}}}`,
		`{"id":"srd-5.2.1:content/spells/mass-suggestion","name":"Mass Suggestion","source_id":"srd-5.2.1:content/spells/mass-suggestion","source_page":148,"scope":"complete-source-record","level":6,"school":"enchantment","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"up-to-count","count":12,"eligibility":"visible-hearing-understanding-creatures-in-range"},"components":{"verbal":true,"somatic":false,"material":{"description":"a snake’s tongue","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"slot-duration","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.mass-suggestion-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.mass-suggestion-charmed","ending.mass-suggestion-damage-or-completion","scaling.duration-by-slot","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-damage-or-activity-completed","mass_suggestion":{"suggestion":"course-of-activity","suggestion_max_words":25,"max_targets":12,"suggestion_requirements":["sounds-achievable","nothing-that-obviously-deals-damage-to-targets-or-allies"],"suggestion_examples":["walk-to-village-and-help-villagers-harvest-crops-until-sunset","drop-weapons-and-dance-stop-in-an-hour"],"save":"wisdom","failure_condition":"charmed","condition_ends":["duration-ended","caster-or-allies-damage-target"],"pursuit":"pursues-suggestion-to-best-of-its-ability","activity_duration":"can-continue-for-entire-duration","completion":"spell-ends-for-target-upon-completing-shorter-activity","slot_durations":[{"min_slot":7,"max_slot":7,"seconds":864000},{"min_slot":8,"max_slot":8,"seconds":2592000},{"min_slot":9,"max_slot":9,"seconds":31622400}]}}}`,
		`{"id":"srd-5.2.1:content/spells/maze","name":"Maze","source_id":"srd-5.2.1:content/spells/maze","source_page":148,"scope":"complete-source-record","level":8,"school":"conjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"banish","mechanic":"banishment.maze","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"return","mechanic":"return.maze","timing":"when-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["banishment.maze","return.maze","check.maze-escape-investigation","concentration.source-duration"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-escape","maze":{"destination":"labyrinthine-demiplane","remains":"for-duration-or-until-escape","escape_action":"study","escape_check":"intelligence-investigation","escape_dc":20,"escape_success":"escapes-and-spell-ends","return_space":"space-it-left","occupied_return":"nearest-unoccupied-space"}}}`,
		`{"id":"srd-5.2.1:content/spells/meld-into-stone","name":"Meld into Stone","source_id":"srd-5.2.1:content/spells/meld-into-stone","source_page":148,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["cleric","druid","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"stone-object-or-surface-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"meld","mechanic":"transformation.meld-into-stone","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"partial","mechanic":"damage.meld-into-stone-partial-expulsion","timing":"when-partly-destroyed-or-reshaped","dice_count":6,"die_sides":6,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"destroyed","mechanic":"damage.meld-into-stone-destroyed-expulsion","timing":"when-destroyed-or-transmuted","dice_count":0,"die_sides":0,"damage_type":"force","bonus":50,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.meld-into-stone","damage.meld-into-stone-partial-expulsion","damage.meld-into-stone-destroyed-expulsion","expulsion.meld-into-stone-prone","cast.ritual"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-left-or-expelled","meld_into_stone":{"stone_requirement":"stone-object-or-surface-large-enough-to-fully-contain-body","merged":["body","equipment"],"touch_requirement":"must-touch-stone","concealment":"nothing-visible-or-detectable-by-nonmagical-senses","outside_sight":"cannot-see","outside_hearing":"disadvantage-on-wisdom-perception-checks-to-hear","time_awareness":"aware-of-passage-of-time","self_casting":"can-cast-spells-on-self","exit_movement_feet":5,"exit":"leave-where-entered-which-ends-spell","other_movement":"cannot-move","minor_damage":"minor-physical-damage-does-not-harm","partial_triggers":["partial-destruction","shape-change-so-caster-no-longer-fits"],"partial_dice_count":6,"partial_die_sides":6,"partial_damage_type":"force","destroyed_triggers":["complete-destruction","transmutation-into-different-substance"],"destroyed_damage":50,"destroyed_damage_type":"force","expelled_space":"unoccupied-space-closest-to-where-first-entered","expelled_condition":"prone"}}}`,
		`{"id":"srd-5.2.1:content/spells/mending","name":"Mending","source_id":"srd-5.2.1:content/spells/mending","source_page":148,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["bard","cleric","druid","sorcerer","wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"two lodestones","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"repair","mechanic":"repair.mending","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["repair.mending","cast.one-minute","components.material"],"slot_levels":[],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"instantaneous","mending":{"repairs":"single-break-or-tear","examples":["broken-chain-link","two-halves-of-broken-key","torn-cloak","leaking-wineskin"],"max_dimension_feet":1,"result":"leaves-no-trace-of-former-damage","magic_item":"can-physically-repair-but-cannot-restore-magic"}}}`,
		`{"id":"srd-5.2.1:content/spells/message","name":"Message","source_id":"srd-5.2.1:content/spells/message","source_page":148,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["bard","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":false,"somatic":true,"material":{"description":"a copper wire","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":6,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"message","mechanic":"communication.message-whisper","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["communication.message-whisper","blocker.message-barriers","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","message":{"source_pages":[148,149],"gesture":"point-toward-creature","delivery":"whisper-heard-only-by-target","reply":"target-can-reply-in-whisper-only-caster-hears","through_objects":"can-cast-through-solid-objects","through_requirements":["familiar-with-target","knows-target-is-beyond-barrier"],"blockers":["magical-silence","1-foot-of-stone-metal-or-wood","thin-sheet-of-lead"],"solid_blocker_feet":1}}}`,
		`{"id":"srd-5.2.1:content/spells/meteor-swarm","name":"Meteor Swarm","source_id":"srd-5.2.1:content/spells/meteor-swarm","source_page":149,"scope":"complete-source-record","level":9,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":5280},"target":{"shape":"spheres","count":4,"eligibility":"creatures-in-spheres-on-visible-points"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"fire","mechanic":"damage.meteor-swarm-fire","timing":"immediate","dice_count":20,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"bludgeoning","mechanic":"damage.meteor-swarm-bludgeoning","timing":"immediate","dice_count":20,"die_sides":6,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.meteor-swarm-fire","damage.meteor-swarm-bludgeoning","save.dexterity-half","object.meteor-swarm-damage-and-ignite"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","meteor_swarm":{"range_miles":1,"orbs":"blazing-orbs-of-fire","point_count":4,"points":"different-points-visible-within-range","sphere_radius_feet":40,"save":"dexterity","fire_dice_count":20,"fire_die_sides":6,"bludgeoning_dice_count":20,"bludgeoning_die_sides":6,"save_success":"half-damage","overlap":"creature-in-more-than-one-sphere-affected-only-once","object_eligibility":"nonmagical-object-not-worn-or-carried","object_effect":"takes-damage-if-in-area","flammable_object":"starts-burning"}}}`,
		`{"id":"srd-5.2.1:content/spells/mind-blank","name":"Mind Blank","source_id":"srd-5.2.1:content/spells/mind-blank","source_page":149,"scope":"complete-source-record","level":8,"school":"abjuration","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"immunity","mechanic":"immunity.mind-blank-psychic-and-charmed","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ward","mechanic":"protection.mind-blank-ward","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["immunity.mind-blank-psychic-and-charmed","protection.mind-blank-ward"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","mind_blank":{"target":"one-willing-creature-touched","damage_immunity":"psychic","condition_immunity":"charmed","unaffected_by":["sensing-emotions-or-alignment","reading-thoughts","magically-detecting-location"],"spells_cannot":["gather-information-about-target","observe-target-remotely","control-targets-mind"],"includes_spell":"wish"}}}`,
		`{"id":"srd-5.2.1:content/spells/mind-spike","name":"Mind Spike","source_id":"srd-5.2.1:content/spells/mind-spike","source_page":149,"scope":"complete-source-record","level":2,"school":"divination","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":false,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.mind-spike","timing":"immediate","dice_count":3,"die_sides":8,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"tracking","mechanic":"tracking.mind-spike-location","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.mind-spike","tracking.mind-spike-location","save.wisdom-half","scaling.per-slot","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","mind_spike":{"target":"one-creature-visible-in-range","save":"wisdom","damage_dice_count":3,"damage_die_sides":8,"damage_type":"psychic","save_success":"half-damage","failure_knowledge":"always-know-targets-location-until-spell-ends","knowledge_limit":"only-while-on-same-plane-of-existence","knowledge_benefits":["target-cannot-become-hidden-from-caster","target-gains-no-benefit-from-invisible-against-caster"],"slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/minor-illusion","name":"Minor Illusion","source_id":"srd-5.2.1:content/spells/minor-illusion","source_page":149,"scope":"complete-source-record","level":0,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"choice","count":1,"eligibility":"sound-or-object-image-in-range"},"components":{"verbal":false,"somatic":true,"material":{"description":"a bit of fleece","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sound","mechanic":"illusion.minor-illusion-sound","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"image","mechanic":"illusion.minor-illusion-image","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.minor-illusion-sound","illusion.minor-illusion-image","choice.minor-illusion-sound-or-image","check.minor-illusion-investigation","ending.minor-illusion-recast","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-cast-again","minor_illusion":{"options":["sound","image"],"recast":"illusion-ends-if-cast-again","discern_action":"study","discern_check":"intelligence-investigation","discern_dc":"caster-spell-save-dc","discerned":"illusion-becomes-faint-to-creature","sound_volume":"whisper-to-scream","sound_examples":["casters-voice","someone-elses-voice","lions-roar","beating-of-drums","any-other-chosen-sound"],"sound_timing":["continues-unabated-throughout-duration","discrete-sounds-at-different-times"],"image_subject":"image-of-an-object","image_examples":["chair","muddy-footprints","small-chest"],"image_max_cube_feet":5,"image_cannot_create":["sound","light","smell","other-sensory-effect"],"image_physical_interaction":"reveals-illusion-because-things-pass-through"}}}`,
		`{"id":"srd-5.2.1:content/spells/mirage-arcane","name":"Mirage Arcane","source_id":"srd-5.2.1:content/spells/mirage-arcane","source_page":149,"scope":"complete-source-record","level":7,"school":"illusion","lists":["bard","druid","wizard"],"casting_time":"10-minutes","range":{"kind":"sight","feet":0},"target":{"shape":"square","count":0,"eligibility":"terrain-area-up-to-1-mile-square"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":864000,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"illusion.mirage-arcane-terrain","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.mirage-arcane-terrain","range.sight","perception.mirage-arcane-truesight","cast.ten-minutes"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"duration-ended","mirage_arcane":{"source_pages":[149,150],"area_max_miles_square":1,"altered_senses":["look","sound","smell","feel"],"disguise":"another-sort-of-terrain","examples":["open-fields-or-road-as-swamp-hill-crevasse-or-other-rough-or-impassable-terrain","pond-as-grassy-meadow","precipice-as-gentle-slope","rock-strewn-gully-as-wide-smooth-road"],"structures":"can-alter-appearance-of-structures-or-add-them","creatures":"does-not-disguise-conceal-or-add-creatures","elements":["audible","visual","tactile","olfactory"],"movement_effect":"can-create-or-remove-difficult-terrain-or-otherwise-impede-movement","removed_pieces":"illusory-pieces-removed-from-area-disappear-immediately","removed_piece_examples":["rock","stick"],"truesight":"sees-through-illusion-to-true-form","truesight_limit":"other-elements-remain-and-can-be-physically-interacted-with"}}}`,
	}
	batch := loadSpells213(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Mass Heal", "Mass Healing Word", "Mass Suggestion", "Maze", "Meld into Stone", "Mending", "Message", "Meteor Swarm", "Mind Blank", "Mind Spike", "Minor Illusion", "Mirage Arcane"}
	pages := []int{147, 148, 148, 148, 148, 148, 148, 149, 149, 149, 149, 149}
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

func TestSpells213CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells213(t)
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
func TestSpells213MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_213.json")
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

func TestSpells213JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_213.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"hit_point_pool": 700`, `"hit_point_pool": 700, "hit_point_pool": 701`},
		{`"hit_point_pool": 700`, `"hit_point_pool": 700, "HIT_POINT_POOL": 701`},
		{`"hit_point_pool": 700`, `"hit_point_pool": 700, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/mass-heal"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "healing"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 60`, `"feet": 61`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
		{`"ritual": true`, `"ritual": false`},
		{`"kind": "sight"`, `"kind": "unlimited"`},
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

func TestSpells213InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"heal pool", 0, func(s *SpellRecord) { s.SourceDetails.MassHeal.HitPointPool = 800 }},
		{"heal conditions", 0, func(s *SpellRecord) {
			s.SourceDetails.MassHeal.ConditionsRemoved = s.SourceDetails.MassHeal.ConditionsRemoved[:2]
		}},
		{"word targets", 1, func(s *SpellRecord) { s.Target.Count = 12 }},
		{"word slot die", 1, func(s *SpellRecord) { s.Scaling.DieSides = 8 }},
		{"suggestion table", 2, func(s *SpellRecord) { s.SourceDetails.MassSuggestion.SlotDurations[2].Seconds = 31536000 }},
		{"suggestion concentration", 2, func(s *SpellRecord) { *s.Concentration = true }},
		{"maze dc", 3, func(s *SpellRecord) { s.SourceDetails.Maze.EscapeDC = 15 }},
		{"meld damage", 4, func(s *SpellRecord) { s.Effects[2].Bonus = 40 }},
		{"meld ritual", 4, func(s *SpellRecord) { *s.Ritual = false }},
		{"mending size", 5, func(s *SpellRecord) { s.SourceDetails.Mending.MaxDimensionFeet = 2 }},
		{"mending casting", 5, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"message pages", 6, func(s *SpellRecord) { s.SourceDetails.Message.SourcePages = []int{148} }},
		{"message duration", 6, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"meteor points", 7, func(s *SpellRecord) { s.SourceDetails.MeteorSwarm.PointCount = 3 }},
		{"meteor range", 7, func(s *SpellRecord) { s.Range.Feet = 1000 }},
		{"blank immunity", 8, func(s *SpellRecord) { s.SourceDetails.MindBlank.ConditionImmunity = "frightened" }},
		{"spike save", 9, func(s *SpellRecord) { s.SourceDetails.MindSpike.Save = "intelligence" }},
		{"spike verbal", 9, func(s *SpellRecord) { *s.Components.Verbal = true }},
		{"illusion options", 10, func(s *SpellRecord) { s.SourceDetails.MinorIllusion.Options = []string{"sound"} }},
		{"illusion cube", 10, func(s *SpellRecord) { s.SourceDetails.MinorIllusion.ImageMaxCubeFeet = 10 }},
		{"illusion scaling", 10, func(s *SpellRecord) { s.Scaling.Kind = "cantrip-upgrade" }},
		{"mirage range", 11, func(s *SpellRecord) { s.Range.Kind = "unlimited" }},
		{"mirage area", 11, func(s *SpellRecord) { s.SourceDetails.MirageArcane.AreaMaxMilesSquare = 10 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.MirageArcane = &SpellMirageArcane{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells213(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells213RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_213.json")
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

func TestSpells213WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_213.json")
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
		`"hit_point_pool": 700, "hit_point_pool": 701`,
		`"hit_point_pool": 700, "HIT_POINT_POOL": 701`,
		`"hit_point_pool": 700, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"hit_point_pool": 700`), []byte(replacement), 1)
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
func TestSpells213PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_213.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"hit_point_pool", 1, 700},
		{"max_targets", 1, 6},
		{"max_targets", 2, 12},
		{"healing_dice_count", 1, 2},
		{"healing_die_sides", 1, 4},
		{"slot_dice_count", 1, 1},
		{"slot_dice_count", 2, 1},
		{"slot_die_sides", 1, 4},
		{"slot_die_sides", 2, 8},
		{"suggestion_max_words", 1, 25},
		{"min_slot", 1, 7},
		{"min_slot", 2, 8},
		{"min_slot", 3, 9},
		{"max_slot", 1, 7},
		{"max_slot", 2, 8},
		{"max_slot", 3, 9},
		{"seconds", 1, 864000},
		{"seconds", 2, 2592000},
		{"seconds", 3, 31622400},
		{"escape_dc", 1, 20},
		{"exit_movement_feet", 1, 5},
		{"partial_dice_count", 1, 6},
		{"partial_die_sides", 1, 6},
		{"destroyed_damage", 1, 50},
		{"max_dimension_feet", 1, 1},
		{"solid_blocker_feet", 1, 1},
		{"range_miles", 1, 1},
		{"point_count", 1, 4},
		{"sphere_radius_feet", 1, 40},
		{"fire_dice_count", 1, 20},
		{"fire_die_sides", 1, 6},
		{"bludgeoning_dice_count", 1, 20},
		{"bludgeoning_die_sides", 1, 6},
		{"damage_dice_count", 1, 3},
		{"damage_die_sides", 1, 8},
		{"image_max_cube_feet", 1, 5},
		{"area_max_miles_square", 1, 1},
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
// inventory assigns to catalogue issue #213.
func TestSpells213MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 213 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells213(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells213ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_213.json")
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
	want := "SRD 5.2.1 PDF pp.147-150 adapted into twelve complete offline spell families for issue #213, including page-break continuation text for Message (pp.148-149) and Mirage Arcane (pp.149-150). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells213PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.MirageArcane = loadSpells213(t).Records[11].SourceDetails.MirageArcane
	if r.Validate() == nil {
		t.Fatal("batch 213 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells213RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells213(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell213Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell213Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
