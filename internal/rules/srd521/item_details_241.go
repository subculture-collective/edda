package srd521

import (
	_ "embed"
	"fmt"
	"slices"
)

//go:embed testdata/items_241.json
var items241JSON []byte

// ItemDetails241 holds exactly one typed payload per variant of catalogue batch
// #241 (magic items 17, SRD 5.2.1 pp. 241-243). Payloads are source facts only.
type ItemDetails241 struct {
	LordlyMight       *ItemRodOfLordlyMight          `json:"lordly_might,omitempty"`
	Resurrection      *ItemRodOfResurrection         `json:"resurrection,omitempty"`
	Rulership         *ItemRodOfRulership            `json:"rulership,omitempty"`
	Security          *ItemRodOfSecurity             `json:"security,omitempty"`
	ClimbingRope      *ItemRopeOfClimbing            `json:"climbing_rope,omitempty"`
	Entanglement      *ItemRopeOfEntanglement        `json:"entanglement,omitempty"`
	Scarab            *ItemScarabOfProtection        `json:"scarab,omitempty"`
	Speed             *ItemScimitarOfSpeed           `json:"speed,omitempty"`
	SendingStones     *ItemSendingStones             `json:"sending_stones,omitempty"`
	Sentinel          *ItemSentinelShield            `json:"sentinel,omitempty"`
	MissileAttraction *ItemShieldOfMissileAttraction `json:"missile_attraction,omitempty"`
	Cavalier          *ItemShieldOfTheCavalier       `json:"cavalier,omitempty"`
}

// ItemRodOfLordlyMight retains typed SRD source facts without implementing them.
type ItemRodOfLordlyMight struct {
	FunctionsAs string                        `json:"functions_as"`
	AttackBonus int                           `json:"attack_bonus"`
	DamageBonus int                           `json:"damage_bonus"`
	ButtonCount int                           `json:"button_count"`
	Buttons     ItemRodOfLordlyMightButtons   `json:"buttons"`
	DrainLife   ItemRodOfLordlyMightDrainLife `json:"drain_life"`
	Paralyze    ItemRodOfLordlyMightParalyze  `json:"paralyze"`
	Terrify     ItemRodOfLordlyMightTerrify   `json:"terrify"`
	SourcePages []int                         `json:"source_pages"`
}

// ItemRodOfLordlyMightButtons is a nested part of a source payload.
type ItemRodOfLordlyMightButtons struct {
	PressAction      string                                  `json:"press_action"`
	EffectLastsUntil []string                                `json:"effect_lasts_until"`
	FieryBlade       ItemRodOfLordlyMightButtonsFieryBlade   `json:"fiery_blade"`
	Battleaxe        ItemRodOfLordlyMightButtonsBattleaxe    `json:"battleaxe"`
	Spear            ItemRodOfLordlyMightButtonsSpear        `json:"spear"`
	ClimbingPole     ItemRodOfLordlyMightButtonsClimbingPole `json:"climbing_pole"`
	BatteringRam     ItemRodOfLordlyMightButtonsBatteringRam `json:"battering_ram"`
	NormalForm       ItemRodOfLordlyMightButtonsNormalForm   `json:"normal_form"`
}

// ItemRodOfLordlyMightButtonsFieryBlade is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsFieryBlade struct {
	Button                 int      `json:"button"`
	BrightLightRadiusFeet  int      `json:"bright_light_radius_feet"`
	DimLightAdditionalFeet int      `json:"dim_light_additional_feet"`
	FunctionsAs            []string `json:"functions_as"`
	ExtraDamage            ItemDice `json:"extra_damage"`
}

// ItemRodOfLordlyMightButtonsBattleaxe is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsBattleaxe struct {
	Button               int    `json:"button"`
	FunctionsAs          string `json:"functions_as"`
	AttackAndDamageBonus int    `json:"attack_and_damage_bonus"`
}

// ItemRodOfLordlyMightButtonsSpear is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsSpear struct {
	Button               int    `json:"button"`
	FunctionsAs          string `json:"functions_as"`
	HaftFeet             int    `json:"haft_feet"`
	AttackAndDamageBonus int    `json:"attack_and_damage_bonus"`
}

