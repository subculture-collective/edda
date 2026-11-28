package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeGrantFeatRow struct {
	scan func(dest ...any) error
}

func (r fakeGrantFeatRow) Scan(dest ...any) error { return r.scan(dest...) }

// fakeGrantFeatDB holds one feat and the player characters of each campaign.
type fakeGrantFeatDB struct {
	featID     uuid.UUID
	campaignID uuid.UUID
	characters map[uuid.UUID]uuid.UUID // character -> campaign
	inserts    [][]any
}

func (d *fakeGrantFeatDB) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	d.inserts = append(d.inserts, args)
	return pgconn.CommandTag{}, nil
}

func (d *fakeGrantFeatDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, pgx.ErrNoRows
}

func (d *fakeGrantFeatDB) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	switch {
	case strings.Contains(sql, "FROM player_characters"):
		id, _ := args[0].(uuid.UUID)
		campaignID, _ := args[1].(uuid.UUID)
		return fakeGrantFeatRow{scan: func(dest ...any) error {
			if owner, ok := d.characters[id]; !ok || owner != campaignID {
				return pgx.ErrNoRows
			}
			*(dest[0].(*uuid.UUID)) = id
			return nil
		}}
	case strings.Contains(sql, "FROM feat_definitions"):
		return fakeGrantFeatRow{scan: func(dest ...any) error {
			*(dest[0].(*uuid.UUID)) = d.featID
			*(dest[1].(*string)) = "Sturdy."
			return nil
		}}
	}
	return fakeGrantFeatRow{scan: func(...any) error { return pgx.ErrNoRows }}
}

func TestGrantFeatCharacterMustBeInCurrentCampaign(t *testing.T) {
	campaignID := uuid.New()
	ownID := uuid.New()
	otherPCID := uuid.New()
	foreignID := uuid.New()
	db := &fakeGrantFeatDB{
		featID:     uuid.New(),
		campaignID: campaignID,
		characters: map[uuid.UUID]uuid.UUID{ownID: campaignID, otherPCID: campaignID, foreignID: uuid.New()},
	}
	h := &GrantFeatHandler{db: db}
	ctx := WithCurrentCampaignID(WithCurrentPlayerCharacterID(context.Background(), ownID), campaignID)

	unknownID := uuid.New()
	_, foreignErr := h.Handle(ctx, map[string]any{"feat_name": "Tough", "character_id": foreignID.String()})
	_, unknownErr := h.Handle(ctx, map[string]any{"feat_name": "Tough", "character_id": unknownID.String()})
	if foreignErr == nil || unknownErr == nil {
		t.Fatalf("foreign err = %v, unknown err = %v, want both errors", foreignErr, unknownErr)
	}
	if strings.ReplaceAll(foreignErr.Error(), foreignID.String(), "X") != strings.ReplaceAll(unknownErr.Error(), unknownID.String(), "X") {
		t.Fatalf("foreign error %q differs from unknown error %q beyond the ID", foreignErr, unknownErr)
	}
	if len(db.inserts) != 0 {
		t.Fatalf("feat inserted for unknown or foreign character: %v", db.inserts)
	}

	if _, err := h.Handle(ctx, map[string]any{"feat_name": "Tough", "character_id": otherPCID.String()}); err != nil {
		t.Fatalf("same-campaign character: %v", err)
	}
	if _, err := h.Handle(ctx, map[string]any{"feat_name": "Tough"}); err != nil {
		t.Fatalf("current character: %v", err)
	}
	if len(db.inserts) != 2 {
		t.Fatalf("inserts = %d, want 2", len(db.inserts))
	}
}
