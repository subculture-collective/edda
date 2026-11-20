package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_236.json
var items236JSON []byte

// ItemDetails236 holds exactly one typed payload per variant of catalogue batch
// #236 (magic items 12, SRD 5.2.1 pp. 233-235). Payloads are source facts only.
type ItemDetails236 struct {
	Sharpness    *ItemOilSharpness     `json:"oil_of_sharpness,omitempty"`
	Slipperiness *ItemOilSlipperiness  `json:"oil_of_slipperiness,omitempty"`
	Pearl        *ItemPearlOfPower     `json:"pearl_of_power,omitempty"`
	Health       *ItemPeriaptHealth    `json:"periapt_of_health,omitempty"`
	Proof        *ItemPeriaptProof     `json:"periapt_of_proof_against_poison,omitempty"`
	WoundClosure *ItemPeriaptClosure   `json:"periapt_of_wound_closure,omitempty"`
	Philter      *ItemPhilterOfLove    `json:"philter_of_love,omitempty"`
	Haunting     *ItemPipesHaunting    `json:"pipes_of_haunting,omitempty"`
	Sewers       *ItemPipesSewers      `json:"pipes_of_the_sewers,omitempty"`
	PlateArmor   *ItemPlateEtherealACN `json:"plate_armor_of_etherealness,omitempty"`
	PortableHole *ItemPortableHole     `json:"portable_hole,omitempty"`
	AnimalFriend *ItemPotionAnimal     `json:"potion_of_animal_friendship,omitempty"`
}

// ItemSave236 is a source saving throw.
type ItemSave236 struct {
	Ability string `json:"ability"`
	DC      int    `json:"dc"`
}

// ItemOilSharpness retains typed SRD source facts without implementing them.
type ItemOilSharpness struct {
	CoatsMeleeWeapons  int      `json:"coats_melee_weapons"`
	CoatsAmmunition    int      `json:"coats_ammunition"`
	RequiresNonmagical bool     `json:"requires_nonmagical"`
	DamageTypes        []string `json:"damage_types"`
	ApplicationMinutes int      `json:"application_minutes"`
	WeaponBonus        int      `json:"weapon_bonus"`
	AmmunitionBonus    int      `json:"ammunition_bonus"`
	Appearance         string   `json:"appearance"`
}

// ItemOilCoverage is what one vial covers.
type ItemOilCoverage struct {
	MaxSize          string `json:"max_size"`
	IncludesEquipped bool   `json:"includes_worn_and_carried_equipment"`
	ExtraVialAbove   string `json:"extra_vial_per_size_above"`
}

// ItemOilSlipperiness retains typed SRD source facts without implementing them.
// The text starts on p. 233 and continues on p. 234.
type ItemOilSlipperiness struct {
	SourcePages        []int           `json:"source_pages"`
	Coverage           ItemOilCoverage `json:"coverage"`
	ApplicationMinutes int             `json:"application_minutes"`
	BodySpell          string          `json:"body_spell"`
	BodyDurationHours  int             `json:"body_duration_hours"`
	PourAction         string          `json:"pour_action"`
	PourAreaFeet       int             `json:"pour_area_square_feet"`
	PourSpell          string          `json:"pour_spell"`
	PourDurationHours  int             `json:"pour_duration_hours"`
	Appearance         string          `json:"appearance"`
}

// ItemPearlOfPower retains typed SRD source facts without implementing them.
type ItemPearlOfPower struct {
	AttunementBy  []string `json:"attunement_by_any_of"`
	Condition     string   `json:"condition"`
	Action        string   `json:"action"`
	SlotsRegained int      `json:"spell_slots_regained"`
	MaxSlotLevel  int      `json:"max_slot_level"`
	UseLimit      int      `json:"uses"`
	Recharge      string   `json:"recharge"`
}

// ItemPeriaptHealth retains typed SRD source facts without implementing them.
type ItemPeriaptHealth struct {
	Condition       string   `json:"condition"`
	Action          string   `json:"action"`
	Healing         ItemDice `json:"healing_dice"`
	HealingBonus    int      `json:"healing_bonus"`
	UseLimit        int      `json:"uses"`
	Recharge        string   `json:"recharge"`
	SaveMode        string   `json:"poisoned_save_mode"`
	SaveApplication []string `json:"poisoned_save_applies_to"`
}

