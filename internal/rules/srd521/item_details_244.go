package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_244.json
var items244JSON []byte

// ItemDetails244 holds exactly one typed payload per variant of catalogue batch
// #244 (magic items 20, SRD 5.2.1 pp. 248-250). Payloads are source facts only.
type ItemDetails244 struct {
	TalismanGood        *ItemTalismanOfPureGood     `json:"talisman_good,omitempty"`
	TalismanEvil        *ItemTalismanOfUltimateEvil `json:"talisman_evil,omitempty"`
	TalismanSphere      *ItemTalismanOfTheSphere    `json:"talisman_sphere,omitempty"`
	ThunderousGreatclub *ItemThunderousGreatclub    `json:"thunderous_greatclub,omitempty"`
	ClearThought        *ItemTomeOfClearThought     `json:"clear_thought,omitempty"`
	Leadership          *ItemTomeOfLeadership       `json:"leadership,omitempty"`
	Understanding       *ItemTomeOfUnderstanding    `json:"understanding,omitempty"`
	FishCommand         *ItemTridentOfFishCommand   `json:"fish_command,omitempty"`
	Solvent             *ItemUniversalSolvent       `json:"solvent,omitempty"`
	Vicious             *ItemViciousWeapon          `json:"vicious,omitempty"`
	Vorpal              *ItemVorpalSword            `json:"vorpal,omitempty"`
	Binding             *ItemWandOfBinding          `json:"binding,omitempty"`
}

// ItemTalismanOfPureGood retains typed SRD source facts without implementing them.
type ItemTalismanOfPureGood struct {
	AttunementByAnyOf []string                          `json:"attunement_by_any_of"`
	TouchDamage       ItemTalismanOfPureGoodTouchDamage `json:"touch_damage"`
	HolySymbol        bool                              `json:"holy_symbol"`
	SpellAttackBonus  int                               `json:"spell_attack_bonus"`
	MaxCharges        int                               `json:"max_charges"`
	Rebuke            ItemTalismanOfPureGoodRebuke      `json:"rebuke"`
	LastCharge        string                            `json:"last_charge"`
	SourcePages       []int                             `json:"source_pages"`
}

// ItemTalismanOfPureGoodTouchDamage is a nested part of a source payload.
type ItemTalismanOfPureGoodTouchDamage struct {
	CreatureTypes []string `json:"creature_types"`
	Damage        ItemDice `json:"damage"`
	Repeats       string   `json:"repeats"`
}

// ItemTalismanOfPureGoodRebuke is a nested part of a source payload.
type ItemTalismanOfPureGoodRebuke struct {
	Action            string   `json:"action"`
	ChargeCost        int      `json:"charge_cost"`
	Target            string   `json:"target"`
	RangeFeet         int      `json:"range_feet"`
	SaveAbility       string   `json:"save_ability"`
	SaveDC            int      `json:"save_dc"`
	DisadvantageTypes []string `json:"disadvantage_types"`
	Failure           string   `json:"failure"`
	SuccessDamage     ItemDice `json:"success_damage"`
}

// ItemTalismanOfUltimateEvil retains typed SRD source facts without implementing them.
type ItemTalismanOfUltimateEvil struct {
	TouchDamage      ItemTalismanOfUltimateEvilTouchDamage `json:"touch_damage"`
	HolySymbol       bool                                  `json:"holy_symbol"`
	SpellAttackBonus int                                   `json:"spell_attack_bonus"`
	MaxCharges       int                                   `json:"max_charges"`
	UltimateEnd      ItemTalismanOfUltimateEvilUltimateEnd `json:"ultimate_end"`
	LastCharge       string                                `json:"last_charge"`
}

// ItemTalismanOfUltimateEvilTouchDamage is a nested part of a source payload.
type ItemTalismanOfUltimateEvilTouchDamage struct {
	ExcludedCreatureTypes []string `json:"excluded_creature_types"`
	Damage                ItemDice `json:"damage"`
	Repeats               string   `json:"repeats"`
}

// ItemTalismanOfUltimateEvilUltimateEnd is a nested part of a source payload.
type ItemTalismanOfUltimateEvilUltimateEnd struct {
	Action            string   `json:"action"`
	ChargeCost        int      `json:"charge_cost"`
	Target            string   `json:"target"`
	RangeFeet         int      `json:"range_feet"`
	SaveAbility       string   `json:"save_ability"`
	SaveDC            int      `json:"save_dc"`
	DisadvantageTypes []string `json:"disadvantage_types"`
	Failure           string   `json:"failure"`
	SuccessDamage     ItemDice `json:"success_damage"`
}

