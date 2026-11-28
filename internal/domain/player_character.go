package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PlayerCharacter struct {
	ID                uuid.UUID
	CampaignID        uuid.UUID
	UserID            uuid.UUID
	Name              string
	Description       string
	Stats             json.RawMessage
	HP                int
	MaxHP             int
	Experience        int
	Level             int
	Status            string
	Abilities         json.RawMessage
	CurrentLocationID *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ErrCharacterNameRequired refuses a player character with a blank name.
var ErrCharacterNameRequired = errors.New("character name is required")

// ValidateCharacterName is the rule for a player character's name: it must hold
// something other than white space.
func ValidateCharacterName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrCharacterNameRequired
	}
	return nil
}

func (pc *PlayerCharacter) Validate() error {
	if err := ValidateCharacterName(pc.Name); err != nil {
		return err
	}
	if pc.CampaignID == uuid.Nil {
		return errors.New("player character campaign_id is required")
	}
	if pc.UserID == uuid.Nil {
		return errors.New("player character user_id is required")
	}
	if pc.Level < 1 {
		return errors.New("player character level must be at least 1")
	}
	return nil
}
