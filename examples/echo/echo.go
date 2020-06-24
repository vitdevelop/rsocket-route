package echo

import (
	"context"
	"fmt"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/vitdevelop/rsocket-route"
)

func init() {
	var err error
	err = route.Add(Routings)

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