// ItemTalismanOfTheSphere retains typed SRD source facts without implementing them.
type ItemTalismanOfTheSphere struct {
	Condition string                       `json:"condition"`
	Check     ItemTalismanOfTheSphereCheck `json:"check"`
	Move      ItemTalismanOfTheSphereMove  `json:"move"`
}

// ItemTalismanOfTheSphereCheck is a nested part of a source payload.
type ItemTalismanOfTheSphereCheck struct {
	Ability string `json:"ability"`
	Skill   string `json:"skill"`
	Mode    string `json:"mode"`
	Purpose string `json:"purpose"`
}

// ItemTalismanOfTheSphereMove is a nested part of a source payload.
type ItemTalismanOfTheSphereMove struct {
	Timing                                string `json:"timing"`
	Action                                string `json:"action"`
	BaseFeet                              int    `json:"base_feet"`
	AdditionalFeetPerIntelligenceModifier int    `json:"additional_feet_per_intelligence_modifier"`
	StraightLineRequired                  bool   `json:"straight_line_required"`
}

// ItemThunderousGreatclub retains typed SRD source facts without implementing them.
type ItemThunderousGreatclub struct {
	Strength            ItemThunderousGreatclubStrength   `json:"strength"`
	CreatureExtraDamage ItemDice                          `json:"creature_extra_damage"`
	ObjectExtraDamage   ItemDice                          `json:"object_extra_damage"`
	ObjectFilter        string                            `json:"object_filter"`
	Clap                ItemThunderousGreatclubClap       `json:"clap"`
	Earthquake          ItemThunderousGreatclubEarthquake `json:"earthquake"`
}

// ItemThunderousGreatclubStrength is a nested part of a source payload.
type ItemThunderousGreatclubStrength struct {
	Score     int    `json:"score"`
	Condition string `json:"condition"`
}

// ItemThunderousGreatclubClap is a nested part of a source payload.
type ItemThunderousGreatclubClap struct {
	Action       string   `json:"action"`
	Strike       string   `json:"strike"`
	AudibleFeet  int      `json:"audible_feet"`
	ConeFeet     int      `json:"cone_feet"`
	SaveAbility  string   `json:"save_ability"`
	SaveDC       int      `json:"save_dc"`
	Condition    string   `json:"condition"`
	ObjectDamage ItemDice `json:"object_damage"`
}

// ItemThunderousGreatclubEarthquake is a nested part of a source payload.
type ItemThunderousGreatclubEarthquake struct {
	Action                   string                                   `json:"action"`
	Strike                   string                                   `json:"strike"`
	CircleRadiusFeet         int                                      `json:"circle_radius_feet"`
	StructureDamage          int                                      `json:"structure_damage"`
	StructureDamageType      string                                   `json:"structure_damage_type"`
	ProneSaveAbility         string                                   `json:"prone_save_ability"`
	ProneSaveDC              int                                      `json:"prone_save_dc"`
	Condition                string                                   `json:"condition"`
	ConcentrationSaveAbility string                                   `json:"concentration_save_ability"`
	ConcentrationSaveDC      int                                      `json:"concentration_save_dc"`
	Fissure                  ItemThunderousGreatclubEarthquakeFissure `json:"fissure"`
	ReuseAfter               string                                   `json:"reuse_after"`
}

// ItemThunderousGreatclubEarthquakeFissure is a nested part of a source payload.
type ItemThunderousGreatclubEarthquakeFissure struct {
	DepthFeet   int    `json:"depth_feet"`
	WidthFeet   int    `json:"width_feet"`
	SaveAbility string `json:"save_ability"`
	SaveDC      int    `json:"save_dc"`
	Failure     string `json:"failure"`
	Success     string `json:"success"`
	Structures  string `json:"structures"`
}

// ItemTomeOfClearThought retains typed SRD source facts without implementing them.
type ItemTomeOfClearThought struct {
	Contents                   string                      `json:"contents"`
	Study                      ItemTomeOfClearThoughtStudy `json:"study"`
	Ability                    string                      `json:"ability"`
	Increase                   int                         `json:"increase"`
	Maximum                    int                         `json:"maximum"`
	MagicAfterUse              string                      `json:"magic_after_use"`
	MagicReturnsAfterCenturies int                         `json:"magic_returns_after_centuries"`
}

