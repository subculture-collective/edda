package srd521

import (
	_ "embed"
	"fmt"
)

//go:embed testdata/items_246.json
var items246JSON []byte

// ItemDetails246 holds exactly one typed payload per variant of catalogue batch
// #246 (magic items 22, SRD 5.2.1 p. 253). Payloads are source facts only.
type ItemDetails246 struct {
	Warning     *ItemWeaponOfWarning  `json:"warning,omitempty"`
	WeaponBonus *ItemWeaponBonus      `json:"weapon_bonus,omitempty"`
	Well        *ItemWellOfManyWorlds `json:"well,omitempty"`
	WindFan     *ItemWindFan          `json:"wind_fan,omitempty"`
	WingedBoots *ItemWingedBoots      `json:"winged_boots,omitempty"`
	Wings       *ItemWingsOfFlying    `json:"wings,omitempty"`
}

// ItemWeaponOfWarning retains typed SRD source facts without implementing them.
type ItemWeaponOfWarning struct {
	Condition             string                                   `json:"condition"`
	RadiusFeet            int                                      `json:"radius_feet"`
	Subjects              string                                   `json:"subjects"`
	Alarm                 ItemWeaponOfWarningAlarm                 `json:"alarm"`
	SupernaturalReadiness ItemWeaponOfWarningSupernaturalReadiness `json:"supernatural_readiness"`
}

// ItemWeaponOfWarningAlarm is a nested part of a source payload.
type ItemWeaponOfWarningAlarm struct {
	Wakes                 string `json:"wakes"`
	When                  string `json:"when"`
	WakesFromMagicalSleep bool   `json:"wakes_from_magical_sleep"`
}

// ItemWeaponOfWarningSupernaturalReadiness is a nested part of a source payload.
type ItemWeaponOfWarningSupernaturalReadiness struct {
	Roll string `json:"roll"`
	Mode string `json:"mode"`
}

// ItemWeaponBonus retains typed SRD source facts without implementing them.
type ItemWeaponBonus struct {
	AttackBonus int    `json:"attack_bonus"`
	DamageBonus int    `json:"damage_bonus"`
	BonusRule   string `json:"bonus_rule"`
}

// ItemWellOfManyWorlds retains typed SRD source facts without implementing them.
type ItemWellOfManyWorlds struct {
	Material             string                     `json:"material"`
	FoldedSize           string                     `json:"folded_size"`
	UnfoldedDiameterFeet int                        `json:"unfolded_diameter_feet"`
	Open                 ItemWellOfManyWorldsOpen   `json:"open"`
	Portal               ItemWellOfManyWorldsPortal `json:"portal"`
	Close                ItemWellOfManyWorldsClose  `json:"close"`
	ReuseAfter           ItemDice                   `json:"reuse_after"`
}

// ItemWellOfManyWorldsOpen is a nested part of a source payload.
type ItemWellOfManyWorldsOpen struct {
	Action    string `json:"action"`
	Placement string `json:"placement"`
}

// ItemWellOfManyWorldsPortal is a nested part of a source payload.
type ItemWellOfManyWorldsPortal struct {
	Direction    string `json:"direction"`
	DiameterFeet int    `json:"diameter_feet"`
	Destination  string `json:"destination"`
	ChosenBy     string `json:"chosen_by"`
}

// ItemWellOfManyWorldsClose is a nested part of a source payload.
type ItemWellOfManyWorldsClose struct {
	Action     string `json:"action"`
	WithinFeet int    `json:"within_feet"`
	Method     string `json:"method"`
}

// ItemWindFan retains typed SRD source facts without implementing them.
type ItemWindFan struct {
	Condition string               `json:"condition"`
	Spell     string               `json:"spell"`
	SaveDC    int                  `json:"save_dc"`
	RepeatUse ItemWindFanRepeatUse `json:"repeat_use"`
}

// ItemWindFanRepeatUse is a nested part of a source payload.
type ItemWindFanRepeatUse struct {
	FailurePercentCumulative int    `json:"failure_percent_cumulative"`
	AppliesTo                string `json:"applies_to"`
	OnFailure                string `json:"on_failure"`
}

// ItemWingedBoots retains typed SRD source facts without implementing them.
type ItemWingedBoots struct {
	MaxCharges    int                     `json:"max_charges"`
	RechargeDice  ItemDice                `json:"recharge_dice"`
	Action        string                  `json:"action"`
	ChargeCost    int                     `json:"charge_cost"`
	FlySpeedFeet  int                     `json:"fly_speed_feet"`
	DurationHours int                     `json:"duration_hours"`
	OnExpiry      ItemWingedBootsOnExpiry `json:"on_expiry"`
}

// ItemWingedBootsOnExpiry is a nested part of a source payload.
type ItemWingedBootsOnExpiry struct {
	Condition           string `json:"condition"`
	DescentFeetPerRound int    `json:"descent_feet_per_round"`
}

// ItemWingsOfFlying retains typed SRD source facts without implementing them.
type ItemWingsOfFlying struct {
	Form                    string   `json:"form"`
	Action                  string   `json:"action"`
	DurationHours           int      `json:"duration_hours"`
	EarlyEndAction          string   `json:"early_end_action"`
	FlySpeedFeet            int      `json:"fly_speed_feet"`
	AloftWhenWingsDisappear string   `json:"aloft_when_wings_disappear"`
	ReuseAfter              ItemDice `json:"reuse_after"`
}

