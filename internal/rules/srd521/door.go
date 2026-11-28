package srd521

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// OpenDoor opens one known, ordinary, unlocked, closed door (#146).
	OpenDoor = "door-open"
	// DoorFactsProfile names the trusted narrative-space door contract.
	DoorFactsProfile = "ordinary-door-v1"
	// DoorReachFeet is the ordinary 5-foot reach (SRD 5.2.1 Reach, p. 186).
	DoorReachFeet = 5
	// DoorOpen and DoorClosed are the only declared door states.
	DoorOpen   = "open"
	DoorClosed = "closed"
)

// DoorSourceIDs cite Interacting with Objects (p. 11) and Reach.
var DoorSourceIDs = []string{"srd-5.2.1:content/exploration/interacting-with-objects", "srd-5.2.1:content/rules-glossary/reach"}

// DoorFacts is one explicitly adjudicated door known to one actor. Every field
// is required. Unknown fields (traps, stuck or barred doors, keys, hit points)
// are rejected rather than ignored.
type DoorFacts struct {
	Profile      string `json:"profile"`
	ActorID      string `json:"actor_id"`
	Hidden       *bool  `json:"hidden"`
	State        string `json:"state"`
	Locked       *bool  `json:"locked"`
	DistanceFeet *int   `json:"distance_feet"`
}

// DoorFor reports whether raw facts describe a door this actor knows about.
// It reads only identity and hidden fields, so a hidden or foreign door is
// indistinguishable from no door at all.
func DoorFor(raw []byte, actorID string) bool {
	var identity struct {
		ActorID string `json:"actor_id"`
		Hidden  *bool  `json:"hidden"`
	}
	return json.Unmarshal(raw, &identity) == nil && identity.ActorID == actorID && identity.Hidden != nil && !*identity.Hidden
}

// ParseDoorFacts accepts only the exact ordinary-door vocabulary. Its errors
// are player-safe for a door the actor already knows.
func ParseDoorFacts(raw []byte) (DoorFacts, error) {
	var facts DoorFacts
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&facts) != nil || facts.Profile != DoorFactsProfile {
		return facts, errors.New("This door declares facts Edda does not support; only an ordinary door with an explicit state, lock ruling and distance is supported.")
	}
	if facts.Hidden == nil || *facts.Hidden {
		return facts, errors.New("This door is not known.")
	}
	if facts.State != DoorOpen && facts.State != DoorClosed {
		return facts, errors.New("Clarify whether the door is open or closed.")
	}
	if facts.Locked == nil {
		return facts, errors.New("Clarify whether the door is locked.")
	}
	if *facts.Locked {
		return facts, errors.New("Locked doors are not supported.")
	}
	if facts.DistanceFeet == nil || *facts.DistanceFeet < 0 {
		return facts, errors.New("Clarify the exact distance to the door.")
	}
	return facts, nil
}

// DoorReceipt records one confirmed door state change.
type DoorReceipt struct {
	RulesVersion string   `json:"rules_version"`
	Kind         string   `json:"kind"`
	DoorID       string   `json:"door_id"`
	Door         string   `json:"door"`
	Before       string   `json:"before"`
	After        string   `json:"after"`
	DistanceFeet int      `json:"distance_feet"`
	ReachFeet    int      `json:"reach_feet"`
	Reason       string   `json:"reason"`
	SourceIDs    []string `json:"source_ids"`
}

// OpenDoorFacts opens a closed, unlocked door within reach outside combat.
// It returns the changed facts and the receipt; no time or action is spent.
func OpenDoorFacts(facts DoorFacts, id, name string) (DoorFacts, *DoorReceipt, error) {
	if facts.Locked == nil || *facts.Locked || facts.DistanceFeet == nil {
		return facts, nil, errors.New("Only a known ordinary unlocked door can be opened.")
	}
	if facts.State != DoorClosed {
		return facts, nil, errors.New("The door is already open.")
	}
	if *facts.DistanceFeet > DoorReachFeet {
		return facts, nil, fmt.Errorf("Out of reach: the door is %d feet away; opening it requires %d feet or less.", *facts.DistanceFeet, DoorReachFeet)
	}
	after := facts
	after.State = DoorOpen
	return after, &DoorReceipt{
		RulesVersion: "srd-5.2.1", Kind: OpenDoor, DoorID: id, Door: name,
		Before: DoorClosed, After: DoorOpen, DistanceFeet: *facts.DistanceFeet, ReachFeet: DoorReachFeet,
		Reason:    fmt.Sprintf("%s is now open. Outside combat, opening an ordinary unlocked door within reach takes no action and no campaign time.", name),
		SourceIDs: append([]string(nil), DoorSourceIDs...),
	}, nil
}
