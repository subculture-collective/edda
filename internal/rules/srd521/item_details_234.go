package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_234.json
var items234JSON []byte

// ItemDetails234 holds exactly one typed payload per variant of catalogue batch
// #234 (magic items 10, SRD 5.2.1 pp. 228-229, with the Manual of Golems table
// on p. 230). Payloads are source facts only.
type ItemDetails234 struct {
	IronBands    *ItemIronBands        `json:"iron_bands,omitempty"`
	IronFlask    *ItemIronFlask        `json:"iron_flask,omitempty"`
	Javelin      *ItemJavelinLightning `json:"javelin_of_lightning,omitempty"`
	Lantern      *ItemLanternRevealing `json:"lantern_of_revealing,omitempty"`
	LuckBlade    *ItemLuckBlade        `json:"luck_blade,omitempty"`
	Disruption   *ItemMaceDisruption   `json:"mace_of_disruption,omitempty"`
	Smiting      *ItemMaceSmiting      `json:"mace_of_smiting,omitempty"`
	Terror       *ItemMaceTerror       `json:"mace_of_terror,omitempty"`
	Mantle       *ItemMantleResistance `json:"mantle_of_spell_resistance,omitempty"`
	BodilyHealth *ItemManualStudy      `json:"manual_of_bodily_health,omitempty"`
	GainfulExerc *ItemManualStudy      `json:"manual_of_gainful_exercise,omitempty"`
	ManualGolems *ItemManualOfGolems   `json:"manual_of_golems,omitempty"`
}

// ItemIronBands retains typed SRD source facts without implementing them.
type ItemIronBands struct {
	Appearance        string    `json:"appearance"`
	DiameterInches    int       `json:"diameter_inches"`
	WeightPounds      int       `json:"weight_pounds"`
	Action            string    `json:"action"`
	ThrowRangeFeet    int       `json:"throw_range_feet"`
	TargetMaxSize     string    `json:"target_max_size"`
	TargetVisibility  string    `json:"target_visibility"`
	OpensInto         string    `json:"opens_into"`
	AttackKind        string    `json:"attack_kind"`
	AttackAbility     string    `json:"attack_ability"`
	AttackAddsProf    bool      `json:"attack_adds_proficiency_bonus"`
	HitCondition      string    `json:"hit_condition"`
	ReleaseAction     string    `json:"release_action"`
	ContractsOn       []string  `json:"contracts_on"`
	BreakAction       string    `json:"break_action"`
	BreakCheck        ItemCheck `json:"break_check"`
	BreakBy           string    `json:"break_by"`
	BreakSuccess      []string  `json:"break_success"`
	BreakFailure      string    `json:"break_failure"`
	BreakLockoutHours int       `json:"break_failure_lockout_hours"`
	Reuse             string    `json:"reuse"`
}

// ItemIronFlask retains typed SRD source facts without implementing them.
type ItemIronFlask struct {
	Condition          string   `json:"condition"`
	Action             string   `json:"action"`
	RangeFeet          int      `json:"range_feet"`
	TargetVisibility   string   `json:"target_visibility"`
	RequiresEmpty      bool     `json:"requires_empty_flask"`
	TargetNativeTo     string   `json:"target_native_to"`
	SaveAbility        string   `json:"save_ability"`
	SaveDC             int      `json:"save_dc"`
	SaveFailure        string   `json:"save_failure"`
	RepeatTrapAdv      string   `json:"previously_trapped_save"`
	TrappedUntil       string   `json:"trapped_until"`
	Capacity           int      `json:"capacity_creatures"`
	TrappedEffects     []string `json:"trapped_effects"`
	ReleaseAction      string   `json:"release_action"`
	ReleaseBy          string   `json:"release_by"`
	ObeysHours         int      `json:"obeys_for_hours"`
	UnderstandsAnyLang bool     `json:"understands_commands_in_any_language"`
	SelfDefenseOn      []string `json:"self_defense_when"`
	SelfDefense        string   `json:"self_defense_behavior"`
	AfterDuration      string   `json:"after_duration"`
	IdentifyReveals    string   `json:"identify_reveals"`
	TypeDeterminedBy   string   `json:"type_determined_by"`
	NewlyFound         string   `json:"newly_found"`
}

