package statedb

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

// IsSceneObject reports whether an npcs row is a scene object rather than a
// character. Doors (#146) are stored as npcs rows with server-owned
// properties.srd521_door facts. A properties.srd521_magic fact alone does not
// make an object, because a creature can bear magic; its bearer field must say
// "object". The NPC list queries apply the same two tests.
func IsSceneObject(properties []byte) bool {
	var props map[string]json.RawMessage
	if json.Unmarshal(properties, &props) != nil {
		return false
	}
	if _, door := props["srd521_door"]; door {
		return true
	}
	// A map keeps the key match exact, as in SQL.
	var magic map[string]json.RawMessage
	if json.Unmarshal(props["srd521_magic"], &magic) != nil {
		return false
	}
	var bearer string
	return json.Unmarshal(magic["bearer"], &bearer) == nil && bearer == "object"
}

// GetCharacterNPC is GetNPCByID for callers that treat the row as a character.
// A scene object returns pgx.ErrNoRows, so it looks the same as a missing NPC.
// Rules code that works on objects uses its own queries.
func GetCharacterNPC(ctx context.Context, q interface {
	GetNPCByID(context.Context, GetNPCByIDParams) (Npc, error)
}, arg GetNPCByIDParams) (Npc, error) {
	npc, err := q.GetNPCByID(ctx, arg)
	if err != nil {
		return Npc{}, err
	}
	if IsSceneObject(npc.Properties) {
		return Npc{}, pgx.ErrNoRows
	}
	return npc, nil
}
