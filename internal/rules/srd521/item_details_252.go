package srd521

import _ "embed"

//go:embed testdata/items_252.json
var items252JSON []byte

// ItemDetails252 holds the armor table row of each record in catalogue batch
// #252 (armor 02, SRD 5.2.1 p. 92). Payloads are source facts only.
type ItemDetails252 struct {
	Armor *ItemArmor `json:"armor,omitempty"`
}

// itemManifest252 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest252 = mundaneSourceManifest{Batch: "data-armor-02", Issue: 252, Families: []mundaneSourceManifestFamily{
	{Slug: "studded-leather-armor", Name: "Studded Leather Armor", Page: 92, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Studded Leather Armor", CostAmount: 45, Currency: "gp",
			Keys: "item.wear armor.light armor.training armor-class.12-plus-dexterity don.minutes-1 doff.minutes-1"},
	}},
}}

// LoadItems252 returns a fresh validated copy of batch #252. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems252() (ItemSourceBatch[ItemDetails252], error) { return DecodeItems252(items252JSON) }

// DecodeItems252 validates supplied JSON as batch #252 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems252(data []byte) (ItemSourceBatch[ItemDetails252], error) {
	return decodeArmorSourceBatch(data, itemManifest252, func(d *ItemDetails252) *ItemArmor { return d.Armor })
}
