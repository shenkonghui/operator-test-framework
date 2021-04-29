package util

import (
	"operator-test-framework/pkg/api"
	"reflect"
	"testing"
)

func Test_ConvertStrToPara(t *testing.T) {
	type args struct {
		str  string
		para []api.Parameter
	}
	tests := []struct {
		name string
		args args
		want []api.Parameter
	}{
		{
			args: args{
				str: "name1=a",
			},
			want: []api.Parameter{
				{
					Name:  "name1",
					Value: "a",
				},
			},
		},
		{
			args: args{
				str: "name1=a,name2=b",
				para: []api.Parameter{
					{
						Name:  "name1",
						Value: "hello",
					},
				},
			},
			want: []api.Parameter{
				{
					Name:  "name1",
					Value: "a",
				},
				{
					Name:  "name2",
					Value: "b",
				},
			},
		},
		{
			args: args{
				str: "name1=a,name2=b",
				para: []api.Parameter{
					{
						Name:  "name3",
						Value: "hello",
					},
				},
			},
			want: []api.Parameter{
				{
					Name:  "name3",
					Value: "hello",
				},
				{
					Name:  "name1",
					Value: "a",
				},
				{
					Name:  "name2",
					Value: "b",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertStrToPara(tt.args.str, tt.args.para); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertStringTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
