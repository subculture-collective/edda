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

func loadItems229(t *testing.T) ItemSourceBatch[ItemDetails229] {
	t.Helper()
	b, err := LoadItems229()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 216-218, in planning inventory order for issue #229.
func TestItems229SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		options    []string
		attunement string
		rarity     string
	}{
		{"Crystal Ball of Telepathy", "crystal-ball-of-telepathy", 216, "wondrous-item", nil, "required", "legendary"},
		{"Crystal Ball of True Seeing", "crystal-ball-of-true-seeing", 216, "wondrous-item", nil, "required", "legendary"},
		{"Cube of Force", "cube-of-force", 216, "wondrous-item", nil, "required", "rare"},
		{"Cubic Gate", "cubic-gate", 216, "wondrous-item", nil, "none", "legendary"},
		{"Dagger of Venom", "dagger-of-venom", 216, "weapon", []string{"dagger"}, "none", "rare"},
		{"Dancing Sword", "dancing-sword", 217, "weapon", []string{"greatsword", "longsword", "rapier", "scimitar", "shortsword"}, "required", "very-rare"},
		{"Decanter of Endless Water", "decanter-of-endless-water", 217, "wondrous-item", nil, "none", "uncommon"},
		{"Deck of Illusions", "deck-of-illusions", 217, "wondrous-item", nil, "none", "uncommon"},
		{"Defender", "defender", 218, "weapon", []string{"any-melee-weapon"}, "required", "legendary"},
		{"Demon Armor", "demon-armor", 218, "armor", []string{"any-light-armor", "any-medium-armor", "any-heavy-armor"}, "required", "very-rare"},
		{"Dimensional Shackles", "dimensional-shackles", 218, "wondrous-item", nil, "none", "rare"},
		{"Dragon Orb", "dragon-orb", 218, "wondrous-item", nil, "required", "artifact"},
	}
	// Only three families list charges; the Cube of Force and Dragon Orb cost
	// varies per spell, so their per-use charge is zero.
	charges := map[string]ItemCharges{
		"cube-of-force": {Mode: "limited", Maximum: 10, PerUse: 0, Recharge: "regain-1d6-at-dawn"},
		"cubic-gate":    {Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		"dragon-orb":    {Mode: "limited", Maximum: 7, PerUse: 0, Recharge: "regain-1d4-plus-3-at-dawn"},
	}
	b := loadItems229(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-05" || b.CatalogueIssue != 229 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #229")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 216–218") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			wantOptions := want.options
			if wantOptions == nil {
				wantOptions = []string{}
			}
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || !slices.Equal(f.Base.Options, wantOptions) || len(f.Base.Excluded) != 0 {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != 1 {
				t.Fatalf("the source lists one form of %s, got %d variants", want.name, len(f.Variants))
			}
			v := f.Variants[0]
			if v.ID != id+"/variant/standard" || v.Name != want.name || v.Rarity != want.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
				t.Fatalf("variant mismatch: %+v", v)
			}
			if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
				t.Fatal("no price is listed; cost must stay unspecified")
			}
			wantCharges, listed := charges[want.slug]
			if !listed {
				wantCharges = ItemCharges{Mode: "none", Recharge: "none"}
			}
			if *v.Charges != wantCharges {
				t.Fatalf("charges mismatch: %+v", *v.Charges)
			}
			if v.MechanicsAvailable() {
				t.Fatal("catalogue record became available")
			}
		})
	}
	// Join every family to the decision inventory the issue cites, in its order.
	data, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Family string `json:"family"`
			Page   int    `json:"source_page"`
			Batch  string `json:"batch_key"`
			Issue  int    `json:"catalogue_issue"`
			Status string `json:"status"`
		} `json:"records"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 229 {
			continue
		}
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Family != "magic-items" || r.Batch != b.Batch || r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 12 {
		t.Fatalf("inventory assigns %d families to issue #229", index)
	}
	b.Families[0].Variants[0].Details.Telepathy.Condition = "changed"
	if fresh := loadItems229(t); fresh.Families[0].Variants[0].Details.Telepathy.Condition != "while-touching" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all twelve variants, transcribed from pinned PDF
// pp. 216-218 (the Deck of Illusions table continues from p. 217 onto p. 218).
func TestItems229SourceValues(t *testing.T) {
	expected := []string{
		`{"telepathy":{"condition":"while-touching","spell":"scrying","save_dc":17,"communication_mode":"telepathy","communication_targets":"creatures-you-can-see-within-range-of-sensor","communication_range_feet":30,"communication_range_from":"spell-sensor","extra_spell":"suggestion","extra_save_dc":17,"extra_targets":"one-of-those-creatures","extra_delivery":"through-the-sensor","extra_concentration_required":false,"extra_ends_with":"scrying-ends","extra_uses":1,"extra_recharge":"next-dawn"}}`,
		`{"true_seeing":{"condition":"while-touching","spell":"scrying","save_dc":17,"sense":"truesight","sense_range_feet":120,"sense_centered_on":"spell-sensor"}}`,
		`{"cube_of_force":{"appearance":"cube","width_inches":1,"markings":"distinct-marking-on-each-face","activation":"press-face","cost":"expend-charges-required-for-face","save_dc":17,"starting_charges":10,"regain_dice":{"dice":1,"sides":6,"type":"charges"},"regain_timing":"dawn","faces":[{"spell":"mage-armor","charge_cost":1},{"spell":"shield","charge_cost":1},{"spell":"tiny-hut","charge_cost":3},{"spell":"private-sanctum","charge_cost":4},{"spell":"resilient-sphere","charge_cost":4},{"spell":"wall-of-force","charge_cost":5}]}}`,
		`{"cubic_gate":{"appearance":"cube","width_inches":3,"radiates":"palpable-magical-energy","sides":6,"material_plane_sides":1,"other_planes_determined_by":"gm","max_charges":3,"regain_dice":{"dice":1,"sides":3,"type":"charges"},"regain_timing":"dawn","action":"magic","charge_cost":1,"gate_spell":"gate","gate_activation":"press-one-side","gate_effect":"opens-portal-to-plane-keyed-to-side","plane_shift_spell":"plane-shift","plane_shift_activation":"press-one-side-twice","plane_shift_effect":"transports-targets-to-plane-keyed-to-side"}}`,
		`{"venom":{"attack_bonus":1,"damage_bonus":1,"apply_action":"bonus-action","coats":"blade","poison_duration_seconds":60,"poison_ends_on":["duration-elapsed","attack-hits-creature"],"save":"constitution","save_dc":15,"failure_damage":{"dice":2,"sides":10,"type":"poison"},"failure_condition":"poisoned","condition_duration_seconds":60,"uses":1,"recharge":"next-dawn"}}`,
		`{"dancing":{"activate_action":"bonus-action","activation":"toss-into-air","hover_fly_feet":30,"attack_range_feet":5,"attack_roll":"uses-your-attack-roll","adds_ability_modifier_to_damage":true,"move_action":"bonus-action","move_feet":30,"move_within_feet_of_you":30,"attack_with_move":true,"attacks_before_return":4,"returns_to":"your-hand","no_hand_free":"falls-to-ground-in-your-space","no_path":"moves-as-close-as-possible-then-falls","hover_ends_on":["you-grasp-it","more-than-30-feet-away"],"hover_end_distance_feet":30}}`,
		`{"decanter":{"appearance":"stoppered-flask","weight_pounds":2,"action":"magic","opens":"remove-stopper","water_types":["fresh","salt"],"water_choice":"your-choice","stops_at":"start-of-your-next-turn","commands":[{"word":"splash","gallons":1},{"word":"fountain","gallons":5},{"word":"geyser","gallons":30}],"geyser_shape":"line","geyser_length_feet":30,"geyser_width_feet":1,"geyser_aim":"one-direction-no-action-if-holding","geyser_targets":"one-creature-in-line","geyser_save":"strength","geyser_save_dc":13,"geyser_damage":{"dice":1,"sides":4,"type":"bludgeoning"},"geyser_failure_condition":"prone","object_condition":"not-worn-or-carried","object_max_pounds":200,"object_effect":"knocked-over"}}`,
		`{"illusions":{"full_deck_cards":34,"creature_cards":32,"mirror_cards":2,"missing_cards_dice":{"dice":1,"sides":20,"type":"cards"},"missing_cards_subtract":1,"functions_only_if":"cards-drawn-at-random","draw_action":"magic","throw_range_feet":30,"illusion":"looks-and-behaves-like-real-creature","can_do_harm":false,"move_action":"magic","move_while_within_feet":120,"move_requires":"can-see-illusion","move_distance_feet":30,"move_from":"its-card","physical_interaction":"reveals-false-objects-pass-through","inspect_action":"study","inspect_check":{"ability":"intelligence","skill":"investigation","dc":15},"inspect_result":"identifies-as-illusion","lasts_until":["card-moved","illusion-dispelled"],"dispelled_by":["dispel-magic","similar-effect"],"after_end":"card-image-disappears","card_reusable":false,"stat_block_source":"monsters","stat_block_exception":"card-drawer","table_die_sides":100,"table_page":218,"illusions":[{"min":1,"max":3,"creature":"adult-red-dragon"},{"min":4,"max":6,"creature":"archmage"},{"min":7,"max":9,"creature":"assassin"},{"min":10,"max":12,"creature":"bandit-captain"},{"min":13,"max":15,"creature":"basilisk"},{"min":16,"max":18,"creature":"berserker"},{"min":19,"max":21,"creature":"bugbear-warrior"},{"min":22,"max":24,"creature":"cloud-giant"},{"min":25,"max":27,"creature":"druid"},{"min":28,"max":30,"creature":"erinyes"},{"min":31,"max":33,"creature":"ettin"},{"min":34,"max":36,"creature":"fire-giant"},{"min":37,"max":39,"creature":"frost-giant"},{"min":40,"max":42,"creature":"gnoll-warrior"},{"min":43,"max":45,"creature":"goblin-warrior"},{"min":46,"max":48,"creature":"guardian-naga"},{"min":49,"max":51,"creature":"hill-giant"},{"min":52,"max":54,"creature":"hobgoblin-warrior"},{"min":55,"max":57,"creature":"incubus"},{"min":58,"max":60,"creature":"iron-golem"},{"min":61,"max":63,"creature":"knight"},{"min":64,"max":66,"creature":"kobold-warrior"},{"min":67,"max":69,"creature":"lich"},{"min":70,"max":72,"creature":"medusa"},{"min":73,"max":75,"creature":"night-hag"},{"min":76,"max":78,"creature":"ogre"},{"min":79,"max":81,"creature":"oni"},{"min":82,"max":84,"creature":"priest"},{"min":85,"max":87,"creature":"succubus"},{"min":88,"max":90,"creature":"troll"},{"min":91,"max":93,"creature":"veteran-warrior"},{"min":94,"max":96,"creature":"wyvern"},{"min":97,"max":100,"creature":"card-drawer"}]}}`,
		`{"defender":{"attack_bonus":3,"damage_bonus":3,"transfer_timing":"first-attack-each-turn","transfer_from":["attack-rolls","damage-rolls"],"transfer_to":"armor-class","transfer_amount":"some-or-all","example_attack_damage_bonus":1,"example_armor_class_bonus":2,"adjustment_lasts":"until-start-of-your-next-turn","must_hold_for_armor_class_bonus":true}}`,
		`{"demon_armor":{"condition":"while-wearing","armor_class_bonus":1,"language":"abyssal","gauntlets":"clawed-gauntlets","unarmed_damage":{"dice":1,"sides":8,"type":"slashing"},"unarmed_replaces_damage_type":"bludgeoning","unarmed_attack_bonus":1,"unarmed_damage_bonus":1,"cursed":true,"curse_trigger":"don-armor","cannot_doff":true,"doff_allowed_if":["targeted-by-remove-curse","targeted-by-similar-magic"],"curse_disadvantage":["attack-rolls-against-demons","saving-throws-against-demon-spells-and-special-abilities"]}}`,
		`{"shackles":{"place_action":"utilize","place_target_condition":"incapacitated","fit_size_min":"small","fit_size_max":"large","prevents_movement":"any-extradimensional-movement","prevented_examples":["teleportation","travel-to-different-plane"],"does_not_prevent":"passing-through-interdimensional-portal","remove_action":"utilize","removable_by":["you","creature-you-designate-when-you-use-the-shackles"],"escape_interval_days":30,"escape_check":{"ability":"strength","skill":"athletics","dc":30},"escape_success":"breaks-free-and-destroys-shackles"}}`,
		`{"dragon_orb":{"appearance":"etched-crystal-globe","diameter_inches":10,"used_diameter_inches":20,"interior":"swirling-mist","peer_action":"magic","peer_requires":"attuned","save":"charisma","save_dc":15,"success_effect":"control-orb-while-attuned","failure_condition":"charmed","failure_lasts":"while-attuned","charmed_cannot_end_attunement":true,"charmed_spell":"suggestion","charmed_spell_frequency":"at-will","charmed_spell_save_dc":18,"charmed_urge":"work-toward-evil-ends-orb-desires","desire_examples":["annihilation-of-a-particular-society-or-organization","freedom-from-the-orb","spread-suffering-in-the-world","advance-worship-of-tiamat","something-else-the-gm-decides"],"max_charges":7,"regain_dice":{"dice":1,"sides":4,"type":"charges"},"regain_bonus":3,"regain_timing":"dawn","spells_require":"control-orb","spells":[{"spell":"cure-wounds","spell_level":9,"save_dc":0,"charge_cost":4},{"spell":"daylight","spell_level":0,"save_dc":0,"charge_cost":1},{"spell":"death-ward","spell_level":0,"save_dc":0,"charge_cost":2},{"spell":"detect-magic","spell_level":0,"save_dc":0,"charge_cost":0},{"spell":"scrying","spell_level":0,"save_dc":18,"charge_cost":3}],"call_action":"magic","call_requires":"control-orb","call_range_miles":40,"call_directions":"all-directions","call_targets":"chromatic-dragons-in-range","call_effect":"compelled-to-come-to-orb-by-most-direct-route","call_unaffected":"dragon-deities-such-as-tiamat","call_attitude":"might-be-hostile-toward-you","call_cooldown_seconds":3600,"armor_class":20,"destroyed_by":["damage-from-plus-3-weapon","disintegrate-spell"],"nothing_else_harms":true}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"touch:cast-spell:touching-crystal-orb:item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 communication.telepathy-within-30-feet-of-sensor spell.cast-suggestion spell.suggestion-save-dc-17 concentration.suggestion-ends-with-scrying item.once-until-next-dawn",
		"touch:cast-spell:touching-crystal-orb:item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 senses.truesight-120-feet-at-sensor",
		"none:press-face:expend-charges-required-for-face:item.attunement item.charges item.press-face spell.cast-face-spell spell.face-save-dc-17 item.face-charge-cost-table",
		"none:magic-action:expend-one-charge:item.charges action.magic item.six-sides-keyed-to-planes spell.cast-gate spell.cast-plane-shift",
		"wield:bonus-action:once-until-next-dawn:item.magic-weapon attack.weapon-bonus damage.weapon-bonus action.bonus-action effect.poison-coating-one-minute save.constitution-dc-15 damage.poison-2d10 condition.poisoned-one-minute item.once-until-next-dawn",
		"wield:bonus-action:toss-into-air hover-flight:item.attunement item.magic-weapon action.bonus-action effect.hover-fly-30-feet attack.uses-your-attack-roll damage.adds-ability-modifier attack.returns-after-fourth-attack effect.hover-ends-when-grasped-or-over-30-feet",
		"none:magic-action:remove-stopper command-word:action.magic item.command-words water.fresh-or-salt water.splash-1-gallon water.fountain-5-gallons water.geyser-30-gallons area.line-30-feet-by-1-foot save.strength-dc-13 damage.bludgeoning-1d4 condition.prone",
		"none:magic-action:draw-card-at-random throw-within-30-feet:action.magic item.deck-34-cards item.cards-drawn-at-random summon.illusory-creature-from-table action.study check.intelligence-investigation-dc-15 effect.illusion-ends-when-card-moved-or-dispelled item.illusion-table",
		"wield:weapon-attack:first-attack-each-turn:item.attunement item.magic-weapon attack.weapon-bonus-3 damage.weapon-bonus-3 defense.transfer-bonus-to-armor-class effect.adjustment-lasts-until-next-turn",
		"wear:none::item.attunement item.wear defense.armor-class-bonus-1 language.abyssal attack.unarmed-strike-slashing-1d8 attack.unarmed-bonus-1 item.curse item.curse-cannot-doff attack.disadvantage-against-demons save.disadvantage-against-demons",
		"none:utilize-action:target-incapacitated-small-to-large:action.utilize condition.incapacitated movement.prevents-extradimensional-movement check.strength-athletics-dc-30 effect.escape-once-every-30-days",
		"none:magic-action:attuned control-orb-for-spells-and-call:item.attunement item.charges action.magic save.charisma-dc-15 condition.charmed spell.suggestion-save-dc-18 item.spell-charge-cost-table action.call-dragons defense.armor-class-20 item.artifact-destruction",
	}
	index := 0
	for _, f := range loadItems229(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails229
				d := json.NewDecoder(strings.NewReader(expected[index]))
				d.DisallowUnknownFields()
				if err := d.Decode(&want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(v.Details, &want) {
					t.Fatalf("PDF p.%d values mismatch: %+v", f.Presentations[0].Pages[0], v.Details)
				}
				got := fmt.Sprintf("%s:%s:%s:%s", v.Equip.Mode, v.Use.Mode, strings.Join(v.Use.Restrictions, " "), strings.Join(v.RequiredMechanicKeys, " "))
				if got != requirements[index] || len(v.Equip.Restrictions) != 0 {
					t.Fatalf("source requirements mismatch: %s", got)
				}
			})
			index++
		}
	}
	if index != len(expected) || index != len(requirements) {
		t.Fatal("fixture count disagrees with the batch")
	}
}

