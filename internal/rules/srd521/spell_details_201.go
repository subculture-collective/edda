package srd521

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellWoodland retains Conjure Woodland Beings' typed SRD source facts without implementing them.
type SpellWoodland struct {
	SourcePages         []int    `json:"source_pages"`
	EmanationFeet       int      `json:"emanation_feet"`
	Origin              string   `json:"origin"`
	Triggers            []string `json:"triggers"`
	TriggerVisibility   string   `json:"trigger_visibility"`
	Save                string   `json:"save"`
	SaveSuccess         string   `json:"save_success"`
	SavesPerTurn        int      `json:"saves_per_turn"`
	BonusAction         string   `json:"bonus_action"`
	BonusActionDuration string   `json:"bonus_action_duration"`
}

// SpellOtherPlane retains Contact Other Plane's typed SRD source facts without implementing them.
type SpellOtherPlane struct {
	Save                   string   `json:"save"`
	SaveDC                 int      `json:"save_dc"`
	SuccessQuestionsMax    int      `json:"success_questions_max"`
	QuestionTiming         string   `json:"question_timing"`
	AnswerFormat           string   `json:"answer_format"`
	AnswerExamples         []string `json:"answer_examples"`
	UnclearMeaning         string   `json:"unclear_meaning"`
	MisleadingAlternative  string   `json:"misleading_alternative"`
	FailureDamageDiceCount int      `json:"failure_damage_dice_count"`
	FailureDamageDieSides  int      `json:"failure_damage_die_sides"`
	FailureDamageType      string   `json:"failure_damage_type"`
	FailureCondition       string   `json:"failure_condition"`
	FailureEnds            []string `json:"failure_ends"`
}

// SpellContagion retains Contagion's typed SRD source facts without implementing them.
type SpellContagion struct {
	Save             string `json:"save"`
	SaveSuccess      string `json:"save_success"`
	AbilityChoice    string `json:"ability_choice"`
	AbilityEffect    string `json:"ability_effect"`
	RepeatTiming     string `json:"repeat_timing"`
	SuccessesToEnd   int    `json:"successes_to_end"`
	SuccessResult    string `json:"success_result"`
	FailuresToLast   int    `json:"failures_to_last"`
	FailureResult    string `json:"failure_result"`
	DurationDays     int    `json:"duration_days"`
	EndEffectSave    string `json:"end_effect_save"`
	EndEffectFailure string `json:"end_effect_failure"`
}

// SpellContingency retains Contingency's typed SRD source facts without implementing them.
type SpellContingency struct {
	ContingentMaxLevel    int    `json:"contingent_max_level"`
	ContingentCastingTime string `json:"contingent_casting_time"`
	ContingentTarget      string `json:"contingent_target"`
	ContingentCaster      string `json:"contingent_caster"`
	SlotsExpended         string `json:"slots_expended"`
	DeferredOnCast        bool   `json:"deferred_on_cast"`
	TriggerDescribed      string `json:"trigger_described"`
	TriggerExampleSpell   string `json:"trigger_example_spell"`
	TriggerExample        string `json:"trigger_example"`
	TakesEffect           string `json:"takes_effect"`
	RegardlessOfWish      bool   `json:"regardless_of_wish"`
	EndsAfterTrigger      bool   `json:"ends_after_trigger"`
	AffectsOnly           string `json:"affects_only"`
	MaxActive             int    `json:"max_active"`
	RecastEnds            string `json:"recast_ends"`
	ComponentRequirement  string `json:"component_requirement"`
}

// SpellFlame retains Continual Flame's typed SRD source facts without implementing them.
type SpellFlame struct {
	Source                 string `json:"source"`
	BrightLightRadiusFeet  int    `json:"bright_light_radius_feet"`
	DimLightAdditionalFeet int    `json:"dim_light_additional_feet"`
	Appearance             string `json:"appearance"`
	Heat                   string `json:"heat"`
	Fuel                   string `json:"fuel"`
	Concealment            string `json:"concealment"`
	Extinguish             string `json:"extinguish"`
}

