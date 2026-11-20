package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_231.json
var items231JSON []byte

// ItemDetails231 holds exactly one typed payload per variant of catalogue batch
// #231 (magic items 07, SRD 5.2.1 pp. 220-223). Payloads are source facts only.
type ItemDetails231 struct {
	EnergyBow     *ItemEnergyBow          `json:"energy_bow,omitempty"`
	Eversmoking   *ItemEversmokingBottle  `json:"eversmoking_bottle,omitempty"`
	EyesCharming  *ItemEyesOfCharming     `json:"eyes_of_charming,omitempty"`
	EyesMinute    *ItemEyesOfMinuteSeeing `json:"eyes_of_minute_seeing,omitempty"`
	EyesEagle     *ItemEyesOfTheEagle     `json:"eyes_of_the_eagle,omitempty"`
	FeatherToken  *ItemFeatherToken       `json:"feather_token,omitempty"`
	Figurine      *ItemFigurine           `json:"figurine,omitempty"`
	FlameTongue   *ItemFlameTongue        `json:"flame_tongue,omitempty"`
	FoldingBoat   *ItemFoldingBoat        `json:"folding_boat,omitempty"`
	FrostBrand    *ItemFrostBrand         `json:"frost_brand,omitempty"`
	Gauntlets     *ItemGauntletsOgrePower `json:"gauntlets_of_ogre_power,omitempty"`
	GemBrightness *ItemGemOfBrightness    `json:"gem_of_brightness,omitempty"`
}

