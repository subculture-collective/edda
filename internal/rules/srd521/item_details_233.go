package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_233.json
var items233JSON []byte

// ItemDetails233 holds exactly one typed payload per variant of catalogue batch
// #233 (magic items 09, SRD 5.2.1 pp. 225-227, plus the Ioun Stone text that
// continues on p. 228). Payloads are source facts only.
type ItemDetails233 struct {
	HelmBrilliance *ItemHelmOfBrilliance             `json:"helm_of_brilliance,omitempty"`
	HelmLanguages  *ItemHelmOfComprehendingLanguages `json:"helm_of_comprehending_languages,omitempty"`
	HelmTelepathy  *ItemHelmOfTelepathy              `json:"helm_of_telepathy,omitempty"`
	HelmTeleport   *ItemHelmOfTeleportation          `json:"helm_of_teleportation,omitempty"`
	HolyAvenger    *ItemHolyAvenger                  `json:"holy_avenger,omitempty"`
	HornBlasting   *ItemHornOfBlasting               `json:"horn_of_blasting,omitempty"`
	HornValhalla   *ItemHornOfValhalla               `json:"horn_of_valhalla,omitempty"`
	HorseshoesFast *ItemHorseshoesOfSpeed            `json:"horseshoes_of_speed,omitempty"`
	HorseshoesAir  *ItemHorseshoesOfAZephyr          `json:"horseshoes_of_a_zephyr,omitempty"`
	ImmovableRod   *ItemImmovableRod                 `json:"immovable_rod,omitempty"`
	Fortress       *ItemInstantFortress              `json:"instant_fortress,omitempty"`
	IounStone      *ItemIounStone                    `json:"ioun_stone,omitempty"`
}

// ItemFact233 is one numeric source fact: a key, a non-negative amount and the
// unit it is printed in.
type ItemFact233 struct {
	Key    string `json:"key"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// ItemHelmOfBrilliance retains typed SRD source facts without implementing them.
// The text starts on p. 225 and continues on p. 226.
type ItemHelmOfBrilliance struct {
	SourcePages          []int                  `json:"source_pages"`
	Gems                 []ItemHelmGem          `json:"gems"`
	GemsCrumbleWhenPried bool                   `json:"gems_crumble_when_pried"`
	LosesMagicWithoutGem bool                   `json:"loses_magic_when_all_gems_gone"`
	Condition            string                 `json:"condition"`
	DiamondLight         ItemHelmDiamondLight   `json:"diamond_light"`
	FireOpalFlames       ItemHelmFireOpalFlames `json:"fire_opal_flames"`
	RubyResistance       ItemHelmRubyResistance `json:"ruby_resistance"`
	Spells               ItemHelmSpells         `json:"spells"`
	FireDamage           ItemHelmFireDamage     `json:"taking_fire_damage"`
}

// ItemHelmGem is one gem type and the dice that set how many the helm has.
type ItemHelmGem struct {
	Gem   string `json:"gem"`
	Dice  int    `json:"dice"`
	Sides int    `json:"sides"`
}

// ItemHelmDiamondLight is the Diamond Light property.
type ItemHelmDiamondLight struct {
	Requires     string   `json:"requires"`
	EmanationFee int      `json:"emanation_feet"`
	DimLightWhen string   `json:"dim_light_when"`
	Trigger      string   `json:"damage_trigger"`
	Damage       ItemDice `json:"damage"`
}

// ItemHelmFireOpalFlames is the Fire Opal Flames property.
type ItemHelmFireOpalFlames struct {
	Requires               string   `json:"requires"`
	Action                 string   `json:"action"`
	Target                 string   `json:"target"`
	BrightLightFeet        int      `json:"bright_light_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	HarmlessTo             []string `json:"harmless_to"`
	ExtraDamage            ItemDice `json:"extra_damage"`
	ExtraDamageTrigger     string   `json:"extra_damage_trigger"`
	EndsOn                 []string `json:"ends_on"`
}

// ItemHelmRubyResistance is the Ruby Resistance property.
type ItemHelmRubyResistance struct {
	Requires   string `json:"requires"`
	Resistance string `json:"resistance"`
}

