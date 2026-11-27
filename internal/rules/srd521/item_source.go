package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// ItemSourceBatch is one catalogue batch of descriptive item families. D is the
// batch's own typed details struct, so a batch adds source payloads in its own
// files without editing a shared record type. It is never inventory or rules.
type ItemSourceBatch[D any] struct {
	SchemaVersion  int                   `json:"schema_version"`
	RulesVersion   string                `json:"rules_version"`
	Batch          string                `json:"batch"`
	CatalogueIssue int                   `json:"catalogue_issue"`
	Source         Source                `json:"source"`
	Families       []ItemSourceFamily[D] `json:"families"`
}

// ItemSourceFamily nests every source-listed variant under one canonical
// identity. Presentations reference planning inventory IDs, not owned items.
type ItemSourceFamily[D any] struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Kind          string                 `json:"kind"`
	Category      string                 `json:"category"`
	Base          *ItemBase              `json:"base"`
	Presentations []ItemReference        `json:"presentations"`
	Variants      []ItemSourceVariant[D] `json:"variants"`
}

// ItemBase records the base items the source category line allows or excludes.
// Both lists are empty for a Wondrous Item.
type ItemBase struct {
	Options  []string `json:"options"`
	Excluded []string `json:"excluded"`
}

// ItemSourceVariant keeps the item contract's requirement fields and replaces
// the Healing-only dice with a typed per-family details payload. It holds no
// runtime quantity, current charges, AC or action grants.
type ItemSourceVariant[D any] struct {
	ID                   string           `json:"id"`
	Name                 string           `json:"name"`
	Rarity               string           `json:"rarity"`
	Attunement           *ItemAttunement  `json:"attunement"`
	Charges              *ItemCharges     `json:"charges"`
	Equip                *ItemRestriction `json:"equip"`
	Use                  *ItemRestriction `json:"use"`
	Cost                 *ItemCost        `json:"cost"`
	RequiredMechanicKeys []string         `json:"required_mechanic_keys"`
	Details              *D               `json:"details"`
}

// MechanicsAvailable is always false: no resolver reads an item source batch,
// so requiring a mechanic key never registers, implements or enables it.
func (ItemSourceVariant[D]) MechanicsAvailable() bool { return false }

// ItemDice records source dice only; it is not passed to a resolver.
type ItemDice struct {
	Dice  int    `json:"dice"`
	Sides int    `json:"sides"`
	Type  string `json:"type"`
}

// ItemCheck records a source ability check and its DC.
type ItemCheck struct {
	Ability string `json:"ability"`
	Skill   string `json:"skill"`
	DC      int    `json:"dc"`
}

// ItemSourceIdentity is the batch-independent identity used to compose batches.
type ItemSourceIdentity struct {
	FamilyID        string
	Name            string
	PresentationIDs []string
	VariantIDs      []string
}

// SourceIdentities lists the batch's family, presentation and variant IDs.
func (b ItemSourceBatch[D]) SourceIdentities() []ItemSourceIdentity {
	out := make([]ItemSourceIdentity, 0, len(b.Families))
	for _, f := range b.Families {
		id := ItemSourceIdentity{FamilyID: f.ID, Name: f.Name}
		for _, p := range f.Presentations {
			id.PresentationIDs = append(id.PresentationIDs, p.ID)
		}
		for _, v := range f.Variants {
			id.VariantIDs = append(id.VariantIDs, v.ID)
		}
		out = append(out, id)
	}
	return out
}

// SourceIdentities lists the contract fixture's identities for composition.
func (c ItemCatalogue) SourceIdentities() []ItemSourceIdentity {
	out := make([]ItemSourceIdentity, 0, len(c.Families))
	for _, f := range c.Families {
		id := ItemSourceIdentity{FamilyID: f.ID, Name: f.Name}
		for _, p := range f.Presentations {
			id.PresentationIDs = append(id.PresentationIDs, p.ID)
		}
		for _, v := range f.Variants {
			id.VariantIDs = append(id.VariantIDs, v.ID)
		}
		out = append(out, id)
	}
	return out
}

