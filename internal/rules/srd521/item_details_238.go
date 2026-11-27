package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_238.json
var items238JSON []byte

// ItemDetails238 holds exactly one typed payload per variant of catalogue batch
// #238 (magic items 14, SRD 5.2.1 pp. 236-238). Payloads are source facts only.
type ItemDetails238 struct {
	Poison           *ItemPotionOfPoison           `json:"poison,omitempty"`
	ResistancePotion *ItemPotionOfResistance       `json:"resistance_potion,omitempty"`
	Speed            *ItemPotionOfSpeed            `json:"speed,omitempty"`
	Vitality         *ItemPotionOfVitality         `json:"vitality,omitempty"`
	WaterBreathing   *ItemPotionOfWaterBreathing   `json:"water_breathing,omitempty"`
	Acrobat          *ItemQuarterstaffOfTheAcrobat `json:"acrobat,omitempty"`
	AnimalInfluence  *ItemRingOfAnimalInfluence    `json:"animal_influence,omitempty"`
	DjinniSummoning  *ItemRingOfDjinniSummoning    `json:"djinni_summoning,omitempty"`
	ElementalCommand *ItemRingOfElementalCommand   `json:"elemental_command,omitempty"`
	Evasion          *ItemRingOfEvasion            `json:"evasion,omitempty"`
	FeatherFalling   *ItemRingOfFeatherFalling     `json:"feather_falling,omitempty"`
}

// ItemPotionOfPoison retains typed SRD source facts without implementing them.
type ItemPotionOfPoison struct {
	Disguise      string   `json:"disguise"`
	MaskedBy      string   `json:"masked_by"`
	RevealedBy    string   `json:"revealed_by"`
	Damage        ItemDice `json:"damage"`
	SaveAbility   string   `json:"save_ability"`
	SaveDC        int      `json:"save_dc"`
	Condition     string   `json:"condition"`
	DurationHours int      `json:"duration_hours"`
}

// ItemPotionOfResistance retains typed SRD source facts without implementing them.
type ItemPotionOfResistance struct {
	Defense       string                             `json:"defense"`
	DurationHours int                                `json:"duration_hours"`
	TypeChosenBy  string                             `json:"type_chosen_by"`
	TableDieSides int                                `json:"table_die_sides"`
	DamageTypes   []ItemPotionOfResistanceDamageType `json:"damage_types"`
}

// ItemPotionOfResistanceDamageType is a nested part of a source payload.
type ItemPotionOfResistanceDamageType struct {
	Roll       int    `json:"roll"`
	DamageType string `json:"damage_type"`
}

// ItemPotionOfSpeed retains typed SRD source facts without implementing them.
type ItemPotionOfSpeed struct {
	Spell                 string `json:"spell"`
	DurationMinutes       int    `json:"duration_minutes"`
	ConcentrationRequired bool   `json:"concentration_required"`
	LethargyWhenItEnds    bool   `json:"lethargy_when_it_ends"`
	Appearance            string `json:"appearance"`
}

// ItemPotionOfVitality retains typed SRD source facts without implementing them.
type ItemPotionOfVitality struct {
	Removes       string `json:"removes"`
	EndsCondition string `json:"ends_condition"`
	HitDieHealing string `json:"hit_die_healing"`
	DurationHours int    `json:"duration_hours"`
	Appearance    string `json:"appearance"`
}

// ItemPotionOfWaterBreathing retains typed SRD source facts without implementing them.
type ItemPotionOfWaterBreathing struct {
	Grants        string `json:"grants"`
	DurationHours int    `json:"duration_hours"`
	Appearance    string `json:"appearance"`
	SourcePages   []int  `json:"source_pages"`
}

