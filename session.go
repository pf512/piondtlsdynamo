package dtls

import (
	"encoding/hex"
	"fmt"
)

type Session struct {
	ID     []byte
	Secret []byte
	Addr   string
}

func (s *Session) String() string {
	return fmt.Sprintf(
		"[id: %s, secret: %s, addr: %s]",
		hex.EncodeToString(s.ID),
		hex.EncodeToString(s.Secret),
		s.Addr,
	)
}

type SessionStore interface {
	// Set save a session to store. And will be fetched by id or addr.
	Set(s *Session) error
	// Get fetch a session by session id or remove addrese
	Get(idOrAddr string) *Session
	// Del rm saved session
	Del(idOrAddr string) error
}
