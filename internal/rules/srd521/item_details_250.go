package srd521

import _ "embed"

//go:embed testdata/items_250.json
var items250JSON []byte

// ItemDetails250 holds the weapon table row of each record in catalogue batch
// #250 (weapons 04, SRD 5.2.1 p. 91). Payloads are source facts only.
type ItemDetails250 struct {
	Weapon *ItemWeapon `json:"weapon,omitempty"`
}

// itemManifest250 is the reviewed manifest for the batch, in planning
// inventory order. Heavy Crossbow follows: the inventory omits it.
var itemManifest250 = weaponSourceManifest{Batch: "data-weapons-04", Issue: 250, Weapons: []weaponSourceManifestEntry{
	{Slug: "whip", Name: "Whip", CostAmount: 2, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.melee damage.1d4-slashing property.finesse property.reach mastery.slow"},
	{Slug: "heavy-crossbow", Name: "Heavy Crossbow", CostAmount: 50, Currency: "gp",
		Keys: "item.wield weapon.attack weapon.martial weapon.ranged damage.1d10-piercing property.ammunition property.heavy property.loading property.two-handed range.100-400 ammunition.bolt mastery.push"},
}}

// LoadItems250 returns a fresh validated copy of batch #250. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems250() (ItemSourceBatch[ItemDetails250], error) { return DecodeItems250(items250JSON) }

// DecodeItems250 validates supplied JSON as batch #250 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems250(data []byte) (ItemSourceBatch[ItemDetails250], error) {
	return decodeWeaponSourceBatch(data, itemManifest250, func(d *ItemDetails250) *ItemWeapon { return d.Weapon })
}