// ItemQuarterstaffOfTheAcrobat retains typed SRD source facts without implementing them.
type ItemQuarterstaffOfTheAcrobat struct {
	AttackBonus      int                                          `json:"attack_bonus"`
	DamageBonus      int                                          `json:"damage_bonus"`
	Light            ItemQuarterstaffOfTheAcrobatLight            `json:"light"`
	AlterForm        ItemQuarterstaffOfTheAcrobatAlterForm        `json:"alter_form"`
	AcrobaticAssist  ItemQuarterstaffOfTheAcrobatAcrobaticAssist  `json:"acrobatic_assist"`
	AttackDeflection ItemQuarterstaffOfTheAcrobatAttackDeflection `json:"attack_deflection"`
	RangedWeapon     ItemQuarterstaffOfTheAcrobatRangedWeapon     `json:"ranged_weapon"`
}

// ItemQuarterstaffOfTheAcrobatLight is a nested part of a source payload.
type ItemQuarterstaffOfTheAcrobatLight struct {
	Condition        string   `json:"condition"`
	Color            string   `json:"color"`
	DimRadiusFeet    int      `json:"dim_radius_feet"`
	Activation       []string `json:"activation"`
	ExtinguishAction string   `json:"extinguish_action"`
}

// ItemQuarterstaffOfTheAcrobatAlterForm is a nested part of a source payload.
type ItemQuarterstaffOfTheAcrobatAlterForm struct {
	Action          string   `json:"action"`
	Forms           []string `json:"forms"`
	RodInches       int      `json:"rod_inches"`
	PoleFeet        int      `json:"pole_feet"`
	ElongationLimit string   `json:"elongation_limit"`
}

// ItemQuarterstaffOfTheAcrobatAcrobaticAssist is a nested part of a source payload.
type ItemQuarterstaffOfTheAcrobatAcrobaticAssist struct {
	Forms        []string `json:"forms"`
	CheckAbility string   `json:"check_ability"`
	CheckSkill   string   `json:"check_skill"`
	CheckMode    string   `json:"check_mode"`
}

// ItemQuarterstaffOfTheAcrobatAttackDeflection is a nested part of a source payload.
type ItemQuarterstaffOfTheAcrobatAttackDeflection struct {
	Form            string `json:"form"`
	Trigger         string `json:"trigger"`
	Action          string `json:"action"`
	ArmorClassBonus int    `json:"armor_class_bonus"`
	PossibleResult  string `json:"possible_result"`
	ReuseAfter      string `json:"reuse_after"`
}

// ItemQuarterstaffOfTheAcrobatRangedWeapon is a nested part of a source payload.
type ItemQuarterstaffOfTheAcrobatRangedWeapon struct {
	Form            string `json:"form"`
	Property        string `json:"property"`
	NormalRangeFeet int    `json:"normal_range_feet"`
	LongRangeFeet   int    `json:"long_range_feet"`
	ReturnsToHand   string `json:"returns_to_hand"`
}

// ItemRingOfAnimalInfluence retains typed SRD source facts without implementing them.
type ItemRingOfAnimalInfluence struct {
	MaxCharges      int      `json:"max_charges"`
	Recharge        ItemDice `json:"recharge"`
	RechargeWhen    string   `json:"recharge_when"`
	Condition       string   `json:"condition"`
	ChargeCost      int      `json:"charge_cost"`
	SaveDC          int      `json:"save_dc"`
	Spells          []string `json:"spells"`
	FearAffectsOnly string   `json:"fear_affects_only"`
}

// ItemRingOfDjinniSummoning retains typed SRD source facts without implementing them.
type ItemRingOfDjinniSummoning struct {
	Condition             string `json:"condition"`
	Action                string `json:"action"`
	Creature              string `json:"creature"`
	Origin                string `json:"origin"`
	SummonRangeFeet       int    `json:"summon_range_feet"`
	ConcentrationRequired bool   `json:"concentration_required"`
	MaxDurationHours      int    `json:"max_duration_hours"`
	AlsoEndsAt            string `json:"also_ends_at"`
	Attitude              string `json:"attitude"`
	ObeysCommands         bool   `json:"obeys_commands"`
	WithoutCommand        string `json:"without_command"`
	CooldownHours         int    `json:"cooldown_hours"`
	NonmagicalWhen        string `json:"nonmagical_when"`
	OftenCreatedBy        string `json:"often_created_by"`
}

