package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// SpellFindPath retains Find the Path's typed SRD source facts without implementing them.
type SpellFindPath struct {
	Destination           string   `json:"destination"`
	FamiliarityRequired   bool     `json:"familiarity_required"`
	RouteSensed           string   `json:"route_sensed"`
	FailingDestinations   []string `json:"failing_destinations"`
	KnownWhileOnSamePlane []string `json:"known_while_on_same_plane"`
	PlaneRequirement      string   `json:"plane_requirement"`
	PathChoice            string   `json:"path_choice"`
}

// SpellFingerOfDeath retains Finger of Death's typed SRD source facts without implementing them.
type SpellFingerOfDeath struct {
	TargetSight     string `json:"target_sight"`
	Save            string `json:"save"`
	DamageType      string `json:"damage_type"`
	DiceCount       int    `json:"dice_count"`
	DieSides        int    `json:"die_sides"`
	FlatBonus       int    `json:"flat_bonus"`
	SuccessDamage   string `json:"success_damage"`
	RaisedCreature  string `json:"raised_creature"`
	RisesWhen       string `json:"rises_when"`
	RisesAs         string `json:"rises_as"`
	ZombieReference string `json:"zombie_reference"`
	ZombieObeys     string `json:"zombie_obeys"`
}

// SpellFireBolt retains Fire Bolt's typed SRD source facts without implementing them.
type SpellFireBolt struct {
	Attack     string               `json:"attack"`
	TargetKind string               `json:"target_kind"`
	DamageType string               `json:"damage_type"`
	Ignites    string               `json:"ignites"`
	Upgrades   []SpellSplashUpgrade `json:"upgrades"`
}

// SpellFireShield retains Fire Shield's typed SRD source facts without implementing them.
type SpellFireShield struct {
	BrightLightRadiusFeet int      `json:"bright_light_radius_feet"`
	DimLightExtraFeet     int      `json:"dim_light_extra_feet"`
	ShieldChoices         []string `json:"shield_choices"`
	WarmResistance        string   `json:"warm_resistance"`
	ChillResistance       string   `json:"chill_resistance"`
	RetaliationTrigger    string   `json:"retaliation_trigger"`
	RetaliationRangeFeet  int      `json:"retaliation_range_feet"`
	RetaliationDiceCount  int      `json:"retaliation_dice_count"`
	RetaliationDieSides   int      `json:"retaliation_die_sides"`
	WarmDamageType        string   `json:"warm_damage_type"`
	ChillDamageType       string   `json:"chill_damage_type"`
}

// SpellFireStorm retains Fire Storm's typed SRD source facts without implementing them.
type SpellFireStorm struct {
	MaxCubes      int    `json:"max_cubes"`
	CubeSizeFeet  int    `json:"cube_size_feet"`
	Arrangement   string `json:"arrangement"`
	Contiguity    string `json:"contiguity"`
	Save          string `json:"save"`
	DamageType    string `json:"damage_type"`
	DiceCount     int    `json:"dice_count"`
	DieSides      int    `json:"die_sides"`
	SuccessDamage string `json:"success_damage"`
	Ignites       string `json:"ignites"`
}

// SpellFireballBurst retains Fireball's typed SRD source facts without implementing them.
type SpellFireballBurst struct {
	Origin         string `json:"origin"`
	Center         string `json:"center"`
	RadiusFeet     int    `json:"radius_feet"`
	Save           string `json:"save"`
	DamageType     string `json:"damage_type"`
	DiceCount      int    `json:"dice_count"`
	DieSides       int    `json:"die_sides"`
	SuccessDamage  string `json:"success_damage"`
	Ignites        string `json:"ignites"`
	SlotAboveLevel int    `json:"slot_above_level"`
	SlotDiceCount  int    `json:"slot_dice_count"`
	SlotDieSides   int    `json:"slot_die_sides"`
}

