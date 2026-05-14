package utils

import "time"

func ParseDOB(dob string) (time.Time, error) {
	timeLayout := "2006-01-02"
	return time.Parse(timeLayout, dob)
}