// ItemRingOfElementalCommand retains typed SRD source facts without implementing them.
type ItemRingOfElementalCommand struct {
	LinkedPlane   string                               `json:"linked_plane"`
	PlaneChosenBy string                               `json:"plane_chosen_by"`
	Bane          ItemRingOfElementalCommandBane       `json:"bane"`
	Compulsion    ItemRingOfElementalCommandCompulsion `json:"compulsion"`
	Focus         ItemRingOfElementalCommandFocus      `json:"focus"`
	SpellSaveDC   int                                  `json:"spell_save_dc"`
	Spells        []ItemRingOfElementalCommandSpell    `json:"spells"`
	MaxCharges    int                                  `json:"max_charges"`
	Recharge      ItemDice                             `json:"recharge"`
	RechargeBonus int                                  `json:"recharge_bonus"`
	RechargeWhen  string                               `json:"recharge_when"`
	SourcePages   []int                                `json:"source_pages"`
}

// ItemRingOfElementalCommandBane is a nested part of a source payload.
type ItemRingOfElementalCommandBane struct {
	TargetCreatureType string `json:"target_creature_type"`
	YourAttacks        string `json:"your_attacks"`
	AttacksAgainstYou  string `json:"attacks_against_you"`
}

// ItemRingOfElementalCommandCompulsion is a nested part of a source payload.
type ItemRingOfElementalCommandCompulsion struct {
	Action        string `json:"action"`
	RangeFeet     int    `json:"range_feet"`
	SaveAbility   string `json:"save_ability"`
	SaveDC        int    `json:"save_dc"`
	Condition     string `json:"condition"`
	ConditionEnds string `json:"condition_ends"`
	Control       string `json:"control"`
}

// ItemRingOfElementalCommandFocus is a nested part of a source payload.
type ItemRingOfElementalCommandFocus struct {
	Language      string   `json:"language"`
	Defense       string   `json:"defense"`
	DamageType    string   `json:"damage_type"`
	SpeedKind     string   `json:"speed_kind"`
	Speed         string   `json:"speed"`
	OtherBenefits []string `json:"other_benefits" srd:"may-be-empty"`
}

// ItemRingOfElementalCommandSpell is a nested part of a source payload.
type ItemRingOfElementalCommandSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemRingOfEvasion retains typed SRD source facts without implementing them.
type ItemRingOfEvasion struct {
	MaxCharges   int      `json:"max_charges"`
	Recharge     ItemDice `json:"recharge"`
	RechargeWhen string   `json:"recharge_when"`
	Trigger      string   `json:"trigger"`
	Action       string   `json:"action"`
	ChargeCost   int      `json:"charge_cost"`
	Result       string   `json:"result"`
}

// ItemRingOfFeatherFalling retains typed SRD source facts without implementing them.
type ItemRingOfFeatherFalling struct {
	Condition           string `json:"condition"`
	DescentFeetPerRound int    `json:"descent_feet_per_round"`
	TakesFallingDamage  bool   `json:"takes_falling_damage"`
}

