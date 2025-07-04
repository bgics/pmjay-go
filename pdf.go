package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

type formData struct {
	name            string
	date            time.Time
	address         string
	dateOfBirth     time.Time
	dateOfAdmission time.Time
	diagnosis       string
}

type textLine struct {
	text string
	x    float64
	y    float64
}

func printPDF(filename string) {
	cmd := exec.Command(".\\PDFtoPrinter.exe", filename)
	cmd.Run()
}

func generatePDF(outFileStr string, fd formData, numDays int) error {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, fontDirStr)

	pdf.AddFont(fontConfig.familyStr, fontConfig.styleStr, fontConfig.fileStr)
	pdf.SetFont(fontConfig.familyStr, fontConfig.styleStr, fontConfig.size)

	for range numDays {
		pdf.AddPage()
		pdf.Image(templateFileStr, 0, 0, a4Width, a4Height, false, "", 0, "")

		textLines, err := convertToTextLines(fd)
		if err != nil {
			return err
		}

		for _, line := range textLines {
			pdf.Text(line.x, line.y-fieldPrintYOffset, line.text)
		}

		fd.date = fd.date.AddDate(0, 0, 1)
	}

	return pdf.OutputFileAndClose(outFileStr)
}

func convertToTextLines(fd formData) ([]textLine, error) {
	if fd.date.Compare(fd.dateOfAdmission) < 0 {
		return nil, fmt.Errorf("date is before date of admission")
	}

	if fd.date.Compare(fd.dateOfBirth) < 0 {
		return nil, fmt.Errorf("date is before date of birth")
	}

	if fd.dateOfAdmission.Compare(fd.dateOfBirth) < 0 {
		return nil, fmt.Errorf("date of admission is before date of birth")
	}

	var output []textLine

	output = append(output, makeTextFieldTextLine(fd.name, NAME))
	output = append(output, makeTextFieldTextLine(fd.diagnosis, DIAGNOSIS))
	output = append(output, makeAddressTextLines(fd.address)...)
	output = append(output, makeDateTextLine(fd.date, DATE))
	output = append(output, makeDateTextLine(fd.dateOfBirth, DATE_OF_BIRTH))
	output = append(output, makeDateTextLine(fd.dateOfAdmission, DATE_OF_ADMISSION))
	output = append(output, makeDayOfAdmissionTextLine(fd.date, fd.dateOfAdmission))
	output = append(output, makeAgeTextLine(fd.date, fd.dateOfBirth))

	return output, nil
}

func makeTextFieldTextLine(fieldString string, cfgKey fieldName) textLine {
	cfg := fieldConfig[cfgKey]

	return textLine{
		text: trim(fieldString, cfg.maxChars),
		x:    cfg.x,
		y:    cfg.y,
	}
}

func makeAddressTextLines(address string) []textLine {
	var output []textLine

	cfgKeys := [3]fieldName{ADDRESS1, ADDRESS2, ADDRESS3}

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

func makeDateTextLine(date time.Time, cfgKey fieldName) textLine {
	cfg := fieldConfig[cfgKey]

	dateString := date.Format("02/01/2006")
	return textLine{
		text: dateString,
		x:    cfg.x,
		y:    cfg.y,
	}
}

func makeDayOfAdmissionTextLine(date, dateOfAdmission time.Time) textLine {
	cfg := fieldConfig[DAY_OF_ADMISSION]

	numDays := int(date.Sub(dateOfAdmission).Hours()/24) + 1
	fieldString := fmt.Sprintf("DAY %d", numDays)

	return textLine{
		text: fieldString,
		x:    cfg.x,
		y:    cfg.y,
	}
}

func makeAgeTextLine(date, dateOfBirth time.Time) textLine {
	cfg := fieldConfig[AGE]

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
		x:    cfg.x,
		y:    cfg.y,
	}
}

func trim(fieldValue string, max int) string {
	if len(fieldValue) > max {
		return fieldValue[:max]
	}

	return fieldValue
}