// ItemJavelinLightning retains typed SRD source facts without implementing them.
type ItemJavelinLightning struct {
	SwapTrigger string                   `json:"damage_swap_trigger"`
	SwapFrom    string                   `json:"damage_swap_from"`
	SwapTo      string                   `json:"damage_swap_to"`
	Bolt        ItemJavelinLightningBolt `json:"lightning_bolt"`
}

// ItemJavelinLightningBolt is the Lightning Bolt property.
type ItemJavelinLightningBolt struct {
	Trigger     string   `json:"trigger"`
	RangeFeet   int      `json:"range_feet"`
	ForgoAttack bool     `json:"forgo_attack_roll"`
	Shape       string   `json:"shape"`
	WidthFeet   int      `json:"width_feet"`
	LineBetween string   `json:"line_between"`
	Affects     []string `json:"affects"`
	SaveAbility string   `json:"save_ability"`
	SaveDC      int      `json:"save_dc"`
	Damage      ItemDice `json:"damage"`
	OnSuccess   string   `json:"on_success"`
	Returns     string   `json:"returns"`
	Reuse       string   `json:"reuse"`
}

// ItemLanternRevealing retains typed SRD source facts without implementing them.
type ItemLanternRevealing struct {
	Condition              string   `json:"condition"`
	BurnHours              int      `json:"burn_hours"`
	OilPints               int      `json:"oil_pints"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	Reveals                []string `json:"reveals"`
	RevealsWhen            string   `json:"reveals_when"`
	HoodAction             string   `json:"hood_action"`
	HoodedLight            string   `json:"hooded_light"`
	HoodedRadiusFeet       int      `json:"hooded_radius_feet"`
}

// ItemLuckBlade retains typed SRD source facts without implementing them.
type ItemLuckBlade struct {
	AttackBonus   int          `json:"attack_bonus"`
	DamageBonus   int          `json:"damage_bonus"`
	SaveBonus     int          `json:"save_bonus"`
	SaveBonusWhen string       `json:"save_bonus_when"`
	Luck          ItemLuckProp `json:"luck"`
	Wish          ItemWishProp `json:"wish"`
}

// ItemLuckProp is the Luck property.
type ItemLuckProp struct {
	Condition string `json:"condition"`
	Action    string `json:"action"`
	Effect    string `json:"effect"`
	Forbidden string `json:"forbidden_while"`
	Reroll    string `json:"reroll_use"`
	Reuse     string `json:"reuse"`
}

// ItemWishProp is the Wish property. The charges start at 1d3.
type ItemWishProp struct {
	InitialCharges ItemDice `json:"initial_charges"`
	Condition      string   `json:"condition"`
	ChargeCost     int      `json:"charge_cost"`
	Spell          string   `json:"spell"`
	Reuse          string   `json:"reuse"`
	LostWhen       string   `json:"property_lost_when"`
}

// ItemMaceDisruption retains typed SRD source facts without implementing them.
type ItemMaceDisruption struct {
	ExtraDamage            ItemDice `json:"extra_damage"`
	ExtraDamageTargets     []string `json:"extra_damage_targets"`
	DestroyAtHitPoints     int      `json:"destroy_at_hit_points_or_fewer"`
	SaveAbility            string   `json:"save_ability"`
	SaveDC                 int      `json:"save_dc"`
	SaveFailure            string   `json:"save_failure"`
	SaveSuccess            string   `json:"save_success_condition"`
	SuccessEnds            string   `json:"success_condition_ends"`
	LightCondition         string   `json:"light_condition"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
}

// ItemMaceSmiting retains typed SRD source facts without implementing them.
type ItemMaceSmiting struct {
	AttackBonus         int    `json:"attack_bonus"`
	DamageBonus         int    `json:"damage_bonus"`
	ConstructBonus      int    `json:"construct_bonus"`
	CriticalTrigger     string `json:"critical_trigger"`
	CriticalDamage      int    `json:"critical_extra_damage"`
	CriticalDamageType  string `json:"critical_damage_type"`
	ConstructCritDamage int    `json:"construct_critical_extra_damage"`
	DestroyAtHitPoints  int    `json:"construct_destroyed_at_hit_points_or_fewer"`
}

