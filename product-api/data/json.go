package data

import (
	"encoding/json"
	"io"
)

// encode data into json instead of marshalling
// NewEncoder() provides better performance than json.Unmasrshal as it does not
// have to buffer the output into as in memory slice of bytes

// this reduces allocations and the overheads of the service
// https://pkg.go.dev/encoding/json
// https://pkg.go.dev/io
// https://pkg.go.dev/time
// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
