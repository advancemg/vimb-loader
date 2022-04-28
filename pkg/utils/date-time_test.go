package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDaysFromMonth(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type args struct {
		year  int
		month time.Month
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "days",
			args: args{
				year:  2020,
				month: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDaysFromMonth(tt.args.year, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysFromMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestGetDaysFromYearMonth(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type args struct {
		yearMonth string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "days",
			args:    args{"201908"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDaysFromYearMonth(tt.args.yearMonth)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysFromYearMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			startDay := fmt.Sprintf("%v", got[0].Format(time.RFC3339))
			fmt.Println(startDay[0 : len(startDay)-1])
			endDay := fmt.Sprintf("%v", got[len(got)-1].Format(time.RFC3339))
			fmt.Println(endDay[0 : len(endDay)-1])
		})
	}
}

func TestGetWeekDayByYearMonth(t *testing.T) {
	for i := 0; i < 12; i++ {
		yearMonth := 201901 + i
		month, _ := GetWeekDayByYearMonth(yearMonth)
		for i, t := range month {
			fmt.Println(i, t)
		}
	}
}
