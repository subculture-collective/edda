package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_239.json
var items239JSON []byte

// ItemDetails239 holds exactly one typed payload per variant of catalogue batch
// #239 (magic items 15, SRD 5.2.1 pp. 238-239). Payloads are source facts only.
type ItemDetails239 struct {
	FreeAction     *ItemRingOfFreeAction    `json:"free_action,omitempty"`
	Invisibility   *ItemRingOfInvisibility  `json:"invisibility,omitempty"`
	Jumping        *ItemRingOfJumping       `json:"jumping,omitempty"`
	MindShielding  *ItemRingOfMindShielding `json:"mind_shielding,omitempty"`
	Protection     *ItemRingOfProtection    `json:"protection,omitempty"`
	Regeneration   *ItemRingOfRegeneration  `json:"regeneration,omitempty"`
	ResistanceRing *ItemRingOfResistance    `json:"resistance_ring,omitempty"`
	ShootingStars  *ItemRingOfShootingStars `json:"shooting_stars,omitempty"`
	SpellStoring   *ItemRingOfSpellStoring  `json:"spell_storing,omitempty"`
	SpellTurning   *ItemRingOfSpellTurning  `json:"spell_turning,omitempty"`
	Swimming       *ItemRingOfSwimming      `json:"swimming,omitempty"`
	Telekinesis    *ItemRingOfTelekinesis   `json:"telekinesis,omitempty"`
}

// ItemRingOfFreeAction retains typed SRD source facts without implementing them.
type ItemRingOfFreeAction struct {
	Condition                          string   `json:"condition"`
	DifficultTerrainCostsExtraMovement bool     `json:"difficult_terrain_costs_extra_movement"`
	MagicCannot                        []string `json:"magic_cannot"`
}

// ItemRingOfInvisibility retains typed SRD source facts without implementing them.
type ItemRingOfInvisibility struct {
	Condition       string   `json:"condition"`
	Action          string   `json:"action"`
	GrantsCondition string   `json:"grants_condition"`
	EndsWhen        []string `json:"ends_when"`
}

// ItemRingOfJumping retains typed SRD source facts without implementing them.
type ItemRingOfJumping struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
	Target    string `json:"target"`
}

// ItemRingOfMindShielding retains typed SRD source facts without implementing them.
type ItemRingOfMindShielding struct {
	Condition                string                               `json:"condition"`
	ImmuneToMagicThat        []string                             `json:"immune_to_magic_that"`
	TelepathyNeedsPermission bool                                 `json:"telepathy_needs_permission"`
	Imperceptible            ItemRingOfMindShieldingImperceptible `json:"imperceptible"`
	Soul                     ItemRingOfMindShieldingSoul          `json:"soul"`
}

// ItemRingOfMindShieldingImperceptible is a nested part of a source payload.
type ItemRingOfMindShieldingImperceptible struct {
	Action   string   `json:"action"`
	EndsWhen []string `json:"ends_when"`
}

// ItemRingOfMindShieldingSoul is a nested part of a source payload.
type ItemRingOfMindShieldingSoul struct {
	Trigger          string   `json:"trigger"`
	BlockedIf        string   `json:"blocked_if"`
	Choices          []string `json:"choices"`
	TelepathyWhile   string   `json:"telepathy_while"`
	TelepathyWith    string   `json:"telepathy_with"`
	WearerCanPrevent bool     `json:"wearer_can_prevent"`
}

// ItemRingOfProtection retains typed SRD source facts without implementing them.
type ItemRingOfProtection struct {
	Condition        string `json:"condition"`
	ArmorClassBonus  int    `json:"armor_class_bonus"`
	SavingThrowBonus int    `json:"saving_throw_bonus"`
}

// ItemRingOfRegeneration retains typed SRD source facts without implementing them.
type ItemRingOfRegeneration struct {
	Condition        string                         `json:"condition"`
	Healing          ItemDice                       `json:"healing"`
	IntervalMinutes  int                            `json:"interval_minutes"`
	MinimumHitPoints int                            `json:"minimum_hit_points"`
	Regrowth         ItemRingOfRegenerationRegrowth `json:"regrowth"`
}

// ItemRingOfRegenerationRegrowth is a nested part of a source payload.
type ItemRingOfRegenerationRegrowth struct {
	Trigger                    string   `json:"trigger"`
	Result                     string   `json:"result"`
	Delay                      ItemDice `json:"delay"`
	DelayBonusDays             int      `json:"delay_bonus_days"`
	RequiresHitPointsWholeTime bool     `json:"requires_hit_points_whole_time"`
}

// ItemRingOfResistance retains typed SRD source facts without implementing them.
type ItemRingOfResistance struct {
	Condition     string                           `json:"condition"`
	Defense       string                           `json:"defense"`
	TypeShownBy   string                           `json:"type_shown_by"`
	TypeChosenBy  string                           `json:"type_chosen_by"`
	TableDieSides int                              `json:"table_die_sides"`
	DamageTypes   []ItemRingOfResistanceDamageType `json:"damage_types"`
}

