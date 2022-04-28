package models

import (
	"github.com/advancemg/badgerhold"
	"testing"
)

func TestMediaplanBadgerQuery_Find(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	type args struct {
		filter *badgerhold.Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "find where",
			args: args{
				filter: badgerhold.Where("Month").Eq(201903),
			},
			wantErr: false,
		},
		{
			name: "find where is nil",
			args: args{
				filter: badgerhold.Where("OrdBegDate").IsNil(),
			},
			wantErr: false,
		},
		{
			name: "find where is nil and other field",
			args: args{
				filter: badgerhold.
					Where("OrdBegDate").IsNil().
					And("ChannelId").Eq(1020269),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &MediaplanBadgerQuery{}
			var result []Mediaplan
			if err := query.Find(&result, tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
			}
			println(len(result))
		})
	}
}
