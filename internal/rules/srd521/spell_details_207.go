package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellFog retains Fog Cloud's typed SRD source facts without implementing them.
type SpellFog struct {
	RadiusFeet                  int    `json:"radius_feet"`
	Shape                       string `json:"shape"`
	Origin                      string `json:"origin"`
	Obscurement                 string `json:"obscurement"`
	Dispersal                   string `json:"dispersal"`
	DispersalExample            string `json:"dispersal_example"`
	AdditionalRadiusFeetPerSlot int    `json:"additional_radius_feet_per_slot"`
}

// SpellForbiddance retains Forbiddance's typed SRD source facts without implementing them.
type SpellForbiddance struct {
	AreaSquareFeet      int      `json:"area_square_feet"`
	HeightFeet          int      `json:"height_feet"`
	HeightFrom          string   `json:"height_from"`
	BlockedTravel       []string `json:"blocked_travel"`
	PortalExample       string   `json:"portal_example"`
	PlanarProofing      []string `json:"planar_proofing"`
	DamagedTypes        []string `json:"damaged_types"`
	TypeSelection       string   `json:"type_selection"`
	DamageTriggers      []string `json:"damage_triggers"`
	DamageDiceCount     int      `json:"damage_dice_count"`
	DamageDieSides      int      `json:"damage_die_sides"`
	DamageTypeChoices   []string `json:"damage_type_choices"`
	DamageTypeSelection string   `json:"damage_type_selection"`
	Password            string   `json:"password"`
	PasswordExemption   string   `json:"password_exemption"`
	OverlapRule         string   `json:"overlap_rule"`
	DailyCastingDays    int      `json:"daily_casting_days"`
	DailyCastingPlace   string   `json:"daily_casting_place"`
	DailyCastingResult  string   `json:"daily_casting_result"`
	MaterialConsumed    string   `json:"material_consumed"`
}

// SpellForcecage retains Forcecage's typed SRD source facts without implementing them.
type SpellForcecage struct {
	Prison         string   `json:"prison"`
	Forms          []string `json:"forms"`
	FormChoice     string   `json:"form_choice"`
	CageMaxSide    int      `json:"cage_max_side_feet"`
	CageBarWidth   string   `json:"cage_bar_diameter"`
	CageBarSpacing string   `json:"cage_bar_spacing"`
	BoxMaxSide     int      `json:"box_max_side_feet"`
	BoxBlocks      []string `json:"box_blocks"`
	TrappedWhen    string   `json:"trapped_when"`
	PushedAway     string   `json:"pushed_away"`
	NonmagicalExit string   `json:"nonmagical_exit"`
	TravelSave     string   `json:"travel_save"`
	TravelTriggers []string `json:"travel_triggers"`
	SaveSuccess    string   `json:"save_success"`
	SaveFailure    string   `json:"save_failure"`
	Ethereal       string   `json:"ethereal"`
	Dispel         string   `json:"dispel"`
}

// SpellForesight retains Foresight's typed SRD source facts without implementing them.
type SpellForesight struct {
	Target             string `json:"target"`
	Grants             string `json:"grants"`
	TargetAdvantage    string `json:"target_advantage"`
	OthersDisadvantage string `json:"others_disadvantage"`
	EndsEarly          string `json:"ends_early"`
}

// SpellFreedom retains Freedom of Movement's typed SRD source facts without implementing them.
type SpellFreedom struct {
	Target                   string   `json:"target"`
	DifficultTerrain         string   `json:"difficult_terrain"`
	BlockedMagic             []string `json:"blocked_magic"`
	SwimSpeed                string   `json:"swim_speed"`
	EscapeCostFeet           int      `json:"escape_cost_feet"`
	EscapeResult             string   `json:"escape_result"`
	EscapeExamples           []string `json:"escape_examples"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellFreezing retains Freezing Sphere's typed SRD source facts without implementing them.
type SpellFreezing struct {
	RadiusFeet           int      `json:"radius_feet"`
	Origin               string   `json:"origin"`
	Save                 string   `json:"save"`
	FailDamage           string   `json:"fail_damage"`
	SuccessDamage        string   `json:"success_damage"`
	DamageDiceCount      int      `json:"damage_dice_count"`
	DamageDieSides       int      `json:"damage_die_sides"`
	DamageType           string   `json:"damage_type"`
	SlotDiceCount        int      `json:"slot_dice_count"`
	SlotDieSides         int      `json:"slot_die_sides"`
	WaterFreezeDepth     int      `json:"water_freeze_depth_inches"`
	WaterFreezeAreaFeet  int      `json:"water_freeze_area_feet"`
	IceSeconds           int      `json:"ice_seconds"`
	TrappedCreatures     string   `json:"trapped_creatures"`
	TrappedCondition     string   `json:"trapped_condition"`
	BreakFreeAction      string   `json:"break_free_action"`
	BreakFreeCheck       string   `json:"break_free_check"`
	BreakFreeDC          string   `json:"break_free_dc"`
	HoldOption           string   `json:"hold_option"`
	HeldGlobeSize        string   `json:"held_globe_size"`
	HeldGlobeTemperature string   `json:"held_globe_temperature"`
	HeldGlobeLocation    string   `json:"held_globe_location"`
	ThrownBy             []string `json:"thrown_by"`
	ThrowRangeFeet       int      `json:"throw_range_feet"`
	SlingRange           string   `json:"sling_range"`
	ImpactResult         string   `json:"impact_result"`
	SetDown              string   `json:"set_down"`
	ExplodeSeconds       int      `json:"explode_seconds"`
	ExplodeCondition     string   `json:"explode_condition"`
}

// SpellGaseous retains Gaseous Form's typed SRD source facts without implementing them.
type SpellGaseous struct {
	SourcePages              []int    `json:"source_pages"`
	Target                   string   `json:"target"`
	Transforms               string   `json:"transforms"`
	Form                     string   `json:"form"`
	EndsWhen                 []string `json:"ends_when"`
	MovementOnly             string   `json:"movement_only"`
	FlySpeedFeet             int      `json:"fly_speed_feet"`
	Hover                    bool     `json:"hover"`
	EnterOccupied            string   `json:"enter_occupied"`
	Resistances              []string `json:"resistances"`
	ConditionImmunity        string   `json:"condition_immunity"`
	SaveAdvantage            []string `json:"save_advantage"`
	NarrowOpenings           string   `json:"narrow_openings"`
	Liquids                  string   `json:"liquids"`
	Prohibited               []string `json:"prohibited"`
	CarriedObjects           string   `json:"carried_objects"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
}

