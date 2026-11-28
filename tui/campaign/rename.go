package campaign

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// guidedFormDefaults is the description the guided character form saved while
// it ignored the player's choices.
const guidedFormDefaults = "human barbarian. Acolyte"

// CharacterRenamedMsg carries the name chosen for a character that had none.
type CharacterRenamedMsg struct{ Name string }

// RenameModel asks the player to name a saved character that has none. It emits
// CharacterRenamedMsg, or BackMsg on escape.
type RenameModel struct {
	input       textinput.Model
	description string
	err         error
	width       int
	height      int
}

// NewRenameModel builds the prompt. description is the saved character's, which
// tells the player which character this is.
func NewRenameModel(description string) RenameModel {
	input := textinput.New()
	input.Placeholder = "Character name"
	input.Focus()
	return RenameModel{input: input, description: description}
}

// SetSize updates the available rendering area.
func (m *RenameModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = styles.InputWidth(m.textWidth(), m.input.Prompt)
	m.input.SetCursor(m.input.Position()) // rescroll typed text to the new width
}

// textWidth is the columns inside the box, or 0 before the first resize.
func (m RenameModel) textWidth() int {
	if m.width <= 0 {
		return 0
	}
	return max(m.width-styles.Container.GetHorizontalFrameSize(), 1)
}

// Init implements tea.Model.
func (m RenameModel) Init() tea.Cmd { return textinput.Blink }

// Update implements tea.Model.
func (m RenameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyEsc:
			return m, func() tea.Msg { return BackMsg{} }
		case tea.KeyEnter:
			name := strings.TrimSpace(m.input.Value())
			if m.err = domain.ValidateCharacterName(name); m.err != nil {
				return m, nil
			}
			return m, func() tea.Msg { return CharacterRenamedMsg{Name: name} }
		}
		m.err = nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m RenameModel) View() string {
	width := m.textWidth()
	wrap := func(style lipgloss.Style, text string) string {
		if width > 0 {
			style = style.Width(width)
		}
		return style.Render(text)
	}
	about := "This saved character has no name."
	if m.description == guidedFormDefaults {
		about += " The guided form that saved it also replaced your race, class, background and alignment with Human, Barbarian, Acolyte and Lawful Good. Naming does not restore your choices."
	}
	saved := strings.Join(strings.Fields(m.description), " ")
	if width > 0 {
		saved = ansi.Truncate("Saved as: "+saved, width, "…")
	} else {
		saved = "Saved as: " + saved
	}

	blocks := []string{
		wrap(styles.Header, "✦ Name Your Character"),
		"",
		wrap(styles.Body, about),
		styles.Muted.Render(saved),
		"",
		m.input.View(),
	}
	if m.err != nil {
		blocks = append(blocks, wrap(styles.StatusError, m.err.Error()))
	}
	blocks = append(blocks, wrap(styles.Muted, "Enter: save · Esc: campaigns"))
	return styles.Frame(styles.Container, m.width, m.height).
		Render(styles.JoinVertical(blocks...))
}
