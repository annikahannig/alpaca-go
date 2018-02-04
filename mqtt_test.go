package alpaca

import (
	"testing"
)

func TestEncodeActionTopic(t *testing.T) {
	topics := Topics{
		"foo": "v1/basement/foo",
		"bar": "v1/bar",
	}

	expected := map[string]string{
		"foo/BAR":  "foo/BAR",
		"@foo/BAR": "v1/basement/foo/BAR",
		"@bar/FOO": "v1/bar/FOO",
		"@bar":     "@bar",
	}

	for actionType, result := range expected {
		encoded := encodeActionType(actionType, topics)
		if encoded != result {
			t.Error(
				"Encoding action type:", actionType,
				"yielded:", encoded, " expected was:", result,
			)
		}
	}
}

func TestDecodeActionTopic(t *testing.T) {

}