// SpellWaterFlood retains Control Water's Flood effect.
type SpellWaterFlood struct {
	StandingWaterRiseFeet int    `json:"standing_water_rise_feet"`
	WaveHeightFeet        int    `json:"wave_height_feet"`
	WaveArea              string `json:"wave_area"`
	WavePath              string `json:"wave_path"`
	VehicleSize           string `json:"vehicle_size"`
	VehicleEffect         string `json:"vehicle_effect"`
	CapsizePercent        int    `json:"capsize_percent"`
	CapsizeTarget         string `json:"capsize_target"`
	ElevatedUntil         string `json:"elevated_until"`
	WaveRepeats           string `json:"wave_repeats"`
}

// SpellWaterPart retains Control Water's Part Water effect.
type SpellWaterPart struct {
	Result string `json:"result"`
	Walls  string `json:"walls"`
	Lasts  string `json:"lasts"`
	Refill string `json:"refill"`
}

// SpellWaterRedirect retains Control Water's Redirect Flow effect.
type SpellWaterRedirect struct {
	Movement   string `json:"movement"`
	Obstacles  string `json:"obstacles"`
	BeyondArea string `json:"beyond_area"`
	Lasts      string `json:"lasts"`
}

// SpellWaterWhirlpool retains Control Water's Whirlpool effect.
type SpellWaterWhirlpool struct {
	MinimumAreaSquareFeet int      `json:"minimum_area_square_feet"`
	MinimumDepthFeet      int      `json:"minimum_depth_feet"`
	Location              string   `json:"location"`
	BaseWidthFeet         int      `json:"base_width_feet"`
	TopWidthFeet          int      `json:"top_width_feet"`
	HeightFeet            int      `json:"height_feet"`
	PullRadiusFeet        int      `json:"pull_radius_feet"`
	PullSubjects          string   `json:"pull_subjects"`
	PullFeet              int      `json:"pull_feet"`
	SaveTriggers          []string `json:"save_triggers"`
	Save                  string   `json:"save"`
	SaveSuccess           string   `json:"save_success"`
	EscapeAction          string   `json:"escape_action"`
	EscapeCheck           string   `json:"escape_check"`
	Lasts                 string   `json:"lasts"`
}

// SpellWaterControl retains Control Water's typed SRD source facts without implementing them.
type SpellWaterControl struct {
	SourcePages   []int               `json:"source_pages"`
	CubeMaxFeet   int                 `json:"cube_max_feet"`
	RepeatAction  string              `json:"repeat_action"`
	RepeatChoices []string            `json:"repeat_choices"`
	EffectChoices []string            `json:"effect_choices"`
	Flood         SpellWaterFlood     `json:"flood"`
	Part          SpellWaterPart      `json:"part"`
	Redirect      SpellWaterRedirect  `json:"redirect"`
	Whirlpool     SpellWaterWhirlpool `json:"whirlpool"`
}

// SpellWeatherStage is one printed row of a Control Weather table.
type SpellWeatherStage struct {
	Stage     int    `json:"stage"`
	Condition string `json:"condition"`
}

// SpellWeather retains Control Weather's typed SRD source facts without implementing them.
type SpellWeather struct {
	RadiusMiles             int                 `json:"radius_miles"`
	RequiresOutdoors        bool                `json:"requires_outdoors"`
	EndsIndoors             bool                `json:"ends_indoors"`
	CurrentConditionsBy     string              `json:"current_conditions_by"`
	Changeable              []string            `json:"changeable"`
	DelayDiceCount          int                 `json:"delay_dice_count"`
	DelayDieSides           int                 `json:"delay_die_sides"`
	DelayMultiplierMinutes  int                 `json:"delay_multiplier_minutes"`
	RechangeAfterEffect     bool                `json:"rechange_after_effect"`
	EndBehavior             string              `json:"end_behavior"`
	StageChange             int                 `json:"stage_change"`
	StageDirections         []string            `json:"stage_directions"`
	WindDirectionChangeable bool                `json:"wind_direction_changeable"`
	Precipitation           []SpellWeatherStage `json:"precipitation"`
	Temperature             []SpellWeatherStage `json:"temperature"`
	Wind                    []SpellWeatherStage `json:"wind"`
}

