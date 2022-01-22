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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"shanhu.io/misc/errcode"
	"shanhu.io/text/lexing"
)

// Decoder is a decoder that is capable of parsing a stream.
type Decoder struct {
	p *parser
}

// NewDecoder creates a new decoder that can parse a stream
// of jsonx objects.
func NewDecoder(r io.Reader) *Decoder {
	p, _ := newParser("", r)
	return &Decoder{p: p}
}

// More returns true if there is more stuff.
func (d *Decoder) More() bool {
	return !(d.p.See(lexing.EOF) || d.p.Token() == nil)
}

// Decode decodes a JSON value from the parser. When there is
// error on parsing JSONx, v is always unchanged.
func (d *Decoder) Decode(v interface{}) []*lexing.Error {
	value := parseValue(d.p)
	if errs := d.p.Errs(); errs != nil {
		return errs
	}
	if d.p.See(tokSemi) {
		d.p.Shift()
	}

	bs, errs := marshalValueLexing(value)
	if errs != nil {
		return errs
	}
	if err := json.Unmarshal(bs, v); err != nil {
		return lexing.SingleErr(err)
	}
	return nil
}

// DecodeSeries decode a typed series.
func (d *Decoder) DecodeSeries(tm TypeMaker) (
	[]*Typed, []*lexing.Error,
) {
	s := parseSeries(d.p)
	if errs := d.p.Errs(); errs != nil {
		return nil, errs
	}

	var res []*Typed

	errList := lexing.NewErrorList()
	for _, entry := range s.entries {
		typ := entry.typ.name
		v := tm(typ)
		if v == nil {
			errList.Add(&lexing.Error{
				Pos:  entry.typ.tok.Pos,
				Err:  fmt.Errorf("type %q unknown", typ),
				Code: "jsonx.unknownType",
			})
		}

		bs, errs := marshalValueLexing(entry.value)
		if errs != nil {
			errList.AddAll(errs)
		}
		if err := json.Unmarshal(bs, v); err != nil {
			errList.Add(&lexing.Error{
				Pos:  entry.typ.tok.Pos,
				Err:  fmt.Errorf("json marshal: %s", err),
				Code: "jsonx.marshalJSON",
			})
		}

		if errList.InJail() {
			errList.BailOut()
			continue
		}

		res = append(res, &Typed{
			Type: entry.typ.name,
			V:    v,
		})
	}

	if errs := errList.Errs(); errs != nil {
		return nil, errs
	}
	return res, nil
}

// Unmarshal unmarshals a value into a JSON object. When there is an error on
// parsing JSONx, v is always unchagned.
func Unmarshal(bs []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(bs))
	if errs := dec.Decode(v); errs != nil {
		return errs[0]
	}
	if dec.More() {
		return errcode.InvalidArgf("expect EOF, got more")
	}
	return nil
}

// ReadFile reads a file and unmarshals the content into the JSON object.
func ReadFile(file string, v interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return Unmarshal(bs, v)
}

// ReadFileMaybeJSON reads a file that might be JSONx or JSON.
func ReadFileMaybeJSON(file string, v interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := Unmarshal(bs, v); err != nil {
		// JSONx fails to decode, maybe it is plain JSON?
		if json.Unmarshal(bs, v) == nil {
			return nil
		}
		return err
	}
	return nil
}

// ReadSeriesFile reads a typed series.
func ReadSeriesFile(file string, tm TypeMaker) ([]*Typed, []*lexing.Error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	defer f.Close()

	dec := NewDecoder(f)
	return dec.DecodeSeries(tm)
}