// ItemMaceTerror retains typed SRD source facts without implementing them.
type ItemMaceTerror struct {
	MaxCharges      int      `json:"max_charges"`
	Regain          ItemDice `json:"regain_expended_charges"`
	RegainWhen      string   `json:"regain_when"`
	Condition       string   `json:"condition"`
	Action          string   `json:"action"`
	ChargeCost      int      `json:"charge_cost"`
	RangeFeet       int      `json:"range_feet"`
	Targets         string   `json:"targets"`
	SaveAbility     string   `json:"save_ability"`
	SaveDC          int      `json:"save_dc"`
	FailCondition   string   `json:"fail_condition"`
	DurationSeconds int      `json:"duration_seconds"`
	Behavior        []string `json:"frightened_behavior"`
	NoMoveAction    string   `json:"no_movement_action"`
	RepeatSave      string   `json:"repeat_save"`
	SuccessEnds     string   `json:"success_ends"`
}

// ItemMantleResistance retains typed SRD source facts without implementing them.
type ItemMantleResistance struct {
	Condition string `json:"condition"`
	Roll      string `json:"roll"`
	Mode      string `json:"mode"`
	Against   string `json:"against"`
}

// ItemManualStudy is the shared rule of the two ability manuals: study for a
// time, raise one ability score, then lose the magic until it returns.
type ItemManualStudy struct {
	StudyHours        int    `json:"study_hours"`
	WithinDays        int    `json:"within_days"`
	Ability           string `json:"ability"`
	Increase          int    `json:"increase"`
	Maximum           int    `json:"maximum"`
	LosesMagic        bool   `json:"then_loses_magic"`
	RegainsAfterYears int    `json:"regains_magic_after_years"`
}

// ItemManualOfGolems retains typed SRD source facts without implementing them.
// The 1d20 table is printed at the top of p. 230.
type ItemManualOfGolems struct {
	SourcePages        []int          `json:"source_pages"`
	TypeChosenBy       string         `json:"type_chosen_by"`
	TableDieSides      int            `json:"table_die_sides"`
	Golems             []ItemGolemRow `json:"golems"`
	DecipherCaster     bool           `json:"decipher_requires_spellcaster"`
	DecipherSlots      int            `json:"decipher_spell_slots"`
	DecipherSlotLevel  int            `json:"decipher_slot_level"`
	UnableDamage       ItemDice       `json:"unable_user_damage"`
	Uninterrupted      bool           `json:"work_uninterrupted"`
	MaxRestHoursPerDay int            `json:"max_rest_hours_per_day"`
	SuppliesCost       string         `json:"supplies_cost"`
	OnCompletion       string         `json:"on_completion"`
	Animation          string         `json:"animation"`
	StatBlockIn        string         `json:"stat_block_in"`
	GolemControl       string         `json:"golem_control"`
}

// ItemGolemRow is one row of the Manual of Golems 1d20 table. Costs are the
// supplies the creator pays, not a price of the item.
type ItemGolemRow struct {
	Golem    string `json:"golem"`
	TableMin int    `json:"table_min"`
	TableMax int    `json:"table_max"`
	Days     int    `json:"time_days"`
	CostGP   int    `json:"cost_gp"`
}

