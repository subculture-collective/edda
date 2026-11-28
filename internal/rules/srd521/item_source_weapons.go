package srd521

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// The Weapons table (SRD 5.2.1 p. 91) lists mundane items: they have no rarity
// or attunement, and each has a listed price. This file holds the payload and
// validation shared by the weapon catalogue batches. A weapon record describes
// the table row. It is not read by any attack resolver, and the fixed weapon
// fixtures in creation and combat code do not read it.

const (
	weaponSourcePage = 91
	weaponProperties = "ammunition|finesse|heavy|light|loading|reach|thrown|two-handed|versatile"
	weaponMasteries  = "cleave|graze|nick|push|sap|slow|topple|vex"
)

// ItemWeapon is one row of the Weapons table. Range is zero without the
// Ammunition or Thrown property, VersatileDamage is zero without Versatile, and
// Ammunition is "none" for a weapon that fires none.
type ItemWeapon struct {
	ProficiencyCategory string                `json:"proficiency_category"`
	AttackKind          string                `json:"attack_kind"`
	Damage              ItemWeaponDamage      `json:"damage"`
	Properties          []string              `json:"properties" srd:"may-be-empty"`
	Range               ItemWeaponRange       `json:"range"`
	Ammunition          string                `json:"ammunition"`
	VersatileDamage     ItemWeaponVersatility `json:"versatile_damage"`
	TwoHandedException  string                `json:"two_handed_exception"`
	Mastery             string                `json:"mastery"`
	Weight              ItemWeaponWeight      `json:"weight"`
}

// ItemWeaponDamage is either dice, or the flat amount the table prints for the
// Blowgun. Exactly one of the two forms is set.
type ItemWeaponDamage struct {
	Dice  int    `json:"dice"`
	Sides int    `json:"sides"`
	Flat  int    `json:"flat"`
	Type  string `json:"type"`
}

// ItemWeaponRange records the printed normal and long range in feet.
type ItemWeaponRange struct {
	NormalFeet int `json:"normal_feet"`
	LongFeet   int `json:"long_feet"`
}

// ItemWeaponVersatility records the two-handed damage dice of a Versatile weapon.
type ItemWeaponVersatility struct {
	Dice  int `json:"dice"`
	Sides int `json:"sides"`
}

// ItemWeaponWeight is the printed weight as a fraction of a pound. Listed is
// false where the table prints a dash, which is not a weight of zero.
type ItemWeaponWeight struct {
	Listed            bool `json:"listed"`
	PoundsNumerator   int  `json:"pounds_numerator"`
	PoundsDenominator int  `json:"pounds_denominator"`
}

// weaponSourceManifest is a weapon batch's reviewed identity, price and
// mechanic key manifest, independent of incoming JSON.
type weaponSourceManifest struct {
	Batch   string
	Issue   int
	Weapons []weaponSourceManifestEntry
}

type weaponSourceManifestEntry struct {
	Slug, Name string
	CostAmount int
	Currency   string
	Keys       string
}

// decodeWeaponSourceBatch applies the shared strict decoding, then the weapon
// manifest. The magic-item validator does not accept these records.
func decodeWeaponSourceBatch[D any](data []byte, m weaponSourceManifest, weapon func(*D) *ItemWeapon) (ItemSourceBatch[D], error) {
	b, err := strictItemSourceBatch[D](data)
	if err != nil {
		return ItemSourceBatch[D]{}, err
	}
	if err := validateWeaponSourceBatch(b, m, weapon); err != nil {
		return ItemSourceBatch[D]{}, err
	}
	return b, nil
}

func validateWeaponSourceBatch[D any](b ItemSourceBatch[D], m weaponSourceManifest, weapon func(*D) *ItemWeapon) error {
	if b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		return fmt.Errorf("unsupported item contract version")
	}
	if err := validateSource(b.Source); err != nil {
		return err
	}
	if b.Batch != m.Batch || b.CatalogueIssue != m.Issue || len(b.Families) != len(m.Weapons) {
		return fmt.Errorf("weapon batch %q does not match its reviewed manifest", b.Batch)
	}
	pages := []int{weaponSourcePage, weaponSourcePage}
	for i, f := range b.Families {
		want := m.Weapons[i]
		id := "srd-5.2.1:content/items/" + want.Slug
		if f.ID != id || f.Name != want.Name || f.Kind != "mundane" || f.Category != "weapon" {
			return fmt.Errorf("unreviewed, missing or reordered weapon %q", f.ID)
		}
		if f.Base == nil || !sameItemKeys(f.Base.Options, "") || !sameItemKeys(f.Base.Excluded, "") {
			return fmt.Errorf("%s: a mundane weapon has no base item restriction", f.Name)
		}
		if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/weapons/"+want.Slug || !slices.Equal(f.Presentations[0].Pages, pages) {
			return fmt.Errorf("%s: invalid, missing or duplicate source presentation", f.Name)
		}
		if len(f.Variants) != 1 {
			return fmt.Errorf("%s: a table weapon has exactly one form", f.Name)
		}
		v := f.Variants[0]
		if v.ID != id+"/variant/standard" || v.Name != want.Name || v.Rarity != "none" {
			return fmt.Errorf("%s: invalid or unreviewed variant %q", f.Name, v.ID)
		}
		if v.Attunement == nil || v.Attunement.Mode != "none" || !sameItemKeys(v.Attunement.Restrictions, "") {
			return fmt.Errorf("%s: a mundane weapon has no attunement", v.Name)
		}
		if v.Charges == nil || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
			return fmt.Errorf("%s: a table weapon has no charges", v.Name)
		}
		if v.Equip == nil || v.Equip.Mode != "wield" || !sameItemKeys(v.Equip.Restrictions, "") {
			return fmt.Errorf("%s: equip restriction disagrees with source", v.Name)
		}
		if v.Use == nil || v.Use.Mode != "weapon-attack" || !sameItemKeys(v.Use.Restrictions, "") {
			return fmt.Errorf("%s: use restriction disagrees with source", v.Name)
		}
		if v.Cost == nil || v.Cost.Kind != "listed" || v.Cost.Amount != want.CostAmount || v.Cost.Currency != want.Currency || !slices.Equal(v.Cost.SourcePages, pages) {
			return fmt.Errorf("%s: listed price disagrees with source", v.Name)
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
		if !slices.Equal(payloads, []string{"weapon"}) {
			return fmt.Errorf("%s: details must hold exactly the %q payload", v.Name, "weapon")
		}
		p := weapon(v.Details)
		if err := checkItemWeapon(p); err != nil {
			return fmt.Errorf("%s: %w", v.Name, err)
		}
		// The keys restate the row's values, so the two cannot disagree.
		if got := strings.Join(itemWeaponKeys(p), " "); got != want.Keys {
			return fmt.Errorf("%s: row values disagree with the mechanic keys: %s", v.Name, got)
		}
	}
	return nil
}

