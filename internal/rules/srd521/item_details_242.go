package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_242.json
var items242JSON []byte

// ItemDetails242 holds exactly one typed payload per variant of catalogue batch
// #242 (magic items 18, SRD 5.2.1 pp. 243-246). Payloads are source facts only.
type ItemDetails242 struct {
	ShieldBonus    *ItemShield12Or3              `json:"shield_bonus,omitempty"`
	SpiderSlippers *ItemSlippersOfSpiderClimbing `json:"spider_slippers,omitempty"`
	Glue           *ItemSovereignGlue            `json:"glue,omitempty"`
	Scroll         *ItemSpellScroll              `json:"scroll,omitempty"`
	Spellguard     *ItemSpellguardShield         `json:"spellguard,omitempty"`
	Annihilation   *ItemSphereOfAnnihilation     `json:"annihilation,omitempty"`
	Charming       *ItemStaffOfCharming          `json:"charming,omitempty"`
	Fire           *ItemStaffOfFire              `json:"fire,omitempty"`
	Frost          *ItemStaffOfFrost             `json:"frost,omitempty"`
	Healing        *ItemStaffOfHealing           `json:"healing,omitempty"`
	Power          *ItemStaffOfPower             `json:"power,omitempty"`
	Striking       *ItemStaffOfStriking          `json:"striking,omitempty"`
}

// ItemShield12Or3 retains typed SRD source facts without implementing them.
type ItemShield12Or3 struct {
	Condition             string `json:"condition"`
	ArmorClassBonus       int    `json:"armor_class_bonus"`
	StacksWithShieldBonus bool   `json:"stacks_with_shield_bonus"`
}

// ItemSlippersOfSpiderClimbing retains typed SRD source facts without implementing them.
type ItemSlippersOfSpiderClimbing struct {
	Condition        string   `json:"condition"`
	Movement         []string `json:"movement"`
	HandsFree        bool     `json:"hands_free"`
	ClimbSpeed       string   `json:"climb_speed"`
	ExcludedSurfaces string   `json:"excluded_surfaces"`
}

// ItemSovereignGlue retains typed SRD source facts without implementing them.
type ItemSovereignGlue struct {
	Appearance         string   `json:"appearance"`
	Bond               string   `json:"bond"`
	Storage            string   `json:"storage"`
	FoundOunces        ItemDice `json:"found_ounces"`
	FoundOuncesBonus   int      `json:"found_ounces_bonus"`
	CoverageOunces     int      `json:"coverage_ounces"`
	CoverageSquareFeet int      `json:"coverage_square_feet"`
	ApplyAction        string   `json:"apply_action"`
	SetMinutes         int      `json:"set_minutes"`
	BrokenOnlyBy       []string `json:"broken_only_by"`
}

// ItemSpellScroll retains typed SRD source facts without implementing them.
type ItemSpellScroll struct {
	SpellLevel                 int                             `json:"spell_level"`
	Rarity                     string                          `json:"rarity"`
	SaveDC                     int                             `json:"save_dc"`
	AttackBonus                int                             `json:"attack_bonus"`
	Bears                      string                          `json:"bears"`
	ReadableIf                 string                          `json:"readable_if"`
	MaterialComponentsRequired bool                            `json:"material_components_required"`
	Otherwise                  string                          `json:"otherwise"`
	CastingTime                string                          `json:"casting_time"`
	AfterCasting               string                          `json:"after_casting"`
	LostIfCastingInterrupted   bool                            `json:"lost_if_casting_interrupted"`
	HigherLevelCheck           ItemSpellScrollHigherLevelCheck `json:"higher_level_check"`
	Copying                    ItemSpellScrollCopying          `json:"copying"`
}

// ItemSpellScrollHigherLevelCheck is a nested part of a source payload.
type ItemSpellScrollHigherLevelCheck struct {
	Trigger string `json:"trigger"`
	Ability string `json:"ability"`
	DCBase  int    `json:"dc_base"`
	DCAdds  string `json:"dc_adds"`
	Failure string `json:"failure"`
}

