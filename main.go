package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"social-media-tools/downloader"
	"social-media-tools/internal/clipboard"
	"social-media-tools/internal/storage"
	"social-media-tools/platform"
	"social-media-tools/platform/youtube"
	"social-media-tools/ui/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateDetecting state = iota
	stateSelecting
	stateDownloading
	stateFinished
)

type model struct {
	state            state
	url              string
	mediaType        string
	usePWD           bool
	progress         progress.Model
	textInput        textinput.Model
	logs             []string
	platforms        []platform.Platform
	err              error
	downloading      bool
	progressChan     chan downloader.ProgressMsg
	width            int
	height           int
	showModal        bool
	clipboardChecked bool
}

func initialModel() model {
	p := progress.New(progress.WithGradient("#2997ff", "#0071e3"))
	p.Width = 60

	ti := textinput.New()
	ti.Placeholder = "https://www.youtube.com/watch?v=..."
	ti.KeyMap.Paste.SetKeys("ctrl+v", "ctrl+shift+v", "shift+insert")

	return model{
		state:        stateDetecting,
		progress:     p,
		textInput:    ti,
		platforms:    []platform.Platform{&youtube.YouTube{}},
		progressChan: make(chan downloader.ProgressMsg),
	}
}

func (m model) Init() tea.Cmd {
	return checkClipboard
}

func checkClipboard() tea.Msg {
	url := clipboard.GetURL()
	if url != "" {
		return foundURLMsg(url)
	}
	return noURLMsg{}
}

type foundURLMsg string
type noURLMsg struct{}
type downloadProgressMsg downloader.ProgressMsg

func (m *model) selectURL(url string) {
	m.url = strings.TrimSpace(url)
	m.textInput.SetValue(m.url)
	m.textInput.CursorEnd()
	m.state = stateSelecting
	m.showModal = false
	m.err = nil
}

func (m *model) openURLModal(initialValue string) tea.Cmd {
	m.showModal = true
	m.textInput.SetValue(strings.TrimSpace(initialValue))
	m.textInput.CursorEnd()
	return m.textInput.Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 20
		return m, nil

	case tea.KeyMsg:
		if m.showModal {
			switch msg.String() {
			case "esc":
				m.showModal = false
				m.textInput.Blur()
				if m.url == "" {
					m.state = stateDetecting
				}
				return m, nil
			case "enter":
				val := strings.TrimSpace(m.textInput.Value())
				if val != "" {
					m.textInput.Blur()
					m.selectURL(val)
					return m, nil
				}
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			if m.state != stateDownloading {
				m.err = nil
				return m, m.openURLModal("")
			}
		case "c", "C":
			if m.state != stateDownloading {
				m.usePWD = !m.usePWD
				return m, nil
			}
		case "1", "2", "3":
			if m.state == stateSelecting {
				switch msg.String() {
				case "1":
					m.mediaType = "video"
				case "2":
					m.mediaType = "audio"
				case "3":
					m.mediaType = "thumbnail"
				}
				m.state = stateDownloading
				m.downloading = true
				m.logs = []string{}
				return m, m.startDownload()
			}
		}

	case foundURLMsg:
		m.clipboardChecked = true
		if m.state == stateDetecting && m.url == "" {
			m.selectURL(string(msg))
		}
		return m, nil

	case noURLMsg:
		m.clipboardChecked = true
		return m, nil

	case downloadProgressMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.state = stateFinished
			m.downloading = false
			return m, nil
		}
		if msg.Done {
			m.state = stateFinished
			m.downloading = false
			cmd := m.progress.SetPercent(1.0)
			return m, cmd
		}
		var cmds []tea.Cmd
		if msg.Percent > 0 {
			cmds = append(cmds, m.progress.SetPercent(msg.Percent))
		}
		if msg.Log != "" {
			m.logs = append(m.logs, msg.Log)
			if len(m.logs) > 10 {
				m.logs = m.logs[1:]
			}
		}
		cmds = append(cmds, m.waitForProgress())
		return m, tea.Batch(cmds...)

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) startDownload() tea.Cmd {
	platformName := "Other"
	for _, p := range m.platforms {
		if p.Match(m.url) {
			platformName = p.GetName()
			break
		}
	}

	outputPath := storage.GetDownloadPath(m.mediaType, platformName, m.usePWD)
	go downloader.Download(m.url, m.mediaType, outputPath, m.progressChan)
	return m.waitForProgress()
}

func (m model) currentDownloadPath() string {
	platformName := "Other"
	for _, p := range m.platforms {
		if p.Match(m.url) {
			platformName = p.GetName()
			break
		}
	}

	return storage.GetDownloadPath(m.mediaType, platformName, m.usePWD)
}

