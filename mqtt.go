package alpaca

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type Dispatch func(Action) error

type Topics map[string]string

/*
 Decode an incoming mqtt message and create an
 action from it's topic and payload
*/
func decodeMessage(msg mqtt.Message) (Action, error) {
	// Decode topic
	tokens := strings.Split(msg.Topic(), "/")
	actionType := tokens[len(tokens)-1]

	// Decode payload
	var payload interface{}
	var err error

	// Make action
	action := Action{
		Type:    actionType,
		Payload: payload,
	}

	return action, err
}

/*
 Encode an outgoing mqtt message payload
*/
func encodeMessagePayload(action Action) ([]byte, error) {
	payload, err := json.Marshal(action.Payload)

	return payload, err
}

/*
 Create dispatch function:
 Encode action for transport and publish to MQTT
*/
func makeDispatch(client mqtt.Client, baseTopic string) Dispatch {
	dispatch := func(action Action) error {
		// Prepare payload
		topic := baseTopic + "/" + action.Type
		payload, err := encodeMessagePayload(action)
		if err != nil {
			return err
		}

		// Send message
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		return nil
	}

	return dispatch
}

/*
 Connect to MQTT broker and create action channel
 and dispatch function.
*/
func Dial(brokerUri string) (chan Action, Dispatch, error) {
	actions := make(chan Action)

	opts := mqtt.NewClientOptions()

	// Basic configuration
	opts.AddBroker(brokerUri)
	opts.SetClientID("daliqtt")

	// Reconnects and timeouts
	opts.SetMaxReconnectInterval(15.0 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetKeepAlive(2 * time.Second)

	// Register handler funcs
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		action, err := decodeMessage(msg)
		if err != nil {
			log.Println("Error while decoding message:", err)
			return
		}

		// Forward to service
		actions <- action
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		// Subscribe to topic
		/*
			topic := config.BaseTopic + "/#"
			if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}

			log.Println("Subscribed to topic:", topic)

			// Subscribe to meta topic
			topic = "_meta/#"
			if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}

			log.Println("Subscribed to topic:", topic)
		*/
		log.Println("Implement this.")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, nil, token.Error()
	}

	// Create dispatch function
	dispatch := makeDispatch(client, "null")

	return actions, dispatch, nil
}