// ItemSpellScrollCopying is a nested part of a source payload.
type ItemSpellScrollCopying struct {
	Eligible                 string `json:"eligible"`
	CheckAbility             string `json:"check_ability"`
	CheckSkill               string `json:"check_skill"`
	DCBase                   int    `json:"dc_base"`
	DCAdds                   string `json:"dc_adds"`
	Success                  string `json:"success"`
	ScrollDestroyedEitherWay bool   `json:"scroll_destroyed_either_way"`
}

// ItemSpellguardShield retains typed SRD source facts without implementing them.
type ItemSpellguardShield struct {
	Condition              string   `json:"condition"`
	SaveMode               string   `json:"save_mode"`
	Against                []string `json:"against"`
	SpellAttacksAgainstYou string   `json:"spell_attacks_against_you"`
}

// ItemSphereOfAnnihilation retains typed SRD source facts without implementing them.
type ItemSphereOfAnnihilation struct {
	DiameterFeet        int                                        `json:"diameter_feet"`
	Nature              string                                     `json:"nature"`
	Obliterates         string                                     `json:"obliterates"`
	ArtifactException   string                                     `json:"artifact_exception"`
	TouchDamage         ItemDice                                   `json:"touch_damage"`
	Control             ItemSphereOfAnnihilationControl            `json:"control"`
	Move                ItemSphereOfAnnihilationMove               `json:"move"`
	InteractionTriggers []string                                   `json:"interaction_triggers"`
	TableDieSides       int                                        `json:"table_die_sides"`
	InteractionTable    []ItemSphereOfAnnihilationInteractionTable `json:"interaction_table"`
	RiftRadiusFeet      int                                        `json:"rift_radius_feet"`
	SourcePages         []int                                      `json:"source_pages"`
}

// ItemSphereOfAnnihilationControl is a nested part of a source payload.
type ItemSphereOfAnnihilationControl struct {
	StationaryUntilControlled bool   `json:"stationary_until_controlled"`
	RangeFeet                 int    `json:"range_feet"`
	Action                    string `json:"action"`
	CheckAbility              string `json:"check_ability"`
	CheckSkill                string `json:"check_skill"`
	DC                        int    `json:"dc"`
	LastsUntil                string `json:"lasts_until"`
	OtherController           string `json:"other_controller"`
	FailureMoveFeet           int    `json:"failure_move_feet"`
}

// ItemSphereOfAnnihilationMove is a nested part of a source payload.
type ItemSphereOfAnnihilationMove struct {
	Action                      string   `json:"action"`
	FeetPerIntelligenceModifier int      `json:"feet_per_intelligence_modifier"`
	MinimumFeet                 int      `json:"minimum_feet"`
	SaveAbility                 string   `json:"save_ability"`
	SaveDC                      int      `json:"save_dc"`
	Damage                      ItemDice `json:"damage"`
	AtZeroHitPoints             string   `json:"at_zero_hit_points"`
}

// ItemSphereOfAnnihilationInteractionTable is a nested part of a source payload.
type ItemSphereOfAnnihilationInteractionTable struct {
	Min    int    `json:"min"`
	Max    int    `json:"max"`
	Result string `json:"result"`
}

// ItemStaffOfCharming retains typed SRD source facts without implementing them.
type ItemStaffOfCharming struct {
	AttunementByAnyOf  []string                              `json:"attunement_by_any_of"`
	Condition          string                                `json:"condition"`
	MaxCharges         int                                   `json:"max_charges"`
	CastSpell          ItemStaffOfCharmingCastSpell          `json:"cast_spell"`
	ReflectEnchantment ItemStaffOfCharmingReflectEnchantment `json:"reflect_enchantment"`
	ResistEnchantment  ItemStaffOfCharmingResistEnchantment  `json:"resist_enchantment"`
	Recharge           ItemStaffOfCharmingRecharge           `json:"recharge"`
	LastCharge         ItemStaffOfCharmingLastCharge         `json:"last_charge"`
}

