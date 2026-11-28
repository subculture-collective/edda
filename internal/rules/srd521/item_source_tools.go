package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// The tool entries (SRD 5.2.1 pp. 93-94) payload shared by the tool catalogue
// batches. A record describes one printed entry. It is not read by any check,
// crafting or proficiency resolver.

// ItemTool is one tool entry. Tasks and craft targets are keys that follow the
// printed wording; they do not link to other records.
type ItemTool struct {
	ToolGroup   string            `json:"tool_group"`
	Ability     string            `json:"ability"`
	Utilize     []ItemToolUtilize `json:"utilize"`
	Craft       []ItemToolCraft   `json:"craft" srd:"may-be-empty"`
	Weight      ItemToolWeight    `json:"weight"`
	SourcePages []int             `json:"source_pages"`
}

// ItemToolUtilize is one thing the Utilize line lists, with its DC.
type ItemToolUtilize struct {
	Task string `json:"task"`
	DC   int    `json:"dc"`
}

// ItemToolCraft is one entry of the Craft line, with the printed exceptions to a
// group such as "Medium armor (except Hide)".
type ItemToolCraft struct {
	Item   string   `json:"item"`
	Except []string `json:"except" srd:"may-be-empty"`
}

// ItemToolWeight is the printed weight in pounds. Listed is false where the
// entry prints a dash, which is not a weight of zero.
type ItemToolWeight struct {
	Listed bool `json:"listed"`
	Pounds int  `json:"pounds"`
}

var toolSourceKind = mundaneSourceKind{Category: "tool", Section: "tools", Use: "utilize", Payload: "tool"}

// decodeToolSourceBatch decodes one tool batch strictly against its manifest.
func decodeToolSourceBatch[D any](data []byte, m mundaneSourceManifest, tool func(*D) *ItemTool) (ItemSourceBatch[D], error) {
	return decodeMundaneSourceBatch(data, m, toolSourceKind, func(d *D, page int) (mundaneSourceRow, error) {
		p := tool(d)
		if err := checkItemTool(p, page); err != nil {
			return mundaneSourceRow{}, err
		}
		return mundaneSourceRow{Equip: "none", Keys: itemToolKeys(p)}, nil
	})
}

// checkItemTool checks that an entry is internally consistent. page is the page
// of the entry's heading; an entry may continue onto the next page.
func checkItemTool(p *ItemTool, page int) error {
	if p.ToolGroup != "artisans-tools" && p.ToolGroup != "other-tools" {
		return fmt.Errorf("unknown tool group")
	}
	if !slices.Contains(abilityNames, p.Ability) {
		return fmt.Errorf("unknown tool ability")
	}
	tasks := make([]string, 0, len(p.Utilize))
	for _, u := range p.Utilize {
		if u.DC < 1 {
			return fmt.Errorf("a Utilize task needs a positive DC")
		}
		tasks = append(tasks, u.Task)
	}
	if !uniqueItemKeys(tasks) {
		return fmt.Errorf("missing or duplicate Utilize task")
	}
	crafted := map[string]bool{}
	for _, c := range p.Craft {
		if crafted[c.Item] || len(c.Except) > 0 && !uniqueItemKeys(c.Except) {
			return fmt.Errorf("duplicate craft entry or exception")
		}
		crafted[c.Item] = true
	}
	if p.Weight.Listed != (p.Weight.Pounds > 0) {
		return fmt.Errorf("invalid tool weight")
	}
	if len(p.SourcePages) != 2 || p.SourcePages[0] != page || p.SourcePages[1] != page && p.SourcePages[1] != page+1 {
		return fmt.Errorf("tool source pages disagree with the presentation")
	}
	return nil
}

// itemToolKeys restates an entry as the mechanic keys a resolver would need.
func itemToolKeys(p *ItemTool) []string {
	keys := []string{"tool.utilize", "tool.proficiency", "tool." + p.ToolGroup, "ability." + p.Ability}
	for _, u := range p.Utilize {
		keys = append(keys, fmt.Sprintf("utilize.%s-dc-%d", u.Task, u.DC))
	}
	for _, c := range p.Craft {
		key := "craft." + c.Item
		if len(c.Except) > 0 {
			key += "-except-" + strings.Join(c.Except, "-")
		}
		keys = append(keys, key)
	}
	return keys
}
