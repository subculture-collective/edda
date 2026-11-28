package game

import (
	"strings"
	"testing"
)

func TestDeclaredFallFailsClosed(t *testing.T) {
	const ok = `"id":"f1","height_feet":30,"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false`
	tests := map[string]struct{ fall, reason string }{
		"supported":           {`{` + ok + `}`, ""},
		"unknown fact":        {`{` + ok + `,"spikes":true}`, "requires adjudication"},
		"missing height":      {`{"id":"f1","landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false}`, "Clarify the exact height"},
		"negative height":     {`{"id":"f1","height_feet":-5,"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false}`, "Clarify the exact height"},
		"missing landing":     {`{"id":"f1","height_feet":30,"flying":false,"onto_creature":false,"fall_prevention":false}`, "Clarify what the faller lands on"},
		"other liquid":        {`{"id":"f1","height_feet":30,"landing":"acid","flying":false,"onto_creature":false,"fall_prevention":false}`, "requires adjudication"},
		"eligibility unknown": {`{"id":"f1","height_feet":30,"landing":"water"}`, "Clarify whether"},
		"flying":              {`{"id":"f1","height_feet":30,"landing":"solid","flying":true,"onto_creature":false,"fall_prevention":false}`, "require adjudication"},
		"onto a creature":     {`{"id":"f1","height_feet":30,"landing":"solid","flying":false,"onto_creature":true,"fall_prevention":false}`, "require adjudication"},
		"magic prevention":    {`{"id":"f1","height_feet":30,"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":true}`, "require adjudication"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fall, declared, err := readDeclaredFall([]byte(`{"srd521_fall":` + tt.fall + `}`))
			if !declared {
				t.Fatal("declaration not found")
			}
			if err == nil {
				err = fall.Validate()
			}
			if tt.reason == "" && err != nil {
				t.Fatal(err)
			}
			if tt.reason != "" && (err == nil || !strings.Contains(err.Error(), tt.reason)) {
				t.Fatalf("got %v, want %q", err, tt.reason)
			}
		})
	}
	if _, declared, _ := readDeclaredFall([]byte(`{"srd521":{}}`)); declared {
		t.Fatal("absent declaration reported")
	}
}

func TestFallSkillChoices(t *testing.T) {
	for _, tt := range []struct {
		landing, skill, ability string
		ok                      bool
	}{
		{"solid", "", "", true},
		{"solid", "athletics", "", false},
		{"water", "", "", true},
		{"water", "athletics", "strength", true},
		{"water", "acrobatics", "dexterity", true},
		{"water", "perception", "", false},
	} {
		if ability, ok := fallSkill(tt.landing, tt.skill); ability != tt.ability || ok != tt.ok {
			t.Errorf("fallSkill(%q,%q)=%q,%t", tt.landing, tt.skill, ability, ok)
		}
	}
}

func TestReadConditions(t *testing.T) {
	if !ReadConditions([]byte(`{"srd521_conditions":{"prone":true}}`)).Prone || ReadConditions([]byte(`{}`)).Prone || ReadConditions([]byte(`not json`)).Prone {
		t.Fatal("conditions")
	}
}