// itemManifest238 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest238 = itemSourceManifest{Batch: "data-magic-items-14", Issue: 238, Families: []itemSourceManifestFamily{
	{Slug: "potion-of-poison", Name: "Potion of Poison", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume item.disguised-as-beneficial-potion damage.4d6-poison save.constitution-dc-13 condition.poisoned duration.1-hour", Payload: "poison",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Poison", "uncommon"}}},
	{Slug: "potion-of-resistance", Name: "Potion of Resistance", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume defense.resistance-one-damage-type table.damage-type-1d10 duration.1-hour", Payload: "resistance_potion",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Resistance", "uncommon"}}},
	{Slug: "potion-of-speed", Name: "Potion of Speed", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume spell.effect-haste duration.1-minute spell.no-concentration spell.no-lethargy", Payload: "speed",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Speed", "very-rare"}}},
	{Slug: "potion-of-vitality", Name: "Potion of Vitality", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume condition.remove-exhaustion-levels condition.end-poisoned hit-dice.maximum-healing duration.24-hours", Payload: "vitality",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Vitality", "very-rare"}}},
	{Slug: "potion-of-water-breathing", Name: "Potion of Water Breathing", Page: 236, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume movement.breathe-underwater duration.24-hours", Payload: "water_breathing",
		Variants: []itemSourceManifestVariant{{"standard", "Potion of Water Breathing", "uncommon"}}},
	{Slug: "quarterstaff-of-the-acrobat", Name: "Quarterstaff of the Acrobat", Page: 237, Category: "weapon", Options: "quarterstaff", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon bonus.attack-and-damage-2 light.green-dim-10-feet item.alter-form check.dexterity-acrobatics-advantage reaction.armor-class-bonus-5 property.thrown-30-120 item.returns-to-hand", Payload: "acrobat",
		Variants: []itemSourceManifestVariant{{"standard", "Quarterstaff of the Acrobat", "very-rare"}}},
	{Slug: "ring-of-animal-influence", Name: "Ring of Animal Influence", Page: 237, Category: "ring", Attunement: "none", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.wear item.charges spell.cast-animal-friendship spell.cast-fear-beasts-only spell.cast-speak-with-animals spell.save-dc-13", Payload: "animal_influence",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Animal Influence", "rare"}}},
	{Slug: "ring-of-djinni-summoning", Name: "Ring of Djinni Summoning", Page: 237, Category: "ring", Attunement: "required", Equip: "wear", Use: "magic-action",
		Keys: "item.attunement item.wear action.magic creature.summon-djinni spell.concentration-up-to-1-hour creature.obeys-commands item.cooldown-24-hours item.nonmagical-if-djinni-dies", Payload: "djinni_summoning",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Djinni Summoning", "legendary"}}},
	{Slug: "ring-of-elemental-command", Name: "Ring of Elemental Command", Page: 237, Category: "ring", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "spell-linked-to-plane expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-plus-1-at-dawn"},
		Keys:    "item.attunement item.wear item.charges attack.advantage-against-elementals action.magic-compel-elemental save.wisdom-dc-18 condition.charmed feature.elemental-focus-by-linked-plane spell.cast-by-linked-plane spell.save-dc-18", Payload: "elemental_command",
		Variants: []itemSourceManifestVariant{
			{"air", "Ring of Elemental Command (air)", "legendary"},
			{"earth", "Ring of Elemental Command (earth)", "legendary"},
			{"fire", "Ring of Elemental Command (fire)", "legendary"},
			{"water", "Ring of Elemental Command (water)", "legendary"},
		}},
	{Slug: "ring-of-evasion", Name: "Ring of Evasion", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "reaction", UseRestrictions: "failed-dexterity-save expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.wear item.charges action.reaction save.dexterity-failure-becomes-success", Payload: "evasion",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Evasion", "rare"}}},
	{Slug: "ring-of-feather-falling", Name: "Ring of Feather Falling", Page: 238, Category: "ring", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear movement.fall-60-feet-per-round damage.no-falling-damage", Payload: "feather_falling",
		Variants: []itemSourceManifestVariant{{"standard", "Ring of Feather Falling", "rare"}}},
}}

// LoadItems238 returns a fresh validated copy of batch #238. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems238() (ItemSourceBatch[ItemDetails238], error) { return DecodeItems238(items238JSON) }

// DecodeItems238 validates supplied JSON as batch #238 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems238(data []byte) (ItemSourceBatch[ItemDetails238], error) {
	return decodeItemSourceBatch(data, itemManifest238, validateItemDetails238)
}