// SpellGate retains Gate's typed SRD source facts without implementing them.
type SpellGate struct {
	PortalShape        string   `json:"portal_shape"`
	DiameterMinFeet    int      `json:"diameter_min_feet"`
	DiameterMaxFeet    int      `json:"diameter_max_feet"`
	Orientation        string   `json:"orientation"`
	Destination        string   `json:"destination"`
	DestinationVisible bool     `json:"destination_visible"`
	Sides              string   `json:"sides"`
	Travel             string   `json:"travel"`
	Arrival            string   `json:"arrival"`
	PlanarRulers       string   `json:"planar_rulers"`
	NamedCreature      string   `json:"named_creature"`
	InvalidNames       []string `json:"invalid_names"`
	NamedOnOtherPlane  string   `json:"named_on_other_plane"`
	CreatureControl    string   `json:"creature_control"`
	CreatureBehavior   string   `json:"creature_behavior"`
	CreatureMay        []string `json:"creature_may"`
}

// SpellGeasSlotEntry retains one higher-level slot duration printed with Geas.
type SpellGeasSlotEntry struct {
	SlotLevel    int    `json:"slot_level"`
	Duration     string `json:"duration"`
	DurationDays int    `json:"duration_days"`
}

// SpellGeas retains Geas' typed SRD source facts without implementing them.
type SpellGeas struct {
	Save            string               `json:"save"`
	AutoSuccess     string               `json:"auto_success"`
	Condition       string               `json:"condition"`
	CommandKinds    []string             `json:"command_kinds"`
	DamageDiceCount int                  `json:"damage_dice_count"`
	DamageDieSides  int                  `json:"damage_die_sides"`
	DamageType      string               `json:"damage_type"`
	DamageTrigger   string               `json:"damage_trigger"`
	DamageLimit     string               `json:"damage_limit"`
	ExcludedCommand string               `json:"excluded_command"`
	SuicidalResult  string               `json:"suicidal_result"`
	EndedBy         []string             `json:"ended_by"`
	SlotDurations   []SpellGeasSlotEntry `json:"slot_durations"`
}

// SpellRepose retains Gentle Repose's typed SRD source facts without implementing them.
type SpellRepose struct {
	Target           string `json:"target"`
	ProtectedFrom    string `json:"protected_from"`
	CannotBecome     string `json:"cannot_become"`
	TimeLimitEffect  string `json:"time_limit_effect"`
	TimeLimitExample string `json:"time_limit_example"`
}

// SpellInsectJab is the Giant Insect Poison Jab action.
type SpellInsectJab struct {
	AttackKind             string `json:"attack_kind"`
	AttackBonus            string `json:"attack_bonus"`
	ReachFeet              int    `json:"reach_feet"`
	PiercingDiceCount      int    `json:"piercing_dice_count"`
	PiercingDieSides       int    `json:"piercing_die_sides"`
	PiercingBonus          int    `json:"piercing_bonus"`
	PiercingAddsSpellLevel bool   `json:"piercing_adds_spell_level"`
	PoisonDiceCount        int    `json:"poison_dice_count"`
	PoisonDieSides         int    `json:"poison_die_sides"`
}

// SpellInsectWeb is the Giant Insect Web Bolt action (Spider only).
type SpellInsectWeb struct {
	Form                      string `json:"form"`
	AttackKind                string `json:"attack_kind"`
	AttackBonus               string `json:"attack_bonus"`
	RangeFeet                 int    `json:"range_feet"`
	BludgeoningDiceCount      int    `json:"bludgeoning_dice_count"`
	BludgeoningDieSides       int    `json:"bludgeoning_die_sides"`
	BludgeoningBonus          int    `json:"bludgeoning_bonus"`
	BludgeoningAddsSpellLevel bool   `json:"bludgeoning_adds_spell_level"`
	SpeedResult               string `json:"speed_result"`
}

