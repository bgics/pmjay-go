package main

import (
	"fmt"
	"log"
	// "log"
	"strings"

	"math/rand"

	"github.com/phpdave11/gofpdf"
)

type Field struct {
	FieldName string
	Lines     []TextLine
}

type TextLine struct {
	StartX float64
	EndX   float64
	Y      float64
}

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	fields = [...]Field{
		{"Name", []TextLine{
			{25.24, 123, 44.53},
		}},
		{"Address", []TextLine{
			{30.10, 123, 53.54},
			{12.43, 123, 62.55},
			{12.43, 123, 71.56},
		}},
		{"Date", []TextLine{
			{136.83, 196.05, 44.53},
		}},
		{"Age", []TextLine{
			{136.83, 152.72, 53.54},
		}},
		{"Date of Birth", []TextLine{
			{152.26, 196.05, 62.55},
		}},
		{"Day of Admission", []TextLine{
			{159.75, 196.05, 71.56},
		}},
		{"Date of Admission", []TextLine{
			{47.21, 123, 80.57},
		}},
		{"Diagnosis", []TextLine{
			{32.93, 196.05, 89.58},
		}},
	}
)

func getMax(pdf *gofpdf.Fpdf, limit float64) int {
	var builder strings.Builder
	for {
		beforeWrite := builder.String()
		builder.WriteByte(charset[rand.Intn(len(charset))])

		if pdf.GetStringWidth(builder.String()) > limit {
			return len(beforeWrite)
		}
	}

}

func genRandStr(numChars int) string {
	var builder strings.Builder

	for range numChars {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}

	return builder.String()
}

func main() {
	pdf := gofpdf.New("P", "mm", "A4", ".")
	pdf.AddPage()

	pdf.AddFont("JBM", "", "JetBrainsMono-Regular.json")
	pdf.SetFont("JBM", "", 12)

	pdf.Image("form_template.png", 0, 0, 210, 297, false, "", 0, "")

	for _, field := range fields {
		fmt.Printf("Field: %s\n", field.FieldName)
		for i, line := range field.Lines {
			lineMax := getMax(pdf, line.EndX-line.StartX)
			fmt.Printf("\t Line[%d] Max is: %d\n", i, lineMax)

			pdf.Text(line.StartX, line.Y, genRandStr(lineMax))
		}
	}

	err := pdf.OutputFileAndClose("output.pdf")

	if err != nil {
		log.Fatal(err)
	}
}
