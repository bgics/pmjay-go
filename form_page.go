package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type formPageModel struct {
	nameInput       textinput.Model
	diagnosisInput  textinput.Model
	addressInput    textinput.Model
	date            time.Time
	dateOfAdmission time.Time
	dateOfBirth     time.Time
	numDays         int

	fieldIndex           int
	dateInteractionIndex int

	errMsg string
}

func newFormPageModel() *formPageModel {
	m := &formPageModel{}

	m.nameInput = makeFormTextInput(true)
	m.addressInput = makeFormTextInput(false)
	m.diagnosisInput = makeFormTextInput(false)

	m.date = time.Now()
	m.dateOfAdmission = time.Now()
	m.dateOfBirth = time.Now()

	m.numDays = 1
	m.errMsg = ""

	return m
}

func makeFormTextInput(focus bool) textinput.Model {
	t := textinput.New()

	t.Cursor.SetMode(cursor.CursorStatic)

	t.Prompt = " "
	t.Width = 40

	if focus {
		t.Focus()
	}

	return t
}

func formPageModelFromRecord(fd formData) *formPageModel {
	m := newFormPageModel()

	m.nameInput.SetValue(fd.name)
	m.addressInput.SetValue(fd.address)
	m.diagnosisInput.SetValue(fd.diagnosis)

	m.date = fd.date
	m.dateOfAdmission = fd.dateOfAdmission
	m.dateOfBirth = fd.dateOfBirth

	return m
}

func (m *formPageModel) update(msg tea.Msg) (tea.Cmd, int) {
	prevFieldIndex := m.fieldIndex

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, START_PAGE
		case "shift+tab":
			m.fieldIndex--
		case "tab":
			m.fieldIndex++
		case "enter":
			if m.fieldIndex == 7 {
				fd, err := m.validateInput()
				if err != nil {
					m.errMsg = err.Error()
					return nil, FORM_PAGE
				}

				return m.generatePrintPDFCmd(fd), START_PAGE
			}
		case "left", "right":
			if 3 > m.fieldIndex || m.fieldIndex > 5 {
				return nil, FORM_PAGE
			}

			if msg.String() == "left" {
				m.dateInteractionIndex--
			} else {
				m.dateInteractionIndex++
			}

			if m.dateInteractionIndex > 2 {
				m.dateInteractionIndex = 0
			} else if m.dateInteractionIndex < 0 {
				m.dateInteractionIndex = 2
			}

			return nil, FORM_PAGE
		case "up", "down":
			if 3 > m.fieldIndex || m.fieldIndex > 6 {
				return nil, FORM_PAGE
			}

			inc := 1
			if msg.String() == "down" {
				inc = -1
			}

			if m.fieldIndex == 6 {
				m.numDays += inc
				if m.numDays < 1 {
					m.numDays = 1
				}

				return nil, FORM_PAGE
			}

			deltaDay := 0
			deltaMonth := 0
			deltaYear := 0

			switch m.dateInteractionIndex {
			case 0:
				deltaDay = inc
			case 1:
				deltaMonth = inc
			case 2:
				deltaYear = inc
			}

			switch m.fieldIndex {
			case 3:
				m.date = m.date.AddDate(deltaYear, deltaMonth, deltaDay)
			case 4:
				m.dateOfAdmission = m.dateOfAdmission.AddDate(deltaYear, deltaMonth, deltaDay)
			case 5:
				m.dateOfBirth = m.dateOfBirth.AddDate(deltaYear, deltaMonth, deltaDay)
			}

			return nil, FORM_PAGE
		}
	}

	if m.fieldIndex > 7 {
		m.fieldIndex = 0
	} else if m.fieldIndex < 0 {
		m.fieldIndex = 7
	}

	if prevFieldIndex != m.fieldIndex {
		var cmd tea.Cmd
		switch m.fieldIndex {
		case 0:
			cmd = m.nameInput.Focus()
			m.addressInput.Blur()
			m.diagnosisInput.Blur()
		case 1:
			m.nameInput.Blur()
			cmd = m.addressInput.Focus()
			m.diagnosisInput.Blur()
		case 2:
			m.nameInput.Blur()
			m.addressInput.Blur()
			cmd = m.diagnosisInput.Focus()
		default:
			m.nameInput.Blur()
			m.addressInput.Blur()
			m.diagnosisInput.Blur()
		}

		return cmd, FORM_PAGE
	}

	cmd := make([]tea.Cmd, 3)

	m.nameInput, cmd[0] = m.nameInput.Update(msg)
	m.addressInput, cmd[1] = m.addressInput.Update(msg)
	m.diagnosisInput, cmd[2] = m.diagnosisInput.Update(msg)

	return tea.Batch(cmd...), FORM_PAGE
}

