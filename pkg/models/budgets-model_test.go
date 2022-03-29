package models

import (
	go_convert "github.com/advancemg/go-convert"
	"testing"
)

func TestGetBudgets_DataConfiguration(t *testing.T) {
	type fields struct {
		UnsortedMap go_convert.UnsortedMap
	}
	type args struct {
		s3Key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "load arch data",
			fields: fields{
				UnsortedMap: go_convert.UnsortedMap{},
			},
			args: args{
				s3Key: "adfadf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &GetBudgets{
				UnsortedMap: tt.fields.UnsortedMap,
			}
			if err := request.DataConfiguration(tt.args.s3Key); (err != nil) != tt.wantErr {
				t.Errorf("DataConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
