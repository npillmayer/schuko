package koanfadapter

import (
	"bytes"

	"github.com/npillmayer/nestext"
	"github.com/npillmayer/nestext/ntenc"
)

type NestedTextParser struct{}

func (nt NestedTextParser) Unmarshal(data []byte) (map[string]interface{}, error) {
	r := bytes.NewReader(data)
	result, err := nestext.Parse(r, nestext.TopLevel("dict"))
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

func (nt NestedTextParser) Marshal(m map[string]interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := ntenc.Encode(m, buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
