package utils

import (
	"encoding/json"
	"io"
)

func DecodeBody[customType any](body io.ReadCloser) (*customType, error) {
	var data customType
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
