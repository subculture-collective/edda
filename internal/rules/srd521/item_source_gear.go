package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// The adventuring gear payload (SRD 5.2.1 pp. 94-100) shared by the gear
// catalogue batches. A record describes one printed entry and its table rows.
// It is not read by any inventory, check, damage, light or spellcasting
// resolver.

const (
	gearActions    = "attack-replacement|attack-action|utilize|bonus-action|action"
	gearCheckKinds = "saving-throw|ability-check"
	gearSkills     = "acrobatics|athletics|sleight-of-hand|none"
	gearDCRules    = "fixed|plus-dexterity-modifier-and-proficiency-bonus"
	gearDamage     = "acid|fire|piercing|poison|radiant"
	gearUnits      = "item|flask|day|sheet"
	gearClasses    = "cleric|druid|paladin|ranger|sorcerer|warlock|wizard"
)

// ItemGear is one adventuring gear entry or one row of its variant table. Names,
// purposes and traits are keys; contents name gear by family slug but do not
// link to or grant other records.
type ItemGear struct {
	Actions           []string          `json:"actions" srd:"may-be-empty"`
	Checks            []ItemGearCheck   `json:"checks" srd:"may-be-empty"`
	Damage            []ItemGearDamage  `json:"damage" srd:"may-be-empty"`
	Measures          []ItemGearMeasure `json:"measures" srd:"may-be-empty"`
	Traits            []string          `json:"traits" srd:"may-be-empty"`
	Contents          []ItemGearContent `json:"contents" srd:"may-be-empty"`
	SpellcastingFocus []string          `json:"spellcasting_focus" srd:"may-be-empty"`
	Weight            ItemGearWeight    `json:"weight"`
	SourcePages       []int             `json:"source_pages"`
}

// ItemGearCheck is a saving throw or ability check the entry names. DCRule
// "fixed" means DC is the whole DC; otherwise the rule adds to it. A saving
// throw has the skill "none".
type ItemGearCheck struct {
	Purpose string `json:"purpose"`
	Kind    string `json:"kind"`
	Ability string `json:"ability"`
	Skill   string `json:"skill"`
	DC      int    `json:"dc"`
	DCRule  string `json:"dc_rule"`
}

// ItemGearDamage is either dice or a flat amount; exactly one form is set.
type ItemGearDamage struct {
	Purpose string `json:"purpose"`
	Dice    int    `json:"dice"`
	Sides   int    `json:"sides"`
	Flat    int    `json:"flat"`
	Type    string `json:"type"`
}

// ItemGearMeasure is a printed quantity as a fraction of its unit. A zero
// numerator is a printed zero, such as a Speed reduced to 0.
type ItemGearMeasure struct {
	Name        string `json:"name"`
	Numerator   int    `json:"numerator"`
	Denominator int    `json:"denominator"`
	Unit        string `json:"unit"`
}

// ItemGearContent is one entry of a pack's printed contents list.
type ItemGearContent struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

// ItemGearWeight is the printed weight as a fraction of a pound. Listed is false
// where the table prints a dash, which is not a weight of zero. Condition is
// "full" where the table prints the weight of a filled container.
type ItemGearWeight struct {
	Listed            bool   `json:"listed"`
	PoundsNumerator   int    `json:"pounds_numerator"`
	PoundsDenominator int    `json:"pounds_denominator"`
	Condition         string `json:"condition"`
}

var gearSourceKind = mundaneSourceKind{Category: "adventuring-gear", Section: "adventuring-gear", Payload: "gear"}

// decodeGearSourceBatch decodes one gear batch strictly against its manifest.
func decodeGearSourceBatch[D any](data []byte, m mundaneSourceManifest, gear func(*D) *ItemGear) (ItemSourceBatch[D], error) {
	return decodeMundaneSourceBatch(data, m, gearSourceKind, func(d *D, page int) (mundaneSourceRow, error) {
		p := gear(d)
		if err := checkItemGear(p, page); err != nil {
			return mundaneSourceRow{}, err
		}
		return mundaneSourceRow{Equip: "none", Use: itemGearUse(p), Keys: itemGearKeys(p)}, nil
	})
}

