package utils

import "time"

type JsonTime struct {
	time.Time
}

func (j JsonTime) format() string {
	return j.Time.Unix()
}

func (j JsonTime) MarshalText() ([]byte, error) {
	return []byte(j.format()), nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

