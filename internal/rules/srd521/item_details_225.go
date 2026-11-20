package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_225.json
var items225JSON []byte

// ItemDetails225 holds exactly one typed payload per variant of catalogue batch
// #225 (magic items 01, SRD 5.2.1 pp. 209–211). Payloads are source facts only.
type ItemDetails225 struct {
	Adamantine      *ItemAdamantine      `json:"adamantine,omitempty"`
	Slaying         *ItemSlaying         `json:"slaying,omitempty"`
	AmmunitionBonus *ItemAmmunitionBonus `json:"ammunition_bonus,omitempty"`
	Health          *ItemAmuletHealth    `json:"health,omitempty"`
	Proof           *ItemAmuletProof     `json:"proof,omitempty"`
	Planes          *ItemAmuletPlanes    `json:"planes,omitempty"`
	Animated        *ItemAnimatedShield  `json:"animated,omitempty"`
	Apparatus       *ItemApparatus       `json:"apparatus,omitempty"`
	Invulnerability *ItemInvulnerability `json:"invulnerability,omitempty"`
	Resistance      *ItemArmorResistance `json:"resistance,omitempty"`
	Vulnerability   *ItemVulnerability   `json:"vulnerability,omitempty"`
	ArmorBonus      *ItemArmorBonus      `json:"armor_bonus,omitempty"`
}

// ItemAdamantine retains typed SRD source facts without implementing them.
type ItemAdamantine struct {
	Condition                string `json:"condition"`
	CriticalHitAgainstWearer string `json:"critical_hit_against_wearer"`
}

// ItemSlayingType is one row of the Ammunition of Slaying 1d100 table.
type ItemSlayingType struct {
	Min          int    `json:"min"`
	Max          int    `json:"max"`
	CreatureType string `json:"creature_type"`
}

// ItemSlaying retains typed SRD source facts without implementing them.
type ItemSlaying struct {
	TypeSelection    []string          `json:"type_selection"`
	TableDieSides    int               `json:"table_die_sides"`
	CreatureTypes    []ItemSlayingType `json:"creature_types"`
	Trigger          string            `json:"trigger"`
	Save             string            `json:"save"`
	SaveDC           int               `json:"save_dc"`
	ExtraDamage      ItemDice          `json:"extra_damage"`
	SaveFailure      string            `json:"save_failure"`
	SaveSuccess      string            `json:"save_success"`
	AfterExtraDamage string            `json:"after_extra_damage"`
}

// ItemAmmunitionBonus retains typed SRD source facts without implementing them.
type ItemAmmunitionBonus struct {
	AttackBonus       int    `json:"attack_bonus"`
	DamageBonus       int    `json:"damage_bonus"`
	BonusDeterminedBy string `json:"bonus_determined_by"`
	AfterHit          string `json:"after_hit"`
	TypicalQuantities []int  `json:"typical_quantities"`
	ValuePieces       int    `json:"value_pieces"`
	ValueEquals       string `json:"value_equals"`
}

// ItemAmuletHealth retains typed SRD source facts without implementing them.
type ItemAmuletHealth struct {
	Condition              string `json:"condition"`
	Ability                string `json:"ability"`
	Score                  int    `json:"score"`
	NoEffectIfScoreAtLeast int    `json:"no_effect_if_score_at_least"`
}

// ItemAmuletProof retains typed SRD source facts without implementing them.
type ItemAmuletProof struct {
	Condition   string   `json:"condition"`
	Protections []string `json:"protections"`
	Exception   string   `json:"exception"`
}

// ItemPlaneRoll is one plane of a destination row's sub-roll.
type ItemPlaneRoll struct {
	Roll  int    `json:"roll"`
	Plane string `json:"plane"`
}

// ItemPlanesDestination is one row of the Amulet of the Planes 1d100 table.
// Rows without a sub-roll use zero sides and an explicit empty plane list.
type ItemPlanesDestination struct {
	Min         int             `json:"min"`
	Max         int             `json:"max"`
	Destination string          `json:"destination"`
	SubDieSides int             `json:"sub_die_sides"`
	Planes      []ItemPlaneRoll `json:"planes" srd:"may-be-empty"`
}

