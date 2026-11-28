package srd521

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// Wolf and Brown Bear are the only animal blocks the bounded encounter runs
// (#179, #180). Their batch files stay outside LoadMonsters, so catalogue
// presence still grants nothing to the other records they contain.
const (
	WolfID      = "srd-5.2.1:content/animals/wolf"
	BrownBearID = "srd-5.2.1:content/animals/brown-bear"
)

//go:embed testdata/monsters_155.json
var monsters155JSON []byte

//go:embed testdata/monsters_177.json
var monsters177JSON []byte

var creatureSizes = []string{"tiny", "small", "medium", "large", "huge", "gargantuan"}

// beastFixture is the per-profile capability set. Every required mechanic key
// of the record must appear here, or the record refuses to load. Skills,
// Darkvision and Climb are satisfied because the encounter derives no check,
// sight ruling or path from them: line of sight and distance are explicit facts.
type beastFixture struct {
	id, action string
	data       []byte
	// maxProne is the largest target size the size-gated Prone effect affects;
	// gate is the exact source clause that states it.
	maxProne, gate string
	attacks        []string
	implemented    []string
}

var beastFixtures = map[string]beastFixture{
	"wolf-v1": {id: WolfID, action: "bite", data: monsters177JSON, maxProne: "medium", gate: "If the target is a Medium or smaller creature", attacks: []string{"bite"},
		implemented: []string{"monster/stat-block", "monster/skills", "sense/darkvision", "trait/pack-tactics", "attack/advantage-with-nearby-ally", "requirement/ally-not-incapacitated", "action/melee-attack", "effect/piercing-damage", "requirement/target-size", "condition/prone"}},
	// The source leaves the Multiattack order open; this fixture rules Bite then Claw.
	"brown-bear-v1": {id: BrownBearID, action: "multiattack", data: monsters155JSON, maxProne: "large", gate: "If the target is a Large or smaller creature", attacks: []string{"bite", "claw"},
		implemented: []string{"monster/stat-block", "monster/skills", "movement/climb", "sense/darkvision", "action/multiattack", "action/melee-attack", "effect/piercing-damage", "effect/slashing-damage", "condition/prone", "requirement/target-size"}},
}

// BeastProfile pins one Wolf or Brown Bear record revision, its configured
// action and the declared size of the player it attacks. The version prefix
// names the profile. It carries no HP or private tactics.
type BeastProfile struct {
	Version    string `json:"version"`
	ActionID   string `json:"action_id"`
	TargetSize string `json:"target_size"`
}

// BeastRecord returns the validated pinned record and content revision for a
// supported profile ("wolf-v1" or "brown-bear-v1").
func BeastRecord(profile string) (Monster, string, error) {
	f, ok := beastFixtures[profile]
	if !ok {
		return Monster{}, "", fmt.Errorf("unsupported beast profile")
	}
	c, err := DecodeMonsters(bytes.NewReader(f.data))
	if err != nil {
		return Monster{}, "", err
	}
	for _, m := range c.Monsters {
		if m.ID != f.id {
			continue
		}
		implemented := map[string]bool{}
		for _, key := range f.implemented {
			implemented[key] = true
		}
		if !m.availableWith(implemented) || !f.shape(m) {
			return Monster{}, "", fmt.Errorf("unsupported %s mechanics", m.Name)
		}
		data, err := json.Marshal(m)
		if err != nil {
			return Monster{}, "", err
		}
		return m, fmt.Sprintf("%s:%x", profile, sha256.Sum256(data)), nil
	}
	return Monster{}, "", fmt.Errorf("missing %s record", f.id)
}

// shape refuses any block whose traits, action list or size gate differ from
// what the runtime resolves, even when every key is known.
func (f beastFixture) shape(m Monster) bool {
	if len(m.BonusActions)+len(m.Reactions)+len(m.LegendaryActions) != 0 || m.Legendary != nil {
		return false
	}
	ids := []string{}
	for _, a := range m.Actions {
		ids = append(ids, a.ID)
	}
	want := f.attacks
	if f.action == "multiattack" {
		want = append([]string{"multiattack"}, f.attacks...)
	}
	if !slices.Equal(ids, want) {
		return false
	}
	gated := false
	for _, id := range f.attacks {
		a := beastAction(m, id)
		if a.Attack == nil || a.Attack.Kind != "melee" || a.Attack.Reach != 5 || a.Attack.Damage.Dice != 1 {
			return false
		}
		for _, e := range a.Effects {
			if e.Key == "requirement/target-size" {
				gated = gated || strings.Contains(e.Text, f.gate)
			}
		}
	}
	if !gated {
		return false
	}
	if f.id == WolfID {
		return len(m.Traits) == 1 && m.Traits[0].Key == "trait/pack-tactics" && slices.Equal(m.Traits[0].Requires, []string{"attack/advantage-with-nearby-ally", "requirement/ally-not-incapacitated"})
	}
	return len(m.Traits) == 0
}

