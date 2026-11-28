package srd521

import _ "embed"

//go:embed testdata/items_261.json
var items261JSON []byte

// ItemDetails261 holds the gear entry of each record in catalogue batch
// #261 (adventuring gear 06, SRD 5.2.1 p. 99). Payloads are source
// facts only.
type ItemDetails261 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest261 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest261 = mundaneSourceManifest{Batch: "data-adventuring-gear-06", Issue: 261, Families: []mundaneSourceManifestFamily{
	{Slug: "pouch", Name: "Pouch", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Pouch", CostAmount: 5, Currency: "sp",
			Keys: "item.gear measure.capacity-weight-6-pounds measure.capacity-volume-1-over-5-cubic-foot"},
	}},
	{Slug: "priests-pack", Name: "Priest’s Pack", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Priest’s Pack", CostAmount: 33, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.blanket-1-item content.holy-water-1-item content.lamp-1-item content.rations-7-day content.robe-1-item content.tinderbox-1-item"},
	}},
	{Slug: "quiver", Name: "Quiver", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Quiver", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.capacity-20-arrows"},
	}},
	{Slug: "ram-portable", Name: "Ram, Portable", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ram, Portable", CostAmount: 4, Currency: "gp",
			Keys: "item.gear measure.strength-check-bonus-4-bonus measure.helpers-1-character trait.used-to-break-down-doors trait.helper-grants-advantage-on-check"},
	}},
	{Slug: "rations", Name: "Rations", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Rations", CostAmount: 5, Currency: "sp",
			Keys: "item.gear trait.travel-ready-food-including-jerky-dried-fruit-hardtack-and-nuts trait.not-eating-risks-malnutrition"},
	}},
	{Slug: "robe", Name: "Robe", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Robe", CostAmount: 1, Currency: "gp",
			Keys: "item.gear trait.vocational-or-ceremonial-significance trait.required-colors-or-symbols-for-admission-to-some-events-and-locations"},
	}},
	{Slug: "rope", Name: "Rope", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Rope", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.utilize action.action ability-check.tie-knot-dexterity-sleight-of-hand-dc-10 ability-check.burst-strength-athletics-dc-20 ability-check.escape-dexterity-acrobatics-dc-15 trait.bind-only-unwilling-grappled-incapacitated-or-restrained-creature trait.bound-legs-restrain-until-escape"},
	}},
	{Slug: "sack", Name: "Sack", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Sack", CostAmount: 1, Currency: "cp",
			Keys: "item.gear measure.capacity-weight-30-pounds measure.capacity-volume-1-cubic-foot"},
	}},
	{Slug: "scholars-pack", Name: "Scholar’s Pack", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Scholar’s Pack", CostAmount: 40, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.book-1-item content.ink-1-item content.ink-pen-1-item content.lamp-1-item content.oil-10-flask content.parchment-10-sheet content.tinderbox-1-item"},
	}},
	{Slug: "shovel", Name: "Shovel", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Shovel", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.digging-time-1-hour measure.hole-side-5-feet trait.digs-hole-in-soil-or-similar-material"},
	}},
	{Slug: "signal-whistle", Name: "Signal Whistle", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Signal Whistle", CostAmount: 5, Currency: "cp",
			Keys: "item.gear action.utilize measure.audible-distance-600-feet trait.blown-to-produce-sound"},
	}},
}}

// LoadItems261 returns a fresh validated copy of batch #261. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems261() (ItemSourceBatch[ItemDetails261], error) { return DecodeItems261(items261JSON) }

// DecodeItems261 validates supplied JSON as batch #261 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems261(data []byte) (ItemSourceBatch[ItemDetails261], error) {
	return decodeGearSourceBatch(data, itemManifest261, itemGear261)
}

func itemGear261(d *ItemDetails261) *ItemGear { return d.Gear }
