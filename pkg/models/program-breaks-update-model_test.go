package models

import (
	"fmt"
	"testing"
	"time"
)

func TestProgramBreaksUpdateRequest_loadFromFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		S3Key string
		Month string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "loadFromFile-ProgramBreaks",
			fields:  fields{"../../dev-test-data/networks.gz", "201903"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &ProgramBreaksUpdateRequest{
				S3Key: tt.fields.S3Key,
				Month: tt.fields.Month,
			}
			start := time.Now()
			if err := request.loadFromFile(); (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(time.Since(start))
		})
	}
}