// itemManifest246 is the reviewed manifest for the batch's families, in
// planning inventory order.
var itemManifest246 = itemSourceManifest{Batch: "data-magic-items-22", Issue: 246, Families: []itemSourceManifestFamily{
	{Slug: "weapon-of-warning", Name: "Weapon of Warning", Page: 253, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "required", Equip: "within-reach", Use: "none",
		Keys: "item.attunement item.magic-weapon effect.wake-natural-sleepers-at-combat-start roll.initiative-advantage aura.allies-within-30-feet", Payload: "warning",
		Variants: []itemSourceManifestVariant{{"standard", "Weapon of Warning", "uncommon"}}},
	{Slug: "weapon-1-2-or-3", Name: "Weapon, +1, +2, or +3", Page: 253, Category: "weapon", Options: "any-simple-weapon any-martial-weapon", Attunement: "none", Equip: "wield", Use: "weapon-attack",
		Keys: "item.magic-weapon bonus.attack-and-damage-by-rarity", Payload: "weapon_bonus",
		Variants: []itemSourceManifestVariant{{"plus-1", "Weapon, +1", "uncommon"}, {"plus-2", "Weapon, +2", "rare"}, {"plus-3", "Weapon, +3", "very-rare"}}},
	{Slug: "well-of-many-worlds", Name: "Well of Many Worlds", Page: 253, Category: "wondrous-item", Attunement: "none", Equip: "none", Use: "magic-action", UseRestrictions: "unfold-and-place-on-solid-surface",
		Keys: "action.magic-unfold effect.two-way-planar-portal item.gm-chooses-destination action.magic-close-portal item.1d8-hour-reuse", Payload: "well",
		Variants: []itemSourceManifestVariant{{"standard", "Well of Many Worlds", "legendary"}}},
	{Slug: "wind-fan", Name: "Wind Fan", Page: 253, Category: "wondrous-item", Attunement: "none", Equip: "hold", Use: "cast-spell", UseRestrictions: "cumulative-failure-chance-before-next-dawn",
		Keys: "item.hold spell.cast-gust-of-wind spell.save-dc-13 item.cumulative-20-percent-failure item.destroyed-on-failure", Payload: "wind_fan",
		Variants: []itemSourceManifestVariant{{"standard", "Wind Fan", "uncommon"}}},
	{Slug: "winged-boots", Name: "Winged Boots", Page: 253, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "expend-one-charge",
		Charges: &ItemCharges{Mode: "limited", Maximum: 4, PerUse: 1, Recharge: "regain-1d4-at-dawn"},
		Keys:    "item.attunement item.wear item.charges action.magic movement.fly-speed-30 duration.1-hour movement.gentle-descent-30-feet-per-round", Payload: "winged_boots",
		Variants: []itemSourceManifestVariant{{"standard", "Winged Boots", "uncommon"}}},
	{Slug: "wings-of-flying", Name: "Wings of Flying", Page: 253, Category: "wondrous-item", Attunement: "required", Equip: "wear", Use: "magic-action", UseRestrictions: "1d12-hour-wait-after-wings-disappear",
		Keys: "item.attunement item.wear action.magic movement.fly-speed-60 duration.1-hour action.magic-end-early item.1d12-hour-reuse", Payload: "wings",
		Variants: []itemSourceManifestVariant{{"standard", "Wings of Flying", "rare"}}},
}}

// LoadItems246 returns a fresh validated copy of batch #246. It has no state,
// inventory, action registry or resolver side effects.
func LoadItems246() (ItemSourceBatch[ItemDetails246], error) { return DecodeItems246(items246JSON) }

// DecodeItems246 validates supplied JSON as batch #246 without network or
// database access. Validation is source validation, never gameplay support.
func DecodeItems246(data []byte) (ItemSourceBatch[ItemDetails246], error) {
	return decodeItemSourceBatch(data, itemManifest246, validateItemDetails246)
}

// validateItemDetails246 checks each payload's internal consistency and the
// source literals that the family's mechanic keys also name.
func validateItemDetails246(family, variant string, d *ItemDetails246) error {
	ok := true
	switch family {
	case "weapon-of-warning":
		p := d.Warning
		ok = p.RadiusFeet == 30 && !p.Alarm.WakesFromMagicalSleep && p.SupernaturalReadiness.Roll == "initiative" && p.SupernaturalReadiness.Mode == "advantage"
	case "weapon-1-2-or-3":
		p := d.WeaponBonus
		bonus := map[string]int{"plus-1": 1, "plus-2": 2, "plus-3": 3}[variant]
		ok = p.AttackBonus == bonus && p.DamageBonus == bonus
	case "well-of-many-worlds":
		p := d.Well
		ok = p.UnfoldedDiameterFeet == 6 && p.Portal.DiameterFeet == 6 && p.Close.WithinFeet == 5 && p.ReuseAfter == ItemDice{Dice: 1, Sides: 8, Type: "hours"}
	case "wind-fan":
		p := d.WindFan
		ok = p.Spell == "gust-of-wind" && p.SaveDC == 13 && p.RepeatUse.FailurePercentCumulative == 20
	case "winged-boots":
		p := d.WingedBoots
		ok = p.MaxCharges == 4 && p.FlySpeedFeet == 30 && p.DurationHours == 1 && p.OnExpiry.DescentFeetPerRound == 30 && p.RechargeDice == ItemDice{Dice: 1, Sides: 4, Type: "charges"}
	case "wings-of-flying":
		p := d.Wings
		ok = p.FlySpeedFeet == 60 && p.DurationHours == 1 && p.ReuseAfter == ItemDice{Dice: 1, Sides: 12, Type: "hours"} && p.AloftWhenWingsDisappear == "fall"
	}
	if !ok {
		return fmt.Errorf("incomplete or duplicated source values")
	}
	return nil
}
