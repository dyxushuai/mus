package utils

import "time"

const Format = "2006-01-02T15:04:05"
const jsonFormat = `"` + Format + `"`

var fixedZone = time.FixedZone("", 0)

type Time time.Time


func New(t time.Time) Time {
	return Time(time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		0,
		fixedZone,
	))
}

func (it Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(it).Format(jsonFormat)), nil
}

func (it *Time) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(jsonFormat, string(data), fixedZone)
	if err == nil {
		*it = Time(t)
	}

	return err
}

func (it Time) String() string {
	return time.Time(it).String()
}
