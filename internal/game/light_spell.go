package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// LightEntryPrefix selects a Light cast: EntryID is this prefix followed by
// one color from srd521.LightColors. The target is one object the caster carries.
const LightEntryPrefix = "light:"

// lightObjectTypes are the item types that can hold the spell. Consumables
// and quest items are excluded: their use or removal is not modelled.
const lightObjectTypes = `('weapon','armor','misc')`

// LightObject is one object the caster carries, revalidated in this snapshot.
type LightObject struct {
	ID        uuid.UUID
	Name      string
	UpdatedAt time.Time
}

// LightObject returns only an object listed in this preparation snapshot.
func (s *StoredAction) LightObject(id uuid.UUID) *LightObject {
	for i := range s.LightObjects {
		if s.LightObjects[i].ID == id {
			return &s.LightObjects[i]
		}
	}
	return nil
}

// IsLightIntent reports whether the intent casts or dismisses Light.
func IsLightIntent(intent ActionIntent) bool {
	return intent.Kind == srd521.DismissLight || isLightCover(intent.Kind) || (intent.Kind == srd521.SpellCast && (intent.EntryID == srd521.LightSpell || strings.HasPrefix(intent.EntryID, LightEntryPrefix)))
}

func isLightCover(kind string) bool {
	return kind == srd521.CoverLight || kind == srd521.UncoverLight
}

// ApplyLight drops an ended Light from the sheet and names the object the
// caster still carries.
func ApplyLight(ctx context.Context, database statedb.DBTX, campaign, actor uuid.UUID, now int, sheet *srd521.CharacterSheet) error {
	if !sheet.Light.Active(now) {
		sheet.Light = nil
		return nil
	}
	item, err := uuid.Parse(sheet.Light.ItemID)
	if err != nil {
		return nil
	}
	sheet.LightItem, _, err = lightItemName(ctx, database, campaign, actor, item, false)
	return err
}

// LightColor returns the explicit color named by a Light cast entry.
func LightColor(entryID string) (string, bool) {
	color, ok := strings.CutPrefix(entryID, LightEntryPrefix)
	return color, ok && slices.Contains(srd521.LightColors, color)
}

// LightEmitter is the observer's active Light on an object they still carry.
type LightEmitter struct {
	EffectID   string `json:"effect_id"`
	ItemID     string `json:"item_id"`
	Item       string `json:"item"`
	BrightFeet int    `json:"bright_feet"`
	DimFeet    int    `json:"dim_feet"`
}

func readLightEffect(stats []byte) *srd521.LightEffect {
	var envelope struct {
		Build srd521.Build `json:"srd521"`
	}
	if json.Unmarshal(stats, &envelope) != nil || envelope.Build.Casting == nil || envelope.Build.Casting.Light.Validate() != nil {
		return nil
	}
	return envelope.Build.Casting.Light
}

func lightItemName(ctx context.Context, database statedb.DBTX, campaign, actor, item uuid.UUID, lock bool) (string, bool, error) {
	locking := ""
	if lock {
		locking = " FOR UPDATE"
	}
	rows, err := database.Query(ctx, `SELECT name FROM items WHERE id=$1 AND campaign_id=$2 AND player_character_id=$3 AND quantity=1 AND item_type IN `+lightObjectTypes+locking, dbutil.ToPgtype(item), dbutil.ToPgtype(campaign), dbutil.ToPgtype(actor))
	if err != nil {
		return "", false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return "", false, rows.Err()
	}
	var name string
	if err := rows.Scan(&name); err != nil {
		return "", false, err
	}
	return name, true, rows.Err()
}

// ReadLightEmitter returns the actor's Light if it lasts at now, is not declared
// covered and its object is still carried by the actor. A covered, dropped,
// given-away or used-up object sheds no light here.
func ReadLightEmitter(ctx context.Context, database statedb.DBTX, campaignID uuid.UUID, pc domain.PlayerCharacter, now int, lock bool) (*LightEmitter, error) {
	light := readLightEffect(pc.Stats)
	if !light.Active(now) || light.Covered {
		return nil, nil
	}
	item, err := uuid.Parse(light.ItemID)
	if err != nil {
		return nil, nil
	}
	name, ok, err := lightItemName(ctx, database, campaignID, pc.ID, item, lock)
	if err != nil || !ok {
		return nil, err
	}
	return &LightEmitter{EffectID: light.ID, ItemID: light.ItemID, Item: name, BrightFeet: light.BrightFeet, DimFeet: light.DimFeet}, nil
}

// illuminate raises the light at a target's space by the observer's carried
// Light, then the unchanged #145 derivation decides sight.
func illuminate(space combat.Space, emitter *LightEmitter) (combat.Space, bool) {
	if emitter == nil {
		return space, false
	}
	return space.Illuminated(emitter.BrightFeet, emitter.DimFeet)
}

