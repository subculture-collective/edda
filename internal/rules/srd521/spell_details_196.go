package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellSplashUpgrade is one source-listed cantrip damage tier.
type SpellSplashUpgrade struct {
	CharacterLevel int `json:"character_level"`
	DiceCount      int `json:"dice_count"`
	DieSides       int `json:"die_sides"`
}

// SpellSplash retains typed SRD source facts without implementing them.
type SpellSplash struct {
	Center      string               `json:"center"`
	RadiusFeet  int                  `json:"radius_feet"`
	Save        string               `json:"save"`
	SaveSuccess string               `json:"save_success"`
	Upgrades    []SpellSplashUpgrade `json:"upgrades"`
	Appearance  string               `json:"appearance"`
}

// SpellAid retains typed SRD source facts without implementing them.
type SpellAid struct {
	BaseTargets                int    `json:"base_targets"`
	TargetChoice               string `json:"target_choice"`
	HitPointMaximumIncrease    int    `json:"hit_point_maximum_increase"`
	CurrentHitPointIncrease    int    `json:"current_hit_point_increase"`
	AdditionalHitPointsPerSlot int    `json:"additional_hit_points_per_slot"`
}

// SpellAlarmMode is one of the two source-listed alarm variants.
type SpellAlarmMode struct {
	Choice          string `json:"choice"`
	Signal          string `json:"signal"`
	SignalSeconds   int    `json:"signal_seconds"`
	DistanceFeet    int    `json:"distance_feet"`
	DistanceMiles   int    `json:"distance_miles"`
	DistanceMeasure string `json:"distance_measure"`
	AwakensCaster   bool   `json:"awakens_caster"`
}

// SpellAlarm retains typed SRD source facts without implementing them.
type SpellAlarm struct {
	WardChoices       []string         `json:"ward_choices"`
	AreaMaxCubeFeet   int              `json:"area_max_cube_feet"`
	Triggers          []string         `json:"triggers"`
	DesignateExempt   bool             `json:"designate_exempt"`
	ExemptionTiming   string           `json:"exemption_timing"`
	ChooseCount       int              `json:"choose_count"`
	Modes             []SpellAlarmMode `json:"modes"`
	RitualAlternative bool             `json:"ritual_alternative"`
	AlertedCreature   string           `json:"alerted_creature"`
}

// SpellAlterAquatic is the Aquatic Adaptation option.
type SpellAlterAquatic struct {
	Choice            string `json:"choice"`
	Appearance        string `json:"appearance"`
	BreatheUnderwater bool   `json:"breathe_underwater"`
	SwimSpeed         string `json:"swim_speed"`
}

// SpellAlterAppearance is the Change Appearance option.
type SpellAlterAppearance struct {
	Choice               string   `json:"choice"`
	Controls             []string `json:"controls"`
	AppearAsOtherSpecies bool     `json:"appear_as_other_species"`
	StatisticsChange     bool     `json:"statistics_change"`
	DifferentSize        bool     `json:"different_size"`
	BasicShapeChanges    bool     `json:"basic_shape_changes"`
	ShapeExample         string   `json:"shape_example"`
	RepeatAction         string   `json:"repeat_action"`
}

// SpellAlterGrowth is one source-listed natural weapon and its damage type.
type SpellAlterGrowth struct {
	Growth     string `json:"growth"`
	DamageType string `json:"damage_type"`
}

// SpellAlterWeapons is the Natural Weapons option.
type SpellAlterWeapons struct {
	Choice               string             `json:"choice"`
	Growths              []SpellAlterGrowth `json:"growths"`
	Trigger              string             `json:"trigger"`
	DiceCount            int                `json:"dice_count"`
	DieSides             int                `json:"die_sides"`
	ReplacesNormalDamage bool               `json:"replaces_normal_damage"`
	AttackDamageAbility  string             `json:"attack_damage_ability"`
	InsteadOf            string             `json:"instead_of"`
}