// ItemRodOfLordlyMightButtonsClimbingPole is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsClimbingPole struct {
	Button             int      `json:"button"`
	MaximumLengthFeet  int      `json:"maximum_length_feet"`
	ButtonsStayInReach bool     `json:"buttons_stay_in_reach"`
	Anchors            []string `json:"anchors"`
	AnchorSurface      string   `json:"anchor_surface"`
	BarLengthInches    int      `json:"bar_length_inches"`
	BarSpacingFeet     int      `json:"bar_spacing_feet"`
	CapacityPounds     int      `json:"capacity_pounds"`
	RevertsOn          []string `json:"reverts_on"`
}

// ItemRodOfLordlyMightButtonsBatteringRam is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsBatteringRam struct {
	Button        int      `json:"button"`
	CheckBonus    int      `json:"check_bonus"`
	CheckAbility  string   `json:"check_ability"`
	CheckSkill    string   `json:"check_skill"`
	CheckSubjects []string `json:"check_subjects"`
}

// ItemRodOfLordlyMightButtonsNormalForm is a nested part of a source payload.
type ItemRodOfLordlyMightButtonsNormalForm struct {
	Button               int      `json:"button"`
	Indicates            string   `json:"indicates"`
	WithoutMagneticNorth string   `json:"without_magnetic_north"`
	Knowledge            []string `json:"knowledge"`
}

// ItemRodOfLordlyMightDrainLife is a nested part of a source payload.
type ItemRodOfLordlyMightDrainLife struct {
	Trigger     string   `json:"trigger"`
	SaveAbility string   `json:"save_ability"`
	SaveDC      int      `json:"save_dc"`
	ExtraDamage ItemDice `json:"extra_damage"`
	YouRegain   string   `json:"you_regain"`
	ReuseAfter  string   `json:"reuse_after"`
}

// ItemRodOfLordlyMightParalyze is a nested part of a source payload.
type ItemRodOfLordlyMightParalyze struct {
	Trigger         string `json:"trigger"`
	SaveAbility     string `json:"save_ability"`
	SaveDC          int    `json:"save_dc"`
	Condition       string `json:"condition"`
	DurationMinutes int    `json:"duration_minutes"`
	RepeatSave      string `json:"repeat_save"`
	ReuseAfter      string `json:"reuse_after"`
}

// ItemRodOfLordlyMightTerrify is a nested part of a source payload.
type ItemRodOfLordlyMightTerrify struct {
	Condition       string `json:"condition"`
	Action          string `json:"action"`
	RangeFeet       int    `json:"range_feet"`
	SaveAbility     string `json:"save_ability"`
	SaveDC          int    `json:"save_dc"`
	DurationMinutes int    `json:"duration_minutes"`
	RepeatSave      string `json:"repeat_save"`
	RepeatSuccess   string `json:"repeat_success"`
	ReuseAfter      string `json:"reuse_after"`
}

// ItemRodOfResurrection retains typed SRD source facts without implementing them.
type ItemRodOfResurrection struct {
	MaxCharges      int                             `json:"max_charges"`
	Condition       string                          `json:"condition"`
	Spells          []ItemRodOfResurrectionSpell    `json:"spells"`
	RechargeCharges int                             `json:"recharge_charges"`
	RechargeWhen    string                          `json:"recharge_when"`
	LastCharge      ItemRodOfResurrectionLastCharge `json:"last_charge"`
}

// ItemRodOfResurrectionSpell is a nested part of a source payload.
type ItemRodOfResurrectionSpell struct {
	Spell      string `json:"spell"`
	ChargeCost int    `json:"charge_cost"`
}

// ItemRodOfResurrectionLastCharge is a nested part of a source payload.
type ItemRodOfResurrectionLastCharge struct {
	DieSides    int    `json:"die_sides"`
	DestroyedOn int    `json:"destroyed_on"`
	Result      string `json:"result"`
}

// ItemRodOfRulership retains typed SRD source facts without implementing them.
type ItemRodOfRulership struct {
	Action        string   `json:"action"`
	Targets       string   `json:"targets"`
	RangeFeet     int      `json:"range_feet"`
	SaveAbility   string   `json:"save_ability"`
	SaveDC        int      `json:"save_dc"`
	Condition     string   `json:"condition"`
	DurationHours int      `json:"duration_hours"`
	RegardsYouAs  string   `json:"regards_you_as"`
	EndsEarlyWhen []string `json:"ends_early_when"`
	ReuseAfter    string   `json:"reuse_after"`
}

