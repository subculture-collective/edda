package combat

import "fmt"

// Light levels and obscurement are declared server rulings for the target's
// space (SRD 5.2.1 "Vision and Light", rules glossary). An empty value means
// unresolved; it never defaults to bright light or a clear view.
const (
	LightBright = "bright"
	LightDim    = "dim"
	LightDark   = "dark"

	ObscurementNone  = "none"
	ObscurementLight = "light"
	ObscurementHeavy = "heavy"
)

// Perception circumstances for a Wisdom (Perception) check that relies on sight.
const (
	SightOrdinary      = "ordinary"
	SightDisadvantage  = "disadvantage"
	SightAutomaticFail = "automatic-failure"
)

// Senses is an observer's declared special senses. A nil *Senses means the
// observer's senses are unresolved; DarkvisionFeet 0 is an explicit "none".
type Senses struct {
	DarkvisionFeet int `json:"darkvision_feet"`
}

// Sight is the derived view of one target from the observing actor.
type Sight struct {
	Light        string   `json:"light"`
	Obscurement  string   `json:"obscurement"`
	Seen         bool     `json:"seen"`
	Circumstance string   `json:"circumstance"`
	Darkvision   bool     `json:"darkvision,omitempty"`
	SourceIDs    []string `json:"source_ids"`
}

var lightObscurement = map[string]int{LightBright: 0, LightDim: 1, LightDark: 2}
var obscurementLevel = map[string]int{ObscurementNone: 0, ObscurementLight: 1, ObscurementHeavy: 2}
var obscurementName = []string{ObscurementNone, ObscurementLight, ObscurementHeavy}

// Sight derives whether the observer can see the target in s.Light, which is
// the light level at the target's space. A later light source (#264) raises
// that level before calling Sight; this derivation does not change.
//
// Dim light is lightly obscured and darkness heavily obscured. Darkvision
// within its range treats dim light as bright and darkness as dim; it does not
// reduce declared non-light obscurement such as fog or foliage. No line of
// sight or total cover blocks sight. Any unresolved fact returns an error.
func (s Space) Sight(observer *Senses) (Sight, error) {
	if err := s.Validate(); err != nil {
		return Sight{}, err
	}
	fromLight, ok := lightObscurement[s.Light]
	if !ok {
		return Sight{}, fmt.Errorf("Clarify the light at the target before deciding whether it can be seen.")
	}
	declared, ok := obscurementLevel[s.Obscurement]
	if !ok {
		return Sight{}, fmt.Errorf("Clarify whether fog, foliage or similar obscures the target.")
	}
	sight := Sight{Light: s.Light, SourceIDs: []string{"srd-5.2.1:content/exploration/vision-and-light"}}
	if fromLight > 0 {
		if observer == nil || observer.DarkvisionFeet < 0 {
			return Sight{}, fmt.Errorf("Clarify the observer's senses before deciding whether the target can be seen.")
		}
		if observer.DarkvisionFeet > 0 && *s.DistanceFeet <= observer.DarkvisionFeet {
			fromLight--
			sight.Darkvision = true
			sight.SourceIDs = append(sight.SourceIDs, "srd-5.2.1:content/rules-glossary/darkvision")
		}
	}
	level := max(fromLight, declared)
	sight.Obscurement = obscurementName[level]
	sight.Seen = level < 2 && *s.LineOfSight && s.Cover != "total"
	switch {
	case !sight.Seen:
		sight.Circumstance = SightAutomaticFail
		sight.SourceIDs = append(sight.SourceIDs, "srd-5.2.1:content/rules-glossary/blinded")
	case level == 1:
		sight.Circumstance = SightDisadvantage
	default:
		sight.Circumstance = SightOrdinary
	}
	return sight, nil
}

// Illuminated returns s with the light at the target's space raised by a light
// source the observer carries, which sheds bright light within brightFeet and
// dim light for dimFeet further (#264). Light never lowers a declared level.
// Unresolved light or distance is returned unchanged, so Sight still fails
// closed. The bool reports whether the source changed the level.
func (s Space) Illuminated(brightFeet, dimFeet int) (Space, bool) {
	current, ok := lightObscurement[s.Light]
	if !ok || s.DistanceFeet == nil || *s.DistanceFeet < 0 || brightFeet < 0 || dimFeet < 0 {
		return s, false
	}
	distance := *s.DistanceFeet
	switch {
	case distance <= brightFeet && current > lightObscurement[LightBright]:
		s.Light = LightBright
	case distance > brightFeet && distance <= brightFeet+dimFeet && current > lightObscurement[LightDim]:
		s.Light = LightDim
	default:
		return s, false
	}
	return s, true
}