// SpellAlterSelf retains typed SRD source facts without implementing them.
type SpellAlterSelf struct {
	ChooseCount   int                  `json:"choose_count"`
	ReplaceAction string               `json:"replace_action"`
	Aquatic       SpellAlterAquatic    `json:"aquatic"`
	Appearance    SpellAlterAppearance `json:"appearance"`
	Weapons       SpellAlterWeapons    `json:"weapons"`
}

// SpellFriendship retains typed SRD source facts without implementing them.
type SpellFriendship struct {
	SourcePages              []int    `json:"source_pages"`
	BaseTargets              int      `json:"base_targets"`
	AdditionalTargetsPerSlot int      `json:"additional_targets_per_slot"`
	TargetType               string   `json:"target_type"`
	Save                     string   `json:"save"`
	SaveSuccess              string   `json:"save_success"`
	Condition                string   `json:"condition"`
	EndTriggers              []string `json:"end_triggers"`
}

// SpellMessenger retains typed SRD source facts without implementing them.
type SpellMessenger struct {
	TargetSize                       string   `json:"target_size"`
	TargetType                       string   `json:"target_type"`
	Save                             string   `json:"save"`
	SaveSuccess                      string   `json:"save_success"`
	AutomaticSuccess                 string   `json:"automatic_success"`
	Destination                      string   `json:"destination"`
	Recipient                        string   `json:"recipient"`
	RecipientExamples                []string `json:"recipient_examples"`
	MaxWords                         int      `json:"max_words"`
	TravelPrecision                  string   `json:"travel_precision"`
	MilesPer24Hours                  int      `json:"miles_per_24_hours"`
	FlyingMilesPer24Hours            int      `json:"flying_miles_per_24_hours"`
	Delivery                         string   `json:"delivery"`
	Failure                          string   `json:"failure"`
	RitualAlternative                bool     `json:"ritual_alternative"`
	AdditionalDurationSecondsPerSlot int      `json:"additional_duration_seconds_per_slot"`
}

// SpellShapes retains typed SRD source facts without implementing them.
type SpellShapes struct {
	TargetChoice             string   `json:"target_choice"`
	FormMaxSize              string   `json:"form_max_size"`
	FormType                 string   `json:"form_type"`
	FormMaxCR                int      `json:"form_max_cr"`
	DifferentFormPerTarget   bool     `json:"different_form_per_target"`
	RetransformAction        string   `json:"retransform_action"`
	RetransformTiming        string   `json:"retransform_timing"`
	ReplaceStatistics        string   `json:"replace_statistics"`
	Retained                 []string `json:"retained"`
	ActionsLimitedBy         string   `json:"actions_limited_by"`
	CanCastSpells            bool     `json:"can_cast_spells"`
	Equipment                string   `json:"equipment"`
	TemporaryHP              string   `json:"temporary_hp"`
	TemporaryHPEndsWithSpell bool     `json:"temporary_hp_ends_with_spell"`
	TargetEndAction          string   `json:"target_end_action"`
}

// SpellAnimateDeadChoice maps a chosen remains kind to the created creature.
type SpellAnimateDeadChoice struct {
	Choice   string `json:"choice"`
	Creature string `json:"creature"`
}

