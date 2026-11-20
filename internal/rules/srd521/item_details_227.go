package srd521

import (
	_ "embed"
	"fmt"
)

//go:embed testdata/items_227.json
var items227JSON []byte

// ItemDetails227 holds exactly one typed payload per variant of catalogue batch
// #227 (magic items 03, SRD 5.2.1 pp. 214–215). Payloads are source facts only.
type ItemDetails227 struct {
	BootsSpeed       *ItemBootsSpeed       `json:"boots_speed,omitempty"`
	BootsStriding    *ItemBootsStriding    `json:"boots_striding,omitempty"`
	BootsWinterlands *ItemBootsWinterlands `json:"boots_winterlands,omitempty"`
	Bowl             *ItemElementalBowl    `json:"bowl,omitempty"`
	BracersArchery   *ItemBracersArchery   `json:"bracers_archery,omitempty"`
	BracersDefense   *ItemBracersDefense   `json:"bracers_defense,omitempty"`
	Brazier          *ItemElementalBrazier `json:"brazier,omitempty"`
	Brooch           *ItemBroochShielding  `json:"brooch,omitempty"`
	Broom            *ItemBroomFlying      `json:"broom,omitempty"`
	Candle           *ItemCandleInvocation `json:"candle,omitempty"`
	Cape             *ItemCapeMountebank   `json:"cape,omitempty"`
	Carpet           *ItemCarpetFlying     `json:"carpet,omitempty"`
}

// ItemBootsSpeed retains typed SRD source facts without implementing them.
// The ten-minute total is a use limit until a Long Rest, not a charge pool.
type ItemBootsSpeed struct {
	Condition                       string `json:"condition"`
	ActivateAction                  string `json:"activate_action"`
	Activation                      string `json:"activation"`
	SpeedMultiplier                 int    `json:"speed_multiplier"`
	OpportunityAttacksAgainstWearer string `json:"opportunity_attacks_against_wearer"`
	EndTrigger                      string `json:"end_trigger"`
	TotalUseSeconds                 int    `json:"total_use_seconds"`
	AfterTotalUse                   string `json:"after_total_use"`
	Recharge                        string `json:"recharge"`
}

// ItemBootsStriding retains typed SRD source facts without implementing them.
type ItemBootsStriding struct {
	Condition            string   `json:"condition"`
	SpeedBecomesFeet     int      `json:"speed_becomes_feet"`
	SpeedUnless          string   `json:"speed_unless"`
	SpeedNotReducedBy    []string `json:"speed_not_reduced_by"`
	JumpFrequency        string   `json:"jump_frequency"`
	JumpDistanceFeet     int      `json:"jump_distance_feet"`
	JumpMovementCostFeet int      `json:"jump_movement_cost_feet"`
}

// ItemWinterlandsColdResistance is the Cold Resistance benefit. The tolerated
// temperature is the printed 0 degrees Fahrenheit, so zero is a source value.
type ItemWinterlandsColdResistance struct {
	Benefit                        string `json:"benefit"`
	Resistance                     string `json:"resistance"`
	ToleratedTemperatureFahrenheit int    `json:"tolerated_temperature_fahrenheit"`
	TemperatureBound               string `json:"temperature_bound"`
	Without                        string `json:"without"`
}

// ItemWinterlandsWinterStrider is the Winter Strider benefit.
type ItemWinterlandsWinterStrider struct {
	Benefit   string   `json:"benefit"`
	Ignores   string   `json:"ignores"`
	CreatedBy []string `json:"created_by"`
}

// ItemBootsWinterlands retains typed SRD source facts without implementing them.
type ItemBootsWinterlands struct {
	Appearance     string                        `json:"appearance"`
	Condition      string                        `json:"condition"`
	ColdResistance ItemWinterlandsColdResistance `json:"cold_resistance"`
	WinterStrider  ItemWinterlandsWinterStrider  `json:"winter_strider"`
}

// ItemBowlBrazierSummon is the summoning property the Bowl of Commanding Water
// Elementals and the Brazier of Commanding Fire Elementals print in the same
// terms. Once until the next dawn is a use limit, not a charge pool.
type ItemBowlBrazierSummon struct {
	WithinFeet      int      `json:"within_feet"`
	Action          string   `json:"action"`
	Summons         string   `json:"summons"`
	Appears         string   `json:"appears"`
	ElementalTraits []string `json:"elemental_traits"`
	DurationSeconds int      `json:"duration_seconds"`
	EndsWhen        []string `json:"ends_when"`
	DismissAction   string   `json:"dismiss_action"`
	Uses            int      `json:"uses"`
	Recharge        string   `json:"recharge"`
}

