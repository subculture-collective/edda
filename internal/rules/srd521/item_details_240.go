package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_240.json
var items240JSON []byte

// ItemDetails240 holds exactly one typed payload per variant of catalogue batch
// #240 (magic items 16, SRD 5.2.1 pp. 239-241). Payloads are source facts only.
type ItemDetails240 struct {
	ThreeWishes   *ItemRingOfThreeWishes         `json:"three_wishes,omitempty"`
	Warmth        *ItemRingOfWarmth              `json:"warmth,omitempty"`
	WaterWalking  *ItemRingOfWaterWalking        `json:"water_walking,omitempty"`
	XRay          *ItemRingOfXRayVision          `json:"x_ray,omitempty"`
	Ram           *ItemRingOfTheRam              `json:"ram,omitempty"`
	Eyes          *ItemRobeOfEyes                `json:"eyes,omitempty"`
	Scintillating *ItemRobeOfScintillatingColors `json:"scintillating,omitempty"`
	Stars         *ItemRobeOfStars               `json:"stars,omitempty"`
	UsefulItems   *ItemRobeOfUsefulItems         `json:"useful_items,omitempty"`
	Archmagi      *ItemRobeOfTheArchmagi         `json:"archmagi,omitempty"`
	Absorption    *ItemRodOfAbsorption           `json:"absorption,omitempty"`
	Alertness     *ItemRodOfAlertness            `json:"alertness,omitempty"`
}

// ItemRingOfThreeWishes retains typed SRD source facts without implementing them.
type ItemRingOfThreeWishes struct {
	Condition       string `json:"condition"`
	Spell           string `json:"spell"`
	MaxCharges      int    `json:"max_charges"`
	ChargeCost      int    `json:"charge_cost"`
	AfterLastCharge string `json:"after_last_charge"`
}

// ItemRingOfWarmth retains typed SRD source facts without implementing them.
type ItemRingOfWarmth struct {
	Trigger                      string   `json:"trigger"`
	DamageType                   string   `json:"damage_type"`
	Reduction                    ItemDice `json:"reduction"`
	Protected                    []string `json:"protected"`
	ProtectedAtOrBelowFahrenheit int      `json:"protected_at_or_below_fahrenheit"`
}

// ItemRingOfWaterWalking retains typed SRD source facts without implementing them.
type ItemRingOfWaterWalking struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
	Target    string `json:"target"`
}

// ItemRingOfXRayVision retains typed SRD source facts without implementing them.
type ItemRingOfXRayVision struct {
	Condition       string                          `json:"condition"`
	Action          string                          `json:"action"`
	RangeFeet       int                             `json:"range_feet"`
	DurationMinutes int                             `json:"duration_minutes"`
	SolidObjects    string                          `json:"solid_objects"`
	Penetration     ItemRingOfXRayVisionPenetration `json:"penetration"`
	Reuse           ItemRingOfXRayVisionReuse       `json:"reuse"`
	SourcePages     []int                           `json:"source_pages"`
}

// ItemRingOfXRayVisionPenetration is a nested part of a source payload.
type ItemRingOfXRayVisionPenetration struct {
	StoneFeet         int      `json:"stone_feet"`
	CommonMetalInches int      `json:"common_metal_inches"`
	WoodOrDirtFeet    int      `json:"wood_or_dirt_feet"`
	BlockedBy         []string `json:"blocked_by"`
}

// ItemRingOfXRayVisionReuse is a nested part of a source payload.
type ItemRingOfXRayVisionReuse struct {
	Trigger          string `json:"trigger"`
	SaveAbility      string `json:"save_ability"`
	SaveDC           int    `json:"save_dc"`
	FailureCondition string `json:"failure_condition"`
	ExhaustionLevels int    `json:"exhaustion_levels"`
}

// ItemRingOfTheRam retains typed SRD source facts without implementing them.
type ItemRingOfTheRam struct {
	MaxCharges     int                         `json:"max_charges"`
	Recharge       ItemDice                    `json:"recharge"`
	RechargeWhen   string                      `json:"recharge_when"`
	Action         string                      `json:"action"`
	MinimumCharges int                         `json:"minimum_charges"`
	MaximumCharges int                         `json:"maximum_charges"`
	Attack         ItemRingOfTheRamAttack      `json:"attack"`
	BreakObject    ItemRingOfTheRamBreakObject `json:"break_object"`
}

