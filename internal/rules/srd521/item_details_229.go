package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_229.json
var items229JSON []byte

// ItemDetails229 holds exactly one typed payload per variant of catalogue batch
// #229 (magic items 05, SRD 5.2.1 pp. 216–218). Payloads are source facts only.
type ItemDetails229 struct {
	Telepathy  *ItemCrystalBallTelepathy  `json:"telepathy,omitempty"`
	TrueSeeing *ItemCrystalBallTrueSeeing `json:"true_seeing,omitempty"`
	CubeForce  *ItemCubeOfForce           `json:"cube_of_force,omitempty"`
	CubicGate  *ItemCubicGate             `json:"cubic_gate,omitempty"`
	Venom      *ItemDaggerOfVenom         `json:"venom,omitempty"`
	Dancing    *ItemDancingSword          `json:"dancing,omitempty"`
	Decanter   *ItemDecanterEndlessWater  `json:"decanter,omitempty"`
	Illusions  *ItemDeckOfIllusions       `json:"illusions,omitempty"`
	Defender   *ItemDefender              `json:"defender,omitempty"`
	DemonArmor *ItemDemonArmor            `json:"demon_armor,omitempty"`
	Shackles   *ItemDimensionalShackles   `json:"shackles,omitempty"`
	DragonOrb  *ItemDragonOrb             `json:"dragon_orb,omitempty"`
}

// ItemCrystalBallTelepathy retains typed SRD source facts without implementing them.
type ItemCrystalBallTelepathy struct {
	Condition              string `json:"condition"`
	Spell                  string `json:"spell"`
	SaveDC                 int    `json:"save_dc"`
	CommunicationMode      string `json:"communication_mode"`
	CommunicationTargets   string `json:"communication_targets"`
	CommunicationRangeFeet int    `json:"communication_range_feet"`
	CommunicationRangeFrom string `json:"communication_range_from"`
	ExtraSpell             string `json:"extra_spell"`
	ExtraSaveDC            int    `json:"extra_save_dc"`
	ExtraTargets           string `json:"extra_targets"`
	ExtraDelivery          string `json:"extra_delivery"`
	ConcentrationRequired  bool   `json:"extra_concentration_required"`
	ExtraEndsWith          string `json:"extra_ends_with"`
	ExtraUses              int    `json:"extra_uses"`
	ExtraRecharge          string `json:"extra_recharge"`
}

// ItemCrystalBallTrueSeeing retains typed SRD source facts without implementing them.
type ItemCrystalBallTrueSeeing struct {
	Condition       string `json:"condition"`
	Spell           string `json:"spell"`
	SaveDC          int    `json:"save_dc"`
	Sense           string `json:"sense"`
	SenseRangeFeet  int    `json:"sense_range_feet"`
	SenseCenteredOn string `json:"sense_centered_on"`
}

// ItemCubeFace is one row of the Cube of Force Faces table.
type ItemCubeFace struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemCubeOfForce retains typed SRD source facts without implementing them.
// The charge cost varies by face, so the variant's per-use charge is zero and
// each cost lives in the Faces table.
type ItemCubeOfForce struct {
	Appearance      string         `json:"appearance"`
	WidthInches     int            `json:"width_inches"`
	Markings        string         `json:"markings"`
	Activation      string         `json:"activation"`
	Cost            string         `json:"cost"`
	SaveDC          int            `json:"save_dc"`
	StartingCharges int            `json:"starting_charges"`
	RegainDice      ItemDice       `json:"regain_dice"`
	RegainTiming    string         `json:"regain_timing"`
	Faces           []ItemCubeFace `json:"faces"`
}

