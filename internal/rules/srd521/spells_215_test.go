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

func loadSpells215(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_215.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.152-155 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells215SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/planar-ally","name":"Planar Ally","source_id":"srd-5.2.1:content/spells/planar-ally","source_page":152,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["cleric"],"casting_time":"10-minutes","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-space-in-range-for-loyal-creature"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"arrival","mechanic":"summon.planar-ally","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.planar-ally","negotiation.planar-ally-service","cast.ten-minutes"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"instantaneous","planar_ally":{"entity_requirement":"being-known-to-caster","entity_examples":["god","demon-prince","other-being-of-cosmic-power"],"creature_type_choices":["celestial","elemental","fiend"],"creature_loyalty":"loyal-to-entity","named_request":"may-speak-known-creature-name-but-gm-may-send-another","compulsion":"none","service_basis":"may-ask-service-for-payment-without-obligation","task_range":"simple-to-complex","task_examples":["fly-caster-across-chasm","help-fight-battle","spy-on-enemies","protect-during-dungeon-foray"],"bargaining_requirement":"caster-able-to-communicate-with-creature","payment_examples":["gold-or-magic-item-donation-to-allied-temple","living-sacrifice","gift-of-treasure","quest-undertaken-by-caster"],"payment_rates":[{"task_unit":"minute","gp_per_unit":100},{"task_unit":"hour","gp_per_unit":1000},{"task_unit":"day","gp_per_unit":10000}],"max_task_days":10,"gm_adjustment":"gm-can-adjust-payment-for-casting-circumstances","aligned_task_payment":"might-be-halved-or-waived","nonhazardous_payment":"typically-half-suggested-payment","dangerous_payment":"might-require-greater-gift","suicidal_tasks":"rarely-accepted","return_triggers":["task-completed","agreed-service-duration-expires"],"return_result":"returns-to-home-plane-after-reporting-to-caster-if-possible","no_agreement_result":"immediately-returns-to-home-plane"}}}`,
		`{"id":"srd-5.2.1:content/spells/planar-binding","name":"Planar Binding","source_id":"srd-5.2.1:content/spells/planar-binding","source_page":152,"scope":"complete-source-record","level":5,"school":"abjuration","lists":["bard","cleric","druid","warlock","wizard"],"casting_time":"1-hour","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"celestial-elemental-fey-or-fiend-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a jewel worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"slot-duration","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["binding"]},"effects":[{"id":"binding","mechanic":"binding.planar-binding","timing":"on-failed-charisma-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["binding.planar-binding","save.charisma-negates","duration.planar-binding-extends-source-spell","scaling.duration-by-slot","cast.one-hour","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"duration-ended","planar_binding":{"source_pages":[152,153],"creature_types":["celestial","elemental","fey","fiend"],"range_requirement":"within-range-for-entire-casting","typical_setup":"summoned-into-center-of-inverted-magic-circle-while-casting","save":"charisma","save_timing":"completion-of-casting","failure_result":"bound-to-serve-for-duration","source_spell_extension":"summoning-or-creating-spell-duration-extended-to-match","command_obedience":"follows-commands-to-best-of-ability","command_examples":["accompany-on-adventure","guard-location","deliver-message"],"hostile_behavior":"strives-to-twist-commands-to-own-objectives","completion_same_plane":"travels-to-caster-to-report","completion_other_plane":"returns-to-binding-place-until-spell-ends","slot_durations":[{"slot":6,"days":10},{"slot":7,"days":30},{"slot":8,"days":180},{"slot":9,"days":366}]}}}`,
		`{"id":"srd-5.2.1:content/spells/plane-shift","name":"Plane Shift","source_id":"srd-5.2.1:content/spells/plane-shift","source_page":153,"scope":"complete-source-record","level":7,"school":"conjuration","lists":["cleric","druid","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"selected-creatures","count":0,"eligibility":"caster-and-willing-creatures-linked-in-circle"},"components":{"verbal":true,"somatic":true,"material":{"description":"a forked, metal rod worth 250+ GP and attuned to a plane of existence","cost_cp":25000,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"transport","mechanic":"teleport.plane-shift","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.plane-shift","components.material"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","plane_shift":{"max_additional_creatures":8,"traveler_requirement":"willing-creatures-who-link-hands-in-circle","destination_kind":"different-plane-of-existence","general_destination":"general-terms-with-arrival-in-or-near-as-gm-determines","destination_examples":["city-on-elemental-plane-of-fire","palace-on-second-level-of-nine-hells"],"circle_destination":"known-sigil-sequence-of-teleportation-circle-on-another-plane","small_circle_overflow":"closest-unoccupied-spaces-next-to-circle"}}}`,
		`{"id":"srd-5.2.1:content/spells/plant-growth","name":"Plant Growth","source_id":"srd-5.2.1:content/spells/plant-growth","source_page":153,"scope":"complete-source-record","level":3,"school":"transmutation","lists":["bard","druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":150},"target":{"shape":"choice","count":1,"eligibility":"overgrowth-or-enrichment-point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"overgrowth","mechanic":"terrain.plant-growth-overgrowth","timing":"immediate-on-action-casting","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"enrichment","mechanic":"growth.plant-growth-enrichment","timing":"immediate-on-eight-hour-casting","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.plant-growth-overgrowth","growth.plant-growth-enrichment","cast.plant-growth-casting-time-choice"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","plant_growth":{"variants":[{"choice":"overgrowth","casting_time":"action","casting_seconds":0},{"choice":"enrichment","casting_time":"8-hours","casting_seconds":28800}],"overgrowth_radius_feet":100,"overgrowth_shape":"sphere","overgrowth_plants":"all-normal-plants-become-thick-and-overgrown","movement_cost_feet":4,"movement_unit_feet":1,"overgrowth_exclusions":"caster-can-exclude-areas-of-any-size","enrichment_radius_feet":2640,"enrichment_days":365,"enrichment_yield_multiplier":2,"enrichment_limit_per_year":1}}}`,
		`{"id":"srd-5.2.1:content/spells/poison-spray","name":"Poison Spray","source_id":"srd-5.2.1:content/spells/poison-spray","source_page":153,"scope":"complete-source-record","level":0,"school":"necromancy","lists":["druid","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"cantrip-upgrade","above_level":0,"dice_count":1,"die_sides":12,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.poison-spray","timing":"on-hit","dice_count":1,"die_sides":12,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.poison-spray","attack.ranged-spell","scaling.cantrip-upgrade"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","poison_spray":{"appearance":"toxic-mist","attack":"ranged-spell-attack","upgrades":[{"level":5,"dice_count":2,"die_sides":12},{"level":11,"dice_count":3,"die_sides":12},{"level":17,"dice_count":4,"die_sides":12}]}}}`,
		`{"id":"srd-5.2.1:content/spells/polymorph","name":"Polymorph","source_id":"srd-5.2.1:content/spells/polymorph","source_page":153,"scope":"complete-source-record","level":4,"school":"transmutation","lists":["bard","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a caterpillar cocoon","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"form","mechanic":"transformation.polymorph","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"temporary-hit-points","mechanic":"healing.polymorph-temporary-hit-points","timing":"with-transformation","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.polymorph","healing.polymorph-temporary-hit-points","save.wisdom-negates","ending.polymorph-no-temporary-hit-points","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-no-temporary-hit-points","polymorph":{"source_pages":[153,154],"save":"wisdom","form_creature_type":"beast","challenge_limit":"challenge-rating-at-most-targets","no_challenge_rating_limit":"targets-level","statistics":"replaced-by-chosen-beast-stat-block","retained":["alignment","personality","creature-type","hit-points","hit-point-dice"],"stat_block_reference":"animals-section-of-monsters","temporary_hit_points":"equal-to-beast-form-hit-points","temporary_hit_points_end":"vanish-if-any-remain-when-spell-ends","early_end":"ends-on-target-with-no-temporary-hit-points-left","action_limit":"limited-by-new-form-anatomy","forbidden":["speaking","casting-spells"],"gear":"melds-into-new-form","gear_use":"cannot-use-or-benefit-from-equipment"}}}`,
		`{"id":"srd-5.2.1:content/spells/power-word-heal","name":"Power Word Heal","source_id":"srd-5.2.1:content/spells/power-word-heal","source_page":154,"scope":"complete-source-record","level":9,"school":"enchantment","lists":["bard","cleric"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"healing","mechanic":"healing.power-word-heal-all-hit-points","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"conditions","mechanic":"ending.power-word-heal-conditions","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.power-word-heal-all-hit-points","ending.power-word-heal-conditions","reaction.power-word-heal-stand"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","power_word_heal":{"hit_points":"target-regains-all-hit-points","conditions_ended":["charmed","frightened","paralyzed","poisoned","stunned"],"prone_option":"prone-target-can-use-reaction-to-stand-up"}}}`,
		`{"id":"srd-5.2.1:content/spells/power-word-kill","name":"Power Word Kill","source_id":"srd-5.2.1:content/spells/power-word-kill","source_page":154,"scope":"complete-source-record","level":9,"school":"enchantment","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"death","mechanic":"death.power-word-kill","timing":"immediate-at-or-below-threshold","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.power-word-kill","timing":"immediate-above-threshold","dice_count":12,"die_sides":12,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["death.power-word-kill","damage.power-word-kill"],"slot_levels":[9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","power_word_kill":{"hit_point_threshold":100,"at_or_below_result":"target-dies","above_dice_count":12,"above_die_sides":12,"above_damage_type":"psychic"}}}`,
		`{"id":"srd-5.2.1:content/spells/power-word-stun","name":"Power Word Stun","source_id":"srd-5.2.1:content/spells/power-word-stun","source_page":154,"scope":"complete-source-record","level":8,"school":"enchantment","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"stun","mechanic":"condition.power-word-stun-stunned","timing":"immediate-at-or-below-threshold","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"speed","mechanic":"movement.power-word-stun-speed-zero","timing":"above-threshold-until-start-of-casters-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.power-word-stun-stunned","movement.power-word-stun-speed-zero","save.constitution-end-on-success"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","power_word_stun":{"hit_point_threshold":150,"at_or_below_condition":"stunned","above_speed_feet":0,"above_duration":"until-start-of-casters-next-turn","repeat_save":"constitution","repeat_save_timing":"end-of-each-of-targets-turns","repeat_success":"ends-condition-on-self"}}}`,
		`{"id":"srd-5.2.1:content/spells/prayer-of-healing","name":"Prayer of Healing","source_id":"srd-5.2.1:content/spells/prayer-of-healing","source_page":154,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["cleric","paladin"],"casting_time":"10-minutes","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":5,"eligibility":"chosen-creatures-in-range-for-entire-casting"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":2,"dice_count":1,"die_sides":8,"effect_ids":["healing"]},"effects":[{"id":"rest","mechanic":"rest.prayer-of-healing-short-rest","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"healing","mechanic":"healing.prayer-of-healing","timing":"immediate","dice_count":2,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["rest.prayer-of-healing-short-rest","healing.prayer-of-healing","limit.prayer-of-healing-per-long-rest","scaling.per-slot","cast.ten-minutes"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"instantaneous","prayer_of_healing":{"max_creatures":5,"presence_requirement":"remain-within-range-for-entire-casting","rest_benefit":"short-rest","healing_dice_count":2,"healing_die_sides":8,"reuse_limit":"not-again-until-creature-finishes-long-rest","slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/prestidigitation","name":"Prestidigitation","source_id":"srd-5.2.1:content/spells/prestidigitation","source_page":154,"scope":"complete-source-record","level":0,"school":"transmutation","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":10},"target":{"shape":"choice","count":1,"eligibility":"one-chosen-minor-effect-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"magic","mechanic":"utility.prestidigitation","timing":"chosen-option","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["utility.prestidigitation","limit.prestidigitation-active-effects"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"up-to-duration","prestidigitation":{"options":[{"choice":"sensory-effect","effect":"instantaneous-harmless-sensory-effect","details":["shower-of-sparks","puff-of-wind","faint-musical-notes","odd-odor"],"lasts":"instantaneous"},{"choice":"fire-play","effect":"light-or-snuff-out-small-flame","details":["candle","torch","small-campfire"],"lasts":"instantaneous"},{"choice":"clean-or-soil","effect":"clean-or-soil-object","details":["no-larger-than-1-cubic-foot"],"lasts":"instantaneous"},{"choice":"minor-sensation","effect":"chill-warm-or-flavor-nonliving-material","details":["up-to-1-cubic-foot"],"lasts":"1-hour"},{"choice":"magic-mark","effect":"color-small-mark-or-symbol-appears","details":["on-object-or-surface"],"lasts":"1-hour"},{"choice":"minor-creation","effect":"nonmagical-trinket-or-illusory-image-fits-in-hand","details":["trinket-deals-no-damage","trinket-has-no-monetary-worth"],"lasts":"until-end-of-casters-next-turn"}],"choose_count":1,"max_active_non_instantaneous":3}}}`,
		`{"id":"srd-5.2.1:content/spells/prismatic-spray","name":"Prismatic Spray","source_id":"srd-5.2.1:content/spells/prismatic-spray","source_page":154,"scope":"complete-source-record","level":7,"school":"evocation","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"cone","count":0,"eligibility":"creatures-in-self-origin-cone"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"red","mechanic":"damage.prismatic-spray-red","timing":"on-ray-1","dice_count":12,"die_sides":6,"damage_type":"fire","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"orange","mechanic":"damage.prismatic-spray-orange","timing":"on-ray-2","dice_count":12,"die_sides":6,"damage_type":"acid","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"yellow","mechanic":"damage.prismatic-spray-yellow","timing":"on-ray-3","dice_count":12,"die_sides":6,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"green","mechanic":"damage.prismatic-spray-green","timing":"on-ray-4","dice_count":12,"die_sides":6,"damage_type":"poison","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"blue","mechanic":"damage.prismatic-spray-blue","timing":"on-ray-5","dice_count":12,"die_sides":6,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"indigo","mechanic":"condition.prismatic-spray-indigo","timing":"on-ray-6","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"violet","mechanic":"condition.prismatic-spray-violet","timing":"on-ray-7","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"special","mechanic":"roll.prismatic-spray-two-rays","timing":"on-ray-8","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.prismatic-spray-red","damage.prismatic-spray-orange","damage.prismatic-spray-yellow","damage.prismatic-spray-green","damage.prismatic-spray-blue","condition.prismatic-spray-indigo","condition.prismatic-spray-violet","roll.prismatic-spray-two-rays","save.dexterity-prismatic-spray"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","prismatic_spray":{"source_pages":[154,155],"cone_feet":60,"ray_count":8,"save":"dexterity","ray_die_sides":8,"damage_rays":[{"roll":1,"color":"red","dice_count":12,"die_sides":6,"damage_type":"fire","success":"half-damage"},{"roll":2,"color":"orange","dice_count":12,"die_sides":6,"damage_type":"acid","success":"half-damage"},{"roll":3,"color":"yellow","dice_count":12,"die_sides":6,"damage_type":"lightning","success":"half-damage"},{"roll":4,"color":"green","dice_count":12,"die_sides":6,"damage_type":"poison","success":"half-damage"},{"roll":5,"color":"blue","dice_count":12,"die_sides":6,"damage_type":"cold","success":"half-damage"}],"indigo_roll":6,"indigo_condition":"restrained","indigo_save":"constitution","indigo_save_timing":"end-of-each-of-its-turns","indigo_successes_to_end":3,"indigo_failures_to_petrify":3,"indigo_petrified_until":"freed-by-effect-like-greater-restoration","indigo_tracking":"need-not-be-consecutive-until-three-of-a-kind","violet_roll":7,"violet_condition":"blinded","violet_save":"wisdom","violet_save_timing":"start-of-casters-next-turn","violet_success":"condition-ends","violet_failure":"condition-ends-and-teleports-to-another-plane-gm-choice","special_roll":8,"special_result":"struck-by-two-rays","special_reroll":"roll-twice-rerolling-any-8"}}}`,
	}
	batch := loadSpells215(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Planar Ally", "Planar Binding", "Plane Shift", "Plant Growth", "Poison Spray", "Polymorph", "Power Word Heal", "Power Word Kill", "Power Word Stun", "Prayer of Healing", "Prestidigitation", "Prismatic Spray"}
	pages := []int{152, 152, 153, 153, 153, 153, 154, 154, 154, 154, 154, 154}
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

func TestSpells215CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells215(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t), loadSpells209(t), loadSpells210(t), loadSpells211(t), loadSpells212(t), loadSpells213(t), loadSpells214(t), loadSpells215(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 252 {
		t.Fatalf("expected 252 records, got %d", len(all.Records))
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
func TestSpells215MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_215.json")
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

func TestSpells215JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_215.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"max_task_days": 10`, `"max_task_days": 10, "max_task_days": 11`},
		{`"max_task_days": 10`, `"max_task_days": 10, "MAX_TASK_DAYS": 11`},
		{`"max_task_days": 10`, `"max_task_days": 10, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/planar-ally"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "binding"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "1-hour"`, `"casting_time": "8-hours"`},
		{`"casting_time": "10-minutes"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"casting_seconds": 3600`, `"casting_seconds": 60`},
		{`"casting_seconds": 600`, `"casting_seconds": 601`},
		{`"feet": 60`, `"feet": 61`},
		{`"cost_cp": 100000`, `"cost_cp": 100001`},
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

func TestSpells215InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"ally payment rate", 0, func(s *SpellRecord) { s.SourceDetails.PlanarAlly.PaymentRates[1].GPPerUnit = 100 }},
		{"ally types", 0, func(s *SpellRecord) {
			s.SourceDetails.PlanarAlly.CreatureTypeChoices = s.SourceDetails.PlanarAlly.CreatureTypeChoices[:2]
		}},
		{"binding table", 1, func(s *SpellRecord) { s.SourceDetails.PlanarBinding.SlotDurations[3].Days = 365 }},
		{"binding scaling kind", 1, func(s *SpellRecord) { s.Scaling.Kind = "slot-targets" }},
		{"shift travelers", 2, func(s *SpellRecord) { s.SourceDetails.PlaneShift.MaxAdditionalCreatures = 9 }},
		{"shift consumed", 2, func(s *SpellRecord) { *s.Components.Material.Consumed = true }},
		{"growth enrichment casting", 3, func(s *SpellRecord) { s.SourceDetails.PlantGrowth.Variants[1].CastingSeconds = 3600 }},
		{"growth yield", 3, func(s *SpellRecord) { s.SourceDetails.PlantGrowth.EnrichmentYieldMultiplier = 3 }},
		{"poison die", 4, func(s *SpellRecord) { s.Effects[0].DieSides = 10 }},
		{"poison upgrade", 4, func(s *SpellRecord) { s.SourceDetails.PoisonSpray.Upgrades[0].Level = 4 }},
		{"polymorph save", 5, func(s *SpellRecord) { s.SourceDetails.Polymorph.Save = "constitution" }},
		{"polymorph retained", 5, func(s *SpellRecord) { s.SourceDetails.Polymorph.Retained = s.SourceDetails.Polymorph.Retained[:4] }},
		{"heal somatic", 6, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"heal conditions", 6, func(s *SpellRecord) { s.SourceDetails.PowerWordHeal.ConditionsEnded[0] = "prone" }},
		{"kill threshold", 7, func(s *SpellRecord) { s.SourceDetails.PowerWordKill.HitPointThreshold = 150 }},
		{"kill dice", 7, func(s *SpellRecord) { s.Effects[1].DiceCount = 10 }},
		{"stun threshold", 8, func(s *SpellRecord) { s.SourceDetails.PowerWordStun.HitPointThreshold = 100 }},
		{"stun speed", 8, func(s *SpellRecord) { s.SourceDetails.PowerWordStun.AboveSpeedFeet = 5 }},
		{"prayer creatures", 9, func(s *SpellRecord) { s.Target.Count = 6 }},
		{"prayer scaling die", 9, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"prestidigitation active", 10, func(s *SpellRecord) { s.SourceDetails.Prestidigitation.MaxActiveNonInstantaneous = 4 }},
		{"prestidigitation options", 10, func(s *SpellRecord) {
			s.SourceDetails.Prestidigitation.Options = s.SourceDetails.Prestidigitation.Options[:5]
		}},
		{"spray cone", 11, func(s *SpellRecord) { s.SourceDetails.PrismaticSpray.ConeFeet = 30 }},
		{"spray ray type", 11, func(s *SpellRecord) { s.SourceDetails.PrismaticSpray.DamageRays[4].DamageType = "force" }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.PrismaticSpray = &SpellPrismaticSpray{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells215(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells215RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_215.json")
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

func TestSpells215WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_215.json")
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
		`"max_task_days": 10, "max_task_days": 11`,
		`"max_task_days": 10, "MAX_TASK_DAYS": 11`,
		`"max_task_days": 10, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"max_task_days": 10`), []byte(replacement), 1)
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
func TestSpells215PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_215.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"gp_per_unit", 1, 100},
		{"gp_per_unit", 2, 1000},
		{"gp_per_unit", 3, 10000},
		{"max_task_days", 1, 10},
		{"slot", 1, 6},
		{"days", 1, 10},
		{"slot", 2, 7},
		{"days", 2, 30},
		{"slot", 3, 8},
		{"days", 3, 180},
		{"slot", 4, 9},
		{"days", 4, 366},
		{"max_additional_creatures", 1, 8},
		{"overgrowth_radius_feet", 1, 100},
		{"movement_cost_feet", 1, 4},
		{"movement_unit_feet", 1, 1},
		{"enrichment_radius_feet", 1, 2640},
		{"enrichment_days", 1, 365},
		{"enrichment_yield_multiplier", 1, 2},
		{"enrichment_limit_per_year", 1, 1},
		{"level", 6, 5},
		{"dice_count", 12, 2},
		{"die_sides", 12, 12},
		{"level", 7, 11},
		{"dice_count", 13, 3},
		{"die_sides", 13, 12},
		{"level", 8, 17},
		{"dice_count", 14, 4},
		{"die_sides", 14, 12},
		{"hit_point_threshold", 1, 100},
		{"above_dice_count", 1, 12},
		{"above_die_sides", 1, 12},
		{"hit_point_threshold", 2, 150},
		{"above_speed_feet", 1, 0},
		{"max_creatures", 1, 5},
		{"healing_dice_count", 1, 2},
		{"healing_die_sides", 1, 8},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 8},
		{"choose_count", 1, 1},
		{"max_active_non_instantaneous", 1, 3},
		{"cone_feet", 1, 60},
		{"ray_count", 1, 8},
		{"ray_die_sides", 1, 8},
		{"roll", 1, 1},
		{"dice_count", 41, 12},
		{"die_sides", 41, 6},
		{"roll", 2, 2},
		{"dice_count", 42, 12},
		{"die_sides", 42, 6},
		{"roll", 3, 3},
		{"dice_count", 43, 12},
		{"die_sides", 43, 6},
		{"roll", 4, 4},
		{"dice_count", 44, 12},
		{"die_sides", 44, 6},
		{"roll", 5, 5},
		{"dice_count", 45, 12},
		{"die_sides", 45, 6},
		{"indigo_roll", 1, 6},
		{"indigo_successes_to_end", 1, 3},
		{"indigo_failures_to_petrify", 1, 3},
		{"violet_roll", 1, 7},
		{"special_roll", 1, 8},
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
// inventory assigns to catalogue issue #215.
func TestSpells215MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 215 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells215(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells215ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_215.json")
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
	want := "SRD 5.2.1 PDF pp.152-155 adapted into twelve complete offline spell families for issue #215, including page-break continuation text for Planar Binding (pp.152-153), Polymorph (pp.153-154) and Prismatic Spray (pp.154-155). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells215PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.PrismaticSpray = loadSpells215(t).Records[11].SourceDetails.PrismaticSpray
	if r.Validate() == nil {
		t.Fatal("batch 215 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells215RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells215(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell215Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell215Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
