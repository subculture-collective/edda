package srd521

import _ "embed"

//go:embed testdata/items_249.json
var items249JSON []byte

// ItemDetails249 holds the weapon table row of each record in catalogue batch
// #249 (weapons 03, SRD 5.2.1 p. 91). Payloads are source facts only.
type ItemDetails249 struct {
	Weapon *ItemWeapon `json:"weapon,omitempty"`
}

// itemManifest249 is the reviewed manifest for the batch's weapons, in planning
// inventory order.
var itemManifest249 = weaponSourceManifest{Batch: "data-weapons-03", Issue: 249, Weapons: []weaponSourceManifestEntry{
	{Slug: "pistol", Name: "Pistol", CostAmount: 250, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1d10-piercing property.ammunition property.loading range.30-90 ammunition.bullet mastery.vex"},
	{Slug: "quarterstaff", Name: "Quarterstaff", CostAmount: 2, Currency: "sp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d6-bludgeoning property.versatile damage.two-handed-1d8 mastery.topple"},
	{Slug: "rapier", Name: "Rapier", CostAmount: 25, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-piercing property.finesse mastery.vex"},
	{Slug: "scimitar", Name: "Scimitar", CostAmount: 25, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d6-slashing property.finesse property.light mastery.nick"},
	{Slug: "shortbow", Name: "Shortbow", CostAmount: 25, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.ranged damage.1d6-piercing property.ammunition property.two-handed range.80-320 ammunition.arrow mastery.vex"},
	{Slug: "shortsword", Name: "Shortsword", CostAmount: 10, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d6-piercing property.finesse property.light mastery.vex"},
	{Slug: "sickle", Name: "Sickle", CostAmount: 1, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d4-slashing property.light mastery.nick"},
	{Slug: "sling", Name: "Sling", CostAmount: 1, Currency: "sp",
		Keys: "item.wield weapon.attack weapon.simple weapon.ranged damage.1d4-bludgeoning property.ammunition range.30-120 ammunition.bullet mastery.slow"},
	{Slug: "spear", Name: "Spear", CostAmount: 1, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.simple weapon.melee damage.1d6-piercing property.thrown property.versatile range.20-60 damage.two-handed-1d8 mastery.sap"},
	{Slug: "trident", Name: "Trident", CostAmount: 5, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-piercing property.thrown property.versatile range.20-60 damage.two-handed-1d10 mastery.topple"},
	{Slug: "war-pick", Name: "War Pick", CostAmount: 5, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-piercing property.versatile damage.two-handed-1d10 mastery.sap"},
	{Slug: "warhammer", Name: "Warhammer", CostAmount: 15, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d8-bludgeoning property.versatile damage.two-handed-1d10 mastery.push"},
}}

// LoadItems249 returns a fresh validated copy of batch #249. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems249() (ItemSourceBatch[ItemDetails249], error) { return DecodeItems249(items249JSON) }

// DecodeItems249 validates supplied JSON as batch #249 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems249(data []byte) (ItemSourceBatch[ItemDetails249], error) {
	return decodeWeaponSourceBatch(data, itemManifest249, func(d *ItemDetails249) *ItemWeapon { return d.Weapon })
}
