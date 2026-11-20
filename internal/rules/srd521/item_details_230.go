package srd521

import (
	_ "embed"
	"fmt"
)

//go:embed testdata/items_230.json
var items230JSON []byte

// ItemDetails230 holds exactly one typed payload per variant of catalogue batch
// #230 (magic items 06, SRD 5.2.1 pp. 218–220). Payloads are source facts only.
type ItemDetails230 struct {
	ScaleMail     *ItemDragonScaleMail     `json:"scale_mail,omitempty"`
	Slayer        *ItemDragonSlayer        `json:"slayer,omitempty"`
	Disappearance *ItemDustDisappearance   `json:"disappearance,omitempty"`
	Dryness       *ItemDustDryness         `json:"dryness,omitempty"`
	Sneezing      *ItemDustSneezingChoking `json:"sneezing,omitempty"`
	Plate         *ItemDwarvenPlate        `json:"plate,omitempty"`
	Thrower       *ItemDwarvenThrower      `json:"thrower,omitempty"`
	Quiver        *ItemEfficientQuiver     `json:"quiver,omitempty"`
	Bottle        *ItemEfreetiBottle       `json:"bottle,omitempty"`
	Gem           *ItemElementalGem        `json:"gem,omitempty"`
	Elixir        *ItemElixirOfHealth      `json:"elixir,omitempty"`
	Chain         *ItemElvenChain          `json:"chain,omitempty"`
}

// ItemDragonScaleMail retains typed SRD source facts without implementing them.
// One variant per row of the source's dragon/resistance table.
type ItemDragonScaleMail struct {
	Condition       string   `json:"condition"`
	ACBonus         int      `json:"ac_bonus"`
	BreathSaveMode  string   `json:"breath_weapon_save_mode"`
	BreathSource    string   `json:"breath_weapon_source"`
	Dragon          string   `json:"dragon"`
	Resistance      string   `json:"resistance"`
	SenseAction     string   `json:"sense_action"`
	SenseTarget     string   `json:"sense_target"`
	SenseRangeMiles int      `json:"sense_range_miles"`
	SenseReveals    []string `json:"sense_reveals"`
	SenseUses       int      `json:"sense_uses"`
	SenseRecharge   string   `json:"sense_recharge"`
}

// ItemDragonSlayer retains typed SRD source facts without implementing them.
type ItemDragonSlayer struct {
	AttackBonus      int      `json:"attack_bonus"`
	DamageBonus      int      `json:"damage_bonus"`
	ExtraDamage      ItemDice `json:"extra_damage"`
	ExtraDamageIfHit string   `json:"extra_damage_target"`
}

// ItemDustDisappearance retains typed SRD source facts without implementing them.
type ItemDustDisappearance struct {
	Appearance         string   `json:"appearance"`
	Uses               int      `json:"uses"`
	Action             string   `json:"action"`
	AreaEmanationFeet  int      `json:"area_emanation_feet"`
	Affects            []string `json:"affects"`
	GrantsCondition    string   `json:"grants_condition"`
	Duration           ItemDice `json:"duration"`
	DurationSameForAll bool     `json:"duration_same_for_all"`
	ConsumedWhen       string   `json:"consumed_when"`
	EndsForCreatureOn  []string `json:"ends_for_creature_after"`
}

// ItemDustDryness retains typed SRD source facts without implementing them.
type ItemDustDryness struct {
	Container       string   `json:"container"`
	Pinches         ItemDice `json:"pinches"`
	PinchesBonus    int      `json:"pinches_bonus"`
	WaterAction     string   `json:"water_action"`
	WaterCubeFeet   int      `json:"water_cube_feet"`
	WaterResult     string   `json:"water_result"`
	PelletLocation  string   `json:"pellet_location"`
	PelletWeight    string   `json:"pellet_weight"`
	SmashAction     string   `json:"smash_action"`
	SmashSurface    string   `json:"smash_surface"`
	SmashResult     string   `json:"smash_result"`
	SmashEnds       string   `json:"smash_ends"`
	ElementalAction string   `json:"elemental_action"`
	ElementalRange  int      `json:"elemental_range_feet"`
	ElementalTarget string   `json:"elemental_target"`
	SaveAbility     string   `json:"save_ability"`
	SaveDC          int      `json:"save_dc"`
	Damage          ItemDice `json:"damage"`
	SaveSuccess     string   `json:"save_success"`
}

