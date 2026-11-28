package srd521

import (
	"errors"
	"testing"
)

func route(feet int, terrain string, road bool) TravelRoute {
	return TravelRoute{Profile: TravelRouteProfile, DistanceFeet: &feet, Terrain: terrain, GoodRoad: &road}
}

func TestPlanTravelPaceBoundaries(t *testing.T) {
	for _, tc := range []struct {
		pace    string
		miles   int
		seconds int
	}{{PaceFast, 4, 3600}, {PaceNormal, 3, 3600}, {PaceSlow, 2, 3600}} {
		feet := tc.miles * 5280
		plan, err := PlanTravel(route(feet, "grassland", false), tc.pace)
		if err != nil || plan.Seconds != tc.seconds || plan.MilesPerHour != tc.miles {
			t.Fatalf("%s at %d feet: %+v %v", tc.pace, feet, plan, err)
		}
		if _, err := PlanTravel(route(feet+1, "grassland", false), tc.pace); err == nil {
			t.Fatalf("%s accepted one foot over one hour", tc.pace)
		}
	}
	// Exact seconds are rounded up, never down.
	plan, err := PlanTravel(route(1, "grassland", false), PaceNormal)
	if err != nil || plan.Seconds != 1 {
		t.Fatalf("one foot: %+v %v", plan, err)
	}
	plan, err = PlanTravel(route(5280, "grassland", false), PaceNormal)
	if err != nil || plan.Seconds != 1200 {
		t.Fatalf("one mile at Normal: %+v %v", plan, err)
	}
}

func TestPlanTravelTerrainMaximumPaceAndRoads(t *testing.T) {
	for _, tc := range []struct {
		terrain string
		road    bool
		maximum string
	}{
		{"grassland", false, PaceFast},
		{"forest", false, PaceNormal},
		{"forest", true, PaceFast},
		{"mountain", false, PaceSlow},
		{"mountain", true, PaceNormal},
		{"swamp", false, PaceSlow},
		{"urban", false, PaceNormal},
	} {
		for _, pace := range TravelPaces {
			plan, err := PlanTravel(route(5280, tc.terrain, tc.road), pace)
			if plan == nil || plan.MaximumPace != tc.maximum {
				t.Fatalf("%s road=%t %s: %+v %v", tc.terrain, tc.road, pace, plan, err)
			}
			allowed := pace == PaceSlow || pace == PaceNormal && tc.maximum != PaceSlow || pace == tc.maximum
			if (err == nil) != allowed {
				t.Fatalf("%s road=%t %s: allowed=%t err=%v", tc.terrain, tc.road, pace, allowed, err)
			}
		}
	}
	// Arctic Fast needs unmodelled equipment, even on a road.
	if _, err := PlanTravel(route(5280, "arctic", true), PaceFast); err == nil {
		t.Fatal("arctic fast pace accepted")
	}
	if _, err := PlanTravel(route(5280, "arctic", false), PaceNormal); err != nil {
		t.Fatal(err)
	}
}

func TestPlanTravelPaceEffects(t *testing.T) {
	for pace, want := range map[string]int{PaceFast: 2, PaceNormal: 1, PaceSlow: 1} {
		plan, err := PlanTravel(route(5280, "grassland", false), pace)
		if err != nil || len(plan.Effects) != want {
			t.Fatalf("%s: %+v %v", pace, plan, err)
		}
	}
}

func TestPlanTravelUnresolvedFactsFailClosed(t *testing.T) {
	feet, road := 5280, false
	zero := 0
	for name, r := range map[string]TravelRoute{
		"empty":      {},
		"profile":    {Profile: "other", DistanceFeet: &feet, Terrain: "forest", GoodRoad: &road},
		"distance":   {Profile: TravelRouteProfile, Terrain: "forest", GoodRoad: &road},
		"zero":       {Profile: TravelRouteProfile, DistanceFeet: &zero, Terrain: "forest", GoodRoad: &road},
		"road":       {Profile: TravelRouteProfile, DistanceFeet: &feet, Terrain: "forest"},
		"terrain":    {Profile: TravelRouteProfile, DistanceFeet: &feet, GoodRoad: &road},
		"unknown":    route(feet, "jungle", false),
		"waterborne": route(feet, "waterborne", false),
	} {
		plan, err := PlanTravel(r, PaceNormal)
		var unresolved TravelUnresolved
		if plan != nil || !errors.As(err, &unresolved) {
			t.Fatalf("%s: %+v %v", name, plan, err)
		}
	}
	if _, err := PlanTravel(route(feet, "forest", false), "forced-march"); err == nil {
		t.Fatal("unknown pace accepted")
	}
}

func TestTravelDayCounter(t *testing.T) {
	day := &TravelDay{Day: 2, Seconds: 7200}
	if day.TravelledOn(86400) != 7200 || day.TravelledOn(86399) != 0 || (*TravelDay)(nil).TravelledOn(0) != 0 {
		t.Fatal("travel day counter")
	}
}
