package handle

import (
	"errors"
	"fmt"
	"github.com/vitdevelop/rsocket-route/decode"

	"github.com/rsocket/rsocket-go"

	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

var Methods []rsocket.OptAbstractSocket

func init() {
	Methods = make([]rsocket.OptAbstractSocket, 0)
	Methods = append(Methods, requestResponse())
	Methods = append(Methods, fireAndForget())
	Methods = append(Methods, requestStream())
	Methods = append(Methods, requestChannel())
	Methods = append(Methods, metadataPush())
}

var rr = make(map[string]func(msg payload.Payload) mono.Mono)
var fnf = make(map[string]func(msg payload.Payload))
var rs = make(map[string]func(msg payload.Payload) flux.Flux)
var rc = make(map[string]func(msgs rx.Publisher) flux.Flux)

func Paths(paths map[string]interface{}) (err error) {
	for path, genericMethod := range paths {
		switch method := genericMethod.(type) {
		case func(msg payload.Payload) mono.Mono:
			if _, ok := rr[path]; ok {
				err = pathError(path)
				continue
			}
			rr[path] = method
		case func(msg payload.Payload):
			if _, ok := fnf[path]; ok {
				err = pathError(path)
				continue
			}
			fnf[path] = method
		case func(msg payload.Payload) flux.Flux:
			if _, ok := rs[path]; ok {
				err = pathError(path)
				continue
			}
			rs[path] = method
		case func(msgs rx.Publisher) flux.Flux:
			if _, ok := rc[path]; ok {
				err = pathError(path)
				continue
			}
			rc[path] = method
		default:
			err = errors.New(fmt.Sprintf("%v for %s is unknown", genericMethod, path))
		}
	}
	return
}

func pathError(path string) error {
	return errors.New(fmt.Sprintf("path: %s duplicated", path))
}

func requestResponse() rsocket.OptAbstractSocket {
	return rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
		if metadata, exists := msg.Metadata(); exists {
			if routings, exists := decode.Routes(metadata); exists {
				for _, rrRoute := range routings {
					if funcRoute, exists := rr[rrRoute]; exists {
						return funcRoute(msg)
					}
				}
			}
		}
		return mono.Error(errors.New("routing not found"))
	})
}

func fireAndForget() rsocket.OptAbstractSocket {
	return rsocket.FireAndForget(func(msg payload.Payload) {
		if metadata, exists := msg.Metadata(); exists {
			if routings, exists := decode.Routes(metadata); exists {
				for _, rrRoute := range routings {
					if funcRoute, exists := fnf[rrRoute]; exists {
						funcRoute(msg)
					}
				}
			}
		}
	})
}

func requestStream() rsocket.OptAbstractSocket {
	return rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
		if metadata, exists := msg.Metadata(); exists {
			if routings, exists := decode.Routes(metadata); exists {
				for _, rrRoute := range routings {
					if funcRoute, exists := rs[rrRoute]; exists {
						return funcRoute(msg)
					}
				}
			}
		}
		return flux.Error(errors.New("routing not found"))
	})
}

func requestChannel() rsocket.OptAbstractSocket {
	return rsocket.RequestChannel(func(msgs rx.Publisher) flux.Flux {
		return msgs.(flux.Flux).
			SwitchOnFirst(func(sig flux.Signal, f flux.Flux) flux.Flux {
				if msg, ok := sig.Value(); ok {
					if metadata, exists := msg.Metadata(); exists {
						if routings, exists := decode.Routes(metadata); exists {
							for _, rrRoute := range routings {
								if funcRoute, exists := rc[rrRoute]; exists {
									return funcRoute(f.(rx.Publisher))
								}
							}
						}
					}
				}
				return flux.Error(errors.New("routing not found"))
			})
	})
}

func metadataPush() rsocket.OptAbstractSocket {
	return rsocket.MetadataPush(func(msg payload.Payload) {
		if metadata, exists := msg.Metadata(); exists {
			decode.MimeType(metadata)
		}
	})
}
