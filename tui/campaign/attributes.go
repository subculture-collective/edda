package campaign

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

type attributeFormResult struct {
	Genre        string
	SettingStyle string
	Tone         string
}

// AttributesModel is a huh.Form-backed sub-model for selecting campaign
// genre, setting style, and tone via guided dropdowns.
type AttributesModel struct {
	form      *huh.Form
	result    attributeFormResult
	errorText string
	errorView viewport.Model
	width     int
	height    int
}

// NewAttributesModel builds the three-step attribute selection form.
func NewAttributesModel() AttributesModel {
	return NewAttributesModelWithState("", "", "", "")
}

// NewAttributesModelWithValues builds the attribute form with preselected
// values. It is used when proposal generation fails or the player navigates
// back from the proposal picker so the previous selections are preserved.
func NewAttributesModelWithValues(genre, settingStyle, tone string) AttributesModel {
	return NewAttributesModelWithState(genre, settingStyle, tone, "")
}

// NewAttributesModelWithState builds the attribute form with preselected values
// and an inline error message displayed above the form body when non-empty.
func NewAttributesModelWithState(genre, settingStyle, tone, errorText string) AttributesModel {
	result := attributeFormResult{
		Genre:        genre,
		SettingStyle: settingStyle,
		Tone:         tone,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Genre").
				Options(
					huh.NewOption("Fantasy", "Fantasy"),
					huh.NewOption("Sci-Fi", "Sci-Fi"),
					huh.NewOption("Horror", "Horror"),
					huh.NewOption("Historical", "Historical"),
					huh.NewOption("Modern", "Modern"),
					huh.NewOption("Post-Apocalyptic", "Post-Apocalyptic"),
					huh.NewOption("Steampunk", "Steampunk"),
				).
				Value(&result.Genre),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setting Style").
				Options(
					huh.NewOption("Open Wilderness", "Open Wilderness"),
					huh.NewOption("Urban Sprawl", "Urban Sprawl"),
					huh.NewOption("Island Archipelago", "Island Archipelago"),
					huh.NewOption("Underground", "Underground"),
					huh.NewOption("War-Torn Kingdom", "War-Torn Kingdom"),
					huh.NewOption("Peaceful Realm", "Peaceful Realm"),
				).
				Value(&result.SettingStyle),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Tone").
				Options(
					huh.NewOption("Gritty and Dark", "Gritty and Dark"),
					huh.NewOption("Light-Hearted", "Light-Hearted"),
					huh.NewOption("Epic and Grand", "Epic and Grand"),
					huh.NewOption("Mysterious", "Mysterious"),
					huh.NewOption("Humorous", "Humorous"),
					huh.NewOption("Tense and Suspenseful", "Tense and Suspenseful"),
				).
				Value(&result.Tone),
		),
	).WithShowHelp(false)

	return AttributesModel{
		form:      form,
		result:    result,
		errorText: errorText,
		errorView: viewport.New(0, 0),
	}
}

// SetSize updates the available rendering area.
func (m *AttributesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.layoutError()
}

// layoutError wraps the error to the box and gives it the rows the form leaves.
// A longer error scrolls, so the form stays on screen.
func (m *AttributesModel) layoutError() {
	if m.errorText == "" || m.width <= 0 || m.height <= 0 {
		return
	}
	frame := styles.Container
	width := max(m.width-frame.GetHorizontalFrameSize(), 1)
	content := lipgloss.NewStyle().Width(width).Render(m.errorMessage())
	// The title, the line under it and the line under the error take three rows.
	rows := max(m.height-frame.GetVerticalFrameSize()-3-lipgloss.Height(m.form.View()), 1)

	m.errorView.Width = width
	m.errorView.Height = min(lipgloss.Height(content), rows)
	m.errorView.SetContent(content)
	m.errorView.SetYOffset(m.errorView.YOffset) // a taller box may need a smaller offset
}

func (m AttributesModel) errorMessage() string {
	return styles.StatusError.Render("Error: " + m.errorText)
}

// errorScrolls reports whether the error needs more rows than it has.
func (m AttributesModel) errorScrolls() bool {
	return m.errorView.Height > 0 && m.errorView.TotalLineCount() > m.errorView.Height
}

// Init implements tea.Model.
func (m AttributesModel) Init() tea.Cmd {
	return m.form.Init()
}

// Update implements tea.Model.
func (m AttributesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// The form does not use these, so they scroll a long error.
	if m.errorScrolls() {
		switch msg := msg.(type) {
		case tea.MouseMsg:
			var cmd tea.Cmd
			m.errorView, cmd = m.errorView.Update(msg)
			return m, cmd
		case tea.KeyMsg:
			if msg.Type == tea.KeyPgUp || msg.Type == tea.KeyPgDown {
				var cmd tea.Cmd
				m.errorView, cmd = m.errorView.Update(msg)
				return m, cmd
			}
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	m.layoutError() // each form step has its own height

	switch m.form.State {
	case huh.StateCompleted:
		return m, func() tea.Msg {
			return AttributesReadyMsg{
				Genre:        m.result.Genre,
				SettingStyle: m.result.SettingStyle,
				Tone:         m.result.Tone,
			}
		}
	case huh.StateAborted:
		return m, func() tea.Msg { return BackMsg{} }
	}

	return m, cmd
}

// View implements tea.Model.
func (m AttributesModel) View() string {
	title := styles.Header.Render("✦ Campaign Attributes")
	body := m.form.View()
	content := styles.JoinVertical(title, "")
	if m.errorText != "" {
		message, spacer := m.errorMessage(), ""
		if m.width > 0 && m.height > 0 {
			message = m.errorView.View()
			if m.errorScrolls() {
				hint := fmt.Sprintf("PgUp/PgDn: scroll the error (%.0f%%)", m.errorView.ScrollPercent()*100)
				spacer = styles.Muted.Render(ansi.Truncate(hint, m.errorView.Width, "…"))
			}
		}
		content = styles.JoinVertical(content, message, spacer)
	}
	content = styles.JoinVertical(content, body)
	return styles.Frame(styles.Container, m.width, m.height).
		Render(content)
}