// ItemStaffOfCharmingCastSpell is a nested part of a source payload.
type ItemStaffOfCharmingCastSpell struct {
	ChargeCost   int      `json:"charge_cost"`
	Spells       []string `json:"spells"`
	SaveDCSource string   `json:"save_dc_source"`
}

// ItemStaffOfCharmingReflectEnchantment is a nested part of a source payload.
type ItemStaffOfCharmingReflectEnchantment struct {
	Trigger    string `json:"trigger"`
	Action     string `json:"action"`
	ChargeCost int    `json:"charge_cost"`
	Result     string `json:"result"`
}

// ItemStaffOfCharmingResistEnchantment is a nested part of a source payload.
type ItemStaffOfCharmingResistEnchantment struct {
	Trigger    string `json:"trigger"`
	Result     string `json:"result"`
	ReuseAfter string `json:"reuse_after"`
}

// ItemStaffOfCharmingRecharge is a nested part of a source payload.
type ItemStaffOfCharmingRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfCharmingLastCharge is a nested part of a source payload.
type ItemStaffOfCharmingLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStaffOfFire retains typed SRD source facts without implementing them.
type ItemStaffOfFire struct {
	AttunementByAnyOf []string                  `json:"attunement_by_any_of"`
	Condition         string                    `json:"condition"`
	Resistance        string                    `json:"resistance"`
	MaxCharges        int                       `json:"max_charges"`
	SaveDCSource      string                    `json:"save_dc_source"`
	Spells            []ItemStaffOfFireSpell    `json:"spells"`
	Recharge          ItemStaffOfFireRecharge   `json:"recharge"`
	LastCharge        ItemStaffOfFireLastCharge `json:"last_charge"`
}

// ItemStaffOfFireSpell is a nested part of a source payload.
type ItemStaffOfFireSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemStaffOfFireRecharge is a nested part of a source payload.
type ItemStaffOfFireRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfFireLastCharge is a nested part of a source payload.
type ItemStaffOfFireLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStaffOfFrost retains typed SRD source facts without implementing them.
type ItemStaffOfFrost struct {
	AttunementByAnyOf []string                   `json:"attunement_by_any_of"`
	Condition         string                     `json:"condition"`
	Resistance        string                     `json:"resistance"`
	MaxCharges        int                        `json:"max_charges"`
	SaveDCSource      string                     `json:"save_dc_source"`
	Spells            []ItemStaffOfFrostSpell    `json:"spells"`
	Recharge          ItemStaffOfFrostRecharge   `json:"recharge"`
	LastCharge        ItemStaffOfFrostLastCharge `json:"last_charge"`
}

// ItemStaffOfFrostSpell is a nested part of a source payload.
type ItemStaffOfFrostSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemStaffOfFrostRecharge is a nested part of a source payload.
type ItemStaffOfFrostRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfFrostLastCharge is a nested part of a source payload.
type ItemStaffOfFrostLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStaffOfHealing retains typed SRD source facts without implementing them.
type ItemStaffOfHealing struct {
	AttunementByAnyOf   []string                     `json:"attunement_by_any_of"`
	Condition           string                       `json:"condition"`
	MaxCharges          int                          `json:"max_charges"`
	SpellModifierSource string                       `json:"spell_modifier_source"`
	CureWounds          ItemStaffOfHealingCureWounds `json:"cure_wounds"`
	Spells              []ItemStaffOfHealingSpell    `json:"spells"`
	Recharge            ItemStaffOfHealingRecharge   `json:"recharge"`
	LastCharge          ItemStaffOfHealingLastCharge `json:"last_charge"`
}

// ItemStaffOfHealingCureWounds is a nested part of a source payload.
type ItemStaffOfHealingCureWounds struct {
	ChargesPerSpellLevel int `json:"charges_per_spell_level"`
	MaximumSpellLevel    int `json:"maximum_spell_level"`
}