// ItemRodOfSecurity retains typed SRD source facts without implementing them.
type ItemRodOfSecurity struct {
	Condition      string                     `json:"condition"`
	Action         string                     `json:"action"`
	Travelers      ItemRodOfSecurityTravelers `json:"travelers"`
	Destination    string                     `json:"destination"`
	FormChosenBy   string                     `json:"form_chosen_by"`
	Demiplane      ItemRodOfSecurityDemiplane `json:"demiplane"`
	Healing        ItemRodOfSecurityHealing   `json:"healing"`
	Stay           ItemRodOfSecurityStay      `json:"stay"`
	EndsWhen       []string                   `json:"ends_when"`
	ReturnTo       string                     `json:"return_to"`
	ReuseAfterDays int                        `json:"reuse_after_days"`
}

// ItemRodOfSecurityTravelers is a nested part of a source payload.
type ItemRodOfSecurityTravelers struct {
	IncludesYou           bool `json:"includes_you"`
	OtherWillingCreatures int  `json:"other_willing_creatures"`
}

// ItemRodOfSecurityDemiplane is a nested part of a source payload.
type ItemRodOfSecurityDemiplane struct {
	Provides              string `json:"provides"`
	EnvironmentCanHarm    bool   `json:"environment_can_harm"`
	ObjectsExistOnlyThere bool   `json:"objects_exist_only_there"`
	Aging                 bool   `json:"aging"`
	TimePassesNormally    bool   `json:"time_passes_normally"`
}

// ItemRodOfSecurityHealing is a nested part of a source payload.
type ItemRodOfSecurityHealing struct {
	Per                   string `json:"per"`
	AsIfHitPointDiceSpent int    `json:"as_if_hit_point_dice_spent"`
}

// ItemRodOfSecurityStay is a nested part of a source payload.
type ItemRodOfSecurityStay struct {
	BaseDays  int    `json:"base_days"`
	DividedBy string `json:"divided_by"`
	Rounding  string `json:"rounding"`
}

// ItemRopeOfClimbing retains typed SRD source facts without implementing them.
type ItemRopeOfClimbing struct {
	LengthFeet            int                       `json:"length_feet"`
	CapacityPounds        int                       `json:"capacity_pounds"`
	Condition             string                    `json:"condition"`
	Action                string                    `json:"action"`
	MoveLimit             string                    `json:"move_limit"`
	MoveFeetPerTurn       int                       `json:"move_feet_per_turn"`
	MovesUntil            []string                  `json:"moves_until"`
	OtherCommands         []string                  `json:"other_commands"`
	Knotted               ItemRopeOfClimbingKnotted `json:"knotted"`
	ArmorClass            int                       `json:"armor_class"`
	HitPoints             int                       `json:"hit_points"`
	Immunities            []string                  `json:"immunities"`
	RegainsHitPoints      int                       `json:"regains_hit_points"`
	RegainIntervalMinutes int                       `json:"regain_interval_minutes"`
	RegainNeedsHitPoints  int                       `json:"regain_needs_hit_points"`
	DestroyedAtHitPoints  int                       `json:"destroyed_at_hit_points"`
}

// ItemRopeOfClimbingKnotted is a nested part of a source payload.
type ItemRopeOfClimbingKnotted struct {
	KnotSpacingFeet int    `json:"knot_spacing_feet"`
	LengthFeet      int    `json:"length_feet"`
	CheckMode       string `json:"check_mode"`
	CheckSubject    string `json:"check_subject"`
}

// ItemRopeOfEntanglement retains typed SRD source facts without implementing them.
type ItemRopeOfEntanglement struct {
	LengthFeet            int                             `json:"length_feet"`
	Condition             string                          `json:"condition"`
	Action                string                          `json:"action"`
	Targets               int                             `json:"targets"`
	RangeFeet             int                             `json:"range_feet"`
	SaveAbility           string                          `json:"save_ability"`
	SaveDC                int                             `json:"save_dc"`
	AppliesCondition      string                          `json:"applies_condition"`
	Release               []ItemRopeOfEntanglementRelease `json:"release"`
	Escape                ItemRopeOfEntanglementEscape    `json:"escape"`
	Recoil                ItemRopeOfEntanglementRecoil    `json:"recoil"`
	ArmorClass            int                             `json:"armor_class"`
	HitPoints             int                             `json:"hit_points"`
	Immunities            []string                        `json:"immunities"`
	RegainsHitPoints      int                             `json:"regains_hit_points"`
	RegainIntervalMinutes int                             `json:"regain_interval_minutes"`
	RegainNeedsHitPoints  int                             `json:"regain_needs_hit_points"`
	DestroyedAtHitPoints  int                             `json:"destroyed_at_hit_points"`
	SourcePages           []int                           `json:"source_pages"`
}