// ItemDustSneezingChoking retains typed SRD source facts without implementing them.
type ItemDustSneezingChoking struct {
	Container         string   `json:"container"`
	Resembles         string   `json:"resembles"`
	IdentifyReveals   bool     `json:"identify_reveals"`
	Uses              int      `json:"uses"`
	Action            string   `json:"action"`
	AreaEmanationFeet int      `json:"area_emanation_feet"`
	AffectsUser       bool     `json:"affects_user"`
	SaveAbility       string   `json:"save_ability"`
	SaveDC            int      `json:"save_dc"`
	AutoSucceed       []string `json:"auto_succeed_creature_types"`
	FailureEffect     string   `json:"failure_effect"`
	FailureConditions []string `json:"failure_conditions"`
	RepeatSave        string   `json:"repeat_save"`
	EndsOn            []string `json:"ends_on"`
}

// ItemDwarvenPlate retains typed SRD source facts without implementing them.
type ItemDwarvenPlate struct {
	Condition         string `json:"condition"`
	ACBonus           int    `json:"ac_bonus"`
	ReactionTrigger   string `json:"reaction_trigger"`
	ReactionEffect    string `json:"reaction_effect"`
	ReduceDistanceMax int    `json:"reduce_distance_up_to_feet"`
}

// ItemDwarvenThrower retains typed SRD source facts without implementing them.
type ItemDwarvenThrower struct {
	AttunementBy      []string `json:"attunement_by_any_of"`
	AttackBonus       int      `json:"attack_bonus"`
	DamageBonus       int      `json:"damage_bonus"`
	Thrown            bool     `json:"thrown_property"`
	NormalRangeFeet   int      `json:"normal_range_feet"`
	LongRangeFeet     int      `json:"long_range_feet"`
	RangedExtraDamage ItemDice `json:"ranged_hit_extra_damage"`
	GiantExtraDamage  ItemDice `json:"giant_target_extra_damage"`
	GiantTarget       string   `json:"giant_target"`
	ReturnsTo         string   `json:"returns_to"`
	ReturnsWhen       string   `json:"returns_when"`
}

// ItemQuiverCompartment is one compartment of the Efficient Quiver, in source order.
type ItemQuiverCompartment struct {
	Length   string   `json:"length"`
	Capacity int      `json:"capacity"`
	Holds    []string `json:"holds"`
}

// ItemEfficientQuiver retains typed SRD source facts without implementing them.
type ItemEfficientQuiver struct {
	Compartments     []ItemQuiverCompartment `json:"compartments"`
	Extradimensional bool                    `json:"extradimensional"`
	MaxWeightPounds  int                     `json:"max_weight_pounds"`
	DrawAs           string                  `json:"draw_as"`
}

// ItemEfreetiRow is one row of the Efreeti Bottle's 1d10 table.
type ItemEfreetiRow struct {
	Min                int      `json:"min"`
	Max                int      `json:"max"`
	Effect             string   `json:"effect"`
	UnderstandsOwner   bool     `json:"understands_owner_languages"`
	ObeysOwner         bool     `json:"obeys_owner_commands"`
	FightRounds        int      `json:"fight_rounds"`
	DurationSeconds    int      `json:"duration_seconds"`
	StopperLockedHours int      `json:"stopper_locked_hours"`
	SameEffectOpenings int      `json:"same_effect_next_openings"`
	EscapesOnOpening   int      `json:"escapes_on_opening"`
	Wishes             int      `json:"wishes"`
	Ends               []string `json:"ends"`
	BottleAfter        string   `json:"bottle_after"`
}

