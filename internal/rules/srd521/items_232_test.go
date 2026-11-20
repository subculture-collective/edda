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

func loadItems232(t *testing.T) ItemSourceBatch[ItemDetails232] {
	t.Helper()
	b, err := LoadItems232()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 223-225, in planning inventory order for issue #232. No family
// in this batch prints several kinds, so each has one "standard" variant.
func TestItems232SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		rarity     string
	}{
		{"Gem of Seeing", "gem-of-seeing", 223, "wondrous-item", nil, "required", "rare"},
		{"Giant Slayer", "giant-slayer", 223, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, "none", "rare"},
		{"Glamoured Studded Leather", "glamoured-studded-leather", 223, "armor", []string{"studded-leather-armor"}, "none", "rare"},
		{"Gloves of Missile Snaring", "gloves-of-missile-snaring", 224, "wondrous-item", nil, "required", "uncommon"},
		{"Gloves of Swimming and Climbing", "gloves-of-swimming-and-climbing", 224, "wondrous-item", nil, "required", "uncommon"},
		{"Gloves of Thievery", "gloves-of-thievery", 224, "wondrous-item", nil, "none", "uncommon"},
		{"Goggles of Night", "goggles-of-night", 224, "wondrous-item", nil, "none", "uncommon"},
		{"Hammer of Thunderbolts", "hammer-of-thunderbolts", 224, "weapon", []string{"maul", "warhammer"}, "required", "legendary"},
		{"Handy Haversack", "handy-haversack", 224, "wondrous-item", nil, "none", "rare"},
		{"Hat of Disguise", "hat-of-disguise", 224, "wondrous-item", nil, "required", "uncommon"},
		{"Hat of Many Spells", "hat-of-many-spells", 224, "wondrous-item", nil, "required", "very-rare"},
		{"Headband of Intellect", "headband-of-intellect", 225, "wondrous-item", nil, "required", "uncommon"},
	}
	b := loadItems232(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-08" || b.CatalogueIssue != 232 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #232")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 223–225") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
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
			// Only the Gem of Seeing (3) and the Hammer of Thunderbolts (5) list charges.
			wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
			switch want.slug {
			case "gem-of-seeing":
				wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}
			case "hammer-of-thunderbolts":
				wantCharges = ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-plus-1-at-dawn"}
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
		if r.Issue != 232 {
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
		t.Fatalf("inventory assigns %d families to issue #232", index)
	}
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "223, 224, 225" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.GemSeeing.Sense = "changed"
	if fresh := loadItems232(t); fresh.Families[0].Variants[0].Details.GemSeeing.Sense != "truesight" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 223-225. The Hat of Many Spells starts on p. 224 and its table continues on