// ItemCubicGate retains typed SRD source facts without implementing them.
type ItemCubicGate struct {
	Appearance              string   `json:"appearance"`
	WidthInches             int      `json:"width_inches"`
	Radiates                string   `json:"radiates"`
	Sides                   int      `json:"sides"`
	MaterialPlaneSides      int      `json:"material_plane_sides"`
	OtherPlanesDeterminedBy string   `json:"other_planes_determined_by"`
	MaxCharges              int      `json:"max_charges"`
	RegainDice              ItemDice `json:"regain_dice"`
	RegainTiming            string   `json:"regain_timing"`
	Action                  string   `json:"action"`
	ChargeCost              int      `json:"charge_cost"`
	GateSpell               string   `json:"gate_spell"`
	GateActivation          string   `json:"gate_activation"`
	GateEffect              string   `json:"gate_effect"`
	PlaneShiftSpell         string   `json:"plane_shift_spell"`
	PlaneShiftActivation    string   `json:"plane_shift_activation"`
	PlaneShiftEffect        string   `json:"plane_shift_effect"`
}

// ItemDaggerOfVenom retains typed SRD source facts without implementing them.
type ItemDaggerOfVenom struct {
	AttackBonus              int      `json:"attack_bonus"`
	DamageBonus              int      `json:"damage_bonus"`
	ApplyAction              string   `json:"apply_action"`
	Coats                    string   `json:"coats"`
	PoisonDurationSeconds    int      `json:"poison_duration_seconds"`
	PoisonEndsOn             []string `json:"poison_ends_on"`
	Save                     string   `json:"save"`
	SaveDC                   int      `json:"save_dc"`
	FailureDamage            ItemDice `json:"failure_damage"`
	FailureCondition         string   `json:"failure_condition"`
	ConditionDurationSeconds int      `json:"condition_duration_seconds"`
	Uses                     int      `json:"uses"`
	Recharge                 string   `json:"recharge"`
}

// ItemDancingSword retains typed SRD source facts without implementing them.
type ItemDancingSword struct {
	ActivateAction              string   `json:"activate_action"`
	Activation                  string   `json:"activation"`
	HoverFlyFeet                int      `json:"hover_fly_feet"`
	AttackRangeFeet             int      `json:"attack_range_feet"`
	AttackRoll                  string   `json:"attack_roll"`
	AddsAbilityModifierToDamage bool     `json:"adds_ability_modifier_to_damage"`
	MoveAction                  string   `json:"move_action"`
	MoveFeet                    int      `json:"move_feet"`
	MoveWithinFeetOfYou         int      `json:"move_within_feet_of_you"`
	AttackWithMove              bool     `json:"attack_with_move"`
	AttacksBeforeReturn         int      `json:"attacks_before_return"`
	ReturnsTo                   string   `json:"returns_to"`
	NoHandFree                  string   `json:"no_hand_free"`
	NoPath                      string   `json:"no_path"`
	HoverEndsOn                 []string `json:"hover_ends_on"`
	HoverEndDistanceFeet        int      `json:"hover_end_distance_feet"`
}

// ItemDecanterCommand is one command word of the Decanter of Endless Water.
type ItemDecanterCommand struct {
	Word    string `json:"word"`
	Gallons int    `json:"gallons"`
}

// ItemDecanterEndlessWater retains typed SRD source facts without implementing them.
type ItemDecanterEndlessWater struct {
	Appearance       string                `json:"appearance"`
	WeightPounds     int                   `json:"weight_pounds"`
	Action           string                `json:"action"`
	Opens            string                `json:"opens"`
	WaterTypes       []string              `json:"water_types"`
	WaterChoice      string                `json:"water_choice"`
	StopsAt          string                `json:"stops_at"`
	Commands         []ItemDecanterCommand `json:"commands"`
	GeyserShape      string                `json:"geyser_shape"`
	GeyserLengthFeet int                   `json:"geyser_length_feet"`
	GeyserWidthFeet  int                   `json:"geyser_width_feet"`
	GeyserAim        string                `json:"geyser_aim"`
	GeyserTargets    string                `json:"geyser_targets"`
	GeyserSave       string                `json:"geyser_save"`
	GeyserSaveDC     int                   `json:"geyser_save_dc"`
	GeyserDamage     ItemDice              `json:"geyser_damage"`
	GeyserCondition  string                `json:"geyser_failure_condition"`
	ObjectCondition  string                `json:"object_condition"`
	ObjectMaxPounds  int                   `json:"object_max_pounds"`
	ObjectEffect     string                `json:"object_effect"`
}