// ItemTomeOfClearThoughtStudy is a nested part of a source payload.
type ItemTomeOfClearThoughtStudy struct {
	Hours      int    `json:"hours"`
	WithinDays int    `json:"within_days"`
	Activity   string `json:"activity"`
}

// ItemTomeOfLeadership retains typed SRD source facts without implementing them.
type ItemTomeOfLeadership struct {
	Contents                   string                    `json:"contents"`
	Study                      ItemTomeOfLeadershipStudy `json:"study"`
	Ability                    string                    `json:"ability"`
	Increase                   int                       `json:"increase"`
	Maximum                    int                       `json:"maximum"`
	MagicAfterUse              string                    `json:"magic_after_use"`
	MagicReturnsAfterCenturies int                       `json:"magic_returns_after_centuries"`
	SourcePages                []int                     `json:"source_pages"`
}

// ItemTomeOfLeadershipStudy is a nested part of a source payload.
type ItemTomeOfLeadershipStudy struct {
	Hours      int    `json:"hours"`
	WithinDays int    `json:"within_days"`
	Activity   string `json:"activity"`
}

// ItemTomeOfUnderstanding retains typed SRD source facts without implementing them.
type ItemTomeOfUnderstanding struct {
	Contents                   string                       `json:"contents"`
	Study                      ItemTomeOfUnderstandingStudy `json:"study"`
	Ability                    string                       `json:"ability"`
	Increase                   int                          `json:"increase"`
	Maximum                    int                          `json:"maximum"`
	MagicAfterUse              string                       `json:"magic_after_use"`
	MagicReturnsAfterCenturies int                          `json:"magic_returns_after_centuries"`
}

// ItemTomeOfUnderstandingStudy is a nested part of a source payload.
type ItemTomeOfUnderstandingStudy struct {
	Hours      int    `json:"hours"`
	WithinDays int    `json:"within_days"`
	Activity   string `json:"activity"`
}

// ItemTridentOfFishCommand retains typed SRD source facts without implementing them.
type ItemTridentOfFishCommand struct {
	MaxCharges   int      `json:"max_charges"`
	RechargeDice ItemDice `json:"recharge_dice"`
	Condition    string   `json:"condition"`
	ChargeCost   int      `json:"charge_cost"`
	Spell        string   `json:"spell"`
	SaveDC       int      `json:"save_dc"`
	Target       string   `json:"target"`
}

// ItemUniversalSolvent retains typed SRD source facts without implementing them.
type ItemUniversalSolvent struct {
	Appearance           string                        `json:"appearance"`
	OuncesWhenFound      ItemDice                      `json:"ounces_when_found"`
	OuncesWhenFoundBonus int                           `json:"ounces_when_found_bonus"`
	Action               string                        `json:"action"`
	MinimumOuncesPerPour int                           `json:"minimum_ounces_per_pour"`
	Target               string                        `json:"target"`
	Dissolves            ItemUniversalSolventDissolves `json:"dissolves"`
}

// ItemUniversalSolventDissolves is a nested part of a source payload.
type ItemUniversalSolventDissolves struct {
	MaximumSquareFeet int    `json:"maximum_square_feet"`
	Timing            string `json:"timing"`
	Includes          string `json:"includes"`
}

// ItemViciousWeapon retains typed SRD source facts without implementing them.
type ItemViciousWeapon struct {
	Trigger        string   `json:"trigger"`
	ExtraDamage    ItemDice `json:"extra_damage"`
	DamageTypeRule string   `json:"damage_type_rule"`
}

// ItemVorpalSword retains typed SRD source facts without implementing them.
type ItemVorpalSword struct {
	AttackBonus         int                 `json:"attack_bonus"`
	DamageBonus         int                 `json:"damage_bonus"`
	IgnoresResistanceTo string              `json:"ignores_resistance_to"`
	Head                ItemVorpalSwordHead `json:"head"`
}

// ItemVorpalSwordHead is a nested part of a source payload.
type ItemVorpalSwordHead struct {
	Target                string                                 `json:"target"`
	TriggerNaturalRoll    int                                    `json:"trigger_natural_roll"`
	Effect                string                                 `json:"effect"`
	Death                 string                                 `json:"death"`
	ImmuneWhen            []string                               `json:"immune_when"`
	ImmuneExtraDamage     int                                    `json:"immune_extra_damage"`
	ImmuneExtraDamageType string                                 `json:"immune_extra_damage_type"`
	LegendaryResistance   ItemVorpalSwordHeadLegendaryResistance `json:"legendary_resistance"`
}

