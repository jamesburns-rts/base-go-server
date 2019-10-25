package util

import "testing"

func TestToSnakeCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "typical1",
			args: args{"NumberOne"},
			want: "number_one",
		},
		{
			name: "no change",
			args: args{"number_one"},
			want: "number_one",
		},
		{
			name: "lower first",
			args: args{"numberOne"},
			want: "number_one",
		},
	}
	for _, tt := range tests {
		if got := CamelToSnakeCase(tt.args.str); got != tt.want {
			t.Errorf("%q. CamelToSnakeCase() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