// p. 225; a printed 00 is recorded as 100.
func TestItems232SourceValues(t *testing.T) {
	expected := []string{
		`{"gem_of_seeing":{"max_charges":3,"action":"magic","charge_cost":1,"sense":"truesight","sense_range_feet":120,"condition":"peering-through-gem","duration_seconds":600,"recharge":{"dice":1,"sides":3,"type":"charges"},"recharge_when":"daily-at-dawn"}}`,
		`{"giant_slayer":{"attack_bonus":1,"damage_bonus":1,"trigger":"hit-giant-with-weapon","target_creature_type":"giant","extra_damage":{"dice":2,"sides":6,"type":"weapon-type"},"save_ability":"strength","save_dc":15,"condition":"prone"}}`,
		`{"glamoured_studded_leather":{"armor_class_bonus":1,"action":"bonus","appearance":"normal-clothing-or-other-armor","chosen_details":["color","style","accessories"],"retains_bulk_and_weight":true,"ends_on":["use-property-again","doff-armor"]}}`,
		`{"gloves_of_missile_snaring":{"trigger":"hit-by-ranged-or-thrown-weapon-attack-roll","action":"reaction","reduction":{"dice":1,"sides":10,"type":"damage-reduction"},"reduction_bonus_ability":"dexterity","requires_free_hand":true,"catch_condition":"damage-reduced-to-0","catchable":"ammunition-or-weapon-small-enough-to-hold-in-that-hand"}}`,
		`{"gloves_of_swimming_and_climbing":{"condition":"while-wearing","climb_speed":"equal-to-speed","swim_speed":"equal-to-speed","check_bonus":5,"check_ability":"strength","check_skill":"athletics","check_subjects":["climb","swim"]}}`,
		`{"gloves_of_thievery":{"appearance":"imperceptible-while-worn","condition":"while-wearing","check_bonus":5,"check_ability":"dexterity","check_skill":"sleight-of-hand"}}`,
		`{"goggles_of_night":{"appearance":"dark-lenses","condition":"while-wearing","darkvision_feet":60,"existing_darkvision_increase_feet":60}}`,
		`{"hammer_of_thunderbolts":{"attack_bonus":1,"damage_bonus":1,"max_charges":5,"hurl":{"charge_cost":1,"thrown_normal_range_feet":20,"thrown_long_range_feet":60,"thunderclap_audible_feet":300,"burst_radius_feet":30,"burst_excludes":"you","save_ability":"constitution","save_dc":17,"condition":"stunned","condition_ends":"end-of-your-next-turn","returns_to_hand":"immediately-after-hitting-or-missing","thunderclap_requires_hit":true},"recharge":{"dice":1,"sides":4,"type":"charges"},"recharge_bonus":1,"recharge_when":"daily-at-dawn","giants_bane":{"requires_attunement_to_weapon":true,"requires_wearing_either":["belt-of-giant-strength","gauntlets-of-ogre-power"],"requires_attunement_to_worn_item":true,"natural_roll":20,"target_creature_type":"giant","save_ability":"constitution","save_dc":17,"failure":"die"},"might_of_giants":{"affects":["belt-of-giant-strength","gauntlets-of-ogre-power"],"strength_increase":4,"strength_maximum":30}}}`,
		`{"handy_haversack":{"weight_pounds":5,"side_pouches":{"count":2,"capacity_pounds":200,"capacity_cubic_feet":25},"central_pouch":{"count":1,"capacity_pounds":500,"capacity_cubic_feet":64},"retrieval_actions":["utilize","bonus-action"],"retrieved_item_on_top":true,"rupture_triggers":["overloaded","pierced","torn"],"contents_lost_when_destroyed":true,"artifact_turns_up_again":true,"inside_out_result":"contents-spill-unharmed-and-haversack-must-be-put-right","air_minutes_per_pouch":10,"air_divided_by":"breathing-creatures-inside","nested_extradimensional_space":{"items":["bag-of-holding","portable-hole","similar-item"],"result":"both-items-destroyed","gate_to":"astral-plane","gate_radius_feet":10,"spared_by":"total-cover","destination":"random-location-on-astral-plane","gate_closes_after":true,"gate_one_way":true,"gate_reopens":false}}}`,
		`{"hat_of_disguise":{"condition":"while-wearing","spell":"disguise-self","ends_when":"hat-removed"}}`,
		`{"hat_of_many_spells":{"source_pages":[224,225],"attunement_class":"wizard","condition":"while-holding","focus_spell_list":"wizard","added_somatic_component":"reach-into-hat-and-pull-spell-out","unknown_spell":{"min_spell_level":1,"spell_list":"wizard","level_limit":"level-you-can-cast","material_cost_limit_gp":1000,"slot_expended":"spell-level","check_ability":"intelligence","check_skill":"arcana","check_dc_base":10,"check_dc_adds":"spell-level","success_casting_time":"normal-casting-time","success_lockout_until":"short-or-long-rest","failure_result":"fail-to-cast-and-random-effect","table_die_sides":100},"spell_stats":"your-spell-save-dc-and-spell-attack-bonus","failure_table":[{"min":1,"max":50,"effect":"cast-random-spell","sub_die_sides":10,"outcomes":[{"min":1,"max":1,"result":"enlarge-reduce-enlarge-effect"},{"min":2,"max":2,"result":"enlarge-reduce-reduce-effect"},{"min":3,"max":3,"result":"faerie-fire"},{"min":4,"max":4,"result":"fireball"},{"min":5,"max":5,"result":"gust-of-wind"},{"min":6,"max":6,"result":"invisibility-cast-on-yourself"},{"min":7,"max":7,"result":"lightning-bolt"},{"min":8,"max":8,"result":"phantasmal-force"},{"min":9,"max":9,"result":"polymorph"},{"min":10,"max":10,"result":"stinking-cloud"}],"facts":[],"traits":[]},{"min":51,"max":55,"effect":"stunned-condition","sub_die_sides":0,"outcomes":[],"facts":[],"traits":["until-end-of-your-next-turn","believing-something-awesome-happened"]},{"min":56,"max":60,"effect":"harmless-swarm-of-butterflies","sub_die_sides":0,"outcomes":[],"facts":[{"key":"cube-size","amount":10,"unit":"feet"},{"key":"range","amount":30,"unit":"feet"},{"key":"disperses-after","amount":1,"unit":"minutes"}],"traits":["swarm-fills-cube-within-range-of-yourself"]},{"min":61,"max":65,"effect":"pull-nonmagical-object","sub_die_sides":4,"outcomes":[{"min":1,"max":1,"result":"vial-of-acid"},{"min":2,"max":2,"result":"flask-of-alchemists-fire"},{"min":3,"max":3,"result":"crowbar"},{"min":4,"max":4,"result":"lit-torch"}],"facts":[],"traits":[]},{"min":66,"max":70,"effect":"magic-sickness-poisoned-condition","sub_die_sides":0,"outcomes":[],"facts":[{"key":"duration","amount":1,"unit":"hours"}],"traits":[]},{"min":71,"max":75,"effect":"petrified-condition","sub_die_sides":0,"outcomes":[],"facts":[],"traits":["until-end-of-your-next-turn"]},{"min":76,"max":80,"effect":"pull-nonmagical-object","sub_die_sides":4,"outcomes":[{"min":1,"max":1,"result":"dagger"},{"min":2,"max":2,"result":"rope-with-grappling-hook-tied-to-one-end"},{"min":3,"max":3,"result":"bag-of-caltrops"},{"min":4,"max":4,"result":"gem-worth-50-gp"}],"facts":[{"key":"gem-value","amount":50,"unit":"gp"}],"traits":[]},{"min":81,"max":85,"effect":"creature-appears","sub_die_sides":4,"outcomes":[{"min":1,"max":1,"result":"camel"},{"min":2,"max":2,"result":"constrictor-snake"},{"min":3,"max":3,"result":"elephant"},{"min":4,"max":4,"result":"mule"}],"facts":[{"key":"duration","amount":1,"unit":"hours"}],"traits":["unoccupied-space-as-close-to-you-as-possible","not-under-your-control","acts-as-it-normally-would","disappears-when-duration-elapses","disappears-at-0-hit-points"]},{"min":86,"max":90,"effect":"hostile-swarm-of-bats","sub_die_sides":0,"outcomes":[],"facts":[],"traits":["flies-out-of-hat","occupies-your-space","attacks-you"]},{"min":91,"max":95,"effect":"two-way-portal-to-another-plane","sub_die_sides":0,"outcomes":[],"facts":[{"key":"portal-diameter","amount":10,"unit":"feet"},{"key":"range","amount":30,"unit":"feet"}],"traits":["vertical","unoccupied-space-within-range","open-until-end-of-your-next-turn","gm-determines-destination"]},{"min":96,"max":100,"effect":"pull-magic-item","sub_die_sides":6,"outcomes":[{"min":1,"max":3,"result":"common"},{"min":4,"max":5,"result":"uncommon"},{"min":6,"max":6,"result":"rare"}],"facts":[{"key":"disappears-after","amount":1,"unit":"hours"}],"traits":["gm-chooses-item","disappears-if-not-consumed-or-destroyed"]}]}}`,
		`{"headband_of_intellect":{"condition":"while-wearing","ability":"intelligence","score":19,"no_effect_at_or_above_score":19}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:magic-action:expend-one-charge:item.attunement item.charges action.magic senses.truesight-120-feet duration.ten-minutes item.peer-through-gem",
		"wield:weapon-attack::item.magic-weapon bonus.attack-and-damage-1 damage.extra-2d6-against-giants save.strength-dc-15 condition.prone",
		"wear:bonus-action:assume-appearance-of-clothing-or-armor:item.wear defense.armor-class-bonus-1 action.bonus-action effect.illusory-appearance item.retains-bulk-and-weight effect.appearance-lasts-until-used-again-or-doffed",
		"wear:reaction:hit-by-ranged-or-thrown-attack free-hand:item.attunement item.wear action.reaction damage.reduce-1d10-plus-dexterity-modifier item.free-hand-required item.catch-ammunition-or-weapon",
		"wear:none::item.attunement item.wear movement.climb-speed-equals-speed movement.swim-speed-equals-speed check.athletics-bonus-5-climb-or-swim",
		"wear:none::item.wear item.imperceptible-while-worn check.sleight-of-hand-bonus-5",
		"wear:none::item.wear vision.darkvision-60-feet vision.darkvision-increase-60-feet",
		"wield:weapon-attack:expend-one-charge thrown-ranged-attack:item.attunement item.magic-weapon bonus.attack-and-damage-1 item.charges property.thrown-20-60 effect.thunderclap-audible-300-feet save.constitution-dc-17 condition.stunned item.returns-to-hand property.giants-bane effect.die-on-failed-save-against-giant ability.strength-score-bonus-4-max-30",
		"none:container:retrieve-with-utilize-or-bonus-action:item.extradimensional-storage item.side-pouches-200-pounds-25-cubic-feet item.central-pouch-500-pounds-64-cubic-feet item.weight-5-pounds item.retrieved-item-on-top item.destroyed-if-overloaded-pierced-or-torn item.contents-lost-if-destroyed item.inside-out-spills-contents item.air-10-minutes-per-pouch item.nested-extradimensional-space-opens-astral-gate",
		"wear:cast-spell:disguise-self:item.attunement item.wear spell.cast-disguise-self item.spell-ends-when-removed",
		"hold:cast-spell:unknown-wizard-spell until-short-or-long-rest-after-success:item.attunement attunement.wizard item.hold spell.focus-wizard-spells spell.added-somatic-component spell.cast-unknown-wizard-spell check.intelligence-arcana-dc-10-plus-spell-level item.spell-slot-expended table.failure-1d100 spell.uses-your-save-dc-and-attack-bonus",
		"wear:none::item.attunement item.wear ability.intelligence-score-19",
	}
	index := 0
	for _, f := range loadItems232(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails232
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

func TestItems232CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems232(t)
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
	sixth, err := LoadItems230()
	if err != nil {
		t.Fatal(err)
	}
	seventh, err := LoadItems231()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), first.SourceIdentities(), second.SourceIdentities(), third.SourceIdentities(), fourth.SourceIdentities(), fifth.SourceIdentities(), sixth.SourceIdentities(), seventh.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 97 || variants != 135 || composed[0].Name != "Potion of Healing" || composed[85].Name != "Gem of Seeing" || composed[96].Name != "Headband of Intellect" {
		t.Fatal("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #232")
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
			if strings.Contains(fixture, "items_232") {
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
func TestItems232RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items232JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 700 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items232JSON, &fresh); err != nil {
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
		if _, err := DecodeItems232(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems232RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails232]
	// v addresses the single variant of a family by index; see TestItems232SourceIdentities.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails232] { return &b.Families[family].Variants[0] }
	const (
		seeing, slayer, glamoured, snaring, swimclimb, thievery, goggles, hammer, haversack, disguise, spells, headband = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-07" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 231 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/helm-of-brilliance"
			b.Families[12].Name = "Helm of Brilliance"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/gem-of-brightness" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Gem of True Seeing" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) {
			b.Families[glamoured].Base.Options = append(b.Families[glamoured].Base.Options, "leather-armor")
		}},
		{"dropped base", func(b *batch) { b.Families[hammer].Base.Options = []string{"maul"} }},
		{"giant slayer base", func(b *batch) { b.Families[slayer].Base.Options = []string{"any-weapon"} }},
		{"invented base", func(b *batch) { b.Families[goggles].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{224, 224} }},
		{"page range", func(b *batch) { b.Families[spells].Presentations[0].Pages = []int{224, 225} }},
		{"headband page", func(b *batch) { b.Families[headband].Presentations[0].Pages = []int{224, 224} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Gem" }},
		{"rarity", func(b *batch) { v(b, 0).Rarity = "uncommon" }},
		{"hammer rarity", func(b *batch) { v(b, hammer).Rarity = "very-rare" }},
		{"hat rarity", func(b *batch) { v(b, spells).Rarity = "rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, seeing).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, haversack).Attunement.Mode = "required" }},
		{"slayer attunement", func(b *batch) { v(b, slayer).Attunement.Mode = "required" }},
		{"wizard is payload data, not a restriction", func(b *batch) { v(b, spells).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, seeing).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, seeing).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"charge maximum", func(b *batch) { v(b, seeing).Charges.Maximum = 4 }},
		{"charge recharge", func(b *batch) { v(b, seeing).Charges.Recharge = "regain-all-at-dawn" }},
		{"hammer charge maximum", func(b *batch) { v(b, hammer).Charges.Maximum = 3 }},
		{"hammer recharge", func(b *batch) { v(b, hammer).Charges.Recharge = "regain-1d4-at-dawn" }},
		{"invented charges", func(b *batch) { v(b, haversack).Charges.Maximum = 10 }},
		{"equip mode", func(b *batch) { v(b, spells).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, snaring).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, snaring).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, snaring).Use.Restrictions = []string{"free-hand"} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails232{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.GiantSlayer = v(b, 1).Details.GiantSlayer }},
		{"blank key", func(b *batch) { v(b, 0).Details.GemSeeing.Sense = " " }},
		{"prose value", func(b *batch) { v(b, 0).Details.GemSeeing.Sense = "True sight" }},
		{"negative value", func(b *batch) { v(b, 0).Details.GemSeeing.ChargeCost = -1 }},
		{"seeing charges", func(b *batch) { v(b, seeing).Details.GemSeeing.MaxCharges = 4 }},
		{"seeing range", func(b *batch) { v(b, seeing).Details.GemSeeing.SenseRangeFeet = 60 }},
		{"seeing duration", func(b *batch) { v(b, seeing).Details.GemSeeing.DurationSeconds = 60 }},
		{"seeing recharge dice", func(b *batch) { v(b, seeing).Details.GemSeeing.Recharge.Sides = 4 }},
		{"slayer bonus", func(b *batch) { v(b, slayer).Details.GiantSlayer.AttackBonus = 0 }},
		{"slayer dc", func(b *batch) { v(b, slayer).Details.GiantSlayer.SaveDC = 17 }},
		{"slayer dice", func(b *batch) { v(b, slayer).Details.GiantSlayer.ExtraDamage.Dice = 0 }},
		{"glamoured bonus", func(b *batch) { v(b, glamoured).Details.Glamoured.ArmorClassBonus = 0 }},
		{"glamoured details", func(b *batch) { v(b, glamoured).Details.Glamoured.ChosenDetails = []string{} }},
		{"glamoured duplicate end", func(b *batch) { v(b, glamoured).Details.Glamoured.EndsOn[1] = "use-property-again" }},
		{"snaring dice", func(b *batch) { v(b, snaring).Details.MissileSnaring.Reduction.Sides = 0 }},
		{"swim climb bonus", func(b *batch) { v(b, swimclimb).Details.SwimClimb.CheckBonus = 3 }},
		{"swim climb subjects", func(b *batch) { v(b, swimclimb).Details.SwimClimb.CheckSubjects = []string{"climb", "climb"} }},
		{"thievery bonus", func(b *batch) { v(b, thievery).Details.Thievery.CheckBonus = 0 }},
		{"goggles range", func(b *batch) { v(b, goggles).Details.Goggles.DarkvisionFeet = 0 }},
		{"goggles increase", func(b *batch) { v(b, goggles).Details.Goggles.ExistingIncreaseFeet = 0 }},
		{"hammer dc", func(b *batch) { v(b, hammer).Details.Hammer.Hurl.SaveDC = 15 }},
		{"hammer bane dc", func(b *batch) { v(b, hammer).Details.Hammer.Bane.SaveDC = 15 }},
		{"hammer long range", func(b *batch) { v(b, hammer).Details.Hammer.Hurl.ThrownLongFeet = 20 }},
		{"hammer charges", func(b *batch) { v(b, hammer).Details.Hammer.MaxCharges = 6 }},
		{"hammer recharge bonus", func(b *batch) { v(b, hammer).Details.Hammer.RechargeBonus = 0 }},
		{"hammer duplicate worn item", func(b *batch) { v(b, hammer).Details.Hammer.Bane.WornWithEither[1] = "belt-of-giant-strength" }},
		{"hammer strength maximum", func(b *batch) { v(b, hammer).Details.Hammer.Might.StrengthMaximum = 0 }},
		{"haversack weight", func(b *batch) { v(b, haversack).Details.Haversack.WeightPounds = 0 }},
		{"haversack side pouches", func(b *batch) { v(b, haversack).Details.Haversack.SidePouches.Count = 1 }},
		{"haversack central pouch", func(b *batch) { v(b, haversack).Details.Haversack.CentralPouch.Count = 2 }},
		{"haversack capacity", func(b *batch) { v(b, haversack).Details.Haversack.CentralPouch.CapacityCubic = 0 }},
		{"haversack duplicate rupture", func(b *batch) { v(b, haversack).Details.Haversack.RuptureTriggers[1] = "overloaded" }},
		{"haversack no retrieval", func(b *batch) { v(b, haversack).Details.Haversack.RetrievalActions = []string{} }},
		{"haversack gate radius", func(b *batch) { v(b, haversack).Details.Haversack.Nested.GateRadiusFt = 0 }},
		{"haversack no nested items", func(b *batch) { v(b, haversack).Details.Haversack.Nested.Items = []string{} }},
		{"hat pages", func(b *batch) { v(b, spells).Details.HatManySpells.SourcePages = []int{224} }},
		{"hat dc", func(b *batch) { v(b, spells).Details.HatManySpells.Unknown.CheckDCBase = 11 }},
		{"hat die", func(b *batch) { v(b, spells).Details.HatManySpells.Unknown.TableDieSides = 20 }},
		{"hat cost limit", func(b *batch) { v(b, spells).Details.HatManySpells.Unknown.MaterialCostLimitGP = 0 }},
		{"hat missing row", func(b *batch) {
			v(b, spells).Details.HatManySpells.FailureTable = v(b, spells).Details.HatManySpells.FailureTable[:10]
		}},
		{"hat row gap", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[1].Min = 52 }},
		{"hat row overlap", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[0].Max = 51 }},
		{"hat last row is 00", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[10].Max = 99 }},
		{"hat sub die without outcomes", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[0].Outcomes = []ItemHatOutcome{} }},
		{"hat outcomes without sub die", func(b *batch) {
			v(b, spells).Details.HatManySpells.FailureTable[1].Outcomes = []ItemHatOutcome{{Min: 1, Max: 1, Result: "x"}}
		}},
		{"hat sub table gap", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[3].Outcomes[1].Min = 3 }},
		{"hat sub die sides", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[10].SubDieSides = 0 }},
		{"hat rarity table", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[10].Outcomes[2].Max = 5 }},
		{"hat duplicate fact", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[2].Facts[1].Key = "cube-size" }},
		{"hat duplicate trait", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[8].Traits[1] = "flies-out-of-hat" }},
		{"hat fact amount zero", func(b *batch) { v(b, spells).Details.HatManySpells.FailureTable[2].Facts[0].Amount = 0 }},
		{"hat duplicate sub-roll result", func(b *batch) {
			o := v(b, spells).Details.HatManySpells.FailureTable[3].Outcomes
			o[1].Result = o[0].Result
		}},
		{"hat attunement class", func(b *batch) { v(b, spells).Details.HatManySpells.AttunementClass = "cleric" }},
		{"headband score", func(b *batch) { v(b, headband).Details.Headband.Score = 20 }},
		{"headband threshold", func(b *batch) { v(b, headband).Details.Headband.NoEffectAtOrAbove = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems232(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems232(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems232RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/gem-of-seeing/variant/standard"`)
	if !bytes.Contains(items232JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items232JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems232(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"gem_of_seeing": {`, `"gem_of_seeing_plus": {`},
		"unknown detail":     {`"sense_range_feet"`, `"sense_range_all"`},
		"string for number":  {`"sense_range_feet": 120`, `"sense_range_feet": "120"`},
		"fractional number":  {`"sense_range_feet": 120`, `"sense_range_feet": 120.5`},
		"number for boolean": {`"retains_bulk_and_weight": true`, `"retains_bulk_and_weight": 1`},
		"null fact list":     {`"facts": []`, `"facts": null`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items232JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items232JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems232(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items232JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems232(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items232JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items232JSON); err == nil {
		t.Fatal("#231 decoder accepted the #232 batch")
	}
	if _, err := DecodeItems232(items231JSON); err == nil {
		t.Fatal("#232 decoder accepted the #231 batch")
	}
}
