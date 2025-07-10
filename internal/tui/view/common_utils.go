package view

import (
	"time"

	"github.com/bgics/pmjay-go/config"
	"github.com/bgics/pmjay-go/internal/tui"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

const (
	TextInputWidth = 40
)

func makeTextInput(focus bool, cfgKeys ...config.FieldName) textinput.Model {
	t := textinput.New()
	t.Cursor.SetMode(cursor.CursorBlink)
	t.Prompt = " "
	t.Width = TextInputWidth
	if focus {
		t.Focus()
	}

	maxChars := 0
	for _, key := range cfgKeys {
		maxChars += config.FieldConfig[key].MaxChars
	}

	t.CharLimit = maxChars

	return t
}

func makeDateInput() datepicker.Model {
	d := datepicker.New(time.Now())
	defaultStyle := datepicker.DefaultStyles()

	defaultStyle.FocusedText = defaultStyle.FocusedText.
		Foreground(tui.DatePickerHighlightColor)

	d.Styles = defaultStyle

	return d
}