// ItemRopeOfEntanglementRelease is a nested part of a source payload.
type ItemRopeOfEntanglementRelease struct {
	Method      string `json:"method"`
	RopeCoilsIn string `json:"rope_coils_in"`
}

// ItemRopeOfEntanglementEscape is a nested part of a source payload.
type ItemRopeOfEntanglementEscape struct {
	Action  string                               `json:"action"`
	Choices []ItemRopeOfEntanglementEscapeChoice `json:"choices"`
	DC      int                                  `json:"dc"`
	Success string                               `json:"success"`
}

// ItemRopeOfEntanglementEscapeChoice is a nested part of a source payload.
type ItemRopeOfEntanglementEscapeChoice struct {
	Ability string `json:"ability"`
	Skill   string `json:"skill"`
}

// ItemRopeOfEntanglementRecoil is a nested part of a source payload.
type ItemRopeOfEntanglementRecoil struct {
	Trigger   string `json:"trigger"`
	Action    string `json:"action"`
	Otherwise string `json:"otherwise"`
}

// ItemScarabOfProtection retains typed SRD source facts without implementing them.
type ItemScarabOfProtection struct {
	Appearance      string                                `json:"appearance"`
	Condition       string                                `json:"condition"`
	BenefitCount    int                                   `json:"benefit_count"`
	ArmorClassBonus int                                   `json:"armor_class_bonus"`
	Preservation    ItemScarabOfProtectionPreservation    `json:"preservation"`
	SpellResistance ItemScarabOfProtectionSpellResistance `json:"spell_resistance"`
}

// ItemScarabOfProtectionPreservation is a nested part of a source payload.
type ItemScarabOfProtectionPreservation struct {
	MaxCharges      int      `json:"max_charges"`
	Triggers        []string `json:"triggers"`
	Action          string   `json:"action"`
	ChargeCost      int      `json:"charge_cost"`
	Result          string   `json:"result"`
	AfterLastCharge string   `json:"after_last_charge"`
}

// ItemScarabOfProtectionSpellResistance is a nested part of a source payload.
type ItemScarabOfProtectionSpellResistance struct {
	SaveMode string `json:"save_mode"`
	Against  string `json:"against"`
}

// ItemScimitarOfSpeed retains typed SRD source facts without implementing them.
type ItemScimitarOfSpeed struct {
	AttackBonus       int    `json:"attack_bonus"`
	DamageBonus       int    `json:"damage_bonus"`
	ExtraAttackAction string `json:"extra_attack_action"`
	ExtraAttacks      int    `json:"extra_attacks"`
	Frequency         string `json:"frequency"`
}

// ItemSendingStones retains typed SRD source facts without implementing them.
type ItemSendingStones struct {
	StonesPerPair              string `json:"stones_per_pair"`
	Matching                   string `json:"matching"`
	Condition                  string `json:"condition"`
	Spell                      string `json:"spell"`
	Target                     string `json:"target"`
	IfNoBearer                 string `json:"if_no_bearer"`
	SharedReuseAfter           string `json:"shared_reuse_after"`
	OtherStoneWhenOneDestroyed string `json:"other_stone_when_one_destroyed"`
}

// ItemSentinelShield retains typed SRD source facts without implementing them.
type ItemSentinelShield struct {
	Condition   string   `json:"condition"`
	AdvantageOn []string `json:"advantage_on"`
	Emblem      string   `json:"emblem"`
}

// ItemShieldOfMissileAttraction retains typed SRD source facts without implementing them.
type ItemShieldOfMissileAttraction struct {
	Condition      string                             `json:"condition"`
	Defense        string                             `json:"defense"`
	DefenseAgainst string                             `json:"defense_against"`
	Curse          ItemShieldOfMissileAttractionCurse `json:"curse"`
}

// ItemShieldOfMissileAttractionCurse is a nested part of a source payload.
type ItemShieldOfMissileAttractionCurse struct {
	BeginsOn                string `json:"begins_on"`
	EndsOn                  string `json:"ends_on"`
	RemovingShieldEndsCurse bool   `json:"removing_shield_ends_curse"`
	Trigger                 string `json:"trigger"`
	RadiusFeet              int    `json:"radius_feet"`
	Result                  string `json:"result"`
}