func beastAction(m Monster, id string) MonsterAction {
	for _, a := range m.Actions {
		if a.ID == id {
			return a
		}
	}
	return MonsterAction{}
}

// Profile returns the profile name the version pins.
func (p BeastProfile) Profile() string {
	name, _, _ := strings.Cut(p.Version, ":")
	return name
}

// PackTactics reports whether the pinned block has the Pack Tactics trait.
func (p BeastProfile) PackTactics() bool { return p.Profile() == "wolf-v1" }

// Pinned validates the revision and configured action, not the target facts.
func (p BeastProfile) Pinned() (Monster, error) {
	m, version, err := BeastRecord(p.Profile())
	if err != nil {
		return Monster{}, err
	}
	if p.Version != version || p.ActionID != beastFixtures[p.Profile()].action {
		return Monster{}, fmt.Errorf("unsupported %s revision or action", m.Name)
	}
	return m, nil
}

// Record validates the revision, configured action and declared target size
// before any roll. It returns the ordered attacks the configured action makes.
func (p BeastProfile) Record() (Monster, []MonsterAction, error) {
	m, err := p.Pinned()
	if err != nil {
		return Monster{}, nil, err
	}
	f := beastFixtures[p.Profile()]
	if !slices.Contains(creatureSizes, p.TargetSize) {
		return Monster{}, nil, fmt.Errorf("Clarify the target's size before rolling.")
	}
	attacks := []MonsterAction{}
	for _, id := range f.attacks {
		attacks = append(attacks, beastAction(m, id))
	}
	return m, attacks, nil
}

// BeastAttack names one pinned attack resolution. Advantage lists the durable
// encounter facts that grant it; TargetProne is the target's saved condition.
type BeastAttack struct {
	ActionID    string
	ArmorClass  int
	Advantage   []string
	TargetProne bool
}

// ResolveBeastAttack rolls one attack of the configured action. A hit against
// a target within the size gate applies Prone without a saving throw; a target
// that is already Prone is not affected twice.
func ResolveBeastAttack(p BeastProfile, attack BeastAttack, dice Dice) (*AttackReceipt, error) {
	m, attacks, err := p.Record()
	if err != nil {
		return nil, err
	}
	var a *MonsterAction
	for i := range attacks {
		if attacks[i].ID == attack.ActionID {
			a = &attacks[i]
		}
	}
	if a == nil || attack.ArmorClass < 1 || dice == nil {
		return nil, fmt.Errorf("unsupported attack")
	}
	source := "animals/" + strings.TrimPrefix(m.ID, prefix+"animals/")
	w := a.Attack
	r, err := resolveAttackRoll("srd-5.2.1", a.Name, derived(term(m.Name+" attack", w.Bonus, source)), w.Damage.Modifier, w.Damage.Sides, w.DamageType, attack.ArmorClass, attack.Advantage, dice)
	if err != nil {
		return nil, err
	}
	r.SourceIDs = []string{m.ID, prefix + "d20-tests/attack-rolls", prefix + "damage/critical-hits"}
	if len(r.Advantage) > 0 {
		r.SourceIDs = append(r.SourceIDs, prefix+"d20-tests/advantage-disadvantage")
	}
	gated := false
	for _, e := range a.Effects {
		gated = gated || e.Key == "requirement/target-size"
	}
	if r.Outcome != "miss" && gated {
		r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/prone")
		switch {
		case slices.Index(creatureSizes, p.TargetSize) > slices.Index(creatureSizes, beastFixtures[p.Profile()].maxProne):
			r.Condition = "size-exempt"
		case attack.TargetProne:
			r.Condition = "already-prone"
		default:
			r.Condition = "prone-applied"
		}
	}
	return r, nil
}
