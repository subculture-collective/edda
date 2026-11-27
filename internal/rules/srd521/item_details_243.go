package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_243.json
var items243JSON []byte

// ItemDetails243 holds exactly one typed payload per variant of catalogue batch
// #243 (magic items 19, SRD 5.2.1 pp. 246-248). Payloads are source facts only.
type ItemDetails243 struct {
	SwarmingInsects  *ItemStaffOfSwarmingInsects            `json:"swarming_insects,omitempty"`
	ThunderLightning *ItemStaffOfThunderAndLightning        `json:"thunder_lightning,omitempty"`
	Withering        *ItemStaffOfWithering                  `json:"withering,omitempty"`
	Magi             *ItemStaffOfTheMagi                    `json:"magi,omitempty"`
	Python           *ItemStaffOfThePython                  `json:"python,omitempty"`
	Woodlands        *ItemStaffOfTheWoodlands               `json:"woodlands,omitempty"`
	EarthStone       *ItemStoneOfControllingEarthElementals `json:"earth_stone,omitempty"`
	Luckstone        *ItemStoneOfGoodLuckLuckstone          `json:"luckstone,omitempty"`
	SunBlade         *ItemSunBlade                          `json:"sun_blade,omitempty"`
	LifeStealing     *ItemSwordOfLifeStealing               `json:"life_stealing,omitempty"`
	Sharpness        *ItemSwordOfSharpness                  `json:"sharpness,omitempty"`
	Wounding         *ItemSwordOfWounding                   `json:"wounding,omitempty"`
}

// ItemStaffOfSwarmingInsects retains typed SRD source facts without implementing them.
type ItemStaffOfSwarmingInsects struct {
	AttunementByAnyOf   []string                              `json:"attunement_by_any_of"`
	MaxCharges          int                                   `json:"max_charges"`
	InsectCloud         ItemStaffOfSwarmingInsectsInsectCloud `json:"insect_cloud"`
	SaveAndAttackSource string                                `json:"save_and_attack_source"`
	Spells              []ItemStaffOfSwarmingInsectsSpell     `json:"spells"`
	Recharge            ItemStaffOfSwarmingInsectsRecharge    `json:"recharge"`
	LastCharge          ItemStaffOfSwarmingInsectsLastCharge  `json:"last_charge"`
}

// ItemStaffOfSwarmingInsectsInsectCloud is a nested part of a source payload.
type ItemStaffOfSwarmingInsectsInsectCloud struct {
	Action          string `json:"action"`
	ChargeCost      int    `json:"charge_cost"`
	Swarm           string `json:"swarm"`
	EmanationFeet   int    `json:"emanation_feet"`
	DurationMinutes int    `json:"duration_minutes"`
	Area            string `json:"area"`
	DispersedBy     string `json:"dispersed_by"`
}

// ItemStaffOfSwarmingInsectsSpell is a nested part of a source payload.
type ItemStaffOfSwarmingInsectsSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemStaffOfSwarmingInsectsRecharge is a nested part of a source payload.
type ItemStaffOfSwarmingInsectsRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfSwarmingInsectsLastCharge is a nested part of a source payload.
type ItemStaffOfSwarmingInsectsLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStaffOfThunderAndLightning retains typed SRD source facts without implementing them.
type ItemStaffOfThunderAndLightning struct {
	WieldedAs              string                                            `json:"wielded_as"`
	AttackBonus            int                                               `json:"attack_bonus"`
	DamageBonus            int                                               `json:"damage_bonus"`
	EachPropertyReuseAfter string                                            `json:"each_property_reuse_after"`
	Lightning              ItemStaffOfThunderAndLightningLightning           `json:"lightning"`
	Thunder                ItemStaffOfThunderAndLightningThunder             `json:"thunder"`
	ThunderAndLightning    ItemStaffOfThunderAndLightningThunderAndLightning `json:"thunder_and_lightning"`
	LightningStrike        ItemStaffOfThunderAndLightningLightningStrike     `json:"lightning_strike"`
	Thunderclap            ItemStaffOfThunderAndLightningThunderclap         `json:"thunderclap"`
	SourcePages            []int                                             `json:"source_pages"`
}

