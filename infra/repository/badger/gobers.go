package badgerevent

import (
	"bytes"
	"encoding/gob"
)

func Encode[T any](toEncode T) ([]byte, error) {
	var buf bytes.Buffer

	if errEncode := gob.NewEncoder(&buf).Encode(toEncode); errEncode != nil {
		return []byte{}, errEncode
	}

	return buf.Bytes(), nil
}

func Decode[T any](toDecode []byte, decodeInTo T) error {
	return gob.NewDecoder(bytes.NewReader(toDecode)).Decode(decodeInTo)
}