// ItemStaffOfHealingSpell is a nested part of a source payload.
type ItemStaffOfHealingSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemStaffOfHealingRecharge is a nested part of a source payload.
type ItemStaffOfHealingRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfHealingLastCharge is a nested part of a source payload.
type ItemStaffOfHealingLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemStaffOfPower retains typed SRD source facts without implementing them.
type ItemStaffOfPower struct {
	AttunementByAnyOf     []string                          `json:"attunement_by_any_of"`
	MaxCharges            int                               `json:"max_charges"`
	WieldedAs             string                            `json:"wielded_as"`
	AttackBonus           int                               `json:"attack_bonus"`
	DamageBonus           int                               `json:"damage_bonus"`
	HoldingBonus          int                               `json:"holding_bonus"`
	HoldingBonusAppliesTo []string                          `json:"holding_bonus_applies_to"`
	SaveDCSource          string                            `json:"save_dc_source"`
	Spells                []ItemStaffOfPowerSpell           `json:"spells"`
	Recharge              ItemStaffOfPowerRecharge          `json:"recharge"`
	LastCharge            ItemStaffOfPowerLastCharge        `json:"last_charge"`
	RetributiveStrike     ItemStaffOfPowerRetributiveStrike `json:"retributive_strike"`
	SourcePages           []int                             `json:"source_pages"`
}

// ItemStaffOfPowerSpell is a nested part of a source payload.
type ItemStaffOfPowerSpell struct {
	Spell        string `json:"spell"`
	PrintedLevel string `json:"printed_level"`
	ChargeCost   int    `json:"charge_cost"`
}

// ItemStaffOfPowerRecharge is a nested part of a source payload.
type ItemStaffOfPowerRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfPowerLastCharge is a nested part of a source payload.
type ItemStaffOfPowerLastCharge struct {
	DieSides          int      `json:"die_sides"`
	LosesPropertiesOn int      `json:"loses_properties_on"`
	Retains           string   `json:"retains"`
	RegainsOn         int      `json:"regains_on"`
	Regains           ItemDice `json:"regains"`
	RegainsBonus      int      `json:"regains_bonus"`
}

