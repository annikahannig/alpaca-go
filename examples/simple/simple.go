package main

import (
	"log"

	"github.com/mhannig/alpaca-go"
)

// Actions
const REVERSE_REQUEST = "@strings/REVERSE_REQUEST"
const REVERSE_SUCCESS = "@strings/REVERSE_SUCCESS"
const REVERSE_ERROR = "@strings/REVERSE_ERROR"

// Action Creators
func ReverseSuccess(str string) alpaca.Action {
	return alpaca.Action{
		Type:    REVERSE_SUCCESS,
		Payload: str,
	}
}

func ReverseError(err error) alpaca.Action {
	return alpaca.Action{
		Type:    REVERSE_ERROR,
		Payload: err.Error(),
	}
}

func handleReverse(action alpaca.Action, dispatch alpaca.Dispatch) {
	var payload string
	err := action.DecodePayload(&payload)
	if err != nil {
		log.Println("Could not decode payload:", err)
		dispatch(ReverseError(err))
		return
	}

	reverse := ""
	for _, c := range payload {
		reverse = string(c) + reverse
	}

	log.Println("Reversed string:", payload, " -> ", reverse)
	dispatch(ReverseSuccess(reverse))
}

func main() {

	actions, dispatch := alpaca.DialMqtt(
		"tcp://localhost:1883",
		alpaca.Routes{
			"strings": "v1/simple/strings",
		})

	for action := range actions {
		switch action.Type {
		case REVERSE_REQUEST:
			handleReverse(action, dispatch)
		}
	}

}
