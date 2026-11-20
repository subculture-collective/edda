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

func loadItems228(t *testing.T) ItemSourceBatch[ItemDetails228] {
	t.Helper()
	b, err := LoadItems228()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 215-216, in planning inventory order for issue #228.
func TestItems228SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		attunement string
		rarity     string
	}{
		{"Censer of Controlling Air Elementals", "censer-of-controlling-air-elementals", 215, "none", "rare"},
		{"Chime of Opening", "chime-of-opening", 215, "none", "rare"},
		{"Circlet of Blasting", "circlet-of-blasting", 215, "none", "uncommon"},
		{"Cloak of Arachnida", "cloak-of-arachnida", 215, "required", "very-rare"},
		{"Cloak of Displacement", "cloak-of-displacement", 215, "required", "rare"},
		{"Cloak of Elvenkind", "cloak-of-elvenkind", 215, "required", "uncommon"},
		{"Cloak of Invisibility", "cloak-of-invisibility", 215, "required", "legendary"},
		{"Cloak of Protection", "cloak-of-protection", 216, "required", "uncommon"},
		{"Cloak of the Bat", "cloak-of-the-bat", 216, "required", "rare"},
		{"Cloak of the Manta Ray", "cloak-of-the-manta-ray", 216, "required", "uncommon"},
		{"Crystal Ball", "crystal-ball", 216, "required", "very-rare"},
		{"Crystal Ball of Mind Reading", "crystal-ball-of-mind-reading", 216, "required", "legendary"},
	}
	b := loadItems228(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-04" || b.CatalogueIssue != 228 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #228")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 215–216") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != "wondrous-item" || len(f.Base.Options) != 0 || len(f.Base.Excluded) != 0 {
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
			// Only the Cloak of Invisibility lists charges (3, regaining 1d3 at dawn).
			wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
			if want.slug == "cloak-of-invisibility" {
				wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}
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
		if r.Issue != 228 {
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
		t.Fatalf("inventory assigns %d families to issue #228", index)
	}
	b.Families[0].Variants[0].Details.Censer.Condition = "changed"
	if fresh := loadItems228(t); fresh.Families[0].Variants[0].Details.Censer.Condition != "gently-swinging" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all twelve variants, transcribed from pinned PDF
