package srd521

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// SpellMassHeal retains typed SRD source facts without implementing them.
type SpellMassHeal struct {
	HitPointPool      int      `json:"hit_point_pool"`
	Division          string   `json:"division"`
	Targets           string   `json:"targets"`
	ConditionsRemoved []string `json:"conditions_removed"`
	ConditionTargets  string   `json:"condition_targets"`
}

// SpellMassHealingWord retains typed SRD source facts without implementing them.
type SpellMassHealingWord struct {
	MaxTargets       int    `json:"max_targets"`
	Targets          string `json:"targets"`
	HealingDiceCount int    `json:"healing_dice_count"`
	HealingDieSides  int    `json:"healing_die_sides"`
	HealingAdds      string `json:"healing_adds"`
	SlotDiceCount    int    `json:"slot_dice_count"`
	SlotDieSides     int    `json:"slot_die_sides"`
}

// SpellMassSuggestionSlotDuration is one source-listed row.
type SpellMassSuggestionSlotDuration struct {
	MinSlot int `json:"min_slot"`
	MaxSlot int `json:"max_slot"`
	Seconds int `json:"seconds"`
}

// SpellMassSuggestion retains typed SRD source facts without implementing them.
type SpellMassSuggestion struct {
	Suggestion             string                            `json:"suggestion"`
	SuggestionMaxWords     int                               `json:"suggestion_max_words"`
	MaxTargets             int                               `json:"max_targets"`
	SuggestionRequirements []string                          `json:"suggestion_requirements"`
	SuggestionExamples     []string                          `json:"suggestion_examples"`
	Save                   string                            `json:"save"`
	FailureCondition       string                            `json:"failure_condition"`
	ConditionEnds          []string                          `json:"condition_ends"`
	Pursuit                string                            `json:"pursuit"`
	ActivityDuration       string                            `json:"activity_duration"`
	Completion             string                            `json:"completion"`
	SlotDurations          []SpellMassSuggestionSlotDuration `json:"slot_durations"`
}

// SpellMaze retains typed SRD source facts without implementing them.
type SpellMaze struct {
	Destination    string `json:"destination"`
	Remains        string `json:"remains"`
	EscapeAction   string `json:"escape_action"`
	EscapeCheck    string `json:"escape_check"`
	EscapeDC       int    `json:"escape_dc"`
	EscapeSuccess  string `json:"escape_success"`
	ReturnSpace    string `json:"return_space"`
	OccupiedReturn string `json:"occupied_return"`
}

// SpellMeldIntoStone retains typed SRD source facts without implementing them.
type SpellMeldIntoStone struct {
	StoneRequirement    string   `json:"stone_requirement"`
	Merged              []string `json:"merged"`
	TouchRequirement    string   `json:"touch_requirement"`
	Concealment         string   `json:"concealment"`
	OutsideSight        string   `json:"outside_sight"`
	OutsideHearing      string   `json:"outside_hearing"`
	TimeAwareness       string   `json:"time_awareness"`
	SelfCasting         string   `json:"self_casting"`
	ExitMovementFeet    int      `json:"exit_movement_feet"`
	Exit                string   `json:"exit"`
	OtherMovement       string   `json:"other_movement"`
	MinorDamage         string   `json:"minor_damage"`
	PartialTriggers     []string `json:"partial_triggers"`
	PartialDiceCount    int      `json:"partial_dice_count"`
	PartialDieSides     int      `json:"partial_die_sides"`
	PartialDamageType   string   `json:"partial_damage_type"`
	DestroyedTriggers   []string `json:"destroyed_triggers"`
	DestroyedDamage     int      `json:"destroyed_damage"`
	DestroyedDamageType string   `json:"destroyed_damage_type"`
	ExpelledSpace       string   `json:"expelled_space"`
	ExpelledCondition   string   `json:"expelled_condition"`
}

// SpellMending retains typed SRD source facts without implementing them.
type SpellMending struct {
	Repairs          string   `json:"repairs"`
	Examples         []string `json:"examples"`
	MaxDimensionFeet int      `json:"max_dimension_feet"`
	Result           string   `json:"result"`
	MagicItem        string   `json:"magic_item"`
}

// SpellMessage retains typed SRD source facts without implementing them.
type SpellMessage struct {
	SourcePages         []int    `json:"source_pages"`
	Gesture             string   `json:"gesture"`
	Delivery            string   `json:"delivery"`
	Reply               string   `json:"reply"`
	ThroughObjects      string   `json:"through_objects"`
	ThroughRequirements []string `json:"through_requirements"`
	Blockers            []string `json:"blockers"`
	SolidBlockerFeet    int      `json:"solid_blocker_feet"`
}

// SpellMeteorSwarm retains typed SRD source facts without implementing them.
type SpellMeteorSwarm struct {
	RangeMiles           int    `json:"range_miles"`
	Orbs                 string `json:"orbs"`
	PointCount           int    `json:"point_count"`
	Points               string `json:"points"`
	SphereRadiusFeet     int    `json:"sphere_radius_feet"`
	Save                 string `json:"save"`
	FireDiceCount        int    `json:"fire_dice_count"`
	FireDieSides         int    `json:"fire_die_sides"`
	BludgeoningDiceCount int    `json:"bludgeoning_dice_count"`
	BludgeoningDieSides  int    `json:"bludgeoning_die_sides"`
	SaveSuccess          string `json:"save_success"`
	Overlap              string `json:"overlap"`
	ObjectEligibility    string `json:"object_eligibility"`
	ObjectEffect         string `json:"object_effect"`
	FlammableObject      string `json:"flammable_object"`
}

