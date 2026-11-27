package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_245.json
var items245JSON []byte

// ItemDetails245 holds exactly one typed payload per variant of catalogue batch
// #245 (magic items 21, SRD 5.2.1 pp. 250-252). Payloads are source facts only.
type ItemDetails245 struct {
	EnemyDetection *ItemWandOfEnemyDetection `json:"enemy_detection,omitempty"`
	Fear           *ItemWandOfFear           `json:"fear,omitempty"`
	Fireballs      *ItemWandOfFireballs      `json:"fireballs,omitempty"`
	LightningBolts *ItemWandOfLightningBolts `json:"lightning_bolts,omitempty"`
	MagicDetection *ItemWandOfMagicDetection `json:"magic_detection,omitempty"`
	MagicMissiles  *ItemWandOfMagicMissiles  `json:"magic_missiles,omitempty"`
	Paralysis      *ItemWandOfParalysis      `json:"paralysis,omitempty"`
	Polymorph      *ItemWandSingleSpell      `json:"polymorph,omitempty"`
	Secrets        *ItemWandOfSecrets        `json:"secrets,omitempty"`
	Web            *ItemWandSingleSpell      `json:"web,omitempty"`
	Wonder         *ItemWandOfWonder         `json:"wonder,omitempty"`
	WarMage        *ItemWandOfTheWarMage     `json:"war_mage,omitempty"`
}

// ItemWandOfEnemyDetection retains typed SRD source facts without implementing them.
type ItemWandOfEnemyDetection struct {
	MaxCharges      int                                `json:"max_charges"`
	Action          string                             `json:"action"`
	ChargeCost      int                                `json:"charge_cost"`
	DurationMinutes int                                `json:"duration_minutes"`
	Reveals         string                             `json:"reveals"`
	RangeFeet       int                                `json:"range_feet"`
	RevealsDistance bool                               `json:"reveals_distance"`
	Senses          []string                           `json:"senses"`
	EndsWhen        string                             `json:"ends_when"`
	Recharge        ItemWandOfEnemyDetectionRecharge   `json:"recharge"`
	LastCharge      ItemWandOfEnemyDetectionLastCharge `json:"last_charge"`
}

// ItemWandOfEnemyDetectionRecharge is a nested part of a source payload.
type ItemWandOfEnemyDetectionRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfEnemyDetectionLastCharge is a nested part of a source payload.
type ItemWandOfEnemyDetectionLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfFear retains typed SRD source facts without implementing them.
type ItemWandOfFear struct {
	MaxCharges int                      `json:"max_charges"`
	SaveDC     int                      `json:"save_dc"`
	Spells     []ItemWandOfFearSpell    `json:"spells"`
	Recharge   ItemWandOfFearRecharge   `json:"recharge"`
	LastCharge ItemWandOfFearLastCharge `json:"last_charge"`
}