// ItemElementalBowl retains typed SRD source facts without implementing them.
type ItemElementalBowl struct {
	VesselCondition       string                `json:"vessel_condition"`
	Summon                ItemBowlBrazierSummon `json:"summon"`
	DimensionsApproximate bool                  `json:"dimensions_approximate"`
	DiameterFeet          int                   `json:"diameter_feet"`
	Depth                 string                `json:"depth"`
	CapacityGallons       int                   `json:"capacity_gallons"`
}

// ItemBracersArchery retains typed SRD source facts without implementing them.
type ItemBracersArchery struct {
	Condition           string   `json:"condition"`
	WeaponProficiencies []string `json:"weapon_proficiencies"`
	DamageBonus         int      `json:"damage_bonus"`
	DamageBonusWeapons  []string `json:"damage_bonus_weapons"`
}

// ItemBracersDefense retains typed SRD source facts without implementing them.
type ItemBracersDefense struct {
	Condition string   `json:"condition"`
	ACBonus   int      `json:"ac_bonus"`
	Requires  []string `json:"requires"`
}

// ItemElementalBrazier retains typed SRD source facts without implementing them.
type ItemElementalBrazier struct {
	Summon ItemBowlBrazierSummon `json:"summon"`
}

// ItemBroochShielding retains typed SRD source facts without implementing them.
type ItemBroochShielding struct {
	Condition            string   `json:"condition"`
	Resistances          []string `json:"resistances"`
	ImmunityToDamageFrom string   `json:"immunity_to_damage_from"`
}

// ItemBroomFlying retains typed SRD source facts without implementing them.
type ItemBroomFlying struct {
	Appearance            string   `json:"appearance"`
	Inactive              string   `json:"inactive"`
	HoverRequirement      string   `json:"hover_requirement"`
	HoverAction           string   `json:"hover_action"`
	FlySpeedFeet          int      `json:"fly_speed_feet"`
	CapacityPounds        int      `json:"capacity_pounds"`
	ReducedFlySpeedFeet   int      `json:"reduced_fly_speed_feet"`
	ReducedWhenOverPounds int      `json:"reduced_when_over_pounds"`
	HoverEnds             []string `json:"hover_ends"`
	SendAction            string   `json:"send_action"`
	SendMode              string   `json:"send_mode"`
	SendRangeMiles        int      `json:"send_range_miles"`
	SendRequirements      []string `json:"send_requirements"`
	RecallAction          string   `json:"recall_action"`
	RecallRequirements    []string `json:"recall_requirements"`
	RecallRangeMiles      int      `json:"recall_range_miles"`
}

// ItemCandlePlane is one row of the Candle of Invocation 1d100 table.
type ItemCandlePlane struct {
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Plane string `json:"plane"`
}

// ItemCandleInvocation retains typed SRD source facts without implementing
// them. The entry starts on p. 214; its Outer Plane table is on TablePage.
// Burn time is a remaining duration in the source, not a charge pool.
type ItemCandleInvocation struct {
	Activation                    string            `json:"activation"`
	ActivateAction                string            `json:"activate_action"`
	BurnSeconds                   int               `json:"burn_seconds"`
	AfterBurn                     string            `json:"after_burn"`
	SnuffEarly                    bool              `json:"snuff_early"`
	BurnDeductionIncrementSeconds int               `json:"burn_deduction_increment_seconds"`
	Light                         string            `json:"light"`
	LightRadiusFeet               int               `json:"light_radius_feet"`
	YouInLight                    string            `json:"you_in_light"`
	FreeCasters                   []string          `json:"free_casters"`
	FreeCastSpellLevel            int               `json:"free_cast_spell_level"`
	FreeCastRequirement           string            `json:"free_cast_requirement"`
	FreeCastCost                  string            `json:"free_cast_cost"`
	GateCondition                 string            `json:"gate_condition"`
	GateCasts                     string            `json:"gate_casts"`
	GateDestroysCandle            bool              `json:"gate_destroys_candle"`
	GateLinksTo                   string            `json:"gate_links_to"`
	PlaneSelection                []string          `json:"plane_selection"`
	TableDieSides                 int               `json:"table_die_sides"`
	TablePage                     int               `json:"table_page"`
	OuterPlanes                   []ItemCandlePlane `json:"outer_planes"`
}