// SpellInsectSpew is the Giant Insect Venomous Spew bonus action (Centipede only).
type SpellInsectSpew struct {
	Form            string `json:"form"`
	Action          string `json:"action"`
	Save            string `json:"save"`
	SaveDC          string `json:"save_dc"`
	TargetRangeFeet int    `json:"target_range_feet"`
	Target          string `json:"target"`
	Failure         string `json:"failure"`
}

// SpellGiantInsect is the stat block printed with Giant Insect (p.136).
type SpellGiantInsect struct {
	Name                string             `json:"name"`
	Size                string             `json:"size"`
	CreatureType        string             `json:"creature_type"`
	Alignment           string             `json:"alignment"`
	ArmorClassBase      int                `json:"armor_class_base"`
	ArmorClassAdds      string             `json:"armor_class_adds"`
	HitPointsBase       int                `json:"hit_points_base"`
	HitPointsPerLevel   int                `json:"hit_points_per_spell_level"`
	HitPointsAboveLevel int                `json:"hit_points_above_level"`
	SpeedFeet           int                `json:"speed_feet"`
	ClimbFeet           int                `json:"climb_feet"`
	FlyFeet             int                `json:"fly_feet"`
	FlyForms            []string           `json:"fly_forms"`
	Abilities           []SpellStatAbility `json:"abilities"`
	DarkvisionFeet      int                `json:"darkvision_feet"`
	PassivePerception   int                `json:"passive_perception"`
	Languages           string             `json:"languages"`
	ChallengeRating     string             `json:"challenge_rating"`
	XP                  int                `json:"xp"`
	ProficiencyBonus    string             `json:"proficiency_bonus"`
	SpiderClimb         string             `json:"spider_climb"`
	Multiattack         string             `json:"multiattack"`
	PoisonJab           SpellInsectJab     `json:"poison_jab"`
	WebBolt             SpellInsectWeb     `json:"web_bolt"`
	VenomousSpew        SpellInsectSpew    `json:"venomous_spew"`
}

// SpellInsect retains Giant Insect's typed SRD source facts without implementing them.
type SpellInsect struct {
	SourcePages      []int            `json:"source_pages"`
	Forms            []string         `json:"forms"`
	FormChoice       string           `json:"form_choice"`
	FormDetails      string           `json:"form_details"`
	Placement        string           `json:"placement"`
	StatBlock        string           `json:"stat_block"`
	Disappears       []string         `json:"disappears"`
	Ally             string           `json:"ally"`
	Initiative       string           `json:"initiative"`
	Commands         string           `json:"commands"`
	NoCommand        string           `json:"no_command"`
	SpellLevelSource string           `json:"spell_level_source"`
	GiantInsect      SpellGiantInsect `json:"giant_insect"`
}

// SpellGlibness retains Glibness' typed SRD source facts without implementing them.
type SpellGlibness struct {
	CheckAbility    string `json:"check_ability"`
	ReplaceRollWith int    `json:"replace_roll_with"`
	ReplaceChoice   string `json:"replace_choice"`
	TruthMagic      string `json:"truth_magic"`
	TruthCondition  string `json:"truth_condition"`
}

func spell207Templates() map[string]any {
	return map[string]any{
		"fog": SpellFog{}, "forbiddance": SpellForbiddance{}, "forcecage": SpellForcecage{}, "foresight": SpellForesight{},
		"freedom": SpellFreedom{}, "freezing": SpellFreezing{}, "gaseous": SpellGaseous{}, "gate": SpellGate{},
		"geas": SpellGeas{}, "repose": SpellRepose{}, "insect": SpellInsect{}, "glibness": SpellGlibness{},
	}
}

func validSpell207Target(t SpellTarget) bool {
	switch t.Shape {
	case "created-effect":
		return t.Count == 1 && slices.Contains([]string{"floor-area-up-to-40000-square-feet", "area-in-range"}, t.Eligibility)
	case "single":
		return t.Count == 1 && t.Eligibility == "corpse-or-remains-touch"
	}
	return false
}

