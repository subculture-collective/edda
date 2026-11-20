package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_235.json
var items235JSON []byte

// ItemDetails235 holds exactly one typed payload per variant of catalogue batch
// #235 (magic items 11, SRD 5.2.1 pp. 230-233). Payloads are source facts only.
type ItemDetails235 struct {
	Manual       *ItemManualOfQuickness    `json:"manual_of_quickness_of_action,omitempty"`
	Pigments     *ItemMarvelousPigments    `json:"marvelous_pigments,omitempty"`
	Medallion    *ItemMedallionOfThoughts  `json:"medallion_of_thoughts,omitempty"`
	Mirror       *ItemMirrorOfLifeTrapping `json:"mirror_of_life_trapping,omitempty"`
	MithralArmor *ItemMithralArmor         `json:"mithral_armor,omitempty"`
	Deck         *ItemMysteriousDeck       `json:"mysterious_deck,omitempty"`
	Adaptation   *ItemNecklaceAdaptation   `json:"necklace_of_adaptation,omitempty"`
	Fireballs    *ItemNecklaceFireballs    `json:"necklace_of_fireballs,omitempty"`
	PrayerBeads  *ItemNecklacePrayerBeads  `json:"necklace_of_prayer_beads,omitempty"`
	NineLives    *ItemNineLivesStealer     `json:"nine_lives_stealer,omitempty"`
	Oathbow      *ItemOathbow              `json:"oathbow,omitempty"`
	Etherealness *ItemOilOfEtherealness    `json:"oil_of_etherealness,omitempty"`
}

// ItemRoll235 is a printed die expression such as 1d4 + 2. It is source data
// only and is never passed to a dice resolver.
type ItemRoll235 struct {
	Dice  int `json:"dice"`
	Sides int `json:"sides"`
	Bonus int `json:"bonus"`
}