// ItemCapeMountebank retains typed SRD source facts without implementing them.
type ItemCapeMountebank struct {
	Scent         string `json:"scent"`
	Condition     string `json:"condition"`
	Casts         string `json:"casts"`
	Action        string `json:"action"`
	Uses          int    `json:"uses"`
	Recharge      string `json:"recharge"`
	OnTeleport    string `json:"on_teleport"`
	VacatedSpace  string `json:"vacated_space"`
	ObscuredUntil string `json:"obscured_until"`
}

// ItemCarpetSize is one row of the Carpet of Flying 1d100 size table.
type ItemCarpetSize struct {
	Min            int `json:"min"`
	Max            int `json:"max"`
	WidthFeet      int `json:"width_feet"`
	LengthFeet     int `json:"length_feet"`
	CapacityPounds int `json:"capacity_pounds"`
	FlySpeedFeet   int `json:"fly_speed_feet"`
}

// ItemCarpetFlying retains typed SRD source facts without implementing them.
// The four sizes are rows of the source table, not separate variants.
type ItemCarpetFlying struct {
	ActivateAction        string           `json:"activate_action"`
	ActivationRequirement string           `json:"activation_requirement"`
	Effect                string           `json:"effect"`
	Control               string           `json:"control"`
	ControlRangeFeet      int              `json:"control_range_feet"`
	SizeCount             int              `json:"size_count"`
	SizeSelection         []string         `json:"size_selection"`
	TableDieSides         int              `json:"table_die_sides"`
	Sizes                 []ItemCarpetSize `json:"sizes"`
	MaxLoadMultiplier     int              `json:"max_load_multiplier"`
	OverloadedFlySpeed    string           `json:"overloaded_fly_speed"`
	OverloadedWhen        string           `json:"overloaded_when"`
}

