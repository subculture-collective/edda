package srd521

import (
	"bytes"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"
)

// itemBatchLoaders lists every catalogue batch in issue order, from #225.
var itemBatchLoaders = []func() ([]ItemSourceIdentity, error){
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
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems245(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems246(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems247(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems248(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems249(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems250(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems251(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems252(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems253(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems254(); return b.SourceIdentities(), err },
	func() ([]ItemSourceIdentity, error) { b, err := LoadItems255(); return b.SourceIdentities(), err },
}

// itemIdentitiesThrough returns the Healing contract and batches #225 to issue.
func itemIdentitiesThrough(t *testing.T, issue int) [][]ItemSourceIdentity {
	t.Helper()
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	identities := [][]ItemSourceIdentity{contract.SourceIdentities()}
	for _, load := range itemBatchLoaders[:issue-224] {
		ids, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, ids)
	}
	return identities
}

// TestItemsBatches225To255ComposeTogether composes the Healing contract with
// every catalogue batch from #225 to #255, magic items, weapons, armor and
// tools, and checks that no family, presentation, variant or name is shared.
func TestItemsBatches225To255ComposeTogether(t *testing.T) {
	composed, err := ComposeItemSourceIdentities(itemIdentitiesThrough(t, 255)...)
	if err != nil {
		t.Fatal(err)
	}
	variants, sections := 0, map[string]int{}
	for _, f := range composed {
		variants += len(f.VariantIDs)
		for _, p := range f.PresentationIDs {
			section, _, _ := strings.Cut(strings.TrimPrefix(p, prefix), "/")
			sections[section]++
		}
	}
	if len(composed) != 334 || variants != 422 || sections["weapons"] != 38 || sections["armor"] != 13 || sections["tools"] != 25 {
		t.Fatalf("unexpected combined composition: %d families, %d variants, %v", len(composed), variants, sections)
	}

	// Every armor and tool record in the planning inventory has a catalogue
	// record, and so does every weapon in it.
	presented := map[string]bool{}
	for _, f := range composed {
		for _, p := range f.PresentationIDs {
			presented[p] = true
		}
	}
	listed := 0
	for _, r := range itemInventory(t) {
		if r.Family != "weapons" && r.Family != "armor" && r.Family != "tools" {
			continue
		}
		listed++
		if !presented[r.ID] {
			t.Fatalf("inventory record %s has no catalogue record", r.ID)
		}
	}
	// The weapons table prints one more weapon than the inventory lists.
	if listed != 75 || !presented["srd-5.2.1:content/weapons/heavy-crossbow"] {
		t.Fatalf("expected 75 inventory records and Heavy Crossbow, got %d records", listed)
	}
}

type itemInventoryRecord struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Family string `json:"family"`
	Page   int    `json:"source_page"`
	Batch  string `json:"batch_key"`
	Issue  int    `json:"catalogue_issue"`
	Status string `json:"status"`
}

func itemInventory(t *testing.T) []itemInventoryRecord {
	t.Helper()
	data, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []itemInventoryRecord `json:"records"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	return inventory.Records
}

// joinItemInventory joins the inventory records of one issue, in order, to the
// batch identities. pages are the reviewed inventory pages.
func joinItemInventory(t *testing.T, ids []ItemSourceIdentity, batch string, issue int, family string, pages []int) {
	t.Helper()
	index := 0
	for _, r := range itemInventory(t) {
		if r.Issue != issue {
			continue
		}
		if r.Family != family || r.Batch != batch || index >= len(ids) || index >= len(pages) {
			t.Fatalf("inventory record %s changed family or batch", r.ID)
		}
		if r.Status != "unsupported" || r.ID != ids[index].PresentationIDs[0] || r.Name != ids[index].Name || r.Page != pages[index] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, ids[index].Name)
		}
		index++
	}
	if index != len(ids) {
		t.Fatalf("inventory assigns %d records to issue #%d, batch holds %d", index, issue, len(ids))
	}
}

func requireInventoryOmits(t *testing.T, name string) {
	t.Helper()
	for _, r := range itemInventory(t) {
		if r.Name == name {
			t.Fatalf("the inventory now lists %s: join it instead", name)
		}
	}
}

type itemVariantWant struct {
	slug, name string
	cost       int
	currency   string
}

type itemFamilyWant struct {
	name, slug string
	page       int
	variants   []itemVariantWant
}

// checkMundaneIdentities compares armor, tool and weapon families with literal
// source identities and listed prices.
func checkMundaneIdentities[D any](t *testing.T, families []ItemSourceFamily[D], category, section string, expected []itemFamilyWant) {
	t.Helper()
	if len(families) != len(expected) {
		t.Fatalf("expected %d families, got %d", len(expected), len(families))
	}
	for i, want := range expected {
		f := families[i]
		id := "srd-5.2.1:content/items/" + want.slug
		pages := []int{want.page, want.page}
		if f.ID != id || f.Name != want.name || f.Kind != "mundane" || f.Category != category || len(f.Base.Options) != 0 || len(f.Base.Excluded) != 0 {
			t.Fatalf("source identity mismatch: %s %s", f.ID, f.Name)
		}
		if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/"+section+"/"+want.slug || !slices.Equal(f.Presentations[0].Pages, pages) {
			t.Fatalf("%s: source presentation mismatch", f.Name)
		}
		if len(f.Variants) != len(want.variants) {
			t.Fatalf("%s: expected %d variants, got %d", f.Name, len(want.variants), len(f.Variants))
		}
		for j, wv := range want.variants {
			v := f.Variants[j]
			if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != "none" || v.Attunement.Mode != "none" || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
				t.Fatalf("%s: variant mismatch: %s %s", f.Name, v.ID, v.Name)
			}
			if v.Cost.Kind != "listed" || v.Cost.Amount != wv.cost || v.Cost.Currency != wv.currency || !slices.Equal(v.Cost.SourcePages, pages) {
				t.Fatalf("%s: listed price mismatch: %+v", v.Name, v.Cost)
			}
			if v.MechanicsAvailable() {
				t.Fatalf("%s: catalogue record became available", v.Name)
			}
		}
	}
}

// checkItemBatchComposition rejects shared identities and confirms that the
// batch added no coverage fixture and moved no coverage boundary.
func checkItemBatchComposition(t *testing.T, ids []ItemSourceIdentity, section, fixtureName string) {
	t.Helper()
	if _, err := ComposeItemSourceIdentities(ids, ids); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := slices.Clone(ids[:1])
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/" + section + "/other"}
	if _, err := ComposeItemSourceIdentities(ids, clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, fixtureName) {
				t.Fatal("catalogue batch must not be a coverage fixture")
			}
		}
	}
	if report := register.Coverage(); report.PlannedSpellItemScopes != 753 || report.FullSupport {
		t.Fatal("catalogue data changed the coverage boundary")
	}
}

// rejectMissingItemMembers deletes each object member in turn; every decode
// must fail. minimum guards against a walk that skips part of the document.
func rejectMissingItemMembers(t *testing.T, source []byte, minimum int, decode func([]byte) error) {
	t.Helper()
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
	if err := json.Unmarshal(source, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) != minimum {
		t.Fatalf("expected to walk %d members, got %d", minimum, len(paths))
	}
	if err := decode(source); err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(source, &fresh); err != nil {
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
		if err := decode(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

type itemJSONChange struct{ name, old, replacement string }

// rejectUnsupportedItemJSON inserts runtime grants and aliases before anchor,
// applies the batch's own replacements and tries ambiguous documents.
func rejectUnsupportedItemJSON(t *testing.T, source []byte, anchor string, decode func([]byte) error, changes []itemJSONChange) {
	t.Helper()
	if !bytes.Contains(source, []byte(anchor)) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"proficient":true,`, `"playable":true,`, `"available":true,`, `"properties":{"attack_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(source, []byte(anchor), []byte(field+anchor), 1)
			if err := decode(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	changes = append(changes,
		itemJSONChange{"null details", `"details": {`, `"details": null, "x": {`},
		itemJSONChange{"null list", `"excluded": []`, `"excluded": null`})
	for _, c := range changes {
		t.Run(c.name, func(t *testing.T) {
			if !bytes.Contains(source, []byte(c.old)) {
				t.Fatal("replacement anchor missing")
			}
			if err := decode(bytes.Replace(source, []byte(c.old), []byte(c.replacement), 1)); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(source), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if err := decode(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
}
