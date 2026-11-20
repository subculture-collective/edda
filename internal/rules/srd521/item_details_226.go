package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_226.json
var items226JSON []byte

// ItemDetails226 holds exactly one typed payload per variant of catalogue batch
// #226 (magic items 02, SRD 5.2.1 pp. 211–213). Payloads are source facts only.
type ItemDetails226 struct {
	ArrowCatching *ItemArrowCatchingShield `json:"arrow_catching,omitempty"`
	Beans         *ItemBagOfBeans          `json:"beans,omitempty"`
	Devouring     *ItemBagOfDevouring      `json:"devouring,omitempty"`
	Holding       *ItemBagOfHolding        `json:"holding,omitempty"`
	Tricks        *ItemBagOfTricks         `json:"tricks,omitempty"`
	Force         *ItemBeadOfForce         `json:"force,omitempty"`
	Nourishment   *ItemBeadOfNourishment   `json:"nourishment,omitempty"`
	Dwarvenkind   *ItemBeltOfDwarvenkind   `json:"dwarvenkind,omitempty"`
	GiantStrength *ItemBeltOfGiantStrength `json:"giant_strength,omitempty"`
	Berserker     *ItemBerserkerAxe        `json:"berserker,omitempty"`
	Elvenkind     *ItemBootsOfElvenkind    `json:"elvenkind,omitempty"`
	Levitation    *ItemBootsOfLevitation   `json:"levitation,omitempty"`
}

// ItemArrowCatchingShield retains typed SRD source facts without implementing them.
type ItemArrowCatchingShield struct {
	Condition         string `json:"condition"`
	ACBonus           int    `json:"ac_bonus"`
	ACBonusAgainst    string `json:"ac_bonus_against"`
	AddsToShieldBonus bool   `json:"adds_to_shield_bonus"`
	RedirectAction    string `json:"redirect_action"`
	RedirectTrigger   string `json:"redirect_trigger"`
	RedirectRangeFeet int    `json:"redirect_range_feet"`
	RedirectEffect    string `json:"redirect_effect"`
}

// ItemBeanQuantity is one source number of a bean effect: dice plus a flat bonus, in a unit.
type ItemBeanQuantity struct {
	Label string `json:"label"`
	Dice  int    `json:"dice"`
	Sides int    `json:"sides"`
	Bonus int    `json:"bonus"`
	Unit  string `json:"unit"`
}

// ItemBeanEffect is one row of the Bag of Beans 1d100 table. Row 96–00 is recorded as 96–100.
type ItemBeanEffect struct {
	Min        int                `json:"min"`
	Max        int                `json:"max"`
	Effect     string             `json:"effect"`
	Creatures  []string           `json:"creatures" srd:"may-be-empty"`
	Quantities []ItemBeanQuantity `json:"quantities"`
	Facts      []string           `json:"facts" srd:"may-be-empty"`
}

// ItemBagOfBeans retains typed SRD source facts without implementing them.
type ItemBagOfBeans struct {
	Material                string           `json:"material"`
	BeansDice               int              `json:"beans_dice"`
	BeansDieSides           int              `json:"beans_die_sides"`
	WeightPoundsNumerator   int              `json:"weight_pounds_numerator"`
	WeightPoundsDenominator int              `json:"weight_pounds_denominator"`
	WeightIgnoresContents   bool             `json:"weight_ignores_contents"`
	WhenEmpty               string           `json:"when_empty"`
	DumpTrigger             string           `json:"dump_trigger"`
	AreaShape               string           `json:"area_shape"`
	AreaRadiusFeet          int              `json:"area_radius_feet"`
	AreaCenter              string           `json:"area_center"`
	DumpedBeansDestroyed    bool             `json:"dumped_beans_destroyed"`
	AffectsUser             bool             `json:"affects_user"`
	Save                    string           `json:"save"`
	SaveDC                  int              `json:"save_dc"`
	Damage                  ItemDice         `json:"damage"`
	SaveSuccess             string           `json:"save_success"`
	PlantSteps              []string         `json:"plant_steps"`
	PlantDelaySeconds       int              `json:"plant_delay_seconds"`
	BeanDisappears          bool             `json:"bean_disappears"`
	EffectOrigin            string           `json:"effect_origin"`
	EffectSelection         []string         `json:"effect_selection"`
	TableDieSides           int              `json:"table_die_sides"`
	Effects                 []ItemBeanEffect `json:"effects"`
}