// itemManifest227 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Every family is a Wondrous Item with one printed
// rarity, so each has the single variant standard; table rows are data.
var itemManifest227 = itemSourceManifest{Batch: "data-magic-items-03", Issue: 227, Families: []itemSourceManifestFamily{
	{Slug: "boots-of-speed", Name: "Boots of Speed", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "bonus-action", UseRestrictions: "ten-minutes-total-until-long-rest",
		Keys: "item.attunement item.wear action.bonus-action speed.doubled attack.opportunity-attacks-against-wearer-disadvantage duration.ends-on-repeat-activation item.ten-minutes-total-until-long-rest", Payload: "boots_speed",
		Variants: []itemSourceManifestVariant{{"standard", "Boots of Speed", "rare"}}},
	{Slug: "boots-of-striding-and-springing", Name: "Boots of Striding and Springing", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear speed.becomes-30-feet-unless-higher speed.ignores-carrying-and-heavy-armor-reduction movement.jump-30-feet-for-10-feet-once-per-turn", Payload: "boots_striding",
		Variants: []itemSourceManifestVariant{{"standard", "Boots of Striding and Springing", "uncommon"}}},
	{Slug: "boots-of-the-winterlands", Name: "Boots of the Winterlands", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.resistance-cold environment.tolerates-0-fahrenheit-or-lower movement.ignores-ice-and-snow-difficult-terrain", Payload: "boots_winterlands",
		Variants: []itemSourceManifestVariant{{"standard", "Boots of the Winterlands", "uncommon"}}},
	{Slug: "bowl-of-commanding-water-elementals", Name: "Bowl of Commanding Water Elementals", Page: 214, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "bowl-filled-with-water within-5-feet-of-bowl once-until-next-dawn",
		Keys: "action.magic summon.water-elemental summon.acts-after-summoner-initiative duration.one-hour-death-or-dismissal action.bonus-action-dismiss item.once-until-next-dawn", Payload: "bowl",
		Variants: []itemSourceManifestVariant{{"standard", "Bowl of Commanding Water Elementals", "rare"}}},
	{Slug: "bracers-of-archery", Name: "Bracers of Archery", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear proficiency.longbow-and-shortbow damage.longbow-and-shortbow-bonus-2", Payload: "bracers_archery",
		Variants: []itemSourceManifestVariant{{"standard", "Bracers of Archery", "uncommon"}}},
	{Slug: "bracers-of-defense", Name: "Bracers of Defense", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.armor-class-bonus-2-without-armor-or-shield", Payload: "bracers_defense",
		Variants: []itemSourceManifestVariant{{"standard", "Bracers of Defense", "rare"}}},
	{Slug: "brazier-of-commanding-fire-elementals", Name: "Brazier of Commanding Fire Elementals", Page: 214, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "within-5-feet-of-brazier once-until-next-dawn",
		Keys: "action.magic summon.fire-elemental summon.acts-after-summoner-initiative duration.one-hour-death-or-dismissal action.bonus-action-dismiss item.once-until-next-dawn", Payload: "brazier",
		Variants: []itemSourceManifestVariant{{"standard", "Brazier of Commanding Fire Elementals", "rare"}}},
	{Slug: "brooch-of-shielding", Name: "Brooch of Shielding", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.resistance-force defense.immunity-magic-missile-damage", Payload: "brooch",
		Variants: []itemSourceManifestVariant{{"standard", "Brooch of Shielding", "uncommon"}}},
	{Slug: "broom-of-flying", Name: "Broom of Flying", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "stand-astride-broom",
		Keys: "item.attunement action.magic movement.ridden-fly-speed-50-feet item.carrying-capacity-400-pounds movement.fly-speed-30-feet-over-200-pounds item.hover-ends-on-landing-or-dismount item.send-to-named-familiar-location-within-1-mile item.recall-by-command-word-within-1-mile", Payload: "broom",
		Variants: []itemSourceManifestVariant{{"standard", "Broom of Flying", "uncommon"}}},
	{Slug: "candle-of-invocation", Name: "Candle of Invocation", Page: 214, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "four-hours-total-burn-time",
		Keys: "item.attunement action.magic item.burn-time-four-hours light.dim-30-feet d20.advantage-while-in-light spell.cleric-or-druid-level-1-without-slot spell.cast-gate-on-first-lighting item.destroyed-by-gate-or-burn item.outer-plane-table", Payload: "candle",
		Variants: []itemSourceManifestVariant{{"standard", "Candle of Invocation", "very-rare"}}},
	{Slug: "cape-of-the-mountebank", Name: "Cape of the Mountebank", Page: 215, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "magic-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.wear action.magic spell.cast-dimension-door item.once-until-next-dawn area.vacated-space-lightly-obscured", Payload: "cape",
		Variants: []itemSourceManifestVariant{{"standard", "Cape of the Mountebank", "rare"}}},
	{Slug: "carpet-of-flying", Name: "Carpet of Flying", Page: 215, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "command-word",
		Keys: "action.magic item.command-word movement.hover-and-fly item.directed-within-30-feet item.carpet-size-table item.carrying-capacity-double-at-half-fly-speed", Payload: "carpet",
		Variants: []itemSourceManifestVariant{{"standard", "Carpet of Flying", "very-rare"}}},
}}

// LoadItems227 returns a fresh validated copy of batch #227. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems227() (ItemSourceBatch[ItemDetails227], error) { return DecodeItems227(items227JSON) }

// DecodeItems227 validates supplied JSON as batch #227 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems227(data []byte) (ItemSourceBatch[ItemDetails227], error) {
	return decodeItemSourceBatch(data, itemManifest227, validateItemDetails227)
}

