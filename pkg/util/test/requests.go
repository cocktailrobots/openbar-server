package test

import (
	"bytes"
	"encoding/json"
	"io"
)

func JsonReaderForObject(obj any) io.Reader {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(data)
}
