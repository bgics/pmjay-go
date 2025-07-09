package store

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bgics/pmjay-go/config"
	"github.com/bgics/pmjay-go/model"
)

var (
	CSVHeader = []string{"Name", "Address", "Diagnosis", "Date", "Date of Admission", "Date of Birth"}
)

type Store struct {
	records []model.FormData
	isValid bool
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) AddRecord(fd model.FormData) error {
	if !s.isValid {
		if err := s.loadRecords(); err != nil {
			return fmt.Errorf("cannot load records: %w", err)
		}
	}

	if len(s.records) >= 10 {
		s.records = s.records[1:]
	}

	index := s.getRecordIndex(fd.Name)

	if index != -1 {
		s.records[index] = fd
	} else {
		s.records = append(s.records, fd)
	}
	s.sortRecords()

	if err := s.storeRecords(); err != nil {
		s.isValid = false
		return fmt.Errorf("cannot save records: %w", err)
	}

	return nil
}

func (s *Store) GetRecordsByName(name string) ([]model.FormData, error) {
	if !s.isValid {
		if err := s.loadRecords(); err != nil {
			return nil, fmt.Errorf("cannot load records: %w", err)
		}
	}

	var output []model.FormData
	for _, record := range s.records {
		if strings.Contains(sanitizeString(record.Name), sanitizeString(name)) {
			output = append(output, record)
		}
	}

	return output, nil
}

func (s *Store) storeRecords() error {
	file, err := os.Create("data.csv")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(CSVHeader)
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

func (s *Store) loadRecords() error {
	file, err := os.Open("data.csv")
	if os.IsNotExist(err) {
		s.isValid = true
		return nil
	}
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	s.records, err = rowsToRecords(rows[1:])
	if err != nil {
		return err
	}

	s.sortRecords()

	s.isValid = true
	return nil
}

func (s *Store) getRecordIndex(name string) int {
	index := -1

	for i, record := range s.records {
		if sanitizeString(record.Name) == sanitizeString(name) {
			return i
		}
	}

	return index
}

func sanitizeString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func (s *Store) sortRecords() {
	slices.SortStableFunc(s.records, func(a, b model.FormData) int {
		return a.Date.Compare(b.Date)
	})
}

func recordsToRows(records []model.FormData) [][]string {
	var output [][]string
	for _, record := range records {
		fields := []string{
			record.Name,
			record.Address,
			record.Diagnosis,
			record.Date.Format(config.DateFormat),
			record.DateOfAdmission.Format(config.DateFormat),
			record.DateOfBirth.Format(config.DateFormat),
		}

		output = append(output, fields)
	}
	return output
}

func rowsToRecords(rows [][]string) ([]model.FormData, error) {
	var output []model.FormData

	for _, row := range rows {
		date, err := time.Parse(config.DateFormat, row[3])
		if err != nil {
			return nil, err
		}

		dateOfAdmission, err := time.Parse(config.DateFormat, row[4])
		if err != nil {
			return nil, err
		}

		dateOfBirth, err := time.Parse(config.DateFormat, row[5])
		if err != nil {
			return nil, err
		}

		record := model.FormData{
			Name:            row[0],
			Address:         row[1],
			Diagnosis:       row[2],
			Date:            date,
			DateOfAdmission: dateOfAdmission,
			DateOfBirth:     dateOfBirth,
		}

		output = append(output, record)
	}
	return output, nil
}
