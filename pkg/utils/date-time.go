package utils

import (
	"strconv"
	"time"
)

type YearMonth struct {
	Month       time.Month `json:"month"`
	Year        int        `json:"year"`
	IntValue    int
	ValueString string
}

func GetActualDays() ([]time.Time, error) {
	actual := time.Now().UTC()
	actualYear := actual.Year()
	actualMonth := actual.Month()
	actualDay := actual.Day()
	begin := time.Date(actualYear, actualMonth, actualDay, 0, 0, 0, 0, time.UTC)
	end := time.Date(actualYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	var result []time.Time
	for begin.Unix() <= end.Unix() {
		result = append(result, begin)
		begin = begin.AddDate(0, 0, 1)
	}
	return result, nil
}

func GetActualMonths() ([]YearMonth, error) {
	actual := time.Now().UTC()
	actualYear := actual.Year()
	actualMonth := actual.Month()
	var result []YearMonth
	begin := time.Date(actualYear, actualMonth, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(actualYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	for begin.Unix() < end.Unix() {
		yearMonth := begin.Format(`200601`)
		yearMonthInt, err := strconv.Atoi(yearMonth)
		if err != nil {
			return nil, err
		}
		result = append(result, YearMonth{
			Month:       begin.Month(),
			Year:        begin.Year(),
			IntValue:    yearMonthInt,
			ValueString: yearMonth,
		})
		begin = begin.AddDate(0, 1, 0)
	}
	return result, nil
}
