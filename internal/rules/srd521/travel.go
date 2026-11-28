package srd521

import (
	"fmt"
	"slices"
)

// TravelHour is one explicitly confirmed travel segment of at most one hour
// along a known connected route with declared distance and terrain.
const TravelHour = "travel-hour"

// TravelRouteProfile marks server-owned route facts on a location connection.
const TravelRouteProfile = "travel-route-v1"

// Travel paces in step order, slowest first (SRD 5.2.1 p. 12).
const (
	PaceSlow   = "slow"
	PaceNormal = "normal"
	PaceFast   = "fast"
)

// TravelPaces lists the paces in step order, slowest first.
var TravelPaces = []string{PaceSlow, PaceNormal, PaceFast}

// travelMilesPerHour is the Hour column of the Travel Pace table (p. 12).
var travelMilesPerHour = map[string]int{PaceFast: 4, PaceNormal: 3, PaceSlow: 2}

// travelPaceEffects is the game effect of each pace (p. 12).
var travelPaceEffects = map[string][]string{
	PaceFast:   {"Disadvantage on Wisdom (Perception or Survival) checks", "Disadvantage on Dexterity (Stealth) checks"},
	PaceNormal: {"Disadvantage on Dexterity (Stealth) checks"},
	PaceSlow:   {"Advantage on Wisdom (Perception or Survival) checks"},
}

// travelMaximumPace is the Maximum Pace column of the Travel Terrain table
// (p. 192). Waterborne travel uses a vehicle's speed and is not listed here.
var travelMaximumPace = map[string]string{
	"arctic": PaceFast, "coastal": PaceNormal, "desert": PaceNormal, "forest": PaceNormal,
	"grassland": PaceFast, "hill": PaceNormal, "mountain": PaceSlow, "swamp": PaceSlow,
	"underdark": PaceNormal, "urban": PaceNormal,
}

// TravelMilesPerHour returns the Hour column for a pace, or 0.
func TravelMilesPerHour(pace string) int { return travelMilesPerHour[pace] }

// TravelPaceEffects returns a copy of the game effects of a pace.
func TravelPaceEffects(pace string) []string { return slices.Clone(travelPaceEffects[pace]) }

// TravelDaySeconds is the 8 hours of travel per day before Extended Travel
// (p. 192) requires Constitution saves. Edda does not model Extended Travel.
const TravelDaySeconds = 8 * 3600

const feetPerMile = 5280

// TravelRoute is the declared fact set for one connection. Nothing is inferred
// from the connection's travel_time prose.
type TravelRoute struct {
	Profile      string `json:"profile"`
	DistanceFeet *int   `json:"distance_feet"`
	Terrain      string `json:"terrain"`
	GoodRoad     *bool  `json:"good_road"`
}

// TravelPlan is the derived, inspectable result of one route at one pace.
type TravelPlan struct {
	Pace         string   `json:"pace"`
	MaximumPace  string   `json:"maximum_pace"`
	DistanceFeet int      `json:"distance_feet"`
	Terrain      string   `json:"terrain"`
	GoodRoad     bool     `json:"good_road"`
	MilesPerHour int      `json:"miles_per_hour"`
	Seconds      int      `json:"seconds"`
	Effects      []string `json:"effects"`
}

// TravelUnresolved is returned when a route's facts are missing or unknown.
type TravelUnresolved struct{ Reason string }

func (e TravelUnresolved) Error() string { return e.Reason }

// PlanTravel derives one segment, failing closed on undeclared or unknown
// facts and on routes that one hour at the chosen pace cannot cover.
func PlanTravel(route TravelRoute, pace string) (*TravelPlan, error) {
	if route.Profile != TravelRouteProfile || route.DistanceFeet == nil || route.GoodRoad == nil || route.Terrain == "" {
		return nil, TravelUnresolved{"Route distance, terrain and road are not declared; travel time prose is not an SRD distance."}
	}
	if *route.DistanceFeet <= 0 {
		return nil, TravelUnresolved{"Route distance must be a positive number of feet."}
	}
	if route.Terrain == "waterborne" {
		return nil, TravelUnresolved{"Waterborne travel depends on a vehicle; vehicles are not supported."}
	}
	maximum, ok := travelMaximumPace[route.Terrain]
	if !ok {
		return nil, TravelUnresolved{fmt.Sprintf("Terrain %q is not a Travel Terrain entry.", route.Terrain)}
	}
	mph, ok := travelMilesPerHour[pace]
	if !ok {
		return nil, fmt.Errorf("choose a fast, normal or slow pace")
	}
	step := slices.Index(TravelPaces, maximum)
	if *route.GoodRoad && step < len(TravelPaces)-1 {
		step++
	}
	maximum = TravelPaces[step]
	plan := &TravelPlan{Pace: pace, MaximumPace: maximum, DistanceFeet: *route.DistanceFeet, Terrain: route.Terrain, GoodRoad: *route.GoodRoad, MilesPerHour: mph, Effects: TravelPaceEffects(pace)}
	if slices.Index(TravelPaces, pace) > step {
		return plan, fmt.Errorf("%s terrain allows at most a %s pace", route.Terrain, maximum)
	}
	if route.Terrain == "arctic" && pace == PaceFast {
		return plan, fmt.Errorf("a fast pace in arctic terrain needs equipment such as skis, which is not modelled")
	}
	if plan.DistanceFeet > mph*feetPerMile {
		return plan, fmt.Errorf("the route is longer than %d miles, one hour at a %s pace", mph, pace)
	}
	// Ceiling of distance / (mph miles per hour), in whole campaign seconds.
	plan.Seconds = (plan.DistanceFeet*3600 + mph*feetPerMile - 1) / (mph * feetPerMile)
	return plan, nil
}

// TravelDay is the server-owned record of travel time on one campaign day,
// stored at stats.srd521_travel.
type TravelDay struct {
	Day     int `json:"day"`
	Seconds int `json:"seconds"`
}

// TravelledOn returns the travel seconds already recorded on the day of now.
func (d *TravelDay) TravelledOn(now int) int {
	if d == nil || d.Day != now/86400+1 {
		return 0
	}
	return d.Seconds
}

// TravelReceipt is the public, saved account of one travel segment.
type TravelReceipt struct {
	RulesVersion    string         `json:"rules_version"`
	Kind            string         `json:"kind"`
	Reason          string         `json:"reason"`
	RouteID         string         `json:"route_id"`
	FromID          string         `json:"from_id"`
	From            string         `json:"from"`
	ToID            string         `json:"to_id"`
	To              string         `json:"to"`
	Plan            TravelPlan     `json:"plan"`
	TimeBefore      int            `json:"time_before"`
	TimeAfter       int            `json:"time_after"`
	DayTravelBefore int            `json:"day_travel_before"`
	DayTravelAfter  int            `json:"day_travel_after"`
	Expired         *Concentration `json:"expired,omitempty"`
	SourceIDs       []string       `json:"source_ids"`
}

// TravelSourceIDs are the pinned sections used by a travel segment.
func TravelSourceIDs() []string {
	return []string{prefix + "playing-the-game/travel-pace", prefix + "gameplay-toolbox/travel-pace", prefix + "gameplay-toolbox/travel-terrain"}
}