// spell207TargetScaling returns the base and per-slot extra targets of the
// batch's two target-scaling spells.
func spell207TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	switch {
	case d.Freedom != nil:
		return d.Freedom.BaseTargets, d.Freedom.AdditionalTargetsPerSlot, true
	case d.Gaseous != nil:
		return d.Gaseous.BaseTargets, d.Gaseous.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

type spell207Effect struct {
	id, mechanic string
	dice, sides  int
	damage       string
}

// spell207Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD pages.
type spell207Row struct {
	slug, name            string
	page, level           int
	school                string
	lists                 []string
	casting               string
	ritual, concentration bool
	rangeKind             string
	rangeFeet             int
	target                SpellTarget
	seconds               int
	endTiming             string
	verbal, somatic       bool
	material              string
	costCP                int
	consumed              bool
	scaling               SpellScaling
	effects               []spell207Effect
	required              []string
}

var spell207NoScaling = SpellScaling{"none", 0, 0, 0, []string{}}

var spell207Rows = map[string]spell207Row{
	"fog": {"fog-cloud", "Fog Cloud", 133, 1, "conjuration", []string{"druid", "ranger", "sorcerer", "wizard"}, "action", false, true, "feet", 120, SpellTarget{"created-effect", 1, "point-in-range"}, 3600, "duration-or-concentration-ended", true, true, "", 0, false,
		SpellScaling{"slot-radius", 1, 0, 0, []string{"fog"}}, []spell207Effect{{"fog", "area.heavily-obscured-fog-sphere-20", 0, 0, ""}, {"disperse", "dispel.strong-wind-disperses-fog", 0, 0, ""}},
		[]string{"area.heavily-obscured-fog-sphere-20", "dispel.strong-wind-disperses-fog", "scaling.radius-per-slot", "concentration.source-duration"}},
	"forbiddance": {"forbiddance", "Forbiddance", 133, 6, "abjuration", []string{"cleric"}, "10-minutes", true, false, "touch", 0, SpellTarget{"created-effect", 1, "floor-area-up-to-40000-square-feet"}, 86400, "duration-ended", true, true, "ruby dust worth 1,000+ GP", 100000, false,
		spell207NoScaling, []spell207Effect{{"ward", "protection.forbiddance-teleport-portal-ward", 0, 0, ""}, {"planar", "protection.forbiddance-planar-travel-proofing", 0, 0, ""}, {"damage", "damage.forbiddance-5d10", 5, 10, "radiant-or-necrotic"}, {"password", "exemption.forbiddance-password", 0, 0, ""}, {"overlap", "limit.forbiddance-no-overlap", 0, 0, ""}, {"permanence", "duration.forbiddance-daily-30-day-permanence", 0, 0, ""}},
		[]string{"protection.forbiddance-teleport-portal-ward", "protection.forbiddance-planar-travel-proofing", "damage.forbiddance-5d10", "exemption.forbiddance-password", "limit.forbiddance-no-overlap", "duration.forbiddance-daily-30-day-permanence", "cast.ten-minutes", "cast.ritual", "components.priced-material", "components.material"}},
	"forcecage": {"forcecage", "Forcecage", 134, 7, "evocation", []string{"bard", "warlock", "wizard"}, "action", false, true, "feet", 100, SpellTarget{"created-effect", 1, "area-in-range"}, 3600, "duration-or-concentration-ended", true, true, "ruby dust worth 1,500+ GP, which the spell consumes", 150000, true,
		spell207NoScaling, []spell207Effect{{"prison", "area.forcecage-prison-cage-or-box", 0, 0, ""}, {"trap", "trap.forcecage-creatures-inside-or-pushed-out", 0, 0, ""}, {"exit", "save.forcecage-teleport-or-planar-exit", 0, 0, ""}, {"ethereal", "block.forcecage-ethereal-travel", 0, 0, ""}, {"dispel", "immunity.forcecage-dispel-magic", 0, 0, ""}},
		[]string{"area.forcecage-prison-cage-or-box", "trap.forcecage-creatures-inside-or-pushed-out", "save.forcecage-teleport-or-planar-exit", "block.forcecage-ethereal-travel", "immunity.forcecage-dispel-magic", "components.priced-material", "components.material", "concentration.source-duration"}},
	"foresight": {"foresight", "Foresight", 134, 9, "divination", []string{"bard", "druid", "warlock", "wizard"}, "1-minute", false, false, "touch", 0, SpellTarget{"single", 1, "willing-creature-touch"}, 28800, "duration-ended-or-recast", true, true, "a hummingbird feather", 0, false,
		spell207NoScaling, []spell207Effect{{"advantage", "buff.foresight-d20-advantage", 0, 0, ""}, {"defense", "buff.foresight-attackers-disadvantage", 0, 0, ""}, {"ending", "ending.foresight-recast", 0, 0, ""}},
		[]string{"buff.foresight-d20-advantage", "buff.foresight-attackers-disadvantage", "ending.foresight-recast", "cast.one-minute", "components.material"}},
	"freedom": {"freedom-of-movement", "Freedom of Movement", 134, 4, "abjuration", []string{"bard", "cleric", "druid", "ranger"}, "action", false, false, "touch", 0, SpellTarget{"single", 1, "willing-creature-touch"}, 3600, "duration-ended", true, true, "a leather strap", 0, false,
		SpellScaling{"slot-targets", 4, 0, 0, []string{"protection"}}, []spell207Effect{{"protection", "protection.freedom-of-movement-terrain-and-magic", 0, 0, ""}, {"swim", "movement.freedom-of-movement-swim-speed", 0, 0, ""}, {"escape", "movement.freedom-of-movement-escape-restraints", 0, 0, ""}},
		[]string{"protection.freedom-of-movement-terrain-and-magic", "movement.freedom-of-movement-swim-speed", "movement.freedom-of-movement-escape-restraints", "scaling.targets-per-slot", "components.material"}},
	"freezing": {"freezing-sphere", "Freezing Sphere", 134, 6, "evocation", []string{"sorcerer", "wizard"}, "action", false, false, "feet", 300, SpellTarget{"sphere", 0, "creatures-in-sphere"}, 0, "instantaneous", true, true, "a miniature crystal sphere", 0, false,
		SpellScaling{"slot", 6, 1, 6, []string{"damage"}}, []spell207Effect{{"damage", "damage.freezing-sphere", 10, 6, "cold"}, {"freeze", "effect.freezing-sphere-ice", 0, 0, ""}, {"break", "check.freezing-sphere-break-free", 0, 0, ""}, {"globe", "object.freezing-sphere-held-globe", 0, 0, ""}},
		[]string{"damage.freezing-sphere", "effect.freezing-sphere-ice", "check.freezing-sphere-break-free", "object.freezing-sphere-held-globe", "area.sphere-60-radius", "save.constitution-half", "scaling.per-slot", "components.material"}},
	"gaseous": {"gaseous-form", "Gaseous Form", 134, 3, "transmutation", []string{"sorcerer", "warlock", "wizard"}, "action", false, true, "touch", 0, SpellTarget{"single", 1, "willing-creature-touch"}, 3600, "duration-or-concentration-ended", true, true, "a bit of gauze", 0, false,
		SpellScaling{"slot-targets", 3, 0, 0, []string{"form"}}, []spell207Effect{{"form", "transformation.gaseous-form", 0, 0, ""}, {"defense", "defense.gaseous-form-resistances-and-advantage", 0, 0, ""}, {"restrictions", "restriction.gaseous-form-no-speech-objects-or-attacks", 0, 0, ""}, {"ending", "ending.gaseous-form-zero-hit-points-or-magic-action", 0, 0, ""}},
		[]string{"transformation.gaseous-form", "defense.gaseous-form-resistances-and-advantage", "restriction.gaseous-form-no-speech-objects-or-attacks", "ending.gaseous-form-zero-hit-points-or-magic-action", "scaling.targets-per-slot", "components.material", "concentration.source-duration"}},
	"gate": {"gate", "Gate", 135, 9, "conjuration", []string{"cleric", "sorcerer", "warlock", "wizard"}, "action", false, true, "feet", 60, SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, 60, "duration-or-concentration-ended", true, true, "a diamond worth 5,000+ GP", 500000, false,
		spell207NoScaling, []spell207Effect{{"portal", "portal.gate-planar", 0, 0, ""}, {"named", "summon.gate-named-creature", 0, 0, ""}, {"rulers", "restriction.gate-planar-ruler-prevention", 0, 0, ""}},
		[]string{"portal.gate-planar", "summon.gate-named-creature", "restriction.gate-planar-ruler-prevention", "components.priced-material", "components.material", "concentration.source-duration"}},
	"geas": {"geas", "Geas", 135, 5, "enchantment", []string{"bard", "cleric", "druid", "paladin", "wizard"}, "1-minute", false, false, "feet", 60, SpellTarget{"single", 1, "visible-creature-in-range"}, 2592000, "duration-ended-suicidal-command-or-listed-spell", true, false, "", 0, false,
		SpellScaling{"slot-duration", 5, 0, 0, []string{"command"}}, []spell207Effect{{"command", "charm.geas-command", 0, 0, ""}, {"damage", "damage.geas-psychic", 5, 10, "psychic"}, {"suicidal", "ending.geas-suicidal-command", 0, 0, ""}, {"ending", "ending.geas-listed-spells", 0, 0, ""}},
		[]string{"charm.geas-command", "damage.geas-psychic", "ending.geas-suicidal-command", "ending.geas-listed-spells", "cast.one-minute", "scaling.geas-slot-duration"}},
	"repose": {"gentle-repose", "Gentle Repose", 135, 2, "necromancy", []string{"cleric", "paladin", "wizard"}, "action", true, false, "touch", 0, SpellTarget{"single", 1, "corpse-or-remains-touch"}, 864000, "duration-ended", true, true, "2 Copper Pieces, which the spell consumes", 2, true,
		spell207NoScaling, []spell207Effect{{"protection", "protection.gentle-repose-decay-and-undeath", 0, 0, ""}, {"time", "time.gentle-repose-extends-raising-limit", 0, 0, ""}},
		[]string{"protection.gentle-repose-decay-and-undeath", "time.gentle-repose-extends-raising-limit", "cast.ritual", "components.material"}},
	"insect": {"giant-insect", "Giant Insect", 135, 4, "conjuration", []string{"druid"}, "action", false, true, "feet", 60, SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, 600, "duration-or-concentration-ended", true, true, "", 0, false,
		SpellScaling{"slot-spell-level", 4, 0, 0, []string{"summon"}}, []spell207Effect{{"summon", "summon.giant-insect", 0, 0, ""}, {"control", "control.giant-insect-verbal-commands", 0, 0, ""}, {"jab", "damage.giant-insect-poison-jab-piercing", 1, 6, "piercing"}, {"jab-poison", "damage.giant-insect-poison-jab-poison", 1, 4, "poison"}, {"web", "damage.giant-insect-web-bolt", 1, 10, "bludgeoning"}, {"spew", "condition.giant-insect-venomous-spew-poisoned", 0, 0, ""}},
		[]string{"summon.giant-insect", "control.giant-insect-verbal-commands", "damage.giant-insect-poison-jab-piercing", "damage.giant-insect-poison-jab-poison", "damage.giant-insect-web-bolt", "condition.giant-insect-venomous-spew-poisoned", "creature.giant-insect-stat-block", "initiative.shared-after-caster", "scaling.stat-block-spell-level", "concentration.source-duration"}},
	"glibness": {"glibness", "Glibness", 136, 8, "enchantment", []string{"bard", "warlock"}, "action", false, false, "self", 0, SpellTarget{"single", 1, "self"}, 3600, "duration-ended", true, false, "", 0, false,
		spell207NoScaling, []spell207Effect{{"check", "check.glibness-charisma-15", 0, 0, ""}, {"truth", "truth.glibness-magic-indicates-truthful", 0, 0, ""}},
		[]string{"check.glibness-charisma-15", "truth.glibness-magic-indicates-truthful"}},
}

func (s SpellRecord) validateSpell207() (bool, error) {
	expected := ""
	for key, row := range spell207Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 207 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"fog": d.Fog, "forbiddance": d.Forbiddance, "forcecage": d.Forcecage, "foresight": d.Foresight,
		"freedom": d.Freedom, "freezing": d.Freezing, "gaseous": d.Gaseous, "gate": d.Gate,
		"geas": d.Geas, "repose": d.Repose, "insect": d.Insect, "glibness": d.Glibness,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 207 payload")
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return true, err
		}
		if err := requireSpell198Fields(raw, reflect.TypeOf(value).Elem()); err != nil {
			return true, err
		}
	}
	if expected == "" && count == 0 {
		return false, nil
	}
	row := spell207Rows[expected]
	castingSeconds := map[string]int{"action": 0, "1-minute": 60, "10-minutes": 600}
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != castingSeconds[row.casting] || d.CastingTriggers == nil || len(d.CastingTriggers) != 0 {
		return true, fmt.Errorf("invalid batch 207 source details")
	}
	if err := s.validateSpell207Header(row); err != nil {
		return true, err
	}
	if err := d.validateSpell207Payload(expected); err != nil {
		return true, err
	}
	return true, nil
}

