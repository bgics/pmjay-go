package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	START_MENU = iota
	FORM_PAGE
	DATE_EDIT
)

type AppState int

var (
	focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("231"))
	unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))

	cursorStyle = focusedStyle

	choices = []string{"New Patient", "Search Records"}
)

type FormPageModel struct {
	NameInput      textinput.Model
	DiagnosisInput textinput.Model
	AddressInput   textinput.Model

	Date            time.Time
	DateOfAdmission time.Time
	DateOfBirth     time.Time

	NumDays int
}

type model struct {
	inputs     []textinput.Model
	focusIndex int

	choiceIndex int

	state AppState
}

func initModel() model {
	m := model{
		inputs: make([]textinput.Model, 7),
		state:  START_MENU,
	}

	m.inputs[0] = makeFormInput("Name", "", NAME, true)
	m.inputs[1] = makeFormInput("Address", "", ADDRESS1, false)
	m.inputs[2] = makeFormInput("Date", time.Now().Format("02/01/2006"), DATE, false)
	m.inputs[3] = makeFormInput("Date of Admission", time.Now().Format("02/01/2006"), DATE_OF_ADMISSION, false)
	m.inputs[4] = makeFormInput("Date of Birth", time.Now().Format("02/01/2006"), DATE_OF_BIRTH, false)
	m.inputs[5] = makeFormInput("Diagnosis", "", DIAGNOSIS, false)

	m.inputs[6] = makeInput("Number of Days", "1", NumDaysCharLimit, false)

	return m
}

func makeFormInput(prompt, defaultValue string, cfgKey FieldName, focus bool) textinput.Model {
	var maxChars int
	if cfgKey == ADDRESS1 {
		maxChars = FieldConfig[ADDRESS1].MaxChars + FieldConfig[ADDRESS2].MaxChars + FieldConfig[ADDRESS3].MaxChars
	} else {
		maxChars = FieldConfig[cfgKey].MaxChars
	}

	return makeInput(prompt, defaultValue, maxChars, focus)
}

func makeInput(prompt, defaultValue string, charLimit int, focus bool) textinput.Model {
	t := textinput.New()

	t.Prompt = fmt.Sprintf("%s: ", prompt)

	t.CharLimit = charLimit

	if focus {
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle
	} else {
		t.PromptStyle = unfocusedStyle
		t.TextStyle = unfocusedStyle
	}

	t.Cursor.Style = cursorStyle
	t.Cursor.SetMode(cursor.CursorStatic)

	t.SetValue(defaultValue)
	t.Width = TextInputWidth

	return t
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down", "left", "right", "enter", "esc":
			switch m.state {
			case START_MENU:
				return m.HandleStartNav(msg.String())
			case FORM_PAGE:
				return m.HandleFormNav(msg.String())
			case DATE_EDIT:
				return m.HandleDateNav(msg.String())
			}
		}
	}

	if m.state == FORM_PAGE {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case START_MENU:
		return m.RenderStartMenu()
	case FORM_PAGE, DATE_EDIT:
		return m.RenderFormPage()
	default:
		// TODO: make an error page
		return "Error"
	}
}

func (m *model) RenderStartMenu() string {
	var output strings.Builder

	output.WriteString("\n\n")
	for i, choice := range choices {
		if i == m.choiceIndex {
			output.WriteString(fmt.Sprintf("  > %s\n\n", focusedStyle.Render(choice)))
		} else {
			output.WriteString(fmt.Sprintf("    %s\n\n", unfocusedStyle.Render(choice)))
		}
	}

	return output.String()
}

func (m *model) RenderFormPage() string {
	var output strings.Builder

	output.WriteString("\n\n")
	for _, input := range m.inputs {
		output.WriteString(fmt.Sprintf("    %s\n\n", input.View()))
	}

	if m.focusIndex == len(m.inputs) {
		output.WriteString("    " + focusedStyle.Render("[ Print ]"))
	} else {
		output.WriteString("    " + unfocusedStyle.Render("[ Print ]"))
	}

	return output.String()
}

func (m model) HandleDateNav(input string) (model, tea.Cmd) {
	switch input {
	case "esc":
		m.inputs[m.focusIndex].Cursor.SetMode(cursor.CursorStatic)
		m.state = FORM_PAGE
	}

	return m, nil
}

func (m model) HandleFormNav(input string) (model, tea.Cmd) {
	switch input {
	case "enter":
		if 2 <= m.focusIndex && m.focusIndex <= 4 {
			m.inputs[m.focusIndex].Cursor.SetMode(cursor.CursorHide)
			m.state = DATE_EDIT

			return m, nil
		}
	case "up", "shift+tab":
		m.focusIndex--
	case "down", "tab":
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	cmd := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focusIndex {
			cmd[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}

		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = unfocusedStyle
		m.inputs[i].TextStyle = unfocusedStyle
	}

	return m, tea.Batch(cmd...)
}

func (m model) HandleStartNav(input string) (model, tea.Cmd) {
	switch input {
	case "enter":
		// TODO: add search page
		m.focusIndex = 0
		m.state = FORM_PAGE

		return m, m.inputs[0].Focus()
	case "up", "shift+tab":
		m.choiceIndex--
	case "down", "tab":
		m.choiceIndex++
	}

	if m.choiceIndex > len(choices)-1 {
		m.choiceIndex = 0
	} else if m.choiceIndex < 0 {
		m.choiceIndex = len(choices) - 1
	}

	return m, nil
}
func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
