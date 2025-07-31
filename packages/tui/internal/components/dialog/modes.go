package dialog

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/sst/opencode-sdk-go"
	"github.com/sst/opencode/internal/app"
	list "github.com/sst/opencode/internal/components/list"
	"github.com/sst/opencode/internal/components/modal"
	"github.com/sst/opencode/internal/layout"
	"github.com/sst/opencode/internal/styles"
	"github.com/sst/opencode/internal/theme"
	"github.com/sst/opencode/internal/util"
)

// ModeSelectedMsg is sent when a mode is selected
type ModeSelectedMsg struct {
	Mode opencode.Mode
}

// ModeDialog interface for the mode switching dialog
type ModeDialog interface {
	layout.Modal
}

type modeDialog struct {
	width  int
	height int

	modal         *modal.Modal
	list          list.List[list.Item]
	app           *app.App
	originalMode  *opencode.Mode
	modeApplied   bool
}

func (m *modeDialog) Init() tea.Cmd {
	return nil
}

func (m *modeDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, idx := m.list.GetSelectedItem(); idx >= 0 {
				if modeItem, ok := item.(modeItem); ok {
					selectedMode := modeItem.mode
					m.modeApplied = true
					return m, tea.Sequence(
						util.CmdHandler(modal.CloseModalMsg{}),
						util.CmdHandler(ModeSelectedMsg{Mode: selectedMode}),
					)
				}
			}
		}
	}

	_, prevIdx := m.list.GetSelectedItem()

	var cmd tea.Cmd
	listModel, cmd := m.list.Update(msg)
	m.list = listModel.(list.List[list.Item])

	if item, newIdx := m.list.GetSelectedItem(); newIdx >= 0 && newIdx != prevIdx {
		if modeItem, ok := item.(modeItem); ok {
			// Preview mode change while navigating
			m.app.Mode = &modeItem.mode
			return m, util.CmdHandler(ModeSelectedMsg{Mode: modeItem.mode})
		}
	}
	return m, cmd
}

func (m *modeDialog) Render(background string) string {
	return m.modal.Render(m.list.View(), background)
}

func (m *modeDialog) Close() tea.Cmd {
	if !m.modeApplied && m.originalMode != nil {
		m.app.Mode = m.originalMode
		return util.CmdHandler(ModeSelectedMsg{Mode: *m.originalMode})
	}
	return nil
}

// modeItem is a custom list item for mode selections
type modeItem struct {
	mode          opencode.Mode
	isCurrentMode bool
}

func (m modeItem) Render(
	selected bool,
	width int,
	baseStyle styles.Style,
) string {
	t := theme.CurrentTheme()

	itemStyle := baseStyle.
		Background(t.BackgroundPanel()).
		Foreground(t.Text())

	if selected {
		itemStyle = itemStyle.Foreground(t.Primary())
	}

	// Show only mode name - clean and simple
	text := m.mode.Name
	
	// Add current mode indicator
	if m.isCurrentMode {
		text = "● " + text
	} else {
		text = "  " + text
	}

	return itemStyle.
		PaddingLeft(1).
		Render(text)
}

func (m modeItem) Selectable() bool {
	return true
}

// NewModeDialog creates a new mode switching dialog
func NewModeDialog(app *app.App) ModeDialog {
	modes := app.Modes
	currentMode := app.Mode

	var selectedIdx int
	for i, mode := range modes {
		if currentMode != nil && mode.Name == currentMode.Name {
			selectedIdx = i
		}
	}

	// Convert modes to list items
	items := make([]list.Item, len(modes))
	for i, mode := range modes {
		isCurrentMode := currentMode != nil && mode.Name == currentMode.Name
		items[i] = modeItem{
			mode:          mode,
			isCurrentMode: isCurrentMode,
		}
	}

	listComponent := list.NewListComponent(
		list.WithItems(items),
		list.WithMaxVisibleHeight[list.Item](10),
		list.WithFallbackMessage[list.Item]("No modes available"),
		list.WithAlphaNumericKeys[list.Item](true),
		list.WithRenderFunc(func(item list.Item, selected bool, width int, baseStyle styles.Style) string {
			return item.Render(selected, width, baseStyle)
		}),
		list.WithSelectableFunc(func(item list.Item) bool {
			return item.Selectable()
		}),
	)

	// Set the initial selection to the current mode
	listComponent.SetSelectedIndex(selectedIdx)

	// Set the max width for the list to match the modal width
	listComponent.SetMaxWidth(36) // 40 (modal max width) - 4 (modal padding)
	return &modeDialog{
		list:         listComponent,
		modal:        modal.New(modal.WithTitle("Select Mode"), modal.WithMaxWidth(40)),
		app:          app,
		originalMode: currentMode,
		modeApplied:  false,
	}
}