// ItemRingOfResistanceDamageType is a nested part of a source payload.
type ItemRingOfResistanceDamageType struct {
	Roll       int    `json:"roll"`
	DamageType string `json:"damage_type"`
	Gemstone   string `json:"gemstone"`
}

// ItemRingOfShootingStars retains typed SRD source facts without implementing them.
type ItemRingOfShootingStars struct {
	FreeSpells       []string                                `json:"free_spells"`
	MaxCharges       int                                     `json:"max_charges"`
	Recharge         ItemDice                                `json:"recharge"`
	RechargeWhen     string                                  `json:"recharge_when"`
	FaerieFire       ItemRingOfShootingStarsFaerieFire       `json:"faerie_fire"`
	LightningSpheres ItemRingOfShootingStarsLightningSpheres `json:"lightning_spheres"`
	ShootingStars    ItemRingOfShootingStarsShootingStars    `json:"shooting_stars"`
	SourcePages      []int                                   `json:"source_pages"`
}

// ItemRingOfShootingStarsFaerieFire is a nested part of a source payload.
type ItemRingOfShootingStarsFaerieFire struct {
	ChargeCost int    `json:"charge_cost"`
	Spell      string `json:"spell"`
}

// ItemRingOfShootingStarsLightningSpheres is a nested part of a source payload.
type ItemRingOfShootingStarsLightningSpheres struct {
	ChargeCost                  int                                                  `json:"charge_cost"`
	Action                      string                                               `json:"action"`
	MaxSpheres                  int                                                  `json:"max_spheres"`
	SphereDiameterFeet          int                                                  `json:"sphere_diameter_feet"`
	PlacementRangeFeet          int                                                  `json:"placement_range_feet"`
	ConcentrationRequired       bool                                                 `json:"concentration_required"`
	MaxDurationMinutes          int                                                  `json:"max_duration_minutes"`
	DimLightRadiusFeet          int                                                  `json:"dim_light_radius_feet"`
	MoveAction                  string                                               `json:"move_action"`
	MoveFeet                    int                                                  `json:"move_feet"`
	LeashFeet                   int                                                  `json:"leash_feet"`
	DischargeDistanceFeet       int                                                  `json:"discharge_distance_feet"`
	DischargeExcludes           []string                                             `json:"discharge_excludes"`
	SphereDisappearsOnDischarge bool                                                 `json:"sphere_disappears_on_discharge"`
	SaveAbility                 string                                               `json:"save_ability"`
	SaveDC                      int                                                  `json:"save_dc"`
	SuccessDamage               string                                               `json:"success_damage"`
	DamageTable                 []ItemRingOfShootingStarsLightningSpheresDamageTable `json:"damage_table"`
}

// ItemRingOfShootingStarsLightningSpheresDamageTable is a nested part of a source payload.
type ItemRingOfShootingStarsLightningSpheresDamageTable struct {
	Spheres int      `json:"spheres"`
	Damage  ItemDice `json:"damage"`
}

// ItemRingOfShootingStarsShootingStars is a nested part of a source payload.
type ItemRingOfShootingStarsShootingStars struct {
	MinimumCharges int      `json:"minimum_charges"`
	MaximumCharges int      `json:"maximum_charges"`
	Action         string   `json:"action"`
	MotesPerCharge int      `json:"motes_per_charge"`
	RangeFeet      int      `json:"range_feet"`
	AreaShape      string   `json:"area_shape"`
	AreaFeet       int      `json:"area_feet"`
	SaveAbility    string   `json:"save_ability"`
	SaveDC         int      `json:"save_dc"`
	Damage         ItemDice `json:"damage"`
	SuccessDamage  string   `json:"success_damage"`
}

// ItemRingOfSpellStoring retains typed SRD source facts without implementing them.
type ItemRingOfSpellStoring struct {
	Holds               string                        `json:"holds"`
	CapacitySpellLevels int                           `json:"capacity_spell_levels"`
	FoundLevels         ItemDice                      `json:"found_levels"`
	FoundLevelsPenalty  int                           `json:"found_levels_penalty"`
	FoundSpellsChosenBy string                        `json:"found_spells_chosen_by"`
	Storing             ItemRingOfSpellStoringStoring `json:"storing"`
	Casting             ItemRingOfSpellStoringCasting `json:"casting"`
}

// ItemRingOfSpellStoringStoring is a nested part of a source payload.
type ItemRingOfSpellStoringStoring struct {
	Caster            string `json:"caster"`
	MinimumSpellLevel int    `json:"minimum_spell_level"`
	MaximumSpellLevel int    `json:"maximum_spell_level"`
	Method            string `json:"method"`
	OtherEffect       string `json:"other_effect"`
	IfItCannotFit     string `json:"if_it_cannot_fit"`
	SpaceUsed         string `json:"space_used"`
}

