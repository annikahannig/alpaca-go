package alpaca

import (
	"encoding/json"
	"testing"
)

type TestPayload struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestPayloadDecoding(t *testing.T) {

	a := Action{
		Type:    "FOO",
		Payload: []byte("{\"foo\": \"fnord\", \"bar\": 42}"),
	}

	var payload TestPayload
	err := a.DecodePayload(&payload)
	if err != nil {
		t.Error(err)
	}

	if payload.Foo != "fnord" {
		t.Error("Expected payload.Foo to be fnord")
	}

	if payload.Bar != 42 {
		t.Error("Expected payload.Bar to be 42")
	}
}

func TestPayloadDecodingError(t *testing.T) {
	a := Action{
		Type:    "FOO",
		Payload: []byte("{\"fo: fnord 42}"),
	}

	var payload TestPayload
	err := a.DecodePayload(&payload)
	if err == nil {
		t.Error("Expected decoding error in payload")
	}
}

func TestPayloadTypeError(t *testing.T) {
	a := Action{
		Type:    "FOO",
		Payload: TestPayload{"fnord", 42},
	}

	// This should trigger a decode error, as the payload
	// is already decoded
	var payload TestPayload
	err := a.DecodePayload(&payload)
	if err == nil {
		t.Error("Expected type error in payload decoding")
	}
}