// ItemStaffOfThunderAndLightningLightning is a nested part of a source payload.
type ItemStaffOfThunderAndLightningLightning struct {
	Trigger        string   `json:"trigger"`
	ExtraDamage    ItemDice `json:"extra_damage"`
	ActionRequired bool     `json:"action_required"`
}

// ItemStaffOfThunderAndLightningThunder is a nested part of a source payload.
type ItemStaffOfThunderAndLightningThunder struct {
	Trigger        string `json:"trigger"`
	AudibleFeet    int    `json:"audible_feet"`
	ActionRequired bool   `json:"action_required"`
	SaveAbility    string `json:"save_ability"`
	SaveDC         int    `json:"save_dc"`
	Condition      string `json:"condition"`
	ConditionEnds  string `json:"condition_ends"`
}

// ItemStaffOfThunderAndLightningThunderAndLightning is a nested part of a source payload.
type ItemStaffOfThunderAndLightningThunderAndLightning struct {
	Trigger              string   `json:"trigger"`
	Action               string   `json:"action"`
	Uses                 []string `json:"uses"`
	SpendsTheirDailyUses bool     `json:"spends_their_daily_uses"`
}

// ItemStaffOfThunderAndLightningLightningStrike is a nested part of a source payload.
type ItemStaffOfThunderAndLightningLightningStrike struct {
	Action         string   `json:"action"`
	LineWidthFeet  int      `json:"line_width_feet"`
	LineLengthFeet int      `json:"line_length_feet"`
	SaveAbility    string   `json:"save_ability"`
	SaveDC         int      `json:"save_dc"`
	Damage         ItemDice `json:"damage"`
	SuccessDamage  string   `json:"success_damage"`
}

// ItemStaffOfThunderAndLightningThunderclap is a nested part of a source payload.
type ItemStaffOfThunderAndLightningThunderclap struct {
	Action           string   `json:"action"`
	AudibleFeet      int      `json:"audible_feet"`
	EmanationFeet    int      `json:"emanation_feet"`
	SaveAbility      string   `json:"save_ability"`
	SaveDC           int      `json:"save_dc"`
	Damage           ItemDice `json:"damage"`
	FailureCondition string   `json:"failure_condition"`
	ConditionMinutes int      `json:"condition_minutes"`
	Success          string   `json:"success"`
}

// ItemStaffOfWithering retains typed SRD source facts without implementing them.
type ItemStaffOfWithering struct {
	WieldedAs    string                      `json:"wielded_as"`
	HitDamage    string                      `json:"hit_damage"`
	MaxCharges   int                         `json:"max_charges"`
	RechargeDice ItemDice                    `json:"recharge_dice"`
	ChargeCost   int                         `json:"charge_cost"`
	ExtraDamage  ItemDice                    `json:"extra_damage"`
	SaveAbility  string                      `json:"save_ability"`
	SaveDC       int                         `json:"save_dc"`
	Failure      ItemStaffOfWitheringFailure `json:"failure"`
}

// ItemStaffOfWitheringFailure is a nested part of a source payload.
type ItemStaffOfWitheringFailure struct {
	Mode          string   `json:"mode"`
	DurationHours int      `json:"duration_hours"`
	AppliesTo     []string `json:"applies_to"`
}

// ItemStaffOfTheMagi retains typed SRD source facts without implementing them.
type ItemStaffOfTheMagi struct {
	AttunementByAnyOf []string                            `json:"attunement_by_any_of"`
	MaxCharges        int                                 `json:"max_charges"`
	WieldedAs         string                              `json:"wielded_as"`
	AttackBonus       int                                 `json:"attack_bonus"`
	DamageBonus       int                                 `json:"damage_bonus"`
	SpellAttackBonus  int                                 `json:"spell_attack_bonus"`
	SpellAbsorption   ItemStaffOfTheMagiSpellAbsorption   `json:"spell_absorption"`
	SaveDCSource      string                              `json:"save_dc_source"`
	Spells            []ItemStaffOfTheMagiSpell           `json:"spells"`
	Recharge          ItemStaffOfTheMagiRecharge          `json:"recharge"`
	LastCharge        ItemStaffOfTheMagiLastCharge        `json:"last_charge"`
	RetributiveStrike ItemStaffOfTheMagiRetributiveStrike `json:"retributive_strike"`
	SourcePages       []int                               `json:"source_pages"`
}

