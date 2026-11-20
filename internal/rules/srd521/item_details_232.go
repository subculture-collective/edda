package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_232.json
var items232JSON []byte

// ItemDetails232 holds exactly one typed payload per variant of catalogue batch
// #232 (magic items 08, SRD 5.2.1 pp. 223-225). Payloads are source facts only.
type ItemDetails232 struct {
	GemSeeing      *ItemGemOfSeeing              `json:"gem_of_seeing,omitempty"`
	GiantSlayer    *ItemGiantSlayer              `json:"giant_slayer,omitempty"`
	Glamoured      *ItemGlamouredStuddedLeather  `json:"glamoured_studded_leather,omitempty"`
	MissileSnaring *ItemGlovesOfMissileSnaring   `json:"gloves_of_missile_snaring,omitempty"`
	SwimClimb      *ItemGlovesOfSwimmingClimbing `json:"gloves_of_swimming_and_climbing,omitempty"`
	Thievery       *ItemGlovesOfThievery         `json:"gloves_of_thievery,omitempty"`
	Goggles        *ItemGogglesOfNight           `json:"goggles_of_night,omitempty"`
	Hammer         *ItemHammerOfThunderbolts     `json:"hammer_of_thunderbolts,omitempty"`
	Haversack      *ItemHandyHaversack           `json:"handy_haversack,omitempty"`
	HatDisguise    *ItemHatOfDisguise            `json:"hat_of_disguise,omitempty"`
	HatManySpells  *ItemHatOfManySpells          `json:"hat_of_many_spells,omitempty"`
	Headband       *ItemHeadbandOfIntellect      `json:"headband_of_intellect,omitempty"`
}

