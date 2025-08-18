package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/junwei890/crawler/src"
)

var (
	focused = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#322C2B"))

	errored = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC3C22")).Bold(true)

	finished = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")).Bold(true)

	bindings = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9")).Italic(true)
)

type errorLogs string

type conversion struct{}

type doneMsg struct {
	err error
}

type model struct {
	uri      textinput.Model
	sites    textarea.Model
	spinner  spinner.Model
	view     viewport.Model
	stream   chan struct{}
	errors   chan string
	count    int
	logs     []string
	crawling bool
	done     bool
	err      error
	width    int
	height   int
}

func InitialModel() model {
	uri := textinput.New()
	uri.Placeholder = "MongoDB URI here"
	uri.Focus()
	uri.Width = 97

	sites := textarea.New()
	sites.Placeholder = "Sites here, each site should be on a newline"
	sites.SetWidth(100)
	sites.SetHeight(10)

	s := spinner.New()
	s.Spinner = spinner.Meter

	v := viewport.New(150, 20)

	return model{
		uri:     uri,
		sites:   sites,
		spinner: s,
		view:    v,
		stream:  make(chan struct{}),
		errors:  make(chan string),
	}
}

func startCrawl(uri string, links []string, stream chan struct{}, errors chan string) tea.Cmd {
	return func() tea.Msg {
		err := src.StartCrawl(uri, links, stream, errors)

		// could be nil
		return doneMsg{err: err}
	}
}

// will be batched with startCrawl
func wait(stream chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return conversion(<-stream)
	}
}

// this too
func waitErrors(errors chan string) tea.Cmd {
	return func() tea.Msg {
		return errorLogs(<-errors)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(wait(m.stream), waitErrors(m.errors), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case tea.KeyMsg:
		if m.crawling {
			if msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}

			return m, nil
		}

		switch msg.Type {
		case tea.KeyTab:
			if m.uri.Focused() {
				m.sites.Focus()
				m.uri.Blur()
			} else if m.sites.Focused() {
				m.uri.Focus()
				m.sites.Blur()
			}

			return m, nil
		case tea.KeyCtrlS:
			if m.uri.Focused() {
				m.sites.Focus()
				m.uri.Blur()

				return m, nil
			} else if m.sites.Focused() {
				m.crawling = true

				return m, tea.Batch(m.spinner.Tick, startCrawl(m.uri.Value(), strings.Fields(m.sites.Value()), m.stream, m.errors))
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case conversion:
		// important, if not program blocks forever
		if m.crawling {
			m.count++
			return m, wait(m.stream)
		}
		return m, nil
	case errorLogs:
		m.logs = append(m.logs, string(msg))
		var logs string
		for _, log := range m.logs {
			logs += fmt.Sprintf("%s\n", log)
		}

		m.view.SetContent(logs)

		// same here
		return m, waitErrors(m.errors)
	case doneMsg:
		m.crawling = false
		m.done = true
		m.err = msg.err

		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.crawling {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		m.view, cmd = m.view.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		if m.uri.Focused() {
			m.uri, cmd = m.uri.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.sites.Focused() {
			m.sites, cmd = m.sites.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var ui string

	if m.done {
		if m.err != nil {
			ui = errored.Render(fmt.Sprintf("Crawl failed: %s", m.err)) + "\n\n" + "Miscellaneous error logs" + "\n" + focused.Render(m.view.View()) + "\n" + bindings.Render("Ctrl-C/Esc to exit")
		} else {
			ui = finished.Render(fmt.Sprintf("Crawl completed! %d pages crawled.", m.count)) + "\n\n" + "Miscellaneous error logs" + "\n" + focused.Render(m.view.View()) + "\n" + bindings.Render("Ctrl-C/Esc to exit")
		}
	} else if m.crawling {
		ui = fmt.Sprintf("%s  Pages crawled: %d", m.spinner.View(), m.count) + "\n\n" + "Miscellaneous error logs" + "\n" + focused.Render(m.view.View()) + "\n" + bindings.Render("Ctrl-C/Esc to exit")
	} else {
		var uriView, sitesView string
		if m.uri.Focused() {
			uriView = focused.Render(m.uri.View())
			sitesView = m.sites.View()
		} else {
			uriView = m.uri.View()
			sitesView = focused.Render(m.sites.View())
		}

		keybinds := bindings.Render("Tab: Switch focus • Ctrl-S: Confirm input field • Ctrl-C/Esc: Exit")
		ui = uriView + "\n" + sitesView + "\n" + keybinds
	}

	// align on every view change
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
	}
	return ui
}
