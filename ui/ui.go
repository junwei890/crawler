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
	labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))

	focusedInputStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#F8F8FF"))
	
	blurredInputStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#778899"))
	
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#778899")).Bold(true)
	
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#778899")).Bold(true)
	
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC143C")).Bold(true)
	
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")).Bold(true)
	
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9")).Italic(true)
)

type crawlDoneMsg struct {
	err error
}

type model struct {
	uri      textinput.Model
	links    textarea.Model
	spinner  spinner.Model
	crawling bool
	done     bool
	err      error
	width    int
	height   int
}

func InitialModel() model {
	uri := textinput.New()
	uri.Focus()
	uri.Width = 97

	links := textarea.New()
	links.SetWidth(100)
	links.SetHeight(20)

	spin := spinner.New()
	spin.Spinner = spinner.Meter
	spin.Style = spinnerStyle

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
				m.links.Focus()
				m.uri.Blur()
			} else if m.links.Focused() {
				m.uri.Focus()
				m.links.Blur()
			}

			return m, nil

		case tea.KeyCtrlS:
			if m.uri.Focused() {
				m.links.Focus()
				m.uri.Blur()

				return m, nil
			} else if m.links.Focused() {
				m.crawling = true

				return m, tea.Batch(m.spinner.Tick, startCrawl(m.uri.Value(), strings.Fields(m.links.Value())))
			}

		case tea.KeyCtrlC, tea.KeyEsc:
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
		if m.uri.Focused() {
			m.uri, cmd = m.uri.Update(msg)
		} else if m.links.Focused() {
			m.links, cmd = m.links.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var content string

	if m.done {
		if m.err != nil {
			content = fmt.Sprintf("Crawling failed: %s\n\nPress Ctrl-C to exit.", m.err)

			content = errorStyle.Render(content)
		} else {
			content = "Crawling completed successfully!\n\nPress Ctrl-C to exit."

			content = successStyle.Render(content)
		}
	} else if m.crawling {
		content = fmt.Sprintf("%s Crawling and indexing, this might take awhile...\n\nPress Ctrl-C to cancel.", m.spinner.View())
		
		content = statusStyle.Render(content)
	} else {
		var uriView, linksView string
		
		if m.uri.Focused() {
			uriView = focusedInputStyle.Render(m.uri.View())
			linksView = blurredInputStyle.Render(m.links.View())
		} else {
			uriView = blurredInputStyle.Render(m.uri.View())
			linksView = focusedInputStyle.Render(m.links.View())
		}

		help := helpStyle.Render("Tab: Switch focus • Ctrl-S: Confirm input field • Ctrl-C/Esc: Exit")
		content = labelStyle.Render("MongoDB URI") + "\n" + uriView + "\n" + labelStyle.Render("Links to Crawl") + "\n" + linksView + "\n" + help
	}

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
