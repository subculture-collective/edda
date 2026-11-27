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

func loadItems241(t *testing.T) ItemSourceBatch[ItemDetails241] {
	t.Helper()
	b, err := LoadItems241()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 241-243, in planning inventory order for issue #241.
func TestItems241SourceIdentities(t *testing.T) {
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
		{"Rod of Lordly Might", "rod-of-lordly-might", 241, "rod", nil, "required", noCharges, []variant{{"standard", "Rod of Lordly Might", "legendary"}}},
		{"Rod of Resurrection", "rod-of-resurrection", 242, "rod", nil, "required", ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1-at-dawn"}, []variant{{"standard", "Rod of Resurrection", "legendary"}}},
		{"Rod of Rulership", "rod-of-rulership", 242, "rod", nil, "required", noCharges, []variant{{"standard", "Rod of Rulership", "rare"}}},
		{"Rod of Security", "rod-of-security", 242, "rod", nil, "none", noCharges, []variant{{"standard", "Rod of Security", "very-rare"}}},
		{"Rope of Climbing", "rope-of-climbing", 242, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Rope of Climbing", "uncommon"}}},
		{"Rope of Entanglement", "rope-of-entanglement", 242, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Rope of Entanglement", "rare"}}},
		{"Scarab of Protection", "scarab-of-protection", 243, "wondrous-item", nil, "required", ItemCharges{Mode: "limited", Maximum: 12, PerUse: 1, Recharge: "none"}, []variant{{"standard", "Scarab of Protection", "legendary"}}},
		{"Scimitar of Speed", "scimitar-of-speed", 243, "weapon", []string{"scimitar"}, "required", noCharges, []variant{{"standard", "Scimitar of Speed", "very-rare"}}},
		{"Sending Stones", "sending-stones", 243, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Sending Stones", "uncommon"}}},
		{"Sentinel Shield", "sentinel-shield", 243, "armor", []string{"shield"}, "none", noCharges, []variant{{"standard", "Sentinel Shield", "uncommon"}}},
		{"Shield of Missile Attraction", "shield-of-missile-attraction", 243, "armor", []string{"shield"}, "required", noCharges, []variant{{"standard", "Shield of Missile Attraction", "rare"}}},
		{"Shield of the Cavalier", "shield-of-the-cavalier", 243, "armor", []string{"shield"}, "required", noCharges, []variant{{"standard", "Shield of the Cavalier", "very-rare"}}},
	}
	b := loadItems241(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-17" || b.CatalogueIssue != 241 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #241")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 241–243") {
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
		if r.Issue != 241 {
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
		t.Fatalf("inventory assigns %d families to issue #241", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "241, 242, 243" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.LordlyMight.FunctionsAs = "changed"
	if fresh := loadItems241(t); fresh.Families[0].Variants[0].Details.LordlyMight.FunctionsAs != "mace" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 241-243.
func TestItems241SourceValues(t *testing.T) {
	expected := []string{
		`{"lordly_might":{"functions_as":"mace","attack_bonus":3,"damage_bonus":3,"button_count":6,"buttons":{"press_action":"bonus","effect_lasts_until":["different-button-pushed","same-button-pushed-again"],"fiery_blade":{"button":1,"bright_light_radius_feet":40,"dim_light_additional_feet":40,"functions_as":["longsword","shortsword"],"extra_damage":{"dice":2,"sides":6,"type":"fire"}},"battleaxe":{"button":2,"functions_as":"battleaxe","attack_and_damage_bonus":3},"spear":{"button":3,"functions_as":"spear","haft_feet":6,"attack_and_damage_bonus":3},"climbing_pole":{"button":4,"maximum_length_feet":50,"buttons_stay_in_reach":true,"anchors":["bottom-spike","three-top-hooks"],"anchor_surface":"as-hard-as-granite","bar_length_inches":3,"bar_spacing_feet":1,"capacity_pounds":4000,"reverts_on":["more-weight","lack-of-solid-anchoring"]},"battering_ram":{"button":5,"check_bonus":10,"check_ability":"strength","check_skill":"athletics","check_subjects":["doors","barricades","other-barriers"]},"normal_form":{"button":6,"indicates":"magnetic-north","without_magnetic_north":"nothing-happens","knowledge":["approximate-depth-below-ground","height-above-ground"]}},"drain_life":{"trigger":"hit-creature-with-melee-attack","save_ability":"constitution","save_dc":17,"extra_damage":{"dice":4,"sides":6,"type":"necrotic"},"you_regain":"half-of-that-necrotic-damage","reuse_after":"next-dawn"},"paralyze":{"trigger":"hit-creature-with-melee-attack","save_ability":"constitution","save_dc":17,"condition":"paralyzed","duration_minutes":1,"repeat_save":"end-of-each-of-its-turns","reuse_after":"next-dawn"},"terrify":{"condition":"frightened","action":"magic","range_feet":30,"save_ability":"wisdom","save_dc":17,"duration_minutes":1,"repeat_save":"end-of-each-of-its-turns","repeat_success":"ends-effect-on-itself","reuse_after":"next-dawn"},"source_pages":[241,242]}}`,
		`{"resurrection":{"max_charges":5,"condition":"while-holding","spells":[{"spell":"heal","charge_cost":1},{"spell":"resurrection","charge_cost":5}],"recharge_charges":1,"recharge_when":"daily-at-dawn","last_charge":{"die_sides":20,"destroyed_on":1,"result":"disappears-in-harmless-burst-of-radiance"}}}`,
		`{"rulership":{"action":"magic","targets":"each-creature-of-your-choice-you-can-see","range_feet":120,"save_ability":"wisdom","save_dc":15,"condition":"charmed","duration_hours":8,"regards_you_as":"trusted-leader","ends_early_when":["harmed-by-you-or-allies","commanded-against-its-nature"],"reuse_after":"next-dawn"}}`,
		`{"security":{"condition":"while-holding","action":"magic","travelers":{"includes_you":true,"other_willing_creatures":199},"destination":"demiplane","form_chosen_by":"you","demiplane":{"provides":"enough-water-and-food","environment_can_harm":false,"objects_exist_only_there":true,"aging":false,"time_passes_normally":true},"healing":{"per":"each-hour-spent-in-demiplane","as_if_hit_point_dice_spent":1},"stay":{"base_days":200,"divided_by":"creatures-present","rounding":"down"},"ends_when":["time-runs-out","magic-action-to-end"],"return_to":"original-or-nearest-unoccupied-space","reuse_after_days":10}}`,
		`{"climbing_rope":{"length_feet":60,"capacity_pounds":3000,"condition":"while-holding-one-end","action":"magic","move_limit":"rope-length-away-from-you","move_feet_per_turn":10,"moves_until":["reaching-destination","told-to-stop"],"other_commands":["fasten","unfasten","knot","unknot","coil"],"knotted":{"knot_spacing_feet":1,"length_feet":50,"check_mode":"advantage","check_subject":"climb-using-the-rope"},"armor_class":20,"hit_points":20,"immunities":["poison","psychic"],"regains_hit_points":1,"regain_interval_minutes":5,"regain_needs_hit_points":1,"destroyed_at_hit_points":0}}`,
		`{"entanglement":{"length_feet":30,"condition":"while-holding-one-end","action":"magic","targets":1,"range_feet":20,"save_ability":"dexterity","save_dc":15,"applies_condition":"restrained","release":[{"method":"let-go-of-your-end","rope_coils_in":"targets-space"},{"method":"bonus-action-to-repeat-command","rope_coils_in":"your-hand"}],"escape":{"action":"action","choices":[{"ability":"strength","skill":"athletics"},{"ability":"dexterity","skill":"acrobatics"}],"dc":15,"success":"no-longer-restrained"},"recoil":{"trigger":"target-escapes-while-you-hold-rope","action":"reaction","otherwise":"rope-coils-in-targets-space"},"armor_class":20,"hit_points":20,"immunities":["poison","psychic"],"regains_hit_points":1,"regain_interval_minutes":5,"regain_needs_hit_points":1,"destroyed_at_hit_points":0,"source_pages":[242,243]}}`,
		`{"scarab":{"appearance":"beetle-shaped-medallion","condition":"while-on-your-person","benefit_count":3,"armor_class_bonus":1,"preservation":{"max_charges":12,"triggers":["fail-save-against-necromancy-spell","fail-save-against-harmful-effect-from-undead"],"action":"reaction","charge_cost":1,"result":"failed-save-becomes-success","after_last_charge":"crumbles-to-powder-and-is-destroyed"},"spell_resistance":{"save_mode":"advantage","against":"spells"}}}`,
		`{"speed":{"attack_bonus":2,"damage_bonus":2,"extra_attack_action":"bonus","extra_attacks":1,"frequency":"each-of-your-turns"}}`,
		`{"sending_stones":{"stones_per_pair":"pair","matching":"carved-to-match","condition":"while-touching-one-stone","spell":"sending","target":"bearer-of-other-stone","if_no_bearer":"you-learn-this-and-do-not-cast","shared_reuse_after":"next-dawn","other_stone_when_one_destroyed":"nonmagical"}}`,
		`{"sentinel":{"condition":"while-holding","advantage_on":["initiative-rolls","wisdom-perception-checks"],"emblem":"eye"}}`,
		`{"missile_attraction":{"condition":"while-holding","defense":"resistance","defense_against":"damage-from-ranged-weapon-attacks","curse":{"begins_on":"attunement","ends_on":"remove-curse-or-similar-magic","removing_shield_ends_curse":false,"trigger":"ranged-weapon-attack-targets-nearby-creature","radius_feet":10,"result":"you-become-the-target"}}}`,
		`{"cavalier":{"condition":"while-holding","armor_class_bonus":2,"stacks_with_shield_bonus":true,"forceful_bash":{"trigger":"attack-action","attacks_replaced":1,"reach_feet":5,"attack_roll_adds":["proficiency-bonus","strength-modifier"],"damage":{"dice":2,"sides":6,"type":"force"},"damage_bonus":2,"damage_adds_ability_modifier":"strength","push_feet":10,"prone_if":"creature-your-size-or-smaller","applies_condition":"prone"},"protective_field":{"action":"reaction","triggers":["you-or-nearby-ally-targeted-by-attack","you-or-nearby-ally-saves-against-area-effect"],"ally_range_feet":5,"emanation_feet":5,"immobile":true,"not_fully_contained":"pushed-to-nearest-unoccupied-space-outside","triggering_effect":"no-effect-inside","concentration_required":true,"max_duration_minutes":1,"passage":"nothing-in-or-out","isolation":["inside-cannot-be-damaged-from-outside","inside-cannot-damage-outside"],"reuse_after":"next-dawn"}}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wield:weapon-attack::item.attunement item.magic-weapon bonus.attack-and-damage-3 item.six-buttons item.weapon-forms item.climbing-pole item.battering-ram item.magnetic-north save.constitution-dc-17 damage.drain-life-4d6-necrotic condition.paralyzed save.wisdom-dc-17 condition.frightened item.once-until-dawn",
		"hold:cast-spell:expend-listed-charges:item.attunement item.hold item.charges spell.cast-heal spell.cast-resurrection item.destroyed-on-1-after-last-charge",
		"hold:magic-action::item.attunement action.magic save.wisdom-dc-15 condition.charmed duration.8-hours item.once-until-dawn",
		"hold:magic-action::item.hold action.magic movement.transport-to-demiplane healing.hit-point-die-per-hour effect.no-aging duration.200-days-divided-by-creatures item.10-day-cooldown",
		"hold:magic-action:holding-one-end:action.magic item.animate-rope item.fasten-knot-and-coil check.climb-advantage-when-knotted item.armor-class-20-hit-points-20 defense.immunity-poison-and-psychic healing.1-hit-point-per-5-minutes",
		"hold:magic-action:holding-one-end:action.magic save.dexterity-dc-15 condition.restrained check.escape-dc-15 action.bonus-release action.reaction-recoil item.armor-class-20-hit-points-20 defense.immunity-poison-and-psychic healing.1-hit-point-per-5-minutes",
		"carry:reaction:failed-save-against-necromancy-or-undead expend-one-charge:item.attunement item.carry defense.armor-class-bonus-1 item.charges action.reaction save.failure-becomes-success-against-necromancy-or-undead item.destroyed-at-last-charge save.advantage-against-spells",
		"wield:weapon-attack::item.attunement item.magic-weapon bonus.attack-and-damage-2 attack.bonus-action-attack",
		"touch:cast-spell:sending:item.paired spell.cast-sending spell.target-other-stone-bearer item.once-until-dawn-per-pair item.nonmagical-if-pair-destroyed",
		"hold:none::item.hold check.initiative-advantage check.perception-advantage",
		"hold:none::item.attunement item.hold defense.resistance-ranged-weapon-damage item.curse item.curse-ends-on-remove-curse item.curse-redirects-ranged-attacks-within-10-feet",
		"hold:reaction:protective-field-trigger:item.attunement item.hold defense.armor-class-bonus-2 attack.forceful-bash damage.force-2d6-plus-2 movement.push-10-feet condition.prone action.reaction-protective-field spell.concentration-up-to-1-minute item.once-until-dawn",
	}
	index := 0
	for _, f := range loadItems241(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails241
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

func TestItems241CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems241(t)
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
	if len(composed) != 204 || variants != 265 || composed[0].Name != "Potion of Healing" || composed[192].Name != "Rod of Lordly Might" || composed[203].Name != "Shield of the Cavalier" {
		t.Fatalf("expected Potion of Healing, batches #225 to #240 and the twelve families of #241: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_241") {
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
func TestItems241RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items241JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items241JSON, &fresh); err != nil {
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
		if _, err := DecodeItems241(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems241RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails241]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails241] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-16" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 240 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/slippers-of-spider-climbing"
			b.Families[12].Name = "Slippers of Spider Climbing"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"scimitar-of-speed widened base", func(b *batch) { b.Families[7].Base.Options = append(b.Families[7].Base.Options, "any-weapon") }},
		{"sentinel-shield widened base", func(b *batch) { b.Families[9].Base.Options = append(b.Families[9].Base.Options, "any-weapon") }},
		{"shield-of-missile-attraction widened base", func(b *batch) { b.Families[10].Base.Options = append(b.Families[10].Base.Options, "any-weapon") }},
		{"shield-of-the-cavalier widened base", func(b *batch) { b.Families[11].Base.Options = append(b.Families[11].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{242, 242} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{241, 242} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"rod-of-lordly-might rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"rod-of-resurrection rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"rod-of-rulership rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"rod-of-security rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"rope-of-climbing rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"rope-of-entanglement rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"scarab-of-protection rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"scimitar-of-speed rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"sending-stones rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"sentinel-shield rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"shield-of-missile-attraction rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"shield-of-the-cavalier rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"rod-of-lordly-might attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"rod-of-resurrection attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"rod-of-rulership attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "none" }},
		{"rod-of-security attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "required" }},
		{"rope-of-climbing attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "required" }},
		{"rope-of-entanglement attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "required" }},
		{"scarab-of-protection attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "none" }},
		{"scimitar-of-speed attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"sending-stones attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "required" }},
		{"sentinel-shield attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "required" }},
		{"shield-of-missile-attraction attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "none" }},
		{"shield-of-the-cavalier attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 10 }},
		{"rod-of-resurrection dropped charges", func(b *batch) { *v(b, 1, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"rod-of-resurrection charge maximum", func(b *batch) { v(b, 1, 0).Charges.Maximum++ }},
		{"rod-of-resurrection charge recharge", func(b *batch) { v(b, 1, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"scarab-of-protection dropped charges", func(b *batch) { *v(b, 6, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"scarab-of-protection charge maximum", func(b *batch) { v(b, 6, 0).Charges.Maximum++ }},
		{"scarab-of-protection charge recharge", func(b *batch) { v(b, 6, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "reaction" }},
		{"dropped use restriction", func(b *batch) { v(b, 1, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails241{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Resurrection = v(b, 1, 0).Details.Resurrection }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.LordlyMight.FunctionsAs = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.LordlyMight.FunctionsAs = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.LordlyMight.AttackBonus = -1 }},
		{"lordly bonus", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.AttackBonus = 2 }},
		{"lordly buttons", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.ButtonCount = 5 }},
		{"lordly button order", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.Buttons.Spear.Button = 4 }},
		{"lordly drain dc", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.DrainLife.SaveDC = 15 }},
		{"lordly drain dice", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.DrainLife.ExtraDamage.Dice = 2 }},
		{"lordly terrify dc", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.Terrify.SaveDC = 15 }},
		{"lordly pages", func(b *batch) { p := v(b, 0, 0).Details.LordlyMight; p.SourcePages = []int{241} }},
		{"resurrection charges", func(b *batch) { p := v(b, 1, 0).Details.Resurrection; p.MaxCharges = 3 }},
		{"resurrection recharge", func(b *batch) { p := v(b, 1, 0).Details.Resurrection; p.RechargeCharges = 2 }},
		{"resurrection cost", func(b *batch) { p := v(b, 1, 0).Details.Resurrection; p.Spells[1].ChargeCost = 4 }},
		{"resurrection die", func(b *batch) { p := v(b, 1, 0).Details.Resurrection; p.LastCharge.DieSides = 6 }},
		{"rulership range", func(b *batch) { p := v(b, 2, 0).Details.Rulership; p.RangeFeet = 60 }},
		{"rulership dc", func(b *batch) { p := v(b, 2, 0).Details.Rulership; p.SaveDC = 17 }},
		{"rulership duration", func(b *batch) { p := v(b, 2, 0).Details.Rulership; p.DurationHours = 1 }},
		{"security travelers", func(b *batch) { p := v(b, 3, 0).Details.Security; p.Travelers.OtherWillingCreatures = 99 }},
		{"security days", func(b *batch) { p := v(b, 3, 0).Details.Security; p.Stay.BaseDays = 100 }},
		{"security cooldown", func(b *batch) { p := v(b, 3, 0).Details.Security; p.ReuseAfterDays = 1 }},
		{"climbing rope length", func(b *batch) { p := v(b, 4, 0).Details.ClimbingRope; p.LengthFeet = 50 }},
		{"climbing rope capacity", func(b *batch) { p := v(b, 4, 0).Details.ClimbingRope; p.CapacityPounds = 300 }},
		{"climbing rope ac", func(b *batch) { p := v(b, 4, 0).Details.ClimbingRope; p.ArmorClass = 15 }},
		{"climbing rope immunity", func(b *batch) { p := v(b, 4, 0).Details.ClimbingRope; p.Immunities[1] = "poison" }},
		{"entanglement range", func(b *batch) { p := v(b, 5, 0).Details.Entanglement; p.RangeFeet = 30 }},
		{"entanglement dc", func(b *batch) { p := v(b, 5, 0).Details.Entanglement; p.SaveDC = 13 }},
		{"entanglement escape", func(b *batch) { p := v(b, 5, 0).Details.Entanglement; p.Escape.DC = 20 }},
		{"entanglement release", func(b *batch) { p := v(b, 5, 0).Details.Entanglement; p.Release = p.Release[:1] }},
		{"entanglement pages", func(b *batch) { p := v(b, 5, 0).Details.Entanglement; p.SourcePages = []int{242} }},
		{"scarab armor", func(b *batch) { p := v(b, 6, 0).Details.Scarab; p.ArmorClassBonus = 2 }},
		{"scarab charges", func(b *batch) { p := v(b, 6, 0).Details.Scarab; p.Preservation.MaxCharges = 10 }},
		{"scarab duplicate trigger", func(b *batch) {
			p := v(b, 6, 0).Details.Scarab
			p.Preservation.Triggers[1] = "fail-save-against-necromancy-spell"
		}},
		{"scimitar bonus", func(b *batch) { p := v(b, 7, 0).Details.Speed; p.AttackBonus = 3 }},
		{"scimitar attacks", func(b *batch) { p := v(b, 7, 0).Details.Speed; p.ExtraAttacks = 2 }},
		{"stones spell", func(b *batch) { p := v(b, 8, 0).Details.SendingStones; p.Spell = "message" }},
		{"stones target", func(b *batch) { p := v(b, 8, 0).Details.SendingStones; p.Target = "any-creature" }},
		{"sentinel dropped advantage", func(b *batch) { p := v(b, 9, 0).Details.Sentinel; p.AdvantageOn = p.AdvantageOn[:1] }},
		{"attraction radius", func(b *batch) { p := v(b, 10, 0).Details.MissileAttraction; p.Curse.RadiusFeet = 5 }},
		{"attraction removal", func(b *batch) { p := v(b, 10, 0).Details.MissileAttraction; p.Curse.RemovingShieldEndsCurse = true }},
		{"cavalier armor", func(b *batch) { p := v(b, 11, 0).Details.Cavalier; p.ArmorClassBonus = 1 }},
		{"cavalier damage", func(b *batch) { p := v(b, 11, 0).Details.Cavalier; p.ForcefulBash.Damage.Sides = 4 }},
		{"cavalier push", func(b *batch) { p := v(b, 11, 0).Details.Cavalier; p.ForcefulBash.PushFeet = 5 }},
		{"cavalier field", func(b *batch) { p := v(b, 11, 0).Details.Cavalier; p.ProtectiveField.EmanationFeet = 10 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems241(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems241(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems241RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/rod-of-lordly-might/variant/standard"`)
	if !bytes.Contains(items241JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items241JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems241(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"lordly_might": {`, `"lordly_might_plus": {`},
		"unknown detail":     {`"attack_bonus": 3`, `"attack_bonus_all": 3`},
		"string for number":  {`"attack_bonus": 3`, `"attack_bonus": "3"`},
		"fractional number":  {`"attack_bonus": 3`, `"attack_bonus": 3.5`},
		"number for boolean": {`"buttons_stay_in_reach": true`, `"buttons_stay_in_reach": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items241JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items241JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems241(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items241JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems241(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items241JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems240(items241JSON); err == nil {
		t.Fatal("#240 decoder accepted the #241 batch")
	}
	if _, err := DecodeItems241(items240JSON); err == nil {
		t.Fatal("#241 decoder accepted the #240 batch")
	}
}
