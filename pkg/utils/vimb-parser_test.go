package utils

import (
	go_convert "github.com/advancemg/go-convert"
	"reflect"
	"testing"
)

func TestVimbResponse_Convert(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type fields struct {
		FilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *go_convert.UnsortedMap
		wantErr bool
	}{
		{
			name:    "testConvert",
			fields:  fields{FilePath: "../../dev-test-data/budgets.gz"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &VimbResponse{
				FilePath: tt.fields.FilePath,
			}
			got, err := req.Convert("BudgetList")
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