// ItemFact231 is one numeric source fact: a key, a non-negative amount and the
// unit it is printed in. The unit "minus" marks an amount printed with a minus
// sign (the Giant Fly's -4 modifiers) because payload numbers are never negative.
type ItemFact231 struct {
	Key    string `json:"key"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// ItemEnergyBow retains typed SRD source facts without implementing them.
type ItemEnergyBow struct {
	AttackBonus            int                   `json:"attack_bonus"`
	DamageBonus            int                   `json:"damage_bonus"`
	HasString              bool                  `json:"has_string"`
	ArrowAppears           string                `json:"arrow_appears"`
	ArrowDamageType        string                `json:"arrow_damage_type"`
	ReplacesDamageType     string                `json:"replaces_damage_type"`
	ArrowEnds              string                `json:"arrow_ends"`
	BrightLightFeet        int                   `json:"bright_light_feet"`
	DimLightAdditionalFeet int                   `json:"dim_light_additional_feet"`
	Restraint              ItemEnergyBowRestrain `json:"arrow_of_restraint"`
	Transport              ItemEnergyBowTransfer `json:"arrow_of_transport"`
	Ladder                 ItemEnergyBowLadder   `json:"energy_ladder"`
}

// ItemEnergyBowRestrain is the Arrow of Restraint property.
type ItemEnergyBowRestrain struct {
	Trigger         string    `json:"trigger"`
	SaveAbility     string    `json:"save_ability"`
	SaveDC          int       `json:"save_dc"`
	Condition       string    `json:"condition"`
	DurationSeconds int       `json:"duration_seconds"`
	EscapeAction    string    `json:"escape_action"`
	EscapeCheck     ItemCheck `json:"escape_check"`
	EscapeEnds      string    `json:"escape_ends"`
}

// ItemEnergyBowTransfer is the Arrow of Transport property.
type ItemEnergyBowTransfer struct {
	Action               string   `json:"action"`
	RangeFeet            int      `json:"range_feet"`
	Targets              []string `json:"targets"`
	ObjectMaxCubeFeet    int      `json:"object_max_cube_feet"`
	Destination          string   `json:"destination"`
	DestinationRangeFeet int      `json:"destination_range_feet"`
}

// ItemEnergyBowLadder is the Energy Ladder property.
type ItemEnergyBowLadder struct {
	Action           string `json:"action"`
	WallRangeFeet    int    `json:"wall_range_feet"`
	LadderLengthFeet int    `json:"ladder_length_feet"`
	DurationSeconds  int    `json:"duration_seconds"`
	EndsWith         string `json:"ends_with"`
}

// ItemEversmokingBottle retains typed SRD source facts without implementing
// them. The text starts on p. 220 and continues on p. 221.
type ItemEversmokingBottle struct {
	SourcePages         []int    `json:"source_pages"`
	Action              string   `json:"action"`
	Operations          []string `json:"operations"`
	EmanationFeet       int      `json:"emanation_feet"`
	Obscurement         string   `json:"obscurement"`
	GrowthFeetPerMinute int      `json:"growth_feet_per_minute"`
	MaxEmanationFeet    int      `json:"max_emanation_feet"`
	ClosedState         string   `json:"closed_state"`
	ClosedDisperseSecs  int      `json:"closed_disperses_after_seconds"`
	WindDisperseSeconds int      `json:"strong_wind_disperses_after_seconds"`
	WindExample         string   `json:"strong_wind_example"`
}

// ItemEyesOfCharming retains typed SRD source facts without implementing them.
type ItemEyesOfCharming struct {
	Appearance           string `json:"appearance"`
	Condition            string `json:"condition"`
	MaxCharges           int    `json:"max_charges"`
	Spell                string `json:"spell"`
	SaveDC               int    `json:"save_dc"`
	BaseSpellLevel       int    `json:"base_spell_level"`
	BaseCharges          int    `json:"base_charges"`
	LevelsPerExtraCharge int    `json:"levels_per_additional_charge"`
	Recharge             string `json:"recharge"`
}

// ItemEyesOfMinuteSeeing retains typed SRD source facts without implementing them.
type ItemEyesOfMinuteSeeing struct {
	Appearance     string `json:"appearance"`
	Condition      string `json:"condition"`
	RangeFeet      int    `json:"range_feet"`
	Darkvision     bool   `json:"darkvision_within_range"`
	CheckAbility   string `json:"check_ability"`
	CheckSkill     string `json:"check_skill"`
	CheckMode      string `json:"check_mode"`
	CheckSubject   string `json:"check_subject"`
	CheckRangeFeet int    `json:"check_range_feet"`
}

// ItemEyesOfTheEagle retains typed SRD source facts without implementing them.
type ItemEyesOfTheEagle struct {
	Appearance       string `json:"appearance"`
	Condition        string `json:"condition"`
	CheckAbility     string `json:"check_ability"`
	CheckSkill       string `json:"check_skill"`
	CheckMode        string `json:"check_mode"`
	CheckRequirement string `json:"check_requirement"`
	DetailCondition  string `json:"detail_condition"`
	DetailSizeFeet   int    `json:"detail_size_feet"`
	DetailSubjects   string `json:"detail_subjects"`
}

// ItemFeatherToken is one Feather Token kind: its row of the 1d100 table and the
// numeric facts of its single-use effect. Facts and lists are key lists, not prose.
type ItemFeatherToken struct {
	Token         string        `json:"token"`
	TableDieSides int           `json:"table_die_sides"`
	TableMin      int           `json:"table_min"`
	TableMax      int           `json:"table_max"`
	UseLimit      string        `json:"use_limit"`
	Action        string        `json:"action"`
	Placement     string        `json:"placement"`
	Effect        string        `json:"effect"`
	DismissAction string        `json:"dismiss_action"`
	Facts         []ItemFact231 `json:"facts"`
	Traits        []string      `json:"traits"`
	EndsOn        []string      `json:"ends_on" srd:"may-be-empty"`
}

// ItemFigurine holds the rules shared by every Figurine of Wondrous Power and
// the forms of one printed figurine (pp. 221-222). A duration of 0 on a form
// means the source states no time limit, only a charge pool.
type ItemFigurine struct {
	SourcePages          []int              `json:"source_pages"`
	Action               string             `json:"action"`
	ThrowRangeFeet       int                `json:"throw_range_feet"`
	ThrowTarget          string             `json:"throw_target"`
	NoSpaceResult        string             `json:"no_space_result"`
	FriendlyTo           []string           `json:"friendly_to"`
	UnderstandsLanguages bool               `json:"understands_owner_languages"`
	ObeysCommands        bool               `json:"obeys_owner_commands"`
	Initiative           string             `json:"initiative"`
	NoCommandBehavior    string             `json:"no_command_behavior"`
	RevertTriggers       []string           `json:"revert_triggers"`
	RevertAction         string             `json:"revert_action"`
	FigurinesInSet       int                `json:"figurines_in_set"`
	Forms                []ItemFigurineForm `json:"forms"`
}

// ItemFigurineForm is one creature a figurine becomes.
type ItemFigurineForm struct {
	Key           string        `json:"key"`
	Creature      string        `json:"creature"`
	DurationHours int           `json:"duration_hours"`
	CooldownDays  int           `json:"cooldown_days"`
	Facts         []ItemFact231 `json:"facts" srd:"may-be-empty"`
	Traits        []string      `json:"traits" srd:"may-be-empty"`
}

// ItemFlameTongue retains typed SRD source facts without implementing them.
type ItemFlameTongue struct {
	Condition              string   `json:"condition"`
	Action                 string   `json:"action"`
	Activation             string   `json:"activation"`
	Effect                 string   `json:"effect"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	ExtraDamage            ItemDice `json:"extra_damage"`
	ExtraDamageCondition   string   `json:"extra_damage_condition"`
	EndsOn                 []string `json:"ends_on"`
}

