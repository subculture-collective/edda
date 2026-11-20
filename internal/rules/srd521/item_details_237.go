package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_237.json
var items237JSON []byte

// ItemDetails237 holds exactly one typed payload per family of catalogue batch
// #237 (magic items 13, SRD 5.2.1 pp. 235-236). Payloads are source facts only.
type ItemDetails237 struct {
	Clairvoyance  *ItemPotionClairvoyance  `json:"clairvoyance,omitempty"`
	Climbing      *ItemPotionClimbing      `json:"climbing,omitempty"`
	Diminution    *ItemPotionDiminution    `json:"diminution,omitempty"`
	Flying        *ItemPotionFlying        `json:"flying,omitempty"`
	GaseousForm   *ItemPotionGaseousForm   `json:"gaseous_form,omitempty"`
	GiantStrength *ItemPotionGiantStrength `json:"giant_strength,omitempty"`
	Growth        *ItemPotionGrowth        `json:"growth,omitempty"`
	Heroism       *ItemPotionHeroism       `json:"heroism,omitempty"`
	Invisibility  *ItemPotionInvisibility  `json:"invisibility,omitempty"`
	Invuln        *ItemPotionInvulnerable  `json:"invulnerability,omitempty"`
	Longevity     *ItemPotionLongevity     `json:"longevity,omitempty"`
	MindReading   *ItemPotionMindReading   `json:"mind_reading,omitempty"`
}

