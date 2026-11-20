package srd521

import (
	_ "embed"
	"fmt"
)

//go:embed testdata/items_228.json
var items228JSON []byte

// ItemDetails228 holds exactly one typed payload per variant of catalogue batch
// #228 (magic items 04, SRD 5.2.1 pp. 215–216). Payloads are source facts only.
type ItemDetails228 struct {
	Censer       *ItemCenserAirElementals   `json:"censer,omitempty"`
	Chime        *ItemChimeOpening          `json:"chime,omitempty"`
	Circlet      *ItemCircletBlasting       `json:"circlet,omitempty"`
	Arachnida    *ItemCloakArachnida        `json:"arachnida,omitempty"`
	Displacement *ItemCloakDisplacement     `json:"displacement,omitempty"`
	Elvenkind    *ItemCloakElvenkind        `json:"elvenkind,omitempty"`
	Invisibility *ItemCloakInvisibility     `json:"invisibility,omitempty"`
	Protection   *ItemCloakProtection       `json:"protection,omitempty"`
	Bat          *ItemCloakBat              `json:"bat,omitempty"`
	MantaRay     *ItemCloakMantaRay         `json:"manta_ray,omitempty"`
	CrystalBall  *ItemCrystalBall           `json:"crystal_ball,omitempty"`
	MindReading  *ItemCrystalBallMindReader `json:"mind_reading,omitempty"`
}

// ItemCenserAirElementals retains typed SRD source facts without implementing them.
type ItemCenserAirElementals struct {
	Condition        string   `json:"condition"`
	Action           string   `json:"action"`
	Summons          string   `json:"summons"`
	Location         string   `json:"location"`
	UnderstandsOwner bool     `json:"understands_owner_languages"`
	ObeysOwner       bool     `json:"obeys_owner_commands"`
	Initiative       string   `json:"initiative"`
	DurationSeconds  int      `json:"duration_seconds"`
	EndTriggers      []string `json:"end_triggers"`
	Uses             int      `json:"uses"`
	Recharge         string   `json:"recharge"`
}

// ItemChimeOpening retains typed SRD source facts without implementing them.
// The ten uses are a use limit, not a charge pool: the source never says charges.
type ItemChimeOpening struct {
	Appearance    string `json:"appearance"`
	LengthFeet    int    `json:"length_feet"`
	WeightPounds  int    `json:"weight_pounds"`
	Action        string `json:"action"`
	Spell         string `json:"spell"`
	SoundReplaces string `json:"sound_replaces"`
	AudibleFeet   int    `json:"audible_feet"`
	Uses          int    `json:"uses"`
	AfterFinalUse string `json:"after_final_use"`
}

// ItemCircletBlasting retains typed SRD source facts without implementing them.
type ItemCircletBlasting struct {
	Condition   string `json:"condition"`
	Spell       string `json:"spell"`
	AttackBonus int    `json:"attack_bonus"`
	Uses        int    `json:"uses"`
	Recharge    string `json:"recharge"`
}

// ItemCloakArachnida retains typed SRD source facts without implementing them.
type ItemCloakArachnida struct {
	Condition         string   `json:"condition"`
	Resistances       []string `json:"resistances"`
	ClimbSpeed        string   `json:"climb_speed"`
	ClimbDirections   []string `json:"climb_directions"`
	ClimbHandsFree    bool     `json:"climb_hands_free"`
	WebCatch          string   `json:"web_catch"`
	WebMovement       string   `json:"web_movement"`
	WebSpell          string   `json:"web_spell"`
	WebSaveDC         int      `json:"web_save_dc"`
	WebAreaMultiplier int      `json:"web_area_multiplier"`
	WebUses           int      `json:"web_uses"`
	WebRecharge       string   `json:"web_recharge"`
}

// ItemCloakDisplacement retains typed SRD source facts without implementing them.
type ItemCloakDisplacement struct {
	Condition      string `json:"condition"`
	Effect         string `json:"effect"`
	AttackRolls    string `json:"attack_rolls_against_wearer"`
	SuspendedBy    string `json:"suspended_by"`
	ResumesAt      string `json:"resumes_at"`
	SuppressedWhen string `json:"suppressed_while"`
}

// ItemCloakElvenkind retains typed SRD source facts without implementing them.
type ItemCloakElvenkind struct {
	Condition       string `json:"condition"`
	PerceiveAbility string `json:"perceive_wearer_ability"`
	PerceiveSkill   string `json:"perceive_wearer_skill"`
	PerceiveMode    string `json:"perceive_wearer_mode"`
	StealthAbility  string `json:"stealth_ability"`
	StealthSkill    string `json:"stealth_skill"`
	StealthMode     string `json:"stealth_mode"`
}

// ItemCloakInvisibility retains typed SRD source facts without implementing
// them. The charge pool is also recorded on the variant's charges field.
type ItemCloakInvisibility struct {
	Condition       string   `json:"condition"`
	Action          string   `json:"action"`
	Gesture         string   `json:"gesture"`
	ChargeCost      int      `json:"charge_cost"`
	GrantsCondition string   `json:"grants_condition"`
	DurationSeconds int      `json:"duration_seconds"`
	EndsEarly       []string `json:"ends_early"`
	HoodDownAction  string   `json:"hood_down_action"`
	MaxCharges      int      `json:"max_charges"`
	RegainDice      ItemDice `json:"regain_dice"`
	RegainTiming    string   `json:"regain_timing"`
}

