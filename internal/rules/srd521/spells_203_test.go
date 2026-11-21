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

func loadSpells203(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_203.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.123-125 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells203SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/dimension-door","name":"Dimension Door","source_id":"srd-5.2.1:content/spells/dimension-door","source_page":123,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":500},"target":{"shape":"choice","count":1,"eligibility":"destination-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"teleport","mechanic":"teleport.dimension-door","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"companion","mechanic":"teleport.dimension-door-companion","timing":"willing-creature-within-5-feet-when-caster-teleports","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"collision","mechanic":"damage.dimension-door-collision","timing":"arrival-space-occupied-or-completely-filled","dice_count":4,"die_sides":6,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["teleport.dimension-door","teleport.dimension-door-companion","damage.dimension-door-collision"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","dimension_door":{"source_pages":[123,124],"destination_choices":["visible-place","visualized-place","described-by-distance-and-direction"],"direction_examples":["200-feet-straight-downward","300-feet-upward-to-the-northwest-at-45-degree-angle"],"arrival":"exactly-the-spot-desired","companion_count":1,"companion_willing":true,"companion_within_feet":5,"companion_arrival_feet":5,"companion_arrival":"space-within-5-feet-of-destination-space","collision_trigger":"arrival-space-occupied-by-creature-or-completely-filled-by-objects","collision_dice_count":4,"collision_die_sides":6,"collision_damage_type":"force","collision_damaged":"caster-and-each-traveling-creature","collision_result":"teleportation-fails"}}}`,
		`{"id":"srd-5.2.1:content/spells/disguise-self","name":"Disguise Self","source_id":"srd-5.2.1:content/spells/disguise-self","source_page":124,"scope":"complete-source-record","level":1,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"disguise","mechanic":"illusion.disguise-self","timing":"until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"inspection","mechanic":"check.disguise-self-investigation","timing":"study-action-to-inspect-appearance","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.disguise-self","check.disguise-self-investigation"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","disguise_self":{"affects":["caster","clothing","armor","weapons","other-belongings-on-person"],"height_change_feet":1,"height_changes":["shorter","taller"],"weight_changes":["heavier","lighter"],"form_requirement":"same-basic-arrangement-of-limbs","other_extent":"up-to-caster","physical_inspection":"changes-fail-to-hold-up","inspection_example":"added-hat-objects-pass-through-and-touch-feels-nothing","discovery_action":"study-action","discovery_check":"intelligence-investigation","discovery_dc":"caster-spell-save-dc","discovery_result":"creature-discerns-disguise"}}}`,
		`{"id":"srd-5.2.1:content/spells/disintegrate","name":"Disintegrate","source_id":"srd-5.2.1:content/spells/disintegrate","source_page":124,"scope":"complete-source-record","level":6,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"visible-creature-nonmagical-object-or-force-creation-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a lodestone and dust","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":6,"dice_count":3,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.disintegrate","timing":"failed-dexterity-save","dice_count":10,"die_sides":6,"damage_type":"force","bonus":40,"add_casting_modifier":false,"half_initial_only":false},{"id":"dust","mechanic":"disintegrate.zero-hit-points-gray-dust","timing":"damage-reduces-target-to-0-hit-points","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"object","mechanic":"disintegrate.automatic-object-or-force-creation","timing":"target-is-large-or-smaller-nonmagical-object-or-force-creation","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.disintegrate","disintegrate.zero-hit-points-gray-dust","disintegrate.automatic-object-or-force-creation","save.dexterity","scaling.per-slot","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","disintegrate":{"target_choices":["creature","nonmagical-object","creation-of-magical-force"],"force_creation_example":"wall-of-force","save":"dexterity","fail_dice_count":10,"fail_die_sides":6,"fail_bonus":40,"damage_type":"force","success_damage":"none-stated","zero_hit_points_result":"disintegrated-into-gray-dust","dust_includes":"target-and-everything-nonmagical-worn-and-carried","revived_only_by":["true-resurrection","wish"],"automatic_size":"large-or-smaller","automatic_targets":["nonmagical-object","creation-of-magical-force"],"oversize_size":"huge-or-larger","oversize_result":"disintegrates-10-foot-cube-portion","slot_dice_count":3,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/dispel-evil-and-good","name":"Dispel Evil and Good","source_id":"srd-5.2.1:content/spells/dispel-evil-and-good","source_page":124,"scope":"complete-source-record","level":5,"school":"abjuration","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"powdered silver and iron","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"disadvantage","mechanic":"disadvantage.attacks-against-caster-by-listed-types","timing":"until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"break","mechanic":"special.break-enchantment","timing":"magic-action-touch","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"dismissal","mechanic":"special.dismissal","timing":"magic-action-target-within-5-feet","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["disadvantage.attacks-against-caster-by-listed-types","special.break-enchantment","special.dismissal","save.charisma","concentration.source-duration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","dispel_evil":{"affected_types":["celestial","elemental","fey","fiend","undead"],"attack_effect":"disadvantage-on-attack-rolls-against-caster","early_end":"either-special-function","break_action":"magic-action","break_target":"touched-creature-possessed-by-or-charmed-or-frightened-by-listed-types","break_results":["no-longer-possessed","no-longer-charmed","no-longer-frightened"],"dismissal_action":"magic-action","dismissal_range_feet":5,"dismissal_target":"visible-creature-of-listed-type","dismissal_save":"charisma","dismissal_failure":"sent-back-to-home-plane-if-not-there","undead_elsewhere":"shadowfell","fey_elsewhere":"feywild"}}}`,
		`{"id":"srd-5.2.1:content/spells/dispel-magic","name":"Dispel Magic","source_id":"srd-5.2.1:content/spells/dispel-magic","source_page":124,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["bard","cleric","druid","paladin","ranger","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"choice","count":1,"eligibility":"creature-object-or-magical-effect-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-auto-end","above_level":3,"dice_count":0,"die_sides":0,"effect_ids":["check"]},"effects":[{"id":"end","mechanic":"dispel.ongoing-spell-level-3-or-lower","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"check","mechanic":"check.dispel-ongoing-spell-level-4-or-higher","timing":"each-ongoing-spell-level-4-or-higher-on-target","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["dispel.ongoing-spell-level-3-or-lower","check.dispel-ongoing-spell-level-4-or-higher","scaling.auto-end-by-slot-level"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","dispel_magic":{"target_choices":["creature","object","magical-effect"],"auto_end_up_to_level":3,"check_from_level":4,"check_ability":"spellcasting-ability","check_dc_base":10,"check_dc_adds":"spell-level","check_success":"spell-ends","slot_auto_end":"spell-level-equal-to-or-less-than-slot-level"}}}`,
		`{"id":"srd-5.2.1:content/spells/dissonant-whispers","name":"Dissonant Whispers","source_id":"srd-5.2.1:content/spells/dissonant-whispers","source_page":124,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.dissonant-whispers","timing":"failed-wisdom-save-half-on-success","dice_count":3,"die_sides":6,"damage_type":"psychic","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"flee","mechanic":"movement.dissonant-whispers-reaction-flee","timing":"failed-wisdom-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.dissonant-whispers","movement.dissonant-whispers-reaction-flee","save.wisdom-half","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","dissonant_whispers":{"target":"one-creature-of-choice-that-caster-can-see-in-range","perception":"hears-discordant-melody-in-its-mind","save":"wisdom","fail_dice_count":3,"fail_die_sides":6,"damage_type":"psychic","fail_movement_trigger":"must-immediately-use-reaction-if-available","fail_movement":"move-as-far-away-from-caster-as-it-can","fail_route":"safest-route","success_damage":"half-damage-only","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/divination","name":"Divination","source_id":"srd-5.2.1:content/spells/divination","source_page":125,"scope":"complete-source-record","level":4,"school":"divination","lists":["cleric","druid","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"incense worth 25+ GP, which the spell consumes","cost_cp":2500,"consumed":true}},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"contact","mechanic":"divination.contact-god","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"answer","mechanic":"divination.cumulative-no-answer-chance","timing":"repeat-casting-before-long-rest","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["divination.contact-god","divination.cumulative-no-answer-chance","cast.ritual","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","divination":{"contact":"god-or-gods-servants","question_count":1,"question_subject":"specific-goal-event-or-activity","question_window_days":7,"reply_by":"gm","reply_truth":"truthful","reply_forms":["short-phrase","cryptic-rhyme"],"accounts_for_changed_circumstances":false,"changed_circumstances_example":"casting-of-other-spells","repeat_reset":"long-rest","no_answer_chance_percent":25,"no_answer_cumulative":true,"no_answer_after":"first-casting"}}}`,
		`{"id":"srd-5.2.1:content/spells/divine-favor","name":"Divine Favor","source_id":"srd-5.2.1:content/spells/divine-favor","source_page":125,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["paladin"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"damage.divine-favor-weapon-hit","timing":"weapon-attack-hits-until-spell-ends","dice_count":1,"die_sides":4,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.divine-favor-weapon-hit"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","divine_favor":{"attack_kind":"attacks-with-weapons","trigger":"hit","extra_dice_count":1,"extra_die_sides":4,"damage_type":"radiant","applies":"until-spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/divine-smite","name":"Divine Smite","source_id":"srd-5.2.1:content/spells/divine-smite","source_page":125,"scope":"complete-source-record","level":1,"school":"evocation","lists":["paladin"],"casting_time":"bonus-action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"hit-target"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.divine-smite","timing":"immediately-after-hit","dice_count":2,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"bonus","mechanic":"damage.divine-smite-fiend-undead","timing":"target-is-fiend-or-undead","dice_count":1,"die_sides":8,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.divine-smite","damage.divine-smite-fiend-undead","trigger.bonus-action-after-melee-or-unarmed-hit","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":["hit-with-melee-weapon","hit-with-unarmed-strike"],"end_timing":"instantaneous","divine_smite":{"trigger_hits":["melee-weapon","unarmed-strike"],"trigger_timing":"immediately-after-hitting-target","extra_dice_count":2,"extra_die_sides":8,"damage_type":"radiant","type_bonus_dice_count":1,"type_bonus_die_sides":8,"type_bonus_types":["fiend","undead"],"slot_dice_count":1,"slot_die_sides":8}}}`,
		`{"id":"srd-5.2.1:content/spells/divine-word","name":"Divine Word","source_id":"srd-5.2.1:content/spells/divine-word","source_page":125,"scope":"complete-source-record","level":7,"school":"evocation","lists":["cleric"],"casting_time":"bonus-action","range":{"kind":"feet","feet":30},"target":{"shape":"selected-creatures","count":0,"eligibility":"creatures-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"table","mechanic":"effect.divine-word-hit-point-table","timing":"failed-charisma-save-target-50-hit-points-or-fewer","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"banish","mechanic":"banish.divine-word-planar","timing":"failed-charisma-save-regardless-of-hit-points","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["effect.divine-word-hit-point-table","banish.divine-word-planar","save.charisma"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","divine_word":{"source_pages":[125],"affects":"each-creature-of-choice-in-range","save":"charisma","hit_point_limit":50,"limit_basis":"current-hit-points","rows":[{"min_hit_points":0,"max_hit_points":20,"dies":true,"blinded":false,"deafened":false,"stunned":false,"duration_seconds":0},{"min_hit_points":21,"max_hit_points":30,"dies":false,"blinded":true,"deafened":true,"stunned":true,"duration_seconds":3600},{"min_hit_points":31,"max_hit_points":40,"dies":false,"blinded":true,"deafened":true,"stunned":false,"duration_seconds":600},{"min_hit_points":41,"max_hit_points":50,"dies":false,"blinded":false,"deafened":true,"stunned":false,"duration_seconds":60}],"banish_types":["celestial","elemental","fey","fiend"],"banish_regardless_of_hit_points":true,"banish_result":"forced-back-to-plane-of-origin-if-not-there","banish_return_hours":24,"banish_return_exception":"wish-spell"}}}`,
		`{"id":"srd-5.2.1:content/spells/dominate-beast","name":"Dominate Beast","source_id":"srd-5.2.1:content/spells/dominate-beast","source_page":125,"scope":"complete-source-record","level":4,"school":"enchantment","lists":["druid","ranger","sorcerer"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-beast-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-duration-table","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.charmed","timing":"failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"link","mechanic":"link.dominate-telepathic-command","timing":"caster-turn-same-plane-no-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"reaction","mechanic":"command.dominate-reaction","timing":"caster-spends-own-reaction","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.charmed","link.dominate-telepathic-command","command.dominate-reaction","save.wisdom","save.advantage-when-fighting","save.repeat-on-damage","scaling.duration-by-slot-table","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","dominate_beast":{"source_pages":[125],"target_kind":"beast","save":"wisdom","failure_condition":"charmed","condition_duration":"spell-duration","save_advantage":"caster-or-allies-fighting-target","repeat_save":"whenever-target-takes-damage","repeat_success":"spell-ends-on-target","link":"telepathic-while-same-plane","command_timing":"caster-turn","command_action":"no-action-required","command_examples":["Attack that creature","Move over there","Fetch that object"],"obey":"target-does-its-best-to-obey-on-its-turn","idle":"acts-and-moves-as-it-likes-focusing-on-protecting-itself","reaction_command":"caster-can-command-target-to-take-reaction","reaction_cost":"caster-must-take-own-reaction","slot_durations":[{"slot_level":5,"duration_seconds":600,"or_higher":false},{"slot_level":6,"duration_seconds":3600,"or_higher":false},{"slot_level":7,"duration_seconds":28800,"or_higher":true}]}}}`,
		`{"id":"srd-5.2.1:content/spells/dominate-monster","name":"Dominate Monster","source_id":"srd-5.2.1:content/spells/dominate-monster","source_page":125,"scope":"complete-source-record","level":8,"school":"enchantment","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-duration-table","above_level":8,"dice_count":0,"die_sides":0,"effect_ids":["charm"]},"effects":[{"id":"charm","mechanic":"condition.charmed","timing":"failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"link","mechanic":"link.dominate-telepathic-command","timing":"caster-turn-same-plane-no-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"reaction","mechanic":"command.dominate-reaction","timing":"caster-spends-own-reaction","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.charmed","link.dominate-telepathic-command","command.dominate-reaction","save.wisdom","save.advantage-when-fighting","save.repeat-on-damage","scaling.duration-by-slot-table","concentration.source-duration"],"slot_levels":[8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","dominate_monster":{"source_pages":[125,126],"target_kind":"creature","save":"wisdom","failure_condition":"charmed","condition_duration":"spell-duration","save_advantage":"caster-or-allies-fighting-target","repeat_save":"whenever-target-takes-damage","repeat_success":"spell-ends-on-target","link":"telepathic-while-same-plane","command_timing":"caster-turn","command_action":"no-action-required","command_examples":["Attack that creature","Move over there","Fetch that object"],"obey":"target-does-its-best-to-obey-on-its-turn","idle":"acts-and-moves-as-it-likes-focusing-on-protecting-itself","reaction_command":"caster-can-command-target-to-take-reaction","reaction_cost":"caster-must-take-own-reaction","slot_durations":[{"slot_level":9,"duration_seconds":28800,"or_higher":false}]}}}`,
	}
	batch := loadSpells203(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Dimension Door", "Disguise Self", "Disintegrate", "Dispel Evil and Good", "Dispel Magic", "Dissonant Whispers", "Divination", "Divine Favor", "Divine Smite", "Divine Word", "Dominate Beast", "Dominate Monster"}
	pages := []int{123, 124, 124, 124, 124, 124, 125, 125, 125, 125, 125, 125}
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
			// All twelve are levelled spells: slots run from the spell level to 9.
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

func TestSpells203CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells203(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), loadSpells203(t), loadSpells204(t), loadSpells205(t), loadSpells206(t), loadSpells207(t), loadSpells208(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 168 {
		t.Fatalf("expected 168 records, got %d", len(all.Records))
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
func TestSpells203MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_203.json")
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

func TestSpells203JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_203.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"companion_within_feet": 5`, `"companion_within_feet": 5, "companion_within_feet": 6`},
		{`"companion_within_feet": 5`, `"companion_within_feet": 5, "COMPANION_WITHIN_FEET": 6`},
		{`"companion_within_feet": 5`, `"companion_within_feet": 5, "unknown": true`},
		{`"banish_types": [
            "celestial",
            "elemental",
            "fey",
            "fiend"
          ]`, `"banish_types": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/dimension-door"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "damage"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [
        4,
        5,
        6,
        7,
        8,
        9
      ],`, ``},
		{`"slot_levels": [
        4,
        5,
        6,
        7,
        8,
        9
      ],`, `"slot_levels": null,`},
		{`"slot_levels": [
        4,
        5,
        6,
        7,
        8,
        9
      ],`, `"slot_levels": [0],`},
		{`"casting_time": "action"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 0`, `"casting_seconds": 60`},
		{`"feet": 500`, `"feet": 501`},
		{`"cost_cp": 2500`, `"cost_cp": 2501`},
		{`"ritual": true`, `"ritual": false`},
		{`"kind": "slot-auto-end"`, `"kind": "slot"`},
		{`"kind": "slot-duration-table"`, `"kind": "slot-targets"`},
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

func TestSpells203InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"door collision dice", 0, func(s *SpellRecord) { s.SourceDetails.DimensionDoor.CollisionDiceCount = 6 }},
		{"door continuation page", 0, func(s *SpellRecord) { s.SourceDetails.DimensionDoor.SourcePages = []int{123} }},
		{"door companion distance", 0, func(s *SpellRecord) { s.SourceDetails.DimensionDoor.CompanionWithinFeet = 10 }},
		{"door range", 0, func(s *SpellRecord) { s.Range.Feet = 400 }},
		{"door somatic", 0, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"disguise height", 1, func(s *SpellRecord) { s.SourceDetails.DisguiseSelf.HeightChangeFeet = 2 }},
		{"disguise check", 1, func(s *SpellRecord) { s.SourceDetails.DisguiseSelf.DiscoveryCheck = "intelligence-arcana" }},
		{"disguise duration", 1, func(s *SpellRecord) { s.DurationSeconds = 600 }},
		{"disintegrate bonus", 2, func(s *SpellRecord) { s.SourceDetails.Disintegrate.FailBonus = 30 }},
		{"disintegrate effect bonus", 2, func(s *SpellRecord) { s.Effects[0].Bonus = 0 }},
		{"disintegrate scaling row", 2, func(s *SpellRecord) { s.Scaling.DiceCount = 2 }},
		{"disintegrate material", 2, func(s *SpellRecord) { s.Components.Material.Description = "a lodestone" }},
		{"disintegrate size", 2, func(s *SpellRecord) { s.SourceDetails.Disintegrate.AutomaticSize = "huge-or-smaller" }},
		{"dispel evil type", 3, func(s *SpellRecord) { s.SourceDetails.DispelEvil.AffectedTypes[4] = "humanoid" }},
		{"dispel evil range", 3, func(s *SpellRecord) { s.SourceDetails.DispelEvil.DismissalRange = 10 }},
		{"dispel evil plane", 3, func(s *SpellRecord) { s.SourceDetails.DispelEvil.UndeadElsewhere = "feywild" }},
		{"dispel evil concentration", 3, func(s *SpellRecord) { *s.Concentration = false }},
		{"dispel magic level", 4, func(s *SpellRecord) { s.SourceDetails.DispelMagic.AutoEndUpToLvl = 4 }},
		{"dispel magic dc", 4, func(s *SpellRecord) { s.SourceDetails.DispelMagic.CheckDCBase = 8 }},
		{"dispel magic scaling kind", 4, func(s *SpellRecord) { s.Scaling.Kind = "slot" }},
		{"dispel magic scaling effect", 4, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"unknown"} }},
		{"dispel magic scaling mechanic", 4, func(s *SpellRecord) { s.RequiredMechanics = s.RequiredMechanics[:2] }},
		{"whispers dice", 5, func(s *SpellRecord) { s.SourceDetails.DissonantWhispers.FailDiceCount = 2 }},
		{"whispers scaling level", 5, func(s *SpellRecord) { s.Scaling.AboveLevel = 2 }},
		{"whispers damage type", 5, func(s *SpellRecord) { s.SourceDetails.DissonantWhispers.DamageType = "force" }},
		{"divination consumed", 6, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"divination cost", 6, func(s *SpellRecord) { *s.Components.Material.CostCP = 2000 }},
		{"divination window", 6, func(s *SpellRecord) { s.SourceDetails.Divination.QuestionWindowDays = 10 }},
		{"divination chance", 6, func(s *SpellRecord) { s.SourceDetails.Divination.NoAnswerChancePercent = 20 }},
		{"divination ritual", 6, func(s *SpellRecord) { *s.Ritual = false }},
		{"favor casting time", 7, func(s *SpellRecord) { s.CastingTime = "action" }},
		{"favor die", 7, func(s *SpellRecord) { s.SourceDetails.DivineFavor.ExtraDieSides = 6 }},
		{"smite triggers", 8, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{} }},
		{"smite types", 8, func(s *SpellRecord) { s.SourceDetails.DivineSmite.TypeBonusTypes[1] = "celestial" }},
		{"smite scaling die", 8, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"smite component", 8, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"word row duration", 9, func(s *SpellRecord) { s.SourceDetails.DivineWord.Rows[1].DurationSeconds = 600 }},
		{"word row death", 9, func(s *SpellRecord) { s.SourceDetails.DivineWord.Rows[0].Dies = false }},
		{"word row bound", 9, func(s *SpellRecord) { s.SourceDetails.DivineWord.Rows[3].MaxHitPoints = 49 }},
		{"word hit point limit", 9, func(s *SpellRecord) { s.SourceDetails.DivineWord.HitPointLimit = 60 }},
		{"word banish hours", 9, func(s *SpellRecord) { s.SourceDetails.DivineWord.BanishReturnHrs = 12 }},
		{"word range", 9, func(s *SpellRecord) { s.Range.Feet = 60 }},
		{"beast slot table", 10, func(s *SpellRecord) { s.SourceDetails.DominateBeast.SlotDurations[2].OrHigher = false }},
		{"beast slot duration", 10, func(s *SpellRecord) { s.SourceDetails.DominateBeast.SlotDurations[0].DurationSeconds = 300 }},
		{"beast target kind", 10, func(s *SpellRecord) { s.SourceDetails.DominateBeast.TargetKind = "creature" }},
		{"beast scaling kind", 10, func(s *SpellRecord) { s.Scaling.Kind = "slot" }},
		{"beast target", 10, func(s *SpellRecord) { s.Target.Eligibility = "visible-creature-in-range" }},
		{"monster continuation page", 11, func(s *SpellRecord) { s.SourceDetails.DominateMonster.SourcePages = []int{125} }},
		{"monster slot level", 11, func(s *SpellRecord) { s.SourceDetails.DominateMonster.SlotDurations[0].SlotLevel = 8 }},
		{"monster duration", 11, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"monster swapped payload", 11, func(s *SpellRecord) {
			s.SourceDetails.DominateBeast, s.SourceDetails.DominateMonster = s.SourceDetails.DominateMonster, nil
		}},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.DispelMagic = &SpellDispelMagic{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells203(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells203RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_203.json")
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

func TestSpells203WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_203.json")
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
		`"slot_level": 5, "slot_level": 6`,
		`"slot_level": 5, "SLOT_LEVEL": 6`,
		`"slot_level": 5, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"slot_level": 5`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells203PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_203.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"companion_within_feet", 1, 5}, {"companion_arrival_feet", 1, 5}, {"collision_dice_count", 1, 4}, {"collision_die_sides", 1, 6},
		{"height_change_feet", 1, 1}, {"fail_dice_count", 1, 10}, {"fail_dice_count", 2, 3}, {"fail_bonus", 1, 40}, {"slot_dice_count", 1, 3}, {"slot_dice_count", 2, 1},
		{"dismissal_range_feet", 1, 5}, {"auto_end_up_to_level", 1, 3}, {"check_from_level", 1, 4}, {"check_dc_base", 1, 10},
		{"question_count", 1, 1}, {"question_window_days", 1, 7}, {"no_answer_chance_percent", 1, 25},
		{"extra_dice_count", 1, 1}, {"extra_die_sides", 1, 4}, {"extra_dice_count", 2, 2}, {"extra_die_sides", 2, 8}, {"type_bonus_dice_count", 1, 1}, {"type_bonus_die_sides", 1, 8},
		{"hit_point_limit", 1, 50}, {"min_hit_points", 2, 21}, {"max_hit_points", 1, 20}, {"max_hit_points", 4, 50}, {"banish_return_hours", 1, 24},
		{"slot_level", 1, 5}, {"slot_level", 2, 6}, {"slot_level", 3, 7}, {"slot_level", 4, 9},
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
// inventory assigns to catalogue issue #203.
func TestSpells203MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 203 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells203(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells203ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_203.json")
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
	want := "SRD 5.2.1 PDF pp.123-125 adapted into twelve complete offline spell families for issue #203, including page-break continuation text for Dimension Door (pp.123-124) and Dominate Monster (pp.125-126). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	if strings.Contains(root.Source.Changes, "Divine Word") {
		t.Fatal("Divine Word sits wholly on p.125 and is not a page-break continuation")
	}
}