// SpellAnimateDead retains typed SRD source facts without implementing them.
type SpellAnimateDead struct {
	TargetChoices              []SpellAnimateDeadChoice `json:"target_choices"`
	RemainsSizes               []string                 `json:"remains_sizes"`
	RemainsType                string                   `json:"remains_type"`
	CreatureType               string                   `json:"creature_type"`
	StatBlocks                 string                   `json:"stat_blocks"`
	CommandAction              string                   `json:"command_action"`
	CommandTiming              string                   `json:"command_timing"`
	CommandRangeFeet           int                      `json:"command_range_feet"`
	CommandMultiple            string                   `json:"command_multiple"`
	CommandChoices             []string                 `json:"command_choices"`
	GeneralCommandExamples     []string                 `json:"general_command_examples"`
	NoCommand                  string                   `json:"no_command"`
	OrderPersistence           string                   `json:"order_persistence"`
	ControlHours               int                      `json:"control_hours"`
	ControlEnd                 string                   `json:"control_end"`
	ReassertBy                 string                   `json:"reassert_by"`
	ReassertHours              int                      `json:"reassert_hours"`
	ReassertCount              int                      `json:"reassert_count"`
	ReassertReplacesAnimating  bool                     `json:"reassert_replaces_animating"`
	AdditionalCreaturesPerSlot int                      `json:"additional_creatures_per_slot"`
	AdditionalSource           string                   `json:"additional_source"`
}

// SpellSizeCount is a per-size count, such as object slots or Hit Points.
type SpellSizeCount struct {
	Size  string `json:"size"`
	Count int    `json:"count"`
}

// SpellSizeDice is a per-size dice expression. The Slam bonus is explicit.
type SpellSizeDice struct {
	Size                    string `json:"size"`
	DiceCount               int    `json:"dice_count"`
	DieSides                int    `json:"die_sides"`
	Bonus                   int    `json:"bonus"`
	AddSpellcastingModifier bool   `json:"add_spellcasting_modifier"`
}

// SpellStatAbility is one stat block ability row. Modifier and save are the
// only signed quantities in this batch; validation pins their exact values.
type SpellStatAbility struct {
	Ability  string `json:"ability"`
	Score    int    `json:"score"`
	Modifier int    `json:"modifier" srd:"signed"`
	Save     int    `json:"save" srd:"signed"`
}

// SpellAnimatedSlam is the Animated Object's only source-listed action.
type SpellAnimatedSlam struct {
	AttackKind  string          `json:"attack_kind"`
	AttackBonus string          `json:"attack_bonus"`
	ReachFeet   int             `json:"reach_feet"`
	DamageType  string          `json:"damage_type"`
	Damage      []SpellSizeDice `json:"damage"`
}

// SpellAnimatedObject is the stat block printed with Animate Objects (p.109).
type SpellAnimatedObject struct {
	Name                string             `json:"name"`
	Sizes               string             `json:"sizes"`
	CreatureType        string             `json:"creature_type"`
	Alignment           string             `json:"alignment"`
	ArmorClass          int                `json:"armor_class"`
	HitPoints           []SpellSizeCount   `json:"hit_points"`
	SpeedFeet           int                `json:"speed_feet"`
	Abilities           []SpellStatAbility `json:"abilities"`
	DamageImmunities    []string           `json:"damage_immunities"`
	ConditionImmunities []string           `json:"condition_immunities"`
	BlindsightFeet      int                `json:"blindsight_feet"`
	PassivePerception   int                `json:"passive_perception"`
	Languages           string             `json:"languages"`
	ChallengeRating     string             `json:"challenge_rating"`
	XP                  int                `json:"xp"`
	ProficiencyBonus    string             `json:"proficiency_bonus"`
	Slam                SpellAnimatedSlam  `json:"slam"`
}

// SpellAnimateObjects retains typed SRD source facts without implementing them.
type SpellAnimateObjects struct {
	SourcePages        []int               `json:"source_pages"`
	ObjectRequirements []string            `json:"object_requirements"`
	MaxObjects         string              `json:"max_objects"`
	SizeCosts          []SpellSizeCount    `json:"size_costs"`
	Appearance         string              `json:"appearance"`
	CreatureType       string              `json:"creature_type"`
	StatBlock          string              `json:"stat_block"`
	ControlEnds        []string            `json:"control_ends"`
	Ally               string              `json:"ally"`
	Initiative         string              `json:"initiative"`
	CommandAction      string              `json:"command_action"`
	CommandRangeFeet   int                 `json:"command_range_feet"`
	CommandMultiple    string              `json:"command_multiple"`
	NoCommand          string              `json:"no_command"`
	AtZeroHitPoints    string              `json:"at_zero_hit_points"`
	SlotDamage         []SpellSizeDice     `json:"slot_damage"`
	AnimatedObject     SpellAnimatedObject `json:"animated_object"`
}