// ItemPeriaptProof retains typed SRD source facts without implementing them.
type ItemPeriaptProof struct {
	Appearance         string   `json:"appearance"`
	Condition          string   `json:"condition"`
	ConditionImmunity  []string `json:"condition_immunities"`
	DamageTypeImmunity []string `json:"damage_immunities"`
}

// ItemPeriaptLifePreservation is the Life Preservation property.
type ItemPeriaptLifePreservation struct {
	Trigger       string `json:"trigger"`
	RollAtOrBelow int    `json:"changes_roll_at_or_below"`
	ChangedTo     int    `json:"changed_to"`
	Result        string `json:"result"`
}

// ItemPeriaptHealingBoost is the Natural Healing Boost property.
type ItemPeriaptHealingBoost struct {
	Trigger    string `json:"trigger"`
	Multiplier int    `json:"multiplier"`
}

// ItemPeriaptClosure retains typed SRD source facts without implementing them.
type ItemPeriaptClosure struct {
	Condition    string                      `json:"condition"`
	Preservation ItemPeriaptLifePreservation `json:"life_preservation"`
	HealingBoost ItemPeriaptHealingBoost     `json:"natural_healing_boost"`
}

// ItemPhilterOfLove retains typed SRD source facts without implementing them.
type ItemPhilterOfLove struct {
	Trigger       string `json:"trigger"`
	WindowMinutes int    `json:"window_minutes"`
	Condition     string `json:"condition"`
	CharmedBy     string `json:"charmed_by"`
	DurationHours int    `json:"duration_hours"`
	Appearance    string `json:"appearance"`
}

// ItemPipesHaunting retains typed SRD source facts without implementing them.
// The text starts in the left column of p. 234 and continues in the right.
type ItemPipesHaunting struct {
	Action          string      `json:"action"`
	ChargeCost      int         `json:"charge_cost"`
	Effect          string      `json:"effect"`
	Targets         string      `json:"targets"`
	RangeFeet       int         `json:"range_feet"`
	Save            ItemSave236 `json:"save"`
	Condition       string      `json:"condition"`
	DurationSeconds int         `json:"duration_seconds"`
	RepeatSave      string      `json:"repeat_save"`
	SuccessEnds     string      `json:"success_ends"`
	ImmunityHours   int         `json:"immunity_hours_after_success"`
}

// ItemPipesSewersPassive is the rats' attitude while the pipes are carried.
type ItemPipesSewersPassive struct {
	Condition string   `json:"condition"`
	Creatures []string `json:"creatures"`
	Attitude  string   `json:"attitude"`
	Unless    string   `json:"unless"`
}

// ItemPipesSewersCall is calling swarms with the pipes.
type ItemPipesSewersCall struct {
	PlayAction      string `json:"play_action"`
	CallAction      string `json:"call_action"`
	MinCharges      int    `json:"min_charges"`
	MaxCharges      int    `json:"max_charges"`
	Swarm           string `json:"swarm"`
	SwarmsPerCharge int    `json:"swarms_per_charge"`
	Range           string `json:"range"`
	Shortfall       string `json:"insufficient_rats"`
	DeterminedBy    string `json:"availability_determined_by"`
	Movement        string `json:"swarm_movement"`
}

// ItemPipesSewersSway is swaying an uncontrolled swarm.
type ItemPipesSewersSway struct {
	Trigger        string      `json:"trigger"`
	TriggerFeet    int         `json:"trigger_range_feet"`
	Save           ItemSave236 `json:"save"`
	SuccessResult  string      `json:"success_result"`
	ImmunityHours  int         `json:"success_immunity_hours"`
	FailureResult  string      `json:"failure_result"`
	FailureWhile   string      `json:"failure_lasts_while"`
	ObeysCommands  bool        `json:"friendly_swarm_obeys_commands"`
	Uncommanded    string      `json:"uncommanded_behavior"`
	ControlEnds    string      `json:"control_ends_when"`
	ControlEndFeet int         `json:"control_end_distance_feet"`
	AfterControl   string      `json:"after_control_ends"`
}

// ItemPipesSewers retains typed SRD source facts without implementing them.
type ItemPipesSewers struct {
	Passive ItemPipesSewersPassive `json:"passive"`
	Call    ItemPipesSewersCall    `json:"call_swarms"`
	Sway    ItemPipesSewersSway    `json:"sway_swarms"`
}

