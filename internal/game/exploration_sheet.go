package game

import (
	"context"
	"slices"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// KnownDoorChoices names the doors the actor knows at their location, in the
// order a client lists them. It copies no fact but an accepted door's state.
func KnownDoorChoices(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter) ([]srd521.SheetDoor, error) {
	doors, err := readKnownDoors(ctx, database, campaign, pc, false)
	if err != nil {
		return nil, err
	}
	choices := make([]srd521.SheetDoor, len(doors))
	for i, door := range doors {
		choices[i] = srd521.SheetDoor{ID: door.ID.String(), Name: door.Name}
		if door.Err == nil {
			choices[i].State = door.Facts.State
		}
	}
	return choices, nil
}

// ApplyExploration adds the choices a client needs to offer exploration
// actions: the known doors, and the travel paces when a known connected route
// has declared facts. Distance, terrain and destinations come from the preview.
func ApplyExploration(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, sheet *srd521.CharacterSheet) error {
	doors, err := KnownDoorChoices(ctx, database, campaign, pc)
	if err != nil {
		return err
	}
	sheet.Doors = doors
	if pc.CurrentLocationID == nil {
		return nil
	}
	routes, err := readTravelRoutes(ctx, database, campaign.ID, *pc.CurrentLocationID, false)
	if err != nil {
		return err
	}
	if slices.ContainsFunc(routes, func(route travelRow) bool { return len(route.Facts) > 0 }) {
		sheet.TravelPaces = slices.Clone(srd521.TravelPaces)
	}
	return nil
}