// ItemStaffOfPowerRetributiveStrike is a nested part of a source payload.
type ItemStaffOfPowerRetributiveStrike struct {
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

// ItemStaffOfStriking retains typed SRD source facts without implementing them.
type ItemStaffOfStriking struct {
	WieldedAs            string                        `json:"wielded_as"`
	AttackBonus          int                           `json:"attack_bonus"`
	DamageBonus          int                           `json:"damage_bonus"`
	MaxCharges           int                           `json:"max_charges"`
	Trigger              string                        `json:"trigger"`
	MaximumChargesPerHit int                           `json:"maximum_charges_per_hit"`
	ExtraDamagePerCharge ItemDice                      `json:"extra_damage_per_charge"`
	Recharge             ItemStaffOfStrikingRecharge   `json:"recharge"`
	LastCharge           ItemStaffOfStrikingLastCharge `json:"last_charge"`
}

// ItemStaffOfStrikingRecharge is a nested part of a source payload.
type ItemStaffOfStrikingRecharge struct {
	Dice  ItemDice `json:"dice"`
	Bonus int      `json:"bonus"`
	When  string   `json:"when"`
}

// ItemStaffOfStrikingLastCharge is a nested part of a source payload.
type ItemStaffOfStrikingLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// itemManifest242 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest242 = itemSourceManifest{Batch: "data-magic-items-18", Issue: 242, Families: []itemSourceManifestFamily{
	{Slug: "shield-1-2-or-3", Name: "Shield, +1, +2, or +3", Page: 243, Category: "armor", Options: "shield", Attunement: "none", Equip: "hold", Use: "none",
		Keys: "item.hold defense.armor-class-bonus-by-rarity", Payload: "shield_bonus",
		Variants: []itemSourceManifestVariant{
			{"plus-1", "Shield, +1", "uncommon"},
			{"plus-2", "Shield, +2", "rare"},
			{"plus-3", "Shield, +3", "very-rare"},
		}},
	{Slug: "slippers-of-spider-climbing", Name: "Slippers of Spider Climbing", Page: 244, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear movement.walk-walls-and-ceilings-hands-free movement.climb-speed-equals-speed movement.not-on-slippery-surfaces", Payload: "spider_slippers",
		Variants: []itemSourceManifestVariant{{"standard", "Slippers of Spider Climbing", "uncommon"}}},
	{Slug: "sovereign-glue", Name: "Sovereign Glue", Page: 244, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "apply-one-ounce",
		Keys: "action.utilize item.permanent-bond item.1-ounce-covers-1-foot-square duration.sets-after-1-minute item.bond-broken-only-by-listed-means", Payload: "glue",
		Variants: []itemSourceManifestVariant{{"standard", "Sovereign Glue", "legendary"}}},
	{Slug: "spell-scroll", Name: "Spell Scroll", Page: 244, Gear: &itemSourceGearListing{Page: 99, PriceGP: map[string]int{"cantrip": 30, "level-1": 50}}, Category: "scroll", Attunement: "none", Equip: "none", Use: "cast-spell", UseRestrictions: "spell-on-your-spell-list",
		Keys: "item.consume spell.cast-from-scroll spell.must-be-on-your-spell-list spell.no-material-components check.spellcasting-ability-dc-10-plus-level table.spell-level-rarity-save-dc-attack-bonus spell.copy-into-spellbook", Payload: "scroll",
		Variants: []itemSourceManifestVariant{
			{"cantrip", "Spell Scroll (cantrip)", "common"},
			{"level-1", "Spell Scroll (level 1)", "common"},
			{"level-2", "Spell Scroll (level 2)", "uncommon"},
			{"level-3", "Spell Scroll (level 3)", "uncommon"},
			{"level-4", "Spell Scroll (level 4)", "rare"},
			{"level-5", "Spell Scroll (level 5)", "rare"},
			{"level-6", "Spell Scroll (level 6)", "very-rare"},
			{"level-7", "Spell Scroll (level 7)", "very-rare"},
			{"level-8", "Spell Scroll (level 8)", "very-rare"},
			{"level-9", "Spell Scroll (level 9)", "legendary"},
		}},
	{Slug: "spellguard-shield", Name: "Spellguard Shield", Page: 244, Category: "armor", Options: "shield", Attunement: "required", Equip: "hold", Use: "none",
		Keys: "item.attunement item.hold save.magic-advantage attack.spell-attacks-disadvantage-against-you", Payload: "spellguard",
		Variants: []itemSourceManifestVariant{{"standard", "Spellguard Shield", "very-rare"}}},
	{Slug: "sphere-of-annihilation", Name: "Sphere of Annihilation", Page: 244, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "within-60-feet",
		Keys: "effect.obliterates-matter damage.force-8d10 check.intelligence-arcana-dc-25-control action.magic action.bonus-move-sphere save.dexterity-dc-19 table.planar-interaction-1d100", Payload: "annihilation",
		Variants: []itemSourceManifestVariant{{"standard", "Sphere of Annihilation", "legendary"}}},
	{Slug: "staff-of-charming", Name: "Staff of Charming", Page: 245, Category: "staff", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d8-plus-2-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.hold item.charges spell.cast-charm-person spell.cast-command spell.cast-comprehend-languages action.reaction-reflect-enchantment save.resist-enchantment-once-until-dawn item.destroyed-on-1-after-last-charge", Payload: "charming",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Charming", "rare"}}},
	{Slug: "staff-of-fire", Name: "Staff of Fire", Page: 245, Category: "staff", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.hold defense.resistance-fire item.charges spell.cast-listed-spells spell.uses-your-save-dc item.destroyed-on-1-after-last-charge", Payload: "fire",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Fire", "very-rare"}}},
	{Slug: "staff-of-frost", Name: "Staff of Frost", Page: 245, Category: "staff", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.hold defense.resistance-cold item.charges spell.cast-listed-spells spell.uses-your-save-dc item.destroyed-on-1-after-last-charge", Payload: "frost",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Frost", "very-rare"}}},
	{Slug: "staff-of-healing", Name: "Staff of Healing", Page: 245, Category: "staff", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.hold item.charges spell.cast-cure-wounds-per-level spell.cast-lesser-restoration spell.cast-mass-cure-wounds spell.uses-your-spellcasting-modifier item.destroyed-on-1-after-last-charge", Payload: "healing",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Healing", "rare"}}},
	{Slug: "staff-of-power", Name: "Staff of Power", Page: 245, Category: "staff", Attunement: "required", Equip: "wield", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 20, PerUse: 1, Recharge: "regain-2d8-plus-4-at-dawn"},
		Keys:    "item.attunement attunement.by-class item.magic-weapon bonus.attack-and-damage-2 bonus.armor-class-saves-and-spell-attack-2 item.charges spell.cast-listed-spells spell.uses-your-save-dc item.last-charge-roll item.retributive-strike damage.force-by-charges save.dexterity-dc-17", Payload: "power",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Power", "very-rare"}}},
	{Slug: "staff-of-striking", Name: "Staff of Striking", Page: 246, Category: "staff", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "expend-up-to-three-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 1, Recharge: "regain-1d6-plus-4-at-dawn"},
		Keys:    "item.attunement item.magic-weapon bonus.attack-and-damage-3 item.charges damage.force-1d6-per-charge item.nonmagical-on-1-after-last-charge", Payload: "striking",
		Variants: []itemSourceManifestVariant{{"standard", "Staff of Striking", "very-rare"}}},
}}

