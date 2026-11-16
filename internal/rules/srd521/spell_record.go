package srd521

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
)

// SpellRecords is an offline, versioned data contract. It never registers actions.
type SpellRecords struct {
	SchemaVersion int           `json:"schema_version"`
	Source        Source        `json:"source"`
	Records       []SpellRecord `json:"records"`
}

// SpellRecord describes source data or an explicitly bounded existing fixture.
// RequiredMechanics are dependencies, not resolver registrations.
type SpellRecord struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	SourceID          string          `json:"source_id"`
	SourcePage        int             `json:"source_page"`
	Scope             string          `json:"scope"`
	Level             int             `json:"level"`
	School            string          `json:"school"`
	Lists             []string        `json:"lists"`
	CastingTime       string          `json:"casting_time"`
	Range             SpellRange      `json:"range"`
	Target            SpellTarget     `json:"target"`
	Components        SpellComponents `json:"components"`
	Ritual            *bool           `json:"ritual"`
	Concentration     *bool           `json:"concentration"`
	DurationSeconds   int             `json:"duration_seconds"`
	Duration          string          `json:"duration"`
	Scaling           SpellScaling    `json:"scaling"`
	Effects           []SpellEffect   `json:"effects"`
	RequiredMechanics []string        `json:"required_mechanics"`
	SlotLevels        []int           `json:"slot_levels"`
}

// SpellRange retains the source range even for a self-only fixture.
type SpellRange struct {
	Kind string `json:"kind"`
	Feet int    `json:"feet"`
}

// SpellTarget declares the shape and eligibility of a selected target.
type SpellTarget struct {
	Shape       string `json:"shape"`
	Count       int    `json:"count"`
	Eligibility string `json:"eligibility"`
}

// SpellComponents retains material cost in copper pieces and consumption.
// A nil Material means no material requirement; false flags must be explicit.
type SpellComponents struct {
	Verbal   *bool          `json:"verbal"`
	Somatic  *bool          `json:"somatic"`
	Material *SpellMaterial `json:"material"`
}

// SpellMaterial distinguishes absent cost from a consumed component.
type SpellMaterial struct {
	Description string `json:"description"`
	CostCP      int    `json:"cost_cp"`
	Consumed    *bool  `json:"consumed"`
}

// SpellScaling identifies per-slot dice additions to each named effect.
// Bounded fixtures use none; their higher-slot catalogue data belongs elsewhere.
type SpellScaling struct {
	Kind       string   `json:"kind"`
	AboveLevel int      `json:"above_level"`
	DiceCount  int      `json:"dice_count"`
	DieSides   int      `json:"die_sides"`
	EffectIDs  []string `json:"effect_ids"`
}

// SpellEffect stores numeric payloads and explicit timing/outcome dependencies.
// HalfInitialOnly means half the initial roll on a miss, with no delayed damage.
type SpellEffect struct {
	ID                 string `json:"id"`
	Mechanic           string `json:"mechanic"`
	Timing             string `json:"timing"`
	DiceCount          int    `json:"dice_count"`
	DieSides           int    `json:"die_sides"`
	DamageType         string `json:"damage_type"`
	Bonus              int    `json:"bonus"`
	AddCastingModifier *bool  `json:"add_casting_modifier"`
	HalfInitialOnly    *bool  `json:"half_initial_only"`
}

