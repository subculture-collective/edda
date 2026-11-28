//go:build integration

package statedb_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/world"
)

const blankNameAudit = "../../../docs/data-checks/blank-character-names.sql"

// The audit must run read-only on the current schema, find what the broken
// guided form saved, and stop listing a character once it is renamed.
func TestIntegrationBlankCharacterNameAudit(t *testing.T) {
	ctx := context.Background()
	audit, err := os.ReadFile(blankNameAudit)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	q := statedb.New(tx)

	user := createUser(t, q, "Player")
	campaign := func(genre, tone string) statedb.Campaign {
		c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{
			Name: genre + " " + tone, Genre: txt(genre), Tone: txt(tone), Status: "active", CreatedBy: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	// save persists a character as creation did before the name was checked.
	save := func(c statedb.Campaign, attrs world.CharacterAttributes) statedb.PlayerCharacter {
		pc, err := world.PersistCharacterProfile(ctx, q, dbutil.FromPgtype(c.ID), dbutil.FromPgtype(user.ID),
			world.BuildCharacterFromAttributes(attrs), nil)
		if err != nil {
			t.Fatal(err)
		}
		return pc
	}
	defaults := world.CharacterAttributes{
		Race: world.D5ERaces[0], Class: world.D5EClasses[0],
		Background: world.D5EBackgrounds[0], Alignment: world.D5EAlignments[0],
	}
	named := defaults
	named.Name = "Mara"
	other := defaults
	other.Name, other.Class = " \t\n", "Wizard"

	broken, kept := campaign("Dark Fantasy", "Gritty"), campaign("Science fiction", "Hopeful")
	guided := save(broken, defaults)
	blank := save(kept, other)
	save(kept, named)

	// run returns each statement's rows as column name to value.
	run := func() [][]map[string]string {
		results, err := tx.Conn().PgConn().Exec(ctx, string(audit)).ReadAll()
		if err != nil {
			t.Fatalf("run audit: %v", err)
		}
		return rowsByColumn(results)
	}
	// flags maps each row's id to its true-or-false column.
	flags := func(rows []map[string]string, id, flag string) map[string]string {
		out := map[string]string{}
		for _, row := range rows {
			out[row[id]] = row[flag]
		}
		return out
	}
	uuidText := func(id pgtype.UUID) string { return dbutil.FromPgtype(id).String() }

	if _, err := tx.Exec(ctx, "SAVEPOINT writable"); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, "SET TRANSACTION READ ONLY"); err != nil {
		t.Fatal(err)
	}
	results := run()
	if len(results) != 2 {
		t.Fatalf("audit has %d statements, want 2", len(results))
	}
	// Other tests roll back or delete their rows, so only these are listed.
	characters := flags(results[0], "character_id", "guided_form_defaults")
	if len(characters) != 2 || characters[uuidText(guided.ID)] != "t" || characters[uuidText(blank.ID)] != "f" {
		t.Errorf("blank-named characters = %v", characters)
	}
	campaigns := flags(results[1], "campaign_id", "has_guided_form_character")
	if len(campaigns) != 1 || campaigns[uuidText(broken.ID)] != "t" {
		t.Errorf("candidate campaigns = %v", campaigns)
	}
	if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT writable"); err != nil {
		t.Fatal(err)
	}

	// The TUI's repair finds the character, refuses another campaign's ID and
	// saves the trimmed name.
	unnamed, err := world.UnnamedCharacter(ctx, q, kept.ID)
	if err != nil || unnamed == nil || unnamed.ID != blank.ID {
		t.Fatalf("unnamed character = %+v, error = %v", unnamed, err)
	}
	if err := world.RenameCharacter(ctx, q, broken.ID, blank.ID, "Iona"); err == nil {
		t.Error("renamed a character through another campaign")
	}
	if err := world.RenameCharacter(ctx, q, kept.ID, blank.ID, " Iona "); err != nil {
		t.Fatal(err)
	}
	renamed, err := q.GetPlayerCharacterByID(ctx, blank.ID)
	if err != nil || renamed.Name != "Iona" {
		t.Fatalf("renamed character = %q, error = %v", renamed.Name, err)
	}
	if unnamed, err = world.UnnamedCharacter(ctx, q, kept.ID); err != nil || unnamed != nil {
		t.Errorf("unnamed character after rename = %+v, error = %v", unnamed, err)
	}
	if characters := flags(run()[0], "character_id", "guided_form_defaults"); len(characters) != 1 {
		t.Errorf("blank-named characters after rename = %v", characters)
	}
}

func rowsByColumn(results []*pgconn.Result) [][]map[string]string {
	out := make([][]map[string]string, len(results))
	for i, result := range results {
		for _, values := range result.Rows {
			row := map[string]string{}
			for column, field := range result.FieldDescriptions {
				row[field.Name] = string(values[column])
			}
			out[i] = append(out[i], row)
		}
	}
	return out
}