// ItemEfreetiBottle retains typed SRD source facts without implementing them.
type ItemEfreetiBottle struct {
	Appearance     string           `json:"appearance"`
	Action         string           `json:"action"`
	SmokeEnds      string           `json:"smoke_ends"`
	SmokeEndEffect string           `json:"smoke_end_effect"`
	Summons        string           `json:"summons"`
	AppearsRange   int              `json:"appears_within_feet"`
	TableAppliesTo string           `json:"table_applies_to"`
	TableDieSides  int              `json:"table_die_sides"`
	Effects        []ItemEfreetiRow `json:"effects"`
}

// ItemElementalGem retains typed SRD source facts without implementing them.
// One variant per row of the source's gem/elemental table.
type ItemElementalGem struct {
	Gem              string   `json:"gem"`
	Elemental        string   `json:"elemental"`
	Action           string   `json:"action"`
	GemAfter         string   `json:"gem_after"`
	Location         string   `json:"location"`
	UnderstandsOwner bool     `json:"understands_owner_languages"`
	ObeysOwner       bool     `json:"obeys_owner_commands"`
	Initiative       string   `json:"initiative"`
	DurationSeconds  int      `json:"duration_seconds"`
	EndTriggers      []string `json:"end_triggers"`
}

// ItemElixirOfHealth retains typed SRD source facts without implementing them.
type ItemElixirOfHealth struct {
	CuresMagicalContagions bool     `json:"cures_magical_contagions"`
	EndsConditions         []string `json:"ends_conditions"`
	Appearance             string   `json:"appearance"`
}

// ItemElvenChain retains typed SRD source facts without implementing them.
type ItemElvenChain struct {
	Condition          string   `json:"condition"`
	ACBonus            int      `json:"ac_bonus"`
	TrainedWithThis    bool     `json:"trained_with_this_armor"`
	TrainedEvenWithout []string `json:"trained_even_without"`
}

