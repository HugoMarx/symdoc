package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"symdoc/internal/config"
	"symdoc/internal/repo"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/charmbracelet/glamour"
)

type model struct {
	width    int
	height   int
	filePath string
	content  string
	viewport viewport.Model
	help     help.Model
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 5).Margin(1, 0)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).BorderStyle(b)
	}()
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

func (m model) View() tea.View {
	view := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	return tea.NewView(view)
}

func main() {
	config := config.Load()
	if !config.MdGenerated {
		if err := repo.BuildMarkdown(config.AssetsDir); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("Already converted !")
	}

	fileName := os.Args[1:][0] // Starts at index 1 since 0 is the program path.

	filePath := filepath.Join(os.ExpandEnv(config.AssetsDir), fileName)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("File not found : ", err.Error())
	}

	output, err := glamour.Render(string(fileContent), "dark")
	if err != nil {
		log.Fatal("Unable to render file content : ", err.Error())
	}

	vp := viewport.New(viewport.WithWidth(90), viewport.WithHeight(30)) // largeur, hauteur
	vp.SetContent(output)

	p := tea.NewProgram(model{
		filePath: filePath,
		content:  output,
		viewport: vp,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func (m model) headerView() string {
	mainTitle := getMainTitle(m.filePath)

	title := titleStyle.Render(mainTitle)
	line := strings.Repeat(" ", max(0, (m.viewport.Width()/2)-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width()-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func getMainTitle(filePath string) string {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return filePath
	}

	pattern := regexp.MustCompile(`^#\s+(.+)`)
	match := pattern.FindSubmatch(fileContent)

	return string(match[1])
}