// ItemVorpalSwordHeadLegendaryResistance is a nested part of a source payload.
type ItemVorpalSwordHeadLegendaryResistance struct {
	Cost   string `json:"cost"`
	Result string `json:"result"`
}

// ItemWandOfBinding retains typed SRD source facts without implementing them.
type ItemWandOfBinding struct {
	MaxCharges int                         `json:"max_charges"`
	SaveDC     int                         `json:"save_dc"`
	Spells     []ItemWandOfBindingSpell    `json:"spells"`
	Recharge   ItemWandOfBindingRecharge   `json:"recharge"`
	LastCharge ItemWandOfBindingLastCharge `json:"last_charge"`
}

// ItemWandOfBindingSpell is a nested part of a source payload.
type ItemWandOfBindingSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemWandOfBindingRecharge is a nested part of a source payload.
type ItemWandOfBindingRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemWandOfBindingLastCharge is a nested part of a source payload.
type ItemWandOfBindingLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// itemManifest244 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest244 = itemSourceManifest{Batch: "data-magic-items-20", Issue: 244, Families: []itemSourceManifestFamily{
	{Slug: "talisman-of-pure-good", Name: "Talisman of Pure Good", Page: 248, Category: "wondrous-item", Attunement: "required", Equip: "wear-or-hold", Use: "magic-action", UseRestrictions: "expend-one-charge target-on-ground-within-120-feet",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "none"},
		Keys:    "item.attunement attunement.by-class damage.8d6-radiant-on-fiend-or-undead-touch item.holy-symbol bonus.spell-attack-2 item.charges action.magic-pure-rebuke save.dexterity-dc-20 effect.destroyed-in-fissure damage.4d6-psychic-on-success item.destroyed-after-last-charge", Payload: "talisman_good",
		Variants: []itemSourceManifestVariant{{"standard", "Talisman of Pure Good", "legendary"}}},
	{Slug: "talisman-of-ultimate-evil", Name: "Talisman of Ultimate Evil", Page: 249, Category: "wondrous-item", Attunement: "required", Equip: "wear-or-hold", Use: "magic-action", UseRestrictions: "expend-one-charge target-on-ground-within-120-feet",
		Charges: &ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "none"},
		Keys:    "item.attunement damage.8d6-necrotic-on-other-creature-touch item.holy-symbol bonus.spell-attack-2 item.charges action.magic-ultimate-end save.dexterity-dc-20 effect.destroyed-in-fissure damage.4d6-psychic-on-success item.destroyed-after-last-charge", Payload: "talisman_evil",
		Variants: []itemSourceManifestVariant{{"standard", "Talisman of Ultimate Evil", "legendary"}}},
	{Slug: "talisman-of-the-sphere", Name: "Talisman of the Sphere", Page: 249, Category: "wondrous-item", Attunement: "required", Equip: "wear-or-hold", Use: "magic-action", UseRestrictions: "controlling-sphere-of-annihilation-at-turn-start",
		Keys: "item.attunement check.advantage-arcana-to-control-sphere action.magic-move-sphere movement.10-feet-plus-10-times-intelligence-modifier", Payload: "talisman_sphere",
		Variants: []itemSourceManifestVariant{{"standard", "Talisman of the Sphere", "legendary"}}},
	{Slug: "thunderous-greatclub", Name: "Thunderous Greatclub", Page: 249, Category: "weapon", Options: "greatclub", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "earthquake-once-until-next-dawn",
		Keys: "item.attunement item.magic-weapon ability.strength-20 damage.extra-1d8-thunder damage.extra-3d8-thunder-to-objects action.magic-clap-of-thunder save.strength-dc-15 condition.prone damage.3d8-thunder-to-objects-in-cone action.magic-earthquake damage.50-bludgeoning-to-structures save.dexterity-dc-20 save.constitution-dc-20-concentration effect.fissure item.earthquake-once-until-dawn", Payload: "thunderous_greatclub",
		Variants: []itemSourceManifestVariant{{"standard", "Thunderous Greatclub", "very-rare"}}},
	{Slug: "tome-of-clear-thought", Name: "Tome of Clear Thought", Page: 249, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "48-hours-within-6-days",
		Keys: "item.study-48-hours-within-6-days ability.increase-intelligence-by-2 ability.maximum-30 item.loses-magic-for-a-century", Payload: "clear_thought",
		Variants: []itemSourceManifestVariant{{"standard", "Tome of Clear Thought", "very-rare"}}},
	{Slug: "tome-of-leadership-and-influence", Name: "Tome of Leadership and Influence", Page: 249, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "48-hours-within-6-days",
		Keys: "item.study-48-hours-within-6-days ability.increase-charisma-by-2 ability.maximum-30 item.loses-magic-for-a-century", Payload: "leadership",
		Variants: []itemSourceManifestVariant{{"standard", "Tome of Leadership and Influence", "very-rare"}}},
	{Slug: "tome-of-understanding", Name: "Tome of Understanding", Page: 250, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "48-hours-within-6-days",
		Keys: "item.study-48-hours-within-6-days ability.increase-wisdom-by-2 ability.maximum-30 item.loses-magic-for-a-century", Payload: "understanding",
		Variants: []itemSourceManifestVariant{{"standard", "Tome of Understanding", "very-rare"}}},
	{Slug: "trident-of-fish-command", Name: "Trident of Fish Command", Page: 250, Category: "weapon", Options: "trident", Attunement: "required", Equip: "carry", Use: "cast-spell", UseRestrictions: "expend-one-charge beast-with-swim-speed",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.magic-weapon item.carry item.charges spell.cast-dominate-beast spell.save-dc-15 spell.target-beast-with-swim-speed", Payload: "fish_command",
		Variants: []itemSourceManifestVariant{{"standard", "Trident of Fish Command", "uncommon"}}},
	{Slug: "universal-solvent", Name: "Universal Solvent", Page: 250, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "pour-onto-surface-within-reach",
		Keys: "item.found-with-1d6-plus-1-ounces action.utilize effect.dissolve-adhesive-1-square-foot-per-ounce", Payload: "solvent",
		Variants: []itemSourceManifestVariant{{"standard", "Universal Solvent", "legendary"}}},
	{Slug: "vicious-weapon", Name: "Vicious Weapon", Page: 250, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "none", Equip: "wield", Use: "weapon-attack",
		Keys: "item.magic-weapon damage.extra-2d6-of-weapon-type", Payload: "vicious",
		Variants: []itemSourceManifestVariant{{"standard", "Vicious Weapon", "rare"}}},
	{Slug: "vorpal-sword", Name: "Vorpal Sword", Page: 250, Category: "weapon", Options: "glaive greatsword longsword scimitar", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon bonus.attack-and-damage-3 damage.ignore-slashing-resistance effect.sever-head-on-natural-20 damage.extra-30-slashing-when-head-kept resistance.legendary-resistance-avoids-head-loss", Payload: "vorpal",
		Variants: []itemSourceManifestVariant{{"standard", "Vorpal Sword", "legendary"}}},
	{Slug: "wand-of-binding", Name: "Wand of Binding", Page: 250, Category: "wand", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 1, Recharge: "regain-1d6-plus-1-at-dawn"},
		Keys:    "item.attunement item.hold item.charges spell.cast-listed-spells spell.save-dc-17 item.destroyed-on-1-after-last-charge", Payload: "binding",
		Variants: []itemSourceManifestVariant{{"standard", "Wand of Binding", "rare"}}},
}}