// itemManifest230 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Dragon Scale Mail has ten variants (the source's
// dragon table) and Elemental Gem four (its gem table); the rest have one.
var itemManifest230 = itemSourceManifest{Batch: "data-magic-items-06", Issue: 230, Families: []itemSourceManifestFamily{
	{Slug: "dragon-scale-mail", Name: "Dragon Scale Mail", Page: 218, Category: "armor", Options: "scale-mail", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.attunement item.wear defense.armor-class-bonus-1 save.advantage-against-dragon-breath defense.resistance-by-dragon-type action.magic sense.closest-same-type-dragon-30-miles item.once-until-next-dawn", Payload: "scale_mail",
		Variants: []itemSourceManifestVariant{
			{"black", "Dragon Scale Mail (black)", "very-rare"}, {"blue", "Dragon Scale Mail (blue)", "very-rare"}, {"brass", "Dragon Scale Mail (brass)", "very-rare"},
			{"bronze", "Dragon Scale Mail (bronze)", "very-rare"}, {"copper", "Dragon Scale Mail (copper)", "very-rare"}, {"gold", "Dragon Scale Mail (gold)", "very-rare"},
			{"green", "Dragon Scale Mail (green)", "very-rare"}, {"red", "Dragon Scale Mail (red)", "very-rare"}, {"silver", "Dragon Scale Mail (silver)", "very-rare"},
			{"white", "Dragon Scale Mail (white)", "very-rare"}}},
	{Slug: "dragon-slayer", Name: "Dragon Slayer", Page: 219, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "none", Equip: "wield", Use: "weapon-attack",
		Keys: "item.magic-weapon attack.weapon-bonus damage.weapon-bonus damage.extra-3d6-against-dragons", Payload: "slayer",
		Variants: []itemSourceManifestVariant{{"standard", "Dragon Slayer", "rare"}}},
	{Slug: "dust-of-disappearance", Name: "Dust of Disappearance", Page: 219, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "throw-into-air one-use",
		Keys: "action.utilize area.emanation-10-feet condition.invisible duration.2d4-minutes effect.ends-after-attack-damage-or-spell item.consumed-on-use", Payload: "disappearance",
		Variants: []itemSourceManifestVariant{{"standard", "Dust of Disappearance", "uncommon"}}},
	{Slug: "dust-of-dryness", Name: "Dust of Dryness", Page: 219, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "sprinkle-one-pinch",
		Keys: "item.pinches-1d6-plus-4 action.utilize effect.water-cube-15-feet-to-pellet action.utilize-smash-pellet effect.pellet-releases-water save.constitution-dc-13 damage.necrotic-10d6 damage.half-on-success", Payload: "dryness",
		Variants: []itemSourceManifestVariant{{"standard", "Dust of Dryness", "uncommon"}}},
	{Slug: "dust-of-sneezing-and-choking", Name: "Dust of Sneezing and Choking", Page: 219, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "throw-into-air one-use",
		Keys: "action.utilize area.emanation-30-feet save.constitution-dc-15 save.automatic-success-by-creature-type condition.incapacitated condition.suffocating save.repeat-end-of-turn effect.ends-with-lesser-restoration", Payload: "sneezing",
		Variants: []itemSourceManifestVariant{{"standard", "Dust of Sneezing and Choking", "uncommon"}}},
	{Slug: "dwarven-plate", Name: "Dwarven Plate", Page: 219, Category: "armor", Options: "half-plate-armor plate-armor", Attunement: "none", Equip: "wear", Use: "reaction", UseRestrictions: "moved-against-will-along-ground",
		Keys: "item.wear defense.armor-class-bonus-2 action.reaction movement.reduce-forced-ground-movement-10-feet", Payload: "plate",
		Variants: []itemSourceManifestVariant{{"standard", "Dwarven Plate", "very-rare"}}},
	{Slug: "dwarven-thrower", Name: "Dwarven Thrower", Page: 219, Category: "weapon", Options: "warhammer", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon attack.weapon-bonus damage.weapon-bonus property.thrown-20-60 damage.extra-force-1d8-ranged-hit damage.extra-force-2d8-against-giants item.returns-to-hand", Payload: "thrower",
		Variants: []itemSourceManifestVariant{{"standard", "Dwarven Thrower", "very-rare"}}},
	{Slug: "efficient-quiver", Name: "Efficient Quiver", Page: 219, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "container",
		Keys: "item.extradimensional-storage item.weight-cap-2-pounds item.compartments-60-18-6 action.draw-as-regular-quiver", Payload: "quiver",
		Variants: []itemSourceManifestVariant{{"standard", "Efficient Quiver", "uncommon"}}},
	{Slug: "efreeti-bottle", Name: "Efreeti Bottle", Page: 220, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "remove-stopper",
		Keys: "action.magic summon.efreeti table.efreeti-1d10-first-opening effect.bottle-loses-magic", Payload: "bottle",
		Variants: []itemSourceManifestVariant{{"standard", "Efreeti Bottle", "very-rare"}}},
	{Slug: "elemental-gem", Name: "Elemental Gem", Page: 220, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "break-gem gem-ceases-to-be-magical",
		Keys: "action.utilize summon.elemental-by-gem-type duration.one-hour effect.gem-ceases-to-be-magical", Payload: "gem",
		Variants: []itemSourceManifestVariant{
			{"blue-sapphire", "Elemental Gem (blue sapphire)", "uncommon"}, {"emerald", "Elemental Gem (emerald)", "uncommon"},
			{"red-corundum", "Elemental Gem (red corundum)", "uncommon"}, {"yellow-diamond", "Elemental Gem (yellow diamond)", "uncommon"}}},
	{Slug: "elixir-of-health", Name: "Elixir of Health", Page: 220, Category: "potion", Attunement: "none", Equip: "none", Use: "consume", UseRestrictions: "drink-one",
		Keys: "item.consume cure.magical-contagions condition.end-blinded condition.end-deafened condition.end-paralyzed condition.end-poisoned", Payload: "elixir",
		Variants: []itemSourceManifestVariant{{"standard", "Elixir of Health", "rare"}}},
	{Slug: "elven-chain", Name: "Elven Chain", Page: 220, Category: "armor", Options: "chain-mail chain-shirt", Attunement: "none", Equip: "wear", Use: "none",
		Keys: "item.wear defense.armor-class-bonus-1 training.trained-with-this-armor", Payload: "chain",
		Variants: []itemSourceManifestVariant{{"standard", "Elven Chain", "rare"}}},
}}