// ItemHelmSpells is the Spells property.
type ItemHelmSpells struct {
	SaveDC         int             `json:"save_dc"`
	Component      string          `json:"component"`
	GemDestroyed   bool            `json:"gem_destroyed_when_cast"`
	GemDisappears  bool            `json:"gem_disappears_from_helm"`
	SpellsWithGems []ItemHelmSpell `json:"spells"`
}

// ItemHelmSpell pairs a spell with the gem type used as its component.
type ItemHelmSpell struct {
	Spell string `json:"spell"`
	Gem   string `json:"gem"`
}

// ItemHelmFireDamage is the Taking Fire Damage property.
type ItemHelmFireDamage struct {
	Trigger       string `json:"trigger"`
	RollDieSides  int    `json:"roll_die_sides"`
	DestroyedOn   int    `json:"destroyed_on_roll"`
	Effect        string `json:"effect"`
	HelmDestroyed bool   `json:"helm_destroyed"`
	EmanationFeet int    `json:"emanation_feet"`
	Origin        string `json:"origin"`
	SaveAbility   string `json:"save_ability"`
	SaveDC        int    `json:"save_dc"`
	DamageType    string `json:"damage_type"`
	Damage        string `json:"damage"`
}

// ItemHelmOfComprehendingLanguages retains typed SRD source facts without implementing them.
type ItemHelmOfComprehendingLanguages struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
}

// ItemHelmOfTelepathy retains typed SRD source facts without implementing them.
type ItemHelmOfTelepathy struct {
	Condition        string   `json:"condition"`
	TelepathyFeet    int      `json:"telepathy_range_feet"`
	Spells           []string `json:"spells"`
	SaveDCSpell      string   `json:"save_dc_spell"`
	SaveDC           int      `json:"save_dc"`
	SpellReuse       string   `json:"spell_reuse"`
	SpellReuseWindow string   `json:"spell_reuse_window"`
}

// ItemHelmOfTeleportation retains typed SRD source facts without implementing them.
type ItemHelmOfTeleportation struct {
	Condition    string   `json:"condition"`
	MaxCharges   int      `json:"max_charges"`
	Spell        string   `json:"spell"`
	ChargeCost   int      `json:"charge_cost"`
	RegainDice   ItemDice `json:"regain_dice"`
	RegainTiming string   `json:"regain_timing"`
}

// ItemHolyAvenger retains typed SRD source facts without implementing them.
type ItemHolyAvenger struct {
	AttunementBy       []string            `json:"attunement_by_any_of"`
	AttackBonus        int                 `json:"attack_bonus"`
	DamageBonus        int                 `json:"damage_bonus"`
	ExtraDamage        ItemDice            `json:"extra_damage"`
	ExtraDamageAgainst []string            `json:"extra_damage_against"`
	ExtraDamageTrigger string              `json:"extra_damage_trigger"`
	Aura               ItemHolyAvengerAura `json:"aura"`
}

// ItemHolyAvengerAura is the Emanation created while the weapon is held drawn.
type ItemHolyAvengerAura struct {
	Condition       string   `json:"condition"`
	EmanationFeet   int      `json:"emanation_feet"`
	Origin          string   `json:"origin"`
	Affected        []string `json:"affected"`
	Benefit         string   `json:"benefit"`
	Against         []string `json:"against"`
	LargerMinLevels int      `json:"larger_emanation_min_paladin_levels"`
	LargerFeet      int      `json:"larger_emanation_feet"`
}

// ItemHornOfBlasting retains typed SRD source facts without implementing them.
type ItemHornOfBlasting struct {
	Action                string   `json:"action"`
	ConeFeet              int      `json:"cone_feet"`
	AudibleFeet           int      `json:"audible_feet"`
	SaveAbility           string   `json:"save_ability"`
	SaveDC                int      `json:"save_dc"`
	Damage                ItemDice `json:"damage"`
	Condition             string   `json:"condition"`
	ConditionSeconds      int      `json:"condition_duration_seconds"`
	SuccessResult         string   `json:"success_result"`
	ObjectTargets         string   `json:"object_targets"`
	ObjectDamage          ItemDice `json:"object_damage"`
	ExplosionPercent      int      `json:"explosion_chance_percent"`
	ExplosionDamage       ItemDice `json:"explosion_damage"`
	ExplosionTarget       string   `json:"explosion_target"`
	ExplosionDestroysHorn bool     `json:"explosion_destroys_horn"`
}