// ItemBagOfDevouring retains typed SRD source facts without implementing them.
type ItemBagOfDevouring struct {
	Resembles           string    `json:"resembles"`
	Nature              string    `json:"nature"`
	ClosesOrifice       string    `json:"closes_orifice"`
	CreatureSenses      string    `json:"creature_senses"`
	Devours             []string  `json:"devours"`
	DevouredOutcome     string    `json:"devoured_outcome"`
	PullTrigger         string    `json:"pull_trigger"`
	PullChancePercent   int       `json:"pull_chance_percent"`
	EscapeAction        string    `json:"escape_action"`
	EscapeCheck         ItemCheck `json:"escape_check"`
	RescueAction        string    `json:"rescue_action"`
	RescueCheck         ItemCheck `json:"rescue_check"`
	RescueCondition     string    `json:"rescue_condition"`
	TurnStartInside     string    `json:"turn_start_inside"`
	Stores              string    `json:"stores"`
	CapacityCubicFeet   int       `json:"capacity_cubic_feet"`
	SwallowFrequency    string    `json:"swallow_frequency"`
	SwallowEffect       string    `json:"swallow_effect"`
	SwallowDeterminedBy string    `json:"swallow_determined_by"`
	DestroyedBy         []string  `json:"destroyed_by"`
	DestroyedEffect     string    `json:"destroyed_effect"`
}

// ItemBagOfHolding retains typed SRD source facts without implementing them.
type ItemBagOfHolding struct {
	InteriorSquareSideFeet int      `json:"interior_square_side_feet"`
	InteriorDepthFeet      int      `json:"interior_depth_feet"`
	InteriorApproximate    bool     `json:"interior_approximate"`
	CapacityPounds         int      `json:"capacity_pounds"`
	CapacityCubicFeet      int      `json:"capacity_cubic_feet"`
	WeightPounds           int      `json:"weight_pounds"`
	WeightIgnoresContents  bool     `json:"weight_ignores_contents"`
	RetrieveAction         string   `json:"retrieve_action"`
	DestroyedBy            []string `json:"destroyed_by"`
	DestroyedEffect        string   `json:"destroyed_effect"`
	InsideOutEffect        string   `json:"inside_out_effect"`
	InsideOutReuse         string   `json:"inside_out_reuse"`
	AirMinutes             int      `json:"air_minutes"`
	AirDividedBy           string   `json:"air_divided_by"`
	NestingSpaces          []string `json:"nesting_spaces"`
	NestingEffects         []string `json:"nesting_effects"`
	GateOrigin             string   `json:"gate_origin"`
	GateShape              string   `json:"gate_shape"`
	GateRadiusFeet         int      `json:"gate_radius_feet"`
	GateEffect             string   `json:"gate_effect"`
	GateCloses             bool     `json:"gate_closes"`
	GateOneWay             bool     `json:"gate_one_way"`
	GateReopenable         bool     `json:"gate_reopenable"`
}

// ItemBagOfTricksCreature is one row of a Bag of Tricks 1d8 table.
type ItemBagOfTricksCreature struct {
	Roll     int    `json:"roll"`
	Creature string `json:"creature"`
}

