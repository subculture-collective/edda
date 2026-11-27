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

func loadItems238(t *testing.T) ItemSourceBatch[ItemDetails238] {
	t.Helper()
	b, err := LoadItems238()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 236-238, in planning inventory order for issue #238.
// Potions of Healing, the twelfth source family of the issue, is held by the
// Healing contract fixture and is joined to the inventory from there.
func TestItems238SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	noCharges := ItemCharges{Mode: "none", Recharge: "none"}
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		charges    ItemCharges
		variants   []variant
	}{
		{"Potion of Poison", "potion-of-poison", 236, "potion", nil, "none", noCharges, []variant{{"standard", "Potion of Poison", "uncommon"}}},
		{"Potion of Resistance", "potion-of-resistance", 236, "potion", nil, "none", noCharges, []variant{{"standard", "Potion of Resistance", "uncommon"}}},
		{"Potion of Speed", "potion-of-speed", 236, "potion", nil, "none", noCharges, []variant{{"standard", "Potion of Speed", "very-rare"}}},
		{"Potion of Vitality", "potion-of-vitality", 236, "potion", nil, "none", noCharges, []variant{{"standard", "Potion of Vitality", "very-rare"}}},
		{"Potion of Water Breathing", "potion-of-water-breathing", 236, "potion", nil, "none", noCharges, []variant{{"standard", "Potion of Water Breathing", "uncommon"}}},
		{"Quarterstaff of the Acrobat", "quarterstaff-of-the-acrobat", 237, "weapon", []string{"quarterstaff"}, "required", noCharges, []variant{{"standard", "Quarterstaff of the Acrobat", "very-rare"}}},
		{"Ring of Animal Influence", "ring-of-animal-influence", 237, "ring", nil, "none", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Ring of Animal Influence", "rare"}}},
		{"Ring of Djinni Summoning", "ring-of-djinni-summoning", 237, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Djinni Summoning", "legendary"}}},
		{"Ring of Elemental Command", "ring-of-elemental-command", 237, "ring", nil, "required", ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-plus-1-at-dawn"}, []variant{{"air", "Ring of Elemental Command (air)", "legendary"}, {"earth", "Ring of Elemental Command (earth)", "legendary"}, {"fire", "Ring of Elemental Command (fire)", "legendary"}, {"water", "Ring of Elemental Command (water)", "legendary"}}},
		{"Ring of Evasion", "ring-of-evasion", 238, "ring", nil, "required", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Ring of Evasion", "rare"}}},
		{"Ring of Feather Falling", "ring-of-feather-falling", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Feather Falling", "rare"}}},
	}
	b := loadItems238(t)
	if len(b.Families) != 11 || b.Batch != "data-magic-items-14" || b.CatalogueIssue != 238 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the eleven families of issue #238")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 236–238") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	total := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || len(f.Base.Excluded) != 0 || !slices.Equal(f.Base.Options, append([]string{}, want.base...)) {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != len(want.variants) {
				t.Fatalf("the source lists %d forms of %s, got %d variants", len(want.variants), want.name, len(f.Variants))
			}
			for j, wv := range want.variants {
				v := f.Variants[j]
				if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != wv.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				if *v.Charges != want.charges {
					t.Fatalf("charges mismatch: %+v", *v.Charges)
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
				total++
			}
		})
	}
	if total != 14 {
		t.Fatalf("expected 14 source-listed variants, got %d", total)
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
		Batches []struct {
			Key     string `json:"key"`
			Pages   string `json:"pages"`
			Members []struct {
				Name string `json:"name"`
				Page int    `json:"source_page"`
			} `json:"members"`
		} `json:"batches"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	healing := contract.Families[0]
	healingJoined := 0
	type member struct {
		name string
		page int
	}
	var members []member
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 238 {
			continue
		}
		if r.Family != "magic-items" || r.Batch != b.Batch {
			t.Fatalf("inventory record %s changed family or batch", r.ID)
		}
		members = append(members, member{r.Name, r.Page})
		if r.ID == "srd-5.2.1:content/magic-items/potions-of-healing" {
			// Held once, by the contract fixture; its four table variants stay there.
			if r.Name != "Potions of Healing" || healing.ID != "srd-5.2.1:content/items/potion-of-healing" || len(healing.Presentations) != 2 || healing.Presentations[1].ID != r.ID ||
				!slices.Equal(healing.Presentations[1].Pages, []int{r.Page, r.Page}) || len(healing.Variants) != 4 {
				t.Fatal("the contract fixture no longer holds Potions of Healing")
			}
			for i, want := range []struct {
				rarity  string
				healing ItemHealing
			}{
				{"common", ItemHealing{2, 4, 2}}, {"uncommon", ItemHealing{4, 4, 4}}, {"rare", ItemHealing{8, 4, 8}}, {"very-rare", ItemHealing{10, 4, 20}},
			} {
				if v := healing.Variants[i]; v.Rarity != want.rarity || *v.Healing != want.healing {
					t.Fatalf("healing variant disagrees with the p. 236 table: %+v", v)
				}
			}
			healingJoined++
			continue
		}
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 11 || len(members) != 12 {
		t.Fatalf("inventory assigns %d families to issue #238", len(members))
	}
	if healingJoined != 1 {
		t.Fatal("Potions of Healing must be joined exactly once")
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "236, 237, 238" || len(batch.Members) != 12 {
			t.Fatalf("planning batch pages or member count changed: %q, %d", batch.Pages, len(batch.Members))
		}
		for i, m := range batch.Members {
			if m.Name != members[i].name || m.Page != members[i].page {
				t.Fatalf("planning batch member %d is %s p.%d", i, m.Name, m.Page)
			}
		}
	}
	if planned != 1 {
		t.Fatalf("expected one planning batch %s, got %d", b.Batch, planned)
	}
	b.Families[0].Variants[0].Details.Poison.Disguise = "changed"
	if fresh := loadItems238(t); fresh.Families[0].Variants[0].Details.Poison.Disguise != "looks-smells-and-tastes-like-healing-or-other-beneficial-potion" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 14 variants, transcribed from pinned PDF pp.
// 236-238.
func TestItems238SourceValues(t *testing.T) {
	expected := []string{
		`{"poison":{"disguise":"looks-smells-and-tastes-like-healing-or-other-beneficial-potion","masked_by":"illusion-magic","revealed_by":"identify","damage":{"dice":4,"sides":6,"type":"poison"},"save_ability":"constitution","save_dc":13,"condition":"poisoned","duration_hours":1}}`,
		`{"resistance_potion":{"defense":"resistance","duration_hours":1,"type_chosen_by":"gm-chooses-or-rolls-randomly","table_die_sides":10,"damage_types":[{"roll":1,"damage_type":"acid"},{"roll":2,"damage_type":"cold"},{"roll":3,"damage_type":"fire"},{"roll":4,"damage_type":"force"},{"roll":5,"damage_type":"lightning"},{"roll":6,"damage_type":"necrotic"},{"roll":7,"damage_type":"poison"},{"roll":8,"damage_type":"psychic"},{"roll":9,"damage_type":"radiant"},{"roll":10,"damage_type":"thunder"}]}}`,
		`{"speed":{"spell":"haste","duration_minutes":1,"concentration_required":false,"lethargy_when_it_ends":false,"appearance":"yellow-fluid-streaked-with-black-that-swirls-on-its-own"}}`,
		`{"vitality":{"removes":"all-exhaustion-levels","ends_condition":"poisoned","hit_die_healing":"maximum-for-any-hit-point-die-spent","duration_hours":24,"appearance":"crimson-liquid-pulses-with-dull-light-like-a-heartbeat"}}`,
		`{"water_breathing":{"grants":"breathe-underwater","duration_hours":24,"appearance":"cloudy-green-fluid-smelling-of-the-sea-with-jellyfish-like-bubble","source_pages":[236,237]}}`,
		`{"acrobat":{"attack_bonus":2,"damage_bonus":2,"light":{"condition":"while-holding","color":"green","dim_radius_feet":10,"activation":["bonus-action","after-rolling-initiative"],"extinguish_action":"bonus"},"alter_form":{"action":"bonus","forms":["six-inch-rod","ten-foot-pole","quarterstaff"],"rod_inches":6,"pole_feet":10,"elongation_limit":"surrounding-space"},"acrobatic_assist":{"forms":["quarterstaff","ten-foot-pole"],"check_ability":"dexterity","check_skill":"acrobatics","check_mode":"advantage"},"attack_deflection":{"form":"quarterstaff","trigger":"hit-by-attack-while-holding","action":"reaction","armor_class_bonus":5,"possible_result":"attack-may-miss-you","reuse_after":"short-or-long-rest"},"ranged_weapon":{"form":"quarterstaff","property":"thrown","normal_range_feet":30,"long_range_feet":120,"returns_to_hand":"immediately-after-ranged-attack"}}}`,
		`{"animal_influence":{"max_charges":3,"recharge":{"dice":1,"sides":3,"type":"charges"},"recharge_when":"daily-at-dawn","condition":"while-wearing","charge_cost":1,"save_dc":13,"spells":["animal-friendship","fear","speak-with-animals"],"fear_affects_only":"beasts"}}`,
		`{"djinni_summoning":{"condition":"while-wearing","action":"magic","creature":"djinni","origin":"elemental-plane-of-air","summon_range_feet":120,"concentration_required":true,"max_duration_hours":1,"also_ends_at":"zero-hit-points","attitude":"friendly-to-you-and-allies","obeys_commands":true,"without_command":"defends-itself-only","cooldown_hours":24,"nonmagical_when":"djinni-dies","often_created_by":"djinn-as-gifts"}}`,
		`{"elemental_command":{"linked_plane":"air","plane_chosen_by":"gm-chooses-or-rolls-randomly","bane":{"target_creature_type":"elemental","your_attacks":"advantage","attacks_against_you":"disadvantage"},"compulsion":{"action":"magic","range_feet":60,"save_ability":"wisdom","save_dc":18,"condition":"charmed","condition_ends":"start-of-your-next-turn","control":"you-determine-its-move-and-action-on-its-next-turn"},"focus":{"language":"auran","defense":"resistance","damage_type":"lightning","speed_kind":"fly","speed":"equal-to-speed","other_benefits":["can-hover"]},"spell_save_dc":18,"spells":[{"spell":"chain-lightning","charge_cost":3},{"spell":"feather-fall","charge_cost":0},{"spell":"gust-of-wind","charge_cost":2},{"spell":"wind-wall","charge_cost":1}],"max_charges":5,"recharge":{"dice":1,"sides":4,"type":"charges"},"recharge_bonus":1,"recharge_when":"daily-at-dawn","source_pages":[237,238]}}`,
		`{"elemental_command":{"linked_plane":"earth","plane_chosen_by":"gm-chooses-or-rolls-randomly","bane":{"target_creature_type":"elemental","your_attacks":"advantage","attacks_against_you":"disadvantage"},"compulsion":{"action":"magic","range_feet":60,"save_ability":"wisdom","save_dc":18,"condition":"charmed","condition_ends":"start-of-your-next-turn","control":"you-determine-its-move-and-action-on-its-next-turn"},"focus":{"language":"terran","defense":"resistance","damage_type":"acid","speed_kind":"none","speed":"none","other_benefits":["rubble-rocks-or-dirt-is-not-difficult-terrain","move-through-solid-earth-or-rock-as-difficult-terrain","passage-does-not-disturb-matter","shunted-out-if-turn-ends-in-solid-earth-or-rock"]},"spell_save_dc":18,"spells":[{"spell":"earthquake","charge_cost":5},{"spell":"stone-shape","charge_cost":2},{"spell":"stoneskin","charge_cost":3},{"spell":"wall-of-stone","charge_cost":3}],"max_charges":5,"recharge":{"dice":1,"sides":4,"type":"charges"},"recharge_bonus":1,"recharge_when":"daily-at-dawn","source_pages":[237,238]}}`,
		`{"elemental_command":{"linked_plane":"fire","plane_chosen_by":"gm-chooses-or-rolls-randomly","bane":{"target_creature_type":"elemental","your_attacks":"advantage","attacks_against_you":"disadvantage"},"compulsion":{"action":"magic","range_feet":60,"save_ability":"wisdom","save_dc":18,"condition":"charmed","condition_ends":"start-of-your-next-turn","control":"you-determine-its-move-and-action-on-its-next-turn"},"focus":{"language":"ignan","defense":"immunity","damage_type":"fire","speed_kind":"none","speed":"none","other_benefits":[]},"spell_save_dc":18,"spells":[{"spell":"burning-hands","charge_cost":1},{"spell":"fireball","charge_cost":2},{"spell":"fire-storm","charge_cost":4},{"spell":"wall-of-fire","charge_cost":3}],"max_charges":5,"recharge":{"dice":1,"sides":4,"type":"charges"},"recharge_bonus":1,"recharge_when":"daily-at-dawn","source_pages":[237,238]}}`,
		`{"elemental_command":{"linked_plane":"water","plane_chosen_by":"gm-chooses-or-rolls-randomly","bane":{"target_creature_type":"elemental","your_attacks":"advantage","attacks_against_you":"disadvantage"},"compulsion":{"action":"magic","range_feet":60,"save_ability":"wisdom","save_dc":18,"condition":"charmed","condition_ends":"start-of-your-next-turn","control":"you-determine-its-move-and-action-on-its-next-turn"},"focus":{"language":"aquan","defense":"none","damage_type":"none","speed_kind":"swim","speed":"60-feet","other_benefits":["can-breathe-underwater"]},"spell_save_dc":18,"spells":[{"spell":"create-or-destroy-water","charge_cost":1},{"spell":"ice-storm","charge_cost":2},{"spell":"tsunami","charge_cost":5},{"spell":"wall-of-ice","charge_cost":3},{"spell":"water-walk","charge_cost":2}],"max_charges":5,"recharge":{"dice":1,"sides":4,"type":"charges"},"recharge_bonus":1,"recharge_when":"daily-at-dawn","source_pages":[237,238]}}`,
		`{"evasion":{"max_charges":3,"recharge":{"dice":1,"sides":3,"type":"charges"},"recharge_when":"daily-at-dawn","trigger":"fail-dexterity-saving-throw-while-wearing","action":"reaction","charge_cost":1,"result":"succeed-on-that-save"}}`,
		`{"feather_falling":{"condition":"while-falling-and-wearing","descent_feet_per_round":60,"takes_falling_damage":false}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:consume:drink-one:item.consume item.disguised-as-beneficial-potion damage.4d6-poison save.constitution-dc-13 condition.poisoned duration.1-hour",
		"none:consume:drink-one:item.consume defense.resistance-one-damage-type table.damage-type-1d10 duration.1-hour",
		"none:consume:drink-one:item.consume spell.effect-haste duration.1-minute spell.no-concentration spell.no-lethargy",
		"none:consume:drink-one:item.consume condition.remove-exhaustion-levels condition.end-poisoned hit-dice.maximum-healing duration.24-hours",
		"none:consume:drink-one:item.consume movement.breathe-underwater duration.24-hours",
		"wield:weapon-attack::item.attunement item.magic-weapon bonus.attack-and-damage-2 light.green-dim-10-feet item.alter-form check.dexterity-acrobatics-advantage reaction.armor-class-bonus-5 property.thrown-30-120 item.returns-to-hand",
		"wear:cast-spell:expend-one-charge:item.wear item.charges spell.cast-animal-friendship spell.cast-fear-beasts-only spell.cast-speak-with-animals spell.save-dc-13",
		"wear:magic-action::item.attunement item.wear action.magic creature.summon-djinni spell.concentration-up-to-1-hour creature.obeys-commands item.cooldown-24-hours item.nonmagical-if-djinni-dies",
		"wear:cast-spell:spell-linked-to-plane expend-listed-charges:item.attunement item.wear item.charges attack.advantage-against-elementals action.magic-compel-elemental save.wisdom-dc-18 condition.charmed feature.elemental-focus-by-linked-plane spell.cast-by-linked-plane spell.save-dc-18",
		"wear:cast-spell:spell-linked-to-plane expend-listed-charges:item.attunement item.wear item.charges attack.advantage-against-elementals action.magic-compel-elemental save.wisdom-dc-18 condition.charmed feature.elemental-focus-by-linked-plane spell.cast-by-linked-plane spell.save-dc-18",
		"wear:cast-spell:spell-linked-to-plane expend-listed-charges:item.attunement item.wear item.charges attack.advantage-against-elementals action.magic-compel-elemental save.wisdom-dc-18 condition.charmed feature.elemental-focus-by-linked-plane spell.cast-by-linked-plane spell.save-dc-18",
		"wear:cast-spell:spell-linked-to-plane expend-listed-charges:item.attunement item.wear item.charges attack.advantage-against-elementals action.magic-compel-elemental save.wisdom-dc-18 condition.charmed feature.elemental-focus-by-linked-plane spell.cast-by-linked-plane spell.save-dc-18",
		"wear:reaction:failed-dexterity-save expend-one-charge:item.attunement item.wear item.charges action.reaction save.dexterity-failure-becomes-success",
		"wear:none::item.attunement item.wear movement.fall-60-feet-per-round damage.no-falling-damage",
	}
	index := 0
	for _, f := range loadItems238(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails238
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

func TestItems238CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems238(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	identities := [][]ItemSourceIdentity{contract.SourceIdentities()}
	for _, load := range []func() ([]ItemSourceIdentity, error){
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems225(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems226(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems227(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems228(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems229(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems230(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems231(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems232(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems233(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems234(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems235(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems236(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems237(); return b.SourceIdentities(), err },
	} {
		ids, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, ids)
	}
	composed, err := ComposeItemSourceIdentities(append(identities, batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 168 || variants != 229 || composed[0].Name != "Potion of Healing" || composed[157].Name != "Potion of Poison" || composed[167].Name != "Ring of Feather Falling" {
		t.Fatalf("expected Potion of Healing, batches #225 to #237 and the eleven families of #238: %d families, %d variants", len(composed), variants)
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	// A second Potions of Healing record would repeat the contract fixture's
	// presentation, which composition rejects.
	second := batch.SourceIdentities()[:1]
	second[0].FamilyID, second[0].Name = "srd-5.2.1:content/items/potions-of-healing", "Potions of Healing"
	second[0].PresentationIDs = []string{"srd-5.2.1:content/magic-items/potions-of-healing"}
	second[0].VariantIDs = []string{"srd-5.2.1:content/items/potions-of-healing/variant/basic"}
	if _, err := ComposeItemSourceIdentities(contract.SourceIdentities(), second); err == nil {
		t.Fatal("second Potions of Healing identity accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_238") {
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
func TestItems238RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items238JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 700 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items238JSON, &fresh); err != nil {
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
		if _, err := DecodeItems238(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems238RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails238]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails238] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-13" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 237 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:10] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[11].ID = "srd-5.2.1:content/items/ring-of-free-action"
			b.Families[11].Name = "Ring of Free Action"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"quarterstaff-of-the-acrobat widened base", func(b *batch) { b.Families[5].Base.Options = append(b.Families[5].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{237, 237} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{236, 237} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"ring-of-elemental-command dropped variant", func(b *batch) { b.Families[8].Variants = b.Families[8].Variants[:3] }},
		{"ring-of-elemental-command reordered variants", func(b *batch) { s := b.Families[8].Variants; s[0], s[1] = s[1], s[0] }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"potion-of-poison rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"potion-of-resistance rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"potion-of-speed rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"potion-of-vitality rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"potion-of-water-breathing rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"quarterstaff-of-the-acrobat rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"ring-of-animal-influence rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"ring-of-djinni-summoning rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"ring-of-elemental-command rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"ring-of-evasion rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"ring-of-feather-falling rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"potion-of-poison attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "required" }},
		{"potion-of-resistance attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "required" }},
		{"potion-of-speed attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "required" }},
		{"potion-of-vitality attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "required" }},
		{"potion-of-water-breathing attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "required" }},
		{"quarterstaff-of-the-acrobat attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "none" }},
		{"ring-of-animal-influence attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "required" }},
		{"ring-of-djinni-summoning attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"ring-of-elemental-command attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "none" }},
		{"ring-of-evasion attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"ring-of-feather-falling attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 10 }},
		{"ring-of-animal-influence dropped charges", func(b *batch) { *v(b, 6, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-animal-influence charge maximum", func(b *batch) { v(b, 6, 0).Charges.Maximum++ }},
		{"ring-of-animal-influence charge recharge", func(b *batch) { v(b, 6, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"ring-of-elemental-command dropped charges", func(b *batch) { *v(b, 8, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-elemental-command charge maximum", func(b *batch) { v(b, 8, 0).Charges.Maximum++ }},
		{"ring-of-elemental-command charge recharge", func(b *batch) { v(b, 8, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"ring-of-evasion dropped charges", func(b *batch) { *v(b, 9, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-evasion charge maximum", func(b *batch) { v(b, 9, 0).Charges.Maximum++ }},
		{"ring-of-evasion charge recharge", func(b *batch) { v(b, 9, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "reaction" }},
		{"dropped use restriction", func(b *batch) { v(b, 0, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails238{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.ResistancePotion = v(b, 1, 0).Details.ResistancePotion }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.Poison.Disguise = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.Poison.Disguise = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.Poison.SaveDC = -1 }},
		{"poison damage", func(b *batch) { p := v(b, 0, 0).Details.Poison; p.Damage.Dice = 2 }},
		{"poison dc", func(b *batch) { p := v(b, 0, 0).Details.Poison; p.SaveDC = 15 }},
		{"poison duration", func(b *batch) { p := v(b, 0, 0).Details.Poison; p.DurationHours = 24 }},
		{"resistance duration", func(b *batch) { p := v(b, 1, 0).Details.ResistancePotion; p.DurationHours = 8 }},
		{"resistance die", func(b *batch) { p := v(b, 1, 0).Details.ResistancePotion; p.TableDieSides = 12 }},
		{"resistance missing row", func(b *batch) { p := v(b, 1, 0).Details.ResistancePotion; p.DamageTypes = p.DamageTypes[:9] }},
		{"resistance duplicate type", func(b *batch) { p := v(b, 1, 0).Details.ResistancePotion; p.DamageTypes[1].DamageType = "acid" }},
		{"resistance row gap", func(b *batch) { p := v(b, 1, 0).Details.ResistancePotion; p.DamageTypes[4].Roll = 6 }},
		{"speed duration", func(b *batch) { p := v(b, 2, 0).Details.Speed; p.DurationMinutes = 10 }},
		{"speed concentration", func(b *batch) { p := v(b, 2, 0).Details.Speed; p.ConcentrationRequired = true }},
		{"speed lethargy", func(b *batch) { p := v(b, 2, 0).Details.Speed; p.LethargyWhenItEnds = true }},
		{"vitality duration", func(b *batch) { p := v(b, 3, 0).Details.Vitality; p.DurationHours = 1 }},
		{"vitality condition", func(b *batch) { p := v(b, 3, 0).Details.Vitality; p.EndsCondition = "charmed" }},
		{"water breathing duration", func(b *batch) { p := v(b, 4, 0).Details.WaterBreathing; p.DurationHours = 1 }},
		{"water breathing pages", func(b *batch) { p := v(b, 4, 0).Details.WaterBreathing; p.SourcePages = []int{236} }},
		{"acrobat bonus", func(b *batch) { p := v(b, 5, 0).Details.Acrobat; p.AttackBonus = 1 }},
		{"acrobat light", func(b *batch) { p := v(b, 5, 0).Details.Acrobat; p.Light.DimRadiusFeet = 20 }},
		{"acrobat deflection", func(b *batch) { p := v(b, 5, 0).Details.Acrobat; p.AttackDeflection.ArmorClassBonus = 2 }},
		{"acrobat range", func(b *batch) { p := v(b, 5, 0).Details.Acrobat; p.RangedWeapon.LongRangeFeet = 60 }},
		{"acrobat duplicate form", func(b *batch) { p := v(b, 5, 0).Details.Acrobat; p.AlterForm.Forms[1] = "six-inch-rod" }},
		{"animal influence charges", func(b *batch) { p := v(b, 6, 0).Details.AnimalInfluence; p.MaxCharges = 4 }},
		{"animal influence dc", func(b *batch) { p := v(b, 6, 0).Details.AnimalInfluence; p.SaveDC = 15 }},
		{"animal influence spell", func(b *batch) { p := v(b, 6, 0).Details.AnimalInfluence; p.Spells[1] = "command" }},
		{"djinni range", func(b *batch) { p := v(b, 7, 0).Details.DjinniSummoning; p.SummonRangeFeet = 60 }},
		{"djinni duration", func(b *batch) { p := v(b, 7, 0).Details.DjinniSummoning; p.MaxDurationHours = 8 }},
		{"djinni cooldown", func(b *batch) { p := v(b, 7, 0).Details.DjinniSummoning; p.CooldownHours = 1 }},
		{"command plane", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.LinkedPlane = "fire" }},
		{"command spell cost", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.Spells[0].ChargeCost = 1 }},
		{"command missing spell", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.Spells = p.Spells[:3] }},
		{"command charges", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.MaxCharges = 6 }},
		{"command dc", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.SpellSaveDC = 17 }},
		{"command compulsion dc", func(b *batch) { p := v(b, 8, 0).Details.ElementalCommand; p.Compulsion.SaveDC = 15 }},
		{"evasion charges", func(b *batch) { p := v(b, 9, 0).Details.Evasion; p.MaxCharges = 5 }},
		{"evasion cost", func(b *batch) { p := v(b, 9, 0).Details.Evasion; p.ChargeCost = 2 }},
		{"evasion recharge", func(b *batch) { p := v(b, 9, 0).Details.Evasion; p.Recharge.Sides = 4 }},
		{"feather falling speed", func(b *batch) { p := v(b, 10, 0).Details.FeatherFalling; p.DescentFeetPerRound = 30 }},
		{"feather falling damage", func(b *batch) { p := v(b, 10, 0).Details.FeatherFalling; p.TakesFallingDamage = true }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems238(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems238(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems238RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/potion-of-poison/variant/standard"`)
	if !bytes.Contains(items238JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items238JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems238(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"poison": {`, `"poison_plus": {`},
		"unknown detail":     {`"dice": 4`, `"dice_all": 4`},
		"string for number":  {`"dice": 4`, `"dice": "4"`},
		"fractional number":  {`"dice": 4`, `"dice": 4.5`},
		"number for boolean": {`"concentration_required": true`, `"concentration_required": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items238JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items238JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems238(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items238JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems238(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items238JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems237(items238JSON); err == nil {
		t.Fatal("#237 decoder accepted the #238 batch")
	}
	if _, err := DecodeItems238(items237JSON); err == nil {
		t.Fatal("#238 decoder accepted the #237 batch")
	}
}
