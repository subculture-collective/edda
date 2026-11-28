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

func loadSpells220(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_220.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.164-166 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells220SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/speak-with-dead","name":"Speak with Dead","source_id":"srd-5.2.1:content/spells/speak-with-dead","source_page":164,"scope":"complete-source-record","level":3,"school":"necromancy","lists":["bard","cleric","wizard"],"casting_time":"action","range":{"kind":"feet","feet":10},"target":{"shape":"single","count":1,"eligibility":"corpse-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"burning incense","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"speech","mechanic":"utility.speak-with-dead","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.speak-with-dead","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","speak_with_dead":{"corpse_requirement":"corpse-must-have-mouth","failure_conditions":["deceased-was-undead-when-it-died","corpse-targeted-by-spell-within-past-10-days"],"repeat_block_days":10,"max_questions":5,"knowledge":"only-what-it-knew-in-life-including-languages","answer_style":["brief","cryptic","repetitive"],"truthfulness":"no-compulsion-to-be-truthful-if-antagonized-or-recognizes-enemy","soul":"returns-animating-spirit-not-soul","limits":["cannot-learn-new-information","does-not-comprehend-events-since-death","cannot-speculate-about-future"]}}}`,
		`{"id":"srd-5.2.1:content/spells/speak-with-plants","name":"Speak with Plants","source_id":"srd-5.2.1:content/spells/speak-with-plants","source_page":164,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["bard","druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"plants-in-immobile-emanation"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"sentience","mechanic":"utility.speak-with-plants","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"terrain","mechanic":"terrain.speak-with-plants-change","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.speak-with-plants","terrain.speak-with-plants-change"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","speak_with_plants":{"emanation_feet":30,"emanation_mobility":"immobile","grants":"limited-sentience-animation-communication-and-simple-commands","question_topics":["creatures-that-passed","weather","other-circumstances"],"question_window_days":1,"terrain_options":["plant-growth-difficult-terrain-to-ordinary","ordinary-terrain-with-plants-to-difficult"],"movement_limit":"plants-cannot-uproot-but-move-branches-tendrils-and-stalks","plant_creatures":"communicate-as-if-sharing-common-language"}}}`,
		`{"id":"srd-5.2.1:content/spells/spider-climb","name":"Spider Climb","source_id":"srd-5.2.1:content/spells/spider-climb","source_page":164,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of bitumen and a spider","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":2,"dice_count":0,"die_sides":0,"effect_ids":["climb"]},"effects":[{"id":"climb","mechanic":"movement.spider-climb","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.spider-climb","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","spider_climb":{"surfaces":["up-down-and-across-vertical-surfaces","along-ceilings"],"hands":"left-free","climb_speed":"equal-to-speed","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/spike-growth","name":"Spike Growth","source_id":"srd-5.2.1:content/spells/spike-growth","source_page":164,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"sphere","count":0,"eligibility":"ground-in-sphere-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"seven thorns","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"terrain.spike-growth","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.spike-growth","timing":"each-5-feet-moved-into-or-within","dice_count":2,"die_sides":4,"damage_type":"piercing","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.spike-growth","damage.spike-growth","check.spike-growth-recognition","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","spike_growth":{"radius_feet":20,"terrain":"difficult-terrain","damage_interval_feet":5,"damage_dice_count":2,"damage_die_sides":4,"camouflage":"looks-natural","recognition_requirement":"did-not-see-area-when-cast","recognition_action":"search","recognition_check":"wisdom-perception-or-survival","recognition_dc":"caster-spell-save-dc","recognition_timing":"before-entering"}}}`,
		`{"id":"srd-5.2.1:content/spells/spirit-guardians","name":"Spirit Guardians","source_id":"srd-5.2.1:content/spells/spirit-guardians","source_page":164,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"undesignated-creatures-in-emanation"},"components":{"verbal":true,"somatic":true,"material":{"description":"a prayer scroll","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":8,"effect_ids":["radiant","necrotic"]},"effects":[{"id":"speed","mechanic":"movement.spirit-guardians-half-speed","timing":"while-in-emanation","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"radiant","mechanic":"damage.spirit-guardians-radiant","timing":"emanation-enters-creature-enters-or-ends-turn","dice_count":3,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"necrotic","mechanic":"damage.spirit-guardians-necrotic","timing":"emanation-enters-creature-enters-or-ends-turn","dice_count":3,"die_sides":8,"damage_type":"necrotic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.spirit-guardians-half-speed","damage.spirit-guardians-radiant","damage.spirit-guardians-necrotic","choice.spirit-guardians-alignment","save.wisdom-half","scaling.per-slot","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","spirit_guardians":{"source_pages":[164,165],"emanation_feet":15,"forms":[{"alignment":"good","appearance_choices":["angelic","fey"],"damage_type":"radiant"},{"alignment":"neutral","appearance_choices":["angelic","fey"],"damage_type":"radiant"},{"alignment":"evil","appearance_choices":["fiendish"],"damage_type":"necrotic"}],"exemptions":"creatures-designated-at-casting","speed_divisor":2,"save_triggers":["emanation-enters-creature-space","creature-enters-emanation","creature-ends-turn-there"],"save":"wisdom","damage_dice_count":3,"damage_die_sides":8,"save_success":"half-damage","saves_per_turn":1,"slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/spiritual-weapon","name":"Spiritual Weapon","source_id":"srd-5.2.1:content/spells/spiritual-weapon","source_page":165,"scope":"complete-source-record","level":2,"school":"evocation","lists":["cleric"],"casting_time":"bonus-action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"chosen-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.spiritual-weapon","timing":"on-hit-at-casting-and-on-later-bonus-actions","dice_count":1,"die_sides":8,"damage_type":"force","bonus":0,"add_casting_modifier":true,"half_initial_only":false},{"id":"movement","mechanic":"movement.spiritual-weapon","timing":"bonus-action-on-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.spiritual-weapon","movement.spiritual-weapon","attack.melee-spell","scaling.per-slot","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","spiritual_weapon":{"appearance":"spectral-weapon-of-casters-choice","attack":"melee-spell-attack","attack_reach_feet":5,"attack_targets":1,"damage_dice_count":1,"damage_die_sides":8,"damage_adds":"spellcasting-ability-modifier","damage_type":"force","later_action":"bonus-action","move_max_feet":20,"repeat_attack":"against-creature-within-reach","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/starry-wisp","name":"Starry Wisp","source_id":"srd-5.2.1:content/spells/starry-wisp","source_page":165,"scope":"complete-source-record","level":0,"school":"evocation","lists":["bard","druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.starry-wisp","timing":"on-hit","dice_count":1,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"revealed","mechanic":"light.starry-wisp-revealed","timing":"on-hit-until-end-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.starry-wisp","light.starry-wisp-revealed","attack.ranged-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","starry_wisp":{"attack":"ranged-spell-attack","target_kinds":["creature","object"],"dim_light_feet":10,"invisibility":"cannot-benefit-from-invisible-condition","effect_ends":"end-of-casters-next-turn","upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
		`{"id":"srd-5.2.1:content/spells/stinking-cloud","name":"Stinking Cloud","source_id":"srd-5.2.1:content/spells/stinking-cloud","source_page":165,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"sphere","count":0,"eligibility":"area-centered-on-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a rotten egg","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"cloud","mechanic":"terrain.stinking-cloud","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"poisoned","mechanic":"condition.stinking-cloud-poisoned","timing":"start-turn-in-sphere-until-end-of-current-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.stinking-cloud","condition.stinking-cloud-poisoned","save.constitution-negates","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-dispersed-by-strong-wind","stinking_cloud":{"radius_feet":20,"appearance":"yellow-nauseating-gas","obscurement":"heavily-obscured","lingers":"for-duration-or-until-strong-wind-disperses-it","wind_example":"gust-of-wind","save":"constitution","save_trigger":"starts-turn-in-sphere","condition":"poisoned","condition_ends":"end-of-current-turn","poisoned_limit":"cannot-take-action-or-bonus-action"}}}`,
		`{"id":"srd-5.2.1:content/spells/stone-shape","name":"Stone Shape","source_id":"srd-5.2.1:content/spells/stone-shape","source_page":165,"scope":"complete-source-record","level":4,"school":"transmutation","lists":["cleric","druid","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"stone-object-or-section-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"soft clay","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"shape","mechanic":"creation.stone-shape","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["creation.stone-shape","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","stone_shape":{"max_object_size":"medium","max_section_dimension_feet":5,"result":"any-shape-caster-likes","examples":["large-rock-into-weapon-statue-or-coffer","passage-through-wall","seal-stone-door-or-frame"],"example_wall_thickness_feet":5,"max_hinges":2,"max_latches":1,"detail_limit":"no-finer-mechanical-detail"}}}`,
		`{"id":"srd-5.2.1:content/spells/stoneskin","name":"Stoneskin","source_id":"srd-5.2.1:content/spells/stoneskin","source_page":165,"scope":"complete-source-record","level":4,"school":"transmutation","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"diamond dust worth 100+ GP, which the spell consumes","cost_cp":10000,"consumed":true}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"resistance","mechanic":"defense.stoneskin-resistance","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.stoneskin-resistance","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","stoneskin":{"resistance_damage_types":["bludgeoning","piercing","slashing"]}}}`,
		`{"id":"srd-5.2.1:content/spells/storm-of-vengeance","name":"Storm of Vengeance","source_id":"srd-5.2.1:content/spells/storm-of-vengeance","source_page":165,"scope":"complete-source-record","level":9,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":5280},"target":{"shape":"sphere","count":0,"eligibility":"creatures-and-objects-under-storm-cloud"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"initial","mechanic":"damage.storm-of-vengeance-thunder","timing":"when-cloud-appears","dice_count":2,"die_sides":6,"damage_type":"thunder","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"deafened","mechanic":"condition.storm-of-vengeance-deafened","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"acid","mechanic":"damage.storm-of-vengeance-acid","timing":"start-of-casters-turn-2","dice_count":4,"die_sides":6,"damage_type":"acid","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"lightning","mechanic":"damage.storm-of-vengeance-lightning","timing":"start-of-casters-turn-3","dice_count":10,"die_sides":6,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"hail","mechanic":"damage.storm-of-vengeance-hail","timing":"start-of-casters-turn-4","dice_count":2,"die_sides":6,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"freezing","mechanic":"damage.storm-of-vengeance-cold","timing":"start-of-casters-turns-5-to-10","dice_count":1,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.storm-of-vengeance-thunder","condition.storm-of-vengeance-deafened","damage.storm-of-vengeance-acid","damage.storm-of-vengeance-lightning","damage.storm-of-vengeance-hail","damage.storm-of-vengeance-cold","save.constitution-storm-of-vengeance","save.dexterity-storm-of-vengeance-lightning","terrain.storm-of-vengeance-gusts","concentration.source-duration"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","storm_of_vengeance":{"source_pages":[165,166],"range_miles":1,"cloud_radius_feet":300,"initial_save":"constitution","initial_dice_count":2,"initial_die_sides":6,"initial_damage_type":"thunder","initial_condition":"deafened","initial_condition_duration":"for-duration","turn_timing":"start-of-each-of-casters-later-turns","turns":[{"first_turn":2,"last_turn":2,"effect":"acidic-rain-on-each-creature-and-object","dice_count":4,"die_sides":6,"damage_type":"acid","save":"none"},{"first_turn":3,"last_turn":3,"effect":"lightning-bolts-strike-chosen-targets","dice_count":10,"die_sides":6,"damage_type":"lightning","save":"dexterity-half"},{"first_turn":4,"last_turn":4,"effect":"hailstones-on-each-creature","dice_count":2,"die_sides":6,"damage_type":"bludgeoning","save":"none"},{"first_turn":5,"last_turn":10,"effect":"gusts-and-freezing-rain-on-each-creature","dice_count":1,"die_sides":6,"damage_type":"cold","save":"none"}],"lightning_bolts":6,"lightning_targets":"six-different-creatures-or-objects","gust_effects":["difficult-terrain","heavily-obscured","ranged-weapon-attacks-impossible","strong-wind"],"gust_effects_end":"spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/suggestion","name":"Suggestion","source_id":"srd-5.2.1:content/spells/suggestion","source_page":166,"scope":"complete-source-record","level":2,"school":"enchantment","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-hearing-understanding-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":{"description":"a drop of honey","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"charm","mechanic":"condition.suggestion-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"compulsion","mechanic":"control.suggestion-activity","timing":"while-charmed","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.suggestion-charmed","control.suggestion-activity","save.wisdom-negates","ending.suggestion-damage-or-completion","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-damage-or-activity-completed","suggestion":{"max_words":25,"requirements":["sounds-achievable","no-obvious-damage-to-target-or-allies"],"examples":["fetch-cult-treasure-vault-key-and-give-it-to-caster","stop-fighting-leave-library-peacefully-and-do-not-return"],"save":"wisdom","failure_condition":"charmed","condition_ends":["duration-ends","caster-or-allies-damage-target"],"pursuit":"pursues-suggestion-to-best-of-ability","activity_duration":"may-continue-for-entire-duration","completion":"spell-ends-for-target-when-activity-completed-sooner"}}}`,
	}
	batch := loadSpells220(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Speak with Dead", "Speak with Plants", "Spider Climb", "Spike Growth", "Spirit Guardians", "Spiritual Weapon", "Starry Wisp", "Stinking Cloud", "Stone Shape", "Stoneskin", "Storm of Vengeance", "Suggestion"}
	pages := []int{164, 164, 164, 164, 164, 165, 165, 165, 165, 165, 165, 166}
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

func TestSpells220CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells220(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t), loadSpells220(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 300 {
		t.Fatalf("expected 300 records, got %d", len(all.Records))
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
func TestSpells220MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_220.json")
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

func TestSpells220JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_220.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"repeat_block_days": 10`, `"repeat_block_days": 10, "repeat_block_days": 11`},
		{`"repeat_block_days": 10`, `"repeat_block_days": 10, "REPEAT_BLOCK_DAYS": 11`},
		{`"repeat_block_days": 10`, `"repeat_block_days": 10, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/speak-with-dead"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "climb"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "bonus-action"`, `"casting_time": "action"`},
		{`"feet": 10`, `"feet": 11`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
		{`"concentration": true`, `"concentration": false`},
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

func TestSpells220InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"dead questions", 0, func(s *SpellRecord) { s.SourceDetails.SpeakWithDead.MaxQuestions = 10 }},
		{"dead block", 0, func(s *SpellRecord) { s.SourceDetails.SpeakWithDead.RepeatBlockDays = 7 }},
		{"plants area", 1, func(s *SpellRecord) { s.SourceDetails.SpeakWithPlants.EmanationFeet = 15 }},
		{"plants terrain", 1, func(s *SpellRecord) {
			s.SourceDetails.SpeakWithPlants.TerrainOptions = s.SourceDetails.SpeakWithPlants.TerrainOptions[:1]
		}},
		{"climb increment", 2, func(s *SpellRecord) { s.SourceDetails.SpiderClimb.AdditionalTargetsPerSlot = 2 }},
		{"climb scaling kind", 2, func(s *SpellRecord) { s.Scaling.Kind = "slot-duration" }},
		{"spike dice", 3, func(s *SpellRecord) { s.Effects[1].DieSides = 6 }},
		{"spike interval", 3, func(s *SpellRecord) { s.SourceDetails.SpikeGrowth.DamageIntervalFeet = 10 }},
		{"guardians radius", 4, func(s *SpellRecord) { s.SourceDetails.SpiritGuardians.EmanationFeet = 30 }},
		{"guardians evil damage", 4, func(s *SpellRecord) { s.SourceDetails.SpiritGuardians.Forms[2].DamageType = "radiant" }},
		{"weapon modifier", 5, func(s *SpellRecord) { *s.Effects[0].AddCastingModifier = false }},
		{"weapon move", 5, func(s *SpellRecord) { s.SourceDetails.SpiritualWeapon.MoveMaxFeet = 30 }},
		{"wisp light", 6, func(s *SpellRecord) { s.SourceDetails.StarryWisp.DimLightFeet = 20 }},
		{"wisp upgrade", 6, func(s *SpellRecord) { s.SourceDetails.StarryWisp.Upgrades[0].DiceCount = 3 }},
		{"cloud save", 7, func(s *SpellRecord) { s.SourceDetails.StinkingCloud.Save = "wisdom" }},
		{"cloud radius", 7, func(s *SpellRecord) { s.SourceDetails.StinkingCloud.RadiusFeet = 30 }},
		{"shape hinges", 8, func(s *SpellRecord) { s.SourceDetails.StoneShape.MaxHinges = 3 }},
		{"shape size", 8, func(s *SpellRecord) { s.SourceDetails.StoneShape.MaxObjectSize = "large" }},
		{"skin types", 9, func(s *SpellRecord) {
			s.SourceDetails.Stoneskin.ResistanceDamageTypes = s.SourceDetails.Stoneskin.ResistanceDamageTypes[:2]
		}},
		{"skin consumed", 9, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"storm turn", 10, func(s *SpellRecord) { s.SourceDetails.StormOfVengeance.Turns[1].DiceCount = 8 }},
		{"storm pages", 10, func(s *SpellRecord) { s.SourceDetails.StormOfVengeance.SourcePages = []int{165} }},
		{"suggestion words", 11, func(s *SpellRecord) { s.SourceDetails.Suggestion.MaxWords = 30 }},
		{"suggestion somatic", 11, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Suggestion = &SpellSuggestion{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells220(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells220RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_220.json")
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

func TestSpells220WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_220.json")
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
		`"repeat_block_days": 10, "repeat_block_days": 11`,
		`"repeat_block_days": 10, "REPEAT_BLOCK_DAYS": 11`,
		`"repeat_block_days": 10, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"repeat_block_days": 10`), []byte(replacement), 1)
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
func TestSpells220PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_220.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"repeat_block_days", 1, 10},
		{"max_questions", 1, 5},
		{"emanation_feet", 1, 30},
		{"question_window_days", 1, 1},
		{"base_targets", 1, 1},
		{"additional_targets_per_slot", 1, 1},
		{"radius_feet", 1, 20},
		{"damage_interval_feet", 1, 5},
		{"damage_dice_count", 1, 2},
		{"damage_die_sides", 1, 4},
		{"emanation_feet", 2, 15},
		{"speed_divisor", 1, 2},
		{"damage_dice_count", 2, 3},
		{"damage_die_sides", 2, 8},
		{"saves_per_turn", 1, 1},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
		{"attack_reach_feet", 1, 5},
		{"attack_targets", 1, 1},
		{"damage_dice_count", 3, 1},
		{"damage_die_sides", 3, 8},
		{"move_max_feet", 1, 20},
		{"slot_dice_count", 2, 1},
		{"slot_die_sides", 2, 8},
		{"dim_light_feet", 1, 10},
		{"level", 8, 5},
		{"dice_count", 21, 2},
		{"die_sides", 21, 8},
		{"level", 9, 11},
		{"dice_count", 22, 3},
		{"die_sides", 22, 8},
		{"level", 10, 17},
		{"dice_count", 23, 4},
		{"die_sides", 23, 8},
		{"radius_feet", 2, 20},
		{"max_section_dimension_feet", 1, 5},
		{"example_wall_thickness_feet", 1, 5},
		{"max_hinges", 1, 2},
		{"max_latches", 1, 1},
		{"range_miles", 1, 1},
		{"cloud_radius_feet", 1, 300},
		{"initial_dice_count", 1, 2},
		{"initial_die_sides", 1, 6},
		{"first_turn", 1, 2},
		{"last_turn", 1, 2},
		{"dice_count", 38, 4},
		{"die_sides", 38, 6},
		{"first_turn", 2, 3},
		{"last_turn", 2, 3},
		{"dice_count", 39, 10},
		{"die_sides", 39, 6},
		{"first_turn", 3, 4},
		{"last_turn", 3, 4},
		{"dice_count", 40, 2},
		{"die_sides", 40, 6},
		{"first_turn", 4, 5},
		{"last_turn", 4, 10},
		{"dice_count", 41, 1},
		{"die_sides", 41, 6},
		{"lightning_bolts", 1, 6},
		{"max_words", 1, 25},
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
// inventory assigns to catalogue issue #220.
func TestSpells220MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 220 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells220(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells220ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_220.json")
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
	want := "SRD 5.2.1 PDF pp.164-166 adapted into twelve complete offline spell families for issue #220, including page-break continuation text for Spirit Guardians (pp.164-165) and Storm of Vengeance (pp.165-166). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells220PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.Suggestion = loadSpells220(t).Records[11].SourceDetails.Suggestion
	if r.Validate() == nil {
		t.Fatal("batch 220 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells220RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells220(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell220Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell220Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