// SpellFlameBlade retains Flame Blade's typed SRD source facts without implementing them.
type SpellFlameBlade struct {
	BladeForm             string `json:"blade_form"`
	BladeShape            string `json:"blade_shape"`
	Release               string `json:"release"`
	ReevokeAction         string `json:"reevoke_action"`
	AttackAction          string `json:"attack_action"`
	Attack                string `json:"attack"`
	DamageType            string `json:"damage_type"`
	DiceCount             int    `json:"dice_count"`
	DieSides              int    `json:"die_sides"`
	DamageModifier        string `json:"damage_modifier"`
	BrightLightRadiusFeet int    `json:"bright_light_radius_feet"`
	DimLightExtraFeet     int    `json:"dim_light_extra_feet"`
	SlotAboveLevel        int    `json:"slot_above_level"`
	SlotDiceCount         int    `json:"slot_dice_count"`
	SlotDieSides          int    `json:"slot_die_sides"`
}

// SpellFlameStrike retains Flame Strike's typed SRD source facts without implementing them.
type SpellFlameStrike struct {
	Origin           string   `json:"origin"`
	Center           string   `json:"center"`
	RadiusFeet       int      `json:"radius_feet"`
	HeightFeet       int      `json:"height_feet"`
	Save             string   `json:"save"`
	FireDiceCount    int      `json:"fire_dice_count"`
	FireDieSides     int      `json:"fire_die_sides"`
	RadiantDiceCount int      `json:"radiant_dice_count"`
	RadiantDieSides  int      `json:"radiant_die_sides"`
	SuccessDamage    string   `json:"success_damage"`
	SlotAboveLevel   int      `json:"slot_above_level"`
	SlotDiceCount    int      `json:"slot_dice_count"`
	SlotDieSides     int      `json:"slot_die_sides"`
	SlotAppliesTo    []string `json:"slot_applies_to"`
}

// SpellFlamingSphere retains Flaming Sphere's typed SRD source facts without implementing them.
type SpellFlamingSphere struct {
	DiameterFeet          int    `json:"diameter_feet"`
	Space                 string `json:"space"`
	DamageTriggerFeet     int    `json:"damage_trigger_feet"`
	DamageTrigger         string `json:"damage_trigger"`
	Save                  string `json:"save"`
	DamageType            string `json:"damage_type"`
	DiceCount             int    `json:"dice_count"`
	DieSides              int    `json:"die_sides"`
	SuccessDamage         string `json:"success_damage"`
	MoveAction            string `json:"move_action"`
	MoveFeet              int    `json:"move_feet"`
	MoveMode              string `json:"move_mode"`
	MoveIntoCreature      string `json:"move_into_creature"`
	BarrierHeightFeet     int    `json:"barrier_height_feet"`
	PitWidthFeet          int    `json:"pit_width_feet"`
	Ignites               string `json:"ignites"`
	BrightLightRadiusFeet int    `json:"bright_light_radius_feet"`
	DimLightExtraFeet     int    `json:"dim_light_extra_feet"`
	SlotAboveLevel        int    `json:"slot_above_level"`
	SlotDiceCount         int    `json:"slot_dice_count"`
	SlotDieSides          int    `json:"slot_die_sides"`
}

// SpellFleshToStone retains Flesh to Stone's typed SRD source facts without implementing them.
type SpellFleshToStone struct {
	TargetSight                    string `json:"target_sight"`
	Save                           string `json:"save"`
	FailureCondition               string `json:"failure_condition"`
	SuccessEffect                  string `json:"success_effect"`
	ConstructsAutomaticallySucceed bool   `json:"constructs_automatically_succeed"`
	RepeatSaveTiming               string `json:"repeat_save_timing"`
	SuccessesToEnd                 int    `json:"successes_to_end"`
	FailuresToPetrify              int    `json:"failures_to_petrify"`
	ResultsConsecutive             bool   `json:"results_consecutive"`
	PetrifiedCondition             string `json:"petrified_condition"`
	FullDurationResult             string `json:"full_duration_result"`
}