// SpellMindBlank retains typed SRD source facts without implementing them.
type SpellMindBlank struct {
	Target            string   `json:"target"`
	DamageImmunity    string   `json:"damage_immunity"`
	ConditionImmunity string   `json:"condition_immunity"`
	UnaffectedBy      []string `json:"unaffected_by"`
	SpellsCannot      []string `json:"spells_cannot"`
	IncludesSpell     string   `json:"includes_spell"`
}

// SpellMindSpike retains typed SRD source facts without implementing them.
type SpellMindSpike struct {
	Target            string   `json:"target"`
	Save              string   `json:"save"`
	DamageDiceCount   int      `json:"damage_dice_count"`
	DamageDieSides    int      `json:"damage_die_sides"`
	DamageType        string   `json:"damage_type"`
	SaveSuccess       string   `json:"save_success"`
	FailureKnowledge  string   `json:"failure_knowledge"`
	KnowledgeLimit    string   `json:"knowledge_limit"`
	KnowledgeBenefits []string `json:"knowledge_benefits"`
	SlotDiceCount     int      `json:"slot_dice_count"`
	SlotDieSides      int      `json:"slot_die_sides"`
}

// SpellMinorIllusion retains typed SRD source facts without implementing them.
type SpellMinorIllusion struct {
	Options                  []string `json:"options"`
	Recast                   string   `json:"recast"`
	DiscernAction            string   `json:"discern_action"`
	DiscernCheck             string   `json:"discern_check"`
	DiscernDC                string   `json:"discern_dc"`
	Discerned                string   `json:"discerned"`
	SoundVolume              string   `json:"sound_volume"`
	SoundExamples            []string `json:"sound_examples"`
	SoundTiming              []string `json:"sound_timing"`
	ImageSubject             string   `json:"image_subject"`
	ImageExamples            []string `json:"image_examples"`
	ImageMaxCubeFeet         int      `json:"image_max_cube_feet"`
	ImageCannotCreate        []string `json:"image_cannot_create"`
	ImagePhysicalInteraction string   `json:"image_physical_interaction"`
}

// SpellMirageArcane retains typed SRD source facts without implementing them.
type SpellMirageArcane struct {
	SourcePages          []int    `json:"source_pages"`
	AreaMaxMilesSquare   int      `json:"area_max_miles_square"`
	AlteredSenses        []string `json:"altered_senses"`
	Disguise             string   `json:"disguise"`
	Examples             []string `json:"examples"`
	Structures           string   `json:"structures"`
	Creatures            string   `json:"creatures"`
	Elements             []string `json:"elements"`
	MovementEffect       string   `json:"movement_effect"`
	RemovedPieces        string   `json:"removed_pieces"`
	RemovedPieceExamples []string `json:"removed_piece_examples"`
	Truesight            string   `json:"truesight"`
	TruesightLimit       string   `json:"truesight_limit"`
}

func spell213Templates() map[string]any {
	return map[string]any{
		"mass_heal":         SpellMassHeal{},
		"mass_healing_word": SpellMassHealingWord{},
		"mass_suggestion":   SpellMassSuggestion{},
		"maze":              SpellMaze{},
		"meld_into_stone":   SpellMeldIntoStone{},
		"mending":           SpellMending{},
		"message":           SpellMessage{},
		"meteor_swarm":      SpellMeteorSwarm{},
		"mind_blank":        SpellMindBlank{},
		"mind_spike":        SpellMindSpike{},
		"minor_illusion":    SpellMinorIllusion{},
		"mirage_arcane":     SpellMirageArcane{},
	}
}

func validSpell213Target(t SpellTarget) bool {
	switch t.Shape {
	case "up-to-count":
		return (t.Count == 6 && t.Eligibility == "visible-creatures-in-range") || (t.Count == 12 && t.Eligibility == "visible-hearing-understanding-creatures-in-range")
	case "single":
		return t.Count == 1 && t.Eligibility == "stone-object-or-surface-touch"
	case "spheres":
		return t.Count == 4 && t.Eligibility == "creatures-in-spheres-on-visible-points"
	case "choice":
		return t.Count == 1 && t.Eligibility == "sound-or-object-image-in-range"
	case "square":
		return t.Count == 0 && t.Eligibility == "terrain-area-up-to-1-mile-square"
	}
	return false
}

// spell213DurationScaling reports whether a typed slot table or increment backs
// this batch's slot-duration scaling.
func spell213DurationScaling(d *SpellSourceDetails, mechanics []string) bool {
	return (d.MassSuggestion != nil && len(d.MassSuggestion.SlotDurations) == 3 && slices.Contains(mechanics, "scaling.duration-by-slot"))
}

// has213Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has213Payload() bool {
	return d != nil && (d.MassHeal != nil || d.MassHealingWord != nil || d.MassSuggestion != nil || d.Maze != nil || d.MeldIntoStone != nil || d.Mending != nil ||
		d.Message != nil || d.MeteorSwarm != nil || d.MindBlank != nil || d.MindSpike != nil || d.MinorIllusion != nil || d.MirageArcane != nil)
}

// validateSpell213 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell213() (bool, error) {
	pin, ok := spell213Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has213Payload() {
			return true, fmt.Errorf("mismatched batch 213 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 213 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 213 record differs from the printed source")
	}
	return true, nil
}
