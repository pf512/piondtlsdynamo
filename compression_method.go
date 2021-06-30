package dtls

import "github.com/pf512/piondtlsdynamo/pkg/protocol"

func defaultCompressionMethods() []*protocol.CompressionMethod {
	return []*protocol.CompressionMethod{
		{},
	}
}
