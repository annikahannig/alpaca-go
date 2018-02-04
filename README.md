

# Alpaca

A library for writing services on a MQTT message bus.

*This is still work in progress!*



## Why?

TODO: figure out.

Something along the lines of: After implementing [DaliQTT](https://github.com/cccb/daliqtt)
I wanted something more generalized to create new services more
easily.

I guess this will later need a python port aswell.

## How to use

Creating a new service in your network
is now as easy as:

    dispatch, actions, err := alpaca.DialMqtt(
        "tcp://user:pass@localhost:1889",
        alpaca.Topics{
            "lights": "v1/upstairs/lights",
            "meta": "v1/_meta/",
        })
    if err != nil {
        panic(err)
    }

    handle(dispatch, actions)
    

With a service handler like:

    const GET_LIGHT_VALUE_REQUEST = "@lights/GET_LIGHT_VALUE_REQUEST"

    type Light struct {
        Id     int `json:"id"`
        Value  int `json:"value"`
    }

    func handle(dispatch alpaca.Dispatch, actions chan alpaca.Action) {
        // Do some more setup stuff...

        // Handle incoming actions
        for action := range actions {
            switch action.Type {
            case SET_LIGHT_VALUE_REQUEST:
                var payload Light
                action.DecodePayload(&payload)
                setLightValue(light.Id, light.Value) 
                dispatch(SetLightValueSuccess(light))
            }
        }
    }