// ItemFoldingBoatCommand is one of the three command words.
type ItemFoldingBoatCommand struct {
	Ordinal  int      `json:"ordinal"`
	Result   string   `json:"result"`
	Requires []string `json:"requires" srd:"may-be-empty"`
}

// ItemFoldingBoat retains typed SRD source facts without implementing them.
type ItemFoldingBoat struct {
	LengthInches      int                      `json:"length_inches"`
	WidthInches       int                      `json:"width_inches"`
	DepthInches       int                      `json:"depth_inches"`
	WeightPounds      int                      `json:"weight_pounds"`
	Floats            bool                     `json:"floats"`
	StoresItems       bool                     `json:"stores_items_when_open"`
	CommandAction     string                   `json:"command_action"`
	Commands          []ItemFoldingBoatCommand `json:"command_words"`
	ObjectsThatFit    string                   `json:"objects_that_fit"`
	ObjectsThatNotFit string                   `json:"objects_that_do_not_fit"`
	VesselWeight      string                   `json:"vessel_weight"`
	StoredItemsInBoat string                   `json:"stored_items_when_vessel"`
	VesselStatistics  []string                 `json:"vessel_statistics_in"`
	DestroyedWhen     string                   `json:"destroyed_when"`
}

// ItemFrostBrand retains typed SRD source facts without implementing them.
type ItemFrostBrand struct {
	ExtraDamage            ItemDice `json:"extra_damage"`
	ExtraDamageTrigger     string   `json:"extra_damage_trigger"`
	ResistanceCondition    string   `json:"resistance_condition"`
	Resistance             string   `json:"resistance"`
	LightCondition         string   `json:"light_condition"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	ExtinguishTrigger      string   `json:"extinguish_trigger"`
	ExtinguishTargets      string   `json:"extinguish_targets"`
	ExtinguishRadiusFeet   int      `json:"extinguish_radius_feet"`
	ExtinguishUses         int      `json:"extinguish_uses"`
	ExtinguishCooldownHrs  int      `json:"extinguish_cooldown_hours"`
}

// ItemGauntletsOgrePower retains typed SRD source facts without implementing them.
type ItemGauntletsOgrePower struct {
	Condition         string `json:"condition"`
	StrengthScore     int    `json:"strength_score"`
	NoEffectAtOrAbove int    `json:"no_effect_at_or_above_score"`
}

// ItemGemOfBrightness retains typed SRD source facts without implementing them.
// The 50 charges are also recorded on the variant's charges field.
type ItemGemOfBrightness struct {
	Condition       string             `json:"condition"`
	Action          string             `json:"action"`
	CommandWords    int                `json:"command_words"`
	MaxCharges      int                `json:"max_charges"`
	First           ItemGemBrightLight `json:"first_command_word"`
	Second          ItemGemBrightBeam  `json:"second_command_word"`
	Third           ItemGemBrightFlare `json:"third_command_word"`
	WhenDepleted    string             `json:"when_depleted"`
	DepletedValueGP int                `json:"depleted_value_gp"`
}

// ItemGemBrightLight is the first command word.
type ItemGemBrightLight struct {
	ChargeCost             int      `json:"charge_cost"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	EndsOn                 []string `json:"ends_on"`
}