// checkItemGear checks that an entry is internally consistent. page is the page
// of the entry's heading; its text and tables may run two pages further.
func checkItemGear(p *ItemGear, page int) error {
	oneOf := func(value, allowed string) bool { return slices.Contains(strings.Split(allowed, "|"), value) }
	if !itemKeys(p.Actions, gearActions) || !itemKeys(p.SpellcastingFocus, gearClasses) {
		return fmt.Errorf("unknown or duplicate action or spellcasting class")
	}
	purposes := map[string]bool{}
	for _, c := range p.Checks {
		saving := c.Kind == "saving-throw"
		if !oneOf(c.Kind, gearCheckKinds) || !slices.Contains(abilityNames, c.Ability) || !oneOf(c.Skill, gearSkills) || saving != (c.Skill == "none") {
			return fmt.Errorf("invalid gear check")
		}
		if c.DC < 1 || !oneOf(c.DCRule, gearDCRules) || c.DCRule != "fixed" && (!saving || c.DC != 8) {
			return fmt.Errorf("invalid gear DC")
		}
		if purposes["check."+c.Purpose] {
			return fmt.Errorf("duplicate gear check")
		}
		purposes["check."+c.Purpose] = true
	}
	for _, d := range p.Damage {
		rolled := d.Flat == 0 && d.Dice > 0 && slices.Contains([]int{4, 6, 8, 10, 12}, d.Sides)
		flat := d.Flat > 0 && d.Dice == 0 && d.Sides == 0
		if !(rolled || flat) || !oneOf(d.Type, gearDamage) || purposes["damage."+d.Purpose] {
			return fmt.Errorf("invalid or duplicate gear damage")
		}
		purposes["damage."+d.Purpose] = true
	}
	for _, q := range p.Measures {
		if q.Denominator < 1 || purposes["measure."+q.Name] {
			return fmt.Errorf("invalid or duplicate gear measure")
		}
		purposes["measure."+q.Name] = true
	}
	for _, c := range p.Contents {
		if c.Quantity < 1 || !oneOf(c.Unit, gearUnits) || purposes["content."+c.Item] {
			return fmt.Errorf("invalid or duplicate pack content")
		}
		purposes["content."+c.Item] = true
	}
	if len(p.Traits) > 0 && !uniqueItemKeys(p.Traits) {
		return fmt.Errorf("duplicate gear trait")
	}
	if len(p.Measures)+len(p.Traits)+len(p.Contents) == 0 {
		return fmt.Errorf("a gear entry needs a measure, trait or content")
	}
	w := p.Weight
	if w.Listed != (w.PoundsNumerator > 0) || w.PoundsDenominator < 1 || !oneOf(w.Condition, "none|full") || !w.Listed && w.Condition != "none" {
		return fmt.Errorf("invalid gear weight")
	}
	if len(p.SourcePages) != 2 || p.SourcePages[0] != page || p.SourcePages[1] < page || p.SourcePages[1] > page+2 {
		return fmt.Errorf("gear source pages disagree with the presentation")
	}
	return nil
}

// itemGearUse is the one action the entry names, "multiple" where it names
// several, or "none" where it names no action.
func itemGearUse(p *ItemGear) string {
	switch len(p.Actions) {
	case 0:
		return "none"
	case 1:
		return p.Actions[0]
	}
	return "multiple"
}

// itemGearKeys restates an entry as the mechanic keys a resolver would need.
// They cover every payload value except weight and source_pages.
func itemGearKeys(p *ItemGear) []string {
	keys := []string{"item.gear"}
	for _, a := range p.Actions {
		keys = append(keys, "action."+a)
	}
	for _, class := range p.SpellcastingFocus {
		keys = append(keys, "spellcasting-focus."+class)
	}
	for _, c := range p.Checks {
		key := c.Kind + "." + c.Purpose + "-" + c.Ability
		if c.Skill != "none" {
			key += "-" + c.Skill
		}
		key += fmt.Sprintf("-dc-%d", c.DC)
		if c.DCRule != "fixed" {
			key += "-" + c.DCRule
		}
		keys = append(keys, key)
	}
	for _, d := range p.Damage {
		amount := fmt.Sprint(d.Flat)
		if d.Dice > 0 {
			amount = fmt.Sprintf("%dd%d", d.Dice, d.Sides)
		}
		keys = append(keys, "damage."+d.Purpose+"-"+amount+"-"+d.Type)
	}
	for _, q := range p.Measures {
		amount := fmt.Sprint(q.Numerator)
		if q.Denominator != 1 {
			amount += fmt.Sprintf("-over-%d", q.Denominator)
		}
		keys = append(keys, "measure."+q.Name+"-"+amount+"-"+q.Unit)
	}
	for _, t := range p.Traits {
		keys = append(keys, "trait."+t)
	}
	for _, c := range p.Contents {
		keys = append(keys, fmt.Sprintf("content.%s-%d-%s", c.Item, c.Quantity, c.Unit))
	}
	return keys
}
