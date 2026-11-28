package srd521

import "testing"

func TestConcentrationLabelNamesEachSpell(t *testing.T) {
	for id, want := range map[string]string{
		"shield-of-faith":               "Shield of Faith",
		DetectMagicID:                   "Detect Magic",
		"bless":                         "Bless",
		"protection-from-evil-and-good": "Protection from Evil and Good",
		"pass-without-trace":            "Pass without Trace",
		"meld-into-stone":               "Meld into Stone",
		"transport-via-plants":          "Transport via Plants",
		"":                              "Unnamed spell",
	} {
		if got := ConcentrationLabel(id); got != want {
			t.Errorf("ConcentrationLabel(%q) = %q, want %q", id, got, want)
		}
	}
}