// checkItemWeapon checks that a row is internally consistent.
func checkItemWeapon(p *ItemWeapon) error {
	oneOf := func(value, allowed string) bool { return slices.Contains(strings.Split(allowed, "|"), value) }
	has := func(property string) bool { return slices.Contains(p.Properties, property) }
	if !oneOf(p.ProficiencyCategory, "simple|martial") || !oneOf(p.AttackKind, "melee|ranged") || !oneOf(p.Mastery, weaponMasteries) {
		return fmt.Errorf("unknown weapon category or mastery")
	}
	rolled := p.Damage.Flat == 0 && p.Damage.Dice > 0 && slices.Contains([]int{4, 6, 8, 10, 12}, p.Damage.Sides)
	flat := p.Damage.Flat > 0 && p.Damage.Dice == 0 && p.Damage.Sides == 0
	if !(rolled || flat) || !oneOf(p.Damage.Type, "bludgeoning|piercing|slashing") {
		return fmt.Errorf("invalid weapon damage")
	}
	// The table prints properties in alphabetical order.
	if !itemKeys(p.Properties, weaponProperties) || !slices.IsSorted(p.Properties) {
		return fmt.Errorf("unknown, duplicate or reordered weapon property")
	}
	ranged := p.Range.NormalFeet > 0 && p.Range.LongFeet > p.Range.NormalFeet
	if (has("ammunition") || has("thrown")) != ranged || !ranged && p.Range != (ItemWeaponRange{}) {
		return fmt.Errorf("range disagrees with the weapon's properties")
	}
	if has("ammunition") != oneOf(p.Ammunition, "arrow|bolt|bullet|needle") || !has("ammunition") && p.Ammunition != "none" {
		return fmt.Errorf("ammunition disagrees with the weapon's properties")
	}
	versatile := p.VersatileDamage.Dice > 0 && p.VersatileDamage.Sides > p.Damage.Sides
	if has("versatile") != versatile || !versatile && p.VersatileDamage != (ItemWeaponVersatility{}) {
		return fmt.Errorf("two-handed damage disagrees with the Versatile property")
	}
	if p.TwoHandedException != "none" && (p.TwoHandedException != "unless-mounted" || !has("two-handed")) {
		return fmt.Errorf("invalid Two-Handed exception")
	}
	if p.Weight.Listed != (p.Weight.PoundsNumerator > 0) || p.Weight.PoundsDenominator < 1 {
		return fmt.Errorf("invalid weapon weight")
	}
	return nil
}

// itemWeaponKeys restates a row as the mechanic keys a resolver would need.
func itemWeaponKeys(p *ItemWeapon) []string {
	damage := strconv.Itoa(p.Damage.Flat)
	if p.Damage.Dice > 0 {
		damage = fmt.Sprintf("%dd%d", p.Damage.Dice, p.Damage.Sides)
	}
	keys := []string{"item.wield", "weapon.attack", "weapon." + p.ProficiencyCategory, "weapon." + p.AttackKind, "damage." + damage + "-" + p.Damage.Type}
	for _, property := range p.Properties {
		keys = append(keys, "property."+property)
	}
	if p.Range.NormalFeet > 0 {
		keys = append(keys, fmt.Sprintf("range.%d-%d", p.Range.NormalFeet, p.Range.LongFeet))
	}
	if p.Ammunition != "none" {
		keys = append(keys, "ammunition."+p.Ammunition)
	}
	if p.VersatileDamage.Dice > 0 {
		keys = append(keys, fmt.Sprintf("damage.two-handed-%dd%d", p.VersatileDamage.Dice, p.VersatileDamage.Sides))
	}
	if p.TwoHandedException != "none" {
		keys = append(keys, "property.two-handed-"+p.TwoHandedException)
	}
	return append(keys, "mastery."+p.Mastery)
}
