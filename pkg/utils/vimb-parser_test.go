package utils

import (
	go_convert "github.com/advancemg/go-convert"
	"reflect"
	"testing"
)

func TestVimbResponse_Convert(t *testing.T) {
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
			fields:  fields{FilePath: "/Users/eminshakh/work/vimb-loader/s3-buckets/storage/vimb/vimb/GetBudgets/201904/budgets.gz"},
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