// ItemShieldOfTheCavalier retains typed SRD source facts without implementing them.
type ItemShieldOfTheCavalier struct {
	Condition             string                                 `json:"condition"`
	ArmorClassBonus       int                                    `json:"armor_class_bonus"`
	StacksWithShieldBonus bool                                   `json:"stacks_with_shield_bonus"`
	ForcefulBash          ItemShieldOfTheCavalierForcefulBash    `json:"forceful_bash"`
	ProtectiveField       ItemShieldOfTheCavalierProtectiveField `json:"protective_field"`
}

// ItemShieldOfTheCavalierForcefulBash is a nested part of a source payload.
type ItemShieldOfTheCavalierForcefulBash struct {
	Trigger                   string   `json:"trigger"`
	AttacksReplaced           int      `json:"attacks_replaced"`
	ReachFeet                 int      `json:"reach_feet"`
	AttackRollAdds            []string `json:"attack_roll_adds"`
	Damage                    ItemDice `json:"damage"`
	DamageBonus               int      `json:"damage_bonus"`
	DamageAddsAbilityModifier string   `json:"damage_adds_ability_modifier"`
	PushFeet                  int      `json:"push_feet"`
	ProneIf                   string   `json:"prone_if"`
	AppliesCondition          string   `json:"applies_condition"`
}

// ItemShieldOfTheCavalierProtectiveField is a nested part of a source payload.
type ItemShieldOfTheCavalierProtectiveField struct {
	Action                string   `json:"action"`
	Triggers              []string `json:"triggers"`
	AllyRangeFeet         int      `json:"ally_range_feet"`
	EmanationFeet         int      `json:"emanation_feet"`
	Immobile              bool     `json:"immobile"`
	NotFullyContained     string   `json:"not_fully_contained"`
	TriggeringEffect      string   `json:"triggering_effect"`
	ConcentrationRequired bool     `json:"concentration_required"`
	MaxDurationMinutes    int      `json:"max_duration_minutes"`
	Passage               string   `json:"passage"`
	Isolation             []string `json:"isolation"`
	ReuseAfter            string   `json:"reuse_after"`
}