// pp. 215-216 (the Cloak of Invisibility text continues across the page break).
func TestItems228SourceValues(t *testing.T) {
	expected := []string{
		`{"censer":{"condition":"gently-swinging","action":"magic","summons":"air-elemental","location":"unoccupied-space-nearest-censer","understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-summoner","duration_seconds":3600,"end_triggers":["duration-elapsed","dies","dismissed-as-bonus-action"],"uses":1,"recharge":"next-dawn"}}`,
		`{"chime":{"appearance":"hollow-metal-tube","length_feet":1,"weight_pounds":1,"action":"magic","spell":"knock","sound_replaces":"customary-knocking-sound","audible_feet":300,"uses":10,"after_final_use":"cracks-and-becomes-useless"}}`,
		`{"circlet":{"condition":"while-wearing","spell":"scorching-ray","attack_bonus":5,"uses":1,"recharge":"next-dawn"}}`,
		`{"arachnida":{"condition":"while-wearing","resistances":["poison"],"climb_speed":"equal-to-speed","climb_directions":["up","down","across-vertical-surfaces","along-ceilings"],"climb_hands_free":true,"web_catch":"cannot-be-caught-in-webs","web_movement":"as-difficult-terrain","web_spell":"web","web_save_dc":13,"web_area_multiplier":2,"web_uses":1,"web_recharge":"next-dawn"}}`,
		`{"displacement":{"condition":"while-wearing","effect":"illusion-appears-near-actual-location","attack_rolls_against_wearer":"disadvantage","suspended_by":"taking-damage","resumes_at":"start-of-next-turn","suppressed_while":"speed-0"}}`,
		`{"elvenkind":{"condition":"while-wearing","perceive_wearer_ability":"wisdom","perceive_wearer_skill":"perception","perceive_wearer_mode":"disadvantage","stealth_ability":"dexterity","stealth_skill":"stealth","stealth_mode":"advantage"}}`,
		`{"invisibility":{"condition":"while-wearing","action":"magic","gesture":"pull-hood-over-head","charge_cost":1,"grants_condition":"invisible","duration_seconds":3600,"ends_early":["hood-pulled-down","cloak-no-longer-worn"],"hood_down_action":"none","max_charges":3,"regain_dice":{"dice":1,"sides":3,"type":"charges"},"regain_timing":"dawn"}}`,
		`{"protection":{"condition":"while-wearing","ac_bonus":1,"saving_throw_bonus":1}}`,
		`{"bat":{"condition":"while-wearing","stealth_ability":"dexterity","stealth_skill":"stealth","stealth_mode":"advantage","fly_speed_feet":40,"fly_light":["dim-light","darkness"],"fly_grip":"grip-edges-of-cloak","fly_ends_on":["fail-to-grip-edges-while-flying","no-longer-in-dim-light-or-darkness"],"polymorph_spell":"polymorph","polymorph_target":"self","polymorph_form":"bat","polymorph_light":["dim-light","darkness"],"retained_scores":["intelligence","wisdom","charisma"],"polymorph_uses":1,"polymorph_recharge":"next-dawn"}}`,
		`{"manta_ray":{"condition":"while-wearing","breathe_underwater":true,"swim_speed_feet":60}}`,
		`{"crystal_ball":{"condition":"while-touching","spell":"scrying","save_dc":17}}`,
		`{"mind_reading":{"condition":"while-touching","spell":"scrying","save_dc":17,"extra_spell":"detect-thoughts","extra_save_dc":17,"extra_targets":"creatures-you-can-see-within-range-of-sensor","extra_range_feet":30,"extra_range_from":"spell-sensor","extra_concentration_required":false,"extra_ends_with":"scrying-ends"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:magic-action:swinging-censer once-until-next-dawn:action.magic summon.air-elemental duration.one-hour item.once-until-next-dawn",
		"none:magic-action:ten-uses:action.magic spell.cast-knock item.limited-uses item.useless-after-final-use",
		"wear:cast-spell:once-until-next-dawn:item.wear spell.cast-scorching-ray spell.attack-bonus-5 item.once-until-next-dawn",
		"wear:cast-spell:once-until-next-dawn:item.attunement item.wear defense.resistance-poison movement.climb-speed-equals-speed movement.web-immunity-and-passage spell.cast-web spell.web-save-dc-13 item.once-until-next-dawn",
		"wear:none::item.attunement item.wear defense.attackers-have-disadvantage effect.suspended-after-damage effect.suppressed-at-speed-0",
		"wear:none::item.attunement item.wear check.perception-to-perceive-wearer-disadvantage check.stealth-advantage",
		"wear:magic-action:expend-one-charge:item.attunement item.wear item.charges action.magic condition.invisible duration.one-hour effect.ends-when-hood-down-or-cloak-removed",
		"wear:none::item.attunement item.wear defense.armor-class-bonus-1 save.bonus-1",
		"wear:cast-spell:dim-light-or-darkness once-until-next-dawn:item.attunement item.wear check.stealth-advantage movement.fly-speed-40-in-dim-light-or-darkness spell.cast-polymorph-self-bat item.once-until-next-dawn",
		"wear:none::item.attunement item.wear movement.breathe-underwater movement.swim-speed-60",
		"touch:cast-spell:touching-crystal-orb:item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17",
		"touch:cast-spell:touching-crystal-orb:item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 spell.cast-detect-thoughts spell.detect-thoughts-save-dc-17 concentration.detect-thoughts-ends-with-scrying",
	}
	index := 0
	for _, f := range loadItems228(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails228
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

func TestItems228CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems228(t)
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
	if len(composed) != 25 || variants != 32 || composed[0].Name != "Potion of Healing" || composed[13].Name != "Censer of Controlling Air Elementals" || composed[24].Name != "Crystal Ball of Mind Reading" {
		t.Fatal("expected Potion of Healing, the #225 families and exactly the twelve families of #228")
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
			if strings.Contains(fixture, "items_228") {
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
func TestItems228RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items228JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 500 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items228JSON, &fresh); err != nil {
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
		if _, err := DecodeItems228(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems228RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails228]
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails228] { return &b.Families[family].Variants[0] }
	// Family indexes follow inventory order; see TestItems228SourceIdentities.
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-03" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 227 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/crystal-ball-of-telepathy"
			b.Families[12].Name = "Crystal Ball of Telepathy"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/cape-of-the-mountebank" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Censer of Air Elementals" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{216, 216} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{215, 216} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Censer" }},
		{"rarity", func(b *batch) { v(b, 6).Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, 3).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, 0).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, 3).Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { v(b, 6).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, 6).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"charge maximum", func(b *batch) { v(b, 6).Charges.Maximum = 4 }},
		{"charge per use", func(b *batch) { v(b, 6).Charges.PerUse = 2 }},
		{"charge recharge", func(b *batch) { v(b, 6).Charges.Recharge = "none" }},
		{"invented charges", func(b *batch) { v(b, 1).Charges.Maximum = 10 }},
		{"chime is not a charge pool", func(b *batch) {
			*v(b, 1).Charges = ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "none"}
		}},
		{"equip mode", func(b *batch) { v(b, 2).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, 2).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0).Use.Mode = "bonus-action" }},
		{"dropped use restriction", func(b *batch) { v(b, 0).Use.Restrictions = []string{"swinging-censer"} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails228{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.Chime = v(b, 1).Details.Chime }},
		{"blank key", func(b *batch) { v(b, 0).Details.Censer.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0).Details.Censer.Condition = "While gently swinging" }},
		{"censer duration", func(b *batch) { v(b, 0).Details.Censer.DurationSeconds = 0 }},
		{"censer end triggers", func(b *batch) { v(b, 0).Details.Censer.EndTriggers = []string{} }},
		{"censer duplicate trigger", func(b *batch) {
			v(b, 0).Details.Censer.EndTriggers[1] = "dies"
			v(b, 0).Details.Censer.EndTriggers[2] = "dies"
		}},
		{"chime uses", func(b *batch) { v(b, 1).Details.Chime.Uses = 0 }},
		{"chime range", func(b *batch) { v(b, 1).Details.Chime.AudibleFeet = 0 }},
		{"circlet bonus", func(b *batch) { v(b, 2).Details.Circlet.AttackBonus = 0 }},
		{"arachnida dc", func(b *batch) { v(b, 3).Details.Arachnida.WebSaveDC = 0 }},
		{"arachnida duplicate direction", func(b *batch) { v(b, 3).Details.Arachnida.ClimbDirections[1] = "up" }},
		{"arachnida negative area", func(b *batch) { v(b, 3).Details.Arachnida.WebAreaMultiplier = -2 }},
		{"invisibility cost", func(b *batch) { v(b, 6).Details.Invisibility.ChargeCost = 4 }},
		{"invisibility dice", func(b *batch) { v(b, 6).Details.Invisibility.RegainDice.Sides = 0 }},
		{"invisibility max charges", func(b *batch) { v(b, 6).Details.Invisibility.MaxCharges = 5 }},
		{"invisibility regain 2d6", func(b *batch) {
			d := &v(b, 6).Details.Invisibility.RegainDice
			d.Dice, d.Sides = 2, 6
		}},
		{"invisibility end triggers", func(b *batch) { v(b, 6).Details.Invisibility.EndsEarly = []string{} }},
		{"protection bonus", func(b *batch) { v(b, 7).Details.Protection.ACBonus = 0 }},
		{"bat speed", func(b *batch) { v(b, 8).Details.Bat.FlySpeedFeet = 0 }},
		{"bat retained scores", func(b *batch) { v(b, 8).Details.Bat.RetainedScores = nil }},
		{"bat duplicate light", func(b *batch) { v(b, 8).Details.Bat.FlyLight[1] = "dim-light" }},
		{"manta ray speed", func(b *batch) { v(b, 9).Details.MantaRay.SwimSpeedFeet = 0 }},
		{"crystal ball dc", func(b *batch) { v(b, 10).Details.CrystalBall.SaveDC = 0 }},
		{"mind reading range", func(b *batch) { v(b, 11).Details.MindReading.ExtraRangeFeet = 0 }},
		{"mind reading dc", func(b *batch) { v(b, 11).Details.MindReading.ExtraSaveDC = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems228(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems228(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems228RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/censer-of-controlling-air-elementals/variant/standard"`)
	if !bytes.Contains(items228JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items228JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems228(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"censer": {`, `"censer_plus": {`},
		"unknown detail":     {`"summons"`, `"summons_all"`},
		"string for number":  {`"save_dc": 17`, `"save_dc": "17"`},
		"fractional number":  {`"save_dc": 17`, `"save_dc": 17.5`},
		"number for boolean": {`"climb_hands_free": true`, `"climb_hands_free": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items228JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items228JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems228(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items228JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems228(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor the #225 decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items228JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems225(items228JSON); err == nil {
		t.Fatal("#225 decoder accepted the #228 batch")
	}
	if _, err := DecodeItems228(items225JSON); err == nil {
		t.Fatal("#228 decoder accepted the #225 batch")
	}
}
