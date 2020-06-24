# rsocket-router

> Tool for [rsocket-go](https://github.com/rsocket/rsocket-go) implementation which help easy create routes for [RSocket](http://rsocket.io/) using composite metadata and make your code pretty.

## Features
 - Simple
 - Minimal configuration
 
## Usage

> main.go
```go
package echo

import (
	"context"
	"fmt"
	"github.com/vitdevelop/rsocket-route"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

func init() {
	err := route.Add(Routings)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func main() {
	err := rsocket.Receive().
		Resume().
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(route.GetHandlers()...), nil
		}).
		Transport("tcp://127.0.0.1:8080").
		Serve(context.Background())
	panic(err)
}
```
> controller.go
```go
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
```

#### Dependencies
 - [rsocket-go](https://github.com/rsocket/rsocket-go)
