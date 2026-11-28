package combat

import (
	"slices"
	"testing"
)

func TestSightLightBoundariesAndDarkvision(t *testing.T) {
	sight, blocked := true, false
	at := func(feet int, light, obscurement, cover string, los *bool) Space {
		return Space{DistanceFeet: &feet, LineOfSight: los, Cover: cover, Light: light, Obscurement: obscurement}
	}
	none, dv60 := &Senses{}, &Senses{DarkvisionFeet: 60}
	for _, tc := range []struct {
		name         string
		space        Space
		senses       *Senses
		seen         bool
		obscurement  string
		circumstance string
		darkvision   bool
	}{
		{"bright", at(30, LightBright, ObscurementNone, "none", &sight), nil, true, "none", SightOrdinary, false},
		{"dim without darkvision", at(30, LightDim, ObscurementNone, "none", &sight), none, true, "light", SightDisadvantage, false},
		{"dark without darkvision", at(30, LightDark, ObscurementNone, "none", &sight), none, false, "heavy", SightAutomaticFail, false},
		{"dim at darkvision edge", at(60, LightDim, ObscurementNone, "none", &sight), dv60, true, "none", SightOrdinary, true},
		{"dark at darkvision edge", at(60, LightDark, ObscurementNone, "none", &sight), dv60, true, "light", SightDisadvantage, true},
		{"dim beyond darkvision", at(65, LightDim, ObscurementNone, "none", &sight), dv60, true, "light", SightDisadvantage, false},
		{"dark beyond darkvision", at(65, LightDark, ObscurementNone, "none", &sight), dv60, false, "heavy", SightAutomaticFail, false},
		{"darkvision ignores fog", at(10, LightDim, ObscurementHeavy, "none", &sight), dv60, false, "heavy", SightAutomaticFail, true},
		{"light fog in bright light", at(10, LightBright, ObscurementLight, "none", &sight), nil, true, "light", SightDisadvantage, false},
		{"total cover", at(10, LightBright, ObscurementNone, "total", &sight), nil, false, "none", SightAutomaticFail, false},
		{"half cover does not block sight", at(10, LightBright, ObscurementNone, "half", &sight), nil, true, "none", SightOrdinary, false},
		{"no line of sight", at(10, LightBright, ObscurementNone, "none", &blocked), nil, false, "none", SightAutomaticFail, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.space.Sight(tc.senses)
			if err != nil {
				t.Fatal(err)
			}
			if got.Seen != tc.seen || got.Obscurement != tc.obscurement || got.Circumstance != tc.circumstance || got.Darkvision != tc.darkvision {
				t.Fatalf("sight=%+v", got)
			}
			if tc.darkvision != slices.Contains(got.SourceIDs, "srd-5.2.1:content/rules-glossary/darkvision") {
				t.Fatalf("darkvision source=%v", got.SourceIDs)
			}
		})
	}
}

func TestSightUnresolvedFactsFailClosed(t *testing.T) {
	feet, sight := 10, true
	base := Space{DistanceFeet: &feet, LineOfSight: &sight, Cover: "none", Light: LightDim, Obscurement: ObscurementNone}
	for name, mutate := range map[string]func(*Space) (*Senses, bool){
		"legacy facts without light": func(s *Space) (*Senses, bool) { s.Light, s.Obscurement = "", ""; return &Senses{}, true },
		"unknown light":              func(s *Space) (*Senses, bool) { s.Light = "twilight"; return &Senses{}, true },
		"unknown obscurement":        func(s *Space) (*Senses, bool) { s.Obscurement = ""; return &Senses{}, true },
		"unknown distance":           func(s *Space) (*Senses, bool) { s.DistanceFeet = nil; return &Senses{}, true },
		"unknown line of sight":      func(s *Space) (*Senses, bool) { s.LineOfSight = nil; return &Senses{}, true },
		"unknown cover":              func(s *Space) (*Senses, bool) { s.Cover = ""; return &Senses{}, true },
		"unknown senses in dim":      func(s *Space) (*Senses, bool) { return nil, true },
		"invalid darkvision range":   func(s *Space) (*Senses, bool) { return &Senses{DarkvisionFeet: -5}, true },
		"senses irrelevant in bright": func(s *Space) (*Senses, bool) {
			s.Light = LightBright
			return nil, false
		},
	} {
		t.Run(name, func(t *testing.T) {
			space := base
			senses, wantErr := mutate(&space)
			if _, err := space.Sight(senses); (err != nil) != wantErr {
				t.Fatalf("err=%v", err)
			}
		})
	}
}
