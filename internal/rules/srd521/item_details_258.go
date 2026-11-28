package srd521

import _ "embed"

//go:embed testdata/items_258.json
var items258JSON []byte

// ItemDetails258 holds the gear entry of each record in catalogue batch
// #258 (adventuring gear 03, SRD 5.2.1 p. 97). Payloads are source
// facts only.
type ItemDetails258 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest258 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest258 = mundaneSourceManifest{Batch: "data-adventuring-gear-03", Issue: 258, Families: []mundaneSourceManifestFamily{
	{Slug: "clothes-fine", Name: "Clothes, Fine", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Clothes, Fine", CostAmount: 15, Currency: "gp",
			Keys: "item.gear trait.expensive-fabrics-with-expertly-crafted-details trait.required-for-admission-to-some-events-and-locations"},
	}},
	{Slug: "clothes-travelers", Name: "Clothes, Traveler’s", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Clothes, Traveler’s", CostAmount: 2, Currency: "gp",
			Keys: "item.gear trait.resilient-garments-for-travel-in-various-environments"},
	}},
	{Slug: "component-pouch", Name: "Component Pouch", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Component Pouch", CostAmount: 25, Currency: "gp",
			Keys: "item.gear trait.watertight trait.compartments-hold-all-free-material-components-of-your-spells"},
	}},
	{Slug: "costume", Name: "Costume", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Costume", CostAmount: 5, Currency: "gp",
			Keys: "item.gear trait.advantage-on-checks-to-impersonate-represented-person-or-type-while-wearing"},
	}},
	{Slug: "crowbar", Name: "Crowbar", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Crowbar", CostAmount: 2, Currency: "gp",
			Keys: "item.gear trait.advantage-on-strength-checks-where-leverage-applies"},
	}},
	{Slug: "diplomats-pack", Name: "Diplomat’s Pack", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Diplomat’s Pack", CostAmount: 39, Currency: "gp",
			Keys: "item.gear content.chest-1-item content.clothes-fine-1-item content.ink-1-item content.ink-pen-5-item content.lamp-1-item content.case-map-or-scroll-2-item content.oil-4-flask content.paper-5-sheet content.parchment-5-sheet content.perfume-1-item content.tinderbox-1-item"},
	}},
	{Slug: "druidic-focus", Name: "Druidic Focus", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "sprig-of-mistletoe", Name: "Sprig of mistletoe", CostAmount: 1, Currency: "gp",
			Keys: "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic"},
		{Slug: "wooden-staff", Name: "Wooden staff", CostAmount: 5, Currency: "gp",
			Keys: "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic trait.also-quarterstaff"},
		{Slug: "yew-wand", Name: "Yew wand", CostAmount: 10, Currency: "gp",
			Keys: "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic"},
	}},
	{Slug: "dungeoneers-pack", Name: "Dungeoneer’s Pack", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Dungeoneer’s Pack", CostAmount: 12, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.caltrops-1-item content.crowbar-1-item content.oil-2-flask content.rations-10-day content.rope-1-item content.tinderbox-1-item content.torch-10-item content.waterskin-1-item"},
	}},
	{Slug: "entertainers-pack", Name: "Entertainer’s Pack", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Entertainer’s Pack", CostAmount: 40, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.bedroll-1-item content.bell-1-item content.lantern-bullseye-1-item content.costume-3-item content.mirror-1-item content.oil-8-flask content.rations-9-day content.tinderbox-1-item content.waterskin-1-item"},
	}},
	{Slug: "explorers-pack", Name: "Explorer’s Pack", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Explorer’s Pack", CostAmount: 10, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.bedroll-1-item content.oil-2-flask content.rations-10-day content.rope-1-item content.tinderbox-1-item content.torch-10-item content.waterskin-1-item"},
	}},
	{Slug: "flask", Name: "Flask", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Flask", CostAmount: 2, Currency: "cp",
			Keys: "item.gear measure.capacity-1-pint"},
	}},
	{Slug: "grappling-hook", Name: "Grappling Hook", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Grappling Hook", CostAmount: 2, Currency: "gp",
			Keys: "item.gear action.utilize ability-check.catch-hook-dexterity-acrobatics-dc-13 measure.throw-range-50-feet trait.throw-at-railing-ledge-or-other-catch trait.tied-rope-can-then-be-climbed"},
	}},
}}

// LoadItems258 returns a fresh validated copy of batch #258. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems258() (ItemSourceBatch[ItemDetails258], error) { return DecodeItems258(items258JSON) }

// DecodeItems258 validates supplied JSON as batch #258 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems258(data []byte) (ItemSourceBatch[ItemDetails258], error) {
	return decodeGearSourceBatch(data, itemManifest258, itemGear258)
}

func itemGear258(d *ItemDetails258) *ItemGear { return d.Gear }
