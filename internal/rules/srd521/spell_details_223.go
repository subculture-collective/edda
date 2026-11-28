package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellWallOfFire retains typed SRD source facts without implementing them.
type SpellWallOfFire struct {
	Surface             string   `json:"surface"`
	MaxLengthFeet       int      `json:"max_length_feet"`
	HeightFeet          int      `json:"height_feet"`
	ThicknessFeet       int      `json:"thickness_feet"`
	RingMaxDiameterFeet int      `json:"ring_max_diameter_feet"`
	RingHeightFeet      int      `json:"ring_height_feet"`
	RingThicknessFeet   int      `json:"ring_thickness_feet"`
	Opacity             string   `json:"opacity"`
	AppearSave          string   `json:"appear_save"`
	AppearDiceCount     int      `json:"appear_dice_count"`
	AppearDieSides      int      `json:"appear_die_sides"`
	DamageType          string   `json:"damage_type"`
	AppearSuccess       string   `json:"appear_success"`
	DamagingSide        string   `json:"damaging_side"`
	SideDiceCount       int      `json:"side_dice_count"`
	SideDieSides        int      `json:"side_die_sides"`
	SideDistanceFeet    int      `json:"side_distance_feet"`
	SideTriggers        []string `json:"side_triggers"`
	OtherSide           string   `json:"other_side"`
	SlotDiceCount       int      `json:"slot_dice_count"`
	SlotDieSides        int      `json:"slot_die_sides"`
}

// SpellWallOfForce retains typed SRD source facts without implementing them.
type SpellWallOfForce struct {
	Visibility             string   `json:"visibility"`
	Orientations           []string `json:"orientations"`
	Support                []string `json:"support"`
	Shapes                 []string `json:"shapes"`
	MaxRadiusFeet          int      `json:"max_radius_feet"`
	PanelCount             int      `json:"panel_count"`
	PanelLengthFeet        int      `json:"panel_length_feet"`
	PanelWidthFeet         int      `json:"panel_width_feet"`
	PanelRule              string   `json:"panel_rule"`
	ThicknessQuarterInches int      `json:"thickness_quarter_inches"`
	CutCreature            string   `json:"cut_creature"`
	Passage                string   `json:"passage"`
	DamageImmunity         string   `json:"damage_immunity"`
	DispelImmunity         string   `json:"dispel_immunity"`
	DestroyedBy            string   `json:"destroyed_by"`
	Ethereal               string   `json:"ethereal"`
}

// SpellWallOfIce retains typed SRD source facts without implementing them.
type SpellWallOfIce struct {
	SourcePages         []int                    `json:"source_pages"`
	Surface             string                   `json:"surface"`
	Shapes              []string                 `json:"shapes"`
	MaxRadiusFeet       int                      `json:"max_radius_feet"`
	PanelCount          int                      `json:"panel_count"`
	PanelSquareFeet     int                      `json:"panel_square_feet"`
	PanelRule           string                   `json:"panel_rule"`
	ThicknessFeet       int                      `json:"thickness_feet"`
	CutCreature         string                   `json:"cut_creature"`
	AppearSave          string                   `json:"appear_save"`
	AppearDiceCount     int                      `json:"appear_dice_count"`
	AppearDieSides      int                      `json:"appear_die_sides"`
	DamageType          string                   `json:"damage_type"`
	AppearSuccess       string                   `json:"appear_success"`
	ArmorClass          int                      `json:"armor_class"`
	HitPointsPerSection int                      `json:"hit_points_per_section"`
	SectionFeet         int                      `json:"section_feet"`
	Immunities          []string                 `json:"immunities"`
	Vulnerabilities     []string                 `json:"vulnerabilities"`
	Breach              string                   `json:"breach"`
	FrigidAirTrigger    string                   `json:"frigid_air_trigger"`
	FrigidAirSave       string                   `json:"frigid_air_save"`
	FrigidAirDiceCount  int                      `json:"frigid_air_dice_count"`
	FrigidAirDieSides   int                      `json:"frigid_air_die_sides"`
	FrigidAirSuccess    string                   `json:"frigid_air_success"`
	SlotDice            []SpellHandSlotDiceEntry `json:"slot_dice"`
}

