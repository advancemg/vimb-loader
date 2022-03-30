package utils

import (
	"errors"
	"fmt"
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

func GetDaysFromMonth(year int, month time.Month) ([]string, error) {
	var days []string
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	for i := firstOfMonth.Day(); i <= lastOfMonth.Day(); i++ {
		if i <= 9 {
			days = append(days, fmt.Sprintf("0%d", i))
		} else {
			days = append(days, fmt.Sprintf("%d", i))
		}
	}
	return days, nil
}

func GetDaysFromYearMonth(yearMonth string) ([]time.Time, error) {
	if len(yearMonth) < 6 || len(yearMonth) > 6 {
		return nil, errors.New("not year month format [200012]")
	}
	year := Int(yearMonth[:4])
	month := Int(yearMonth[4:6])
	actualYear := year
	actualMonth := time.Month(month)
	begin := time.Date(actualYear, actualMonth, 1, 0, 0, 0, 0, time.UTC)
	end := begin.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
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

func GetActualStartEndDays(month int) (string, string) {
	days, _ := GetActualDays()
	var tmpMonth string
	var startDay string
	var endDay string
	for i := 0; i < len(days); i++ {
		if month == int(days[i].Month()) {
			nextDay := i + 1
			if days[i].Month().String() != tmpMonth {
				if days[i].Day() <= 9 {
					startDay = fmt.Sprintf("%d%d%s", days[i].Year(), int(days[i].Month()), fmt.Sprintf("0%d", days[i].Day()))
				} else {
					startDay = fmt.Sprintf("%d%d%d", days[i].Year(), int(days[i].Month()), days[i].Day())
				}
			}
			if nextDay < len(days) {
				if days[i].Month().String() != days[nextDay].Month().String() {
					endDay = fmt.Sprintf("%d%d%d", days[i].Year(), int(days[i].Month()), days[i].Day())
				}
				if (len(days) - nextDay) == 1 {
					endDay = fmt.Sprintf("%d%d%d", days[nextDay].Year(), int(days[nextDay].Month()), days[nextDay].Day())
				}
			}
			tmpMonth = days[i].Month().String()
		}
	}
	return startDay, endDay
}

func GetActualStartEndDaysForTest(year, month int) (string, string) {
	days, _ := GetActualDays()
	var tmpMonth string
	var startDay string
	var endDay string
	for i := 0; i < len(days); i++ {
		if month == int(days[i].Month()) {
			nextDay := i + 1
			if days[i].Month().String() != tmpMonth {
				if days[i].Day() <= 9 {
					startDay = fmt.Sprintf("%d%d%s", year, int(days[i].Month()), fmt.Sprintf("0%d", days[i].Day()))
				} else {
					startDay = fmt.Sprintf("%d%d%d", year, int(days[i].Month()), days[i].Day())
				}
			}
			if nextDay < len(days) {
				if days[i].Month().String() != days[nextDay].Month().String() {
					endDay = fmt.Sprintf("%d%d%d", year, int(days[i].Month()), days[i].Day())
				}
				if (len(days) - nextDay) == 1 {
					endDay = fmt.Sprintf("%d%d%d", year, int(days[nextDay].Month()), days[nextDay].Day())
				}
			}
			tmpMonth = days[i].Month().String()
		}
	}
	return startDay, endDay
}
