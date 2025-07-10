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
	datepicker "github.com/ethanefung/bubble-datepicker"
)

const (
	nameIndex = iota
	addressIndex
	diagnosisIndex
	dateIndex
	doaIndex
	dobIndex
	genderIndex
	numDaysIndex
	printBtnIndex
	saveBtnIndex
)

type FormPageModel struct {
	nameInput      textinput.Model
	addressInput   textinput.Model
	diagnosisInput textinput.Model

	gender model.Gender

	date            time.Time
	dateOfAdmission time.Time
	dateOfBirth     time.Time

	numDays    int
	fieldIndex int

	datePicker     datepicker.Model
	datePickerMode bool

	// TODO: make error a field in shared state and then show error at bottom on each page
	errMsg      string
	sharedState *tui.SharedState
}

func NewFormPageModel(sharedState *tui.SharedState) *FormPageModel {
	m := &FormPageModel{}

	m.nameInput = makeTextInput(true, config.NAME)
	m.addressInput = makeTextInput(false, config.ADDRESS1, config.ADDRESS2, config.ADDRESS3)
	m.diagnosisInput = makeTextInput(false, config.DIAGNOSIS)

	m.numDays = 1

	m.datePicker = makeDateInput()
	m.datePickerMode = false

	m.sharedState = sharedState

	if m.sharedState.LastPageIndex == tui.SEARCH_PAGE {
		record := m.sharedState.SelectedRecord
		m.setFormWithRecord(record)
	} else {
		m.date = time.Now()
		m.dateOfAdmission = time.Now()
		m.dateOfBirth = time.Now()
		m.gender = model.Male
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

	m.gender = record.Gender
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
		case "tab", "down", "shift+tab", "up":
			if m.datePickerMode {
				return m.handleDatePicker(msg)
			}

			return m.handleFormNav(msg.String())
		case "left", "right":
			if m.datePickerMode {
				return m.handleDatePicker(msg)
			}

			if m.fieldIndex == numDaysIndex {
				return m.handleNumDaysInput(msg.String())
			}

			if m.fieldIndex == genderIndex {
				return m.handleGenderInput()
			}

			return m, nil
		case "enter":
			if dateIndex <= m.fieldIndex && m.fieldIndex <= dobIndex {
				if m.datePickerMode {
					m.datePickerMode = false
					m.blurDatePicker()
				} else {
					m.datePickerMode = true
					m.focusDatePicker()
				}

				return m, nil
			}
		}
	}

	return m.handleFormInput(msg)
}

func (m *FormPageModel) View() string {
	var output strings.Builder

	output.WriteString("\n")

	inputFields := m.renderTextInputs()

	dateFields := m.renderDateInputs()
	genderField := m.renderGenderField()
	numDaysField := m.renderNumDaysField()
	formButtons := m.renderButtons()

	errorMsg := m.renderError()

	bottomFields := lipgloss.JoinVertical(
		lipgloss.Left,
		dateFields,
		genderField,
		numDaysField,
		formButtons,
	)

	datePicker := m.renderDatePicker()

	bottomForm := lipgloss.JoinHorizontal(
		lipgloss.Left,
		bottomFields,
		datePicker,
	)

	output.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		inputFields,
		bottomForm,
		errorMsg,
	))

	return output.String()
}

func (m *FormPageModel) handleGenderInput() (tea.Model, tea.Cmd) {
	if m.gender == model.Male {
		m.gender = model.Female
	} else {
		m.gender = model.Male
	}

	return m, nil
}

func (m *FormPageModel) focusDatePicker() {
	switch m.fieldIndex {
	case dateIndex:
		m.datePicker.SetTime(m.date)
	case doaIndex:
		m.datePicker.SetTime(m.dateOfAdmission)
	case dobIndex:
		m.datePicker.SetTime(m.dateOfBirth)
	}

	m.datePicker.SelectDate()
	m.datePicker.SetFocus(datepicker.FocusCalendar)
}

func (m *FormPageModel) blurDatePicker() {
	m.datePicker.UnselectDate()
	m.datePicker.SetFocus(datepicker.FocusNone)
}