// ItemFact235 is one numeric source fact: a key, a non-negative amount and the
// unit it is printed in. The unit "minus" marks an amount printed with a minus
// sign (the Euryale card's -2 penalty) because payload numbers are never negative.
// A unit such as "d4" marks a die count: amount 1 and unit "d4" is 1d4.
type ItemFact235 struct {
	Key    string `json:"key"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// ItemManualOfQuickness retains typed SRD source facts without implementing them.
type ItemManualOfQuickness struct {
	StudyHours        int    `json:"study_hours"`
	StudyWithinDays   int    `json:"study_within_days"`
	Ability           string `json:"ability"`
	Increase          int    `json:"increase"`
	MaxScore          int    `json:"max_score"`
	LosesMagic        bool   `json:"loses_magic_after_use"`
	RegainsAfterYears int    `json:"regains_magic_after_years"`
}

// ItemMarvelousPigments retains typed SRD source facts without implementing them.
type ItemMarvelousPigments struct {
	Container           string      `json:"container"`
	Pots                ItemRoll235 `json:"pots"`
	BrushWeightPounds   int         `json:"brush_and_box_weight_pounds"`
	PotsExpended        int         `json:"pots_expended_per_painting"`
	CubeFeet            int         `json:"cube_feet"`
	PaintingMinutes     int         `json:"painting_minutes"`
	RequiresConcentrate bool        `json:"requires_concentration"`
	MustRemainInCube    bool        `json:"must_remain_in_cube"`
	FailureTriggers     []string    `json:"failure_triggers"`
	FailureResult       string      `json:"failure_result"`
	CompletionResult    string      `json:"completion_result"`
	ObjectMaxValueGP    int         `json:"object_max_value_gp"`
	TotalMaxValueGP     int         `json:"total_max_value_per_pot_gp"`
	OverValueResult     string      `json:"over_value_result"`
	EnergyResult        string      `json:"energy_result"`
}

// ItemMedallionOfThoughts retains typed SRD source facts without implementing them.
type ItemMedallionOfThoughts struct {
	Condition      string      `json:"condition"`
	MaxCharges     int         `json:"max_charges"`
	ChargesPerCast int         `json:"charges_per_cast"`
	Spell          string      `json:"spell"`
	SaveDC         int         `json:"save_dc"`
	Recharge       ItemRoll235 `json:"recharge_amount"`
	RechargeTime   string      `json:"recharge_time"`
}

// ItemMirrorOfLifeTrapping retains typed SRD source facts without implementing them.
type ItemMirrorOfLifeTrapping struct {
	HeightFeet            int      `json:"height_feet"`
	WidthFeet             int      `json:"width_feet"`
	WeightPounds          int      `json:"weight_pounds"`
	AC                    int      `json:"ac"`
	HP                    int      `json:"hp"`
	Immunities            []string `json:"immunities"`
	Vulnerabilities       []string `json:"vulnerabilities"`
	DestroyedWhen         string   `json:"destroyed_when"`
	ActivateAction        string   `json:"activate_action"`
	ActivateRequires      []string `json:"activate_requires"`
	ActivateRangeFeet     int      `json:"activate_range_feet"`
	DeactivateAction      string   `json:"deactivate_action"`
	ReflectionRangeFeet   int      `json:"reflection_range_feet"`
	SaveAbility           string   `json:"save_ability"`
	SaveDC                int      `json:"save_dc"`
	SaveFailure           string   `json:"save_failure"`
	SaveAdvantageWhen     string   `json:"save_advantage_when"`
	SaveAutoSucceed       []string `json:"save_auto_succeed"`
	Cells                 int      `json:"cells"`
	CellDescription       string   `json:"cell_description"`
	CellVisibilityFeet    int      `json:"cell_visibility_feet"`
	TrappedAge            bool     `json:"trapped_creatures_age"`
	TrappedNeed           bool     `json:"trapped_creatures_need_food_drink_or_sleep"`
	EscapeMethod          string   `json:"escape_method"`
	OtherwiseConfined     string   `json:"otherwise_confined"`
	FullCellsResult       string   `json:"full_cells_result"`
	FreedAppears          string   `json:"freed_creature_appears"`
	ShatterResult         string   `json:"shatter_result"`
	NameAction            string   `json:"name_action"`
	NameRangeFeet         int      `json:"name_range_feet"`
	NameResult            string   `json:"name_result"`
	FreeAction            string   `json:"free_action"`
	FreeCommandOrdinal    int      `json:"free_command_word_ordinal"`
	FreeResult            string   `json:"free_result"`
	ExtradimensionalItems []string `json:"extradimensional_items_destroyed"`
	GateTo                string   `json:"gate_to"`
	GateRadiusFeet        int      `json:"gate_radius_feet"`
	GateAffected          string   `json:"gate_affected"`
	GateAfter             string   `json:"gate_after"`
	GateOneWay            bool     `json:"gate_one_way"`
	GateCanReopen         bool     `json:"gate_can_reopen"`
}

// ItemMithralArmor retains typed SRD source facts without implementing them.
type ItemMithralArmor struct {
	Wearable            string `json:"wearable"`
	StealthDisadvantage string `json:"stealth_disadvantage"`
	StrengthRequirement string `json:"strength_requirement"`
}

// ItemDeckRow is one card row of a Mysterious Deck table column. A printed 00
// is stored as 100.
type ItemDeckRow struct {
	Card string `json:"card"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

// ItemDeckCard is the effect of one card as keys and numeric facts.
type ItemDeckCard struct {
	Card    string        `json:"card"`
	Effects []string      `json:"effects"`
	Facts   []ItemFact235 `json:"facts" srd:"may-be-empty"`
}

// ItemAvatarAbility is one ability line of the Avatar of Death stat block.
type ItemAvatarAbility struct {
	Ability string `json:"ability"`
	Score   int    `json:"score"`
	Mod     int    `json:"mod"`
	Save    int    `json:"save"`
}

// ItemAvatarDamage is a printed damage line: its average, dice, bonus and type.
type ItemAvatarDamage struct {
	Average int    `json:"average"`
	Dice    int    `json:"dice"`
	Sides   int    `json:"sides"`
	Bonus   int    `json:"bonus"`
	Type    string `json:"type"`
}

// ItemAvatarOfDeath is the stat block printed on p. 232 for the Skull card.
// "Half the HP maximum of its summoner" and a Proficiency Bonus equal to the
// summoner's are rules, recorded as keys.
type ItemAvatarOfDeath struct {
	Size                string              `json:"size"`
	CreatureType        string              `json:"creature_type"`
	Alignment           string              `json:"alignment"`
	AC                  int                 `json:"ac"`
	InitiativeModifier  int                 `json:"initiative_modifier"`
	InitiativeScore     int                 `json:"initiative_score"`
	HitPoints           string              `json:"hit_points"`
	SpeedFeet           int                 `json:"speed_feet"`
	FlySpeedFeet        int                 `json:"fly_speed_feet"`
	Hover               bool                `json:"hover"`
	Abilities           []ItemAvatarAbility `json:"abilities"`
	DamageImmunities    []string            `json:"damage_immunities"`
	ConditionImmunities []string            `json:"condition_immunities"`
	TruesightFeet       int                 `json:"truesight_feet"`
	PassivePerception   int                 `json:"passive_perception"`
	Languages           string              `json:"languages"`
	ChallengeRating     string              `json:"challenge_rating"`
	XP                  int                 `json:"xp"`
	ProficiencyBonus    string              `json:"proficiency_bonus"`
	IncorporealMovement string              `json:"incorporeal_movement"`
	IncorporealDamage   ItemAvatarDamage    `json:"incorporeal_damage_ending_turn_inside_object"`
	MultiattackCount    string              `json:"multiattack_count"`
	MultiattackAttack   string              `json:"multiattack_attack"`
	ScytheAttack        string              `json:"reaping_scythe_attack"`
	ScytheAutomaticHit  bool                `json:"reaping_scythe_automatic_hit"`
	ScytheReachFeet     int                 `json:"reaping_scythe_reach_feet"`
	ScytheDamage        ItemAvatarDamage    `json:"reaping_scythe_damage"`
	ScytheExtraDamage   ItemAvatarDamage    `json:"reaping_scythe_extra_damage"`
}

// ItemMysteriousDeck retains typed SRD source facts without implementing them.
// The text starts on p. 231 and continues on p. 232. The two table columns are
// two lists of rows, each covering 1..100 exactly; a card absent from a column
// (printed as a dash) has no row in it.
type ItemMysteriousDeck struct {
	SourcePages         []int             `json:"source_pages"`
	Materials           []string          `json:"card_materials"`
	DeckSizes           []int             `json:"deck_sizes"`
	ThirteenCardPercent int               `json:"thirteen_card_percent"`
	TableDieSides       int               `json:"table_die_sides"`
	Deck13              []ItemDeckRow     `json:"deck_13_cards"`
	Deck22              []ItemDeckRow     `json:"deck_22_cards"`
	DeclareDraws        bool              `json:"declare_draw_count_before_drawing"`
	ExcessDrawsNoEffect bool              `json:"excess_draws_have_no_effect"`
	MaxHoursBetween     int               `json:"max_hours_between_draws"`
	UndrawnResult       string            `json:"undrawn_cards_result"`
	DrawnCardDisappears bool              `json:"drawn_card_disappears"`
	DoesNotReappear     []string          `json:"does_not_reappear"`
	RerollWhenGone      []string          `json:"reroll_when_gone"`
	Cards               []ItemDeckCard    `json:"cards"`
	Avatar              ItemAvatarOfDeath `json:"avatar_of_death"`
}

// ItemNecklaceAdaptation retains typed SRD source facts without implementing them.
type ItemNecklaceAdaptation struct {
	Condition        string   `json:"condition"`
	BreatheIn        string   `json:"breathe_normally_in"`
	AdvantageSaves   []string `json:"advantage_on_saves_to"`
	AdvantageAgainst string   `json:"advantage_against_condition"`
}

// ItemNecklaceFireballs retains typed SRD source facts without implementing them.
type ItemNecklaceFireballs struct {
	Beads              ItemRoll235 `json:"beads"`
	DetachAction       string      `json:"detach_action"`
	ThrowRangeFeet     int         `json:"throw_range_feet"`
	Spell              string      `json:"spell"`
	SpellLevel         int         `json:"spell_level"`
	SaveDC             int         `json:"save_dc"`
	MultipleBeads      string      `json:"multiple_beads"`
	ExtraDamagePerBead ItemRoll235 `json:"extra_damage_per_bead_after_first"`
	MaximumDamage      ItemRoll235 `json:"maximum_damage"`
}

// ItemPrayerBead is one row of the Necklace of Prayer Beads 1d20 table. A zero
// spell level means the source gives none.
type ItemPrayerBead struct {
	Min        int    `json:"min"`
	Max        int    `json:"max"`
	Bead       string `json:"bead"`
	Spell      string `json:"spell"`
	SpellLevel int    `json:"spell_level"`
}

// ItemNecklacePrayerBeads retains typed SRD source facts without implementing them.
type ItemNecklacePrayerBeads struct {
	AttunementBy      []string         `json:"attunement_by_any_of"`
	MagicBeads        ItemRoll235      `json:"magic_beads"`
	MagicMaterials    []string         `json:"magic_bead_materials"`
	NonmagicalBeads   []string         `json:"nonmagical_bead_examples"`
	RemovedLosesMagic bool             `json:"removed_bead_loses_magic"`
	BeadTypes         int              `json:"bead_types"`
	TypeSelection     []string         `json:"type_selection"`
	AllowsRepeats     bool             `json:"necklace_may_have_repeated_types"`
	RequiresWearing   bool             `json:"requires_wearing"`
	CastAction        string           `json:"cast_action"`
	SaveDC            string           `json:"save_dc"`
	UsesPerBead       int              `json:"uses_per_bead"`
	BeadRecharge      string           `json:"bead_recharge"`
	TableDieSides     int              `json:"table_die_sides"`
	Beads             []ItemPrayerBead `json:"beads_table"`
}

// ItemNineLivesStealer retains typed SRD source facts without implementing them.
// The charge count is rolled when the weapon is found, so the variant's charge
// maximum is the highest result of the printed dice.
type ItemNineLivesStealer struct {
	AttackBonus    int         `json:"attack_bonus"`
	DamageBonus    int         `json:"damage_bonus"`
	Property       string      `json:"property"`
	Charges        ItemRoll235 `json:"charges"`
	TargetHPBelow  int         `json:"target_hit_points_below"`
	TriggerRoll    int         `json:"trigger_d20_roll"`
	SaveAbility    string      `json:"save_ability"`
	SaveDC         int         `json:"save_dc"`
	SaveFailure    string      `json:"save_failure"`
	AutoSucceed    []string    `json:"save_auto_succeed"`
	ChargeLostWhen string      `json:"charge_lost_when"`
	ChargesLost    int         `json:"charges_lost"`
	WhenNoCharges  string      `json:"when_no_charges"`
}

// ItemOathbow retains typed SRD source facts without implementing them.
type ItemOathbow struct {
	NockedEffect       string   `json:"nocked_effect"`
	CommandTrigger     string   `json:"command_trigger"`
	CommandMethods     []string `json:"command_methods"`
	SwornEnemy         string   `json:"sworn_enemy"`
	SwornEnemyEnds     []string `json:"sworn_enemy_ends"`
	SwornEnemyDays     int      `json:"sworn_enemy_days"`
	MaxSwornEnemies    int      `json:"max_sworn_enemies"`
	NewEnemyWhen       string   `json:"new_enemy_when"`
	AttackAdvantage    bool     `json:"ranged_attack_advantage_against_sworn_enemy"`
	CoverIgnored       []string `json:"cover_ignored"`
	NoLongRangePenalty bool     `json:"no_long_range_disadvantage"`
	ExtraDamage        ItemDice `json:"extra_damage"`
	ExtraDamageWhen    string   `json:"extra_damage_when"`
	OtherWeapons       string   `json:"other_weapons"`
	OtherWeaponsWhile  string   `json:"other_weapons_while"`
}

// ItemOilOfEtherealness retains typed SRD source facts without implementing them.
type ItemOilOfEtherealness struct {
	Covers                 string `json:"covers"`
	ExtraVialsPerSizeAbove int    `json:"extra_vials_per_size_above_medium"`
	ApplyMinutes           int    `json:"apply_minutes"`
	Spell                  string `json:"spell"`
	DurationHours          int    `json:"duration_hours"`
	Appearance             string `json:"appearance"`
}

// itemManifest235 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Each family has the one source-listed form
// "standard"; table rows and card effects are data inside a payload.
var itemManifest235 = itemSourceManifest{Batch: "data-magic-items-11", Issue: 235, Families: []itemSourceManifestFamily{
	{Slug: "manual-of-quickness-of-action", Name: "Manual of Quickness of Action", Page: 230, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "study", UseRestrictions: "48-hours-within-6-days",
		Keys: "study.48-hours-within-6-days ability.dexterity-increase-2 ability.dexterity-maximum-30 item.loses-magic-after-use item.regains-magic-after-a-century", Payload: "manual_of_quickness_of_action",
		Variants: []itemSourceManifestVariant{{"standard", "Manual of Quickness of Action", "very-rare"}}},
	{Slug: "marvelous-pigments", Name: "Marvelous Pigments", Page: 230, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "paint", UseRestrictions: "expend-one-pot concentration remain-in-cube",
		Keys: "item.pots-1d4 action.expend-one-pot area.cube-20-foot time.10-minutes concentration.required effect.paintings-become-real object.value-limit-25-gp object.total-value-limit-500-gp effect.energy-dissipates", Payload: "marvelous_pigments",
		Variants: []itemSourceManifestVariant{{"standard", "Marvelous Pigments", "very-rare"}}},
	{Slug: "medallion-of-thoughts", Name: "Medallion of Thoughts", Page: 230, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-at-dawn"},
		Keys:    "item.attunement item.wear item.charges spell.cast-detect-thoughts spell.detect-thoughts-save-dc-13", Payload: "medallion_of_thoughts",
		Variants: []itemSourceManifestVariant{{"standard", "Medallion of Thoughts", "uncommon"}}},
	{Slug: "mirror-of-life-trapping", Name: "Mirror of Life Trapping", Page: 230, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "command-word hanging-on-vertical-surface within-5-feet",
		Keys: "action.magic item.command-word save.charisma-dc-15 effect.trapped-in-extradimensional-cell area.cells-12 effect.free-one-trapped-creature effect.destroys-bag-of-holding-and-portable-hole effect.gate-to-astral-plane", Payload: "mirror_of_life_trapping",
		Variants: []itemSourceManifestVariant{{"standard", "Mirror of Life Trapping", "very-rare"}}},
	{Slug: "mithral-armor", Name: "Mithral Armor", Page: 231, Category: "armor", Options: "any-medium-armor any-heavy-armor", Excluded: "hide-armor", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear item.worn-under-clothes armor.no-stealth-disadvantage armor.no-strength-requirement", Payload: "mithral_armor",
		Variants: []itemSourceManifestVariant{{"standard", "Mithral Armor", "uncommon"}}},
	{Slug: "mysterious-deck", Name: "Mysterious Deck", Page: 231, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "draw-card", UseRestrictions: "declare-draw-count draw-within-1-hour-of-previous",
		Keys: "item.deck-13-or-22-cards table.mysterious-deck-1d100 action.declare-draws card.effects-22 card.fool-and-jester-leave-deck creature.avatar-of-death", Payload: "mysterious_deck",
		Variants: []itemSourceManifestVariant{{"standard", "Mysterious Deck", "legendary"}}},
	{Slug: "necklace-of-adaptation", Name: "Necklace of Adaptation", Page: 232, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear breath.any-environment save.advantage-poisoned", Payload: "necklace_of_adaptation",
		Variants: []itemSourceManifestVariant{{"standard", "Necklace of Adaptation", "uncommon"}}},
	{Slug: "necklace-of-fireballs", Name: "Necklace of Fireballs", Page: 233, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "detach-bead throw-within-60-feet",
		Keys: "item.beads-1d6-plus-3 action.magic spell.fireball-level-3 save.dexterity-dc-15 damage.extra-1d6-per-extra-bead damage.maximum-12d6", Payload: "necklace_of_fireballs",
		Variants: []itemSourceManifestVariant{{"standard", "Necklace of Fireballs", "rare"}}},
	{Slug: "necklace-of-prayer-beads", Name: "Necklace of Prayer Beads", Page: 233, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "bonus-action once-per-bead-until-next-dawn",
		Keys: "item.attunement-by-cleric-druid-or-paladin item.wear item.beads-1d4-plus-2 table.bead-type-1d20 spell.cast-from-bead item.bead-loses-magic-when-removed", Payload: "necklace_of_prayer_beads",
		Variants: []itemSourceManifestVariant{{"standard", "Necklace of Prayer Beads", "rare"}}},
	{Slug: "nine-lives-stealer", Name: "Nine Lives Stealer", Page: 233, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "natural-20-against-creature-under-100-hit-points",
		Charges: &ItemCharges{Mode: "limited", Maximum: 9, PerUse: 1, Recharge: "none"},
		Keys:    "item.attunement item.charges-1d8-plus-1 bonus.attack-and-damage-2 save.constitution-dc-15 effect.slain-instantly effect.constructs-and-undead-succeed effect.loses-property-without-charges", Payload: "nine_lives_stealer",
		Variants: []itemSourceManifestVariant{{"standard", "Nine Lives Stealer", "very-rare"}}},
	{Slug: "oathbow", Name: "Oathbow", Page: 233, Category: "weapon", Options: "longbow shortbow", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "command-words-on-ranged-attack",
		Keys: "item.attunement effect.sworn-enemy attack.advantage-against-sworn-enemy cover.ignored-against-sworn-enemy damage.extra-3d6 attack.disadvantage-with-other-weapons", Payload: "oathbow",
		Variants: []itemSourceManifestVariant{{"standard", "Oathbow", "very-rare"}}},
	{Slug: "oil-of-etherealness", Name: "Oil of Etherealness", Page: 233, Category: "potion", Attunement: "none", Equip: "none", Use: "apply", UseRestrictions: "10-minutes one-vial-covers-medium-or-smaller",
		Keys: "item.vials-per-size-above-medium time.apply-10-minutes spell.etherealness duration.1-hour", Payload: "oil_of_etherealness",
		Variants: []itemSourceManifestVariant{{"standard", "Oil of Etherealness", "rare"}}},
}}

// LoadItems235 returns a fresh validated copy of batch #235. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems235() (ItemSourceBatch[ItemDetails235], error) { return DecodeItems235(items235JSON) }

// DecodeItems235 validates supplied JSON as batch #235 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems235(data []byte) (ItemSourceBatch[ItemDetails235], error) {
	return decodeItemSourceBatch(data, itemManifest235, validateItemDetails235)
}

