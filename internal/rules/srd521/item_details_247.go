package srd521

import _ "embed"

//go:embed testdata/items_247.json
var items247JSON []byte

// ItemDetails247 holds the weapon table row of each record in catalogue batch
// #247 (weapons 01, SRD 5.2.1 p. 91). Payloads are source facts only.
type ItemDetails247 struct {
	Weapon *ItemWeapon `json:"weapon,omitempty"`
}

// itemManifest247 is the reviewed manifest for the batch's weapons, in planning
// inventory order.
var itemManifest247 = weaponSourceManifest{Batch: "data-weapons-01", Issue: 247, Weapons: []weaponSourceManifestEntry{
	{Slug: "battleaxe", Name: "Battleaxe", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-slashing property.versatile damage.two-handed-1d10 mastery.topple"},
	{Slug: "blowgun", Name: "Blowgun", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1-piercing property.ammunition property.loading range.25-100 ammunition.needle mastery.vex"},
	{Slug: "club", Name: "Club", CostAmount: 1, Currency: "sp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d4-bludgeoning property.light mastery.slow"},
	{Slug: "dagger", Name: "Dagger", CostAmount: 2, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d4-piercing property.finesse property.light property.thrown range.20-60 mastery.nick"},
	{Slug: "dart", Name: "Dart", CostAmount: 5, Currency: "cp",
		Keys: "item.wield weapon.attack weapon.simple weapon.ranged damage.1d4-piercing property.finesse property.thrown range.20-60 mastery.vex"},
	{Slug: "flail", Name: "Flail", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-bludgeoning mastery.sap"},
	{Slug: "glaive", Name: "Glaive", CostAmount: 20, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d10-slashing property.heavy property.reach property.two-handed mastery.graze"},
	{Slug: "greataxe", Name: "Greataxe", CostAmount: 30, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d12-slashing property.heavy property.two-handed mastery.cleave"},
	{Slug: "greatclub", Name: "Greatclub", CostAmount: 2, Currency: "sp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d8-bludgeoning property.two-handed mastery.push"},
	{Slug: "greatsword", Name: "Greatsword", CostAmount: 50, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.2d6-slashing property.heavy property.two-handed mastery.graze"},
	{Slug: "halberd", Name: "Halberd", CostAmount: 20, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d10-slashing property.heavy property.reach property.two-handed mastery.cleave"},
	{Slug: "hand-crossbow", Name: "Hand Crossbow", CostAmount: 75, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1d6-piercing property.ammunition property.light property.loading range.30-120 ammunition.bolt mastery.vex"},
}}

// LoadItems247 returns a fresh validated copy of batch #247. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems247() (ItemSourceBatch[ItemDetails247], error) { return DecodeItems247(items247JSON) }

// DecodeItems247 validates supplied JSON as batch #247 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems247(data []byte) (ItemSourceBatch[ItemDetails247], error) {
	return decodeWeaponSourceBatch(data, itemManifest247, func(d *ItemDetails247) *ItemWeapon { return d.Weapon })
}