// SpellCounter retains Counterspell's typed SRD source facts without implementing them.
type SpellCounter struct {
	TriggerRangeFeet  int      `json:"trigger_range_feet"`
	TriggerComponents []string `json:"trigger_components"`
	TriggerVisibility string   `json:"trigger_visibility"`
	Save              string   `json:"save"`
	FailureResult     string   `json:"failure_result"`
	WastedCosts       []string `json:"wasted_costs"`
	SlotRefund        string   `json:"slot_refund"`
}

// SpellProvisions retains Create Food and Water's typed SRD source facts without implementing them.
type SpellProvisions struct {
	FoodPounds     int      `json:"food_pounds"`
	WaterGallons   int      `json:"water_gallons"`
	Placement      []string `json:"placement"`
	FoodQuality    string   `json:"food_quality"`
	FoodForm       string   `json:"food_form"`
	WaterQuality   string   `json:"water_quality"`
	SpoilHours     int      `json:"spoil_hours"`
	SpoilCondition string   `json:"spoil_condition"`
}

// SpellUndeadOptionsEntry is one creature choice at a higher-level slot.
type SpellUndeadOptionsEntry struct {
	Creature string `json:"creature"`
	Count    int    `json:"count"`
}

// SpellUndeadSlotOptionsEntry lists the choices printed for one higher-level slot.
type SpellUndeadSlotOptionsEntry struct {
	SlotLevel int                       `json:"slot_level"`
	Options   []SpellUndeadOptionsEntry `json:"options"`
}

// SpellUndead retains Create Undead's typed SRD source facts without implementing them.
type SpellUndead struct {
	NightOnly             bool                          `json:"night_only"`
	CorpseMax             int                           `json:"corpse_max"`
	CorpseSizes           []string                      `json:"corpse_sizes"`
	CorpseType            string                        `json:"corpse_type"`
	Result                string                        `json:"result"`
	StatBlockReference    string                        `json:"stat_block_reference"`
	CommandAction         string                        `json:"command_action"`
	CommandRangeFeet      int                           `json:"command_range_feet"`
	CommandScope          []string                      `json:"command_scope"`
	CommandDecides        []string                      `json:"command_decides"`
	GeneralCommandExample string                        `json:"general_command_example"`
	NoCommandBehavior     string                        `json:"no_command_behavior"`
	OrderPersistence      string                        `json:"order_persistence"`
	ControlHours          int                           `json:"control_hours"`
	ControlEnds           string                        `json:"control_ends"`
	Reassert              string                        `json:"reassert"`
	ReassertCreaturesMax  int                           `json:"reassert_creatures_max"`
	ReassertInstead       string                        `json:"reassert_instead"`
	OnyxPerCorpse         bool                          `json:"onyx_per_corpse"`
	OnyxMinimumCostCP     int                           `json:"onyx_minimum_cost_cp"`
	SlotOptions           []SpellUndeadSlotOptionsEntry `json:"slot_options"`
}

// SpellWaterCreate retains Create or Destroy Water's typed SRD source facts without implementing them.
type SpellWaterCreate struct {
	Choices                   []string `json:"choices"`
	CreateGallons             int      `json:"create_gallons"`
	CreateContainer           string   `json:"create_container"`
	CreateAlternative         string   `json:"create_alternative"`
	DestroyGallons            int      `json:"destroy_gallons"`
	DestroyContainer          string   `json:"destroy_container"`
	DestroyAlternative        string   `json:"destroy_alternative"`
	CubeFeet                  int      `json:"cube_feet"`
	AdditionalGallonsPerSlot  int      `json:"additional_gallons_per_slot"`
	AdditionalCubeFeetPerSlot int      `json:"additional_cube_feet_per_slot"`
}

// SpellCreationMaterialsEntry is one row of Creation's Materials table.
type SpellCreationMaterialsEntry struct {
	Material        string `json:"material"`
	DurationSeconds int    `json:"duration_seconds"`
}

