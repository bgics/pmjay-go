package main

type fieldName int

const (
	NAME = iota
	ADDRESS1
	ADDRESS2
	ADDRESS3
	DATE
	AGE
	DATE_OF_BIRTH
	DAY_OF_ADMISSION
	DATE_OF_ADMISSION
	DIAGNOSIS
)

const (
	templateFileStr = "./assets/form_template.png"
	fontDirStr      = "./assets"

	a4Width  = 210
	a4Height = 297

	fieldPrintYOffset = 0.5
)

var fieldConfig = map[fieldName]struct {
	x        float64
	y        float64
	maxChars int
}{
	NAME:              {25.24, 44.53, 41},
	ADDRESS1:          {30.10, 53.54, 39},
	ADDRESS2:          {12.43, 62.55, 47},
	ADDRESS3:          {12.43, 71.56, 47},
	DATE:              {136.83, 44.53, 25},
	AGE:               {136, 53.54, 7},
	DATE_OF_BIRTH:     {152.26, 62.55, 18},
	DAY_OF_ADMISSION:  {159.75, 71.56, 15},
	DATE_OF_ADMISSION: {47.21, 80.57, 32},
	DIAGNOSIS:         {32.93, 89.58, 70},
}

var fontConfig = struct {
	familyStr string
	styleStr  string
	fileStr   string
	size      float64
}{
	"JetBrainsMono",
	"",
	"JetBrainsMono-Regular.json",
	11,
}
