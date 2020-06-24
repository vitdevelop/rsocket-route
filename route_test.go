package route

import (
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github.com/vitdevelop/rsocket-route/internal/handle"
	"testing"
)

func TestGetHandlers(t *testing.T) {
	if len(handle.Methods) != 5 {
		t.Error("handle methods doesn't have 5 methods")
	}
}

func TestAdd(t *testing.T) {
	echoFunc := func(msg payload.Payload) mono.Mono {
		return mono.Just(msg)
	}
	echo := make(map[string]interface{})
	echo["/echo"] = echoFunc
	if err := Add(echo); err != nil {
		t.Error(err)
	}
}