// itemManifest234 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. No family has more than one printed variant: the
// four golems are rows of a table inside one payload.
var itemManifest234 = itemSourceManifest{Batch: "data-magic-items-10", Issue: 234, Families: []itemSourceManifestFamily{
	{Slug: "iron-bands", Name: "Iron Bands", Page: 228, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "throw-at-huge-or-smaller-within-60-feet once-until-dawn",
		Keys: "action.magic attack.ranged-dexterity-plus-proficiency condition.restrained action.bonus-release check.strength-athletics-dc-20 effect.destroyed-on-successful-break item.once-until-dawn", Payload: "iron_bands",
		Variants: []itemSourceManifestVariant{{"standard", "Iron Bands", "rare"}}},
	{Slug: "iron-flask", Name: "Iron Flask", Page: 228, Category: "wondrous-item", Attunement: "none", Equip: "hold", Use: "magic-action", UseRestrictions: "empty-flask-target-from-another-plane",
		Keys: "action.magic save.wisdom-dc-17 effect.trapped-in-flask effect.obeys-commands-1-hour item.holds-one-creature", Payload: "iron_flask",
		Variants: []itemSourceManifestVariant{{"standard", "Iron Flask", "legendary"}}},
	{Slug: "javelin-of-lightning", Name: "Javelin of Lightning", Page: 228, Category: "weapon", Options: "javelin", Attunement: "none", Equip: "wield", Use: "weapon-attack", UseRestrictions: "lightning-bolt-once-until-dawn",
		Keys: "damage.lightning-instead-of-piercing property.lightning-bolt area.line-5-foot save.dexterity-dc-13 damage.lightning-4d6 item.once-until-dawn", Payload: "javelin_of_lightning",
		Variants: []itemSourceManifestVariant{{"standard", "Javelin of Lightning", "uncommon"}}},
	{Slug: "lantern-of-revealing", Name: "Lantern of Revealing", Page: 228, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "lower-hood",
		Keys: "light.bright-30-dim-30 effect.reveals-invisible-in-bright-light item.burns-6-hours-per-pint action.utilize light.hooded-dim-5", Payload: "lantern_of_revealing",
		Variants: []itemSourceManifestVariant{{"standard", "Lantern of Revealing", "uncommon"}}},
	{Slug: "luck-blade", Name: "Luck Blade", Page: 229, Category: "weapon", Options: "glaive greatsword longsword rapier scimitar sickle shortsword", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "luck-once-until-dawn wish-once-until-dawn",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "none"},
		Keys:    "item.attunement bonus.attack-and-damage-1 bonus.saving-throw-1 property.luck action.none property.wish item.charges-1d3 item.once-until-dawn", Payload: "luck_blade",
		Variants: []itemSourceManifestVariant{{"standard", "Luck Blade", "legendary"}}},
	{Slug: "mace-of-disruption", Name: "Mace of Disruption", Page: 229, Category: "weapon", Options: "mace", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement damage.radiant-2d6 target.fiend-or-undead save.wisdom-dc-15 effect.destroyed-at-25-hit-points condition.frightened light.bright-20-dim-20", Payload: "mace_of_disruption",
		Variants: []itemSourceManifestVariant{{"standard", "Mace of Disruption", "rare"}}},
	{Slug: "mace-of-smiting", Name: "Mace of Smiting", Page: 229, Category: "weapon", Options: "mace", Attunement: "none", Equip: "wield", Use: "weapon-attack",
		Keys: "bonus.attack-and-damage-1 bonus.construct-3 damage.bludgeoning-7-on-natural-20 damage.bludgeoning-14-on-construct-natural-20 effect.construct-destroyed-at-25-hit-points", Payload: "mace_of_smiting",
		Variants: []itemSourceManifestVariant{{"standard", "Mace of Smiting", "rare"}}},
	{Slug: "mace-of-terror", Name: "Mace of Terror", Page: 229, Category: "weapon", Options: "mace", Attunement: "required", Equip: "wield", Use: "magic-action", UseRestrictions: "expend-1-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.charges action.magic area.range-30-foot save.wisdom-dc-15 condition.frightened duration.1-minute effect.frightened-must-flee", Payload: "mace_of_terror",
		Variants: []itemSourceManifestVariant{{"standard", "Mace of Terror", "rare"}}},
	{Slug: "mantle-of-spell-resistance", Name: "Mantle of Spell Resistance", Page: 229, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear save.advantage-against-spells", Payload: "mantle_of_spell_resistance",
		Variants: []itemSourceManifestVariant{{"standard", "Mantle of Spell Resistance", "rare"}}},
	{Slug: "manual-of-bodily-health", Name: "Manual of Bodily Health", Page: 229, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "study-48-hours-within-6-days",
		Keys: "item.study-48-hours-in-6-days ability.constitution-plus-2-maximum-30 item.loses-magic-regains-in-century", Payload: "manual_of_bodily_health",
		Variants: []itemSourceManifestVariant{{"standard", "Manual of Bodily Health", "very-rare"}}},
	{Slug: "manual-of-gainful-exercise", Name: "Manual of Gainful Exercise", Page: 229, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "study-48-hours-within-6-days",
		Keys: "item.study-48-hours-in-6-days ability.strength-plus-2-maximum-30 item.loses-magic-regains-in-century", Payload: "manual_of_gainful_exercise",
		Variants: []itemSourceManifestVariant{{"standard", "Manual of Gainful Exercise", "very-rare"}}},
	{Slug: "manual-of-golems", Name: "Manual of Golems", Page: 229, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "spellcaster-with-two-level-5-slots",
		Keys: "item.golem-type-by-gm-or-1d20-table creation.time-and-cost-by-golem requirement.two-level-5-spell-slots damage.psychic-6d6-unable-reader item.consumed-in-eldritch-flames creature.golem-obeys-creator", Payload: "manual_of_golems",
		Variants: []itemSourceManifestVariant{{"standard", "Manual of Golems", "very-rare"}}},
}}

