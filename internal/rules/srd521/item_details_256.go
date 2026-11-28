package srd521

import _ "embed"

//go:embed testdata/items_256.json
var items256JSON []byte

// ItemDetails256 holds the gear entry of each record in catalogue batch
// #256 (adventuring gear 01, SRD 5.2.1 pp. 94-96). Payloads are source
// facts only.
type ItemDetails256 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest256 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest256 = mundaneSourceManifest{Batch: "data-adventuring-gear-01", Issue: 256, Families: []mundaneSourceManifestFamily{
	{Slug: "acid", Name: "Acid", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Acid", CostAmount: 25, Currency: "gp",
			Keys: "item.gear action.attack-replacement saving-throw.avoid-damage-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-vial-2d6-acid measure.throw-range-20-feet trait.throw-vial-in-place-of-one-attack trait.target-one-creature-or-object-you-can-see"},
	}},
	{Slug: "alchemists-fire", Name: "Alchemist’s Fire", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Alchemist’s Fire", CostAmount: 50, Currency: "gp",
			Keys: "item.gear action.attack-replacement saving-throw.avoid-damage-and-burning-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-flask-1d4-fire measure.throw-range-20-feet trait.throw-flask-in-place-of-one-attack trait.target-one-creature-or-object-you-can-see trait.failed-save-starts-burning"},
	}},
	{Slug: "ammunition", Name: "Ammunition", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "arrows", Name: "Arrows", CostAmount: 1, Currency: "gp", CostPage: 96,
			Keys: "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-quiver"},
		{Slug: "bolts", Name: "Bolts", CostAmount: 1, Currency: "gp", CostPage: 96,
			Keys: "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-case"},
		{Slug: "bullets-firearm", Name: "Bullets, Firearm", CostAmount: 3, Currency: "gp", CostPage: 96,
			Keys: "item.gear measure.bundle-amount-10-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
		{Slug: "bullets-sling", Name: "Bullets, Sling", CostAmount: 4, Currency: "cp", CostPage: 96,
			Keys: "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
		{Slug: "needles", Name: "Needles", CostAmount: 1, Currency: "gp", CostPage: 96,
			Keys: "item.gear measure.bundle-amount-50-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
	}},
	{Slug: "antitoxin", Name: "Antitoxin", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Antitoxin", CostAmount: 50, Currency: "gp",
			Keys: "item.gear action.bonus-action measure.benefit-duration-1-hour trait.drink-one-vial trait.advantage-on-saves-to-avoid-or-end-poisoned"},
	}},
	{Slug: "arcane-focus", Name: "Arcane Focus", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "crystal", Name: "Crystal", CostAmount: 10, Currency: "gp",
			Keys: "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{Slug: "orb", Name: "Orb", CostAmount: 20, Currency: "gp",
			Keys: "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{Slug: "rod", Name: "Rod", CostAmount: 10, Currency: "gp",
			Keys: "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{Slug: "staff", Name: "Staff", CostAmount: 5, Currency: "gp",
			Keys: "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic trait.also-quarterstaff"},
		{Slug: "wand", Name: "Wand", CostAmount: 10, Currency: "gp",
			Keys: "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
	}},
	{Slug: "backpack", Name: "Backpack", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Backpack", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.capacity-weight-30-pounds measure.capacity-volume-1-cubic-foot trait.can-serve-as-saddlebag"},
	}},
	{Slug: "ball-bearings", Name: "Ball Bearings", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ball Bearings", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.utilize saving-throw.avoid-prone-dexterity-dc-10 measure.covered-square-side-10-feet measure.placement-range-10-feet measure.recovery-time-10-minutes trait.spill-from-pouch-over-level-area trait.triggers-on-first-entry-on-a-turn trait.failed-save-causes-prone"},
	}},
	{Slug: "barrel", Name: "Barrel", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Barrel", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.liquid-capacity-40-gallons measure.dry-goods-capacity-4-cubic-feet trait.holds-liquid-or-dry-goods"},
	}},
	{Slug: "basket", Name: "Basket", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Basket", CostAmount: 4, Currency: "sp",
			Keys: "item.gear measure.capacity-weight-40-pounds measure.capacity-volume-2-cubic-feet"},
	}},
	{Slug: "bedroll", Name: "Bedroll", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Bedroll", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.occupants-1-creature trait.occupant-is-small-or-medium trait.automatic-success-on-saves-against-extreme-cold-while-inside"},
	}},
	{Slug: "bell", Name: "Bell", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Bell", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.utilize measure.audible-distance-60-feet trait.rung-to-produce-sound"},
	}},
	{Slug: "blanket", Name: "Blanket", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Blanket", CostAmount: 5, Currency: "sp",
			Keys: "item.gear trait.advantage-on-saves-against-extreme-cold-while-wrapped"},
	}},
}}

// LoadItems256 returns a fresh validated copy of batch #256. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems256() (ItemSourceBatch[ItemDetails256], error) { return DecodeItems256(items256JSON) }

// DecodeItems256 validates supplied JSON as batch #256 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems256(data []byte) (ItemSourceBatch[ItemDetails256], error) {
	return decodeGearSourceBatch(data, itemManifest256, itemGear256)
}

func itemGear256(d *ItemDetails256) *ItemGear { return d.Gear }
