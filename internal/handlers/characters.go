package handlers

import (
	"fmt"
	"net/http"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// GetCharacter returns the player character for a campaign.
func (h *CharacterHandlers) GetCharacter(w http.ResponseWriter, r *http.Request) {
	campaignID, err := campaignIDFromURL(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid campaign id: %v", err))
		return
	}

	pcs, err := h.Queries.GetPlayerCharacterByCampaign(r.Context(), dbutil.ToPgtype(campaignID))
	if err != nil {
		h.Logger.Error("get character for campaign", "campaign_id", campaignID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get character")
		return
	}
	if len(pcs) == 0 {
		writeError(w, http.StatusNotFound, "no character found for campaign")
		return
	}

	campaign, err := h.Queries.GetCampaignByID(r.Context(), dbutil.ToPgtype(campaignID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get character rules version")
		return
	}
	response := playerCharacterToResponse(pcs[len(pcs)-1])
	response.RulesVersion = campaign.RulesVersion
	if campaign.RulesVersion == string(domain.RulesVersionSRD521) {
		pc := pcs[len(pcs)-1]
		response.Sheet, err = srd521.Project(pc.Stats, int(pc.Level), int(pc.Hp), int(pc.MaxHp), int(pc.Experience))
		if err == nil {
			nameConcentrations(response.Sheet.Concentration)
			response.Sheet.Fall = game.PendingFall(pc.Stats)
			response.Sheet.DeathSaves = game.ZeroHPState(pc.Stats, int(pc.Hp))
			response.Sheet.Prone = game.ReadConditions(pc.Stats).Prone
		}
		if err == nil && h.Pool != nil {
			now, clockErr := game.ReadSpellClock(r.Context(), h.Pool, campaignID)
			if clockErr != nil {
				err = clockErr
			} else {
				err = game.ApplyFaithAC(r.Context(), h.Pool, campaignID, dbutil.FromPgtype(pc.ID), now, response.Sheet)
				if err == nil {
					err = game.ApplyLight(r.Context(), h.Pool, campaignID, dbutil.FromPgtype(pc.ID), now, response.Sheet)
				}
				if err == nil {
					actor := domain.PlayerCharacter{ID: dbutil.FromPgtype(pc.ID)}
					if pc.CurrentLocationID.Valid {
						location := dbutil.FromPgtype(pc.CurrentLocationID)
						actor.CurrentLocationID = &location
					}
					err = game.ApplyExploration(r.Context(), h.Pool, domain.Campaign{ID: campaignID}, actor, response.Sheet)
				}
			}
		}
		if err != nil {
			response.Sheet = nil
			response.SheetUnavailable = err.Error()
		}
	}
	writeJSON(w, http.StatusOK, response)
}

// GetCharacterInventory returns the inventory items for the player character.
func (h *CharacterHandlers) GetCharacterInventory(w http.ResponseWriter, r *http.Request) {
	campaignID, err := campaignIDFromURL(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid campaign id: %v", err))
		return
	}

	pcs, err := h.Queries.GetPlayerCharacterByCampaign(r.Context(), dbutil.ToPgtype(campaignID))
	if err != nil || len(pcs) == 0 {
		writeError(w, http.StatusNotFound, "no character found for campaign")
		return
	}

	items, err := h.Queries.ListItemsByPlayer(r.Context(), statedb.ListItemsByPlayerParams{
		CampaignID:        dbutil.ToPgtype(campaignID),
		PlayerCharacterID: pcs[len(pcs)-1].ID,
	})
	if err != nil {
		h.Logger.Error("list inventory for campaign", "campaign_id", campaignID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list inventory")
		return
	}

	resp := make([]api.ItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, itemToResponse(item))
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetCharacterAbilities returns the abilities for the player character.
func (h *CharacterHandlers) GetCharacterAbilities(w http.ResponseWriter, r *http.Request) {
	campaignID, err := campaignIDFromURL(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid campaign id: %v", err))
		return
	}

	pcs, err := h.Queries.GetPlayerCharacterByCampaign(r.Context(), dbutil.ToPgtype(campaignID))
	if err != nil || len(pcs) == 0 {
		writeError(w, http.StatusNotFound, "no character found for campaign")
		return
	}

	writeJSON(w, http.StatusOK, parseAbilities(pcs[len(pcs)-1].Abilities))
}
