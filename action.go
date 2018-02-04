package alpaca

import (
	"encoding/json"
	"fmt"
)

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (a Action) DecodePayload(v interface{}) error {
	data, ok := a.Payload.([]byte)
	if !ok {
		return fmt.Errorf("Expected payload to be []byte")
	}

	return json.Unmarshal(data, v)
}
