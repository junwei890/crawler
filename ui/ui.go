package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/junwei890/crawler/src"
)

type crawlDoneMsg struct {
	err error
}

type model struct {
	uri      textinput.Model
	links    textarea.Model
	spinner  spinner.Model
	focus    int
	crawling bool
	done     bool
	err      error
}

func InitialModel() model {
	uri := textinput.New()
	uri.Focus()
	uri.Width = 100

	links := textarea.New()
	links.SetWidth(100)
	links.SetHeight(20)

	spin := spinner.New()
	spin.Spinner = spinner.Meter

	return model{
		uri:     uri,
		links:   links,
		spinner: spin,
	}
}

func startCrawl(uri string, links []string) tea.Cmd {
	return func() tea.Msg {
		err := src.StartCrawl(uri, links)
		return crawlDoneMsg{err: err}
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.crawling {
			if msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyTab:
			m.focus = (m.focus + 1) % 2
			switch m.focus {
			case 0:
				m.uri.Focus()
				m.links.Blur()
			case 1:
				m.links.Focus()
				m.uri.Blur()
			}
			return m, nil
		case tea.KeyCtrlS:
			switch m.focus {
			case 0:
				m.focus = 1
				m.links.Focus()
				m.uri.Blur()
				return m, nil
			case 1:
				m.crawling = true
				return m, tea.Batch(m.spinner.Tick, startCrawl(m.uri.Value(), strings.Fields(m.links.Value())))
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case crawlDoneMsg:
		m.crawling = false
		m.done = true
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	if m.crawling {
		m.spinner, cmd = m.spinner.Update(msg)
	} else {
		switch m.focus {
		case 0:
			m.uri, cmd = m.uri.Update(msg)
		case 1:
			m.links, cmd = m.links.Update(msg)
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.done {
		if m.err != nil {
			return fmt.Sprintf(`Crawling failed with %s.
Ctrl-C to exit.`, m.err)
		}
		return `Crawling done!
Ctrl-C to exit.`
	}

	if m.crawling {
		return fmt.Sprintf("%s Crawling... There may be sites with long crawl delays, this might take awhile...", m.spinner.View())
	}

	return fmt.Sprintf("MongoDB URI\n%s\nLinks\n%s\nSwitch focus: Tab | Confirm input: Ctrl-S | Exit: Ctrl-C", m.uri.View(), m.links.View())
}
