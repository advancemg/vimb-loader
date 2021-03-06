package models

import "testing"

func TestChannelsUpdateRequest_loadFromFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		S3Key              string
		SellingDirectionID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "loadFromFile-Channels",
			fields:  fields{"../../dev-test-data/channels.gz", "23"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &ChannelsUpdateRequest{
				S3Key:              tt.fields.S3Key,
				SellingDirectionID: tt.fields.SellingDirectionID,
			}
			if err := request.loadFromFile(); (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