// ItemCloakProtection retains typed SRD source facts without implementing them.
type ItemCloakProtection struct {
	Condition        string `json:"condition"`
	ACBonus          int    `json:"ac_bonus"`
	SavingThrowBonus int    `json:"saving_throw_bonus"`
}

// ItemCloakBat retains typed SRD source facts without implementing them.
type ItemCloakBat struct {
	Condition         string   `json:"condition"`
	StealthAbility    string   `json:"stealth_ability"`
	StealthSkill      string   `json:"stealth_skill"`
	StealthMode       string   `json:"stealth_mode"`
	FlySpeedFeet      int      `json:"fly_speed_feet"`
	FlyLight          []string `json:"fly_light"`
	FlyGrip           string   `json:"fly_grip"`
	FlyEndsOn         []string `json:"fly_ends_on"`
	PolymorphSpell    string   `json:"polymorph_spell"`
	PolymorphTarget   string   `json:"polymorph_target"`
	PolymorphForm     string   `json:"polymorph_form"`
	PolymorphLight    []string `json:"polymorph_light"`
	RetainedScores    []string `json:"retained_scores"`
	PolymorphUses     int      `json:"polymorph_uses"`
	PolymorphRecharge string   `json:"polymorph_recharge"`
}

// ItemCloakMantaRay retains typed SRD source facts without implementing them.
type ItemCloakMantaRay struct {
	Condition         string `json:"condition"`
	BreatheUnderwater bool   `json:"breathe_underwater"`
	SwimSpeedFeet     int    `json:"swim_speed_feet"`
}

// ItemCrystalBall retains typed SRD source facts without implementing them.
type ItemCrystalBall struct {
	Condition string `json:"condition"`
	Spell     string `json:"spell"`
	SaveDC    int    `json:"save_dc"`
}

// ItemCrystalBallMindReader retains typed SRD source facts without implementing them.
type ItemCrystalBallMindReader struct {
	Condition             string `json:"condition"`
	Spell                 string `json:"spell"`
	SaveDC                int    `json:"save_dc"`
	ExtraSpell            string `json:"extra_spell"`
	ExtraSaveDC           int    `json:"extra_save_dc"`
	ExtraTargets          string `json:"extra_targets"`
	ExtraRangeFeet        int    `json:"extra_range_feet"`
	ExtraRangeFrom        string `json:"extra_range_from"`
	ConcentrationRequired bool   `json:"extra_concentration_required"`
	ExtraEndsWith         string `json:"extra_ends_with"`
}