// itemManifest241 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest241 = itemSourceManifest{Batch: "data-magic-items-17", Issue: 241, Families: []itemSourceManifestFamily{
	{Slug: "rod-of-lordly-might", Name: "Rod of Lordly Might", Page: 241, Category: "rod", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon bonus.attack-and-damage-3 item.six-buttons item.weapon-forms item.climbing-pole item.battering-ram item.magnetic-north save.constitution-dc-17 damage.drain-life-4d6-necrotic condition.paralyzed save.wisdom-dc-17 condition.frightened item.once-until-dawn", Payload: "lordly_might",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Lordly Might", "legendary"}}},
	{Slug: "rod-of-resurrection", Name: "Rod of Resurrection", Page: 242, Category: "rod", Attunement: "required", Equip: "hold", Use: "cast-spell", UseRestrictions: "expend-listed-charges",
		Charges: &ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1-at-dawn"},
		Keys:    "item.attunement item.hold item.charges spell.cast-heal spell.cast-resurrection item.destroyed-on-1-after-last-charge", Payload: "resurrection",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Resurrection", "legendary"}}},
	{Slug: "rod-of-rulership", Name: "Rod of Rulership", Page: 242, Category: "rod", Attunement: "required", Equip: "hold", Use: "magic-action",
		Keys: "item.attunement action.magic save.wisdom-dc-15 condition.charmed duration.8-hours item.once-until-dawn", Payload: "rulership",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Rulership", "rare"}}},
	{Slug: "rod-of-security", Name: "Rod of Security", Page: 242, Category: "rod", Attunement: "none", Equip: "hold", Use: "magic-action",
		Keys: "item.hold action.magic movement.transport-to-demiplane healing.hit-point-die-per-hour effect.no-aging duration.200-days-divided-by-creatures item.10-day-cooldown", Payload: "security",
		Variants: []itemSourceManifestVariant{{"standard", "Rod of Security", "very-rare"}}},
	{Slug: "rope-of-climbing", Name: "Rope of Climbing", Page: 242, Category: "wondrous-item", Attunement: "none", Equip: "hold", Use: "magic-action", UseRestrictions: "holding-one-end",
		Keys: "action.magic item.animate-rope item.fasten-knot-and-coil check.climb-advantage-when-knotted item.armor-class-20-hit-points-20 defense.immunity-poison-and-psychic healing.1-hit-point-per-5-minutes", Payload: "climbing_rope",
		Variants: []itemSourceManifestVariant{{"standard", "Rope of Climbing", "uncommon"}}},
	{Slug: "rope-of-entanglement", Name: "Rope of Entanglement", Page: 242, Category: "wondrous-item", Attunement: "none", Equip: "hold", Use: "magic-action", UseRestrictions: "holding-one-end",
		Keys: "action.magic save.dexterity-dc-15 condition.restrained check.escape-dc-15 action.bonus-release action.reaction-recoil item.armor-class-20-hit-points-20 defense.immunity-poison-and-psychic healing.1-hit-point-per-5-minutes", Payload: "entanglement",
		Variants: []itemSourceManifestVariant{{"standard", "Rope of Entanglement", "rare"}}},
	{Slug: "scarab-of-protection", Name: "Scarab of Protection", Page: 243, Category: "wondrous-item", Attunement: "required", Equip: "carry", Use: "reaction", UseRestrictions: "failed-save-against-necromancy-or-undead expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 12, PerUse: 1, Recharge: "none"},
		Keys:    "item.attunement item.carry defense.armor-class-bonus-1 item.charges action.reaction save.failure-becomes-success-against-necromancy-or-undead item.destroyed-at-last-charge save.advantage-against-spells", Payload: "scarab",
		Variants: []itemSourceManifestVariant{{"standard", "Scarab of Protection", "legendary"}}},
	{Slug: "scimitar-of-speed", Name: "Scimitar of Speed", Page: 243, Category: "weapon", Options: "scimitar", Attunement: "required", Equip: "wield", Use: "weapon-attack",
		Keys: "item.attunement item.magic-weapon bonus.attack-and-damage-2 attack.bonus-action-attack", Payload: "speed",
		Variants: []itemSourceManifestVariant{{"standard", "Scimitar of Speed", "very-rare"}}},
	{Slug: "sending-stones", Name: "Sending Stones", Page: 243, Category: "wondrous-item", Attunement: "none", Equip: "touch", Use: "cast-spell", UseRestrictions: "sending",
		Keys: "item.paired spell.cast-sending spell.target-other-stone-bearer item.once-until-dawn-per-pair item.nonmagical-if-pair-destroyed", Payload: "sending_stones",
		Variants: []itemSourceManifestVariant{{"standard", "Sending Stones", "uncommon"}}},
	{Slug: "sentinel-shield", Name: "Sentinel Shield", Page: 243, Category: "armor", Options: "shield", Attunement: "none", Equip: "hold", Use: "none",
		Keys: "item.hold check.initiative-advantage check.perception-advantage", Payload: "sentinel",
		Variants: []itemSourceManifestVariant{{"standard", "Sentinel Shield", "uncommon"}}},
	{Slug: "shield-of-missile-attraction", Name: "Shield of Missile Attraction", Page: 243, Category: "armor", Options: "shield", Attunement: "required", Equip: "hold", Use: "none",
		Keys: "item.attunement item.hold defense.resistance-ranged-weapon-damage item.curse item.curse-ends-on-remove-curse item.curse-redirects-ranged-attacks-within-10-feet", Payload: "missile_attraction",
		Variants: []itemSourceManifestVariant{{"standard", "Shield of Missile Attraction", "rare"}}},
	{Slug: "shield-of-the-cavalier", Name: "Shield of the Cavalier", Page: 243, Category: "armor", Options: "shield", Attunement: "required", Equip: "hold", Use: "reaction", UseRestrictions: "protective-field-trigger",
		Keys: "item.attunement item.hold defense.armor-class-bonus-2 attack.forceful-bash damage.force-2d6-plus-2 movement.push-10-feet condition.prone action.reaction-protective-field spell.concentration-up-to-1-minute item.once-until-dawn", Payload: "cavalier",
		Variants: []itemSourceManifestVariant{{"standard", "Shield of the Cavalier", "very-rare"}}},
}}

// LoadItems241 returns a fresh validated copy of batch #241. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems241() (ItemSourceBatch[ItemDetails241], error) { return DecodeItems241(items241JSON) }

