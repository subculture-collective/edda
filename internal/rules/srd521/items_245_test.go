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

func loadItems245(t *testing.T) ItemSourceBatch[ItemDetails245] {
	t.Helper()
	b, err := LoadItems245()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 250-252, in planning inventory order for issue #245.
func TestItems245SourceIdentities(t *testing.T) {
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
		{"Wand of Enemy Detection", "wand-of-enemy-detection", 250, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Enemy Detection", "rare"}}},
		{"Wand of Fear", "wand-of-fear", 250, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Fear", "rare"}}},
		{"Wand of Fireballs", "wand-of-fireballs", 250, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Fireballs", "rare"}}},
		{"Wand of Lightning Bolts", "wand-of-lightning-bolts", 251, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Lightning Bolts", "rare"}}},
		{"Wand of Magic Detection", "wand-of-magic-detection", 251, "wand", nil, "none", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Wand of Magic Detection", "uncommon"}}},
		{"Wand of Magic Missiles", "wand-of-magic-missiles", 251, "wand", nil, "none", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Magic Missiles", "uncommon"}}},
		{"Wand of Paralysis", "wand-of-paralysis", 251, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Paralysis", "rare"}}},
		{"Wand of Polymorph", "wand-of-polymorph", 251, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Polymorph", "very-rare"}}},
		{"Wand of Secrets", "wand-of-secrets", 251, "wand", nil, "none", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Wand of Secrets", "uncommon"}}},
		{"Wand of Web", "wand-of-web", 251, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Web", "uncommon"}}},
		{"Wand of Wonder", "wand-of-wonder", 251, "wand", nil, "required", ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"}, []variant{{"standard", "Wand of Wonder", "rare"}}},
		{"Wand of the War Mage, +1, +2, or +3", "wand-of-the-war-mage-1-2-or-3", 251, "wand", nil, "required", noCharges, []variant{{"plus-1", "Wand of the War Mage, +1", "uncommon"}, {"plus-2", "Wand of the War Mage, +2", "rare"}, {"plus-3", "Wand of the War Mage, +3", "very-rare"}}},
	}
	b := loadItems245(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-21" || b.CatalogueIssue != 245 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #245")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 250–252") {
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
	type member struct {
		name string
		page int
	}
	var members []member
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 245 {
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
		t.Fatalf("inventory assigns %d families to issue #245", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "250, 251" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.EnemyDetection.Action = "changed"
	if fresh := loadItems245(t); fresh.Families[0].Variants[0].Details.EnemyDetection.Action != "magic" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 14 variants, transcribed from pinned PDF
// pp. 250-252.
func TestItems245SourceValues(t *testing.T) {
	expected := []string{
		`{"enemy_detection":{"max_charges":7,"action":"magic","charge_cost":1,"duration_minutes":1,"reveals":"direction-of-nearest-hostile-creature","range_feet":60,"reveals_distance":false,"senses":["invisible","ethereal","disguised","hidden","in-plain-sight"],"ends_when":"you-stop-holding-the-wand","recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"fear":{"max_charges":7,"save_dc":15,"spells":[{"spell":"command","limit":"flee-or-grovel-only","charge_cost":1},{"spell":"fear","limit":"60-foot-cone","charge_cost":3}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"fireballs":{"attunement_by_any_of":["spellcaster"],"max_charges":7,"maximum_charges_per_cast":3,"spell":"fireball","save_dc":15,"base_charge_cost":1,"base_spell_level":3,"levels_per_additional_charge":1,"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"},"source_pages":[250,251]}}`,
		`{"lightning_bolts":{"attunement_by_any_of":["spellcaster"],"max_charges":7,"maximum_charges_per_cast":3,"spell":"lightning-bolt","save_dc":15,"base_charge_cost":1,"base_spell_level":3,"levels_per_additional_charge":1,"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"magic_detection":{"max_charges":3,"charge_cost":1,"spell":"detect-magic","recharge_dice":{"dice":1,"sides":3,"type":"charges"}}}`,
		`{"magic_missiles":{"max_charges":7,"maximum_charges_per_cast":3,"spell":"magic-missile","base_charge_cost":1,"base_spell_level":1,"levels_per_additional_charge":1,"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"paralysis":{"attunement_by_any_of":["spellcaster"],"max_charges":7,"action":"magic","charge_cost":1,"target":"one-visible-creature","range_feet":60,"save_ability":"constitution","save_dc":15,"condition":"paralyzed","duration_minutes":1,"repeat_save":"end-of-each-of-its-turns","repeat_success":"ends-effect-on-itself","recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"polymorph":{"attunement_by_any_of":["spellcaster"],"max_charges":7,"charge_cost":1,"spell":"polymorph","save_dc":15,"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"secrets":{"max_charges":3,"recharge_dice":{"dice":1,"sides":3,"type":"charges"},"action":"magic","charge_cost":1,"detects":["secret-door","trap"],"range_feet":60,"result":"pulses-and-points-at-nearest"}}`,
		`{"web":{"attunement_by_any_of":["spellcaster"],"max_charges":7,"charge_cost":1,"spell":"web","save_dc":13,"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-ashes"}}}`,
		`{"wonder":{"max_charges":7,"action":"magic","charge_cost":1,"point_range_feet":120,"spell_save_dc":15,"minimum_spell_range_feet":120,"multiple_subjects":"gm-determines-randomly","table_die_sides":100,"effects":[{"roll_from":1,"roll_to":20,"effect":"cast-spell-from-1d10-table"},{"roll_from":21,"roll_to":25,"effect":"you-are-stunned"},{"roll_from":26,"roll_to":30,"effect":"cast-gust-of-wind"},{"roll_from":31,"roll_to":35,"effect":"you-take-psychic-damage"},{"roll_from":36,"roll_to":40,"effect":"heavy-rain"},{"roll_from":41,"roll_to":45,"effect":"butterfly-cloud"},{"roll_from":46,"roll_to":50,"effect":"cast-lightning-bolt"},{"roll_from":51,"roll_to":55,"effect":"enlarge-closest-creature"},{"roll_from":56,"roll_to":60,"effect":"creature-appears-from-1d4-table"},{"roll_from":61,"roll_to":64,"effect":"grass-grows"},{"roll_from":65,"roll_to":68,"effect":"object-moves-to-ethereal-plane"},{"roll_from":69,"roll_to":72,"effect":"you-shrink"},{"roll_from":73,"roll_to":77,"effect":"leaves-grow-on-nearest-creature"},{"roll_from":78,"roll_to":82,"effect":"blinding-light-burst"},{"roll_from":83,"roll_to":87,"effect":"you-turn-invisible"},{"roll_from":88,"roll_to":92,"effect":"gem-stream"},{"roll_from":93,"roll_to":97,"effect":"polymorph-closest-creature-from-1d4-table"},{"roll_from":98,"roll_to":100,"effect":"petrify-closest-creature"}],"spell_table":{"die_sides":10,"rows":[{"roll_from":1,"roll_to":2,"result":"darkness"},{"roll_from":3,"roll_to":4,"result":"faerie-fire"},{"roll_from":5,"roll_to":6,"result":"fireball"},{"roll_from":7,"roll_to":8,"result":"slow"},{"roll_from":9,"roll_to":10,"result":"stinking-cloud"}]},"stunned":{"condition":"stunned","until":"start-of-your-next-turn","belief":"something-awesome-just-happened"},"gust_line":"from-you-to-point-of-origin","psychic_damage":{"dice":1,"sides":6,"type":"psychic"},"rain":{"duration_minutes":1,"height_feet":120,"radius_feet":60,"area":"lightly-obscured"},"butterflies":{"count":600,"height_feet":60,"radius_feet":30,"duration_minutes":10,"area":"heavily-obscured"},"lightning_line":"from-you-to-point-of-origin","enlarge":{"effect":"as-enlarge-reduce","fallback":"you-if-target-cannot-be-affected"},"creature":{"placement":"unoccupied-space-near-point-of-origin","control":"not-under-your-control","behavior":"acts-as-it-normally-would","duration_hours":1,"also_disappears_at_hit_points":0,"table":{"die_sides":4,"rows":[{"roll_from":1,"roll_to":1,"result":"rhinoceros"},{"roll_from":2,"roll_to":2,"result":"elephant"},{"roll_from":3,"roll_to":4,"result":"rat"}]}},"grass":{"radius_feet":60,"center":"as-close-to-point-of-origin-as-possible","size_multiplier":10,"overgrown_minutes":1},"ethereal":{"object":"gm-choice-not-worn-or-carried","range_feet":120,"maximum_dimension_feet":10,"without_object":"nothing-happens"},"shrink":{"effect":"as-enlarge-reduce-on-yourself","duration_minutes":1},"leaves":{"fall_off_after_hours":24,"early_removal":"picked-off"},"light_burst":{"emanation_feet":30,"save_ability":"constitution","save_dc":15,"condition":"blinded","duration_minutes":1,"repeat_save":"end-of-each-of-its-turns"},"invisibility_target":"yourself","gems":{"count_dice":{"dice":1,"sides":4,"type":"gems"},"count_multiplier":10,"value_gp_each":1,"line_length_feet":30,"line_width_feet":5,"direction":"toward-point-of-origin","damage_per_gem":1,"damage_type":"bludgeoning","damage_division":"total-divided-equally-in-line"},"polymorph_table":{"die_sides":4,"rows":[{"roll_from":1,"roll_to":1,"result":"black-bear"},{"roll_from":2,"roll_to":2,"result":"giant-wasp"},{"roll_from":3,"roll_to":4,"result":"frog"}]},"petrify":{"save_ability":"constitution","save_dc":15,"first_failure":"restrained-and-turning-to-stone","repeat_save":"end-of-its-next-turn","repeat_success":"effect-ends","repeat_failure":"petrified-instead-of-restrained","petrified_until":"freed-by-greater-restoration-or-similar-magic"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-dust"},"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":1,"when":"daily-at-dawn"},"source_pages":[251,252]}}`,
		`{"war_mage":{"attunement_by_any_of":["spellcaster"],"condition":"while-holding","spell_attack_bonus":1,"ignores_half_cover":true}}`,
		`{"war_mage":{"attunement_by_any_of":["spellcaster"],"condition":"while-holding","spell_attack_bonus":2,"ignores_half_cover":true}}`,
		`{"war_mage":{"attunement_by_any_of":["spellcaster"],"condition":"while-holding","spell_attack_bonus":3,"ignores_half_cover":true}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"hold:magic-action:expend-one-charge:item.attunement item.hold item.charges action.magic sense.direction-of-nearest-hostile-within-60-feet duration.1-minute item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-listed-charges:item.attunement item.hold item.charges spell.cast-listed-spells spell.save-dc-15 item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-up-to-three-charges:item.attunement attunement.spellcaster item.hold item.charges spell.cast-fireball spell.save-dc-15 spell.level-3-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-up-to-three-charges:item.attunement attunement.spellcaster item.hold item.charges spell.cast-lightning-bolt spell.save-dc-15 spell.level-3-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-one-charge:item.hold item.charges spell.cast-detect-magic",
		"hold:cast-spell:expend-up-to-three-charges:item.hold item.charges spell.cast-magic-missile spell.level-1-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge",
		"hold:magic-action:expend-one-charge target-within-60-feet:item.attunement attunement.spellcaster item.hold item.charges action.magic-ray save.constitution-dc-15 condition.paralyzed duration.1-minute save.repeat-at-turn-end item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-one-charge:item.attunement attunement.spellcaster item.hold item.charges spell.cast-polymorph spell.save-dc-15 item.destroyed-on-1-after-last-charge",
		"hold:magic-action:expend-one-charge:item.hold item.charges action.magic sense.nearest-secret-door-or-trap-within-60-feet",
		"hold:cast-spell:expend-one-charge:item.attunement attunement.spellcaster item.hold item.charges spell.cast-web spell.save-dc-13 item.destroyed-on-1-after-last-charge",
		"hold:magic-action:expend-one-charge point-within-120-feet:item.attunement item.hold item.charges action.magic table.wand-of-wonder-effects spell.save-dc-15 spell.minimum-range-120-feet item.destroyed-on-1-after-last-charge",
		"hold:none::item.attunement attunement.spellcaster item.hold bonus.spell-attack-by-rarity attack.ignore-half-cover-on-spell-attacks",
		"hold:none::item.attunement attunement.spellcaster item.hold bonus.spell-attack-by-rarity attack.ignore-half-cover-on-spell-attacks",
		"hold:none::item.attunement attunement.spellcaster item.hold bonus.spell-attack-by-rarity attack.ignore-half-cover-on-spell-attacks",
	}
	index := 0
	for _, f := range loadItems245(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails245
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

func TestItems245CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems245(t)
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
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems244(); return b.SourceIdentities(), err },
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
	if len(composed) != 252 || variants != 326 || composed[0].Name != "Potion of Healing" || composed[240].Name != "Wand of Enemy Detection" || composed[251].Name != "Wand of the War Mage, +1, +2, or +3" {
		t.Fatalf("expected Potion of Healing, batches #225 to #244 and the twelve families of #245: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_245") {
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
func TestItems245RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items245JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 800 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items245JSON, &fresh); err != nil {
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
		if _, err := DecodeItems245(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems245RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails245]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails245] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-20" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 244 }},
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
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{251, 251} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{250, 251} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"wand-of-enemy-detection rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"wand-of-enemy-detection attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"wand-of-fear rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"wand-of-fear attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"wand-of-fireballs rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"wand-of-fireballs attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "none" }},
		{"wand-of-lightning-bolts rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"wand-of-lightning-bolts attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "none" }},
		{"wand-of-magic-detection rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"wand-of-magic-detection attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "required" }},
		{"wand-of-magic-missiles rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"wand-of-magic-missiles attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "required" }},
		{"wand-of-paralysis rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"wand-of-paralysis attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "none" }},
		{"wand-of-polymorph rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"wand-of-polymorph attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"wand-of-secrets rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"wand-of-secrets attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "required" }},
		{"wand-of-web rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"wand-of-web attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"wand-of-wonder rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"wand-of-wonder attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-1 rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-1 attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-2 rarity", func(b *batch) { v(b, 11, 1).Rarity = "common" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-2 attunement", func(b *batch) { v(b, 11, 1).Attunement.Mode = "none" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-3 rarity", func(b *batch) { v(b, 11, 2).Rarity = "common" }},
		{"wand-of-the-war-mage-1-2-or-3 plus-3 attunement", func(b *batch) { v(b, 11, 2).Attunement.Mode = "none" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 11, 0).Charges.Maximum = 10 }},
		{"wand-of-enemy-detection dropped charges", func(b *batch) { *v(b, 0, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-enemy-detection charge maximum", func(b *batch) { v(b, 0, 0).Charges.Maximum++ }},
		{"wand-of-enemy-detection charge recharge", func(b *batch) { v(b, 0, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-fear dropped charges", func(b *batch) { *v(b, 1, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-fear charge maximum", func(b *batch) { v(b, 1, 0).Charges.Maximum++ }},
		{"wand-of-fear charge recharge", func(b *batch) { v(b, 1, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-fireballs dropped charges", func(b *batch) { *v(b, 2, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-fireballs charge maximum", func(b *batch) { v(b, 2, 0).Charges.Maximum++ }},
		{"wand-of-fireballs charge recharge", func(b *batch) { v(b, 2, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-lightning-bolts dropped charges", func(b *batch) { *v(b, 3, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-lightning-bolts charge maximum", func(b *batch) { v(b, 3, 0).Charges.Maximum++ }},
		{"wand-of-lightning-bolts charge recharge", func(b *batch) { v(b, 3, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-magic-detection dropped charges", func(b *batch) { *v(b, 4, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-magic-detection charge maximum", func(b *batch) { v(b, 4, 0).Charges.Maximum++ }},
		{"wand-of-magic-detection charge recharge", func(b *batch) { v(b, 4, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-magic-missiles dropped charges", func(b *batch) { *v(b, 5, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-magic-missiles charge maximum", func(b *batch) { v(b, 5, 0).Charges.Maximum++ }},
		{"wand-of-magic-missiles charge recharge", func(b *batch) { v(b, 5, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-paralysis dropped charges", func(b *batch) { *v(b, 6, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-paralysis charge maximum", func(b *batch) { v(b, 6, 0).Charges.Maximum++ }},
		{"wand-of-paralysis charge recharge", func(b *batch) { v(b, 6, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-polymorph dropped charges", func(b *batch) { *v(b, 7, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-polymorph charge maximum", func(b *batch) { v(b, 7, 0).Charges.Maximum++ }},
		{"wand-of-polymorph charge recharge", func(b *batch) { v(b, 7, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-secrets dropped charges", func(b *batch) { *v(b, 8, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-secrets charge maximum", func(b *batch) { v(b, 8, 0).Charges.Maximum++ }},
		{"wand-of-secrets charge recharge", func(b *batch) { v(b, 8, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-web dropped charges", func(b *batch) { *v(b, 9, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-web charge maximum", func(b *batch) { v(b, 9, 0).Charges.Maximum++ }},
		{"wand-of-web charge recharge", func(b *batch) { v(b, 9, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"wand-of-wonder dropped charges", func(b *batch) { *v(b, 10, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"wand-of-wonder charge maximum", func(b *batch) { v(b, 10, 0).Charges.Maximum++ }},
		{"wand-of-wonder charge recharge", func(b *batch) { v(b, 10, 0).Charges.Recharge = "regain-all-at-dawn" }},
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
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails245{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Fear = v(b, 1, 0).Details.Fear }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.EnemyDetection.Action = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.EnemyDetection.Action = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.EnemyDetection.MaxCharges = -1 }},
		{"detection range", func(b *batch) { p := v(b, 0, 0).Details.EnemyDetection; p.RangeFeet = 120 }},
		{"detection distance", func(b *batch) { p := v(b, 0, 0).Details.EnemyDetection; p.RevealsDistance = true }},
		{"detection senses", func(b *batch) { p := v(b, 0, 0).Details.EnemyDetection; p.Senses = p.Senses[:4] }},
		{"fear dc", func(b *batch) { p := v(b, 1, 0).Details.Fear; p.SaveDC = 17 }},
		{"fear limit", func(b *batch) { p := v(b, 1, 0).Details.Fear; p.Spells[0].Limit = "any-command" }},
		{"fear cost", func(b *batch) { p := v(b, 1, 0).Details.Fear; p.Spells[1].ChargeCost = 2 }},
		{"fireball dc", func(b *batch) { p := v(b, 2, 0).Details.Fireballs; p.SaveDC = 17 }},
		{"fireball level", func(b *batch) { p := v(b, 2, 0).Details.Fireballs; p.BaseSpellLevel = 5 }},
		{"fireball cast limit", func(b *batch) { p := v(b, 2, 0).Details.Fireballs; p.MaximumChargesPerCast = 7 }},
		{"lightning-bolt dc", func(b *batch) { p := v(b, 3, 0).Details.LightningBolts; p.SaveDC = 17 }},
		{"lightning-bolt level", func(b *batch) { p := v(b, 3, 0).Details.LightningBolts; p.BaseSpellLevel = 5 }},
		{"lightning-bolt cast limit", func(b *batch) { p := v(b, 3, 0).Details.LightningBolts; p.MaximumChargesPerCast = 7 }},
		{"magic detection charges", func(b *batch) { p := v(b, 4, 0).Details.MagicDetection; p.MaxCharges = 7 }},
		{"magic detection dice", func(b *batch) { p := v(b, 4, 0).Details.MagicDetection; p.RechargeDice.Sides = 6 }},
		{"missiles level", func(b *batch) { p := v(b, 5, 0).Details.MagicMissiles; p.BaseSpellLevel = 3 }},
		{"missiles cast limit", func(b *batch) { p := v(b, 5, 0).Details.MagicMissiles; p.MaximumChargesPerCast = 7 }},
		{"paralysis dc", func(b *batch) { p := v(b, 6, 0).Details.Paralysis; p.SaveDC = 17 }},
		{"paralysis range", func(b *batch) { p := v(b, 6, 0).Details.Paralysis; p.RangeFeet = 120 }},
		{"paralysis duration", func(b *batch) { p := v(b, 6, 0).Details.Paralysis; p.DurationMinutes = 10 }},
		{"polymorph dc", func(b *batch) { p := v(b, 7, 0).Details.Polymorph; p.SaveDC = 17 }},
		{"polymorph cost", func(b *batch) { p := v(b, 7, 0).Details.Polymorph; p.ChargeCost = 2 }},
		{"polymorph class", func(b *batch) { p := v(b, 7, 0).Details.Polymorph; p.AttunementByAnyOf = []string{"wizard"} }},
		{"secrets range", func(b *batch) { p := v(b, 8, 0).Details.Secrets; p.RangeFeet = 30 }},
		{"secrets targets", func(b *batch) { p := v(b, 8, 0).Details.Secrets; p.Detects = p.Detects[:1] }},
		{"web dc", func(b *batch) { p := v(b, 9, 0).Details.Web; p.SaveDC = 17 }},
		{"web cost", func(b *batch) { p := v(b, 9, 0).Details.Web; p.ChargeCost = 2 }},
		{"web class", func(b *batch) { p := v(b, 9, 0).Details.Web; p.AttunementByAnyOf = []string{"wizard"} }},
		{"wonder dc", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.SpellSaveDC = 13 }},
		{"wonder table gap", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.Effects = append(p.Effects[:3], p.Effects[4:]...) }},
		{"wonder row overlap", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.Effects[1].RollFrom = 20 }},
		{"wonder spell table", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.SpellTable.Rows = p.SpellTable.Rows[:4] }},
		{"wonder creature table", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.Creature.Table.Rows[2].RollTo = 3 }},
		{"wonder form table", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.PolymorphTable.DieSides = 6 }},
		{"wonder duplicate effect", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.Effects[1].Effect = p.Effects[0].Effect }},
		{"wonder gems", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.Gems.CountMultiplier = 100 }},
		{"wonder pages", func(b *batch) { p := v(b, 10, 0).Details.Wonder; p.SourcePages = []int{251} }},
		{"war mage bonus", func(b *batch) { p := v(b, 11, 0).Details.WarMage; p.SpellAttackBonus = 3 }},
		{"war mage cover", func(b *batch) { p := v(b, 11, 0).Details.WarMage; p.IgnoresHalfCover = false }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems245(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems245(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems245RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/wand-of-enemy-detection/variant/standard"`)
	if !bytes.Contains(items245JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items245JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems245(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"enemy_detection": {`, `"enemy_detection_plus": {`},
		"unknown detail":     {`"max_charges": 7`, `"max_charges_all": 7`},
		"string for number":  {`"max_charges": 7`, `"max_charges": "7"`},
		"fractional number":  {`"max_charges": 7`, `"max_charges": 7.5`},
		"number for boolean": {`"reveals_distance": false`, `"reveals_distance": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items245JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items245JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems245(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items245JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems245(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items245JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems244(items245JSON); err == nil {
		t.Fatal("#244 decoder accepted the #245 batch")
	}
	if _, err := DecodeItems245(items244JSON); err == nil {
		t.Fatal("#245 decoder accepted the #244 batch")
	}
}
