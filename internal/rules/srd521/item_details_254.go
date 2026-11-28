package srd521

import _ "embed"

//go:embed testdata/items_254.json
var items254JSON []byte

// ItemDetails254 holds the tool entry of each record in catalogue batch
// #254 (tools 02, SRD 5.2.1 pp. 93–94). Payloads are source facts only.
type ItemDetails254 struct {
	Tool *ItemTool `json:"tool,omitempty"`
}

// itemManifest254 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest254 = mundaneSourceManifest{Batch: "data-tools-02", Issue: 254, Families: []mundaneSourceManifestFamily{
	{Slug: "jewelers-tools", Name: "Jeweler’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Jeweler’s Tools", CostAmount: 25, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-gems-value-dc-15 craft.arcane-focus craft.holy-symbol"},
	}},
	{Slug: "leatherworkers-tools", Name: "Leatherworker’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Leatherworker’s Tools", CostAmount: 5, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.add-design-to-leather-item-dc-10 craft.sling craft.whip craft.hide-armor craft.leather-armor craft.studded-leather-armor craft.backpack craft.crossbow-bolt-case craft.map-or-scroll-case craft.parchment craft.pouch craft.quiver craft.waterskin"},
	}},
	{Slug: "masons-tools", Name: "Mason’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Mason’s Tools", CostAmount: 10, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.chisel-symbol-or-hole-in-stone-dc-10 craft.block-and-tackle"},
	}},
	{Slug: "musical-instrument", Name: "Musical Instrument", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "bagpipes", Name: "Bagpipes", CostAmount: 30, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "drum", Name: "Drum", CostAmount: 6, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "dulcimer", Name: "Dulcimer", CostAmount: 25, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "flute", Name: "Flute", CostAmount: 2, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "horn", Name: "Horn", CostAmount: 3, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "lute", Name: "Lute", CostAmount: 35, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "lyre", Name: "Lyre", CostAmount: 30, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "pan-flute", Name: "Pan flute", CostAmount: 12, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "shawm", Name: "Shawm", CostAmount: 2, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{Slug: "viol", Name: "Viol", CostAmount: 30, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
	}},
	{Slug: "navigators-tools", Name: "Navigator’s Tools", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Navigator’s Tools", CostAmount: 25, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.plot-course-dc-10 utilize.determine-position-by-stargazing-dc-15"},
	}},
	{Slug: "painters-supplies", Name: "Painter’s Supplies", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Painter’s Supplies", CostAmount: 10, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.paint-recognizable-image-of-something-youve-seen-dc-10 craft.druidic-focus craft.holy-symbol"},
	}},
	{Slug: "poisoners-kit", Name: "Poisoner’s Kit", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Poisoner’s Kit", CostAmount: 50, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.intelligence utilize.detect-poisoned-object-dc-10 craft.basic-poison"},
	}},
	{Slug: "potters-tools", Name: "Potter’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Potter’s Tools", CostAmount: 10, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-what-ceramic-object-held-in-past-24-hours-dc-15 craft.jug craft.lamp"},
	}},
	{Slug: "smiths-tools", Name: "Smith’s Tools", Page: 93, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Smith’s Tools", CostAmount: 20, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.pry-open-door-or-container-dc-20 craft.any-melee-weapon-except-club-greatclub-quarterstaff-whip craft.medium-armor-except-hide craft.heavy-armor craft.ball-bearings craft.bucket craft.caltrops craft.chain craft.crowbar craft.firearm-bullets craft.grappling-hook craft.iron-pot craft.iron-spikes craft.sling-bullets"},
	}},
	{Slug: "thieves-tools", Name: "Thieves’ Tools", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Thieves’ Tools", CostAmount: 25, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.other-tools ability.dexterity utilize.pick-lock-dc-15 utilize.disarm-trap-dc-15"},
	}},
	{Slug: "tinkers-tools", Name: "Tinker’s Tools", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Tinker’s Tools", CostAmount: 50, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.assemble-tiny-item-composed-of-scrap-which-falls-apart-in-1-minute-dc-20 craft.musket craft.pistol craft.bell craft.bullseye-lantern craft.flask craft.hooded-lantern craft.hunting-trap craft.lock craft.manacles craft.mirror craft.shovel craft.signal-whistle craft.tinderbox"},
	}},
	{Slug: "weavers-tools", Name: "Weaver’s Tools", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Weaver’s Tools", CostAmount: 1, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.mend-tear-in-clothing-dc-10 utilize.sew-tiny-design-dc-10 craft.padded-armor craft.basket craft.bedroll craft.blanket craft.fine-clothes craft.net craft.robe craft.rope craft.sack craft.string craft.tent craft.travelers-clothes"},
	}},
}}

// LoadItems254 returns a fresh validated copy of batch #254. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems254() (ItemSourceBatch[ItemDetails254], error) { return DecodeItems254(items254JSON) }

// DecodeItems254 validates supplied JSON as batch #254 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems254(data []byte) (ItemSourceBatch[ItemDetails254], error) {
	return decodeToolSourceBatch(data, itemManifest254, func(d *ItemDetails254) *ItemTool { return d.Tool })
}