// ItemPlateEtherealACN retains typed SRD source facts without implementing them.
type ItemPlateEtherealACN struct {
	Condition  string   `json:"condition"`
	Action     string   `json:"action"`
	Activation string   `json:"activation"`
	Spell      string   `json:"spell"`
	EndsOn     []string `json:"ends_on"`
	UseLimit   int      `json:"uses"`
	Recharge   string   `json:"recharge"`
}

// ItemHoleUnfold is placing and opening the hole.
type ItemHoleUnfold struct {
	Action     string `json:"action"`
	Surface    string `json:"surface"`
	DepthFeet  int    `json:"depth_feet"`
	Shape      string `json:"shape"`
	Plane      string `json:"plane"`
	NoPassages bool   `json:"cannot_create_open_passages"`
	Exit       string `json:"exit"`
}

// ItemHoleClose is folding the hole up.
type ItemHoleClose struct {
	Action   string `json:"action"`
	Contents string `json:"contents"`
	Weight   string `json:"weight"`
}

// ItemHoleEscape is leaving a folded hole.
type ItemHoleEscape struct {
	Action     string    `json:"action"`
	Check      ItemCheck `json:"check"`
	Success    string    `json:"success"`
	WithinFeet int       `json:"appears_within_feet"`
	AirHours   int       `json:"closed_air_hours"`
	AirDivide  string    `json:"air_divided_by"`
}

// ItemHoleNesting is placing the hole inside another extradimensional space.
type ItemHoleNesting struct {
	Containers  []string `json:"containers"`
	Result      string   `json:"result"`
	GateOrigin  string   `json:"gate_origin"`
	PullFeet    int      `json:"pulls_in_within_feet"`
	Exempt      string   `json:"exempt"`
	Destination string   `json:"destination"`
	GateCloses  bool     `json:"gate_closes_afterward"`
	OneWay      bool     `json:"one_way"`
	CanReopen   bool     `json:"can_be_reopened"`
}

// ItemPortableHole retains typed SRD source facts without implementing them.
// The text starts on p. 234 and continues on p. 235.
type ItemPortableHole struct {
	SourcePages []int           `json:"source_pages"`
	Cloth       string          `json:"cloth"`
	DiameterFt  int             `json:"diameter_feet"`
	Unfold      ItemHoleUnfold  `json:"unfold"`
	Close       ItemHoleClose   `json:"close"`
	Escape      ItemHoleEscape  `json:"escape_when_folded"`
	Nesting     ItemHoleNesting `json:"inside_extradimensional_space"`
}

// ItemPotionAnimal retains typed SRD source facts without implementing them.
type ItemPotionAnimal struct {
	Trigger         string   `json:"trigger"`
	Spell           string   `json:"spell"`
	SpellLevel      int      `json:"spell_level"`
	SaveDC          int      `json:"save_dc"`
	Appearance      string   `json:"appearance"`
	AppearanceItems []string `json:"appearance_items"`
}

