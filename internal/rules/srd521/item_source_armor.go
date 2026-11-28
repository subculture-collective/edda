package srd521

import (
	"fmt"
	"slices"
)

// The Armor table (SRD 5.2.1 p. 92) payload shared by the armor catalogue
// batches. A record describes the table row and its heading's don and doff
// times. It is not read by the equipment or character projection code, which
// keeps its own fixed Chain Mail and Shield values.

// ItemArmor is one row of the Armor table.
type ItemArmor struct {
	ArmorCategory       string           `json:"armor_category"`
	ArmorClass          ItemArmorClass   `json:"armor_class"`
	StrengthRequirement int              `json:"strength_requirement"`
	StealthDisadvantage bool             `json:"stealth_disadvantage"`
	DonDoff             ItemArmorDonDoff `json:"don_doff"`
	WeightPounds        int              `json:"weight_pounds"`
}

// ItemArmorClass is either a base Armor Class, with the printed Dexterity
// modifier rule, or the bonus the table prints for a Shield.
type ItemArmorClass struct {
	Base             int    `json:"base"`
	Bonus            int    `json:"bonus"`
	Dexterity        string `json:"dexterity"`
	DexterityMaximum int    `json:"dexterity_maximum"`
}

// ItemArmorDonDoff is the time printed in the row's category heading. A Shield
// takes the Utilize action instead of a number of minutes.
type ItemArmorDonDoff struct {
	Action      string `json:"action"`
	DonMinutes  int    `json:"don_minutes"`
	DoffMinutes int    `json:"doff_minutes"`
}

var armorSourceKind = mundaneSourceKind{Category: "armor", Section: "armor", Use: "none", Payload: "armor"}

// decodeArmorSourceBatch decodes one armor batch strictly against its manifest.
func decodeArmorSourceBatch[D any](data []byte, m mundaneSourceManifest, armor func(*D) *ItemArmor) (ItemSourceBatch[D], error) {
	return decodeMundaneSourceBatch(data, m, armorSourceKind, func(d *D, _ int) (mundaneSourceRow, error) {
		p := armor(d)
		if err := checkItemArmor(p); err != nil {
			return mundaneSourceRow{}, err
		}
		equip := "wear"
		if p.ArmorCategory == "shield" {
			equip = "wield"
		}
		return mundaneSourceRow{Equip: equip, Keys: itemArmorKeys(p)}, nil
	})
}

// checkItemArmor checks that a row is internally consistent.
func checkItemArmor(p *ItemArmor) error {
	if !slices.Contains([]string{"light", "medium", "heavy", "shield"}, p.ArmorCategory) {
		return fmt.Errorf("unknown armor category")
	}
	shield := p.ArmorCategory == "shield"
	class := p.ArmorClass
	if shield != (class.Bonus > 0) || shield == (class.Base > 0) {
		return fmt.Errorf("a suit of armor has a base Armor Class and a Shield has a bonus")
	}
	if !slices.Contains([]string{"full", "capped", "none"}, class.Dexterity) || (class.Dexterity == "capped") != (class.DexterityMaximum > 0) || shield && class.Dexterity != "none" {
		return fmt.Errorf("invalid Dexterity modifier rule")
	}
	timed := p.DonDoff.Action == "none" && p.DonDoff.DonMinutes > 0 && p.DonDoff.DoffMinutes > 0
	utilize := p.DonDoff.Action == "utilize" && p.DonDoff.DonMinutes == 0 && p.DonDoff.DoffMinutes == 0
	if shield && !utilize || !shield && !timed {
		return fmt.Errorf("invalid don and doff time")
	}
	if p.WeightPounds < 1 {
		return fmt.Errorf("invalid armor weight")
	}
	return nil
}

// itemArmorKeys restates a row as the mechanic keys a resolver would need.
func itemArmorKeys(p *ItemArmor) []string {
	equip := "item.wear"
	class := fmt.Sprintf("armor-class.%d", p.ArmorClass.Base)
	switch p.ArmorClass.Dexterity {
	case "full":
		class += "-plus-dexterity"
	case "capped":
		class += fmt.Sprintf("-plus-dexterity-max-%d", p.ArmorClass.DexterityMaximum)
	}
	don, doff := fmt.Sprintf("don.minutes-%d", p.DonDoff.DonMinutes), fmt.Sprintf("doff.minutes-%d", p.DonDoff.DoffMinutes)
	if p.ArmorCategory == "shield" {
		equip = "item.wield"
		class = fmt.Sprintf("armor-class.bonus-%d", p.ArmorClass.Bonus)
		don, doff = "don.utilize-action", "doff.utilize-action"
	}
	keys := []string{equip, "armor." + p.ArmorCategory, "armor.training", class}
	if p.StrengthRequirement > 0 {
		keys = append(keys, fmt.Sprintf("strength.%d", p.StrengthRequirement))
	}
	if p.StealthDisadvantage {
		keys = append(keys, "stealth.disadvantage")
	}
	return append(keys, don, doff)
}
