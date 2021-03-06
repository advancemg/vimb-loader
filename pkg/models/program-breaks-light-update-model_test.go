package models

import (
	"testing"
)

func TestProgramBreaksLightUpdateRequest_loadFromFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		S3Key string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "loadFromFile-ProgramBreaksLight",
			fields:  fields{"../../dev-test-data/networks_light_new.gz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &ProgramBreaksLightUpdateRequest{
				S3Key: tt.fields.S3Key,
			}
			if err := request.loadFromFile(); (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
