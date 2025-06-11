package util

import "testing"

func TestIsValidTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "empty", args: args{tag: ""}, want: false},
		{name: "abc", args: args{tag: "abc"}, want: false},
		{name: "v112.3", args: args{tag: "v12.3"}, want: false},
		{name: "v1.2.a", args: args{tag: "v1.2.a"}, want: false},
		{name: "v1.2.3", args: args{tag: "v1.2.33"}, want: true},
		{name: "v11.22.33", args: args{tag: "v11.22.33"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTag(tt.args.tag); got != tt.want {
				t.Errorf("IsValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
