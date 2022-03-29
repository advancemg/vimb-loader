package utils

import (
	"testing"
)

func TestGetDaysFromYearMonth(t *testing.T) {
	type args struct {
		yearMonth string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get days by month",
			args: args{
				yearMonth: "202002",
			},
			wantErr: false,
		},
		{
			name: "get days by month",
			args: args{
				yearMonth: "20200",
			},
			wantErr: true,
		},
		{
			name: "get days by month",
			args: args{
				yearMonth: "20201231",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDaysFromYearMonth(tt.args.yearMonth)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysFromYearMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetPeriodFromYearMonths(t *testing.T) {
	type args struct {
		beginMonth string
		endMonth   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get month by period",
			args: args{
				beginMonth: "202203",
				endMonth:   "202205",
			},
			wantErr: false,
		},
		{
			name: "get month by period",
			args: args{
				beginMonth: "202204",
				endMonth:   "202203",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetPeriodFromYearMonths(tt.args.beginMonth, tt.args.endMonth)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPeriodFromYearMonths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetDaysByPeriod(t *testing.T) {
	type args struct {
		beginDay string
		endDay   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get days by period",
			args: args{
				beginDay: "20200201",
				endDay:   "20200229",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDaysByPeriod(tt.args.beginDay, tt.args.endDay)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysByPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
