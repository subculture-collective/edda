// Package character provides the character sheet view for the TUI. It shows
// the saved character; an SRD 5.2.1 character's values come from the sheet the
// server derives.
package character

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/progression"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
	"git.subcult.tv/subculture-collective/edda/tui/msgs"
	"git.subcult.tv/subculture-collective/edda/tui/scrollbox"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// Model is the Bubble Tea model for the character sheet view.
type Model struct {
	width, height int
	box           scrollbox.Model

	character *game.CharacterRecord
	loaded    bool
	notice    string
}

// New returns a character Model that has nothing to show yet.
func New() Model {
	return Model{}
}

// Show replaces the sheet with the saved character. Nil means the campaign has
// none.
func (m *Model) Show(character *game.CharacterRecord) {
	m.character, m.loaded, m.notice = character, true, ""
}

// Notify says why the sheet is missing or may be out of date. A character that
// is already shown stays.
func (m *Model) Notify(notice string) { m.notice = notice }

// SetSize updates the viewport dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd { return nil }

// Update implements tea.Model. Pressing Escape emits NavigateBackMsg so the
// parent App can return focus to the narrative view. Other messages scroll a
// sheet that does not fit.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEsc {
		return m, func() tea.Msg { return msgs.NavigateBackMsg{} }
	}
	return m, m.box.Update(msg, m.content(), m.width, m.height)
}

// View implements tea.Model and renders the character sheet.
func (m Model) View() string {
	return m.box.View(m.content(), m.width, m.height)
}

// content renders everything inside the box.
func (m Model) content() string {
	lines := []string{styles.Header.Render("⚔️  Character Sheet"), ""}
	if m.notice != "" {
		lines = append(lines, styles.SystemMessage.Render(m.notice))
	}
	switch {
	case !m.loaded:
	case m.character == nil:
		lines = append(lines, styles.Muted.Render("No character yet."))
	default:
		lines = append(lines, sheetLines(m.character)...)
	}
	lines = append(lines, "", styles.Muted.Render("Press Esc to return to the narrative view."))
	return styles.JoinVertical(lines...)
}

func sheetLines(c *game.CharacterRecord) []string {
	lines := []string{styles.SubHeader.Render(c.Name)}
	if c.Description != "" {
		lines = append(lines, styles.Body.Render(c.Description))
	}
	nextXP := progression.NextLevelThreshold(c.Level)
	if c.Sheet != nil {
		nextXP = c.Sheet.Progression.NextLevelXP
	}
	vitals := fmt.Sprintf("Level %d · %s · HP %d/%d · XP %d", c.Level, words(c.Status), c.HP, c.MaxHP, c.Experience)
	// A rules character without a sheet has no known XP threshold.
	if c.Sheet != nil || c.RulesVersion != rulesVersionSRD {
		vitals += fmt.Sprintf("/%d", nextXP)
	}
	lines = append(lines, styles.Body.Render(vitals))
	if c.Location != "" {
		lines = append(lines, styles.Body.Render("Location: "+c.Location))
	}

	switch {
	case c.Sheet != nil:
		return append(lines, rulesLines(c.Sheet)...)
	case errors.Is(c.SheetError, game.ErrSheetLoadFailed):
		return append(lines, "", styles.StatusWarning.Render("The character sheet could not be loaded. Open this tab again to retry; the Logs tab says why."))
	case c.SheetError != nil:
		// As on the web: it says what is wrong with the saved build.
		return append(lines, "", styles.StatusWarning.Render("Rules sheet unavailable: "+c.SheetError.Error()))
	}
	return append(lines, freeformLines(c)...)
}

const rulesVersionSRD = string(domain.RulesVersionSRD521)

// freeformLines are the stats and abilities of a character without SRD rules.
func freeformLines(c *game.CharacterRecord) []string {
	var lines []string
	var stats map[string]any
	// Unreadable stats or abilities show as none, as they do in the web app.
	if json.Unmarshal(c.Stats, &stats) == nil && len(stats) > 0 {
		lines = section(lines, "Stats")
		for _, key := range slices.Sorted(maps.Keys(stats)) {
			lines = append(lines, styles.Body.Render("  "+words(key)+": "+statValue(stats[key])))
		}
	}
	var abilities []api.CharacterAbility
	if json.Unmarshal(c.Abilities, &abilities) == nil && len(abilities) > 0 {
		lines = section(lines, "Abilities")
		for _, ability := range abilities {
			line := "  ◆ " + ability.Name
			if ability.Description != "" {
				line += ": " + ability.Description
			}
			lines = append(lines, styles.Body.Render(line))
		}
	}
	return lines
}

var abilityOrder = []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"}

