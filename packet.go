package dtls

import "github.com/pf512/piondtlsdynamo/pkg/protocol/recordlayer"

type packet struct {
	record                   *recordlayer.RecordLayer
	shouldEncrypt            bool
	resetLocalSequenceNumber bool
}