// ItemStaffOfTheMagiSpellAbsorption is a nested part of a source payload.
type ItemStaffOfTheMagiSpellAbsorption struct {
	SaveMode             string `json:"save_mode"`
	Trigger              string `json:"trigger"`
	Action               string `json:"action"`
	Result               string `json:"result"`
	ChargesGained        string `json:"charges_gained"`
	ExplodesAboveCharges int    `json:"explodes_above_charges"`
	Explosion            string `json:"explosion"`
}

// ItemStaffOfTheMagiSpell is a nested part of a source payload.
type ItemStaffOfTheMagiSpell struct {
	Spell        string `json:"spell"`
	PrintedLevel string `json:"printed_level"`
	ChargeCost   int    `json:"charge_cost"`
}

// ItemStaffOfTheMagiRecharge is a nested part of a source payload.
type ItemStaffOfTheMagiRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfTheMagiLastCharge is a nested part of a source payload.
type ItemStaffOfTheMagiLastCharge struct {
	DieSides     int      `json:"die_sides"`
	RegainsOn    int      `json:"regains_on"`
	Regains      ItemDice `json:"regains"`
	RegainsBonus int      `json:"regains_bonus"`
}

// ItemStaffOfTheMagiRetributiveStrike is a nested part of a source payload.
type ItemStaffOfTheMagiRetributiveStrike struct {
	Action               string `json:"action"`
	AreaEmanationFeet    int    `json:"area_emanation_feet"`
	PlanarTravelPercent  int    `json:"planar_travel_percent"`
	YourDamagePerCharge  int    `json:"your_damage_per_charge"`
	OtherDamagePerCharge int    `json:"other_damage_per_charge"`
	DamageType           string `json:"damage_type"`
	SaveAbility          string `json:"save_ability"`
	SaveDC               int    `json:"save_dc"`
	SuccessDamage        string `json:"success_damage"`
}

// ItemStaffOfThePython retains typed SRD source facts without implementing them.
type ItemStaffOfThePython struct {
	Action                 string                      `json:"action"`
	ThrowRangeFeet         int                         `json:"throw_range_feet"`
	Creature               string                      `json:"creature"`
	Control                string                      `json:"control"`
	Initiative             string                      `json:"initiative"`
	Command                ItemStaffOfThePythonCommand `json:"command"`
	Revert                 ItemStaffOfThePythonRevert  `json:"revert"`
	AtZeroHitPoints        string                      `json:"at_zero_hit_points"`
	HitPointsOnEarlyRevert string                      `json:"hit_points_on_early_revert"`
}

// ItemStaffOfThePythonCommand is a nested part of a source payload.
type ItemStaffOfThePythonCommand struct {
	Timing             string   `json:"timing"`
	Method             string   `json:"method"`
	ActionRequired     bool     `json:"action_required"`
	RangeFeet          int      `json:"range_feet"`
	BlockedByCondition string   `json:"blocked_by_condition"`
	Options            []string `json:"options"`
	WithoutCommands    string   `json:"without_commands"`
}

// ItemStaffOfThePythonRevert is a nested part of a source payload.
type ItemStaffOfThePythonRevert struct {
	Action          string `json:"action"`
	Location        string `json:"location"`
	ReuseAfterHours int    `json:"reuse_after_hours"`
}

// ItemStaffOfTheWoodlands retains typed SRD source facts without implementing them.
type ItemStaffOfTheWoodlands struct {
	AttunementByAnyOf []string                          `json:"attunement_by_any_of"`
	MaxCharges        int                               `json:"max_charges"`
	WieldedAs         string                            `json:"wielded_as"`
	AttackBonus       int                               `json:"attack_bonus"`
	DamageBonus       int                               `json:"damage_bonus"`
	SpellAttackBonus  int                               `json:"spell_attack_bonus"`
	SaveDCSource      string                            `json:"save_dc_source"`
	Spells            []ItemStaffOfTheWoodlandsSpell    `json:"spells"`
	TreeForm          ItemStaffOfTheWoodlandsTreeForm   `json:"tree_form"`
	RechargeDice      ItemDice                          `json:"recharge_dice"`
	LastCharge        ItemStaffOfTheWoodlandsLastCharge `json:"last_charge"`
}

