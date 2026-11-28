package srd521

import _ "embed"

//go:embed testdata/items_260.json
var items260JSON []byte

// ItemDetails260 holds the gear entry of each record in catalogue batch
// #260 (adventuring gear 05, SRD 5.2.1 pp. 98-99). Payloads are source
// facts only.
type ItemDetails260 struct {
	Gear *ItemGear `json:"gear,omitempty"`
}

// itemManifest260 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest260 = mundaneSourceManifest{Batch: "data-adventuring-gear-05", Issue: 260, Families: []mundaneSourceManifestFamily{
	{Slug: "magnifying-glass", Name: "Magnifying Glass", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Magnifying Glass", CostAmount: 100, Currency: "gp",
			Keys: "item.gear measure.fire-ignition-time-5-minutes trait.advantage-on-checks-to-appraise-or-inspect-highly-detailed-item trait.lighting-fire-requires-sunlight-bright-light-and-tinder"},
	}},
	{Slug: "manacles", Name: "Manacles", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Manacles", CostAmount: 2, Currency: "gp",
			Keys: "item.gear action.utilize action.action ability-check.bind-creature-dexterity-sleight-of-hand-dc-13 ability-check.escape-dexterity-sleight-of-hand-dc-20 ability-check.burst-strength-athletics-dc-25 ability-check.pick-lock-without-key-dexterity-sleight-of-hand-dc-15 measure.bind-range-5-feet trait.target-unwilling-small-or-medium-grappled-incapacitated-or-restrained-creature trait.bound-creature-has-disadvantage-on-attack-rolls trait.restrained-if-attached-to-fixed-chain-or-hook trait.comes-with-key trait.thieves-tools-can-pick-without-key"},
	}},
	{Slug: "map", Name: "Map", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Map", CostAmount: 1, Currency: "gp",
			Keys: "item.gear measure.wayfinding-check-bonus-5-bonus trait.bonus-requires-consulting-accurate-map trait.bonus-to-wisdom-survival-checks-to-find-way-in-mapped-place"},
	}},
	{Slug: "mirror", Name: "Mirror", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Mirror", CostAmount: 5, Currency: "gp",
			Keys: "item.gear trait.handheld-steel trait.useful-for-cosmetics-peeking-around-corners-and-signals"},
	}},
	{Slug: "net", Name: "Net", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Net", CostAmount: 1, Currency: "gp",
			Keys: "item.gear action.attack-replacement action.action saving-throw.avoid-restrained-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus ability-check.escape-strength-athletics-dc-10 measure.throw-range-15-feet measure.escape-helper-range-5-feet measure.armor-class-10-armor-class measure.hit-points-5-hit-points trait.throw-net-in-place-of-one-attack trait.target-one-creature-you-can-see trait.failed-save-restrains-until-escape trait.huge-or-larger-target-succeeds-automatically trait.successful-check-frees-restrained-creature trait.immune-to-bludgeoning-poison-and-psychic-damage trait.destroying-net-frees-target-and-ends-effect"},
	}},
	{Slug: "oil", Name: "Oil", Page: 98, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Oil", CostAmount: 1, Currency: "sp",
			Keys: "item.gear action.attack-replacement action.utilize saving-throw.avoid-being-covered-in-oil-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.extra-from-burning-oil-5-fire damage.burning-area-5-fire measure.throw-range-20-feet measure.drying-time-1-minute measure.covered-square-side-5-feet measure.pour-range-5-feet measure.area-burn-rounds-2-rounds measure.area-burn-seconds-12-seconds measure.fuel-burn-duration-6-hours trait.douse-creature-object-or-space-or-use-as-fuel trait.throw-flask-in-place-of-one-attack trait.target-one-creature-or-object trait.extra-damage-if-fire-damage-taken-before-oil-dries trait.pour-flask-on-level-ground trait.area-damages-on-entry-or-turn-end-once-per-turn trait.fuel-for-lamps-and-lanterns trait.fuel-duration-need-not-be-consecutive trait.extinguish-and-rekindle-until-total-duration-burned"},
	}},
	{Slug: "paper", Name: "Paper", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Paper", CostAmount: 2, Currency: "sp",
			Keys: "item.gear measure.handwritten-words-per-sheet-250-words"},
	}},
	{Slug: "parchment", Name: "Parchment", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Parchment", CostAmount: 1, Currency: "sp",
			Keys: "item.gear measure.handwritten-words-per-sheet-250-words"},
	}},
	{Slug: "perfume", Name: "Perfume", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Perfume", CostAmount: 5, Currency: "gp",
			Keys: "item.gear measure.vial-volume-4-ounces measure.benefit-duration-1-hour measure.influence-range-5-feet trait.advantage-on-charisma-persuasion-to-influence-indifferent-humanoid"},
	}},
	{Slug: "poison-basic", Name: "Poison, Basic", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Poison, Basic", CostAmount: 100, Currency: "gp",
			Keys: "item.gear action.bonus-action damage.extra-from-poisoned-weapon-or-ammunition-1d4-poison measure.coated-weapons-1-weapon measure.coated-ammunition-3-pieces measure.potency-duration-1-minute trait.coat-one-weapon-or-ammunition-from-vial trait.applies-on-piercing-or-slashing-damage trait.potency-ends-early-when-damage-is-dealt"},
	}},
	{Slug: "pole", Name: "Pole", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Pole", CostAmount: 5, Currency: "cp",
			Keys: "item.gear measure.length-10-feet measure.touch-range-10-feet trait.advantage-when-vaulting-in-high-or-long-jump-athletics-check"},
	}},
	{Slug: "pot-iron", Name: "Pot, Iron", Page: 99, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Pot, Iron", CostAmount: 2, Currency: "gp",
			Keys: "item.gear measure.capacity-1-gallon"},
	}},
}}

// LoadItems260 returns a fresh validated copy of batch #260. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems260() (ItemSourceBatch[ItemDetails260], error) { return DecodeItems260(items260JSON) }

// DecodeItems260 validates supplied JSON as batch #260 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems260(data []byte) (ItemSourceBatch[ItemDetails260], error) {
	return decodeGearSourceBatch(data, itemManifest260, itemGear260)
}

func itemGear260(d *ItemDetails260) *ItemGear { return d.Gear }