// LoadItems234 returns a fresh validated copy of batch #234. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems234() (ItemSourceBatch[ItemDetails234], error) { return DecodeItems234(items234JSON) }

// DecodeItems234 validates supplied JSON as batch #234 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems234(data []byte) (ItemSourceBatch[ItemDetails234], error) {
	return decodeItemSourceBatch(data, itemManifest234, validateItemDetails234)
}

// validateItemDetails234 checks the internal consistency of each payload and a
// few source literals that the family's mechanic keys also name.
func validateItemDetails234(family, variant string, d *ItemDetails234) error {
	positive := func(values ...int) bool {
		for _, v := range values {
			if v <= 0 {
				return false
			}
		}
		return true
	}
	dice := func(x ItemDice) bool { return x.Dice > 0 && x.Sides > 0 }
	ok := true
	switch family {
	case "iron-bands":
		b := d.IronBands
		ok = positive(b.DiameterInches, b.WeightPounds, b.ThrowRangeFeet, b.BreakCheck.DC, b.BreakLockoutHours) && uniqueItemKeys(b.ContractsOn) && uniqueItemKeys(b.BreakSuccess) &&
			b.BreakCheck.DC == 20 // p. 228; named by the mechanic keys
	case "iron-flask":
		f := d.IronFlask
		ok = positive(f.RangeFeet, f.SaveDC, f.Capacity, f.ObeysHours) && uniqueItemKeys(f.TrappedEffects) && uniqueItemKeys(f.SelfDefenseOn) && f.SaveDC == 17
	case "javelin-of-lightning":
		j := d.Javelin
		ok = positive(j.Bolt.RangeFeet, j.Bolt.WidthFeet, j.Bolt.SaveDC) && dice(j.Bolt.Damage) && uniqueItemKeys(j.Bolt.Affects) && j.Bolt.SaveDC == 13
	case "lantern-of-revealing":
		l := d.Lantern
		ok = positive(l.BurnHours, l.OilPints, l.BrightLightFeet, l.DimLightAdditionalFeet, l.HoodedRadiusFeet) && uniqueItemKeys(l.Reveals)
	case "luck-blade":
		b := d.LuckBlade
		ok = positive(b.AttackBonus, b.DamageBonus, b.SaveBonus, b.Wish.ChargeCost) && dice(b.Wish.InitialCharges)
	case "mace-of-disruption":
		m := d.Disruption
		ok = dice(m.ExtraDamage) && positive(m.DestroyAtHitPoints, m.SaveDC, m.BrightLightFeet, m.DimLightAdditionalFeet) && uniqueItemKeys(m.ExtraDamageTargets) && m.SaveDC == 15
	case "mace-of-smiting":
		m := d.Smiting
		ok = positive(m.AttackBonus, m.DamageBonus, m.CriticalDamage, m.DestroyAtHitPoints) && m.ConstructBonus > m.AttackBonus && m.ConstructCritDamage > m.CriticalDamage
	case "mace-of-terror":
		m := d.Terror
		ok = positive(m.MaxCharges, m.ChargeCost, m.RangeFeet, m.SaveDC, m.DurationSeconds) && dice(m.Regain) && uniqueItemKeys(m.Behavior) && m.MaxCharges == 3 && m.SaveDC == 15
	case "mantle-of-spell-resistance":
		ok = d.Mantle.Mode == "advantage"
	case "manual-of-bodily-health", "manual-of-gainful-exercise":
		m := d.BodilyHealth
		want := "constitution"
		if family == "manual-of-gainful-exercise" {
			m, want = d.GainfulExerc, "strength"
		}
		ok = m.Ability == want && m.StudyHours == 48 && m.WithinDays == 6 && m.Increase == 2 && m.Maximum == 30 && m.RegainsAfterYears == 100
	case "manual-of-golems":
		g := d.ManualGolems
		rows := [][2]int{}
		names := []string{}
		for _, r := range g.Golems {
			rows = append(rows, [2]int{r.TableMin, r.TableMax})
			names = append(names, r.Golem)
			ok = ok && positive(r.Days, r.CostGP)
		}
		ok = ok && slices.Equal(g.SourcePages, []int{229, 230}) && itemRollTable(rows, g.TableDieSides) && uniqueItemKeys(names) && len(names) == 4 &&
			positive(g.DecipherSlots, g.DecipherSlotLevel, g.MaxRestHoursPerDay) && dice(g.UnableDamage)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
