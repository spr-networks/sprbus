package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tidwall/gjson"

	"fmt"
	"strings"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
		//Foreground(lipgloss.Color("129")).
		Padding(0, 0)

	titleStyleViewport = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	statusMessageStyle = lipgloss.NewStyle().
		//Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).MarginLeft(2).
		Foreground(lipgloss.Color("#3C3C3C")).
		Render
)

const useHighPerformanceRenderer = false

type item struct {
	title       string
	description string
	json        string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) Json() string        { return i.json }
func (i item) FilterValue() string { return i.title + i.description }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	title        string
	content      string
	viewport     viewport.Model
	ready        bool
}

func NewModel(items []list.Item) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	eventList := list.New(items, delegate, 0, 0)
	eventList.Title = "sprbus event viewer"
	eventList.Styles.Title = titleStyle
	eventList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         eventList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		title:        "",
		content:      "",
		ready:        false,
	}
}

func (m model) Init() tea.Cmd {
	//return checkServer
	//return tea.Batch(downloadAndInstall(m.packages[m.index]), m.spinner.Tick)
	//NOTE do we need to start HandleEvent here?
	return tea.EnterAltScreen
}

type EventMsg struct {
	title       string
	description string
	json        string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent("viwport content is heah")
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case EventMsg:
		newItem := item{title: msg.title, description: msg.description, json: msg.json}
		insCmd := m.list.InsertItem(0, newItem)
		//statusCmd := m.list.NewStatusMessage(statusMessageStyle("Event: " + newItem.Title()))
		//return m, tea.Batch(insCmd, statusCmd)
		return m, insCmd

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case msg.String() == "enter":
			//set viewport content
			i := m.list.SelectedItem()
			i2, ok := i.(item)
			if ok {
				m.title = fmt.Sprintf("%s %s", i2.Title(), i2.Description())
				jsonMarkdown := fmt.Sprintf("```json\n%s\n```", gjson.Get(i2.Json(), "@pretty").String())
				jsonFormatted, _ := glamour.Render(jsonMarkdown, "dark")
				m.content = jsonFormatted
				m.viewport.SetContent(m.content)
			}

			//cmd := m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("select:%v", idx)))

			return m, nil

		case msg.String() == "q":
			if m.title != "" {
				m.title = ""
				m.content = ""
				//cmd := m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("close")))

				return m, nil
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			newItem := item{title: "fish", description: "gosi"} //m.itemGenerator.next()
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))

			return m, tea.Batch(insCmd, statusCmd)

		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// show item
	if m.title != "" {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}

	return appStyle.Render(m.list.View())
}

func (m model) headerView() string {
	title := titleStyleViewport.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
