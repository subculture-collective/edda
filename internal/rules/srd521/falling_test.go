package srd521

import "testing"

// fallDice returns each queued d6 face (1-6) and records how many were rolled.
type fallDice struct {
	faces []int
	calls int
}

func (d *fallDice) Intn(n int) int {
	d.calls++
	if n != 6 || len(d.faces) == 0 {
		return -1
	}
	face := d.faces[0]
	d.faces = d.faces[1:]
	return face - 1
}

func faces(n, face int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = face
	}
	return out
}

func TestFallDiceBoundariesAndCap(t *testing.T) {
	for height, want := range map[int]int{0: 0, 9: 0, 10: 1, 19: 1, 20: 2, 199: 19, 200: 20, 210: 20, 5000: 20} {
		if got := FallDice(height); got != want {
			t.Errorf("FallDice(%d)=%d, want %d", height, got, want)
		}
	}
}

func TestResolveFall(t *testing.T) {
	tests := []struct {
		name               string
		fall               Fall
		faces              []int
		rolled, damage, hp int
		prone              bool
	}{
		{"under 10 feet deals no damage and no Prone", Fall{HeightFeet: 9, Landing: FallSolid, HP: 12}, nil, 0, 0, 12, false},
		{"10 feet is one die and Prone", Fall{HeightFeet: 10, Landing: FallSolid, HP: 12}, []int{4}, 4, 4, 8, true},
		{"35 feet rounds down to 3 dice", Fall{HeightFeet: 35, Landing: FallSolid, HP: 12}, []int{1, 2, 3}, 6, 6, 6, true},
		{"maximum 20 dice; HP stops at 0", Fall{HeightFeet: 500, Landing: FallSolid, HP: 30}, faces(20, 6), 120, 120, 0, true},
		{"water success halves, rounding down", Fall{HeightFeet: 30, Landing: FallWater, Halved: true, HP: 12}, []int{6, 2, 3}, 11, 5, 7, true},
		{"halved 1 avoids all damage and Prone", Fall{HeightFeet: 10, Landing: FallWater, Halved: true, HP: 12}, []int{1}, 1, 0, 12, false},
		{"water without the check is full damage", Fall{HeightFeet: 20, Landing: FallWater, HP: 12}, []int{5, 5}, 10, 10, 2, true},
		{"already Prone stays Prone without damage", Fall{HeightFeet: 5, Landing: FallSolid, HP: 12, ProneBefore: true}, nil, 0, 0, 12, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dice := &fallDice{faces: tt.faces}
			r, err := ResolveFall(tt.fall, dice)
			if err != nil {
				t.Fatal(err)
			}
			if dice.calls != len(tt.faces) || r.DiceCount != len(tt.faces) || len(r.Dice) != len(tt.faces) {
				t.Fatalf("dice: calls=%d receipt=%+v", dice.calls, r)
			}
			if r.Rolled != tt.rolled || r.Damage != tt.damage || r.HPAfter != tt.hp || r.ProneAfter != tt.prone || r.DamageType != "bludgeoning" || r.HPBefore != tt.fall.HP {
				t.Fatalf("receipt: %+v", r)
			}
			if r.SourceIDs[0] != "srd-5.2.1:content/rules-glossary/falling" {
				t.Fatalf("sources: %v", r.SourceIDs)
			}
		})
	}
}

func TestResolveFallRejectsUnsupportedInputs(t *testing.T) {
	for name, f := range map[string]Fall{
		"negative height": {HeightFeet: -1, Landing: FallSolid, HP: 5},
		"0 HP":            {HeightFeet: 20, Landing: FallSolid},
		"unknown landing": {HeightFeet: 20, Landing: "lava", HP: 5},
		"halved on solid": {HeightFeet: 20, Landing: FallSolid, Halved: true, HP: 5},
	} {
		if _, err := ResolveFall(f, &fallDice{faces: faces(2, 3)}); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	if _, err := ResolveFall(Fall{HeightFeet: 20, Landing: FallSolid, HP: 5}, &fallDice{faces: []int{3}}); err == nil {
		t.Error("invalid second die accepted")
	}
}