// SpellFloatingDisk retains Floating Disk's typed SRD source facts without implementing them.
type SpellFloatingDisk struct {
	Plane                    string `json:"plane"`
	DiameterFeet             int    `json:"diameter_feet"`
	ThicknessInches          int    `json:"thickness_inches"`
	HeightAboveGroundFeet    int    `json:"height_above_ground_feet"`
	Space                    string `json:"space"`
	CapacityPounds           int    `json:"capacity_pounds"`
	OverloadResult           string `json:"overload_result"`
	ImmobileWithinFeet       int    `json:"immobile_within_feet"`
	FollowDistanceFeet       int    `json:"follow_distance_feet"`
	Terrain                  string `json:"terrain"`
	ElevationChangeLimitFeet int    `json:"elevation_change_limit_feet"`
	SpellEndsBeyondFeet      int    `json:"spell_ends_beyond_feet"`
}

// SpellFly retains Fly's typed SRD source facts without implementing them.
type SpellFly struct {
	Target                   string `json:"target"`
	FlySpeedFeet             int    `json:"fly_speed_feet"`
	Hover                    bool   `json:"hover"`
	Ending                   string `json:"ending"`
	BaseTargets              int    `json:"base_targets"`
	AdditionalTargetsPerSlot int    `json:"additional_targets_per_slot"`
}

func spell206Templates() map[string]any {
	return map[string]any{
		"find_path":       SpellFindPath{},
		"finger_of_death": SpellFingerOfDeath{},
		"fire_bolt":       SpellFireBolt{},
		"fire_shield":     SpellFireShield{},
		"fire_storm":      SpellFireStorm{},
		"fireball_burst":  SpellFireballBurst{},
		"flame_blade":     SpellFlameBlade{},
		"flame_strike":    SpellFlameStrike{},
		"flaming_sphere":  SpellFlamingSphere{},
		"flesh_to_stone":  SpellFleshToStone{},
		"floating_disk":   SpellFloatingDisk{},
		"fly":             SpellFly{},
	}
}

// validSpell206Target lists the target shapes first used by this batch. They
// are source descriptions only and select no runtime targeting behavior.
func validSpell206Target(t SpellTarget) bool {
	switch t.Shape {
	case "choice":
		return t.Count == 1 && t.Eligibility == "creature-or-object-in-range"
	case "cubes":
		return t.Count == 10 && t.Eligibility == "creatures-in-contiguous-cubes"
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "creatures-in-cylinder"
	case "created-effect":
		return t.Count == 1 && t.Eligibility == "unoccupied-ground-space-in-range"
	}
	return false
}

type spell206Effect struct {
	id, mechanic, timing string
	dice, sides          int
	damage               string
	bonus                int
	modifier             bool
}

// spell206Row pins every header value, scaling row, effect and dependency a
// batch record must carry, by equality with the printed SRD page.
type spell206Row struct {
	slug, name            string
	page, level           int
	school                string
	lists                 []string
	casting               string
	ritual, concentration bool
	rangeKind             string
	rangeFeet             int
	target                SpellTarget
	duration              string
	seconds               int
	endTiming             string
	castingSeconds        int
	material              string
	costCP                int
	scaling               SpellScaling
	effects               []spell206Effect
	required              []string
}

