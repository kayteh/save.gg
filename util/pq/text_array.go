package pq

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strings"
)

// A text array, identified by Postgres as text[].
type TextArray struct {
	data []string
}

// Scan is used by the sql driver to "scan" in relevant values.
//
// TextArray.Scan in particular converts the []byte to a []string.
//
// Also see TextArray.UnmarshalJSON for another conversion.
func (s *TextArray) Scan(value interface{}) (err error) {
	var b []byte

	switch value.(type) {
	case []byte:
		b = value.([]byte)
	default:
		return ErrTypeMismatch
	}

	if len(b) <= 2 {
		return nil
	}

	p := bytes.Split(b[1:len(b)-1], []byte(","))

	for _, r := range p {
		if len(r) == 0 {
			continue
		} else {
			s.data = append(s.data, string(r))
		}
	}

	return nil

}

// Value is used by the sql driver to "get" a representation of the type.
//
// TextArray.Value in particular converts the []string by joining it, separated by `,`,
// and wrapping it in `{}`, which is exactly as Postgres expects.
func (s *TextArray) Value() (o driver.Value, err error) {
	return fmt.Sprintf("{%s}", strings.Join(s.Copy(), ",")), nil
}

// Appends a new string to the underlying []string.
func (s *TextArray) Append(v string) {
	s.data = append(s.data, v)
}

// Set the underlying []string. This is destructive.
func (s *TextArray) Set(v []string) {
	s.data = v
}

// Check if the []string includes a value.
func (s *TextArray) Has(v string) bool {
	for _, i := range s.data {
		if i == v {
			return true
		}
	}
	return false
}

// Removes all values that match v from the underlying []string.
func (s *TextArray) Remove(v string) {
	var n []string

	for _, i := range s.data {
		if i != v {
			n = append(n, i)
		}
	}

	s.Set(n)
}

// Returns a copy of the underlying []string. Edits to what this returns are not preserved.
func (s *TextArray) Copy() []string {
	var o []string

	if len(s.data) > 0 {
		copy(s.data, o)
	}

	return o
}

// Modification of Scan to de-JSON-ify JSON.
//TODO(kkz): Properly implement this
func (s *TextArray) UnmarshalJSON(b []byte) error {

	if len(b) <= 2 {
		return nil
	}

	p := bytes.Split(b[2:len(b)-2], []byte(`","`))

	for _, r := range p {
		if len(r) == 0 {
			continue
		} else {
			s.data = append(s.data, string(r))
		}
	}

	return nil
}

// Modification of Value to output a JSON []byte.
func (s *TextArray) MarshalJSON() ([]byte, error) {
	if len(s.data) == 0 {
		return []byte("[]"), nil
	}

	return []byte(fmt.Sprintf(`["%s"]`, strings.Join(s.data, `","`))), nil
}