// ItemStaffOfTheWoodlandsSpell is a nested part of a source payload.
type ItemStaffOfTheWoodlandsSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemStaffOfTheWoodlandsTreeForm is a nested part of a source payload.
type ItemStaffOfTheWoodlandsTreeForm struct {
	Action            string `json:"action"`
	ChargeCost        int    `json:"charge_cost"`
	Placement         string `json:"placement"`
	TreeHeightFeet    int    `json:"tree_height_feet"`
	TrunkDiameterFeet int    `json:"trunk_diameter_feet"`
	BranchRadiusFeet  int    `json:"branch_radius_feet"`
	Aura              string `json:"aura"`
	Revert            string `json:"revert"`
	CreaturesInTree   string `json:"creatures_in_tree"`
}

// ItemStaffOfTheWoodlandsLastCharge is a nested part of a source payload.
type ItemStaffOfTheWoodlandsLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStoneOfControllingEarthElementals retains typed SRD source facts without implementing them.
type ItemStoneOfControllingEarthElementals struct {
	WeightPounds       int      `json:"weight_pounds"`
	Condition          string   `json:"condition"`
	Action             string   `json:"action"`
	Creature           string   `json:"creature"`
	SummonRangeFeet    int      `json:"summon_range_feet"`
	ObeysCommands      bool     `json:"obeys_commands"`
	Initiative         string   `json:"initiative"`
	DurationHours      int      `json:"duration_hours"`
	AlsoDisappearsWhen []string `json:"also_disappears_when"`
	ReuseAfter         string   `json:"reuse_after"`
}

// ItemStoneOfGoodLuckLuckstone retains typed SRD source facts without implementing them.
type ItemStoneOfGoodLuckLuckstone struct {
	Appearance        string `json:"appearance"`
	Condition         string `json:"condition"`
	AbilityCheckBonus int    `json:"ability_check_bonus"`
	SavingThrowBonus  int    `json:"saving_throw_bonus"`
}

// ItemSunBlade retains typed SRD source facts without implementing them.
type ItemSunBlade struct {
	Appearance         string               `json:"appearance"`
	Blade              ItemSunBladeBlade    `json:"blade"`
	AttackBonus        int                  `json:"attack_bonus"`
	DamageBonus        int                  `json:"damage_bonus"`
	DamageType         string               `json:"damage_type"`
	ReplacesDamageType string               `json:"replaces_damage_type"`
	UndeadExtraDamage  ItemDice             `json:"undead_extra_damage"`
	Sunlight           ItemSunBladeSunlight `json:"sunlight"`
}

// ItemSunBladeBlade is a nested part of a source payload.
type ItemSunBladeBlade struct {
	Condition                   string   `json:"condition"`
	Action                      string   `json:"action"`
	Options                     []string `json:"options"`
	FunctionsAs                 string   `json:"functions_as"`
	Property                    string   `json:"property"`
	ProficiencyIfProficientWith []string `json:"proficiency_if_proficient_with"`
}

// ItemSunBladeSunlight is a nested part of a source payload.
type ItemSunBladeSunlight struct {
	BrightRadiusFeet  int    `json:"bright_radius_feet"`
	DimAdditionalFeet int    `json:"dim_additional_feet"`
	LightIsSunlight   bool   `json:"light_is_sunlight"`
	AdjustAction      string `json:"adjust_action"`
	AdjustStepFeet    int    `json:"adjust_step_feet"`
	MaximumEachFeet   int    `json:"maximum_each_feet"`
	MinimumEachFeet   int    `json:"minimum_each_feet"`
}

// ItemSwordOfLifeStealing retains typed SRD source facts without implementing them.
type ItemSwordOfLifeStealing struct {
	TriggerNaturalRoll    int      `json:"trigger_natural_roll"`
	Target                string   `json:"target"`
	ExtraDamage           int      `json:"extra_damage"`
	DamageType            string   `json:"damage_type"`
	ExcludedCreatureTypes []string `json:"excluded_creature_types"`
	TemporaryHitPoints    string   `json:"temporary_hit_points"`
}