// ItemFact232 is one numeric source fact: a key, a non-negative amount and the
// unit it is printed in.
type ItemFact232 struct {
	Key    string `json:"key"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// ItemGemOfSeeing retains typed SRD source facts without implementing them.
type ItemGemOfSeeing struct {
	MaxCharges      int      `json:"max_charges"`
	Action          string   `json:"action"`
	ChargeCost      int      `json:"charge_cost"`
	Sense           string   `json:"sense"`
	SenseRangeFeet  int      `json:"sense_range_feet"`
	Condition       string   `json:"condition"`
	DurationSeconds int      `json:"duration_seconds"`
	Recharge        ItemDice `json:"recharge"`
	RechargeWhen    string   `json:"recharge_when"`
}

// ItemGiantSlayer retains typed SRD source facts without implementing them.
type ItemGiantSlayer struct {
	AttackBonus  int      `json:"attack_bonus"`
	DamageBonus  int      `json:"damage_bonus"`
	Trigger      string   `json:"trigger"`
	CreatureType string   `json:"target_creature_type"`
	ExtraDamage  ItemDice `json:"extra_damage"`
	SaveAbility  string   `json:"save_ability"`
	SaveDC       int      `json:"save_dc"`
	Condition    string   `json:"condition"`
}

// ItemGlamouredStuddedLeather retains typed SRD source facts without implementing them.
type ItemGlamouredStuddedLeather struct {
	ArmorClassBonus int      `json:"armor_class_bonus"`
	Action          string   `json:"action"`
	Appearance      string   `json:"appearance"`
	ChosenDetails   []string `json:"chosen_details"`
	KeepsBulk       bool     `json:"retains_bulk_and_weight"`
	EndsOn          []string `json:"ends_on"`
}

// ItemGlovesOfMissileSnaring retains typed SRD source facts without implementing them.
type ItemGlovesOfMissileSnaring struct {
	Trigger        string   `json:"trigger"`
	Action         string   `json:"action"`
	Reduction      ItemDice `json:"reduction"`
	BonusAbility   string   `json:"reduction_bonus_ability"`
	NeedsFreeHand  bool     `json:"requires_free_hand"`
	CatchCondition string   `json:"catch_condition"`
	Catchable      string   `json:"catchable"`
}

// ItemGlovesOfSwimmingClimbing retains typed SRD source facts without implementing them.
type ItemGlovesOfSwimmingClimbing struct {
	Condition     string   `json:"condition"`
	ClimbSpeed    string   `json:"climb_speed"`
	SwimSpeed     string   `json:"swim_speed"`
	CheckBonus    int      `json:"check_bonus"`
	CheckAbility  string   `json:"check_ability"`
	CheckSkill    string   `json:"check_skill"`
	CheckSubjects []string `json:"check_subjects"`
}

// ItemGlovesOfThievery retains typed SRD source facts without implementing them.
type ItemGlovesOfThievery struct {
	Appearance   string `json:"appearance"`
	Condition    string `json:"condition"`
	CheckBonus   int    `json:"check_bonus"`
	CheckAbility string `json:"check_ability"`
	CheckSkill   string `json:"check_skill"`
}

// ItemGogglesOfNight retains typed SRD source facts without implementing them.
type ItemGogglesOfNight struct {
	Appearance           string `json:"appearance"`
	Condition            string `json:"condition"`
	DarkvisionFeet       int    `json:"darkvision_feet"`
	ExistingIncreaseFeet int    `json:"existing_darkvision_increase_feet"`
}

// ItemHammerOfThunderbolts retains typed SRD source facts without implementing them.
type ItemHammerOfThunderbolts struct {
	AttackBonus   int                    `json:"attack_bonus"`
	DamageBonus   int                    `json:"damage_bonus"`
	MaxCharges    int                    `json:"max_charges"`
	Hurl          ItemHammerHurl         `json:"hurl"`
	Recharge      ItemDice               `json:"recharge"`
	RechargeBonus int                    `json:"recharge_bonus"`
	RechargeWhen  string                 `json:"recharge_when"`
	Bane          ItemHammerGiantsBane   `json:"giants_bane"`
	Might         ItemHammerMightOfGiant `json:"might_of_giants"`
}

// ItemHammerHurl is the charge-powered thrown attack and its thunderclap.
type ItemHammerHurl struct {
	ChargeCost        int    `json:"charge_cost"`
	ThrownNormalFeet  int    `json:"thrown_normal_range_feet"`
	ThrownLongFeet    int    `json:"thrown_long_range_feet"`
	ThunderclapFeet   int    `json:"thunderclap_audible_feet"`
	BurstRadiusFeet   int    `json:"burst_radius_feet"`
	BurstExcludes     string `json:"burst_excludes"`
	SaveAbility       string `json:"save_ability"`
	SaveDC            int    `json:"save_dc"`
	Condition         string `json:"condition"`
	ConditionEnds     string `json:"condition_ends"`
	ReturnsToHand     string `json:"returns_to_hand"`
	RequiresHitToClap bool   `json:"thunderclap_requires_hit"`
}

// ItemHammerGiantsBane is the Giants' Bane benefit of the Giant's Bane property.
type ItemHammerGiantsBane struct {
	AttunedToWeapon bool     `json:"requires_attunement_to_weapon"`
	WornWithEither  []string `json:"requires_wearing_either"`
	AttunedToWorn   bool     `json:"requires_attunement_to_worn_item"`
	NaturalRoll     int      `json:"natural_roll"`
	CreatureType    string   `json:"target_creature_type"`
	SaveAbility     string   `json:"save_ability"`
	SaveDC          int      `json:"save_dc"`
	Failure         string   `json:"failure"`
}

// ItemHammerMightOfGiant is the Might of Giants benefit.
type ItemHammerMightOfGiant struct {
	Affects          []string `json:"affects"`
	StrengthIncrease int      `json:"strength_increase"`
	StrengthMaximum  int      `json:"strength_maximum"`
}

// ItemHandyHaversack retains typed SRD source facts without implementing them.
type ItemHandyHaversack struct {
	WeightPounds     int                 `json:"weight_pounds"`
	SidePouches      ItemHaversackPouch  `json:"side_pouches"`
	CentralPouch     ItemHaversackPouch  `json:"central_pouch"`
	RetrievalActions []string            `json:"retrieval_actions"`
	RetrievedOnTop   bool                `json:"retrieved_item_on_top"`
	RuptureTriggers  []string            `json:"rupture_triggers"`
	ContentsLost     bool                `json:"contents_lost_when_destroyed"`
	ArtifactReturns  bool                `json:"artifact_turns_up_again"`
	InsideOut        string              `json:"inside_out_result"`
	AirMinutes       int                 `json:"air_minutes_per_pouch"`
	AirDividedBy     string              `json:"air_divided_by"`
	Nested           ItemHaversackNested `json:"nested_extradimensional_space"`
}

// ItemHaversackPouch is one kind of pouch: how many and what each holds.
type ItemHaversackPouch struct {
	Count         int `json:"count"`
	CapacityLbs   int `json:"capacity_pounds"`
	CapacityCubic int `json:"capacity_cubic_feet"`
}

// ItemHaversackNested is the result of placing the haversack inside another
// extradimensional space.
type ItemHaversackNested struct {
	Items         []string `json:"items"`
	Result        string   `json:"result"`
	GateTo        string   `json:"gate_to"`
	GateRadiusFt  int      `json:"gate_radius_feet"`
	SparedByCover string   `json:"spared_by"`
	Destination   string   `json:"destination"`
	GateCloses    bool     `json:"gate_closes_after"`
	OneWay        bool     `json:"gate_one_way"`
	Reopens       bool     `json:"gate_reopens"`
}

// ItemHatOfDisguise retains typed SRD source facts without implementing them.
type ItemHatOfDisguise struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
	EndsWhen  string `json:"ends_when"`
}

// ItemHatOfManySpells retains typed SRD source facts without implementing them.
// The text starts on p. 224 and its table continues on p. 225. The attunement
// class is a payload fact because the shared attunement field has no class slot.
type ItemHatOfManySpells struct {
	SourcePages      []int               `json:"source_pages"`
	AttunementClass  string              `json:"attunement_class"`
	Condition        string              `json:"condition"`
	FocusSpellList   string              `json:"focus_spell_list"`
	SomaticComponent string              `json:"added_somatic_component"`
	Unknown          ItemHatUnknownSpell `json:"unknown_spell"`
	SpellStats       string              `json:"spell_stats"`
	FailureTable     []ItemHatRow        `json:"failure_table"`
}

// ItemHatUnknownSpell is the Unknown Spell property.
type ItemHatUnknownSpell struct {
	MinLevel            int    `json:"min_spell_level"`
	SpellList           string `json:"spell_list"`
	LevelLimit          string `json:"level_limit"`
	MaterialCostLimitGP int    `json:"material_cost_limit_gp"`
	SlotExpended        string `json:"slot_expended"`
	CheckAbility        string `json:"check_ability"`
	CheckSkill          string `json:"check_skill"`
	CheckDCBase         int    `json:"check_dc_base"`
	CheckDCAdds         string `json:"check_dc_adds"`
	SuccessCasting      string `json:"success_casting_time"`
	SuccessLockout      string `json:"success_lockout_until"`
	FailureResult       string `json:"failure_result"`
	TableDieSides       int    `json:"table_die_sides"`
}

// ItemHatRow is one row of the failure table. A printed 00 is stored as 100.
// SubDieSides is 0 when the row has no further roll; otherwise Outcomes covers
// that die exactly once.
type ItemHatRow struct {
	Min         int              `json:"min"`
	Max         int              `json:"max"`
	Effect      string           `json:"effect"`
	SubDieSides int              `json:"sub_die_sides"`
	Outcomes    []ItemHatOutcome `json:"outcomes" srd:"may-be-empty"`
	Facts       []ItemFact232    `json:"facts" srd:"may-be-empty"`
	Traits      []string         `json:"traits" srd:"may-be-empty"`
}

// ItemHatOutcome is one result of a row's further roll.
type ItemHatOutcome struct {
	Min    int    `json:"min"`
	Max    int    `json:"max"`
	Result string `json:"result"`
}

// ItemHeadbandOfIntellect retains typed SRD source facts without implementing them.
type ItemHeadbandOfIntellect struct {
	Condition         string `json:"condition"`
	Ability           string `json:"ability"`
	Score             int    `json:"score"`
	NoEffectAtOrAbove int    `json:"no_effect_at_or_above_score"`
}

// itemManifest232 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. No family lists several kinds, so each has one
// "standard" variant; the Hat of Many Spells' tables are payload data.
var itemManifest232 = itemSourceManifest{Batch: "data-magic-items-08", Issue: 232, Families: []itemSourceManifestFamily{
	{Slug: "gem-of-seeing", Name: "Gem of Seeing", Page: 223, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.charges action.magic senses.truesight-120-feet duration.ten-minutes item.peer-through-gem", Payload: "gem_of_seeing",
		Variants: []itemSourceManifestVariant{{"standard", "Gem of Seeing", "rare"}}},
	{Slug: "giant-slayer", Name: "Giant Slayer", Page: 223, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "none", Equip: "wield", Use: "weapon-attack",
		Keys: "item.magic-weapon bonus.attack-and-damage-1 damage.extra-2d6-against-giants save.strength-dc-15 condition.prone", Payload: "giant_slayer",
		Variants: []itemSourceManifestVariant{{"standard", "Giant Slayer", "rare"}}},
	{Slug: "glamoured-studded-leather", Name: "Glamoured Studded Leather", Page: 223, Category: "armor", Options: "studded-leather-armor", Attunement: "none", Equip: "wear", Use: "bonus-action", UseRestrictions: "assume-appearance-of-clothing-or-armor",
		Keys: "item.wear defense.armor-class-bonus-1 action.bonus-action effect.illusory-appearance item.retains-bulk-and-weight effect.appearance-lasts-until-used-again-or-doffed", Payload: "glamoured_studded_leather",
		Variants: []itemSourceManifestVariant{{"standard", "Glamoured Studded Leather", "rare"}}},
	{Slug: "gloves-of-missile-snaring", Name: "Gloves of Missile Snaring", Page: 224, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "reaction", UseRestrictions: "hit-by-ranged-or-thrown-attack free-hand",
		Keys: "item.attunement item.wear action.reaction damage.reduce-1d10-plus-dexterity-modifier item.free-hand-required item.catch-ammunition-or-weapon", Payload: "gloves_of_missile_snaring",
		Variants: []itemSourceManifestVariant{{"standard", "Gloves of Missile Snaring", "uncommon"}}},
	{Slug: "gloves-of-swimming-and-climbing", Name: "Gloves of Swimming and Climbing", Page: 224, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear movement.climb-speed-equals-speed movement.swim-speed-equals-speed check.athletics-bonus-5-climb-or-swim", Payload: "gloves_of_swimming_and_climbing",
		Variants: []itemSourceManifestVariant{{"standard", "Gloves of Swimming and Climbing", "uncommon"}}},
	{Slug: "gloves-of-thievery", Name: "Gloves of Thievery", Page: 224, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear item.imperceptible-while-worn check.sleight-of-hand-bonus-5", Payload: "gloves_of_thievery",
		Variants: []itemSourceManifestVariant{{"standard", "Gloves of Thievery", "uncommon"}}},
	{Slug: "goggles-of-night", Name: "Goggles of Night", Page: 224, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear vision.darkvision-60-feet vision.darkvision-increase-60-feet", Payload: "goggles_of_night",
		Variants: []itemSourceManifestVariant{{"standard", "Goggles of Night", "uncommon"}}},
	{Slug: "hammer-of-thunderbolts", Name: "Hammer of Thunderbolts", Page: 224, Category: "weapon", Options: "maul warhammer", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "expend-one-charge thrown-ranged-attack",
		Charges: &ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-plus-1-at-dawn"},
		Keys:    "item.attunement item.magic-weapon bonus.attack-and-damage-1 item.charges property.thrown-20-60 effect.thunderclap-audible-300-feet save.constitution-dc-17 condition.stunned item.returns-to-hand property.giants-bane effect.die-on-failed-save-against-giant ability.strength-score-bonus-4-max-30", Payload: "hammer_of_thunderbolts",
		Variants: []itemSourceManifestVariant{{"standard", "Hammer of Thunderbolts", "legendary"}}},
	{Slug: "handy-haversack", Name: "Handy Haversack", Page: 224, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "container", UseRestrictions: "retrieve-with-utilize-or-bonus-action",
		Keys: "item.extradimensional-storage item.side-pouches-200-pounds-25-cubic-feet item.central-pouch-500-pounds-64-cubic-feet item.weight-5-pounds item.retrieved-item-on-top item.destroyed-if-overloaded-pierced-or-torn item.contents-lost-if-destroyed item.inside-out-spills-contents item.air-10-minutes-per-pouch item.nested-extradimensional-space-opens-astral-gate", Payload: "handy_haversack",
		Variants: []itemSourceManifestVariant{{"standard", "Handy Haversack", "rare"}}},
	{Slug: "hat-of-disguise", Name: "Hat of Disguise", Page: 224, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "disguise-self",
		Keys: "item.attunement item.wear spell.cast-disguise-self item.spell-ends-when-removed", Payload: "hat_of_disguise",
		Variants: []itemSourceManifestVariant{{"standard", "Hat of Disguise", "uncommon"}}},
	{Slug: "hat-of-many-spells", Name: "Hat of Many Spells", Page: 224, Category: "wondrous-item", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "unknown-wizard-spell until-short-or-long-rest-after-success",
		Keys: "item.attunement attunement.wizard item.hold spell.focus-wizard-spells spell.added-somatic-component spell.cast-unknown-wizard-spell check.intelligence-arcana-dc-10-plus-spell-level item.spell-slot-expended table.failure-1d100 spell.uses-your-save-dc-and-attack-bonus", Payload: "hat_of_many_spells",
		Variants: []itemSourceManifestVariant{{"standard", "Hat of Many Spells", "very-rare"}}},
	{Slug: "headband-of-intellect", Name: "Headband of Intellect", Page: 225, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear ability.intelligence-score-19", Payload: "headband_of_intellect",
		Variants: []itemSourceManifestVariant{{"standard", "Headband of Intellect", "uncommon"}}},
}}

// LoadItems232 returns a fresh validated copy of batch #232. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems232() (ItemSourceBatch[ItemDetails232], error) { return DecodeItems232(items232JSON) }

// DecodeItems232 validates supplied JSON as batch #232 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems232(data []byte) (ItemSourceBatch[ItemDetails232], error) {
	return decodeItemSourceBatch(data, itemManifest232, validateItemDetails232)
}

// validateItemDetails232 checks the internal consistency of each payload and a
// few source literals that the family's mechanic keys also name.
func validateItemDetails232(family, variant string, d *ItemDetails232) error {
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
	case "gem-of-seeing":
		g := d.GemSeeing
		ok = positive(g.MaxCharges, g.ChargeCost, g.SenseRangeFeet, g.DurationSeconds, g.Recharge.Dice, g.Recharge.Sides) &&
			g.MaxCharges == 3 && g.SenseRangeFeet == 120 && g.DurationSeconds == 600 && g.Recharge.Dice == 1 && g.Recharge.Sides == 3 // p. 223; named by the mechanic keys
	case "giant-slayer":
		g := d.GiantSlayer
		ok = positive(g.AttackBonus, g.DamageBonus, g.ExtraDamage.Dice, g.ExtraDamage.Sides, g.SaveDC) && g.SaveDC == 15
	case "glamoured-studded-leather":
		g := d.Glamoured
		ok = positive(g.ArmorClassBonus) && uniqueItemKeys(g.ChosenDetails) && uniqueItemKeys(g.EndsOn)
	case "gloves-of-missile-snaring":
		g := d.MissileSnaring
		ok = positive(g.Reduction.Dice, g.Reduction.Sides)
	case "gloves-of-swimming-and-climbing":
		g := d.SwimClimb
		ok = positive(g.CheckBonus) && uniqueItemKeys(g.CheckSubjects) && g.CheckBonus == 5
	case "gloves-of-thievery":
		ok = positive(d.Thievery.CheckBonus) && d.Thievery.CheckBonus == 5
	case "goggles-of-night":
		ok = positive(d.Goggles.DarkvisionFeet, d.Goggles.ExistingIncreaseFeet)
	case "hammer-of-thunderbolts":
		h := d.Hammer
		ok = positive(h.AttackBonus, h.DamageBonus, h.MaxCharges, h.Hurl.ChargeCost, h.Hurl.ThrownNormalFeet, h.Hurl.ThrownLongFeet, h.Hurl.ThunderclapFeet, h.Hurl.BurstRadiusFeet,
			h.Hurl.SaveDC, h.Recharge.Dice, h.Recharge.Sides, h.RechargeBonus, h.Bane.NaturalRoll, h.Bane.SaveDC, h.Might.StrengthIncrease, h.Might.StrengthMaximum) &&
			h.Hurl.ThrownLongFeet > h.Hurl.ThrownNormalFeet && h.Hurl.SaveDC == 17 && h.Bane.SaveDC == 17 && h.MaxCharges == 5 &&
			uniqueItemKeys(h.Bane.WornWithEither) && uniqueItemKeys(h.Might.Affects)
	case "handy-haversack":
		h := d.Haversack
		ok = positive(h.WeightPounds, h.SidePouches.Count, h.SidePouches.CapacityLbs, h.SidePouches.CapacityCubic, h.CentralPouch.Count, h.CentralPouch.CapacityLbs, h.CentralPouch.CapacityCubic,
			h.AirMinutes, h.Nested.GateRadiusFt) && uniqueItemKeys(h.RetrievalActions) && uniqueItemKeys(h.RuptureTriggers) && uniqueItemKeys(h.Nested.Items) &&
			h.SidePouches.Count == 2 && h.CentralPouch.Count == 1
	case "hat-of-many-spells":
		h := d.HatManySpells
		u := h.Unknown
		ok = slices.Equal(h.SourcePages, []int{224, 225}) && positive(u.MinLevel, u.MaterialCostLimitGP, u.CheckDCBase, u.TableDieSides) && u.TableDieSides == 100 && u.CheckDCBase == 10 && len(h.FailureTable) == 11
		rows := [][2]int{}
		for _, r := range h.FailureTable {
			rows = append(rows, [2]int{r.Min, r.Max})
			sub := [][2]int{}
			for _, o := range r.Outcomes {
				sub = append(sub, [2]int{o.Min, o.Max})
			}
			// A row either has no further roll or a table covering its die exactly.
			ok = ok && (r.SubDieSides == 0 && len(r.Outcomes) == 0 || itemRollTable(sub, r.SubDieSides)) && validateItemFacts232(r.Facts) && (len(r.Traits) == 0 || uniqueItemKeys(r.Traits))
		}
		ok = ok && itemRollTable(rows, 100)
	case "headband-of-intellect":
		ok = positive(d.Headband.Score, d.Headband.NoEffectAtOrAbove) && d.Headband.Score == 19
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}

// validateItemFacts232 requires distinct, non-empty fact keys.
func validateItemFacts232(facts []ItemFact232) bool {
	keys := []string{}
	for _, f := range facts {
		keys = append(keys, f.Key)
	}
	return len(keys) == 0 || uniqueItemKeys(keys)
}