// SpellCreation retains Creation's typed SRD source facts without implementing them.
type SpellCreation struct {
	Origin                    string                        `json:"origin"`
	MatterKinds               []string                      `json:"matter_kinds"`
	VegetableExamples         []string                      `json:"vegetable_examples"`
	MineralExamples           []string                      `json:"mineral_examples"`
	MaxCubeFeet               int                           `json:"max_cube_feet"`
	FormRequirement           string                        `json:"form_requirement"`
	DurationRule              string                        `json:"duration_rule"`
	MultipleMaterials         string                        `json:"multiple_materials"`
	MaterialComponentUse      string                        `json:"material_component_use"`
	Materials                 []SpellCreationMaterialsEntry `json:"materials"`
	AdditionalCubeFeetPerSlot int                           `json:"additional_cube_feet_per_slot"`
}

func spell201Templates() map[string]any {
	return map[string]any{
		"woodland":      SpellWoodland{},
		"other_plane":   SpellOtherPlane{},
		"contagion":     SpellContagion{},
		"contingency":   SpellContingency{},
		"flame":         SpellFlame{},
		"water_control": SpellWaterControl{},
		"weather":       SpellWeather{},
		"counter":       SpellCounter{},
		"provisions":    SpellProvisions{},
		"undead":        SpellUndead{},
		"water_create":  SpellWaterCreate{},
		"creation":      SpellCreation{},
	}
}

// validSpell201Scaling checks the two scaling kinds introduced by this batch.
// Each references a typed increment or table and a named, unavailable dependency.
func (s SpellRecord) validSpell201Scaling(effects map[string]bool) bool {
	sc, d := s.Scaling, s.SourceDetails
	if d == nil || s.Level == 0 || sc.AboveLevel != s.Level || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) == 0 || !uniqueNames(sc.EffectIDs) {
		return false
	}
	for _, id := range sc.EffectIDs {
		if !effects[id] {
			return false
		}
	}
	switch sc.Kind {
	case "slot-increment":
		water := d.WaterCreate != nil && d.WaterCreate.AdditionalGallonsPerSlot > 0 && d.WaterCreate.AdditionalCubeFeetPerSlot > 0
		creation := d.Creation != nil && d.Creation.AdditionalCubeFeetPerSlot > 0
		return (water || creation) && slices.Contains(s.RequiredMechanics, "scaling.per-slot-increment")
	case "slot-creature-table":
		return d.Undead != nil && len(d.Undead.SlotOptions) > 0 && slices.Contains(s.RequiredMechanics, "scaling.creature-count-by-slot")
	}
	return false
}

