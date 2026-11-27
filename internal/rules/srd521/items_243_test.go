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

func loadItems243(t *testing.T) ItemSourceBatch[ItemDetails243] {
	t.Helper()
	b, err := LoadItems243()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 246-248, in planning inventory order for issue #243.
func TestItems243SourceIdentities(t *testing.T) {
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
		{"Staff of Swarming Insects", "staff-of-swarming-insects", 246, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"}, []variant{{"standard", "Staff of Swarming Insects", "rare"}}},
		{"Staff of Thunder and Lightning", "staff-of-thunder-and-lightning", 247, "staff", nil, "required", noCharges, []variant{{"standard", "Staff of Thunder and Lightning", "very-rare"}}},
		{"Staff of Withering", "staff-of-withering", 248, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}, []variant{{"standard", "Staff of Withering", "rare"}}},
		{"Staff of the Magi", "staff-of-the-magi", 246, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 50, PerUse: 1, Recharge: "regain-4d6-plus-2-at-dawn"}, []variant{{"standard", "Staff of the Magi", "legendary"}}},
		{"Staff of the Python", "staff-of-the-python", 247, "staff", nil, "required", noCharges, []variant{{"standard", "Staff of the Python", "uncommon"}}},
		{"Staff of the Woodlands", "staff-of-the-woodlands", 247, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "regain-1d6-at-dawn"}, []variant{{"standard", "Staff of the Woodlands", "rare"}}},
		{"Stone of Controlling Earth Elementals", "stone-of-controlling-earth-elementals", 248, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Stone of Controlling Earth Elementals", "rare"}}},
		{"Stone of Good Luck (Luckstone)", "stone-of-good-luck-luckstone", 248, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Stone of Good Luck (Luckstone)", "uncommon"}}},
		{"Sun Blade", "sun-blade", 248, "weapon", []string{"longsword"}, "required", noCharges, []variant{{"standard", "Sun Blade", "rare"}}},
		{"Sword of Life Stealing", "sword-of-life-stealing", 248, "weapon", []string{"glaive", "greatsword", "longsword", "rapier", "scimitar", "shortsword"}, "required", noCharges, []variant{{"standard", "Sword of Life Stealing", "rare"}}},
		{"Sword of Sharpness", "sword-of-sharpness", 248, "weapon", []string{"glaive", "greatsword", "longsword", "scimitar"}, "required", noCharges, []variant{{"standard", "Sword of Sharpness", "very-rare"}}},
		{"Sword of Wounding", "sword-of-wounding", 248, "weapon", []string{"glaive", "greatsword", "longsword", "rapier", "scimitar", "shortsword"}, "required", noCharges, []variant{{"standard", "Sword of Wounding", "rare"}}},
	}
	b := loadItems243(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-19" || b.CatalogueIssue != 243 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #243")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 246–248") {
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
		if r.Issue != 243 {
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
		t.Fatalf("inventory assigns %d families to issue #243", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "246, 247, 248" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.SwarmingInsects.SaveAndAttackSource = "changed"
	if fresh := loadItems243(t); fresh.Families[0].Variants[0].Details.SwarmingInsects.SaveAndAttackSource != "your-spell-save-dc-and-spell-attack-modifier" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 246-248.
func TestItems243SourceValues(t *testing.T) {
	expected := []string{
		`{"swarming_insects":{"attunement_by_any_of":["bard","cleric","druid","sorcerer","warlock","wizard"],"max_charges":10,"insect_cloud":{"action":"magic","charge_cost":1,"swarm":"harmless-flying-insects","emanation_feet":30,"duration_minutes":10,"area":"heavily-obscured-for-creatures-other-than-you","dispersed_by":"strong-wind"},"save_and_attack_source":"your-spell-save-dc-and-spell-attack-modifier","spells":[{"spell":"giant-insect","charge_cost":4},{"spell":"insect-plague","charge_cost":5}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"insect-swarm-consumes-staff-then-disperses"}}}`,
		`{"thunder_lightning":{"wielded_as":"quarterstaff","attack_bonus":2,"damage_bonus":2,"each_property_reuse_after":"next-dawn","lightning":{"trigger":"hit-with-melee-attack","extra_damage":{"dice":2,"sides":6,"type":"lightning"},"action_required":false},"thunder":{"trigger":"hit-with-melee-attack","audible_feet":300,"action_required":false,"save_ability":"constitution","save_dc":17,"condition":"stunned","condition_ends":"end-of-your-next-turn"},"thunder_and_lightning":{"trigger":"immediately-after-melee-hit","action":"bonus","uses":["lightning","thunder"],"spends_their_daily_uses":false},"lightning_strike":{"action":"magic","line_width_feet":5,"line_length_feet":120,"save_ability":"dexterity","save_dc":17,"damage":{"dice":9,"sides":6,"type":"lightning"},"success_damage":"half"},"thunderclap":{"action":"magic","audible_feet":600,"emanation_feet":60,"save_ability":"constitution","save_dc":17,"damage":{"dice":2,"sides":6,"type":"thunder"},"failure_condition":"deafened","condition_minutes":1,"success":"half-damage-only"},"source_pages":[247,248]}}`,
		`{"withering":{"wielded_as":"quarterstaff","hit_damage":"as-normal-quarterstaff","max_charges":3,"recharge_dice":{"dice":1,"sides":3,"type":"charges"},"charge_cost":1,"extra_damage":{"dice":2,"sides":10,"type":"necrotic"},"save_ability":"constitution","save_dc":15,"failure":{"mode":"disadvantage","duration_hours":1,"applies_to":["strength-checks-and-saves","constitution-checks-and-saves"]}}}`,
		`{"magi":{"attunement_by_any_of":["sorcerer","warlock","wizard"],"max_charges":50,"wielded_as":"quarterstaff","attack_bonus":2,"damage_bonus":2,"spell_attack_bonus":2,"spell_absorption":{"save_mode":"advantage","trigger":"spell-from-another-creature-targets-only-you","action":"reaction","result":"spell-canceled-and-absorbed","charges_gained":"absorbed-spell-level","explodes_above_charges":50,"explosion":"as-retributive-strike"},"save_dc_source":"your-spell-save-dc","spells":[{"spell":"arcane-lock","printed_level":"none","charge_cost":0},{"spell":"conjure-elemental","printed_level":"none","charge_cost":7},{"spell":"detect-magic","printed_level":"none","charge_cost":0},{"spell":"dispel-magic","printed_level":"none","charge_cost":3},{"spell":"enlarge-reduce","printed_level":"none","charge_cost":0},{"spell":"fireball","printed_level":"level-7-version","charge_cost":7},{"spell":"flaming-sphere","printed_level":"none","charge_cost":2},{"spell":"ice-storm","printed_level":"none","charge_cost":4},{"spell":"invisibility","printed_level":"none","charge_cost":2},{"spell":"knock","printed_level":"none","charge_cost":2},{"spell":"light","printed_level":"none","charge_cost":0},{"spell":"lightning-bolt","printed_level":"level-7-version","charge_cost":7},{"spell":"mage-hand","printed_level":"none","charge_cost":0},{"spell":"passwall","printed_level":"none","charge_cost":5},{"spell":"plane-shift","printed_level":"none","charge_cost":7},{"spell":"protection-from-evil-and-good","printed_level":"none","charge_cost":0},{"spell":"telekinesis","printed_level":"none","charge_cost":5},{"spell":"wall-of-fire","printed_level":"none","charge_cost":4},{"spell":"web","printed_level":"none","charge_cost":2}],"recharge":{"dice":{"dice":4,"sides":6,"type":"charges"},"bonus":2,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"regains_on":20,"regains":{"dice":1,"sides":12,"type":"charges"},"regains_bonus":1},"retributive_strike":{"action":"magic","area_emanation_feet":30,"planar_travel_percent":50,"your_damage_per_charge":16,"other_damage_per_charge":6,"damage_type":"force","save_ability":"dexterity","save_dc":17,"success_damage":"half"},"source_pages":[246,247]}}`,
		`{"python":{"action":"magic","throw_range_feet":10,"creature":"giant-constrictor-snake","control":"under-your-control","initiative":"shares-yours-and-acts-immediately-after","command":{"timing":"your-turn","method":"mentally","action_required":false,"range_feet":60,"blocked_by_condition":"incapacitated","options":["decide-action-and-movement","general-command"],"without_commands":"defends-itself"},"revert":{"action":"bonus","location":"snakes-current-space","reuse_after_hours":1},"at_zero_hit_points":"dies-reverts-then-staff-shatters","hit_points_on_early_revert":"all-regained"}}`,
		`{"woodlands":{"attunement_by_any_of":["druid"],"max_charges":6,"wielded_as":"quarterstaff","attack_bonus":2,"damage_bonus":2,"spell_attack_bonus":2,"save_dc_source":"your-spell-save-dc","spells":[{"spell":"animal-friendship","charge_cost":1},{"spell":"awaken","charge_cost":5},{"spell":"barkskin","charge_cost":2},{"spell":"locate-animals-or-plants","charge_cost":2},{"spell":"pass-without-trace","charge_cost":2},{"spell":"speak-with-animals","charge_cost":1},{"spell":"speak-with-plants","charge_cost":3},{"spell":"wall-of-thorns","charge_cost":6}],"tree_form":{"action":"magic","charge_cost":1,"placement":"one-end-in-earth-in-unoccupied-space","tree_height_feet":60,"trunk_diameter_feet":5,"branch_radius_feet":20,"aura":"faint-transmutation-discerned-with-detect-magic","revert":"touch-tree-and-use-magic-action","creatures_in_tree":"fall"},"recharge_dice":{"dice":1,"sides":6,"type":"charges"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"becomes-nonmagical-quarterstaff"}}}`,
		`{"earth_stone":{"weight_pounds":5,"condition":"while-touching-stone-to-ground","action":"magic","creature":"earth-elemental","summon_range_feet":30,"obeys_commands":true,"initiative":"acts-immediately-after-you-on-your-count","duration_hours":1,"also_disappears_when":["it-dies","dismissed-as-bonus-action"],"reuse_after":"next-dawn"}}`,
		`{"luckstone":{"appearance":"polished-agate","condition":"while-on-your-person","ability_check_bonus":1,"saving_throw_bonus":1}}`,
		`{"sun_blade":{"appearance":"sword-hilt","blade":{"condition":"while-grasping-hilt","action":"bonus","options":["appear","disappear"],"functions_as":"longsword","property":"finesse","proficiency_if_proficient_with":["longswords","shortswords"]},"attack_bonus":2,"damage_bonus":2,"damage_type":"radiant","replaces_damage_type":"slashing","undead_extra_damage":{"dice":1,"sides":8,"type":"radiant"},"sunlight":{"bright_radius_feet":15,"dim_additional_feet":15,"light_is_sunlight":true,"adjust_action":"magic","adjust_step_feet":5,"maximum_each_feet":30,"minimum_each_feet":10}}}`,
		`{"life_stealing":{"trigger_natural_roll":20,"target":"creature","extra_damage":15,"damage_type":"necrotic","excluded_creature_types":["construct","undead"],"temporary_hit_points":"equal-to-necrotic-damage-taken"}}`,
		`{"sharpness":{"object_hit":"maximize-weapon-damage-dice","trigger_natural_roll":20,"target":"creature","extra_damage":14,"damage_type":"slashing","applies_condition":"exhaustion","exhaustion_levels":1}}`,
		`{"wounding":{"trigger":"hit-creature-with-attack","extra_damage":{"dice":2,"sides":6,"type":"necrotic"},"save_ability":"constitution","save_dc":15,"failure":"cannot-regain-hit-points","duration_hours":1,"repeat_save":"end-of-each-of-its-turns","repeat_success":"ends-effect-on-itself"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"hold:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.hold item.charges action.magic-insect-cloud effect.heavily-obscured spell.cast-giant-insect spell.cast-insect-plague spell.uses-your-save-dc-and-attack-modifier item.destroyed-on-1-after-last-charge",
		"wield:weapon-attack:each-property-once-until-dawn:item.attunement item.magic-weapon bonus.attack-and-damage-2 damage.extra-2d6-lightning save.constitution-dc-17 condition.stunned action.bonus-both-properties save.dexterity-dc-17 damage.lightning-strike-9d6 damage.thunderclap-2d6 condition.deafened item.each-property-once-until-dawn",
		"wield:weapon-attack:expend-one-charge-on-hit:item.attunement item.magic-weapon item.charges damage.extra-2d10-necrotic save.constitution-dc-15 effect.disadvantage-strength-and-constitution duration.1-hour",
		"wield:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.magic-weapon bonus.attack-and-damage-2 bonus.spell-attack-2 item.charges save.advantage-against-spells action.reaction-spell-absorption spell.cast-listed-spells spell.uses-your-save-dc item.retributive-strike damage.force-by-charges save.dexterity-dc-17",
		"hold:magic-action:throw-within-10-feet:item.attunement action.magic-throw-staff creature.giant-constrictor-snake creature.under-your-control creature.initiative-after-owner action.bonus-revert item.1-hour-reuse item.destroyed-if-snake-dies",
		"wield:cast-spell:expend-listed-charges:item.attunement attunement.druid item.magic-weapon bonus.attack-and-damage-2 bonus.spell-attack-2 item.charges spell.cast-listed-spells spell.uses-your-save-dc action.magic-tree-form item.nonmagical-on-1-after-last-charge",
		"touch:magic-action:touch-stone-to-ground:action.magic creature.summon-earth-elemental creature.obeys-commands creature.initiative-after-owner duration.1-hour item.once-until-dawn",
		"carry:none::item.attunement item.carry check.all-bonus-1 save.all-bonus-1",
		"wield:weapon-attack::item.attunement item.magic-weapon item.blade-of-radiance property.finesse bonus.attack-and-damage-2 damage.radiant-replaces-slashing damage.extra-1d8-radiant-against-undead light.sunlight",
		"wield:weapon-attack::item.attunement item.magic-weapon damage.extra-15-necrotic-on-natural-20 hit-points.temporary-equal-to-necrotic-damage",
		"wield:weapon-attack::item.attunement item.magic-weapon damage.maximize-weapon-dice-against-objects damage.extra-14-slashing-on-natural-20 condition.exhaustion",
		"wield:weapon-attack::item.attunement item.magic-weapon damage.extra-2d6-necrotic save.constitution-dc-15 effect.cannot-regain-hit-points duration.1-hour save.repeat-at-turn-end",
	}
	index := 0
	for _, f := range loadItems243(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails243
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

func TestItems243CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems243(t)
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
	if len(composed) != 228 || variants != 300 || composed[0].Name != "Potion of Healing" || composed[216].Name != "Staff of Swarming Insects" || composed[227].Name != "Sword of Wounding" {
		t.Fatalf("expected Potion of Healing, batches #225 to #242 and the twelve families of #243: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_243") {
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
func TestItems243RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items243JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 700 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items243JSON, &fresh); err != nil {
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
		if _, err := DecodeItems243(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems243RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails243]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails243] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-18" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 242 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/talisman-of-pure-good"
			b.Families[12].Name = "Talisman of Pure Good"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"sun-blade widened base", func(b *batch) { b.Families[8].Base.Options = append(b.Families[8].Base.Options, "any-weapon") }},
		{"sword-of-life-stealing widened base", func(b *batch) { b.Families[9].Base.Options = append(b.Families[9].Base.Options, "any-weapon") }},
		{"sword-of-sharpness widened base", func(b *batch) { b.Families[10].Base.Options = append(b.Families[10].Base.Options, "any-weapon") }},
		{"sword-of-wounding widened base", func(b *batch) { b.Families[11].Base.Options = append(b.Families[11].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{247, 247} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{246, 247} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"staff-of-swarming-insects rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"staff-of-thunder-and-lightning rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"staff-of-withering rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"staff-of-the-magi rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"staff-of-the-python rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"staff-of-the-woodlands rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"stone-of-controlling-earth-elementals rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"stone-of-good-luck-luckstone rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"sun-blade rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"sword-of-life-stealing rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"sword-of-sharpness rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"sword-of-wounding rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"staff-of-swarming-insects attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"staff-of-thunder-and-lightning attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"staff-of-withering attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "none" }},
		{"staff-of-the-magi attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "none" }},
		{"staff-of-the-python attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "none" }},
		{"staff-of-the-woodlands attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "none" }},
		{"stone-of-controlling-earth-elementals attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "required" }},
		{"stone-of-good-luck-luckstone attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"sun-blade attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "none" }},
		{"sword-of-life-stealing attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"sword-of-sharpness attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"sword-of-wounding attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 1, 0).Charges.Maximum = 10 }},
		{"staff-of-swarming-insects dropped charges", func(b *batch) { *v(b, 0, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-swarming-insects charge maximum", func(b *batch) { v(b, 0, 0).Charges.Maximum++ }},
		{"staff-of-swarming-insects charge recharge", func(b *batch) { v(b, 0, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-withering dropped charges", func(b *batch) { *v(b, 2, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-withering charge maximum", func(b *batch) { v(b, 2, 0).Charges.Maximum++ }},
		{"staff-of-withering charge recharge", func(b *batch) { v(b, 2, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-the-magi dropped charges", func(b *batch) { *v(b, 3, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-the-magi charge maximum", func(b *batch) { v(b, 3, 0).Charges.Maximum++ }},
		{"staff-of-the-magi charge recharge", func(b *batch) { v(b, 3, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-the-woodlands dropped charges", func(b *batch) { *v(b, 5, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-the-woodlands charge maximum", func(b *batch) { v(b, 5, 0).Charges.Maximum++ }},
		{"staff-of-the-woodlands charge recharge", func(b *batch) { v(b, 5, 0).Charges.Recharge = "regain-all-at-dawn" }},
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
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails243{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.ThunderLightning = v(b, 1, 0).Details.ThunderLightning }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.SwarmingInsects.SaveAndAttackSource = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.SwarmingInsects.SaveAndAttackSource = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.SwarmingInsects.MaxCharges = -1 }},
		{"insects charges", func(b *batch) { p := v(b, 0, 0).Details.SwarmingInsects; p.MaxCharges = 6 }},
		{"insects cloud", func(b *batch) { p := v(b, 0, 0).Details.SwarmingInsects; p.InsectCloud.EmanationFeet = 20 }},
		{"insects cost", func(b *batch) { p := v(b, 0, 0).Details.SwarmingInsects; p.Spells[0].ChargeCost = 5 }},
		{"insects class", func(b *batch) { p := v(b, 0, 0).Details.SwarmingInsects; p.AttunementByAnyOf = p.AttunementByAnyOf[:1] }},
		{"thunder bonus", func(b *batch) { p := v(b, 1, 0).Details.ThunderLightning; p.AttackBonus = 3 }},
		{"thunder stun dc", func(b *batch) { p := v(b, 1, 0).Details.ThunderLightning; p.Thunder.SaveDC = 15 }},
		{"thunder strike dice", func(b *batch) { p := v(b, 1, 0).Details.ThunderLightning; p.LightningStrike.Damage.Dice = 8 }},
		{"thunder clap dc", func(b *batch) { p := v(b, 1, 0).Details.ThunderLightning; p.Thunderclap.SaveDC = 15 }},
		{"thunder pages", func(b *batch) { p := v(b, 1, 0).Details.ThunderLightning; p.SourcePages = []int{247} }},
		{"withering charges", func(b *batch) { p := v(b, 2, 0).Details.Withering; p.MaxCharges = 10 }},
		{"withering dc", func(b *batch) { p := v(b, 2, 0).Details.Withering; p.SaveDC = 17 }},
		{"withering dice", func(b *batch) { p := v(b, 2, 0).Details.Withering; p.ExtraDamage.Sides = 6 }},
		{"magi charges", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.MaxCharges = 20 }},
		{"magi absorb limit", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.SpellAbsorption.ExplodesAboveCharges = 49 }},
		{"magi spell cost", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.Spells[1].ChargeCost = 5 }},
		{"magi missing spell", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.Spells = p.Spells[:18] }},
		{"magi multiplier", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.RetributiveStrike.OtherDamagePerCharge = 4 }},
		{"magi pages", func(b *batch) { p := v(b, 3, 0).Details.Magi; p.SourcePages = []int{246} }},
		{"python throw", func(b *batch) { p := v(b, 4, 0).Details.Python; p.ThrowRangeFeet = 30 }},
		{"python command", func(b *batch) { p := v(b, 4, 0).Details.Python; p.Command.RangeFeet = 120 }},
		{"python reuse", func(b *batch) { p := v(b, 4, 0).Details.Python; p.Revert.ReuseAfterHours = 24 }},
		{"woodlands charges", func(b *batch) { p := v(b, 5, 0).Details.Woodlands; p.MaxCharges = 10 }},
		{"woodlands bonus", func(b *batch) { p := v(b, 5, 0).Details.Woodlands; p.SpellAttackBonus = 1 }},
		{"woodlands tree", func(b *batch) { p := v(b, 5, 0).Details.Woodlands; p.TreeForm.TreeHeightFeet = 30 }},
		{"woodlands cost", func(b *batch) { p := v(b, 5, 0).Details.Woodlands; p.Spells[7].ChargeCost = 5 }},
		{"earth stone range", func(b *batch) { p := v(b, 6, 0).Details.EarthStone; p.SummonRangeFeet = 120 }},
		{"earth stone duration", func(b *batch) { p := v(b, 6, 0).Details.EarthStone; p.DurationHours = 8 }},
		{"luckstone checks", func(b *batch) { p := v(b, 7, 0).Details.Luckstone; p.AbilityCheckBonus = 2 }},
		{"luckstone saves", func(b *batch) { p := v(b, 7, 0).Details.Luckstone; p.SavingThrowBonus = 0 }},
		{"sun blade bonus", func(b *batch) { p := v(b, 8, 0).Details.SunBlade; p.AttackBonus = 1 }},
		{"sun blade undead", func(b *batch) { p := v(b, 8, 0).Details.SunBlade; p.UndeadExtraDamage.Sides = 6 }},
		{"sun blade light", func(b *batch) { p := v(b, 8, 0).Details.SunBlade; p.Sunlight.BrightRadiusFeet = 20 }},
		{"sun blade maximum", func(b *batch) { p := v(b, 8, 0).Details.SunBlade; p.Sunlight.MaximumEachFeet = 60 }},
		{"life stealing damage", func(b *batch) { p := v(b, 9, 0).Details.LifeStealing; p.ExtraDamage = 10 }},
		{"life stealing roll", func(b *batch) { p := v(b, 9, 0).Details.LifeStealing; p.TriggerNaturalRoll = 19 }},
		{"life stealing types", func(b *batch) {
			p := v(b, 9, 0).Details.LifeStealing
			p.ExcludedCreatureTypes = p.ExcludedCreatureTypes[:1]
		}},
		{"sharpness damage", func(b *batch) { p := v(b, 10, 0).Details.Sharpness; p.ExtraDamage = 15 }},
		{"sharpness exhaustion", func(b *batch) { p := v(b, 10, 0).Details.Sharpness; p.ExhaustionLevels = 2 }},
		{"sharpness roll", func(b *batch) { p := v(b, 10, 0).Details.Sharpness; p.TriggerNaturalRoll = 19 }},
		{"wounding dc", func(b *batch) { p := v(b, 11, 0).Details.Wounding; p.SaveDC = 13 }},
		{"wounding dice", func(b *batch) { p := v(b, 11, 0).Details.Wounding; p.ExtraDamage.Dice = 1 }},
		{"wounding duration", func(b *batch) { p := v(b, 11, 0).Details.Wounding; p.DurationHours = 24 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems243(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems243(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems243RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/staff-of-swarming-insects/variant/standard"`)
	if !bytes.Contains(items243JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items243JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems243(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"swarming_insects": {`, `"swarming_insects_plus": {`},
		"unknown detail":     {`"max_charges": 10`, `"max_charges_all": 10`},
		"string for number":  {`"max_charges": 10`, `"max_charges": "10"`},
		"fractional number":  {`"max_charges": 10`, `"max_charges": 10.5`},
		"number for boolean": {`"obeys_commands": true`, `"obeys_commands": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items243JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items243JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems243(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items243JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems243(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items243JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems242(items243JSON); err == nil {
		t.Fatal("#242 decoder accepted the #243 batch")
	}
	if _, err := DecodeItems243(items242JSON); err == nil {
		t.Fatal("#243 decoder accepted the #242 batch")
	}
}
