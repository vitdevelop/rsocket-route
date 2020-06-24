package decode

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/rsocket/rsocket-go/extension"
)

func Routes(metadataRaw []byte) ([]string, bool) {
	headers := MimeType(metadataRaw)
	if routingTags, exists := headers[extension.MessageRouting.String()]; exists {
		tags, err := extension.ParseRoutingTags(routingTags)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if err == nil && tags != nil && len(tags) > 0 {
			return tags, true
		}
	}
	return nil, false
}

func CompositeMetadata(metadataRaw []byte) (headers map[string][]byte) {
	headers = make(map[string][]byte)
	metadata := extension.NewCompositeMetadataBytes(metadataRaw).Scanner()
	for metadata.Scan() {
		mimeType, payloadMetadata, err := metadata.Metadata()
		if err != nil {
			_ = fmt.Errorf("%v\n", err)
			continue
		}
		headers[mimeType] = payloadMetadata
	}
	return
}

func MimeType(raw []byte) (headers map[string][]byte) {
	var mimeType string
	m := raw[0]
	idOrLen := (m << 1) >> 1
	if m&0x80 == 0x80 {
		mimeType = extension.MIME(idOrLen).String()
	} else {
		mimeTypeLen := int(idOrLen) + 1
		if cap(raw) < 1+mimeTypeLen {
			mimeType = string(raw)
		} else {
			mimeType = string(raw[1 : 1+mimeTypeLen])
		}
	}

	switch mimeType {
	case extension.MessageCompositeMetadata.String():
	case extension.MessageRouting.String():
		headers = CompositeMetadata(raw)
	default:
		headers = make(map[string][]byte)
		headers[mimeType] = []byte{}
	}
	return
}

func Metadata(metadataMimeType string, metadata []byte) (response map[string]string) {
	response = make(map[string]string)
	mime, ok := extension.ParseMIME(metadataMimeType)
	if !ok {
		return
	}

	switch mime {
	case extension.ApplicationJSON:
		_ = json.Unmarshal(metadata, &response)
	case extension.ApplicationCBOR:
		_ = cbor.Unmarshal(metadata, &response)
	}
	return
}