// rulesLines shows the server-derived values; it computes none of them.
func rulesLines(sheet *srd521.CharacterSheet) []string {
	var lines []string
	add := func(format string, args ...any) {
		lines = append(lines, styles.Body.Render("  "+fmt.Sprintf(format, args...)))
	}

	defenses := make([]string, 0, len(sheet.Defenses)+1)
	for _, name := range slices.Sorted(maps.Keys(sheet.Defenses)) {
		format := "%s %d"
		if name == "initiative" {
			format = "%s %+d"
		}
		defenses = append(defenses, fmt.Sprintf(format, words(name), sheet.Defenses[name].Value))
	}
	defenses = append(defenses, fmt.Sprintf("Proficiency %+d", sheet.Proficiency.Value))
	lines = append(lines, styles.Body.Render(strings.Join(defenses, " · ")))

	if saves := sheet.DeathSaves; saves != nil {
		switch {
		case saves.Dead:
			lines = section(lines, "Dead")
		case saves.Stable:
			lines = section(lines, "At 0 HP: Unconscious")
			add("Stable at 0 HP; still Unconscious.")
		default:
			lines = section(lines, "At 0 HP: Unconscious")
			add("Death Saving Throws: %d successes, %d failures.", saves.Successes, saves.Failures)
		}
	}
	if fall := sheet.Fall; fall != nil {
		lines = section(lines, "Declared fall")
		if fall.HeightFeet == nil {
			add("A declared fall is waiting. /fall shows why it requires adjudication.")
		} else {
			landing := "solid ground"
			if fall.Landing == "water" {
				landing = "water"
			}
			add("A fall of %d feet onto %s is waiting. Preview it with /fall.", *fall.HeightFeet, landing)
		}
	}
	if sheet.Prone {
		lines = section(lines, "Conditions")
		add("Prone")
	}

	lines = section(lines, "Abilities and saves")
	for _, name := range abilityOrder {
		ability := sheet.Abilities[name]
		add("%s %d (%+d), save %+d", words(name), ability.Score, ability.Modifier.Value, sheet.Saves[name].Value)
	}
	lines = section(lines, "Skills")
	skills := make([]string, 0, len(sheet.Skills))
	for _, name := range slices.Sorted(maps.Keys(sheet.Skills)) {
		skills = append(skills, fmt.Sprintf("%s %+d", words(name), sheet.Skills[name].Value))
	}
	add("%s", strings.Join(skills, " · "))

	if len(sheet.Resources) > 0 {
		lines = section(lines, "Resources")
		for _, resource := range sheet.Resources {
			add("%s %d/%d", resource.Label, resource.Current, resource.Maximum)
		}
	}
	if len(sheet.Capabilities) > 0 {
		lines = section(lines, "Capabilities")
		for _, capability := range sheet.Capabilities {
			support := "read only"
			switch {
			case capability.ID == "tactical-mind":
				support = "saved failed-check window only"
			case capability.Playable:
				support = "supported"
			}
			add("%s (%s)", capability.Label, support)
		}
	}
	if casting := sheet.Spellcasting; casting != nil {
		lines = section(lines, "Spellcasting")
		add("Casting ability: %s (%+d)", words(casting.Ability), casting.Modifier.Value)
		light := false
		for _, spell := range casting.Spells {
			known, prepared, use := "Unknown", "Not prepared", "Effect unavailable"
			if spell.Known {
				known = "Known"
			}
			if spell.Prepared {
				prepared = "Prepared"
			}
			if spell.Playable {
				use = "Supported"
				light = light || spell.ID == "light"
			}
			add("%s · %s · %s · %s", spell.Label, known, prepared, use)
		}
		if c := sheet.Concentration; c != nil {
			add("Concentration: %s, expires %s", c.SpellLabel, game.FormatCampaignTime(c.ExpiresAt))
		} else {
			add("No concentration effect active.")
		}
		switch active := sheet.Light; {
		case active != nil:
			item := sheet.LightItem
			if item == "" {
				item = "an object you no longer carry"
			}
			covered := ""
			if active.Covered {
				covered = " Covered: it sheds no light until uncovered."
			}
			add("Light on %s (%s), expires %s.%s", item, active.Color, game.FormatCampaignTime(active.ExpiresAt), covered)
		case light:
			add("No Light active.")
		}
	}
	if len(sheet.Doors) > 0 {
		lines = section(lines, "Doors here")
		for _, door := range sheet.Doors {
			if door.State == "" {
				add("%s", door.Name)
			} else {
				add("%s (%s)", door.Name, door.State)
			}
		}
	}
	if len(sheet.TravelPaces) > 0 {
		lines = section(lines, "Travel")
		add("A known route with declared facts leads from here.")
	}
	return append(lines, "", styles.Muted.Render("SRD 5.2.1 support is partial. Edda does not claim Full Rules support."))
}

// section appends a blank line and a heading.
func section(lines []string, title string) []string {
	return append(lines, "", styles.SubHeader.Render(title))
}

// words turns a saved key such as "armor-class" into "Armor class".
func words(key string) string {
	text := []rune(strings.NewReplacer("-", " ", "_", " ").Replace(key))
	if len(text) > 0 {
		text[0] = unicode.ToUpper(text[0])
	}
	return string(text)
}

// statValue is a saved free-form stat as text.
func statValue(value any) string {
	switch value := value.(type) {
	case string:
		if value != "" {
			return value
		}
	case nil:
	case []any:
		parts := make([]string, len(value))
		for i, part := range value {
			parts[i] = statValue(part)
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ")
		}
	default:
		if data, err := json.Marshal(value); err == nil {
			return string(data)
		}
	}
	return "—"
}