// ItemIllusionRow is one row of the Deck of Illusions 1d100 table. Row 97–00 is
// recorded as 97–100.
type ItemIllusionRow struct {
	Min      int    `json:"min"`
	Max      int    `json:"max"`
	Creature string `json:"creature"`
}

// ItemDeckOfIllusions retains typed SRD source facts without implementing
// them. The entry starts on p. 217; its table ends on TablePage.
type ItemDeckOfIllusions struct {
	FullDeckCards        int               `json:"full_deck_cards"`
	CreatureCards        int               `json:"creature_cards"`
	MirrorCards          int               `json:"mirror_cards"`
	MissingCardsDice     ItemDice          `json:"missing_cards_dice"`
	MissingCardsSubtract int               `json:"missing_cards_subtract"`
	FunctionsOnlyIf      string            `json:"functions_only_if"`
	DrawAction           string            `json:"draw_action"`
	ThrowRangeFeet       int               `json:"throw_range_feet"`
	Illusion             string            `json:"illusion"`
	CanDoHarm            bool              `json:"can_do_harm"`
	MoveAction           string            `json:"move_action"`
	MoveWhileWithinFeet  int               `json:"move_while_within_feet"`
	MoveRequires         string            `json:"move_requires"`
	MoveDistanceFeet     int               `json:"move_distance_feet"`
	MoveFrom             string            `json:"move_from"`
	PhysicalInteraction  string            `json:"physical_interaction"`
	InspectAction        string            `json:"inspect_action"`
	InspectCheck         ItemCheck         `json:"inspect_check"`
	InspectResult        string            `json:"inspect_result"`
	LastsUntil           []string          `json:"lasts_until"`
	DispelledBy          []string          `json:"dispelled_by"`
	AfterEnd             string            `json:"after_end"`
	CardReusable         bool              `json:"card_reusable"`
	StatBlockSource      string            `json:"stat_block_source"`
	StatBlockException   string            `json:"stat_block_exception"`
	TableDieSides        int               `json:"table_die_sides"`
	TablePage            int               `json:"table_page"`
	Illusions            []ItemIllusionRow `json:"illusions"`
}

// ItemDefender retains typed SRD source facts without implementing them.
type ItemDefender struct {
	AttackBonus                int      `json:"attack_bonus"`
	DamageBonus                int      `json:"damage_bonus"`
	TransferTiming             string   `json:"transfer_timing"`
	TransferFrom               []string `json:"transfer_from"`
	TransferTo                 string   `json:"transfer_to"`
	TransferAmount             string   `json:"transfer_amount"`
	ExampleAttackDamageBonus   int      `json:"example_attack_damage_bonus"`
	ExampleArmorClassBonus     int      `json:"example_armor_class_bonus"`
	AdjustmentLasts            string   `json:"adjustment_lasts"`
	MustHoldForArmorClassBonus bool     `json:"must_hold_for_armor_class_bonus"`
}

// ItemDemonArmor retains typed SRD source facts without implementing them.
type ItemDemonArmor struct {
	Condition             string   `json:"condition"`
	ArmorClassBonus       int      `json:"armor_class_bonus"`
	Language              string   `json:"language"`
	Gauntlets             string   `json:"gauntlets"`
	UnarmedDamage         ItemDice `json:"unarmed_damage"`
	UnarmedReplacesDamage string   `json:"unarmed_replaces_damage_type"`
	UnarmedAttackBonus    int      `json:"unarmed_attack_bonus"`
	UnarmedDamageBonus    int      `json:"unarmed_damage_bonus"`
	Cursed                bool     `json:"cursed"`
	CurseTrigger          string   `json:"curse_trigger"`
	CannotDoff            bool     `json:"cannot_doff"`
	DoffAllowedIf         []string `json:"doff_allowed_if"`
	CurseDisadvantage     []string `json:"curse_disadvantage"`
}

