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
	"testing"
)

func loadSpells200(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_200.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.116-118 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells200SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/command","name":"Command","source_id":"srd-5.2.1:content/spells/command","source_page":116,"scope":"complete-source-record","level":1,"school":"enchantment","lists":["bard","cleric","paladin"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"up-to-count","count":1,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["command"]},"effects":[{"id":"command","mechanic":"control.command-one-word","timing":"failed-save-follows-on-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["control.command-one-word","save.wisdom-negates","target.visible-creature","choice.command-option","scaling.targets-per-slot"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","command":{"words":1,"save":"wisdom","save_success":"no-effect","follow_timing":"target-next-turn","options":[{"choice":"approach","effect":"moves-toward-caster-by-shortest-most-direct-route","condition":"not-applicable","ends_turn":"if-it-moves-within-5-feet-of-caster"},{"choice":"drop","effect":"drops-whatever-it-is-holding","condition":"not-applicable","ends_turn":"after-dropping"},{"choice":"flee","effect":"spends-turn-moving-away-from-caster-by-fastest-available-means","condition":"not-applicable","ends_turn":"not-stated"},{"choice":"grovel","effect":"has-prone-condition","condition":"prone","ends_turn":"after-becoming-prone"},{"choice":"halt","effect":"does-not-move-and-takes-no-action-or-bonus-action","condition":"not-applicable","ends_turn":"not-stated"}],"base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/commune","name":"Commune","source_id":"srd-5.2.1:content/spells/commune","source_page":116,"scope":"complete-source-record","level":5,"school":"divination","lists":["cleric"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"incense","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"answers","mechanic":"divination.commune-yes-no","timing":"ask-before-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"no-answer","mechanic":"divination.commune-cumulative-no-answer","timing":"each-casting-after-first-before-long-rest","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["divination.commune-yes-no","divination.commune-cumulative-no-answer","target.deity-or-divine-proxy","cast.one-minute","cast.ritual","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"duration-ended","commune":{"contact":"deity-or-divine-proxy","question_limit":3,"question_kind":"yes-or-no","ask_timing":"before-spell-ends","answer":"correct-for-each-question","unclear_answer":"unclear-if-beyond-deity-knowledge","short_phrase_answer":"gm-may-give-short-phrase-if-one-word-misleading-or-contrary-to-deity-interests","repeat_cumulative_percent":25,"repeat_applies_from_cast":2,"repeat_resets_on":"long-rest","repeat_effect":"no-answer","ritual_alternative":true}}}`,
		`{"id":"srd-5.2.1:content/spells/commune-with-nature","name":"Commune with Nature","source_id":"srd-5.2.1:content/spells/commune-with-nature","source_page":116,"scope":"complete-source-record","level":5,"school":"divination","lists":["druid","ranger"],"casting_time":"1-minute","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"facts","mechanic":"knowledge.nature-area-facts","timing":"on-cast-choose-three-facts","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["knowledge.nature-area-facts","choice.three-area-facts","area.outdoor-3-miles-underground-300-feet","environment.nature-not-replaced-by-construction","cast.one-minute","cast.ritual"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"instantaneous","nature":{"source_pages":[116,117],"outdoor_radius_miles":3,"underground_radius_feet":300,"underground_setting":"caves-and-other-natural-underground-settings","fails_where":"nature-replaced-by-construction","failure_examples":["castles","settlements"],"choose_count":3,"facts":["settlement-locations","portal-locations-to-other-planes","location-of-one-challenge-rating-10-plus-celestial-elemental-fey-fiend-or-undead","most-prevalent-plant-mineral-or-beast","body-of-water-locations"],"creature_minimum_cr":10,"creature_types":["celestial","elemental","fey","fiend","undead"],"creature_choice":"gm","prevalent_kinds":["plant","mineral","beast"],"prevalent_choice":"caster","ritual_alternative":true}}}`,
		`{"id":"srd-5.2.1:content/spells/comprehend-languages","name":"Comprehend Languages","source_id":"srd-5.2.1:content/spells/comprehend-languages","source_page":117,"scope":"complete-source-record","level":1,"school":"divination","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"a pinch of soot and salt","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"understand","mechanic":"language.comprehend-literal-meaning","timing":"while-active","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["language.comprehend-literal-meaning","language.written-requires-touch","limit.no-decoding","cast.ritual","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","comprehend":{"understands":["literal-meaning-of-language-heard-or-seen-signed","literal-meaning-of-written-language-seen"],"written_requires":"touching-surface-with-words","reading_minutes_per_page":1,"decodes_symbols":false,"decodes_secret_messages":false,"ritual_alternative":true}}}`,
		`{"id":"srd-5.2.1:content/spells/compulsion","name":"Compulsion","source_id":"srd-5.2.1:content/spells/compulsion","source_page":117,"scope":"complete-source-record","level":4,"school":"enchantment","lists":["bard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"selected-creatures","count":0,"eligibility":"visible-creatures-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"charmed","mechanic":"condition.compulsion-charmed","timing":"failed-save-until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"direction","mechanic":"control.compulsion-direction","timing":"bonus-action-each-turn-then-charmed-target-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.compulsion-charmed","control.compulsion-direction","save.wisdom-negates","save.repeat-after-forced-move","target.visible-creatures-of-choice","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","compulsion":{"target_choice":"creatures-of-choice-caster-can-see-in-range","save":"wisdom","condition":"charmed","direction_action":"bonus-action","direction":"horizontal-to-caster","movement":"uses-as-much-movement-as-possible-in-direction-on-next-turn","route":"safest","repeat_save":"after-moving-in-direction","repeat_success":"ends-spell-on-itself"}}}`,
		`{"id":"srd-5.2.1:content/spells/cone-of-cold","name":"Cone of Cold","source_id":"srd-5.2.1:content/spells/cone-of-cold","source_page":117,"scope":"complete-source-record","level":5,"school":"evocation","lists":["druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"cone","count":0,"eligibility":"creatures-in-cone"},"components":{"verbal":true,"somatic":true,"material":{"description":"a small crystal or glass cone","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":8,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.cone-of-cold","timing":"on-cast-save-half-on-success","dice_count":8,"die_sides":8,"damage_type":"cold","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"statue","mechanic":"state.frozen-statue-until-thaws","timing":"creature-killed-by-spell","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.cone-of-cold","state.frozen-statue-until-thaws","area.cone-60","save.constitution-half","scaling.per-slot","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","cone":{"cone_feet":60,"origin":"caster","save":"constitution","save_success":"half-damage","killed_state":"frozen-statue-until-it-thaws"}}}`,
		`{"id":"srd-5.2.1:content/spells/confusion","name":"Confusion","source_id":"srd-5.2.1:content/spells/confusion","source_page":117,"scope":"complete-source-record","level":4,"school":"enchantment","lists":["bard","druid","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"sphere","count":0,"eligibility":"creatures-in-sphere"},"components":{"verbal":true,"somatic":true,"material":{"description":"three nut shells","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-radius","above_level":4,"dice_count":0,"die_sides":0,"effect_ids":["behavior"]},"effects":[{"id":"behavior","mechanic":"condition.confusion-random-behavior","timing":"failed-save-start-of-each-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.confusion-random-behavior","area.sphere-10-radius","save.wisdom-negates","save.repeat-end-of-turn","restriction.no-bonus-actions-or-reactions","table.confusion-behavior","scaling.radius-per-slot","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","confusion":{"radius_feet":10,"center":"chosen-point-in-range","save":"wisdom","save_success":"no-effect","restrictions":["cannot-take-bonus-actions","cannot-take-reactions"],"behavior_die_sides":10,"behavior_timing":"start-of-each-affected-turn","behaviors":[{"roll_min":1,"roll_max":1,"behavior":"takes-no-action-uses-all-movement-in-random-direction"},{"roll_min":2,"roll_max":6,"behavior":"does-not-move-or-take-actions"},{"roll_min":7,"roll_max":8,"behavior":"does-not-move-and-takes-attack-action-for-one-melee-attack-against-random-creature-within-reach-else-takes-no-action"},{"roll_min":9,"roll_max":10,"behavior":"target-chooses-its-behavior"}],"direction_die_sides":4,"directions":["north","east","south","west"],"repeat_save":"end-of-each-affected-turn","repeat_success":"ends-spell-on-itself","additional_radius_feet_per_slot":5}}}`,
		`{"id":"srd-5.2.1:content/spells/conjure-animals","name":"Conjure Animals","source_id":"srd-5.2.1:content/spells/conjure-animals","source_page":117,"scope":"complete-source-record","level":3,"school":"conjuration","lists":["druid","ranger"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":10,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.conjure-animals-pack","timing":"failed-save-when-pack-moves-near-or-creature-enters-or-ends-turn-near-pack","dice_count":3,"die_sides":10,"damage_type":"slashing","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"pack-move","mechanic":"movement.conjure-animals-pack","timing":"when-caster-moves-on-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"advantage","mechanic":"defense.conjure-animals-strength-advantage","timing":"caster-within-5-feet-of-pack","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.conjure-animals-pack","movement.conjure-animals-pack","defense.conjure-animals-strength-advantage","summon.spectral-pack-large","save.dexterity-no-damage","save.once-per-turn","scaling.per-slot","concentration.source-duration"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","animals":{"source_pages":[117,118],"size":"large","form":"pack-of-spectral-intangible-animals","placement":"unoccupied-space-visible-in-range","form_choice":"caster-chooses-animal-form","form_examples":["wolves","serpents","birds"],"save_advantage":"strength-saving-throws-while-within-5-feet-of-pack","advantage_range_feet":5,"move_feet":30,"move_timing":"when-caster-moves-on-turn","move_destination":"unoccupied-space-caster-can-see","trigger_range_feet":10,"triggers":["pack-moves-within-10-feet-of-visible-creature","visible-creature-enters-space-within-10-feet-of-pack","visible-creature-ends-turn-within-10-feet-of-pack"],"save":"dexterity","save_success":"no-effect","saves_per_turn":1}}}`,
		`{"id":"srd-5.2.1:content/spells/conjure-celestial","name":"Conjure Celestial","source_id":"srd-5.2.1:content/spells/conjure-celestial","source_page":118,"scope":"complete-source-record","level":7,"school":"conjuration","lists":["cleric"],"casting_time":"action","range":{"kind":"feet","feet":90},"target":{"shape":"created-effect","count":1,"eligibility":"point-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":7,"dice_count":1,"die_sides":12,"effect_ids":["healing","damage"]},"effects":[{"id":"healing","mechanic":"healing.conjure-celestial-healing-light","timing":"bathe-visible-creature-in-healing-light","dice_count":4,"die_sides":12,"damage_type":"","bonus":0,"add_casting_modifier":true,"half_initial_only":false},{"id":"damage","mechanic":"damage.conjure-celestial-searing-light","timing":"bathe-visible-creature-in-searing-light-save-half-on-success","dice_count":6,"die_sides":12,"damage_type":"radiant","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"light","mechanic":"light.conjure-celestial-bright-light","timing":"while-cylinder-lasts","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["healing.conjure-celestial-healing-light","damage.conjure-celestial-searing-light","light.conjure-celestial-bright-light","area.cylinder-10-radius-40-high","summon.pillar-of-light","choice.light-per-creature","save.dexterity-half","scaling.per-slot","concentration.source-duration"],"slot_levels":[7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","celestial":{"appearance":"pillar-of-light","radius_feet":10,"height_feet":40,"center":"point-in-range","choice_per_target":"each-creature-caster-can-see-in-cylinder","lights":[{"choice":"healing-light","effect":"regains-hit-points","save":"not-applicable","save_success":"not-applicable"},{"choice":"searing-light","effect":"radiant-damage","save":"dexterity","save_success":"half-damage"}],"illumination":"bright-light-fills-cylinder","move_feet":30,"move_timing":"when-caster-moves-on-turn","triggers":["cylinder-moves-into-space-of-visible-creature","visible-creature-enters-cylinder","visible-creature-ends-turn-in-cylinder"],"affected_per_turn":1}}}`,
		`{"id":"srd-5.2.1:content/spells/conjure-elemental","name":"Conjure Elemental","source_id":"srd-5.2.1:content/spells/conjure-elemental","source_page":118,"scope":"complete-source-record","level":5,"school":"conjuration","lists":["druid","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":5,"dice_count":1,"die_sides":8,"effect_ids":["damage","repeat-damage"]},"effects":[{"id":"damage","mechanic":"damage.conjure-elemental-initial","timing":"spirit-has-no-restrained-creature-failed-save","dice_count":8,"die_sides":8,"damage_type":"spirit-element","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"restrained","mechanic":"condition.conjure-elemental-restrained","timing":"failed-save-until-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"repeat-damage","mechanic":"damage.conjure-elemental-repeat","timing":"start-of-restrained-target-turn-failed-save","dice_count":4,"die_sides":8,"damage_type":"spirit-element","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.conjure-elemental-initial","damage.conjure-elemental-repeat","condition.conjure-elemental-restrained","summon.elemental-spirit-large","choice.element-damage-type","save.dexterity-negates","save.repeat-start-of-turn","scaling.per-slot","concentration.source-duration"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","elemental":{"size":"large","form":"intangible-spirit-from-elemental-planes","placement":"unoccupied-space-in-range","elements":[{"element":"air","damage_type":"lightning"},{"element":"earth","damage_type":"thunder"},{"element":"fire","damage_type":"fire"},{"element":"water","damage_type":"cold"}],"triggers":["visible-creature-enters-spirit-space","visible-creature-starts-turn-within-5-feet-of-spirit"],"trigger_range_feet":5,"trigger_condition":"spirit-has-no-creature-restrained","save":"dexterity","condition":"restrained","condition_ends":"spell-ends","repeat_save":"start-of-each-restrained-target-turn","repeat_success":"not-restrained-by-spirit"}}}`,
		`{"id":"srd-5.2.1:content/spells/conjure-fey","name":"Conjure Fey","source_id":"srd-5.2.1:content/spells/conjure-fey","source_page":118,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["druid"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"created-effect","count":1,"eligibility":"visible-unoccupied-space-in-range"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":6,"dice_count":1,"die_sides":12,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.conjure-fey-psychic","timing":"melee-spell-attack-hit","dice_count":3,"die_sides":12,"damage_type":"psychic","bonus":0,"add_casting_modifier":true,"half_initial_only":false},{"id":"fear","mechanic":"condition.conjure-fey-frightened","timing":"melee-spell-attack-hit-until-start-of-caster-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"teleport","mechanic":"movement.conjure-fey-teleport","timing":"bonus-action-on-later-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.conjure-fey-psychic","condition.conjure-fey-frightened","movement.conjure-fey-teleport","summon.fey-spirit-medium","attack.melee-spell","action.bonus-action-repeat-attack","scaling.per-slot","concentration.source-duration"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","fey":{"size":"medium","origin":"feywild","form":"looks-like-fey-creature-of-caster-choice","placement":"unoccupied-space-visible-in-range","attack":"melee-spell-attack-against-creature-within-5-feet-of-spirit","attack_range_feet":5,"condition":"frightened","condition_ends":"start-of-caster-next-turn","fear_sources":["caster","spirit"],"later_action":"bonus-action","teleport_feet":30,"teleport_to":"unoccupied-space-caster-can-see-within-30-feet-of-space-spirit-left","later_attack":"attack-against-creature-within-5-feet-after-teleport"}}}`,
		`{"id":"srd-5.2.1:content/spells/conjure-minor-elementals","name":"Conjure Minor Elementals","source_id":"srd-5.2.1:content/spells/conjure-minor-elementals","source_page":118,"scope":"complete-source-record","level":4,"school":"conjuration","lists":["druid","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"emanation","count":0,"eligibility":"creatures-in-emanation"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"slot","above_level":4,"dice_count":1,"die_sides":8,"effect_ids":["extra-damage"]},"effects":[{"id":"extra-damage","mechanic":"damage.conjure-minor-elementals-extra","timing":"caster-attack-hits-creature-in-emanation","dice_count":2,"die_sides":8,"damage_type":"chosen-type","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"terrain","mechanic":"terrain.conjure-minor-elementals-difficult","timing":"while-emanation-lasts","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.conjure-minor-elementals-extra","terrain.conjure-minor-elementals-difficult","area.emanation-15","damage.chosen-type","summon.elemental-spirits","scaling.per-slot","concentration.source-duration"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","minor":{"origin":"elemental-planes","emanation_feet":15,"trigger":"caster-attack-hits-creature-in-emanation","damage_types":["acid","cold","fire","lightning"],"type_choice":"when-caster-makes-attack","terrain":"difficult-terrain","terrain_for":"caster-enemies"}}}`,
	}
	batch := loadSpells200(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Command", "Commune", "Commune with Nature", "Comprehend Languages", "Compulsion", "Cone of Cold", "Confusion", "Conjure Animals", "Conjure Celestial", "Conjure Elemental", "Conjure Fey", "Conjure Minor Elementals"}
	pages := []int{116, 116, 116, 117, 117, 117, 117, 117, 118, 118, 118, 118}
	// The exact family list and pages come from the issue-80 decision inventory.
	inventory, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var decision struct {
		Records []struct {
			ID, Name       string
			SourcePage     int    `json:"source_page"`
			CatalogueIssue int    `json:"catalogue_issue"`
			BatchKey       string `json:"batch_key"`
		} `json:"records"`
	}
	if err := json.Unmarshal(inventory, &decision); err != nil {
		t.Fatal(err)
	}
	listed := 0
	for _, entry := range decision.Records {
		if entry.CatalogueIssue != 200 {
			continue
		}
		if listed >= len(batch.Records) || entry.BatchKey != "data-spells-05" || entry.ID != batch.Records[listed].ID || entry.Name != names[listed] || entry.SourcePage != pages[listed] {
			t.Fatalf("inventory entry %q differs from the batch", entry.Name)
		}
		listed++
	}
	if listed != 12 {
		t.Fatalf("inventory lists %d families", listed)
	}
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

func TestSpells200CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells200(t)
	// #196 restates Acid Arrow, so the legacy fixture cannot be composed beside it.
	bounded := legacy
	bounded.Records = legacy.Records[1:]
	all, err := ComposeSpellSourceCatalogue(bounded, loadSpells196(t), loadSpells197(t), loadSpells218(t), loadSpells198(t), loadSpells199(t), batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 72 {
		t.Fatalf("expected six twelve-family batches, got %d", len(all.Records))
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
func TestSpells200MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_200.json")
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

func TestSpells200JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_200.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"cone_feet": 60`, `"cone_feet": 60, "cone_feet": 61`},
		{`"cone_feet": 60`, `"cone_feet": 60, "CONE_FEET": 61`},
		{`"cone_feet": 60`, `"cone_feet": 60, "unknown": true`},
		{"\"directions\": [\n            \"north\",\n            \"east\",\n            \"south\",\n            \"west\"\n          ]", `"directions": [ ]`},
		{`"source_id": "srd-5.2.1:content/spells/command"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "command"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [
        1,`, `"slot_levels": [
        0,`},
		{`"casting_time": "1-minute"`, `"casting_time": "action"`},
		{`"casting_seconds": 60`, `"casting_seconds": 0`},
		{`"ritual": true`, `"ritual": false`},
		{`"duration_seconds": 3600`, `"duration_seconds": 600`},
		{"\"source_pages\": [\n            116,\n            117\n          ]", "\"source_pages\": [\n            116\n          ]"},
		{`"damage_type": "psychic"`, `"damage_type": "force"`},
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

func TestSpells200InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"command option", 0, func(s *SpellRecord) { s.SourceDetails.Command.Options[4].Choice = "unknown" }},
		{"command words", 0, func(s *SpellRecord) { s.SourceDetails.Command.Words = 2 }},
		{"command zero increment", 0, func(s *SpellRecord) { s.SourceDetails.Command.AdditionalTargetsPerSlot = 0 }},
		{"command base targets", 0, func(s *SpellRecord) { s.SourceDetails.Command.BaseTargets = 2 }},
		{"command scaling dependency", 0, func(s *SpellRecord) { s.RequiredMechanics[4] = "scaling.per-slot" }},
		{"commune question limit", 1, func(s *SpellRecord) { s.SourceDetails.Commune.QuestionLimit = 4 }},
		{"commune repeat chance", 1, func(s *SpellRecord) { s.SourceDetails.Commune.RepeatCumulativePct = 50 }},
		{"commune not ritual", 1, func(s *SpellRecord) { *s.Ritual = false }},
		{"commune casting seconds", 1, func(s *SpellRecord) { s.SourceDetails.CastingSeconds = 0 }},
		{"nature radius", 2, func(s *SpellRecord) { s.SourceDetails.Nature.UndergroundRadiusFt = 30 }},
		{"nature fact count", 2, func(s *SpellRecord) { s.SourceDetails.Nature.Facts = s.SourceDetails.Nature.Facts[:4] }},
		{"nature continuation page", 2, func(s *SpellRecord) { s.SourceDetails.Nature.SourcePages = []int{116} }},
		{"nature creature type", 2, func(s *SpellRecord) { s.SourceDetails.Nature.CreatureTypes[4] = "dragon" }},
		{"comprehend decodes", 3, func(s *SpellRecord) { s.SourceDetails.Comprehend.DecodesSecrets = true }},
		{"comprehend duration", 3, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"compulsion condition", 4, func(s *SpellRecord) { s.SourceDetails.Compulsion.Condition = "frightened" }},
		{"compulsion concentration", 4, func(s *SpellRecord) { *s.Concentration = false }},
		{"cone area", 5, func(s *SpellRecord) { s.SourceDetails.Cone.ConeFeet = 30 }},
		{"cone damage dice", 5, func(s *SpellRecord) { s.Effects[0].DiceCount = 6 }},
		{"cone scaling die", 5, func(s *SpellRecord) { s.Scaling.DieSides = 6 }},
		{"confusion table row", 6, func(s *SpellRecord) { s.SourceDetails.Confusion.Behaviors[1].RollMax = 5 }},
		{"confusion radius scaling", 6, func(s *SpellRecord) { s.SourceDetails.Confusion.AdditionalRadiusFeetSlot = 0 }},
		{"confusion scaling kind", 6, func(s *SpellRecord) { s.Scaling.Kind = "none" }},
		{"confusion direction", 6, func(s *SpellRecord) { s.SourceDetails.Confusion.Directions[3] = "up" }},
		{"animals pack range", 7, func(s *SpellRecord) { s.SourceDetails.Animals.TriggerRangeFeet = 5 }},
		{"animals continuation page", 7, func(s *SpellRecord) { s.SourceDetails.Animals.SourcePages = []int{117} }},
		{"animals trigger count", 7, func(s *SpellRecord) { s.SourceDetails.Animals.Triggers = s.SourceDetails.Animals.Triggers[:2] }},
		{"animals damage type", 7, func(s *SpellRecord) { s.Effects[0].DamageType = "piercing" }},
		{"celestial cylinder", 8, func(s *SpellRecord) { s.SourceDetails.Celestial.HeightFeet = 30 }},
		{"celestial light", 8, func(s *SpellRecord) { s.SourceDetails.Celestial.Lights[1].Choice = "unknown" }},
		{"celestial healing modifier", 8, func(s *SpellRecord) { *s.Effects[0].AddCastingModifier = false }},
		{"celestial scaling effects", 8, func(s *SpellRecord) { s.Scaling.EffectIDs = []string{"damage"} }},
		{"elemental element", 9, func(s *SpellRecord) { s.SourceDetails.Elemental.Elements[0].DamageType = "thunder" }},
		{"elemental restraint", 9, func(s *SpellRecord) { s.SourceDetails.Elemental.Condition = "grappled" }},
		{"elemental repeat dice", 9, func(s *SpellRecord) { s.Effects[2].DiceCount = 8 }},
		{"fey teleport", 10, func(s *SpellRecord) { s.SourceDetails.Fey.TeleportFeet = 60 }},
		{"fey fear sources", 10, func(s *SpellRecord) { s.SourceDetails.Fey.FearSources = []string{"caster"} }},
		{"fey modifier", 10, func(s *SpellRecord) { *s.Effects[0].AddCastingModifier = false }},
		{"minor emanation", 11, func(s *SpellRecord) { s.SourceDetails.Minor.EmanationFeet = 10 }},
		{"minor damage type", 11, func(s *SpellRecord) { s.SourceDetails.Minor.DamageTypes[3] = "thunder" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.Cone = &SpellCone{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Spray = &SpellSpray{} }},
		{"payload on other spell", 5, func(s *SpellRecord) { s.SourceDetails.Cone, s.SourceDetails.Minor = nil, &SpellMinor{} }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
		{"adopted target", 0, func(s *SpellRecord) { s.Target = SpellTarget{"selected-creatures", 0, "visible-creatures-in-range"} }},
		{"adopted range", 7, func(s *SpellRecord) { s.Range.Feet = 90 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells200(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells200RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_200.json")
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

func TestSpells200WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_200.json")
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
		`"choice": "approach", "choice": "unknown"`,
		`"choice": "approach", "CHOICE": "unknown"`,
		`"choice": "approach", "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"choice": "approach"`), []byte(replacement), 1)
		if _, err := DecodeSpellRecords(bytes.NewReader(b)); err == nil {
			t.Fatal("nested duplicate/alias/unknown accepted")
		}
	}
}

func TestSpells200RejectContradictoryTargetScaling(t *testing.T) {
	for _, invalid := range []string{"zero increment", "mismatched base"} {
		t.Run(invalid, func(t *testing.T) {
			r := loadSpells200(t)
			d := r.Records[0].SourceDetails.Command
			if invalid == "zero increment" {
				d.AdditionalTargetsPerSlot = 0
			} else {
				d.BaseTargets = 99
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

// Each printed measurement is pinned by equality, so a positive but wrong value
// is rejected by Decode on its own. n selects the occurrence in the fixture.
func TestSpells200PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_200.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"words", 1, 1}, {"question_limit", 1, 3}, {"repeat_cumulative_percent", 1, 25},
		{"outdoor_radius_miles", 1, 3}, {"underground_radius_feet", 1, 300}, {"choose_count", 1, 3},
		{"creature_minimum_cr", 1, 10}, {"reading_minutes_per_page", 1, 1}, {"cone_feet", 1, 60},
		{"radius_feet", 1, 10}, {"behavior_die_sides", 1, 10}, {"direction_die_sides", 1, 4},
		{"additional_radius_feet_per_slot", 1, 5}, {"advantage_range_feet", 1, 5}, {"move_feet", 1, 30},
		{"trigger_range_feet", 1, 10}, {"trigger_range_feet", 2, 5}, {"saves_per_turn", 1, 1},
		{"radius_feet", 2, 10}, {"height_feet", 1, 40}, {"move_feet", 2, 30}, {"affected_per_turn", 1, 1},
		{"attack_range_feet", 1, 5}, {"teleport_feet", 1, 30}, {"emanation_feet", 1, 15},
		{"base_targets", 1, 1}, {"additional_targets_per_slot", 1, 1}, {"roll_max", 1, 1},
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