// SpellWallOfStone retains typed SRD source facts without implementing them.
type SpellWallOfStone struct {
	Composition                 string   `json:"composition"`
	PanelCount                  int      `json:"panel_count"`
	ThicknessInches             int      `json:"thickness_inches"`
	PanelLengthFeet             int      `json:"panel_length_feet"`
	PanelHeightFeet             int      `json:"panel_height_feet"`
	ThinPanelLengthFeet         int      `json:"thin_panel_length_feet"`
	ThinPanelHeightFeet         int      `json:"thin_panel_height_feet"`
	ThinPanelThicknessInches    int      `json:"thin_panel_thickness_inches"`
	PanelRule                   string   `json:"panel_rule"`
	CutCreature                 string   `json:"cut_creature"`
	EnclosureTrigger            string   `json:"enclosure_trigger"`
	EnclosureSave               string   `json:"enclosure_save"`
	EnclosureSuccess            string   `json:"enclosure_success"`
	Shape                       string   `json:"shape"`
	Orientation                 string   `json:"orientation"`
	Anchor                      string   `json:"anchor"`
	Uses                        []string `json:"uses"`
	LongSpanOverFeet            int      `json:"long_span_over_feet"`
	LongSpanPanelSizeDivisor    int      `json:"long_span_panel_size_divisor"`
	LongSpanPurpose             string   `json:"long_span_purpose"`
	CrudeShaping                []string `json:"crude_shaping"`
	ArmorClass                  int      `json:"armor_class"`
	HitPointsPerInchOfThickness int      `json:"hit_points_per_inch_of_thickness"`
	Immunities                  []string `json:"immunities"`
	Breach                      string   `json:"breach"`
	Collapse                    string   `json:"collapse"`
	FullConcentrationSeconds    int      `json:"full_concentration_seconds"`
	FullConcentrationResult     string   `json:"full_concentration_result"`
	Otherwise                   string   `json:"otherwise"`
}

// SpellWallOfThorns retains typed SRD source facts without implementing them.
type SpellWallOfThorns struct {
	Composition             string   `json:"composition"`
	Surface                 string   `json:"surface"`
	MaxLengthFeet           int      `json:"max_length_feet"`
	StraightHeightFeet      int      `json:"straight_height_feet"`
	StraightThicknessFeet   int      `json:"straight_thickness_feet"`
	CircleDiameterFeet      int      `json:"circle_diameter_feet"`
	CircleMaxHeightFeet     int      `json:"circle_max_height_feet"`
	CircleThicknessFeet     int      `json:"circle_thickness_feet"`
	Sight                   string   `json:"sight"`
	AppearSave              string   `json:"appear_save"`
	AppearDiceCount         int      `json:"appear_dice_count"`
	AppearDieSides          int      `json:"appear_die_sides"`
	AppearDamageType        string   `json:"appear_damage_type"`
	AppearSuccess           string   `json:"appear_success"`
	MovementCostFeetPerFoot int      `json:"movement_cost_feet_per_foot"`
	PassageTriggers         []string `json:"passage_triggers"`
	PassageSave             string   `json:"passage_save"`
	PassageDiceCount        int      `json:"passage_dice_count"`
	PassageDieSides         int      `json:"passage_die_sides"`
	PassageDamageType       string   `json:"passage_damage_type"`
	PassageSuccess          string   `json:"passage_success"`
	PassageSaveLimit        string   `json:"passage_save_limit"`
	SlotDiceCount           int      `json:"slot_dice_count"`
	SlotDieSides            int      `json:"slot_die_sides"`
}

