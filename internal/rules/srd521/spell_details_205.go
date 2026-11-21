package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellEtherealness retains Etherealness' typed SRD source facts without implementing them.
type SpellEtherealness struct {
	Region                   string   `json:"region"`
	Movement                 string   `json:"movement"`
	VerticalExtraFeetPerFoot int      `json:"vertical_extra_feet_per_foot"`
	PerceivedPlaneAppearance string   `json:"perceived_plane_appearance"`
	PerceivedPlaneLimitFeet  int      `json:"perceived_plane_limit_feet"`
	AffectsAndAffectedBy     string   `json:"affects_and_affected_by"`
	OutsideCreatures         string   `json:"outside_creatures"`
	ReturnPlace              string   `json:"return_place"`
	OccupiedResult           string   `json:"occupied_result"`
	ShuntDamageType          string   `json:"shunt_damage_type"`
	ShuntDamagePerFootMoved  int      `json:"shunt_damage_per_foot_moved"`
	InstantEndConditions     []string `json:"instant_end_conditions"`
	InstantEndExample        string   `json:"instant_end_example"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
	SlotTargetKind           string   `json:"slot_target_kind"`
	SlotTargetRangeFeet      int      `json:"slot_target_range_feet"`
}

// SpellRetreat retains Expeditious Retreat's typed SRD source facts without implementing them.
type SpellRetreat struct {
	CastAction       string `json:"cast_action"`
	RepeatAction     string `json:"repeat_action"`
	RepeatActionType string `json:"repeat_action_type"`
	RepeatWhile      string `json:"repeat_while"`
}

// SpellEyebite retains Eyebite's typed SRD source facts without implementing them.
type SpellEyebite struct {
	Appearance              string   `json:"appearance"`
	TargetRangeFeet         int      `json:"target_range_feet"`
	TargetOrigin            string   `json:"target_origin"`
	TargetVisibility        string   `json:"target_visibility"`
	Save                    string   `json:"save"`
	Failure                 string   `json:"failure"`
	Options                 []string `json:"options"`
	RetargetAction          string   `json:"retarget_action"`
	RetargetTiming          string   `json:"retarget_timing"`
	RetargetExclusion       string   `json:"retarget_exclusion"`
	AsleepCondition         string   `json:"asleep_condition"`
	AsleepWakesWhen         []string `json:"asleep_wakes_when"`
	PanickedCondition       string   `json:"panicked_condition"`
	PanickedTurnAction      string   `json:"panicked_turn_action"`
	PanickedEndDistanceFeet int      `json:"panicked_end_distance_feet"`
	PanickedEndWhen         string   `json:"panicked_end_when"`
	SickenedCondition       string   `json:"sickened_condition"`
}

// SpellFabricate retains Fabricate's typed SRD source facts without implementing them.
type SpellFabricate struct {
	Input                 string   `json:"input"`
	Output                string   `json:"output"`
	Examples              []string `json:"examples"`
	ObjectSize            string   `json:"object_size"`
	ObjectCubeFeet        int      `json:"object_cube_feet"`
	ConnectedCubeCount    int      `json:"connected_cube_count"`
	ConnectedCubeFeet     int      `json:"connected_cube_feet"`
	Quantity              string   `json:"quantity"`
	MineralKinds          []string `json:"mineral_kinds"`
	MineralObjectSize     string   `json:"mineral_object_size"`
	MineralCubeFeet       int      `json:"mineral_cube_feet"`
	Quality               string   `json:"quality"`
	CannotCreate          []string `json:"cannot_create"`
	SkilledItems          string   `json:"skilled_items"`
	SkilledItemsException string   `json:"skilled_items_exception"`
}

// SpellFaerieFire retains Faerie Fire's typed SRD source facts without implementing them.
type SpellFaerieFire struct {
	CubeFeet         int      `json:"cube_feet"`
	ObjectsOutlined  string   `json:"objects_outlined"`
	CreatureOutline  string   `json:"creature_outline"`
	Save             string   `json:"save"`
	OutlineColors    []string `json:"outline_colors"`
	ColorChoice      string   `json:"color_choice"`
	LightKind        string   `json:"light_kind"`
	LightRadiusFeet  int      `json:"light_radius_feet"`
	LightSource      string   `json:"light_source"`
	BlockedCondition string   `json:"blocked_condition"`
	AttackRollEffect string   `json:"attack_roll_effect"`
	AttackRollTarget string   `json:"attack_roll_target"`
}

// SpellHound retains Faithful Hound's typed SRD source facts without implementing them.
type SpellHound struct {
	Form                       string   `json:"form"`
	Placement                  string   `json:"placement"`
	EndsWhen                   []string `json:"ends_when"`
	SeparationFeet             int      `json:"separation_feet"`
	VisibleTo                  string   `json:"visible_to"`
	Traits                     []string `json:"traits"`
	BarkTrigger                string   `json:"bark_trigger"`
	BarkTriggerFeet            int      `json:"bark_trigger_feet"`
	Password                   string   `json:"password"`
	BarkResult                 string   `json:"bark_result"`
	TruesightFeet              int      `json:"truesight_feet"`
	BiteTiming                 string   `json:"bite_timing"`
	BiteTarget                 string   `json:"bite_target"`
	BiteReachFeet              int      `json:"bite_reach_feet"`
	BiteSave                   string   `json:"bite_save"`
	BiteFailureDamageDiceCount int      `json:"bite_failure_damage_dice_count"`
	BiteFailureDamageDieSides  int      `json:"bite_failure_damage_die_sides"`
	BiteFailureDamageType      string   `json:"bite_failure_damage_type"`
	MoveAction                 string   `json:"move_action"`
	MoveTiming                 string   `json:"move_timing"`
	MoveFeet                   int      `json:"move_feet"`
}

// SpellFalseLife retains False Life's typed SRD source facts without implementing them.
type SpellFalseLife struct {
	Recipient                           string `json:"recipient"`
	TemporaryHitPointsDiceCount         int    `json:"temporary_hit_points_dice_count"`
	TemporaryHitPointsDieSides          int    `json:"temporary_hit_points_die_sides"`
	TemporaryHitPointsBonus             int    `json:"temporary_hit_points_bonus"`
	AdditionalTemporaryHitPointsPerSlot int    `json:"additional_temporary_hit_points_per_slot"`
}

// SpellFear retains Fear's typed SRD source facts without implementing them.
type SpellFear struct {
	ConeFeet            int      `json:"cone_feet"`
	Origin              string   `json:"origin"`
	Save                string   `json:"save"`
	FailureEffects      []string `json:"failure_effects"`
	Condition           string   `json:"condition"`
	TurnAction          string   `json:"turn_action"`
	TurnActionException string   `json:"turn_action_exception"`
	EndTurnTrigger      string   `json:"end_turn_trigger"`
	EndTurnSave         string   `json:"end_turn_save"`
	EndTurnSaveSuccess  string   `json:"end_turn_save_success"`
}

// SpellFeatherFall retains Feather Fall's typed SRD source facts without implementing them.
type SpellFeatherFall struct {
	TriggerRangeFeet    int    `json:"trigger_range_feet"`
	Trigger             string `json:"trigger"`
	MaxTargets          int    `json:"max_targets"`
	TargetState         string `json:"target_state"`
	DescentFeetPerRound int    `json:"descent_feet_per_round"`
	DescentUntil        string `json:"descent_until"`
	LandingResult       string `json:"landing_result"`
}

// SpellFamiliar retains Find Familiar's typed SRD source facts without implementing them.
type SpellFamiliar struct {
	FormChoices            []string `json:"form_choices"`
	Appearance             string   `json:"appearance"`
	Statistics             string   `json:"statistics"`
	CreatureTypeChoices    []string `json:"creature_type_choices"`
	CreatureTypeReplaces   string   `json:"creature_type_replaces"`
	Independence           string   `json:"independence"`
	TelepathyRangeFeet     int      `json:"telepathy_range_feet"`
	SeeThroughAction       string   `json:"see_through_action"`
	SeeThroughDuration     string   `json:"see_through_duration"`
	SeeThroughBenefit      string   `json:"see_through_benefit"`
	TouchDelivery          string   `json:"touch_delivery"`
	TouchDeliveryRangeFeet int      `json:"touch_delivery_range_feet"`
	TouchDeliveryAction    string   `json:"touch_delivery_action"`
	TouchDeliveryTiming    string   `json:"touch_delivery_timing"`
	CombatRole             string   `json:"combat_role"`
	CombatInitiative       string   `json:"combat_initiative"`
	AttackRestriction      string   `json:"attack_restriction"`
	ZeroHitPoints          string   `json:"zero_hit_points"`
	ReappearsWhen          string   `json:"reappears_when"`
	DismissAction          string   `json:"dismiss_action"`
	DismissOptions         []string `json:"dismiss_options"`
	RecallAction           string   `json:"recall_action"`
	RecallRangeFeet        int      `json:"recall_range_feet"`
	RecallSpace            string   `json:"recall_space"`
	LeavesBehind           string   `json:"leaves_behind"`
	LeavesBehindWhen       []string `json:"leaves_behind_when"`
	FamiliarLimit          string   `json:"familiar_limit"`
	RecastResult           string   `json:"recast_result"`
}

// SpellSteed retains Find Steed's typed SRD source facts without implementing them.
type SpellSteed struct {
	Summons               string   `json:"summons"`
	Placement             string   `json:"placement"`
	StatBlock             string   `json:"stat_block"`
	ReplacesExisting      string   `json:"replaces_existing"`
	Appearance            string   `json:"appearance"`
	AppearanceExamples    []string `json:"appearance_examples"`
	CreatureTypeChoices   []string `json:"creature_type_choices"`
	CreatureTypeTiming    string   `json:"creature_type_timing"`
	CombatRole            string   `json:"combat_role"`
	Initiative            string   `json:"initiative"`
	MountFunction         string   `json:"mount_function"`
	IncapacitatedBehavior string   `json:"incapacitated_behavior"`
	DisappearsWhen        []string `json:"disappears_when"`
	LeavesBehind          string   `json:"leaves_behind"`
	RecastChoice          string   `json:"recast_choice"`
	SlotLevelUse          string   `json:"slot_level_use"`
}

// SpellTraps retains Find Traps' typed SRD source facts without implementing them.
type SpellTraps struct {
	Senses            string   `json:"senses"`
	TrapDefinition    string   `json:"trap_definition"`
	SensedExamples    []string `json:"sensed_examples"`
	NotSensedExamples []string `json:"not_sensed_examples"`
	Reveals           string   `json:"reveals"`
	DoesNotReveal     string   `json:"does_not_reveal"`
	Learns            string   `json:"learns"`
}

func spell205Templates() map[string]any {
	return map[string]any{
		"etherealness": SpellEtherealness{},
		"retreat":      SpellRetreat{},
		"eyebite":      SpellEyebite{},
		"fabricate":    SpellFabricate{},
		"faerie_fire":  SpellFaerieFire{},
		"hound":        SpellHound{},
		"false_life":   SpellFalseLife{},
		"fear":         SpellFear{},
		"feather_fall": SpellFeatherFall{},
		"familiar":     SpellFamiliar{},
		"steed":        SpellSteed{},
		"traps":        SpellTraps{},
	}
}

// spell205Pins holds every printed value of each payload, by equality with PDF pp.128-131.
var spell205Pins = map[string]any{
	"etherealness": SpellEtherealness{Region: "border-ethereal-where-it-overlaps-your-current-plane", Movement: "any-direction", VerticalExtraFeetPerFoot: 1, PerceivedPlaneAppearance: "gray", PerceivedPlaneLimitFeet: 60, AffectsAndAffectedBy: "only-creatures-objects-and-effects-on-the-ethereal-plane", OutsideCreatures: "cannot-perceive-or-interact-unless-a-feature-allows", ReturnPlace: "plane-you-left-in-space-corresponding-to-your-border-ethereal-space", OccupiedResult: "shunted-to-nearest-unoccupied-space", ShuntDamageType: "force", ShuntDamagePerFootMoved: 2, InstantEndConditions: []string{"cast-while-on-the-ethereal-plane", "cast-on-plane-that-does-not-border-the-ethereal-plane"}, InstantEndExample: "outer-planes", BaseTargets: 1, AdditionalTargetsPerSlot: 3, SlotTargetKind: "willing-creatures-including-you", SlotTargetRangeFeet: 10},
	"retreat":      SpellRetreat{CastAction: "dash", RepeatAction: "dash", RepeatActionType: "bonus-action", RepeatWhile: "until-spell-ends"},
	"eyebite":      SpellEyebite{Appearance: "your-eyes-become-an-inky-void", TargetRangeFeet: 60, TargetOrigin: "within-60-feet-of-you", TargetVisibility: "you-can-see", Save: "wisdom", Failure: "affected-by-one-chosen-effect-for-the-duration", Options: []string{"asleep", "panicked", "sickened"}, RetargetAction: "magic-action", RetargetTiming: "each-of-your-turns-until-the-spell-ends", RetargetExclusion: "creature-that-succeeded-on-a-save-against-this-casting", AsleepCondition: "unconscious", AsleepWakesWhen: []string{"it-takes-any-damage", "another-creature-takes-an-action-to-shake-it-awake"}, PanickedCondition: "frightened", PanickedTurnAction: "take-dash-action-and-move-away-from-you-by-safest-and-shortest-route", PanickedEndDistanceFeet: 60, PanickedEndWhen: "moves-to-a-space-at-least-60-feet-away-where-it-cannot-see-you", SickenedCondition: "poisoned"},
	"fabricate":    SpellFabricate{Input: "raw-materials-you-can-see-within-range", Output: "products-of-the-same-material", Examples: []string{"wooden-bridge-from-a-clump-of-trees", "rope-from-a-patch-of-hemp", "clothes-from-flax-or-wool"}, ObjectSize: "large-or-smaller", ObjectCubeFeet: 10, ConnectedCubeCount: 8, ConnectedCubeFeet: 5, Quantity: "sufficient-quantity-of-material", MineralKinds: []string{"metal", "stone", "other-mineral-substance"}, MineralObjectSize: "medium-or-smaller", MineralCubeFeet: 5, Quality: "based-on-the-quality-of-the-raw-materials", CannotCreate: []string{"creatures", "magic-items"}, SkilledItems: "items-requiring-a-high-degree-of-skill-such-as-weapons-and-armor", SkilledItemsException: "proficiency-with-the-type-of-artisans-tools-used-to-craft-such-objects"},
	"faerie_fire":  SpellFaerieFire{CubeFeet: 20, ObjectsOutlined: "all-objects-in-the-cube", CreatureOutline: "each-creature-in-the-cube-that-fails-the-save", Save: "dexterity", OutlineColors: []string{"blue", "green", "violet"}, ColorChoice: "your-choice", LightKind: "dim-light", LightRadiusFeet: 10, LightSource: "objects-and-affected-creatures", BlockedCondition: "invisible", AttackRollEffect: "advantage-if-the-attacker-can-see-it", AttackRollTarget: "affected-creature-or-object"},
	"hound":        SpellHound{Form: "phantom-watchdog", Placement: "unoccupied-space-you-can-see-within-range", EndsWhen: []string{"the-duration-ends", "you-and-the-hound-are-more-than-300-feet-apart"}, SeparationFeet: 300, VisibleTo: "only-you", Traits: []string{"intangible", "invulnerable"}, BarkTrigger: "small-or-larger-creature-comes-within-30-feet-without-first-speaking-the-password", BarkTriggerFeet: 30, Password: "specified-by-you-when-you-cast-the-spell", BarkResult: "hound-starts-barking-loudly", TruesightFeet: 30, BiteTiming: "start-of-each-of-your-turns", BiteTarget: "one-enemy-within-5-feet-of-the-hound", BiteReachFeet: 5, BiteSave: "dexterity", BiteFailureDamageDiceCount: 4, BiteFailureDamageDieSides: 8, BiteFailureDamageType: "force", MoveAction: "magic-action", MoveTiming: "your-later-turns", MoveFeet: 30},
	"false_life":   SpellFalseLife{Recipient: "you", TemporaryHitPointsDiceCount: 2, TemporaryHitPointsDieSides: 4, TemporaryHitPointsBonus: 4, AdditionalTemporaryHitPointsPerSlot: 5},
	"fear":         SpellFear{ConeFeet: 30, Origin: "you", Save: "wisdom", FailureEffects: []string{"drop-whatever-it-is-holding", "frightened-condition-for-the-duration"}, Condition: "frightened", TurnAction: "take-dash-action-and-move-away-from-you-by-the-safest-route", TurnActionException: "nowhere-to-move", EndTurnTrigger: "ends-its-turn-in-a-space-without-line-of-sight-to-you", EndTurnSave: "wisdom", EndTurnSaveSuccess: "spell-ends-on-that-creature"},
	"feather_fall": SpellFeatherFall{TriggerRangeFeet: 60, Trigger: "you-or-a-creature-you-can-see-within-range-falls", MaxTargets: 5, TargetState: "falling-creatures-within-range", DescentFeetPerRound: 60, DescentUntil: "spell-ends", LandingResult: "no-fall-damage-and-spell-ends-for-that-creature"},
	"familiar":     SpellFamiliar{FormChoices: []string{"bat", "cat", "frog", "hawk", "lizard", "octopus", "owl", "rat", "raven", "spider", "weasel", "other-beast-with-challenge-rating-0"}, Appearance: "unoccupied-space-within-range", Statistics: "stat-block-of-the-chosen-form", CreatureTypeChoices: []string{"celestial", "fey", "fiend"}, CreatureTypeReplaces: "beast", Independence: "acts-independently-but-obeys-your-commands", TelepathyRangeFeet: 100, SeeThroughAction: "bonus-action", SeeThroughDuration: "until-the-start-of-your-next-turn", SeeThroughBenefit: "benefits-of-any-special-senses-it-has", TouchDelivery: "familiar-delivers-spells-with-a-range-of-touch", TouchDeliveryRangeFeet: 100, TouchDeliveryAction: "reaction", TouchDeliveryTiming: "when-you-cast-the-spell", CombatRole: "ally-to-you-and-your-allies", CombatInitiative: "rolls-its-own-initiative-and-acts-on-its-own-turn", AttackRestriction: "cannot-attack-but-can-take-other-actions", ZeroHitPoints: "disappears", ReappearsWhen: "you-cast-this-spell-again", DismissAction: "magic-action", DismissOptions: []string{"temporarily-to-a-pocket-dimension", "forever"}, RecallAction: "magic-action-while-temporarily-dismissed", RecallRangeFeet: 30, RecallSpace: "unoccupied-space", LeavesBehind: "anything-it-was-wearing-or-carrying", LeavesBehindWhen: []string{"drops-to-0-hit-points", "disappears-into-the-pocket-dimension"}, FamiliarLimit: "cannot-have-more-than-one-familiar-at-a-time", RecastResult: "existing-familiar-adopts-a-new-eligible-form"},
	"steed":        SpellSteed{Summons: "otherworldly-being-that-appears-as-a-loyal-steed", Placement: "unoccupied-space-of-your-choice-within-range", StatBlock: "otherworldly-steed", ReplacesExisting: "steed-from-this-spell-is-replaced-by-the-new-one", Appearance: "large-rideable-animal-of-your-choice", AppearanceExamples: []string{"horse", "camel", "dire-wolf", "elk"}, CreatureTypeChoices: []string{"celestial", "fey", "fiend"}, CreatureTypeTiming: "whenever-you-cast-the-spell", CombatRole: "ally-to-you-and-your-allies", Initiative: "shares-your-initiative-count", MountFunction: "controlled-mount-while-you-ride-it", IncapacitatedBehavior: "takes-its-turn-immediately-after-yours-and-acts-independently-focusing-on-protecting-you", DisappearsWhen: []string{"drops-to-0-hit-points", "you-die"}, LeavesBehind: "anything-it-was-wearing-or-carrying", RecastChoice: "summon-the-steed-that-disappeared-or-a-different-one", SlotLevelUse: "spell-slot-level-is-the-spell-level-in-the-stat-block"},
	"traps":        SpellTraps{Senses: "any-trap-within-range-and-line-of-sight", TrapDefinition: "object-or-mechanism-created-to-cause-damage-or-other-danger", SensedExamples: []string{"alarm-spell", "glyph-of-warding-spell", "mechanical-pit-trap"}, NotSensedExamples: []string{"natural-weakness-in-the-floor", "unstable-ceiling", "hidden-sinkhole"}, Reveals: "that-a-trap-is-present", DoesNotReveal: "trap-location", Learns: "general-nature-of-the-danger-posed-by-a-sensed-trap"},
}

var spell205Rows = map[string]spell205Row{
	"etherealness": {slug: "etherealness", name: "Etherealness", page: 128, level: 7, school: "conjuration", lists: []string{"bard", "cleric", "sorcerer", "warlock", "wizard"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: false, rangeKind: "self", rangeFeet: 0, target: SpellTarget{"single", 1, "self"}, duration: "timed", seconds: 28800, endTiming: "duration-ended", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"slot-targets", 7, 0, 0, []string{"targets"}}, effects: []spell205Effect{{"border", "plane.border-ethereal", "for-duration", 0, 0, "", 0}, {"interaction", "interaction.ethereal-only", "while-on-ethereal-plane", 0, 0, "", 0}, {"return", "return.shunt-force-damage", "when-spell-ends-in-occupied-space", 0, 0, "", 0}, {"ending", "ending.instant-on-unsupported-plane", "when-cast-on-ethereal-plane-or-plane-that-does-not-border-it", 0, 0, "", 0}, {"targets", "targets.willing-creatures-per-slot", "slot-level-above-7", 0, 0, "", 0}}, required: []string{"plane.border-ethereal", "interaction.ethereal-only", "return.shunt-force-damage", "ending.instant-on-unsupported-plane", "targets.willing-creatures-per-slot", "scaling.targets-per-slot"}},
	"retreat":      {slug: "expeditious-retreat", name: "Expeditious Retreat", page: 129, level: 1, school: "transmutation", lists: []string{"sorcerer", "warlock", "wizard"}, castingTime: "bonus-action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: true, rangeKind: "self", rangeFeet: 0, target: SpellTarget{"single", 1, "self"}, duration: "timed", seconds: 600, endTiming: "duration-or-concentration-ended", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"dash", "action.dash-on-cast", "immediate", 0, 0, "", 0}, {"repeat", "action.dash-as-bonus-action", "bonus-action-each-turn-until-spell-ends", 0, 0, "", 0}}, required: []string{"action.dash-on-cast", "action.dash-as-bonus-action", "concentration.source-duration"}},
	"eyebite":      {slug: "eyebite", name: "Eyebite", page: 129, level: 6, school: "necromancy", lists: []string{"bard", "sorcerer", "warlock", "wizard"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: true, rangeKind: "self", rangeFeet: 0, target: SpellTarget{"single", 1, "visible-creature-in-range"}, duration: "timed", seconds: 60, endTiming: "duration-or-concentration-ended", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"save", "save.eyebite-wisdom", "on-cast", 0, 0, "", 0}, {"asleep", "condition.eyebite-unconscious", "failed-save-if-asleep-chosen", 0, 0, "", 0}, {"panicked", "condition.eyebite-frightened-dash-away", "failed-save-if-panicked-chosen", 0, 0, "", 0}, {"sickened", "condition.eyebite-poisoned", "failed-save-if-sickened-chosen", 0, 0, "", 0}, {"retarget", "target.eyebite-magic-action-retarget", "each-of-your-turns-until-spell-ends", 0, 0, "", 0}}, required: []string{"save.eyebite-wisdom", "condition.eyebite-unconscious", "condition.eyebite-frightened-dash-away", "condition.eyebite-poisoned", "target.eyebite-magic-action-retarget", "concentration.source-duration"}},
	"fabricate":    {slug: "fabricate", name: "Fabricate", page: 129, level: 4, school: "transmutation", lists: []string{"wizard"}, castingTime: "10-minutes", castingSeconds: 600, triggers: []string{}, ritual: false, concentration: false, rangeKind: "feet", rangeFeet: 120, target: SpellTarget{"choice", 1, "visible-raw-materials-in-range"}, duration: "instantaneous", seconds: 0, endTiming: "instantaneous", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"convert", "craft.fabricate-same-material", "immediate", 0, 0, "", 0}, {"size", "limit.fabricate-object-size", "on-cast", 0, 0, "", 0}, {"restriction", "restriction.fabricate-skilled-items-need-tool-proficiency", "on-cast", 0, 0, "", 0}}, required: []string{"craft.fabricate-same-material", "limit.fabricate-object-size", "restriction.fabricate-skilled-items-need-tool-proficiency", "cast.ten-minutes"}},
	"faerie_fire":  {slug: "faerie-fire", name: "Faerie Fire", page: 129, level: 1, school: "evocation", lists: []string{"bard", "druid"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: true, rangeKind: "feet", rangeFeet: 60, target: SpellTarget{"cube", 0, "objects-and-creatures-in-cube"}, duration: "timed", seconds: 60, endTiming: "duration-or-concentration-ended", verbal: true, somatic: false, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"outline", "outline.faerie-fire-objects-and-failed-save-creatures", "on-cast", 0, 0, "", 0}, {"light", "light.dim-10-foot-radius", "for-duration", 0, 0, "", 0}, {"invisibility", "condition.faerie-fire-blocks-invisible", "for-duration", 0, 0, "", 0}, {"advantage", "roll.attack-advantage-if-attacker-sees-target", "for-duration", 0, 0, "", 0}}, required: []string{"outline.faerie-fire-objects-and-failed-save-creatures", "light.dim-10-foot-radius", "condition.faerie-fire-blocks-invisible", "roll.attack-advantage-if-attacker-sees-target", "save.dexterity", "concentration.source-duration"}},
	"hound":        {slug: "faithful-hound", name: "Faithful Hound", page: 129, level: 4, school: "conjuration", lists: []string{"wizard"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: false, rangeKind: "feet", rangeFeet: 30, target: SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, duration: "timed", seconds: 28800, endTiming: "duration-ended-or-caster-and-hound-more-than-300-feet-apart", verbal: true, somatic: true, material: "a silver whistle", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"hound", "summon.faithful-hound", "on-cast", 0, 0, "", 0}, {"alarm", "alarm.faithful-hound-bark", "small-or-larger-creature-within-30-feet-without-password", 0, 0, "", 0}, {"bite", "damage.faithful-hound-bite", "start-of-each-of-your-turns", 4, 8, "force", 0}, {"move", "movement.faithful-hound-magic-action", "magic-action-on-your-later-turns", 0, 0, "", 0}}, required: []string{"summon.faithful-hound", "alarm.faithful-hound-bark", "damage.faithful-hound-bite", "movement.faithful-hound-magic-action", "save.dexterity", "components.material"}},
	"false_life":   {slug: "false-life", name: "False Life", page: 129, level: 1, school: "necromancy", lists: []string{"sorcerer", "wizard"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: false, rangeKind: "self", rangeFeet: 0, target: SpellTarget{"single", 1, "self"}, duration: "instantaneous", seconds: 0, endTiming: "instantaneous", verbal: true, somatic: true, material: "a drop of alcohol", costCP: 0, consumed: false, scaling: SpellScaling{"slot-bonus", 1, 0, 0, []string{"temporary-hit-points"}}, effects: []spell205Effect{{"temporary-hit-points", "temporary-hit-points.false-life", "immediate", 2, 4, "", 4}}, required: []string{"temporary-hit-points.false-life", "scaling.bonus-per-slot", "components.material"}},
	"fear":         {slug: "fear", name: "Fear", page: 130, level: 3, school: "illusion", lists: []string{"bard", "sorcerer", "warlock", "wizard"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: true, rangeKind: "self", rangeFeet: 0, target: SpellTarget{"cone", 0, "creatures-in-cone"}, duration: "timed", seconds: 60, endTiming: "duration-or-concentration-ended", verbal: true, somatic: true, material: "a white feather", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"save", "save.fear-wisdom", "on-cast", 0, 0, "", 0}, {"drop", "drop.fear-held-items", "failed-save-on-cast", 0, 0, "", 0}, {"frightened", "condition.fear-frightened-dash-away", "failed-save-for-duration", 0, 0, "", 0}, {"recheck", "save.fear-end-turn-without-line-of-sight", "end-of-creature-turn-without-line-of-sight", 0, 0, "", 0}}, required: []string{"save.fear-wisdom", "drop.fear-held-items", "condition.fear-frightened-dash-away", "save.fear-end-turn-without-line-of-sight", "area.cone-30", "concentration.source-duration", "components.material"}},
	"feather_fall": {slug: "feather-fall", name: "Feather Fall", page: 130, level: 1, school: "transmutation", lists: []string{"bard", "sorcerer", "wizard"}, castingTime: "reaction", castingSeconds: 0, triggers: []string{"you-or-a-creature-you-can-see-within-60-feet-of-you-falls"}, ritual: false, concentration: false, rangeKind: "feet", rangeFeet: 60, target: SpellTarget{"up-to-count", 5, "falling-creatures-in-range"}, duration: "timed", seconds: 60, endTiming: "duration-ended-or-creature-lands", verbal: true, somatic: false, material: "a small feather or piece of down", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"descent", "descent.feather-fall-60-feet-per-round", "for-duration", 0, 0, "", 0}, {"landing", "landing.feather-fall-no-damage-ends-for-creature", "creature-lands-before-spell-ends", 0, 0, "", 0}}, required: []string{"descent.feather-fall-60-feet-per-round", "landing.feather-fall-no-damage-ends-for-creature", "cast.reaction-trigger", "components.material"}},
	"familiar":     {slug: "find-familiar", name: "Find Familiar", page: 130, level: 1, school: "conjuration", lists: []string{"wizard"}, castingTime: "1-hour", castingSeconds: 3600, triggers: []string{}, ritual: true, concentration: false, rangeKind: "feet", rangeFeet: 10, target: SpellTarget{"created-effect", 1, "unoccupied-space-in-range"}, duration: "instantaneous", seconds: 0, endTiming: "instantaneous", verbal: true, somatic: true, material: "burning incense worth 10+ GP, which the spell consumes", costCP: 1000, consumed: true, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"familiar", "summon.find-familiar", "on-cast", 0, 0, "", 0}, {"telepathy", "communication.familiar-telepathy-100-feet", "while-familiar-within-100-feet", 0, 0, "", 0}, {"senses", "sense.familiar-eyes-and-ears-bonus-action", "bonus-action-until-start-of-your-next-turn", 0, 0, "", 0}, {"touch", "delivery.familiar-touch-spells", "reaction-when-you-cast-a-touch-spell-with-familiar-within-100-feet", 0, 0, "", 0}, {"combat", "combat.familiar-ally-own-initiative", "while-familiar-exists", 0, 0, "", 0}, {"disappearance", "state.familiar-disappears-at-zero-hit-points", "familiar-drops-to-0-hit-points", 0, 0, "", 0}, {"dismissal", "action.familiar-dismiss-and-recall", "magic-action", 0, 0, "", 0}, {"limit", "limit.one-familiar-only", "when-cast-while-you-have-a-familiar", 0, 0, "", 0}}, required: []string{"summon.find-familiar", "communication.familiar-telepathy-100-feet", "sense.familiar-eyes-and-ears-bonus-action", "delivery.familiar-touch-spells", "combat.familiar-ally-own-initiative", "state.familiar-disappears-at-zero-hit-points", "action.familiar-dismiss-and-recall", "limit.one-familiar-only", "cast.one-hour", "cast.ritual", "components.consumed-priced-material", "components.material"}},
	"steed":        {slug: "find-steed", name: "Find Steed", page: 130, level: 2, school: "conjuration", lists: []string{"paladin"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: false, rangeKind: "feet", rangeFeet: 30, target: SpellTarget{"created-effect", 1, "unoccupied-space-in-range"}, duration: "instantaneous", seconds: 0, endTiming: "instantaneous", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"steed", "summon.find-steed", "on-cast", 0, 0, "", 0}, {"replacement", "limit.one-steed-replaced-by-new-one", "when-cast-while-you-have-a-steed", 0, 0, "", 0}, {"combat", "combat.steed-ally-shared-initiative-controlled-mount", "while-steed-exists", 0, 0, "", 0}, {"disappearance", "state.steed-disappears-at-zero-hit-points-or-caster-death", "steed-drops-to-0-hit-points-or-you-die", 0, 0, "", 0}, {"slot-level", "stat-block.slot-level-as-spell-level", "on-cast", 0, 0, "", 0}}, required: []string{"summon.find-steed", "limit.one-steed-replaced-by-new-one", "combat.steed-ally-shared-initiative-controlled-mount", "state.steed-disappears-at-zero-hit-points-or-caster-death", "stat-block.otherworldly-steed", "stat-block.slot-level-as-spell-level"}},
	"traps":        {slug: "find-traps", name: "Find Traps", page: 131, level: 2, school: "divination", lists: []string{"cleric", "druid", "ranger"}, castingTime: "action", castingSeconds: 0, triggers: []string{}, ritual: false, concentration: false, rangeKind: "feet", rangeFeet: 120, target: SpellTarget{"selected-objects", 0, "traps-in-line-of-sight-in-range"}, duration: "instantaneous", seconds: 0, endTiming: "instantaneous", verbal: true, somatic: true, material: "", costCP: 0, consumed: false, scaling: SpellScaling{"none", 0, 0, 0, []string{}}, effects: []spell205Effect{{"presence", "sense.find-traps-presence", "immediate", 0, 0, "", 0}, {"nature", "sense.find-traps-general-nature-of-danger", "immediate", 0, 0, "", 0}}, required: []string{"sense.find-traps-presence", "sense.find-traps-general-nature-of-danger"}},
}

func validSpell205Target(t SpellTarget) bool {
	switch t.Shape {
	case "choice":
		return t.Count == 1 && t.Eligibility == "visible-raw-materials-in-range"
	case "cube":
		return t.Count == 0 && t.Eligibility == "objects-and-creatures-in-cube"
	case "up-to-count":
		return t.Count == 5 && t.Eligibility == "falling-creatures-in-range"
	case "selected-objects":
		return t.Count == 0 && t.Eligibility == "traps-in-line-of-sight-in-range"
	}
	return false
}

// spell205TargetScaling returns Etherealness' typed base and per-slot target counts.
func spell205TargetScaling(d *SpellSourceDetails) (base, additional int, ok bool) {
	if d.Etherealness != nil {
		return d.Etherealness.BaseTargets, d.Etherealness.AdditionalTargetsPerSlot, true
	}
	return 0, 0, false
}

type spell205Effect struct {
	id, mechanic, timing string
	dice, sides          int
	damage               string
	bonus                int
}

// spell205Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD page.
type spell205Row struct {
	slug, name            string
	page, level           int
	school                string
	lists                 []string
	castingTime           string
	castingSeconds        int
	triggers              []string
	ritual, concentration bool
	rangeKind             string
	rangeFeet             int
	target                SpellTarget
	duration              string
	seconds               int
	endTiming             string
	verbal, somatic       bool
	material              string
	costCP                int
	consumed              bool
	scaling               SpellScaling
	effects               []spell205Effect
	required              []string
}

func (s SpellRecord) validateSpell205() (bool, error) {
	expected := ""
	for key, row := range spell205Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 205 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"etherealness": d.Etherealness, "retreat": d.Retreat, "eyebite": d.Eyebite, "fabricate": d.Fabricate,
		"faerie_fire": d.FaerieFire, "hound": d.Hound, "false_life": d.FalseLife, "fear": d.Fear,
		"feather_fall": d.FeatherFall, "familiar": d.Familiar, "steed": d.Steed, "traps": d.Traps,
	}
	count := 0
	var selected any
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 205 payload")
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return true, err
		}
		if err := requireSpell198Fields(raw, reflect.TypeOf(value).Elem()); err != nil {
			return true, err
		}
		selected = reflect.ValueOf(value).Elem().Interface()
	}
	if expected == "" && count == 0 {
		return false, nil
	}
	row := spell205Rows[expected]
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != row.castingSeconds || d.CastingTriggers == nil || !slices.Equal(d.CastingTriggers, row.triggers) {
		return true, fmt.Errorf("invalid batch 205 source details")
	}
	if err := s.validateSpell205Header(row); err != nil {
		return true, err
	}
	if !reflect.DeepEqual(selected, spell205Pins[expected]) {
		return true, fmt.Errorf("invalid batch 205 %s payload", expected)
	}
	return true, nil
}

func (s SpellRecord) validateSpell205Header(row spell205Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != row.castingTime || flags(s.Ritual) != row.ritual || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != row.duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || flags(c.Verbal) != row.verbal || flags(c.Somatic) != row.somatic {
		return fmt.Errorf("batch 205 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 205 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed != row.consumed {
		return fmt.Errorf("batch 205 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 205 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 205 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.Timing != want.timing || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != want.bonus ||
			e.AddCastingModifier == nil || *e.AddCastingModifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 205 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 205 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 205 slot levels differ from source")
	}
	return nil
}