// LoadItems242 returns a fresh validated copy of batch #242. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems242() (ItemSourceBatch[ItemDetails242], error) { return DecodeItems242(items242JSON) }

// DecodeItems242 validates supplied JSON as batch #242 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems242(data []byte) (ItemSourceBatch[ItemDetails242], error) {
	return decodeItemSourceBatch(data, itemManifest242, validateItemDetails242)
}

// validateItemDetails242 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails242(family, variant string, d *ItemDetails242) error {
	ok := true
	switch family {
	case "shield-1-2-or-3":
		p := d.ShieldBonus
		bonus := map[string]int{"plus-1": 1, "plus-2": 2, "plus-3": 3}
		ok = p.ArmorClassBonus == bonus[variant] && p.StacksWithShieldBonus
	case "slippers-of-spider-climbing":
		p := d.SpiderSlippers
		ok = p.HandsFree && p.ClimbSpeed == "equal-to-speed" && uniqueItemKeys(p.Movement)
	case "sovereign-glue":
		p := d.Glue
		ok = p.FoundOunces == ItemDice{Dice: 1, Sides: 6, Type: "ounces"} && p.FoundOuncesBonus == 1 && p.CoverageOunces == 1 && p.SetMinutes == 1 && uniqueItemKeys(p.BrokenOnlyBy) && len(p.BrokenOnlyBy) == 3
	case "spell-scroll":
		p := d.Scroll
		rows := map[string][3]int{
			"cantrip": {0, 13, 5}, "level-1": {1, 13, 5}, "level-2": {2, 13, 5}, "level-3": {3, 15, 7}, "level-4": {4, 15, 7},
			"level-5": {5, 17, 9}, "level-6": {6, 17, 9}, "level-7": {7, 18, 10}, "level-8": {8, 18, 10}, "level-9": {9, 19, 11},
		}
		rarities := map[string]string{
			"cantrip": "common", "level-1": "common", "level-2": "uncommon", "level-3": "uncommon", "level-4": "rare",
			"level-5": "rare", "level-6": "very-rare", "level-7": "very-rare", "level-8": "very-rare", "level-9": "legendary",
		}
		ok = [3]int{p.SpellLevel, p.SaveDC, p.AttackBonus} == rows[variant] && p.Rarity == rarities[variant] && p.HigherLevelCheck.DCBase == 10 && p.Copying.DCBase == 10 && !p.MaterialComponentsRequired
	case "spellguard-shield":
		p := d.Spellguard
		ok = p.SaveMode == "advantage" && p.SpellAttacksAgainstYou == "disadvantage" && uniqueItemKeys(p.Against)
	case "sphere-of-annihilation":
		p := d.Annihilation
		rows := [][2]int{}
		for _, r := range p.InteractionTable {
			rows = append(rows, [2]int{r.Min, r.Max})
		}
		ok = p.Control.DC == 25 && p.Move.SaveDC == 19 && p.TouchDamage == ItemDice{Dice: 8, Sides: 10, Type: "force"} && p.Move.Damage == p.TouchDamage && p.TableDieSides == 100 && len(rows) == 3 && itemRollTable(rows, 100) && slices.Equal(p.SourcePages, []int{244, 245})
	case "staff-of-charming":
		p := d.Charming
		ok = p.MaxCharges == 10 && p.CastSpell.ChargeCost == 1 && p.ReflectEnchantment.ChargeCost == 1 && uniqueItemKeys(p.CastSpell.Spells) && len(p.AttunementByAnyOf) == 6 && uniqueItemKeys(p.AttunementByAnyOf) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 8, Type: "charges"} && p.Recharge.Bonus == 2
	case "staff-of-fire":
		p := d.Fire
		ok = p.MaxCharges == 10 && p.Resistance == "fire" && slices.Equal(p.Spells, []ItemStaffOfFireSpell{{Spell: "burning-hands", ChargeCost: 1}, {Spell: "fireball", ChargeCost: 3}, {Spell: "wall-of-fire", ChargeCost: 4}}) && p.Recharge.Dice == ItemDice{Dice: 1, Sides: 6, Type: "charges"} && p.Recharge.Bonus == 4
	case "staff-of-frost":
		p := d.Frost
		ok = p.MaxCharges == 10 && p.Resistance == "cold" && slices.Equal(p.Spells, []ItemStaffOfFrostSpell{{Spell: "cone-of-cold", ChargeCost: 5}, {Spell: "fog-cloud", ChargeCost: 1}, {Spell: "ice-storm", ChargeCost: 4}, {Spell: "wall-of-ice", ChargeCost: 4}}) && p.Recharge.Bonus == 4
	case "staff-of-healing":
		p := d.Healing
		ok = p.MaxCharges == 10 && p.CureWounds.ChargesPerSpellLevel == 1 && p.CureWounds.MaximumSpellLevel == 4 && slices.Equal(p.Spells, []ItemStaffOfHealingSpell{{Spell: "lesser-restoration", ChargeCost: 2}, {Spell: "mass-cure-wounds", ChargeCost: 5}}) && len(p.AttunementByAnyOf) == 3
	case "staff-of-power":
		p := d.Power
		r := p.RetributiveStrike
		costs := []int{}
		names := []string{}
		for _, s := range p.Spells {
			costs, names = append(costs, s.ChargeCost), append(names, s.Spell)
		}
		ok = p.MaxCharges == 20 && p.AttackBonus == 2 && p.HoldingBonus == 2 && slices.Equal(costs, []int{5, 5, 6, 5, 2, 5, 1, 1, 5}) && uniqueItemKeys(names) && r.SaveDC == 17 && r.YourDamagePerCharge == 16 && r.OtherDamagePerCharge == 4 && p.Recharge.Dice == ItemDice{Dice: 2, Sides: 8, Type: "charges"} && slices.Equal(p.SourcePages, []int{245, 246})
	case "staff-of-striking":
		p := d.Striking
		ok = p.MaxCharges == 10 && p.AttackBonus == 3 && p.DamageBonus == 3 && p.MaximumChargesPerHit == 3 && p.ExtraDamagePerCharge == ItemDice{Dice: 1, Sides: 6, Type: "force"}
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