// ItemBagOfTricks retains typed SRD source facts without implementing them.
type ItemBagOfTricks struct {
	Color           string                    `json:"color"`
	ClothColors     []string                  `json:"cloth_colors"`
	Appearance      string                    `json:"appearance"`
	Contents        string                    `json:"contents"`
	Action          string                    `json:"action"`
	ThrowFeet       int                       `json:"throw_feet"`
	Transforms      string                    `json:"transforms"`
	TableDieSides   int                       `json:"table_die_sides"`
	Creatures       []ItemBagOfTricksCreature `json:"creatures"`
	StatBlockSource string                    `json:"stat_block_source"`
	Vanishes        []string                  `json:"vanishes"`
	Attitude        string                    `json:"attitude"`
	Initiative      string                    `json:"initiative"`
	CommandAction   string                    `json:"command_action"`
	CommandScope    string                    `json:"command_scope"`
	WithoutOrders   string                    `json:"without_orders"`
	Uses            int                       `json:"uses"`
	UsesReset       string                    `json:"uses_reset"`
}

// ItemBeadOfForce retains typed SRD source facts without implementing them.
type ItemBeadOfForce struct {
	SourcePages               []int    `json:"source_pages"`
	Appearance                string   `json:"appearance"`
	DiameterInchesNumerator   int      `json:"diameter_inches_numerator"`
	DiameterInchesDenominator int      `json:"diameter_inches_denominator"`
	WeightOunces              int      `json:"weight_ounces"`
	FoundDice                 int      `json:"found_dice"`
	FoundDieSides             int      `json:"found_die_sides"`
	FoundBonus                int      `json:"found_bonus"`
	Action                    string   `json:"action"`
	ThrowFeet                 int      `json:"throw_feet"`
	AreaShape                 string   `json:"area_shape"`
	AreaRadiusFeet            int      `json:"area_radius_feet"`
	DestroyedOnImpact         bool     `json:"destroyed_on_impact"`
	Save                      string   `json:"save"`
	SaveDC                    int      `json:"save_dc"`
	Damage                    ItemDice `json:"damage"`
	SaveSuccess               string   `json:"save_success"`
	SphereDurationSeconds     int      `json:"sphere_duration_seconds"`
	TrappedIfAll              []string `json:"trapped_if_all"`
	PushedIfAny               []string `json:"pushed_if_any"`
	PushEffect                string   `json:"push_effect"`
	PassesThrough             []string `json:"passes_through"`
	Blocked                   []string `json:"blocked"`
	MoveAction                string   `json:"move_action"`
	MoveActor                 string   `json:"move_actor"`
	MoveDistance              string   `json:"move_distance"`
	CanBePickedUp             bool     `json:"can_be_picked_up"`
	SphereWeightPounds        int      `json:"sphere_weight_pounds"`
	WeightIgnoresCreatures    bool     `json:"weight_ignores_creatures"`
}

// ItemBeadOfNourishment retains typed SRD source facts without implementing them.
type ItemBeadOfNourishment struct {
	Appearance      string `json:"appearance"`
	ConsumedBy      string `json:"consumed_by"`
	NourishmentDays int    `json:"nourishment_days"`
	Equivalent      string `json:"equivalent"`
}

// ItemBeltOfDwarvenkind retains typed SRD source facts without implementing them.
type ItemBeltOfDwarvenkind struct {
	Condition            string   `json:"condition"`
	Benefits             []string `json:"benefits"`
	Language             string   `json:"language"`
	CheckAbility         string   `json:"check_ability"`
	CheckSkill           string   `json:"check_skill"`
	CheckAudience        []string `json:"check_audience"`
	ConstitutionIncrease int      `json:"constitution_increase"`
	ConstitutionMaximum  int      `json:"constitution_maximum"`
	BeardCondition       string   `json:"beard_condition"`
	BeardChancePercent   int      `json:"beard_chance_percent"`
	BeardTiming          string   `json:"beard_timing"`
	BeardEffects         []string `json:"beard_effects"`
	AdditionalCondition  string   `json:"additional_condition"`
	AdditionalBenefits   []string `json:"additional_benefits"`
	DarkvisionFeet       int      `json:"darkvision_feet"`
	Resistance           string   `json:"resistance"`
	SaveAdvantage        string   `json:"save_advantage"`
}