// ItemRingOfSpellStoringCasting is a nested part of a source payload.
type ItemRingOfSpellStoringCasting struct {
	Condition              string   `json:"condition"`
	UsesFromOriginalCaster []string `json:"uses_from_original_caster"`
	Otherwise              string   `json:"otherwise"`
	FreesSpace             bool     `json:"frees_space"`
}

// ItemRingOfSpellTurning retains typed SRD source facts without implementing them.
type ItemRingOfSpellTurning struct {
	Condition           string                        `json:"condition"`
	SaveMode            string                        `json:"save_mode"`
	NegateMaxSpellLevel int                           `json:"negate_max_spell_level"`
	NegateResult        string                        `json:"negate_result"`
	Deflect             ItemRingOfSpellTurningDeflect `json:"deflect"`
}

// ItemRingOfSpellTurningDeflect is a nested part of a source payload.
type ItemRingOfSpellTurningDeflect struct {
	Requires       []string `json:"requires"`
	Action         string   `json:"action"`
	Target         string   `json:"target"`
	CasterSaveUses string   `json:"caster_save_uses"`
}

// ItemRingOfSwimming retains typed SRD source facts without implementing them.
type ItemRingOfSwimming struct {
	Condition     string `json:"condition"`
	SwimSpeedFeet int    `json:"swim_speed_feet"`
}

// ItemRingOfTelekinesis retains typed SRD source facts without implementing them.
type ItemRingOfTelekinesis struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
}

// itemManifest239 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest239 = itemSourceManifest{Batch: "data-magic-items-15", Issue: 239, Families: []itemSourceManifestFamily{
	{Slug: "ring-of-free-action", Name: "Ring of Free Action", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear movement.ignore-difficult-terrain-cost defense.magic-cannot-reduce-speed defense.magic-cannot-paralyze-or-restrain", Payload: "free_action",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Free Action", "rare"}}},
	{Slug: "ring-of-invisibility", Name: "Ring of Invisibility", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "magic-action",
		Keys: "item.attunement item.wear action.magic condition.invisible effect.ends-on-removal-or-bonus-action", Payload: "invisibility",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Invisibility", "legendary"}}},
	{Slug: "ring-of-jumping", Name: "Ring of Jumping", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "jump-targeting-yourself",
		Keys: "item.attunement item.wear spell.cast-jump spell.target-self-only", Payload: "jumping",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Jumping", "uncommon"}}},
	{Slug: "ring-of-mind-shielding", Name: "Ring of Mind Shielding", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.immunity-mind-reading-magic communication.telepathy-only-if-allowed item.imperceptible-toggle item.soul-vessel-on-death communication.soul-telepathy", Payload: "mind_shielding",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Mind Shielding", "uncommon"}}},
	{Slug: "ring-of-protection", Name: "Ring of Protection", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.armor-class-bonus-1 save.all-bonus-1", Payload: "protection",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Protection", "rare"}}},
	{Slug: "ring-of-regeneration", Name: "Ring of Regeneration", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear healing.1d6-every-10-minutes healing.regrow-body-part healing.requires-1-hit-point", Payload: "regeneration",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Regeneration", "very-rare"}}},
	{Slug: "ring-of-resistance", Name: "Ring of Resistance", Page: 238, Category: "ring", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear defense.resistance-one-damage-type table.damage-type-and-gemstone-1d10", Payload: "resistance_ring",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Resistance", "rare"}}},
	{Slug: "ring-of-shooting-stars", Name: "Ring of Shooting Stars", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "regain-1d6-at-dawn"},
		Keys:    "item.attunement item.wear item.charges spell.cast-dancing-lights spell.cast-light spell.cast-faerie-fire action.magic-lightning-spheres save.dexterity-dc-15 table.sphere-damage action.magic-shooting-stars damage.radiant-5d4-per-charge", Payload: "shooting_stars",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Shooting Stars", "very-rare"}}},
	{Slug: "ring-of-spell-storing", Name: "Ring of Spell Storing", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "spell-stored-in-ring",
		Keys: "item.attunement item.wear item.store-spells item.capacity-5-spell-levels spell.cast-stored-spell spell.uses-original-caster-statistics", Payload: "spell_storing",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Spell Storing", "rare"}}},
	{Slug: "ring-of-spell-turning", Name: "Ring of Spell Turning", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "reaction", UseRestrictions: "successful-save-against-qualifying-spell",
		Keys: "item.attunement item.wear save.advantage-against-spells spell.negate-level-7-or-lower-on-successful-save action.reaction-deflect-spell", Payload: "spell_turning",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Spell Turning", "legendary"}}},
	{Slug: "ring-of-swimming", Name: "Ring of Swimming", Page: 239, Category: "ring", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear movement.swim-speed-40-feet", Payload: "swimming",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Swimming", "uncommon"}}},
	{Slug: "ring-of-telekinesis", Name: "Ring of Telekinesis", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "telekinesis",
		Keys: "item.attunement item.wear spell.cast-telekinesis", Payload: "telekinesis",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Telekinesis", "very-rare"}}},
}}