func (m *FormPageModel) handleDatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	prev := m.datePicker.Time

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.datePicker.Focused == datepicker.FocusCalendar {
				m.datePicker.SetFocus(datepicker.FocusHeaderMonth)
			} else {
				m.datePicker, cmd = m.datePicker.Update(msg)
			}
		case "shift+tab":
			if m.datePicker.Focused == datepicker.FocusHeaderMonth {
				m.datePicker.SetFocus(datepicker.FocusCalendar)
			} else {
				m.datePicker, cmd = m.datePicker.Update(msg)
			}
		default:
			m.datePicker, cmd = m.datePicker.Update(msg)
		}
	}

	if m.datePicker.Time != prev {
		switch m.fieldIndex {
		case dateIndex:
			m.date = m.datePicker.Time
		case doaIndex:
			m.dateOfAdmission = m.datePicker.Time
		case dobIndex:
			m.dateOfBirth = m.datePicker.Time
		}
	}

	return m, cmd
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
	case "right":
		m.numDays++
	case "left":
		m.numDays--
		if m.numDays < 1 {
			m.numDays = 1
		}
	}

	return m, nil
}

func (m *FormPageModel) handleFormNav(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "down", "tab":
		m.fieldIndex = cyclicAdjust(m.fieldIndex+1, nameIndex, saveBtnIndex)
	case "up", "shift+tab":
		m.fieldIndex = cyclicAdjust(m.fieldIndex-1, nameIndex, saveBtnIndex)
	}

	if m.fieldIndex <= dateIndex || m.fieldIndex == saveBtnIndex {
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
	case nameIndex:
		cmd = m.nameInput.Focus()
	case addressIndex:
		cmd = m.addressInput.Focus()
	case diagnosisIndex:
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

func (m *FormPageModel) renderGenderField() string {
	if m.fieldIndex == genderIndex {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			tui.FieldNameActiveStyle.Render("> GENDER"),
			tui.SimpleFieldActiveStyle.Render(string(m.gender)),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameInactiveStyle.Render("  GENDER"),
		tui.SimpleFieldInactiveStyle.Render(string(m.gender)),
	)
}

func (m *FormPageModel) renderDatePicker() string {
	if !m.datePickerMode {
		return ""
	}
	return tui.DatePickerStyle.Render(m.datePicker.View())
}

func (m *FormPageModel) renderTextInputs() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		makeTextField("NAME", m.nameInput.View(), m.fieldIndex == nameIndex),
		makeTextField("ADDRESS", m.addressInput.View(), m.fieldIndex == addressIndex),
		makeTextField("DIAGNOSIS", m.diagnosisInput.View(), m.fieldIndex == diagnosisIndex),
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
		m.renderDateField("DATE", m.date, dateIndex),
		m.renderDateField("DOA", m.dateOfAdmission, doaIndex),
		m.renderDateField("DOB", m.dateOfBirth, dobIndex),
	)
}

func (m *FormPageModel) renderDateField(fieldName string, date time.Time, index int) string {
	prefix := "  "
	if m.fieldIndex == index {
		prefix = "> "
	}

	return makeDateField(prefix+fieldName, date, index == m.fieldIndex)
}

func makeDateField(fieldName string, date time.Time, active bool) string {
	fieldNameStyle := tui.FieldNameActiveStyle
	if !active {
		fieldNameStyle = tui.FieldNameInactiveStyle
	}

	dateStyle := tui.DateFieldActiveStyle
	if !active {
		dateStyle = tui.DateFieldInactiveStyle
	}

	dateLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fieldNameStyle.Render(fieldName),
		dateStyle.Render(date.Format(config.DateFormat)),
	)

	return dateLine + "\n"
}

func (m *FormPageModel) renderNumDaysField() string {
	if m.fieldIndex == numDaysIndex {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			tui.FieldNameActiveStyle.Render("> DAYS"),
			tui.SimpleFieldActiveStyle.Render(strconv.Itoa(m.numDays)),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		tui.FieldNameInactiveStyle.Render("  DAYS"),
		tui.SimpleFieldInactiveStyle.Render(strconv.Itoa(m.numDays)),
	)
}

func (m *FormPageModel) renderButtons() string {
	printBtn := m.renderButton("PRINT", printBtnIndex)
	saveBtn := m.renderButton("SAVE", saveBtnIndex)

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
