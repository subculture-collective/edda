package srd521

import _ "embed"

//go:embed testdata/items_255.json
var items255JSON []byte

// ItemDetails255 holds the tool entry of each record in catalogue batch
// #255 (tools 03, SRD 5.2.1 p. 94). Payloads are source facts only.
type ItemDetails255 struct {
	Tool *ItemTool `json:"tool,omitempty"`
}

// itemManifest255 is the reviewed manifest for the batch, in planning
// inventory order.
var itemManifest255 = mundaneSourceManifest{Batch: "data-tools-03", Issue: 255, Families: []mundaneSourceManifestFamily{
	{Slug: "woodcarvers-tools", Name: "Woodcarver’s Tools", Page: 94, Variants: []mundaneSourceManifestVariant{
		{Slug: "standard", Name: "Woodcarver’s Tools", CostAmount: 1, Currency: "gp",
			Keys: "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.carve-pattern-in-wood-dc-10 craft.club craft.greatclub craft.quarterstaff craft.ranged-weapons-except-pistol-musket-sling craft.arcane-focus craft.arrows craft.bolts craft.druidic-focus craft.ink-pen craft.needles"},
	}},
}}

// LoadItems255 returns a fresh validated copy of batch #255. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems255() (ItemSourceBatch[ItemDetails255], error) { return DecodeItems255(items255JSON) }

// DecodeItems255 validates supplied JSON as batch #255 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems255(data []byte) (ItemSourceBatch[ItemDetails255], error) {
	return decodeToolSourceBatch(data, itemManifest255, func(d *ItemDetails255) *ItemTool { return d.Tool })
}