// ItemHornOfValhalla is the shared text of the Horn of Valhalla and the table
// row of one horn type. A variant is one metal; a printed 00 is stored as 100.
type ItemHornOfValhalla struct {
	Action          string   `json:"action"`
	SpawnRangeFeet  int      `json:"spawn_range_feet"`
	SpiritsFrom     string   `json:"spirits_from_plane"`
	SpiritStatBlock string   `json:"spirit_stat_block"`
	SpiritAppears   string   `json:"spirit_appearance"`
	DurationHours   int      `json:"duration_hours"`
	ReturnTriggers  []string `json:"return_triggers"`
	Immunities      []string `json:"spirit_immunities"`
	CooldownDays    int      `json:"cooldown_days"`
	TypesKnown      int      `json:"types_known"`
	Type            string   `json:"type"`
	TableDieSides   int      `json:"table_die_sides"`
	TableMin        int      `json:"table_min"`
	TableMax        int      `json:"table_max"`
	Spirits         int      `json:"spirits"`
	Requirement     string   `json:"requirement"`
	UnmetResult     string   `json:"unmet_requirement_result"`
	MetResult       []string `json:"met_requirement_result"`
}

// ItemHorseshoesOfSpeed retains typed SRD source facts without implementing them.
type ItemHorseshoesOfSpeed struct {
	SetSize         int    `json:"set_size"`
	AffixAction     string `json:"affix_action"`
	AffixTarget     string `json:"affix_target"`
	RemoveAction    string `json:"remove_action"`
	RequiresAllFour bool   `json:"requires_all_affixed_to_same_creature"`
	SpeedBonusFeet  int    `json:"speed_bonus_feet"`
}

// ItemHorseshoesOfAZephyr retains typed SRD source facts without implementing them.
type ItemHorseshoesOfAZephyr struct {
	SetSize         int      `json:"set_size"`
	AffixAction     string   `json:"affix_action"`
	AffixTarget     string   `json:"affix_target"`
	RemoveAction    string   `json:"remove_action"`
	RequiresAllFour bool     `json:"requires_all_affixed"`
	FloatInches     int      `json:"float_inches"`
	MovesNormally   bool     `json:"moves_normally"`
	CrossSurfaces   string   `json:"cross_surfaces"`
	SurfaceExamples []string `json:"surface_examples"`
	LeavesTracks    bool     `json:"leaves_tracks"`
	IgnoresDiffTerr bool     `json:"ignores_difficult_terrain"`
	TravelHours     int      `json:"travel_hours_per_day"`
	NoExhaustionFor string   `json:"no_exhaustion_from"`
}

// ItemImmovableRod retains typed SRD source facts without implementing them.
type ItemImmovableRod struct {
	Material        string    `json:"material"`
	ButtonOn        string    `json:"button_on"`
	ButtonAction    string    `json:"button_action"`
	FixedEffect     string    `json:"fixed_effect"`
	IgnoresGravity  bool      `json:"holds_against_gravity"`
	ReleasedBy      []string  `json:"released_by"`
	MaxWeightPounds int       `json:"max_weight_pounds"`
	OverloadResult  string    `json:"overload_result"`
	MoveAction      string    `json:"move_action"`
	MoveCheck       ItemCheck `json:"move_check"`
	MoveFeet        int       `json:"move_feet"`
	MoveOn          string    `json:"move_on"`
}

// ItemInstantFortress retains typed SRD source facts without implementing them.
type ItemInstantFortress struct {
	StatuetteInches    int                `json:"statuette_inches"`
	StatuetteMaterial  string             `json:"statuette_material"`
	Action             string             `json:"action"`
	Activation         string             `json:"activation"`
	PlacedOn           string             `json:"placed_on"`
	Tower              ItemFortressTower  `json:"tower"`
	RevertActivation   string             `json:"revert_activation"`
	RevertCondition    string             `json:"revert_condition"`
	DisplacedCreatures string             `json:"displaced_creatures"`
	DisplacedObjects   string             `json:"displaced_objects"`
	DamageRepairedBy   ItemFortressRepair `json:"damage_repair"`
}

