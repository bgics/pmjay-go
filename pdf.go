package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

type TextLine struct {
	Text string
	X    float64
	Y    float64
}

func GeneratePDF(outFileStr string, formData FormData, numDays int) error {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, FontDirStr)

	pdf.AddFont(FontConfig.FamilyStr, FontConfig.StyleStr, FontConfig.FileStr)
	pdf.SetFont(FontConfig.FamilyStr, FontConfig.StyleStr, FontConfig.Size)

	for range numDays {
		pdf.AddPage()
		pdf.Image(TemplateFileStr, 0, 0, A4Width, A4Height, false, "", 0, "")

		textLines, err := convertToTextLines(formData)
		if err != nil {
			return err
		}

		for _, line := range textLines {
			pdf.Text(line.X, line.Y-FieldPrintYOffset, line.Text)
		}

		formData.Date = formData.Date.Add(24 * time.Hour)
	}

	return pdf.OutputFileAndClose(outFileStr)
}

func convertToTextLines(formData FormData) ([]TextLine, error) {
	if !formData.Date.After(formData.DateOfAdmission) {
		return nil, fmt.Errorf("date is before date of admission")
	}

	if !formData.Date.After(formData.DateOfBirth) {
		return nil, fmt.Errorf("date is before date of birth")
	}

	if !formData.DateOfAdmission.After(formData.DateOfBirth) {
		return nil, fmt.Errorf("date of admission is before date of birth")
	}

	var output []TextLine

	output = append(output, makeTextFieldTextLine(formData.Name, NAME))
	output = append(output, makeTextFieldTextLine(formData.Diagnosis, DIAGNOSIS))
	output = append(output, makeAddressTextLines(formData.Address)...)
	output = append(output, makeDateTextLine(formData.Date, DATE))
	output = append(output, makeDateTextLine(formData.DateOfBirth, DATE_OF_BIRTH))
	output = append(output, makeDateTextLine(formData.DateOfAdmission, DATE_OF_ADMISSION))
	output = append(output, makeDayOfAdmissionTextLine(formData.Date, formData.DateOfAdmission))
	output = append(output, makeAgeTextLine(formData.Date, formData.DateOfBirth))

	return output, nil
}

func makeTextFieldTextLine(fieldString string, cfgKey FieldName) TextLine {
	cfg := FieldConfig[cfgKey]

	return TextLine{
		Text: trim(fieldString, cfg.MaxChars),
		X:    cfg.X,
		Y:    cfg.Y,
	}
}

func makeAddressTextLines(address string) []TextLine {
	var output []TextLine

	cfgKeys := [3]FieldName{ADDRESS1, ADDRESS2, ADDRESS3}

	for i := range 3 {
		line := makeTextFieldTextLine(address, cfgKeys[i])
		output = append(output, line)
		if len(address) <= len(line.Text) {
			return output
		}
		address = strings.TrimPrefix(address, line.Text)
	}

	return output
}

func makeDateTextLine(date time.Time, cfgKey FieldName) TextLine {
	cfg := FieldConfig[cfgKey]

	dateString := date.Format("02/01/2006")
	return TextLine{
		Text: dateString,
		X:    cfg.X,
		Y:    cfg.Y,
	}
}

func makeDayOfAdmissionTextLine(date, dateOfAdmission time.Time) TextLine {
	cfg := FieldConfig[DAY_OF_ADMISSION]

	numDays := int(date.Sub(dateOfAdmission).Hours()/24) + 1
	fieldString := fmt.Sprintf("DAY %d", numDays)

	return TextLine{
		Text: fieldString,
		X:    cfg.X,
		Y:    cfg.Y,
	}
}

func makeAgeTextLine(date, dateOfBirth time.Time) TextLine {
	cfg := FieldConfig[AGE]

	age := int(date.Sub(dateOfBirth).Hours()/24) + 1
	var suffix string
	if age > 1 {
		suffix = "DAYS"
	} else {
		suffix = "DAY"
	}
	fieldString := fmt.Sprintf("%d %s", age, suffix)

	return TextLine{
		Text: fieldString,
		X:    cfg.X,
		Y:    cfg.Y,
	}
}

func trim(fieldValue string, max int) string {
	if len(fieldValue) > max {
		return fieldValue[:max]
	}

	return fieldValue
}
