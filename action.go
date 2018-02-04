package alpaca

import (
	"encoding/json"
	"fmt"
)

type Payload interface{}

type Action struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Actions chan Action

/*
 Decode the payload.
 If an action is received from MQTT, the payload
 should be raw bytes. We decode this then on
 demand, as the handler function should know what
 kind of payload it requires.
*/
func (a Action) DecodePayload(v interface{}) error {
	data, ok := a.Payload.([]byte)
	if !ok {
		return fmt.Errorf("Expected payload to be []byte")
	}

	return json.Unmarshal(data, v)
}
