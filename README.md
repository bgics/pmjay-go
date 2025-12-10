# PMJAY-Go 🏥

A terminal-based CLI tool designed to streamline daily reporting for PMJAY (Pradhan Mantri Jan Arogya Yojana) scheme patients.

## About

This tool was created to help doctors manage the tedious task of filling out daily patient reports under the PMJAY scheme. Instead of manually filling forms for each patient every day, enter the patient details once and let the tool automatically generate forms for all required days.

### Key Features

- **📝 Quick Patient Entry**: Input patient information through an intuitive terminal interface
- **🔍 Patient Search**: Search and retrieve previously entered patient records
- **📄 Automated PDF Generation**: Automatically generates daily report forms for multiple days
- **💾 Patient Data Storage**: Stores patient entries in CSV format for easy access
- **🖨️ Direct Printing**: Integrates with PDFtoPrinter for seamless printing

## Prerequisites

- **Go 1.24.4** or later
- **Task** - Task runner (install from [taskfile.dev](https://taskfile.dev/))
- **PDFtoPrinter.exe** - For printing functionality (Windows)
- **Windows OS** - This tool is designed exclusively for Windows

## Installation

1. Clone the repository:
```bash
git clone https://github.com/bgics/pmjay-go.git
cd pmjay-go
```

2. Install Go dependencies:
```bash
go mod download
```

3. Place `PDFtoPrinter.exe` in the project root directory (for printing support)

## Building

Build the application using Task:

```bash
task build
```

This will:
- Embed the application icon
- Compile the executable (`pmjay.exe`)
- Create a distributable ZIP package with all necessary assets

To clean build artifacts:
```bash
task clean
```

## Usage

Run the application:

```bash
./pmjay.exe
```

### Workflow

1. **Start Screen**: Choose between "New Patient" or "Search Records"
2. **New Patient**: Fill in patient details including:
   - Name
   - Address
   - Date of Birth
   - Date of Admission
   - Diagnosis
   - Gender
   - Number of days for report generation
3. **Search Records**: Look up previously entered patients by name
4. **PDF Generation**: Forms are automatically generated and can be printed directly

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

## Technologies Used

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[gofpdf](https://github.com/phpdave11/gofpdf)** - PDF generation
- **[Task](https://taskfile.dev/)** - Modern task runner

## License

This project is a personal tool created to help streamline medical documentation workflow.

## Notes

- Patient records are stored locally in CSV format
- The tool maintains up to 10 recent patient records for quick access
- PDF output file is named `output.pdf` by default
- Form template and fonts are bundled in the `assets` directory

---

*Made with ❤️ to make daily medical reporting easier*
