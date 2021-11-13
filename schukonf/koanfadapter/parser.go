package koanfadapter

import (
	"bytes"

	"github.com/npillmayer/nestext"
	"github.com/npillmayer/nestext/ntenc"
)

// NestedTextParser is a thin wrapper on top of npillmayer/nestext to enable using
// NestedText (see https://nestedtext.org) as a configuration format.
// Koanf needs parsers to implement the koanf.Parser interface.
type NestedTextParser struct{}

// Unmarshal is part of the koanf.Parser interface.
func (nt NestedTextParser) Unmarshal(data []byte) (map[string]interface{}, error) {
	r := bytes.NewReader(data)
	result, err := nestext.Parse(r, nestext.TopLevel("dict"))
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// Marshal is part of the koanf.Parser interface.
func (nt NestedTextParser) Marshal(m map[string]interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := ntenc.Encode(m, buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
