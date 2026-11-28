package statedb

import "encoding/json"

// IsSceneObject reports whether an npcs row is a scene object rather than a
// character. Doors (#146) are stored as npcs rows with server-owned
// properties.srd521_door facts; the NPC list queries exclude them with the
// same test. srd521_magic is not a discriminator: a creature can carry magic.
func IsSceneObject(properties []byte) bool {
	var props map[string]json.RawMessage
	if json.Unmarshal(properties, &props) != nil {
		return false
	}
	_, door := props["srd521_door"]
	return door
}
