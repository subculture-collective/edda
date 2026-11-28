// Package quest provides the quest log view for the TUI. It shows the
// campaign's quests and their objectives.
package quest

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/tui/msgs"
	"git.subcult.tv/subculture-collective/edda/tui/scrollbox"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// Model is the Bubble Tea model for the quest log view.
type Model struct {
	width, height int
	box           scrollbox.Model

	quests []game.QuestRecord
	loaded bool
	notice string
}

// New returns a quest Model that has nothing to show yet.
func New() Model {
	return Model{}
}

// Show replaces the log with quests. Active quests come first; otherwise the
// given order stays.
func (m *Model) Show(quests []game.QuestRecord) {
	quests = slices.Clone(quests)
	slices.SortStableFunc(quests, func(a, b game.QuestRecord) int { return finished(a) - finished(b) })
	m.quests, m.loaded, m.notice = quests, true, ""
}

func finished(q game.QuestRecord) int {
	if q.Status == domain.QuestStatusActive {
		return 0
	}
	return 1
}

// Notify says why the log is missing or may be out of date. Quests that are
// already shown stay.
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
// log that does not fit.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEsc {
		return m, func() tea.Msg { return msgs.NavigateBackMsg{} }
	}
	return m, m.box.Update(msg, m.content(), m.width, m.height)
}

// View implements tea.Model and renders the quest log.
func (m Model) View() string {
	return m.box.View(m.content(), m.width, m.height)
}

// content renders everything inside the box.
func (m Model) content() string {
	lines := []string{styles.Header.Render("📜 Quest Log"), ""}
	if m.notice != "" {
		lines = append(lines, styles.SystemMessage.Render(m.notice))
	}
	switch {
	case !m.loaded:
	case len(m.quests) == 0:
		lines = append(lines, styles.Muted.Render("No quests yet."))
	default:
		sections := make([]string, len(m.quests))
		for i, q := range m.quests {
			sections[i] = renderQuest(q)
		}
		lines = append(lines, strings.Join(sections, "\n\n"))
	}
	return styles.JoinVertical(lines...)
}

func renderQuest(q game.QuestRecord) string {
	name := styles.SubHeader.Render("◆ " + q.Title)
	if q.Status != domain.QuestStatusActive {
		name += styles.Muted.Render("  (" + string(q.Status) + ")")
	}
	lines := []string{name}
	if q.Description != "" {
		lines = append(lines, styles.Body.Render("  "+q.Description))
	}
	for _, objective := range q.Objectives {
		lines = append(lines, renderObjective(objective, q.Status))
	}
	return styles.JoinVertical(lines...)
}

// renderObjective marks an objective done, still open, or lost with a failed
// quest.
func renderObjective(objective domain.QuestObjective, status domain.QuestStatus) string {
	switch {
	case objective.Completed:
		return styles.StatusSuccess.Render("  ✓ ") + styles.Muted.Render(objective.Description)
	case status == domain.QuestStatusFailed:
		return styles.StatusError.Render("  ✗ ") + styles.Muted.Render(objective.Description)
	default:
		return styles.Muted.Render("  ○ ") + styles.Body.Render(objective.Description)
	}
}
