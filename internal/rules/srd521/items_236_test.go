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

func loadItems236(t *testing.T) ItemSourceBatch[ItemDetails236] {
	t.Helper()
	b, err := LoadItems236()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 233-235, in planning inventory order for issue #236. No family
// prints more than one form, so each has the single variant "standard".
func TestItems236SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		rarity     string
		charged    bool
	}{
		{"Oil of Sharpness", "oil-of-sharpness", 233, "potion", nil, "none", "very-rare", false},
		{"Oil of Slipperiness", "oil-of-slipperiness", 233, "potion", nil, "none", "uncommon", false},
		{"Pearl of Power", "pearl-of-power", 234, "wondrous-item", nil, "required", "uncommon", false},
		{"Periapt of Health", "periapt-of-health", 234, "wondrous-item", nil, "required", "uncommon", false},
		{"Periapt of Proof against Poison", "periapt-of-proof-against-poison", 234, "wondrous-item", nil, "required", "rare", false},
		{"Periapt of Wound Closure", "periapt-of-wound-closure", 234, "wondrous-item", nil, "required", "uncommon", false},
		{"Philter of Love", "philter-of-love", 234, "potion", nil, "none", "uncommon", false},
		{"Pipes of Haunting", "pipes-of-haunting", 234, "wondrous-item", nil, "none", "uncommon", true},
		{"Pipes of the Sewers", "pipes-of-the-sewers", 234, "wondrous-item", nil, "required", "uncommon", true},
		{"Plate Armor of Etherealness", "plate-armor-of-etherealness", 234, "armor", []string{"half-plate-armor", "plate-armor"}, "required", "legendary", false},
		{"Portable Hole", "portable-hole", 234, "wondrous-item", nil, "none", "rare", false},
		{"Potion of Animal Friendship", "potion-of-animal-friendship", 235, "potion", nil, "none", "uncommon", false},
	}
	b := loadItems236(t)
	if len(b.Families) != len(expected) {
		t.Fatalf("expected %d families, got %d", len(expected), len(b.Families))
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
			// Only the two pipes list charges: 3, regained 1d3 at dawn.
			wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
			if want.charged {
				wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}
			}
			if *v.Charges != wantCharges {
				t.Fatalf("charges mismatch: %+v", *v.Charges)
			}
			if v.MechanicsAvailable() {
				t.Fatal("catalogue record became available")
			}
			total++
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
		if r.Issue != 236 {
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
		t.Fatalf("inventory assigns %d families to issue #236", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "233, 234, 235" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.Sharpness.Appearance = "changed"
	if fresh := loadItems236(t); fresh.Families[0].Variants[0].Details.Sharpness.Appearance != "clear-gelatinous-oil-with-ultrathin-silver-shards" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 233-235. The Oil of Slipperiness continues from p. 233 to p. 234, the Pipes of
// Haunting across the columns of p. 234 and the Portable Hole from p. 234 to p. 235.
func TestItems236SourceValues(t *testing.T) {
	expected := []string{
		`{"oil_of_sharpness":{"coats_melee_weapons":1,"coats_ammunition":20,"requires_nonmagical":true,"damage_types":["slashing","piercing"],"application_minutes":1,"weapon_bonus":3,"ammunition_bonus":3,"appearance":"clear-gelatinous-oil-with-ultrathin-silver-shards"}}`,
		`{"oil_of_slipperiness":{"source_pages":[233,234],"coverage":{"max_size":"medium","includes_worn_and_carried_equipment":true,"extra_vial_per_size_above":"medium"},"application_minutes":10,"body_spell":"freedom-of-movement","body_duration_hours":8,"pour_action":"magic","pour_area_square_feet":10,"pour_spell":"grease","pour_duration_hours":8,"appearance":"sticky-black-unguent-thick-and-heavy-flows-quickly-when-poured"}}`,
		`{"pearl_of_power":{"attunement_by_any_of":["spellcaster"],"condition":"on-your-person","action":"magic","spell_slots_regained":1,"max_slot_level":3,"uses":1,"recharge":"next-dawn"}}`,
		`{"periapt_of_health":{"condition":"while-wearing","action":"magic","healing_dice":{"dice":2,"sides":4,"type":"healing"},"healing_bonus":2,"uses":1,"recharge":"next-dawn","poisoned_save_mode":"advantage","poisoned_save_applies_to":["avoid-poisoned","end-poisoned"]}}`,
		`{"periapt_of_proof_against_poison":{"appearance":"delicate-silver-chain-with-brilliant-cut-black-gem-pendant","condition":"while-wearing","condition_immunities":["poisoned"],"damage_immunities":["poison"]}}`,
		`{"periapt_of_wound_closure":{"condition":"while-wearing","life_preservation":{"trigger":"make-death-saving-throw","changes_roll_at_or_below":9,"changed_to":10,"result":"failed-save-becomes-successful"},"natural_healing_boost":{"trigger":"roll-hit-point-die-to-regain-hit-points","multiplier":2}}}`,
		`{"philter_of_love":{"trigger":"next-creature-you-see-after-drinking","window_minutes":10,"condition":"charmed","charmed_by":"that-creature","duration_hours":1,"appearance":"rose-hued-effervescent-liquid-with-one-heart-shaped-bubble"}}`,
		`{"pipes_of_haunting":{"action":"magic","charge_cost":1,"effect":"eerie-spellbinding-tune","targets":"each-creature-of-your-choice","range_feet":30,"save":{"ability":"wisdom","dc":15},"condition":"frightened","duration_seconds":60,"repeat_save":"end-of-each-of-its-turns","success_ends":"effect-on-itself","immunity_hours_after_success":24}}`,
		`{"pipes_of_the_sewers":{"passive":{"condition":"on-your-person","creatures":["ordinary-rat","giant-rat"],"attitude":"indifferent","unless":"you-threaten-or-harm-them"},"call_swarms":{"play_action":"magic","call_action":"bonus","min_charges":1,"max_charges":3,"swarm":"swarm-of-rats","swarms_per_charge":1,"range":"half-a-mile","insufficient_rats":"charge-is-wasted","availability_determined_by":"gm","swarm_movement":"toward-music-by-shortest-available-route-not-otherwise-controlled"},"sway_swarms":{"trigger":"uncontrolled-swarm-of-rats-within-30-feet-while-you-play","trigger_range_feet":30,"save":{"ability":"wisdom","dc":15},"success_result":"behaves-normally-and-cannot-be-swayed","success_immunity_hours":24,"failure_result":"friendly-to-you-and-your-allies","failure_lasts_while":"you-play-the-pipes-each-round-as-magic-action","friendly_swarm_obeys_commands":true,"uncommanded_behavior":"defends-itself-and-takes-no-other-actions","control_ends_when":"friendly-swarm-starts-turn-more-than-30-feet-away","control_end_distance_feet":30,"after_control_ends":"behaves-normally-and-cannot-be-swayed-for-24-hours"}}}`,
		`{"plate_armor_of_etherealness":{"condition":"while-wearing","action":"magic","activation":"command-word","spell":"etherealness","ends_on":["armor-removed","magic-action-repeat-command-word"],"uses":1,"recharge":"next-dawn"}}`,
		`{"portable_hole":{"source_pages":[234,235],"cloth":"fine-black-silk-soft-folds-to-handkerchief-size","diameter_feet":6,"unfold":{"action":"magic","surface":"on-or-against-solid-surface","depth_feet":10,"shape":"cylindrical","plane":"different-plane-of-existence","cannot_create_open_passages":true,"exit":"creature-inside-open-hole-can-climb-out"},"close":{"action":"magic","contents":"creatures-and-objects-remain-in-extradimensional-space","weight":"next-to-nothing"},"escape_when_folded":{"action":"action","check":{"ability":"strength","skill":"athletics","dc":10},"success":"forces-its-way-out","appears_within_feet":5,"closed_air_hours":1,"air_divided_by":"number-of-breathing-creatures-inside"},"inside_extradimensional_space":{"containers":["bag-of-holding","handy-haversack","similar-item"],"result":"both-items-destroyed-and-gate-to-astral-plane-opens","gate_origin":"where-one-item-was-placed-inside-the-other","pulls_in_within_feet":10,"exempt":"creature-behind-total-cover","destination":"random-location-on-astral-plane","gate_closes_afterward":true,"one_way":true,"can_be_reopened":false}}}`,
		`{"potion_of_animal_friendship":{"trigger":"drink-this-potion","spell":"animal-friendship","spell_level":3,"save_dc":13,"appearance":"muddy-liquid","appearance_items":["fish-scale","hummingbird-feather","cat-claw","squirrel-hair"]}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:apply:one-vial-coats-one-melee-weapon-or-twenty-ammunition nonmagical-slashing-or-piercing-only:item.consume action.apply-1-minute target.nonmagical-slashing-or-piercing effect.plus-3-weapon effect.plus-3-ammunition",
		"none:apply:one-vial-covers-one-medium-or-smaller-creature pour-as-magic-action:item.consume action.apply-10-minutes spell.effect-freedom-of-movement duration.8-hours action.magic area.square-10-foot spell.effect-grease",
		"carry:magic-action:once-until-next-dawn:item.attunement item.carry action.magic spell-slot.regain-one-level-3-or-lower item.once-until-next-dawn",
		"wear:magic-action:once-until-next-dawn:item.attunement item.wear action.magic healing.hit-points-2d4-plus-2 item.once-until-next-dawn save.advantage-avoid-or-end-poisoned",
		"wear:none::item.attunement item.wear immunity.condition-poisoned immunity.damage-poison",
		"wear:none::item.attunement item.wear death-save.change-9-or-lower-to-10 healing.double-hit-point-die",
		"none:consume:drink-one:item.consume condition.charmed duration.1-hour window.10-minutes-after-drinking",
		"none:magic-action:expend-one-charge:item.charges action.magic area.30-feet save.wisdom-dc-15 condition.frightened duration.1-minute save.repeat-end-of-turn immunity.24-hours-after-success",
		"carry:magic-action:bonus-action-expend-one-to-three-charges:item.attunement item.carry item.charges attitude.rats-indifferent action.magic action.bonus summon.swarm-of-rats-per-charge save.wisdom-dc-15 condition.friendly immunity.24-hours-after-success",
		"wear:magic-action:command-word once-until-next-dawn:item.attunement item.wear action.magic item.command-word spell.effect-etherealness item.once-until-next-dawn",
		"none:magic-action:unfold-or-fold:action.magic item.extradimensional-hole-10-feet-deep check.strength-athletics-dc-10 item.air-1-hour-divided hazard.nested-extradimensional-space-destroys-both-and-opens-astral-gate",
		"none:consume:drink-one:item.consume spell.cast-animal-friendship spell.level-3-version spell.save-dc-13",
	}
	index := 0
	for _, f := range loadItems236(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails236
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

func TestItems236CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems236(t)
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
		earlier, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, earlier)
	}
	composed, err := ComposeItemSourceIdentities(append(identities, batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 97 || variants != 135 || composed[0].Name != "Potion of Healing" || composed[84].Name != "Gem of Brightness" || composed[85].Name != "Oil of Sharpness" || composed[96].Name != "Potion of Animal Friendship" {
		t.Fatal("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #236")
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
			if strings.Contains(fixture, "items_236") {
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
func TestItems236RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items236JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 500 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items236JSON, &fresh); err != nil {
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
		if _, err := DecodeItems236(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems236RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails236]
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails236] { return &b.Families[family].Variants[0] }
	const (
		sharpness, slipperiness, pearl, health, proof, closure, philter, haunting, sewers, plate, hole, animal = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-11" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 235 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/potion-of-clairvoyance"
			b.Families[12].Name = "Potion of Clairvoyance"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/oil-of-etherealness" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Oil of Sharpening" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[sharpness].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) { b.Families[plate].Base.Options = append(b.Families[plate].Base.Options, "chain-mail") }},
		{"dropped base", func(b *batch) { b.Families[plate].Base.Options = []string{"plate-armor"} }},
		{"invented base", func(b *batch) { b.Families[hole].Base.Options = []string{"any-armor"} }},
		{"excluded base", func(b *batch) { b.Families[plate].Base.Excluded = []string{"half-plate-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{234, 234} }},
		{"page range", func(b *batch) { b.Families[slipperiness].Presentations[0].Pages = []int{233, 234} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Oil" }},
		{"rarity", func(b *batch) { v(b, sharpness).Rarity = "rare" }},
		{"plate rarity", func(b *batch) { v(b, plate).Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"invented attunement", func(b *batch) { v(b, sharpness).Attunement.Mode = "required" }},
		{"dropped attunement", func(b *batch) { v(b, pearl).Attunement.Mode = "none" }},
		{"dropped sewers attunement", func(b *batch) { v(b, sewers).Attunement.Mode = "none" }},
		{"haunting attunement", func(b *batch) { v(b, haunting).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, pearl).Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { v(b, haunting).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, haunting).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"charge maximum", func(b *batch) { v(b, haunting).Charges.Maximum = 4 }},
		{"sewers charge recharge", func(b *batch) { v(b, sewers).Charges.Recharge = "none" }},
		{"sewers per use", func(b *batch) { v(b, sewers).Charges.PerUse = 3 }},
		{"invented charges", func(b *batch) { v(b, pearl).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, plate).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, plate).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, plate).Use.Mode = "bonus-action" }},
		{"dropped use restriction", func(b *batch) { v(b, plate).Use.Restrictions = []string{"command-word"} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails236{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.Slipperiness = v(b, 1).Details.Slipperiness }},
		{"blank key", func(b *batch) { v(b, sharpness).Details.Sharpness.Appearance = " " }},
		{"prose value", func(b *batch) { v(b, sharpness).Details.Sharpness.Appearance = "A clear oil" }},
		{"negative value", func(b *batch) { v(b, sharpness).Details.Sharpness.ApplicationMinutes = -1 }},
		{"sharpness weapon bonus", func(b *batch) { v(b, sharpness).Details.Sharpness.WeaponBonus = 2 }},
		{"sharpness ammunition bonus", func(b *batch) { v(b, sharpness).Details.Sharpness.AmmunitionBonus = 1 }},
		{"sharpness ammunition count", func(b *batch) { v(b, sharpness).Details.Sharpness.CoatsAmmunition = 0 }},
		{"sharpness damage types", func(b *batch) { v(b, sharpness).Details.Sharpness.DamageTypes = []string{} }},
		{"sharpness duplicate damage type", func(b *batch) { v(b, sharpness).Details.Sharpness.DamageTypes[1] = "slashing" }},
		{"slipperiness pages", func(b *batch) { v(b, slipperiness).Details.Slipperiness.SourcePages = []int{233} }},
		{"slipperiness application", func(b *batch) { v(b, slipperiness).Details.Slipperiness.ApplicationMinutes = 1 }},
		{"slipperiness body duration", func(b *batch) { v(b, slipperiness).Details.Slipperiness.BodyDurationHours = 1 }},
		{"slipperiness pour duration", func(b *batch) { v(b, slipperiness).Details.Slipperiness.PourDurationHours = 0 }},
		{"slipperiness pour area", func(b *batch) { v(b, slipperiness).Details.Slipperiness.PourAreaFeet = 0 }},
		{"pearl slot level", func(b *batch) { v(b, pearl).Details.Pearl.MaxSlotLevel = 4 }},
		{"pearl slots", func(b *batch) { v(b, pearl).Details.Pearl.SlotsRegained = 0 }},
		{"pearl no attunement qualifier", func(b *batch) { v(b, pearl).Details.Pearl.AttunementBy = []string{} }},
		{"pearl duplicate qualifier", func(b *batch) { v(b, pearl).Details.Pearl.AttunementBy = []string{"spellcaster", "spellcaster"} }},
		{"health dice", func(b *batch) { v(b, health).Details.Health.Healing.Dice = 4 }},
		{"health bonus", func(b *batch) { v(b, health).Details.Health.HealingBonus = 0 }},
		{"health save applications", func(b *batch) { v(b, health).Details.Health.SaveApplication = []string{} }},
		{"health duplicate save application", func(b *batch) { v(b, health).Details.Health.SaveApplication[1] = "avoid-poisoned" }},
		{"proof no condition immunity", func(b *batch) { v(b, proof).Details.Proof.ConditionImmunity = []string{} }},
		{"proof no damage immunity", func(b *batch) { v(b, proof).Details.Proof.DamageTypeImmunity = []string{} }},
		{"closure threshold", func(b *batch) { v(b, closure).Details.WoundClosure.Preservation.RollAtOrBelow = 10 }},
		{"closure changed roll", func(b *batch) { v(b, closure).Details.WoundClosure.Preservation.ChangedTo = 9 }},
		{"closure multiplier", func(b *batch) { v(b, closure).Details.WoundClosure.HealingBoost.Multiplier = 3 }},
		{"philter window", func(b *batch) { v(b, philter).Details.Philter.WindowMinutes = 1 }},
		{"philter duration", func(b *batch) { v(b, philter).Details.Philter.DurationHours = 10 }},
		{"haunting dc", func(b *batch) { v(b, haunting).Details.Haunting.Save.DC = 13 }},
		{"haunting charge cost", func(b *batch) { v(b, haunting).Details.Haunting.ChargeCost = 2 }},
		{"haunting range", func(b *batch) { v(b, haunting).Details.Haunting.RangeFeet = 0 }},
		{"haunting duration", func(b *batch) { v(b, haunting).Details.Haunting.DurationSeconds = 0 }},
		{"haunting immunity", func(b *batch) { v(b, haunting).Details.Haunting.ImmunityHours = 0 }},
		{"sewers creatures", func(b *batch) { v(b, sewers).Details.Sewers.Passive.Creatures = []string{} }},
		{"sewers duplicate creature", func(b *batch) { v(b, sewers).Details.Sewers.Passive.Creatures[1] = "ordinary-rat" }},
		{"sewers charge range", func(b *batch) { v(b, sewers).Details.Sewers.Call.MaxCharges = 0 }},
		{"sewers inverted charge range", func(b *batch) { v(b, sewers).Details.Sewers.Call.MinCharges = 4 }},
		{"sewers dc", func(b *batch) { v(b, sewers).Details.Sewers.Sway.Save.DC = 13 }},
		{"sewers sway range", func(b *batch) { v(b, sewers).Details.Sewers.Sway.TriggerFeet = 0 }},
		{"sewers immunity", func(b *batch) { v(b, sewers).Details.Sewers.Sway.ImmunityHours = 0 }},
		{"plate no end", func(b *batch) { v(b, plate).Details.PlateArmor.EndsOn = []string{} }},
		{"plate duplicate end", func(b *batch) { v(b, plate).Details.PlateArmor.EndsOn[1] = "armor-removed" }},
		{"plate uses", func(b *batch) { v(b, plate).Details.PlateArmor.UseLimit = 0 }},
		{"hole pages", func(b *batch) { v(b, hole).Details.PortableHole.SourcePages = []int{234} }},
		{"hole diameter", func(b *batch) { v(b, hole).Details.PortableHole.DiameterFt = 0 }},
		{"hole depth", func(b *batch) { v(b, hole).Details.PortableHole.Unfold.DepthFeet = 5 }},
		{"hole escape dc", func(b *batch) { v(b, hole).Details.PortableHole.Escape.Check.DC = 15 }},
		{"hole air", func(b *batch) { v(b, hole).Details.PortableHole.Escape.AirHours = 0 }},
		{"hole containers", func(b *batch) { v(b, hole).Details.PortableHole.Nesting.Containers = []string{} }},
		{"hole duplicate container", func(b *batch) { v(b, hole).Details.PortableHole.Nesting.Containers[1] = "bag-of-holding" }},
		{"hole gate range", func(b *batch) { v(b, hole).Details.PortableHole.Nesting.PullFeet = 0 }},
		{"animal level", func(b *batch) { v(b, animal).Details.AnimalFriend.SpellLevel = 1 }},
		{"animal dc", func(b *batch) { v(b, animal).Details.AnimalFriend.SaveDC = 15 }},
		{"animal no appearance items", func(b *batch) { v(b, animal).Details.AnimalFriend.AppearanceItems = []string{} }},
		{"animal duplicate appearance item", func(b *batch) { v(b, animal).Details.AnimalFriend.AppearanceItems[1] = "fish-scale" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems236(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems236(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems236RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/oil-of-sharpness/variant/standard"`)
	if !bytes.Contains(items236JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items236JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems236(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"oil_of_sharpness": {`, `"oil_of_sharpness_plus": {`},
		"unknown detail":     {`"weapon_bonus"`, `"weapon_bonus_all"`},
		"string for number":  {`"save_dc": 13`, `"save_dc": "13"`},
		"fractional number":  {`"save_dc": 13`, `"save_dc": 13.5`},
		"number for boolean": {`"requires_nonmagical": true`, `"requires_nonmagical": 1`},
		"null damage types":  {`"damage_types": [`, `"damage_types": null, "x": [`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items236JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items236JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems236(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items236JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems236(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items236JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items236JSON); err == nil {
		t.Fatal("#231 decoder accepted the #236 batch")
	}
	if _, err := DecodeItems236(items231JSON); err == nil {
		t.Fatal("#236 decoder accepted the #231 batch")
	}
}
