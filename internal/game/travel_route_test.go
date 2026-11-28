package game

import (
	"strings"
	"testing"
)

func TestParseTravelRouteMatchesThePreview(t *testing.T) {
	const ok = `"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest","good_road":false`
	tests := map[string]struct{ facts, reason string }{
		"supported":       {`{` + ok + `}`, ""},
		"too long is ok":  {`{"profile":"travel-route-v1","distance_feet":52800,"terrain":"hill","good_road":true}`, ""},
		"unknown key":     {`{` + ok + `,"mount":"horse"}`, "unknown or malformed"},
		"trailing data":   {`{` + ok + `} {}`, "unknown or malformed"},
		"trailing junk":   {`{` + ok + `} x`, "unknown or malformed"},
		"wrong type":      {`{"profile":"travel-route-v1","distance_feet":"far","terrain":"forest","good_road":false}`, "unknown or malformed"},
		"missing road":    {`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest"}`, "not declared"},
		"other profile":   {`{"profile":"travel-route-v2","distance_feet":5280,"terrain":"forest","good_road":false}`, "not declared"},
		"null":            {`null`, "not declared"},
		"zero distance":   {`{"profile":"travel-route-v1","distance_feet":0,"terrain":"forest","good_road":false}`, "positive"},
		"unknown terrain": {`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"lava","good_road":false}`, "Travel Terrain"},
		"waterborne":      {`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"waterborne","good_road":false}`, "vehicle"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ParseTravelRoute([]byte(test.facts))
			if test.reason == "" && err != nil || test.reason != "" && (err == nil || !strings.Contains(err.Error(), test.reason)) {
				t.Fatalf("want %q, got %v", test.reason, err)
			}
		})
	}
}
