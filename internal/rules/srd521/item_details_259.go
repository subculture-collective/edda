package srd521

import _ "embed"

//go:embed testdata/items_259.json
var items259JSON []byte

// ItemDetails259 holds the gear entry of each record in catalogue batch
// #259 (adventuring gear 04, SRD 5.2.1 pp. 97-98). Payloads are source
// facts only.
type ItemDetails259 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest259 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest259 = mundaneSourceManifest{Batch: "data-adventuring-gear-04", Issue: 259, Families: []mundaneSourceManifestFamily{
	{Slug: "healers-kit", Name: "Healer’s Kit", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Healer’s Kit", CostAmount: 5, Currency: "gp",
			Keys: "item.gear action.utilize measure.uses-10-uses measure.stabilize-target-hit-points-0-hit-points trait.expend-one-use-to-stabilize-unconscious-creature trait.no-wisdom-medicine-check-needed"},
	}},
	{Slug: "holy-symbol", Name: "Holy Symbol", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "amulet", Name: "Amulet", CostAmount: 5, Currency: "gp",
			Keys: "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.worn-or-held"},
		{Slug: "emblem", Name: "Emblem", CostAmount: 5, Currency: "gp",
			Keys: "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.borne-on-fabric-or-shield"},
		{Slug: "reliquary", Name: "Reliquary", CostAmount: 5, Currency: "gp",
			Keys: "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.held"},
	}},
	{Slug: "holy-water", Name: "Holy Water", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Holy Water", CostAmount: 25, Currency: "gp",
			Keys: "item.gear action.attack-replacement saving-throw.avoid-damage-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-flask-against-fiend-or-undead-2d8-radiant measure.throw-range-20-feet trait.throw-flask-in-place-of-one-attack trait.target-one-creature-you-can-see trait.damages-only-fiend-or-undead"},
	}},
	{Slug: "hunting-trap", Name: "Hunting Trap", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Hunting Trap", CostAmount: 5, Currency: "gp",
			Keys: "item.gear action.utilize action.action saving-throw.avoid-damage-and-halt-dexterity-dc-13 ability-check.free-creature-strength-athletics-dc-13 damage.failed-save-1d4-piercing damage.each-failed-escape-check-1-piercing measure.failed-save-speed-0-speed measure.typical-chain-length-3-feet trait.sawtooth-steel-ring-snaps-shut-on-pressure-plate trait.chain-affixed-to-immobile-object trait.speed-penalty-lasts-until-start-of-next-turn trait.movement-limited-by-chain-until-free trait.check-frees-self-or-creature-within-reach"},
	}},
	{Slug: "ink", Name: "Ink", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ink", CostAmount: 10, Currency: "gp",
			Keys: "item.gear measure.bottle-volume-1-ounce measure.writable-pages-500-pages"},
	}},
	{Slug: "ink-pen", Name: "Ink Pen", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ink Pen", CostAmount: 2, Currency: "cp",
			Keys: "item.gear trait.writes-or-draws-with-ink"},
	}},
	{Slug: "jug", Name: "Jug", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Jug", CostAmount: 2, Currency: "cp",
			Keys: "item.gear measure.capacity-1-gallon"},
	}},
	{Slug: "ladder", Name: "Ladder", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Ladder", CostAmount: 1, Currency: "sp",
			Keys: "item.gear measure.height-10-feet trait.must-climb-to-move-up-or-down"},
	}},
	{Slug: "lamp", Name: "Lamp", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Lamp", CostAmount: 5, Currency: "sp",
			Keys: "item.gear measure.bright-light-radius-15-feet measure.additional-dim-light-30-feet trait.burns-oil-as-fuel"},
	}},
	{Slug: "lantern-bullseye", Name: "Lantern, Bullseye", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Lantern, Bullseye", CostAmount: 10, Currency: "gp",
			Keys: "item.gear measure.bright-light-cone-60-feet measure.additional-dim-light-60-feet trait.burns-oil-as-fuel"},
	}},
	{Slug: "lantern-hooded", Name: "Lantern, Hooded", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Lantern, Hooded", CostAmount: 5, Currency: "gp",
			Keys: "item.gear action.bonus-action measure.bright-light-radius-30-feet measure.additional-dim-light-30-feet measure.lowered-hood-dim-light-radius-5-feet trait.burns-oil-as-fuel trait.hood-can-be-lowered-or-raised"},
	}},
	{Slug: "lock", Name: "Lock", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Lock", CostAmount: 10, Currency: "gp",
			Keys: "item.gear ability-check.pick-lock-without-key-dexterity-sleight-of-hand-dc-15 trait.comes-with-key trait.thieves-tools-can-pick-without-key"},
	}},
}}

// LoadItems259 returns a fresh validated copy of batch #259. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems259() (ItemSourceBatch[ItemDetails259], error) { return DecodeItems259(items259JSON) }

// DecodeItems259 validates supplied JSON as batch #259 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems259(data []byte) (ItemSourceBatch[ItemDetails259], error) {
	return decodeGearSourceBatch(data, itemManifest259, itemGear259)
}

func itemGear259(d *ItemDetails259) *ItemGear { return d.Gear }