// SpellWardingBond retains typed SRD source facts without implementing them.
type SpellWardingBond struct {
	TargetRequirement  string   `json:"target_requirement"`
	Connection         string   `json:"connection"`
	BenefitRangeFeet   int      `json:"benefit_range_feet"`
	ArmorClassBonus    int      `json:"armor_class_bonus"`
	SavingThrowBonus   int      `json:"saving_throw_bonus"`
	Resistance         string   `json:"resistance"`
	SharedDamage       string   `json:"shared_damage"`
	RingCount          int      `json:"ring_count"`
	RingMinimumCostCP  int      `json:"ring_minimum_cost_cp"`
	RingWearers        []string `json:"ring_wearers"`
	SeparationOverFeet int      `json:"separation_over_feet"`
	EarlyEnd           []string `json:"early_end"`
}

// SpellWaterBreathing retains typed SRD source facts without implementing them.
type SpellWaterBreathing struct {
	MaxCreatures      int    `json:"max_creatures"`
	Grants            string `json:"grants"`
	NormalRespiration string `json:"normal_respiration"`
}

// SpellWaterWalk retains typed SRD source facts without implementing them.
type SpellWaterWalk struct {
	SourcePages      []int    `json:"source_pages"`
	MaxCreatures     int      `json:"max_creatures"`
	Surfaces         string   `json:"surfaces"`
	SurfaceExamples  []string `json:"surface_examples"`
	AsIf             string   `json:"as_if"`
	MoltenLava       string   `json:"molten_lava"`
	TransitionAction string   `json:"transition_action"`
	Transitions      []string `json:"transitions"`
	Falling          string   `json:"falling"`
}

// SpellWeb retains typed SRD source facts without implementing them.
type SpellWeb struct {
	CubeFeet             int      `json:"cube_feet"`
	Terrain              string   `json:"terrain"`
	Obscurement          string   `json:"obscurement"`
	Anchors              []string `json:"anchors"`
	SolidMassExamples    []string `json:"solid_mass_examples"`
	Unanchored           string   `json:"unanchored"`
	FlatSurfaceDepthFeet int      `json:"flat_surface_depth_feet"`
	Save                 string   `json:"save"`
	SaveTriggers         []string `json:"save_triggers"`
	FailureCondition     string   `json:"failure_condition"`
	ConditionEnds        string   `json:"condition_ends"`
	EscapeAction         string   `json:"escape_action"`
	EscapeCheck          string   `json:"escape_check"`
	EscapeDC             string   `json:"escape_dc"`
	EscapeSuccess        string   `json:"escape_success"`
	FlammableCubeFeet    int      `json:"flammable_cube_feet"`
	BurnAwayRounds       int      `json:"burn_away_rounds"`
	FireDiceCount        int      `json:"fire_dice_count"`
	FireDieSides         int      `json:"fire_die_sides"`
	FireDamageType       string   `json:"fire_damage_type"`
	FireTargets          string   `json:"fire_targets"`
}

// SpellWeird retains typed SRD source facts without implementing them.
type SpellWeird struct {
	Illusion          string `json:"illusion"`
	RadiusFeet        int    `json:"radius_feet"`
	Save              string `json:"save"`
	DamageDiceCount   int    `json:"damage_dice_count"`
	DamageDieSides    int    `json:"damage_die_sides"`
	DamageType        string `json:"damage_type"`
	FailureCondition  string `json:"failure_condition"`
	ConditionDuration string `json:"condition_duration"`
	SaveSuccess       string `json:"save_success"`
	RepeatSave        string `json:"repeat_save"`
	RepeatTiming      string `json:"repeat_timing"`
	RepeatDiceCount   int    `json:"repeat_dice_count"`
	RepeatDieSides    int    `json:"repeat_die_sides"`
	RepeatSuccess     string `json:"repeat_success"`
}