// ItemRingOfTheRamAttack is a nested part of a source payload.
type ItemRingOfTheRamAttack struct {
	Kind              string   `json:"kind"`
	Targets           int      `json:"targets"`
	RangeFeet         int      `json:"range_feet"`
	Form              string   `json:"form"`
	AttackBonus       int      `json:"attack_bonus"`
	DamagePerCharge   ItemDice `json:"damage_per_charge"`
	PushFeetPerCharge int      `json:"push_feet_per_charge"`
}

// ItemRingOfTheRamBreakObject is a nested part of a source payload.
type ItemRingOfTheRamBreakObject struct {
	Action              string   `json:"action"`
	Target              []string `json:"target"`
	RangeFeet           int      `json:"range_feet"`
	CheckAbility        string   `json:"check_ability"`
	CheckBonusPerCharge int      `json:"check_bonus_per_charge"`
}

// ItemRobeOfEyes retains typed SRD source facts without implementing them.
type ItemRobeOfEyes struct {
	Appearance      string                        `json:"appearance"`
	AllAroundVision ItemRobeOfEyesAllAroundVision `json:"all_around_vision"`
	Senses          []string                      `json:"senses"`
	SenseRangeFeet  int                           `json:"sense_range_feet"`
	Drawbacks       ItemRobeOfEyesDrawbacks       `json:"drawbacks"`
}

// ItemRobeOfEyesAllAroundVision is a nested part of a source payload.
type ItemRobeOfEyesAllAroundVision struct {
	CheckAbility     string `json:"check_ability"`
	CheckSkill       string `json:"check_skill"`
	CheckMode        string `json:"check_mode"`
	CheckRequirement string `json:"check_requirement"`
}

// ItemRobeOfEyesDrawbacks is a nested part of a source payload.
type ItemRobeOfEyesDrawbacks struct {
	Condition       string                           `json:"condition"`
	DurationMinutes int                              `json:"duration_minutes"`
	Triggers        []ItemRobeOfEyesDrawbacksTrigger `json:"triggers"`
	SaveAbility     string                           `json:"save_ability"`
	SaveWhen        string                           `json:"save_when"`
	Success         string                           `json:"success"`
}

// ItemRobeOfEyesDrawbacksTrigger is a nested part of a source payload.
type ItemRobeOfEyesDrawbacksTrigger struct {
	Spell     string `json:"spell"`
	Placement string `json:"placement"`
	SaveDC    int    `json:"save_dc"`
}

// ItemRobeOfScintillatingColors retains typed SRD source facts without implementing them.
type ItemRobeOfScintillatingColors struct {
	MaxCharges             int                               `json:"max_charges"`
	Recharge               ItemDice                          `json:"recharge"`
	RechargeWhen           string                            `json:"recharge_when"`
	Condition              string                            `json:"condition"`
	Action                 string                            `json:"action"`
	ChargeCost             int                               `json:"charge_cost"`
	Effect                 string                            `json:"effect"`
	Duration               string                            `json:"duration"`
	BrightLightRadiusFeet  int                               `json:"bright_light_radius_feet"`
	DimLightAdditionalFeet int                               `json:"dim_light_additional_feet"`
	AttacksAgainstYou      string                            `json:"attacks_against_you"`
	Stun                   ItemRobeOfScintillatingColorsStun `json:"stun"`
}

// ItemRobeOfScintillatingColorsStun is a nested part of a source payload.
type ItemRobeOfScintillatingColorsStun struct {
	Affects       string `json:"affects"`
	SaveAbility   string `json:"save_ability"`
	SaveDC        int    `json:"save_dc"`
	Condition     string `json:"condition"`
	ConditionEnds string `json:"condition_ends"`
}

// ItemRobeOfStars retains typed SRD source facts without implementing them.
type ItemRobeOfStars struct {
	Appearance       string                      `json:"appearance"`
	SavingThrowBonus int                         `json:"saving_throw_bonus"`
	Missiles         ItemRobeOfStarsMissiles     `json:"missiles"`
	AstralTravel     ItemRobeOfStarsAstralTravel `json:"astral_travel"`
}

