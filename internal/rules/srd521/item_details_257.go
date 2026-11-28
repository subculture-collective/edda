package srd521

import _ "embed"

//go:embed testdata/items_257.json
var items257JSON []byte

// ItemDetails257 holds the gear entry of each record in catalogue batch
// #257 (adventuring gear 02, SRD 5.2.1 pp. 96-97). Payloads are source
// facts only.
type ItemDetails257 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest257 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest257 = mundaneSourceManifest{Batch: "data-adventuring-gear-02", Issue: 257, Families: []mundaneSourceManifestFamily{
	{Slug: "block-and-tackle", Name: "Block and Tackle", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Block and Tackle", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.hoist-multiplier-4-times-normal-lift"},
	}},
	{Slug: "book", Name: "Book", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Book", CostAmount: 25, Currency: "gp",
			Keys: "item.gear measure.knowledge-check-bonus-5-bonus trait.contains-fiction-or-nonfiction trait.bonus-requires-consulting-accurate-nonfiction-about-topic trait.bonus-to-intelligence-arcana-history-nature-or-religion-checks"},
	}},
	{Slug: "bottle-glass", Name: "Bottle, Glass", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Bottle, Glass", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.capacity-3-over-2-pints"},
	}},
	{Slug: "bucket", Name: "Bucket", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Bucket", CostAmount: 5, Currency: "cp",
			Keys: "item.gear measure.capacity-volume-1-over-2-cubic-foot"},
	}},
	{Slug: "burglars-pack", Name: "Burglar’s Pack", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Burglar’s Pack", CostAmount: 16, Currency: "gp",
			Keys: "item.gear content.backpack-1-item content.ball-bearings-1-item content.bell-1-item content.candle-10-item content.crowbar-1-item content.lantern-hooded-1-item content.oil-7-flask content.rations-5-day content.rope-1-item content.tinderbox-1-item content.waterskin-1-item"},
	}},
	{Slug: "caltrops", Name: "Caltrops", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Caltrops", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.utilize saving-throw.avoid-damage-and-halt-dexterity-dc-15 damage.failed-save-1-piercing measure.covered-square-side-5-feet measure.placement-range-5-feet measure.failed-save-speed-0-speed measure.recovery-time-10-minutes trait.spread-from-bag trait.triggers-on-first-entry-on-a-turn trait.speed-penalty-lasts-until-start-of-next-turn"},
	}},
	{Slug: "candle", Name: "Candle", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Candle", CostAmount: 1, Currency: "cp",
			Keys: "item.gear measure.burn-duration-1-hour measure.bright-light-radius-5-feet measure.additional-dim-light-5-feet trait.sheds-light-while-lit"},
	}},
	{Slug: "case-crossbow-bolt", Name: "Case, Crossbow Bolt", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Case, Crossbow Bolt", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.capacity-20-bolts"},
	}},
	{Slug: "case-map-or-scroll", Name: "Case, Map or Scroll", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Case, Map or Scroll", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.paper-capacity-10-sheets measure.parchment-capacity-5-sheets trait.holds-paper-or-parchment"},
	}},
	{Slug: "chain", Name: "Chain", Page: 96, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Chain", CostAmount: 5, Currency: "gp",
			Keys: "item.gear action.utilize action.action ability-check.wrap-creature-strength-athletics-dc-13 ability-check.escape-dexterity-acrobatics-dc-18 ability-check.burst-strength-athletics-dc-20 measure.wrap-range-5-feet trait.target-unwilling-grappled-incapacitated-or-restrained-creature trait.bound-legs-restrain-until-escape"},
	}},
	{Slug: "chest", Name: "Chest", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Chest", CostAmount: 5, Currency: "gp",
			Keys: "item.gear measure.capacity-volume-12-cubic-feet"},
	}},
	{Slug: "climbers-kit", Name: "Climber’s Kit", Page: 97, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Climber’s Kit", CostAmount: 25, Currency: "gp",
			Keys: "item.gear action.utilize action.bonus-action measure.maximum-fall-from-anchor-25-feet measure.maximum-movement-from-anchor-25-feet trait.includes-boot-tips-gloves-pitons-and-harness trait.anchor-yourself trait.moving-farther-requires-undoing-anchor"},
	}},
}}

// LoadItems257 returns a fresh validated copy of batch #257. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems257() (ItemSourceBatch[ItemDetails257], error) { return DecodeItems257(items257JSON) }

// DecodeItems257 validates supplied JSON as batch #257 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems257(data []byte) (ItemSourceBatch[ItemDetails257], error) {
	return decodeGearSourceBatch(data, itemManifest257, itemGear257)
}

func itemGear257(d *ItemDetails257) *ItemGear { return d.Gear }
