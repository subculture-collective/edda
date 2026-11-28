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

func loadItems242(t *testing.T) ItemSourceBatch[ItemDetails242] {
	t.Helper()
	b, err := LoadItems242()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 243-246, in planning inventory order for issue #242.
func TestItems242SourceIdentities(t *testing.T) {
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
		{"Shield, +1, +2, or +3", "shield-1-2-or-3", 243, "armor", []string{"shield"}, "none", noCharges, []variant{{"plus-1", "Shield, +1", "uncommon"}, {"plus-2", "Shield, +2", "rare"}, {"plus-3", "Shield, +3", "very-rare"}}},
		{"Slippers of Spider Climbing", "slippers-of-spider-climbing", 244, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Slippers of Spider Climbing", "uncommon"}}},
		{"Sovereign Glue", "sovereign-glue", 244, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Sovereign Glue", "legendary"}}},
		{"Spell Scroll", "spell-scroll", 244, "scroll", nil, "none", noCharges, []variant{{"cantrip", "Spell Scroll (cantrip)", "common"}, {"level-1", "Spell Scroll (level 1)", "common"}, {"level-2", "Spell Scroll (level 2)", "uncommon"}, {"level-3", "Spell Scroll (level 3)", "uncommon"}, {"level-4", "Spell Scroll (level 4)", "rare"}, {"level-5", "Spell Scroll (level 5)", "rare"}, {"level-6", "Spell Scroll (level 6)", "very-rare"}, {"level-7", "Spell Scroll (level 7)", "very-rare"}, {"level-8", "Spell Scroll (level 8)", "very-rare"}, {"level-9", "Spell Scroll (level 9)", "legendary"}}},
		{"Spellguard Shield", "spellguard-shield", 244, "armor", []string{"shield"}, "required", noCharges, []variant{{"standard", "Spellguard Shield", "very-rare"}}},
		{"Sphere of Annihilation", "sphere-of-annihilation", 244, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Sphere of Annihilation", "legendary"}}},
		{"Staff of Charming", "staff-of-charming", 245, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d8-plus-2-at-dawn"}, []variant{{"standard", "Staff of Charming", "rare"}}},
		{"Staff of Fire", "staff-of-fire", 245, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"}, []variant{{"standard", "Staff of Fire", "very-rare"}}},
		{"Staff of Frost", "staff-of-frost", 245, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"}, []variant{{"standard", "Staff of Frost", "very-rare"}}},
		{"Staff of Healing", "staff-of-healing", 245, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"}, []variant{{"standard", "Staff of Healing", "rare"}}},
		{"Staff of Power", "staff-of-power", 245, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 20, PerUse: 1, Recharge: "regain-2d8-plus-4-at-dawn"}, []variant{{"standard", "Staff of Power", "very-rare"}}},
		{"Staff of Striking", "staff-of-striking", 246, "staff", nil, "required", ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"}, []variant{{"standard", "Staff of Striking", "very-rare"}}},
	}
	b := loadItems242(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-18" || b.CatalogueIssue != 242 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #242")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 243–246") {
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
			presentations := []ItemReference{{ID: "srd-5.2.1:content/magic-items/" + want.slug, Pages: []int{want.page, want.page}}}
			if want.slug == "spell-scroll" {
				presentations = append(presentations, ItemReference{ID: spellScrollGearID, Pages: []int{99, 99}})
			}
			if !reflect.DeepEqual(f.Presentations, presentations) {
				t.Fatalf("source presentation mismatch: %+v", f.Presentations)
			}
			if len(f.Variants) != len(want.variants) {
				t.Fatalf("the source lists %d forms of %s, got %d variants", len(want.variants), want.name, len(f.Variants))
			}
			for j, wv := range want.variants {
				v := f.Variants[j]
				if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != wv.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				// Only the p. 99 gear entry prints a price, for two scrolls.
				cost := ItemCost{Kind: "unspecified", Currency: "none", SourcePages: []int{want.page, want.page}}
				if gp := map[string]int{"Spell Scroll (cantrip)": 30, "Spell Scroll (level 1)": 50}[wv.name]; gp != 0 {
					cost = ItemCost{Kind: "listed", Amount: gp, Currency: "gp", SourcePages: []int{99, 99}}
				}
				if !reflect.DeepEqual(*v.Cost, cost) {
					t.Fatalf("price mismatch: %+v", *v.Cost)
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
	if total != 23 {
		t.Fatalf("expected 23 source-listed variants, got %d", total)
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
		if r.Issue != 242 {
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
		t.Fatalf("inventory assigns %d families to issue #242", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "243, 244, 245, 246" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.ShieldBonus.Condition = "changed"
	if fresh := loadItems242(t); fresh.Families[0].Variants[0].Details.ShieldBonus.Condition != "while-holding" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 23 variants, transcribed from pinned PDF pp.
// 243-246.
func TestItems242SourceValues(t *testing.T) {
	expected := []string{
		`{"shield_bonus":{"condition":"while-holding","armor_class_bonus":1,"stacks_with_shield_bonus":true}}`,
		`{"shield_bonus":{"condition":"while-holding","armor_class_bonus":2,"stacks_with_shield_bonus":true}}`,
		`{"shield_bonus":{"condition":"while-holding","armor_class_bonus":3,"stacks_with_shield_bonus":true}}`,
		`{"spider_slippers":{"condition":"while-wearing","movement":["up-down-and-across-vertical-surfaces","along-ceilings"],"hands_free":true,"climb_speed":"equal-to-speed","excluded_surfaces":"slippery-such-as-ice-or-oil"}}`,
		`{"glue":{"appearance":"viscous-milky-white-substance","bond":"permanent-adhesive-between-two-objects","storage":"jar-or-flask-coated-with-oil-of-slipperiness","found_ounces":{"dice":1,"sides":6,"type":"ounces"},"found_ounces_bonus":1,"coverage_ounces":1,"coverage_square_feet":1,"apply_action":"utilize","set_minutes":1,"broken_only_by":["universal-solvent","oil-of-etherealness","wish"]}}`,
		`{"scroll":{"spell_level":0,"rarity":"common","save_dc":13,"attack_bonus":5,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":1,"rarity":"common","save_dc":13,"attack_bonus":5,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":2,"rarity":"uncommon","save_dc":13,"attack_bonus":5,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":3,"rarity":"uncommon","save_dc":15,"attack_bonus":7,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":4,"rarity":"rare","save_dc":15,"attack_bonus":7,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":5,"rarity":"rare","save_dc":17,"attack_bonus":9,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":6,"rarity":"very-rare","save_dc":17,"attack_bonus":9,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":7,"rarity":"very-rare","save_dc":18,"attack_bonus":10,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":8,"rarity":"very-rare","save_dc":18,"attack_bonus":10,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"scroll":{"spell_level":9,"rarity":"legendary","save_dc":19,"attack_bonus":11,"bears":"one-spell","readable_if":"spell-on-your-spell-list","material_components_required":false,"otherwise":"unintelligible","casting_time":"spells-normal-casting-time","after_casting":"crumbles-to-dust","lost_if_casting_interrupted":false,"higher_level_check":{"trigger":"higher-level-than-you-can-normally-cast","ability":"your-spellcasting-ability","dc_base":10,"dc_adds":"spell-level","failure":"spell-disappears-with-no-other-effect"},"copying":{"eligible":"wizard-spell","check_ability":"intelligence","check_skill":"arcana","dc_base":10,"dc_adds":"spell-level","success":"spell-copied","scroll_destroyed_either_way":true}}}`,
		`{"spellguard":{"condition":"while-holding","save_mode":"advantage","against":["spells","other-magical-effects"],"spell_attacks_against_you":"disadvantage"}}`,
		`{"annihilation":{"diameter_feet":2,"nature":"hole-in-the-multiverse","obliterates":"all-matter-passing-through-or-passed-through","artifact_exception":"artifacts-unscathed-unless-susceptible","touch_damage":{"dice":8,"sides":10,"type":"force"},"control":{"stationary_until_controlled":true,"range_feet":60,"action":"magic","check_ability":"intelligence","check_skill":"arcana","dc":25,"lasts_until":"start-of-your-next-turn","other_controller":"loses-control","failure_move_feet":10},"move":{"action":"bonus","feet_per_intelligence_modifier":5,"minimum_feet":5,"save_ability":"dexterity","save_dc":19,"damage":{"dice":8,"sides":10,"type":"force"},"at_zero_hit_points":"obliterated-leaving-only-possessions"},"interaction_triggers":["planar-portal","extradimensional-space"],"table_die_sides":100,"interaction_table":[{"min":1,"max":50,"result":"sphere-destroyed"},{"min":51,"max":85,"result":"sphere-moves-through"},{"min":86,"max":100,"result":"spatial-rift-sends-all-to-random-plane"}],"rift_radius_feet":180,"source_pages":[244,245]}}`,
		`{"charming":{"attunement_by_any_of":["bard","cleric","druid","sorcerer","warlock","wizard"],"condition":"while-holding","max_charges":10,"cast_spell":{"charge_cost":1,"spells":["charm-person","command","comprehend-languages"],"save_dc_source":"your-spell-save-dc"},"reflect_enchantment":{"trigger":"succeed-on-save-against-enchantment-targeting-only-you","action":"reaction","charge_cost":1,"result":"turn-spell-back-on-caster"},"resist_enchantment":{"trigger":"fail-save-against-enchantment-targeting-only-you","result":"failed-save-becomes-success","reuse_after":"next-dawn"},"recharge":{"dice":{"dice":1,"sides":8,"type":"charges"},"bonus":2,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-to-dust"}}}`,
		`{"fire":{"attunement_by_any_of":["druid","sorcerer","warlock","wizard"],"condition":"while-holding","resistance":"fire","max_charges":10,"save_dc_source":"your-spell-save-dc","spells":[{"spell":"burning-hands","charge_cost":1},{"spell":"fireball","charge_cost":3},{"spell":"wall-of-fire","charge_cost":4}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"crumbles-into-cinders"}}}`,
		`{"frost":{"attunement_by_any_of":["druid","sorcerer","warlock","wizard"],"condition":"while-holding","resistance":"cold","max_charges":10,"save_dc_source":"your-spell-save-dc","spells":[{"spell":"cone-of-cold","charge_cost":5},{"spell":"fog-cloud","charge_cost":1},{"spell":"ice-storm","charge_cost":4},{"spell":"wall-of-ice","charge_cost":4}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"turns-to-water"}}}`,
		`{"healing":{"attunement_by_any_of":["bard","cleric","druid"],"condition":"while-holding","max_charges":10,"spell_modifier_source":"your-spellcasting-ability-modifier","cure_wounds":{"charges_per_spell_level":1,"maximum_spell_level":4},"spells":[{"spell":"lesser-restoration","charge_cost":2},{"spell":"mass-cure-wounds","charge_cost":5}],"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"vanishes-in-flash-of-light"}}}`,
		`{"power":{"attunement_by_any_of":["sorcerer","warlock","wizard"],"max_charges":20,"wielded_as":"quarterstaff","attack_bonus":2,"damage_bonus":2,"holding_bonus":2,"holding_bonus_applies_to":["armor-class","saving-throws","spell-attack-rolls"],"save_dc_source":"your-spell-save-dc","spells":[{"spell":"cone-of-cold","printed_level":"none","charge_cost":5},{"spell":"fireball","printed_level":"level-5-version","charge_cost":5},{"spell":"globe-of-invulnerability","printed_level":"none","charge_cost":6},{"spell":"hold-monster","printed_level":"none","charge_cost":5},{"spell":"levitate","printed_level":"none","charge_cost":2},{"spell":"lightning-bolt","printed_level":"level-5-version","charge_cost":5},{"spell":"magic-missile","printed_level":"none","charge_cost":1},{"spell":"ray-of-enfeeblement","printed_level":"none","charge_cost":1},{"spell":"wall-of-force","printed_level":"none","charge_cost":5}],"recharge":{"dice":{"dice":2,"sides":8,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"loses_properties_on":1,"retains":"attack-and-damage-bonus-only","regains_on":20,"regains":{"dice":1,"sides":8,"type":"charges"},"regains_bonus":2},"retributive_strike":{"action":"magic","area_emanation_feet":30,"planar_travel_percent":50,"your_damage_per_charge":16,"other_damage_per_charge":4,"damage_type":"force","save_ability":"dexterity","save_dc":17,"success_damage":"half"},"source_pages":[245,246]}}`,
		`{"striking":{"wielded_as":"quarterstaff","attack_bonus":3,"damage_bonus":3,"max_charges":10,"trigger":"hit-with-melee-attack","maximum_charges_per_hit":3,"extra_damage_per_charge":{"dice":1,"sides":6,"type":"force"},"recharge":{"dice":{"dice":1,"sides":6,"type":"charges"},"bonus":4,"when":"daily-at-dawn"},"last_charge":{"die_sides":20,"destroyed_on":1,"result":"becomes-nonmagical-quarterstaff"}}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"hold:none::item.hold defense.armor-class-bonus-by-rarity",
		"hold:none::item.hold defense.armor-class-bonus-by-rarity",
		"hold:none::item.hold defense.armor-class-bonus-by-rarity",
		"wear:none::item.attunement item.wear movement.walk-walls-and-ceilings-hands-free movement.climb-speed-equals-speed movement.not-on-slippery-surfaces",
		"none:utilize-action:apply-one-ounce:action.utilize item.permanent-bond item.1-ounce-covers-1-foot-square duration.sets-after-1-minute item.bond-broken-only-by-listed-means",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"none:cast-spell:spell-on-your-spell-list:item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook",
		"hold:none::item.attunement item.hold save.magic-advantage attack.spell-attacks-disadvantage-against-you",
		"none:magic-action:within-60-feet:effect.obliterates-matter damage.force-8d10 check.intelligence-arcana-dc-25-control action.magic action.bonus-move-sphere save.dexterity-dc-19 table.planar-interaction-1d100",
		"hold:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.hold item.charges spell.cast-charm-person spell.cast-command spell.cast-comprehend-languages action.reaction-reflect-enchantment save.resist-enchantment-once-until-dawn item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.hold defense.resistance-fire item.charges spell.cast-listed-spells spell.uses-your-save-dc item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.hold defense.resistance-cold item.charges spell.cast-listed-spells spell.uses-your-save-dc item.destroyed-on-1-after-last-charge",
		"hold:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.hold item.charges spell.cast-cure-wounds-per-level spell.cast-lesser-restoration spell.cast-mass-cure-wounds spell.uses-your-spellcasting-modifier item.destroyed-on-1-after-last-charge",
		"wield:cast-spell:expend-listed-charges:item.attunement attunement.by-class item.magic-weapon bonus.attack-and-damage-2 bonus.armor-class-saves-and-spell-attack-2 item.charges spell.cast-listed-spells spell.uses-your-save-dc item.last-charge-roll item.retributive-strike damage.force-by-charges save.dexterity-dc-17",
		"wield:weapon-attack:expend-up-to-three-charges:item.attunement item.magic-weapon bonus.attack-and-damage-3 item.charges damage.force-1d6-per-charge item.nonmagical-on-1-after-last-charge",
	}
	index := 0
	for _, f := range loadItems242(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails242
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

func TestItems242CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems242(t)
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
	if len(composed) != 216 || variants != 288 || composed[0].Name != "Potion of Healing" || composed[204].Name != "Shield, +1, +2, or +3" || composed[215].Name != "Staff of Striking" {
		t.Fatalf("expected Potion of Healing, batches #225 to #241 and the twelve families of #242: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_242") {
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
func TestItems242RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items242JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 1200 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items242JSON, &fresh); err != nil {
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
		if _, err := DecodeItems242(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems242RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails242]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails242] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-17" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 241 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/staff-of-swarming-insects"
			b.Families[12].Name = "Staff of Swarming Insects"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "weapon" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"shield-1-2-or-3 widened base", func(b *batch) { b.Families[0].Base.Options = append(b.Families[0].Base.Options, "any-weapon") }},
		{"spellguard-shield widened base", func(b *batch) { b.Families[4].Base.Options = append(b.Families[4].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{244, 244} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{243, 244} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"shield-1-2-or-3 dropped variant", func(b *batch) { b.Families[0].Variants = b.Families[0].Variants[:2] }},
		{"shield-1-2-or-3 reordered variants", func(b *batch) { s := b.Families[0].Variants; s[0], s[1] = s[1], s[0] }},
		{"spell-scroll dropped variant", func(b *batch) { b.Families[3].Variants = b.Families[3].Variants[:9] }},
		{"spell-scroll reordered variants", func(b *batch) { s := b.Families[3].Variants; s[0], s[1] = s[1], s[0] }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"shield-1-2-or-3 rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"slippers-of-spider-climbing rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"sovereign-glue rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"spell-scroll rarity", func(b *batch) { v(b, 3, 0).Rarity = "rare" }},
		{"spellguard-shield rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"sphere-of-annihilation rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"staff-of-charming rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"staff-of-fire rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"staff-of-frost rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"staff-of-healing rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"staff-of-power rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"staff-of-striking rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"shield-1-2-or-3 attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "required" }},
		{"slippers-of-spider-climbing attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"sovereign-glue attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "required" }},
		{"spell-scroll dropped gear listing", func(b *batch) { b.Families[3].Presentations = b.Families[3].Presentations[:1] }},
		{"spell-scroll only gear listing", func(b *batch) { b.Families[3].Presentations = b.Families[3].Presentations[1:] }},
		{"spell-scroll reordered listings", func(b *batch) { p := b.Families[3].Presentations; p[0], p[1] = p[1], p[0] }},
		{"spell-scroll gear page", func(b *batch) { b.Families[3].Presentations[1].Pages = []int{95, 95} }},
		{"gear listing on another family", func(b *batch) {
			f := &b.Families[0]
			f.Presentations = append(f.Presentations, ItemReference{ID: "srd-5.2.1:content/adventuring-gear/shield-1-2-or-3", Pages: []int{99, 99}})
		}},
		{"spell-scroll dropped price", func(b *batch) {
			*v(b, 3, 0).Cost = ItemCost{Kind: "unspecified", Currency: "none", SourcePages: []int{244, 244}}
		}},
		{"spell-scroll swapped price", func(b *batch) { v(b, 3, 0).Cost.Amount = 50 }},
		{"spell-scroll price currency", func(b *batch) { v(b, 3, 1).Cost.Currency = "sp" }},
		{"spell-scroll price page", func(b *batch) { v(b, 3, 1).Cost.SourcePages = []int{244, 244} }},
		{"spell-scroll invented level 2 price", func(b *batch) {
			*v(b, 3, 2).Cost = ItemCost{Kind: "listed", Amount: 100, Currency: "gp", SourcePages: []int{99, 99}}
		}},
		{"spell-scroll attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "required" }},
		{"spellguard-shield attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "none" }},
		{"sphere-of-annihilation attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "required" }},
		{"staff-of-charming attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "none" }},
		{"staff-of-fire attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"staff-of-frost attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "none" }},
		{"staff-of-healing attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"staff-of-power attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"staff-of-striking attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 10 }},
		{"staff-of-charming dropped charges", func(b *batch) { *v(b, 6, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-charming charge maximum", func(b *batch) { v(b, 6, 0).Charges.Maximum++ }},
		{"staff-of-charming charge recharge", func(b *batch) { v(b, 6, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-fire dropped charges", func(b *batch) { *v(b, 7, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-fire charge maximum", func(b *batch) { v(b, 7, 0).Charges.Maximum++ }},
		{"staff-of-fire charge recharge", func(b *batch) { v(b, 7, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-frost dropped charges", func(b *batch) { *v(b, 8, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-frost charge maximum", func(b *batch) { v(b, 8, 0).Charges.Maximum++ }},
		{"staff-of-frost charge recharge", func(b *batch) { v(b, 8, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-healing dropped charges", func(b *batch) { *v(b, 9, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-healing charge maximum", func(b *batch) { v(b, 9, 0).Charges.Maximum++ }},
		{"staff-of-healing charge recharge", func(b *batch) { v(b, 9, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-power dropped charges", func(b *batch) { *v(b, 10, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-power charge maximum", func(b *batch) { v(b, 10, 0).Charges.Maximum++ }},
		{"staff-of-power charge recharge", func(b *batch) { v(b, 10, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"staff-of-striking dropped charges", func(b *batch) { *v(b, 11, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"staff-of-striking charge maximum", func(b *batch) { v(b, 11, 0).Charges.Maximum++ }},
		{"staff-of-striking charge recharge", func(b *batch) { v(b, 11, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "reaction" }},
		{"dropped use restriction", func(b *batch) { v(b, 2, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails242{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.SpiderSlippers = v(b, 1, 0).Details.SpiderSlippers }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.ShieldBonus.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.ShieldBonus.Condition = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.ShieldBonus.ArmorClassBonus = -1 }},
		{"shield bonus swapped", func(b *batch) { p := v(b, 0, 0).Details.ShieldBonus; p.ArmorClassBonus = 2 }},
		{"shield stacking", func(b *batch) { p := v(b, 0, 0).Details.ShieldBonus; p.StacksWithShieldBonus = false }},
		{"slippers hands", func(b *batch) { p := v(b, 1, 0).Details.SpiderSlippers; p.HandsFree = false }},
		{"slippers speed", func(b *batch) { p := v(b, 1, 0).Details.SpiderSlippers; p.ClimbSpeed = "half-speed" }},
		{"glue dice", func(b *batch) { p := v(b, 2, 0).Details.Glue; p.FoundOunces.Sides = 4 }},
		{"glue set time", func(b *batch) { p := v(b, 2, 0).Details.Glue; p.SetMinutes = 10 }},
		{"glue solvents", func(b *batch) { p := v(b, 2, 0).Details.Glue; p.BrokenOnlyBy = p.BrokenOnlyBy[:2] }},
		{"scroll level", func(b *batch) { p := v(b, 3, 0).Details.Scroll; p.SpellLevel = 1 }},
		{"scroll dc", func(b *batch) { p := v(b, 3, 0).Details.Scroll; p.SaveDC = 15 }},
		{"scroll bonus", func(b *batch) { p := v(b, 3, 0).Details.Scroll; p.AttackBonus = 7 }},
		{"scroll rarity row", func(b *batch) { p := v(b, 3, 0).Details.Scroll; p.Rarity = "rare" }},
		{"scroll copy dc", func(b *batch) { p := v(b, 3, 0).Details.Scroll; p.Copying.DCBase = 12 }},
		{"spellguard saves", func(b *batch) { p := v(b, 4, 0).Details.Spellguard; p.SaveMode = "normal" }},
		{"spellguard attacks", func(b *batch) { p := v(b, 4, 0).Details.Spellguard; p.SpellAttacksAgainstYou = "normal" }},
		{"sphere control dc", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.Control.DC = 20 }},
		{"sphere save dc", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.Move.SaveDC = 17 }},
		{"sphere damage", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.TouchDamage.Dice = 4 }},
		{"sphere table gap", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.InteractionTable[1].Min = 52 }},
		{"sphere table end", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.InteractionTable[2].Max = 99 }},
		{"sphere pages", func(b *batch) { p := v(b, 5, 0).Details.Annihilation; p.SourcePages = []int{244} }},
		{"charming charges", func(b *batch) { p := v(b, 6, 0).Details.Charming; p.MaxCharges = 7 }},
		{"charming recharge", func(b *batch) { p := v(b, 6, 0).Details.Charming; p.Recharge.Bonus = 4 }},
		{"charming class", func(b *batch) { p := v(b, 6, 0).Details.Charming; p.AttunementByAnyOf = p.AttunementByAnyOf[:5] }},
		{"charming duplicate spell", func(b *batch) { p := v(b, 6, 0).Details.Charming; p.CastSpell.Spells[1] = "charm-person" }},
		{"fire charges", func(b *batch) { p := v(b, 7, 0).Details.Fire; p.MaxCharges = 20 }},
		{"fire spell cost", func(b *batch) { p := v(b, 7, 0).Details.Fire; p.Spells[1].ChargeCost = 5 }},
		{"fire resistance", func(b *batch) { p := v(b, 7, 0).Details.Fire; p.Resistance = "cold" }},
		{"frost charges", func(b *batch) { p := v(b, 8, 0).Details.Frost; p.MaxCharges = 6 }},
		{"frost spell cost", func(b *batch) { p := v(b, 8, 0).Details.Frost; p.Spells[0].ChargeCost = 4 }},
		{"frost missing spell", func(b *batch) { p := v(b, 8, 0).Details.Frost; p.Spells = p.Spells[:3] }},
		{"healing staff maximum", func(b *batch) { p := v(b, 9, 0).Details.Healing; p.CureWounds.MaximumSpellLevel = 5 }},
		{"healing staff cost", func(b *batch) { p := v(b, 9, 0).Details.Healing; p.Spells[1].ChargeCost = 4 }},
		{"healing staff charges", func(b *batch) { p := v(b, 9, 0).Details.Healing; p.MaxCharges = 8 }},
		{"power charges", func(b *batch) { p := v(b, 10, 0).Details.Power; p.MaxCharges = 50 }},
		{"power bonus", func(b *batch) { p := v(b, 10, 0).Details.Power; p.AttackBonus = 3 }},
		{"power spell cost", func(b *batch) { p := v(b, 10, 0).Details.Power; p.Spells[2].ChargeCost = 5 }},
		{"power strike multiplier", func(b *batch) { p := v(b, 10, 0).Details.Power; p.RetributiveStrike.OtherDamagePerCharge = 6 }},
		{"power strike dc", func(b *batch) { p := v(b, 10, 0).Details.Power; p.RetributiveStrike.SaveDC = 15 }},
		{"power pages", func(b *batch) { p := v(b, 10, 0).Details.Power; p.SourcePages = []int{245} }},
		{"striking bonus", func(b *batch) { p := v(b, 11, 0).Details.Striking; p.AttackBonus = 2 }},
		{"striking per hit", func(b *batch) { p := v(b, 11, 0).Details.Striking; p.MaximumChargesPerHit = 5 }},
		{"striking dice", func(b *batch) { p := v(b, 11, 0).Details.Striking; p.ExtraDamagePerCharge.Sides = 8 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems242(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems242(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems242RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/shield-1-2-or-3/variant/plus-1"`)
	if !bytes.Contains(items242JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items242JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems242(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"shield_bonus": {`, `"shield_bonus_plus": {`},
		"unknown detail":     {`"armor_class_bonus": 1`, `"armor_class_bonus_all": 1`},
		"string for number":  {`"armor_class_bonus": 1`, `"armor_class_bonus": "1"`},
		"fractional number":  {`"armor_class_bonus": 1`, `"armor_class_bonus": 1.5`},
		"number for boolean": {`"stacks_with_shield_bonus": true`, `"stacks_with_shield_bonus": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items242JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items242JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems242(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items242JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems242(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items242JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems241(items242JSON); err == nil {
		t.Fatal("#241 decoder accepted the #242 batch")
	}
	if _, err := DecodeItems242(items241JSON); err == nil {
		t.Fatal("#242 decoder accepted the #241 batch")
	}
}