// ComposeItemSourceIdentities joins validated batches in input order. A repeated
// family, presentation, variant or name is an error, never a last-wins merge.
// The result is an identity index; it grants nothing.
func ComposeItemSourceIdentities(batches ...[]ItemSourceIdentity) ([]ItemSourceIdentity, error) {
	if len(batches) == 0 {
		return nil, fmt.Errorf("item source batches required")
	}
	seen := map[string]bool{}
	claim := func(kind, id string) error {
		if id == "" || seen[kind+"\x00"+id] {
			return fmt.Errorf("duplicate or empty item source %s %q", kind, id)
		}
		seen[kind+"\x00"+id] = true
		return nil
	}
	var out []ItemSourceIdentity
	for _, batch := range batches {
		for _, f := range batch {
			if len(f.PresentationIDs) == 0 || len(f.VariantIDs) == 0 {
				return nil, fmt.Errorf("item source family %q lacks presentations or variants", f.FamilyID)
			}
			if err := claim("family", f.FamilyID); err != nil {
				return nil, err
			}
			if err := claim("name", f.Name); err != nil {
				return nil, err
			}
			for _, p := range f.PresentationIDs {
				if err := claim("presentation", p); err != nil {
					return nil, err
				}
			}
			for _, v := range f.VariantIDs {
				if err := claim("variant", v); err != nil {
					return nil, err
				}
			}
			out = append(out, f)
		}
	}
	return out, nil
}

// itemSourceManifest is a batch's reviewed identity and requirement manifest.
// It lives in code, independently of incoming JSON, so omitting a family or a
// variant, or widening a requirement, cannot pass validation.
type itemSourceManifest struct {
	Batch    string
	Issue    int
	Families []itemSourceManifestFamily
}

// itemSourceManifestFamily pins the fields shared by a family's variants.
// List fields are space-separated keys in source order.
type itemSourceManifestFamily struct {
	Slug, Name      string
	Page            int
	Category        string
	Options         string
	Excluded        string
	Attunement      string
	Equip, Use      string
	UseRestrictions string
	Charges         *ItemCharges // nil: the family has no charge pool
	Keys            string
	Payload         string
	Variants        []itemSourceManifestVariant
}

type itemSourceManifestVariant struct{ Slug, Name, Rarity string }

