package srd521

import _ "embed"

//go:embed testdata/items_253.json
var items253JSON []byte

// ItemDetails253 holds the tool entry of each record in catalogue batch
// #253 (tools 01, SRD 5.2.1 pp. 93–94). Payloads are source facts only.
type ItemDetails253 struct {
	Tool *ItemTool `json:"tool,omitempty"`
}

// itemManifest253 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest253 = mundaneSourceManifest{Batch: "data-tools-01", Issue: 253, Families: []mundaneSourceManifestFamily{
	{Slug: "alchemists-supplies", Name: "Alchemist’s Supplies", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Alchemist’s Supplies", CostAmount: 50, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.identify-substance-dc-15 utilize.start-fire-dc-15 craft.acid craft.alchemists-fire craft.component-pouch craft.oil craft.paper craft.perfume"},
	}},
	{Slug: "brewers-supplies", Name: "Brewer’s Supplies", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Brewer’s Supplies", CostAmount: 20, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.detect-poisoned-drink-dc-15 utilize.identify-alcohol-dc-10 craft.antitoxin"},
	}},
	{Slug: "calligraphers-supplies", Name: "Calligrapher’s Supplies", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Calligrapher’s Supplies", CostAmount: 10, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.write-text-with-impressive-flourishes-that-guard-against-forgery-dc-15 craft.ink craft.spell-scroll"},
	}},
	{Slug: "carpenters-tools", Name: "Carpenter’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Carpenter’s Tools", CostAmount: 8, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.seal-or-pry-open-door-or-container-dc-20 craft.club craft.greatclub craft.quarterstaff craft.barrel craft.chest craft.ladder craft.pole craft.portable-ram craft.torch"},
	}},
	{Slug: "cartographers-tools", Name: "Cartographer’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Cartographer’s Tools", CostAmount: 15, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.draft-map-of-small-area-dc-15 craft.map"},
	}},
	{Slug: "cobblers-tools", Name: "Cobbler’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Cobbler’s Tools", CostAmount: 5, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.modify-footwear-to-give-advantage-on-wearers-next-dexterity-acrobatics-check-dc-10 craft.climbers-kit"},
	}},
	{Slug: "cooks-utensils", Name: "Cook’s Utensils", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Cook’s Utensils", CostAmount: 1, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.improve-foods-flavor-dc-10 utilize.detect-spoiled-or-poisoned-food-dc-15 craft.rations"},
	}},
	{Slug: "disguise-kit", Name: "Disguise Kit", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Disguise Kit", CostAmount: 25, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.apply-makeup-dc-10 craft.costume"},
	}},
	{Slug: "forgery-kit", Name: "Forgery Kit", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Forgery Kit", CostAmount: 15, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.dexterity utilize.mimic-10-or-fewer-words-of-someone-elses-handwriting-dc-15 utilize.duplicate-wax-seal-dc-20"},
	}},
	{Slug: "gaming-set", Name: "Gaming Set", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "dice", Name: "Dice", CostAmount: 1, Currency: "sp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{Slug: "dragonchess", Name: "Dragonchess", CostAmount: 1, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{Slug: "playing-cards", Name: "Playing cards", CostAmount: 5, Currency: "sp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{Slug: "three-dragon-ante", Name: "Three-dragon ante", CostAmount: 1, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
	}},
	{Slug: "glassblowers-tools", Name: "Glassblower’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Glassblower’s Tools", CostAmount: 30, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-what-glass-object-held-in-past-24-hours-dc-15 craft.glass-bottle craft.magnifying-glass craft.spyglass craft.vial"},
	}},
	{Slug: "herbalism-kit", Name: "Herbalism Kit", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Herbalism Kit", CostAmount: 5, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.intelligence utilize.identify-plant-dc-10 craft.antitoxin craft.candle craft.healers-kit craft.potion-of-healing"},
	}},
}}

// LoadItems253 returns a fresh validated copy of batch #253. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems253() (ItemSourceBatch[ItemDetails253], error) { return DecodeItems253(items253JSON) }

// DecodeItems253 validates supplied JSON as batch #253 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems253(data []byte) (ItemSourceBatch[ItemDetails253], error) {
	return decodeToolSourceBatch(data, itemManifest253, func(d *ItemDetails253) *ItemTool { return d.Tool })
}