// ItemPotionClairvoyance retains typed SRD source facts without implementing them.
type ItemPotionClairvoyance struct {
	Spell                 string `json:"spell"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// ItemPotionClimbing retains typed SRD source facts without implementing them.
type ItemPotionClimbing struct {
	SpeedKind     string `json:"speed_kind"`
	SpeedEquals   string `json:"speed_equals"`
	DurationHours int    `json:"duration_hours"`
	CheckAbility  string `json:"check_ability"`
	CheckSkill    string `json:"check_skill"`
	CheckMode     string `json:"check_mode"`
	CheckSubject  string `json:"check_subject"`
	Appearance    string `json:"appearance"`
}

// ItemPotionDiminution retains typed SRD source facts without implementing them.
type ItemPotionDiminution struct {
	Spell                 string `json:"spell"`
	SpellEffect           string `json:"spell_effect"`
	DurationDice          int    `json:"duration_dice"`
	DurationDieSides      int    `json:"duration_die_sides"`
	DurationUnit          string `json:"duration_unit"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// ItemPotionFlying retains typed SRD source facts without implementing them.
type ItemPotionFlying struct {
	SpeedKind       string `json:"speed_kind"`
	SpeedEquals     string `json:"speed_equals"`
	DurationHours   int    `json:"duration_hours"`
	CanHover        bool   `json:"can_hover"`
	FallsWhenItEnds string `json:"falls_when_it_ends"`
	FallException   string `json:"fall_exception"`
	Appearance      string `json:"appearance"`
}

// ItemPotionGaseousForm retains typed SRD source facts without implementing them.
type ItemPotionGaseousForm struct {
	Spell                 string `json:"spell"`
	DurationHours         int    `json:"duration_hours"`
	EndsEarlyAction       string `json:"ends_early_action"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// ItemPotionGiantStrength is one row of the source's giant table: a variant per
// row. A row printed "frost or stone" lists two giants and is one variant.
type ItemPotionGiantStrength struct {
	Giants                 []string `json:"giants"`
	StrengthScore          int      `json:"strength_score"`
	DurationHours          int      `json:"duration_hours"`
	ScoreChange            string   `json:"score_change"`
	NoEffectAtOrAboveScore int      `json:"no_effect_at_or_above_score"`
	Appearance             string   `json:"appearance"`
}

// ItemPotionGrowth retains typed SRD source facts without implementing them.
type ItemPotionGrowth struct {
	Spell                 string `json:"spell"`
	SpellEffect           string `json:"spell_effect"`
	DurationMinutes       int    `json:"duration_minutes"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// ItemPotionHeroism retains typed SRD source facts without implementing them.
type ItemPotionHeroism struct {
	TemporaryHitPoints    int    `json:"temporary_hit_points"`
	DurationHours         int    `json:"duration_hours"`
	Spell                 string `json:"spell"`
	SpellDuration         string `json:"spell_duration"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// ItemPotionInvisibility retains typed SRD source facts without implementing them.
type ItemPotionInvisibility struct {
	Condition     string   `json:"condition"`
	DurationHours int      `json:"duration_hours"`
	EndsEarlyOn   []string `json:"ends_early_on"`
	Appearance    string   `json:"appearance"`
}

// ItemPotionInvulnerable retains typed SRD source facts without implementing them.
type ItemPotionInvulnerable struct {
	Resistance      string `json:"resistance"`
	DurationMinutes int    `json:"duration_minutes"`
	Appearance      string `json:"appearance"`
}

// ItemAgeRoll is a source roll of years: dice, die sides and a flat bonus.
type ItemAgeRoll struct {
	Dice  int `json:"dice"`
	Sides int `json:"sides"`
	Bonus int `json:"bonus"`
}

// ItemPotionLongevity retains typed SRD source facts without implementing them.
type ItemPotionLongevity struct {
	AgeReduction          ItemAgeRoll `json:"age_reduction_years"`
	MinimumAgeYears       int         `json:"minimum_age_years"`
	AgingChancePercent    int         `json:"aging_chance_percent"`
	AgingChanceCumulative bool        `json:"aging_chance_cumulative"`
	AgingChanceTrigger    string      `json:"aging_chance_trigger"`
	AgingInstead          ItemAgeRoll `json:"aging_instead_years"`
	Appearance            string      `json:"appearance"`
}

// ItemPotionMindReading retains typed SRD source facts without implementing them.
type ItemPotionMindReading struct {
	Spell                 string `json:"spell"`
	SaveDC                int    `json:"save_dc"`
	DurationMinutes       int    `json:"duration_minutes"`
	ConcentrationRequired bool   `json:"concentration_required"`
	Appearance            string `json:"appearance"`
}

// itemManifest237 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Potion of Giant Strength has five variants, one per
// row of its table; every other family has the single variant `standard`.
var itemManifest237 = itemSourceManifest{Batch: "data-magic-items-13", Issue: 237, Families: []itemSourceManifestFamily{
	{Slug: "potion-of-clairvoyance", Name: "Potion of Clairvoyance", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-clairvoyance spell.no-concentration", Payload: "clairvoyance",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Clairvoyance", "rare"}}},
	{Slug: "potion-of-climbing", Name: "Potion of Climbing", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume speed.climb-equal-to-speed duration.1-hour check.strength-athletics-advantage-climbing", Payload: "climbing",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Climbing", "common"}}},
	{Slug: "potion-of-diminution", Name: "Potion of Diminution", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-enlarge-reduce-reduce duration.1d4-hours spell.no-concentration", Payload: "diminution",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Diminution", "rare"}}},
	{Slug: "potion-of-flying", Name: "Potion of Flying", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume speed.fly-equal-to-speed duration.1-hour movement.hover movement.fall-when-ends-in-air", Payload: "flying",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Flying", "very-rare"}}},
	{Slug: "potion-of-gaseous-form", Name: "Potion of Gaseous Form", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-gaseous-form duration.1-hour action.bonus-end-effect spell.no-concentration", Payload: "gaseous_form",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Gaseous Form", "rare"}}},
	{Slug: "potion-of-giant-strength", Name: "Potion of Giant Strength", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score", Payload: "giant_strength",
		Variants: []itemSourceManifestVariant{
			{"hill", "Potion of Giant Strength (hill)", "uncommon"}, {"frost-or-stone", "Potion of Giant Strength (frost or stone)", "rare"},
			{"fire", "Potion of Giant Strength (fire)", "rare"}, {"cloud", "Potion of Giant Strength (cloud)", "very-rare"}, {"storm", "Potion of Giant Strength (storm)", "legendary"}}},
	{Slug: "potion-of-growth", Name: "Potion of Growth", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-enlarge-reduce-enlarge duration.10-minutes spell.no-concentration", Payload: "growth",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Growth", "uncommon"}}},
	{Slug: "potion-of-heroism", Name: "Potion of Heroism", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume hit-points.temporary-10 duration.1-hour spell.effect-bless spell.no-concentration", Payload: "heroism",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Heroism", "rare"}}},
	{Slug: "potion-of-invisibility", Name: "Potion of Invisibility", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume condition.invisible duration.1-hour effect.ends-early-on-attack-damage-or-spell", Payload: "invisibility",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Invisibility", "rare"}}},
	{Slug: "potion-of-invulnerability", Name: "Potion of Invulnerability", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume defense.resistance-all-damage duration.1-minute", Payload: "invulnerability",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Invulnerability", "rare"}}},
	{Slug: "potion-of-longevity", Name: "Potion of Longevity", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume age.reduce-1d6-plus-6-years age.minimum-13-years age.cumulative-10-percent-chance-to-age-instead", Payload: "longevity",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Longevity", "very-rare"}}},
	{Slug: "potion-of-mind-reading", Name: "Potion of Mind Reading", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-detect-thoughts spell.detect-thoughts-save-dc-13 duration.10-minutes spell.no-concentration", Payload: "mind_reading",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Mind Reading", "rare"}}},
}}

// LoadItems237 returns a fresh validated copy of batch #237. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems237() (ItemSourceBatch[ItemDetails237], error) { return DecodeItems237(items237JSON) }

// DecodeItems237 validates supplied JSON as batch #237 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems237(data []byte) (ItemSourceBatch[ItemDetails237], error) {
	return decodeItemSourceBatch(data, itemManifest237, validateItemDetails237)
}

// validateItemDetails237 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails237(family, variant string, d *ItemDetails237) error {
	ok := true
	switch family {
	case "potion-of-clairvoyance":
		ok = !d.Clairvoyance.ConcentrationRequired
	case "potion-of-climbing":
		ok = d.Climbing.DurationHours == 1
	case "potion-of-diminution":
		p := d.Diminution
		ok = p.DurationDice == 1 && p.DurationDieSides == 4 && !p.ConcentrationRequired
	case "potion-of-flying":
		ok = d.Flying.DurationHours == 1 && d.Flying.CanHover
	case "potion-of-gaseous-form":
		ok = d.GaseousForm.DurationHours == 1 && !d.GaseousForm.ConcentrationRequired
	case "potion-of-giant-strength":
		g := d.GiantStrength
		// One row of the p. 235 table per variant: giants, Strength score.
		rows := map[string]struct {
			giants []string
			score  int
		}{
			"hill": {[]string{"hill"}, 21}, "frost-or-stone": {[]string{"frost", "stone"}, 23}, "fire": {[]string{"fire"}, 25},
			"cloud": {[]string{"cloud"}, 27}, "storm": {[]string{"storm"}, 29},
		}
		row := rows[variant]
		ok = slices.Equal(g.Giants, row.giants) && g.StrengthScore == row.score && g.NoEffectAtOrAboveScore == row.score && g.DurationHours == 1
	case "potion-of-growth":
		ok = d.Growth.DurationMinutes == 10 && !d.Growth.ConcentrationRequired
	case "potion-of-heroism":
		p := d.Heroism
		ok = p.TemporaryHitPoints == 10 && p.DurationHours == 1 && !p.ConcentrationRequired
	case "potion-of-invisibility":
		ok = d.Invisibility.DurationHours == 1 && uniqueItemKeys(d.Invisibility.EndsEarlyOn)
	case "potion-of-invulnerability":
		ok = d.Invuln.DurationMinutes == 1
	case "potion-of-longevity":
		p := d.Longevity
		ok = p.AgeReduction == ItemAgeRoll{Dice: 1, Sides: 6, Bonus: 6} && p.AgingInstead == p.AgeReduction && p.MinimumAgeYears == 13 && p.AgingChancePercent == 10 && p.AgingChanceCumulative
	case "potion-of-mind-reading":
		p := d.MindReading
		ok = p.SaveDC == 13 && p.DurationMinutes == 10 && !p.ConcentrationRequired
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