// ItemBeltOfGiantStrength retains typed SRD source facts without implementing them.
type ItemBeltOfGiantStrength struct {
	Condition              string   `json:"condition"`
	GiantTypes             []string `json:"giant_types"`
	Ability                string   `json:"ability"`
	Score                  int      `json:"score"`
	ScoreDeterminedBy      string   `json:"score_determined_by"`
	NoEffectIfScoreAtLeast int      `json:"no_effect_if_score_at_least"`
}

// ItemBerserkerAxe retains typed SRD source facts without implementing them.
type ItemBerserkerAxe struct {
	AttackBonus             int      `json:"attack_bonus"`
	DamageBonus             int      `json:"damage_bonus"`
	HitPointMaximumPerLevel int      `json:"hit_point_maximum_per_level"`
	HitPointCondition       string   `json:"hit_point_condition"`
	Cursed                  bool     `json:"cursed"`
	CurseTrigger            string   `json:"curse_trigger"`
	CurseEffects            []string `json:"curse_effects"`
	BerserkTrigger          string   `json:"berserk_trigger"`
	BerserkSave             string   `json:"berserk_save"`
	BerserkSaveDC           int      `json:"berserk_save_dc"`
	BerserkSaveFailure      string   `json:"berserk_save_failure"`
	BerserkEnds             string   `json:"berserk_ends"`
	BerserkEndRangeFeet     int      `json:"berserk_end_range_feet"`
	BerserkTarget           string   `json:"berserk_target"`
	BerserkTieBreak         string   `json:"berserk_tie_break"`
	BerserkCompulsion       []string `json:"berserk_compulsion"`
	BerserkUnreachable      string   `json:"berserk_unreachable"`
	BerserkRetargetTriggers []string `json:"berserk_retarget_triggers"`
	BerserkRetarget         string   `json:"berserk_retarget"`
}

// ItemBootsOfElvenkind retains typed SRD source facts without implementing them.
type ItemBootsOfElvenkind struct {
	Condition        string `json:"condition"`
	SilentSteps      bool   `json:"silent_steps"`
	AnySurface       bool   `json:"any_surface"`
	AdvantageAbility string `json:"advantage_ability"`
	AdvantageSkill   string `json:"advantage_skill"`
}

// ItemBootsOfLevitation retains typed SRD source facts without implementing them.
type ItemBootsOfLevitation struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
	Target    string `json:"target"`
}

