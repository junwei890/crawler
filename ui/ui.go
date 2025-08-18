package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/junwei890/crawler/src"
)

var (
	labels = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C562AF"))

	focused = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#B33791"))

	blurred = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FEC5F6"))

	status = lipgloss.NewStyle().Foreground(lipgloss.Color("#FEC5F6")).Bold(true)

	errored = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC3C22")).Bold(true)

	finished = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")).Bold(true)

	bindings = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9")).Italic(true)
)

type conversion struct{}

type doneMsg struct {
	err error
}

type model struct {
	uri      textinput.Model
	sites    textarea.Model
	spinner  spinner.Model
	stream   chan struct{}
	count    int
	crawling bool
	done     bool
	err      error
	width    int
	height   int
}

func InitialModel() model {
	uri := textinput.New()
	uri.Placeholder = "Paste your MongoDB URI here"
	uri.Focus()
	uri.Width = 97

	sites := textarea.New()
	sites.Placeholder = "Paste sites here, making sure each site is on a newline"
	sites.SetWidth(100)
	sites.SetHeight(20)

	s := spinner.New()
	s.Spinner = spinner.Meter

	return model{
		uri:     uri,
		sites:   sites,
		spinner: s,
		stream:  make(chan struct{}),
	}
}

func startCrawl(uri string, links []string, stream chan struct{}) tea.Cmd {
	return func() tea.Msg {
		err := src.StartCrawl(uri, links, stream)

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

func (m model) Init() tea.Cmd {
	return textinput.Blink
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

				return m, tea.Batch(m.spinner.Tick, wait(m.stream), startCrawl(m.uri.Value(), strings.Fields(m.sites.Value()), m.stream))
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case conversion:
		m.count++

		// important, if not program blocks forever
		return m, wait(m.stream)
	case doneMsg:
		m.crawling = false
		m.done = true
		m.err = msg.err

		return m, nil
	}

	var cmd tea.Cmd
	if m.crawling {
		m.spinner, cmd = m.spinner.Update(msg)
	} else {
		if m.uri.Focused() {
			m.uri, cmd = m.uri.Update(msg)
		} else if m.sites.Focused() {
			m.sites, cmd = m.sites.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var ui string

	if m.done {
		if m.err != nil {
			ui = errored.Render(fmt.Sprintf("Crawling failed: %s\n\nCtrl-C to exit.", m.err))
		} else {
			ui = finished.Render(fmt.Sprintf("Crawling completed! %d pages crawled.\n\nCtrl-C to exit.", m.count))
		}
	} else if m.crawling {
		ui = status.Render(fmt.Sprintf("%s  Pages crawled: %d", m.spinner.View(), m.count))
	} else {
		var uriView, sitesView string
		if m.uri.Focused() {
			uriView = focused.Render(m.uri.View())
			sitesView = blurred.Render(m.sites.View())
		} else {
			uriView = blurred.Render(m.uri.View())
			sitesView = focused.Render(m.sites.View())
		}

		keybinds := bindings.Render("Tab: Switch focus • Ctrl-S: Confirm input field • Ctrl-C/Esc: Exit")
		ui = labels.Render("MongoDB URI") + "\n" + uriView + "\n" + labels.Render("Sites") + "\n" + sitesView + "\n" + keybinds
	}

	// align on every view change
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
	}
	return ui
}