// DecodeItems241 validates supplied JSON as batch #241 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems241(data []byte) (ItemSourceBatch[ItemDetails241], error) {
	return decodeItemSourceBatch(data, itemManifest241, validateItemDetails241)
}

// validateItemDetails241 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails241(family, variant string, d *ItemDetails241) error {
	ok := true
	switch family {
	case "rod-of-lordly-might":
		p := d.LordlyMight
		b := p.Buttons
		ok = p.AttackBonus == 3 && p.DamageBonus == 3 && p.ButtonCount == 6 && b.FieryBlade.Button == 1 && b.Battleaxe.Button == 2 && b.Spear.Button == 3 && b.ClimbingPole.Button == 4 && b.BatteringRam.Button == 5 && b.NormalForm.Button == 6 &&
			p.DrainLife.SaveDC == 17 && p.Paralyze.SaveDC == 17 && p.Terrify.SaveDC == 17 && p.DrainLife.ExtraDamage == ItemDice{Dice: 4, Sides: 6, Type: "necrotic"} && slices.Equal(p.SourcePages, []int{241, 242})
	case "rod-of-resurrection":
		p := d.Resurrection
		ok = p.MaxCharges == 5 && p.RechargeCharges == 1 && len(p.Spells) == 2 && p.Spells[0].ChargeCost == 1 && p.Spells[1].ChargeCost == 5 && p.LastCharge.DieSides == 20 && p.LastCharge.DestroyedOn == 1
	case "rod-of-rulership":
		p := d.Rulership
		ok = p.RangeFeet == 120 && p.SaveDC == 15 && p.DurationHours == 8 && p.Condition == "charmed"
	case "rod-of-security":
		p := d.Security
		ok = p.Travelers.OtherWillingCreatures == 199 && p.Stay.BaseDays == 200 && p.Healing.AsIfHitPointDiceSpent == 1 && p.ReuseAfterDays == 10
	case "rope-of-climbing":
		p := d.ClimbingRope
		ok = p.LengthFeet == 60 && p.CapacityPounds == 3000 && p.Knotted.LengthFeet == 50 && p.ArmorClass == 20 && p.HitPoints == 20 && uniqueItemKeys(p.Immunities) && uniqueItemKeys(p.OtherCommands)
	case "rope-of-entanglement":
		p := d.Entanglement
		ok = p.LengthFeet == 30 && p.RangeFeet == 20 && p.SaveDC == 15 && p.Escape.DC == 15 && p.ArmorClass == 20 && p.HitPoints == 20 && len(p.Release) == 2 && len(p.Escape.Choices) == 2 && uniqueItemKeys(p.Immunities) && slices.Equal(p.SourcePages, []int{242, 243})
	case "scarab-of-protection":
		p := d.Scarab
		ok = p.BenefitCount == 3 && p.ArmorClassBonus == 1 && p.Preservation.MaxCharges == 12 && p.Preservation.ChargeCost == 1 && uniqueItemKeys(p.Preservation.Triggers)
	case "scimitar-of-speed":
		p := d.Speed
		ok = p.AttackBonus == 2 && p.DamageBonus == 2 && p.ExtraAttacks == 1 && p.ExtraAttackAction == "bonus"
	case "sending-stones":
		p := d.SendingStones
		ok = p.Spell == "sending" && p.Target == "bearer-of-other-stone"
	case "sentinel-shield":
		p := d.Sentinel
		ok = uniqueItemKeys(p.AdvantageOn) && len(p.AdvantageOn) == 2
	case "shield-of-missile-attraction":
		p := d.MissileAttraction
		ok = p.Defense == "resistance" && p.Curse.RadiusFeet == 10 && !p.Curse.RemovingShieldEndsCurse
	case "shield-of-the-cavalier":
		p := d.Cavalier
		b, f := p.ForcefulBash, p.ProtectiveField
		ok = p.ArmorClassBonus == 2 && p.StacksWithShieldBonus && b.Damage == ItemDice{Dice: 2, Sides: 6, Type: "force"} && b.DamageBonus == 2 && b.PushFeet == 10 && b.ReachFeet == 5 && f.EmanationFeet == 5 && f.MaxDurationMinutes == 1 && uniqueItemKeys(f.Triggers) && uniqueItemKeys(f.Isolation)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