// itemManifest226 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Bag of Tricks has one variant per cloth colour table
// and Belt of Giant Strength one per printed table row; table rows inside a
// payload are data. Berserker Axe's three base weapons are base options.
var itemManifest226 = itemSourceManifest{Batch: "data-magic-items-02", Issue: 226, Families: []itemSourceManifestFamily{
	{Slug: "arrow-catching-shield", Name: "Arrow-Catching Shield", Page: 211, Category: "armor", Options: "shield", Attunement: "required", Equip: "wield", Use: "reaction", UseRestrictions: "ranged-attack-roll-against-target-within-5-feet",
		Keys: "item.attunement item.wield-shield defense.armor-class-bonus-against-ranged-attacks reaction.become-target-of-ranged-attack", Payload: "arrow_catching",
		Variants: []itemSourceManifestVariant{{"standard", "Arrow-Catching Shield", "rare"}}},
	{Slug: "bag-of-beans", Name: "Bag of Beans", Page: 211, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "dump-or-plant", UseRestrictions: "dump-one-or-more-beans plant-in-dirt-or-sand-then-water",
		Keys: "item.contains-3d4-beans item.nonmagical-when-empty area.sphere-10-foot-radius save.dexterity-dc-15-half damage.force-5d4 item.planted-bean-delayed-effect item.bean-effect-table", Payload: "beans",
		Variants: []itemSourceManifestVariant{{"standard", "Bag of Beans", "rare"}}},
	{Slug: "bag-of-devouring", Name: "Bag of Devouring", Page: 212, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "container",
		Keys: "item.extradimensional-creature-orifice item.devours-animal-and-vegetable-matter chance.pulled-inside-50-percent check.strength-athletics-dc-15 check.strength-athletics-dc-20 item.devours-creature-at-turn-start item.capacity-one-cubic-foot item.daily-planar-ejection item.destroyed-if-pierced-or-torn transport.random-astral-plane-location", Payload: "devouring",
		Variants: []itemSourceManifestVariant{{"standard", "Bag of Devouring", "very-rare"}}},
	{Slug: "bag-of-holding", Name: "Bag of Holding", Page: 212, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "retrieve-item",
		Keys: "item.extradimensional-storage item.capacity-500-pounds-64-cubic-feet item.fixed-weight-5-pounds action.utilize item.destroyed-if-overloaded-pierced-or-torn item.inside-out-spills-contents item.air-supply item.extradimensional-nesting-opens-astral-gate area.sphere-10-foot-radius transport.random-astral-plane-location", Payload: "holding",
		Variants: []itemSourceManifestVariant{{"standard", "Bag of Holding", "uncommon"}}},
	{Slug: "bag-of-tricks", Name: "Bag of Tricks", Page: 212, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "three-objects-until-next-dawn",
		Keys: "action.magic item.throw-object-20-feet item.creature-table-by-color summon.creature-until-dawn-or-0-hit-points summon.friendly-acts-after-user action.bonus-action item.three-uses-until-next-dawn", Payload: "tricks",
		Variants: []itemSourceManifestVariant{{"gray", "Gray Bag of Tricks", "uncommon"}, {"rust", "Rust Bag of Tricks", "uncommon"}, {"tan", "Tan Bag of Tricks", "uncommon"}}},
	{Slug: "bead-of-force", Name: "Bead of Force", Page: 212, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "throw-up-to-60-feet bead-destroyed",
		Keys: "action.magic item.throw-60-feet item.destroyed-on-use area.sphere-10-foot-radius save.dexterity-dc-15 damage.force-5d4 effect.force-sphere-one-minute effect.trapped-or-pushed-out action.utilize", Payload: "force",
		Variants: []itemSourceManifestVariant{{"standard", "Bead of Force", "rare"}}},
	{Slug: "bead-of-nourishment", Name: "Bead of Nourishment", Page: 213, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "consume-one",
		Keys: "item.consume sustenance.one-day-of-rations", Payload: "nourishment",
		Variants: []itemSourceManifestVariant{{"standard", "Bead of Nourishment", "common"}}},
	{Slug: "belt-of-dwarvenkind", Name: "Belt of Dwarvenkind", Page: 213, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear language.dwarvish check.advantage-persuasion-with-dwarves-and-duergar ability.constitution-plus-2-maximum-20 chance.daily-beard-50-percent species.not-dwarf-or-duergar senses.darkvision-60-feet defense.resistance-poison save.advantage-against-poisoned", Payload: "dwarvenkind",
		Variants: []itemSourceManifestVariant{{"standard", "Belt of Dwarvenkind", "rare"}}},
	{Slug: "belt-of-giant-strength", Name: "Belt of Giant Strength", Page: 213, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower", Payload: "giant_strength",
		Variants: []itemSourceManifestVariant{{"hill", "Belt of Giant Strength (hill)", "rare"}, {"frost-or-stone", "Belt of Giant Strength (frost or stone)", "very-rare"}, {"fire", "Belt of Giant Strength (fire)", "very-rare"}, {"cloud", "Belt of Giant Strength (cloud)", "legendary"}, {"storm", "Belt of Giant Strength (storm)", "legendary"}}},
	{Slug: "berserker-axe", Name: "Berserker Axe", Page: 213, Category: "weapon", Options: "battleaxe greataxe halberd", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon attack.weapon-bonus damage.weapon-bonus hit-points.maximum-plus-1-per-level item.curse item.curse-extends-on-attunement item.curse-unwilling-to-part attack.disadvantage-with-other-weapons save.wisdom-dc-15 condition.berserk-compelled-attack", Payload: "berserker",
		Variants: []itemSourceManifestVariant{{"standard", "Berserker Axe", "rare"}}},
	{Slug: "boots-of-elvenkind", Name: "Boots of Elvenkind", Page: 213, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear movement.silent-steps check.advantage-dexterity-stealth", Payload: "elvenkind",
		Variants: []itemSourceManifestVariant{{"standard", "Boots of Elvenkind", "uncommon"}}},
	{Slug: "boots-of-levitation", Name: "Boots of Levitation", Page: 213, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "levitate self-only",
		Keys: "item.attunement item.wear spell.cast-levitate-on-self", Payload: "levitation",
		Variants: []itemSourceManifestVariant{{"standard", "Boots of Levitation", "rare"}}},
}}