var (
	itemToken       = regexp.MustCompile(`^[a-z0-9]+([.-][a-z0-9]+)*$`)
	itemMechanicKey = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*\.[a-z0-9]+(-[a-z0-9]+)*$`)
)

const (
	itemCategories  = "armor|potion|ring|rod|scroll|staff|wand|weapon|wondrous-item"
	itemRarities    = "common|uncommon|rare|very-rare|legendary|artifact"
	itemAttunements = "none|required"
)

// decodeItemSourceBatch decodes one magic-item batch strictly and validates it
// against its reviewed manifest.
func decodeItemSourceBatch[D any](data []byte, m itemSourceManifest, check func(family, variant string, d *D) error) (ItemSourceBatch[D], error) {
	b, err := strictItemSourceBatch[D](data)
	if err != nil {
		return ItemSourceBatch[D]{}, err
	}
	if err := validateItemSourceBatch(b, m, check); err != nil {
		return ItemSourceBatch[D]{}, err
	}
	return b, nil
}

// strictItemSourceBatch decodes without source validation. Duplicate or
// mixed-case members, unknown fields, null, trailing data and any omitted field
// fail: the typed value must re-encode to exactly the supplied JSON tree.
func strictItemSourceBatch[D any](data []byte) (ItemSourceBatch[D], error) {
	var b ItemSourceBatch[D]
	if err := uniqueCatalogueJSON(json.NewDecoder(bytes.NewReader(data)), true); err != nil {
		return b, err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&b); err != nil {
		return ItemSourceBatch[D]{}, fmt.Errorf("decode item source batch: %w", err)
	}
	canonical, err := json.Marshal(b)
	if err != nil {
		return ItemSourceBatch[D]{}, err
	}
	var got, want any
	if err := json.Unmarshal(data, &got); err != nil {
		return ItemSourceBatch[D]{}, err
	}
	if err := json.Unmarshal(canonical, &want); err != nil {
		return ItemSourceBatch[D]{}, err
	}
	if !reflect.DeepEqual(got, want) {
		return ItemSourceBatch[D]{}, fmt.Errorf("item source batch omits a required field")
	}
	return b, nil
}

// validateItemSourceBatch checks a batch against its reviewed manifest, then
// each details payload structurally and with the batch's own source checks.
func validateItemSourceBatch[D any](b ItemSourceBatch[D], m itemSourceManifest, check func(family, variant string, d *D) error) error {
	if b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		return fmt.Errorf("unsupported item contract version")
	}
	if err := validateSource(b.Source); err != nil {
		return err
	}
	if b.Batch != m.Batch || b.CatalogueIssue != m.Issue || len(b.Families) != len(m.Families) {
		return fmt.Errorf("item batch %q does not match its reviewed manifest", b.Batch)
	}
	for i, f := range b.Families {
		want := m.Families[i]
		id := "srd-5.2.1:content/items/" + want.Slug
		if f.ID != id || f.Name != want.Name || f.Kind != "magic" || f.Category != want.Category || !slices.Contains(strings.Split(itemCategories, "|"), f.Category) {
			return fmt.Errorf("unreviewed, missing or reordered item family %q", f.ID)
		}
		if f.Base == nil || !sameItemKeys(f.Base.Options, want.Options) || !sameItemKeys(f.Base.Excluded, want.Excluded) {
			return fmt.Errorf("%s: base item restriction disagrees with source", f.Name)
		}
		if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.Slug || !slices.Equal(f.Presentations[0].Pages, []int{want.Page, want.Page}) {
			return fmt.Errorf("%s: invalid, missing or duplicate source presentation", f.Name)
		}
		if len(f.Variants) != len(want.Variants) {
			return fmt.Errorf("%s: missing or unreviewed source-listed variants", f.Name)
		}
		for j, v := range f.Variants {
			wv := want.Variants[j]
			if v.ID != id+"/variant/"+wv.Slug || v.Name != wv.Name || v.Rarity != wv.Rarity || !slices.Contains(strings.Split(itemRarities, "|"), v.Rarity) {
				return fmt.Errorf("%s: invalid, duplicate or unreviewed variant %q", f.Name, v.ID)
			}
			if v.Attunement == nil || v.Attunement.Mode != want.Attunement || !slices.Contains(strings.Split(itemAttunements, "|"), v.Attunement.Mode) || !sameItemKeys(v.Attunement.Restrictions, "") {
				return fmt.Errorf("%s: attunement disagrees with source", v.Name)
			}
			wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
			if want.Charges != nil {
				wantCharges = *want.Charges
			}
			if v.Charges == nil || *v.Charges != wantCharges {
				return fmt.Errorf("%s: unreviewed charges/recharge", v.Name)
			}
			if v.Equip == nil || v.Equip.Mode != want.Equip || !sameItemKeys(v.Equip.Restrictions, "") {
				return fmt.Errorf("%s: equip restriction disagrees with source", v.Name)
			}
			if v.Use == nil || v.Use.Mode != want.Use || !sameItemKeys(v.Use.Restrictions, want.UseRestrictions) {
				return fmt.Errorf("%s: use restriction disagrees with source", v.Name)
			}
			// No batch family lists a price. Unspecified never means free.
			if v.Cost == nil || v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.Page, want.Page}) {
				return fmt.Errorf("%s: unlisted cost must remain unspecified", v.Name)
			}
			if !sameItemKeys(v.RequiredMechanicKeys, want.Keys) || len(v.RequiredMechanicKeys) == 0 {
				return fmt.Errorf("%s: missing, unknown or reordered mechanic key", v.Name)
			}
			for _, key := range v.RequiredMechanicKeys {
				if !itemMechanicKey.MatchString(key) {
					return fmt.Errorf("%s: invalid mechanic key %q", v.Name, key)
				}
			}
			if v.Details == nil {
				return fmt.Errorf("%s: source details required", v.Name)
			}
			payloads, err := itemDetailPayloads(reflect.ValueOf(v.Details).Elem())
			if err != nil {
				return fmt.Errorf("%s: %w", v.Name, err)
			}
			if !slices.Equal(payloads, []string{want.Payload}) {
				return fmt.Errorf("%s: details must hold exactly the %q payload", v.Name, want.Payload)
			}
			if err := check(want.Slug, wv.Slug, v.Details); err != nil {
				return fmt.Errorf("%s: %w", v.Name, err)
			}
		}
	}
	return nil
}

// sameItemKeys requires an explicit list equal to the space-separated keys.
func sameItemKeys(got []string, want string) bool {
	return got != nil && slices.Equal(got, strings.Fields(want))
}

// itemDetailPayloads names the non-nil payloads of a batch details struct and
// validates each one's tree. Every details field must be a payload pointer.
func itemDetailPayloads(details reflect.Value) ([]string, error) {
	names := []string{}
	for i := 0; i < details.NumField(); i++ {
		field := details.Field(i)
		if field.Kind() != reflect.Pointer {
			return nil, fmt.Errorf("item details fields must be payload pointers")
		}
		if field.IsNil() {
			continue
		}
		name, _, _ := strings.Cut(details.Type().Field(i).Tag.Get("json"), ",")
		if err := validateItemDetailTree(field, false); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		names = append(names, name)
	}
	return names, nil
}

// validateItemDetailTree fails closed on values a payload cannot express as
// source data: nil members, blank or free-form strings, negative numbers and
// empty lists (unless the field is tagged srd:"may-be-empty").
func validateItemDetailTree(v reflect.Value, mayBeEmpty bool) error {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return fmt.Errorf("missing required value")
		}
		return validateItemDetailTree(v.Elem(), false)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if err := validateItemDetailTree(v.Field(i), field.Tag.Get("srd") == "may-be-empty"); err != nil {
				return fmt.Errorf("%s: %w", field.Tag.Get("json"), err)
			}
		}
	case reflect.Slice:
		if v.IsNil() || v.Len() == 0 && !mayBeEmpty {
			return fmt.Errorf("empty required list")
		}
		for i := 0; i < v.Len(); i++ {
			if err := validateItemDetailTree(v.Index(i), false); err != nil {
				return err
			}
		}
	case reflect.String:
		if !itemToken.MatchString(v.String()) {
			return fmt.Errorf("invalid source key %q", v.String())
		}
	case reflect.Int:
		if v.Int() < 0 {
			return fmt.Errorf("negative source value")
		}
	case reflect.Bool:
	default:
		return fmt.Errorf("unsupported details type %s", v.Kind())
	}
	return nil
}

// itemRollTable reports whether rows cover 1..sides exactly once, in order.
func itemRollTable(rows [][2]int, sides int) bool {
	next := 1
	for _, row := range rows {
		if row[0] != next || row[1] < row[0] {
			return false
		}
		next = row[1] + 1
	}
	return sides > 0 && next == sides+1
}

// uniqueItemKeys reports whether keys are non-empty and distinct.
func uniqueItemKeys(keys []string) bool {
	seen := map[string]bool{}
	for _, key := range keys {
		if key == "" || seen[key] {
			return false
		}
		seen[key] = true
	}
	return len(keys) > 0
}