func (m *formPageModel) view() string {
	var output strings.Builder

	output.WriteString("\n")

	nameField := makeTextField("NAME", m.nameInput.View())
	addressField := makeTextField("ADDRESS", m.addressInput.View())
	diagnosisField := makeTextField("DIAGNOSIS", m.diagnosisInput.View())

	var dateField, doaField, dobField string
	switch m.fieldIndex {
	case 3:
		dateField = makeActiveDateField("DATE", m.date, m.dateInteractionIndex)
		doaField = makeInactiveDateField("DOA", m.dateOfAdmission)
		dobField = makeInactiveDateField("DOB", m.dateOfBirth)
	case 4:
		dateField = makeInactiveDateField("DATE", m.date)
		doaField = makeActiveDateField("DOA", m.dateOfAdmission, m.dateInteractionIndex)
		dobField = makeInactiveDateField("DOB", m.dateOfBirth)
	case 5:
		dateField = makeInactiveDateField("DATE", m.date)
		doaField = makeInactiveDateField("DOA", m.dateOfAdmission)
		dobField = makeActiveDateField("DOB", m.dateOfBirth, m.dateInteractionIndex)
	default:
		dateField = makeInactiveDateField("DATE", m.date)
		doaField = makeInactiveDateField("DOA", m.dateOfAdmission)
		dobField = makeInactiveDateField("DOB", m.dateOfBirth)
	}

	var numDaysField string
	if m.fieldIndex == 6 {
		numDaysField = lipgloss.JoinHorizontal(
			lipgloss.Center,
			formFieldNameStyle.Render("> DAYS"),
			numDaysFieldStyle.Render(strconv.Itoa(m.numDays)),
		)
	} else {
		numDaysField = lipgloss.JoinHorizontal(
			lipgloss.Center,
			formFieldNameStyle.Render("  DAYS"),
			numDaysFieldStyle.Render(strconv.Itoa(m.numDays)),
		)
	}

	var printBtn string
	if m.fieldIndex == 7 {
		printBtn = printBtnStyle.Render("> PRINT")
	} else {
		printBtn = printBtnStyle.Render("  PRINT")
	}

	var errMsg string
	if len(m.errMsg) > 0 {
		errMsg = errStyle.Render("[ERROR]: " + m.errMsg)
	}

	output.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		nameField,
		addressField,
		diagnosisField,
		dateField,
		doaField,
		dobField,
		numDaysField,
		printBtn,
		errMsg,
	))

	return output.String()
}

func (m *formPageModel) validateInput() (formData, error) {
	if len(m.nameInput.Value()) == 0 {
		return formData{}, fmt.Errorf("name field is empty")
	}
	if len(m.addressInput.Value()) == 0 {
		return formData{}, fmt.Errorf("address field is empty")
	}
	if len(m.diagnosisInput.Value()) == 0 {
		return formData{}, fmt.Errorf("diagnosis field is empty")
	}

	if m.date.Compare(m.dateOfBirth) < 0 {
		return formData{}, fmt.Errorf("dob is after date")
	}
	if m.date.Compare(m.dateOfAdmission) < 0 {
		return formData{}, fmt.Errorf("doa is after date")
	}
	if m.dateOfAdmission.Compare(m.dateOfBirth) < 0 {
		return formData{}, fmt.Errorf("dob is after doa")
	}

	return formData{
		name:            m.nameInput.Value(),
		address:         m.addressInput.Value(),
		diagnosis:       m.diagnosisInput.Value(),
		date:            m.date,
		dateOfAdmission: m.dateOfAdmission,

		dateOfBirth: m.dateOfBirth,
	}, nil
}

func (m *formPageModel) generatePrintPDFCmd(fd formData) tea.Cmd {
	return func() tea.Msg {
		s := store{}
		s.loadRecords()
		s.addRecord(fd)
		s.storeRecords()

		err := generatePDF("output.pdf", fd, m.numDays)
		if err == nil {
			printPDF("output.pdf")
		}

		return nil
	}
}

func makeTextField(fieldName, inputView string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		formFieldNameStyle.Render(fieldName),
		borderStyle.
			Render(inputView),
	)

}

func makeActiveDateField(fieldName string, date time.Time, dateInteractionIndex int) string {
	var indicators string
	switch dateInteractionIndex {
	case 0:
		indicators = "^^"
	case 1:
		indicators = "   ^^"
	default:
		indicators = "      ^^^^"
	}

	dateLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formFieldNameStyle.Render(fieldName),
		dateFieldStyle.Render(date.Format("02/01/2006")),
	)

	interactionLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formFieldNameStyle.Render(""),
		dateIndicatorStyle.Render(indicators),
	)

	return lipgloss.JoinVertical(lipgloss.Left, dateLine, interactionLine)
}

func makeInactiveDateField(fieldName string, date time.Time) string {
	dateLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formFieldNameStyle.Render(fieldName),
		dateFieldStyle.Render(date.Format("02/01/2006")),
	)

	return dateLine + "\n"
}
