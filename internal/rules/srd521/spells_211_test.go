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

func loadSpells211(t *testing.T) SpellRecords {
	t.Helper()
	b, err := os.ReadFile("testdata/spells_211.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// Literal PDF pp.143-145 value fixtures include every choice, effect and scaling
// value, independently of decoder defaults and structural validation.
func TestSpells211SourceValues(t *testing.T) {
	expected := []string{
		`{"id":"srd-5.2.1:content/spells/instant-summons","name":"Instant Summons","source_id":"srd-5.2.1:content/spells/instant-summons","source_page":143,"scope":"complete-source-record","level":6,"school":"conjuration","lists":["wizard"],"casting_time":"1-minute","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"sapphire-and-eligible-object-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a sapphire worth 1,000+ GP","cost_cp":100000,"consumed":false}},"ritual":true,"concentration":false,"duration":"until-dispelled","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"mark","mechanic":"mark.instant-summons","timing":"until-dispelled","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"summon","mechanic":"summon.instant-summons-object","timing":"magic-action-speak-name-and-crush-sapphire","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["mark.instant-summons","summon.instant-summons-object","cast.one-minute","cast.ritual","components.material"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":60,"casting_triggers":[],"end_timing":"dispelled-or-sapphire-crushed","instant_summons":{"object_max_weight_pounds":10,"object_max_longest_dimension_feet":6,"object_mark":"invisible-mark-on-object","sapphire_inscription":"object-name-invisibly-inscribed","sapphire_reuse":"different-sapphire-each-casting","activation_action":"magic-action","activation_steps":["speak-object-name","crush-sapphire"],"result":"object-instantly-appears-in-casters-hand","result_distance":"regardless-of-physical-or-planar-distances","after_result":"spell-ends","held_by_another":"not-transported","held_knowledge":["who-that-creature-is","where-that-creature-is-currently-located"]}}}`,
		`{"id":"srd-5.2.1:content/spells/invisibility","name":"Invisibility","source_id":"srd-5.2.1:content/spells/invisibility","source_page":143,"scope":"complete-source-record","level":2,"school":"illusion","lists":["bard","sorcerer","warlock","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"an eyelash in gum arabic","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"slot-targets","above_level":2,"dice_count":0,"die_sides":0,"effect_ids":["invisible"]},"effects":[{"id":"invisible","mechanic":"condition.invisibility-invisible","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["condition.invisibility-invisible","ending.invisibility-attack-damage-or-spell","scaling.targets-per-slot","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-early-end-trigger","invisibility":{"target":"creature-touched","condition":"invisible","early_end_triggers":["target-makes-attack-roll","target-deals-damage","target-casts-spell"],"early_end_timing":"immediately-after","base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/irresistible-dance","name":"Irresistible Dance","source_id":"srd-5.2.1:content/spells/irresistible-dance","source_page":143,"scope":"complete-source-record","level":6,"school":"enchantment","lists":["bard","wizard"],"casting_time":"action","range":{"kind":"feet","feet":30},"target":{"shape":"single","count":1,"eligibility":"visible-creature-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":60,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"success-dance","mechanic":"movement.irresistible-dance-on-success","timing":"on-successful-save-until-end-of-targets-next-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"charm","mechanic":"condition.irresistible-dance-charmed","timing":"on-failed-wisdom-save-for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.irresistible-dance-on-success","condition.irresistible-dance-charmed","save.irresistible-dance-action-repeat","concentration.source-duration"],"slot_levels":[6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-concentration-or-successful-save","irresistible_dance":{"save":"wisdom","success_effect":"dances-comically-until-end-of-its-next-turn","success_movement":"must-spend-all-movement-to-dance-in-place","failure_condition":"charmed","charmed_effects":["dances-comically","must-use-all-movement-to-dance-in-place","disadvantage-on-dexterity-saving-throws","disadvantage-on-attack-rolls","other-creatures-have-advantage-on-attack-rolls-against-it"],"repeat_action":"action-on-each-of-its-turns","repeat_save":"collect-itself-and-repeat-save","repeat_success":"spell-ends-on-target"}}}`,
		`{"id":"srd-5.2.1:content/spells/jump","name":"Jump","source_id":"srd-5.2.1:content/spells/jump","source_page":143,"scope":"complete-source-record","level":1,"school":"transmutation","lists":["druid","ranger","sorcerer","wizard"],"casting_time":"bonus-action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"willing-creature-touch"},"components":{"verbal":true,"somatic":true,"material":{"description":"a grasshopper’s hind leg","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":60,"scaling":{"kind":"slot-targets","above_level":1,"dice_count":0,"die_sides":0,"effect_ids":["jump"]},"effects":[{"id":"jump","mechanic":"movement.jump","timing":"once-on-each-of-targets-turns","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.jump","scaling.targets-per-slot","components.material"],"slot_levels":[1,2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended","jump":{"target":"willing-creature-touched","jump_max_feet":30,"movement_cost_feet":10,"uses_per_turn":1,"base_targets":1,"additional_targets_per_slot":1}}}`,
		`{"id":"srd-5.2.1:content/spells/knock","name":"Knock","source_id":"srd-5.2.1:content/spells/knock","source_page":143,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["bard","sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-access-prevented-object-in-range"},"components":{"verbal":true,"somatic":false,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"open","mechanic":"unlock.knock","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"noise","mechanic":"sound.knock","timing":"when-spell-cast","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["unlock.knock","sound.knock","suppression.knock-arcane-lock"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","knock":{"source_pages":[143,144],"object_examples":["door","box","chest","set-of-manacles","padlock","another-object-with-mundane-or-magical-means-preventing-access"],"mundane_states":["held-shut-by-mundane-lock","stuck","barred"],"mundane_results":["unlocked","unstuck","unbarred"],"locks_unlocked_if_multiple":1,"arcane_lock_suppression_minutes":10,"arcane_lock_result":"target-can-be-opened-and-closed","sound":"loud-knock-emanates-from-target","sound_audible_feet":300}}}`,
		`{"id":"srd-5.2.1:content/spells/legend-lore","name":"Legend Lore","source_id":"srd-5.2.1:content/spells/legend-lore","source_page":144,"scope":"complete-source-record","level":5,"school":"divination","lists":["bard","cleric","wizard"],"casting_time":"10-minutes","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"incense worth 250+ GP, which the spell consumes, and four ivory strips worth 50+ GP each","cost_cp":45000,"consumed":true}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"lore","mechanic":"knowledge.legend-lore","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["knowledge.legend-lore","response.gm-described-lore","cast.ten-minutes","components.material"],"slot_levels":[5,6,7,8,9],"source_details":{"casting_seconds":600,"casting_triggers":[],"end_timing":"instantaneous","legend_lore":{"incense_min_cost_cp":25000,"incense_consumed":true,"ivory_strip_count":4,"ivory_strip_min_cost_cp":5000,"ivory_strips_consumed":false,"subject_input":"name-or-describe","subjects":["famous-person","famous-place","famous-object"],"result":"brief-summary-of-significant-lore","described_by":"gm","possible_content":["important-details","amusing-revelations","secret-lore-never-widely-known"],"detail_scaling":"more-prior-knowledge-gives-more-precise-and-detailed-information","accuracy":"accurate-but-might-be-couched-in-figurative-language-or-poetry","accuracy_determined_by":"gm","not_famous":"sad-musical-trombone-notes-and-spell-fails"}}}`,
		`{"id":"srd-5.2.1:content/spells/lesser-restoration","name":"Lesser Restoration","source_id":"srd-5.2.1:content/spells/lesser-restoration","source_page":144,"scope":"complete-source-record","level":2,"school":"abjuration","lists":["bard","cleric","druid","paladin","ranger"],"casting_time":"bonus-action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"creature-touch"},"components":{"verbal":true,"somatic":true,"material":null},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"remove","mechanic":"remove.lesser-restoration","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["remove.lesser-restoration"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","lesser_restoration":{"target":"creature-touched","conditions_ended_count":1,"condition_choices":["blinded","deafened","paralyzed","poisoned"]}}}`,
		`{"id":"srd-5.2.1:content/spells/levitate","name":"Levitate","source_id":"srd-5.2.1:content/spells/levitate","source_page":144,"scope":"complete-source-record","level":2,"school":"transmutation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"feet","feet":60},"target":{"shape":"single","count":1,"eligibility":"visible-creature-or-loose-object-in-range"},"components":{"verbal":true,"somatic":true,"material":{"description":"a metal spring","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"rise","mechanic":"movement.levitate-rise","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"altitude","mechanic":"control.levitate-altitude","timing":"on-casters-turn","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false},{"id":"descent","mechanic":"movement.levitate-gentle-descent","timing":"when-spell-ends","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["movement.levitate-rise","control.levitate-altitude","movement.levitate-gentle-descent","save.constitution-negates","concentration.source-duration","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","levitate":{"rise_max_feet":20,"suspension":"remains-suspended-for-duration","object_max_weight_pounds":500,"unwilling_save":"constitution","unwilling_success":"unaffected","movement":"only-by-pushing-or-pulling-against-fixed-object-or-surface-within-reach","movement_surface_examples":["wall","ceiling"],"movement_mode":"as-if-climbing","altitude_change_max_feet":20,"altitude_directions":["up","down"],"self_altitude":"move-up-or-down-as-part-of-move","other_altitude_action":"magic-action","other_altitude_limit":"target-must-remain-within-spell-range","ending":"floats-gently-to-ground-if-still-aloft"}}}`,
		`{"id":"srd-5.2.1:content/spells/light","name":"Light","source_id":"srd-5.2.1:content/spells/light","source_page":144,"scope":"complete-source-record","level":0,"school":"evocation","lists":["bard","cleric","sorcerer","wizard"],"casting_time":"action","range":{"kind":"touch","feet":0},"target":{"shape":"single","count":1,"eligibility":"unattended-large-or-smaller-object-touch"},"components":{"verbal":true,"somatic":false,"material":{"description":"a firefly or phosphorescent moss","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"light","mechanic":"light.light-cantrip","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["light.light-cantrip","ending.light-recast","components.material"],"slot_levels":[],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-ended-or-cast-again","light":{"object_max_size":"large","object_requirement":"not-worn-or-carried-by-someone-else","bright_light_radius_feet":20,"dim_light_additional_feet":20,"color":"colored-as-caster-likes","blocking":"covering-with-something-opaque-blocks-light","recast":"spell-ends-if-cast-again"}}}`,
		`{"id":"srd-5.2.1:content/spells/lightning-bolt","name":"Lightning Bolt","source_id":"srd-5.2.1:content/spells/lightning-bolt","source_page":144,"scope":"complete-source-record","level":3,"school":"evocation","lists":["sorcerer","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"line","count":0,"eligibility":"creatures-in-line"},"components":{"verbal":true,"somatic":true,"material":{"description":"a bit of fur and a crystal rod","cost_cp":0,"consumed":false}},"ritual":false,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"slot","above_level":3,"dice_count":1,"die_sides":6,"effect_ids":["damage"]},"effects":[{"id":"damage","mechanic":"damage.lightning-bolt","timing":"immediate","dice_count":8,"die_sides":6,"damage_type":"lightning","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["damage.lightning-bolt","save.dexterity-half","scaling.per-slot","components.material"],"slot_levels":[3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","lightning_bolt":{"line_length_feet":100,"line_width_feet":5,"direction":"chosen-by-caster","save":"dexterity","damage_dice_count":8,"damage_die_sides":6,"damage_type":"lightning","save_success":"half-damage","slot_dice_count":1,"slot_die_sides":6}}}`,
		`{"id":"srd-5.2.1:content/spells/locate-animals-or-plants","name":"Locate Animals or Plants","source_id":"srd-5.2.1:content/spells/locate-animals-or-plants","source_page":144,"scope":"complete-source-record","level":2,"school":"divination","lists":["bard","druid","ranger"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"fur from a bloodhound","cost_cp":0,"consumed":false}},"ritual":true,"concentration":false,"duration":"instantaneous","duration_seconds":0,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"locate","mechanic":"detect.locate-animals-or-plants","timing":"immediate","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["detect.locate-animals-or-plants","cast.ritual","components.material"],"slot_levels":[2,3,4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"instantaneous","locate_animals_or_plants":{"query":"describe-or-name-specific-kind","kind_categories":["beast","plant-creature","nonmagical-plant"],"learns":["direction","distance"],"result":"closest-creature-or-plant-of-that-kind","search_radius_miles":5,"requirement":"only-if-any-are-present"}}}`,
		`{"id":"srd-5.2.1:content/spells/locate-creature","name":"Locate Creature","source_id":"srd-5.2.1:content/spells/locate-creature","source_page":145,"scope":"complete-source-record","level":4,"school":"divination","lists":["bard","cleric","druid","paladin","ranger","wizard"],"casting_time":"action","range":{"kind":"self","feet":0},"target":{"shape":"single","count":1,"eligibility":"self"},"components":{"verbal":true,"somatic":true,"material":{"description":"fur from a bloodhound","cost_cp":0,"consumed":false}},"ritual":false,"concentration":true,"duration":"timed","duration_seconds":3600,"scaling":{"kind":"none","above_level":0,"dice_count":0,"die_sides":0,"effect_ids":[]},"effects":[{"id":"locate","mechanic":"detect.locate-creature","timing":"for-duration","dice_count":0,"die_sides":0,"damage_type":"","bonus":0,"add_casting_modifier":false,"half_initial_only":false}],"required_mechanics":["detect.locate-creature","blocker.locate-creature-form-or-lead","concentration.source-duration","components.material"],"slot_levels":[4,5,6,7,8,9],"source_details":{"casting_seconds":0,"casting_triggers":[],"end_timing":"duration-or-concentration-ended","locate_creature":{"query":"describe-or-name-familiar-creature","senses":"direction-to-creatures-location","sense_radius_feet":1000,"moving_target":"know-direction-of-movement","targets":["specific-creature-known-to-caster","nearest-creature-of-specific-kind"],"kind_requirement":"seen-such-creature-up-close","kind_examples":["human","unicorn"],"close_distance_feet":30,"close_sightings_minimum":1,"different_form":"not-located","different_form_examples":["flesh-to-stone","polymorph"],"lead_blocker":"any-thickness-of-lead-blocking-direct-path"}}}`,
	}
	batch := loadSpells211(t)
	if len(batch.Records) != 12 {
		t.Fatal("exact twelve families required")
	}
	names := []string{"Instant Summons", "Invisibility", "Irresistible Dance", "Jump", "Knock", "Legend Lore", "Lesser Restoration", "Levitate", "Light", "Lightning Bolt", "Locate Animals or Plants", "Locate Creature"}
	pages := []int{143, 143, 143, 143, 143, 144, 144, 144, 144, 144, 144, 145}
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

func TestSpells211CompositionBoundaries(t *testing.T) {
	legacy := loadSpellRecords(t)
	batch := loadSpells211(t)
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
func TestSpells211MalformedFields(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_211.json")
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

func TestSpells211JSONGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_211.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{`"object_max_weight_pounds": 10`, `"object_max_weight_pounds": 10, "object_max_weight_pounds": 11`},
		{`"object_max_weight_pounds": 10`, `"object_max_weight_pounds": 10, "OBJECT_MAX_WEIGHT_POUNDS": 11`},
		{`"object_max_weight_pounds": 10`, `"object_max_weight_pounds": 10, "unknown": true`},
		{`"source_id": "srd-5.2.1:content/spells/instant-summons"`, `"source_id": "srd-5.2.1:content/spells/unknown"`},
		{`"effect_ids": [
          "invisible"
        ]`, `"effect_ids": ["unknown"]`},
		{`"slot_levels": [],`, ``},
		{`"slot_levels": [],`, `"slot_levels": null,`},
		{`"slot_levels": [],`, `"slot_levels": [0],`},
		{`"casting_time": "1-minute"`, `"casting_time": "1-hour"`},
		{`"casting_seconds": 60`, `"casting_seconds": 3600`},
		{`"feet": 30`, `"feet": 31`},
		{`"cost_cp": 100000`, `"cost_cp": 100001`},
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

func TestSpells211InvalidTypedVariants(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  int
		change func(*SpellRecord)
	}{
		{"summons weight", 0, func(s *SpellRecord) { s.SourceDetails.InstantSummons.ObjectMaxWeightPounds = 20 }},
		{"summons ritual", 0, func(s *SpellRecord) { *s.Ritual = false }},
		{"invisibility triggers", 1, func(s *SpellRecord) {
			s.SourceDetails.Invisibility.EarlyEndTriggers = s.SourceDetails.Invisibility.EarlyEndTriggers[:2]
		}},
		{"invisibility duration", 1, func(s *SpellRecord) { s.DurationSeconds = 60 }},
		{"dance penalty", 2, func(s *SpellRecord) {
			s.SourceDetails.IrresistibleDance.CharmedEffects[2] = "disadvantage-on-wisdom-saving-throws"
		}},
		{"dance somatic", 2, func(s *SpellRecord) { *s.Components.Somatic = true }},
		{"jump distance", 3, func(s *SpellRecord) { s.SourceDetails.Jump.JumpMaxFeet = 40 }},
		{"jump concentration", 3, func(s *SpellRecord) { *s.Concentration = true }},
		{"knock pages", 4, func(s *SpellRecord) { s.SourceDetails.Knock.SourcePages = []int{143} }},
		{"knock suppression", 4, func(s *SpellRecord) { s.SourceDetails.Knock.ArcaneLockSuppressionMinutes = 60 }},
		{"lore total cost", 5, func(s *SpellRecord) { *s.Components.Material.CostCP = 25000 }},
		{"lore strips", 5, func(s *SpellRecord) { s.SourceDetails.LegendLore.IvoryStripsConsumed = true }},
		{"restoration choices", 6, func(s *SpellRecord) { s.SourceDetails.LesserRestoration.ConditionChoices[2] = "petrified" }},
		{"levitate weight", 7, func(s *SpellRecord) { s.SourceDetails.Levitate.ObjectMaxWeightPounds = 1000 }},
		{"levitate save", 7, func(s *SpellRecord) { s.SourceDetails.Levitate.UnwillingSave = "strength" }},
		{"light slots", 8, func(s *SpellRecord) { s.SlotLevels = []int{1} }},
		{"light radius", 8, func(s *SpellRecord) { s.SourceDetails.Light.BrightLightRadiusFeet = 30 }},
		{"bolt width", 9, func(s *SpellRecord) { s.SourceDetails.LightningBolt.LineWidthFeet = 10 }},
		{"bolt shape", 9, func(s *SpellRecord) { s.Target.Shape = "cone" }},
		{"plants radius", 10, func(s *SpellRecord) { s.SourceDetails.LocateAnimalsOrPlants.SearchRadiusMiles = 3 }},
		{"plants ritual", 10, func(s *SpellRecord) { *s.Ritual = false }},
		{"creature blocker", 11, func(s *SpellRecord) { s.SourceDetails.LocateCreature.LeadBlocker = "none" }},
		{"creature sightings", 11, func(s *SpellRecord) { s.SourceDetails.LocateCreature.CloseSightingsMinimum = 2 }},
		{"required mechanic", 11, func(s *SpellRecord) { s.RequiredMechanics[0] = "unknown" }},
		{"effect list", 8, func(s *SpellRecord) { s.Effects = append(s.Effects, s.Effects[0]) }},
		{"effect timing", 9, func(s *SpellRecord) { s.Effects[0].Timing = "never" }},
		{"multiple payloads", 0, func(s *SpellRecord) { s.SourceDetails.LocateCreature = &SpellLocateCreature{} }},
		{"earlier batch payload", 0, func(s *SpellRecord) { s.SourceDetails.Blink = &SpellBlink{} }},
		{"missing details", 0, func(s *SpellRecord) { s.SourceDetails = nil }},
		{"bounded scope", 0, func(s *SpellRecord) { s.Scope = "bounded-self-level-1" }},
		{"source link", 0, func(s *SpellRecord) { s.SourceID = "unknown" }},
		{"effect graph", 0, func(s *SpellRecord) { s.Effects[0].Mechanic = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := loadSpells211(t)
			tc.change(&r.Records[tc.index])
			if r.Validate() == nil || r.Records[tc.index].MechanicsAvailable() {
				t.Fatal("invalid source accepted")
			}
		})
	}
}

func TestSpells211RequiredNumbersAndFlags(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_211.json")
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

func TestSpells211WhitespaceEmptyListsAndNestedGuards(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_211.json")
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
		`"object_max_weight_pounds": 10, "object_max_weight_pounds": 11`,
		`"object_max_weight_pounds": 10, "OBJECT_MAX_WEIGHT_POUNDS": 11`,
		`"object_max_weight_pounds": 10, "unknown": true`,
	} {
		b := bytes.Replace(original, []byte(`"object_max_weight_pounds": 10`), []byte(replacement), 1)
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
func TestSpells211PinnedPrintedValues(t *testing.T) {
	original, err := os.ReadFile("testdata/spells_211.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		n    int
		want int
	}{
		{"object_max_weight_pounds", 1, 10},
		{"object_max_weight_pounds", 2, 500},
		{"object_max_longest_dimension_feet", 1, 6},
		{"base_targets", 1, 1},
		{"base_targets", 2, 1},
		{"additional_targets_per_slot", 1, 1},
		{"additional_targets_per_slot", 2, 1},
		{"jump_max_feet", 1, 30},
		{"movement_cost_feet", 1, 10},
		{"uses_per_turn", 1, 1},
		{"locks_unlocked_if_multiple", 1, 1},
		{"arcane_lock_suppression_minutes", 1, 10},
		{"sound_audible_feet", 1, 300},
		{"incense_min_cost_cp", 1, 25000},
		{"ivory_strip_count", 1, 4},
		{"ivory_strip_min_cost_cp", 1, 5000},
		{"conditions_ended_count", 1, 1},
		{"rise_max_feet", 1, 20},
		{"altitude_change_max_feet", 1, 20},
		{"bright_light_radius_feet", 1, 20},
		{"dim_light_additional_feet", 1, 20},
		{"line_length_feet", 1, 100},
		{"line_width_feet", 1, 5},
		{"damage_dice_count", 1, 8},
		{"damage_die_sides", 1, 6},
		{"slot_dice_count", 1, 1},
		{"slot_die_sides", 1, 6},
		{"search_radius_miles", 1, 5},
		{"sense_radius_feet", 1, 1000},
		{"close_distance_feet", 1, 30},
		{"close_sightings_minimum", 1, 1},
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
// inventory assigns to catalogue issue #211.
func TestSpells211MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueItem == 211 {
			want[r.ID+"|"+r.Name] = r.SourcePage
		}
	}
	got := map[string]int{}
	for _, r := range loadSpells211(t).Records {
		got[r.ID+"|"+r.Name] = r.SourcePage
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from the decision inventory: got %v want %v", got, want)
	}
}

func TestSpells211ChangesNoteNamesOnlyPageBreakContinuations(t *testing.T) {
	b, err := os.ReadFile("testdata/spells_211.json")
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
	want := "SRD 5.2.1 PDF pp.143-145 adapted into twelve complete offline spell families for issue #211, including page-break continuation text for Knock (pp.143-144). No executable effects or availability grants."
	if root.Source.Changes != want {
		t.Fatalf("changes note = %q", root.Source.Changes)
	}
}

// A batch payload cannot ride on a record of an earlier batch.
func TestSpells211PayloadOnEarlierBatchRecord(t *testing.T) {
	r := loadSpells202(t)
	r.Records[0].SourceDetails.Cure = nil
	r.Records[0].SourceDetails.LocateCreature = loadSpells211(t).Records[11].SourceDetails.LocateCreature
	if r.Validate() == nil {
		t.Fatal("batch 211 payload accepted on a Cure Wounds record")
	}
}

// Each record round-trips exactly, so the pinned printed record and the fixture agree.
func TestSpells211RoundTripsPinnedRecords(t *testing.T) {
	for _, r := range loadSpells211(t).Records {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != spell211Pins[strings.TrimPrefix(r.ID, prefix+"spells/")] {
			t.Fatalf("%s differs from its pinned record", r.Name)
		}
	}
	if len(spell211Pins) != 12 {
		t.Fatal("twelve pinned records required")
	}
}
