package srd521

import (
	"fmt"
	"reflect"
	"slices"
)

// Armor and tools are mundane items with listed prices, like weapons, but a
// tool may list several variants. This file holds the identity and requirement
// checks their catalogue batches share; each payload has its own file.

// mundaneSourceManifest is a batch's reviewed identity, price and mechanic key
// manifest, independent of incoming JSON.
type mundaneSourceManifest struct {
	Batch    string
	Issue    int
	Families []mundaneSourceManifestFamily
}

type mundaneSourceManifestFamily struct {
	Slug, Name string
	Page       int
	Variants   []mundaneSourceManifestVariant
}

type mundaneSourceManifestVariant struct {
	Slug, Name string
	CostAmount int
	Currency   string
	Keys       string
}

// mundaneSourceKind names what one kind of mundane record must say: its
// category, its presentation section, its use mode and its only payload.
type mundaneSourceKind struct{ Category, Section, Use, Payload string }

// mundaneSourceRow is what a checked payload requires of its variant.
type mundaneSourceRow struct {
	Equip string
	Keys  []string
}

// decodeMundaneSourceBatch applies the shared strict decoding, then the
// manifest. row checks one payload against its family's start page.
func decodeMundaneSourceBatch[D any](data []byte, m mundaneSourceManifest, kind mundaneSourceKind, row func(d *D, page int) (mundaneSourceRow, error)) (ItemSourceBatch[D], error) {
	b, err := strictItemSourceBatch[D](data)
	if err != nil {
		return ItemSourceBatch[D]{}, err
	}
	if err := validateMundaneSourceBatch(b, m, kind, row); err != nil {
		return ItemSourceBatch[D]{}, err
	}
	return b, nil
}

func validateMundaneSourceBatch[D any](b ItemSourceBatch[D], m mundaneSourceManifest, kind mundaneSourceKind, row func(d *D, page int) (mundaneSourceRow, error)) error {
	if b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		return fmt.Errorf("unsupported item contract version")
	}
	if err := validateSource(b.Source); err != nil {
		return err
	}
	if b.Batch != m.Batch || b.CatalogueIssue != m.Issue || len(b.Families) != len(m.Families) {
		return fmt.Errorf("%s batch %q does not match its reviewed manifest", kind.Category, b.Batch)
	}
	for i, f := range b.Families {
		want := m.Families[i]
		id := "srd-5.2.1:content/items/" + want.Slug
		pages := []int{want.Page, want.Page}
		if f.ID != id || f.Name != want.Name || f.Kind != "mundane" || f.Category != kind.Category {
			return fmt.Errorf("unreviewed, missing or reordered %s %q", kind.Category, f.ID)
		}
		if f.Base == nil || !sameItemKeys(f.Base.Options, "") || !sameItemKeys(f.Base.Excluded, "") {
			return fmt.Errorf("%s: a mundane item has no base item restriction", f.Name)
		}
		if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/"+kind.Section+"/"+want.Slug || !slices.Equal(f.Presentations[0].Pages, pages) {
			return fmt.Errorf("%s: invalid, missing or duplicate source presentation", f.Name)
		}
		if len(f.Variants) != len(want.Variants) {
			return fmt.Errorf("%s: missing or unreviewed source-listed variants", f.Name)
		}
		for j, v := range f.Variants {
			wv := want.Variants[j]
			if v.ID != id+"/variant/"+wv.Slug || v.Name != wv.Name || v.Rarity != "none" {
				return fmt.Errorf("%s: invalid, duplicate or unreviewed variant %q", f.Name, v.ID)
			}
			if v.Attunement == nil || v.Attunement.Mode != "none" || !sameItemKeys(v.Attunement.Restrictions, "") {
				return fmt.Errorf("%s: a mundane item has no attunement", v.Name)
			}
			if v.Charges == nil || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
				return fmt.Errorf("%s: a mundane item has no charges", v.Name)
			}
			if v.Cost == nil || v.Cost.Kind != "listed" || v.Cost.Amount != wv.CostAmount || v.Cost.Currency != wv.Currency || !slices.Equal(v.Cost.SourcePages, pages) {
				return fmt.Errorf("%s: listed price disagrees with source", v.Name)
			}
			if v.Details == nil {
				return fmt.Errorf("%s: source details required", v.Name)
			}
			payloads, err := itemDetailPayloads(reflect.ValueOf(v.Details).Elem())
			if err != nil {
				return fmt.Errorf("%s: %w", v.Name, err)
			}
			if !slices.Equal(payloads, []string{kind.Payload}) {
				return fmt.Errorf("%s: details must hold exactly the %q payload", v.Name, kind.Payload)
			}
			checked, err := row(v.Details, want.Page)
			if err != nil {
				return fmt.Errorf("%s: %w", v.Name, err)
			}
			if v.Equip == nil || v.Equip.Mode != checked.Equip || !sameItemKeys(v.Equip.Restrictions, "") {
				return fmt.Errorf("%s: equip restriction disagrees with source", v.Name)
			}
			if v.Use == nil || v.Use.Mode != kind.Use || !sameItemKeys(v.Use.Restrictions, "") {
				return fmt.Errorf("%s: use restriction disagrees with source", v.Name)
			}
			if !sameItemKeys(v.RequiredMechanicKeys, wv.Keys) || len(v.RequiredMechanicKeys) == 0 {
				return fmt.Errorf("%s: missing, unknown or reordered mechanic key", v.Name)
			}
			for _, key := range v.RequiredMechanicKeys {
				if !itemMechanicKey.MatchString(key) {
					return fmt.Errorf("%s: invalid mechanic key %q", v.Name, key)
				}
			}
			// The keys restate the payload's values, so the two cannot disagree.
			if !slices.Equal(v.RequiredMechanicKeys, checked.Keys) {
				return fmt.Errorf("%s: values disagree with the mechanic keys: %v", v.Name, checked.Keys)
			}
		}
	}
	return nil
}
