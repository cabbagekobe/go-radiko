package internal

import (
	"time"
)

const (
	dateLayout     = "20060102"
	datetimeLayout = "20060102150405"

	// Always use Asia/Tokyo timezone.
	tz = "Asia/Tokyo"
)

var location *time.Location

func init() {
	var err error

	location, err = time.LoadLocation(tz)
	if err != nil {
		// panic happens in init function.
		panic(err)
	}
}

// Date returns a textual representation of the time value
// formatted in dateLayout.
func Date(t time.Time) string {
	localTime := t.In(location)
	return localTime.Format(dateLayout)
}

// Datetime returns a textual representation of the time value
// formatted in datetimeLayout.
func Datetime(t time.Time) string {
	localTime := t.In(location)
	return localTime.Format(datetimeLayout)
}