// SpellAntilife retains typed SRD source facts without implementing them.
type SpellAntilife struct {
	EmanationFeet int      `json:"emanation_feet"`
	Origin        string   `json:"origin"`
	Prevents      []string `json:"prevents"`
	ExemptTypes   []string `json:"exempt_types"`
	AffectedMay   []string `json:"affected_may"`
	EndTrigger    string   `json:"end_trigger"`
}

// SpellAntimagic retains typed SRD source facts without implementing them.
type SpellAntimagic struct {
	EmanationFeet               int      `json:"emanation_feet"`
	Origin                      string   `json:"origin"`
	Prohibited                  []string `json:"prohibited"`
	OutsideEffects              string   `json:"outside_effects"`
	MagicItems                  string   `json:"magic_items"`
	AreasOfEffect               string   `json:"areas_of_effect"`
	Teleportation               string   `json:"teleportation"`
	PlanarTravel                string   `json:"planar_travel"`
	Portals                     string   `json:"portals"`
	OngoingSpells               string   `json:"ongoing_spells"`
	SuppressionExceptions       []string `json:"suppression_exceptions"`
	SuppressedFunctions         bool     `json:"suppressed_functions"`
	SuppressedTimeCountsAgainst bool     `json:"suppressed_time_counts_against_duration"`
	DispelMagic                 string   `json:"dispel_magic"`
	OtherAntimagicFields        string   `json:"other_antimagic_fields"`
}

func spell196Templates() map[string]any {
	return map[string]any{
		"splash":     SpellSplash{},
		"aid":        SpellAid{},
		"alarm":      SpellAlarm{},
		"alter":      SpellAlterSelf{},
		"friendship": SpellFriendship{},
		"messenger":  SpellMessenger{},
		"shapes":     SpellShapes{},
		"dead":       SpellAnimateDead{},
		"objects":    SpellAnimateObjects{},
		"antilife":   SpellAntilife{},
		"antimagic":  SpellAntimagic{},
	}
}

// spell196Scaling pins each payload to its source scaling shape, so a typed
// increment cannot be contradicted by a different common scaling kind.
var spell196Scaling = map[string]string{
	"splash":     "cantrip-upgrade",
	"aid":        "slot-bonus",
	"alarm":      "none",
	"alter":      "none",
	"friendship": "slot-targets",
	"messenger":  "slot-duration",
	"shapes":     "none",
	"dead":       "slot-creatures",
	"objects":    "slot-dice-by-size",
	"antilife":   "none",
	"antimagic":  "none",
}

// spell196Header is the printed header of one #196 family (PDF pp.107-109).
type spell196Header struct {
	name            string
	page, level     int
	rangeFeet       int
	durationSeconds int
	concentration   bool
}

// spell196Headers binds each #196 ID to its printed name, page, level, range,
// duration and concentration, so a structurally valid record cannot carry
// another spell's header under one of these IDs.
var spell196Headers = map[string]spell196Header{
	"acid-arrow":        {"Acid Arrow", 107, 2, 90, 0, false},
	"acid-splash":       {"Acid Splash", 107, 0, 60, 0, false},
	"aid":               {"Aid", 107, 2, 30, 28800, false},
	"alarm":             {"Alarm", 107, 1, 30, 28800, false},
	"alter-self":        {"Alter Self", 107, 2, 0, 3600, true},
	"animal-friendship": {"Animal Friendship", 107, 1, 30, 86400, false},
	"animal-messenger":  {"Animal Messenger", 108, 2, 30, 86400, false},
	"animal-shapes":     {"Animal Shapes", 108, 8, 30, 86400, false},
	"animate-dead":      {"Animate Dead", 108, 3, 10, 0, false},
	"animate-objects":   {"Animate Objects", 108, 5, 120, 60, true},
	"antilife-shell":    {"Antilife Shell", 109, 5, 0, 3600, true},
	"antimagic-field":   {"Antimagic Field", 109, 8, 0, 3600, true},
}

