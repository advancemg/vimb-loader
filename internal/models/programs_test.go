package models

import "testing"

func TestProgramUpdateRequest_Update(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestPrograms",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &ProgramUpdateRequest{}
			if err := request.Update(); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
