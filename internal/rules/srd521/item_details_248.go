package srd521

import _ "embed"

//go:embed testdata/items_248.json
var items248JSON []byte

// ItemDetails248 holds the weapon table row of each record in catalogue batch
// #248 (weapons 02, SRD 5.2.1 p. 91). Payloads are source facts only.
type ItemDetails248 struct {
	Weapon *ItemWeapon `json:"weapon,omitempty"`
}

// itemManifest248 is the reviewed manifest for the batch's weapons, in planning
// inventory order.
var itemManifest248 = weaponSourceManifest{Batch: "data-weapons-02", Issue: 248, Weapons: []weaponSourceManifestEntry{
	{Slug: "handaxe", Name: "Handaxe", CostAmount: 5, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d6-slashing property.light property.thrown range.20-60 mastery.vex"},
	{Slug: "javelin", Name: "Javelin", CostAmount: 5, Currency: "sp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d6-piercing property.thrown range.30-120 mastery.slow"},
	{Slug: "lance", Name: "Lance", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d10-piercing property.heavy property.reach property.two-handed property.two-handed-unless-mounted mastery.topple"},
	{Slug: "light-crossbow", Name: "Light Crossbow", CostAmount: 25, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.ranged damage.1d8-piercing property.ammunition property.loading property.two-handed range.80-320 ammunition.bolt mastery.slow"},
	{Slug: "light-hammer", Name: "Light Hammer", CostAmount: 2, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d4-bludgeoning property.light property.thrown range.20-60 mastery.nick"},
	{Slug: "longbow", Name: "Longbow", CostAmount: 50, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1d8-piercing property.ammunition property.heavy property.two-handed range.150-600 ammunition.arrow mastery.slow"},
	{Slug: "longsword", Name: "Longsword", CostAmount: 15, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-slashing property.versatile damage.two-handed-1d10 mastery.sap"},
	{Slug: "mace", Name: "Mace", CostAmount: 5, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d6-bludgeoning mastery.sap"},
	{Slug: "maul", Name: "Maul", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.2d6-bludgeoning property.heavy property.two-handed mastery.topple"},
	{Slug: "morningstar", Name: "Morningstar", CostAmount: 15, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-piercing mastery.sap"},
	{Slug: "musket", Name: "Musket", CostAmount: 500, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1d12-piercing property.ammunition property.loading property.two-handed range.40-120 ammunition.bullet mastery.slow"},
	{Slug: "pike", Name: "Pike", CostAmount: 5, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d10-piercing property.heavy property.reach property.two-handed mastery.push"},
}}

// LoadItems248 returns a fresh validated copy of batch #248. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems248() (ItemSourceBatch[ItemDetails248], error) { return DecodeItems248(items248JSON) }

// DecodeItems248 validates supplied JSON as batch #248 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems248(data []byte) (ItemSourceBatch[ItemDetails248], error) {
	return decodeWeaponSourceBatch(data, itemManifest248, func(d *ItemDetails248) *ItemWeapon { return d.Weapon })
}
