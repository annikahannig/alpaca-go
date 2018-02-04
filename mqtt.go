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

type Client struct {
	client mqtt.Client
	topics Topics
}

/*
 Decode an incoming mqtt message and create an
 action from it's topic and payload
*/
func decodeMessage(msg mqtt.Message, topics Topics) Action {
	// Make action
	action := Action{
		Type:    decodeTopic(msg.Topic(), topics),
		Payload: msg.Payload(),
	}

	return action
}

/*
 Encode an outgoing mqtt message payload
*/
func encodeMessagePayload(action Action) ([]byte, error) {
	payload, err := json.Marshal(action.Payload)

	return payload, err
}

/*
 Encode topic from action type:
 In case actions are prefixed with an @ we make a lookup
 on the topics registry and expand the topic:

 Example:

    @lights/SET_VALUE

 will expand to

    v1/upstairs/lights/SET_VALUE
*/
func encodeActionType(actionType string, topics Topics) string {
	tokens := strings.SplitN(actionType, "/", 2)
	if len(tokens) == 1 {
		return actionType // Nothing to do here
	}

	if !strings.HasPrefix(tokens[0], "@") {
		return actionType // Still nothing to do here
	}

	route, ok := topics[tokens[0][1:]]
	if !ok {
		log.Println("Warning: Could not expand route for", actionType)
		return actionType
	}

	return route + "/" + strings.Join(tokens[1:], "/")
}

/*
 Decode topic from message (for use in action Type)
*/
func decodeTopic(topic string, topics Topics) string {
	tokens := strings.Split(topic, "/")
	if len(tokens) == 1 {
		return topic // Nothing to do here
	}

	for handle, route := range topics {
		if strings.HasPrefix(topic, route) {
			return "@" + strings.Replace(topic, route, handle, 1)
		}
	}

	return topic
}

/*
 Create dispatch function:
 Encode action for transport and publish to MQTT
*/
func makeDispatch(client mqtt.Client, topics Topics) Dispatch {
	dispatch := func(action Action) error {
		// Prepare payload
		topic := encodeActionType(action.Type, topics)
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

func makeOnConnectHandler(topics Topics) mqtt.OnConnectHandler {
	handler := func(client mqtt.Client) {
		// Subscribe to topics
		for _, base := range topics {
			// We are interested in all messages on this topic
			topic := base + "/#"

			token := client.Subscribe(topic, 0, nil)
			if token.Wait() && token.Error() != nil {
				panic(token.Error())
			}

			log.Println("Subscribed to topic:", topic)
		}
	}

	return handler
}

/*
 Create message handler for receiving messages, decoding the actions and
 dispatching them into the actions channel.
*/
func makeMessageHandler(actions Actions, topics Topics) mqtt.MessageHandler {
	handler := func(client mqtt.Client, msg mqtt.Message) {
		action := decodeMessage(msg, topics)
		// Forward to handler
		actions <- action
	}

	return handler
}

/*
 Connect to MQTT broker and create action channel
 and dispatch function.
*/
func DialMqtt(brokerUri string, topics Topics) (Actions, Dispatch) {
	opts := mqtt.NewClientOptions()

	// Basic configuration
	opts.AddBroker(brokerUri)

	// Reconnects and timeouts
	opts.SetMaxReconnectInterval(15.0 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetKeepAlive(2 * time.Second)

	return Connect(opts, topics)
}

/*
 Connect to MQTT broker like DialMqtt, but give the user
 more control over the client options.
*/
func Connect(opts *mqtt.ClientOptions, topics Topics) (Actions, Dispatch) {

	// Create actions channel
	actions := make(Actions)

	// Register handler funcs
	opts.SetOnConnectHandler(makeOnConnectHandler(topics))
	opts.SetDefaultPublishHandler(makeMessageHandler(actions, topics))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Create dispatch function
	dispatch := makeDispatch(client, topics)

	return actions, dispatch
}
