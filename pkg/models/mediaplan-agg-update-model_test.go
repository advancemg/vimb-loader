package models

import "testing"

func TestMediaplanAggUpdateRequest_Update(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		Month        int64
		ChannelId    int64
		MediaplanId  int64
		AdvertiserId int64
		AgreementId  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "update agg mediaplan",
			fields: fields{
				Month:        201902,
				ChannelId:    1020582,
				MediaplanId:  14825353,
				AdvertiserId: 700068653,
				AgreementId:  81024,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MediaplanAggUpdateRequest{
				Month:        tt.fields.Month,
				ChannelId:    tt.fields.ChannelId,
				MediaplanId:  tt.fields.MediaplanId,
				AdvertiserId: tt.fields.AdvertiserId,
				AgreementId:  tt.fields.AgreementId,
			}
			if err := request.Update(); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
