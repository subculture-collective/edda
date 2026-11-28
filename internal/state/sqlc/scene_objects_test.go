package statedb

import "testing"

func TestIsSceneObject(t *testing.T) {
	for props, want := range map[string]bool{
		`{"srd521_door":{"profile":"ordinary-door-v1"}}`:           true,
		`{"srd521_magic":{"profile":"detect-magic-v1"}}`:           false,
		`{"srd521_target":{"profile":"bandit-v1"},"door":"front"}`: false,
		`{}`: false,
		``:   false,
	} {
		if got := IsSceneObject([]byte(props)); got != want {
			t.Errorf("IsSceneObject(%s) = %v, want %v", props, got, want)
		}
	}
}
