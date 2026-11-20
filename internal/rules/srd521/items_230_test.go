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

func loadItems230(t *testing.T) ItemSourceBatch[ItemDetails230] {
	t.Helper()
	b, err := LoadItems230()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 218-220, in planning inventory order for issue #230. Variants
// are the rows of the source's own dragon (p. 219) and gem (p. 220) tables.
func TestItems230SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug, category string
		page                 int
		options              []string
		attunement, rarity   string
		variants             []string
	}{
		{"Dragon Scale Mail", "dragon-scale-mail", "armor", 218, []string{"scale-mail"}, "required", "very-rare", []string{"black", "blue", "brass", "bronze", "copper", "gold", "green", "red", "silver", "white"}},
		{"Dragon Slayer", "dragon-slayer", "weapon", 219, []string{"any-simple-weapon", "any-martial-weapon"}, "none", "rare", []string{"standard"}},
		{"Dust of Disappearance", "dust-of-disappearance", "wondrous-item", 219, nil, "none", "uncommon", []string{"standard"}},
		{"Dust of Dryness", "dust-of-dryness", "wondrous-item", 219, nil, "none", "uncommon", []string{"standard"}},
		{"Dust of Sneezing and Choking", "dust-of-sneezing-and-choking", "wondrous-item", 219, nil, "none", "uncommon", []string{"standard"}},
		{"Dwarven Plate", "dwarven-plate", "armor", 219, []string{"half-plate-armor", "plate-armor"}, "none", "very-rare", []string{"standard"}},
		{"Dwarven Thrower", "dwarven-thrower", "weapon", 219, []string{"warhammer"}, "required", "very-rare", []string{"standard"}},
		{"Efficient Quiver", "efficient-quiver", "wondrous-item", 219, nil, "none", "uncommon", []string{"standard"}},
		{"Efreeti Bottle", "efreeti-bottle", "wondrous-item", 220, nil, "none", "very-rare", []string{"standard"}},
		{"Elemental Gem", "elemental-gem", "wondrous-item", 220, nil, "none", "uncommon", []string{"blue-sapphire", "emerald", "red-corundum", "yellow-diamond"}},
		{"Elixir of Health", "elixir-of-health", "potion", 220, nil, "none", "rare", []string{"standard"}},
		{"Elven Chain", "elven-chain", "armor", 220, []string{"chain-mail", "chain-shirt"}, "none", "rare", []string{"standard"}},
	}
	b := loadItems230(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-06" || b.CatalogueIssue != 230 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #230")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 218–220") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	total := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || len(f.Base.Excluded) != 0 || len(f.Base.Options) != len(want.options) || (len(want.options) > 0 && !slices.Equal(f.Base.Options, want.options)) {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			got := []string{}
			for _, v := range f.Variants {
				got = append(got, strings.TrimPrefix(v.ID, id+"/variant/"))
				if v.Rarity != want.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				// None of the twelve lists charges (the Dust of Dryness pinches are a quantity).
				if *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
					t.Fatalf("charges mismatch: %+v", *v.Charges)
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
			}
			if !slices.Equal(got, want.variants) {
				t.Fatalf("source-listed variants %v, want %v", got, want.variants)
			}
			total += len(got)
		})
	}
	if total != 24 {
		t.Fatalf("expected 24 variants, got %d", total)
	}
	names := map[string]string{
		"Dragon Scale Mail (black)": "black", "Elemental Gem (blue sapphire)": "blue-sapphire", "Elemental Gem (yellow diamond)": "yellow-diamond",
	}
	for _, f := range b.Families {
		for _, v := range f.Variants {
			if want, ok := names[v.Name]; ok && !strings.HasSuffix(v.ID, "/variant/"+want) {
				t.Fatalf("variant name %q does not match its identity %q", v.Name, v.ID)
			}
		}
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
		if r.Issue != 230 {
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
		t.Fatalf("inventory assigns %d families to issue #230", index)
	}
	b.Families[0].Variants[0].Details.ScaleMail.Condition = "changed"
	if fresh := loadItems230(t); fresh.Families[0].Variants[0].Details.ScaleMail.Condition != "while-wearing" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all twenty-four variants, transcribed from pinned
// PDF pp. 218-220 (Dragon Scale Mail and Dust of Dryness continue across a page
// or column break; the dragon and gem rows come from the source's tables).
func TestItems230SourceValues(t *testing.T) {
	expected := []string{
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"black","resistance":"acid","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"blue","resistance":"lightning","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"brass","resistance":"fire","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"bronze","resistance":"lightning","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"copper","resistance":"acid","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"gold","resistance":"fire","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"green","resistance":"poison","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"red","resistance":"fire","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"silver","resistance":"cold","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"scale_mail":{"condition":"while-wearing","ac_bonus":1,"breath_weapon_save_mode":"advantage","breath_weapon_source":"dragons","dragon":"white","resistance":"cold","sense_action":"magic","sense_target":"closest-dragon-of-same-type","sense_range_miles":30,"sense_reveals":["distance","direction"],"sense_uses":1,"sense_recharge":"next-dawn"}}`,
		`{"slayer":{"attack_bonus":1,"damage_bonus":1,"extra_damage":{"dice":3,"sides":6,"type":"weapon-type"},"extra_damage_target":"dragon"}}`,
		`{"disappearance":{"appearance":"fine-sand","uses":1,"action":"utilize","area_emanation_feet":10,"affects":["you","each-creature","each-object"],"grants_condition":"invisible","duration":{"dice":2,"sides":4,"type":"minutes"},"duration_same_for_all":true,"consumed_when":"magic-takes-effect","ends_for_creature_after":["makes-attack-roll","deals-damage","casts-spell"]}}`,
		`{"dryness":{"container":"small-packet","pinches":{"dice":1,"sides":6,"type":"pinches"},"pinches_bonus":4,"water_action":"utilize","water_cube_feet":15,"water_result":"marble-sized-pellet","pellet_location":"floats-or-rests-near-where-sprinkled","pellet_weight":"negligible","smash_action":"utilize","smash_surface":"hard-surface","smash_result":"pellet-shatters-and-releases-absorbed-water","smash_ends":"destroys-pellet-and-ends-its-magic","elemental_action":"utilize","elemental_range_feet":5,"elemental_target":"elemental-composed-mostly-of-water","save_ability":"constitution","save_dc":13,"damage":{"dice":10,"sides":6,"type":"necrotic"},"save_success":"half-damage"}}`,
		`{"sneezing":{"container":"small-container","resembles":"dust-of-disappearance","identify_reveals":true,"uses":1,"action":"utilize","area_emanation_feet":30,"affects_user":true,"save_ability":"constitution","save_dc":15,"auto_succeed_creature_types":["constructs","elementals","oozes","plants","undead"],"failure_effect":"sneezes-uncontrollably","failure_conditions":["incapacitated","suffocating"],"repeat_save":"end-of-each-of-its-turns","ends_on":["successful-save","lesser-restoration"]}}`,
		`{"plate":{"condition":"while-wearing","ac_bonus":2,"reaction_trigger":"effect-moves-you-against-your-will-along-ground","reaction_effect":"reduce-distance-moved","reduce_distance_up_to_feet":10}}`,
		`{"thrower":{"attunement_by_any_of":["dwarf","creature-attuned-to-belt-of-dwarvenkind"],"attack_bonus":3,"damage_bonus":3,"thrown_property":true,"normal_range_feet":20,"long_range_feet":60,"ranged_hit_extra_damage":{"dice":1,"sides":8,"type":"force"},"giant_target_extra_damage":{"dice":2,"sides":8,"type":"force"},"giant_target":"giant","returns_to":"your-hand","returns_when":"immediately-after-hitting-or-missing"}}`,
		`{"quiver":{"compartments":[{"length":"shortest","capacity":60,"holds":["arrows","bolts","similar-objects"]},{"length":"midsize","capacity":18,"holds":["javelins","similar-objects"]},{"length":"longest","capacity":6,"holds":["long-objects","bows","quarterstaffs","spears"]}],"extradimensional":true,"max_weight_pounds":2,"draw_as":"regular-quiver-or-scabbard"}}`,
		`{"bottle":{"appearance":"painted-brass-bottle","action":"magic","smoke_ends":"end-of-your-turn","smoke_end_effect":"flash-of-harmless-fire","summons":"efreeti","appears_within_feet":30,"table_applies_to":"first-opening","table_die_sides":10,"effects":[{"min":1,"max":1,"effect":"efreeti-attacks-you","understands_owner_languages":false,"obeys_owner_commands":false,"fight_rounds":5,"duration_seconds":0,"stopper_locked_hours":0,"same_effect_next_openings":0,"escapes_on_opening":0,"wishes":0,"ends":["after-fighting-rounds"],"bottle_after":"bottle-loses-magic"},{"min":2,"max":9,"effect":"efreeti-serves-you","understands_owner_languages":true,"obeys_owner_commands":true,"fight_rounds":0,"duration_seconds":3600,"stopper_locked_hours":24,"same_effect_next_openings":2,"escapes_on_opening":4,"wishes":0,"ends":["duration-elapsed"],"bottle_after":"returns-to-bottle-with-new-stopper"},{"min":10,"max":10,"effect":"efreeti-casts-wish-once","understands_owner_languages":true,"obeys_owner_commands":false,"fight_rounds":0,"duration_seconds":3600,"stopper_locked_hours":0,"same_effect_next_openings":0,"escapes_on_opening":0,"wishes":1,"ends":["grants-wish","duration-elapsed"],"bottle_after":"bottle-loses-magic"}]}}`,
		`{"gem":{"gem":"blue-sapphire","elemental":"air-elemental","action":"utilize","gem_after":"ceases-to-be-magical","location":"unoccupied-space-nearest-broken-gem","understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-summoner","duration_seconds":3600,"end_triggers":["duration-elapsed","dies","dismissed-as-bonus-action"]}}`,
		`{"gem":{"gem":"emerald","elemental":"water-elemental","action":"utilize","gem_after":"ceases-to-be-magical","location":"unoccupied-space-nearest-broken-gem","understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-summoner","duration_seconds":3600,"end_triggers":["duration-elapsed","dies","dismissed-as-bonus-action"]}}`,
		`{"gem":{"gem":"red-corundum","elemental":"fire-elemental","action":"utilize","gem_after":"ceases-to-be-magical","location":"unoccupied-space-nearest-broken-gem","understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-summoner","duration_seconds":3600,"end_triggers":["duration-elapsed","dies","dismissed-as-bonus-action"]}}`,
		`{"gem":{"gem":"yellow-diamond","elemental":"earth-elemental","action":"utilize","gem_after":"ceases-to-be-magical","location":"unoccupied-space-nearest-broken-gem","understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-summoner","duration_seconds":3600,"end_triggers":["duration-elapsed","dies","dismissed-as-bonus-action"]}}`,
		`{"elixir":{"cures_magical_contagions":true,"ends_conditions":["blinded","deafened","paralyzed","poisoned"],"appearance":"clear-red-liquid-with-tiny-bubbles-of-light"}}`,
		`{"chain":{"condition":"while-wearing","ac_bonus":1,"trained_with_this_armor":true,"trained_even_without":["medium-armor","heavy-armor"]}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn",
		"wield:weapon-attack::item.magic-weapon attack.weapon-bonus damage.weapon-bonus damage.extra-3d6-against-dragons",
		"none:utilize-action:throw-into-air one-use:action.utilize area.emanation-10-feet condition.invisible duration.2d4-minutes effect.ends-after-attack-damage-or-spell item.consumed-on-use",
		"none:utilize-action:sprinkle-one-pinch:item.pinches-1d6-plus-4 action.utilize effect.water-cube-15-feet-to-pellet action.utilize-smash-pellet effect.pellet-releases-water save.constitution-dc-13 damage.necrotic-10d6 damage.half-on-success",
		"none:utilize-action:throw-into-air one-use:action.utilize area.emanation-30-feet save.constitution-dc-15 save.automatic-success-by-creature-type condition.incapacitated condition.suffocating save.repeat-end-of-turn effect.ends-with-lesser-restoration",
		"wear:reaction:moved-against-will-along-ground:item.wear defense.armor-class-bonus-2 action.reaction movement.reduce-forced-ground-movement-10-feet",
		"wield:weapon-attack::item.attunement item.magic-weapon attack.weapon-bonus damage.weapon-bonus property.thrown-20-60 damage.extra-force-1d8-ranged-hit damage.extra-force-2d8-against-giants item.returns-to-hand",
		"none:container::item.extradimensional-storage item.weight-cap-2-pounds item.compartments-60-18-6 action.draw-as-regular-quiver",
		"none:magic-action:remove-stopper:action.magic summon.efreeti table.efreeti-1d10-first-opening effect.bottle-loses-magic",
		"none:utilize-action:break-gem gem-ceases-to-be-magical:action.utilize summon.elemental-by-gem-type duration.one-hour effect.gem-ceases-to-be-magical",
		"none:utilize-action:break-gem gem-ceases-to-be-magical:action.utilize summon.elemental-by-gem-type duration.one-hour effect.gem-ceases-to-be-magical",
		"none:utilize-action:break-gem gem-ceases-to-be-magical:action.utilize summon.elemental-by-gem-type duration.one-hour effect.gem-ceases-to-be-magical",
		"none:utilize-action:break-gem gem-ceases-to-be-magical:action.utilize summon.elemental-by-gem-type duration.one-hour effect.gem-ceases-to-be-magical",
		"none:consume:drink-one:item.consume cure.magical-contagions condition.end-blinded condition.end-deafened condition.end-paralyzed condition.end-poisoned",
		"wear:none::item.wear defense.armor-class-bonus-1 training.trained-with-this-armor",
	}
	index := 0
	for _, f := range loadItems230(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails230
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

func TestItems230CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems230(t)
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
	fifth, err := LoadItems229()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), first.SourceIdentities(), second.SourceIdentities(), third.SourceIdentities(), fourth.SourceIdentities(), fifth.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 73 || variants != 98 || composed[0].Name != "Potion of Healing" || composed[61].Name != "Dragon Scale Mail" || composed[72].Name != "Elven Chain" {
		t.Fatal("expected Potion of Healing, the #225-#229 families and exactly the twelve families of #230")
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
			if strings.Contains(fixture, "items_230") {
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
func TestItems230RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items230JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 1000 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items230JSON, &fresh); err != nil {
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
		if _, err := DecodeItems230(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems230RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails230]
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails230] {
		return &b.Families[family].Variants[variant]
	}
	// Family indexes follow inventory order; see TestItems230SourceIdentities.
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-05" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 229 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/energy-bow"
			b.Families[12].Name = "Energy Bow"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/dragon-hide-mail" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Dragon Mail" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"scale-mail", "plate-armor"} }},
		{"dropped base", func(b *batch) { b.Families[5].Base.Options = []string{"plate-armor"} }},
		{"excluded base", func(b *batch) { b.Families[11].Base.Excluded = []string{"chain-shirt"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{219, 219} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{218, 219} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"missing dragon variant", func(b *batch) { b.Families[0].Variants = b.Families[0].Variants[:9] }},
		{"invented dragon variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"reordered dragon variants", func(b *batch) {
			vs := b.Families[0].Variants
			vs[0], vs[1] = vs[1], vs[0]
		}},
		{"missing gem variant", func(b *batch) { b.Families[9].Variants = b.Families[9].Variants[:3] }},
		{"invented slayer variant", func(b *batch) { b.Families[1].Variants = append(b.Families[1].Variants, b.Families[1].Variants[0]) }},
		{"no variants", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Dragon Scale Mail" }},
		{"rarity", func(b *batch) { v(b, 0, 0).Rarity = "rare" }},
		{"slayer rarity", func(b *batch) { v(b, 1, 0).Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "required" }},
		{"thrower attunement dropped", func(b *batch) { v(b, 6, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 6, 0).Attunement.Restrictions = []string{"dwarf"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) {
			*v(b, 3, 0).Charges = ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "none"}
		}},
		{"pinches are not charges", func(b *batch) { v(b, 3, 0).Charges.Maximum = 5 }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "bonus-action" }},
		{"dropped use restriction", func(b *batch) { v(b, 0, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails230{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Slayer = v(b, 1, 0).Details.Slayer }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.ScaleMail.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.ScaleMail.Condition = "While wearing this armor" }},
		{"scale mail bonus", func(b *batch) { v(b, 0, 0).Details.ScaleMail.ACBonus = 0 }},
		{"scale mail range", func(b *batch) { v(b, 0, 0).Details.ScaleMail.SenseRangeMiles = 0 }},
		{"black resists lightning", func(b *batch) { v(b, 0, 0).Details.ScaleMail.Resistance = "lightning" }},
		{"gold resists cold", func(b *batch) { v(b, 0, 5).Details.ScaleMail.Resistance = "cold" }},
		{"dragon disagrees with variant", func(b *batch) { v(b, 0, 1).Details.ScaleMail.Dragon = "black" }},
		{"duplicate sense", func(b *batch) { v(b, 0, 0).Details.ScaleMail.SenseReveals[1] = "distance" }},
		{"slayer bonus", func(b *batch) { v(b, 1, 0).Details.Slayer.AttackBonus = 0 }},
		{"slayer dice", func(b *batch) { v(b, 1, 0).Details.Slayer.ExtraDamage.Sides = 0 }},
		{"disappearance duration", func(b *batch) { v(b, 2, 0).Details.Disappearance.Duration.Dice = 0 }},
		{"disappearance area", func(b *batch) { v(b, 2, 0).Details.Disappearance.AreaEmanationFeet = 0 }},
		{"disappearance duplicate end", func(b *batch) { v(b, 2, 0).Details.Disappearance.EndsForCreatureOn[1] = "makes-attack-roll" }},
		{"dryness pinches", func(b *batch) { v(b, 3, 0).Details.Dryness.Pinches.Dice = 0 }},
		{"dryness bonus", func(b *batch) { v(b, 3, 0).Details.Dryness.PinchesBonus = 0 }},
		{"dryness dc", func(b *batch) { v(b, 3, 0).Details.Dryness.SaveDC = 0 }},
		{"dryness damage", func(b *batch) { v(b, 3, 0).Details.Dryness.Damage.Dice = 0 }},
		{"sneezing dc", func(b *batch) { v(b, 4, 0).Details.Sneezing.SaveDC = 0 }},
		{"sneezing creature types", func(b *batch) { v(b, 4, 0).Details.Sneezing.AutoSucceed = []string{} }},
		{"sneezing duplicate condition", func(b *batch) { v(b, 4, 0).Details.Sneezing.FailureConditions[1] = "incapacitated" }},
		{"plate bonus", func(b *batch) { v(b, 5, 0).Details.Plate.ACBonus = 0 }},
		{"plate distance", func(b *batch) { v(b, 5, 0).Details.Plate.ReduceDistanceMax = 0 }},
		{"thrower range order", func(b *batch) { v(b, 6, 0).Details.Thrower.LongRangeFeet = 20 }},
		{"thrower giant dice", func(b *batch) { v(b, 6, 0).Details.Thrower.GiantExtraDamage.Dice = 0 }},
		{"thrower attunement list", func(b *batch) { v(b, 6, 0).Details.Thrower.AttunementBy = []string{} }},
		{"quiver compartment dropped", func(b *batch) { v(b, 7, 0).Details.Quiver.Compartments = v(b, 7, 0).Details.Quiver.Compartments[:2] }},
		{"quiver capacity", func(b *batch) { v(b, 7, 0).Details.Quiver.Compartments[0].Capacity = 0 }},
		{"quiver weight", func(b *batch) { v(b, 7, 0).Details.Quiver.MaxWeightPounds = 0 }},
		{"bottle table gap", func(b *batch) { v(b, 8, 0).Details.Bottle.Effects[1].Min = 3 }},
		{"bottle table die", func(b *batch) { v(b, 8, 0).Details.Bottle.TableDieSides = 20 }},
		{"bottle dropped row", func(b *batch) { v(b, 8, 0).Details.Bottle.Effects = v(b, 8, 0).Details.Bottle.Effects[:2] }},
		{"bottle range", func(b *batch) { v(b, 8, 0).Details.Bottle.AppearsRange = 0 }},
		{"bottle row ends", func(b *batch) { v(b, 8, 0).Details.Bottle.Effects[0].Ends = []string{} }},
		{"gem disagrees with variant", func(b *batch) { v(b, 9, 0).Details.Gem.Gem = "emerald" }},
		{"sapphire summons water", func(b *batch) { v(b, 9, 0).Details.Gem.Elemental = "water-elemental" }},
		{"diamond summons fire", func(b *batch) { v(b, 9, 3).Details.Gem.Elemental = "fire-elemental" }},
		{"gem duration", func(b *batch) { v(b, 9, 0).Details.Gem.DurationSeconds = 0 }},
		{"gem duplicate trigger", func(b *batch) {
			v(b, 9, 0).Details.Gem.EndTriggers[1] = "dies"
			v(b, 9, 0).Details.Gem.EndTriggers[2] = "dies"
		}},
		{"elixir conditions", func(b *batch) { v(b, 10, 0).Details.Elixir.EndsConditions = []string{} }},
		{"elixir duplicate condition", func(b *batch) { v(b, 10, 0).Details.Elixir.EndsConditions[1] = "blinded" }},
		{"chain bonus", func(b *batch) { v(b, 11, 0).Details.Chain.ACBonus = 0 }},
		{"chain training", func(b *batch) { v(b, 11, 0).Details.Chain.TrainedEvenWithout = nil }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems230(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems230(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems230RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/dragon-slayer/variant/standard"`)
	if !bytes.Contains(items230JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items230JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems230(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"slayer": {`, `"slayer_plus": {`},
		"unknown detail":     {`"extra_damage_target"`, `"extra_damage_targets"`},
		"string for number":  {`"save_dc": 13`, `"save_dc": "13"`},
		"fractional number":  {`"save_dc": 13`, `"save_dc": 13.5`},
		"number for boolean": {`"thrown_property": true`, `"thrown_property": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items230JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items230JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems230(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items230JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems230(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items230JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems228(items230JSON); err == nil {
		t.Fatal("#228 decoder accepted the #230 batch")
	}
	if _, err := DecodeItems230(items228JSON); err == nil {
		t.Fatal("#230 decoder accepted the #228 batch")
	}
}