// LoadItems239 returns a fresh validated copy of batch #239. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems239() (ItemSourceBatch[ItemDetails239], error) { return DecodeItems239(items239JSON) }

// DecodeItems239 validates supplied JSON as batch #239 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems239(data []byte) (ItemSourceBatch[ItemDetails239], error) {
	return decodeItemSourceBatch(data, itemManifest239, validateItemDetails239)
}

// validateItemDetails239 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails239(family, variant string, d *ItemDetails239) error {
	ok := true
	switch family {
	case "ring-of-free-action":
		p := d.FreeAction
		ok = !p.DifficultTerrainCostsExtraMovement && uniqueItemKeys(p.MagicCannot) && len(p.MagicCannot) == 3
	case "ring-of-invisibility":
		p := d.Invisibility
		ok = p.Action == "magic" && p.GrantsCondition == "invisible" && uniqueItemKeys(p.EndsWhen) && len(p.EndsWhen) == 2
	case "ring-of-jumping":
		p := d.Jumping
		ok = p.Spell == "jump" && p.Target == "yourself-only"
	case "ring-of-mind-shielding":
		p := d.MindShielding
		ok = uniqueItemKeys(p.ImmuneToMagicThat) && len(p.ImmuneToMagicThat) == 4 && p.TelepathyNeedsPermission && uniqueItemKeys(p.Imperceptible.EndsWhen) && uniqueItemKeys(p.Soul.Choices) && !p.Soul.WearerCanPrevent
	case "ring-of-protection":
		p := d.Protection
		ok = p.ArmorClassBonus == 1 && p.SavingThrowBonus == 1
	case "ring-of-regeneration":
		p := d.Regeneration
		ok = p.Healing == ItemDice{Dice: 1, Sides: 6, Type: "hit-points"} && p.IntervalMinutes == 10 && p.MinimumHitPoints == 1 && p.Regrowth.DelayBonusDays == 1 && p.Regrowth.RequiresHitPointsWholeTime
	case "ring-of-resistance":
		p := d.ResistanceRing
		rows := [][2]int{}
		types, gems := []string{}, []string{}
		for _, r := range p.DamageTypes {
			rows = append(rows, [2]int{r.Roll, r.Roll})
			types, gems = append(types, r.DamageType), append(gems, r.Gemstone)
		}
		ok = p.TableDieSides == 10 && itemRollTable(rows, 10) && uniqueItemKeys(types) && uniqueItemKeys(gems)
	case "ring-of-shooting-stars":
		p := d.ShootingStars
		rows := [][2]int{}
		for _, r := range p.LightningSpheres.DamageTable {
			rows = append(rows, [2]int{r.Spheres, r.Spheres})
		}
		l, s := p.LightningSpheres, p.ShootingStars
		ok = p.MaxCharges == 6 && p.Recharge == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && uniqueItemKeys(p.FreeSpells) && p.FaerieFire.ChargeCost == 1 && l.ChargeCost == 2 && l.SaveDC == 15 && l.MaxSpheres == 4 && itemRollTable(rows, 4) &&
			uniqueItemKeys(l.DischargeExcludes) && s.SaveDC == 15 && s.Damage == ItemDice{Dice: 5, Sides: 4, Type: "radiant"} && s.MinimumCharges == 1 && s.MaximumCharges == 3 && slices.Equal(p.SourcePages, []int{238, 239})
	case "ring-of-spell-storing":
		p := d.SpellStoring
		s := p.Storing
		ok = p.CapacitySpellLevels == 5 && p.FoundLevels == ItemDice{Dice: 1, Sides: 6, Type: "stored-spell-levels"} && p.FoundLevelsPenalty == 1 && s.MinimumSpellLevel == 1 && s.MaximumSpellLevel == 5 && uniqueItemKeys(p.Casting.UsesFromOriginalCaster) && p.Casting.FreesSpace
	case "ring-of-spell-turning":
		p := d.SpellTurning
		ok = p.SaveMode == "advantage" && p.NegateMaxSpellLevel == 7 && p.Deflect.Action == "reaction" && uniqueItemKeys(p.Deflect.Requires)
	case "ring-of-swimming":
		p := d.Swimming
		ok = p.SwimSpeedFeet == 40
	case "ring-of-telekinesis":
		p := d.Telekinesis
		ok = p.Spell == "telekinesis"
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
