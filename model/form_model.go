package model

import "time"

type Gender string

const (
	Male   = "M"
	Female = "F"
)

type FormData struct {
	Name            string
	Address         string
	Diagnosis       string
	Gender          Gender
	Date            time.Time
	DateOfBirth     time.Time
	DateOfAdmission time.Time
}
