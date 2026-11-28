package game

import (
	"context"
	"testing"

	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

var statResolverCampaignID = uuid.New()

var _ domain.StatModifierResolver = (*statModifierResolver)(nil)

func TestStatModifierResolverReturnsModifier(t *testing.T) {
	q := newMockQuerier()
	pcID := uuid.New()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(statResolverCampaignID),
		Stats:      []byte(`{"strength":16,"dexterity":12}`),
	}
	resolver := NewStatModifierResolver(q)

	// strength 16 → modifier (16-10)/2 = 3
	mod, err := resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "strength")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != 3 {
		t.Fatalf("modifier = %d, want 3 for strength 16", mod)
	}

	// dexterity 12 → modifier (12-10)/2 = 1
	mod, err = resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "dexterity")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != 1 {
		t.Fatalf("modifier = %d, want 1 for dexterity 12", mod)
	}
}

func TestStatModifierResolverCaseInsensitive(t *testing.T) {
	q := newMockQuerier()
	pcID := uuid.New()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(statResolverCampaignID),
		Stats:      []byte(`{"Strength":14}`),
	}
	resolver := NewStatModifierResolver(q)

	mod, err := resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "strength")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != 2 {
		t.Fatalf("modifier = %d, want 2 for Strength 14", mod)
	}
}

func TestStatModifierResolverMissingStatReturnsZero(t *testing.T) {
	q := newMockQuerier()
	pcID := uuid.New()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(statResolverCampaignID),
		Stats:      []byte(`{"strength":10}`),
	}
	resolver := NewStatModifierResolver(q)

	mod, err := resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "wisdom")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != 0 {
		t.Fatalf("modifier = %d, want 0 for missing stat", mod)
	}
}

func TestStatModifierResolverEmptyStatsReturnsZero(t *testing.T) {
	q := newMockQuerier()
	pcID := uuid.New()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(statResolverCampaignID),
	}
	resolver := NewStatModifierResolver(q)

	mod, err := resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "strength")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != 0 {
		t.Fatalf("modifier = %d, want 0 for empty stats", mod)
	}
}

func TestStatModifierResolverLowScoreNegativeModifier(t *testing.T) {
	q := newMockQuerier()
	pcID := uuid.New()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(statResolverCampaignID),
		Stats:      []byte(`{"strength":8}`),
	}
	resolver := NewStatModifierResolver(q)

	// strength 8 → modifier (8-10)/2 = -1
	mod, err := resolver.GetStatModifier(tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID), pcID, "strength")
	if err != nil {
		t.Fatalf("GetStatModifier() error = %v", err)
	}
	if mod != -1 {
		t.Fatalf("modifier = %d, want -1 for strength 8", mod)
	}
}

func TestStatModifierResolverForeignCharacterMatchesUnknown(t *testing.T) {
	pcID := uuid.New()
	q := newMockQuerier()
	q.playerCharacter = statedb.PlayerCharacter{
		ID:         dbutil.ToPgtype(pcID),
		CampaignID: dbutil.ToPgtype(uuid.New()),
		Stats:      []byte(`{"strength":18}`),
	}
	resolver := NewStatModifierResolver(q)
	ctx := tools.WithCurrentCampaignID(context.Background(), statResolverCampaignID)

	mod, err := resolver.GetStatModifier(ctx, pcID, "strength")
	if !errors.Is(err, pgx.ErrNoRows) || mod != 0 {
		t.Fatalf("GetStatModifier() = %d, %v; want 0 and pgx.ErrNoRows like an unknown character", mod, err)
	}
	if _, err := resolver.GetStatModifier(context.Background(), pcID, "strength"); err == nil {
		t.Fatal("GetStatModifier() without a current campaign should fail")
	}
}