// ItemRobeOfStarsMissiles is a nested part of a source payload.
type ItemRobeOfStarsMissiles struct {
	LargeStars  int      `json:"large_stars"`
	Action      string   `json:"action"`
	Spell       string   `json:"spell"`
	SpellLevel  int      `json:"spell_level"`
	StarsReturn ItemDice `json:"stars_return"`
	ReturnWhen  string   `json:"return_when"`
}

// ItemRobeOfStarsAstralTravel is a nested part of a source payload.
type ItemRobeOfStarsAstralTravel struct {
	Action      string `json:"action"`
	Destination string `json:"destination"`
	Takes       string `json:"takes"`
	RemainUntil string `json:"remain_until"`
	ReappearAt  string `json:"reappear_at"`
}

// ItemRobeOfUsefulItems retains typed SRD source facts without implementing them.
type ItemRobeOfUsefulItems struct {
	Condition              string                            `json:"condition"`
	Action                 string                            `json:"action"`
	PatchBecomes           string                            `json:"patch_becomes"`
	AfterLastPatch         string                            `json:"after_last_patch"`
	CopiesOfEachFixedPatch int                               `json:"copies_of_each_fixed_patch"`
	FixedPatches           []string                          `json:"fixed_patches"`
	OtherPatches           ItemDice                          `json:"other_patches"`
	OtherPatchesChosenBy   string                            `json:"other_patches_chosen_by"`
	TableDieSides          int                               `json:"table_die_sides"`
	PatchTable             []ItemRobeOfUsefulItemsPatchTable `json:"patch_table"`
	PatchFacts             []ItemRobeOfUsefulItemsPatchFact  `json:"patch_facts"`
	PatchRules             []ItemRobeOfUsefulItemsPatchRule  `json:"patch_rules"`
	SourcePages            []int                             `json:"source_pages"`
}

// ItemRobeOfUsefulItemsPatchTable is a nested part of a source payload.
type ItemRobeOfUsefulItemsPatchTable struct {
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Patch string `json:"patch"`
}

