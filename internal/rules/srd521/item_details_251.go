package srd521

import _ "embed"

//go:embed testdata/items_251.json
var items251JSON []byte

// ItemDetails251 holds the armor table row of each record in catalogue batch
// #251 (armor 01, SRD 5.2.1 p. 92). Payloads are source facts only.
type ItemDetails251 struct {
	Armor *ItemArmor `json:"armor,omitempty"`
}

// itemManifest251 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest251 = mundaneSourceManifest{Batch: "data-armor-01", Issue: 251, Families: []mundaneSourceManifestFamily{
	{Slug: "breastplate", Name: "Breastplate", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Breastplate", CostAmount: 400, Currency: "gp",
			Keys: "item.wear armor.medium armor.training armor-class.14-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
	}},
	{Slug: "chain-mail", Name: "Chain Mail", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Chain Mail", CostAmount: 75, Currency: "gp",
			Keys: "item.wear armor.heavy armor.training armor-class.16 strength.13 stealth.disadvantage don.minutes-10 doff.minutes-5"},
	}},
	{Slug: "chain-shirt", Name: "Chain Shirt", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Chain Shirt", CostAmount: 50, Currency: "gp",
			Keys: "item.wear armor.medium armor.training armor-class.13-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
	}},
	{Slug: "half-plate-armor", Name: "Half Plate Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Half Plate Armor", CostAmount: 750, Currency: "gp",
			Keys: "item.wear armor.medium armor.training armor-class.15-plus-dexterity-max-2 stealth.disadvantage don.minutes-5 doff.minutes-1"},
	}},
	{Slug: "hide-armor", Name: "Hide Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Hide Armor", CostAmount: 10, Currency: "gp",
			Keys: "item.wear armor.medium armor.training armor-class.12-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
	}},
	{Slug: "leather-armor", Name: "Leather Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Leather Armor", CostAmount: 10, Currency: "gp",
			Keys: "item.wear armor.light armor.training armor-class.11-plus-dexterity don.minutes-1 doff.minutes-1"},
	}},
	{Slug: "padded-armor", Name: "Padded Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Padded Armor", CostAmount: 5, Currency: "gp",
			Keys: "item.wear armor.light armor.training armor-class.11-plus-dexterity stealth.disadvantage don.minutes-1 doff.minutes-1"},
	}},
	{Slug: "plate-armor", Name: "Plate Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Plate Armor", CostAmount: 1500, Currency: "gp",
			Keys: "item.wear armor.heavy armor.training armor-class.18 strength.15 stealth.disadvantage don.minutes-10 doff.minutes-5"},
	}},
	{Slug: "ring-mail", Name: "Ring Mail", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ring Mail", CostAmount: 30, Currency: "gp",
			Keys: "item.wear armor.heavy armor.training armor-class.14 stealth.disadvantage don.minutes-10 doff.minutes-5"},
	}},
	{Slug: "scale-mail", Name: "Scale Mail", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Scale Mail", CostAmount: 50, Currency: "gp",
			Keys: "item.wear armor.medium armor.training armor-class.14-plus-dexterity-max-2 stealth.disadvantage don.minutes-5 doff.minutes-1"},
	}},
	{Slug: "shield", Name: "Shield", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Shield", CostAmount: 10, Currency: "gp",
			Keys: "item.wield armor.shield armor.training armor-class.bonus-2 don.utilize-action doff.utilize-action"},
	}},
	{Slug: "splint-armor", Name: "Splint Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Splint Armor", CostAmount: 200, Currency: "gp",
			Keys: "item.wear armor.heavy armor.training armor-class.17 strength.15 stealth.disadvantage don.minutes-10 doff.minutes-5"},
	}},
}}

// LoadItems251 returns a fresh validated copy of batch #251. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems251() (ItemSourceBatch[ItemDetails251], error) { return DecodeItems251(items251JSON) }

// DecodeItems251 validates supplied JSON as batch #251 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems251(data []byte) (ItemSourceBatch[ItemDetails251], error) {
	return decodeArmorSourceBatch(data, itemManifest251, func(d *ItemDetails251) *ItemArmor { return d.Armor })
}
