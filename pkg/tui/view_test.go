package tui

import (
	"testing"

	"github.com/Nekroze/bolton/pkg/boltons"
	"github.com/Nekroze/bolton/pkg/hardpoints"
)

func Test_diffPreview(t *testing.T) {
	type args struct {
		input string
		b     boltons.Bolton
		h     *hardpoints.Point
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy_path",
			args: args{
				input: `Foo
# HARDPOINT: shell_function
foobar() {
	true
}`,
				b: boltons.Bolton{Path: "test", Name: "test"},
				h: &hardpoints.Point{Path: "", Line: 2, Tags: []string{"shell_function"}},
			},
			want: `Foo
# HARDPOINT: shell_function
[green]+true[-]
[green]+[-]
foobar() {
	true
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diffPreview(tt.args.input, tt.args.b, tt.args.h); got != tt.want {
				t.Errorf("diffPreview() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apply(t *testing.T) {
	type args struct {
		input string
		b     boltons.Bolton
		h     *hardpoints.Point
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy_path",
			args: args{
				input: `Foo
# HARDPOINT: shell_function
foobar() {
	true
}`,
				b: boltons.Bolton{Path: "test", Name: "test"},
				h: &hardpoints.Point{Path: "", Line: 2, Tags: []string{"shell_function"}},
			},
			want: `Foo
# HARDPOINT: shell_function
true

foobar() {
	true
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := apply(tt.args.input, tt.args.b, tt.args.h); got != tt.want {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