func prepareLight(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	p, err := PrepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "wisdom"})
	if err != nil {
		return nil, err
	}
	p.Intent = intent
	result := &StoredAction{Preview: p}
	if !p.Available {
		return result, nil
	}
	p.Available, p.Modifier = false, nil
	p.Reason = "Choose Light with an explicit color on one object you carry."
	if intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		return result, nil
	}
	if result.Clock, err = ReadSpellClock(ctx, database, campaign.ID); err != nil {
		return nil, err
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	sheet, err := srd521.Project(pc.Stats, pc.Level, pc.HP, pc.MaxHP, pc.Experience)
	if err != nil || sheet.Spellcasting == nil {
		return result, nil
	}
	var spell *srd521.Spell
	for i := range sheet.Spellcasting.Spells {
		if sheet.Spellcasting.Spells[i].ID == srd521.LightSpell {
			spell = &sheet.Spellcasting.Spells[i]
		}
	}
	if spell == nil {
		p.Reason = "This caster does not know Light."
		return result, nil
	}
	current := readLightEffect(pc.Stats)
	if !current.Active(result.Clock) {
		current = nil
	}
	currentItem, carried := "", false
	if current != nil {
		if id, err := uuid.Parse(current.ItemID); err == nil {
			name, ok, err := lightItemName(ctx, database, campaign.ID, pc.ID, id, lock)
			if err != nil {
				return nil, err
			}
			if ok {
				currentItem, carried = name, true
			}
		}
		if currentItem == "" {
			currentItem = "an object you no longer carry"
		}
	}
	result.Light, result.LightItem = current, currentItem
	if intent.Kind == srd521.SpellCast {
		result.LightObjects, err = readLightObjects(ctx, database, campaign.ID, pc.ID, lock)
		if err != nil {
			return nil, err
		}
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, result.Clock, encounter, current, currentItem, result.LightObjects})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("light-v1:%x", sha256.Sum256(encoded))
	if encounter != nil && encounter.Status == "active" {
		p.Reason = "Light is supported outside encounters."
		return result, nil
	}
	p.CostStatus = "fixture-only"
	p.SourceIDs = []string{"srd-5.2.1:content/spells/light", "srd-5.2.1:content/feats/magic-initiate"}
	if isLightCover(intent.Kind) {
		if intent.EntryID != "" {
			return result, nil
		}
		covered := intent.Kind == srd521.CoverLight
		switch {
		case current == nil:
			p.Reason = "No active Light to cover or uncover."
		case !carried:
			p.Reason = "Only the holder can cover or uncover the lit object; you no longer carry it."
		case current.Covered && covered:
			p.Reason = "Your Light is already covered."
		case !current.Covered && !covered:
			p.Reason = "Your Light is not covered."
		default:
			option := "Declare that nothing opaque covers your Light on %s any more; it sheds light again."
			if covered {
				option = "Declare that something opaque now covers your Light on %s; it sheds no light until uncovered."
			}
			p.Available, p.Reason = true, ""
			p.Options = []string{fmt.Sprintf(option, currentItem) + " Expiry unchanged; no action cost is tracked for this declaration."}
		}
		return result, nil
	}
	if intent.Kind == srd521.DismissLight {
		if intent.EntryID != "" {
			return result, nil
		}
		if current == nil {
			p.Reason = "No active Light to dismiss."
			return result, nil
		}
		p.Available, p.Reason = true, ""
		p.Options = []string{fmt.Sprintf("Dismiss your Light on %s (%s, until campaign second %d); no action or slot cost.", currentItem, current.Color, current.ExpiresAt)}
		return result, nil
	}
	color, ok := LightColor(intent.EntryID)
	if !ok {
		p.Reason = "Choose an explicit Light color: " + strings.Join(srd521.LightColors, ", ") + "."
		return result, nil
	}
	if !spell.ComponentsAvailable {
		p.Reason = "Light requires its verbal component and a focus for its material component."
		return result, nil
	}
	if len(result.LightObjects) == 0 {
		p.Reason = "Light needs one object you carry. Objects held by anyone else are unavailable."
		return result, nil
	}
	for _, object := range result.LightObjects {
		p.Targets = append(p.Targets, ActionTarget{ID: object.ID, Label: object.Name + " (carried; touch)"})
	}
	option := fmt.Sprintf("Light (Magic Initiate cantrip); Magic action; verbal component and Holy Symbol focus; no slot. Touch one object you carry: %s bright light in a 20-foot radius and dim light 20 feet further for one hour (%d campaign seconds), or until you cast it again or dismiss it.", color, srd521.LightDurationSeconds)
	if current != nil {
		option += fmt.Sprintf(" Ends your current Light on %s.", currentItem)
	}
	p.Options = []string{option}
	p.Costs = []ActionCost{{ResourceID: "magic-action", Amount: 1}}
	p.Available, p.Reason = true, ""
	return result, nil
}

func readLightObjects(ctx context.Context, database statedb.DBTX, campaign, actor uuid.UUID, lock bool) ([]LightObject, error) {
	locking := ""
	if lock {
		locking = " FOR UPDATE"
	}
	rows, err := database.Query(ctx, `SELECT id,name,updated_at FROM items WHERE campaign_id=$1 AND player_character_id=$2 AND quantity=1 AND item_type IN `+lightObjectTypes+` ORDER BY id`+locking, dbutil.ToPgtype(campaign), dbutil.ToPgtype(actor))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	objects := []LightObject{}
	for rows.Next() {
		var object LightObject
		if err := rows.Scan(&object.ID, &object.Name, &object.UpdatedAt); err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, rows.Err()
}
