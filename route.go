package route

import (
	"github.com/rsocket/rsocket-go"
	"github.com/vitdevelop/rsocket-route/internal/handle"
)

func GetHandlers() []rsocket.OptAbstractSocket {
	return handle.Methods
}

func Add(routes map[string]interface{}) error {
	return handle.Paths(routes)
}
