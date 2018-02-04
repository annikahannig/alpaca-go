

# Alpaca

A library for writing services on a MQTT message bus.

## Why?

TODO: figure out.

## How?

Creating a new service in your network
is now as easy as:

    conn, err := alpaca.Dial("mqtt://user:pass@localhost:1889")
    if err != nil {
        panic(err)
    }

    dispatch, actions := conn.Join(alpaca.Topics{
        "lights": "v1/upstairs/lights",
        "meta": "v1/_meta/",
    })
    
    handle(dispatch, actions)
    


With a service handler
    const GET_LIGHT_VALUE_REQUEST = "@lights/GET_LIGHT_VALUE_REQUEST"

    func handle(dispatch alpaca.Dispatch, actions chan alpaca.Action) {
        // Do some more setup stuff...

        // Handle incoming actions
        for action := range actions {
            switch action.Type {
                case GET_LIGHT_VALUE_REQUEST:
                    dispatch(GetLightValueSuccess(lightValue(action.Payload
            }
        }
    }