// validSpell196Header pins the printed header of a complete #196 record. Acid
// Arrow has no typed payload, so its printed dice are pinned here as well.
func (s SpellRecord) validSpell196Header() bool {
	h, ok := spell196Headers[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok || !strings.HasPrefix(s.ID, prefix+"spells/") || s.Scope != "complete-source-record" {
		return true
	}
	if s.Name != h.name || s.SourcePage != h.page || s.Level != h.level || s.Range.Feet != h.rangeFeet || s.DurationSeconds != h.durationSeconds || s.Concentration == nil || *s.Concentration != h.concentration {
		return false
	}
	if h.name != "Acid Arrow" {
		return true
	}
	return len(s.Effects) == 2 && s.Effects[0].DiceCount == 4 && s.Effects[0].DieSides == 4 && s.Effects[1].DiceCount == 2 && s.Effects[1].DieSides == 4 && s.Scaling.DiceCount == 1 && s.Scaling.DieSides == 4
}

// validateSpell196 checks the #196 payloads. Acid Arrow, the twelfth family,
// uses only the common record fields and therefore has no typed payload.
func (s SpellRecord) validateSpell196() (bool, error) {
	if !s.validSpell196Header() {
		return true, fmt.Errorf("batch 196 header differs from source")
	}
	names := map[string]string{"splash": "acid-splash", "aid": "aid", "alarm": "alarm", "alter": "alter-self", "friendship": "animal-friendship", "messenger": "animal-messenger", "shapes": "animal-shapes", "dead": "animate-dead", "objects": "animate-objects", "antilife": "antilife-shell", "antimagic": "antimagic-field"}
	expected := ""
	for key, name := range names {
		if s.ID == prefix+"spells/"+name {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 196 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"splash":     d.Splash,
		"aid":        d.Aid,
		"alarm":      d.Alarm,
		"alter":      d.Alter,
		"friendship": d.Friendship,
		"messenger":  d.Messenger,
		"shapes":     d.Shapes,
		"dead":       d.Dead,
		"objects":    d.Objects,
		"antilife":   d.Antilife,
		"antimagic":  d.Antimagic,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 196 payload")
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
	castingSeconds := 0
	if expected == "alarm" || expected == "dead" {
		castingSeconds = 60
	}
	if count != 1 || s.Scope != "complete-source-record" || d.CastingSeconds != castingSeconds || len(d.CastingTriggers) != 0 || d.CastingTriggers == nil || strings.TrimSpace(d.EndTiming) == "" || s.Scaling.Kind != spell196Scaling[expected] {
		return true, fmt.Errorf("invalid batch 196 source details")
	}
	if d.Barkskin != nil || d.Beacon != nil || d.Befuddlement != nil || d.Curse != nil || d.Tentacles != nil || d.Barrier != nil || d.Bless != nil || d.Blight != nil || d.Blindness != nil || d.Blink != nil || d.Blur != nil || d.Burning != nil || d.Rays != nil || d.Scrying != nil || d.Smite != nil || d.Chest != nil || d.Vision != nil || d.Seeming != nil || d.Sending != nil || d.Sequester != nil || d.Shapechange != nil || d.Shatter != nil || d.Shield != nil || d.Faith != nil {
		return true, fmt.Errorf("multiple source payloads")
	}
	sizes := func(rows []SpellSizeCount, counts ...int) bool {
		return slices.Equal(rows, []SpellSizeCount{{"medium-or-smaller", counts[0]}, {"large", counts[1]}, {"huge", counts[2]}})
	}
	switch {
	case d.Splash != nil:
		if d.Splash.RadiusFeet != 5 || d.Splash.Save != "dexterity" || s.Scaling.DiceCount != 1 || s.Scaling.DieSides != 6 || s.Level != 0 || !slices.Equal(d.Splash.Upgrades, []SpellSplashUpgrade{{5, 2, 6}, {11, 3, 6}, {17, 4, 6}}) {
			return true, fmt.Errorf("invalid cantrip area or upgrades")
		}
	case d.Aid != nil:
		if d.Aid.BaseTargets != s.Target.Count || d.Aid.BaseTargets != 3 || d.Aid.HitPointMaximumIncrease != 5 || d.Aid.AdditionalHitPointsPerSlot != 5 || d.Aid.CurrentHitPointIncrease != d.Aid.HitPointMaximumIncrease {
			return true, fmt.Errorf("invalid aid targets or increase")
		}
	case d.Alarm != nil:
		a := d.Alarm
		if !slices.Equal(a.WardChoices, []string{"door", "window", "area"}) || a.ChooseCount != 1 || len(a.Modes) != 2 || a.Modes[0].Choice != "audible" || a.Modes[1].Choice != "mental" || a.AreaMaxCubeFeet != 20 || a.Modes[0].SignalSeconds != 10 || a.Modes[0].DistanceFeet != 60 || a.Modes[0].DistanceMiles != 0 || a.Modes[1].SignalSeconds != 0 || a.Modes[1].DistanceFeet != 0 || a.Modes[1].DistanceMiles != 1 || !*s.Ritual || !a.RitualAlternative {
			return true, fmt.Errorf("invalid alarm ward or mode variants")
		}
	case d.Alter != nil:
		a := d.Alter
		growths := []SpellAlterGrowth{{"claws", "slashing"}, {"fangs", "piercing"}, {"horns", "piercing"}, {"hooves", "bludgeoning"}}
		if a.ChooseCount != 1 || a.Aquatic.Choice != "aquatic-adaptation" || a.Appearance.Choice != "change-appearance" || a.Weapons.Choice != "natural-weapons" || !slices.Equal(a.Weapons.Growths, growths) || a.Weapons.DiceCount != 1 || a.Weapons.DieSides != 6 {
			return true, fmt.Errorf("invalid alter self options")
		}
	case d.Friendship != nil:
		if d.Friendship.Save != "wisdom" || d.Friendship.BaseTargets != 1 || d.Friendship.AdditionalTargetsPerSlot != 1 || !slices.Equal(d.Friendship.SourcePages, []int{107, 108}) || len(d.Friendship.EndTriggers) != 2 {
			return true, fmt.Errorf("invalid friendship pages or endings")
		}
	case d.Messenger != nil:
		m := d.Messenger
		if m.Save != "charisma" || m.MaxWords != 25 || m.MilesPer24Hours != 25 || m.FlyingMilesPer24Hours != 50 || m.AdditionalDurationSecondsPerSlot != 172800 || !*s.Ritual || !m.RitualAlternative {
			return true, fmt.Errorf("invalid messenger limits")
		}
	case d.Shapes != nil:
		if d.Shapes.FormMaxCR != 4 || d.Shapes.CanCastSpells || len(d.Shapes.Retained) != 8 {
			return true, fmt.Errorf("invalid animal form limits")
		}
	case d.Dead != nil:
		a := d.Dead
		if !slices.Equal(a.TargetChoices, []SpellAnimateDeadChoice{{"pile-of-bones", "skeleton"}, {"corpse", "zombie"}}) || !slices.Equal(a.RemainsSizes, []string{"medium", "small"}) || len(a.CommandChoices) != 2 || a.CommandRangeFeet != 60 || a.ControlHours != 24 || a.ReassertHours != 24 || a.ReassertCount != 4 || a.AdditionalCreaturesPerSlot != 2 {
			return true, fmt.Errorf("invalid animate dead choices or control")
		}
	case d.Objects != nil:
		o := d.Objects
		block := o.AnimatedObject
		abilities := []SpellStatAbility{{"strength", 16, 3, 3}, {"dexterity", 10, 0, 0}, {"constitution", 10, 0, 0}, {"intelligence", 3, -4, -4}, {"wisdom", 3, -4, -4}, {"charisma", 1, -5, -5}}
		slot := []SpellSizeDice{{"medium-or-smaller", 1, 4, 0, false}, {"large", 1, 6, 0, false}, {"huge", 1, 12, 0, false}}
		slam := []SpellSizeDice{{"medium-or-smaller", 1, 4, 3, false}, {"large", 2, 6, 3, true}, {"huge", 2, 12, 3, true}}
		if !slices.Equal(o.SourcePages, []int{108, 109}) || !sizes(o.SizeCosts, 1, 2, 3) || !slices.Equal(o.SlotDamage, slot) || o.CommandRangeFeet != 500 {
			return true, fmt.Errorf("invalid animate objects pages, sizes or scaling")
		}
		if !slices.Equal(block.Abilities, abilities) || !sizes(block.HitPoints, 10, 20, 40) || !slices.Equal(block.Slam.Damage, slam) || block.ArmorClass != 15 || block.SpeedFeet != 30 || block.Slam.ReachFeet != 5 || block.BlindsightFeet != 30 || block.PassivePerception != 6 || block.XP != 0 || !slices.Equal(block.DamageImmunities, []string{"poison", "psychic"}) || !slices.Equal(block.ConditionImmunities, []string{"charmed", "exhaustion", "frightened", "paralyzed", "poisoned"}) {
			return true, fmt.Errorf("invalid animated object stat block")
		}
	case d.Antilife != nil:
		if d.Antilife.EmanationFeet != 10 || !slices.Equal(d.Antilife.ExemptTypes, []string{"construct", "undead"}) {
			return true, fmt.Errorf("invalid antilife aura or exemptions")
		}
	case d.Antimagic != nil:
		if d.Antimagic.EmanationFeet != 10 || d.Antimagic.SuppressedFunctions || !slices.Equal(d.Antimagic.SuppressionExceptions, []string{"cast-by-artifact", "cast-by-deity"}) {
			return true, fmt.Errorf("invalid antimagic aura or exceptions")
		}
	}
	return true, nil
}

// validSpell196Scaling checks the scaling kinds introduced by this batch. Each
// references a typed increment and a named, still unavailable, dependency.
func (s SpellRecord) validSpell196Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || len(sc.EffectIDs) == 0 || !uniqueNames(sc.EffectIDs) {
		return false
	}
	for _, id := range sc.EffectIDs {
		if !effects[id] {
			return false
		}
	}
	if sc.Kind == "cantrip-upgrade" {
		return s.Level == 0 && sc.AboveLevel == 0 && sc.DiceCount > 0 && slices.Contains([]int{4, 6, 8, 10, 12, 20}, sc.DieSides) && (d.Splash != nil || d.Chill != nil || d.FireBolt != nil) && slices.Contains(s.RequiredMechanics, "scaling.cantrip-upgrade")
	}
	if s.Level == 0 || sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 {
		return false
	}
	switch sc.Kind {
	case "slot-bonus":
		return ((d.Aid != nil && d.Aid.AdditionalHitPointsPerSlot > 0) || (d.FalseLife != nil && d.FalseLife.AdditionalTemporaryHitPointsPerSlot > 0)) && slices.Contains(s.RequiredMechanics, "scaling.bonus-per-slot")
	case "slot-creatures":
		return d.Dead != nil && d.Dead.AdditionalCreaturesPerSlot > 0 && slices.Contains(s.RequiredMechanics, "scaling.creatures-per-slot")
	case "slot-dice-by-size":
		return d.Objects != nil && len(d.Objects.SlotDamage) == 3 && slices.Contains(s.RequiredMechanics, "scaling.per-slot-by-size")
	}
	return false
}