// DecodeSpellRecords rejects malformed data and unknown JSON fields offline.
// Unknown mechanic names are retained but remain unavailable.
func DecodeSpellRecords(r io.Reader) (SpellRecords, error) {
	var records SpellRecords
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	if err := d.Decode(&records); err != nil {
		return records, fmt.Errorf("decode spell records: %w", err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return SpellRecords{}, fmt.Errorf("spell records require one JSON object")
	}
	if err := records.Validate(); err != nil {
		return SpellRecords{}, err
	}
	return records, nil
}

// Validate checks structure, provenance and internal dependencies, not conformance.
func (r SpellRecords) Validate() error {
	if r.SchemaVersion != 1 {
		return fmt.Errorf("unsupported spell schema %d", r.SchemaVersion)
	}
	if err := validateSource(r.Source); err != nil {
		return err
	}
	if len(r.Records) == 0 {
		return fmt.Errorf("spell records required")
	}
	seen := map[string]bool{}
	for _, s := range r.Records {
		if seen[s.ID] || !strings.HasPrefix(s.ID, prefix+"spells/") || !contentID.MatchString(s.ID) || s.SourceID != s.ID || s.SourcePage < 107 || s.SourcePage > 175 || strings.TrimSpace(s.Name) == "" {
			return fmt.Errorf("invalid, duplicate or source-less spell %q", s.ID)
		}
		seen[s.ID] = true
		if err := s.validate(); err != nil {
			return fmt.Errorf("%s: %w", s.ID, err)
		}
	}
	return nil
}

func (s SpellRecord) validate() error {
	if !slices.Contains([]string{"complete-source-record", "bounded-self-level-1"}, s.Scope) || s.Level < 0 || s.Level > 9 || !slices.Contains(strings.Fields("abjuration conjuration divination enchantment evocation illusion necromancy transmutation"), s.School) || len(s.Lists) == 0 || !uniqueNames(s.Lists) {
		return fmt.Errorf("invalid scope, level, school or lists")
	}
	for _, class := range s.Lists {
		if !slices.Contains(strings.Fields("bard cleric druid paladin ranger sorcerer warlock wizard"), class) {
			return fmt.Errorf("invalid spell list")
		}
	}
	if !slices.Contains([]string{"action", "bonus-action"}, s.CastingTime) || !slices.Contains([]string{"feet", "touch", "self"}, s.Range.Kind) || s.Range.Feet < 0 || (s.Range.Kind == "feet" && s.Range.Feet == 0) || (s.Range.Kind != "feet" && s.Range.Feet != 0) || s.Target.Shape != "single" || s.Target.Count != 1 || !slices.Contains([]string{"target-in-range", "creature-in-range", "conscious-self"}, s.Target.Eligibility) {
		return fmt.Errorf("invalid casting, range or target shape")
	}
	if s.Components.Verbal == nil || s.Components.Somatic == nil || s.Ritual == nil || s.Concentration == nil {
		return fmt.Errorf("explicit component, ritual and concentration flags required")
	}
	if m := s.Components.Material; m != nil {
		if strings.TrimSpace(m.Description) == "" || m.CostCP < 0 || m.Consumed == nil {
			return fmt.Errorf("invalid material")
		}
	}
	if !slices.Contains([]string{"instantaneous", "timed"}, s.Duration) || s.DurationSeconds < 0 || (s.Duration == "instantaneous" && (s.DurationSeconds != 0 || *s.Concentration)) || (s.Duration == "timed" && s.DurationSeconds == 0) {
		return fmt.Errorf("invalid duration")
	}
	if len(s.RequiredMechanics) == 0 || !uniqueNames(s.RequiredMechanics) || len(s.Effects) == 0 || len(s.SlotLevels) == 0 {
		return fmt.Errorf("mechanics, effects and slots required")
	}
	last := -1
	for _, level := range s.SlotLevels {
		if level < s.Level || level > 9 || level <= last {
			return fmt.Errorf("invalid slot levels")
		}
		last = level
	}
	effects := map[string]bool{}
	for _, e := range s.Effects {
		if e.ID == "" || effects[e.ID] || !slices.Contains(s.RequiredMechanics, e.Mechanic) || e.Timing == "" || e.DiceCount < 0 || (e.DiceCount > 0 && !slices.Contains([]int{4, 6, 8, 10, 12, 20}, e.DieSides)) || (e.DiceCount == 0 && e.DieSides != 0) || e.AddCastingModifier == nil || e.HalfInitialOnly == nil {
			return fmt.Errorf("invalid effect")
		}
		effects[e.ID] = true
	}
	sc := s.Scaling
	switch sc.Kind {
	case "none":
		if sc.AboveLevel != 0 || sc.DiceCount != 0 || sc.DieSides != 0 || len(sc.EffectIDs) != 0 {
			return fmt.Errorf("unexpected scaling payload")
		}
	case "slot":
		if sc.AboveLevel != s.Level || sc.DiceCount <= 0 || !slices.Contains([]int{4, 6, 8, 10, 12, 20}, sc.DieSides) || !uniqueNames(sc.EffectIDs) || len(sc.EffectIDs) == 0 || !slices.Contains(s.RequiredMechanics, "scaling.per-slot") {
			return fmt.Errorf("invalid slot scaling")
		}
		for _, id := range sc.EffectIDs {
			if !effects[id] {
				return fmt.Errorf("unknown scaling effect")
			}
		}
	default:
		return fmt.Errorf("unsupported scaling shape")
	}
	if s.Scope == "bounded-self-level-1" && (s.Level != 1 || s.Target.Eligibility != "conscious-self" || !slices.Equal(s.SlotLevels, []int{1}) || sc.Kind != "none") {
		return fmt.Errorf("bounded fixture cannot widen targets or slots")
	}
	return nil
}

func uniqueNames(names []string) bool {
	seen := map[string]bool{}
	for _, name := range names {
		if strings.TrimSpace(name) == "" || seen[name] {
			return false
		}
		seen[name] = true
	}
	return true
}

// MechanicsAvailable reports only the two existing bounded fixture dependencies.
// It is not permission to cast. Complete source records and unknown/unsupported
// keys always return false, independently of their presence in a catalogue.
func (s SpellRecord) MechanicsAvailable() bool {
	if s.validate() != nil || s.Scope != "bounded-self-level-1" {
		return false
	}
	var required []string
	switch s.ID {
	case prefix + "spells/cure-wounds":
		required = []string{"cast.cleric-self-slot-1", "healing.cure-wounds-self"}
	case prefix + "spells/shield-of-faith":
		required = []string{"cast.cleric-self-slot-1", "defense.shield-of-faith-self", "concentration.self-600-seconds"}
	default:
		return false
	}
	return s.SourceID == s.ID && s.SourcePage >= 107 && s.SourcePage <= 175 && slices.Equal(s.RequiredMechanics, required)
}