// itemDeckCards235 lists the 22 cards of the Mysterious Deck in printed order.
var itemDeckCards235 = []string{"balance", "comet", "donjon", "euryale", "fates", "flames", "fool", "gem", "jester", "key", "knight", "moon", "puzzle", "rogue", "ruin", "sage", "skull", "star", "sun", "talons", "throne", "void"}

// validateItemDeckRows235 requires rows to name distinct cards of the deck, in
// printed order, and to cover the 1d100 exactly once.
func validateItemDeckRows235(rows []ItemDeckRow, size int) bool {
	rolls := [][2]int{}
	cards := []string{}
	for _, r := range rows {
		rolls = append(rolls, [2]int{r.Min, r.Max})
		cards = append(cards, r.Card)
	}
	if len(rows) != size || !itemRollTable(rolls, 100) || !uniqueItemKeys(cards) {
		return false
	}
	last := -1
	for _, c := range cards {
		i := slices.Index(itemDeckCards235, c)
		if i <= last {
			return false
		}
		last = i
	}
	return true
}

// validateItemDetails235 checks the internal consistency of each payload and a
// few source literals that the family's mechanic keys also name.
func validateItemDetails235(family, variant string, d *ItemDetails235) error {
	positive := func(values ...int) bool {
		for _, v := range values {
			if v <= 0 {
				return false
			}
		}
		return true
	}
	roll := func(r ItemRoll235) bool { return positive(r.Dice, r.Sides) }
	ok := true
	switch family {
	case "manual-of-quickness-of-action":
		m := d.Manual
		ok = positive(m.StudyHours, m.StudyWithinDays, m.Increase, m.MaxScore, m.RegainsAfterYears) && m.StudyHours == 48 && m.Increase == 2 && m.MaxScore == 30 && m.RegainsAfterYears == 100
	case "marvelous-pigments":
		p := d.Pigments
		ok = roll(p.Pots) && positive(p.BrushWeightPounds, p.PotsExpended, p.CubeFeet, p.PaintingMinutes, p.ObjectMaxValueGP, p.TotalMaxValueGP) && uniqueItemKeys(p.FailureTriggers) &&
			p.CubeFeet == 20 && p.ObjectMaxValueGP == 25 && p.TotalMaxValueGP == 500 && p.TotalMaxValueGP > p.ObjectMaxValueGP
	case "medallion-of-thoughts":
		m := d.Medallion
		ok = roll(m.Recharge) && positive(m.MaxCharges, m.ChargesPerCast, m.SaveDC) && m.MaxCharges == 5 && m.SaveDC == 13
	case "mirror-of-life-trapping":
		m := d.Mirror
		ok = positive(m.HeightFeet, m.WidthFeet, m.WeightPounds, m.AC, m.HP, m.ActivateRangeFeet, m.ReflectionRangeFeet, m.SaveDC, m.Cells, m.CellVisibilityFeet, m.NameRangeFeet, m.FreeCommandOrdinal, m.GateRadiusFeet) &&
			uniqueItemKeys(m.Immunities) && uniqueItemKeys(m.Vulnerabilities) && uniqueItemKeys(m.ActivateRequires) && uniqueItemKeys(m.SaveAutoSucceed) && uniqueItemKeys(m.ExtradimensionalItems) &&
			m.SaveDC == 15 && m.Cells == 12 && m.FreeCommandOrdinal == 2
	case "mithral-armor":
		ok = d.MithralArmor.Wearable != ""
	case "mysterious-deck":
		k := d.Deck
		ok = slices.Equal(k.SourcePages, []int{231, 232}) && uniqueItemKeys(k.Materials) && slices.Equal(k.DeckSizes, []int{13, 22}) && k.ThirteenCardPercent == 75 && k.TableDieSides == 100 &&
			positive(k.MaxHoursBetween) && uniqueItemKeys(k.DoesNotReappear) && uniqueItemKeys(k.RerollWhenGone) &&
			validateItemDeckRows235(k.Deck13, 13) && validateItemDeckRows235(k.Deck22, 22) && len(k.Cards) == len(itemDeckCards235)
		for i, row := range k.Deck22 {
			ok = ok && row.Card == itemDeckCards235[i]
		}
		for i, c := range k.Cards {
			facts := []string{}
			for _, f := range c.Facts {
				facts = append(facts, f.Key)
			}
			ok = ok && c.Card == itemDeckCards235[i] && uniqueItemKeys(c.Effects) && (len(facts) == 0 || uniqueItemKeys(facts))
		}
		a := k.Avatar
		abilities := []string{}
		for _, line := range a.Abilities {
			abilities = append(abilities, line.Ability)
		}
		ok = ok && slices.Equal(abilities, []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"}) && uniqueItemKeys(a.DamageImmunities) && uniqueItemKeys(a.ConditionImmunities) &&
			positive(a.AC, a.InitiativeScore, a.SpeedFeet, a.FlySpeedFeet, a.TruesightFeet, a.PassivePerception, a.ScytheReachFeet, a.ScytheDamage.Dice, a.ScytheDamage.Sides, a.ScytheExtraDamage.Dice, a.ScytheExtraDamage.Sides, a.IncorporealDamage.Dice, a.IncorporealDamage.Sides) &&
			a.AC == 20 && a.InitiativeScore == 13
	case "necklace-of-adaptation":
		n := d.Adaptation
		ok = uniqueItemKeys(n.AdvantageSaves)
	case "necklace-of-fireballs":
		n := d.Fireballs
		ok = roll(n.Beads) && roll(n.ExtraDamagePerBead) && roll(n.MaximumDamage) && positive(n.ThrowRangeFeet, n.SpellLevel, n.SaveDC) && n.SpellLevel == 3 && n.SaveDC == 15 && n.ThrowRangeFeet == 60
	case "necklace-of-prayer-beads":
		n := d.PrayerBeads
		rows := [][2]int{}
		beads := []string{}
		for _, b := range n.Beads {
			rows = append(rows, [2]int{b.Min, b.Max})
			beads = append(beads, b.Bead)
		}
		ok = roll(n.MagicBeads) && positive(n.BeadTypes, n.UsesPerBead) && uniqueItemKeys(n.AttunementBy) && uniqueItemKeys(n.MagicMaterials) && uniqueItemKeys(n.NonmagicalBeads) && uniqueItemKeys(n.TypeSelection) &&
			n.BeadTypes == 6 && len(n.Beads) == n.BeadTypes && uniqueItemKeys(beads) && n.TableDieSides == 20 && itemRollTable(rows, n.TableDieSides)
	case "nine-lives-stealer":
		n := d.NineLives
		ok = roll(n.Charges) && positive(n.AttackBonus, n.DamageBonus, n.TargetHPBelow, n.TriggerRoll, n.SaveDC, n.ChargesLost) && uniqueItemKeys(n.AutoSucceed) &&
			n.AttackBonus == 2 && n.DamageBonus == 2 && n.SaveDC == 15 && n.TargetHPBelow == 100 && n.TriggerRoll == 20
	case "oathbow":
		o := d.Oathbow
		ok = positive(o.SwornEnemyDays, o.MaxSwornEnemies, o.ExtraDamage.Dice, o.ExtraDamage.Sides) && uniqueItemKeys(o.CommandMethods) && uniqueItemKeys(o.SwornEnemyEnds) && uniqueItemKeys(o.CoverIgnored) && o.SwornEnemyDays == 7
	case "oil-of-etherealness":
		o := d.Etherealness
		ok = positive(o.ExtraVialsPerSizeAbove, o.ApplyMinutes, o.DurationHours) && o.ApplyMinutes == 10 && o.DurationHours == 1
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
