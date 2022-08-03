package util

import (
	"testing"
)

func TestCamel2Case(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{name: "DbName"}, want: "db_name"},
		{name: "case2", args: args{name: "TaskName"}, want: "task_name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Camel2Case(tt.args.name); got != tt.want {
				t.Errorf("Camel2Case() = %v, want %v", got, tt.want)
			}
		})
	}
}
