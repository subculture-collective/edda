package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// The exact issue #197 families, in source order.
var spells197Names = []string{"Antipathy/Sympathy", "Arcane Eye", "Arcane Hand", "Arcane Lock", "Arcane Sword", "Arcanist’s Magic Aura", "Astral Projection", "Augury", "Aura of Life", "Awaken", "Bane", "Banishment"}

func loadSpells197(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_197.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.109-112 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells197SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/antipathy-sympathy","name":"Antipathy/Sympathy","source_id":"srd-5.2.1:content/spells/antipathy-sympathy","source_page":109,"scope":"complete-source-record","level":8,"school":"enchantment","lists":["bard","druid","wizard"],"casting_time":"1-hour","range":{"kind":"feet","feet":60},"target":{"shape":"choice","count":1,"eligibility":"huge-or-smaller-creature-or-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a mix of vinegar and honey","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":864000,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"antipathy","mechanic":"condition.antipathy-frightened-flee","timing":"failed-save-while-antipathy-chosen","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"sympathy","mechanic":"condition.sympathy-charmed-approach","timing":"failed-save-while-sympathy-chosen","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.antipathy-frightened-flee","condition.sympathy-charmed-approach","choice.antipathy-or-sympathy","target.creature-or-object-huge-or-smaller","target.specified-creature-kind","save.wisdom-within-120-feet","save.wisdom-end-turn-beyond-120-feet","save.wisdom-on-target-damage-sympathy","immunity.one-minute-after-successful-save","cast.one-hour","components.material"],"slot_levels":[8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"duration-ended","antipathy":{"source_pages":[109,110],"modes":[{"choice":"antipathy","condition":"frightened","movement":"must-use-movement-on-its-turns-to-get-as-far-away-as-possible-by-safest-route","within_feet":0,"within_restriction":"not-applicable","damage_save":"not-applicable"},{"choice":"sympathy","condition":"charmed","movement":"must-use-movement-on-its-turns-to-get-as-close-as-possible-by-safest-route","within_feet":5,"within_restriction":"cannot-willingly-move-away","damage_save":"wisdom-save-to-end-effect-when-target-damages-charmed-creature"}],"choose_count":1,"target_choices":["creature","object"],"maximum_size":"huge","creature_kind":"caster-specified-kind-of-creature","kind_examples":["red dragons","goblins","vampires"],"save":"wisdom","save_trigger":"comes-within-distance-of-target","save_distance_feet":120,"end_save_timing":"ends-turn-more-than-distance-from-target","end_save_distance_feet":120,"end_save_success":"no-longer-affected-by-target","immunity_seconds":60,"after_immunity":"can-be-affected-again"}}}`,
		`{"id":"srd-5.2.1:content/spells/arcane-eye","name":"Arcane Eye","source_id":"srd-5.2.1:content/spells/arcane-eye","source_page":110,"scope":"complete-source-record","level":4,"school":"divination","lists":["wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of bat fur","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"eye","mechanic":"sensor.arcane-eye","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["sensor.arcane-eye","senses.every-direction-and-darkvision-30","movement.bonus-action-30-feet","movement.one-inch-opening","condition.invisible","components.material","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","eye":{"invisible":true,"invulnerable":true,"hovers":true,"vision":"caster-mentally-receives-visual-information-in-every-direction","darkvision_feet":30,"move_action":"bonus-action","move_feet":30,"move_direction":"any","blocked":"solid-barrier","minimum_opening_inches":1}}}`,
		`{"id":"srd-5.2.1:content/spells/arcane-hand","name":"Arcane Hand","source_id":"srd-5.2.1:content/spells/arcane-hand","source_page":110,"scope":"complete-source-record","level":5,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"an eggshell and a glove","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-effect-dice","above_level":5,"dice_count":0,"die_sides":0,"effect_ids":["clenched-fist","grasping-crush"]},"effects":[{"id":"clenched-fist","mechanic":"attack.arcane-hand-clenched-fist","timing":"on-cast-or-bonus-action-melee-spell-attack-hit","dice_count":5,"die_sides":8,"damage_type":"force","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"forceful-hand","mechanic":"movement.arcane-hand-push","timing":"on-cast-or-bonus-action-failed-strength-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"grasping-hand","mechanic":"condition.arcane-hand-grapple","timing":"on-cast-or-bonus-action-failed-dexterity-save","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"grasping-crush","mechanic":"damage.arcane-hand-crush","timing":"bonus-action-while-hand-grapples-target","dice_count":4,"die_sides":6,"damage_type":"bludgeoning","bonus":0,"add_casting_modifier":true,"half_initial_only":false},{"id":"interposing-hand","mechanic":"cover.arcane-hand-interposing","timing":"on-cast-or-bonus-action","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["attack.arcane-hand-clenched-fist","movement.arcane-hand-push","condition.arcane-hand-grapple","damage.arcane-hand-crush","cover.arcane-hand-interposing","object.arcane-hand-ac-20-caster-hit-point-maximum","movement.bonus-action-60-feet","choice.four-hand-effects","attack.melee-spell","save.strength","save.dexterity","condition.grappled-spell-save-dc","cover.half","terrain.difficult-for-enemies","scaling.per-slot-by-effect","components.material","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-hand-hit-points-ended","hand":{"size":"large","placement":"visible-unoccupied-space-in-range","appearance":"hand-of-shimmering-magical-energy","object_ac":20,"hit_points":"caster-hit-point-maximum","zero_hit_points":"spell-ends","occupies_space":false,"command_timings":["when-cast","bonus-action-on-later-turns"],"move_feet":60,"choose_count":1,"options":[{"choice":"clenched-fist","reach_feet":5,"maximum_size":"not-applicable","resolution":"melee-spell-attack","dice_count":5,"die_sides":8,"damage_type":"force","add_casting_modifier":false,"push_base_feet":0,"push_feet_per_modifier":0,"outcome":"damage-on-hit","followup":"none"},{"choice":"forceful-hand","reach_feet":5,"maximum_size":"huge","resolution":"strength-save","dice_count":0,"die_sides":0,"damage_type":"none","add_casting_modifier":false,"push_base_feet":5,"push_feet_per_modifier":5,"outcome":"pushed-up-to-base-plus-per-spellcasting-modifier-feet-on-failed-save","followup":"hand-moves-with-target-remaining-within-5-feet"},{"choice":"grasping-hand","reach_feet":5,"maximum_size":"huge","resolution":"dexterity-save","dice_count":4,"die_sides":6,"damage_type":"bludgeoning","add_casting_modifier":true,"push_base_feet":0,"push_feet_per_modifier":0,"outcome":"grappled-with-escape-dc-equal-to-spell-save-dc-on-failed-save","followup":"bonus-action-crush-damage-while-hand-grapples-target"},{"choice":"interposing-hand","reach_feet":0,"maximum_size":"not-applicable","resolution":"none","dice_count":0,"die_sides":0,"damage_type":"none","add_casting_modifier":false,"push_base_feet":0,"push_feet_per_modifier":0,"outcome":"half-cover-for-caster-against-attacks-and-effects-originating-from-or-passing-through-its-space","followup":"its-space-is-difficult-terrain-for-caster-enemies"}],"slot_dice":[{"effect_id":"clenched-fist","dice_count":2,"die_sides":8},{"effect_id":"grasping-crush","dice_count":2,"die_sides":6}]}}}`,
		`{"id":"srd-5.2.1:content/spells/arcane-lock","name":"Arcane Lock","source_id":"srd-5.2.1:content/spells/arcane-lock","source_page":110,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"closed-door-window-gate-container-or-hatch-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"gold dust worth 25+ GP, which the spell consumes","cost_cp":2500,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"lock","mechanic":"object.arcane-lock","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["object.arcane-lock","target.closed-object-touch","object.no-nonmagical-unlock","object.caster-and-designated-creatures-bypass","object.password-unlock-one-minute","components.priced-consumed-material","components.material","duration.until-dispelled"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"dispelled","lock":{"objects":["door","window","gate","container","hatch"],"must_be_closed":true,"unlock":"no-nonmagical-means","bypass":["caster","creatures-designated-when-cast"],"password_optional":true,"password_distance_feet":5,"password_unlock_seconds":60}}}`,
		`{"id":"srd-5.2.1:content/spells/arcane-sword","name":"Arcane Sword","source_id":"srd-5.2.1:content/spells/arcane-sword","source_page":110,"scope":"complete-source-record","level":7,"school":"evocation","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a miniature sword worth 250+ GP","cost_cp":25000,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"damage","mechanic":"attack.arcane-sword","timing":"melee-spell-attack-hit-on-appearance-or-bonus-action","dice_count":4,"die_sides":12,"damage_type":"force","bonus":0,"add_casting_modifier":true,"half_initial_only":false}],"required_mechanics":["attack.arcane-sword","attack.melee-spell","damage.add-spellcasting-modifier","movement.bonus-action-30-feet","components.priced-material","components.material","concentration.source-duration"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","sword":{"appearance":"spectral-sword-hovering-within-range","attack":"melee-spell-attack","attack_timings":["when-sword-appears","bonus-action-on-later-turns"],"reach_feet":5,"move_action":"bonus-action","move_feet":30,"move_destination":"visible-spot","repeat_targets":["same-target","different-target"]}}}`,
		`{"id":"srd-5.2.1:content/spells/arcanists-magic-aura","name":"Arcanist’s Magic Aura","source_id":"srd-5.2.1:content/spells/arcanists-magic-aura","source_page":110,"scope":"complete-source-record","level":2,"school":"illusion","lists":["wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"willing-creature-or-unattended-object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a small square of silk","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"mask","mechanic":"illusion.magic-aura-mask","timing":"for-duration-on-willing-creature","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"false-aura","mechanic":"illusion.magic-aura-false-aura","timing":"for-duration-on-object-not-worn-or-carried","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["illusion.magic-aura-mask","illusion.magic-aura-false-aura","target.willing-creature-or-unattended-object","choice.creature-type","choice.false-aura-appearance","duration.until-dispelled-after-30-daily-casts","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","magic_aura":{"source_pages":[110,111],"variants":[{"choice":"mask","target":"willing-creature","effect":"spells-and-other-magical-effects-treat-target-as-chosen-creature-type","options":["creature-type-other-than-actual-type"]},{"choice":"false-aura","target":"object-not-worn-or-carried","effect":"changes-appearance-to-spells-and-magical-effects-that-detect-magical-auras","options":["nonmagical-object-appears-magical","magic-item-appears-nonmagical","aura-appears-to-belong-to-chosen-school-of-magic"]}],"detection_example":"Detect Magic","repeat_days":30,"repeat_requirement":"cast-on-same-target-every-day","repeat_result":"illusion-lasts-until-dispelled"}}}`,
		`{"id":"srd-5.2.1:content/spells/astral-projection","name":"Astral Projection","source_id":"srd-5.2.1:content/spells/astral-projection","source_page":111,"scope":"complete-source-record","level":9,"school":"necromancy","lists":["cleric","warlock","wizard"],"casting_time":"1-hour","range":{"kind":"feet","feet":10},"target":{"shape":"self-and-up-to-count","count":8,"eligibility":"willing-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"for each of the spell’s targets, one jacinth worth 1,000+ GP and one silver bar worth 100+ GP, all of which the spell consumes","cost_cp":110000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"projection","mechanic":"plane.astral-projection","timing":"until-dispelled-or-ended-for-target","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["plane.astral-projection","target.self-and-up-to-eight-willing","condition.unconscious-suspended-animation","plane.silver-cord","death.cord-cut","plane.leave-astral-reenter-body","effect.separate-body-and-astral-form","end.zero-hit-points-per-target","end.magic-action-dismiss","cast.one-hour","components.priced-consumed-material-per-target","components.material","duration.until-dispelled"],"slot_levels":[9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"dispelled-dismissed-or-target-at-zero-hit-points","astral":{"maximum_other_targets":8,"includes_caster":true,"willing":true,"jacinth_minimum_cost_cp":100000,"silver_bar_minimum_cost_cp":10000,"materials_per_target":true,"destination_plane":"astral","already_on_plane":"spell-ends-instantly","body_state":"suspended-animation","body_condition":"unconscious","unneeded_sustenance":["food","air"],"body_ages":false,"astral_form":"resembles-body-replicating-game-statistics-and-possessions","cord":"silvery-cord-trailing-from-between-shoulder-blades","cord_visible_feet":1,"cord_cut_condition":"only-when-an-effect-states-that-it-does-so","cord_cut_result":"body-and-astral-form-both-die","leaving_plane":"body-and-possessions-travel-along-cord-and-target-reenters-body-on-new-plane","separate_effects":true,"zero_hit_points":["body","astral-form"],"zero_hit_points_result":"spell-ends-for-that-target","dismiss_action":"magic-action","dismiss_result":"spell-ends-for-all-targets","end_result":"target-not-dead-reappears-in-body-and-exits-suspended-animation"}}}`,
		`{"id":"srd-5.2.1:content/spells/augury","name":"Augury","source_id":"srd-5.2.1:content/spells/augury","source_page":111,"scope":"complete-source-record","level":2,"school":"divination","lists":["cleric","druid","wizard"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"specially marked sticks, bones, cards, or other divinatory tokens worth 25+ GP","cost_cp":2500,"consumed":false}},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"omen","mechanic":"divination.augury-omen","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["divination.augury-omen","divination.gm-chosen-omen-table","chance.cumulative-25-percent-no-answer","cast.one-minute","cast.ritual","components.priced-material","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"instantaneous","augury":{"subject":"results-of-course-of-action-planned","window_minutes":30,"omen_source":"otherworldly-entity","chosen_by":"gm","omens":[{"omen":"Weal","results":"Good"},{"omen":"Woe","results":"Bad"},{"omen":"Weal and woe","results":"Good and bad"},{"omen":"Indifference","results":"Neither good nor bad"}],"limitation":"does-not-account-for-circumstances-that-might-change-results","repeat_scope":"each-casting-after-first-before-finishing-long-rest","cumulative_no_answer_percent":25}}}`,
		`{"id":"srd-5.2.1:content/spells/aura-of-life","name":"Aura of Life","source_id":"srd-5.2.1:content/spells/aura-of-life","source_page":111,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["cleric","paladin"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"self-and-allies-in-aura"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"aura","mechanic":"defense.aura-of-life","timing":"while-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"revive","mechanic":"healing.aura-of-life-revive","timing":"ally-with-0-hit-points-starts-turn-in-aura","dice_count":0,"die_sides":0,"damage_type":"","bonus":1,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["defense.aura-of-life","healing.aura-of-life-revive","area.emanation-30","resistance.necrotic","hit-points.maximum-cannot-be-reduced","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","life":{"emanation_feet":30,"affects":["caster","allies"],"resistance":"necrotic","hit_point_maximum_reduction":"prevented","revive_trigger":"ally-with-0-hit-points-starts-turn-in-aura","revive_hit_points":1}}}`,
		`{"id":"srd-5.2.1:content/spells/awaken","name":"Awaken","source_id":"srd-5.2.1:content/spells/awaken","source_page":111,"scope":"complete-source-record","level":5,"school":"transmutation","lists":["bard","druid"],"casting_time":"8-hours","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"low-intelligence-beast-or-plant-creature-or-natural-plant-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"an agate worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"awaken","mechanic":"transformation.awaken","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"charm","mechanic":"condition.awaken-charmed","timing":"30-days-or-until-caster-or-allies-deal-damage","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transformation.awaken","condition.awaken-charmed","target.beast-or-plant-intelligence-3-or-natural-plant","ability.intelligence-10","language.one-known-by-caster","creature.natural-plant-becomes-plant-creature","statistics.gm-chosen","attitude.chosen-after-charm","cast.eight-hours","components.priced-consumed-material","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":28800,"casting_triggers":[],"end_timing":"instantaneous","awaken":{"source_pages":[111,112],"target_choices":["beast-creature","plant-creature","natural-noncreature-plant"],"maximum_intelligence":3,"gained_intelligence":10,"language":"speak-one-language-caster-knows","plant_becomes":"plant-creature","plant_gains":["move-limbs-roots-vines-creepers","senses-similar-to-human"],"plant_statistics":"gm-chooses-appropriate-statistics","statistics_examples":["Awakened Shrub","Awakened Tree"],"condition":"charmed","condition_days":30,"condition_ends":["after-condition-days","caster-or-allies-deal-damage-to-it"],"after_condition":"awakened-creature-chooses-attitude-toward-caster"}}}`,
		`{"id":"srd-5.2.1:content/spells/bane","name":"Bane","source_id":"srd-5.2.1:content/spells/bane","source_page":112,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","cleric","warlock"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":3,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a drop of blood","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["penalty"]},"effects":[{"id":"penalty","mechanic":"roll.bane","timing":"each-failed-target-attack-roll-or-saving-throw","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["roll.bane","save.charisma","roll.attack-and-save-subtract-d4","scaling.targets-per-slot","components.material","concentration.source-duration"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","bane":{"base_targets":3,"additional_targets_per_slot":1,"save":"charisma","save_success":"unaffected","rolls":["attack-rolls","saving-throws"],"die_sides":4,"roll_modifier":"subtract","target_choice":"up-to-count-visible-creatures-in-range"}}}`,
		`{"id":"srd-5.2.1:content/spells/banishment","name":"Banishment","source_id":"srd-5.2.1:content/spells/banishment","source_page":112,"scope":"complete-source-record","level":4,"school":"abjuration","lists":["cleric","paladin","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"up-to-count","count":1,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pentacle","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["banish"]},"effects":[{"id":"banish","mechanic":"transport.banishment","timing":"failed-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"no-return","mechanic":"transport.banishment-associated-plane","timing":"listed-creature-type-and-spell-lasts-1-minute","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["transport.banishment","transport.banishment-associated-plane","save.charisma","condition.incapacitated","transport.return-space-or-nearest-unoccupied","target.creature-type","scaling.targets-per-slot","components.material","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","banishment":{"base_targets":1,"additional_targets_per_slot":1,"save":"charisma","save_success":"unaffected","destination":"harmless-demiplane","condition":"incapacitated","return":"space-it-left-or-nearest-unoccupied-space-if-occupied","no_return_types":["aberration","celestial","elemental","fey","fiend"],"no_return_condition":"spell-lasts-for-full-duration","no_return_seconds":60,"no_return_destination":"random-location-on-plane-associated-with-creature-type-gm-choice"}}}`,
	}
	batch := loadSpells197(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := spells197Names
	for i, r := range batch.Records {
		t.Run(names[i], func(t *testing.T) {
			var want SpellRecord
			if err := json.Unmarshal([]byte(expected[i]), &want); err != nil {
				t.Fatal(err)
			}
			if r.Name != names[i] || !reflect.DeepEqual(r, want) {
				t.Fatalf("PDF p.%d source values differ", want.SourcePage)
			}
			if r.MechanicsAvailable() {
				t.Fatal("source effect available")
			}
			slots := []int{}
			for level := r.Level; level <= 9; level++ {
				slots = append(slots, level)
			}
			if !slices.Equal(r.SlotLevels, slots) {
				t.Fatal("source slot metadata incomplete")
			}
		})
	}
}

func TestSpells197CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells197(t)
	all, err := ComposeSpellSourceCatalogue(legacy, loadSpells218(t), loadSpells198(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 37 {
		t.Fatal("expected Acid Arrow and three twelve-family batches")
	}
	for _, r := range all.Records {
		if r.Scope != "complete-source-record" || r.MechanicsAvailable() {
			t.Fatal("source granted availability")
		}
	}
	if _, err := ComposeSpellSourceCatalogue(batch, batch); err == nil {
		t.Fatal("duplicate canonical IDs accepted")
	}
	for _, r := range legacy.Records[1:] {
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
// The valid input covers explicit false and zero values.
func TestSpells197MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_197.json")
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

func TestSpells197JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_197.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"darkvision_feet": 30`, `"darkvision_feet": 30, "darkvision_feet": 60`},
		{`"darkvision_feet": 30`, `"darkvision_feet": 30, "DARKVISION_FEET": 60`},
		{`"darkvision_feet": 30`, `"darkvision_feet": 30, "unknown": true`},
		{`"affects": [
            "caster",
            "allies"
          ]`, `"affects": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/bane"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "penalty"
        ]`, `"effect_ids": ["unknown"]`},
		{`"casting_time": "8-hours"`, `"casting_time": "action"`},
		{`"casting_seconds": 3600`, `"casting_seconds": 600`},
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

func TestSpells197InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"antipathy unknown mode", 0, func(s *SpellRecord) { s.SourceDetails.Antipathy.Modes[1].Choice = "unknown" }},
		{"antipathy continuation page", 0, func(s *SpellRecord) { s.SourceDetails.Antipathy.SourcePages = []int{109} }},
		{"antipathy casting seconds", 0, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"eye casting seconds", 1, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 60 }},
		{"hand unknown option", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options[3].Choice = "unknown" }},
		{"hand omitted option", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options = s.SourceDetails.Hand.Options[:3] }},
		{"hand slot dice effect", 2, func(s *SpellRecord) { s.SourceDetails.Hand.SlotDice[1].EffectID = "forceful-hand" }},
		{"hand slot dice count", 2, func(s *SpellRecord) { s.SourceDetails.Hand.SlotDice[0].DiceCount = 0 }},
		{"hand slot dice sides", 2, func(s *SpellRecord) { s.SourceDetails.Hand.SlotDice[0].DieSides = 7 }},
		{"hand slot dice omitted", 2, func(s *SpellRecord) { s.SourceDetails.Hand.SlotDice = s.SourceDetails.Hand.SlotDice[:1] }},
		{"hand scaling dropped", 2, func(s *SpellRecord) { s.Scaling = SpellScaling{Kind: "none", EffectIDs: []string{}} }},
		{"hand scaling dependency", 2, func(s *SpellRecord) {
			s.RequiredMechanics = slices.DeleteFunc(slices.Clone(s.RequiredMechanics), func(m string) bool { return m == "scaling.per-slot-by-effect" })
		}},
		{"per-effect scaling without hand", 4, func(s *SpellRecord) {
			s.Scaling = SpellScaling{Kind: "slot-effect-dice", AboveLevel: 7, EffectIDs: []string{"damage"}}
			s.RequiredMechanics = append(s.RequiredMechanics, "scaling.per-slot-by-effect")
		}},
		{"magic aura variant", 5, func(s *SpellRecord) { s.SourceDetails.MagicAura.Variants[0].Choice = "unknown" }},
		{"magic aura object options", 5, func(s *SpellRecord) {
			s.SourceDetails.MagicAura.Variants[1].Options = s.SourceDetails.MagicAura.Variants[1].Options[:2]
		}},
		{"astral target count", 6, func(s *SpellRecord) { s.SourceDetails.Astral.MaximumOtherTargets = 7 }},
		{"astral component minimum", 6, func(s *SpellRecord) { s.SourceDetails.Astral.SilverBarMinimumCostCP = 0 }},
		{"astral unconsumed", 6, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"augury omen", 7, func(s *SpellRecord) { s.SourceDetails.Augury.Omens[2].Omen = "unknown" }},
		{"augury chance", 7, func(s *SpellRecord) { s.SourceDetails.Augury.CumulativeNoAnswerPercent = 101 }},
		{"augury ritual", 7, func(s *SpellRecord) { *s.Ritual = false }},
		{"aura of life healing", 8, func(s *SpellRecord) { s.SourceDetails.Life.ReviveHitPoints = 0 }},
		{"awaken target", 9, func(s *SpellRecord) { s.SourceDetails.Awaken.TargetChoices[2] = "unknown" }},
		{"awaken continuation page", 9, func(s *SpellRecord) { s.SourceDetails.Awaken.SourcePages = []int{111} }},
		{"bane save", 10, func(s *SpellRecord) { s.SourceDetails.Bane.Save = "wisdom" }},
		{"bane scaling dropped", 10, func(s *SpellRecord) { s.Scaling = SpellScaling{Kind: "none", EffectIDs: []string{}} }},
		{"banishment type", 11, func(s *SpellRecord) { s.SourceDetails.Banishment.NoReturnTypes[0] = "unknown" }},
		{"banishment full duration", 11, func(s *SpellRecord) { s.SourceDetails.Banishment.NoReturnSeconds = 600 }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Eye = &SpellEye{} }},
		{"earlier batch payload", 1, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"payload on other spell", 1, func(s *SpellRecord) {
			s.SourceDetails = &SpellSourceDetails{CastingTriggers: []string{}, EndTiming: "x", Lock: &SpellLock{}}
		}},
		{"missing details", 3, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"casting trigger", 1, func(s *SpellRecord) { s.SourceDetails.CastingTriggers = []string{"unknown"} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
		{"bounded scope", 10, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells197(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

// Values a record stores twice, or that its required mechanics depend on,
// must agree; PDF pp.109-111 literals are in TestSpells197SourceValues.
func TestSpells197RejectInconsistentDuplicates(t *testing.T) {
	dropMechanic := func(s *SpellRecord, key string) {
		s.RequiredMechanics = slices.DeleteFunc(slices.Clone(s.RequiredMechanics), func(m string) bool { return m == key })
	}
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"hand clenched fist option dice", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options[0].DiceCount = 6 }},
		{"hand clenched fist effect dice", 2, func(s *SpellRecord) { s.Effects[0].DiceCount = 6 }},
		{"hand clenched fist option damage type", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options[0].DamageType = "bludgeoning" }},
		{"hand grasping option sides", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options[2].DieSides = 8 }},
		{"hand grasping option modifier", 2, func(s *SpellRecord) { s.SourceDetails.Hand.Options[2].AddCastingModifier = false }},
		{"bane effect die", 10, func(s *SpellRecord) { s.Effects[0].DieSides = 6 }},
		{"bane effect dice count", 10, func(s *SpellRecord) { s.Effects[0].DiceCount = 2 }},
		{"lock unconsumed", 3, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"lock free material", 3, func(s *SpellRecord) { *s.Components.Material.CostCP = 0 }},
		{"lock priced mechanic dropped", 3, func(s *SpellRecord) { dropMechanic(s, "components.priced-consumed-material") }},
		{"sword free material", 4, func(s *SpellRecord) { *s.Components.Material.CostCP = 0 }},
		{"sword consumed", 4, func(s *SpellRecord) { *s.Components.Material.Consumed = true }},
		{"awaken free material", 9, func(s *SpellRecord) { *s.Components.Material.CostCP = 0 }},
		{"eye priced material", 1, func(s *SpellRecord) { *s.Components.Material.CostCP = 100 }},
		{"lock timed duration", 3, func(s *SpellRecord) { s.Duration, s.DurationSeconds = "timed", 60 }},
		{"lock duration mechanic dropped", 3, func(s *SpellRecord) { dropMechanic(s, "duration.until-dispelled") }},
		{"eye until dispelled", 1, func(s *SpellRecord) {
			s.Duration, s.DurationSeconds, *s.Concentration = "until-dispelled", 0, false
		}},
		{"eye darkvision", 1, func(s *SpellRecord) { s.SourceDetails.Eye.DarkvisionFeet = 0 }},
		{"eye move", 1, func(s *SpellRecord) { s.SourceDetails.Eye.MoveFeet = 0 }},
		{"eye opening", 1, func(s *SpellRecord) { s.SourceDetails.Eye.MinimumOpeningInches = 0 }},
		{"sword move", 4, func(s *SpellRecord) { s.SourceDetails.Sword.MoveFeet = 0 }},
		{"sword reach", 4, func(s *SpellRecord) { s.SourceDetails.Sword.ReachFeet = 0 }},
		{"eye slot levels", 1, func(s *SpellRecord) { s.SlotLevels = []int{4, 5} }},
		{"astral level", 6, func(s *SpellRecord) { s.Level, s.SlotLevels = 8, []int{8, 9} }},
		{"bane target count", 10, func(s *SpellRecord) { s.Target.Count = 1 }},
		{"lock target from another record", 3, func(s *SpellRecord) {
			s.Target = SpellTarget{Shape: "created-effect", Count: 1, Eligibility: "point-in-range"}
		}},
		{"life target", 8, func(s *SpellRecord) { s.Target = SpellTarget{Shape: "single", Count: 1, Eligibility: "self"} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells197(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil {
				t.Fatal("inconsistent record validated")
			}
			b, err := json.Marshal(r)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
				t.Fatal("inconsistent record decoded")
			}
		})
	}
}

// Earlier batches cannot adopt a long casting time without its exact length.
func TestSpells197CastingTimeOnEarlierBatch(t *testing.T) {
	r := loadSpells198(t)
	r.Records[0].CastingTime = "1-hour"
	if r.Validate() == nil {
		t.Fatal("long casting time accepted without casting seconds")
	}
}

func TestSpells197RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_197.json")
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

func TestSpells197WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_197.json")
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
		`"choice": "clenched-fist", "choice": "unknown"`,
		`"choice": "clenched-fist", "CHOICE": "unknown"`,
		`"choice": "clenched-fist", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"choice": "clenched-fist"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

func TestSpells197RejectContradictoryTargetScaling(t *testing.T) {
	for _, index := range []int{10, 11} {
		for _, invalid := range []string{"zero increment", "mismatched base"} {
			t.Run(fmt.Sprintf("%d/%s", index, invalid), func(t *testing.T) {
				r := loadSpells197(t)
				s := &r.Records[index]
				if d := s.SourceDetails.Bane; d != nil {
					if invalid == "zero increment" {
						d.AdditionalTargetsPerSlot = 0
					} else {
						d.BaseTargets = 99
					}
				} else {
					d := s.SourceDetails.Banishment
					if invalid == "zero increment" {
						d.AdditionalTargetsPerSlot = 0
					} else {
						d.BaseTargets = 99
					}
				}
				if err := r.Validate(); err == nil {
					t.Error("contradictory target scaling validated")
				}
				b, err := json.Marshal(r)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
					t.Error("contradictory target scaling decoded")
				}
			})
		}
	}
}

// The decision inventory names exactly these twelve identities and first pages
// for catalogue issue #197; no adjacent batch may leak into the fixture.
func TestSpells197MatchesDecisionInventory(t *testing.T) {
	b, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			Family         string `json:"family"`
			Name           string `json:"name"`
			SourcePage     int    `json:"source_page"`
			ID             string `json:"id"`
			CatalogueIssue int    `json:"catalogue_issue"`
		} `json:"records"`
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	type identity struct {
		id, name string
		page     int
	}
	var want, got []identity
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 197 {
			if r.Family != "spells" {
				t.Fatalf("issue #197 owns non-spell family %q", r.Family)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
		}
	}
	batch := loadSpells197(t)
	for _, r := range batch.Records {
		got = append(got, identity{r.ID, r.Name, r.SourcePage})
		if r.SourcePage < 109 || r.SourcePage > 112 {
			t.Fatalf("%s outside PDF pp.109-112", r.Name)
		}
	}
	if len(want) != 12 || !slices.Equal(got, want) {
		t.Fatalf("fixture identities differ from decision inventory:\n got %v\nwant %v", got, want)
	}
	for i, r := range batch.Records {
		if r.Name != spells197Names[i] {
			t.Fatalf("family %d is %q, want %q", i, r.Name, spells197Names[i])
		}
	}
	src := batch.Source
	if src.SHA256 != artifactSHA256 || src.License != "CC-BY-4.0" || src.Attribution != Attribution || !strings.Contains(src.Changes, "issue #197") || !strings.Contains(src.Changes, "adapted") {
		t.Fatal("source pin, license, attribution or adaptation notice differs")
	}
}

// Continuation text and source-listed variants stay explicit data.
func TestSpells197VariantsAndContinuations(t *testing.T) {
	batch := loadSpells197(t)
	d := func(i int) *SpellSourceDetails { return batch.Records[i].SourceDetails }
	if !slices.Equal(d(0).Antipathy.SourcePages, []int{109, 110}) || !slices.Equal(d(5).MagicAura.SourcePages, []int{110, 111}) || !slices.Equal(d(9).Awaken.SourcePages, []int{111, 112}) {
		t.Fatal("continuation pages lost")
	}
	if len(d(0).Antipathy.Modes) != 2 || len(d(2).Hand.Options) != 4 || len(d(5).MagicAura.Variants) != 2 || len(d(7).Augury.Omens) != 4 || len(d(9).Awaken.TargetChoices) != 3 || len(d(11).Banishment.NoReturnTypes) != 5 {
		t.Fatal("source-listed variant missing")
	}
	seconds := map[string]int{"Antipathy/Sympathy": 3600, "Astral Projection": 3600, "Augury": 60, "Awaken": 28800}
	for _, r := range batch.Records {
		if r.SourceDetails.CastingSeconds != seconds[r.Name] {
			t.Fatalf("%s casting seconds differ", r.Name)
		}
		if *r.Ritual != (r.Name == "Augury") {
			t.Fatalf("%s ritual flag differs", r.Name)
		}
	}
}
