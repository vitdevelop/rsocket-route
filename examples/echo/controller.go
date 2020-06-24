package echo

import (
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

var Routings = make(map[string]interface{})

func init() {
	Routings["/request-response"] = echoRR
	Routings["/request-stream"] = echoRS
	Routings["/request-channel"] = echoRC
}

func echoRR(msg payload.Payload) mono.Mono {
	return mono.Just(msg)
}

func echoRS(msg payload.Payload) flux.Flux {
	return flux.Just(msg)
}

func echoRC(msgs rx.Publisher) flux.Flux {
	return msgs.(flux.Flux)
}