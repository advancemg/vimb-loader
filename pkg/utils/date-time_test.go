package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDaysFromMonth(t *testing.T) {
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
				year:  2022,
				month: 4,
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