// ItemWandOfFearSpell is a nested part of a source payload.
type ItemWandOfFearSpell struct {
	Spell      string `json:"spell"`
	Limit      string `json:"limit"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemWandOfFearRecharge is a nested part of a source payload.
type ItemWandOfFearRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfFearLastCharge is a nested part of a source payload.
type ItemWandOfFearLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfFireballs retains typed SRD source facts without implementing them.
type ItemWandOfFireballs struct {
	AttunementByAnyOf         []string                      `json:"attunement_by_any_of"`
	MaxCharges                int                           `json:"max_charges"`
	MaximumChargesPerCast     int                           `json:"maximum_charges_per_cast"`
	Spell                     string                        `json:"spell"`
	SaveDC                    int                           `json:"save_dc"`
	BaseChargeCost            int                           `json:"base_charge_cost"`
	BaseSpellLevel            int                           `json:"base_spell_level"`
	LevelsPerAdditionalCharge int                           `json:"levels_per_additional_charge"`
	Recharge                  ItemWandOfFireballsRecharge   `json:"recharge"`
	LastCharge                ItemWandOfFireballsLastCharge `json:"last_charge"`
	SourcePages               []int                         `json:"source_pages"`
}

// ItemWandOfFireballsRecharge is a nested part of a source payload.
type ItemWandOfFireballsRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfFireballsLastCharge is a nested part of a source payload.
type ItemWandOfFireballsLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfLightningBolts retains typed SRD source facts without implementing them.
type ItemWandOfLightningBolts struct {
	AttunementByAnyOf         []string                           `json:"attunement_by_any_of"`
	MaxCharges                int                                `json:"max_charges"`
	MaximumChargesPerCast     int                                `json:"maximum_charges_per_cast"`
	Spell                     string                             `json:"spell"`
	SaveDC                    int                                `json:"save_dc"`
	BaseChargeCost            int                                `json:"base_charge_cost"`
	BaseSpellLevel            int                                `json:"base_spell_level"`
	LevelsPerAdditionalCharge int                                `json:"levels_per_additional_charge"`
	Recharge                  ItemWandOfLightningBoltsRecharge   `json:"recharge"`
	LastCharge                ItemWandOfLightningBoltsLastCharge `json:"last_charge"`
}

// ItemWandOfLightningBoltsRecharge is a nested part of a source payload.
type ItemWandOfLightningBoltsRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfLightningBoltsLastCharge is a nested part of a source payload.
type ItemWandOfLightningBoltsLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfMagicDetection retains typed SRD source facts without implementing them.
type ItemWandOfMagicDetection struct {
	MaxCharges   int      `json:"max_charges"`
	ChargeCost   int      `json:"charge_cost"`
	Spell        string   `json:"spell"`
	RechargeDice ItemDice `json:"recharge_dice"`
}

// ItemWandOfMagicMissiles retains typed SRD source facts without implementing them.
type ItemWandOfMagicMissiles struct {
	MaxCharges                int                               `json:"max_charges"`
	MaximumChargesPerCast     int                               `json:"maximum_charges_per_cast"`
	Spell                     string                            `json:"spell"`
	BaseChargeCost            int                               `json:"base_charge_cost"`
	BaseSpellLevel            int                               `json:"base_spell_level"`
	LevelsPerAdditionalCharge int                               `json:"levels_per_additional_charge"`
	Recharge                  ItemWandOfMagicMissilesRecharge   `json:"recharge"`
	LastCharge                ItemWandOfMagicMissilesLastCharge `json:"last_charge"`
}

// ItemWandOfMagicMissilesRecharge is a nested part of a source payload.
type ItemWandOfMagicMissilesRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfMagicMissilesLastCharge is a nested part of a source payload.
type ItemWandOfMagicMissilesLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfParalysis retains typed SRD source facts without implementing them.
type ItemWandOfParalysis struct {
	AttunementByAnyOf []string                      `json:"attunement_by_any_of"`
	MaxCharges        int                           `json:"max_charges"`
	Action            string                        `json:"action"`
	ChargeCost        int                           `json:"charge_cost"`
	Target            string                        `json:"target"`
	RangeFeet         int                           `json:"range_feet"`
	SaveAbility       string                        `json:"save_ability"`
	SaveDC            int                           `json:"save_dc"`
	Condition         string                        `json:"condition"`
	DurationMinutes   int                           `json:"duration_minutes"`
	RepeatSave        string                        `json:"repeat_save"`
	RepeatSuccess     string                        `json:"repeat_success"`
	Recharge          ItemWandOfParalysisRecharge   `json:"recharge"`
	LastCharge        ItemWandOfParalysisLastCharge `json:"last_charge"`
}

// ItemWandOfParalysisRecharge is a nested part of a source payload.
type ItemWandOfParalysisRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfParalysisLastCharge is a nested part of a source payload.
type ItemWandOfParalysisLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandSingleSpell retains typed SRD source facts without implementing them.
type ItemWandSingleSpell struct {
	AttunementByAnyOf []string                      `json:"attunement_by_any_of"`
	MaxCharges        int                           `json:"max_charges"`
	ChargeCost        int                           `json:"charge_cost"`
	Spell             string                        `json:"spell"`
	SaveDC            int                           `json:"save_dc"`
	Recharge          ItemWandSingleSpellRecharge   `json:"recharge"`
	LastCharge        ItemWandSingleSpellLastCharge `json:"last_charge"`
}

// ItemWandSingleSpellRecharge is a nested part of a source payload.
type ItemWandSingleSpellRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandSingleSpellLastCharge is a nested part of a source payload.
type ItemWandSingleSpellLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfSecrets retains typed SRD source facts without implementing them.
type ItemWandOfSecrets struct {
	MaxCharges   int      `json:"max_charges"`
	RechargeDice ItemDice `json:"recharge_dice"`
	Action       string   `json:"action"`
	ChargeCost   int      `json:"charge_cost"`
	Detects      []string `json:"detects"`
	RangeFeet    int      `json:"range_feet"`
	Result       string   `json:"result"`
}

// ItemWandOfWonder retains typed SRD source facts without implementing them.
type ItemWandOfWonder struct {
	MaxCharges            int                         `json:"max_charges"`
	Action                string                      `json:"action"`
	ChargeCost            int                         `json:"charge_cost"`
	PointRangeFeet        int                         `json:"point_range_feet"`
	SpellSaveDC           int                         `json:"spell_save_dc"`
	MinimumSpellRangeFeet int                         `json:"minimum_spell_range_feet"`
	MultipleSubjects      string                      `json:"multiple_subjects"`
	TableDieSides         int                         `json:"table_die_sides"`
	Effects               []ItemWandOfWonderEffect    `json:"effects"`
	SpellTable            ItemWandRollTable           `json:"spell_table"`
	Stunned               ItemWandOfWonderStunned     `json:"stunned"`
	GustLine              string                      `json:"gust_line"`
	PsychicDamage         ItemDice                    `json:"psychic_damage"`
	Rain                  ItemWandOfWonderRain        `json:"rain"`
	Butterflies           ItemWandOfWonderButterflies `json:"butterflies"`
	LightningLine         string                      `json:"lightning_line"`
	Enlarge               ItemWandOfWonderEnlarge     `json:"enlarge"`
	Creature              ItemWandOfWonderCreature    `json:"creature"`
	Grass                 ItemWandOfWonderGrass       `json:"grass"`
	Ethereal              ItemWandOfWonderEthereal    `json:"ethereal"`
	Shrink                ItemWandOfWonderShrink      `json:"shrink"`
	Leaves                ItemWandOfWonderLeaves      `json:"leaves"`
	LightBurst            ItemWandOfWonderLightBurst  `json:"light_burst"`
	InvisibilityTarget    string                      `json:"invisibility_target"`
	Gems                  ItemWandOfWonderGems        `json:"gems"`
	PolymorphTable        ItemWandRollTable           `json:"polymorph_table"`
	Petrify               ItemWandOfWonderPetrify     `json:"petrify"`
	LastCharge            ItemWandOfWonderLastCharge  `json:"last_charge"`
	Recharge              ItemWandOfWonderRecharge    `json:"recharge"`
	SourcePages           []int                       `json:"source_pages"`
}

// ItemWandOfWonderEffect is a nested part of a source payload.
type ItemWandOfWonderEffect struct {
	RollFrom int    `json:"roll_from"`
	RollTo   int    `json:"roll_to"`
	Effect   string `json:"effect"`
}

// ItemWandRollTable is a nested part of a source payload.
type ItemWandRollTable struct {
	DieSides int               `json:"die_sides"`
	Rows     []ItemWandRollRow `json:"rows"`
}

// ItemWandRollRow is a nested part of a source payload.
type ItemWandRollRow struct {
	RollFrom int    `json:"roll_from"`
	RollTo   int    `json:"roll_to"`
	Result   string `json:"result"`
}

// ItemWandOfWonderStunned is a nested part of a source payload.
type ItemWandOfWonderStunned struct {
	Condition string `json:"condition"`
	Until     string `json:"until"`
	Belief    string `json:"belief"`
}

// ItemWandOfWonderRain is a nested part of a source payload.
type ItemWandOfWonderRain struct {
	DurationMinutes int    `json:"duration_minutes"`
	HeightFeet      int    `json:"height_feet"`
	RadiusFeet      int    `json:"radius_feet"`
	Area            string `json:"area"`
}

// ItemWandOfWonderButterflies is a nested part of a source payload.
type ItemWandOfWonderButterflies struct {
	Count           int    `json:"count"`
	HeightFeet      int    `json:"height_feet"`
	RadiusFeet      int    `json:"radius_feet"`
	DurationMinutes int    `json:"duration_minutes"`
	Area            string `json:"area"`
}

// ItemWandOfWonderEnlarge is a nested part of a source payload.
type ItemWandOfWonderEnlarge struct {
	Effect   string `json:"effect"`
	Fallback string `json:"fallback"`
}

// ItemWandOfWonderCreature is a nested part of a source payload.
type ItemWandOfWonderCreature struct {
	Placement                 string            `json:"placement"`
	Control                   string            `json:"control"`
	Behavior                  string            `json:"behavior"`
	DurationHours             int               `json:"duration_hours"`
	AlsoDisappearsAtHitPoints int               `json:"also_disappears_at_hit_points"`
	Table                     ItemWandRollTable `json:"table"`
}

// ItemWandOfWonderGrass is a nested part of a source payload.
type ItemWandOfWonderGrass struct {
	RadiusFeet       int    `json:"radius_feet"`
	Center           string `json:"center"`
	SizeMultiplier   int    `json:"size_multiplier"`
	OvergrownMinutes int    `json:"overgrown_minutes"`
}

// ItemWandOfWonderEthereal is a nested part of a source payload.
type ItemWandOfWonderEthereal struct {
	Object               string `json:"object"`
	RangeFeet            int    `json:"range_feet"`
	MaximumDimensionFeet int    `json:"maximum_dimension_feet"`
	WithoutObject        string `json:"without_object"`
}

// ItemWandOfWonderShrink is a nested part of a source payload.
type ItemWandOfWonderShrink struct {
	Effect          string `json:"effect"`
	DurationMinutes int    `json:"duration_minutes"`
}

// ItemWandOfWonderLeaves is a nested part of a source payload.
type ItemWandOfWonderLeaves struct {
	FallOffAfterHours int    `json:"fall_off_after_hours"`
	EarlyRemoval      string `json:"early_removal"`
}

// ItemWandOfWonderLightBurst is a nested part of a source payload.
type ItemWandOfWonderLightBurst struct {
	EmanationFeet   int    `json:"emanation_feet"`
	SaveAbility     string `json:"save_ability"`
	SaveDC          int    `json:"save_dc"`
	Condition       string `json:"condition"`
	DurationMinutes int    `json:"duration_minutes"`
	RepeatSave      string `json:"repeat_save"`
}

// ItemWandOfWonderGems is a nested part of a source payload.
type ItemWandOfWonderGems struct {
	CountDice       ItemDice `json:"count_dice"`
	CountMultiplier int      `json:"count_multiplier"`
	ValueGPEach     int      `json:"value_gp_each"`
	LineLengthFeet  int      `json:"line_length_feet"`
	LineWidthFeet   int      `json:"line_width_feet"`
	Direction       string   `json:"direction"`
	DamagePerGem    int      `json:"damage_per_gem"`
	DamageType      string   `json:"damage_type"`
	DamageDivision  string   `json:"damage_division"`
}

// ItemWandOfWonderPetrify is a nested part of a source payload.
type ItemWandOfWonderPetrify struct {
	SaveAbility    string `json:"save_ability"`
	SaveDC         int    `json:"save_dc"`
	FirstFailure   string `json:"first_failure"`
	RepeatSave     string `json:"repeat_save"`
	RepeatSuccess  string `json:"repeat_success"`
	RepeatFailure  string `json:"repeat_failure"`
	PetrifiedUntil string `json:"petrified_until"`
}

// ItemWandOfWonderLastCharge is a nested part of a source payload.
type ItemWandOfWonderLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemWandOfWonderRecharge is a nested part of a source payload.
type ItemWandOfWonderRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfTheWarMage retains typed SRD source facts without implementing them.
type ItemWandOfTheWarMage struct {
	AttunementByAnyOf []string `json:"attunement_by_any_of"`
	Condition         string   `json:"condition"`
	SpellAttackBonus  int      `json:"spell_attack_bonus"`
	IgnoresHalfCover  bool     `json:"ignores_half_cover"`
}

// GemsCountMultiplierCheck keeps the gem stream's printed dice and multiplier together.
func (p *ItemWandOfWonder) GemsCountMultiplierCheck() bool {
	return p.Gems.CountDice == ItemDice{Dice: 1, Sides: 4, Type: "gems"} && p.Gems.CountMultiplier == 10
}

func itemWonderRows245(rows []ItemWandOfWonderEffect) [][2]int {
	out := make([][2]int, 0, len(rows))
	for _, r := range rows {
		out = append(out, [2]int{r.RollFrom, r.RollTo})
	}
	return out
}

func itemWonderEffects245(rows []ItemWandOfWonderEffect) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Effect)
	}
	return out
}

func itemSubRows245(rows []ItemWandRollRow) [][2]int {
	out := make([][2]int, 0, len(rows))
	for _, r := range rows {
		out = append(out, [2]int{r.RollFrom, r.RollTo})
	}
	return out
}

// itemManifest245 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest245 = itemSourceManifest{Batch: "data-magic-items-21", Issue: 245, Families: []itemSourceManifestFamily{
	{Slug: "wand-of-enemy-detection", Name: "Wand of Enemy Detection", Page: 250, Category: "wand", Attunement: "required", Equip: "hold", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement item.hold item.charges action.magic sense.direction-of-nearest-hostile-within-60-feet duration.1-minute item.destroyed-on-1-after-last-charge", Payload: "enemy_detection",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Enemy Detection", "rare"}}},
	{Slug: "wand-of-fear", Name: "Wand of Fear", Page: 250, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement item.hold item.charges spell.cast-listed-spells spell.save-dc-15 item.destroyed-on-1-after-last-charge", Payload: "fear",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Fear", "rare"}}},
	{Slug: "wand-of-fireballs", Name: "Wand of Fireballs", Page: 250, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-up-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement attunement.spellcaster item.hold item.charges spell.cast-fireball spell.save-dc-15 spell.level-3-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge", Payload: "fireballs",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Fireballs", "rare"}}},
	{Slug: "wand-of-lightning-bolts", Name: "Wand of Lightning Bolts", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-up-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement attunement.spellcaster item.hold item.charges spell.cast-lightning-bolt spell.save-dc-15 spell.level-3-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge", Payload: "lightning_bolts",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Lightning Bolts", "rare"}}},
	{Slug: "wand-of-magic-detection", Name: "Wand of Magic Detection", Page: 251, Category: "wand", Attunement: "none", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.hold item.charges spell.cast-detect-magic", Payload: "magic_detection",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Magic Detection", "uncommon"}}},
	{Slug: "wand-of-magic-missiles", Name: "Wand of Magic Missiles", Page: 251, Category: "wand", Attunement: "none", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-up-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.hold item.charges spell.cast-magic-missile spell.level-1-plus-1-per-extra-charge item.destroyed-on-1-after-last-charge", Payload: "magic_missiles",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Magic Missiles", "uncommon"}}},
	{Slug: "wand-of-paralysis", Name: "Wand of Paralysis", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "magic-action", UseRestrictions: "expend-one-charge target-within-60-feet",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement attunement.spellcaster item.hold item.charges action.magic-ray save.constitution-dc-15 condition.paralyzed duration.1-minute save.repeat-at-turn-end item.destroyed-on-1-after-last-charge", Payload: "paralysis",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Paralysis", "rare"}}},
	{Slug: "wand-of-polymorph", Name: "Wand of Polymorph", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement attunement.spellcaster item.hold item.charges spell.cast-polymorph spell.save-dc-15 item.destroyed-on-1-after-last-charge", Payload: "polymorph",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Polymorph", "very-rare"}}},
	{Slug: "wand-of-secrets", Name: "Wand of Secrets", Page: 251, Category: "wand", Attunement: "none", Equip: "hold", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.hold item.charges action.magic sense.nearest-secret-door-or-trap-within-60-feet", Payload: "secrets",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Secrets", "uncommon"}}},
	{Slug: "wand-of-web", Name: "Wand of Web", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement attunement.spellcaster item.hold item.charges spell.cast-web spell.save-dc-13 item.destroyed-on-1-after-last-charge", Payload: "web",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Web", "uncommon"}}},
	{Slug: "wand-of-wonder", Name: "Wand of Wonder", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "magic-action", UseRestrictions: "expend-one-charge point-within-120-feet",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement item.hold item.charges action.magic table.wand-of-wonder-effects spell.save-dc-15 spell.minimum-range-120-feet item.destroyed-on-1-after-last-charge", Payload: "wonder",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Wonder", "rare"}}},
	{Slug: "wand-of-the-war-mage-1-2-or-3", Name: "Wand of the War Mage, +1, +2, or +3", Page: 251, Category: "wand", Attunement: "required", Equip: "hold", Use: "none",
		Keys: "item.attunement attunement.spellcaster item.hold bonus.spell-attack-by-rarity attack.ignore-half-cover-on-spell-attacks", Payload: "war_mage",
		Variants: []itemSourceManifestVariant{{"plus-1", "Wand of the War Mage, +1", "uncommon"}, {"plus-2", "Wand of the War Mage, +2", "rare"}, {"plus-3", "Wand of the War Mage, +3", "very-rare"}}},
}}

// LoadItems245 returns a fresh validated copy of batch #245. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems245() (ItemSourceBatch[ItemDetails245], error) { return DecodeItems245(items245JSON) }

// DecodeItems245 validates supplied JSON as batch #245 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems245(data []byte) (ItemSourceBatch[ItemDetails245], error) {
	return decodeItemSourceBatch(data, itemManifest245, validateItemDetails245)
}

// validateItemDetails245 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails245(family, variant string, d *ItemDetails245) error {
	ok := true
	switch family {
	case "wand-of-enemy-detection":
		p := d.EnemyDetection
		ok = p.MaxCharges == 7 && p.RangeFeet == 60 && p.DurationMinutes == 1 && !p.RevealsDistance && len(p.Senses) == 5 && uniqueItemKeys(p.Senses) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-fear":
		p := d.Fear
		ok = p.MaxCharges == 7 && p.SaveDC == 15 && slices.Equal(p.Spells, []ItemWandOfFearSpell{{Spell: "command", Limit: "flee-or-grovel-only", ChargeCost: 1}, {Spell: "fear", Limit: "60-foot-cone", ChargeCost: 3}}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-fireballs":
		p := d.Fireballs
		ok = p.MaxCharges == 7 && p.MaximumChargesPerCast == 3 && p.Spell == "fireball" && p.SaveDC == 15 && p.BaseSpellLevel == 3 && p.LevelsPerAdditionalCharge == 1 && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20 && slices.Equal(p.SourcePages, []int{250, 251})
	case "wand-of-lightning-bolts":
		p := d.LightningBolts
		ok = p.MaxCharges == 7 && p.MaximumChargesPerCast == 3 && p.Spell == "lightning-bolt" && p.SaveDC == 15 && p.BaseSpellLevel == 3 && p.LevelsPerAdditionalCharge == 1 && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-magic-detection":
		p := d.MagicDetection
		ok = p.MaxCharges == 3 && p.ChargeCost == 1 && p.Spell == "detect-magic" && p.RechargeDice == ItemDice{Dice: 1, Sides: 3, Type: "charges"}
	case "wand-of-magic-missiles":
		p := d.MagicMissiles
		ok = p.MaxCharges == 7 && p.MaximumChargesPerCast == 3 && p.BaseSpellLevel == 1 && p.Spell == "magic-missile" && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-paralysis":
		p := d.Paralysis
		ok = p.MaxCharges == 7 && p.RangeFeet == 60 && p.SaveDC == 15 && p.Condition == "paralyzed" && p.DurationMinutes == 1 && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-polymorph":
		p := d.Polymorph
		ok = p.MaxCharges == 7 && p.ChargeCost == 1 && p.Spell == "polymorph" && p.SaveDC == 15 && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-secrets":
		p := d.Secrets
		ok = p.MaxCharges == 3 && p.RangeFeet == 60 && slices.Equal(p.Detects, []string{"secret-door", "trap"}) && p.RechargeDice == ItemDice{Dice: 1, Sides: 3, Type: "charges"}
	case "wand-of-web":
		p := d.Web
		ok = p.MaxCharges == 7 && p.ChargeCost == 1 && p.Spell == "web" && p.SaveDC == 13 && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20
	case "wand-of-wonder":
		p := d.Wonder
		ok = p.MaxCharges == 7 && p.SpellSaveDC == 15 && p.TableDieSides == 100 && len(p.Effects) == 18 && itemRollTable(itemWonderRows245(p.Effects), 100) &&
			itemRollTable(itemSubRows245(p.SpellTable.Rows), p.SpellTable.DieSides) && itemRollTable(itemSubRows245(p.Creature.Table.Rows), p.Creature.Table.DieSides) && itemRollTable(itemSubRows245(p.PolymorphTable.Rows), p.PolymorphTable.DieSides) &&
			p.SpellTable.DieSides == 10 && p.Creature.Table.DieSides == 4 && p.PolymorphTable.DieSides == 4 && uniqueItemKeys(itemWonderEffects245(p.Effects)) &&
			p.LightBurst.SaveDC == 15 && p.Petrify.SaveDC == 15 && p.Butterflies.Count == 600 && p.GemsCountMultiplierCheck() && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1 && p.LastCharge.DieSides == 20 && slices.Equal(p.SourcePages, []int{251, 252})
	case "wand-of-the-war-mage-1-2-or-3":
		p := d.WarMage
		ok = p.SpellAttackBonus == map[string]int{"plus-1": 1, "plus-2": 2, "plus-3": 3}[variant] && p.IgnoresHalfCover && slices.Equal(p.AttunementByAnyOf, []string{"spellcaster"})
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
