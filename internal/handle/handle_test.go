package handle

import (
	"testing"

	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

func TestPaths(t *testing.T) {
	echoFunc := func(msg payload.Payload) mono.Mono {
		return mono.Just(msg)
	}
	echo := make(map[string]interface{})
	echo["/echo"] = echoFunc
	if err := Paths(echo); err != nil {
		t.Error(err)
	}
}