// ItemDimensionalShackles retains typed SRD source facts without implementing them.
type ItemDimensionalShackles struct {
	PlaceAction          string    `json:"place_action"`
	PlaceTargetCondition string    `json:"place_target_condition"`
	FitSizeMin           string    `json:"fit_size_min"`
	FitSizeMax           string    `json:"fit_size_max"`
	PreventsMovement     string    `json:"prevents_movement"`
	PreventedExamples    []string  `json:"prevented_examples"`
	DoesNotPrevent       string    `json:"does_not_prevent"`
	RemoveAction         string    `json:"remove_action"`
	RemovableBy          []string  `json:"removable_by"`
	EscapeIntervalDays   int       `json:"escape_interval_days"`
	EscapeCheck          ItemCheck `json:"escape_check"`
	EscapeSuccess        string    `json:"escape_success"`
}

// ItemOrbSpell is one row of the Dragon Orb spell table. SpellLevel and SaveDC
// are zero when the source gives none.
type ItemOrbSpell struct {
	Spell      string `json:"spell"`
	SpellLevel int    `json:"spell_level"`
	SaveDC     int    `json:"save_dc"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemDragonOrb retains typed SRD source facts without implementing them. The
// charge cost varies by spell, so the variant's per-use charge is zero and each
// cost lives in the Spells table.
type ItemDragonOrb struct {
	Appearance                 string         `json:"appearance"`
	DiameterInches             int            `json:"diameter_inches"`
	UsedDiameterInches         int            `json:"used_diameter_inches"`
	Interior                   string         `json:"interior"`
	PeerAction                 string         `json:"peer_action"`
	PeerRequires               string         `json:"peer_requires"`
	Save                       string         `json:"save"`
	SaveDC                     int            `json:"save_dc"`
	SuccessEffect              string         `json:"success_effect"`
	FailureCondition           string         `json:"failure_condition"`
	FailureLasts               string         `json:"failure_lasts"`
	CharmedCannotEndAttunement bool           `json:"charmed_cannot_end_attunement"`
	CharmedSpell               string         `json:"charmed_spell"`
	CharmedSpellFrequency      string         `json:"charmed_spell_frequency"`
	CharmedSpellSaveDC         int            `json:"charmed_spell_save_dc"`
	CharmedUrge                string         `json:"charmed_urge"`
	DesireExamples             []string       `json:"desire_examples"`
	MaxCharges                 int            `json:"max_charges"`
	RegainDice                 ItemDice       `json:"regain_dice"`
	RegainBonus                int            `json:"regain_bonus"`
	RegainTiming               string         `json:"regain_timing"`
	SpellsRequire              string         `json:"spells_require"`
	Spells                     []ItemOrbSpell `json:"spells"`
	CallAction                 string         `json:"call_action"`
	CallRequires               string         `json:"call_requires"`
	CallRangeMiles             int            `json:"call_range_miles"`
	CallDirections             string         `json:"call_directions"`
	CallTargets                string         `json:"call_targets"`
	CallEffect                 string         `json:"call_effect"`
	CallUnaffected             string         `json:"call_unaffected"`
	CallAttitude               string         `json:"call_attitude"`
	CallCooldownSeconds        int            `json:"call_cooldown_seconds"`
	ArmorClass                 int            `json:"armor_class"`
	DestroyedBy                []string       `json:"destroyed_by"`
	NothingElseHarms           bool           `json:"nothing_else_harms"`
}

const itemAnyArmor229 = "any-light-armor any-medium-armor any-heavy-armor"

// itemManifest229 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Every family has the single variant standard. The
// Cube of Force and Dragon Orb charge costs vary per spell, so their per-use
// charge is zero and the costs are table rows in the payload.
var itemManifest229 = itemSourceManifest{Batch: "data-magic-items-05", Issue: 229, Families: []itemSourceManifestFamily{
	{Slug: "crystal-ball-of-telepathy", Name: "Crystal Ball of Telepathy", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "touch", Use: "cast-spell", UseRestrictions: "touching-crystal-orb",
		Keys: "item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 communication.telepathy-within-30-feet-of-sensor spell.cast-suggestion spell.suggestion-save-dc-17 concentration.suggestion-ends-with-scrying item.once-until-next-dawn", Payload: "telepathy",
		Variants: []itemSourceManifestVariant{{"standard", "Crystal Ball of Telepathy", "legendary"}}},
	{Slug: "crystal-ball-of-true-seeing", Name: "Crystal Ball of True Seeing", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "touch", Use: "cast-spell", UseRestrictions: "touching-crystal-orb",
		Keys: "item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 senses.truesight-120-feet-at-sensor", Payload: "true_seeing",
		Variants: []itemSourceManifestVariant{{"standard", "Crystal Ball of True Seeing", "legendary"}}},
	{Slug: "cube-of-force", Name: "Cube of Force", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "press-face", UseRestrictions: "expend-charges-required-for-face",
		Charges: &ItemCharges{Mode: "limited", Maximum: 10, PerUse: 0, Recharge: "regain-1d6-at-dawn"},
		Keys:    "item.attunement item.charges item.press-face spell.cast-face-spell spell.face-save-dc-17 item.face-charge-cost-table", Payload: "cube_of_force",
		Variants: []itemSourceManifestVariant{{"standard", "Cube of Force", "rare"}}},
	{Slug: "cubic-gate", Name: "Cubic Gate", Page: 216, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.charges action.magic item.six-sides-keyed-to-planes spell.cast-gate spell.cast-plane-shift", Payload: "cubic_gate",
		Variants: []itemSourceManifestVariant{{"standard", "Cubic Gate", "legendary"}}},
	{Slug: "dagger-of-venom", Name: "Dagger of Venom", Page: 216, Category: "weapon", Options: "dagger", Attunement: "none", Equip: "wield", Use: "bonus-action", UseRestrictions: "once-until-next-dawn",
		Keys: "item.magic-weapon attack.weapon-bonus damage.weapon-bonus action.bonus-action effect.poison-coating-one-minute save.constitution-dc-15 damage.poison-2d10 condition.poisoned-one-minute item.once-until-next-dawn", Payload: "venom",
		Variants: []itemSourceManifestVariant{{"standard", "Dagger of Venom", "rare"}}},
	{Slug: "dancing-sword", Name: "Dancing Sword", Page: 217, Category: "weapon", Options: "greatsword longsword rapier scimitar shortsword", Attunement: "required", Equip: "wield", Use: "bonus-action", UseRestrictions: "toss-into-air hover-flight",
		Keys: "item.attunement item.magic-weapon action.bonus-action effect.hover-fly-30-feet attack.uses-your-attack-roll damage.adds-ability-modifier attack.returns-after-fourth-attack effect.hover-ends-when-grasped-or-over-30-feet", Payload: "dancing",
		Variants: []itemSourceManifestVariant{{"standard", "Dancing Sword", "very-rare"}}},
	{Slug: "decanter-of-endless-water", Name: "Decanter of Endless Water", Page: 217, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "remove-stopper command-word",
		Keys: "action.magic item.command-words water.fresh-or-salt water.splash-1-gallon water.fountain-5-gallons water.geyser-30-gallons area.line-30-feet-by-1-foot save.strength-dc-13 damage.bludgeoning-1d4 condition.prone", Payload: "decanter",
		Variants: []itemSourceManifestVariant{{"standard", "Decanter of Endless Water", "uncommon"}}},
	{Slug: "deck-of-illusions", Name: "Deck of Illusions", Page: 217, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "draw-card-at-random throw-within-30-feet",
		Keys: "action.magic item.deck-34-cards item.cards-drawn-at-random summon.illusory-creature-from-table action.study check.intelligence-investigation-dc-15 effect.illusion-ends-when-card-moved-or-dispelled item.illusion-table", Payload: "illusions",
		Variants: []itemSourceManifestVariant{{"standard", "Deck of Illusions", "uncommon"}}},
	{Slug: "defender", Name: "Defender", Page: 218, Category: "weapon", Options: "any-melee-weapon", Attunement: "required", Equip: "wield", Use: "weapon-attack", UseRestrictions: "first-attack-each-turn",
		Keys: "item.attunement item.magic-weapon attack.weapon-bonus-3 damage.weapon-bonus-3 defense.transfer-bonus-to-armor-class effect.adjustment-lasts-until-next-turn", Payload: "defender",
		Variants: []itemSourceManifestVariant{{"standard", "Defender", "legendary"}}},
	{Slug: "demon-armor", Name: "Demon Armor", Page: 218, Category: "armor", Options: itemAnyArmor229, Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.armor-class-bonus-1 language.abyssal attack.unarmed-strike-slashing-1d8 attack.unarmed-bonus-1 item.curse item.curse-cannot-doff attack.disadvantage-against-demons save.disadvantage-against-demons", Payload: "demon_armor",
		Variants: []itemSourceManifestVariant{{"standard", "Demon Armor", "very-rare"}}},
	{Slug: "dimensional-shackles", Name: "Dimensional Shackles", Page: 218, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "utilize-action", UseRestrictions: "target-incapacitated-small-to-large",
		Keys: "action.utilize condition.incapacitated movement.prevents-extradimensional-movement check.strength-athletics-dc-30 effect.escape-once-every-30-days", Payload: "shackles",
		Variants: []itemSourceManifestVariant{{"standard", "Dimensional Shackles", "rare"}}},
	{Slug: "dragon-orb", Name: "Dragon Orb", Page: 218, Category: "wondrous-item", Attunement: "required", Equip: "none", Use: "magic-action", UseRestrictions: "attuned control-orb-for-spells-and-call",
		Charges: &ItemCharges{Mode: "limited", Maximum: 7, PerUse: 0, Recharge: "regain-1d4-plus-3-at-dawn"},
		Keys:    "item.attunement item.charges action.magic save.charisma-dc-15 condition.charmed spell.suggestion-save-dc-18 item.spell-charge-cost-table action.call-dragons defense.armor-class-20 item.artifact-destruction", Payload: "dragon_orb",
		Variants: []itemSourceManifestVariant{{"standard", "Dragon Orb", "artifact"}}},
}}

// LoadItems229 returns a fresh validated copy of batch #229. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems229() (ItemSourceBatch[ItemDetails229], error) { return DecodeItems229(items229JSON) }

// DecodeItems229 validates supplied JSON as batch #229 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems229(data []byte) (ItemSourceBatch[ItemDetails229], error) {
	return decodeItemSourceBatch(data, itemManifest229, validateItemDetails229)
}

// validateItemDetails229 checks the internal consistency of each payload: the
// printed tables must be complete and in order, and DCs, dice and counts must
// match the source values the family's mechanic keys name.
func validateItemDetails229(family, _ string, d *ItemDetails229) error {
	positive := func(values ...int) bool {
		for _, v := range values {
			if v <= 0 {
				return false
			}
		}
		return true
	}
	dice := func(x ItemDice, n, sides int) bool { return x.Dice == n && x.Sides == sides }
	ok := true
	switch family {
	case "crystal-ball-of-telepathy":
		c := d.Telepathy
		ok = c.SaveDC == 17 && c.ExtraSaveDC == 17 && positive(c.CommunicationRangeFeet, c.ExtraUses)
	case "crystal-ball-of-true-seeing":
		c := d.TrueSeeing
		ok = c.SaveDC == 17 && positive(c.SenseRangeFeet)
	case "cube-of-force":
		c := d.CubeForce
		// The Cube of Force Faces table, p. 216, in printed order.
		want := []ItemCubeFace{{"mage-armor", 1}, {"shield", 1}, {"tiny-hut", 3}, {"private-sanctum", 4}, {"resilient-sphere", 4}, {"wall-of-force", 5}}
		ok = c.SaveDC == 17 && c.StartingCharges == 10 && positive(c.WidthInches) && dice(c.RegainDice, 1, 6) && c.RegainDice.Type == "charges" && slices.Equal(c.Faces, want)
	case "cubic-gate":
		c := d.CubicGate
		ok = c.MaxCharges == 3 && c.ChargeCost == 1 && c.Sides == 6 && c.MaterialPlaneSides == 1 && positive(c.WidthInches) && dice(c.RegainDice, 1, 3) && c.RegainDice.Type == "charges"
	case "dagger-of-venom":
		c := d.Venom
		ok = positive(c.AttackBonus, c.DamageBonus, c.PoisonDurationSeconds, c.ConditionDurationSeconds, c.Uses) && c.SaveDC == 15 && dice(c.FailureDamage, 2, 10) && c.FailureDamage.Type == "poison" && uniqueItemKeys(c.PoisonEndsOn)
	case "dancing-sword":
		c := d.Dancing
		ok = positive(c.HoverFlyFeet, c.AttackRangeFeet, c.MoveFeet, c.MoveWithinFeetOfYou, c.HoverEndDistanceFeet) && c.AttacksBeforeReturn == 4 && uniqueItemKeys(c.HoverEndsOn)
	case "decanter-of-endless-water":
		c := d.Decanter
		want := []ItemDecanterCommand{{"splash", 1}, {"fountain", 5}, {"geyser", 30}}
		ok = positive(c.WeightPounds, c.GeyserLengthFeet, c.GeyserWidthFeet, c.ObjectMaxPounds) && c.GeyserSaveDC == 13 && dice(c.GeyserDamage, 1, 4) && c.GeyserDamage.Type == "bludgeoning" &&
			slices.Equal(c.Commands, want) && slices.Equal(c.WaterTypes, []string{"fresh", "salt"})
	case "deck-of-illusions":
		c := d.Illusions
		rows, creatures := [][2]int{}, []string{}
		for _, r := range c.Illusions {
			rows, creatures = append(rows, [2]int{r.Min, r.Max}), append(creatures, r.Creature)
		}
		// The 1d100 table, pp. 217–218: 32 creature rows of three, then the card drawer 97–00.
		ok = itemRollTable(rows, c.TableDieSides) && uniqueItemKeys(creatures) && len(creatures) == c.CreatureCards+1 && creatures[len(creatures)-1] == "card-drawer" &&
			c.FullDeckCards == c.CreatureCards+c.MirrorCards && c.FullDeckCards == 34 && c.TablePage == 218 && dice(c.MissingCardsDice, 1, 20) && c.MissingCardsDice.Type == "cards" &&
			c.InspectCheck.DC == 15 && positive(c.MissingCardsSubtract, c.ThrowRangeFeet, c.MoveWhileWithinFeet, c.MoveDistanceFeet) && uniqueItemKeys(c.LastsUntil) && uniqueItemKeys(c.DispelledBy)
	case "defender":
		c := d.Defender
		ok = c.AttackBonus == 3 && c.DamageBonus == 3 && positive(c.ExampleAttackDamageBonus, c.ExampleArmorClassBonus) && c.ExampleAttackDamageBonus+c.ExampleArmorClassBonus == c.AttackBonus && uniqueItemKeys(c.TransferFrom)
	case "demon-armor":
		c := d.DemonArmor
		ok = positive(c.ArmorClassBonus, c.UnarmedAttackBonus, c.UnarmedDamageBonus) && dice(c.UnarmedDamage, 1, 8) && c.UnarmedDamage.Type == "slashing" && uniqueItemKeys(c.DoffAllowedIf) && uniqueItemKeys(c.CurseDisadvantage)
	case "dimensional-shackles":
		c := d.Shackles
		ok = c.EscapeCheck.DC == 30 && positive(c.EscapeIntervalDays) && uniqueItemKeys(c.PreventedExamples) && uniqueItemKeys(c.RemovableBy)
	case "dragon-orb":
		c := d.DragonOrb
		// The orb's spell table, p. 218, in printed order.
		want := []ItemOrbSpell{{"cure-wounds", 9, 0, 4}, {"daylight", 0, 0, 1}, {"death-ward", 0, 0, 2}, {"detect-magic", 0, 0, 0}, {"scrying", 0, 18, 3}}
		ok = c.SaveDC == 15 && c.CharmedSpellSaveDC == 18 && c.MaxCharges == 7 && dice(c.RegainDice, 1, 4) && c.RegainDice.Type == "charges" && c.RegainBonus == 3 &&
			slices.Equal(c.Spells, want) && positive(c.DiameterInches, c.UsedDiameterInches, c.CallRangeMiles, c.CallCooldownSeconds, c.ArmorClass) && uniqueItemKeys(c.DesireExamples) && uniqueItemKeys(c.DestroyedBy)
	}
	if !ok {
		return fmt.Errorf("incomplete or inconsistent source values")
	}
	return nil
}
