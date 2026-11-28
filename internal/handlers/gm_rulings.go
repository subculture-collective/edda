package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// rulingEngine declares server-owned facts for the campaign owner acting as
// game master. These routes are mounted behind RequireOwnership.
type rulingEngine interface {
	DeclareRoute(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, []byte) (*engine.RouteRuling, error)
	DeclareFall(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, game.DeclaredFall) (*game.DeclaredFall, error)
}

// DeclareRoute handles PUT /campaigns/{id}/locations/{lid}/routes/{to}: the
// body is the route fact object stored in location_connections.srd521_route.
func (h *ActionHandlers) DeclareRoute(w http.ResponseWriter, r *http.Request) {
	campaignID, user, rulings, ok := h.rulingRequest(w, r)
	if !ok {
		return
	}
	fromID, errFrom := uuid.Parse(chi.URLParam(r, "lid"))
	toID, errTo := uuid.Parse(chi.URLParam(r, "to"))
	if errFrom != nil || errTo != nil {
		writeError(w, http.StatusBadRequest, "invalid location id")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid route facts")
		return
	}
	ruling, err := rulings.DeclareRoute(r.Context(), campaignID, user, fromID, toID, body)
	if err != nil {
		h.writeRulingError(w, err, "connection not found")
		return
	}
	writeJSON(w, http.StatusOK, ruling)
}

// DeclareFall handles POST /campaigns/{id}/character/fall: actor_id plus the
// declared fall facts. The server assigns the declaration's id.
func (h *ActionHandlers) DeclareFall(w http.ResponseWriter, r *http.Request) {
	campaignID, user, rulings, ok := h.rulingRequest(w, r)
	if !ok {
		return
	}
	var request struct {
		ActorID uuid.UUID `json:"actor_id"`
		game.DeclaredFall
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err != nil || decoder.Decode(&request) != nil || decoder.Decode(new(any)) != io.EOF {
		writeError(w, http.StatusUnprocessableEntity, "invalid fall declaration")
		return
	}
	if request.ActorID == uuid.Nil {
		writeError(w, http.StatusUnprocessableEntity, "actor_id is required")
		return
	}
	fall, err := rulings.DeclareFall(r.Context(), campaignID, user, request.ActorID, request.DeclaredFall)
	if err != nil {
		h.writeRulingError(w, err, "character not found")
		return
	}
	writeJSON(w, http.StatusOK, fall)
}

func (h *ActionHandlers) rulingRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, rulingEngine, bool) {
	campaignID, err := campaignIDFromURL(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign id")
		return uuid.Nil, uuid.Nil, nil, false
	}
	user, ok := auth.UserFromContext(r.Context())
	if !ok || user == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return uuid.Nil, uuid.Nil, nil, false
	}
	rulings, ok := h.Engine.(rulingEngine)
	if !ok {
		writeError(w, http.StatusServiceUnavailable, "GM rulings unavailable")
		return uuid.Nil, uuid.Nil, nil, false
	}
	return campaignID, user, rulings, true
}

func (h *ActionHandlers) writeRulingError(w http.ResponseWriter, err error, missing string) {
	var refusal engine.RulingRefusal
	switch {
	case errors.Is(err, engine.ErrRulingNotFound):
		writeError(w, http.StatusNotFound, missing)
	case errors.As(err, &refusal) && refusal.Conflict:
		writeError(w, http.StatusConflict, refusal.Reason)
	case errors.As(err, &refusal):
		writeError(w, http.StatusUnprocessableEntity, refusal.Reason)
	case errors.Is(err, engine.ErrTurnInProgress), errors.Is(err, engine.ErrTurnRecoveryRequired):
		writeError(w, http.StatusConflict, err.Error())
	default:
		h.Logger.Errorf("GM ruling: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to record the ruling")
	}
}
