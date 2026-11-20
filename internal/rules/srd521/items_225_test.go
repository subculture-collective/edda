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

func loadItems225(t *testing.T) ItemSourceBatch[ItemDetails225] {
	t.Helper()
	b, err := LoadItems225()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 209–211, in planning inventory order for issue #225.
func TestItems225SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       string
		excluded   string
		attunement string
		variants   string // slug=rarity, in source order
	}{
		{"Adamantine Armor", "adamantine-armor", 209, "armor", "any-medium-armor any-heavy-armor", "hide-armor", "none", "standard=uncommon"},
		{"Ammunition of Slaying", "ammunition-of-slaying", 209, "weapon", "any-ammunition", "", "none", "standard=very-rare"},
		{"Ammunition, +1, +2, or +3", "ammunition-1-2-or-3", 209, "weapon", "any-ammunition", "", "none", "plus-1=uncommon plus-2=rare plus-3=very-rare"},
		{"Amulet of Health", "amulet-of-health", 209, "wondrous-item", "", "", "required", "standard=rare"},
		{"Amulet of Proof against Detection and Location", "amulet-of-proof-against-detection-and-location", 209, "wondrous-item", "", "", "required", "standard=uncommon"},
		{"Amulet of the Planes", "amulet-of-the-planes", 209, "wondrous-item", "", "", "required", "standard=very-rare"},
		{"Animated Shield", "animated-shield", 209, "armor", "shield", "", "required", "standard=very-rare"},
		{"Apparatus of the Crab", "apparatus-of-the-crab", 210, "wondrous-item", "", "", "none", "standard=legendary"},
		{"Armor of Invulnerability", "armor-of-invulnerability", 210, "armor", "plate-armor", "", "required", "standard=legendary"},
		{"Armor of Resistance", "armor-of-resistance", 210, "armor", "any-light-armor any-medium-armor any-heavy-armor", "", "required", "standard=rare"},
		{"Armor of Vulnerability", "armor-of-vulnerability", 211, "armor", "any-light-armor any-medium-armor any-heavy-armor", "", "required", "standard=rare"},
		{"Armor, +1, +2, or +3", "armor-1-2-or-3", 210, "armor", "any-light-armor any-medium-armor any-heavy-armor", "", "none", "plus-1=rare plus-2=very-rare plus-3=legendary"},
	}
	b := loadItems225(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-01" || b.CatalogueIssue != 225 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #225")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 209–211") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	variants := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || !slices.Equal(f.Base.Options, strings.Fields(want.base)) || !slices.Equal(f.Base.Excluded, strings.Fields(want.excluded)) {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			got := []string{}
			for _, v := range f.Variants {
				slug := strings.TrimPrefix(v.ID, id+"/variant/")
				got = append(got, slug+"="+v.Rarity)
				if v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
					t.Fatal("attunement or charges mismatch")
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
			}
			if !slices.Equal(got, strings.Fields(want.variants)) {
				t.Fatalf("source-listed variants mismatch: %v", got)
			}
			variants += len(f.Variants)
		})
	}
	if variants != 16 {
		t.Fatalf("expected sixteen explicit variants, got %d", variants)
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
		if r.Issue != 225 {
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
		t.Fatalf("inventory assigns %d families to issue #225", index)
	}
	b.Families[0].Variants[0].Details.Adamantine.Condition = "changed"
	if fresh := loadItems225(t); fresh.Families[0].Variants[0].Details.Adamantine.Condition != "while-wearing" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all sixteen variants, transcribed from pinned PDF
// pp. 209–211 (the lever table spans both columns of p. 210). Table rows and
// d100 ranges are data here, never separate variants.
func TestItems225SourceValues(t *testing.T) {
	expected := []string{
		`{"adamantine":{"condition":"while-wearing","critical_hit_against_wearer":"becomes-normal-hit"}}`,
		`{"slaying":{"type_selection":["gm-chooses","random-table"],"table_die_sides":100,"creature_types":[{"min":1,"max":10,"creature_type":"aberrations"},{"min":11,"max":15,"creature_type":"beasts"},{"min":16,"max":20,"creature_type":"celestials"},{"min":21,"max":25,"creature_type":"constructs"},{"min":26,"max":35,"creature_type":"dragons"},{"min":36,"max":45,"creature_type":"elementals"},{"min":46,"max":50,"creature_type":"humanoids"},{"min":51,"max":60,"creature_type":"fey"},{"min":61,"max":70,"creature_type":"fiends"},{"min":71,"max":75,"creature_type":"giants"},{"min":76,"max":80,"creature_type":"monstrosities"},{"min":81,"max":85,"creature_type":"oozes"},{"min":86,"max":90,"creature_type":"plants"},{"min":91,"max":100,"creature_type":"undead"}],"trigger":"creature-of-type-takes-damage-from-ammunition","save":"constitution","save_dc":17,"extra_damage":{"dice":6,"sides":10,"type":"force"},"save_failure":"full-extra-damage","save_success":"half-extra-damage","after_extra_damage":"becomes-nonmagical"}}`,
		`{"ammunition_bonus":{"attack_bonus":1,"damage_bonus":1,"bonus_determined_by":"rarity","after_hit":"no-longer-magical","typical_quantities":[10,20],"value_pieces":10,"value_equals":"potion-of-same-rarity"}}`,
		`{"ammunition_bonus":{"attack_bonus":2,"damage_bonus":2,"bonus_determined_by":"rarity","after_hit":"no-longer-magical","typical_quantities":[10,20],"value_pieces":10,"value_equals":"potion-of-same-rarity"}}`,
		`{"ammunition_bonus":{"attack_bonus":3,"damage_bonus":3,"bonus_determined_by":"rarity","after_hit":"no-longer-magical","typical_quantities":[10,20],"value_pieces":10,"value_equals":"potion-of-same-rarity"}}`,
		`{"health":{"condition":"while-wearing","ability":"constitution","score":19,"no_effect_if_score_at_least":19}}`,
		`{"proof":{"condition":"while-wearing","protections":["cannot-be-targeted-by-divination-spells","cannot-be-perceived-through-magical-scrying-sensors"],"exception":"wearer-allows-it"}}`,
		`{"planes":{"condition":"while-wearing","action":"magic","named_location":"familiar-location-on-another-plane","check":{"ability":"intelligence","skill":"arcana","dc":15},"success_casts":"plane-shift","failure_affects":["wearer","each-creature-within-radius","each-object-within-radius"],"failure_radius_feet":15,"table_die_sides":100,"destinations":[{"min":1,"max":60,"destination":"random-location-on-named-plane","sub_die_sides":0,"planes":[]},{"min":61,"max":70,"destination":"random-location-on-inner-plane","sub_die_sides":6,"planes":[{"roll":1,"plane":"plane-of-air"},{"roll":2,"plane":"plane-of-earth"},{"roll":3,"plane":"plane-of-fire"},{"roll":4,"plane":"plane-of-water"},{"roll":5,"plane":"feywild"},{"roll":6,"plane":"shadowfell"}]},{"min":71,"max":80,"destination":"random-location-on-outer-plane","sub_die_sides":8,"planes":[{"roll":1,"plane":"arborea"},{"roll":2,"plane":"arcadia"},{"roll":3,"plane":"beastlands"},{"roll":4,"plane":"bytopia"},{"roll":5,"plane":"elysium"},{"roll":6,"plane":"mechanus"},{"roll":7,"plane":"mount-celestia"},{"roll":8,"plane":"ysgard"}]},{"min":81,"max":90,"destination":"random-location-on-outer-plane","sub_die_sides":8,"planes":[{"roll":1,"plane":"abyss"},{"roll":2,"plane":"acheron"},{"roll":3,"plane":"carceri"},{"roll":4,"plane":"gehenna"},{"roll":5,"plane":"hades"},{"roll":6,"plane":"limbo"},{"roll":7,"plane":"nine-hells"},{"roll":8,"plane":"pandemonium"}]},{"min":91,"max":100,"destination":"random-location-on-astral-plane","sub_die_sides":0,"planes":[]}]}}`,
		`{"animated":{"condition":"while-holding","activate_action":"bonus-action","effect":"hovers-in-wielder-space-protecting-as-if-wielded","hands_free":true,"duration_seconds":60,"end_triggers":["duration-elapsed","bonus-action-to-end","wielder-dies","wielder-incapacitated"],"on_end":["falls-to-ground","into-free-hand"]}}`,
		`{"apparatus":{"appearance":"sealed-iron-barrel","weight_pounds":500,"catch_check":{"ability":"intelligence","skill":"investigation","dc":20},"occupants":2,"occupant_max_size":"medium","lever_count":10,"lever_positions":["up","neutral","down"],"transformed_appearance":"giant-lobster","size":"large","ac":20,"hp":200,"speed_feet":30,"swim_feet":30,"speed_without_legs_feet":0,"immunities":["poison","psychic"],"pilots":1,"closed_hatch":["airtight","watertight"],"air_hours":10,"air_divided_by":"breathing-creatures-inside","floats":true,"max_depth_feet":900,"pressure_damage":{"dice":2,"sides":6,"type":"bludgeoning"},"pressure_interval_seconds":60,"operate_action":"utilize","operator":"creature-in-compartment","max_levers_per_action":2,"lever_reset":"neutral-after-each-use","lever_order":"left-to-right","levers":[{"lever":1,"up":"legs-extend-allowing-walk-and-swim","down":"legs-retract-speed-and-swim-speed-0-no-speed-bonuses"},{"lever":2,"up":"forward-window-shutter-opens","down":"forward-window-shutter-closes"},{"lever":3,"up":"side-window-shutters-open-two-per-side","down":"side-window-shutters-close-two-per-side"},{"lever":4,"up":"two-claws-extend-from-front","down":"claws-retract"},{"lever":5,"up":"each-extended-claw-melee-attack-bludgeoning","down":"each-extended-claw-melee-attack-grapple"},{"lever":6,"up":"walks-or-swims-forward-if-legs-extended","down":"walks-or-swims-backward-if-legs-extended"},{"lever":7,"up":"turns-90-degrees-counterclockwise-if-legs-extended","down":"turns-90-degrees-clockwise-if-legs-extended"},{"lever":8,"up":"eyelike-fixtures-emit-light","down":"light-turns-off"},{"lever":9,"up":"sinks-up-to-20-feet-if-in-liquid","down":"rises-up-to-20-feet-if-in-liquid"},{"lever":10,"up":"rear-hatch-unseals-and-opens","down":"rear-hatch-closes-and-seals"}],"claws":2,"claw_attack_bonus":8,"claw_reach_feet":5,"claw_up_average":7,"claw_up_damage":{"dice":2,"sides":6,"type":"bludgeoning"},"claw_down_condition":"grappled","claw_escape_dc":15,"bright_light_feet":30,"dim_light_additional_feet":30,"vertical_feet":20,"turn_degrees":90,"side_shutters_per_side":2}}`,
		`{"invulnerability":{"condition":"while-wearing","resistances":["bludgeoning","piercing","slashing"],"property":"metal-shell","property_action":"magic","property_immunities":["bludgeoning","piercing","slashing"],"property_duration_seconds":600,"property_ends_early":"no-longer-wearing-armor","property_uses":1,"property_recharge":"next-dawn"}}`,
		`{"resistance":{"condition":"while-wearing","resistance_count":1,"type_selection":["gm-chooses","random-table"],"table_die_sides":10,"damage_types":[{"roll":1,"damage_type":"acid"},{"roll":2,"damage_type":"cold"},{"roll":3,"damage_type":"fire"},{"roll":4,"damage_type":"force"},{"roll":5,"damage_type":"lightning"},{"roll":6,"damage_type":"necrotic"},{"roll":7,"damage_type":"poison"},{"roll":8,"damage_type":"psychic"},{"roll":9,"damage_type":"radiant"},{"roll":10,"damage_type":"thunder"}]}}`,
		`{"vulnerability":{"condition":"while-wearing","resistance_count":1,"damage_types":["bludgeoning","piercing","slashing"],"type_selection":["gm-chooses","random"],"cursed":true,"curse_revealed_by":["identify-spell-cast-on-armor","attuning"],"curse_trigger":"attuning","curse_ends_by":["targeted-by-remove-curse-spell","similar-magic"],"removing_armor_ends_curse":false,"vulnerability_count":2,"vulnerability_types":"associated-types-other-than-resisted"}}`,
		`{"armor_bonus":{"condition":"while-wearing","ac_bonus":1,"bonus_determined_by":"rarity"}}`,
		`{"armor_bonus":{"condition":"while-wearing","ac_bonus":2,"bonus_determined_by":"rarity"}}`,
		`{"armor_bonus":{"condition":"while-wearing","ac_bonus":3,"bonus_determined_by":"rarity"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:none::item.wear-armor defense.critical-hit-becomes-normal-hit",
		"none:ammunition-attack:creature-type-determined-by-gm-or-table:item.magic-ammunition item.creature-type-table save.constitution-dc-17 damage.extra-force-half-on-save item.nonmagical-after-extra-damage",
		"none:ammunition-attack::item.magic-ammunition attack.ammunition-bonus damage.ammunition-bonus item.no-longer-magical-after-hit",
		"none:ammunition-attack::item.magic-ammunition attack.ammunition-bonus damage.ammunition-bonus item.no-longer-magical-after-hit",
		"none:ammunition-attack::item.magic-ammunition attack.ammunition-bonus damage.ammunition-bonus item.no-longer-magical-after-hit",
		"wear:none::item.attunement item.wear ability.constitution-score-19",
		"wear:none::item.attunement item.wear protection.untargetable-by-divination-spells protection.imperceptible-to-scrying-sensors protection.wearer-may-allow",
		"wear:magic-action:name-familiar-location-on-another-plane:item.attunement item.wear action.magic check.intelligence-arcana-dc-15 spell.cast-plane-shift transport.random-planar-destination-table area.creatures-and-objects-within-15-feet",
		"hold:bonus-action::item.attunement item.hold-shield action.bonus-action defense.shield-without-wielding duration.one-minute duration.ends-on-bonus-action-death-or-incapacitated",
		"none:utilize-action:creature-in-compartment requires-one-pilot up-to-two-levers-per-action:object.statistics vehicle.pilot vehicle.air-supply vehicle.depth-pressure-damage check.intelligence-investigation-dc-20 action.utilize item.lever-table attack.claw-melee condition.grappled-escape-dc-15 light.bright-and-dim",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear-armor defense.resistance-bludgeoning-piercing-slashing action.magic defense.immunity-bludgeoning-piercing-slashing duration.ten-minutes-or-armor-removed item.once-until-next-dawn",
		"wear:none::item.attunement item.wear-armor item.damage-type-table defense.resistance-chosen-damage-type",
		"wear:none::item.attunement item.wear-armor defense.resistance-chosen-damage-type item.curse item.curse-revealed-by-identify-or-attunement item.curse-ends-on-remove-curse defense.vulnerability-other-two-damage-types",
		"wear:none::item.wear-armor defense.armor-class-bonus",
		"wear:none::item.wear-armor defense.armor-class-bonus",
		"wear:none::item.wear-armor defense.armor-class-bonus",
	}
	index := 0
	for _, f := range loadItems225(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails225
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

func TestItems225CompositionWithContractFixture(t *testing.T) {
	batch := loadItems225(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 13 || variants != 20 || composed[0].Name != "Potion of Healing" || composed[12].Name != "Armor, +1, +2, or +3" {
		t.Fatal("expected Potion of Healing plus the exact twelve families and sixteen variants")
	}
	if len(contract.Families) != 1 || len(contract.Families[0].Variants) != 4 {
		t.Fatal("the Healing contract fixture changed")
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	if _, err := ComposeItemSourceIdentities(); err == nil {
		t.Fatal("empty composition accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_225") {
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
func TestItems225RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items225JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 700 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items225JSON, &fresh); err != nil {
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
		if _, err := DecodeItems225(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems225RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails225]
	// Family indexes follow inventory order; see TestItems225SourceIdentities.
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-02" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 226 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) { b.Families = append(b.Families, b.Families[0]) }},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/arrow-catching-shield" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Adamantine Armour" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) { b.Families[0].Base.Options = append(b.Families[0].Base.Options, "any-light-armor") }},
		{"dropped exclusion", func(b *batch) { b.Families[0].Base.Excluded = []string{} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{210, 210} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{209, 211} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"duplicate presentation", func(b *batch) {
			b.Families[0].Presentations = append(b.Families[0].Presentations, b.Families[0].Presentations[0])
		}},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"missing rarity variant", func(b *batch) { b.Families[2].Variants = b.Families[2].Variants[:2] }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"duplicate variant id", func(b *batch) { b.Families[2].Variants[1].ID = b.Families[2].Variants[0].ID }},
		{"missing variant id", func(b *batch) { b.Families[0].Variants[0].ID = "" }},
		{"variant name", func(b *batch) { b.Families[11].Variants[0].Name = "Armor, +4" }},
		{"rarity", func(b *batch) { b.Families[11].Variants[0].Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { b.Families[0].Variants[0].Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { b.Families[3].Variants[0].Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { b.Families[0].Variants[0].Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { b.Families[3].Variants[0].Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { b.Families[0].Variants[0].Charges = nil }},
		{"invented charges", func(b *batch) { b.Families[8].Variants[0].Charges.Maximum = 1 }},
		{"recharge", func(b *batch) { b.Families[8].Variants[0].Charges.Recharge = "dawn" }},
		{"equip mode", func(b *batch) { b.Families[6].Variants[0].Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { b.Families[0].Variants[0].Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { b.Families[6].Variants[0].Use.Mode = "free-action" }},
		{"dropped use restriction", func(b *batch) { b.Families[8].Variants[0].Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { b.Families[0].Variants[0].Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { b.Families[0].Variants[0].Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { b.Families[0].Variants[0].Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { b.Families[0].Variants[0].Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys = b.Families[0].Variants[0].RequiredMechanicKeys[:1]
		}},
		{"no mechanics", func(b *batch) { b.Families[0].Variants[0].RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { b.Families[0].Variants[0].RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"duplicate mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys[1] = b.Families[0].Variants[0].RequiredMechanicKeys[0]
		}},
		{"missing details", func(b *batch) { b.Families[0].Variants[0].Details = nil }},
		{"empty details", func(b *batch) { b.Families[0].Variants[0].Details = &ItemDetails225{} }},
		{"mismatched payload", func(b *batch) { b.Families[0].Variants[0].Details = b.Families[11].Variants[0].Details }},
		{"second payload", func(b *batch) {
			b.Families[0].Variants[0].Details.ArmorBonus = b.Families[11].Variants[0].Details.ArmorBonus
		}},
		{"blank key", func(b *batch) { b.Families[0].Variants[0].Details.Adamantine.Condition = " " }},
		{"prose value", func(b *batch) { b.Families[0].Variants[0].Details.Adamantine.Condition = "While you're wearing it" }},
		{"slaying table gap", func(b *batch) { b.Families[1].Variants[0].Details.Slaying.CreatureTypes[1].Min = 12 }},
		{"slaying table short", func(b *batch) {
			b.Families[1].Variants[0].Details.Slaying.CreatureTypes = b.Families[1].Variants[0].Details.Slaying.CreatureTypes[:13]
		}},
		{"slaying duplicate type", func(b *batch) {
			b.Families[1].Variants[0].Details.Slaying.CreatureTypes[1].CreatureType = "aberrations"
		}},
		{"slaying dc", func(b *batch) { b.Families[1].Variants[0].Details.Slaying.SaveDC = 0 }},
		{"slaying dice", func(b *batch) { b.Families[1].Variants[0].Details.Slaying.ExtraDamage.Dice = -6 }},
		{"ammunition bonus", func(b *batch) { b.Families[2].Variants[0].Details.AmmunitionBonus.AttackBonus = 3 }},
		{"ammunition damage bonus", func(b *batch) { b.Families[2].Variants[2].Details.AmmunitionBonus.DamageBonus = 1 }},
		{"ammunition quantities", func(b *batch) { b.Families[2].Variants[0].Details.AmmunitionBonus.TypicalQuantities = []int{} }},
		{"health score", func(b *batch) { b.Families[3].Variants[0].Details.Health.Score = 31 }},
		{"health threshold", func(b *batch) { b.Families[3].Variants[0].Details.Health.NoEffectIfScoreAtLeast = 20 }},
		{"proof protections", func(b *batch) { b.Families[4].Variants[0].Details.Proof.Protections = nil }},
		{"planes table gap", func(b *batch) { b.Families[5].Variants[0].Details.Planes.Destinations[4].Max = 99 }},
		{"planes sub-roll", func(b *batch) {
			b.Families[5].Variants[0].Details.Planes.Destinations[1].Planes = b.Families[5].Variants[0].Details.Planes.Destinations[1].Planes[:5]
		}},
		{"planes invented sub-roll", func(b *batch) { b.Families[5].Variants[0].Details.Planes.Destinations[0].SubDieSides = 6 }},
		{"planes check", func(b *batch) { b.Families[5].Variants[0].Details.Planes.Check.DC = 0 }},
		{"shield duration", func(b *batch) { b.Families[6].Variants[0].Details.Animated.DurationSeconds = 0 }},
		{"shield end triggers", func(b *batch) { b.Families[6].Variants[0].Details.Animated.EndTriggers = []string{} }},
		{"apparatus lever missing", func(b *batch) {
			b.Families[7].Variants[0].Details.Apparatus.Levers = b.Families[7].Variants[0].Details.Apparatus.Levers[:9]
		}},
		{"apparatus lever order", func(b *batch) { b.Families[7].Variants[0].Details.Apparatus.Levers[4].Lever = 6 }},
		{"apparatus hit points", func(b *batch) { b.Families[7].Variants[0].Details.Apparatus.HP = 0 }},
		{"apparatus negative speed", func(b *batch) { b.Families[7].Variants[0].Details.Apparatus.SpeedFeet = -30 }},
		{"invulnerability uses", func(b *batch) { b.Families[8].Variants[0].Details.Invulnerability.PropertyUses = 0 }},
		{"invulnerability duplicate type", func(b *batch) { b.Families[8].Variants[0].Details.Invulnerability.Resistances[1] = "bludgeoning" }},
		{"resistance table short", func(b *batch) {
			b.Families[9].Variants[0].Details.Resistance.DamageTypes = b.Families[9].Variants[0].Details.Resistance.DamageTypes[:9]
		}},
		{"resistance duplicate type", func(b *batch) { b.Families[9].Variants[0].Details.Resistance.DamageTypes[9].DamageType = "acid" }},
		{"resistance count", func(b *batch) { b.Families[9].Variants[0].Details.Resistance.ResistanceCount = 2 }},
		{"vulnerability count", func(b *batch) { b.Families[10].Variants[0].Details.Vulnerability.VulnerabilityCount = 3 }},
		{"vulnerability uncursed", func(b *batch) { b.Families[10].Variants[0].Details.Vulnerability.Cursed = false }},
		{"armor bonus", func(b *batch) { b.Families[11].Variants[1].Details.ArmorBonus.ACBonus = 3 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems225(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems225(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems225RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/adamantine-armor/variant/standard"`)
	if !bytes.Contains(items225JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items225JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems225(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"adamantine": {`, `"adamantine_plus": {`},
		"unknown detail":     {`"critical_hit_against_wearer"`, `"critical_hit_immunity"`},
		"string for number":  {`"save_dc": 17`, `"save_dc": "17"`},
		"fractional number":  {`"save_dc": 17`, `"save_dc": 17.5`},
		"number for boolean": {`"hands_free": true`, `"hands_free": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items225JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items225JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems225(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items225JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems225(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// The Healing contract decoder must not accept a batch, nor the reverse.
	if _, err := DecodeItemCatalogue(items225JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems225(itemFamiliesJSON); err == nil {
		t.Fatal("batch decoder accepted the contract fixture")
	}
}