// ItemGemBrightBeam is the second command word.
type ItemGemBrightBeam struct {
	ChargeCost      int    `json:"charge_cost"`
	Target          string `json:"target"`
	RangeFeet       int    `json:"range_feet"`
	SaveAbility     string `json:"save_ability"`
	SaveDC          int    `json:"save_dc"`
	Condition       string `json:"condition"`
	DurationSeconds int    `json:"duration_seconds"`
	RepeatSave      string `json:"repeat_save"`
	SuccessEnds     string `json:"success_ends"`
}

// ItemGemBrightFlare is the third command word.
type ItemGemBrightFlare struct {
	ChargeCost int    `json:"charge_cost"`
	Shape      string `json:"shape"`
	SizeFeet   int    `json:"size_feet"`
	Targets    string `json:"targets"`
	Save       string `json:"save"`
}

// itemManifest231 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Feather Token has one variant per kind and Figurine
// of Wondrous Power one per printed figurine; table rows and the goats of the
// Ivory Goats are data inside a payload.
var itemManifest231 = itemSourceManifest{Batch: "data-magic-items-07", Issue: 231, Families: []itemSourceManifestFamily{
	{Slug: "energy-bow", Name: "Energy Bow", Page: 220, Category: "weapon", Options: "longbow shortbow", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement bonus.attack-and-damage-1 weapon.no-string weapon.energy-arrow-force-damage light.bright-20-dim-20 property.arrow-of-restraint save.strength-dc-15 check.strength-athletics-dc-20 property.arrow-of-transport action.magic property.energy-ladder", Payload: "energy_bow",
		Variants: []itemSourceManifestVariant{{"standard", "Energy Bow", "very-rare"}}},
	{Slug: "eversmoking-bottle", Name: "Eversmoking Bottle", Page: 220, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "open-or-close",
		Keys: "action.magic area.emanation-60-foot obscurement.heavily-obscured area.grows-10-feet-per-minute-to-120 effect.fixed-cloud-disperses-after-10-minutes effect.strong-wind-disperses-after-1-minute", Payload: "eversmoking_bottle",
		Variants: []itemSourceManifestVariant{{"standard", "Eversmoking Bottle", "uncommon"}}},
	{Slug: "eyes-of-charming", Name: "Eyes of Charming", Page: 221, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-one-or-more-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-all-at-dawn"},
		Keys:    "item.attunement item.wear item.charges spell.cast-charm-person spell.charm-person-save-dc-13 spell.level-increases-per-extra-charge", Payload: "eyes_of_charming",
		Variants: []itemSourceManifestVariant{{"standard", "Eyes of Charming", "uncommon"}}},
	{Slug: "eyes-of-minute-seeing", Name: "Eyes of Minute Seeing", Page: 221, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear vision.darkvision-1-foot check.investigation-advantage-within-1-foot", Payload: "eyes_of_minute_seeing",
		Variants: []itemSourceManifestVariant{{"standard", "Eyes of Minute Seeing", "uncommon"}}},
	{Slug: "eyes-of-the-eagle", Name: "Eyes of the Eagle", Page: 221, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear check.perception-sight-advantage vision.make-out-details-2-feet-across", Payload: "eyes_of_the_eagle",
		Variants: []itemSourceManifestVariant{{"standard", "Eyes of the Eagle", "uncommon"}}},
	{Slug: "feather-token", Name: "Feather Token", Page: 221, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "single-use",
		Keys: "action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table", Payload: "feather_token",
		Variants: []itemSourceManifestVariant{{"anchor", "Anchor", "uncommon"}, {"bird", "Bird", "rare"}, {"fan", "Fan", "uncommon"}, {"swan-boat", "Swan Boat", "rare"}, {"tree", "Tree", "uncommon"}, {"whip", "Whip", "rare"}}},
	{Slug: "figurine-of-wondrous-power", Name: "Figurine of Wondrous Power", Page: 221, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "throw-to-point-within-60-feet use-cooldown",
		Keys: "action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown", Payload: "figurine",
		Variants: []itemSourceManifestVariant{{"bronze-griffon", "Bronze Griffon", "rare"}, {"ebony-fly", "Ebony Fly", "rare"}, {"golden-lions", "Golden Lions", "rare"}, {"ivory-goats", "Ivory Goats", "rare"}, {"marble-elephant", "Marble Elephant", "rare"}, {"obsidian-steed", "Obsidian Steed", "very-rare"}, {"onyx-dog", "Onyx Dog", "rare"}, {"serpentine-owl", "Serpentine Owl", "rare"}, {"silver-raven", "Silver Raven", "uncommon"}}},
	{Slug: "flame-tongue", Name: "Flame Tongue", Page: 223, Category: "weapon", Options: "any-melee-weapon", Attunement: "required", Equip: "wield", Use: "bonus-action", UseRestrictions: "command-word",
		Keys: "item.attunement action.bonus item.command-word damage.fire-2d6 light.bright-40-dim-40", Payload: "flame_tongue",
		Variants: []itemSourceManifestVariant{{"standard", "Flame Tongue", "rare"}}},
	{Slug: "folding-boat", Name: "Folding Boat", Page: 223, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "command-word",
		Keys: "action.magic item.command-words-three item.floats item.container vessel.rowboat vessel.keelboat item.destroyed-with-vessel", Payload: "folding_boat",
		Variants: []itemSourceManifestVariant{{"standard", "Folding Boat", "rare"}}},
	{Slug: "frost-brand", Name: "Frost Brand", Page: 223, Category: "weapon", Options: "glaive greatsword longsword rapier scimitar shortsword", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "extinguish-once-per-hour",
		Keys: "item.attunement damage.cold-1d6 defense.resistance-fire light.bright-10-dim-10 effect.extinguish-nonmagical-flames-30-feet item.once-per-hour", Payload: "frost_brand",
		Variants: []itemSourceManifestVariant{{"standard", "Frost Brand", "very-rare"}}},
	{Slug: "gauntlets-of-ogre-power", Name: "Gauntlets of Ogre Power", Page: 223, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear ability.strength-score-19", Payload: "gauntlets_of_ogre_power",
		Variants: []itemSourceManifestVariant{{"standard", "Gauntlets of Ogre Power", "uncommon"}}},
	{Slug: "gem-of-brightness", Name: "Gem of Brightness", Page: 223, Category: "wondrous-item", Attunement: "none", Equip: "hold", Use: "magic-action", UseRestrictions: "command-word",
		Charges: &ItemCharges{Mode: "limited", Maximum: 50, PerUse: 1, Recharge: "none"},
		Keys:    "action.magic item.charges item.command-words-three light.bright-30-dim-30 save.constitution-dc-15 condition.blinded area.cone-30-foot item.nonmagical-jewel-when-expended", Payload: "gem_of_brightness",
		Variants: []itemSourceManifestVariant{{"standard", "Gem of Brightness", "uncommon"}}},
}}