// ItemSwordOfSharpness retains typed SRD source facts without implementing them.
type ItemSwordOfSharpness struct {
	ObjectHit          string `json:"object_hit"`
	TriggerNaturalRoll int    `json:"trigger_natural_roll"`
	Target             string `json:"target"`
	ExtraDamage        int    `json:"extra_damage"`
	DamageType         string `json:"damage_type"`
	AppliesCondition   string `json:"applies_condition"`
	ExhaustionLevels   int    `json:"exhaustion_levels"`
}

// ItemSwordOfWounding retains typed SRD source facts without implementing them.
type ItemSwordOfWounding struct {
	Trigger       string   `json:"trigger"`
	ExtraDamage   ItemDice `json:"extra_damage"`
	SaveAbility   string   `json:"save_ability"`
	SaveDC        int      `json:"save_dc"`
	Failure       string   `json:"failure"`
	DurationHours int      `json:"duration_hours"`
	RepeatSave    string   `json:"repeat_save"`
	RepeatSuccess string   `json:"repeat_success"`
}

// itemManifest243 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest243 = itemSourceManifest{Batch: "data-magic-items-19", Issue: 243, Families: []itemSourceManifestFamily{
	{Slug: "staff-of-swarming-insects", Name: "Staff of Swarming Insects", Page: 246, Category: "staff", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.hold item.charges action.magic-insect-cloud effect.heavily-obscured spell.cast-giant-insect spell.cast-insect-plague spell.uses-your-save-dc-and-attack-modifier item.destroyed-on-1-after-last-charge", Payload: "swarming_insects",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Swarming Insects", "rare"}}},
	{Slug: "staff-of-thunder-and-lightning", Name: "Staff of Thunder and Lightning", Page: 247, Category: "staff", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "each-property-once-until-dawn",
		Keys: "item.attunement item.magic-weapon bonus.attack-and-damage-2 damage.extra-2d6-lightning save.constitution-dc-17 condition.stunned action.bonus-both-properties save.dexterity-dc-17 damage.lightning-strike-9d6 damage.thunderclap-2d6 condition.deafened item.each-property-once-until-dawn", Payload: "thunder_lightning",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Thunder and Lightning", "very-rare"}}},
	{Slug: "staff-of-withering", Name: "Staff of Withering", Page: 248, Category: "staff", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "expend-one-charge-on-hit",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.magic-weapon item.charges damage.extra-2d10-necrotic save.constitution-dc-15 effect.disadvantage-strength-and-constitution duration.1-hour", Payload: "withering",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Withering", "rare"}}},
	{Slug: "staff-of-the-magi", Name: "Staff of the Magi", Page: 246, Category: "staff", Attunement: "required", Equip: "wield", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 50, PerUse: 1, Recharge: "regain-4d6-plus-2-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.magic-weapon bonus.attack-and-damage-2 bonus.spell-attack-2 item.charges save.advantage-against-spells action.reaction-spell-absorption spell.cast-listed-spells spell.uses-your-save-dc item.retributive-strike damage.force-by-charges save.dexterity-dc-17", Payload: "magi",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of the Magi", "legendary"}}},
	{Slug: "staff-of-the-python", Name: "Staff of the Python", Page: 247, Category: "staff", Attunement: "required", Equip: "hold", Use: "magic-action", UseRestrictions: "throw-within-10-feet",
		Keys: "item.attunement action.magic-throw-staff creature.giant-constrictor-snake creature.under-your-control creature.initiative-after-owner action.bonus-revert item.1-hour-reuse item.destroyed-if-snake-dies", Payload: "python",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of the Python", "uncommon"}}},
	{Slug: "staff-of-the-woodlands", Name: "Staff of the Woodlands", Page: 247, Category: "staff", Attunement: "required", Equip: "wield", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "regain-1d6-at-dawn"},
		Keys:    "item.attunement attunement.druid item.magic-weapon bonus.attack-and-damage-2 bonus.spell-attack-2 item.charges spell.cast-listed-spells spell.uses-your-save-dc action.magic-tree-form item.nonmagical-on-1-after-last-charge", Payload: "woodlands",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of the Woodlands", "rare"}}},
	{Slug: "stone-of-controlling-earth-elementals", Name: "Stone of Controlling Earth Elementals", Page: 248, Category: "wondrous-item", Attunement: "none", Equip: "touch", Use: "magic-action", UseRestrictions: "touch-stone-to-ground",
		Keys: "action.magic creature.summon-earth-elemental creature.obeys-commands creature.initiative-after-owner duration.1-hour item.once-until-dawn", Payload: "earth_stone",
		Variants: []itemSourceManifestVariant{{"standard", "Stone of Controlling Earth Elementals", "rare"}}},
	{Slug: "stone-of-good-luck-luckstone", Name: "Stone of Good Luck (Luckstone)", Page: 248, Category: "wondrous-item", Attunement: "required", Equip: "carry", Use: "none",
		Keys: "item.attunement item.carry check.all-bonus-1 save.all-bonus-1", Payload: "luckstone",
		Variants: []itemSourceManifestVariant{{"standard", "Stone of Good Luck (Luckstone)", "uncommon"}}},
	{Slug: "sun-blade", Name: "Sun Blade", Page: 248, Category: "weapon", Options: "longsword", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon item.blade-of-radiance property.finesse bonus.attack-and-damage-2 damage.radiant-replaces-slashing damage.extra-1d8-radiant-against-undead light.sunlight", Payload: "sun_blade",
		Variants: []itemSourceManifestVariant{{"standard", "Sun Blade", "rare"}}},
	{Slug: "sword-of-life-stealing", Name: "Sword of Life Stealing", Page: 248, Category: "weapon", Options: "glaive greatsword longsword rapier scimitar shortsword", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon damage.extra-15-necrotic-on-natural-20 hit-points.temporary-equal-to-necrotic-damage", Payload: "life_stealing",
		Variants: []itemSourceManifestVariant{{"standard", "Sword of Life Stealing", "rare"}}},
	{Slug: "sword-of-sharpness", Name: "Sword of Sharpness", Page: 248, Category: "weapon", Options: "glaive greatsword longsword scimitar", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon damage.maximize-weapon-dice-against-objects damage.extra-14-slashing-on-natural-20 condition.exhaustion", Payload: "sharpness",
		Variants: []itemSourceManifestVariant{{"standard", "Sword of Sharpness", "very-rare"}}},
	{Slug: "sword-of-wounding", Name: "Sword of Wounding", Page: 248, Category: "weapon", Options: "glaive greatsword longsword rapier scimitar shortsword", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon damage.extra-2d6-necrotic save.constitution-dc-15 effect.cannot-regain-hit-points duration.1-hour save.repeat-at-turn-end", Payload: "wounding",
		Variants: []itemSourceManifestVariant{{"standard", "Sword of Wounding", "rare"}}},
}}

