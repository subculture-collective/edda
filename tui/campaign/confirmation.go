package campaign

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"git.subcult.tv/subculture-collective/edda/internal/world"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// confirmationFooterRows are the line above the buttons and the buttons.
const confirmationFooterRows = 2

// ConfirmationModel shows campaign + character details for final review.
type ConfirmationModel struct {
	name             string
	summary          string
	profile          world.CampaignProfile
	characterProfile world.CharacterProfile
	selected         int // 0 = Begin Adventure, 1 = Go Back
	body             viewport.Model
	width            int
	height           int
}

// NewConfirmationModel builds the confirmation review model.
func NewConfirmationModel(name, summary string, profile world.CampaignProfile, charProfile world.CharacterProfile) ConfirmationModel {
	return ConfirmationModel{
		body:             viewport.New(0, 0),
		name:             name,
		summary:          summary,
		profile:          profile,
		characterProfile: charProfile,
	}
}

// SetSize updates the layout dimensions.
func (m *ConfirmationModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.layoutBody()
}

// layoutBody wraps the details to the box and gives them the rows the buttons
// leave. The details scroll when they need more.
func (m *ConfirmationModel) layoutBody() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	frame := styles.Container
	width := max(m.width-frame.GetHorizontalFrameSize(), 1)
	content := lipgloss.NewStyle().Width(width).Render(m.details())
	rows := max(m.height-frame.GetVerticalFrameSize()-confirmationFooterRows, 1)

	m.body.Width = width
	m.body.Height = min(lipgloss.Height(content), rows)
	m.body.SetContent(content)
	m.body.SetYOffset(m.body.YOffset) // a taller box may need a smaller offset
}

// scrolls reports whether the details need more rows than the box has.
func (m ConfirmationModel) scrolls() bool {
	return m.body.Height > 0 && m.body.TotalLineCount() > m.body.Height
}

// Init implements tea.Model.
func (m ConfirmationModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		var cmd tea.Cmd
		m.body, cmd = m.body.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			m.selected = 0
			return m, nil
		case tea.KeyRight:
			m.selected = 1
			return m, nil
		case tea.KeyEnter:
			if m.selected == 0 {
				return m, func() tea.Msg {
					return ConfirmedMsg{
						Name:             m.name,
						Summary:          m.summary,
						Profile:          m.profile,
						CharacterProfile: m.characterProfile,
					}
				}
			}
			return m, func() tea.Msg { return ChangeMsg{} }
		case tea.KeyEsc:
			return m, func() tea.Msg { return ChangeMsg{} }
		}

		// Up and down scroll long details, and otherwise pick a button.
		if m.scrolls() {
			var cmd tea.Cmd
			m.body, cmd = m.body.Update(verticalKey(msg))
			return m, cmd
		}
		switch msg.String() {
		case "up", "k":
			m.selected = 0
		case "down", "j":
			m.selected = 1
		}
	}
	return m, nil
}

// verticalKey passes on the keys that scroll up and down and drops the rest of
// the viewport's default bindings.
func verticalKey(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "up", "k", "down", "j", "pgup", "pgdown":
		return msg
	default:
		return nil
	}
}

// details renders everything above the buttons.
func (m ConfirmationModel) details() string {
	var b strings.Builder

	b.WriteString(styles.Header.Render("✦ Campaign Summary"))
	b.WriteString("\n\n")
	b.WriteString(styles.Body.Render(fmt.Sprintf("Name: %s", m.name)))
	b.WriteString("\n")
	b.WriteString(styles.Body.Render(fmt.Sprintf("Summary: %s", m.summary)))
	b.WriteString("\n")
	b.WriteString(styles.Body.Render(
		fmt.Sprintf("Genre: %s · Tone: %s", m.profile.Genre, m.profile.Tone),
	))
	b.WriteString("\n")
	if len(m.profile.Themes) > 0 {
		b.WriteString(styles.Body.Render(
			fmt.Sprintf("Themes: %s", strings.Join(m.profile.Themes, ", ")),
		))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.Header.Render("✦ Your Character"))
	b.WriteString("\n\n")
	b.WriteString(styles.Body.Render(fmt.Sprintf("Name: %s", m.characterProfile.Name)))
	b.WriteString("\n")
	b.WriteString(styles.Body.Render(fmt.Sprintf("Concept: %s", m.characterProfile.Concept)))
	b.WriteString("\n")
	b.WriteString(styles.Body.Render(fmt.Sprintf("Background: %s", m.characterProfile.Background)))
	return b.String()
}

// View implements tea.Model.
func (m ConfirmationModel) View() string {
	body := m.details()
	spacer := ""
	if m.width > 0 && m.height > 0 {
		body = m.body.View()
		if m.scrolls() {
			hint := fmt.Sprintf("↑/↓ PgUp/PgDn: scroll (%.0f%%) · ←/→: choose", m.body.ScrollPercent()*100)
			spacer = styles.Muted.Render(ansi.Truncate(hint, m.body.Width, "…"))
		}
	}

	// Buttons
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorAccent)
	inactiveStyle := styles.Muted

	beginLabel := "  Begin Adventure  "
	backLabel := "  Go Back  "

	var beginBtn, backBtn string
	if m.selected == 0 {
		beginBtn = activeStyle.Render("[▸" + beginLabel + "]")
		backBtn = inactiveStyle.Render("[ " + backLabel + "]")
	} else {
		beginBtn = inactiveStyle.Render("[ " + beginLabel + "]")
		backBtn = activeStyle.Render("[▸" + backLabel + "]")
	}

	buttons := styles.JoinHorizontal(beginBtn, "    ", backBtn)

	return styles.Frame(styles.Container, m.width, m.height).
		Render(body + "\n" + spacer + "\n" + buttons)
}