// ItemFortressTower is the tower the statuette becomes.
type ItemFortressTower struct {
	Shape           string   `json:"shape"`
	SideFeet        int      `json:"side_feet"`
	HeightFeet      int      `json:"height_feet"`
	Features        []string `json:"features"`
	Floors          int      `json:"floors"`
	FloorConnection []string `json:"floor_connection_your_choice"`
	RoofAccess      string   `json:"roof_access"`
	Doors           int      `json:"doors"`
	DoorLocation    string   `json:"door_location"`
	DoorOpensBy     string   `json:"door_opens_by"`
	DoorOpenAction  string   `json:"door_open_action"`
	DoorImmuneTo    string   `json:"door_immune_to"`
	CannotBeTipped  bool     `json:"cannot_be_tipped"`
	DurableParts    []string `json:"durable_parts"`
	ArmorClass      int      `json:"armor_class"`
	HitPoints       int      `json:"hit_points"`
	ImmuneTo        []string `json:"immune_to"`
	ImmuneException string   `json:"immunity_exception"`
	ResistantTo     string   `json:"resistant_to"`
}

// ItemFortressRepair records how the tower's damage is repaired.
type ItemFortressRepair struct {
	ShrinkingRepairs bool   `json:"shrinking_repairs_damage"`
	OnlySpell        string `json:"only_spell"`
	SpellLevelAtMost int    `json:"counts_as_spell_level_at_most"`
	EachCasting      string `json:"each_casting_restores"`
}

// ItemIounStone is the shared Ioun Stone rules and one stone type. SourcePages
// is the page that prints the type: p. 227, or p. 228 from Mastery on. Ability
// is "none" for a stone without an ability score. Facts and lists are keys.
type ItemIounStone struct {
	SourcePages []int          `json:"source_pages"`
	Shared      ItemIounShared `json:"shared_rules"`
	Color       string         `json:"color"`
	Shape       string         `json:"shape"`
	Benefit     string         `json:"benefit"`
	Action      string         `json:"action"`
	Ability     string         `json:"ability"`
	Facts       []ItemFact233  `json:"facts" srd:"may-be-empty"`
	Traits      []string       `json:"traits"`
}

// ItemIounShared holds the rules every Ioun Stone type shares.
type ItemIounShared struct {
	Size           string   `json:"size"`
	Activation     string   `json:"activation_action"`
	Method         string   `json:"method"`
	OrbitDiceCount int      `json:"orbit_distance_dice"`
	OrbitDiceSides int      `json:"orbit_distance_sides"`
	MaxOrbiting    int      `json:"max_orbiting"`
	CountsAs       string   `json:"orbiting_stone_counts_as"`
	OrbitTraits    []string `json:"orbit_traits"`
	StowAction     string   `json:"stow_action"`
	StowCount      string   `json:"stow_count"`
	AttunementEnds string   `json:"attunement_ends_result"`
}