// LoadItems244 returns a fresh validated copy of batch #244. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems244() (ItemSourceBatch[ItemDetails244], error) { return DecodeItems244(items244JSON) }

// DecodeItems244 validates supplied JSON as batch #244 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems244(data []byte) (ItemSourceBatch[ItemDetails244], error) {
	return decodeItemSourceBatch(data, itemManifest244, validateItemDetails244)
}

// validateItemDetails244 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails244(family, _ string, d *ItemDetails244) error {
	ok := true
	switch family {
	case "talisman-of-pure-good":
		p := d.TalismanGood
		ok = p.MaxCharges == 7 && p.SpellAttackBonus == 2 && p.TouchDamage.Damage == ItemDice{Dice: 8, Sides: 6, Type: "radiant"} && p.Rebuke.SaveDC == 20 && p.Rebuke.RangeFeet == 120 && p.Rebuke.SuccessDamage == ItemDice{Dice: 4, Sides: 6, Type: "psychic"} && slices.Equal(p.AttunementByAnyOf, []string{"cleric", "paladin"}) && slices.Equal(p.Rebuke.DisadvantageTypes, []string{"fiend", "undead"}) && slices.Equal(p.SourcePages, []int{248, 249})
	case "talisman-of-ultimate-evil":
		p := d.TalismanEvil
		ok = p.MaxCharges == 6 && p.SpellAttackBonus == 2 && p.TouchDamage.Damage == ItemDice{Dice: 8, Sides: 6, Type: "necrotic"} && p.UltimateEnd.SaveDC == 20 && p.UltimateEnd.RangeFeet == 120 && p.UltimateEnd.SuccessDamage == ItemDice{Dice: 4, Sides: 6, Type: "psychic"} && slices.Equal(p.UltimateEnd.DisadvantageTypes, []string{"celestial"}) && slices.Equal(p.TouchDamage.ExcludedCreatureTypes, []string{"fiend", "undead"})
	case "talisman-of-the-sphere":
		p := d.TalismanSphere
		ok = p.Check.Skill == "arcana" && p.Move.BaseFeet == 10 && p.Move.AdditionalFeetPerIntelligenceModifier == 10 && !p.Move.StraightLineRequired
	case "thunderous-greatclub":
		p := d.ThunderousGreatclub
		ok = p.Strength.Score == 20 && p.CreatureExtraDamage == ItemDice{Dice: 1, Sides: 8, Type: "thunder"} && p.ObjectExtraDamage == ItemDice{Dice: 3, Sides: 8, Type: "thunder"} && p.Clap.SaveDC == 15 && p.Clap.ConeFeet == 30 && p.Earthquake.StructureDamage == 50 && p.Earthquake.ProneSaveDC == 20 && p.Earthquake.ConcentrationSaveDC == 20 && p.Earthquake.Fissure.SaveDC == 20
	case "tome-of-clear-thought":
		p := d.ClearThought
		ok = p.Study.Hours == 48 && p.Study.WithinDays == 6 && p.Ability == "intelligence" && p.Increase == 2 && p.Maximum == 30 && p.MagicReturnsAfterCenturies == 1
	case "tome-of-leadership-and-influence":
		p := d.Leadership
		ok = p.Study.Hours == 48 && p.Study.WithinDays == 6 && p.Ability == "charisma" && p.Increase == 2 && p.Maximum == 30 && p.MagicReturnsAfterCenturies == 1 && slices.Equal(p.SourcePages, []int{249, 250})
	case "tome-of-understanding":
		p := d.Understanding
		ok = p.Study.Hours == 48 && p.Study.WithinDays == 6 && p.Ability == "wisdom" && p.Increase == 2 && p.Maximum == 30 && p.MagicReturnsAfterCenturies == 1
	case "trident-of-fish-command":
		p := d.FishCommand
		ok = p.MaxCharges == 3 && p.SaveDC == 15 && p.ChargeCost == 1 && p.Spell == "dominate-beast" && p.RechargeDice == ItemDice{Dice: 1, Sides: 3, Type: "charges"}
	case "universal-solvent":
		p := d.Solvent
		ok = p.OuncesWhenFound == ItemDice{Dice: 1, Sides: 6, Type: "ounces"} && p.OuncesWhenFoundBonus == 1 && p.Dissolves.MaximumSquareFeet == 1 && p.Dissolves.Includes == "sovereign-glue"
	case "vicious-weapon":
		p := d.Vicious
		ok = p.ExtraDamage == ItemDice{Dice: 2, Sides: 6, Type: "same-as-weapon"} && p.Trigger == "hit-creature"
	case "vorpal-sword":
		p := d.Vorpal
		ok = p.AttackBonus == 3 && p.DamageBonus == 3 && p.IgnoresResistanceTo == "slashing" && p.Head.TriggerNaturalRoll == 20 && p.Head.ImmuneExtraDamage == 30 && len(p.Head.ImmuneWhen) == 3 && uniqueItemKeys(p.Head.ImmuneWhen)
	case "wand-of-binding":
		p := d.Binding
		ok = p.MaxCharges == 7 && p.SaveDC == 17 && slices.Equal(p.Spells, []ItemWandOfBindingSpell{{Spell: "hold-monster", ChargeCost: 5}, {Spell: "hold-person", ChargeCost: 2}}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 1 && p.LastCharge.DestroyedOn == 1
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