// validateItemDetails238 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails238(family, variant string, d *ItemDetails238) error {
	ok := true
	switch family {
	case "potion-of-poison":
		p := d.Poison
		ok = p.Damage == ItemDice{Dice: 4, Sides: 6, Type: "poison"} && p.SaveDC == 13 && p.DurationHours == 1 && p.Condition == "poisoned"
	case "potion-of-resistance":
		p := d.ResistancePotion
		rows := [][2]int{}
		types := []string{}
		for _, r := range p.DamageTypes {
			rows = append(rows, [2]int{r.Roll, r.Roll})
			types = append(types, r.DamageType)
		}
		ok = p.DurationHours == 1 && p.TableDieSides == 10 && itemRollTable(rows, 10) && uniqueItemKeys(types)
	case "potion-of-speed":
		p := d.Speed
		ok = p.Spell == "haste" && p.DurationMinutes == 1 && !p.ConcentrationRequired && !p.LethargyWhenItEnds
	case "potion-of-vitality":
		p := d.Vitality
		ok = p.DurationHours == 24 && p.EndsCondition == "poisoned"
	case "potion-of-water-breathing":
		p := d.WaterBreathing
		ok = p.DurationHours == 24 && slices.Equal(p.SourcePages, []int{236, 237})
	case "quarterstaff-of-the-acrobat":
		p := d.Acrobat
		ok = p.AttackBonus == 2 && p.DamageBonus == 2 && p.Light.DimRadiusFeet == 10 && p.AttackDeflection.ArmorClassBonus == 5 && p.RangedWeapon.NormalRangeFeet == 30 && p.RangedWeapon.LongRangeFeet == 120 && uniqueItemKeys(p.AlterForm.Forms) && uniqueItemKeys(p.AcrobaticAssist.Forms) && uniqueItemKeys(p.Light.Activation)
	case "ring-of-animal-influence":
		p := d.AnimalInfluence
		ok = p.MaxCharges == 3 && p.SaveDC == 13 && p.ChargeCost == 1 && p.Recharge == ItemDice{Dice: 1, Sides: 3, Type: "charges"} && slices.Equal(p.Spells, []string{"animal-friendship", "fear", "speak-with-animals"})
	case "ring-of-djinni-summoning":
		p := d.DjinniSummoning
		ok = p.SummonRangeFeet == 120 && p.MaxDurationHours == 1 && p.CooldownHours == 24 && p.ConcentrationRequired && p.Creature == "djinni"
	case "ring-of-elemental-command":
		p := d.ElementalCommand
		wantSpells := map[string][]ItemRingOfElementalCommandSpell{
			"air":   {{Spell: "chain-lightning", ChargeCost: 3}, {Spell: "feather-fall", ChargeCost: 0}, {Spell: "gust-of-wind", ChargeCost: 2}, {Spell: "wind-wall", ChargeCost: 1}},
			"earth": {{Spell: "earthquake", ChargeCost: 5}, {Spell: "stone-shape", ChargeCost: 2}, {Spell: "stoneskin", ChargeCost: 3}, {Spell: "wall-of-stone", ChargeCost: 3}},
			"fire":  {{Spell: "burning-hands", ChargeCost: 1}, {Spell: "fireball", ChargeCost: 2}, {Spell: "fire-storm", ChargeCost: 4}, {Spell: "wall-of-fire", ChargeCost: 3}},
			"water": {{Spell: "create-or-destroy-water", ChargeCost: 1}, {Spell: "ice-storm", ChargeCost: 2}, {Spell: "tsunami", ChargeCost: 5}, {Spell: "wall-of-ice", ChargeCost: 3}, {Spell: "water-walk", ChargeCost: 2}},
		}
		ok = p.LinkedPlane == variant && slices.Equal(p.Spells, wantSpells[variant]) && p.MaxCharges == 5 && p.SpellSaveDC == 18 && p.Compulsion.SaveDC == 18 && p.Compulsion.RangeFeet == 60 &&
			(len(p.Focus.OtherBenefits) == 0 || uniqueItemKeys(p.Focus.OtherBenefits)) && slices.Equal(p.SourcePages, []int{237, 238})
	case "ring-of-evasion":
		p := d.Evasion
		ok = p.MaxCharges == 3 && p.ChargeCost == 1 && p.Recharge == ItemDice{Dice: 1, Sides: 3, Type: "charges"}
	case "ring-of-feather-falling":
		p := d.FeatherFalling
		ok = p.DescentFeetPerRound == 60 && !p.TakesFallingDamage
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
