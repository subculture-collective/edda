package srd521

import (
	"strings"
	"testing"
)

const doorActor = "11111111-1111-1111-1111-111111111111"

func door(fields string) []byte {
	return []byte(`{"profile":"ordinary-door-v1","actor_id":"` + doorActor + `"` + fields + `}`)
}

func TestDoorForHidesForeignAndHiddenDoors(t *testing.T) {
	known := `,"hidden":false,"state":"closed","locked":false,"distance_feet":5`
	if !DoorFor(door(known), doorActor) {
		t.Fatal("known door rejected")
	}
	for name, raw := range map[string][]byte{
		"foreign":        door(known),
		"hidden":         door(`,"hidden":true,"state":"closed","locked":false,"distance_feet":5`),
		"hidden missing": door(`,"state":"closed","locked":false,"distance_feet":5`),
		"malformed":      []byte(`[`),
	} {
		actor := doorActor
		if name == "foreign" {
			actor = "22222222-2222-2222-2222-222222222222"
		}
		if DoorFor(raw, actor) {
			t.Errorf("%s door disclosed", name)
		}
	}
}

func TestParseDoorFactsFailsClosed(t *testing.T) {
	for name, tc := range map[string]struct {
		fields, reason string
	}{
		"unknown field":   {`,"hidden":false,"state":"closed","locked":false,"distance_feet":5,"trap":"needle"`, "does not support"},
		"stuck":           {`,"hidden":false,"state":"closed","locked":false,"distance_feet":5,"stuck":true`, "does not support"},
		"state missing":   {`,"hidden":false,"locked":false,"distance_feet":5`, "open or closed"},
		"state unknown":   {`,"hidden":false,"state":"ajar","locked":false,"distance_feet":5`, "open or closed"},
		"lock unknown":    {`,"hidden":false,"state":"closed","distance_feet":5`, "Clarify whether the door is locked"},
		"locked":          {`,"hidden":false,"state":"closed","locked":true,"distance_feet":5`, "Locked doors are not supported"},
		"distance absent": {`,"hidden":false,"state":"closed","locked":false`, "exact distance"},
		"distance < 0":    {`,"hidden":false,"state":"closed","locked":false,"distance_feet":-5`, "exact distance"},
	} {
		if _, err := ParseDoorFacts(door(tc.fields)); err == nil || !strings.Contains(err.Error(), tc.reason) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := ParseDoorFacts([]byte(`{"profile":"secret-door-v1","actor_id":"` + doorActor + `","hidden":false,"state":"closed","locked":false,"distance_feet":5}`)); err == nil {
		t.Error("unknown profile accepted")
	}
}

func TestOpenDoorReachAndState(t *testing.T) {
	for distance, ok := range map[string]bool{"0": true, "5": true, "10": false} {
		facts, err := ParseDoorFacts(door(`,"hidden":false,"state":"closed","locked":false,"distance_feet":` + distance))
		if err != nil {
			t.Fatal(err)
		}
		after, receipt, err := OpenDoorFacts(facts, "door-1", "Oak door")
		if !ok {
			if err == nil || !strings.Contains(err.Error(), "Out of reach") || after.State != DoorClosed {
				t.Fatalf("%s feet opened: %v", distance, err)
			}
			continue
		}
		if err != nil || after.State != DoorOpen || facts.State != DoorClosed || receipt.Before != DoorClosed || receipt.After != DoorOpen || receipt.ReachFeet != 5 || len(receipt.SourceIDs) != 2 {
			t.Fatalf("%s feet: %+v %+v %v", distance, after, receipt, err)
		}
	}
	open, err := ParseDoorFacts(door(`,"hidden":false,"state":"open","locked":false,"distance_feet":5`))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := OpenDoorFacts(open, "door-1", "Oak door"); err == nil || !strings.Contains(err.Error(), "already open") {
		t.Fatalf("open door reopened: %v", err)
	}
	if _, _, err := OpenDoorFacts(DoorFacts{State: DoorClosed}, "door-1", "Oak door"); err == nil {
		t.Fatal("unvalidated facts opened")
	}
}