// ItemAmuletPlanes retains typed SRD source facts without implementing them.
type ItemAmuletPlanes struct {
	Condition         string                  `json:"condition"`
	Action            string                  `json:"action"`
	NamedLocation     string                  `json:"named_location"`
	Check             ItemCheck               `json:"check"`
	SuccessCasts      string                  `json:"success_casts"`
	FailureAffects    []string                `json:"failure_affects"`
	FailureRadiusFeet int                     `json:"failure_radius_feet"`
	TableDieSides     int                     `json:"table_die_sides"`
	Destinations      []ItemPlanesDestination `json:"destinations"`
}

// ItemAnimatedShield retains typed SRD source facts without implementing them.
type ItemAnimatedShield struct {
	Condition       string   `json:"condition"`
	ActivateAction  string   `json:"activate_action"`
	Effect          string   `json:"effect"`
	HandsFree       bool     `json:"hands_free"`
	DurationSeconds int      `json:"duration_seconds"`
	EndTriggers     []string `json:"end_triggers"`
	OnEnd           []string `json:"on_end"`
}

// ItemApparatusLever is one row of the Apparatus of the Crab Levers table.
type ItemApparatusLever struct {
	Lever int    `json:"lever"`
	Up    string `json:"up"`
	Down  string `json:"down"`
}

// ItemApparatus retains typed SRD source facts without implementing them.
type ItemApparatus struct {
	Appearance              string               `json:"appearance"`
	WeightPounds            int                  `json:"weight_pounds"`
	CatchCheck              ItemCheck            `json:"catch_check"`
	Occupants               int                  `json:"occupants"`
	OccupantMaxSize         string               `json:"occupant_max_size"`
	LeverCount              int                  `json:"lever_count"`
	LeverPositions          []string             `json:"lever_positions"`
	TransformedAppearance   string               `json:"transformed_appearance"`
	Size                    string               `json:"size"`
	AC                      int                  `json:"ac"`
	HP                      int                  `json:"hp"`
	SpeedFeet               int                  `json:"speed_feet"`
	SwimFeet                int                  `json:"swim_feet"`
	SpeedWithoutLegsFeet    int                  `json:"speed_without_legs_feet"`
	Immunities              []string             `json:"immunities"`
	Pilots                  int                  `json:"pilots"`
	ClosedHatch             []string             `json:"closed_hatch"`
	AirHours                int                  `json:"air_hours"`
	AirDividedBy            string               `json:"air_divided_by"`
	Floats                  bool                 `json:"floats"`
	MaxDepthFeet            int                  `json:"max_depth_feet"`
	PressureDamage          ItemDice             `json:"pressure_damage"`
	PressureIntervalSeconds int                  `json:"pressure_interval_seconds"`
	OperateAction           string               `json:"operate_action"`
	Operator                string               `json:"operator"`
	MaxLeversPerAction      int                  `json:"max_levers_per_action"`
	LeverReset              string               `json:"lever_reset"`
	LeverOrder              string               `json:"lever_order"`
	Levers                  []ItemApparatusLever `json:"levers"`
	Claws                   int                  `json:"claws"`
	ClawAttackBonus         int                  `json:"claw_attack_bonus"`
	ClawReachFeet           int                  `json:"claw_reach_feet"`
	ClawUpAverage           int                  `json:"claw_up_average"`
	ClawUpDamage            ItemDice             `json:"claw_up_damage"`
	ClawDownCondition       string               `json:"claw_down_condition"`
	ClawEscapeDC            int                  `json:"claw_escape_dc"`
	BrightLightFeet         int                  `json:"bright_light_feet"`
	DimLightAdditionalFeet  int                  `json:"dim_light_additional_feet"`
	VerticalFeet            int                  `json:"vertical_feet"`
	TurnDegrees             int                  `json:"turn_degrees"`
	SideShuttersPerSide     int                  `json:"side_shutters_per_side"`
}

// ItemInvulnerability retains typed SRD source facts without implementing them.
// The once-per-dawn Metal Shell property is a use limit, not a charge pool.
type ItemInvulnerability struct {
	Condition               string   `json:"condition"`
	Resistances             []string `json:"resistances"`
	Property                string   `json:"property"`
	PropertyAction          string   `json:"property_action"`
	PropertyImmunities      []string `json:"property_immunities"`
	PropertyDurationSeconds int      `json:"property_duration_seconds"`
	PropertyEndsEarly       string   `json:"property_ends_early"`
	PropertyUses            int      `json:"property_uses"`
	PropertyRecharge        string   `json:"property_recharge"`
}