var spell206Rows = map[string]spell206Row{
	"find_path": {"find-the-path", "Find the Path", 131, 6, "divination", []string{"bard", "cleric", "druid"}, "1-minute", false, true, "self", 0,
		SpellTarget{"single", 1, "self"}, "timed", 86400, "duration-or-concentration-ended", 60, "a set of divination tools—such as cards or runes—worth 100+ GP", 10000,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"route", "sense.find-the-path-route", "duration-on-same-plane", 0, 0, "", 0, false}},
		[]string{"sense.find-the-path-route", "sense.find-the-path-invalid-destinations", "cast.one-minute", "concentration.source-duration", "components.material"}},
	"finger_of_death": {"finger-of-death", "Finger of Death", 131, 7, "necromancy", []string{"sorcerer", "warlock", "wizard"}, "action", false, false, "feet", 60,
		SpellTarget{"single", 1, "visible-creature-in-range"}, "instantaneous", 0, "instantaneous", 0, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"damage", "damage.finger-of-death", "failed-save-full-damage-success-half", 7, 8, "necrotic", 30, false}, {"zombie", "summon.finger-of-death-zombie", "humanoid-killed-rises-at-start-of-caster-next-turn", 0, 0, "", 0, false}},
		[]string{"damage.finger-of-death", "damage.necrotic", "save.constitution-half", "summon.finger-of-death-zombie"}},
	"fire_bolt": {"fire-bolt", "Fire Bolt", 132, 0, "evocation", []string{"sorcerer", "wizard"}, "action", false, false, "feet", 120,
		SpellTarget{"choice", 1, "creature-or-object-in-range"}, "instantaneous", 0, "instantaneous", 0, "", 0,
		SpellScaling{"cantrip-upgrade", 0, 1, 10, []string{"damage"}},
		[]spell206Effect{{"damage", "damage.fire-bolt", "ranged-spell-attack-hit", 1, 10, "fire", 0, false}, {"ignite", "ignite.flammable-objects", "flammable-object-hit-not-worn-or-carried", 0, 0, "", 0, false}},
		[]string{"damage.fire-bolt", "damage.fire", "attack.ranged-spell", "ignite.flammable-objects", "scaling.cantrip-upgrade"}},
	"fire_shield": {"fire-shield", "Fire Shield", 132, 4, "evocation", []string{"druid", "sorcerer", "wizard"}, "action", false, false, "self", 0,
		SpellTarget{"single", 1, "self"}, "timed", 600, "duration-ended", 0, "a bit of phosphorus or a firefly", 0,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"light", "light.fire-shield", "duration", 0, 0, "", 0, false}, {"resistance", "resistance.fire-shield", "duration-per-chosen-shield", 0, 0, "", 0, false}, {"warm-retaliation", "damage.fire-shield-warm-retaliation", "creature-within-5-feet-hits-you-with-melee-attack-roll", 2, 8, "fire", 0, false}, {"chill-retaliation", "damage.fire-shield-chill-retaliation", "creature-within-5-feet-hits-you-with-melee-attack-roll", 2, 8, "cold", 0, false}},
		[]string{"light.fire-shield", "resistance.fire-shield", "damage.fire-shield-warm-retaliation", "damage.fire-shield-chill-retaliation", "components.material"}},
	"fire_storm": {"fire-storm", "Fire Storm", 132, 7, "evocation", []string{"cleric", "druid", "sorcerer"}, "action", false, false, "feet", 150,
		SpellTarget{"cubes", 10, "creatures-in-contiguous-cubes"}, "instantaneous", 0, "instantaneous", 0, "", 0,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"damage", "damage.fire-storm", "failed-save-full-damage-success-half", 7, 10, "fire", 0, false}, {"ignite", "ignite.flammable-objects", "flammable-objects-in-area-not-worn-or-carried", 0, 0, "", 0, false}},
		[]string{"damage.fire-storm", "damage.fire", "area.fire-storm-cubes", "save.dexterity-half", "ignite.flammable-objects"}},
	"fireball_burst": {"fireball", "Fireball", 131, 3, "evocation", []string{"sorcerer", "wizard"}, "action", false, false, "feet", 150,
		SpellTarget{"sphere", 0, "creatures-in-sphere"}, "instantaneous", 0, "instantaneous", 0, "a ball of bat guano and sulfur", 0,
		SpellScaling{"slot", 3, 1, 6, []string{"damage"}},
		[]spell206Effect{{"damage", "damage.fireball", "failed-save-full-damage-success-half", 8, 6, "fire", 0, false}, {"ignite", "ignite.flammable-objects", "flammable-objects-in-area-not-worn-or-carried", 0, 0, "", 0, false}},
		[]string{"damage.fireball", "damage.fire", "area.sphere-20-radius", "save.dexterity-half", "ignite.flammable-objects", "scaling.per-slot", "components.material"}},
	"flame_blade": {"flame-blade", "Flame Blade", 132, 2, "evocation", []string{"druid", "sorcerer"}, "bonus-action", false, true, "self", 0,
		SpellTarget{"single", 1, "self"}, "timed", 600, "duration-or-concentration-ended", 0, "a sumac leaf", 0,
		SpellScaling{"slot", 2, 1, 6, []string{"damage"}},
		[]spell206Effect{{"blade", "weapon.flame-blade", "duration", 0, 0, "", 0, false}, {"damage", "damage.flame-blade", "melee-spell-attack-hit", 3, 6, "fire", 0, true}, {"light", "light.flame-blade", "duration", 0, 0, "", 0, false}},
		[]string{"weapon.flame-blade", "damage.flame-blade", "damage.fire", "attack.melee-spell", "light.flame-blade", "scaling.per-slot", "concentration.source-duration", "components.material"}},
	"flame_strike": {"flame-strike", "Flame Strike", 132, 5, "evocation", []string{"cleric"}, "action", false, false, "feet", 60,
		SpellTarget{"cylinder", 0, "creatures-in-cylinder"}, "instantaneous", 0, "instantaneous", 0, "a pinch of sulfur", 0,
		SpellScaling{"slot", 5, 1, 6, []string{"fire", "radiant"}},
		[]spell206Effect{{"fire", "damage.flame-strike-fire", "failed-save-full-damage-success-half", 5, 6, "fire", 0, false}, {"radiant", "damage.flame-strike-radiant", "failed-save-full-damage-success-half", 5, 6, "radiant", 0, false}},
		[]string{"damage.flame-strike-fire", "damage.flame-strike-radiant", "area.cylinder-10-radius-40-high", "save.dexterity-half", "scaling.per-slot", "components.material"}},
	"flaming_sphere": {"flaming-sphere", "Flaming Sphere", 132, 2, "conjuration", []string{"druid", "sorcerer", "wizard"}, "action", false, true, "feet", 60,
		SpellTarget{"created-effect", 1, "unoccupied-ground-space-in-range"}, "timed", 60, "duration-or-concentration-ended", 0, "a ball of wax", 0,
		SpellScaling{"slot", 2, 1, 6, []string{"damage"}},
		[]spell206Effect{{"sphere", "object.flaming-sphere", "duration", 0, 0, "", 0, false}, {"damage", "damage.flaming-sphere", "creature-ends-turn-within-5-feet-or-sphere-moves-into-space", 2, 6, "fire", 0, false}, {"move", "movement.flaming-sphere-bonus-action", "bonus-action", 0, 0, "", 0, false}, {"ignite", "ignite.flammable-objects", "flammable-objects-not-worn-or-carried-touched-by-sphere", 0, 0, "", 0, false}, {"light", "light.flaming-sphere", "duration", 0, 0, "", 0, false}},
		[]string{"object.flaming-sphere", "damage.flaming-sphere", "damage.fire", "save.dexterity-half", "movement.flaming-sphere-bonus-action", "ignite.flammable-objects", "light.flaming-sphere", "scaling.per-slot", "concentration.source-duration", "components.material"}},
	"flesh_to_stone": {"flesh-to-stone", "Flesh to Stone", 133, 6, "transmutation", []string{"druid", "sorcerer", "wizard"}, "action", false, true, "feet", 60,
		SpellTarget{"single", 1, "visible-creature-in-range"}, "timed", 60, "duration-or-concentration-ended", 0, "a cockatrice feather", 0,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"restrain", "condition.flesh-to-stone-restrained", "failed-save-for-duration", 0, 0, "", 0, false}, {"slow", "speed.flesh-to-stone-zero", "successful-save-until-start-of-caster-next-turn", 0, 0, "", 0, false}, {"petrify", "condition.flesh-to-stone-petrified", "three-failed-saves", 0, 0, "", 0, false}},
		[]string{"condition.flesh-to-stone-restrained", "speed.flesh-to-stone-zero", "condition.flesh-to-stone-petrified", "save.constitution-repeating", "construct.automatic-success", "concentration.source-duration", "components.material"}},
	"floating_disk": {"floating-disk", "Floating Disk", 133, 1, "conjuration", []string{"wizard"}, "action", true, false, "feet", 30,
		SpellTarget{"created-effect", 1, "visible-unoccupied-space-in-range"}, "timed", 3600, "duration-ended", 0, "a drop of mercury", 0,
		SpellScaling{"none", 0, 0, 0, []string{}},
		[]spell206Effect{{"disk", "object.floating-disk", "duration", 0, 0, "", 0, false}, {"follow", "movement.floating-disk-follow", "caster-moves-more-than-20-feet-away", 0, 0, "", 0, false}, {"ending", "ending.floating-disk-overload-or-distance", "weight-over-500-pounds-or-caster-more-than-100-feet-away", 0, 0, "", 0, false}},
		[]string{"object.floating-disk", "movement.floating-disk-follow", "ending.floating-disk-overload-or-distance", "cast.ritual", "components.material"}},
	"fly": {"fly", "Fly", 133, 3, "transmutation", []string{"sorcerer", "warlock", "wizard"}, "action", false, true, "touch", 0,
		SpellTarget{"single", 1, "willing-creature-touch"}, "timed", 600, "duration-or-concentration-ended", 0, "a feather", 0,
		SpellScaling{"slot-targets", 3, 0, 0, []string{"fly"}},
		[]spell206Effect{{"fly", "speed.fly-60-hover", "duration", 0, 0, "", 0, false}},
		[]string{"speed.fly-60-hover", "scaling.targets-per-slot", "concentration.source-duration", "components.material"}},
}