func formatDownloadPath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil {
		if rel, relErr := filepath.Rel(home, path); relErr == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			path = filepath.Join("~", rel)
		}
	}

	cleaned := filepath.Clean(path)
	parts := strings.Split(cleaned, string(filepath.Separator))
	if len(parts) == 0 {
		return cleaned
	}

	start := 0
	if parts[0] == "" {
		start = 1
	}
	if len(parts)-start <= 3 {
		return cleaned
	}

	return filepath.Join(append([]string{"..."}, parts[len(parts)-3:]...)...)
}

func (m model) waitForProgress() tea.Cmd {
	return func() tea.Msg {
		return downloadProgressMsg(<-m.progressChan)
	}
}

func (m model) View() string {
	var innerContent string

	header := lipgloss.NewStyle().Width(m.width-10).Align(lipgloss.Center).Render(styles.Title.Render(" SOCIAL MEDIA TOOLS ")) + "\n\n"

	switch m.state {
	case stateDetecting:
		status := "Checking clipboard for URL..."
		if m.clipboardChecked {
			status = "No URL found on clipboard."
		}
		innerContent = lipgloss.NewStyle().Width(m.width - 10).Align(lipgloss.Center).Render(
			status + "\n\n" +
				styles.ShortcutHint.Render("Press 'n' to paste or type a URL"),
		)

	case stateSelecting:
		locationLabel := "Default Path"
		if m.usePWD {
			locationLabel = "Current Directory"
		}
		displayPath := formatDownloadPath(m.currentDownloadPath())
		innerContent = lipgloss.NewStyle().Width(m.width - 10).Align(lipgloss.Center).Render(
			fmt.Sprintf("URL: %s\n\n", m.url) +
				"Choose Format:\n\n" +
				styles.Button("1. Video") + " " +
				styles.Button("2. Audio") + " " +
				styles.Button("3. Thumbnail") + "\n\n" +
				fmt.Sprintf("Download to: %s\n%s\n\n", locationLabel, displayPath) +
				styles.ShortcutHint.Render("Press 1, 2, or 3 • Press 'c' to toggle path • Press 'n' for a different URL"),
		)

	case stateDownloading:
		innerContent = lipgloss.NewStyle().Width(m.width - 10).Align(lipgloss.Center).Render(
			fmt.Sprintf("Downloading %s...\n\n", m.mediaType) +
				m.progress.View() + "\n\n" +
				"Logs:\n" +
				styles.LogStyle.Width(m.width-20).Render(strings.Join(m.logs, "\n")),
		)

	case stateFinished:
		if m.err != nil {
			errMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3b30")).Bold(true).Render("✗ Download Failed")
			errDetail := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")).Render(fmt.Sprintf("%v", m.err))
			innerContent = lipgloss.NewStyle().Width(m.width - 10).Align(lipgloss.Center).Render(
				errMsg + "\n\n" +
				errDetail + "\n\n" +
				styles.ShortcutHint.Render("Press 'n' to try again • Press 'q' to quit"),
			)
		} else {
			successIcon := lipgloss.NewStyle().Foreground(lipgloss.Color("#30d158")).Bold(true).Render("✓ Download Complete!")
			pathInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("#86868b")).Render(
				fmt.Sprintf("Saved to: %s", formatDownloadPath(m.currentDownloadPath())),
			)
			innerContent = lipgloss.NewStyle().Width(m.width - 10).Align(lipgloss.Center).Render(
				successIcon + "\n\n" +
				pathInfo + "\n\n" +
				m.progress.View() + "\n\n" +
				styles.ShortcutHint.Render("Press 'q' to quit • Press 'n' to download another"),
			)
		}
	}

	if m.usePWD && m.state != stateSelecting {
		footer := "\n\n" + lipgloss.NewStyle().Width(m.width-10).Align(lipgloss.Center).Render(styles.ButtonInactive("Saving to: Current Directory"))
		innerContent += footer
	}

	fullView := styles.WindowStyle.
		Width(m.width - 4).
		Height(m.height - 4).
		Render(lipgloss.JoinVertical(lipgloss.Center, header, innerContent))

	if m.showModal {
		modalContent := lipgloss.JoinVertical(lipgloss.Center,
			styles.Title.Render(" ENTER URL "),
			"\n",
			m.textInput.View(),
			"\n\n",
			styles.ShortcutHint.Render("Enter to confirm • Ctrl+V / Ctrl+Shift+V to paste • Esc to cancel"),
		)
		modal := styles.ModalStyle.Render(modalContent)

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fullView)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
