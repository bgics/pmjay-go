package model

import "time"

type FormData struct {
	Name            string
	Address         string
	Diagnosis       string
	Date            time.Time
	DateOfBirth     time.Time
	DateOfAdmission time.Time
}