func TestItems229CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems229(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	first, err := LoadItems225()
	if err != nil {
		t.Fatal(err)
	}
	second, err := LoadItems226()
	if err != nil {
		t.Fatal(err)
	}
	third, err := LoadItems227()
	if err != nil {
		t.Fatal(err)
	}
	fourth, err := LoadItems228()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), first.SourceIdentities(), second.SourceIdentities(), third.SourceIdentities(), fourth.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 61 || variants != 74 || composed[0].Name != "Potion of Healing" || composed[37].Name != "Censer of Controlling Air Elementals" || composed[49].Name != "Crystal Ball of Telepathy" || composed[60].Name != "Dragon Orb" {
		t.Fatal("expected Potion of Healing, every earlier family and exactly the twelve families of #229")
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_229") {
				t.Fatal("catalogue batch must not be a coverage fixture")
			}
		}
	}
	if report := register.Coverage(); report.PlannedSpellItemScopes != 753 || report.FullSupport {
		t.Fatal("catalogue data changed the coverage boundary")
	}
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems229RejectMissingFields(t *testing.T) {
	var paths [][]any
	var walk func(node any, path []any)
	walk = func(node any, path []any) {
		switch n := node.(type) {
		case map[string]any:
			for key, child := range n {
				p := append(slices.Clone(path), key)
				paths = append(paths, p)
				walk(child, p)
			}
		case []any:
			for i, child := range n {
				walk(child, append(slices.Clone(path), i))
			}
		}
	}
	var root any
	if err := json.Unmarshal(items229JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 500 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items229JSON, &fresh); err != nil {
			t.Fatal(err)
		}
		node := fresh
		for _, step := range path[:len(path)-1] {
			switch s := step.(type) {
			case string:
				node = node.(map[string]any)[s]
			case int:
				node = node.([]any)[s]
			}
		}
		delete(node.(map[string]any), path[len(path)-1].(string))
		data, err := json.Marshal(fresh)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeItems229(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems229RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails229]
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails229] { return &b.Families[family].Variants[0] }
	// Family indexes follow inventory order; see TestItems229SourceIdentities:
	// 0 Telepathy, 1 True Seeing, 2 Cube of Force, 3 Cubic Gate, 4 Dagger of
	// Venom, 5 Dancing Sword, 6 Decanter, 7 Deck, 8 Defender, 9 Demon Armor,
	// 10 Shackles, 11 Dragon Orb.
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-04" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 228 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/dragon-scale-mail"
			b.Families[12].Name = "Dragon Scale Mail"
		}},
		{"previous batch family", func(b *batch) {
			b.Families[0].ID = "srd-5.2.1:content/items/crystal-ball-of-mind-reading"
			b.Families[0].Name = "Crystal Ball of Mind Reading"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/cape-of-the-mountebank" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Crystal Ball of Telepathic Scrying" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"weapon category", func(b *batch) { b.Families[4].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"dropped weapon base", func(b *batch) { b.Families[4].Base.Options = []string{} }},
		{"widened weapon base", func(b *batch) { b.Families[4].Base.Options = []string{"dagger", "shortsword"} }},
		{"dropped sword base", func(b *batch) { b.Families[5].Base.Options = b.Families[5].Base.Options[:4] }},
		{"narrowed melee base", func(b *batch) { b.Families[8].Base.Options = []string{"longsword"} }},
		{"narrowed armor base", func(b *batch) { b.Families[9].Base.Options = []string{"any-heavy-armor"} }},
		{"invented exclusion", func(b *batch) { b.Families[9].Base.Excluded = []string{"hide-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{215, 215} }},
		{"page range", func(b *batch) { b.Families[7].Presentations[0].Pages = []int{217, 218} }},
		{"deck page", func(b *batch) { b.Families[7].Presentations[0].Pages = []int{218, 218} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Crystal Ball" }},
		{"rarity", func(b *batch) { v(b, 11).Rarity = "legendary" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, 0).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, 3).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, 0).Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { v(b, 2).Charges = nil }},
		{"dropped cube charges", func(b *batch) { *v(b, 2).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"cube maximum", func(b *batch) { v(b, 2).Charges.Maximum = 3 }},
		{"cube per use", func(b *batch) { v(b, 2).Charges.PerUse = 1 }},
		{"cube recharge", func(b *batch) { v(b, 2).Charges.Recharge = "regain-1d3-at-dawn" }},
		{"gate maximum", func(b *batch) { v(b, 3).Charges.Maximum = 10 }},
		{"gate per use", func(b *batch) { v(b, 3).Charges.PerUse = 0 }},
		{"gate recharge", func(b *batch) { v(b, 3).Charges.Recharge = "none" }},
		{"orb maximum", func(b *batch) { v(b, 11).Charges.Maximum = 10 }},
		{"orb per use", func(b *batch) { v(b, 11).Charges.PerUse = 1 }},
		{"orb recharge", func(b *batch) { v(b, 11).Charges.Recharge = "regain-1d4-at-dawn" }},
		{"invented charges", func(b *batch) { v(b, 4).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, 0).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 4).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, 5).Use.Restrictions = []string{"toss-into-air"} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails229{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.TrueSeeing = v(b, 1).Details.TrueSeeing }},
		{"blank key", func(b *batch) { v(b, 0).Details.Telepathy.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0).Details.Telepathy.Condition = "While touching this crystal orb" }},
		{"telepathy dc", func(b *batch) { v(b, 0).Details.Telepathy.SaveDC = 15 }},
		{"telepathy suggestion dc", func(b *batch) { v(b, 0).Details.Telepathy.ExtraSaveDC = 0 }},
		{"telepathy range", func(b *batch) { v(b, 0).Details.Telepathy.CommunicationRangeFeet = 0 }},
		{"telepathy uses", func(b *batch) { v(b, 0).Details.Telepathy.ExtraUses = 0 }},
		{"true seeing range", func(b *batch) { v(b, 1).Details.TrueSeeing.SenseRangeFeet = 0 }},
		{"true seeing dc", func(b *batch) { v(b, 1).Details.TrueSeeing.SaveDC = 18 }},
		{"cube dc", func(b *batch) { v(b, 2).Details.CubeForce.SaveDC = 15 }},
		{"cube starting charges", func(b *batch) { v(b, 2).Details.CubeForce.StartingCharges = 3 }},
		{"cube regain dice", func(b *batch) { v(b, 2).Details.CubeForce.RegainDice.Sides = 3 }},
		{"cube missing face", func(b *batch) { v(b, 2).Details.CubeForce.Faces = v(b, 2).Details.CubeForce.Faces[:5] }},
		{"cube invented face", func(b *batch) {
			f := v(b, 2).Details.CubeForce
			f.Faces = append(f.Faces, ItemCubeFace{Spell: "wall-of-fire", ChargeCost: 5})
		}},
		{"cube reordered faces", func(b *batch) {
			f := v(b, 2).Details.CubeForce.Faces
			f[2], f[3] = f[3], f[2]
		}},
		{"cube face cost", func(b *batch) { v(b, 2).Details.CubeForce.Faces[2].ChargeCost = 4 }},
		{"cube width", func(b *batch) { v(b, 2).Details.CubeForce.WidthInches = 0 }},
		{"gate max charges", func(b *batch) { v(b, 3).Details.CubicGate.MaxCharges = 6 }},
		{"gate cost", func(b *batch) { v(b, 3).Details.CubicGate.ChargeCost = 2 }},
		{"gate sides", func(b *batch) { v(b, 3).Details.CubicGate.Sides = 4 }},
		{"gate regain dice", func(b *batch) { v(b, 3).Details.CubicGate.RegainDice.Sides = 6 }},
		{"gate width", func(b *batch) { v(b, 3).Details.CubicGate.WidthInches = 0 }},
		{"venom dc", func(b *batch) { v(b, 4).Details.Venom.SaveDC = 13 }},
		{"venom damage", func(b *batch) { v(b, 4).Details.Venom.FailureDamage.Sides = 8 }},
		{"venom damage type", func(b *batch) { v(b, 4).Details.Venom.FailureDamage.Type = "necrotic" }},
		{"venom duplicate end", func(b *batch) { v(b, 4).Details.Venom.PoisonEndsOn[1] = "duration-elapsed" }},
		{"venom bonus", func(b *batch) { v(b, 4).Details.Venom.AttackBonus = 0 }},
		{"venom uses", func(b *batch) { v(b, 4).Details.Venom.Uses = 0 }},
		{"sword attacks before return", func(b *batch) { v(b, 5).Details.Dancing.AttacksBeforeReturn = 3 }},
		{"sword range", func(b *batch) { v(b, 5).Details.Dancing.HoverFlyFeet = 0 }},
		{"sword duplicate end", func(b *batch) { v(b, 5).Details.Dancing.HoverEndsOn[1] = "you-grasp-it" }},
		{"decanter gallons", func(b *batch) { v(b, 6).Details.Decanter.Commands[2].Gallons = 25 }},
		{"decanter missing command", func(b *batch) { v(b, 6).Details.Decanter.Commands = v(b, 6).Details.Decanter.Commands[:2] }},
		{"decanter dc", func(b *batch) { v(b, 6).Details.Decanter.GeyserSaveDC = 15 }},
		{"decanter damage", func(b *batch) { v(b, 6).Details.Decanter.GeyserDamage.Dice = 2 }},
		{"decanter water types", func(b *batch) { v(b, 6).Details.Decanter.WaterTypes = []string{"fresh"} }},
		{"decanter weight", func(b *batch) { v(b, 6).Details.Decanter.WeightPounds = 0 }},
		{"decanter object weight", func(b *batch) { v(b, 6).Details.Decanter.ObjectMaxPounds = 0 }},
		{"deck missing row", func(b *batch) { v(b, 7).Details.Illusions.Illusions = v(b, 7).Details.Illusions.Illusions[:32] }},
		{"deck gap", func(b *batch) { v(b, 7).Details.Illusions.Illusions[10].Min = 32 }},
		{"deck row range", func(b *batch) { v(b, 7).Details.Illusions.Illusions[32].Max = 99 }},
		{"deck duplicate creature", func(b *batch) { v(b, 7).Details.Illusions.Illusions[1].Creature = "adult-red-dragon" }},
		{"deck renamed creature", func(b *batch) { v(b, 7).Details.Illusions.Illusions[21].Creature = "kobold" }},
		{"deck reordered creatures", func(b *batch) {
			r := v(b, 7).Details.Illusions.Illusions
			r[0].Creature, r[1].Creature = r[1].Creature, r[0].Creature
		}},
		{"deck drawer row", func(b *batch) { v(b, 7).Details.Illusions.Illusions[32].Creature = "wyvern" }},
		{"deck die", func(b *batch) { v(b, 7).Details.Illusions.TableDieSides = 20 }},
		{"deck table page", func(b *batch) { v(b, 7).Details.Illusions.TablePage = 217 }},
		{"deck cards", func(b *batch) { v(b, 7).Details.Illusions.FullDeckCards = 33 }},
		{"deck mirror cards", func(b *batch) { v(b, 7).Details.Illusions.MirrorCards = 3 }},
		{"deck missing dice", func(b *batch) { v(b, 7).Details.Illusions.MissingCardsDice.Sides = 10 }},
		{"deck check dc", func(b *batch) { v(b, 7).Details.Illusions.InspectCheck.DC = 13 }},
		{"deck range", func(b *batch) { v(b, 7).Details.Illusions.MoveWhileWithinFeet = 0 }},
		{"deck duplicate end", func(b *batch) { v(b, 7).Details.Illusions.LastsUntil[1] = "card-moved" }},
		{"defender bonus", func(b *batch) { v(b, 8).Details.Defender.AttackBonus = 2 }},
		{"defender damage bonus", func(b *batch) { v(b, 8).Details.Defender.DamageBonus = 2 }},
		{"defender example", func(b *batch) { v(b, 8).Details.Defender.ExampleArmorClassBonus = 3 }},
		{"defender duplicate source", func(b *batch) { v(b, 8).Details.Defender.TransferFrom[1] = "attack-rolls" }},
		{"demon armor die", func(b *batch) { v(b, 9).Details.DemonArmor.UnarmedDamage.Sides = 6 }},
		{"demon armor damage type", func(b *batch) { v(b, 9).Details.DemonArmor.UnarmedDamage.Type = "bludgeoning" }},
		{"demon armor bonus", func(b *batch) { v(b, 9).Details.DemonArmor.ArmorClassBonus = 0 }},
		{"demon armor doff", func(b *batch) { v(b, 9).Details.DemonArmor.DoffAllowedIf = []string{} }},
		{"demon armor duplicate disadvantage", func(b *batch) { v(b, 9).Details.DemonArmor.CurseDisadvantage[1] = "attack-rolls-against-demons" }},
		{"shackles dc", func(b *batch) { v(b, 10).Details.Shackles.EscapeCheck.DC = 20 }},
		{"shackles interval", func(b *batch) { v(b, 10).Details.Shackles.EscapeIntervalDays = 0 }},
		{"shackles removers", func(b *batch) { v(b, 10).Details.Shackles.RemovableBy = nil }},
		{"shackles duplicate example", func(b *batch) { v(b, 10).Details.Shackles.PreventedExamples[1] = "teleportation" }},
		{"orb save dc", func(b *batch) { v(b, 11).Details.DragonOrb.SaveDC = 18 }},
		{"orb suggestion dc", func(b *batch) { v(b, 11).Details.DragonOrb.CharmedSpellSaveDC = 15 }},
		{"orb charges", func(b *batch) { v(b, 11).Details.DragonOrb.MaxCharges = 10 }},
		{"orb regain dice", func(b *batch) { v(b, 11).Details.DragonOrb.RegainDice.Sides = 6 }},
		{"orb regain bonus", func(b *batch) { v(b, 11).Details.DragonOrb.RegainBonus = 0 }},
		{"orb missing spell", func(b *batch) { v(b, 11).Details.DragonOrb.Spells = v(b, 11).Details.DragonOrb.Spells[:4] }},
		{"orb spell cost", func(b *batch) { v(b, 11).Details.DragonOrb.Spells[0].ChargeCost = 3 }},
		{"orb free spell cost", func(b *batch) { v(b, 11).Details.DragonOrb.Spells[3].ChargeCost = 1 }},
		{"orb cure wounds level", func(b *batch) { v(b, 11).Details.DragonOrb.Spells[0].SpellLevel = 0 }},
		{"orb scrying dc", func(b *batch) { v(b, 11).Details.DragonOrb.Spells[4].SaveDC = 0 }},
		{"orb reordered spells", func(b *batch) {
			s := v(b, 11).Details.DragonOrb.Spells
			s[1], s[2] = s[2], s[1]
		}},
		{"orb call range", func(b *batch) { v(b, 11).Details.DragonOrb.CallRangeMiles = 0 }},
		{"orb armor class", func(b *batch) { v(b, 11).Details.DragonOrb.ArmorClass = 0 }},
		{"orb diameter", func(b *batch) { v(b, 11).Details.DragonOrb.UsedDiameterInches = 0 }},
		{"orb duplicate destroyer", func(b *batch) { v(b, 11).Details.DragonOrb.DestroyedBy[1] = "damage-from-plus-3-weapon" }},
		{"orb duplicate desire", func(b *batch) {
			v(b, 11).Details.DragonOrb.DesireExamples[1] = "annihilation-of-a-particular-society-or-organization"
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems229(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems229(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems229RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/crystal-ball-of-telepathy/variant/standard"`)
	if !bytes.Contains(items229JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items229JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems229(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"telepathy": {`, `"telepathy_plus": {`},
		"unknown detail":     {`"communication_mode"`, `"communication_mode_all"`},
		"string for number":  {`"save_dc": 17`, `"save_dc": "17"`},
		"fractional number":  {`"save_dc": 17`, `"save_dc": 17.5`},
		"number for boolean": {`"extra_concentration_required": false`, `"extra_concentration_required": 0`},
		"null table row":     {`"min": 1,`, `"min": null,`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items229JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items229JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems229(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items229JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems229(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch, and it accepts none of theirs.
	if _, err := DecodeItemCatalogue(items229JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems225(items229JSON); err == nil {
		t.Fatal("#225 decoder accepted the #229 batch")
	}
	if _, err := DecodeItems228(items229JSON); err == nil {
		t.Fatal("#228 decoder accepted the #229 batch")
	}
	if _, err := DecodeItems229(items228JSON); err == nil {
		t.Fatal("#229 decoder accepted the #228 batch")
	}
}
