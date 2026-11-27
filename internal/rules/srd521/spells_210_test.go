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

func loadSpells210(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_210.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.141-143 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells210SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/hold-person","name":"Hold Person","source_id":"srd-5.2.1:content/spells/hold-person","source_page":141,"scope":"complete-source-record","level":2,"school":"enchantment","lists":["bard","cleric","druid","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-humanoid-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a straight piece of iron","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":2,"dice_count":0,"die_sides":0,"effect_ids":["paralyzed","repeat"]},"effects":[{"id":"paralyzed","mechanic":"condition.hold-person-paralyzed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat","mechanic":"save.hold-person-repeat","timing":"end-of-each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.hold-person-paralyzed","save.hold-person-repeat","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","hold_person":{"target":"humanoid-visible-in-range","save":"wisdom","failure_condition":"paralyzed","repeat_timing":"end-of-each-of-its-turns","repeat_success":"spell-ends-on-target","base_targets":1,"additional_targets_per_slot":1,"additional_target_type":"humanoid"}}}`,
		`{"id":"srd-5.2.1:content/spells/holy-aura","name":"Holy Aura","source_id":"srd-5.2.1:content/spells/holy-aura","source_page":141,"scope":"complete-source-record","level":8,"school":"abjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"chosen-creatures-in-aura"},"components":{"verbal":true,"somatic":true,"material":{"description":"a reliquary worth 1,000+ GP","cost_cp":100000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"saves","mechanic":"advantage.holy-aura-saving-throws","timing":"while-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"attacks","mechanic":"disadvantage.holy-aura-attacks-against","timing":"while-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blinding","mechanic":"condition.holy-aura-blinded","timing":"when-fiend-or-undead-melee-hits-affected-creature","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["advantage.holy-aura-saving-throws","disadvantage.holy-aura-attacks-against","condition.holy-aura-blinded","save.constitution-negates","concentration.source-duration","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","holy_aura":{"emanation_feet":30,"beneficiaries":"creatures-of-casters-choice-while-in-aura","save_benefit":"advantage-on-all-saving-throws","attack_penalty":"other-creatures-have-disadvantage-on-attack-rolls-against-them","retaliation_types":["fiend","undead"],"retaliation_trigger":"hits-affected-creature-with-melee-attack-roll","retaliation_save":"constitution","retaliation_condition":"blinded","retaliation_duration":"until-end-of-attackers-next-turn"}}}`,
		`{"id":"srd-5.2.1:content/spells/hunters-mark","name":"Hunter’s Mark","source_id":"srd-5.2.1:content/spells/hunters-mark","source_page":141,"scope":"complete-source-record","level":1,"school":"divination","lists":["ranger"],"casting_time":"bonus-action","range":{"kind":"feet","feet":90},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-duration","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["damage","tracking"]},"effects":[{"id":"damage","mechanic":"damage.hunters-mark-extra","timing":"whenever-caster-hits-target-with-attack-roll","dice_count":1,"die_sides":6,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"tracking","mechanic":"advantage.hunters-mark-find","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.hunters-mark-extra","advantage.hunters-mark-find","retarget.hunters-mark-bonus-action","scaling.duration-by-slot","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","hunters_mark":{"quarry":"one-creature-visible-in-range","extra_damage_dice_count":1,"extra_damage_die_sides":6,"damage_type":"force","damage_trigger":"whenever-caster-hits-target-with-attack-roll","find_checks":["wisdom-perception","wisdom-survival"],"find_benefit":"advantage-on-checks-to-find-target","retarget_condition":"target-drops-to-0-hit-points-before-spell-ends","retarget_action":"bonus-action","retarget_result":"move-mark-to-new-creature-visible-in-range","slot_durations":[{"min_slot":3,"max_slot":4,"seconds":28800},{"min_slot":5,"max_slot":9,"seconds":86400}]}}}`,
		`{"id":"srd-5.2.1:content/spells/hypnotic-pattern","name":"Hypnotic Pattern","source_id":"srd-5.2.1:content/spells/hypnotic-pattern","source_page":141,"scope":"complete-source-record","level":3,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"cube","count":0,"eligibility":"creatures-seeing-pattern-in-cube"},"components":{"verbal":false,"somatic":true,"material":{"description":"a pinch of confetti","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"charm","mechanic":"condition.hypnotic-pattern-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.hypnotic-pattern-charmed","ending.hypnotic-pattern-damage-or-shake","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","hypnotic_pattern":{"cube_feet":30,"pattern":"twisting-pattern-of-colors-appears-for-a-moment-and-vanishes","affected":"each-creature-in-area-that-can-see-pattern","save":"wisdom","failure_condition":"charmed","charmed_condition":"incapacitated","charmed_speed_feet":0,"ends_for_creature":["takes-any-damage","someone-else-uses-action-to-shake-it-out-of-stupor"]}}}`,
		`{"id":"srd-5.2.1:content/spells/ice-knife","name":"Ice Knife","source_id":"srd-5.2.1:content/spells/ice-knife","source_page":141,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":false,"somatic":true,"material":{"description":"a drop of water or a piece of ice","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":6,"effect_ids":["cold"]},"effects":[{"id":"attack","mechanic":"attack.ice-knife-ranged-spell","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"piercing","mechanic":"damage.ice-knife-piercing","timing":"on-hit","dice_count":1,"die_sides":10,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"cold","mechanic":"damage.ice-knife-cold","timing":"explosion-on-hit-or-miss","dice_count":2,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["attack.ice-knife-ranged-spell","damage.ice-knife-piercing","damage.ice-knife-cold","save.dexterity-negates","scaling.per-slot","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","ice_knife":{"attack":"ranged-spell-attack","piercing_dice_count":1,"piercing_die_sides":10,"explosion_timing":"hit-or-miss","explosion_targets":"target-and-each-creature-within-radius","explosion_radius_feet":5,"save":"dexterity","cold_dice_count":2,"cold_die_sides":6,"save_success":"no-cold-damage","slot_cold_dice_count":1,"slot_cold_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/ice-storm","name":"Ice Storm","source_id":"srd-5.2.1:content/spells/ice-storm","source_page":141,"scope":"complete-source-record","level":4,"school":"evocation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":300},"target":{"shape":"cylinder","count":0,"eligibility":"creatures-in-cylinder-centered-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a mitten","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":4,"dice_count":1,"die_sides":10,"effect_ids":["bludgeoning"]},"effects":[{"id":"bludgeoning","mechanic":"damage.ice-storm-bludgeoning","timing":"immediate","dice_count":2,"die_sides":10,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"cold","mechanic":"damage.ice-storm-cold","timing":"immediate","dice_count":4,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"terrain","mechanic":"terrain.ice-storm-difficult","timing":"until-end-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.ice-storm-bludgeoning","damage.ice-storm-cold","terrain.ice-storm-difficult","save.dexterity-half","scaling.per-slot","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","ice_storm":{"source_pages":[141,142],"cylinder_radius_feet":20,"cylinder_height_feet":40,"center":"point-within-range","save":"dexterity","bludgeoning_dice_count":2,"bludgeoning_die_sides":10,"cold_dice_count":4,"cold_die_sides":6,"save_success":"half-damage","ground":"difficult-terrain","ground_duration":"until-end-of-casters-next-turn","slot_bludgeoning_dice_count":1,"slot_bludgeoning_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/identify","name":"Identify","source_id":"srd-5.2.1:content/spells/identify","source_page":142,"scope":"complete-source-record","level":1,"school":"divination","lists":["bard","wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"object-or-creature-touched-throughout-casting"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pearl worth 100+ GP","cost_cp":10000,"consumed":false}},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"identify","mechanic":"knowledge.identify","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["knowledge.identify","cast.one-minute","cast.ritual","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"instantaneous","identify":{"touch":"throughout-casting","object_eligibility":"magic-item-or-some-other-magical-object","object_reveals":["properties","how-to-use-them","whether-it-requires-attunement","number-of-charges-if-any","ongoing-spells-affecting-item-and-what-they-are","name-of-spell-that-created-item"],"creature_reveals":"ongoing-spells-currently-affecting-creature"}}}`,
		`{"id":"srd-5.2.1:content/spells/illusory-script","name":"Illusory Script","source_id":"srd-5.2.1:content/spells/illusory-script","source_page":142,"scope":"complete-source-record","level":1,"school":"illusion","lists":["bard","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"writing-material-touched"},"components":{"verbal":false,"somatic":true,"material":{"description":"ink worth 10+ GP, which the spell consumes","cost_cp":1000,"consumed":true}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":864000,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"script","mechanic":"illusion.illusory-script","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.illusory-script","perception.illusory-script-truesight","cast.one-minute","cast.ritual","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended-or-dispelled","illusory_script":{"materials":["parchment","paper","another-suitable-material"],"designated_readers":"caster-and-creatures-designated-when-cast","designated_view":["writing-appears-normal","seems-written-in-casters-hand","conveys-meaning-intended-when-written"],"others_view":"unknown-or-magical-script-that-is-unintelligible","alternative_alters":["meaning","handwriting","language"],"alternative_language":"must-be-one-caster-knows","dispelled":"original-script-and-illusion-both-disappear","truesight":"can-read-hidden-message"}}}`,
		`{"id":"srd-5.2.1:content/spells/imprisonment","name":"Imprisonment","source_id":"srd-5.2.1:content/spells/imprisonment","source_page":142,"scope":"complete-source-record","level":9,"school":"abjuration","lists":["warlock","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a statuette of the target worth 5,000+ GP","cost_cp":500000,"consumed":false}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"restraint","mechanic":"restraint.imprisonment","timing":"on-failed-wisdom-save-until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["restraint.imprisonment","save.wisdom-negates","choice.imprisonment-effect","ending.imprisonment-trigger","dispel.imprisonment-level-9-only","cast.one-minute","components.material"],"slot_levels":[9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"dispelled-or-specified-trigger","imprisonment":{"save":"wisdom","save_success":"unaffected-and-immune-to-this-spell","success_immunity_hours":24,"unneeded":["breathe","eat","drink"],"ages":false,"divination":"divination-spells-cannot-locate-or-perceive-target","teleport":"target-cannot-teleport","effect_choice_count":1,"effects":[{"name":"burial","rules":["entombed-beneath-earth","hollow-globe-of-magical-force-just-large-enough-to-contain-target","nothing-can-pass-into-or-out-of-globe"]},{"name":"chaining","rules":["chains-firmly-rooted-in-ground-hold-target","restrained-condition","cannot-be-moved-by-any-means"]},{"name":"hedged-prison","rules":["trapped-in-demiplane","demiplane-warded-against-teleportation-and-planar-travel","caster-chooses-labyrinth-cage-tower-or-the-like"]},{"name":"minimus-containment","rules":["target-becomes-1-inch-tall","trapped-in-indestructible-gemstone-or-similar-object","light-passes-through-gemstone-both-ways","nothing-else-passes-through-by-any-means"]},{"name":"slumber","rules":["unconscious-condition","cannot-be-awoken"]}],"minimus_height_inches":1,"trigger_specified":"when-spell-cast","trigger_complexity":"as-simple-or-as-elaborate-as-caster-chooses","trigger_agreement":"gm-agrees-high-likelihood-of-happening-within-window","trigger_window_years":10,"trigger_requirement":"observable-action","trigger_examples":["particular-offering-at-temple-of-casters-god","saving-casters-true-love","defeating-specific-monster"],"dispel_magic_slot_level":9,"dispel_magic_targets":["prison","component-used-to-create-it"]}}}`,
		`{"id":"srd-5.2.1:content/spells/incendiary-cloud","name":"Incendiary Cloud","source_id":"srd-5.2.1:content/spells/incendiary-cloud","source_page":142,"scope":"complete-source-record","level":8,"school":"conjuration","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"cloud","mechanic":"obscurement.incendiary-cloud","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.incendiary-cloud","timing":"appear-moves-into-enter-or-end-turn","dice_count":10,"die_sides":8,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"drift","mechanic":"movement.incendiary-cloud-drift","timing":"start-of-each-of-casters-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["obscurement.incendiary-cloud","damage.incendiary-cloud","movement.incendiary-cloud-drift","save.dexterity-half","ending.incendiary-cloud-strong-wind","concentration.source-duration"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-strong-wind","incendiary_cloud":{"source_pages":[142,143],"sphere_radius_feet":20,"contents":["embers","smoke"],"center":"point-within-range","obscurement":"heavily-obscured","dispersed_by":"strong-wind","wind_example":"gust-of-wind","save":"dexterity","damage_dice_count":10,"damage_die_sides":8,"damage_type":"fire","save_success":"half-damage","save_triggers":["cloud-appears-with-creature-in-it","sphere-moves-into-creature-space","creature-enters-sphere","creature-ends-turn-in-sphere"],"saves_per_turn":1,"drift_feet":10,"drift_direction":"away-from-caster-in-chosen-direction","drift_timing":"start-of-each-of-casters-turns"}}}`,
		`{"id":"srd-5.2.1:content/spells/inflict-wounds","name":"Inflict Wounds","source_id":"srd-5.2.1:content/spells/inflict-wounds","source_page":143,"scope":"complete-source-record","level":1,"school":"necromancy","lists":["cleric"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.inflict-wounds","timing":"immediate","dice_count":2,"die_sides":10,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.inflict-wounds","save.constitution-half","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","inflict_wounds":{"target":"creature-touched","save":"constitution","damage_dice_count":2,"damage_die_sides":10,"damage_type":"necrotic","save_success":"half-damage","slot_dice_count":1,"slot_die_sides":10}}}`,
		`{"id":"srd-5.2.1:content/spells/insect-plague","name":"Insect Plague","source_id":"srd-5.2.1:content/spells/insect-plague","source_page":143,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["cleric","druid","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":300},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"a locust","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"swarm","mechanic":"terrain.insect-plague-obscured-difficult","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.insect-plague","timing":"appear-first-entry-or-end-turn","dice_count":4,"die_sides":10,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.insect-plague-obscured-difficult","damage.insect-plague","save.constitution-half","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","insect_plague":{"sphere_radius_feet":20,"contents":"swarming-locusts","center":"point-chosen-within-range","obscurement":"lightly-obscured","terrain":"difficult-terrain","save":"constitution","damage_dice_count":4,"damage_die_sides":10,"damage_type":"piercing","save_success":"half-damage","save_triggers":["swarm-appears-with-creature-in-it","creature-enters-area-first-time-on-turn","creature-ends-turn-in-area"],"saves_per_turn":1,"slot_dice_count":1,"slot_die_sides":10}}}`,
	}
	batch := loadSpells210(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Hold Person", "Holy Aura", "Hunter’s Mark", "Hypnotic Pattern", "Ice Knife", "Ice Storm", "Identify", "Illusory Script", "Imprisonment", "Incendiary Cloud", "Inflict Wounds", "Insect Plague"}
	pages := []int{141, 141, 141, 141, 141, 141, 142, 142, 142, 142, 143, 143}
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

func TestSpells210CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells210(t)
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
func TestSpells210MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_210.json")
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

func TestSpells210JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_210.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"base_targets": 1`, `"base_targets": 1, "base_targets": 2`},
		{`"base_targets": 1`, `"base_targets": 1, "BASE_TARGETS": 2`},
		{`"base_targets": 1`, `"base_targets": 1, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/hold-person"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "paralyzed",
          "repeat"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 60`, `"feet": 61`},
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

func TestSpells210InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"hold person eligibility", 0, func(s *SpellRecord) { s.Target.Eligibility = "visible-creature-in-range" }},
		{"hold person list", 0, func(s *SpellRecord) { s.Lists = s.Lists[:5] }},
		{"holy aura radius", 1, func(s *SpellRecord) { s.SourceDetails.HolyAura.EmanationFeet = 10 }},
		{"holy aura consumption", 1, func(s *SpellRecord) { *s.Components.Material.Consumed = true }},
		{"mark table", 2, func(s *SpellRecord) { s.SourceDetails.HuntersMark.SlotDurations[0].Seconds = 14400 }},
		{"mark damage type", 2, func(s *SpellRecord) { s.Effects[0].DamageType = "necrotic" }},
		{"pattern verbal", 3, func(s *SpellRecord) { *s.Components.Verbal = true }},
		{"pattern speed", 3, func(s *SpellRecord) { s.SourceDetails.HypnoticPattern.CharmedSpeedFeet = 5 }},
		{"ice knife radius", 4, func(s *SpellRecord) { s.SourceDetails.IceKnife.ExplosionRadiusFeet = 10 }},
		{"ice knife scaling effect", 4, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"piercing"} }},
		{"ice storm pages", 5, func(s *SpellRecord) { s.SourceDetails.IceStorm.SourcePages = []int{141} }},
		{"ice storm height", 5, func(s *SpellRecord) { s.SourceDetails.IceStorm.CylinderHeightFeet = 20 }},
		{"identify ritual", 6, func(s *SpellRecord) { *s.Ritual = false }},
		{"identify casting", 6, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"script duration", 7, func(s *SpellRecord) { s.DurationSeconds = 86400 }},
		{"script truesight", 7, func(s *SpellRecord) { s.SourceDetails.IllusoryScript.Truesight = "none" }},
		{"imprisonment variants", 8, func(s *SpellRecord) { s.SourceDetails.Imprisonment.Effects = s.SourceDetails.Imprisonment.Effects[:4] }},
		{"imprisonment variant rule", 8, func(s *SpellRecord) { s.SourceDetails.Imprisonment.Effects[4].Rules[0] = "stunned-condition" }},
		{"imprisonment aging", 8, func(s *SpellRecord) { s.SourceDetails.Imprisonment.Ages = true }},
		{"imprisonment duration", 8, func(s *SpellRecord) { s.Duration = "instantaneous" }},
		{"cloud dice", 9, func(s *SpellRecord) { s.Effects[1].DiceCount = 8 }},
		{"cloud drift", 9, func(s *SpellRecord) { s.SourceDetails.IncendiaryCloud.DriftFeet = 15 }},
		{"wounds save", 10, func(s *SpellRecord) { s.SourceDetails.InflictWounds.Save = "dexterity" }},
		{"plague duration", 11, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"plague slot die", 11, func(s *SpellRecord) { s.SourceDetails.InsectPlague.SlotDieSides = 8 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.InsectPlague = &SpellInsectPlague{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells210(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells210RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_210.json")
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

func TestSpells210WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_210.json")
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
		`"base_targets": 1, "base_targets": 2`,
		`"base_targets": 1, "BASE_TARGETS": 2`,
		`"base_targets": 1, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"base_targets": 1`), []byte(replacement), 1)
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
func TestSpells210PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_210.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"base_targets", 1, 1},
		{"additional_targets_per_slot", 1, 1},
		{"emanation_feet", 1, 30},
		{"extra_damage_dice_count", 1, 1},
		{"extra_damage_die_sides", 1, 6},
		{"min_slot", 1, 3},
		{"min_slot", 2, 5},
		{"max_slot", 1, 4},
		{"max_slot", 2, 9},
		{"seconds", 1, 28800},
		{"seconds", 2, 86400},
		{"cube_feet", 1, 30},
		{"charmed_speed_feet", 1, 0},
		{"piercing_dice_count", 1, 1},
		{"piercing_die_sides", 1, 10},
		{"explosion_radius_feet", 1, 5},
		{"cold_dice_count", 1, 2},
		{"cold_dice_count", 2, 4},
		{"cold_die_sides", 1, 6},
		{"cold_die_sides", 2, 6},
		{"slot_cold_dice_count", 1, 1},
		{"slot_cold_die_sides", 1, 6},
		{"cylinder_radius_feet", 1, 20},
		{"cylinder_height_feet", 1, 40},
		{"bludgeoning_dice_count", 1, 2},
		{"bludgeoning_die_sides", 1, 10},
		{"slot_bludgeoning_dice_count", 1, 1},
		{"slot_bludgeoning_die_sides", 1, 10},
		{"success_immunity_hours", 1, 24},
		{"effect_choice_count", 1, 1},
		{"minimus_height_inches", 1, 1},
		{"trigger_window_years", 1, 10},
		{"dispel_magic_slot_level", 1, 9},
		{"sphere_radius_feet", 1, 20},
		{"sphere_radius_feet", 2, 20},
		{"damage_dice_count", 1, 10},
		{"damage_dice_count", 2, 2},
		{"damage_dice_count", 3, 4},
		{"damage_die_sides", 1, 8},
		{"damage_die_sides", 2, 10},
		{"damage_die_sides", 3, 10},
		{"saves_per_turn", 1, 1},
		{"saves_per_turn", 2, 1},
		{"drift_feet", 1, 10},
		{"slot_dice_count", 1, 1},
		{"slot_dice_count", 2, 1},
		{"slot_die_sides", 1, 10},
		{"slot_die_sides", 2, 10},
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
// inventory assigns to catalogue issue #210.
func TestSpells210MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 210 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells210(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells210ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_210.json")
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
	want := "SRD 5.2.1 PDF pp.141-143 adapted into twelve complete offline spell families for issue #210, including page-break continuation text for Ice Storm (pp.141-142) and Incendiary Cloud (pp.142-143). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells210PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.InsectPlague = loadSpells210(t).Records[11].SourceDetails.InsectPlague
	if r.Validate() == nil {
		t.Fatal("batch 210 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells210RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells210(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell210Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell210Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