// validateItemDetails227 checks the internal consistency of each payload:
// tables must cover their die, lists must be distinct and every printed
// distance, duration, bonus and limit must be positive.
func validateItemDetails227(family, _ string, d *ItemDetails227) error {
	summon := func(s ItemBowlBrazierSummon) bool {
		return s.WithinFeet > 0 && s.DurationSeconds > 0 && s.Uses > 0 && uniqueItemKeys(s.ElementalTraits) && uniqueItemKeys(s.EndsWhen)
	}
	switch family {
	case "boots-of-speed":
		if b := d.BootsSpeed; b.SpeedMultiplier < 2 || b.TotalUseSeconds <= 0 {
			return fmt.Errorf("invalid speed multiplier or total use")
		}
	case "boots-of-striding-and-springing":
		if b := d.BootsStriding; b.SpeedBecomesFeet <= 0 || b.JumpDistanceFeet <= 0 || b.JumpMovementCostFeet <= 0 || b.JumpMovementCostFeet >= b.JumpDistanceFeet || !uniqueItemKeys(b.SpeedNotReducedBy) {
			return fmt.Errorf("invalid speed, jump or reduction list")
		}
	case "boots-of-the-winterlands":
		if b := d.BootsWinterlands; b.ColdResistance.Benefit == b.WinterStrider.Benefit || !uniqueItemKeys(b.WinterStrider.CreatedBy) {
			return fmt.Errorf("duplicate benefit or terrain source")
		}
	case "bowl-of-commanding-water-elementals":
		if b := d.Bowl; !summon(b.Summon) || b.DiameterFeet <= 0 || b.CapacityGallons <= 0 {
			return fmt.Errorf("incomplete summon or dimensions")
		}
	case "bracers-of-archery":
		if b := d.BracersArchery; b.DamageBonus <= 0 || !uniqueItemKeys(b.WeaponProficiencies) || !uniqueItemKeys(b.DamageBonusWeapons) {
			return fmt.Errorf("invalid damage bonus or weapon list")
		}
	case "bracers-of-defense":
		if b := d.BracersDefense; b.ACBonus <= 0 || !uniqueItemKeys(b.Requires) {
			return fmt.Errorf("invalid Armor Class bonus or requirements")
		}
	case "brazier-of-commanding-fire-elementals":
		if !summon(d.Brazier.Summon) {
			return fmt.Errorf("incomplete summon")
		}
	case "brooch-of-shielding":
		if !uniqueItemKeys(d.Brooch.Resistances) {
			return fmt.Errorf("duplicate resistance")
		}
	case "broom-of-flying":
		b := d.Broom
		if b.FlySpeedFeet <= 0 || b.ReducedFlySpeedFeet <= 0 || b.ReducedFlySpeedFeet >= b.FlySpeedFeet || b.CapacityPounds <= 0 || b.ReducedWhenOverPounds <= 0 || b.ReducedWhenOverPounds >= b.CapacityPounds ||
			b.SendRangeMiles <= 0 || b.RecallRangeMiles <= 0 || !uniqueItemKeys(b.HoverEnds) || !uniqueItemKeys(b.SendRequirements) || !uniqueItemKeys(b.RecallRequirements) {
			return fmt.Errorf("invalid speed, capacity, range or requirements")
		}
	case "candle-of-invocation":
		c := d.Candle
		rows, planes := [][2]int{}, []string{}
		for _, row := range c.OuterPlanes {
			rows, planes = append(rows, [2]int{row.Min, row.Max}), append(planes, row.Plane)
		}
		if !itemRollTable(rows, c.TableDieSides) || !uniqueItemKeys(planes) || !uniqueItemKeys(c.PlaneSelection) || !uniqueItemKeys(c.FreeCasters) ||
			c.BurnSeconds <= 0 || c.BurnDeductionIncrementSeconds <= 0 || c.BurnDeductionIncrementSeconds > c.BurnSeconds || c.LightRadiusFeet <= 0 || c.FreeCastSpellLevel <= 0 || c.TablePage <= 0 {
			return fmt.Errorf("incomplete Outer Plane table, burn time, light or casting")
		}
	case "cape-of-the-mountebank":
		if d.Cape.Uses <= 0 {
			return fmt.Errorf("invalid uses")
		}
	case "carpet-of-flying":
		c := d.Carpet
		rows := [][2]int{}
		for _, row := range c.Sizes {
			rows = append(rows, [2]int{row.Min, row.Max})
			if row.WidthFeet <= 0 || row.LengthFeet < row.WidthFeet || row.CapacityPounds <= 0 || row.FlySpeedFeet <= 0 {
				return fmt.Errorf("incomplete size row")
			}
		}
		if !itemRollTable(rows, c.TableDieSides) || len(c.Sizes) != c.SizeCount || c.ControlRangeFeet <= 0 || c.MaxLoadMultiplier < 2 || !uniqueItemKeys(c.SizeSelection) {
			return fmt.Errorf("size table disagrees with its die or size count")
		}
	}
	return nil
}