// ItemResistanceType is one row of the Armor of Resistance 1d10 table.
type ItemResistanceType struct {
	Roll       int    `json:"roll"`
	DamageType string `json:"damage_type"`
}

// ItemArmorResistance retains typed SRD source facts without implementing them.
type ItemArmorResistance struct {
	Condition       string               `json:"condition"`
	ResistanceCount int                  `json:"resistance_count"`
	TypeSelection   []string             `json:"type_selection"`
	TableDieSides   int                  `json:"table_die_sides"`
	DamageTypes     []ItemResistanceType `json:"damage_types"`
}

// ItemVulnerability retains typed SRD source facts without implementing them.
type ItemVulnerability struct {
	Condition             string   `json:"condition"`
	ResistanceCount       int      `json:"resistance_count"`
	DamageTypes           []string `json:"damage_types"`
	TypeSelection         []string `json:"type_selection"`
	Cursed                bool     `json:"cursed"`
	CurseRevealedBy       []string `json:"curse_revealed_by"`
	CurseTrigger          string   `json:"curse_trigger"`
	CurseEndsBy           []string `json:"curse_ends_by"`
	RemovingArmorEndCurse bool     `json:"removing_armor_ends_curse"`
	VulnerabilityCount    int      `json:"vulnerability_count"`
	VulnerabilityTypes    string   `json:"vulnerability_types"`
}

// ItemArmorBonus retains typed SRD source facts without implementing them.
type ItemArmorBonus struct {
	Condition         string `json:"condition"`
	ACBonus           int    `json:"ac_bonus"`
	BonusDeterminedBy string `json:"bonus_determined_by"`
}

const itemAnyArmor225 = "any-light-armor any-medium-armor any-heavy-armor"

