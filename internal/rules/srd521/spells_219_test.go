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

func loadSpells219(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_219.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.162-164 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells219SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/shillelagh","name":"Shillelagh","source_id":"srd-5.2.1:content/spells/shillelagh","source_page":162,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["druid"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"held-club-or-quarterstaff"},"components":{"verbal":true,"somatic":true,"material":{"description":"mistletoe","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"cantrip-die-table","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":["weapon-die"]},"effects":[{"id":"weapon-die","mechanic":"weapon.shillelagh-damage-die","timing":"for-duration","dice_count":1,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ability","mechanic":"weapon.shillelagh-spellcasting-ability","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage-type","mechanic":"weapon.shillelagh-damage-type-choice","timing":"when-attack-deals-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["weapon.shillelagh-damage-die","weapon.shillelagh-spellcasting-ability","weapon.shillelagh-damage-type-choice","ending.shillelagh-recast-or-release","scaling.cantrip-die-upgrade","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-recast-or-weapon-released","shillelagh":{"weapon_choices":["club","quarterstaff"],"weapon_requirement":"held-by-caster","ability_substitution":"may-use-spellcasting-ability-instead-of-strength","rolls":["melee-attack-rolls","melee-damage-rolls"],"damage_dice_count":1,"damage_die_sides":8,"damage_type_choices":["force","weapons-normal-damage-type"],"early_end":["spell-cast-again","caster-lets-go-of-weapon"],"die_upgrades":[{"level":5,"dice_count":1,"die_sides":10},{"level":11,"dice_count":1,"die_sides":12},{"level":17,"dice_count":2,"die_sides":6}]}}}`,
		`{"id":"srd-5.2.1:content/spells/shining-smite","name":"Shining Smite","source_id":"srd-5.2.1:content/spells/shining-smite","source_page":162,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["paladin"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"hit-target"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.shining-smite","timing":"as-hit-extra-attack-damage","dice_count":2,"die_sides":6,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"light","mechanic":"light.shining-smite","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"revealed","mechanic":"condition.shining-smite-revealed","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.shining-smite","light.shining-smite","condition.shining-smite-revealed","cast.bonus-action-after-hit","scaling.per-slot","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["immediately-after-hit-with-melee-weapon","immediately-after-hit-with-unarmed-strike"],"end_timing":"duration-or-concentration-ended","shining_smite":{"extra_damage_dice_count":2,"extra_damage_die_sides":6,"damage_type":"radiant","bright_light_feet":5,"attack_benefit":"attack-rolls-against-target-have-advantage","invisibility":"cannot-benefit-from-invisible-condition","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/shocking-grasp","name":"Shocking Grasp","source_id":"srd-5.2.1:content/spells/shocking-grasp","source_page":162,"scope":"complete-source-record","level":0,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.shocking-grasp","timing":"on-hit","dice_count":1,"die_sides":8,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"no-opportunity","mechanic":"restriction.shocking-grasp-opportunity-attacks","timing":"on-hit-until-start-of-targets-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.shocking-grasp","restriction.shocking-grasp-opportunity-attacks","attack.melee-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","shocking_grasp":{"attack":"melee-spell-attack","restriction":"cannot-make-opportunity-attacks","restriction_ends":"start-of-targets-next-turn","upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
		`{"id":"srd-5.2.1:content/spells/silence","name":"Silence","source_id":"srd-5.2.1:content/spells/silence","source_page":162,"scope":"complete-source-record","level":2,"school":"illusion","lists":["bard","cleric","ranger"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"sphere","count":0,"eligibility":"area-centered-on-chosen-point"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"silence","mechanic":"zone.silence-no-sound","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"immunity","mechanic":"defense.silence-thunder-immunity","timing":"while-entirely-inside","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"deafened","mechanic":"condition.silence-deafened","timing":"while-entirely-inside","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["zone.silence-no-sound","defense.silence-thunder-immunity","condition.silence-deafened","restriction.silence-verbal-spells","cast.ritual","concentration.source-duration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","silence":{"radius_feet":20,"sound":"none-created-within-or-passes-through","immunity_targets":["creature","object"],"immunity_damage_type":"thunder","inside_requirement":"entirely-inside-sphere","condition":"deafened","verbal_casting":"casting-spell-with-verbal-component-impossible"}}}`,
		`{"id":"srd-5.2.1:content/spells/silent-image","name":"Silent Image","source_id":"srd-5.2.1:content/spells/silent-image","source_page":162,"scope":"complete-source-record","level":1,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-phenomenon-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of fleece","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"image","mechanic":"illusion.silent-image","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"movement","mechanic":"control.silent-image-magic-action","timing":"magic-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.silent-image","control.silent-image-magic-action","check.silent-image-investigation","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","silent_image":{"phenomena":["object","creature","other-visible-phenomenon"],"max_cube_feet":15,"sensory_limit":"purely-visual-without-sound-smell-or-other-effects","move_action":"magic-action","move_destination":"any-spot-within-range","move_appearance":"may-alter-appearance-so-movement-appears-natural","example":"moved-creature-image-appears-to-walk","physical_interaction":"reveals-illusion-because-things-pass-through","examine_action":"study","examine_check":"intelligence-investigation","examine_dc":"caster-spell-save-dc","discerned":"can-see-through-image"}}}`,
		`{"id":"srd-5.2.1:content/spells/simulacrum","name":"Simulacrum","source_id":"srd-5.2.1:content/spells/simulacrum","source_page":162,"scope":"complete-source-record","level":7,"school":"illusion","lists":["wizard"],"casting_time":"12-hours","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"beast-or-humanoid-within-10-feet-for-entire-casting"},"components":{"verbal":true,"somatic":true,"material":{"description":"powdered ruby worth 1,500+ GP, which the spell consumes","cost_cp":150000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"simulacrum","mechanic":"creation.simulacrum","timing":"on-completion","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["creation.simulacrum","repair.simulacrum","cast.twelve-hours","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":43200,"casting_triggers":[],"end_timing":"dispelled-recast-or-simulacrum-drops-to-0-hit-points","simulacrum":{"source_pages":[162,163],"creature_types":["beast","humanoid"],"casting_distance_feet":10,"completion":"touch-creature-and-same-size-pile-of-ice-or-snow","result_kind":"creature","statistics":"original-creatures-statistics-at-casting","result_creature_type":"construct","hit_point_maximum_divisor":2,"forbidden_casting":"cannot-cast-this-spell","attitude":"friendly-to-caster-and-designated-creatures","control":"obeys-casters-commands-and-acts-on-casters-turn","limits":["cannot-gain-levels","cannot-take-short-or-long-rests"],"repair_requirement":"caster-repairs-during-long-rest","repair_gp_per_hit_point":100,"repair_distance_feet":5,"destruction":"at-0-hit-points-reverts-to-snow-and-melts","recast":"existing-simulacrum-from-this-spell-instantly-destroyed"}}}`,
		`{"id":"srd-5.2.1:content/spells/sleep","name":"Sleep","source_id":"srd-5.2.1:content/spells/sleep","source_page":163,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"sphere","count":0,"eligibility":"chosen-creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of sand or rose petals","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"drowsy","mechanic":"condition.sleep-incapacitated","timing":"on-failed-save-until-end-of-targets-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"asleep","mechanic":"condition.sleep-unconscious","timing":"on-second-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.sleep-incapacitated","condition.sleep-unconscious","save.wisdom-sleep-twice","ending.sleep-damage-or-shake","immunity.sleep-automatic-success","concentration.source-duration","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-damage-or-shaken","sleep":{"radius_feet":5,"save":"wisdom","first_condition":"incapacitated","first_condition_ends":"end-of-its-next-turn","repeat_timing":"when-first-condition-ends","second_failure_condition":"unconscious","ends_on_target":["target-takes-damage","creature-within-5-feet-uses-action-to-shake-it"],"wake_distance_feet":5,"automatic_success":["creatures-that-do-not-sleep","immunity-to-exhaustion"],"no_sleep_example":"elves"}}}`,
		`{"id":"srd-5.2.1:content/spells/sleet-storm","name":"Sleet Storm","source_id":"srd-5.2.1:content/spells/sleet-storm","source_page":163,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"cylinder","count":0,"eligibility":"area-centered-on-chosen-point"},"components":{"verbal":true,"somatic":true,"material":{"description":"a miniature umbrella","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"obscurement","mechanic":"terrain.sleet-storm","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"prone","mechanic":"condition.sleet-storm-prone","timing":"enter-first-time-on-turn-or-start-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.sleet-storm","condition.sleet-storm-prone","save.dexterity-sleet-storm","ending.sleet-storm-concentration-loss","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","sleet_storm":{"height_feet":40,"radius_feet":20,"obscurement":"heavily-obscured","flames":"exposed-flames-doused","ground":"difficult-terrain","save_triggers":["enters-cylinder-first-time-on-turn","starts-turn-there"],"save":"dexterity","failure":["prone","loses-concentration"]}}}`,
		`{"id":"srd-5.2.1:content/spells/slow","name":"Slow","source_id":"srd-5.2.1:content/spells/slow","source_page":163,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"up-to-count","count":6,"eligibility":"chosen-creatures-in-cube-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of molasses","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"slow","mechanic":"debuff.slow","timing":"on-failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["debuff.slow","save.wisdom-negates","save.slow-repeat","failure.slow-somatic-spell","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-successful-save","slow":{"max_targets":6,"cube_feet":40,"save":"wisdom","speed_divisor":2,"armor_class_penalty":2,"dexterity_save_penalty":2,"reactions":"cannot-take-reactions","turn_limit":"action-or-bonus-action-not-both","attack_limit":1,"somatic_failure_percent":25,"somatic_failure_cause":"gestures-too-slow","repeat_timing":"end-of-each-of-its-turns","repeat_success":"ends-spell-on-self"}}}`,
		`{"id":"srd-5.2.1:content/spells/sorcerous-burst","name":"Sorcerous Burst","source_id":"srd-5.2.1:content/spells/sorcerous-burst","source_page":163,"scope":"complete-source-record","level":0,"school":"evocation","lists":["sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.sorcerous-burst","timing":"on-hit","dice_count":1,"die_sides":8,"damage_type":"chosen","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.sorcerous-burst","attack.ranged-spell","damage.sorcerous-burst-extra-die","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","sorcerous_burst":{"attack":"ranged-spell-attack","target_kinds":["creature","object"],"damage_type_choices":["acid","cold","fire","lightning","poison","psychic","thunder"],"choose_count":1,"extra_die_trigger_roll":8,"extra_die_sides":8,"extra_die_limit":"spellcasting-ability-modifier","upgrades":[{"level":5,"dice_count":2,"die_sides":8},{"level":11,"dice_count":3,"die_sides":8},{"level":17,"dice_count":4,"die_sides":8}]}}}`,
		`{"id":"srd-5.2.1:content/spells/spare-the-dying","name":"Spare the Dying","source_id":"srd-5.2.1:content/spells/spare-the-dying","source_page":163,"scope":"complete-source-record","level":0,"school":"necromancy","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"feet","feet":15},"target":{"shape":"single","count":1,"eligibility":"creature-with-0-hit-points-not-dead"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-range","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":["stabilize"]},"effects":[{"id":"stabilize","mechanic":"healing.spare-the-dying-stable","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.spare-the-dying-stable","scaling.cantrip-range-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","spare_the_dying":{"required_hit_points":0,"requirement":"not-dead","result":"creature-becomes-stable","base_range_feet":15,"range_upgrades":[{"level":5,"feet":30},{"level":11,"feet":60},{"level":17,"feet":120}]}}}`,
		`{"id":"srd-5.2.1:content/spells/speak-with-animals","name":"Speak with Animals","source_id":"srd-5.2.1:content/spells/speak-with-animals","source_page":164,"scope":"complete-source-record","level":1,"school":"divination","lists":["bard","druid","ranger","warlock"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"speech","mechanic":"utility.speak-with-animals","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.speak-with-animals","cast.ritual"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","speak_with_animals":{"creature_type":"beast","communication":"comprehend-and-verbally-communicate","skill_options":"any-influence-action-skill-option","typical_topics":"little-beyond-survival-or-companionship","minimum_information":["nearby-locations","nearby-monsters","perceptions-within-past-day"],"perception_window_days":1}}}`,
	}
	batch := loadSpells219(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Shillelagh", "Shining Smite", "Shocking Grasp", "Silence", "Silent Image", "Simulacrum", "Sleep", "Sleet Storm", "Slow", "Sorcerous Burst", "Spare the Dying", "Speak with Animals"}
	pages := []int{162, 162, 162, 162, 162, 162, 163, 163, 163, 163, 163, 164}
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

func TestSpells219CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells219(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t), loadSpells216(t), loadSpells217(t), loadSpells219(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 288 {
		t.Fatalf("expected 288 records, got %d", len(all.Records))
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
func TestSpells219MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_219.json")
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

func TestSpells219JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_219.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"damage_dice_count": 1`, `"damage_dice_count": 1, "damage_dice_count": 2`},
		{`"damage_dice_count": 1`, `"damage_dice_count": 1, "DAMAGE_DICE_COUNT": 2`},
		{`"damage_dice_count": 1`, `"damage_dice_count": 1, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/shillelagh"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "weapon-die"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "12-hours"`, `"casting_time": "24-hours"`},
		{`"casting_seconds": 43200`, `"casting_seconds": 3600`},
		{`"casting_time": "bonus-action"`, `"casting_time": "action"`},
		{`"feet": 120`, `"feet": 121`},
		{`"cost_cp": 0`, `"cost_cp": 1`},
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

func TestSpells219InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"shillelagh die table", 0, func(s *SpellRecord) { s.SourceDetails.Shillelagh.DieUpgrades[2].DieSides = 8 }},
		{"shillelagh scaling kind", 0, func(s *SpellRecord) { s.Scaling.Kind = "cantrip-range" }},
		{"smite trigger", 1, func(s *SpellRecord) {
			s.SourceDetails.CastingTriggers = []string{"immediately-after-hit-with-melee-weapon"}
		}},
		{"smite dice", 1, func(s *SpellRecord) { s.Effects[0].DiceCount = 3 }},
		{"grasp attack", 2, func(s *SpellRecord) { s.SourceDetails.ShockingGrasp.Attack = "ranged-spell-attack" }},
		{"grasp upgrade", 2, func(s *SpellRecord) { s.SourceDetails.ShockingGrasp.Upgrades[0].DieSides = 6 }},
		{"silence radius", 3, func(s *SpellRecord) { s.SourceDetails.Silence.RadiusFeet = 10 }},
		{"silence ritual", 3, func(s *SpellRecord) { *s.Ritual = false }},
		{"image cube", 4, func(s *SpellRecord) { s.SourceDetails.SilentImage.MaxCubeFeet = 20 }},
		{"image check", 4, func(s *SpellRecord) { s.SourceDetails.SilentImage.ExamineCheck = "wisdom-perception" }},
		{"simulacrum casting", 5, func(s *SpellRecord) { s.CastingTime = "1-hour" }},
		{"simulacrum repair", 5, func(s *SpellRecord) { s.SourceDetails.Simulacrum.RepairGPPerHitPoint = 10 }},
		{"sleep save", 6, func(s *SpellRecord) { s.SourceDetails.Sleep.Save = "constitution" }},
		{"sleep condition", 6, func(s *SpellRecord) { s.SourceDetails.Sleep.SecondFailureCondition = "prone" }},
		{"sleet height", 7, func(s *SpellRecord) { s.SourceDetails.SleetStorm.HeightFeet = 20 }},
		{"sleet concentration", 7, func(s *SpellRecord) { *s.Concentration = false }},
		{"slow targets", 8, func(s *SpellRecord) { s.Target.Count = 5 }},
		{"slow percent", 8, func(s *SpellRecord) { s.SourceDetails.Slow.SomaticFailurePercent = 50 }},
		{"burst types", 9, func(s *SpellRecord) {
			s.SourceDetails.SorcerousBurst.DamageTypeChoices = s.SourceDetails.SorcerousBurst.DamageTypeChoices[:6]
		}},
		{"burst trigger", 9, func(s *SpellRecord) { s.SourceDetails.SorcerousBurst.ExtraDieTriggerRoll = 7 }},
		{"spare range table", 10, func(s *SpellRecord) { s.SourceDetails.SpareTheDying.RangeUpgrades[1].Feet = 45 }},
		{"spare scaling kind", 10, func(s *SpellRecord) { s.Scaling.Kind = "cantrip-upgrade" }},
		{"animals window", 11, func(s *SpellRecord) { s.SourceDetails.SpeakWithAnimals.PerceptionWindowDays = 7 }},
		{"animals duration", 11, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.SpeakWithAnimals = &SpellSpeakWithAnimals{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells219(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells219RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_219.json")
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

func TestSpells219WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_219.json")
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
		`"damage_dice_count": 1, "damage_dice_count": 2`,
		`"damage_dice_count": 1, "DAMAGE_DICE_COUNT": 2`,
		`"damage_dice_count": 1, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"damage_dice_count": 1`), []byte(replacement), 1)
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
func TestSpells219PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_219.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"damage_dice_count", 1, 1},
		{"damage_die_sides", 1, 8},
		{"level", 2, 5},
		{"dice_count", 5, 1},
		{"die_sides", 5, 10},
		{"level", 3, 11},
		{"dice_count", 6, 1},
		{"die_sides", 6, 12},
		{"level", 4, 17},
		{"dice_count", 7, 2},
		{"die_sides", 7, 6},
		{"extra_damage_dice_count", 1, 2},
		{"extra_damage_die_sides", 1, 6},
		{"bright_light_feet", 1, 5},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 6},
		{"level", 7, 5},
		{"dice_count", 15, 2},
		{"die_sides", 15, 8},
		{"level", 8, 11},
		{"dice_count", 16, 3},
		{"die_sides", 16, 8},
		{"level", 9, 17},
		{"dice_count", 17, 4},
		{"die_sides", 17, 8},
		{"radius_feet", 1, 20},
		{"max_cube_feet", 1, 15},
		{"casting_distance_feet", 1, 10},
		{"hit_point_maximum_divisor", 1, 2},
		{"repair_gp_per_hit_point", 1, 100},
		{"repair_distance_feet", 1, 5},
		{"radius_feet", 2, 5},
		{"wake_distance_feet", 1, 5},
		{"height_feet", 1, 40},
		{"radius_feet", 3, 20},
		{"max_targets", 1, 6},
		{"cube_feet", 1, 40},
		{"speed_divisor", 1, 2},
		{"armor_class_penalty", 1, 2},
		{"dexterity_save_penalty", 1, 2},
		{"attack_limit", 1, 1},
		{"somatic_failure_percent", 1, 25},
		{"choose_count", 1, 1},
		{"extra_die_trigger_roll", 1, 8},
		{"extra_die_sides", 1, 8},
		{"level", 17, 5},
		{"dice_count", 37, 2},
		{"die_sides", 37, 8},
		{"level", 18, 11},
		{"dice_count", 38, 3},
		{"die_sides", 38, 8},
		{"level", 19, 17},
		{"dice_count", 39, 4},
		{"die_sides", 39, 8},
		{"required_hit_points", 1, 0},
		{"base_range_feet", 1, 15},
		{"level", 21, 5},
		{"feet", 12, 30},
		{"level", 22, 11},
		{"feet", 13, 60},
		{"level", 23, 17},
		{"feet", 14, 120},
		{"perception_window_days", 1, 1},
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
// inventory assigns to catalogue issue #219.
func TestSpells219MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 219 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells219(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells219ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_219.json")
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
	want := "SRD 5.2.1 PDF pp.162-164 adapted into twelve complete offline spell families for issue #219, including page-break continuation text for Simulacrum (pp.162-163). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells219PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.SpeakWithAnimals = loadSpells219(t).Records[11].SourceDetails.SpeakWithAnimals
	if r.Validate() == nil {
		t.Fatal("batch 219 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells219RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells219(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell219Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell219Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
