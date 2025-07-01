package main

import "time"

type FormData struct {
	Name            string
	Date            time.Time
	Address         string
	DateOfBirth     time.Time
	DateOfAdmission time.Time
	Diagnosis       string
}