// itemManifest225 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Rarity variants are explicit; table rows are data.
var itemManifest225 = itemSourceManifest{Batch: "data-magic-items-01", Issue: 225, Families: []itemSourceManifestFamily{
	{Slug: "adamantine-armor", Name: "Adamantine Armor", Page: 209, Category: "armor", Options: "any-medium-armor any-heavy-armor", Excluded: "hide-armor", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear-armor defense.critical-hit-becomes-normal-hit", Payload: "adamantine",
		Variants: []itemSourceManifestVariant{{"standard", "Adamantine Armor", "uncommon"}}},
	{Slug: "ammunition-of-slaying", Name: "Ammunition of Slaying", Page: 209, Category: "weapon", Options: "any-ammunition", Attunement: "none", Equip: "none", Use: "ammunition-attack", UseRestrictions: "creature-type-determined-by-gm-or-table",
		Keys: "item.magic-ammunition item.creature-type-table save.constitution-dc-17 damage.extra-force-half-on-save item.nonmagical-after-extra-damage", Payload: "slaying",
		Variants: []itemSourceManifestVariant{{"standard", "Ammunition of Slaying", "very-rare"}}},
	{Slug: "ammunition-1-2-or-3", Name: "Ammunition, +1, +2, or +3", Page: 209, Category: "weapon", Options: "any-ammunition", Attunement: "none", Equip: "none", Use: "ammunition-attack",
		Keys: "item.magic-ammunition attack.ammunition-bonus damage.ammunition-bonus item.no-longer-magical-after-hit", Payload: "ammunition_bonus",
		Variants: []itemSourceManifestVariant{{"plus-1", "Ammunition, +1", "uncommon"}, {"plus-2", "Ammunition, +2", "rare"}, {"plus-3", "Ammunition, +3", "very-rare"}}},
	{Slug: "amulet-of-health", Name: "Amulet of Health", Page: 209, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear ability.constitution-score-19", Payload: "health",
		Variants: []itemSourceManifestVariant{{"standard", "Amulet of Health", "rare"}}},
	{Slug: "amulet-of-proof-against-detection-and-location", Name: "Amulet of Proof against Detection and Location", Page: 209, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear protection.untargetable-by-divination-spells protection.imperceptible-to-scrying-sensors protection.wearer-may-allow", Payload: "proof",
		Variants: []itemSourceManifestVariant{{"standard", "Amulet of Proof against Detection and Location", "uncommon"}}},
	{Slug: "amulet-of-the-planes", Name: "Amulet of the Planes", Page: 209, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "name-familiar-location-on-another-plane",
		Keys: "item.attunement item.wear action.magic check.intelligence-arcana-dc-15 spell.cast-plane-shift transport.random-planar-destination-table area.creatures-and-objects-within-15-feet", Payload: "planes",
		Variants: []itemSourceManifestVariant{{"standard", "Amulet of the Planes", "very-rare"}}},
	{Slug: "animated-shield", Name: "Animated Shield", Page: 209, Category: "armor", Options: "shield", Attunement: "required", Equip: "hold", Use: "bonus-action",
		Keys: "item.attunement item.hold-shield action.bonus-action defense.shield-without-wielding duration.one-minute duration.ends-on-bonus-action-death-or-incapacitated", Payload: "animated",
		Variants: []itemSourceManifestVariant{{"standard", "Animated Shield", "very-rare"}}},
	{Slug: "apparatus-of-the-crab", Name: "Apparatus of the Crab", Page: 210, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "creature-in-compartment requires-one-pilot up-to-two-levers-per-action",
		Keys: "object.statistics vehicle.pilot vehicle.air-supply vehicle.depth-pressure-damage check.intelligence-investigation-dc-20 action.utilize item.lever-table attack.claw-melee condition.grappled-escape-dc-15 light.bright-and-dim", Payload: "apparatus",
		Variants: []itemSourceManifestVariant{{"standard", "Apparatus of the Crab", "legendary"}}},
	{Slug: "armor-of-invulnerability", Name: "Armor of Invulnerability", Page: 210, Category: "armor", Options: "plate-armor", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.attunement item.wear-armor defense.resistance-bludgeoning-piercing-slashing action.magic defense.immunity-bludgeoning-piercing-slashing duration.ten-minutes-or-armor-removed item.once-until-next-dawn", Payload: "invulnerability",
		Variants: []itemSourceManifestVariant{{"standard", "Armor of Invulnerability", "legendary"}}},
	{Slug: "armor-of-resistance", Name: "Armor of Resistance", Page: 210, Category: "armor", Options: itemAnyArmor225, Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear-armor item.damage-type-table defense.resistance-chosen-damage-type", Payload: "resistance",
		Variants: []itemSourceManifestVariant{{"standard", "Armor of Resistance", "rare"}}},
	{Slug: "armor-of-vulnerability", Name: "Armor of Vulnerability", Page: 211, Category: "armor", Options: itemAnyArmor225, Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear-armor defense.resistance-chosen-damage-type item.curse item.curse-revealed-by-identify-or-attunement item.curse-ends-on-remove-curse defense.vulnerability-other-two-damage-types", Payload: "vulnerability",
		Variants: []itemSourceManifestVariant{{"standard", "Armor of Vulnerability", "rare"}}},
	{Slug: "armor-1-2-or-3", Name: "Armor, +1, +2, or +3", Page: 210, Category: "armor", Options: itemAnyArmor225, Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear-armor defense.armor-class-bonus", Payload: "armor_bonus",
		Variants: []itemSourceManifestVariant{{"plus-1", "Armor, +1", "rare"}, {"plus-2", "Armor, +2", "very-rare"}, {"plus-3", "Armor, +3", "legendary"}}},
}}

// LoadItems225 returns a fresh validated copy of batch #225. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems225() (ItemSourceBatch[ItemDetails225], error) { return DecodeItems225(items225JSON) }

// DecodeItems225 validates supplied JSON as batch #225 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems225(data []byte) (ItemSourceBatch[ItemDetails225], error) {
	return decodeItemSourceBatch(data, itemManifest225, validateItemDetails225)
}