// LoadItems231 returns a fresh validated copy of batch #231. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems231() (ItemSourceBatch[ItemDetails231], error) { return DecodeItems231(items231JSON) }

// DecodeItems231 validates supplied JSON as batch #231 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems231(data []byte) (ItemSourceBatch[ItemDetails231], error) {
	return decodeItemSourceBatch(data, itemManifest231, validateItemDetails231)
}

// validateItemFacts231 requires distinct, non-empty fact keys.
func validateItemFacts231(facts []ItemFact231) bool {
	keys := []string{}
	for _, f := range facts {
		keys = append(keys, f.Key)
	}
	return len(keys) == 0 || uniqueItemKeys(keys)
}

// validateItemDetails231 checks the internal consistency of each payload and a
// few source literals that the family's mechanic keys also name.
func validateItemDetails231(family, variant string, d *ItemDetails231) error {
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
	case "energy-bow":
		b := d.EnergyBow
		ok = positive(b.AttackBonus, b.DamageBonus, b.BrightLightFeet, b.DimLightAdditionalFeet, b.Restraint.SaveDC, b.Restraint.DurationSeconds, b.Restraint.EscapeCheck.DC,
			b.Transport.RangeFeet, b.Transport.ObjectMaxCubeFeet, b.Transport.DestinationRangeFeet, b.Ladder.WallRangeFeet, b.Ladder.LadderLengthFeet, b.Ladder.DurationSeconds) &&
			uniqueItemKeys(b.Transport.Targets) && b.Restraint.SaveDC == 15 && b.Restraint.EscapeCheck.DC == 20 // p. 220; named by the mechanic keys
	case "eversmoking-bottle":
		b := d.Eversmoking
		ok = slices.Equal(b.SourcePages, []int{220, 221}) && uniqueItemKeys(b.Operations) &&
			positive(b.EmanationFeet, b.GrowthFeetPerMinute, b.MaxEmanationFeet, b.ClosedDisperseSecs, b.WindDisperseSeconds) && b.MaxEmanationFeet > b.EmanationFeet
	case "eyes-of-charming":
		e := d.EyesCharming
		ok = positive(e.MaxCharges, e.SaveDC, e.BaseSpellLevel, e.BaseCharges, e.LevelsPerExtraCharge) && e.MaxCharges == 3 && e.SaveDC == 13
	case "eyes-of-minute-seeing":
		ok = positive(d.EyesMinute.RangeFeet, d.EyesMinute.CheckRangeFeet)
	case "eyes-of-the-eagle":
		ok = positive(d.EyesEagle.DetailSizeFeet)
	case "feather-token":
		t := d.FeatherToken
		// The six rows of the 1d100 table, in order; 00 is stored as 100.
		rows := map[string][2]int{"anchor": {1, 20}, "bird": {21, 35}, "fan": {36, 50}, "swan-boat": {51, 65}, "tree": {66, 90}, "whip": {91, 100}}
		ok = t.Token == variant && t.TableDieSides == 100 && [2]int{t.TableMin, t.TableMax} == rows[variant] && validateItemFacts231(t.Facts) &&
			len(t.Facts) > 0 && uniqueItemKeys(t.Traits) && (len(t.EndsOn) == 0 || uniqueItemKeys(t.EndsOn))
	case "figurine-of-wondrous-power":
		f := d.Figurine
		forms := map[string]int{"bronze-griffon": 1, "ebony-fly": 1, "golden-lions": 1, "ivory-goats": 3, "marble-elephant": 1, "obsidian-steed": 1, "onyx-dog": 1, "serpentine-owl": 1, "silver-raven": 1}
		ok = slices.Equal(f.SourcePages, []int{221, 222}) && positive(f.ThrowRangeFeet, f.FigurinesInSet) && uniqueItemKeys(f.FriendlyTo) && uniqueItemKeys(f.RevertTriggers) &&
			len(f.Forms) == forms[variant]
		for _, form := range f.Forms {
			ok = ok && form.CooldownDays > 0 && validateItemFacts231(form.Facts) && (len(form.Traits) == 0 || uniqueItemKeys(form.Traits)) &&
				(form.DurationHours > 0 || form.Key == "goat-of-traveling")
		}
	case "flame-tongue":
		f := d.FlameTongue
		ok = positive(f.BrightLightFeet, f.DimLightAdditionalFeet, f.ExtraDamage.Dice, f.ExtraDamage.Sides) && uniqueItemKeys(f.EndsOn)
	case "folding-boat":
		b := d.FoldingBoat
		ok = positive(b.LengthInches, b.WidthInches, b.DepthInches, b.WeightPounds) && uniqueItemKeys(b.VesselStatistics) && len(b.Commands) == 3
		for i, c := range b.Commands {
			ok = ok && c.Ordinal == i+1 && (len(c.Requires) == 0 || uniqueItemKeys(c.Requires))
		}
	case "frost-brand":
		f := d.FrostBrand
		ok = positive(f.ExtraDamage.Dice, f.ExtraDamage.Sides, f.BrightLightFeet, f.DimLightAdditionalFeet, f.ExtinguishRadiusFeet, f.ExtinguishUses, f.ExtinguishCooldownHrs)
	case "gauntlets-of-ogre-power":
		ok = positive(d.Gauntlets.StrengthScore, d.Gauntlets.NoEffectAtOrAbove) && d.Gauntlets.StrengthScore == 19
	case "gem-of-brightness":
		g := d.GemBrightness
		ok = g.CommandWords == 3 && g.MaxCharges == 50 && g.First.ChargeCost == 0 && g.Second.ChargeCost == 1 && g.Third.ChargeCost == 5 &&
			positive(g.First.BrightLightFeet, g.First.DimLightAdditionalFeet, g.Second.RangeFeet, g.Second.SaveDC, g.Second.DurationSeconds, g.Third.SizeFeet, g.DepletedValueGP) && uniqueItemKeys(g.First.EndsOn)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
