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

func GetDaysFromYearMonthInt(yearMonth int) ([]time.Time, error) {
	return GetDaysFromYearMonth(fmt.Sprintf("%d", yearMonth))
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

func GetPeriodFromYearMonths(beginMonth string, endMonth string) ([]YearMonth, error) {
	if len(beginMonth) < 6 || len(beginMonth) > 6 {
		return nil, errors.New("begin - not year month format [200012]")
	}
	if len(endMonth) < 6 || len(endMonth) > 6 {
		return nil, errors.New("endMonth - not year month format [200012]")
	}
	beginYear := Int(beginMonth[:4])
	beginYearMonth := Int(beginMonth[4:6])
	endYear := Int(endMonth[:4])
	endYearMonth := Int(endMonth[4:6])
	if (beginYear + beginYearMonth) > (endYear + endYearMonth) {
		return nil, errors.New("begin month large end month")
	}
	var result []YearMonth
	begin := time.Date(beginYear, time.Month(beginYearMonth), 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(endYear, time.Month(endYearMonth), 1, 0, 0, 0, 0, time.UTC)
	end := endTime.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
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

func GetDaysByPeriod(beginDay string, endDay string) ([]time.Time, error) {
	if len(beginDay) < 8 || len(beginDay) > 8 {
		return nil, errors.New("beginDay - not format [20001201]")
	}
	if len(endDay) < 8 || len(endDay) > 8 {
		return nil, errors.New("endDay - not format [20001201]")
	}
	beginYear := Int(beginDay[:4])
	beginYearMonth := Int(beginDay[4:6])
	beginYearDay := Int(beginDay[6:8])
	endYear := Int(endDay[:4])
	endYearMonth := Int(endDay[4:6])
	endYearDay := Int(endDay[6:8])
	if (beginYear + beginYearMonth + beginYearDay) > (endYear + endYearMonth + endYearDay) {
		return nil, errors.New("begin day large end day")
	}
	var result []time.Time
	begin := time.Date(beginYear, time.Month(beginYearMonth), beginYearDay, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, time.Month(endYearMonth), endYearDay, 0, 0, 0, 0, time.UTC)
	for begin.Unix() <= end.Unix() {
		result = append(result, begin)
		begin = begin.AddDate(0, 0, 1)
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

func GetWeekDayByYearMonth(yearMonth int) (map[int]time.Time, error) {
	weekDays := map[int]time.Time{}
	monthInt, err := GetDaysFromYearMonthInt(yearMonth)
	if err != nil {
		return nil, err
	}
	if monthInt[0].Weekday() != time.Sunday && monthInt[0].Weekday() != time.Saturday {
		_, week := monthInt[0].ISOWeek()
		weekDays[week] = monthInt[0]
	}
	for _, date := range monthInt {
		if date.Weekday() == time.Monday {
			_, weekDay := date.ISOWeek()
			weekDays[weekDay] = date
		}
	}
	return weekDays, nil
}