// validateItemDetails225 checks the internal consistency of each payload: dice,
// DCs and tables must be complete, and a rarity variant must carry its own bonus.
func validateItemDetails225(family, variant string, d *ItemDetails225) error {
	bonus := map[string]int{"plus-1": 1, "plus-2": 2, "plus-3": 3}[variant]
	dice := func(x ItemDice) bool { return x.Dice > 0 && x.Sides > 0 }
	switch family {
	case "ammunition-of-slaying":
		s := d.Slaying
		rows, types := [][2]int{}, []string{}
		for _, row := range s.CreatureTypes {
			rows, types = append(rows, [2]int{row.Min, row.Max}), append(types, row.CreatureType)
		}
		if !itemRollTable(rows, s.TableDieSides) || !uniqueItemKeys(types) || !uniqueItemKeys(s.TypeSelection) || s.SaveDC <= 0 || !dice(s.ExtraDamage) {
			return fmt.Errorf("incomplete creature type table, save or damage")
		}
	case "ammunition-1-2-or-3":
		a := d.AmmunitionBonus
		if bonus == 0 || a.AttackBonus != bonus || a.DamageBonus != bonus || a.ValuePieces <= 0 || slices.Contains(a.TypicalQuantities, 0) {
			return fmt.Errorf("bonus disagrees with its rarity variant")
		}
	case "amulet-of-health":
		if h := d.Health; h.Score < 1 || h.Score > 30 || h.NoEffectIfScoreAtLeast != h.Score {
			return fmt.Errorf("invalid ability score")
		}
	case "amulet-of-proof-against-detection-and-location":
		if !uniqueItemKeys(d.Proof.Protections) {
			return fmt.Errorf("duplicate protection")
		}
	case "amulet-of-the-planes":
		p := d.Planes
		rows := [][2]int{}
		for _, row := range p.Destinations {
			rows = append(rows, [2]int{row.Min, row.Max})
			planes := []string{}
			for i, plane := range row.Planes {
				if plane.Roll != i+1 {
					return fmt.Errorf("incomplete plane sub-roll")
				}
				planes = append(planes, plane.Plane)
			}
			if len(planes) != row.SubDieSides || len(planes) > 0 && !uniqueItemKeys(planes) {
				return fmt.Errorf("plane list disagrees with its sub-roll die")
			}
		}
		if !itemRollTable(rows, p.TableDieSides) || p.Check.DC <= 0 || p.FailureRadiusFeet <= 0 || !uniqueItemKeys(p.FailureAffects) {
			return fmt.Errorf("incomplete destination table, check or area")
		}
	case "animated-shield":
		if a := d.Animated; a.DurationSeconds <= 0 || !uniqueItemKeys(a.EndTriggers) || !uniqueItemKeys(a.OnEnd) {
			return fmt.Errorf("invalid duration or end triggers")
		}
	case "apparatus-of-the-crab":
		a := d.Apparatus
		for i, lever := range a.Levers {
			if lever.Lever != i+1 || lever.Up == lever.Down {
				return fmt.Errorf("incomplete lever table")
			}
		}
		if len(a.Levers) != a.LeverCount || a.LeverCount == 0 || a.MaxLeversPerAction <= 0 || a.MaxLeversPerAction > a.LeverCount ||
			a.CatchCheck.DC <= 0 || a.ClawEscapeDC <= 0 || a.AC <= 0 || a.HP <= 0 || a.Pilots <= 0 || a.Occupants <= 0 ||
			!dice(a.PressureDamage) || !dice(a.ClawUpDamage) || a.PressureIntervalSeconds <= 0 || !uniqueItemKeys(a.Immunities) || !uniqueItemKeys(a.LeverPositions) {
			return fmt.Errorf("incomplete statistics, levers or attacks")
		}
	case "armor-of-invulnerability":
		if i := d.Invulnerability; !uniqueItemKeys(i.Resistances) || !uniqueItemKeys(i.PropertyImmunities) || i.PropertyDurationSeconds <= 0 || i.PropertyUses <= 0 {
			return fmt.Errorf("invalid damage types, duration or uses")
		}
	case "armor-of-resistance":
		r := d.Resistance
		types := []string{}
		for i, row := range r.DamageTypes {
			if row.Roll != i+1 {
				return fmt.Errorf("incomplete damage type table")
			}
			types = append(types, row.DamageType)
		}
		if len(types) != r.TableDieSides || !uniqueItemKeys(types) || r.ResistanceCount != 1 || !uniqueItemKeys(r.TypeSelection) {
			return fmt.Errorf("damage type table disagrees with its die")
		}
	case "armor-of-vulnerability":
		v := d.Vulnerability
		if !uniqueItemKeys(v.DamageTypes) || v.ResistanceCount <= 0 || v.VulnerabilityCount <= 0 || v.ResistanceCount+v.VulnerabilityCount != len(v.DamageTypes) ||
			!v.Cursed || !uniqueItemKeys(v.CurseRevealedBy) || !uniqueItemKeys(v.CurseEndsBy) || !uniqueItemKeys(v.TypeSelection) {
			return fmt.Errorf("invalid damage types or curse")
		}
	case "armor-1-2-or-3":
		if bonus == 0 || d.ArmorBonus.ACBonus != bonus {
			return fmt.Errorf("bonus disagrees with its rarity variant")
		}
	}
	return nil
}