// spell206Payloads holds each family's printed typed values, so a positive but
// wrong value in a payload is rejected by equality.
var spell206Payloads = map[string]string{
	"find_path":       `{"destination":"location-you-name","familiarity_required":true,"route_sensed":"most-direct-physical-route","failing_destinations":["destination-on-another-plane","moving-destination","unspecific-destination"],"known_while_on_same_plane":["distance","direction"],"plane_requirement":"same-plane-as-destination","path_choice":"know-which-path-is-most-direct-at-each-choice"}`,
	"finger_of_death": `{"target_sight":"creature-you-can-see-in-range","save":"constitution","damage_type":"necrotic","dice_count":7,"die_sides":8,"flat_bonus":30,"success_damage":"half-damage","raised_creature":"humanoid-killed-by-spell","rises_when":"start-of-caster-next-turn","rises_as":"zombie","zombie_reference":"monsters-section","zombie_obeys":"caster-verbal-orders"}`,
	"fire_bolt":       `{"attack":"ranged-spell-attack","target_kind":"creature-or-object-in-range","damage_type":"fire","ignites":"flammable-object-hit-that-is-not-worn-or-carried","upgrades":[{"character_level":5,"dice_count":2,"die_sides":10},{"character_level":11,"dice_count":3,"die_sides":10},{"character_level":17,"dice_count":4,"die_sides":10}]}`,
	"fire_shield":     `{"bright_light_radius_feet":10,"dim_light_extra_feet":10,"shield_choices":["warm-shield","chill-shield"],"warm_resistance":"cold","chill_resistance":"fire","retaliation_trigger":"creature-within-5-feet-hits-you-with-melee-attack-roll","retaliation_range_feet":5,"retaliation_dice_count":2,"retaliation_die_sides":8,"warm_damage_type":"fire","chill_damage_type":"cold"}`,
	"fire_storm":      `{"max_cubes":10,"cube_size_feet":10,"arrangement":"as-you-like","contiguity":"each-cube-contiguous-with-at-least-one-other-cube","save":"dexterity","damage_type":"fire","dice_count":7,"die_sides":10,"success_damage":"half-damage","ignites":"flammable-objects-in-area-not-worn-or-carried"}`,
	"fireball_burst":  `{"origin":"bright-streak-from-caster-to-chosen-point-in-range","center":"point-in-range","radius_feet":20,"save":"dexterity","damage_type":"fire","dice_count":8,"die_sides":6,"success_damage":"half-damage","ignites":"flammable-objects-in-area-not-worn-or-carried","slot_above_level":3,"slot_dice_count":1,"slot_die_sides":6}`,
	"flame_blade":     `{"blade_form":"fiery-blade-in-free-hand","blade_shape":"similar-in-size-and-shape-to-scimitar","release":"disappears-if-you-let-go","reevoke_action":"bonus-action","attack_action":"magic-action","attack":"melee-spell-attack","damage_type":"fire","dice_count":3,"die_sides":6,"damage_modifier":"spellcasting-ability-modifier","bright_light_radius_feet":10,"dim_light_extra_feet":10,"slot_above_level":2,"slot_dice_count":1,"slot_die_sides":6}`,
	"flame_strike":    `{"origin":"vertical-column-of-fire-from-above","center":"point-in-range","radius_feet":10,"height_feet":40,"save":"dexterity","fire_dice_count":5,"fire_die_sides":6,"radiant_dice_count":5,"radiant_die_sides":6,"success_damage":"half-damage","slot_above_level":5,"slot_dice_count":1,"slot_die_sides":6,"slot_applies_to":["fire","radiant"]}`,
	"flaming_sphere":  `{"diameter_feet":5,"space":"unoccupied-space-on-the-ground-in-range","damage_trigger_feet":5,"damage_trigger":"creature-ends-its-turn-within-distance-of-sphere","save":"dexterity","damage_type":"fire","dice_count":2,"die_sides":6,"success_damage":"half-damage","move_action":"bonus-action","move_feet":30,"move_mode":"rolling-along-the-ground","move_into_creature":"creature-makes-the-save-and-sphere-stops-moving-for-the-turn","barrier_height_feet":5,"pit_width_feet":10,"ignites":"flammable-objects-not-worn-or-carried-touched-by-sphere","bright_light_radius_feet":20,"dim_light_extra_feet":20,"slot_above_level":2,"slot_dice_count":1,"slot_die_sides":6}`,
	"flesh_to_stone":  `{"target_sight":"one-creature-you-can-see-in-range","save":"constitution","failure_condition":"restrained","success_effect":"speed-0-until-start-of-caster-next-turn","constructs_automatically_succeed":true,"repeat_save_timing":"end-of-each-restrained-target-turn","successes_to_end":3,"failures_to_petrify":3,"results_consecutive":false,"petrified_condition":"petrified-for-the-duration","full_duration_result":"petrified-until-ended-by-greater-restoration-or-similar-magic"}`,
	"floating_disk":   `{"plane":"circular-horizontal-plane-of-force","diameter_feet":3,"thickness_inches":1,"height_above_ground_feet":3,"space":"unoccupied-space-you-can-see-in-range","capacity_pounds":500,"overload_result":"spell-ends-and-everything-on-the-disk-falls","immobile_within_feet":20,"follow_distance_feet":20,"terrain":"uneven-terrain-stairs-and-slopes","elevation_change_limit_feet":10,"spell_ends_beyond_feet":100}`,
	"fly":             `{"target":"willing-creature-touched","fly_speed_feet":60,"hover":true,"ending":"falls-if-still-aloft-unless-it-can-stop-the-fall","base_targets":1,"additional_targets_per_slot":1}`,
}

