package dtls

import (
	"bytes"
	"context"
	"encoding/hex"

	"github.com/pion/dtls/v2/pkg/protocol"
	"github.com/pion/dtls/v2/pkg/protocol/alert"
	"github.com/pion/dtls/v2/pkg/protocol/handshake"
	"github.com/pion/dtls/v2/pkg/protocol/recordlayer"
)

func flight2Parse(ctx context.Context, c flightConn, state *State, cache *handshakeCache, cfg *handshakeConfig) (flightVal, *alert.Alert, error) {
	seq, msgs, ok := cache.fullPullMap(state.handshakeRecvSequence,
		handshakeCachePullRule{handshake.TypeClientHello, cfg.initialEpoch, true, false},
	)
	if !ok {
		// Client may retransmit the first ClientHello when HelloVerifyRequest is dropped.
		// Parse as flight 0 in this case.
		return flight0Parse(ctx, c, state, cache, cfg)
	}
	state.handshakeRecvSequence = seq

	var clientHello *handshake.MessageClientHello

	// Validate type
	if clientHello, ok = msgs[handshake.TypeClientHello].(*handshake.MessageClientHello); !ok {
		return 0, &alert.Alert{Level: alert.Fatal, Description: alert.InternalError}, nil
	}

	if !clientHello.Version.Equal(protocol.Version1_2) {
		return 0, &alert.Alert{Level: alert.Fatal, Description: alert.ProtocolVersion}, errUnsupportedProtocolVersion
	}

	if len(clientHello.Cookie) == 0 {
		return 0, nil, nil
	}
	if !bytes.Equal(state.cookie, clientHello.Cookie) {
		return 0, &alert.Alert{Level: alert.Fatal, Description: alert.AccessDenied}, errCookieMismatch
	}

	if len(clientHello.SessionID) > 0 && cfg.sessionStore != nil {
		id := hex.EncodeToString(clientHello.SessionID)
		if s := cfg.sessionStore.Get(id); s != nil {
			cfg.log.Tracef("[handshake] resume session for: %s", id)

			state.masterSecret = s.Secret
			state.SessionID = clientHello.SessionID

			return flight4b, nil, nil
		}
	}

	return flight4, nil, nil
}

func flight2Generate(c flightConn, state *State, cache *handshakeCache, cfg *handshakeConfig) ([]*packet, *alert.Alert, error) {
	state.handshakeSendSequence = 0
	return []*packet{
		{
			record: &recordlayer.RecordLayer{
				Header: recordlayer.Header{
					Version: protocol.Version1_2,
				},
				Content: &handshake.Handshake{
					Message: &handshake.MessageHelloVerifyRequest{
						Version: protocol.Version1_2,
						Cookie:  state.cookie,
					},
				},
			},
		},
	}, nil, nil
}
