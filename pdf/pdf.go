package pdf

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bgics/pmjay-go/config"
	"github.com/bgics/pmjay-go/model"
	"github.com/phpdave11/gofpdf"
)

type textLine struct {
	text string
	x    float64
	y    float64
}

func PrintPDF(filename string) error {
	cmd := exec.Command(config.ExeName, filename)
	return cmd.Run()
}

func GeneratePDF(outFileStr string, fd model.FormData, numDays int) error {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, config.FontDirStr)

	pdf.AddFont(config.FontConfig.FamilyStr, config.FontConfig.StyleStr, config.FontConfig.FileStr)
	pdf.SetFont(config.FontConfig.FamilyStr, config.FontConfig.StyleStr, config.FontConfig.Size)

	for range numDays {
		pdf.AddPage()
		pdf.Image(config.TemplateFileStr, 0, 0, config.A4WidthMM, config.A4HeightMM, false, "", 0, "")

		textLines, err := convertToTextLines(fd)
		if err != nil {
			return err
		}

		for _, line := range textLines {
			pdf.Text(line.x, line.y-config.FieldYOffset, line.text)
		}

		fd.Date = fd.Date.AddDate(0, 0, 1)
	}

	return pdf.OutputFileAndClose(outFileStr)
}

func convertToTextLines(fd model.FormData) ([]textLine, error) {
	if fd.Date.Compare(fd.DateOfAdmission) < 0 {
		return nil, fmt.Errorf("date is before date of admission")
	}

	if fd.Date.Compare(fd.DateOfBirth) < 0 {
		return nil, fmt.Errorf("date is before date of birth")
	}

	if fd.DateOfAdmission.Compare(fd.DateOfBirth) < 0 {
		return nil, fmt.Errorf("date of admission is before date of birth")
	}

	output := []textLine{
		makeTextFieldTextLine(fd.Name, config.NAME),
		makeTextFieldTextLine(fd.Diagnosis, config.DIAGNOSIS),
		makeDateTextLine(fd.Date, config.DATE),
		makeDateTextLine(fd.DateOfBirth, config.DATE_OF_BIRTH),
		makeDateTextLine(fd.DateOfAdmission, config.DATE_OF_ADMISSION),
		makeDayOfAdmissionTextLine(fd.Date, fd.DateOfAdmission),
		makeAgeTextLine(fd.Date, fd.DateOfBirth),
	}

	output = append(output, makeAddressTextLines(fd.Address)...)

	return output, nil
}

func makeTextFieldTextLine(fieldString string, cfgKey config.FieldName) textLine {
	cfg := config.FieldConfig[cfgKey]

	return textLine{
		text: trim(fieldString, cfg.MaxChars),
		x:    cfg.X,
		y:    cfg.Y,
	}
}

func makeAddressTextLines(address string) []textLine {
	var output []textLine

	cfgKeys := [3]config.FieldName{config.ADDRESS1, config.ADDRESS2, config.ADDRESS3}

	for i := range 3 {
		line := makeTextFieldTextLine(address, cfgKeys[i])
		output = append(output, line)
		if len(address) <= len(line.text) {
			return output
		}
		address = strings.TrimPrefix(address, line.text)
	}

	return output
}

func makeDateTextLine(date time.Time, cfgKey config.FieldName) textLine {
	cfg := config.FieldConfig[cfgKey]

	dateString := date.Format("02/01/2006")
	return textLine{
		text: dateString,
		x:    cfg.X,
		y:    cfg.Y,
	}
}

func makeDayOfAdmissionTextLine(date, dateOfAdmission time.Time) textLine {
	cfg := config.FieldConfig[config.DAY_OF_ADMISSION]

	numDays := int(date.Sub(dateOfAdmission).Hours()/24) + 1
	fieldString := fmt.Sprintf("DAY %d", numDays)

	return textLine{
		text: fieldString,
		x:    cfg.X,
		y:    cfg.Y,
	}
}

func makeAgeTextLine(date, dateOfBirth time.Time) textLine {
	cfg := config.FieldConfig[config.AGE]

	age := int(date.Sub(dateOfBirth).Hours()/24) + 1
	var suffix string
	if age > 1 {
		suffix = "DAYS"
	} else {
		suffix = "DAY"
	}
	fieldString := fmt.Sprintf("%d %s", age, suffix)

	return textLine{
		text: fieldString,
		x:    cfg.X,
		y:    cfg.Y,
	}
}

func trim(fieldValue string, max int) string {
	if len(fieldValue) > max {
		return fieldValue[:max]
	}

	return fieldValue
}
