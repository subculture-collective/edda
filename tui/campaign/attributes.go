package campaign

import (
	"fmt"
	"slices"
	"strings"

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
	result    *attributeFormResult // the form writes the choices here
	errorText string
	errorView viewport.Model
	gap       bool // blank line under the title
	spacer    bool // line under the error
	formRows  int  // rows the form has; 0 means all it needs
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
	result := &attributeFormResult{
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
		gap:       true,
		spacer:    errorText != "",
	}
}

// SetSize updates the available rendering area.
func (m *AttributesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.layout()
}

// layout shares the box's rows between the title, the error and the form. When
// they do not all fit, it drops the blank lines, then scrolls the error and the
// form's options.
func (m *AttributesModel) layout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	frame := styles.Container
	width := max(m.width-frame.GetHorizontalFrameSize(), 1)
	rows := m.height - frame.GetVerticalFrameSize()
	formHeight := lipgloss.Height(m.form.View()) // each form step has its own height

	var content string
	needed := 2 + formHeight // the title and the line under it
	if m.errorText != "" {
		content = lipgloss.NewStyle().Width(width).Render(m.errorMessage())
		needed += 2 // at least one line of error and the line under it
	}
	m.gap = needed <= rows
	used := 1 // the title
	if m.gap {
		used++
	}

	m.errorView.Width = width
	m.errorView.Height = 0
	m.spacer = false
	if m.errorText != "" {
		errorHeight := lipgloss.Height(content)
		m.errorView.Height = min(max(rows-used-1-formHeight, 1), errorHeight)
		// A short box keeps the line under the error only for the scroll hint.
		m.spacer = m.gap || m.errorView.Height < errorHeight
		m.errorView.SetContent(content)
		m.errorView.SetYOffset(m.errorView.YOffset) // a taller box may need a smaller offset
		used += m.errorView.Height
		if m.spacer {
			used++
		}
	}
	m.formRows = max(rows-used, 1)
}

// formView renders the form in at most formRows rows. When it needs more, the
// field title stays and the options scroll with the cursor. This cuts the
// rendered form because huh's Height moves the chosen option to the top and
// hides the ones above it.
func (m AttributesModel) formView() string {
	view := m.form.View()
	lines := strings.Split(view, "\n")
	if m.formRows <= 0 || len(lines) <= m.formRows {
		return view
	}
	// huh marks the option under the cursor with "> ".
	cursor := max(slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(ansi.Strip(line), "> ")
	}), 0)
	if m.formRows == 1 {
		return lines[cursor]
	}
	options, shown := lines[1:], m.formRows-1
	first := min(max(cursor-1-shown/2, 0), len(options)-shown)
	title := strings.TrimRight(lines[0], " ") + styles.Muted.Render("  ↑/↓: more")
	title = ansi.Truncate(title, m.errorView.Width, "…")
	return strings.Join(append([]string{title}, options[first:first+shown]...), "\n")
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
	m.layout()

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
	blocks := []string{styles.Header.Render("✦ Campaign Attributes")}
	if m.gap {
		blocks = append(blocks, "")
	}
	if m.errorText != "" {
		message, spacer := m.errorMessage(), ""
		if m.width > 0 && m.height > 0 {
			message = m.errorView.View()
			if m.errorScrolls() {
				hint := fmt.Sprintf("PgUp/PgDn: scroll the error (%.0f%%)", m.errorView.ScrollPercent()*100)
				spacer = styles.Muted.Render(ansi.Truncate(hint, m.errorView.Width, "…"))
			}
		}
		blocks = append(blocks, message)
		if m.spacer {
			blocks = append(blocks, spacer)
		}
	}
	blocks = append(blocks, m.formView())
	return styles.Frame(styles.Container, m.width, m.height).
		Render(styles.JoinVertical(blocks...))
}
