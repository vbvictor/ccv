package complexity

import (
	"encoding/json"
	"io"
)

type churnJSON struct {
	Files []*ChurnChunk `json:"files"`
}

func ReadChurn(r io.Reader) ([]*ChurnChunk, error) {
	var data churnJSON
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	return data.Files, nil
}