// itemManifest233 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Horn of Valhalla has one variant per metal and Ioun
// Stone one per stone type, as the source lists them.
var itemManifest233 = itemSourceManifest{Batch: "data-magic-items-09", Issue: 233, Families: []itemSourceManifestFamily{
	{Slug: "helm-of-brilliance", Name: "Helm of Brilliance", Page: 225, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "gem-required",
		Keys: "item.attunement item.wear item.gems-1d10-diamonds-2d10-rubies-3d10-fire-opals-4d10-opals item.loses-magic-without-gems light.emanation-30-foot damage.radiant-1d6-undead action.magic damage.fire-1d6-blazing-weapon light.bright-10-dim-10 defense.resistance-fire spell.save-dc-18 spell.daylight-fireball-prismatic-spray-wall-of-fire item.gem-destroyed-on-cast save.dexterity-dc-17 damage.radiant-per-gem", Payload: "helm_of_brilliance",
		Variants: []itemSourceManifestVariant{{"standard", "Helm of Brilliance", "very-rare"}}},
	{Slug: "helm-of-comprehending-languages", Name: "Helm of Comprehending Languages", Page: 226, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "cast-spell",
		Keys: "item.wear spell.cast-comprehend-languages", Payload: "helm_of_comprehending_languages",
		Variants: []itemSourceManifestVariant{{"standard", "Helm of Comprehending Languages", "uncommon"}}},
	{Slug: "helm-of-telepathy", Name: "Helm of Telepathy", Page: 226, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "each-spell-once-until-next-dawn",
		Keys: "item.attunement item.wear sense.telepathy-30-feet spell.cast-detect-thoughts spell.cast-suggestion spell.suggestion-save-dc-13 item.once-until-next-dawn", Payload: "helm_of_telepathy",
		Variants: []itemSourceManifestVariant{{"standard", "Helm of Telepathy", "uncommon"}}},
	{Slug: "helm-of-teleportation", Name: "Helm of Teleportation", Page: 226, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.wear item.charges spell.cast-teleport item.regain-1d3-charges-at-dawn", Payload: "helm_of_teleportation",
		Variants: []itemSourceManifestVariant{{"standard", "Helm of Teleportation", "rare"}}},
	{Slug: "holy-avenger", Name: "Holy Avenger", Page: 226, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement attunement.paladin bonus.attack-and-damage-3 damage.radiant-2d10-fiend-undead aura.emanation-10-foot save.advantage-spells-and-magical-effects aura.emanation-30-foot-paladin-17", Payload: "holy_avenger",
		Variants: []itemSourceManifestVariant{{"standard", "Holy Avenger", "legendary"}}},
	{Slug: "horn-of-blasting", Name: "Horn of Blasting", Page: 226, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "explosion-chance-20-percent",
		Keys: "action.magic area.cone-30-foot sound.audible-600-feet save.constitution-dc-15 damage.thunder-5d8 condition.deafened damage.thunder-10d8-glass-crystal damage.force-10d6-explosion item.destroyed-by-explosion", Payload: "horn_of_blasting",
		Variants: []itemSourceManifestVariant{{"standard", "Horn of Blasting", "rare"}}},
	{Slug: "horn-of-valhalla", Name: "Horn of Valhalla", Page: 226, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "use-cooldown",
		Keys: "action.magic creature.spirits-berserker-stat-block creature.returns-to-ysgard immunity.charmed-and-frightened item.use-cooldown item.type-by-gm-or-1d100-table creature.attack-you-without-requirement creature.friendly-with-requirement", Payload: "horn_of_valhalla",
		Variants: []itemSourceManifestVariant{{"silver", "Silver", "rare"}, {"brass", "Brass", "rare"}, {"bronze", "Bronze", "very-rare"}, {"iron", "Iron", "legendary"}}},
	{Slug: "horseshoes-of-speed", Name: "Horseshoes of Speed", Page: 227, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "touch-to-hoof",
		Keys: "action.magic item.set-of-four item.affixes-to-hoof speed.bonus-30-feet-with-all-four", Payload: "horseshoes_of_speed",
		Variants: []itemSourceManifestVariant{{"standard", "Horseshoes of Speed", "rare"}}},
	{Slug: "horseshoes-of-a-zephyr", Name: "Horseshoes of a Zephyr", Page: 226, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "touch-to-hoof",
		Keys: "action.magic item.set-of-four item.affixes-to-hoof movement.float-4-inches terrain.cross-nonsolid-surfaces movement.leaves-no-tracks terrain.ignores-difficult travel.12-hours-no-exhaustion", Payload: "horseshoes_of_a_zephyr",
		Variants: []itemSourceManifestVariant{{"standard", "Horseshoes of a Zephyr", "very-rare"}}},
	{Slug: "immovable-rod", Name: "Immovable Rod", Page: 227, Category: "rod", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "press-button",
		Keys: "action.utilize item.fixed-in-place item.holds-8000-pounds check.strength-athletics-dc-30 movement.rod-10-feet", Payload: "immovable_rod",
		Variants: []itemSourceManifestVariant{{"standard", "Immovable Rod", "uncommon"}}},
	{Slug: "instant-fortress", Name: "Instant Fortress", Page: 227, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "command-word",
		Keys: "item.attunement action.magic item.command-word object.tower-20-foot-30-high defense.armor-class-20 defense.hit-points-100 defense.immune-nonsiege-weapon-damage defense.resistance-other-damage spell.wish-repairs-tower", Payload: "instant_fortress",
		Variants: []itemSourceManifestVariant{{"standard", "Instant Fortress", "rare"}}},
	{Slug: "ioun-stone", Name: "Ioun Stone", Page: 227, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "toss-into-air up-to-three-orbiting",
		Keys: "item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone", Payload: "ioun_stone",
		Variants: []itemSourceManifestVariant{
			{"absorption", "Absorption", "very-rare"}, {"agility", "Agility", "very-rare"}, {"awareness", "Awareness", "rare"}, {"fortitude", "Fortitude", "very-rare"},
			{"greater-absorption", "Greater Absorption", "legendary"}, {"insight", "Insight", "very-rare"}, {"intellect", "Intellect", "very-rare"}, {"leadership", "Leadership", "very-rare"},
			{"mastery", "Mastery", "legendary"}, {"protection", "Protection", "rare"}, {"regeneration", "Regeneration", "legendary"}, {"reserve", "Reserve", "rare"},
			{"strength", "Strength", "very-rare"}, {"sustenance", "Sustenance", "rare"}}},
}}

