package main

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_version(t *testing.T) {
	v := version{1, 2, 3}

	tests := []struct {
		name   string
		method func() string
		want   string
	}{
		{name: "string", method: v.String, want: "1.2.3"},
		{name: "next_major", method: v.nextMajor, want: "2.0.0"},
		{name: "next_minor", method: v.nextMinor, want: "1.3.0"},
		{name: "next_patch", method: v.nextPatch, want: "1.2.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(); got != tt.want {
				t.Errorf("got %s; want %s", got, tt.want)
			}
		})
	}
}

func Test_parseVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  version
		err   error
	}{
		{name: "valid version", input: "1.2.3", want: version{1, 2, 3}, err: nil},
		{name: "prefixed version", input: "v1.2.3", want: version{}, err: strconv.ErrSyntax},
		{name: "non numeric version", input: "x.y.z", want: version{}, err: strconv.ErrSyntax},
		{name: "wrong numbers count", input: "1.2.3.4", want: version{}, err: errWrongNumbersCount},
		{name: "negative number", input: "-1.2.3", want: version{}, err: errNegativeNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v; want %v", got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v; want %v", err, tt.err)
			}
		})
	}
}
