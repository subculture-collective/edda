package srd521

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SpellWishBasicUse is one source-listed row or nested block.
type SpellWishBasicUse struct {
	Duplicates    string `json:"duplicates"`
	MaxSpellLevel int    `json:"max_spell_level"`
	Requirements  string `json:"requirements"`
	Result        string `json:"result"`
}

// SpellWishObjectCreation is one source-listed row or nested block.
type SpellWishObjectCreation struct {
	ObjectCount      int    `json:"object_count"`
	MaxValueCP       int    `json:"max_value_cp"`
	Excluded         string `json:"excluded"`
	MaxDimensionFeet int    `json:"max_dimension_feet"`
	Placement        string `json:"placement"`
}

// SpellWishInstantHealth is one source-listed row or nested block.
type SpellWishInstantHealth struct {
	Includes          string `json:"includes"`
	MaxOtherCreatures int    `json:"max_other_creatures"`
	Requirement       string `json:"requirement"`
	Healing           string `json:"healing"`
	Ends              string `json:"ends"`
}

// SpellWishResistance is one source-listed row or nested block.
type SpellWishResistance struct {
	MaxCreatures      int    `json:"max_creatures"`
	Requirement       string `json:"requirement"`
	DamageTypeCount   int    `json:"damage_type_count"`
	DamageTypeChooser string `json:"damage_type_chooser"`
	Duration          string `json:"duration"`
}

// SpellWishSpellImmunity is one source-listed row or nested block.
type SpellWishSpellImmunity struct {
	MaxCreatures    int    `json:"max_creatures"`
	Requirement     string `json:"requirement"`
	Immunity        string `json:"immunity"`
	DurationSeconds int    `json:"duration_seconds"`
}

// SpellWishSuddenLearning is one source-listed row or nested block.
type SpellWishSuddenLearning struct {
	ReplacedFeats int    `json:"replaced_feats"`
	Replacement   string `json:"replacement"`
	OldFeat       string `json:"old_feat"`
	NewFeat       string `json:"new_feat"`
	Restriction   string `json:"restriction"`
}

// SpellWishRollRedo is one source-listed row or nested block.
type SpellWishRollRedo struct {
	Undoes         string   `json:"undoes"`
	Method         string   `json:"method"`
	WindowRounds   int      `json:"window_rounds"`
	WindowIncludes string   `json:"window_includes"`
	Reality        string   `json:"reality"`
	Examples       []string `json:"examples"`
	RerollOptions  []string `json:"reroll_options"`
	ResultChoice   string   `json:"result_choice"`
}

// SpellWishReshapeReality is one source-listed row or nested block.
type SpellWishReshapeReality struct {
	Scope            string   `json:"scope"`
	Statement        string   `json:"statement"`
	Ruling           string   `json:"ruling"`
	Risk             string   `json:"risk"`
	PossibleFailures []string `json:"possible_failures"`
	WordingExample   string   `json:"wording_example"`
	LegendaryWish    string   `json:"legendary_wish"`
	WideConsequences string   `json:"wide_consequences"`
	GodWish          string   `json:"god_wish"`
	MultiverseWish   string   `json:"multiverse_wish"`
}

// SpellWishStress is one source-listed row or nested block.
type SpellWishStress struct {
	Applies                           string `json:"applies"`
	DamageTrigger                     string `json:"damage_trigger"`
	DamageDiceCount                   int    `json:"damage_dice_count"`
	DamageDieSides                    int    `json:"damage_die_sides"`
	DamageType                        string `json:"damage_type"`
	DamagePer                         string `json:"damage_per"`
	DamagePrevention                  string `json:"damage_prevention"`
	StrengthScore                     int    `json:"strength_score"`
	StrengthDaysDiceCount             int    `json:"strength_days_dice_count"`
	StrengthDaysDieSides              int    `json:"strength_days_die_sides"`
	LightActivity                     string `json:"light_activity"`
	LightActivityRecoveryDecreaseDays int    `json:"light_activity_recovery_decrease_days"`
	NeverCastAgainPercent             int    `json:"never_cast_again_percent"`
}

// SpellWish retains typed SRD source facts without implementing them.
type SpellWish struct {
	BasicUse               SpellWishBasicUse       `json:"basic_use"`
	AlternativeChooseCount int                     `json:"alternative_choose_count"`
	ObjectCreation         SpellWishObjectCreation `json:"object_creation"`
	InstantHealth          SpellWishInstantHealth  `json:"instant_health"`
	Resistance             SpellWishResistance     `json:"resistance"`
	SpellImmunity          SpellWishSpellImmunity  `json:"spell_immunity"`
	SuddenLearning         SpellWishSuddenLearning `json:"sudden_learning"`
	RollRedo               SpellWishRollRedo       `json:"roll_redo"`
	ReshapeReality         SpellWishReshapeReality `json:"reshape_reality"`
	Stress                 SpellWishStress         `json:"stress"`
}

// SpellWordOfRecall retains typed SRD source facts without implementing them.
type SpellWordOfRecall struct {
	MaxOtherCreatures int      `json:"max_other_creatures"`
	WithinFeet        int      `json:"within_feet"`
	Destination       string   `json:"destination"`
	Arrival           string   `json:"arrival"`
	Unprepared        string   `json:"unprepared"`
	Preparation       string   `json:"preparation"`
	SanctuaryExamples []string `json:"sanctuary_examples"`
}

// SpellZoneOfTruth retains typed SRD source facts without implementing them.
type SpellZoneOfTruth struct {
	RadiusFeet      int      `json:"radius_feet"`
	SaveTriggers    []string `json:"save_triggers"`
	Save            string   `json:"save"`
	SaveFailure     string   `json:"save_failure"`
	CasterKnowledge string   `json:"caster_knowledge"`
	TargetAwareness string   `json:"target_awareness"`
	TargetEvasion   string   `json:"target_evasion"`
	EvasionLimit    string   `json:"evasion_limit"`
}

func spell224Templates() map[string]any {
	return map[string]any{
		"wish":           SpellWish{},
		"word_of_recall": SpellWordOfRecall{},
		"zone_of_truth":  SpellZoneOfTruth{},
	}
}

func validSpell224Target(t SpellTarget) bool {
	switch t.Shape {
	case "choice":
		return t.Count == 1 && t.Eligibility == "one-chosen-wish-use"
	case "self-and-up-to-count":
		return t.Count == 5 && t.Eligibility == "willing-creatures-within-range"
	case "sphere":
		return t.Count == 0 && t.Eligibility == "creatures-in-truth-zone-on-point-in-range"
	}
	return false
}

// has224Payload reports whether any typed payload of this batch is selected.
func (d *SpellSourceDetails) has224Payload() bool {
	return d != nil && (d.Wish != nil || d.WordOfRecall != nil || d.ZoneOfTruth != nil)
}

// validateSpell224 requires each batch record to equal its pinned printed record
// and to carry exactly its own typed payload.
func (s SpellRecord) validateSpell224() (bool, error) {
	pin, ok := spell224Pins[strings.TrimPrefix(s.ID, prefix+"spells/")]
	if !ok {
		if s.SourceDetails.has224Payload() {
			return true, fmt.Errorf("mismatched batch 224 payload")
		}
		return false, nil
	}
	if d := s.SourceDetails; d == nil || d.sourcePayloadCount() != 1 || s.Scope != "complete-source-record" {
		return true, fmt.Errorf("batch 224 details required")
	}
	got, err := json.Marshal(s)
	if err != nil {
		return true, err
	}
	if string(got) != pin {
		return true, fmt.Errorf("batch 224 record differs from the printed source")
	}
	return true, nil
}
