package view

import (
	"github.com/bgics/pmjay-go/config"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
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
