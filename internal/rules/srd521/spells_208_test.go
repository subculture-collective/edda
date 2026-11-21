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

func loadSpells208(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_208.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.136-139 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells208SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/globe-of-invulnerability","name":"Globe of Invulnerability","source_id":"srd-5.2.1:content/spells/globe-of-invulnerability","source_page":136,"scope":"complete-source-record","level":6,"school":"abjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"area-within-barrier"},"components":{"verbal":true,"somatic":true,"material":{"description":"a glass bead","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-limit","above_level":6,"dice_count":0,"die_sides":0,"effect_ids":["barrier"]},"effects":[{"id":"barrier","mechanic":"barrier.globe-of-invulnerability","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["barrier.globe-of-invulnerability","scaling.barrier-level-per-slot","concentration.source-duration","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","globe":{"emanation_feet":10,"barrier_mobility":"immobile","barrier_appearance":"shimmering","blocked_spell_max_level":5,"blocked_spell_origin":"cast-from-outside-barrier","spell_targeting":"may-target-creatures-and-objects-within-but-no-effect","area_exclusion":"area-within-barrier-excluded-from-areas-of-effect","slot_blocked_level_increment":1}}}`,
		`{"id":"srd-5.2.1:content/spells/glyph-of-warding","name":"Glyph of Warding","source_id":"srd-5.2.1:content/spells/glyph-of-warding","source_page":136,"scope":"complete-source-record","level":3,"school":"abjuration","lists":["bard","cleric","wizard"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"choice","count":1,"eligibility":"surface-or-closable-object"},"components":{"verbal":true,"somatic":true,"material":{"description":"powdered diamond worth 200+ GP, which the spell consumes","cost_cp":20000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":8,"effect_ids":["explosion"]},"effects":[{"id":"explosion","mechanic":"damage.glyph-explosive-rune","timing":"when-glyph-triggered","dice_count":5,"die_sides":8,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"stored-spell","mechanic":"storage.glyph-spell","timing":"when-glyph-triggered","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.glyph-explosive-rune","storage.glyph-spell","scaling.per-slot","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"dispelled-or-triggered","glyph":{"source_pages":[136,137],"inscribe_locations":["surface","within-closable-object"],"max_diameter_feet":10,"break_distance_feet":10,"break_result":"glyph-broken-and-spell-ends-without-triggering","detection_check":"wisdom-perception","detection_dc":"caster-spell-save-dc","trigger_set":"when-glyph-inscribed","trigger_types":["explosive-rune","spell-glyph"],"surface_triggers":["touching-or-stepping-on-glyph","removing-object-covering-glyph","approaching-within-distance"],"object_triggers":["opening-object","seeing-glyph"],"trigger_ends_spell":true,"trigger_refinements":["only-chosen-creature-types","exempt-creatures-meeting-conditions"],"refinement_examples":["aberrations","password"],"rune_shape":"sphere","rune_radius_feet":20,"rune_save":"dexterity","rune_dice_count":5,"rune_die_sides":8,"rune_damage_types":["acid","cold","fire","lightning","thunder"],"rune_damage_choice":"when-glyph-created","rune_success":"half-damage","spell_max_level":3,"spell_prepared":true,"spell_targets":["single-creature","area"],"spell_stored_effect":"no-immediate-effect","spell_target":"creature-that-triggered-glyph","spell_area_center":"creature-that-triggered-glyph","spell_hostile_summons":"appear-closest-to-intruder-and-attack-it","spell_concentration":"lasts-until-end-of-full-duration","slot_rune_dice_count":1,"slot_rune_die_sides":8,"slot_spell_level":"up-to-slot-level"}}}`,
		`{"id":"srd-5.2.1:content/spells/goodberry","name":"Goodberry","source_id":"srd-5.2.1:content/spells/goodberry","source_page":137,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a sprig of mistletoe","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"berries","mechanic":"create.goodberry-berries","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"eat","mechanic":"healing.goodberry-eat","timing":"bonus-action-eat-berry","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["create.goodberry-berries","healing.goodberry-eat","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","goodberry":{"berry_count":10,"berry_location":"caster-hand","eat_action":"bonus-action","berry_hit_points":1,"nourishment_days":1,"uneaten_berries":"disappear-when-spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/grease","name":"Grease","source_id":"srd-5.2.1:content/spells/grease","source_page":137,"scope":"complete-source-record","level":1,"school":"conjuration","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"square","count":0,"eligibility":"creatures-in-area-centered-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of pork rind or butter","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"terrain","mechanic":"terrain.grease-difficult","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"save","mechanic":"save.grease-dexterity-prone","timing":"appear-enter-or-end-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["terrain.grease-difficult","save.grease-dexterity-prone","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","grease":{"square_feet":10,"center":"point-within-range","flammability":"nonflammable","terrain":"difficult-terrain","save":"dexterity","save_triggers":["grease-appears-for-creature-standing-in-area","creature-enters-area","creature-ends-turn-in-area"],"appear_failure":"prone-condition","enter_or_end_failure":"fall-prone"}}}`,
		`{"id":"srd-5.2.1:content/spells/greater-invisibility","name":"Greater Invisibility","source_id":"srd-5.2.1:content/spells/greater-invisibility","source_page":137,"scope":"complete-source-record","level":4,"school":"illusion","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"invisible","mechanic":"condition.greater-invisibility-invisible","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.greater-invisibility-invisible","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","greater_invisibility":{"target":"creature-touched","condition":"invisible","condition_ends":"spell-ends"}}}`,
		`{"id":"srd-5.2.1:content/spells/greater-restoration","name":"Greater Restoration","source_id":"srd-5.2.1:content/spells/greater-restoration","source_page":137,"scope":"complete-source-record","level":5,"school":"abjuration","lists":["bard","cleric","druid","paladin","ranger"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"diamond dust worth 100+ GP, which the spell consumes","cost_cp":10000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"remove","mechanic":"remove.greater-restoration","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["remove.greater-restoration","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","greater_restoration":{"target":"creature-touched","choice_count":1,"exhaustion_levels_removed":1,"conditions_removed":["charmed","petrified"],"curse_includes_attunement":true,"ability_score_reduction":"any-reduction-to-one-ability-score","hit_point_maximum_reduction":"any-reduction"}}}`,
		`{"id":"srd-5.2.1:content/spells/guardian-of-faith","name":"Guardian of Faith","source_id":"srd-5.2.1:content/spells/guardian-of-faith","source_page":137,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-visible-space-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":28800,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"guardian","mechanic":"summon.guardian-of-faith-guardian","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ward","mechanic":"damage.guardian-of-faith-ward","timing":"enemy-moves-within-10-feet-first-time-or-starts-turn-there","dice_count":0,"die_sides":0,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["summon.guardian-of-faith-guardian","damage.guardian-of-faith-ward","save.dexterity-half","ending.guardian-total-damage"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-total-damage-dealt","guardian":{"guardian_size":"large","guardian_nature":"spectral","guardian_posture":"hovers","guardian_space":"unoccupied-space-visible-in-range","occupies_space":true,"invulnerable":true,"appearance":"form-appropriate-for-deity-or-pantheon","trigger_radius_feet":10,"trigger_targets":"enemy","triggers":["moves-within-radius-first-time-on-turn","starts-turn-within-radius"],"save":"dexterity","damage_flat":20,"damage_type":"radiant","save_success":"half-damage","vanishes_total_damage":60}}}`,
		`{"id":"srd-5.2.1:content/spells/guards-and-wards","name":"Guards and Wards","source_id":"srd-5.2.1:content/spells/guards-and-wards","source_page":137,"scope":"complete-source-record","level":6,"school":"abjuration","lists":["bard","wizard"],"casting_time":"1-hour","range":{"kind":"touch","feet":0},"target":{"shape":"created-effect","count":1,"eligibility":"warded-floor-area"},"components":{"verbal":true,"somatic":true,"material":{"description":"a silver rod worth 10+ GP","cost_cp":1000,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":86400,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"ward","mechanic":"ward.guards-and-wards-area","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"corridors","mechanic":"area.guards-and-wards-corridors","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"doors","mechanic":"lock.guards-and-wards-doors","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"stairs","mechanic":"web.guards-and-wards-stairs","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"other","mechanic":"spell.guards-and-wards-other-effect","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["ward.guards-and-wards-area","area.guards-and-wards-corridors","lock.guards-and-wards-doors","web.guards-and-wards-stairs","spell.guards-and-wards-other-effect","dispel.guards-and-wards-effects","extension.guards-and-wards-365-days","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":3600,"casting_triggers":[],"end_timing":"duration-ended-or-all-effects-dispelled","guards":{"source_pages":[137,138],"floor_area_square_feet":2500,"height_feet":20,"shapes":["one-50-foot-square","one-hundred-contiguous-5-foot-squares","twenty-five-contiguous-10-foot-squares"],"unaffected":"individuals-specified-when-cast","password_distance_feet":5,"password_result":"speaker-immune-to-effects","dispel_magic_on_spell":"no-effect","dispellable_effects":4,"all_dispelled":"spell-ends","extension_days":365,"extension_result":"lasts-until-all-effects-dispelled","corridors_obscurement":"heavily-obscured","corridors_misdirection_percent":50,"corridors_misdirection":"creature-other-than-caster-believes-opposite-direction","doors_lock":"arcane-lock","doors_illusion_count":10,"doors_illusion":"appear-as-plain-sections-of-wall","stairs_webs":"web-spell-top-to-bottom","stairs_regrow_minutes":10,"other_effect_count":1,"other_effects":["dancing-lights-four-corridors-repeating-program","magic-mouth-two-locations","stinking-cloud-two-locations-vapors-return-within-10-minutes","gust-of-wind-one-corridor-or-room-continuous","suggestion-one-5-foot-square-entering-creature-mentally"]}}}`,
		`{"id":"srd-5.2.1:content/spells/guidance","name":"Guidance","source_id":"srd-5.2.1:content/spells/guidance","source_page":138,"scope":"complete-source-record","level":0,"school":"divination","lists":["cleric","druid"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"check","mechanic":"bonus.guidance-ability-check","timing":"while-concentrating","dice_count":1,"die_sides":4,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["bonus.guidance-ability-check","concentration.source-duration"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","guidance":{"target":"willing-creature-touched","choice":"skill","bonus_dice_count":1,"bonus_die_sides":4,"applies_to":"ability-check-using-chosen-skill"}}}`,
		`{"id":"srd-5.2.1:content/spells/guiding-bolt","name":"Guiding Bolt","source_id":"srd-5.2.1:content/spells/guiding-bolt","source_page":138,"scope":"complete-source-record","level":1,"school":"evocation","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":120},"target":{"shape":"single","count":1,"eligibility":"creature-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":6,"scaling":{"kind":"slot","above_level":1,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"attack","mechanic":"attack.guiding-bolt-ranged-spell","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"damage","mechanic":"damage.guiding-bolt","timing":"on-hit","dice_count":4,"die_sides":6,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"advantage","mechanic":"advantage.guiding-bolt-next-attack","timing":"before-end-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["attack.guiding-bolt-ranged-spell","damage.guiding-bolt","advantage.guiding-bolt-next-attack","scaling.per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"end-of-caster-next-turn","guiding_bolt":{"attack":"ranged-spell-attack","target":"creature-in-range","damage_dice_count":4,"damage_die_sides":6,"damage_type":"radiant","advantage_on":"next-attack-roll-against-target","advantage_window":"before-end-of-caster-next-turn","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/gust-of-wind","name":"Gust of Wind","source_id":"srd-5.2.1:content/spells/gust-of-wind","source_page":138,"scope":"complete-source-record","level":2,"school":"evocation","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"line","count":0,"eligibility":"creatures-in-line"},"components":{"verbal":true,"somatic":true,"material":{"description":"a legume seed","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"gust","mechanic":"push.gust-of-wind","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"movement","mechanic":"movement.gust-closer-double-cost","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"environment","mechanic":"environment.gust-of-wind-gas-and-flames","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"redirect","mechanic":"direction.gust-bonus-action","timing":"bonus-action-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["push.gust-of-wind","movement.gust-closer-double-cost","environment.gust-of-wind-gas-and-flames","direction.gust-bonus-action","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","gust":{"line_length_feet":60,"line_width_feet":10,"direction":"chosen-by-caster","save":"strength","push_feet":15,"push_direction":"away-from-caster-along-line","end_turn_save":"same-save-when-creature-ends-turn-in-line","approach_movement_cost_feet":2,"approach_movement_per_feet":1,"disperses":"gas-or-vapor","extinguishes":"candles-and-similar-unprotected-flames","protected_flames_example":"lanterns","protected_flames_effect":"dance-wildly","protected_flames_extinguish_percent":50,"redirect_action":"bonus-action-on-later-turns"}}}`,
		`{"id":"srd-5.2.1:content/spells/hallow","name":"Hallow","source_id":"srd-5.2.1:content/spells/hallow","source_page":138,"scope":"complete-source-record","level":5,"school":"abjuration","lists":["cleric"],"casting_time":"24-hours","range":{"kind":"touch","feet":0},"target":{"shape":"created-effect","count":0,"eligibility":"area-around-touched-point"},"components":{"verbal":true,"somatic":true,"material":{"description":"incense worth 1,000+ GP, which the spell consumes","cost_cp":100000,"consumed":true}},"ritual":false,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"area","mechanic":"area.hallow","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"ward","mechanic":"ward.hallow-hallowed-ward","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"extra","mechanic":"effect.hallow-extra","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["area.hallow","ward.hallow-hallowed-ward","effect.hallow-extra","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":86400,"casting_triggers":[],"end_timing":"dispelled","hallow":{"source_pages":[138,139],"touch_target":"point","power":"holy-or-unholy","max_radius_feet":60,"overlap_result":"spell-fails-if-radius-includes-area-under-hallow","ward_types":["aberration","celestial","elemental","fey","fiend","undead"],"ward_choice":"any-of-listed-types","ward_entry":"chosen-types-cannot-willingly-enter","ward_possession_and_conditions":["possessed","charmed","frightened"],"ward_possession_result":"not-possessed-charmed-or-frightened-by-chosen-types-while-in-area","extra_effect_count":1,"extra_effects":["courage","darkness","daylight","peaceful-rest","extradimensional-interference","fear","resistance","silence","tongues","vulnerability"],"courage":"chosen-types-cannot-gain-frightened-condition-in-area","darkness":"darkness-fills-area-normal-light-and-magical-light-from-lower-level-spells-cannot-illuminate","daylight":"bright-light-fills-area-magical-darkness-from-lower-level-spells-cannot-extinguish","peaceful_rest":"interred-dead-bodies-cannot-be-turned-into-undead","extradimensional_interference":"chosen-types-cannot-enter-or-exit-by-teleportation-or-interplanar-travel","fear":"chosen-types-have-frightened-condition-in-area","resistance":"chosen-types-have-resistance-to-one-chosen-damage-type-in-area","silence":"no-sound-emanates-from-or-reaches-into-area","tongues":"chosen-types-communicate-with-any-creature-in-area-without-shared-language","vulnerability":"chosen-types-have-vulnerability-to-one-chosen-damage-type-in-area"}}}`,
	}
	batch := loadSpells208(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Globe of Invulnerability", "Glyph of Warding", "Goodberry", "Grease", "Greater Invisibility", "Greater Restoration", "Guardian of Faith", "Guards and Wards", "Guidance", "Guiding Bolt", "Gust of Wind", "Hallow"}
	pages := []int{136, 136, 137, 137, 137, 137, 137, 137, 138, 138, 138, 138}
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
			// Guidance is a cantrip: it has no slot levels at all.
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

func TestSpells208CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells208(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), loadSpells200(t), loadSpells201(t), loadSpells202(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 108 {
		t.Fatalf("expected nine twelve-family batches, got %d", len(all.Records))
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
func TestSpells208MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_208.json")
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

func TestSpells208JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_208.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"emanation_feet": 10`, `"emanation_feet": 10, "emanation_feet": 11`},
		{`"emanation_feet": 10`, `"emanation_feet": 10, "EMANATION_FEET": 11`},
		{`"emanation_feet": 10`, `"emanation_feet": 10, "unknown": true`},
		{`"rune_damage_types": [
            "acid",
            "cold",
            "fire",
            "lightning",
            "thunder"
          ]`, `"rune_damage_types": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/goodberry"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "explosion"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "24-hours"`, `"casting_time": "8-hours"`},
		{`"casting_seconds": 86400`, `"casting_seconds": 3600`},
		{`"feet": 120`, `"feet": 121`},
		{`"cost_cp": 100000`, `"cost_cp": 100001`},
		{`"concentration": true`, `"concentration": false`},
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

func TestSpells208InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"globe radius", 0, func(s *SpellRecord) { s.SourceDetails.Globe.EmanationFeet = 15 }},
		{"globe blocked level", 0, func(s *SpellRecord) { s.SourceDetails.Globe.BlockedSpellMaxLevel = 6 }},
		{"globe scaling kind", 0, func(s *SpellRecord) { s.Scaling.Kind = "none" }},
		{"globe scaling increment", 0, func(s *SpellRecord) { s.SourceDetails.Globe.SlotBlockedLevelIncrement = 2 }},
		{"globe duration", 0, func(s *SpellRecord) { s.DurationSeconds = 600 }},
		{"glyph continuation page", 1, func(s *SpellRecord) { s.SourceDetails.Glyph.SourcePages = []int{136} }},
		{"glyph rune damage", 1, func(s *SpellRecord) { s.SourceDetails.Glyph.RuneDiceCount = 4 }},
		{"glyph damage type", 1, func(s *SpellRecord) { s.SourceDetails.Glyph.RuneDamageTypes[4] = "necrotic" }},
		{"glyph stored level", 1, func(s *SpellRecord) { s.SourceDetails.Glyph.SpellMaxLevel = 4 }},
		{"glyph material cost", 1, func(s *SpellRecord) { *s.Components.Material.CostCP = 2000 }},
		{"glyph consumed", 1, func(s *SpellRecord) { *s.Components.Material.Consumed = false }},
		{"glyph casting time", 1, func(s *SpellRecord) { s.CastingTime = "10-minutes" }},
		{"glyph scaling die", 1, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"goodberry count", 2, func(s *SpellRecord) { s.SourceDetails.Goodberry.BerryCount = 12 }},
		{"goodberry duration", 2, func(s *SpellRecord) { s.DurationSeconds = 3600 }},
		{"grease square", 3, func(s *SpellRecord) { s.SourceDetails.Grease.SquareFeet = 20 }},
		{"grease trigger", 3, func(s *SpellRecord) { s.SourceDetails.Grease.SaveTriggers[1] = "unknown" }},
		{"greater invisibility condition", 4, func(s *SpellRecord) { s.SourceDetails.GreaterInvisibility.Condition = "hidden" }},
		{"greater invisibility concentration", 4, func(s *SpellRecord) { *s.Concentration = false }},
		{"restoration condition", 5, func(s *SpellRecord) { s.SourceDetails.GreaterRestoration.ConditionsRemoved[1] = "stunned" }},
		{"restoration attunement", 5, func(s *SpellRecord) { s.SourceDetails.GreaterRestoration.CurseIncludesAttune = false }},
		{"restoration material", 5, func(s *SpellRecord) { *s.Components.Material.CostCP = 20000 }},
		{"guardian damage", 6, func(s *SpellRecord) { s.SourceDetails.Guardian.DamageFlat = 10 }},
		{"guardian vanish", 6, func(s *SpellRecord) { s.SourceDetails.Guardian.VanishesTotalDamage = 50 }},
		{"guardian components", 6, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"guards area", 7, func(s *SpellRecord) { s.SourceDetails.Guards.FloorAreaSquareFeet = 2000 }},
		{"guards continuation page", 7, func(s *SpellRecord) { s.SourceDetails.Guards.SourcePages = []int{137} }},
		{"guards other effect", 7, func(s *SpellRecord) { s.SourceDetails.Guards.OtherEffects[3] = "unknown" }},
		{"guards extension", 7, func(s *SpellRecord) { s.SourceDetails.Guards.ExtensionDays = 300 }},
		{"guidance die", 8, func(s *SpellRecord) { s.SourceDetails.Guidance.BonusDieSides = 6 }},
		{"guidance cantrip slots", 8, func(s *SpellRecord) { s.SlotLevels = []int{0} }},
		{"guiding bolt damage", 9, func(s *SpellRecord) { s.Effects[1].DiceCount = 3 }},
		{"guiding bolt advantage window", 9, func(s *SpellRecord) { s.SourceDetails.GuidingBolt.AdvantageWindow = "unknown" }},
		{"guiding bolt scaling", 9, func(s *SpellRecord) { s.Scaling.AboveLevel = 2 }},
		{"gust line", 10, func(s *SpellRecord) { s.SourceDetails.Gust.LineLengthFeet = 30 }},
		{"gust push", 10, func(s *SpellRecord) { s.SourceDetails.Gust.PushFeet = 10 }},
		{"gust lantern chance", 10, func(s *SpellRecord) { s.SourceDetails.Gust.ProtectedFlamesPercent = 25 }},
		{"hallow continuation page", 11, func(s *SpellRecord) { s.SourceDetails.Hallow.SourcePages = []int{138} }},
		{"hallow radius", 11, func(s *SpellRecord) { s.SourceDetails.Hallow.MaxRadiusFeet = 30 }},
		{"hallow extra effect", 11, func(s *SpellRecord) { s.SourceDetails.Hallow.ExtraEffects[9] = "unknown" }},
		{"hallow ward type", 11, func(s *SpellRecord) { s.SourceDetails.Hallow.WardTypes[5] = "humanoid" }},
		{"hallow casting seconds", 11, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 3600 }},
		{"hallow material cost", 11, func(s *SpellRecord) { *s.Components.Material.CostCP = 10000 }},
		{"hallow list", 11, func(s *SpellRecord) { s.Lists = []string{"cleric", "druid"} }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Gust = &SpellGust{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells208(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells208RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_208.json")
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

func TestSpells208WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_208.json")
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
		`"berry_count": 10, "berry_count": 12`,
		`"berry_count": 10, "BERRY_COUNT": 12`,
		`"berry_count": 10, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"berry_count": 10`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells208PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_208.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"emanation_feet", 1, 10}, {"blocked_spell_max_level", 1, 5}, {"slot_blocked_level_increment", 1, 1},
		{"max_diameter_feet", 1, 10}, {"break_distance_feet", 1, 10}, {"rune_radius_feet", 1, 20}, {"rune_dice_count", 1, 5},
		{"rune_die_sides", 1, 8}, {"spell_max_level", 1, 3}, {"berry_count", 1, 10}, {"berry_hit_points", 1, 1},
		{"square_feet", 1, 10}, {"trigger_radius_feet", 1, 10}, {"damage_flat", 1, 20}, {"vanishes_total_damage", 1, 60},
		{"floor_area_square_feet", 1, 2500}, {"height_feet", 1, 20}, {"password_distance_feet", 1, 5},
		{"extension_days", 1, 365}, {"corridors_misdirection_percent", 1, 50}, {"doors_illusion_count", 1, 10},
		{"stairs_regrow_minutes", 1, 10}, {"bonus_die_sides", 1, 4}, {"damage_dice_count", 1, 4},
		{"line_length_feet", 1, 60}, {"line_width_feet", 1, 10}, {"push_feet", 1, 15},
		{"protected_flames_extinguish_percent", 1, 50}, {"max_radius_feet", 1, 60},
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
// inventory assigns to catalogue issue #208.
func TestSpells208MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 208 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells208(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells208ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_208.json")
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
	want := "SRD 5.2.1 PDF pp.136-139 adapted into twelve complete offline spell families for issue #208, including page-break continuation text for Glyph of Warding (pp.136-137), Guards and Wards (pp.137-138) and Hallow (pp.138-139). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
	if strings.Contains(root.Source.Changes, "Gust of Wind") {
		t.Fatal("Gust of Wind sits wholly on p.138 and is not a page-break continuation")
	}
}

// Hallow's page 138 text says only "an area around it" with a radius up to 60
// feet, so the target shape must stay source-neutral and never read "sphere".
func TestSpells208HallowAreaIsNotASphere(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_208.json")
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		t.Fatal(err)
	}
	for _, record := range root["records"].([]any) {
		r := record.(map[string]any)
		if r["name"] != "Hallow" {
			continue
		}
		target := r["target"].(map[string]any)
		if target["shape"] != "created-effect" {
			t.Fatalf("Hallow target shape %v", target["shape"])
		}
		target["shape"] = "sphere"
		bad, err := json.Marshal(root)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeSpellRecords(bytes.NewReader(bad)); err == nil {
			t.Fatal("accepted Hallow target shape sphere")
		}
		return
	}
	t.Fatal("Hallow record missing")
}