// LoadItems243 returns a fresh validated copy of batch #243. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems243() (ItemSourceBatch[ItemDetails243], error) { return DecodeItems243(items243JSON) }

// DecodeItems243 validates supplied JSON as batch #243 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems243(data []byte) (ItemSourceBatch[ItemDetails243], error) {
	return decodeItemSourceBatch(data, itemManifest243, validateItemDetails243)
}

// validateItemDetails243 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails243(family, variant string, d *ItemDetails243) error {
	ok := true
	switch family {
	case "staff-of-swarming-insects":
		p := d.SwarmingInsects
		ok = p.MaxCharges == 10 && p.InsectCloud.EmanationFeet == 30 && p.InsectCloud.DurationMinutes == 10 && slices.Equal(p.Spells, []ItemStaffOfSwarmingInsectsSpell{{Spell: "giant-insect", ChargeCost: 4}, {Spell: "insect-plague", ChargeCost: 5}}) && len(p.AttunementByAnyOf) == 6
	case "staff-of-thunder-and-lightning":
		p := d.ThunderLightning
		ok = p.AttackBonus == 2 && p.Thunder.SaveDC == 17 && p.LightningStrike.SaveDC == 17 && p.Thunderclap.SaveDC == 17 && p.Lightning.ExtraDamage == ItemDice{Dice: 2, Sides: 6, Type: "lightning"} && p.LightningStrike.Damage == ItemDice{Dice: 9, Sides: 6, Type: "lightning"} && p.Thunderclap.Damage == ItemDice{Dice: 2, Sides: 6, Type: "thunder"} && slices.Equal(p.SourcePages, []int{247, 248})
	case "staff-of-withering":
		p := d.Withering
		ok = p.MaxCharges == 3 && p.SaveDC == 15 && p.ExtraDamage == ItemDice{Dice: 2, Sides: 10, Type: "necrotic"} && p.RechargeDice == ItemDice{Dice: 1, Sides: 3, Type: "charges"} && p.Failure.DurationHours == 1
	case "staff-of-the-magi":
		p := d.Magi
		r := p.RetributiveStrike
		costs := []int{}
		names := []string{}
		for _, s := range p.Spells {
			costs, names = append(costs, s.ChargeCost), append(names, s.Spell)
		}
		ok = p.MaxCharges == 50 && p.AttackBonus == 2 && p.SpellAttackBonus == 2 && p.SpellAbsorption.ExplodesAboveCharges == 50 && slices.Equal(costs, []int{0, 7, 0, 3, 0, 7, 2, 4, 2, 2, 0, 7, 0, 5, 7, 0, 5, 4, 2}) && uniqueItemKeys(names) &&
			r.SaveDC == 17 && r.YourDamagePerCharge == 16 && r.OtherDamagePerCharge == 6 && p.Recharge.Dice == ItemDice{Dice: 4, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 2 && slices.Equal(p.SourcePages, []int{246, 247})
	case "staff-of-the-python":
		p := d.Python
		ok = p.ThrowRangeFeet == 10 && p.Command.RangeFeet == 60 && p.Revert.ReuseAfterHours == 1 && p.Creature == "giant-constrictor-snake"
	case "staff-of-the-woodlands":
		p := d.Woodlands
		ok = p.MaxCharges == 6 && p.AttackBonus == 2 && p.SpellAttackBonus == 2 && p.TreeForm.TreeHeightFeet == 60 && p.RechargeDice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && slices.Equal(p.Spells, []ItemStaffOfTheWoodlandsSpell{{Spell: "animal-friendship", ChargeCost: 1}, {Spell: "awaken", ChargeCost: 5}, {Spell: "barkskin", ChargeCost: 2}, {Spell: "locate-animals-or-plants", ChargeCost: 2}, {Spell: "pass-without-trace", ChargeCost: 2}, {Spell: "speak-with-animals", ChargeCost: 1}, {Spell: "speak-with-plants", ChargeCost: 3}, {Spell: "wall-of-thorns", ChargeCost: 6}}) && slices.Equal(p.AttunementByAnyOf, []string{"druid"})
	case "stone-of-controlling-earth-elementals":
		p := d.EarthStone
		ok = p.WeightPounds == 5 && p.SummonRangeFeet == 30 && p.DurationHours == 1 && p.Creature == "earth-elemental" && uniqueItemKeys(p.AlsoDisappearsWhen)
	case "stone-of-good-luck-luckstone":
		p := d.Luckstone
		ok = p.AbilityCheckBonus == 1 && p.SavingThrowBonus == 1
	case "sun-blade":
		p := d.SunBlade
		s := p.Sunlight
		ok = p.AttackBonus == 2 && p.DamageBonus == 2 && p.UndeadExtraDamage == ItemDice{Dice: 1, Sides: 8, Type: "radiant"} && s.BrightRadiusFeet == 15 && s.MaximumEachFeet == 30 && s.MinimumEachFeet == 10 && s.LightIsSunlight && uniqueItemKeys(p.Blade.Options)
	case "sword-of-life-stealing":
		p := d.LifeStealing
		ok = p.TriggerNaturalRoll == 20 && p.ExtraDamage == 15 && slices.Equal(p.ExcludedCreatureTypes, []string{"construct", "undead"})
	case "sword-of-sharpness":
		p := d.Sharpness
		ok = p.TriggerNaturalRoll == 20 && p.ExtraDamage == 14 && p.ExhaustionLevels == 1
	case "sword-of-wounding":
		p := d.Wounding
		ok = p.SaveDC == 15 && p.DurationHours == 1 && p.ExtraDamage == ItemDice{Dice: 2, Sides: 6, Type: "necrotic"}
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
