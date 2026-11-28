package srd521

import "fmt"

// Falling (SRD 5.2.1 Rules Glossary, "Falling [Hazard]", pinned
// PDF, p. 182): 1d6 Bludgeoning damage for every 10 feet fallen, to a maximum of 20d6;
// the creature lands Prone unless it avoids taking any damage. A creature that
// falls into water can use its Reaction for a DC 15 Strength (Athletics) or
// Dexterity (Acrobatics) check; on a success the damage is halved.
const (
	FallResolve  = "fall-resolve"
	FallMaxDice  = 20
	FallWaterDC  = 15
	FallSolid    = "solid"
	FallWater    = "water"
	fallSourceID = "srd-5.2.1:content/rules-glossary/falling"
)

// FallDice is the number of d6s for a fall: one per full 10 feet, capped at 20.
// Falls shorter than 10 feet deal no damage.
func FallDice(heightFeet int) int { return min(max(heightFeet, 0)/10, FallMaxDice) }

// PendingFall is the declared fall waiting for the actor. A client can only
// name it by ID and, for water, choose the entry check. HeightFeet and Landing
// are set only for a declaration the supported rules accept.
type PendingFall struct {
	ID         string `json:"id"`
	HeightFeet *int   `json:"height_feet,omitempty"`
	Landing    string `json:"landing,omitempty"`
}

// Fall is one declared fall at the moment of landing. Halved is the trusted
// result of the water-entry check; callers never pass a damage amount.
type Fall struct {
	HeightFeet  int
	Landing     string
	Halved      bool
	HP          int
	ProneBefore bool
}

// FallReceipt is the player-visible, inspectable result of one fall.
type FallReceipt struct {
	RulesVersion string        `json:"rules_version"`
	HeightFeet   int           `json:"height_feet"`
	Landing      string        `json:"landing"`
	DiceCount    int           `json:"dice_count"`
	Dice         []int         `json:"dice"`
	Rolled       int           `json:"rolled"`
	Halved       bool          `json:"halved"`
	Damage       int           `json:"damage"`
	DamageType   string        `json:"damage_type"`
	HPBefore     int           `json:"hp_before"`
	HPAfter      int           `json:"hp_after"`
	ProneBefore  bool          `json:"prone_before"`
	ProneAfter   bool          `json:"prone_after"`
	WaterCheck   *CheckReceipt `json:"water_check,omitempty"`
	// Concentration is the caster's Constitution save after the damage.
	Concentration *EffectReceipt `json:"concentration,omitempty"`
	// ZeroHP records dropping to 0 HP: Unconscious, or death by massive damage.
	ZeroHP    *ZeroHPReceipt `json:"zero_hp,omitempty"`
	SourceIDs []string       `json:"source_ids"`
}

// ResolveFall rolls the capped falling dice and derives damage and Prone.
// Halving rounds down, so a halved 1 deals no damage and leaves the creature
// standing. HP stops at 0; what 0 HP means is not decided here.
func ResolveFall(f Fall, dice Dice) (*FallReceipt, error) {
	if f.HeightFeet < 0 || f.HP <= 0 || (f.Landing != FallSolid && f.Landing != FallWater) || (f.Halved && f.Landing != FallWater) {
		return nil, fmt.Errorf("unsupported fall")
	}
	r := &FallReceipt{RulesVersion: "srd-5.2.1", HeightFeet: f.HeightFeet, Landing: f.Landing, DiceCount: FallDice(f.HeightFeet), Dice: []int{}, Halved: f.Halved, DamageType: "bludgeoning", HPBefore: f.HP, ProneBefore: f.ProneBefore, SourceIDs: []string{fallSourceID}}
	for range r.DiceCount {
		die := dice.Intn(6) + 1
		if die < 1 || die > 6 {
			return nil, fmt.Errorf("invalid falling die")
		}
		r.Dice = append(r.Dice, die)
		r.Rolled += die
	}
	r.Damage = r.Rolled
	if f.Halved {
		r.Damage /= 2
	}
	r.HPAfter = max(0, f.HP-r.Damage)
	r.ProneAfter = f.ProneBefore || r.Damage > 0
	if r.ProneAfter {
		r.SourceIDs = append(r.SourceIDs, "srd-5.2.1:content/rules-glossary/prone")
	}
	return r, nil
}
