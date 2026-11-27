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

func loadItems240(t *testing.T) ItemSourceBatch[ItemDetails240] {
	t.Helper()
	b, err := LoadItems240()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 239-241, in planning inventory order for issue #240.
func TestItems240SourceIdentities(t *testing.T) {
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
		{"Ring of Three Wishes", "ring-of-three-wishes", 239, "ring", nil, "none", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "none"}, []variant{{"standard", "Ring of Three Wishes", "legendary"}}},
		{"Ring of Warmth", "ring-of-warmth", 239, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Warmth", "uncommon"}}},
		{"Ring of Water Walking", "ring-of-water-walking", 239, "ring", nil, "none", noCharges, []variant{{"standard", "Ring of Water Walking", "uncommon"}}},
		{"Ring of X-ray Vision", "ring-of-x-ray-vision", 239, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of X-ray Vision", "rare"}}},
		{"Ring of the Ram", "ring-of-the-ram", 239, "ring", nil, "required", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Ring of the Ram", "rare"}}},
		{"Robe of Eyes", "robe-of-eyes", 240, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Robe of Eyes", "rare"}}},
		{"Robe of Scintillating Colors", "robe-of-scintillating-colors", 240, "wondrous-item", nil, "required", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Robe of Scintillating Colors", "very-rare"}}},
		{"Robe of Stars", "robe-of-stars", 240, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Robe of Stars", "very-rare"}}},
		{"Robe of Useful Items", "robe-of-useful-items", 240, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Robe of Useful Items", "uncommon"}}},
		{"Robe of the Archmagi", "robe-of-the-archmagi", 240, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Robe of the Archmagi", "legendary"}}},
		{"Rod of Absorption", "rod-of-absorption", 241, "rod", nil, "required", noCharges, []variant{{"standard", "Rod of Absorption", "very-rare"}}},
		{"Rod of Alertness", "rod-of-alertness", 241, "rod", nil, "required", noCharges, []variant{{"standard", "Rod of Alertness", "very-rare"}}},
	}
	b := loadItems240(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-16" || b.CatalogueIssue != 240 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #240")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 239–241") {
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
		if r.Issue != 240 {
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
		t.Fatalf("inventory assigns %d families to issue #240", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "239, 240, 241" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.ThreeWishes.Condition = "changed"
	if fresh := loadItems240(t); fresh.Families[0].Variants[0].Details.ThreeWishes.Condition != "while-wearing" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 239-241.
func TestItems240SourceValues(t *testing.T) {
	expected := []string{
		`{"three_wishes":{"condition":"while-wearing","spell":"wish","max_charges":3,"charge_cost":1,"after_last_charge":"becomes-nonmagical"}}`,
		`{"warmth":{"trigger":"take-cold-damage-while-wearing","damage_type":"cold","reduction":{"dice":2,"sides":8,"type":"damage-reduction"},"protected":["you","everything-you-wear-and-carry"],"protected_at_or_below_fahrenheit":0}}`,
		`{"water_walking":{"condition":"while-wearing","spell":"water-walk","target":"yourself-only"}}`,
		`{"x_ray":{"condition":"while-wearing","action":"magic","range_feet":30,"duration_minutes":1,"solid_objects":"appear-transparent-and-pass-light","penetration":{"stone_feet":1,"common_metal_inches":1,"wood_or_dirt_feet":3,"blocked_by":["thicker-substances","thin-sheet-of-lead"]},"reuse":{"trigger":"use-again-before-long-rest","save_ability":"constitution","save_dc":15,"failure_condition":"exhaustion","exhaustion_levels":1},"source_pages":[239,240]}}`,
		`{"ram":{"max_charges":3,"recharge":{"dice":1,"sides":3,"type":"charges"},"recharge_when":"daily-at-dawn","action":"magic","minimum_charges":1,"maximum_charges":3,"attack":{"kind":"ranged-spell-attack","targets":1,"range_feet":60,"form":"spectral-rams-head","attack_bonus":7,"damage_per_charge":{"dice":2,"sides":10,"type":"force"},"push_feet_per_charge":5},"break_object":{"action":"magic","target":["nonmagical-object","not-worn-or-carried"],"range_feet":60,"check_ability":"strength","check_bonus_per_charge":5}}}`,
		`{"eyes":{"appearance":"eyelike-patterns","all_around_vision":{"check_ability":"wisdom","check_skill":"perception","check_mode":"advantage","check_requirement":"relies-on-sight"},"senses":["darkvision","truesight"],"sense_range_feet":120,"drawbacks":{"condition":"blinded","duration_minutes":1,"triggers":[{"spell":"light","placement":"cast-on-the-robe","save_dc":11},{"spell":"daylight","placement":"cast-within-5-feet-of-the-robe","save_dc":15}],"save_ability":"constitution","save_when":"end-of-each-of-your-turns","success":"ends-condition-on-yourself"}}}`,
		`{"scintillating":{"max_charges":3,"recharge":{"dice":1,"sides":3,"type":"charges"},"recharge_when":"daily-at-dawn","condition":"while-wearing","action":"magic","charge_cost":1,"effect":"shifting-pattern-of-dazzling-hues","duration":"until-end-of-your-next-turn","bright_light_radius_feet":30,"dim_light_additional_feet":30,"attacks_against_you":"disadvantage","stun":{"affects":"creature-in-bright-light-that-can-see-you-on-activation","save_ability":"wisdom","save_dc":15,"condition":"stunned","condition_ends":"when-effect-ends"}}}`,
		`{"stars":{"appearance":"black-or-dark-blue-robe-with-small-white-or-silver-stars","saving_throw_bonus":1,"missiles":{"large_stars":6,"action":"magic","spell":"magic-missile","spell_level":5,"stars_return":{"dice":1,"sides":6,"type":"stars"},"return_when":"daily-at-dusk"},"astral_travel":{"action":"magic","destination":"astral-plane","takes":"everything-worn-and-carried","remain_until":"magic-action-to-return","reappear_at":"last-space-or-nearest-unoccupied-space"}}}`,
		`{"useful_items":{"condition":"while-wearing","action":"magic","patch_becomes":"object-or-creature-it-represents","after_last_patch":"ordinary-garment","copies_of_each_fixed_patch":2,"fixed_patches":["bullseye-lantern-filled-and-lit","dagger","mirror","pole","rope-coiled","sack"],"other_patches":{"dice":4,"sides":4,"type":"patches"},"other_patches_chosen_by":"gm-chooses-or-rolls-randomly","table_die_sides":100,"patch_table":[{"min":1,"max":8,"patch":"bag-of-100-gp"},{"min":9,"max":15,"patch":"silver-coffer"},{"min":16,"max":22,"patch":"iron-door"},{"min":23,"max":30,"patch":"gems-worth-100-gp-each"},{"min":31,"max":44,"patch":"wooden-ladder"},{"min":45,"max":51,"patch":"riding-horse-with-riding-saddle"},{"min":52,"max":59,"patch":"open-pit"},{"min":60,"max":68,"patch":"potions-of-healing"},{"min":69,"max":75,"patch":"rowboat"},{"min":76,"max":83,"patch":"spell-scroll-of-level-1-2-or-3"},{"min":84,"max":90,"patch":"mastiffs"},{"min":91,"max":96,"patch":"window"},{"min":97,"max":100,"patch":"portable-ram"}],"patch_facts":[{"patch":"gems-worth-100-gp-each","key":"count","amount":10,"unit":"gems"},{"patch":"potions-of-healing","key":"count","amount":4,"unit":"potions"},{"patch":"mastiffs","key":"count","amount":2,"unit":"mastiffs"},{"patch":"spell-scroll-of-level-1-2-or-3","key":"minimum-spell-level","amount":1,"unit":"spell-level"},{"patch":"spell-scroll-of-level-1-2-or-3","key":"maximum-spell-level","amount":3,"unit":"spell-level"},{"patch":"bag-of-100-gp","key":"contents-value","amount":100,"unit":"gp"},{"patch":"silver-coffer","key":"length","amount":1,"unit":"feet"},{"patch":"silver-coffer","key":"width-and-depth","amount":6,"unit":"inches"},{"patch":"silver-coffer","key":"value","amount":500,"unit":"gp"},{"patch":"iron-door","key":"maximum-width","amount":10,"unit":"feet"},{"patch":"iron-door","key":"maximum-height","amount":10,"unit":"feet"},{"patch":"gems-worth-100-gp-each","key":"value-each","amount":100,"unit":"gp"},{"patch":"wooden-ladder","key":"length","amount":24,"unit":"feet"},{"patch":"open-pit","key":"cube-size","amount":10,"unit":"feet"},{"patch":"open-pit","key":"placement-range","amount":10,"unit":"feet"},{"patch":"rowboat","key":"length","amount":12,"unit":"feet"},{"patch":"window","key":"width","amount":2,"unit":"feet"},{"patch":"window","key":"height","amount":4,"unit":"feet"},{"patch":"window","key":"maximum-depth","amount":2,"unit":"feet"}],"patch_rules":[{"patch":"iron-door","rule":"barred-on-one-side-of-your-choice"},{"patch":"iron-door","rule":"placed-in-an-opening-you-can-reach"},{"patch":"iron-door","rule":"conforms-attaches-and-hinges-itself"},{"patch":"window","rule":"placed-on-a-vertical-surface-you-can-reach"}],"source_pages":[240,241]}}`,
		`{"archmagi":{"attunement_by_any_of":["sorcerer","warlock","wizard"],"armor":{"requirement":"not-wearing-armor","base_armor_class":15,"adds_ability_modifier":"dexterity"},"magic_resistance":{"save_mode":"advantage","against":["spells","other-magical-effects"]},"war_mage":{"increases":["spell-save-dc","spell-attack-bonus"],"bonus":2}}}`,
		`{"absorption":{"condition":"while-holding","absorb":{"action":"reaction","requires":["spell-targets-only-you","spell-creates-no-area-of-effect"],"result":"spell-canceled-and-energy-stored","energy_level":"level-of-spell-as-cast","casting_resources":"wasted","if_it_cannot_store":"rod-has-no-effect"},"lifetime_capacity_levels":50,"known_on_attunement":["lifetime-levels-absorbed","levels-currently-stored"],"convert":{"requirement":"spellcaster-holding-rod","result":"spell-slots-for-prepared-or-known-spells","level_limit":"equal-to-or-lower-than-your-own-slots","maximum_slot_level":5,"casting":"stored-levels-replace-slots-otherwise-cast-as-normal"},"found_stored_levels":{"dice":1,"sides":10,"type":"spell-energy-levels"},"nonmagical_when":["cannot-absorb-more","no-energy-remaining"]}}`,
		`{"alertness":{"condition":"while-holding","advantage_on":["wisdom-perception-checks","initiative-rolls"],"spells":["detect-evil-and-good","detect-magic","detect-poison-and-disease","see-invisibility"],"protective_aura":{"action":"magic","bright_light_radius_feet":60,"dim_light_additional_feet":60,"beneficiaries":["you","your-allies"],"benefit_requirement":"while-in-the-bright-light","armor_class_bonus":1,"saving_throw_bonus":1,"senses":"location-of-invisible-creatures-in-the-bright-light","duration_minutes":10,"ends_early_when":"creature-takes-magic-action-to-pull-rod","reuse_after":"next-dawn"}}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:cast-spell:expend-one-charge:item.wear item.charges spell.cast-wish item.nonmagical-after-last-charge",
		"wear:none::item.attunement item.wear damage.reduce-cold-by-2d8 defense.unharmed-at-0-degrees-fahrenheit-or-lower",
		"wear:cast-spell:water-walk-targeting-yourself:item.wear spell.cast-water-walk spell.target-self-only",
		"wear:magic-action::item.attunement item.wear action.magic senses.x-ray-vision-30-feet duration.1-minute save.constitution-dc-15 condition.exhaustion",
		"wear:magic-action:expend-one-to-three-charges:item.attunement item.wear item.charges action.magic attack.ranged-spell-plus-7 damage.force-2d10-per-charge movement.push-5-feet-per-charge check.break-object-strength-plus-5-per-charge",
		"wear:none::item.attunement item.wear check.perception-sight-advantage senses.darkvision-120-feet senses.truesight-120-feet effect.blinded-by-light-or-daylight save.constitution-end-blinded",
		"wear:magic-action:expend-one-charge:item.attunement item.wear item.charges action.magic light.bright-30-feet-dim-30-feet attack.disadvantage-against-wearer save.wisdom-dc-15 condition.stunned",
		"wear:magic-action::item.attunement item.wear save.all-bonus-1 spell.cast-magic-missile-level-5 item.removable-stars movement.enter-and-leave-astral-plane",
		"wear:magic-action:detach-one-patch:item.wear action.magic item.detachable-patches table.other-patches-1d100 item.ordinary-after-last-patch",
		"wear:none::item.attunement attunement.sorcerer-warlock-or-wizard item.wear defense.base-armor-class-15-plus-dexterity-modifier save.magic-advantage bonus.spell-save-dc-and-attack-2",
		"hold:reaction:spell-targets-only-you-without-area:item.attunement item.hold action.reaction-absorb-spell item.lifetime-capacity-50-spell-levels spell.convert-stored-energy-to-slots item.nonmagical-when-exhausted",
		"hold:cast-spell:listed-rod-spells:item.attunement item.hold check.perception-advantage check.initiative-advantage spell.cast-listed-spells action.magic-protective-aura bonus.armor-class-and-saves-1 effect.sense-invisible-creatures item.once-until-dawn",
	}
	index := 0
	for _, f := range loadItems240(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails240
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

func TestItems240CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems240(t)
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
	if len(composed) != 192 || variants != 253 || composed[0].Name != "Potion of Healing" || composed[180].Name != "Ring of Three Wishes" || composed[191].Name != "Rod of Alertness" {
		t.Fatalf("expected Potion of Healing, batches #225 to #239 and the twelve families of #240: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_240") {
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
func TestItems240RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items240JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 700 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items240JSON, &fresh); err != nil {
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
		if _, err := DecodeItems240(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems240RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails240]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails240] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-15" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 239 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/rod-of-lordly-might"
			b.Families[12].Name = "Rod of Lordly Might"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{240, 240} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{239, 240} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"ring-of-three-wishes rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"ring-of-warmth rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"ring-of-water-walking rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"ring-of-x-ray-vision rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"ring-of-the-ram rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"robe-of-eyes rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"robe-of-scintillating-colors rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"robe-of-stars rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"robe-of-useful-items rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"robe-of-the-archmagi rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"rod-of-absorption rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"rod-of-alertness rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"ring-of-three-wishes attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "required" }},
		{"ring-of-warmth attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"ring-of-water-walking attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "required" }},
		{"ring-of-x-ray-vision attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "none" }},
		{"ring-of-the-ram attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "none" }},
		{"robe-of-eyes attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "none" }},
		{"robe-of-scintillating-colors attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "none" }},
		{"robe-of-stars attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"robe-of-useful-items attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "required" }},
		{"robe-of-the-archmagi attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"rod-of-absorption attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"rod-of-alertness attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 1, 0).Charges.Maximum = 10 }},
		{"ring-of-three-wishes dropped charges", func(b *batch) { *v(b, 0, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-three-wishes charge maximum", func(b *batch) { v(b, 0, 0).Charges.Maximum++ }},
		{"ring-of-three-wishes charge recharge", func(b *batch) { v(b, 0, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"ring-of-the-ram dropped charges", func(b *batch) { *v(b, 4, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-the-ram charge maximum", func(b *batch) { v(b, 4, 0).Charges.Maximum++ }},
		{"ring-of-the-ram charge recharge", func(b *batch) { v(b, 4, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"robe-of-scintillating-colors dropped charges", func(b *batch) { *v(b, 6, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"robe-of-scintillating-colors charge maximum", func(b *batch) { v(b, 6, 0).Charges.Maximum++ }},
		{"robe-of-scintillating-colors charge recharge", func(b *batch) { v(b, 6, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "hold" }},
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
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails240{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Warmth = v(b, 1, 0).Details.Warmth }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.ThreeWishes.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.ThreeWishes.Condition = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.ThreeWishes.MaxCharges = -1 }},
		{"wishes count", func(b *batch) { p := v(b, 0, 0).Details.ThreeWishes; p.MaxCharges = 4 }},
		{"wishes cost", func(b *batch) { p := v(b, 0, 0).Details.ThreeWishes; p.ChargeCost = 3 }},
		{"wishes spell", func(b *batch) { p := v(b, 0, 0).Details.ThreeWishes; p.Spell = "limited-wish" }},
		{"warmth dice", func(b *batch) { p := v(b, 1, 0).Details.Warmth; p.Reduction.Sides = 6 }},
		{"warmth temperature", func(b *batch) { p := v(b, 1, 0).Details.Warmth; p.ProtectedAtOrBelowFahrenheit = 32 }},
		{"water walking spell", func(b *batch) { p := v(b, 2, 0).Details.WaterWalking; p.Spell = "water-breathing" }},
		{"water walking target", func(b *batch) { p := v(b, 2, 0).Details.WaterWalking; p.Target = "any-creature" }},
		{"x-ray range", func(b *batch) { p := v(b, 3, 0).Details.XRay; p.RangeFeet = 60 }},
		{"x-ray duration", func(b *batch) { p := v(b, 3, 0).Details.XRay; p.DurationMinutes = 10 }},
		{"x-ray dc", func(b *batch) { p := v(b, 3, 0).Details.XRay; p.Reuse.SaveDC = 10 }},
		{"x-ray pages", func(b *batch) { p := v(b, 3, 0).Details.XRay; p.SourcePages = []int{239} }},
		{"ram charges", func(b *batch) { p := v(b, 4, 0).Details.Ram; p.MaxCharges = 6 }},
		{"ram bonus", func(b *batch) { p := v(b, 4, 0).Details.Ram; p.Attack.AttackBonus = 5 }},
		{"ram damage", func(b *batch) { p := v(b, 4, 0).Details.Ram; p.Attack.DamagePerCharge.Sides = 6 }},
		{"ram push", func(b *batch) { p := v(b, 4, 0).Details.Ram; p.Attack.PushFeetPerCharge = 10 }},
		{"ram break", func(b *batch) { p := v(b, 4, 0).Details.Ram; p.BreakObject.CheckBonusPerCharge = 7 }},
		{"eyes range", func(b *batch) { p := v(b, 5, 0).Details.Eyes; p.SenseRangeFeet = 60 }},
		{"eyes light dc", func(b *batch) { p := v(b, 5, 0).Details.Eyes; p.Drawbacks.Triggers[0].SaveDC = 15 }},
		{"eyes dropped trigger", func(b *batch) { p := v(b, 5, 0).Details.Eyes; p.Drawbacks.Triggers = p.Drawbacks.Triggers[:1] }},
		{"eyes duplicate sense", func(b *batch) { p := v(b, 5, 0).Details.Eyes; p.Senses[1] = "darkvision" }},
		{"colors charges", func(b *batch) { p := v(b, 6, 0).Details.Scintillating; p.MaxCharges = 4 }},
		{"colors light", func(b *batch) { p := v(b, 6, 0).Details.Scintillating; p.BrightLightRadiusFeet = 60 }},
		{"colors dc", func(b *batch) { p := v(b, 6, 0).Details.Scintillating; p.Stun.SaveDC = 17 }},
		{"stars bonus", func(b *batch) { p := v(b, 7, 0).Details.Stars; p.SavingThrowBonus = 2 }},
		{"stars count", func(b *batch) { p := v(b, 7, 0).Details.Stars; p.Missiles.LargeStars = 7 }},
		{"stars level", func(b *batch) { p := v(b, 7, 0).Details.Stars; p.Missiles.SpellLevel = 3 }},
		{"stars return", func(b *batch) { p := v(b, 7, 0).Details.Stars; p.Missiles.StarsReturn.Sides = 4 }},
		{"patches copies", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.CopiesOfEachFixedPatch = 1 }},
		{"patches dice", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.OtherPatches.Sides = 6 }},
		{"patches missing row", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchTable = p.PatchTable[:12] }},
		{"patches row gap", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchTable[1].Min = 10 }},
		{"patches last is 00", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchTable[12].Max = 99 }},
		{"patches duplicate", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchTable[1].Patch = "bag-of-100-gp" }},
		{"patches orphan fact", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchFacts[0].Patch = "anvil" }},
		{"patches zero fact", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchFacts[0].Amount = 0 }},
		{"patches orphan rule", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.PatchRules[0].Patch = "anvil" }},
		{"patches fixed", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.FixedPatches = p.FixedPatches[:5] }},
		{"patches pages", func(b *batch) { p := v(b, 8, 0).Details.UsefulItems; p.SourcePages = []int{240} }},
		{"archmagi armor", func(b *batch) { p := v(b, 9, 0).Details.Archmagi; p.Armor.BaseArmorClass = 13 }},
		{"archmagi bonus", func(b *batch) { p := v(b, 9, 0).Details.Archmagi; p.WarMage.Bonus = 1 }},
		{"archmagi class", func(b *batch) {
			p := v(b, 9, 0).Details.Archmagi
			p.AttunementByAnyOf = append(p.AttunementByAnyOf, "cleric")
		}},
		{"absorption capacity", func(b *batch) { p := v(b, 10, 0).Details.Absorption; p.LifetimeCapacityLevels = 100 }},
		{"absorption slot", func(b *batch) { p := v(b, 10, 0).Details.Absorption; p.Convert.MaximumSlotLevel = 9 }},
		{"absorption found", func(b *batch) { p := v(b, 10, 0).Details.Absorption; p.FoundStoredLevels.Sides = 8 }},
		{"alertness light", func(b *batch) { p := v(b, 11, 0).Details.Alertness; p.ProtectiveAura.BrightLightRadiusFeet = 30 }},
		{"alertness bonus", func(b *batch) { p := v(b, 11, 0).Details.Alertness; p.ProtectiveAura.ArmorClassBonus = 2 }},
		{"alertness spells", func(b *batch) { p := v(b, 11, 0).Details.Alertness; p.Spells = p.Spells[:3] }},
		{"alertness duration", func(b *batch) { p := v(b, 11, 0).Details.Alertness; p.ProtectiveAura.DurationMinutes = 1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems240(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems240(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems240RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/ring-of-three-wishes/variant/standard"`)
	if !bytes.Contains(items240JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items240JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems240(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":      {`"details": {`, `"details": null, "x": {`},
		"null list":         {`"excluded": []`, `"excluded": null`},
		"unknown payload":   {`"three_wishes": {`, `"three_wishes_plus": {`},
		"unknown detail":    {`"max_charges": 3`, `"max_charges_all": 3`},
		"string for number": {`"max_charges": 3`, `"max_charges": "3"`},
		"fractional number": {`"max_charges": 3`, `"max_charges": 3.5`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items240JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items240JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems240(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items240JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems240(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items240JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems239(items240JSON); err == nil {
		t.Fatal("#239 decoder accepted the #240 batch")
	}
	if _, err := DecodeItems240(items239JSON); err == nil {
		t.Fatal("#240 decoder accepted the #239 batch")
	}
}
