package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

type sheetQuerier struct {
	statedb.Querier
	stats []byte
}

func (q sheetQuerier) GetPlayerCharacterByCampaign(context.Context, pgtype.UUID) ([]statedb.PlayerCharacter, error) {
	return []statedb.PlayerCharacter{{Name: "Fixture", Stats: q.stats, Hp: 9, MaxHp: 12, Level: 1}}, nil
}

func (q sheetQuerier) GetCampaignByID(context.Context, pgtype.UUID) (statedb.Campaign, error) {
	return statedb.Campaign{RulesVersion: "srd-5.2.1"}, nil
}

// sheetDB fails the spell clock read with clockErr, and every later query with
// queryErr.
type sheetDB struct {
	db.DBTX
	clockErr, queryErr error
}

func (d sheetDB) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return errRow{d.clockErr}
}

func (d sheetDB) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, d.queryErr
}

type errRow struct{ err error }

func (r errRow) Scan(...any) error { return r.err }

// getSheet calls the character route and returns the reply and the log text.
// The logger here does not redact, so the log shows what the handler passed.
func getSheet(t *testing.T, stats []byte, pool db.DBTX) (api.CharacterResponse, string, string) {
	t.Helper()
	var logs bytes.Buffer
	h := &CharacterHandlers{Queries: sheetQuerier{stats: stats}, Logger: log.New(&logs), Pool: pool}
	r := chi.NewRouter()
	r.Get("/campaigns/{id}/character", h.GetCharacter)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaigns/"+uuid.NewString()+"/character", nil))
	var response api.CharacterResponse
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &response) != nil {
		t.Fatalf("response %d %s", rec.Code, rec.Body.String())
	}
	return response, rec.Body.String(), logs.String()
}

func TestGetCharacterWithholdsSheetDatabaseErrors(t *testing.T) {
	stats, err := os.ReadFile("../rules/srd521/testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	failure := errors.New(`ERROR: relation "concentration_casts" does not exist at postgres://edda:fake-pass@db.internal/edda`)
	for name, pool := range map[string]db.DBTX{
		"spell clock": sheetDB{clockErr: failure},
		"later query": sheetDB{queryErr: failure},
	} {
		t.Run(name, func(t *testing.T) {
			response, body, logs := getSheet(t, stats, pool)
			if response.Sheet != nil {
				t.Error("a partial sheet was returned")
			}
			if want := "The character sheet could not be loaded. Try again."; response.SheetUnavailable != want {
				t.Errorf("sheet_unavailable = %q, want %q", response.SheetUnavailable, want)
			}
			for _, private := range []string{"concentration_casts", "db.internal", "fake-pass"} {
				if strings.Contains(body, private) {
					t.Errorf("reply holds %q: %s", private, body)
				}
			}
			if !strings.Contains(logs, "concentration_casts") || !strings.Contains(logs, "db.internal") {
				t.Errorf("log does not hold the cause: %s", logs)
			}
			if strings.Contains(logs, "fake-pass") {
				t.Errorf("log holds a credential: %s", logs)
			}
		})
	}
}

func TestGetCharacterKeepsUnsupportedBuildReason(t *testing.T) {
	for name, tc := range map[string]struct{ stats, want string }{
		"unsupported class": {`{"srd521":{"class":"wizard"}}`, "supported projection: Fighter or Cleric with supported ordinary equipment and training"},
		"unreadable stats":  {`[]`, "invalid SRD build data"},
	} {
		t.Run(name, func(t *testing.T) {
			response, _, logs := getSheet(t, []byte(tc.stats), nil)
			if response.SheetUnavailable != tc.want {
				t.Errorf("sheet_unavailable = %q, want %q", response.SheetUnavailable, tc.want)
			}
			if logs != "" {
				t.Errorf("an unsupported build was logged: %s", logs)
			}
		})
	}
}
