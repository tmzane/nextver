package main

import (
	"errors"
	"reflect"
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
		{name: "next major", method: v.nextMajor, want: "2.0.0"},
		{name: "next minor", method: v.nextMinor, want: "1.3.0"},
		{name: "next patch", method: v.nextPatch, want: "1.2.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(); got != tt.want {
				t.Errorf("got %s; want %s", got, tt.want)
			}
		})
	}
}

func Test_parseVersion_valid(t *testing.T) {
	test := func(input string, want version) {
		t.Run(input, func(t *testing.T) {
			got, err := parseVersion(input)
			if !reflect.DeepEqual(got, &want) {
				t.Errorf("got %s; want %s", got, &want)
			}
			if err != nil {
				t.Errorf("got %v; want nil", err)
			}
		})
	}

	// data source: https://regex101.com/r/vkijKf/1/
	test("0.0.4", version{0, 0, 4})
	test("1.2.3", version{1, 2, 3})
	test("10.20.30", version{10, 20, 30})
	test("1.1.2-prerelease+meta", version{1, 1, 2})
	test("1.1.2+meta", version{1, 1, 2})
	test("1.1.2+meta-valid", version{1, 1, 2})
	test("1.0.0-alpha", version{1, 0, 0})
	test("1.0.0-beta", version{1, 0, 0})
	test("1.0.0-alpha.beta", version{1, 0, 0})
	test("1.0.0-alpha.beta.1", version{1, 0, 0})
	test("1.0.0-alpha.1", version{1, 0, 0})
	test("1.0.0-alpha0.valid", version{1, 0, 0})
	test("1.0.0-alpha.0valid", version{1, 0, 0})
	test("1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", version{1, 0, 0})
	test("1.0.0-rc.1+build.1", version{1, 0, 0})
	test("2.0.0-rc.1+build.123", version{2, 0, 0})
	test("1.2.3-beta", version{1, 2, 3})
	test("10.2.3-DEV-SNAPSHOT", version{10, 2, 3})
	test("1.2.3-SNAPSHOT-123", version{1, 2, 3})
	test("1.0.0", version{1, 0, 0})
	test("2.0.0", version{2, 0, 0})
	test("1.1.7", version{1, 1, 7})
	test("2.0.0+build.1848", version{2, 0, 0})
	test("2.0.1-alpha.1227", version{2, 0, 1})
	test("1.0.0-alpha+beta", version{1, 0, 0})
	test("1.2.3----RC-SNAPSHOT.12.9.1--.12+788", version{1, 2, 3})
	test("1.2.3----R-S.12.9.1--.12+meta", version{1, 2, 3})
	test("1.2.3----RC-SNAPSHOT.12.9.1--.12", version{1, 2, 3})
	test("1.0.0+0.build.1-rc.10000aaa-kk-0.1", version{1, 0, 0})
	test("1.0.0-0A.is.legal", version{1, 0, 0})
}

func Test_parseVersion_invalid(t *testing.T) {
	test := func(input string) {
		t.Run(input, func(t *testing.T) {
			got, err := parseVersion(input)
			if got != nil {
				t.Errorf("got %s; want nil", got)
			}
			if !errors.Is(err, errInvalidFormat) {
				t.Errorf("got %v; want %v", err, errInvalidFormat)
			}
		})
	}

	// data source: https://regex101.com/r/vkijKf/1/
	test("1")
	test("1.2")
	test("1.2.3-0123")
	test("1.2.3-0123.0123")
	test("1.1.2+.123")
	test("+invalid")
	test("-invalid")
	test("-invalid+invalid")
	test("-invalid.01")
	test("alpha")
	test("alpha.beta")
	test("alpha.beta.1")
	test("alpha.1")
	test("alpha+beta")
	test("alpha_beta")
	test("alpha.")
	test("alpha..")
	test("beta")
	test("1.0.0-alpha_beta")
	test("-alpha.")
	test("1.0.0-alpha..")
	test("1.0.0-alpha..1")
	test("1.0.0-alpha...1")
	test("1.0.0-alpha....1")
	test("1.0.0-alpha.....1")
	test("1.0.0-alpha......1")
	test("1.0.0-alpha.......1")
	test("01.1.1")
	test("1.01.1")
	test("1.1.01")
	test("1.2")
	test("1.2.3.DEV")
	test("1.2-SNAPSHOT")
	test("1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788")
	test("1.2-RC-SNAPSHOT")
	test("-1.0.3-gamma+b7718")
	test("+justmeta")
	test("9.8.7+meta+meta")
	test("9.8.7-whatever+meta+meta")
}
