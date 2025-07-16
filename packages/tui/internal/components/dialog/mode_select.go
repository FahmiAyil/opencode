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
	"github.com/muesli/reflow/truncate"
)

// ModeSelectedMsg is sent when a mode is selected from the popup
type ModeSelectedMsg struct {
	Mode opencode.Mode
}

type modeSelectDialog struct {
	app   *app.App
	modal *modal.Modal
	list  list.List[list.Item]
}

type modeItem struct {
	mode     opencode.Mode
	isActive bool
}

func (m modeItem) Render(selected bool, width int, baseStyle styles.Style) string {
	t := theme.CurrentTheme()
	
	indicator := "  "
	if m.isActive {
		indicator = "● " // Current mode indicator
	}
	
	text := indicator + m.mode.Name
	truncatedText := truncate.StringWithTail(text, uint(width-1), "...")
	
	var itemStyle styles.Style
	if selected {
		itemStyle = baseStyle.
			Background(t.Primary()).
			Foreground(t.BackgroundElement()).
			Width(width).
			PaddingLeft(1)
	} else {
		if m.isActive {
			// Highlight current mode even when not selected
			itemStyle = baseStyle.
				Foreground(t.Primary()).
				PaddingLeft(1)
		} else {
			itemStyle = baseStyle.
				Foreground(t.TextMuted()).
				PaddingLeft(1)
		}
	}
	
	return itemStyle.Render(truncatedText)
}

func (m modeItem) Selectable() bool {
	return true
}

func NewModeSelectDialog(app *app.App) layout.Modal {
	modes := app.Modes
	items := make([]list.Item, len(modes))
	currentModeIndex := 0
	
	for i, mode := range modes {
		isActive := app.Mode.Name == mode.Name
		if isActive {
			currentModeIndex = i
		}
		items[i] = modeItem{
			mode:     mode,
			isActive: isActive,
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
	listComponent.SetSelectedIndex(currentModeIndex)
	
	// Set the max width for the list to match the modal width
	listComponent.SetMaxWidth(36) // 40 (modal max width) - 4 (modal padding)
	
	return &modeSelectDialog{
		app:   app,
		list:  listComponent,
		modal: modal.New(modal.WithTitle("Select Mode"), modal.WithMaxWidth(40)),
	}
}

func (m *modeSelectDialog) Init() tea.Cmd {
	return m.list.Init()
}

func (m *modeSelectDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selectedItem, idx := m.list.GetSelectedItem(); idx >= 0 {
				if modeItem, ok := selectedItem.(modeItem); ok {
					return m, func() tea.Msg {
						return ModeSelectedMsg{Mode: modeItem.mode}
					}
				}
			}
		case "esc", "q":
			return m, m.Close()
		}
	}
	
	// Pass all messages to the list for navigation handling
	var cmd tea.Cmd
	listModel, cmd := m.list.Update(msg)
	m.list = listModel.(list.List[list.Item])
	
	return m, cmd
}

func (m *modeSelectDialog) View() string {
	return m.list.View()
}

func (m *modeSelectDialog) Render(background string) string {
	return m.modal.Render(m.list.View(), background)
}

func (m *modeSelectDialog) Close() tea.Cmd {
	return util.CmdHandler(modal.CloseModalMsg{})
}