// LoadItems230 returns a fresh validated copy of batch #230. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems230() (ItemSourceBatch[ItemDetails230], error) { return DecodeItems230(items230JSON) }

// DecodeItems230 validates supplied JSON as batch #230 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems230(data []byte) (ItemSourceBatch[ItemDetails230], error) {
	return decodeItemSourceBatch(data, itemManifest230, validateItemDetails230)
}

// validateItemDetails230 checks each payload against the source's own tables
// and internal consistency: dice, DCs, ranges and key lists.
func validateItemDetails230(family, variant string, d *ItemDetails230) error {
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
	case "dragon-scale-mail":
		// The dragon table, p. 219: each dragon's resistance.
		resistance := map[string]string{"black": "acid", "blue": "lightning", "brass": "fire", "bronze": "lightning", "copper": "acid",
			"gold": "fire", "green": "poison", "red": "fire", "silver": "cold", "white": "cold"}
		c := d.ScaleMail
		ok = c.Dragon == variant && c.Resistance == resistance[variant] && positive(c.ACBonus, c.SenseRangeMiles, c.SenseUses) && uniqueItemKeys(c.SenseReveals)
	case "dragon-slayer":
		c := d.Slayer
		ok = positive(c.AttackBonus, c.DamageBonus) && dice(c.ExtraDamage)
	case "dust-of-disappearance":
		c := d.Disappearance
		ok = positive(c.Uses, c.AreaEmanationFeet) && dice(c.Duration) && uniqueItemKeys(c.Affects) && uniqueItemKeys(c.EndsForCreatureOn)
	case "dust-of-dryness":
		c := d.Dryness
		ok = dice(c.Pinches) && positive(c.PinchesBonus, c.WaterCubeFeet, c.ElementalRange, c.SaveDC) && dice(c.Damage)
	case "dust-of-sneezing-and-choking":
		c := d.Sneezing
		ok = positive(c.Uses, c.AreaEmanationFeet, c.SaveDC) && uniqueItemKeys(c.AutoSucceed) && uniqueItemKeys(c.FailureConditions) && uniqueItemKeys(c.EndsOn)
	case "dwarven-plate":
		c := d.Plate
		ok = positive(c.ACBonus, c.ReduceDistanceMax)
	case "dwarven-thrower":
		c := d.Thrower
		ok = positive(c.AttackBonus, c.DamageBonus, c.NormalRangeFeet, c.LongRangeFeet) && c.LongRangeFeet > c.NormalRangeFeet && dice(c.RangedExtraDamage) && dice(c.GiantExtraDamage) && uniqueItemKeys(c.AttunementBy)
	case "efficient-quiver":
		c := d.Quiver
		ok = positive(c.MaxWeightPounds) && len(c.Compartments) == 3
		for _, compartment := range c.Compartments {
			ok = ok && positive(compartment.Capacity) && uniqueItemKeys(compartment.Holds)
		}
	case "efreeti-bottle":
		c := d.Bottle
		rows := [][2]int{}
		for _, row := range c.Effects {
			rows = append(rows, [2]int{row.Min, row.Max})
			ok = ok && uniqueItemKeys(row.Ends)
		}
		ok = ok && positive(c.AppearsRange) && itemRollTable(rows, c.TableDieSides)
	case "elemental-gem":
		// The gem table, p. 220: each gem's elemental.
		elemental := map[string]string{"blue-sapphire": "air-elemental", "emerald": "water-elemental", "red-corundum": "fire-elemental", "yellow-diamond": "earth-elemental"}
		c := d.Gem
		ok = c.Gem == variant && c.Elemental == elemental[variant] && positive(c.DurationSeconds) && uniqueItemKeys(c.EndTriggers)
	case "elixir-of-health":
		ok = uniqueItemKeys(d.Elixir.EndsConditions)
	case "elven-chain":
		c := d.Chain
		ok = positive(c.ACBonus) && uniqueItemKeys(c.TrainedEvenWithout)
	}
	if !ok {
		return fmt.Errorf("incomplete or inconsistent source values")
	}
	return nil
}
