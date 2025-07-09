package view

import (
	"strconv"
	"strings"
	"time"

	"github.com/bgics/pmjay-go/config"
	"github.com/bgics/pmjay-go/internal/tui"
	"github.com/bgics/pmjay-go/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name = iota
	address
	diagnosis
	date
	doa
	dob
	numDays
	printBtn
	saveBtn
)

const (
	day = iota
	month
	year
)

type FormPageModel struct {
	nameInput       textinput.Model
	addressInput    textinput.Model
	diagnosisInput  textinput.Model
	date            time.Time
	dateOfAdmission time.Time
	dateOfBirth     time.Time
	numDays         int

	fieldIndex           int
	dateInteractionIndex int

	// TODO: make error a field in shared state and then show error at bottom on each page
	errMsg      string
	sharedState *tui.SharedState
}

func NewFormPageModel(sharedState *tui.SharedState) *FormPageModel {
	m := &FormPageModel{}

	m.nameInput = makeTextInput(true, config.NAME)
	m.addressInput = makeTextInput(false, config.ADDRESS1, config.ADDRESS2, config.ADDRESS3)
	m.diagnosisInput = makeTextInput(false, config.DIAGNOSIS)

	m.sharedState = sharedState

	m.numDays = 1

	if m.sharedState.LastPageIndex == tui.SEARCH_PAGE {
		record := m.sharedState.SelectedRecord
		m.setFormWithRecord(record)
	} else {
		m.date = time.Now()
		m.dateOfAdmission = time.Now()
		m.dateOfBirth = time.Now()
	}

	return m
}

func (m *FormPageModel) setFormWithRecord(record model.FormData) {
	m.nameInput.SetValue(record.Name)
	m.addressInput.SetValue(record.Address)
	m.diagnosisInput.SetValue(record.Diagnosis)

	m.date = record.Date
	m.dateOfAdmission = record.DateOfAdmission
	m.dateOfBirth = record.DateOfBirth
}

func (m *FormPageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *FormPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.sharedState.LastPageIndex = tui.FORM_PAGE
			return m, tui.ChangePageCmd(tui.START_PAGE)
		case "tab", "shift+tab":
			return m.handleFormNav(msg.String())
		case "up", "down", "left", "right":
			if date <= m.fieldIndex && m.fieldIndex <= dob {
				return m.handleDateInput(msg.String())
			}

			if m.fieldIndex == numDays {
				return m.handleNumDaysInput(msg.String())
			}

			return m, nil
		}
	}

	return m.handleFormInput(msg)
}

func (m *FormPageModel) View() string {
	var output strings.Builder

	output.WriteString("\n")

	inputFields := m.renderTextInputs()
	dateFields := m.renderDateInputs()
	numDaysField := m.renderNumDaysField()
	formButtons := m.renderButtons()
	errorMsg := m.renderError()

	output.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		inputFields,
		dateFields,
		numDaysField,
		formButtons,
		errorMsg,
	))

	return output.String()
}

func (m *FormPageModel) handleFormInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := make([]tea.Cmd, 3)

	m.nameInput, cmd[0] = m.nameInput.Update(msg)
	m.addressInput, cmd[1] = m.addressInput.Update(msg)
	m.diagnosisInput, cmd[2] = m.diagnosisInput.Update(msg)

	return m, tea.Batch(cmd...)
}

func (m *FormPageModel) handleNumDaysInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up":
		m.numDays++
	case "down":
		m.numDays--
		if m.numDays < 1 {
			m.numDays = 1
		}
	}

	return m, nil
}

func (m *FormPageModel) handleDateInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "right":
		m.dateInteractionIndex = cyclicAdjust(m.dateInteractionIndex+1, day, year)
	case "left":
		m.dateInteractionIndex = cyclicAdjust(m.dateInteractionIndex-1, day, year)
	case "up", "down":
		m.handleDateAdjust(key)
	}

	return m, nil
}

func (m *FormPageModel) handleDateAdjust(key string) {
	var delta int

	if key == "up" {
		delta = 1
	} else {
		delta = -1
	}

	switch m.fieldIndex {
	case date:
		m.date = adjustDate(m.date, m.dateInteractionIndex, delta)
	case doa:
		m.dateOfAdmission = adjustDate(m.dateOfAdmission, m.dateInteractionIndex, delta)
	case dob:
		m.dateOfBirth = adjustDate(m.dateOfBirth, m.dateInteractionIndex, delta)
	}
}