func (s SpellRecord) validateSpell201() (bool, error) {
	names := map[string]string{"woodland": "conjure-woodland-beings", "other_plane": "contact-other-plane", "contagion": "contagion", "contingency": "contingency", "flame": "continual-flame", "water_control": "control-water", "weather": "control-weather", "counter": "counterspell", "provisions": "create-food-and-water", "undead": "create-undead", "water_create": "create-or-destroy-water", "creation": "creation"}
	expected := ""
	for key, name := range names {
		if s.ID == prefix+"spells/"+name {
			expected = key
		}
	}
	if s.SourceDetails == nil {
		if expected != "" {
			return true, fmt.Errorf("batch 201 details required")
		}
		return false, nil
	}
	d := s.SourceDetails
	payloads := map[string]any{
		"woodland":      d.Woodland,
		"other_plane":   d.OtherPlane,
		"contagion":     d.Contagion,
		"contingency":   d.Contingency,
		"flame":         d.Flame,
		"water_control": d.WaterControl,
		"weather":       d.Weather,
		"counter":       d.Counter,
		"provisions":    d.Provisions,
		"undead":        d.Undead,
		"water_create":  d.WaterCreate,
		"creation":      d.Creation,
	}
	count := 0
	for key, value := range payloads {
		if reflect.ValueOf(value).IsNil() {
			continue
		}
		count++
		if expected != key {
			return true, fmt.Errorf("mismatched batch 201 payload")
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
	if count != 1 || d.sourcePayloadCount() != 1 {
		return true, fmt.Errorf("exactly one batch 201 payload required")
	}
	castingSeconds := map[string]int{"action": 0, "reaction": 0, "1-minute": 60, "10-minutes": 600}
	seconds, ok := castingSeconds[s.CastingTime]
	if !ok || s.Scope != "complete-source-record" || d.CastingSeconds != seconds || d.CastingTriggers == nil || strings.TrimSpace(d.EndTiming) == "" || (s.CastingTime == "reaction") != (len(d.CastingTriggers) == 1) {
		return true, fmt.Errorf("invalid batch 201 source details")
	}
	for _, trigger := range d.CastingTriggers {
		if strings.TrimSpace(trigger) == "" {
			return true, fmt.Errorf("blank casting trigger")
		}
	}
	// A special duration is never concentration and has no fixed length.
	if (s.Duration == "special") != (expected == "creation") {
		return true, fmt.Errorf("invalid special duration")
	}
	if w := d.Woodland; w != nil {
		if !slices.Equal(w.SourcePages, []int{118, 119}) || w.EmanationFeet != 10 || w.Save != "wisdom" || w.SavesPerTurn != 1 || w.BonusAction != "disengage" || !slices.Equal(w.Triggers, []string{"emanation-enters-creature-space", "creature-enters-emanation", "creature-ends-turn-in-emanation"}) {
			return true, fmt.Errorf("invalid woodland pages, emanation, save or triggers")
		}
	}
	if p := d.OtherPlane; p != nil {
		if p.Save != "intelligence" || p.SaveDC != 15 || p.SuccessQuestionsMax != 5 || p.FailureDamageDiceCount != 6 || p.FailureDamageDieSides != 6 || p.FailureDamageType != "psychic" || !s.ritualFlag() || !slices.Equal(p.AnswerExamples, []string{"yes", "no", "maybe", "never", "irrelevant", "unclear"}) || !slices.Equal(p.FailureEnds, []string{"finish-long-rest", "greater-restoration-cast-on-you"}) {
			return true, fmt.Errorf("invalid contact other plane save, damage or answers")
		}
	}
	if c := d.Contagion; c != nil {
		if c.Save != "constitution" || c.EndEffectSave != "constitution" || c.SuccessesToEnd != 3 || c.FailuresToLast != 3 || c.DurationDays != 7 || s.DurationSeconds != c.DurationDays*86400 {
			return true, fmt.Errorf("invalid contagion save, thresholds or duration")
		}
	}
	if c := d.Contingency; c != nil {
		m := s.Components.Material
		if c.ContingentMaxLevel != 5 || c.ContingentCastingTime != "action" || c.MaxActive != 1 || !c.DeferredOnCast || !c.RegardlessOfWish || !c.EndsAfterTrigger || m == nil || *m.CostCP != 150000 || *m.Consumed || s.DurationSeconds != 10*86400 {
			return true, fmt.Errorf("invalid contingency limits, statuette or duration")
		}
	}
	if f := d.Flame; f != nil {
		m := s.Components.Material
		if f.BrightLightRadiusFeet != 20 || f.DimLightAdditionalFeet != 20 || s.Duration != "until-dispelled" || m == nil || *m.CostCP != 5000 || !*m.Consumed {
			return true, fmt.Errorf("invalid continual flame light, duration or ruby dust")
		}
	}
	if w := d.WaterControl; w != nil {
		f, l := w.Flood, w.Whirlpool
		if !slices.Equal(w.SourcePages, []int{119, 120}) || w.CubeMaxFeet != 100 || w.RepeatAction != "magic" || !slices.Equal(w.RepeatChoices, []string{"same-effect", "different-effect"}) || !slices.Equal(w.EffectChoices, []string{"flood", "part-water", "redirect-flow", "whirlpool"}) || f.StandingWaterRiseFeet != 20 || f.WaveHeightFeet != 20 || f.CapsizePercent != 25 || l.MinimumAreaSquareFeet != 50 || l.MinimumDepthFeet != 25 || l.BaseWidthFeet != 5 || l.TopWidthFeet != 50 || l.HeightFeet != 25 || l.PullRadiusFeet != 25 || l.PullFeet != 10 || l.Save != "strength" || !slices.Equal(l.SaveTriggers, []string{"first-entry-on-turn", "ends-turn-in-whirlpool"}) || s.Range != (SpellRange{"feet", 300}) {
			return true, fmt.Errorf("invalid control water cube, flood or whirlpool")
		}
	}
	if w := d.Weather; w != nil {
		rows := func(got []SpellWeatherStage, want ...string) bool {
			if len(got) != len(want) {
				return false
			}
			for i, row := range got {
				if row.Stage != i+1 || row.Condition != want[i] {
					return false
				}
			}
			return true
		}
		if w.RadiusMiles != 5 || !w.RequiresOutdoors || !w.EndsIndoors || w.DelayDiceCount != 1 || w.DelayDieSides != 4 || w.DelayMultiplierMinutes != 10 || w.StageChange != 1 || !slices.Equal(w.Changeable, []string{"precipitation", "temperature", "wind"}) || !slices.Equal(w.StageDirections, []string{"up", "down"}) || s.DurationSeconds != 8*3600 ||
			!rows(w.Precipitation, "clear", "light-clouds", "overcast-or-ground-fog", "rain-hail-or-snow", "torrential-rain-driving-hail-or-blizzard") ||
			!rows(w.Temperature, "heat-wave", "hot", "warm", "cool", "cold", "freezing") ||
			!rows(w.Wind, "calm", "moderate-wind", "strong-wind", "gale", "storm") {
			return true, fmt.Errorf("invalid control weather area, delay or tables")
		}
	}
	if c := d.Counter; c != nil {
		if c.TriggerRangeFeet != 60 || c.Save != "constitution" || s.CastingTime != "reaction" || s.Range != (SpellRange{"feet", 60}) || !slices.Equal(c.TriggerComponents, []string{"verbal", "somatic", "material"}) || !slices.Equal(c.WastedCosts, []string{"action", "bonus-action", "reaction"}) {
			return true, fmt.Errorf("invalid counterspell trigger, save or costs")
		}
	}
	if p := d.Provisions; p != nil {
		if p.FoodPounds != 45 || p.WaterGallons != 30 || p.SpoilHours != 24 || !slices.Equal(p.Placement, []string{"ground", "containers-in-range"}) {
			return true, fmt.Errorf("invalid create food and water quantities")
		}
	}
	if u := d.Undead; u != nil {
		m := s.Components.Material
		want := []SpellUndeadSlotOptionsEntry{
			{7, []SpellUndeadOptionsEntry{{"ghoul", 4}}},
			{8, []SpellUndeadOptionsEntry{{"ghoul", 5}, {"ghast-or-wight", 2}}},
			{9, []SpellUndeadOptionsEntry{{"ghoul", 6}, {"ghast-or-wight", 3}, {"mummy", 2}}},
		}
		if !u.NightOnly || u.CorpseMax != 3 || u.CommandAction != "bonus-action" || u.CommandRangeFeet != 120 || u.ControlHours != 24 || u.ReassertCreaturesMax != 3 || !u.OnyxPerCorpse || m == nil || *m.CostCP != u.OnyxMinimumCostCP || u.OnyxMinimumCostCP != 15000 || *m.Consumed || !slices.Equal(u.CorpseSizes, []string{"medium", "small"}) || !reflect.DeepEqual(u.SlotOptions, want) {
			return true, fmt.Errorf("invalid create undead corpses, control or slot options")
		}
	}
	if w := d.WaterCreate; w != nil {
		if !slices.Equal(w.Choices, []string{"create-water", "destroy-water"}) || w.CreateGallons != 10 || w.DestroyGallons != 10 || w.CubeFeet != 30 || w.AdditionalGallonsPerSlot != 10 || w.AdditionalCubeFeetPerSlot != 5 {
			return true, fmt.Errorf("invalid create or destroy water quantities")
		}
	}
	if c := d.Creation; c != nil {
		want := []SpellCreationMaterialsEntry{{"vegetable-matter", 86400}, {"stone-or-crystal", 43200}, {"precious-metals", 3600}, {"gems", 600}, {"adamantine-or-mithral", 60}}
		if c.MaxCubeFeet != 5 || c.AdditionalCubeFeetPerSlot != 5 || !slices.Equal(c.MatterKinds, []string{"vegetable-matter", "mineral-matter"}) || !reflect.DeepEqual(c.Materials, want) || !slices.Contains(s.RequiredMechanics, "duration.special-by-material") {
			return true, fmt.Errorf("invalid creation cube or materials table")
		}
	}
	return true, nil
}

func (s SpellRecord) ritualFlag() bool { return s.Ritual != nil && *s.Ritual }
