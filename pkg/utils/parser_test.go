package utils

import (
	"reflect"
	"testing"
)

func TestStringI(t *testing.T) {
	value := "Value1"
	data := map[string]interface{}{
		"Field1": value,
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "set map",
			args: args{
				data["Field1"],
			},
			want: &value,
		},
		{
			name: "set map",
			args: args{
				data["Field2"],
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringI(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringI() = %v, want %v", got, tt.want)
			}
		})
	}
}