func (s SpellRecord) validateSpell207Header(row spell207Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	duration := "timed"
	if row.seconds == 0 {
		duration = "instantaneous"
	}
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != row.casting || flags(s.Ritual) != row.ritual || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || flags(c.Verbal) != row.verbal || flags(c.Somatic) != row.somatic {
		return fmt.Errorf("batch 207 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 207 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed != row.consumed {
		return fmt.Errorf("batch 207 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 207 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 207 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != 0 || e.AddCastingModifier == nil || *e.AddCastingModifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 207 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 207 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 207 slot levels differ from source")
	}
	return nil
}

var spell207GeasSlots = []SpellGeasSlotEntry{{7, "days", 365}, {8, "days", 365}, {9, "until-ended-by-listed-spell", 0}}

func (d *SpellSourceDetails) validateSpell207Payload(key string) error {
	bad := func() error { return fmt.Errorf("invalid batch 207 %s payload", key) }
	str := func(a ...string) []string { return a }
	switch key {
	case "fog":
		x := d.Fog
		if x.RadiusFeet != 20 || x.Shape != "sphere" || x.Origin != "point-in-range" || x.Obscurement != "heavily-obscured" || x.Dispersal != "strong-wind" || x.DispersalExample != "gust-of-wind" || x.AdditionalRadiusFeetPerSlot != 20 {
			return bad()
		}
	case "forbiddance":
		x := d.Forbiddance
		if x.AreaSquareFeet != 40000 || x.HeightFeet != 30 || x.HeightFrom != "above-the-floor" || !slices.Equal(x.BlockedTravel, str("teleporting-into-area", "using-portals-to-enter-area")) || x.PortalExample != "gate" ||
			!slices.Equal(x.PlanarProofing, str("astral-plane", "ethereal-plane", "feywild", "shadowfell", "plane-shift")) || !slices.Equal(x.DamagedTypes, str("aberrations", "celestials", "elementals", "fey", "fiends", "undead")) ||
			x.TypeSelection != "choose-one-or-more-when-casting" || !slices.Equal(x.DamageTriggers, str("enters-area-first-time-on-turn", "ends-turn-in-area")) || x.DamageDiceCount != 5 || x.DamageDieSides != 10 ||
			!slices.Equal(x.DamageTypeChoices, str("radiant", "necrotic")) || x.DamageTypeSelection != "chosen-when-casting" || x.Password != "may-designate-when-casting" || x.PasswordExemption != "speaks-password-as-it-enters-takes-no-damage" ||
			x.OverlapRule != "area-cannot-overlap-another-forbiddance" || x.DailyCastingDays != 30 || x.DailyCastingPlace != "same-location" || x.DailyCastingResult != "lasts-until-dispelled" || x.MaterialConsumed != "on-the-last-casting" {
			return bad()
		}
	case "forcecage":
		x := d.Forcecage
		if x.Prison != "immobile-invisible-cube-of-magical-force" || !slices.Equal(x.Forms, str("cage", "solid-box")) || x.FormChoice != "caster-choice" || x.CageMaxSide != 20 || x.CageBarWidth != "1/2-inch" || x.CageBarSpacing != "1/2-inch" ||
			x.BoxMaxSide != 10 || !slices.Equal(x.BoxBlocks, str("any-matter-passing-through", "spells-cast-into-or-out-of-area")) || x.TrappedWhen != "completely-inside-cage-area-when-cast" ||
			x.PushedAway != "partially-inside-or-too-large-pushed-away-from-center-until-completely-outside" || x.NonmagicalExit != "cannot-leave-by-nonmagical-means" || x.TravelSave != "charisma" ||
			!slices.Equal(x.TravelTriggers, str("teleportation", "interplanar-travel")) || x.SaveSuccess != "can-use-that-magic-to-exit" || x.SaveFailure != "does-not-exit-and-wastes-spell-or-effect" ||
			x.Ethereal != "cage-extends-into-ethereal-plane-blocking-ethereal-travel" || x.Dispel != "cannot-be-dispelled-by-dispel-magic" {
			return bad()
		}
	case "foresight":
		x := d.Foresight
		if x.Target != "willing-creature-touched" || x.Grants != "limited-ability-to-see-into-immediate-future" || x.TargetAdvantage != "d20-tests" || x.OthersDisadvantage != "attack-rolls-against-target" || x.EndsEarly != "caster-casts-it-again" {
			return bad()
		}
	case "freedom":
		x := d.Freedom
		if x.Target != "willing-creature-touched" || x.DifficultTerrain != "movement-unaffected" || !slices.Equal(x.BlockedMagic, str("reduce-speed", "paralyzed-condition", "restrained-condition")) || x.SwimSpeed != "equal-to-speed" ||
			x.EscapeCostFeet != 5 || x.EscapeResult != "automatically-escape-nonmagical-restraints" || !slices.Equal(x.EscapeExamples, str("manacles", "creature-imposing-the-grappled-condition")) || x.BaseTargets != 1 || x.AdditionalTargetsPerSlot != 1 {
			return bad()
		}
	case "freezing":
		x := d.Freezing
		if x.RadiusFeet != 60 || x.Origin != "globe-streaks-from-caster-to-chosen-point-in-range" || x.Save != "constitution" || x.FailDamage != "full" || x.SuccessDamage != "half" || x.DamageDiceCount != 10 || x.DamageDieSides != 6 || x.DamageType != "cold" ||
			x.SlotDiceCount != 1 || x.SlotDieSides != 6 || x.WaterFreezeDepth != 6 || x.WaterFreezeAreaFeet != 30 || x.IceSeconds != 60 || x.TrappedCreatures != "creatures-swimming-on-surface-of-frozen-water" || x.TrappedCondition != "restrained" ||
			x.BreakFreeAction != "action" || x.BreakFreeCheck != "strength-athletics" || x.BreakFreeDC != "caster-spell-save-dc" || x.HoldOption != "refrain-from-firing-after-completing-casting" || x.HeldGlobeSize != "sling-bullet" ||
			x.HeldGlobeTemperature != "cool-to-the-touch" || x.HeldGlobeLocation != "caster-hand" || !slices.Equal(x.ThrownBy, str("caster", "creature-given-the-globe")) || x.ThrowRangeFeet != 40 || x.SlingRange != "slings-normal-range" ||
			x.ImpactResult != "shatters-with-same-effect-as-normal-casting" || x.SetDown != "can-set-down-without-shattering" || x.ExplodeSeconds != 60 || x.ExplodeCondition != "globe-has-not-already-shattered" {
			return bad()
		}
	case "gaseous":
		x := d.Gaseous
		if !slices.Equal(x.SourcePages, []int{134, 135}) || x.Target != "willing-creature-touched" || x.Transforms != "target-and-everything-it-is-wearing-and-carrying" || x.Form != "misty-cloud" ||
			!slices.Equal(x.EndsWhen, str("target-drops-to-0-hit-points", "target-takes-magic-action-to-end-spell-on-itself")) || x.MovementOnly != "fly-speed" || x.FlySpeedFeet != 10 || !x.Hover ||
			x.EnterOccupied != "can-enter-and-occupy-space-of-another-creature" || !slices.Equal(x.Resistances, str("bludgeoning", "piercing", "slashing")) || x.ConditionImmunity != "prone" ||
			!slices.Equal(x.SaveAdvantage, str("strength", "dexterity", "constitution")) || x.NarrowOpenings != "can-pass-through-narrow-open-openings" || x.Liquids != "treated-as-solid-surfaces" ||
			!slices.Equal(x.Prohibited, str("talk", "manipulate-objects", "attack", "cast-spells")) || x.CarriedObjects != "cannot-be-dropped-used-or-interacted-with" || x.BaseTargets != 1 || x.AdditionalTargetsPerSlot != 1 {
			return bad()
		}
	case "gate":
		x := d.Gate
		if x.PortalShape != "circular-opening" || x.DiameterMinFeet != 5 || x.DiameterMaxFeet != 20 || x.Orientation != "any-direction-caster-chooses" || x.Destination != "precise-location-on-different-plane" || !x.DestinationVisible ||
			x.Sides != "front-and-back-on-each-plane" || x.Travel != "only-by-moving-through-front" || x.Arrival != "instantly-transported-to-unoccupied-space-nearest-portal-on-other-plane" ||
			x.PlanarRulers != "can-prevent-portals-opening-in-presence-or-domain" || x.NamedCreature != "may-speak-name-of-specific-creature-when-casting" || !slices.Equal(x.InvalidNames, str("pseudonym", "title", "nickname")) ||
			x.NamedOnOtherPlane != "portal-opens-next-to-named-creature-and-transports-it-to-nearest-unoccupied-space-on-casters-side" || x.CreatureControl != "no-special-power-over-creature" ||
			x.CreatureBehavior != "free-to-act-as-gm-deems-appropriate" || !slices.Equal(x.CreatureMay, str("leave", "attack-caster", "help-caster")) {
			return bad()
		}
	case "geas":
		x := d.Geas
		if x.Save != "wisdom" || x.AutoSuccess != "target-cannot-understand-command" || x.Condition != "charmed" || !slices.Equal(x.CommandKinds, str("carry-out-service", "refrain-from-action-or-course-of-activity")) ||
			x.DamageDiceCount != 5 || x.DamageDieSides != 10 || x.DamageType != "psychic" || x.DamageTrigger != "acts-directly-counter-to-command" || x.DamageLimit != "no-more-than-once-each-day" ||
			x.ExcludedCommand != "activity-that-would-result-in-certain-death" || x.SuicidalResult != "spell-ends" || !slices.Equal(x.EndedBy, str("remove-curse", "greater-restoration", "wish")) || !slices.Equal(x.SlotDurations, spell207GeasSlots) {
			return bad()
		}
	case "repose":
		x := d.Repose
		if x.Target != "corpse-or-other-remains-touched" || x.ProtectedFrom != "decay" || x.CannotBecome != "undead" || x.TimeLimitEffect != "days-under-spell-do-not-count-against-time-limit" || x.TimeLimitExample != "raise-dead" {
			return bad()
		}
	case "insect":
		x := d.Insect
		if !slices.Equal(x.SourcePages, []int{135, 136}) || !slices.Equal(x.Forms, str("centipede", "spider", "wasp")) || x.FormChoice != "chosen-when-casting" || x.FormDetails != "form-determines-certain-stat-block-details" ||
			x.Placement != "unoccupied-space-visible-in-range" || x.StatBlock != "giant-insect" || !slices.Equal(x.Disappears, str("drops-to-0-hit-points", "spell-ends")) || x.Ally != "caster-and-caster-allies" ||
			x.Initiative != "shares-caster-initiative-count-takes-turn-immediately-after-caster" || x.Commands != "verbal-no-action-required" || x.NoCommand != "takes-dodge-action-and-uses-movement-to-avoid-danger" || x.SpellLevelSource != "spell-slot-level" {
			return bad()
		}
		b := x.GiantInsect
		wantAbilities := []SpellStatAbility{{"strength", 17, 3, 3}, {"dexterity", 13, 1, 1}, {"constitution", 15, 2, 2}, {"intelligence", 4, -3, -3}, {"wisdom", 14, 2, 2}, {"charisma", 3, -4, -4}}
		wantJab := SpellInsectJab{"melee-attack-roll", "caster-spell-attack-modifier", 10, 1, 6, 3, true, 1, 4}
		wantWeb := SpellInsectWeb{"spider", "ranged-attack-roll", "caster-spell-attack-modifier", 60, 1, 10, 3, true, "speed-reduced-to-0-until-start-of-insects-next-turn"}
		wantSpew := SpellInsectSpew{"centipede", "bonus-action", "constitution", "caster-spell-save-dc", 10, "one-creature-the-insect-can-see", "poisoned-condition-until-start-of-insects-next-turn"}
		if b.Name != "Giant Insect" || b.Size != "large" || b.CreatureType != "beast" || b.Alignment != "unaligned" || b.ArmorClassBase != 11 || b.ArmorClassAdds != "spell-level" || b.HitPointsBase != 30 || b.HitPointsPerLevel != 10 || b.HitPointsAboveLevel != 4 ||
			b.SpeedFeet != 40 || b.ClimbFeet != 40 || b.FlyFeet != 40 || !slices.Equal(b.FlyForms, str("wasp")) || !slices.Equal(b.Abilities, wantAbilities) || b.DarkvisionFeet != 60 || b.PassivePerception != 12 ||
			b.Languages != "understands-languages-caster-knows" || b.ChallengeRating != "none" || b.XP != 0 || b.ProficiencyBonus != "equals-caster-proficiency-bonus" ||
			b.SpiderClimb != "climbs-difficult-surfaces-including-ceilings-without-ability-check" || b.Multiattack != "attacks-equal-to-half-spell-level-rounded-down" || b.PoisonJab != wantJab || b.WebBolt != wantWeb || b.VenomousSpew != wantSpew {
			return bad()
		}
	case "glibness":
		x := d.Glibness
		if x.CheckAbility != "charisma" || x.ReplaceRollWith != 15 || x.ReplaceChoice != "optional" || x.TruthMagic != "magic-that-determines-truthfulness-indicates-truthful" || x.TruthCondition != "no-matter-what-caster-says" {
			return bad()
		}
	}
	return nil
}

// validSpell207Scaling checks the one batch scaling kind that carries no dice
// or count: the Giant Insect stat block takes the slot's level as its spell level.
func (s SpellRecord) validSpell207Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	return d != nil && d.Insect != nil && sc.Kind == "slot-spell-level" && sc.AboveLevel == s.Level && sc.DiceCount == 0 && sc.DieSides == 0 && len(sc.EffectIDs) == 1 && effects[sc.EffectIDs[0]] &&
		d.Insect.SpellLevelSource == "spell-slot-level" && slices.Contains(s.RequiredMechanics, "scaling.stat-block-spell-level")
}