// ItemRobeOfUsefulItemsPatchFact is a nested part of a source payload.
type ItemRobeOfUsefulItemsPatchFact struct {
	Patch  string `json:"patch"`
	Key    string `json:"key"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// ItemRobeOfUsefulItemsPatchRule is a nested part of a source payload.
type ItemRobeOfUsefulItemsPatchRule struct {
	Patch string `json:"patch"`
	Rule  string `json:"rule"`
}

// ItemRobeOfTheArchmagi retains typed SRD source facts without implementing them.
type ItemRobeOfTheArchmagi struct {
	AttunementByAnyOf []string                             `json:"attunement_by_any_of"`
	Armor             ItemRobeOfTheArchmagiArmor           `json:"armor"`
	MagicResistance   ItemRobeOfTheArchmagiMagicResistance `json:"magic_resistance"`
	WarMage           ItemRobeOfTheArchmagiWarMage         `json:"war_mage"`
}

// ItemRobeOfTheArchmagiArmor is a nested part of a source payload.
type ItemRobeOfTheArchmagiArmor struct {
	Requirement         string `json:"requirement"`
	BaseArmorClass      int    `json:"base_armor_class"`
	AddsAbilityModifier string `json:"adds_ability_modifier"`
}

// ItemRobeOfTheArchmagiMagicResistance is a nested part of a source payload.
type ItemRobeOfTheArchmagiMagicResistance struct {
	SaveMode string   `json:"save_mode"`
	Against  []string `json:"against"`
}

// ItemRobeOfTheArchmagiWarMage is a nested part of a source payload.
type ItemRobeOfTheArchmagiWarMage struct {
	Increases []string `json:"increases"`
	Bonus     int      `json:"bonus"`
}

// ItemRodOfAbsorption retains typed SRD source facts without implementing them.
type ItemRodOfAbsorption struct {
	Condition              string                     `json:"condition"`
	Absorb                 ItemRodOfAbsorptionAbsorb  `json:"absorb"`
	LifetimeCapacityLevels int                        `json:"lifetime_capacity_levels"`
	KnownOnAttunement      []string                   `json:"known_on_attunement"`
	Convert                ItemRodOfAbsorptionConvert `json:"convert"`
	FoundStoredLevels      ItemDice                   `json:"found_stored_levels"`
	NonmagicalWhen         []string                   `json:"nonmagical_when"`
}

// ItemRodOfAbsorptionAbsorb is a nested part of a source payload.
type ItemRodOfAbsorptionAbsorb struct {
	Action           string   `json:"action"`
	Requires         []string `json:"requires"`
	Result           string   `json:"result"`
	EnergyLevel      string   `json:"energy_level"`
	CastingResources string   `json:"casting_resources"`
	IfItCannotStore  string   `json:"if_it_cannot_store"`
}

// ItemRodOfAbsorptionConvert is a nested part of a source payload.
type ItemRodOfAbsorptionConvert struct {
	Requirement      string `json:"requirement"`
	Result           string `json:"result"`
	LevelLimit       string `json:"level_limit"`
	MaximumSlotLevel int    `json:"maximum_slot_level"`
	Casting          string `json:"casting"`
}

// ItemRodOfAlertness retains typed SRD source facts without implementing them.
type ItemRodOfAlertness struct {
	Condition      string                           `json:"condition"`
	AdvantageOn    []string                         `json:"advantage_on"`
	Spells         []string                         `json:"spells"`
	ProtectiveAura ItemRodOfAlertnessProtectiveAura `json:"protective_aura"`
}

// ItemRodOfAlertnessProtectiveAura is a nested part of a source payload.
type ItemRodOfAlertnessProtectiveAura struct {
	Action                 string   `json:"action"`
	BrightLightRadiusFeet  int      `json:"bright_light_radius_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	Beneficiaries          []string `json:"beneficiaries"`
	BenefitRequirement     string   `json:"benefit_requirement"`
	ArmorClassBonus        int      `json:"armor_class_bonus"`
	SavingThrowBonus       int      `json:"saving_throw_bonus"`
	Senses                 string   `json:"senses"`
	DurationMinutes        int      `json:"duration_minutes"`
	EndsEarlyWhen          string   `json:"ends_early_when"`
	ReuseAfter             string   `json:"reuse_after"`
}