// LoadItems233 returns a fresh validated copy of batch #233. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems233() (ItemSourceBatch[ItemDetails233], error) { return DecodeItems233(items233JSON) }

// DecodeItems233 validates supplied JSON as batch #233 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems233(data []byte) (ItemSourceBatch[ItemDetails233], error) {
	return decodeItemSourceBatch(data, itemManifest233, validateItemDetails233)
}

// validateItemFacts233 requires distinct, non-empty fact keys.
func validateItemFacts233(facts []ItemFact233) bool {
	keys := []string{}
	for _, f := range facts {
		keys = append(keys, f.Key)
	}
	return len(keys) == 0 || uniqueItemKeys(keys)
}

// validateItemDetails233 checks the internal consistency of each payload and a
// few source literals that the family's mechanic keys also name.
func validateItemDetails233(family, variant string, d *ItemDetails233) error {
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
	case "helm-of-brilliance":
		h := d.HelmBrilliance
		gems := map[string][2]int{"diamond": {1, 10}, "ruby": {2, 10}, "fire-opal": {3, 10}, "opal": {4, 10}} // p. 225
		ok = slices.Equal(h.SourcePages, []int{225, 226}) && len(h.Gems) == len(gems) && positive(h.DiamondLight.EmanationFee, h.DiamondLight.Damage.Dice, h.DiamondLight.Damage.Sides,
			h.FireOpalFlames.BrightLightFeet, h.FireOpalFlames.DimLightAdditionalFeet, h.FireOpalFlames.ExtraDamage.Dice, h.FireOpalFlames.ExtraDamage.Sides,
			h.Spells.SaveDC, h.FireDamage.RollDieSides, h.FireDamage.DestroyedOn, h.FireDamage.EmanationFeet, h.FireDamage.SaveDC) &&
			uniqueItemKeys(h.FireOpalFlames.HarmlessTo) && uniqueItemKeys(h.FireOpalFlames.EndsOn) && len(h.Spells.SpellsWithGems) == 4 &&
			h.Spells.SaveDC == 18 && h.FireDamage.SaveDC == 17 // p. 225-226; named by the mechanic keys
		for _, g := range h.Gems {
			want, found := gems[g.Gem]
			ok = ok && found && [2]int{g.Dice, g.Sides} == want
		}
		spells := []string{}
		for _, s := range h.Spells.SpellsWithGems {
			spells = append(spells, s.Spell)
		}
		ok = ok && uniqueItemKeys(spells)
	case "helm-of-comprehending-languages":
		ok = d.HelmLanguages.Spell == "comprehend-languages"
	case "helm-of-telepathy":
		h := d.HelmTelepathy
		ok = positive(h.TelepathyFeet, h.SaveDC) && uniqueItemKeys(h.Spells) && h.SaveDC == 13 && slices.Contains(h.Spells, h.SaveDCSpell)
	case "helm-of-teleportation":
		h := d.HelmTeleport
		ok = positive(h.MaxCharges, h.ChargeCost, h.RegainDice.Dice, h.RegainDice.Sides) && h.MaxCharges == 3
	case "holy-avenger":
		h := d.HolyAvenger
		ok = positive(h.AttackBonus, h.DamageBonus, h.ExtraDamage.Dice, h.ExtraDamage.Sides, h.Aura.EmanationFeet, h.Aura.LargerMinLevels, h.Aura.LargerFeet) &&
			h.AttackBonus == 3 && h.DamageBonus == 3 && uniqueItemKeys(h.AttunementBy) && uniqueItemKeys(h.ExtraDamageAgainst) && uniqueItemKeys(h.Aura.Affected) &&
			uniqueItemKeys(h.Aura.Against) && h.Aura.LargerFeet > h.Aura.EmanationFeet
	case "horn-of-blasting":
		h := d.HornBlasting
		ok = positive(h.ConeFeet, h.AudibleFeet, h.SaveDC, h.Damage.Dice, h.Damage.Sides, h.ConditionSeconds, h.ObjectDamage.Dice, h.ObjectDamage.Sides, h.ExplosionPercent,
			h.ExplosionDamage.Dice, h.ExplosionDamage.Sides) && h.SaveDC == 15 && h.ExplosionPercent == 20
	case "horn-of-valhalla":
		h := d.HornValhalla
		// The four rows of the 1d100 table, in order; 00 is stored as 100.
		rows := map[string][3]int{"silver": {1, 40, 2}, "brass": {41, 75, 3}, "bronze": {76, 90, 4}, "iron": {91, 100, 5}}
		ok = h.Type == variant && h.TableDieSides == 100 && [3]int{h.TableMin, h.TableMax, h.Spirits} == rows[variant] && positive(h.SpawnRangeFeet, h.DurationHours, h.CooldownDays, h.TypesKnown) &&
			h.TypesKnown == len(rows) && uniqueItemKeys(h.ReturnTriggers) && uniqueItemKeys(h.Immunities) && uniqueItemKeys(h.MetResult)
	case "horseshoes-of-speed":
		ok = positive(d.HorseshoesFast.SetSize, d.HorseshoesFast.SpeedBonusFeet) && d.HorseshoesFast.SetSize == 4 && d.HorseshoesFast.SpeedBonusFeet == 30
	case "horseshoes-of-a-zephyr":
		h := d.HorseshoesAir
		ok = positive(h.SetSize, h.FloatInches, h.TravelHours) && h.SetSize == 4 && h.FloatInches == 4 && h.TravelHours == 12 && uniqueItemKeys(h.SurfaceExamples)
	case "immovable-rod":
		r := d.ImmovableRod
		ok = positive(r.MaxWeightPounds, r.MoveCheck.DC, r.MoveFeet) && r.MaxWeightPounds == 8000 && r.MoveCheck.DC == 30 && uniqueItemKeys(r.ReleasedBy)
	case "instant-fortress":
		f := d.Fortress
		t := f.Tower
		ok = positive(f.StatuetteInches, t.SideFeet, t.HeightFeet, t.Floors, t.Doors, t.ArmorClass, t.HitPoints, f.DamageRepairedBy.SpellLevelAtMost) &&
			t.ArmorClass == 20 && t.HitPoints == 100 && t.SideFeet == 20 && t.HeightFeet == 30 && uniqueItemKeys(t.Features) && uniqueItemKeys(t.FloorConnection) &&
			uniqueItemKeys(t.DurableParts) && uniqueItemKeys(t.ImmuneTo)
	case "ioun-stone":
		s := d.IounStone
		sh := s.Shared
		// The stones from Mastery on are printed on p. 228; the shared text and the others are on p. 227.
		page := 227
		if variant == "mastery" || variant == "protection" || variant == "regeneration" || variant == "reserve" || variant == "strength" || variant == "sustenance" {
			page = 228
		}
		ok = slices.Equal(s.SourcePages, []int{page}) && positive(sh.OrbitDiceCount, sh.OrbitDiceSides, sh.MaxOrbiting) && sh.MaxOrbiting == 3 &&
			uniqueItemKeys(sh.OrbitTraits) && validateItemFacts233(s.Facts) && uniqueItemKeys(s.Traits)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