func adjustDate(date time.Time, datePart int, delta int) time.Time {
	switch datePart {
	case day:
		date = date.AddDate(0, 0, delta)
	case month:
		date = date.AddDate(0, delta, 0)
	case year:
		date = date.AddDate(delta, 0, 0)
	}

	return date
}

func (m *FormPageModel) handleFormNav(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab":
		m.fieldIndex = cyclicAdjust(m.fieldIndex+1, name, saveBtn)
	case "shift+tab":
		m.fieldIndex = cyclicAdjust(m.fieldIndex-1, name, saveBtn)
	}

	if m.fieldIndex <= date || m.fieldIndex == saveBtn {
		cmd := m.updateFocus()
		return m, cmd
	}
	return m, nil
}

func (m *FormPageModel) updateFocus() tea.Cmd {
	m.nameInput.Blur()
	m.addressInput.Blur()
	m.diagnosisInput.Blur()

	var cmd tea.Cmd

	switch m.fieldIndex {
	case name:
		cmd = m.nameInput.Focus()
	case address:
		cmd = m.addressInput.Focus()
	case diagnosis:
		cmd = m.diagnosisInput.Focus()
	}

	return cmd
}

func cyclicAdjust(val, min, max int) int {
	if val > max {
		return min
	} else if val < min {
		return max
	}

	return val
}

func (m *FormPageModel) renderTextInputs() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		makeTextField("NAME", m.nameInput.View(), m.fieldIndex == name),
		makeTextField("ADDRESS", m.addressInput.View(), m.fieldIndex == address),
		makeTextField("DIAGNOSIS", m.diagnosisInput.View(), m.fieldIndex == diagnosis),
	)
}

func makeTextField(fieldName, inputView string, active bool) string {
	fieldNameStyle := tui.FieldNameActiveStyle
	if !active {
		fieldNameStyle = tui.FieldNameInactiveStyle
	}

	fieldInputStyle := tui.InputActiveBorderStyle
	if !active {
		fieldInputStyle = tui.InputInactiveBorderStyle
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		fieldNameStyle.Render(fieldName),
		fieldInputStyle.Render(inputView),
	)
}

func (m *FormPageModel) renderDateInputs() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderDateField("DATE", m.date, date),
		m.renderDateField("DOA", m.dateOfAdmission, doa),
		m.renderDateField("DOB", m.dateOfBirth, dob),
	)
}

func (m *FormPageModel) renderDateField(fieldName string, date time.Time, index int) string {
	if index == m.fieldIndex {
		return makeActiveDateField(fieldName, date, m.dateInteractionIndex)
	}
	return makeInactiveDateField(fieldName, date)
}

func makeActiveDateField(fieldName string, date time.Time, dateInteractionIndex int) string {
	var indicators string
	switch dateInteractionIndex {
	case day:
		indicators = "^^"
	case month:
		indicators = "   ^^"
	case year:
		indicators = "      ^^^^"
	}

	dateLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameActiveStyle.Render(fieldName),
		tui.DateFieldActiveStyle.Render(date.Format(config.DateFormat)),
	)

	interactionLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameActiveStyle.Render(""),
		tui.DateIndicatorStyle.Render(indicators),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		dateLine,
		interactionLine,
	)
}

func makeInactiveDateField(fieldName string, date time.Time) string {
	dateLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameInactiveStyle.Render(fieldName),
		tui.DateFieldInactiveStyle.Render(date.Format(config.DateFormat)),
	)

	return dateLine + "\n"
}

func (m *FormPageModel) renderNumDaysField() string {
	if m.fieldIndex == numDays {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			tui.FieldNameActiveStyle.Render("> DAYS"),
			tui.NumDaysFieldActiveStyle.Render(strconv.Itoa(m.numDays)),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameInactiveStyle.Render("  DAYS"),
		tui.NumDaysFieldInactiveStyle.Render(strconv.Itoa(m.numDays)),
	)
}

func (m *FormPageModel) renderButtons() string {
	printBtn := m.renderButton("PRINT", printBtn)
	saveBtn := m.renderButton("SAVE", saveBtn)

	return lipgloss.NewStyle().
		MarginLeft(6).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				printBtn,
				saveBtn,
			),
		)
}

func (m *FormPageModel) renderButton(btnName string, index int) string {
	if m.fieldIndex == index {
		return tui.BtnActiveStyle.Render(btnName)
	}
	return tui.BtnInactiveStyle.Render(btnName)
}

func (m *FormPageModel) renderError() string {
	if len(m.errMsg) > 0 {
		return tui.ErrStyle.Render("[ERROR]: " + m.errMsg)
	}
	return ""
}
