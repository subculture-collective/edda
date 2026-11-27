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

func loadItems244(t *testing.T) ItemSourceBatch[ItemDetails244] {
	t.Helper()
	b, err := LoadItems244()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 248-250, in planning inventory order for issue #244.
func TestItems244SourceIdentities(t *testing.T) {
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
		{"Talisman of Pure Good", "talisman-of-pure-good", 248, "wondrous-item", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "none"}, []variant{{"standard", "Talisman of Pure Good", "legendary"}}},
		{"Talisman of Ultimate Evil", "talisman-of-ultimate-evil", 249, "wondrous-item", nil, "required", ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "none"}, []variant{{"standard", "Talisman of Ultimate Evil", "legendary"}}},
		{"Talisman of the Sphere", "talisman-of-the-sphere", 249, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Talisman of the Sphere", "legendary"}}},
		{"Thunderous Greatclub", "thunderous-greatclub", 249, "weapon", []string{"greatclub"}, "required", noCharges, []variant{{"standard", "Thunderous Greatclub", "very-rare"}}},
		{"Tome of Clear Thought", "tome-of-clear-thought", 249, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Tome of Clear Thought", "very-rare"}}},
		{"Tome of Leadership and Influence", "tome-of-leadership-and-influence", 249, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Tome of Leadership and Influence", "very-rare"}}},
		{"Tome of Understanding", "tome-of-understanding", 250, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Tome of Understanding", "very-rare"}}},
		{"Trident of Fish Command", "trident-of-fish-command", 250, "weapon", []string{"trident"}, "required", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Trident of Fish Command", "uncommon"}}},
		{"Universal Solvent", "universal-solvent", 250, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Universal Solvent", "legendary"}}},
		{"Vicious Weapon", "vicious-weapon", 250, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, "none", noCharges, []variant{{"standard", "Vicious Weapon", "rare"}}},
		{"Vorpal Sword", "vorpal-sword", 250, "weapon", []string{"glaive", "greatsword", "longsword", "scimitar"}, "required", noCharges, []variant{{"standard", "Vorpal Sword", "legendary"}}},
		{"Wand of Binding", "wand-of-binding", 250, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Binding", "rare"}}},
	}
	b := loadItems244(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-20" || b.CatalogueIssue != 244 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #244")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 248–250") {
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
	if total != 12 {
		t.Fatalf("expected 12 source-listed variants, got %d", total)
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
	type member struct {
		name string
		page int
	}
	var members []member
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 244 {
			continue
		}
		if r.Family != "magic-items" || r.Batch != b.Batch {
			t.Fatalf("inventory record %s changed family or batch", r.ID)
		}
		members = append(members, member{r.Name, r.Page})
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 12 || len(members) != 12 {
		t.Fatalf("inventory assigns %d families to issue #244", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "248, 249, 250" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.TalismanGood.TouchDamage.Repeats = "changed"
	if fresh := loadItems244(t); fresh.Families[0].Variants[0].Details.TalismanGood.TouchDamage.Repeats != "end-of-each-turn-holding-or-carrying" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF
// pp. 248-250.
func TestItems244SourceValues(t *testing.T) {
	expected := []string{
		`{"talisman_good":{"attunement_by_any_of":["cleric","paladin"],"touch_damage":{"creature_types":["fiend","undead"],"damage":{"dice":8,"sides":6,"type":"radiant"},"repeats":"end-of-each-turn-holding-or-carrying"},"holy_symbol":true,"spell_attack_bonus":2,"max_charges":7,"rebuke":{"action":"magic","charge_cost":1,"target":"one-visible-creature-on-ground","range_feet":120,"save_ability":"dexterity","save_dc":20,"disadvantage_types":["fiend","undead"],"failure":"destroyed-leaving-no-remains","success_damage":{"dice":4,"sides":6,"type":"psychic"}},"last_charge":"disperses-into-golden-light-and-is-destroyed","source_pages":[248,249]}}`,
		`{"talisman_evil":{"touch_damage":{"excluded_creature_types":["fiend","undead"],"damage":{"dice":8,"sides":6,"type":"necrotic"},"repeats":"end-of-each-turn-holding-or-carrying"},"holy_symbol":true,"spell_attack_bonus":2,"max_charges":6,"ultimate_end":{"action":"magic","charge_cost":1,"target":"one-visible-creature-on-ground","range_feet":120,"save_ability":"dexterity","save_dc":20,"disadvantage_types":["celestial"],"failure":"destroyed-leaving-no-remains","success_damage":{"dice":4,"sides":6,"type":"psychic"}},"last_charge":"dissolves-into-slime-and-is-destroyed"}}`,
		`{"talisman_sphere":{"condition":"while-holding-or-wearing","check":{"ability":"intelligence","skill":"arcana","mode":"advantage","purpose":"control-sphere-of-annihilation"},"move":{"timing":"start-of-turn-while-in-control","action":"magic","base_feet":10,"additional_feet_per_intelligence_modifier":10,"straight_line_required":false}}}`,
		`{"thunderous_greatclub":{"strength":{"score":20,"condition":"while-attuned-unless-already-equal-or-higher"},"creature_extra_damage":{"dice":1,"sides":8,"type":"thunder"},"object_extra_damage":{"dice":3,"sides":8,"type":"thunder"},"object_filter":"not-worn-or-carried","clap":{"action":"magic","strike":"hard-surface","audible_feet":300,"cone_feet":30,"save_ability":"strength","save_dc":15,"condition":"prone","object_damage":{"dice":3,"sides":8,"type":"thunder"}},"earthquake":{"action":"magic","strike":"ground","circle_radius_feet":50,"structure_damage":50,"structure_damage_type":"bludgeoning","prone_save_ability":"dexterity","prone_save_dc":20,"condition":"prone","concentration_save_ability":"constitution","concentration_save_dc":20,"fissure":{"depth_feet":30,"width_feet":10,"save_ability":"dexterity","save_dc":20,"failure":"falls-into-fissure","success":"moves-with-fissure-edge","structures":"collapse-into-fissure"},"reuse_after":"next-dawn"}}}`,
		`{"clear_thought":{"contents":"memory-and-logic-exercises","study":{"hours":48,"within_days":6,"activity":"study-contents-and-practice-guidelines"},"ability":"intelligence","increase":2,"maximum":30,"magic_after_use":"lost","magic_returns_after_centuries":1}}`,
		`{"leadership":{"contents":"guidelines-for-influencing-and-charming","study":{"hours":48,"within_days":6,"activity":"study-contents-and-practice-guidelines"},"ability":"charisma","increase":2,"maximum":30,"magic_after_use":"lost","magic_returns_after_centuries":1,"source_pages":[249,250]}}`,
		`{"understanding":{"contents":"intuition-and-insight-exercises","study":{"hours":48,"within_days":6,"activity":"study-contents-and-practice-guidelines"},"ability":"wisdom","increase":2,"maximum":30,"magic_after_use":"lost","magic_returns_after_centuries":1}}`,
		`{"fish_command":{"max_charges":3,"recharge_dice":{"dice":1,"sides":3,"type":"charges"},"condition":"while-carrying","charge_cost":1,"spell":"dominate-beast","save_dc":15,"target":"beast-with-swim-speed"}}`,
		`{"solvent":{"appearance":"milky-liquid-with-strong-alcohol-smell","ounces_when_found":{"dice":1,"sides":6,"type":"ounces"},"ounces_when_found_bonus":1,"action":"utilize","minimum_ounces_per_pour":1,"target":"surface-within-reach","dissolves":{"maximum_square_feet":1,"timing":"instantly","includes":"sovereign-glue"}}}`,
		`{"vicious":{"trigger":"hit-creature","extra_damage":{"dice":2,"sides":6,"type":"same-as-weapon"},"damage_type_rule":"same-as-weapon-normal-damage"}}`,
		`{"vorpal":{"attack_bonus":3,"damage_bonus":3,"ignores_resistance_to":"slashing","head":{"target":"creature-with-at-least-one-head","trigger_natural_roll":20,"effect":"cut-off-one-head","death":"dies-if-it-cannot-survive-without-head","immune_when":["immunity-to-slashing-damage","no-head-or-no-need-for-one","gm-decides-too-big"],"immune_extra_damage":30,"immune_extra_damage_type":"slashing","legendary_resistance":{"cost":"one-daily-use","result":"keeps-head-and-takes-extra-damage"}}}}`,
		`{"binding":{"max_charges":7,"save_dc":17,"spells":[{"spell":"hold-monster","charge_cost":5},{"spell":"hold-person","charge_cost":2}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear-or-hold:magic-action:expend-one-charge target-on-ground-within-120-feet:item.attunement attunement.by-class damage.8d6-radiant-on-fiend-or-undead-touch item.holy-symbol bonus.spell-attack-2 item.charges action.magic-pure-rebuke save.dexterity-dc-20 effect.destroyed-in-fissure damage.4d6-psychic-on-success item.destroyed-after-last-charge",
		"wear-or-hold:magic-action:expend-one-charge target-on-ground-within-120-feet:item.attunement damage.8d6-necrotic-on-other-creature-touch item.holy-symbol bonus.spell-attack-2 item.charges action.magic-ultimate-end save.dexterity-dc-20 effect.destroyed-in-fissure damage.4d6-psychic-on-success item.destroyed-after-last-charge",
		"wear-or-hold:magic-action:controlling-sphere-of-annihilation-at-turn-start:item.attunement check.advantage-arcana-to-control-sphere action.magic-move-sphere movement.10-feet-plus-10-times-intelligence-modifier",
		"wield:weapon-attack:earthquake-once-until-next-dawn:item.attunement item.magic-weapon ability.strength-20 damage.extra-1d8-thunder damage.extra-3d8-thunder-to-objects action.magic-clap-of-thunder save.strength-dc-15 condition.prone damage.3d8-thunder-to-objects-in-cone action.magic-earthquake damage.50-bludgeoning-to-structures save.dexterity-dc-20 save.constitution-dc-20-concentration effect.fissure item.earthquake-once-until-dawn",
		"none:study:48-hours-within-6-days:item.study-48-hours-within-6-days ability.increase-intelligence-by-2 ability.maximum-30 item.loses-magic-for-a-century",
		"none:study:48-hours-within-6-days:item.study-48-hours-within-6-days ability.increase-charisma-by-2 ability.maximum-30 item.loses-magic-for-a-century",
		"none:study:48-hours-within-6-days:item.study-48-hours-within-6-days ability.increase-wisdom-by-2 ability.maximum-30 item.loses-magic-for-a-century",
		"carry:cast-spell:expend-one-charge beast-with-swim-speed:item.attunement item.magic-weapon item.carry item.charges spell.cast-dominate-beast spell.save-dc-15 spell.target-beast-with-swim-speed",
		"none:utilize-action:pour-onto-surface-within-reach:item.found-with-1d6-plus-1-ounces action.utilize effect.dissolve-adhesive-1-square-foot-per-ounce",
		"wield:weapon-attack::item.magic-weapon damage.extra-2d6-of-weapon-type",
		"wield:weapon-attack::item.attunement item.magic-weapon bonus.attack-and-damage-3 damage.ignore-slashing-resistance effect.sever-head-on-natural-20 damage.extra-30-slashing-when-head-kept resistance.legendary-resistance-avoids-head-loss",
		"hold:cast-spell:expend-listed-charges:item.attunement item.hold item.charges spell.cast-listed-spells spell.save-dc-17 item.destroyed-on-1-after-last-charge",
	}
	index := 0
	for _, f := range loadItems244(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails244
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

func TestItems244CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems244(t)
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
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems238(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems239(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems240(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems241(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems242(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems243(); return b.SourceIdentities(), err },
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
	if len(composed) != 240 || variants != 312 || composed[0].Name != "Potion of Healing" || composed[228].Name != "Talisman of Pure Good" || composed[239].Name != "Wand of Binding" {
		t.Fatalf("expected Potion of Healing, batches #225 to #243 and the twelve families of #244: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_244") {
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
func TestItems244RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items244JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items244JSON, &fresh); err != nil {
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
		if _, err := DecodeItems244(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems244RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails244]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails244] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-19" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 243 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/adamantine-armor"
			b.Families[12].Name = "Adamantine Armor"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"thunderous-greatclub widened base", func(b *batch) { b.Families[3].Base.Options = append(b.Families[3].Base.Options, "any-weapon") }},
		{"trident-of-fish-command widened base", func(b *batch) { b.Families[7].Base.Options = append(b.Families[7].Base.Options, "any-weapon") }},
		{"vicious-weapon widened base", func(b *batch) { b.Families[9].Base.Options = append(b.Families[9].Base.Options, "any-weapon") }},
		{"vorpal-sword widened base", func(b *batch) { b.Families[10].Base.Options = append(b.Families[10].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{249, 249} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{248, 249} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"talisman-of-pure-good rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"talisman-of-pure-good attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"talisman-of-ultimate-evil rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"talisman-of-ultimate-evil attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"talisman-of-the-sphere rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"talisman-of-the-sphere attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "none" }},
		{"thunderous-greatclub rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"thunderous-greatclub attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "none" }},
		{"tome-of-clear-thought rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"tome-of-clear-thought attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "required" }},
		{"tome-of-leadership-and-influence rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"tome-of-leadership-and-influence attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "required" }},
		{"tome-of-understanding rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"tome-of-understanding attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "required" }},
		{"trident-of-fish-command rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"trident-of-fish-command attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"universal-solvent rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"universal-solvent attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "required" }},
		{"vicious-weapon rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"vicious-weapon attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "required" }},
		{"vorpal-sword rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"vorpal-sword attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"wand-of-binding rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"wand-of-binding attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 2, 0).Charges.Maximum = 10 }},
		{"talisman-of-pure-good dropped charges", func(b *batch) { *v(b, 0, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"talisman-of-pure-good charge maximum", func(b *batch) { v(b, 0, 0).Charges.Maximum++ }},
		{"talisman-of-pure-good charge recharge", func(b *batch) { v(b, 0, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"talisman-of-ultimate-evil dropped charges", func(b *batch) { *v(b, 1, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"talisman-of-ultimate-evil charge maximum", func(b *batch) { v(b, 1, 0).Charges.Maximum++ }},
		{"talisman-of-ultimate-evil charge recharge", func(b *batch) { v(b, 1, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"trident-of-fish-command dropped charges", func(b *batch) { *v(b, 7, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"trident-of-fish-command charge maximum", func(b *batch) { v(b, 7, 0).Charges.Maximum++ }},
		{"trident-of-fish-command charge recharge", func(b *batch) { v(b, 7, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-binding dropped charges", func(b *batch) { *v(b, 11, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-binding charge maximum", func(b *batch) { v(b, 11, 0).Charges.Maximum++ }},
		{"wand-of-binding charge recharge", func(b *batch) { v(b, 11, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "unreviewed" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "unreviewed" }},
		{"dropped use restriction", func(b *batch) { v(b, 0, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails244{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.TalismanEvil = v(b, 1, 0).Details.TalismanEvil }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.TalismanGood.TouchDamage.Repeats = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.TalismanGood.TouchDamage.Repeats = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.TalismanGood.SpellAttackBonus = -1 }},
		{"good charges", func(b *batch) { p := v(b, 0, 0).Details.TalismanGood; p.MaxCharges = 6 }},
		{"good save", func(b *batch) { p := v(b, 0, 0).Details.TalismanGood; p.Rebuke.SaveDC = 17 }},
		{"good touch dice", func(b *batch) { p := v(b, 0, 0).Details.TalismanGood; p.TouchDamage.Damage.Dice = 6 }},
		{"good disadvantage", func(b *batch) {
			p := v(b, 0, 0).Details.TalismanGood
			p.Rebuke.DisadvantageTypes = []string{"celestial"}
		}},
		{"good pages", func(b *batch) { p := v(b, 0, 0).Details.TalismanGood; p.SourcePages = []int{248} }},
		{"evil charges", func(b *batch) { p := v(b, 1, 0).Details.TalismanEvil; p.MaxCharges = 7 }},
		{"evil range", func(b *batch) { p := v(b, 1, 0).Details.TalismanEvil; p.UltimateEnd.RangeFeet = 60 }},
		{"evil disadvantage", func(b *batch) {
			p := v(b, 1, 0).Details.TalismanEvil
			p.UltimateEnd.DisadvantageTypes = []string{"fiend", "undead"}
		}},
		{"evil exclusions", func(b *batch) {
			p := v(b, 1, 0).Details.TalismanEvil
			p.TouchDamage.ExcludedCreatureTypes = p.TouchDamage.ExcludedCreatureTypes[:1]
		}},
		{"sphere skill", func(b *batch) { p := v(b, 2, 0).Details.TalismanSphere; p.Check.Skill = "history" }},
		{"sphere base", func(b *batch) { p := v(b, 2, 0).Details.TalismanSphere; p.Move.BaseFeet = 5 }},
		{"sphere line", func(b *batch) { p := v(b, 2, 0).Details.TalismanSphere; p.Move.StraightLineRequired = true }},
		{"greatclub strength", func(b *batch) { p := v(b, 3, 0).Details.ThunderousGreatclub; p.Strength.Score = 19 }},
		{"greatclub object dice", func(b *batch) { p := v(b, 3, 0).Details.ThunderousGreatclub; p.ObjectExtraDamage.Dice = 1 }},
		{"greatclub clap dc", func(b *batch) { p := v(b, 3, 0).Details.ThunderousGreatclub; p.Clap.SaveDC = 20 }},
		{"greatclub structure", func(b *batch) { p := v(b, 3, 0).Details.ThunderousGreatclub; p.Earthquake.StructureDamage = 30 }},
		{"greatclub fissure", func(b *batch) { p := v(b, 3, 0).Details.ThunderousGreatclub; p.Earthquake.Fissure.SaveDC = 15 }},
		{"intelligence tome ability", func(b *batch) { p := v(b, 4, 0).Details.ClearThought; p.Ability = "strength" }},
		{"intelligence tome maximum", func(b *batch) { p := v(b, 4, 0).Details.ClearThought; p.Maximum = 24 }},
		{"intelligence tome hours", func(b *batch) { p := v(b, 4, 0).Details.ClearThought; p.Study.Hours = 24 }},
		{"charisma tome ability", func(b *batch) { p := v(b, 5, 0).Details.Leadership; p.Ability = "strength" }},
		{"charisma tome maximum", func(b *batch) { p := v(b, 5, 0).Details.Leadership; p.Maximum = 24 }},
		{"charisma tome hours", func(b *batch) { p := v(b, 5, 0).Details.Leadership; p.Study.Hours = 24 }},
		{"wisdom tome ability", func(b *batch) { p := v(b, 6, 0).Details.Understanding; p.Ability = "strength" }},
		{"wisdom tome maximum", func(b *batch) { p := v(b, 6, 0).Details.Understanding; p.Maximum = 24 }},
		{"wisdom tome hours", func(b *batch) { p := v(b, 6, 0).Details.Understanding; p.Study.Hours = 24 }},
		{"trident charges", func(b *batch) { p := v(b, 7, 0).Details.FishCommand; p.MaxCharges = 7 }},
		{"trident dc", func(b *batch) { p := v(b, 7, 0).Details.FishCommand; p.SaveDC = 13 }},
		{"trident spell", func(b *batch) { p := v(b, 7, 0).Details.FishCommand; p.Spell = "dominate-monster" }},
		{"solvent dice", func(b *batch) { p := v(b, 8, 0).Details.Solvent; p.OuncesWhenFound.Sides = 4 }},
		{"solvent bonus", func(b *batch) { p := v(b, 8, 0).Details.Solvent; p.OuncesWhenFoundBonus = 0 }},
		{"solvent area", func(b *batch) { p := v(b, 8, 0).Details.Solvent; p.Dissolves.MaximumSquareFeet = 5 }},
		{"vicious dice", func(b *batch) { p := v(b, 9, 0).Details.Vicious; p.ExtraDamage.Dice = 1 }},
		{"vicious type", func(b *batch) { p := v(b, 9, 0).Details.Vicious; p.ExtraDamage.Type = "necrotic" }},
		{"vorpal bonus", func(b *batch) { p := v(b, 10, 0).Details.Vorpal; p.AttackBonus = 2 }},
		{"vorpal roll", func(b *batch) { p := v(b, 10, 0).Details.Vorpal; p.Head.TriggerNaturalRoll = 19 }},
		{"vorpal extra", func(b *batch) { p := v(b, 10, 0).Details.Vorpal; p.Head.ImmuneExtraDamage = 14 }},
		{"vorpal immunity", func(b *batch) { p := v(b, 10, 0).Details.Vorpal; p.Head.ImmuneWhen = p.Head.ImmuneWhen[:2] }},
		{"binding dc", func(b *batch) { p := v(b, 11, 0).Details.Binding; p.SaveDC = 15 }},
		{"binding cost", func(b *batch) { p := v(b, 11, 0).Details.Binding; p.Spells[1].ChargeCost = 3 }},
		{"binding recharge", func(b *batch) { p := v(b, 11, 0).Details.Binding; p.Recharge.Bonus = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems244(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems244(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems244RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/talisman-of-pure-good/variant/standard"`)
	if !bytes.Contains(items244JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items244JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems244(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"talisman_good": {`, `"talisman_good_plus": {`},
		"unknown detail":     {`"spell_attack_bonus": 2`, `"spell_attack_bonus_all": 2`},
		"string for number":  {`"spell_attack_bonus": 2`, `"spell_attack_bonus": "2"`},
		"fractional number":  {`"spell_attack_bonus": 2`, `"spell_attack_bonus": 2.5`},
		"number for boolean": {`"holy_symbol": true`, `"holy_symbol": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items244JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items244JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems244(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items244JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems244(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items244JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems243(items244JSON); err == nil {
		t.Fatal("#243 decoder accepted the #244 batch")
	}
	if _, err := DecodeItems244(items243JSON); err == nil {
		t.Fatal("#244 decoder accepted the #243 batch")
	}
}
