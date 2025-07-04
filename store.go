package main

import (
	"encoding/csv"
	"os"
	"slices"
	"strings"
	"time"
)

type store struct {
	records []formData
}

func (s *store) storeRecords() error {
	file, err := os.Create("data.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Name", "Address", "Diagnosis", "Date", "Date of Admission", "Date of Birth"})
	if err != nil {
		return err
	}

	data := recordsToRows(s.records)
	for _, record := range data {
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) loadRecords() error {
	file, err := os.Open("data.csv")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	s.records = rowsToRecords(records[1:])
	s.sortRecords()

	return nil
}

func (s *store) getRecordsByName(name string) []formData {
	var output []formData
	for _, record := range s.records {
		if strings.Contains(strings.TrimSpace(strings.ToLower(record.name)), strings.TrimSpace(strings.ToLower(name))) {
			output = append(output, record)
		}
	}

	return output
}

func (s *store) addRecord(fd formData) {
	if len(s.records) >= 10 {
		s.records = s.records[1:]
	}

	index := s.getRecordIndex(fd.name)

	if index != -1 {
		s.records[index] = fd
	} else {
		s.records = append(s.records, fd)
	}
	s.sortRecords()
}

func (s *store) getRecordIndex(name string) int {
	index := -1

	for i, record := range s.records {
		if strings.TrimSpace(strings.ToLower(record.name)) == strings.TrimSpace(strings.ToLower(name)) {
			return i
		}
	}

	return index
}

func (s *store) sortRecords() {
	slices.SortStableFunc(s.records, func(a, b formData) int {
		return a.date.Compare(b.date)
	})
}

func recordsToRows(records []formData) [][]string {
	var output [][]string
	for _, record := range records {
		var fields []string

		fields = append(fields, record.name)
		fields = append(fields, record.address)
		fields = append(fields, record.diagnosis)

		fields = append(fields, record.date.Format("02/01/2006"))
		fields = append(fields, record.dateOfAdmission.Format("02/01/2006"))
		fields = append(fields, record.dateOfBirth.Format("02/01/2006"))

		output = append(output, fields)
	}
	return output
}

func rowsToRecords(rows [][]string) []formData {
	var output []formData

	for _, row := range rows {
		date, _ := time.Parse("02/01/2006", row[3])
		dateOfAdmission, _ := time.Parse("02/01/2006", row[4])
		dateOfBirth, _ := time.Parse("02/01/2006", row[5])

		record := formData{
			name:            row[0],
			address:         row[1],
			diagnosis:       row[2],
			date:            date,
			dateOfAdmission: dateOfAdmission,
			dateOfBirth:     dateOfBirth,
		}

		output = append(output, record)
	}
	return output
}