// itemManifest236 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. No family prints more than one variant.
var itemManifest236 = itemSourceManifest{Batch: "data-magic-items-12", Issue: 236, Families: []itemSourceManifestFamily{
	{Slug: "oil-of-sharpness", Name: "Oil of Sharpness", Page: 233, Category: "potion", Attunement: "none", Equip: "none", Use: "apply", UseRestrictions: "one-vial-coats-one-melee-weapon-or-twenty-ammunition nonmagical-slashing-or-piercing-only",
		Keys: "item.consume action.apply-1-minute target.nonmagical-slashing-or-piercing effect.plus-3-weapon effect.plus-3-ammunition", Payload: "oil_of_sharpness",
		Variants: []itemSourceManifestVariant{{"standard", "Oil of Sharpness", "very-rare"}}},
	{Slug: "oil-of-slipperiness", Name: "Oil of Slipperiness", Page: 233, Category: "potion", Attunement: "none", Equip: "none", Use: "apply", UseRestrictions: "one-vial-covers-one-medium-or-smaller-creature pour-as-magic-action",
		Keys: "item.consume action.apply-10-minutes spell.effect-freedom-of-movement duration.8-hours action.magic area.square-10-foot spell.effect-grease", Payload: "oil_of_slipperiness",
		Variants: []itemSourceManifestVariant{{"standard", "Oil of Slipperiness", "uncommon"}}},
	{Slug: "pearl-of-power", Name: "Pearl of Power", Page: 234, Category: "wondrous-item", Attunement: "required", Equip: "carry", Use: "magic-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.attunement item.carry action.magic spell-slot.regain-one-level-3-or-lower item.once-until-next-dawn", Payload: "pearl_of_power",
		Variants: []itemSourceManifestVariant{{"standard", "Pearl of Power", "uncommon"}}},
	{Slug: "periapt-of-health", Name: "Periapt of Health", Page: 234, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.attunement item.wear action.magic healing.hit-points-2d4-plus-2 item.once-until-next-dawn save.advantage-avoid-or-end-poisoned", Payload: "periapt_of_health",
		Variants: []itemSourceManifestVariant{{"standard", "Periapt of Health", "uncommon"}}},
	{Slug: "periapt-of-proof-against-poison", Name: "Periapt of Proof against Poison", Page: 234, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear immunity.condition-poisoned immunity.damage-poison", Payload: "periapt_of_proof_against_poison",
		Variants: []itemSourceManifestVariant{{"standard", "Periapt of Proof against Poison", "rare"}}},
	{Slug: "periapt-of-wound-closure", Name: "Periapt of Wound Closure", Page: 234, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear death-save.change-9-or-lower-to-10 healing.double-hit-point-die", Payload: "periapt_of_wound_closure",
		Variants: []itemSourceManifestVariant{{"standard", "Periapt of Wound Closure", "uncommon"}}},
	{Slug: "philter-of-love", Name: "Philter of Love", Page: 234, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume condition.charmed duration.1-hour window.10-minutes-after-drinking", Payload: "philter_of_love",
		Variants: []itemSourceManifestVariant{{"standard", "Philter of Love", "uncommon"}}},
	{Slug: "pipes-of-haunting", Name: "Pipes of Haunting", Page: 234, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.charges action.magic area.30-feet save.wisdom-dc-15 condition.frightened duration.1-minute save.repeat-end-of-turn immunity.24-hours-after-success", Payload: "pipes_of_haunting",
		Variants: []itemSourceManifestVariant{{"standard", "Pipes of Haunting", "uncommon"}}},
	{Slug: "pipes-of-the-sewers", Name: "Pipes of the Sewers", Page: 234, Category: "wondrous-item", Attunement: "required", Equip: "carry", Use: "magic-action", UseRestrictions: "bonus-action-expend-one-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.carry item.charges attitude.rats-indifferent action.magic action.bonus summon.swarm-of-rats-per-charge save.wisdom-dc-15 condition.friendly immunity.24-hours-after-success", Payload: "pipes_of_the_sewers",
		Variants: []itemSourceManifestVariant{{"standard", "Pipes of the Sewers", "uncommon"}}},
	{Slug: "plate-armor-of-etherealness", Name: "Plate Armor of Etherealness", Page: 234, Category: "armor", Options: "half-plate-armor plate-armor", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "command-word once-until-next-dawn",
		Keys: "item.attunement item.wear action.magic item.command-word spell.effect-etherealness item.once-until-next-dawn", Payload: "plate_armor_of_etherealness",
		Variants: []itemSourceManifestVariant{{"standard", "Plate Armor of Etherealness", "legendary"}}},
	{Slug: "portable-hole", Name: "Portable Hole", Page: 234, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "unfold-or-fold",
		Keys: "action.magic item.extradimensional-hole-10-feet-deep check.strength-athletics-dc-10 item.air-1-hour-divided hazard.nested-extradimensional-space-destroys-both-and-opens-astral-gate", Payload: "portable_hole",
		Variants: []itemSourceManifestVariant{{"standard", "Portable Hole", "rare"}}},
	{Slug: "potion-of-animal-friendship", Name: "Potion of Animal Friendship", Page: 235, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.cast-animal-friendship spell.level-3-version spell.save-dc-13", Payload: "potion_of_animal_friendship",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Animal Friendship", "uncommon"}}},
}}

