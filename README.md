# PMJAY-Go 🏥

A simple CLI tool I made to help my dad with his daily patient reporting.

## Why I Built This

My dad is a doctor who treats patients under the PMJAY (Pradhan Mantri Jan Arogya Yojana) scheme. Every day, he had to fill out the same patient information forms repeatedly for each patient's daily report. It was tedious and time-consuming.

So I built this tool for him - you enter a patient's details once, tell it how many days you need reports for, and it generates all the forms automatically. It also saves patient information so you can quickly look up previous patients without retyping everything.

## What It Does

- Enter patient information once through a simple terminal interface
- Automatically generates PDF forms for multiple days
- Saves patient records locally (CSV file) so you can search and reuse them
- Can print the generated PDFs directly (if you have PDFtoPrinter.exe)

## Quick Start (Pre-built)

If you just want to use the tool without building from source:

1. Download the latest `.zip` file from the [Releases](https://github.com/bgics/pmjay-go/releases) page
2. Extract it anywhere on your computer
3. (Optional) If you want printing support, put `PDFtoPrinter.exe` in the same folder as `pmjay.exe`
4. Run `pmjay.exe`

## Building from Source

If you want to build it yourself:

### Requirements

- Go 1.24.4 or later
- Task runner (from [taskfile.dev](https://taskfile.dev/))
- PDFtoPrinter.exe (optional, for printing)
- Windows (this only works on Windows)

### Build Steps

1. Clone this repo:
```bash
git clone https://github.com/bgics/pmjay-go.git
cd pmjay-go
```

2. Download dependencies:
```bash
go mod download
```

3. (Optional) Put `PDFtoPrinter.exe` in the project folder (same directory as the .exe) if you want to print directly

4. Build it:
```bash
task build
```

This creates `pmjay.exe` and packages everything into a ZIP file.

To clean up build files:
```bash
task clean
```

## How to Use

Just run `pmjay.exe` and follow the prompts:

1. Choose "New Patient" to enter someone new, or "Search Records" to find a previous patient
2. Fill in the patient details (name, address, dates, diagnosis, etc.)
3. Enter how many days you need reports for
4. The tool generates a PDF with all the daily forms

The patient info gets saved automatically, so next time you can search for them by name instead of typing everything again.

## Project Structure

```
pmjay-go/
├── assets/           # Application icon, fonts, and form template
├── config/           # Configuration for PDF generation and form fields
├── internal/tui/     # Terminal UI components using Bubble Tea
├── model/            # Data models (patient information)
├── pdf/              # PDF generation logic
├── store/            # CSV-based patient data storage
└── main.go           # Application entry point
```

## Tech Stack

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - for the terminal UI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - for styling
- [gofpdf](https://github.com/phpdave11/gofpdf) - for generating PDFs
- [Task](https://taskfile.dev/) - for build automation

## License

This is a personal project I made for my dad. Feel free to use it if you find it helpful.

## A Few Notes

- Patient data is saved in a CSV file locally
- It keeps the last 10 patients for quick search
- Generated PDF is saved as `output.pdf`
- The form template and fonts are in the `assets` folder