// LoadItems226 returns a fresh validated copy of batch #226. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems226() (ItemSourceBatch[ItemDetails226], error) { return DecodeItems226(items226JSON) }

// DecodeItems226 validates supplied JSON as batch #226 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems226(data []byte) (ItemSourceBatch[ItemDetails226], error) {
	return decodeItemSourceBatch(data, itemManifest226, validateItemDetails226)
}

// validateItemDetails226 checks the internal consistency of each payload: dice,
// DCs and tables must be complete, and a table variant must carry its own row.
func validateItemDetails226(family, variant string, d *ItemDetails226) error {
	dice := func(x ItemDice) bool { return x.Dice > 0 && x.Sides > 0 }
	switch family {
	case "arrow-catching-shield":
		if a := d.ArrowCatching; a.ACBonus <= 0 || a.RedirectRangeFeet <= 0 {
			return fmt.Errorf("invalid bonus or range")
		}
	case "bag-of-beans":
		b := d.Beans
		rows, effects := [][2]int{}, []string{}
		for _, row := range b.Effects {
			rows, effects = append(rows, [2]int{row.Min, row.Max}), append(effects, row.Effect)
			labels := []string{}
			for _, q := range row.Quantities {
				if q.Dice > 0 != (q.Sides > 0) || q.Dice == 0 && q.Bonus == 0 {
					return fmt.Errorf("incomplete bean effect quantity")
				}
				labels = append(labels, q.Label)
			}
			if !uniqueItemKeys(labels) || len(row.Creatures) > 0 && !uniqueItemKeys(row.Creatures) || len(row.Facts) > 0 && !uniqueItemKeys(row.Facts) {
				return fmt.Errorf("duplicate bean effect quantity, creature or fact")
			}
		}
		if !itemRollTable(rows, b.TableDieSides) || !uniqueItemKeys(effects) || b.BeansDice <= 0 || b.BeansDieSides <= 0 ||
			b.WeightPoundsNumerator <= 0 || b.WeightPoundsDenominator <= 0 || b.AreaRadiusFeet <= 0 || b.SaveDC <= 0 || !dice(b.Damage) ||
			b.PlantDelaySeconds <= 0 || !uniqueItemKeys(b.PlantSteps) || !uniqueItemKeys(b.EffectSelection) {
			return fmt.Errorf("incomplete bean effect table, save or damage")
		}
	case "bag-of-devouring":
		v := d.Devouring
		if v.PullChancePercent <= 0 || v.PullChancePercent > 100 || v.EscapeCheck.DC <= 0 || v.RescueCheck.DC <= 0 || v.CapacityCubicFeet <= 0 ||
			!uniqueItemKeys(v.Devours) || !uniqueItemKeys(v.DestroyedBy) {
			return fmt.Errorf("invalid chance, checks or capacity")
		}
	case "bag-of-holding":
		h := d.Holding
		if h.InteriorSquareSideFeet <= 0 || h.InteriorDepthFeet <= 0 || h.CapacityPounds <= 0 || h.CapacityCubicFeet <= 0 || h.WeightPounds <= 0 ||
			h.AirMinutes <= 0 || h.GateRadiusFeet <= 0 || !uniqueItemKeys(h.DestroyedBy) || !uniqueItemKeys(h.NestingSpaces) || !uniqueItemKeys(h.NestingEffects) {
			return fmt.Errorf("invalid dimensions, capacity or gate")
		}
	case "bag-of-tricks":
		t := d.Tricks
		creatures := []string{}
		for i, row := range t.Creatures {
			if row.Roll != i+1 {
				return fmt.Errorf("incomplete creature table")
			}
			creatures = append(creatures, row.Creature)
		}
		if t.Color != variant || !slices.Contains(t.ClothColors, variant) || !uniqueItemKeys(t.ClothColors) || len(creatures) != t.TableDieSides ||
			!uniqueItemKeys(creatures) || t.ThrowFeet <= 0 || t.Uses <= 0 || !uniqueItemKeys(t.Vanishes) {
			return fmt.Errorf("creature table disagrees with its colour variant or die")
		}
	case "bead-of-force":
		f := d.Force
		if !slices.Equal(f.SourcePages, []int{212, 213}) || f.DiameterInchesNumerator <= 0 || f.DiameterInchesDenominator <= 0 || f.WeightOunces <= 0 ||
			f.FoundDice <= 0 || f.FoundDieSides <= 0 || f.ThrowFeet <= 0 || f.AreaRadiusFeet <= 0 || f.SaveDC <= 0 || !dice(f.Damage) ||
			f.SphereDurationSeconds <= 0 || f.SphereWeightPounds <= 0 || !uniqueItemKeys(f.TrappedIfAll) || !uniqueItemKeys(f.PushedIfAny) ||
			!uniqueItemKeys(f.PassesThrough) || !uniqueItemKeys(f.Blocked) {
			return fmt.Errorf("invalid continuation pages, save, damage or sphere")
		}
	case "bead-of-nourishment":
		if d.Nourishment.NourishmentDays <= 0 {
			return fmt.Errorf("invalid nourishment")
		}
	case "belt-of-dwarvenkind":
		b := d.Dwarvenkind
		if b.ConstitutionIncrease <= 0 || b.ConstitutionMaximum < b.ConstitutionIncrease || b.ConstitutionMaximum > 30 || b.BeardChancePercent <= 0 || b.BeardChancePercent > 100 ||
			b.DarkvisionFeet <= 0 || !uniqueItemKeys(b.Benefits) || !uniqueItemKeys(b.AdditionalBenefits) || !uniqueItemKeys(b.CheckAudience) || !uniqueItemKeys(b.BeardEffects) {
			return fmt.Errorf("invalid benefits, score or chance")
		}
	case "belt-of-giant-strength":
		// The belt table, p. 213: each row's giant types, score and rarity slug.
		rows := map[string]struct {
			giants []string
			score  int
		}{"hill": {[]string{"hill"}, 21}, "frost-or-stone": {[]string{"frost", "stone"}, 23}, "fire": {[]string{"fire"}, 25}, "cloud": {[]string{"cloud"}, 27}, "storm": {[]string{"storm"}, 29}}
		row, ok := rows[variant]
		if g := d.GiantStrength; !ok || !slices.Equal(g.GiantTypes, row.giants) || g.Score != row.score || g.NoEffectIfScoreAtLeast != g.Score {
			return fmt.Errorf("score disagrees with its table row variant")
		}
	case "berserker-axe":
		a := d.Berserker
		if a.AttackBonus <= 0 || a.DamageBonus != a.AttackBonus || a.HitPointMaximumPerLevel <= 0 || !a.Cursed || a.BerserkSaveDC <= 0 || a.BerserkEndRangeFeet <= 0 ||
			!uniqueItemKeys(a.CurseEffects) || !uniqueItemKeys(a.BerserkCompulsion) || !uniqueItemKeys(a.BerserkRetargetTriggers) {
			return fmt.Errorf("invalid bonus, curse or save")
		}
	}
	return nil
}
