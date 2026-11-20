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

func loadItems227(t *testing.T) ItemSourceBatch[ItemDetails227] {
	t.Helper()
	b, err := LoadItems227()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 214–215, in planning inventory order for issue #227. Every
// family is a Wondrous Item with one printed rarity and no base item.
func TestItems227SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		attunement string
		rarity     string
	}{
		{"Boots of Speed", "boots-of-speed", 214, "required", "rare"},
		{"Boots of Striding and Springing", "boots-of-striding-and-springing", 214, "required", "uncommon"},
		{"Boots of the Winterlands", "boots-of-the-winterlands", 214, "required", "uncommon"},
		{"Bowl of Commanding Water Elementals", "bowl-of-commanding-water-elementals", 214, "none", "rare"},
		{"Bracers of Archery", "bracers-of-archery", 214, "required", "uncommon"},
		{"Bracers of Defense", "bracers-of-defense", 214, "required", "rare"},
		{"Brazier of Commanding Fire Elementals", "brazier-of-commanding-fire-elementals", 214, "none", "rare"},
		{"Brooch of Shielding", "brooch-of-shielding", 214, "required", "uncommon"},
		{"Broom of Flying", "broom-of-flying", 214, "required", "uncommon"},
		{"Candle of Invocation", "candle-of-invocation", 214, "required", "very-rare"},
		{"Cape of the Mountebank", "cape-of-the-mountebank", 215, "none", "rare"},
		{"Carpet of Flying", "carpet-of-flying", 215, "none", "very-rare"},
	}
	b := loadItems227(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-03" || b.CatalogueIssue != 227 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #227")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 214–215") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	variants := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != "wondrous-item" || f.Base.Options == nil || len(f.Base.Options) != 0 || f.Base.Excluded == nil || len(f.Base.Excluded) != 0 {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			// One printed rarity per heading: the single explicit variant standard.
			if len(f.Variants) != 1 {
				t.Fatalf("source-listed variants mismatch: %d", len(f.Variants))
			}
			v := f.Variants[0]
			if v.ID != id+"/variant/standard" || v.Name != want.name || v.Rarity != want.rarity {
				t.Fatalf("variant identity mismatch: %s %s", v.ID, v.Rarity)
			}
			if v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
				t.Fatal("attunement or charges mismatch")
			}
			if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
				t.Fatal("no price is listed; cost must stay unspecified")
			}
			if v.MechanicsAvailable() {
				t.Fatal("catalogue record became available")
			}
			variants++
		})
	}
	if variants != 12 {
		t.Fatalf("expected twelve explicit variants, got %d", variants)
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
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 227 {
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
		t.Fatalf("inventory assigns %d families to issue #227", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "214, 215" || len(batch.Members) != 12 {
			t.Fatalf("planning batch pages or member count changed: %q, %d", batch.Pages, len(batch.Members))
		}
		for i, m := range batch.Members {
			if m.Name != expected[i].name || m.Page != expected[i].page {
				t.Fatalf("planning batch member %d is %s p.%d", i, m.Name, m.Page)
			}
		}
	}
	if planned != 1 {
		t.Fatalf("expected one planning batch %s, got %d", b.Batch, planned)
	}
	b.Families[0].Variants[0].Details.BootsSpeed.Condition = "changed"
	if fresh := loadItems227(t); fresh.Families[0].Variants[0].Details.BootsSpeed.Condition != "while-wearing" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all twelve variants, transcribed from pinned PDF
// pp. 214–215. Candle of Invocation starts on p. 214 and its Outer Plane table
// continues on p. 215. Table rows and d100 ranges are data here, never separate
// variants; a printed 00 is recorded as 100.
func TestItems227SourceValues(t *testing.T) {
	expected := []string{
		`{"boots_speed":{"condition":"while-wearing","activate_action":"bonus-action","activation":"click-heels-together","speed_multiplier":2,"opportunity_attacks_against_wearer":"disadvantage","end_trigger":"click-heels-together-again","total_use_seconds":600,"after_total_use":"ceases-to-function-for-wearer","recharge":"long-rest"}}`,
		`{"boots_striding":{"condition":"while-wearing","speed_becomes_feet":30,"speed_unless":"speed-is-higher","speed_not_reduced_by":["carrying-weight-over-carrying-capacity","wearing-heavy-armor"],"jump_frequency":"once-on-each-of-your-turns","jump_distance_feet":30,"jump_movement_cost_feet":10}}`,
		`{"boots_winterlands":{"appearance":"furred-boots","condition":"while-wearing","cold_resistance":{"benefit":"cold-resistance","resistance":"cold","tolerated_temperature_fahrenheit":0,"temperature_bound":"or-lower","without":"additional-protection"},"winter_strider":{"benefit":"winter-strider","ignores":"difficult-terrain","created_by":["ice","snow"]}}}`,
		`{"bowl":{"vessel_condition":"filled-with-water","summon":{"within_feet":5,"action":"magic","summons":"water-elemental","appears":"unoccupied-space-as-close-to-bowl-as-possible","elemental_traits":["understands-your-languages","obeys-your-commands","acts-immediately-after-you-on-your-initiative-count"],"duration_seconds":3600,"ends_when":["duration-elapsed","elemental-dies","dismissed-as-bonus-action"],"dismiss_action":"bonus-action","uses":1,"recharge":"next-dawn"},"dimensions_approximate":true,"diameter_feet":1,"depth":"half-of-diameter","capacity_gallons":3}}`,
		`{"bracers_archery":{"condition":"while-wearing","weapon_proficiencies":["longbow","shortbow"],"damage_bonus":2,"damage_bonus_weapons":["longbow","shortbow"]}}`,
		`{"bracers_defense":{"condition":"while-wearing","ac_bonus":2,"requires":["wearing-no-armor","using-no-shield"]}}`,
		`{"brazier":{"summon":{"within_feet":5,"action":"magic","summons":"fire-elemental","appears":"unoccupied-space-as-close-to-brazier-as-possible","elemental_traits":["understands-your-languages","obeys-your-commands","acts-immediately-after-you-on-your-initiative-count"],"duration_seconds":3600,"ends_when":["duration-elapsed","elemental-dies","dismissed-as-bonus-action"],"dismiss_action":"bonus-action","uses":1,"recharge":"next-dawn"}}}`,
		`{"brooch":{"condition":"while-wearing","resistances":["force"],"immunity_to_damage_from":"magic-missile-spell"}}`,
		`{"broom":{"appearance":"wooden-broom","inactive":"functions-like-mundane-broom","hover_requirement":"stand-astride","hover_action":"magic","fly_speed_feet":50,"capacity_pounds":400,"reduced_fly_speed_feet":30,"reduced_when_over_pounds":200,"hover_ends":["rider-lands","no-longer-riding"],"send_action":"magic","send_mode":"travels-alone","send_range_miles":1,"send_requirements":["name-the-location","familiar-with-location"],"recall_action":"magic","recall_requirements":["command-word","broom-still-within-range"],"recall_range_miles":1}}`,
		`{"candle":{"activation":"lit","activate_action":"magic","burn_seconds":14400,"after_burn":"destroyed","snuff_early":true,"burn_deduction_increment_seconds":60,"light":"dim","light_radius_feet":30,"attuned_in_light":"advantage-on-d20-tests","free_casters":["cleric","druid"],"free_cast_spell_level":1,"free_cast_requirement":"prepared","free_cast_cost":"no-spell-slot","gate_condition":"lit-for-first-time","gate_casts":"gate","gate_destroys_candle":true,"gate_links_to":"outer-plane","plane_selection":["gm-chooses","random-table"],"table_die_sides":100,"table_page":215,"outer_planes":[{"min":1,"max":5,"plane":"abyss"},{"min":6,"max":10,"plane":"acheron"},{"min":11,"max":17,"plane":"arborea"},{"min":18,"max":25,"plane":"arcadia"},{"min":26,"max":33,"plane":"beastlands"},{"min":34,"max":41,"plane":"bytopia"},{"min":42,"max":46,"plane":"carceri"},{"min":47,"max":54,"plane":"elysium"},{"min":55,"max":59,"plane":"gehenna"},{"min":60,"max":64,"plane":"hades"},{"min":65,"max":69,"plane":"limbo"},{"min":70,"max":77,"plane":"mechanus"},{"min":78,"max":85,"plane":"mount-celestia"},{"min":86,"max":90,"plane":"nine-hells"},{"min":91,"max":95,"plane":"pandemonium"},{"min":96,"max":100,"plane":"ysgard"}]}}`,
		`{"cape":{"scent":"brimstone","condition":"while-wearing","casts":"dimension-door","action":"magic","uses":1,"recharge":"next-dawn","on_teleport":"leaves-cloud-of-smoke","vacated_space":"lightly-obscured","obscured_until":"end-of-your-next-turn"}}`,
		`{"carpet":{"activate_action":"magic","activation_requirement":"command-word","effect":"hovers-and-flies","control":"moves-according-to-your-directions","control_range_feet":30,"size_count":4,"size_selection":["gm-chooses","random-table"],"table_die_sides":100,"sizes":[{"min":1,"max":20,"width_feet":3,"length_feet":5,"capacity_pounds":200,"fly_speed_feet":80},{"min":21,"max":55,"width_feet":4,"length_feet":6,"capacity_pounds":400,"fly_speed_feet":60},{"min":56,"max":80,"width_feet":5,"length_feet":7,"capacity_pounds":600,"fly_speed_feet":40},{"min":81,"max":100,"width_feet":6,"length_feet":9,"capacity_pounds":800,"fly_speed_feet":30}],"max_load_multiplier":2,"overloaded_fly_speed":"halved","overloaded_when":"carrying-more-than-normal-capacity"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:bonus-action:ten-minutes-total-until-long-rest:item.attunement item.wear action.bonus-action speed.doubled attack.opportunity-attacks-against-wearer-disadvantage duration.ends-on-repeat-activation item.ten-minutes-total-until-long-rest",
		"wear:none::item.attunement item.wear speed.becomes-30-feet-unless-higher speed.ignores-carrying-and-heavy-armor-reduction movement.jump-30-feet-for-10-feet-once-per-turn",
		"wear:none::item.attunement item.wear defense.resistance-cold environment.tolerates-0-fahrenheit-or-lower movement.ignores-ice-and-snow-difficult-terrain",
		"none:magic-action:bowl-filled-with-water within-5-feet-of-bowl once-until-next-dawn:action.magic summon.water-elemental summon.acts-after-summoner-initiative duration.one-hour-death-or-dismissal action.bonus-action-dismiss item.once-until-next-dawn",
		"wear:none::item.attunement item.wear proficiency.longbow-and-shortbow damage.longbow-and-shortbow-bonus-2",
		"wear:none::item.attunement item.wear defense.armor-class-bonus-2-without-armor-or-shield",
		"none:magic-action:within-5-feet-of-brazier once-until-next-dawn:action.magic summon.fire-elemental summon.acts-after-summoner-initiative duration.one-hour-death-or-dismissal action.bonus-action-dismiss item.once-until-next-dawn",
		"wear:none::item.attunement item.wear defense.resistance-force defense.immunity-magic-missile-damage",
		"none:magic-action:stand-astride-broom:item.attunement action.magic movement.ridden-fly-speed-50-feet item.carrying-capacity-400-pounds movement.fly-speed-30-feet-over-200-pounds item.hover-ends-on-landing-or-dismount item.send-to-named-familiar-location-within-1-mile item.recall-by-command-word-within-1-mile",
		"none:magic-action:four-hours-total-burn-time:item.attunement action.magic item.burn-time-four-hours light.dim-30-feet d20.advantage-while-in-light spell.cleric-or-druid-level-1-without-slot spell.cast-gate-on-first-lighting item.destroyed-by-gate-or-burn item.outer-plane-table",
		"wear:magic-action:once-until-next-dawn:item.wear action.magic spell.cast-dimension-door item.once-until-next-dawn area.vacated-space-lightly-obscured",
		"none:magic-action:command-word:action.magic item.command-word movement.hover-and-fly item.directed-within-30-feet item.carpet-size-table item.carrying-capacity-double-at-half-fly-speed",
	}
	index := 0
	for _, f := range loadItems227(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails227
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

// Earlier catalogues are the Healing contract fixture from #195 and magic items
// 01 from #225. Batches are composed in merge order, not page order.
func TestItems227CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems227(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	first, err := LoadItems225()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), first.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 25 || variants != 32 || composed[0].Name != "Potion of Healing" || composed[1].Name != "Adamantine Armor" || composed[12].Name != "Armor, +1, +2, or +3" || composed[13].Name != "Boots of Speed" || composed[24].Name != "Carpet of Flying" {
		t.Fatal("expected Potion of Healing, the twelve #225 families and the exact twelve #227 families with 32 variants")
	}
	earlier := 0
	for _, f := range first.Families {
		earlier += len(f.Variants)
	}
	if len(first.Families) != 12 || earlier != 16 {
		t.Fatal("magic items 01 changed")
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
			if strings.Contains(fixture, "items_227") {
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
func TestItems227RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items227JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) != 645 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items227JSON, &fresh); err != nil {
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
		if _, err := DecodeItems227(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems227RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails227]
	// Family indexes follow inventory order; see TestItems227SourceIdentities.
	details := func(b *batch, family int) *ItemDetails227 { return b.Families[family].Variants[0].Details }
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
		{"unknown family", func(b *batch) { b.Families[11].ID = "srd-5.2.1:content/items/censer-of-controlling-air-elementals" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Boots of Haste" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[5].Base.Options = []string{"any-light-armor"} }},
		{"invented exclusion", func(b *batch) { b.Families[5].Base.Excluded = []string{"shield"} }},
		{"source page", func(b *batch) { b.Families[10].Presentations[0].Pages = []int{214, 214} }},
		{"page range", func(b *batch) { b.Families[9].Presentations[0].Pages = []int{214, 215} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"duplicate presentation", func(b *batch) {
			b.Families[0].Presentations = append(b.Families[0].Presentations, b.Families[0].Presentations[0])
		}},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"missing variant", func(b *batch) { b.Families[11].Variants = nil }},
		{"invented size variant", func(b *batch) { b.Families[11].Variants = append(b.Families[11].Variants, b.Families[11].Variants[0]) }},
		{"variant slug", func(b *batch) {
			b.Families[11].Variants[0].ID = "srd-5.2.1:content/items/carpet-of-flying/variant/3-ft-by-5-ft"
		}},
		{"missing variant id", func(b *batch) { b.Families[0].Variants[0].ID = "" }},
		{"variant name", func(b *batch) { b.Families[0].Variants[0].Name = "Boots of Speed, +1" }},
		{"rarity", func(b *batch) { b.Families[8].Variants[0].Rarity = "rare" }},
		{"unknown rarity", func(b *batch) { b.Families[0].Variants[0].Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { b.Families[8].Variants[0].Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { b.Families[11].Variants[0].Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { b.Families[9].Variants[0].Attunement.Restrictions = []string{"cleric-or-druid"} }},
		{"missing charges", func(b *batch) { b.Families[0].Variants[0].Charges = nil }},
		{"invented charges", func(b *batch) { b.Families[0].Variants[0].Charges.Maximum = 10 }},
		{"recharge", func(b *batch) { b.Families[3].Variants[0].Charges.Recharge = "dawn" }},
		{"equip mode", func(b *batch) { b.Families[8].Variants[0].Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { b.Families[0].Variants[0].Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { b.Families[0].Variants[0].Use.Mode = "free-action" }},
		{"dropped use restriction", func(b *batch) {
			b.Families[3].Variants[0].Use.Restrictions = b.Families[3].Variants[0].Use.Restrictions[1:]
		}},
		{"listed price", func(b *batch) { b.Families[0].Variants[0].Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { b.Families[0].Variants[0].Cost.Amount = 4000 }},
		{"cost currency", func(b *batch) { b.Families[0].Variants[0].Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { b.Families[9].Variants[0].Cost.SourcePages = []int{215, 215} }},
		{"missing mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys = b.Families[0].Variants[0].RequiredMechanicKeys[:1]
		}},
		{"no mechanics", func(b *batch) { b.Families[0].Variants[0].RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { b.Families[5].Variants[0].RequiredMechanicKeys[2] = "item.grant-ac" }},
		{"duplicate mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys[1] = b.Families[0].Variants[0].RequiredMechanicKeys[0]
		}},
		{"missing details", func(b *batch) { b.Families[0].Variants[0].Details = nil }},
		{"empty details", func(b *batch) { b.Families[0].Variants[0].Details = &ItemDetails227{} }},
		{"mismatched payload", func(b *batch) { b.Families[3].Variants[0].Details = details(b, 6) }},
		{"second payload", func(b *batch) { details(b, 4).BracersDefense = details(b, 5).BracersDefense }},
		{"blank key", func(b *batch) { details(b, 0).BootsSpeed.Condition = " " }},
		{"prose value", func(b *batch) { details(b, 0).BootsSpeed.Activation = "Click the boots' heels together" }},
		{"speed multiplier", func(b *batch) { details(b, 0).BootsSpeed.SpeedMultiplier = 1 }},
		{"speed total use", func(b *batch) { details(b, 0).BootsSpeed.TotalUseSeconds = 0 }},
		{"striding speed", func(b *batch) { details(b, 1).BootsStriding.SpeedBecomesFeet = 0 }},
		{"striding jump cost", func(b *batch) { details(b, 1).BootsStriding.JumpMovementCostFeet = 30 }},
		{"striding reductions", func(b *batch) { details(b, 1).BootsStriding.SpeedNotReducedBy = []string{} }},
		{"winterlands negative temperature", func(b *batch) { details(b, 2).BootsWinterlands.ColdResistance.ToleratedTemperatureFahrenheit = -20 }},
		{"winterlands duplicate benefit", func(b *batch) { details(b, 2).BootsWinterlands.WinterStrider.Benefit = "cold-resistance" }},
		{"winterlands duplicate terrain", func(b *batch) { details(b, 2).BootsWinterlands.WinterStrider.CreatedBy[1] = "ice" }},
		{"bowl range", func(b *batch) { details(b, 3).Bowl.Summon.WithinFeet = 0 }},
		{"bowl uses", func(b *batch) { details(b, 3).Bowl.Summon.Uses = 0 }},
		{"bowl capacity", func(b *batch) { details(b, 3).Bowl.CapacityGallons = 0 }},
		{"bowl end triggers", func(b *batch) { details(b, 3).Bowl.Summon.EndsWhen = nil }},
		{"archery bonus", func(b *batch) { details(b, 4).BracersArchery.DamageBonus = 0 }},
		{"archery duplicate weapon", func(b *batch) { details(b, 4).BracersArchery.WeaponProficiencies[1] = "longbow" }},
		{"defense bonus", func(b *batch) { details(b, 5).BracersDefense.ACBonus = -2 }},
		{"defense requirements", func(b *batch) { details(b, 5).BracersDefense.Requires = []string{} }},
		{"brazier duration", func(b *batch) { details(b, 6).Brazier.Summon.DurationSeconds = 0 }},
		{"brazier duplicate trait", func(b *batch) { details(b, 6).Brazier.Summon.ElementalTraits[2] = "obeys-your-commands" }},
		{"brooch resistances", func(b *batch) { details(b, 7).Brooch.Resistances = []string{} }},
		{"broom reduced speed", func(b *batch) { details(b, 8).Broom.ReducedFlySpeedFeet = 50 }},
		{"broom reduced threshold", func(b *batch) { details(b, 8).Broom.ReducedWhenOverPounds = 400 }},
		{"broom range", func(b *batch) { details(b, 8).Broom.SendRangeMiles = 0 }},
		{"broom hover ends", func(b *batch) { details(b, 8).Broom.HoverEnds = nil }},
		{"candle table gap", func(b *batch) { details(b, 9).Candle.OuterPlanes[2].Min = 12 }},
		{"candle table short", func(b *batch) { details(b, 9).Candle.OuterPlanes = details(b, 9).Candle.OuterPlanes[:15] }},
		{"candle table end", func(b *batch) { details(b, 9).Candle.OuterPlanes[15].Max = 99 }},
		{"candle duplicate plane", func(b *batch) { details(b, 9).Candle.OuterPlanes[1].Plane = "abyss" }},
		{"candle die", func(b *batch) { details(b, 9).Candle.TableDieSides = 20 }},
		{"candle burn time", func(b *batch) { details(b, 9).Candle.BurnSeconds = 0 }},
		{"candle spell level", func(b *batch) { details(b, 9).Candle.FreeCastSpellLevel = 0 }},
		{"candle casters", func(b *batch) { details(b, 9).Candle.FreeCasters = []string{"cleric", "cleric"} }},
		{"cape uses", func(b *batch) { details(b, 10).Cape.Uses = 0 }},
		{"carpet table gap", func(b *batch) { details(b, 11).Carpet.Sizes[1].Min = 22 }},
		{"carpet missing size", func(b *batch) { details(b, 11).Carpet.Sizes = details(b, 11).Carpet.Sizes[:3] }},
		{"carpet size count", func(b *batch) { details(b, 11).Carpet.SizeCount = 5 }},
		{"carpet capacity", func(b *batch) { details(b, 11).Carpet.Sizes[0].CapacityPounds = 0 }},
		{"carpet negative speed", func(b *batch) { details(b, 11).Carpet.Sizes[3].FlySpeedFeet = -30 }},
		{"carpet load multiplier", func(b *batch) { details(b, 11).Carpet.MaxLoadMultiplier = 1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems227(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems227(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems227RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/boots-of-speed/variant/standard"`)
	if !bytes.Contains(items227JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":1,`, `"ac":99,`, `"speed":60,`, `"actions":["dash"],`, `"resources":{"minutes":10},`, `"current_charges":1,`, `"playable":true,`, `"available":true,`, `"properties":{"speed_multiplier":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"mixed-case-duplicate",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items227JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems227(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"boots_speed": {`, `"boots_speed_plus": {`},
		"unknown detail":     {`"speed_multiplier"`, `"speed_bonus"`},
		"string for number":  {`"speed_multiplier": 2`, `"speed_multiplier": "2"`},
		"fractional number":  {`"speed_multiplier": 2`, `"speed_multiplier": 2.5`},
		"number for boolean": {`"snuff_early": true`, `"snuff_early": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items227JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items227JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems227(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items227JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems227(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// The Healing contract decoder must not accept a batch, nor the reverse.
	if _, err := DecodeItemCatalogue(items227JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems227(itemFamiliesJSON); err == nil {
		t.Fatal("batch decoder accepted the contract fixture")
	}
}
