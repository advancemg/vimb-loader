package models

import "testing"

func TestAdvertiserUpdateRequest_loadFromFile(t *testing.T) {
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
			name:    "loadFromFile-Advertiser",
			fields:  fields{"../../dev-test-data/advertisers.gz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &AdvertiserUpdateRequest{
				S3Key: tt.fields.S3Key,
			}
			if err := request.loadFromFile(); (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