// itemManifest240 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest240 = itemSourceManifest{Batch: "data-magic-items-16", Issue: 240, Families: []itemSourceManifestFamily{
	{Slug: "ring-of-three-wishes", Name: "Ring of Three Wishes", Page: 239, Category: "ring", Attunement: "none", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "none"},
		Keys:    "item.wear item.charges spell.cast-wish item.nonmagical-after-last-charge", Payload: "three_wishes",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Three Wishes", "legendary"}}},
	{Slug: "ring-of-warmth", Name: "Ring of Warmth", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear damage.reduce-cold-by-2d8 defense.unharmed-at-0-degrees-fahrenheit-or-lower", Payload: "warmth",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Warmth", "uncommon"}}},
	{Slug: "ring-of-water-walking", Name: "Ring of Water Walking", Page: 239, Category: "ring", Attunement: "none", Equip: "wear", Use: "cast-spell", UseRestrictions: "water-walk-targeting-yourself",
		Keys: "item.wear spell.cast-water-walk spell.target-self-only", Payload: "water_walking",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Water Walking", "uncommon"}}},
	{Slug: "ring-of-x-ray-vision", Name: "Ring of X-ray Vision", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "magic-action",
		Keys: "item.attunement item.wear action.magic senses.x-ray-vision-30-feet duration.1-minute save.constitution-dc-15 condition.exhaustion", Payload: "x_ray",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of X-ray Vision", "rare"}}},
	{Slug: "ring-of-the-ram", Name: "Ring of the Ram", Page: 239, Category: "ring", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "expend-one-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.wear item.charges action.magic attack.ranged-spell-plus-7 damage.force-2d10-per-charge movement.push-5-feet-per-charge check.break-object-strength-plus-5-per-charge", Payload: "ram",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of the Ram", "rare"}}},
	{Slug: "robe-of-eyes", Name: "Robe of Eyes", Page: 240, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear check.perception-sight-advantage senses.darkvision-120-feet senses.truesight-120-feet effect.blinded-by-light-or-daylight save.constitution-end-blinded", Payload: "eyes",
		Variants: []itemSourceManifestVariant{{"standard", "Robe of Eyes", "rare"}}},
	{Slug: "robe-of-scintillating-colors", Name: "Robe of Scintillating Colors", Page: 240, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.wear item.charges action.magic light.bright-30-feet-dim-30-feet attack.disadvantage-against-wearer save.wisdom-dc-15 condition.stunned", Payload: "scintillating",
		Variants: []itemSourceManifestVariant{{"standard", "Robe of Scintillating Colors", "very-rare"}}},
	{Slug: "robe-of-stars", Name: "Robe of Stars", Page: 240, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action",
		Keys: "item.attunement item.wear save.all-bonus-1 spell.cast-magic-missile-level-5 item.removable-stars movement.enter-and-leave-astral-plane", Payload: "stars",
		Variants: []itemSourceManifestVariant{{"standard", "Robe of Stars", "very-rare"}}},
	{Slug: "robe-of-useful-items", Name: "Robe of Useful Items", Page: 240, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "magic-action", UseRestrictions: "detach-one-patch",
		Keys: "item.wear action.magic item.detachable-patches table.other-patches-1d100 item.ordinary-after-last-patch", Payload: "useful_items",
		Variants: []itemSourceManifestVariant{{"standard", "Robe of Useful Items", "uncommon"}}},
	{Slug: "robe-of-the-archmagi", Name: "Robe of the Archmagi", Page: 240, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement attunement.sorcerer-warlock-or-wizard item.wear defense.base-armor-class-15-plus-dexterity-modifier save.magic-advantage bonus.spell-save-dc-and-attack-2", Payload: "archmagi",
		Variants: []itemSourceManifestVariant{{"standard", "Robe of the Archmagi", "legendary"}}},
	{Slug: "rod-of-absorption", Name: "Rod of Absorption", Page: 241, Category: "rod", Attunement: "required", Equip: "hold", Use: "reaction", UseRestrictions: "spell-targets-only-you-without-area",
		Keys: "item.attunement item.hold action.reaction-absorb-spell item.lifetime-capacity-50-spell-levels spell.convert-stored-energy-to-slots item.nonmagical-when-exhausted", Payload: "absorption",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Absorption", "very-rare"}}},
	{Slug: "rod-of-alertness", Name: "Rod of Alertness", Page: 241, Category: "rod", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "listed-rod-spells",
		Keys: "item.attunement item.hold check.perception-advantage check.initiative-advantage spell.cast-listed-spells action.magic-protective-aura bonus.armor-class-and-saves-1 effect.sense-invisible-creatures item.once-until-dawn", Payload: "alertness",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Alertness", "very-rare"}}},
}}

// LoadItems240 returns a fresh validated copy of batch #240. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems240() (ItemSourceBatch[ItemDetails240], error) { return DecodeItems240(items240JSON) }

// DecodeItems240 validates supplied JSON as batch #240 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems240(data []byte) (ItemSourceBatch[ItemDetails240], error) {
	return decodeItemSourceBatch(data, itemManifest240, validateItemDetails240)
}

// validateItemDetails240 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails240(family, variant string, d *ItemDetails240) error {
	ok := true
	switch family {
	case "ring-of-three-wishes":
		p := d.ThreeWishes
		ok = p.Spell == "wish" && p.MaxCharges == 3 && p.ChargeCost == 1
	case "ring-of-warmth":
		p := d.Warmth
		ok = p.Reduction == ItemDice{Dice: 2, Sides: 8, Type: "damage-reduction"} && p.ProtectedAtOrBelowFahrenheit == 0 && uniqueItemKeys(p.Protected)
	case "ring-of-water-walking":
		p := d.WaterWalking
		ok = p.Spell == "water-walk" && p.Target == "yourself-only"
	case "ring-of-x-ray-vision":
		p := d.XRay
		ok = p.RangeFeet == 30 && p.DurationMinutes == 1 && p.Reuse.SaveDC == 15 && p.Reuse.ExhaustionLevels == 1 && uniqueItemKeys(p.Penetration.BlockedBy) && slices.Equal(p.SourcePages, []int{239, 240})
	case "ring-of-the-ram":
		p := d.Ram
		ok = p.MaxCharges == 3 && p.Recharge == ItemDice{Dice: 1, Sides: 3, Type: "charges"} && p.Attack.AttackBonus == 7 && p.Attack.DamagePerCharge == ItemDice{Dice: 2, Sides: 10, Type: "force"} && p.Attack.PushFeetPerCharge == 5 && p.BreakObject.CheckBonusPerCharge == 5 && uniqueItemKeys(p.BreakObject.Target)
	case "robe-of-eyes":
		p := d.Eyes
		d := p.Drawbacks
		ok = p.SenseRangeFeet == 120 && slices.Equal(p.Senses, []string{"darkvision", "truesight"}) && d.DurationMinutes == 1 && len(d.Triggers) == 2 && d.Triggers[0].SaveDC == 11 && d.Triggers[1].SaveDC == 15
	case "robe-of-scintillating-colors":
		p := d.Scintillating
		ok = p.MaxCharges == 3 && p.Recharge == ItemDice{Dice: 1, Sides: 3, Type: "charges"} && p.BrightLightRadiusFeet == 30 && p.DimLightAdditionalFeet == 30 && p.Stun.SaveDC == 15 && p.Stun.Condition == "stunned"
	case "robe-of-stars":
		p := d.Stars
		m := p.Missiles
		ok = p.SavingThrowBonus == 1 && m.LargeStars == 6 && m.SpellLevel == 5 && m.Spell == "magic-missile" && m.StarsReturn == ItemDice{Dice: 1, Sides: 6, Type: "stars"}
	case "robe-of-useful-items":
		p := d.UsefulItems
		rows := [][2]int{}
		patches := []string{}
		for _, r := range p.PatchTable {
			rows = append(rows, [2]int{r.Min, r.Max})
			patches = append(patches, r.Patch)
		}
		facts := []string{}
		for _, f := range p.PatchFacts {
			ok = ok && f.Amount > 0 && slices.Contains(patches, f.Patch)
			facts = append(facts, f.Patch+"/"+f.Key)
		}
		rules := []string{}
		for _, r := range p.PatchRules {
			ok = ok && slices.Contains(patches, r.Patch)
			rules = append(rules, r.Rule)
		}
		ok = ok && uniqueItemKeys(rules) && p.CopiesOfEachFixedPatch == 2 && p.OtherPatches == ItemDice{Dice: 4, Sides: 4, Type: "patches"} && p.TableDieSides == 100 && len(rows) == 13 && itemRollTable(rows, 100) &&
			uniqueItemKeys(patches) && uniqueItemKeys(facts) && uniqueItemKeys(p.FixedPatches) && len(p.FixedPatches) == 6 && slices.Equal(p.SourcePages, []int{240, 241})
	case "robe-of-the-archmagi":
		p := d.Archmagi
		ok = p.Armor.BaseArmorClass == 15 && p.WarMage.Bonus == 2 && slices.Equal(p.AttunementByAnyOf, []string{"sorcerer", "warlock", "wizard"})
	case "rod-of-absorption":
		p := d.Absorption
		ok = p.LifetimeCapacityLevels == 50 && p.Convert.MaximumSlotLevel == 5 && p.FoundStoredLevels == ItemDice{Dice: 1, Sides: 10, Type: "spell-energy-levels"} && uniqueItemKeys(p.Absorb.Requires) && uniqueItemKeys(p.KnownOnAttunement) && uniqueItemKeys(p.NonmagicalWhen)
	case "rod-of-alertness":
		p := d.Alertness
		a := p.ProtectiveAura
		ok = a.BrightLightRadiusFeet == 60 && a.DimLightAdditionalFeet == 60 && a.ArmorClassBonus == 1 && a.SavingThrowBonus == 1 && a.DurationMinutes == 10 && uniqueItemKeys(p.Spells) && len(p.Spells) == 4 && uniqueItemKeys(p.AdvantageOn)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
