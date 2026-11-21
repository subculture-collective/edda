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

func loadItems234(t *testing.T) ItemSourceBatch[ItemDetails234] {
	t.Helper()
	b, err := LoadItems234()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 228-229, in planning inventory order for issue #234. No family
// prints more than one variant: the Manual of Golems table rows are payload data.
func TestItems234SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		variants   []variant
	}{
		{"Iron Bands", "iron-bands", 228, "wondrous-item", nil, "none", []variant{{"standard", "Iron Bands", "rare"}}},
		{"Iron Flask", "iron-flask", 228, "wondrous-item", nil, "none", []variant{{"standard", "Iron Flask", "legendary"}}},
		{"Javelin of Lightning", "javelin-of-lightning", 228, "weapon", []string{"javelin"}, "none", []variant{{"standard", "Javelin of Lightning", "uncommon"}}},
		{"Lantern of Revealing", "lantern-of-revealing", 228, "wondrous-item", nil, "none", []variant{{"standard", "Lantern of Revealing", "uncommon"}}},
		{"Luck Blade", "luck-blade", 229, "weapon", []string{"glaive", "greatsword", "longsword", "rapier", "scimitar", "sickle", "shortsword"}, "required", []variant{{"standard", "Luck Blade", "legendary"}}},
		{"Mace of Disruption", "mace-of-disruption", 229, "weapon", []string{"mace"}, "required", []variant{{"standard", "Mace of Disruption", "rare"}}},
		{"Mace of Smiting", "mace-of-smiting", 229, "weapon", []string{"mace"}, "none", []variant{{"standard", "Mace of Smiting", "rare"}}},
		{"Mace of Terror", "mace-of-terror", 229, "weapon", []string{"mace"}, "required", []variant{{"standard", "Mace of Terror", "rare"}}},
		{"Mantle of Spell Resistance", "mantle-of-spell-resistance", 229, "wondrous-item", nil, "required", []variant{{"standard", "Mantle of Spell Resistance", "rare"}}},
		{"Manual of Bodily Health", "manual-of-bodily-health", 229, "wondrous-item", nil, "none", []variant{{"standard", "Manual of Bodily Health", "very-rare"}}},
		{"Manual of Gainful Exercise", "manual-of-gainful-exercise", 229, "wondrous-item", nil, "none", []variant{{"standard", "Manual of Gainful Exercise", "very-rare"}}},
		{"Manual of Golems", "manual-of-golems", 229, "wondrous-item", nil, "none", []variant{{"standard", "Manual of Golems", "very-rare"}}},
	}
	b := loadItems234(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-10" || b.CatalogueIssue != 234 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #234")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 228–229") {
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
				// Only the Luck Blade (1d3 charges, none regained) and the Mace of Terror (3, regains 1d3 at dawn) list charges.
				wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
				switch want.slug {
				case "luck-blade":
					wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "none"}
				case "mace-of-terror":
					wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}
				}
				if *v.Charges != wantCharges {
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
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 234 {
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
		t.Fatalf("inventory assigns %d families to issue #234", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "228, 229" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.IronBands.HitCondition = "changed"
	if fresh := loadItems234(t); fresh.Families[0].Variants[0].Details.IronBands.HitCondition != "restrained" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 228-229 and the Manual of Golems table at the top of p. 230. A printed
// "1 minute" is 60 seconds and "a century" is 100 years.
func TestItems234SourceValues(t *testing.T) {
	expected := []string{
		`{"iron_bands":{"appearance":"rusty-iron-sphere","diameter_inches":3,"weight_pounds":1,"action":"magic","throw_range_feet":60,"target_max_size":"huge","target_visibility":"can-see","opens_into":"tangle-of-metal-bands","attack_kind":"ranged","attack_ability":"dexterity","attack_adds_proficiency_bonus":true,"hit_condition":"restrained","release_action":"bonus","contracts_on":["release-command","missed-attack"],"break_action":"action","break_check":{"ability":"strength","skill":"athletics","dc":20},"break_by":"creature-that-can-touch-bands","break_success":["item-destroyed","restrained-creature-freed"],"break_failure":"further-attempts-by-that-creature-fail","break_failure_lockout_hours":24,"reuse":"next-dawn"}}`,
		`{"iron_flask":{"condition":"while-holding","action":"magic","range_feet":60,"target_visibility":"can-see","requires_empty_flask":true,"target_native_to":"other-plane-than-yours","save_ability":"wisdom","save_dc":17,"save_failure":"trapped-in-flask","previously_trapped_save":"advantage","trapped_until":"released","capacity_creatures":1,"trapped_effects":["does-not-age","does-not-need-to-breathe","does-not-need-to-eat","does-not-need-to-drink"],"release_action":"magic","release_by":"remove-stopper","obeys_for_hours":1,"understands_commands_in_any_language":true,"self_defense_when":["no-commands-issued","command-likely-to-cause-death-or-imprisonment"],"self_defense_behavior":"defends-itself-takes-no-other-actions","after_duration":"acts-per-normal-disposition-and-alignment","identify_reveals":"whether-flask-contains-creature","type_determined_by":"opening-flask","newly_found":"might-contain-creature-chosen-by-gm"}}`,
		`{"javelin_of_lightning":{"damage_swap_trigger":"attack-roll-hit","damage_swap_from":"piercing","damage_swap_to":"lightning","lightning_bolt":{"trigger":"throw-at-target-within-120-feet","range_feet":120,"forgo_attack_roll":true,"shape":"line","width_feet":5,"line_between":"you-and-target","affects":["target","each-other-creature-in-line-excluding-you"],"save_ability":"dexterity","save_dc":13,"damage":{"dice":4,"sides":6,"type":"lightning"},"on_success":"half-damage","returns":"weapon-reappears-in-hand-after-damage","reuse":"next-dawn"}}}`,
		`{"lantern_of_revealing":{"condition":"while-lit","burn_hours":6,"oil_pints":1,"bright_light_feet":30,"dim_light_additional_feet":30,"reveals":["invisible-creatures","invisible-objects"],"reveals_when":"in-bright-light","hood_action":"utilize","hooded_light":"dim-light","hooded_radius_feet":5}}`,
		`{"luck_blade":{"attack_bonus":1,"damage_bonus":1,"save_bonus":1,"save_bonus_when":"weapon-on-your-person","luck":{"condition":"weapon-on-your-person","action":"none","effect":"reroll-one-failed-d20-test-must-use-second-roll","forbidden_while":"incapacitated","reroll_use":"must-use-second-roll","reuse":"next-dawn"},"wish":{"initial_charges":{"dice":1,"sides":3,"type":"charges"},"condition":"while-holding","charge_cost":1,"spell":"wish","reuse":"next-dawn","property_lost_when":"no-charges"}}}`,
		`{"mace_of_disruption":{"extra_damage":{"dice":2,"sides":6,"type":"radiant"},"extra_damage_targets":["fiend","undead"],"destroy_at_hit_points_or_fewer":25,"save_ability":"wisdom","save_dc":15,"save_failure":"destroyed","save_success_condition":"frightened","success_condition_ends":"end-of-your-next-turn","light_condition":"while-holding","bright_light_feet":20,"dim_light_additional_feet":20}}`,
		`{"mace_of_smiting":{"attack_bonus":1,"damage_bonus":1,"construct_bonus":3,"critical_trigger":"natural-20-on-attack-roll","critical_extra_damage":7,"critical_damage_type":"bludgeoning","construct_critical_extra_damage":14,"construct_destroyed_at_hit_points_or_fewer":25}}`,
		`{"mace_of_terror":{"max_charges":3,"regain_expended_charges":{"dice":1,"sides":3,"type":"charges"},"regain_when":"daily-at-dawn","condition":"while-holding","action":"magic","charge_cost":1,"range_feet":30,"targets":"each-creature-of-your-choice","save_ability":"wisdom","save_dc":15,"fail_condition":"frightened","duration_seconds":60,"frightened_behavior":["must-move-as-far-from-you-as-possible","cannot-make-opportunity-attacks","action-only-dash-or-escape-movement-preventing-effect"],"no_movement_action":"dodge","repeat_save":"end-of-each-of-its-turns","success_ends":"effect-on-itself"}}`,
		`{"mantle_of_spell_resistance":{"condition":"while-wearing","roll":"saving-throw","mode":"advantage","against":"spells"}}`,
		`{"manual_of_bodily_health":{"study_hours":48,"within_days":6,"ability":"constitution","increase":2,"maximum":30,"then_loses_magic":true,"regains_magic_after_years":100}}`,
		`{"manual_of_gainful_exercise":{"study_hours":48,"within_days":6,"ability":"strength","increase":2,"maximum":30,"then_loses_magic":true,"regains_magic_after_years":100}}`,
		`{"manual_of_golems":{"source_pages":[229,230],"type_chosen_by":"gm-or-1d20-table","table_die_sides":20,"golems":[{"golem":"clay-golem","table_min":1,"table_max":5,"time_days":30,"cost_gp":65000},{"golem":"flesh-golem","table_min":6,"table_max":17,"time_days":60,"cost_gp":50000},{"golem":"iron-golem","table_min":18,"table_max":18,"time_days":120,"cost_gp":100000},{"golem":"stone-golem","table_min":19,"table_max":20,"time_days":90,"cost_gp":80000}],"decipher_requires_spellcaster":true,"decipher_spell_slots":2,"decipher_slot_level":5,"unable_user_damage":{"dice":6,"sides":6,"type":"psychic"},"work_uninterrupted":true,"max_rest_hours_per_day":8,"supplies_cost":"pay-table-cost","on_completion":"book-consumed-in-eldritch-flames","animation":"ashes-sprinkled-on-golem","stat_block_in":"monsters","golem_control":"under-your-control-understands-and-obeys"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:magic-action:throw-at-huge-or-smaller-within-60-feet once-until-dawn:action.magic attack.ranged-dexterity-plus-proficiency condition.restrained action.bonus-release check.strength-athletics-dc-20 effect.destroyed-on-successful-break item.once-until-dawn",
		"hold:magic-action:empty-flask-target-from-another-plane:action.magic save.wisdom-dc-17 effect.trapped-in-flask effect.obeys-commands-1-hour item.holds-one-creature",
		"wield:weapon-attack:lightning-bolt-once-until-dawn:damage.lightning-instead-of-piercing property.lightning-bolt area.line-5-foot save.dexterity-dc-13 damage.lightning-4d6 item.once-until-dawn",
		"none:utilize-action:lower-hood:light.bright-30-dim-30 effect.reveals-invisible-in-bright-light item.burns-6-hours-per-pint action.utilize light.hooded-dim-5",
		"wield:weapon-attack:luck-once-until-dawn wish-once-until-dawn:item.attunement bonus.attack-and-damage-1 bonus.saving-throw-1 property.luck action.none property.wish item.charges-1d3 item.once-until-dawn",
		"wield:weapon-attack::item.attunement damage.radiant-2d6 target.fiend-or-undead save.wisdom-dc-15 effect.destroyed-at-25-hit-points condition.frightened light.bright-20-dim-20",
		"wield:weapon-attack::bonus.attack-and-damage-1 bonus.construct-3 damage.bludgeoning-7-on-natural-20 damage.bludgeoning-14-on-construct-natural-20 effect.construct-destroyed-at-25-hit-points",
		"wield:magic-action:expend-1-charge:item.attunement item.charges action.magic area.range-30-foot save.wisdom-dc-15 condition.frightened duration.1-minute effect.frightened-must-flee",
		"wear:none::item.attunement item.wear save.advantage-against-spells",
		"none:study:study-48-hours-within-6-days:item.study-48-hours-in-6-days ability.constitution-plus-2-maximum-30 item.loses-magic-regains-in-century",
		"none:study:study-48-hours-within-6-days:item.study-48-hours-in-6-days ability.strength-plus-2-maximum-30 item.loses-magic-regains-in-century",
		"none:study:spellcaster-with-two-level-5-slots:item.golem-type-by-gm-or-1d20-table creation.time-and-cost-by-golem requirement.two-level-5-spell-slots damage.psychic-6d6-unable-reader item.consumed-in-eldritch-flames creature.golem-obeys-creator",
	}
	index := 0
	for _, f := range loadItems234(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails234
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

func TestItems234CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems234(t)
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
	if len(composed) != 97 || variants != 135 || composed[0].Name != "Potion of Healing" || composed[85].Name != "Iron Bands" || composed[96].Name != "Manual of Golems" {
		t.Fatalf("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #234, got %d families and %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_234") {
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
func TestItems234RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items234JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items234JSON, &fresh); err != nil {
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
		if _, err := DecodeItems234(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems234RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails234]
	// v addresses the single variant of a family by index; see TestItems234SourceIdentities.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails234] { return &b.Families[family].Variants[0] }
	const (
		bands, flask, javelin, lantern, luck, disruption, smiting, terror, mantle, health, exercise, golems = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-09" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 233 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/manual-of-quickness-of-action"
			b.Families[12].Name = "Manual of Quickness of Action"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/iron-band" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Iron Band" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) { b.Families[javelin].Base.Options = []string{"javelin", "spear"} }},
		{"dropped sickle", func(b *batch) { b.Families[luck].Base.Options = b.Families[luck].Base.Options[:5] }},
		{"reordered luck base", func(b *batch) {
			o := b.Families[luck].Base.Options
			o[0], o[1] = o[1], o[0]
		}},
		{"mace base", func(b *batch) { b.Families[terror].Base.Options = []string{"flail"} }},
		{"invented base", func(b *batch) { b.Families[bands].Base.Options = []string{"any-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{229, 229} }},
		{"page range", func(b *batch) { b.Families[golems].Presentations[0].Pages = []int{229, 230} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"golem rows as variants", func(b *batch) {
			f := &b.Families[golems]
			f.Variants = append(f.Variants, f.Variants[0], f.Variants[0], f.Variants[0])
		}},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Bands" }},
		{"rarity", func(b *batch) { v(b, 0).Rarity = "uncommon" }},
		{"flask rarity", func(b *batch) { v(b, flask).Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"invented attunement", func(b *batch) { v(b, smiting).Attunement.Mode = "required" }},
		{"dropped attunement", func(b *batch) { v(b, disruption).Attunement.Mode = "none" }},
		{"mantle attunement", func(b *batch) { v(b, mantle).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, luck).Attunement.Restrictions = []string{"spellcaster"} }},
		{"golem attunement", func(b *batch) { v(b, golems).Attunement.Mode = "required" }},
		{"missing charges", func(b *batch) { v(b, terror).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, terror).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"terror maximum", func(b *batch) { v(b, terror).Charges.Maximum = 4 }},
		{"terror recharge", func(b *batch) { v(b, terror).Charges.Recharge = "regain-all-at-dawn" }},
		{"luck charges dropped", func(b *batch) { *v(b, luck).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"luck recharges", func(b *batch) { v(b, luck).Charges.Recharge = "regain-1d3-at-dawn" }},
		{"invented charges", func(b *batch) { v(b, disruption).Charges.Maximum = 3 }},
		{"flask charges", func(b *batch) { v(b, flask).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, flask).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, luck).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, lantern).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, bands).Use.Restrictions = []string{} }},
		{"study restriction", func(b *batch) { v(b, health).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"golem supplies are not a price", func(b *batch) {
			*v(b, golems).Cost = ItemCost{Kind: "listed", Amount: 65000, Currency: "gp", SourcePages: []int{229, 229}}
		}},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails234{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"swapped manuals", func(b *batch) { v(b, health).Details = v(b, exercise).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.Mantle = v(b, mantle).Details.Mantle }},
		{"blank key", func(b *batch) { v(b, 0).Details.IronBands.HitCondition = " " }},
		{"prose value", func(b *batch) { v(b, 0).Details.IronBands.HitCondition = "Restrained" }},
		{"negative value", func(b *batch) { v(b, 0).Details.IronBands.ThrowRangeFeet = -1 }},
		{"bands break dc", func(b *batch) { v(b, bands).Details.IronBands.BreakCheck.DC = 15 }},
		{"bands lockout", func(b *batch) { v(b, bands).Details.IronBands.BreakLockoutHours = 0 }},
		{"bands size", func(b *batch) { v(b, bands).Details.IronBands.DiameterInches = 0 }},
		{"bands duplicate end", func(b *batch) { v(b, bands).Details.IronBands.ContractsOn[1] = "release-command" }},
		{"bands no success", func(b *batch) { v(b, bands).Details.IronBands.BreakSuccess = []string{} }},
		{"flask dc", func(b *batch) { v(b, flask).Details.IronFlask.SaveDC = 13 }},
		{"flask capacity", func(b *batch) { v(b, flask).Details.IronFlask.Capacity = 0 }},
		{"flask duplicate effect", func(b *batch) { v(b, flask).Details.IronFlask.TrappedEffects[1] = "does-not-age" }},
		{"flask no defence", func(b *batch) { v(b, flask).Details.IronFlask.SelfDefenseOn = []string{} }},
		{"flask obey hours", func(b *batch) { v(b, flask).Details.IronFlask.ObeysHours = 0 }},
		{"javelin dc", func(b *batch) { v(b, javelin).Details.Javelin.Bolt.SaveDC = 15 }},
		{"javelin dice", func(b *batch) { v(b, javelin).Details.Javelin.Bolt.Damage.Dice = 0 }},
		{"javelin range", func(b *batch) { v(b, javelin).Details.Javelin.Bolt.RangeFeet = 0 }},
		{"javelin duplicate target", func(b *batch) { v(b, javelin).Details.Javelin.Bolt.Affects[1] = "target" }},
		{"lantern light", func(b *batch) { v(b, lantern).Details.Lantern.BrightLightFeet = 0 }},
		{"lantern hood", func(b *batch) { v(b, lantern).Details.Lantern.HoodedRadiusFeet = 0 }},
		{"lantern duplicate reveal", func(b *batch) { v(b, lantern).Details.Lantern.Reveals[1] = "invisible-creatures" }},
		{"luck bonus", func(b *batch) { v(b, luck).Details.LuckBlade.SaveBonus = 0 }},
		{"luck charge dice", func(b *batch) { v(b, luck).Details.LuckBlade.Wish.InitialCharges.Sides = 0 }},
		{"luck charge cost", func(b *batch) { v(b, luck).Details.LuckBlade.Wish.ChargeCost = 0 }},
		{"disruption dice", func(b *batch) { v(b, disruption).Details.Disruption.ExtraDamage.Sides = 0 }},
		{"disruption dc", func(b *batch) { v(b, disruption).Details.Disruption.SaveDC = 13 }},
		{"disruption threshold", func(b *batch) { v(b, disruption).Details.Disruption.DestroyAtHitPoints = 0 }},
		{"disruption targets", func(b *batch) { v(b, disruption).Details.Disruption.ExtraDamageTargets[1] = "fiend" }},
		{"disruption light", func(b *batch) { v(b, disruption).Details.Disruption.DimLightAdditionalFeet = 0 }},
		{"smiting bonus", func(b *batch) { v(b, smiting).Details.Smiting.AttackBonus = 0 }},
		{"smiting construct bonus", func(b *batch) { v(b, smiting).Details.Smiting.ConstructBonus = 1 }},
		{"smiting construct crit", func(b *batch) { v(b, smiting).Details.Smiting.ConstructCritDamage = 7 }},
		{"smiting threshold", func(b *batch) { v(b, smiting).Details.Smiting.DestroyAtHitPoints = 0 }},
		{"terror dc", func(b *batch) { v(b, terror).Details.Terror.SaveDC = 13 }},
		{"terror charges", func(b *batch) { v(b, terror).Details.Terror.MaxCharges = 4 }},
		{"terror duration", func(b *batch) { v(b, terror).Details.Terror.DurationSeconds = 0 }},
		{"terror regain", func(b *batch) { v(b, terror).Details.Terror.Regain.Dice = 0 }},
		{"terror behaviour", func(b *batch) { v(b, terror).Details.Terror.Behavior = []string{} }},
		{"mantle mode", func(b *batch) { v(b, mantle).Details.Mantle.Mode = "disadvantage" }},
		{"health ability", func(b *batch) { v(b, health).Details.BodilyHealth.Ability = "strength" }},
		{"health hours", func(b *batch) { v(b, health).Details.BodilyHealth.StudyHours = 24 }},
		{"health maximum", func(b *batch) { v(b, health).Details.BodilyHealth.Maximum = 20 }},
		{"exercise ability", func(b *batch) { v(b, exercise).Details.GainfulExerc.Ability = "constitution" }},
		{"exercise increase", func(b *batch) { v(b, exercise).Details.GainfulExerc.Increase = 1 }},
		{"exercise century", func(b *batch) { v(b, exercise).Details.GainfulExerc.RegainsAfterYears = 10 }},
		{"golem pages", func(b *batch) { v(b, golems).Details.ManualGolems.SourcePages = []int{229} }},
		{"golem die", func(b *batch) { v(b, golems).Details.ManualGolems.TableDieSides = 100 }},
		{"golem missing row", func(b *batch) { g := v(b, golems).Details.ManualGolems; g.Golems = g.Golems[:3] }},
		{"golem row gap", func(b *batch) { v(b, golems).Details.ManualGolems.Golems[1].TableMax = 16 }},
		{"golem row overlap", func(b *batch) { v(b, golems).Details.ManualGolems.Golems[1].TableMin = 5 }},
		{"golem duplicate row", func(b *batch) { v(b, golems).Details.ManualGolems.Golems[3].Golem = "iron-golem" }},
		{"golem time", func(b *batch) { v(b, golems).Details.ManualGolems.Golems[0].Days = 0 }},
		{"golem cost", func(b *batch) { v(b, golems).Details.ManualGolems.Golems[0].CostGP = 0 }},
		{"golem slots", func(b *batch) { v(b, golems).Details.ManualGolems.DecipherSlots = 0 }},
		{"golem damage", func(b *batch) { v(b, golems).Details.ManualGolems.UnableDamage.Dice = 0 }},
		{"golem rest", func(b *batch) { v(b, golems).Details.ManualGolems.MaxRestHoursPerDay = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems234(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems234(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems234RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/iron-bands/variant/standard"`)
	if !bytes.Contains(items234JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items234JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems234(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"iron_bands": {`, `"iron_bands_plus": {`},
		"unknown detail":     {`"hit_condition"`, `"hit_condition_all"`},
		"string for number":  {`"save_dc": 17`, `"save_dc": "17"`},
		"fractional number":  {`"save_dc": 17`, `"save_dc": 17.5`},
		"number for boolean": {`"requires_empty_flask": true`, `"requires_empty_flask": 1`},
		"null golem name":    {`"golem": "clay-golem"`, `"golem": null`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items234JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items234JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems234(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items234JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems234(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items234JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items234JSON); err == nil {
		t.Fatal("#231 decoder accepted the #234 batch")
	}
	if _, err := DecodeItems234(items231JSON); err == nil {
		t.Fatal("#234 decoder accepted the #231 batch")
	}
}