// itemManifest228 is the reviewed manifest for the issue's twelve families, in
// planning inventory order. Every family has the single variant standard.
var itemManifest228 = itemSourceManifest{Batch: "data-magic-items-04", Issue: 228, Families: []itemSourceManifestFamily{
	{Slug: "censer-of-controlling-air-elementals", Name: "Censer of Controlling Air Elementals", Page: 215, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "swinging-censer once-until-next-dawn",
		Keys: "action.magic summon.air-elemental duration.one-hour item.once-until-next-dawn", Payload: "censer",
		Variants: []itemSourceManifestVariant{{"standard", "Censer of Controlling Air Elementals", "rare"}}},
	{Slug: "chime-of-opening", Name: "Chime of Opening", Page: 215, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "ten-uses",
		Keys: "action.magic spell.cast-knock item.limited-uses item.useless-after-final-use", Payload: "chime",
		Variants: []itemSourceManifestVariant{{"standard", "Chime of Opening", "rare"}}},
	{Slug: "circlet-of-blasting", Name: "Circlet of Blasting", Page: 215, Category: "wondrous-item", Attunement: "none", Equip: "wear", Use: "cast-spell", UseRestrictions: "once-until-next-dawn",
		Keys: "item.wear spell.cast-scorching-ray spell.attack-bonus-5 item.once-until-next-dawn", Payload: "circlet",
		Variants: []itemSourceManifestVariant{{"standard", "Circlet of Blasting", "uncommon"}}},
	{Slug: "cloak-of-arachnida", Name: "Cloak of Arachnida", Page: 215, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "once-until-next-dawn",
		Keys: "item.attunement item.wear defense.resistance-poison movement.climb-speed-equals-speed movement.web-immunity-and-passage spell.cast-web spell.web-save-dc-13 item.once-until-next-dawn", Payload: "arachnida",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of Arachnida", "very-rare"}}},
	{Slug: "cloak-of-displacement", Name: "Cloak of Displacement", Page: 215, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.attackers-have-disadvantage effect.suspended-after-damage effect.suppressed-at-speed-0", Payload: "displacement",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of Displacement", "rare"}}},
	{Slug: "cloak-of-elvenkind", Name: "Cloak of Elvenkind", Page: 215, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear check.perception-to-perceive-wearer-disadvantage check.stealth-advantage", Payload: "elvenkind",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of Elvenkind", "uncommon"}}},
	{Slug: "cloak-of-invisibility", Name: "Cloak of Invisibility", Page: 215, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"},
		Keys:    "item.attunement item.wear item.charges action.magic condition.invisible duration.one-hour effect.ends-when-hood-down-or-cloak-removed", Payload: "invisibility",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of Invisibility", "legendary"}}},
	{Slug: "cloak-of-protection", Name: "Cloak of Protection", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear defense.armor-class-bonus-1 save.bonus-1", Payload: "protection",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of Protection", "uncommon"}}},
	{Slug: "cloak-of-the-bat", Name: "Cloak of the Bat", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "cast-spell", UseRestrictions: "dim-light-or-darkness once-until-next-dawn",
		Keys: "item.attunement item.wear check.stealth-advantage movement.fly-speed-40-in-dim-light-or-darkness spell.cast-polymorph-self-bat item.once-until-next-dawn", Payload: "bat",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of the Bat", "rare"}}},
	{Slug: "cloak-of-the-manta-ray", Name: "Cloak of the Manta Ray", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "none",
		Keys: "item.attunement item.wear movement.breathe-underwater movement.swim-speed-60", Payload: "manta_ray",
		Variants: []itemSourceManifestVariant{{"standard", "Cloak of the Manta Ray", "uncommon"}}},
	{Slug: "crystal-ball", Name: "Crystal Ball", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "touch", Use: "cast-spell", UseRestrictions: "touching-crystal-orb",
		Keys: "item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17", Payload: "crystal_ball",
		Variants: []itemSourceManifestVariant{{"standard", "Crystal Ball", "very-rare"}}},
	{Slug: "crystal-ball-of-mind-reading", Name: "Crystal Ball of Mind Reading", Page: 216, Category: "wondrous-item", Attunement: "required", Equip: "touch", Use: "cast-spell", UseRestrictions: "touching-crystal-orb",
		Keys: "item.attunement item.touch spell.cast-scrying spell.scrying-save-dc-17 spell.cast-detect-thoughts spell.detect-thoughts-save-dc-17 concentration.detect-thoughts-ends-with-scrying", Payload: "mind_reading",
		Variants: []itemSourceManifestVariant{{"standard", "Crystal Ball of Mind Reading", "legendary"}}},
}}

// LoadItems228 returns a fresh validated copy of batch #228. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems228() (ItemSourceBatch[ItemDetails228], error) { return DecodeItems228(items228JSON) }

// DecodeItems228 validates supplied JSON as batch #228 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems228(data []byte) (ItemSourceBatch[ItemDetails228], error) {
	return decodeItemSourceBatch(data, itemManifest228, validateItemDetails228)
}

// validateItemDetails228 checks the internal consistency of each payload: DCs,
// uses, speeds and dice must be positive and key lists free of duplicates.
func validateItemDetails228(family, _ string, d *ItemDetails228) error {
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
	case "censer-of-controlling-air-elementals":
		c := d.Censer
		ok = positive(c.DurationSeconds, c.Uses) && uniqueItemKeys(c.EndTriggers)
	case "chime-of-opening":
		c := d.Chime
		ok = positive(c.LengthFeet, c.WeightPounds, c.AudibleFeet, c.Uses)
	case "circlet-of-blasting":
		c := d.Circlet
		ok = positive(c.AttackBonus, c.Uses)
	case "cloak-of-arachnida":
		c := d.Arachnida
		ok = positive(c.WebSaveDC, c.WebAreaMultiplier, c.WebUses) && uniqueItemKeys(c.Resistances) && uniqueItemKeys(c.ClimbDirections)
	case "cloak-of-invisibility":
		c := d.Invisibility
		ok = positive(c.ChargeCost, c.DurationSeconds, c.MaxCharges, c.RegainDice.Dice, c.RegainDice.Sides) && c.ChargeCost <= c.MaxCharges && c.MaxCharges == 3 && c.RegainDice.Dice == 1 && c.RegainDice.Sides == 3 && c.RegainDice.Type == "charges" && uniqueItemKeys(c.EndsEarly)
	case "cloak-of-protection":
		c := d.Protection
		ok = positive(c.ACBonus, c.SavingThrowBonus)
	case "cloak-of-the-bat":
		c := d.Bat
		ok = positive(c.FlySpeedFeet, c.PolymorphUses) && uniqueItemKeys(c.FlyLight) && uniqueItemKeys(c.FlyEndsOn) && uniqueItemKeys(c.PolymorphLight) && uniqueItemKeys(c.RetainedScores)
	case "cloak-of-the-manta-ray":
		ok = positive(d.MantaRay.SwimSpeedFeet)
	case "crystal-ball":
		ok = positive(d.CrystalBall.SaveDC)
	case "crystal-ball-of-mind-reading":
		c := d.MindReading
		ok = positive(c.SaveDC, c.ExtraSaveDC, c.ExtraRangeFeet)
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
