package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/world"
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
	response.NeedsName = domain.ValidateCharacterName(response.Name) != nil
	if campaign.RulesVersion == string(domain.RulesVersionSRD521) {
		pc := pcs[len(pcs)-1]
		sheet, err := srd521.Project(pc.Stats, int(pc.Level), int(pc.Hp), int(pc.MaxHp), int(pc.Experience))
		switch {
		case err != nil && !errors.Is(err, srd521.ErrCoverageUnavailable):
			// What is wrong with the saved build. It holds no internal detail.
			response.SheetUnavailable = err.Error()
		case err == nil:
			if err = h.applyCurrentState(r.Context(), campaignID, pc, sheet); err == nil {
				response.Sheet = sheet
			}
		}
		if err != nil && response.SheetUnavailable == "" {
			h.Logger.Error("load character sheet", "campaign_id", campaignID, "error", logging.Redact(err.Error()))
			response.SheetUnavailable = sheetLoadFailed
		}
	}
	writeJSON(w, http.StatusOK, response)
}

// RenameCharacter saves a new name for one of the campaign's characters. The
// caller must have checked that the user owns the campaign.
func (h *CharacterHandlers) RenameCharacter(w http.ResponseWriter, r *http.Request) {
	campaignID, err := campaignIDFromURL(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid campaign id: %v", err))
		return
	}
	var req api.CharacterRenameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	characterID, err := uuid.Parse(req.CharacterID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	err = world.RenameCharacter(r.Context(), h.Queries, dbutil.ToPgtype(campaignID), dbutil.ToPgtype(characterID), req.Name)
	switch {
	case errors.Is(err, domain.ErrCharacterNameRequired):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, world.ErrCharacterNotFound):
		writeError(w, http.StatusNotFound, "character not found")
	case err != nil:
		h.Logger.Error("rename character", "campaign_id", campaignID, "character_id", characterID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to rename character")
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

// sheetLoadFailed is the sheet_unavailable reason for a failure inside the
// server. The cause is logged and not sent, because database and game errors
// can name tables, hosts and credentials.
const sheetLoadFailed = "The character sheet could not be loaded. Try again."

// applyCurrentState adds what the saved build alone does not hold: conditions,
// and with a database, the effects and choices at the campaign's current time.
func (h *CharacterHandlers) applyCurrentState(ctx context.Context, campaignID uuid.UUID, pc statedb.PlayerCharacter, sheet *srd521.CharacterSheet) error {
	nameConcentrations(sheet.Concentration)
	sheet.Fall = game.PendingFall(pc.Stats)
	sheet.DeathSaves = game.ZeroHPState(pc.Stats, int(pc.Hp))
	sheet.Prone = game.ReadConditions(pc.Stats).Prone
	if h.Pool == nil {
		return nil
	}
	now, err := game.ReadSpellClock(ctx, h.Pool, campaignID)
	if err != nil {
		return fmt.Errorf("read spell clock: %w", err)
	}
	actor := domain.PlayerCharacter{ID: dbutil.FromPgtype(pc.ID)}
	if pc.CurrentLocationID.Valid {
		location := dbutil.FromPgtype(pc.CurrentLocationID)
		actor.CurrentLocationID = &location
	}
	if err := game.ApplyFaithAC(ctx, h.Pool, campaignID, actor.ID, now, sheet); err != nil {
		return fmt.Errorf("apply Shield of Faith: %w", err)
	}
	if err := game.ApplyLight(ctx, h.Pool, campaignID, actor.ID, now, sheet); err != nil {
		return fmt.Errorf("apply Light: %w", err)
	}
	if err := game.ApplyExploration(ctx, h.Pool, domain.Campaign{ID: campaignID}, actor, sheet); err != nil {
		return fmt.Errorf("read exploration choices: %w", err)
	}
	return nil
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
