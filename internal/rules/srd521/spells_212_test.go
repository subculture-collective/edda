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

func loadSpells212(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_212.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.145-147 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells212SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/locate-object","name":"Locate Object","source_id":"srd-5.2.1:content/spells/locate-object","source_page":145,"scope":"complete-source-record","level":2,"school":"divination","lists":["bard","cleric","druid","paladin","ranger","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a forked twig","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"locate","mechanic":"detect.locate-object","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["detect.locate-object","blocker.locate-object-lead","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","locate_object":{"query":"describe-or-name-familiar-object","senses":"direction-to-objects-location","sense_radius_feet":1000,"moving_target":"know-direction-of-motion","specific_requirement":"seen-up-close","close_distance_feet":30,"close_sightings_minimum":1,"alternative":"nearest-object-of-particular-kind","kind_examples":["apparel","jewelry","furniture","tool","weapon"],"lead_blocker":"any-thickness-of-lead-blocking-direct-path"}}}`,
		`{"id":"srd-5.2.1:content/spells/longstrider","name":"Longstrider","source_id":"srd-5.2.1:content/spells/longstrider","source_page":145,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["bard","druid","ranger","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of dirt","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["speed"]},"effects":[{"id":"speed","mechanic":"speed.longstrider","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":10,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["speed.longstrider","scaling.targets-per-slot","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","longstrider":{"target":"creature-touched","speed_increase_feet":10,"base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/mage-armor","name":"Mage Armor","source_id":"srd-5.2.1:content/spells/mage-armor","source_page":145,"scope":"complete-source-record","level":1,"school":"abjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-unarmored-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a piece of cured leather","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"armor","mechanic":"defense.mage-armor-base-armor-class","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.mage-armor-base-armor-class","ending.mage-armor-don-armor","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-target-dons-armor","mage_armor":{"target":"willing-creature-not-wearing-armor","base_armor_class":13,"armor_class_modifier":"target-dexterity-modifier","early_end":"target-dons-armor"}}}`,
		`{"id":"srd-5.2.1:content/spells/mage-hand","name":"Mage Hand","source_id":"srd-5.2.1:content/spells/mage-hand","source_page":145,"scope":"complete-source-record","level":0,"school":"conjuration","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"hand","mechanic":"summon.mage-hand","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"control","mechanic":"control.mage-hand-magic-action","timing":"magic-action-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.mage-hand","control.mage-hand-magic-action","ending.mage-hand-distance-or-recast"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-distance-or-cast-again","mage_hand":{"hand":"spectral-floating-hand","vanish_distance_feet":30,"vanish_triggers":["ever-farther-than-vanish-distance-from-caster","spell-cast-again"],"uses":["manipulate-object","open-unlocked-door-or-container","stow-or-retrieve-item-from-open-container","pour-contents-out-of-vial"],"later_action":"magic-action","later_move_feet":30,"forbidden":["attack","activate-magic-items","carry-more-than-weight-limit"],"carry_limit_pounds":10}}}`,
		`{"id":"srd-5.2.1:content/spells/magic-circle","name":"Magic Circle","source_id":"srd-5.2.1:content/spells/magic-circle","source_page":145,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["cleric","paladin","warlock","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":10},"target":{"shape":"cylinder","count":0,"eligibility":"circle-on-visible-ground-point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"salt and powdered silver worth 100+ GP, which the spell consumes","cost_cp":10000,"consumed":true}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-duration","above_level":3,"dice_count":0,"die_sides":0,"effect_ids":["entry","attacks","protection"]},"effects":[{"id":"entry","mechanic":"barrier.magic-circle-entry","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"attacks","mechanic":"disadvantage.magic-circle-attacks","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"protection","mechanic":"protection.magic-circle-possession-and-conditions","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.magic-circle-entry","disadvantage.magic-circle-attacks","protection.magic-circle-possession-and-conditions","save.magic-circle-charisma-travel","variant.magic-circle-reverse","scaling.duration-per-slot","cast.one-minute","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","magic_circle":{"source_pages":[145,146],"cylinder_radius_feet":10,"cylinder_height_feet":20,"center":"point-on-ground-visible-in-range","runes":"glowing-where-cylinder-intersects-floor-or-other-surface","creature_types":["celestial","elemental","fey","fiend","undead"],"type_choice":"one-or-more","entry_rule":"cannot-willingly-enter-by-nonmagical-means","travel_methods":["teleportation","interplanar-travel"],"travel_save":"charisma","travel_save_result":"must-succeed-before-using-travel-to-enter","attack_rule":"disadvantage-on-attack-rolls-against-targets-within","protections":["cannot-be-possessed-by-creature","cannot-gain-charmed-condition-from-creature","cannot-gain-frightened-condition-from-creature"],"reverse_option":"each-casting-may-reverse-direction","reverse_effect":"prevents-chosen-type-leaving-and-protects-targets-outside","additional_duration_seconds_per_slot":3600}}}`,
		`{"id":"srd-5.2.1:content/spells/magic-jar","name":"Magic Jar","source_id":"srd-5.2.1:content/spells/magic-jar","source_page":146,"scope":"complete-source-record","level":6,"school":"necromancy","lists":["wizard"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a gem, crystal, or reliquary worth 500+ GP","cost_cp":50000,"consumed":false}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"soul","mechanic":"soul.magic-jar-transfer","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"possession","mechanic":"possession.magic-jar","timing":"project-soul-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"return","mechanic":"return.magic-jar","timing":"magic-action-from-host","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["soul.magic-jar-transfer","possession.magic-jar","return.magic-jar","save.magic-jar-charisma","ending.magic-jar-soul-return-or-death","cast.one-minute","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"dispelled-container-destroyed-or-caster-returns","magic_jar":{"body_state":"catatonic","soul_destination":"container-used-as-material-component","container_awareness":"as-if-in-containers-space","container_restrictions":["cannot-move","cannot-take-reactions"],"container_only_action":"project-soul-out-of-container","projection_distance_feet":100,"projection_choices":["return-to-living-body-and-end-spell","attempt-to-possess-humanoid-body"],"possession_target":"humanoid-visible-within-distance","possession_distance_feet":100,"possession_wards":["protection-from-evil-and-good","magic-circle"],"possession_save":"charisma","possession_failure":"caster-soul-enters-body-and-target-soul-trapped-in-container","possession_success":"target-resists-and-cannot-be-targeted-again","possession_retry_hours":24,"host_control":"caster-controls-possessed-body","replaced_statistics":["hit-points","hit-point-dice","strength","dexterity","constitution","speed","senses"],"other_statistics":"caster-keeps-own-game-statistics","trapped_soul":["perceives-from-container-with-own-senses","cannot-move","incapacitated"],"return_action":"magic-action","return_distance_feet":100,"return_result":"caster-returns-to-container-and-host-soul-returns-to-body","host_death":"host-creature-dies-and-caster-makes-charisma-save","host_death_save_dc":"casters-own-spellcasting-dc","host_death_success":"return-to-container-if-within-distance","host_death_otherwise":"caster-dies","container_end_triggers":["container-destroyed","spell-ends"],"container_end_result":"caster-soul-returns-to-body","caster_death_conditions":["body-more-than-100-feet-away","body-dead"],"trapped_soul_on_destruction":"returns-to-its-body-if-body-alive-and-within-distance","trapped_soul_otherwise":"creature-dies","spell_end_container":"container-destroyed"}}}`,
		`{"id":"srd-5.2.1:content/spells/magic-missile","name":"Magic Missile","source_id":"srd-5.2.1:content/spells/magic-missile","source_page":146,"scope":"complete-source-record","level":1,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"darts","count":3,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-count","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["dart"]},"effects":[{"id":"dart","mechanic":"damage.magic-missile-dart","timing":"each-dart-simultaneously","dice_count":1,"die_sides":4,"damage_type":"force","bonus":1,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.magic-missile-dart","scaling.per-slot-dart"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","magic_missile":{"dart_count":3,"dart_dice_count":1,"dart_die_sides":4,"dart_bonus":1,"damage_type":"force","strike_timing":"all-darts-strike-simultaneously","distribution":["one-creature","several-creatures"],"additional_darts_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/magic-mouth","name":"Magic Mouth","source_id":"srd-5.2.1:content/spells/magic-mouth","source_page":146,"scope":"complete-source-record","level":2,"school":"illusion","lists":["bard","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-unworn-uncarried-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"jade dust worth 10+ GP, which the spell consumes","cost_cp":1000,"consumed":true}},"ritual":true,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"message","mechanic":"message.magic-mouth","timing":"when-trigger-occurs","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["message.magic-mouth","choice.magic-mouth-trigger-and-repeat","cast.one-minute","cast.ritual","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"dispelled-or-after-message-if-chosen","magic_mouth":{"message_max_words":25,"message_delivery_max_minutes":10,"trigger_definition":"circumstance-determined-by-caster","mouth":"magical-mouth-appears-on-object","recitation":"casters-voice-at-same-volume-spoken","mouth_placement":"appears-at-existing-mouth-or-mouthlike-feature","mouth_example":"mouth-of-statue","after_delivery_choices":["spell-ends-after-delivering-message","remains-and-repeats-whenever-trigger-occurs"],"trigger_detail":"as-general-or-as-detailed-as-caster-likes","trigger_basis":["visual-conditions","audible-conditions"],"trigger_range_feet":30,"trigger_examples":["any-creature-moves-within-30-feet-of-object","silver-bell-rings-within-30-feet-of-object"]}}}`,
		`{"id":"srd-5.2.1:content/spells/magic-weapon","name":"Magic Weapon","source_id":"srd-5.2.1:content/spells/magic-weapon","source_page":146,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["paladin","ranger","sorcerer","wizard"],"casting_time":"bonus-action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"nonmagical-weapon-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-bonus-table","above_level":2,"dice_count":0,"die_sides":0,"effect_ids":["weapon"]},"effects":[{"id":"weapon","mechanic":"enchantment.magic-weapon","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["enchantment.magic-weapon","ending.magic-weapon-recast","scaling.bonus-by-slot"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-cast-again","magic_weapon":{"source_pages":[146,147],"target":"nonmagical-weapon-touched","becomes":"magic-weapon","attack_roll_bonus":1,"damage_roll_bonus":1,"early_end":"spell-cast-again","slot_bonuses":[{"min_slot":3,"max_slot":5,"bonus":2},{"min_slot":6,"max_slot":9,"bonus":3}]}}}`,
		`{"id":"srd-5.2.1:content/spells/magnificent-mansion","name":"Magnificent Mansion","source_id":"srd-5.2.1:content/spells/magnificent-mansion","source_page":147,"scope":"complete-source-record","level":7,"school":"conjuration","lists":["bard","wizard"],"casting_time":"1-minute","range":{"kind":"feet","feet":300},"target":{"shape":"created-effect","count":1,"eligibility":"door-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a miniature door worth 15+ GP","cost_cp":1500,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"door","mechanic":"create.magnificent-mansion-door","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dwelling","mechanic":"create.magnificent-mansion-dwelling","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"servants","mechanic":"summon.magnificent-mansion-servants","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["create.magnificent-mansion-door","create.magnificent-mansion-dwelling","summon.magnificent-mansion-servants","expulsion.magnificent-mansion","cast.one-minute","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","magnificent_mansion":{"door_appearance":"shimmering","door_width_feet":5,"door_height_feet":10,"door_destination":"extradimensional-dwelling","entrants":"caster-and-creatures-designated-when-cast","entry_requirement":"door-remains-open","door_control_action":"no-action-required","door_control_distance_feet":30,"closed_door":"imperceptible","foyer":"magnificent-foyer-with-numerous-chambers-beyond","atmosphere":["clean","fresh","warm"],"floor_plan":"any-floor-plan-caster-likes","max_contiguous_cubes":50,"cube_feet":10,"furnishing":"furnished-and-decorated-as-caster-chooses","banquet_courses":9,"banquet_people":100,"removed_objects":"dissipate-into-smoke","servant_count":100,"servant_appearance":"near-transparent-with-appearance-and-attire-determined-by-caster","servant_traits":["invulnerable","obey-casters-commands"],"servant_tasks":"any-task-a-human-could-perform","servant_task_examples":["fetch-things","clean","mend","fold-clothes","light-fires","serve-food","pour-wine"],"servant_limits":["cannot-attack","cannot-take-action-that-directly-harms-another-creature","cannot-leave-dwelling"],"end_expulsion":"creatures-and-objects-inside-expelled-into-unoccupied-spaces-nearest-entrance"}}}`,
		`{"id":"srd-5.2.1:content/spells/major-image","name":"Major Image","source_id":"srd-5.2.1:content/spells/major-image","source_page":147,"scope":"complete-source-record","level":3,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"created-effect","count":1,"eligibility":"visible-spot-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of fleece","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot-until-dispelled","above_level":3,"dice_count":0,"die_sides":0,"effect_ids":["image"]},"effects":[{"id":"image","mechanic":"illusion.major-image","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"movement","mechanic":"control.major-image-magic-action","timing":"magic-action-while-in-range","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.major-image","control.major-image-magic-action","check.major-image-investigation","scaling.until-dispelled-by-slot","concentration.source-duration","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","major_image":{"image_subjects":["object","creature","other-visible-phenomenon"],"max_cube_feet":20,"qualities":["seems-real","sounds","smells","temperature-appropriate-to-thing-depicted"],"limits":["cannot-deal-damage","cannot-cause-conditions"],"move_requirement":"caster-within-range-of-illusion","move_action":"magic-action","move_destination":"any-other-spot-within-range","movement_appearance":"appearance-can-be-altered-so-movement-appears-natural","sound_changes":"different-sounds-at-different-times-including-conversation","physical_interaction":"reveals-illusion-because-things-pass-through","discern_action":"study","discern_check":"intelligence-investigation","discern_dc":"caster-spell-save-dc","discerned":["can-see-through-image","other-sensory-qualities-become-faint"],"until_dispelled_min_slot":4,"until_dispelled_concentration":"not-required"}}}`,
		`{"id":"srd-5.2.1:content/spells/mass-cure-wounds","name":"Mass Cure Wounds","source_id":"srd-5.2.1:content/spells/mass-cure-wounds","source_page":147,"scope":"complete-source-record","level":5,"school":"abjuration","lists":["bard","cleric","druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"up-to-count","count":6,"eligibility":"creatures-in-sphere-on-visible-point"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":8,"effect_ids":["healing"]},"effects":[{"id":"healing","mechanic":"healing.mass-cure-wounds","timing":"immediate","dice_count":5,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["healing.mass-cure-wounds","scaling.per-slot"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","mass_cure_wounds":{"center":"point-visible-within-range","sphere_radius_feet":30,"max_targets":6,"healing_dice_count":5,"healing_die_sides":8,"healing_adds":"spellcasting-ability-modifier","slot_dice_count":1,"slot_die_sides":8}}}`,
	}
	batch := loadSpells212(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Locate Object", "Longstrider", "Mage Armor", "Mage Hand", "Magic Circle", "Magic Jar", "Magic Missile", "Magic Mouth", "Magic Weapon", "Magnificent Mansion", "Major Image", "Mass Cure Wounds"}
	pages := []int{145, 145, 145, 145, 145, 146, 146, 146, 146, 147, 147, 147}
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

func TestSpells212CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells212(t)
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
func TestSpells212MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_212.json")
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

func TestSpells212JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_212.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"sense_radius_feet": 1000`, `"sense_radius_feet": 1000, "sense_radius_feet": 1001`},
		{`"sense_radius_feet": 1000`, `"sense_radius_feet": 1000, "SENSE_RADIUS_FEET": 1001`},
		{`"sense_radius_feet": 1000`, `"sense_radius_feet": 1000, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/locate-object"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "speed"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
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

func TestSpells212InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"object radius", 0, func(s *SpellRecord) { s.SourceDetails.LocateObject.SenseRadiusFeet = 500 }},
		{"object duration", 0, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"longstrider speed", 1, func(s *SpellRecord) { s.Effects[0].Bonus = 15 }},
		{"mage armor base", 2, func(s *SpellRecord) { s.SourceDetails.MageArmor.BaseArmorClass = 15 }},
		{"mage armor target", 2, func(s *SpellRecord) { s.Target.Eligibility = "creature-touch" }},
		{"mage hand weight", 3, func(s *SpellRecord) { s.SourceDetails.MageHand.CarryLimitPounds = 20 }},
		{"mage hand slots", 3, func(s *SpellRecord) { s.SlotLevels = []int{1} }},
		{"circle pages", 4, func(s *SpellRecord) { s.SourceDetails.MagicCircle.SourcePages = []int{145} }},
		{"circle types", 4, func(s *SpellRecord) {
			s.SourceDetails.MagicCircle.CreatureTypes = s.SourceDetails.MagicCircle.CreatureTypes[:4]
		}},
		{"circle slot duration", 4, func(s *SpellRecord) { s.SourceDetails.MagicCircle.AdditionalDurationSecondsPerSlot = 0 }},
		{"jar statistics", 5, func(s *SpellRecord) {
			s.SourceDetails.MagicJar.ReplacedStatistics = s.SourceDetails.MagicJar.ReplacedStatistics[:6]
		}},
		{"jar save", 5, func(s *SpellRecord) { s.SourceDetails.MagicJar.PossessionSave = "wisdom" }},
		{"missile count", 6, func(s *SpellRecord) { s.Target.Count = 4 }},
		{"missile slot darts", 6, func(s *SpellRecord) { s.SourceDetails.MagicMissile.AdditionalDartsPerSlot = 0 }},
		{"mouth words", 7, func(s *SpellRecord) { s.SourceDetails.MagicMouth.MessageMaxWords = 30 }},
		{"mouth consumption", 7, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"weapon table", 8, func(s *SpellRecord) { s.SourceDetails.MagicWeapon.SlotBonuses[1].Bonus = 4 }},
		{"weapon scaling kind", 8, func(s *SpellRecord) { s.Scaling.Kind = "slot-bonus" }},
		{"mansion servants", 9, func(s *SpellRecord) { s.SourceDetails.MagnificentMansion.ServantCount = 50 }},
		{"mansion cubes", 9, func(s *SpellRecord) { s.SourceDetails.MagnificentMansion.MaxContiguousCubes = 100 }},
		{"image slot", 10, func(s *SpellRecord) { s.SourceDetails.MajorImage.UntilDispelledMinSlot = 5 }},
		{"image scaling effect", 10, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"unknown"} }},
		{"cure targets", 11, func(s *SpellRecord) { s.SourceDetails.MassCureWounds.MaxTargets = 8 }},
		{"cure modifier", 11, func(s *SpellRecord) { *s.Effects[0].AddCastingModifier = false }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.MassCureWounds = &SpellMassCureWounds{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells212(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells212RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_212.json")
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

func TestSpells212WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_212.json")
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
		`"sense_radius_feet": 1000, "sense_radius_feet": 1001`,
		`"sense_radius_feet": 1000, "SENSE_RADIUS_FEET": 1001`,
		`"sense_radius_feet": 1000, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"sense_radius_feet": 1000`), []byte(replacement), 1)
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
func TestSpells212PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_212.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"sense_radius_feet", 1, 1000},
		{"close_distance_feet", 1, 30},
		{"close_sightings_minimum", 1, 1},
		{"speed_increase_feet", 1, 10},
		{"base_targets", 1, 1},
		{"additional_targets_per_slot", 1, 1},
		{"base_armor_class", 1, 13},
		{"vanish_distance_feet", 1, 30},
		{"later_move_feet", 1, 30},
		{"carry_limit_pounds", 1, 10},
		{"cylinder_radius_feet", 1, 10},
		{"cylinder_height_feet", 1, 20},
		{"additional_duration_seconds_per_slot", 1, 3600},
		{"projection_distance_feet", 1, 100},
		{"possession_distance_feet", 1, 100},
		{"possession_retry_hours", 1, 24},
		{"return_distance_feet", 1, 100},
		{"dart_count", 1, 3},
		{"dart_dice_count", 1, 1},
		{"dart_die_sides", 1, 4},
		{"dart_bonus", 1, 1},
		{"additional_darts_per_slot", 1, 1},
		{"message_max_words", 1, 25},
		{"message_delivery_max_minutes", 1, 10},
		{"trigger_range_feet", 1, 30},
		{"attack_roll_bonus", 1, 1},
		{"damage_roll_bonus", 1, 1},
		{"min_slot", 1, 3},
		{"min_slot", 2, 6},
		{"max_slot", 1, 5},
		{"max_slot", 2, 9},
		{"bonus", 1, 0},
		{"bonus", 2, 10},
		{"bonus", 3, 0},
		{"bonus", 4, 0},
		{"bonus", 5, 0},
		{"bonus", 6, 0},
		{"bonus", 7, 0},
		{"bonus", 8, 0},
		{"bonus", 9, 0},
		{"bonus", 10, 0},
		{"bonus", 11, 0},
		{"bonus", 12, 1},
		{"bonus", 13, 0},
		{"bonus", 14, 1},
		{"bonus", 15, 2},
		{"bonus", 16, 3},
		{"bonus", 17, 0},
		{"bonus", 18, 0},
		{"bonus", 19, 0},
		{"bonus", 20, 0},
		{"bonus", 21, 0},
		{"bonus", 22, 0},
		{"door_width_feet", 1, 5},
		{"door_height_feet", 1, 10},
		{"door_control_distance_feet", 1, 30},
		{"max_contiguous_cubes", 1, 50},
		{"cube_feet", 1, 10},
		{"banquet_courses", 1, 9},
		{"banquet_people", 1, 100},
		{"servant_count", 1, 100},
		{"max_cube_feet", 1, 20},
		{"until_dispelled_min_slot", 1, 4},
		{"sphere_radius_feet", 1, 30},
		{"max_targets", 1, 6},
		{"healing_dice_count", 1, 5},
		{"healing_die_sides", 1, 8},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
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
// inventory assigns to catalogue issue #212.
func TestSpells212MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 212 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells212(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells212ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_212.json")
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
	want := "SRD 5.2.1 PDF pp.145-147 adapted into twelve complete offline spell families for issue #212, including page-break continuation text for Magic Circle (pp.145-146) and Magic Weapon (pp.146-147). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells212PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.MassCureWounds = loadSpells212(t).Records[11].SourceDetails.MassCureWounds
	if r.Validate() == nil {
		t.Fatal("batch 212 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells212RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells212(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell212Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell212Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