// LoadItems236 returns a fresh validated copy of batch #236. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems236() (ItemSourceBatch[ItemDetails236], error) { return DecodeItems236(items236JSON) }

// DecodeItems236 validates supplied JSON as batch #236 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems236(data []byte) (ItemSourceBatch[ItemDetails236], error) {
	return decodeItemSourceBatch(data, itemManifest236, validateItemDetails236)
}

// validateItemDetails236 checks the internal consistency of each payload and
// the source literals that the family's mechanic keys also name.
func validateItemDetails236(family, _ string, d *ItemDetails236) error {
	positive := func(values ...int) bool {
		for _, v := range values {
			if v <= 0 {
				return false
			}
		}
		return true
	}
	ok := true
	switch family {
	case "oil-of-sharpness":
		o := d.Sharpness
		ok = positive(o.CoatsMeleeWeapons, o.CoatsAmmunition, o.ApplicationMinutes) && uniqueItemKeys(o.DamageTypes) && o.WeaponBonus == 3 && o.AmmunitionBonus == 3 // p. 233; named by the mechanic keys
	case "oil-of-slipperiness":
		o := d.Slipperiness
		ok = slices.Equal(o.SourcePages, []int{233, 234}) && positive(o.ApplicationMinutes, o.BodyDurationHours, o.PourAreaFeet, o.PourDurationHours) && o.ApplicationMinutes == 10 && o.BodyDurationHours == 8 && o.PourDurationHours == 8
	case "pearl-of-power":
		p := d.Pearl
		ok = uniqueItemKeys(p.AttunementBy) && positive(p.SlotsRegained, p.MaxSlotLevel, p.UseLimit) && p.MaxSlotLevel == 3
	case "periapt-of-health":
		p := d.Health
		ok = positive(p.Healing.Dice, p.Healing.Sides, p.HealingBonus, p.UseLimit) && uniqueItemKeys(p.SaveApplication) && p.Healing.Dice == 2 && p.Healing.Sides == 4 && p.HealingBonus == 2
	case "periapt-of-proof-against-poison":
		p := d.Proof
		ok = uniqueItemKeys(p.ConditionImmunity) && uniqueItemKeys(p.DamageTypeImmunity)
	case "periapt-of-wound-closure":
		p := d.WoundClosure
		ok = positive(p.Preservation.RollAtOrBelow, p.Preservation.ChangedTo, p.HealingBoost.Multiplier) && p.Preservation.RollAtOrBelow == 9 && p.Preservation.ChangedTo == 10 && p.HealingBoost.Multiplier == 2
	case "philter-of-love":
		p := d.Philter
		ok = positive(p.WindowMinutes, p.DurationHours) && p.WindowMinutes == 10 && p.DurationHours == 1
	case "pipes-of-haunting":
		p := d.Haunting
		ok = positive(p.ChargeCost, p.RangeFeet, p.Save.DC, p.DurationSeconds, p.ImmunityHours) && p.Save.DC == 15 && p.ChargeCost == 1
	case "pipes-of-the-sewers":
		p := d.Sewers
		ok = uniqueItemKeys(p.Passive.Creatures) && positive(p.Call.MinCharges, p.Call.MaxCharges, p.Call.SwarmsPerCharge, p.Sway.TriggerFeet, p.Sway.Save.DC, p.Sway.ImmunityHours, p.Sway.ControlEndFeet) &&
			p.Call.MaxCharges >= p.Call.MinCharges && p.Sway.Save.DC == 15
	case "plate-armor-of-etherealness":
		p := d.PlateArmor
		ok = uniqueItemKeys(p.EndsOn) && positive(p.UseLimit)
	case "portable-hole":
		p := d.PortableHole
		ok = slices.Equal(p.SourcePages, []int{234, 235}) && positive(p.DiameterFt, p.Unfold.DepthFeet, p.Escape.Check.DC, p.Escape.WithinFeet, p.Escape.AirHours, p.Nesting.PullFeet) && uniqueItemKeys(p.Nesting.Containers) && p.Escape.Check.DC == 10 && p.Unfold.DepthFeet == 10
	case "potion-of-animal-friendship":
		p := d.AnimalFriend
		ok = positive(p.SpellLevel, p.SaveDC) && uniqueItemKeys(p.AppearanceItems) && p.SpellLevel == 3 && p.SaveDC == 13
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