// SpellWindWalk retains typed SRD source facts without implementing them.
type SpellWindWalk struct {
	MaxOtherCreatures   int      `json:"max_other_creatures"`
	Form                string   `json:"form"`
	FlySpeedFeet        int      `json:"fly_speed_feet"`
	Hover               string   `json:"hover"`
	ConditionImmunities []string `json:"condition_immunities"`
	DamageResistances   []string `json:"damage_resistances"`
	AllowedActions      []string `json:"allowed_actions"`
	RevertSeconds       int      `json:"revert_seconds"`
	RevertCondition     string   `json:"revert_condition"`
	ReturnToCloud       string   `json:"return_to_cloud"`
	ReturnSeconds       int      `json:"return_seconds"`
	DescentTrigger      string   `json:"descent_trigger"`
	DescentFeetPerRound int      `json:"descent_feet_per_round"`
	DescentSeconds      int      `json:"descent_seconds"`
	SafeLanding         string   `json:"safe_landing"`
	NoLanding           string   `json:"no_landing"`
}

// SpellWindWall retains typed SRD source facts without implementing them.
type SpellWindWall struct {
	Origin           string   `json:"origin"`
	MaxLengthFeet    int      `json:"max_length_feet"`
	MaxHeightFeet    int      `json:"max_height_feet"`
	ThicknessFeet    int      `json:"thickness_feet"`
	Shape            string   `json:"shape"`
	Save             string   `json:"save"`
	DamageDiceCount  int      `json:"damage_dice_count"`
	DamageDieSides   int      `json:"damage_die_sides"`
	DamageType       string   `json:"damage_type"`
	SaveSuccess      string   `json:"save_success"`
	KeptAtBay        []string `json:"kept_at_bay"`
	Blocked          []string `json:"blocked"`
	LooseMaterials   string   `json:"loose_materials"`
	Deflected        []string `json:"deflected"`
	Deflection       string   `json:"deflection"`
	Unaffected       []string `json:"unaffected"`
	GaseousCreatures string   `json:"gaseous_creatures"`
}

func spell223Templates() map[string]any {
	return map[string]any{
		"wall_of_fire":    SpellWallOfFire{},
		"wall_of_force":   SpellWallOfForce{},
		"wall_of_ice":     SpellWallOfIce{},
		"wall_of_stone":   SpellWallOfStone{},
		"wall_of_thorns":  SpellWallOfThorns{},
		"warding_bond":    SpellWardingBond{},
		"water_breathing": SpellWaterBreathing{},
		"water_walk":      SpellWaterWalk{},
		"web":             SpellWeb{},
		"weird":           SpellWeird{},
		"wind_walk":       SpellWindWalk{},
		"wind_wall":       SpellWindWall{},
	}
}

func validSpell223Target(t SpellTarget) bool {
	switch t.Shape {
	case "wall":
		return t.Count == 0 && slices.Contains([]string{"creatures-in-wall-of-fire", "force-wall-on-point-in-range", "creatures-in-ice-wall", "stone-wall-on-point-in-range", "creatures-in-wall-of-thorns", "creatures-in-wind-wall"}, t.Eligibility)
	case "single":
		return t.Count == 1 && t.Eligibility == "other-willing-creature-touch"
	case "up-to-count":
		return t.Count == 10 && t.Eligibility == "chosen-willing-creatures-in-range"
	case "self-and-up-to-count":
		return t.Count == 10 && t.Eligibility == "chosen-willing-creatures-in-range"
	case "cube":
		return t.Count == 0 && t.Eligibility == "webs-in-cube-on-point-in-range"
	case "selected-creatures":
		return t.Count == 0 && t.Eligibility == "chosen-creatures-in-sphere-on-point-in-range"
	}
	return false
}

// has223Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has223Payload() bool {
	return d != nil && (d.WallOfFire != nil || d.WallOfForce != nil || d.WallOfIce != nil || d.WallOfStone != nil || d.WallOfThorns != nil || d.WardingBond != nil || d.WaterBreathing != nil || d.WaterWalk != nil || d.Web != nil || d.Weird != nil || d.WindWalk != nil || d.WindWall != nil)
}

// validateSpell223 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell223() (bool, error) {
	pin, ok := spell223Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has223Payload() {
			return true, fmt.Errorf("mismatched batch 223 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 223 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 223 record differs from the printed source")
	}
	return true, nil
}
