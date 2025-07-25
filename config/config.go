package config

type FieldName int

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
	TemplateFileStr = "./assets/form_template.png"
	FontDirStr      = "./assets"

	A4WidthMM  = 210
	A4HeightMM = 297

	FieldYOffset = 0.5

	ExeName    = ".\\PDFtoPrinter.exe"
	DateFormat = "02/01/2006"

	OutputFileName = "output.pdf"

	GenderStrLen = 3
)

var FieldConfig = map[FieldName]struct {
	X        float64
	Y        float64
	MaxChars int
}{
	NAME:              {25.24, 44.53, 41 - GenderStrLen},
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

var FontConfig = struct {
	FamilyStr string
	StyleStr  string
	FileStr   string
	Size      float64
}{
	"JetBrainsMono",
	"",
	"JetBrainsMono-Bold.json",
	11,
}
