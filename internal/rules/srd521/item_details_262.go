package srd521

import _ "embed"

//go:embed testdata/items_262.json
var items262JSON []byte

// ItemDetails262 holds the gear entry of each record in catalogue batch
// #262 (adventuring gear 07, SRD 5.2.1 p. 100). Payloads are source
// facts only.
type ItemDetails262 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest262 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest262 = mundaneSourceManifest{Batch: "data-adventuring-gear-07", Issue: 262, Families: []mundaneSourceManifestFamily{
	{Slug: "spikes-iron", Name: "Spikes, Iron", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Spikes, Iron", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.utilize measure.bundle-amount-10-spikes trait.hammer-with-blunt-object-into-wood-earth-or-similar-material trait.jam-door-shut-or-anchor-rope-or-chain"},
	}},
	{Slug: "spyglass", Name: "Spyglass", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Spyglass", CostAmount: 1000, Currency: "gp",
			Keys: "item.gear measure.magnification-2-times-size"},
	}},
	{Slug: "string", Name: "String", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "String", CostAmount: 1, Currency: "sp",
			Keys: "item.gear action.utilize measure.length-10-feet trait.can-be-knotted"},
	}},
	{Slug: "tent", Name: "Tent", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Tent", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.occupants-2-creatures trait.occupants-are-small-or-medium"},
	}},
	{Slug: "tinderbox", Name: "Tinderbox", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Tinderbox", CostAmount: 5, Currency: "sp",
			Keys: "item.gear action.bonus-action measure.other-fire-lighting-time-1-minute trait.small-container-of-flint-fire-steel-and-tinder trait.lights-candle-lamp-lantern-torch-or-exposed-fuel"},
	}},
	{Slug: "torch", Name: "Torch", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Torch", CostAmount: 1, Currency: "cp",
			Keys: "item.gear action.attack-action damage.hit-1-fire measure.burn-duration-1-hour measure.bright-light-radius-20-feet measure.additional-dim-light-20-feet trait.attack-with-torch-as-simple-melee-weapon"},
	}},
	{Slug: "vial", Name: "Vial", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Vial", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.capacity-4-ounces"},
	}},
	{Slug: "waterskin", Name: "Waterskin", Page: 100, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Waterskin", CostAmount: 2, Currency: "sp",
			Keys: "item.gear measure.capacity-4-pints trait.insufficient-water-risks-dehydration"},
	}},
}}

// LoadItems262 returns a fresh validated copy of batch #262. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems262() (ItemSourceBatch[ItemDetails262], error) { return DecodeItems262(items262JSON) }

// DecodeItems262 validates supplied JSON as batch #262 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems262(data []byte) (ItemSourceBatch[ItemDetails262], error) {
	return decodeGearSourceBatch(data, itemManifest262, itemGear262)
}

func itemGear262(d *ItemDetails262) *ItemGear { return d.Gear }