func (s SpellRecord) validateSpell206() (bool, error) {
	expected := ""
	for key, row := range spell206Rows {
		if s.ID == prefix+"spells/"+row.slug {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" && s.Scope == "complete-source-record" {
			return true, fmt.Errorf("batch 206 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"find_path":       d.FindPath,
		"finger_of_death": d.FingerOfDeath,
		"fire_bolt":       d.FireBolt,
		"fire_shield":     d.FireShield,
		"fire_storm":      d.FireStorm,
		"fireball_burst":  d.FireballBurst,
		"flame_blade":     d.FlameBlade,
		"flame_strike":    d.FlameStrike,
		"flaming_sphere":  d.FlamingSphere,
		"flesh_to_stone":  d.FleshToStone,
		"floating_disk":   d.FloatingDisk,
		"fly":             d.Fly,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 206 payload")
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return true, err
		}
		if err := requireSpell198Fields(raw, reflect.TypeOf(value).Elem()); err != nil {
			return true, err
		}
		want := reflect.New(reflect.TypeOf(value).Elem()).Interface()
		if err := json.Unmarshal([]byte(spell206Payloads[key]), want); err != nil {
			return true, err
		}
		if !reflect.DeepEqual(value, want) {
			return true, fmt.Errorf("invalid batch 206 %s payload", key)
		}
	}
	if expected == "" && count == 0 {
		return false, nil
	}
	row := spell206Rows[expected]
	if count != 1 || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != row.castingSeconds || d.CastingTriggers == nil || len(d.CastingTriggers) != 0 {
		return true, fmt.Errorf("invalid batch 206 source details")
	}
	return true, s.validateSpell206Header(row)
}

func (s SpellRecord) validateSpell206Header(row spell206Row) error {
	flags := func(p *bool) bool { return p != nil && *p }
	c := s.Components
	if s.Name != row.name || s.SourcePage != row.page || s.Level != row.level || s.School != row.school || !slices.Equal(s.Lists, row.lists) ||
		s.CastingTime != row.casting || flags(s.Ritual) != row.ritual || flags(s.Concentration) != row.concentration ||
		s.Range != (SpellRange{row.rangeKind, row.rangeFeet}) || s.Target != row.target || s.Duration != row.duration || s.DurationSeconds != row.seconds ||
		s.SourceDetails.EndTiming != row.endTiming || !flags(c.Verbal) || !flags(c.Somatic) {
		return fmt.Errorf("batch 206 header differs from source")
	}
	if m := c.Material; row.material == "" {
		if m != nil {
			return fmt.Errorf("batch 206 material differs from source")
		}
	} else if m == nil || m.Description != row.material || m.CostCP == nil || *m.CostCP != row.costCP || m.Consumed == nil || *m.Consumed {
		return fmt.Errorf("batch 206 material differs from source")
	}
	sc := s.Scaling
	if sc.Kind != row.scaling.Kind || sc.AboveLevel != row.scaling.AboveLevel || sc.DiceCount != row.scaling.DiceCount || sc.DieSides != row.scaling.DieSides || !slices.Equal(sc.EffectIDs, row.scaling.EffectIDs) {
		return fmt.Errorf("batch 206 scaling differs from source")
	}
	if len(s.Effects) != len(row.effects) {
		return fmt.Errorf("batch 206 effects differ from source")
	}
	for i, want := range row.effects {
		e := s.Effects[i]
		if e.ID != want.id || e.Mechanic != want.mechanic || e.Timing != want.timing || e.DiceCount != want.dice || e.DieSides != want.sides || e.DamageType != want.damage || e.Bonus != want.bonus || e.AddCastingModifier == nil || *e.AddCastingModifier != want.modifier || e.HalfInitialOnly == nil || *e.HalfInitialOnly {
			return fmt.Errorf("batch 206 effect %s differs from source", want.id)
		}
	}
	if !slices.Equal(s.RequiredMechanics, row.required) {
		return fmt.Errorf("batch 206 required mechanics differ from source")
	}
	slots := []int{}
	for level := s.Level; level <= 9 && s.Level > 0; level++ {
		slots = append(slots, level)
	}
	if !slices.Equal(s.SlotLevels, slots) {
		return fmt.Errorf("batch 206 slot levels differ from source")
	}
	return nil
}
