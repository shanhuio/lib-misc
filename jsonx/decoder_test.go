// Copyright (C) 2022  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package jsonx

import (
	"testing"

	"reflect"
	"strings"
)

func TestDecoder(t *testing.T) {
	input := strings.NewReader(`"a""b";"c"`)

	dec := NewDecoder(input)
	var got []string
	for dec.More() {
		var s string
		if err := dec.Decode(&s); err != nil {
			t.Fatal(err)
		}
		got = append(got, s)
	}

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnmarshal(t *testing.T) {
	var v int
	if err := Unmarshal([]byte("1234"), &v); err != nil {
		t.Fatal(err)
	}
	if v != 1234 {
		t.Errorf("got %d, want 1234", v)
	}
}

func TestDecoder_series(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`str "string"`,
		`num 3`,
		`struct {Field: "value"}`,
	}, "\n"))

	type structType struct {
		Field string
	}

	dec := NewDecoder(input)
	tm := func(t string) interface{} {
		switch t {
		case "str":
			return new(string)
		case "num":
			return new(int)
		case "struct":
			return new(structType)
		}
		return nil
	}
	list, errs := dec.DecodeSeries(tm)
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}

	strVal := "string"
	numVal := 3
	want := []*Typed{
		{Type: "str", V: &strVal},
		{Type: "num", V: &numVal},
		{Type: "struct", V: &structType{Field: "value"}},
	}
	if !reflect.DeepEqual(list, want) {
		t.Errorf("got %v, want %v", list, want)
	}
}
