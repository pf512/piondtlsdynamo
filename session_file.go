package dtls

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type FileSessionStore struct {
	Root string
	TTL  time.Duration
}

type hexSession struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Addr     string `json:"addr"`
	ExpireAt time.Time
}

func (fs *FileSessionStore) Set(s *Session, isClient bool) {
	d := hexSession{
		ID:     hex.EncodeToString(s.ID),
		Secret: hex.EncodeToString(s.Secret),
		Addr:   s.Addr,
	}

	if fs.TTL > 0 {
		d.ExpireAt = time.Now().Add(fs.TTL)
	}

	idPath := path.Join(fs.Root, hex.EncodeToString(s.ID))
	f, err := os.OpenFile(idPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Println("open file error", err)
		return
	}

	if err = json.NewEncoder(f).Encode(d); err != nil {
		log.Println("encode error", err)
		return
	}

	if !isClient {
		return
	}

	addrPath := path.Join(fs.Root, s.Addr)
	if err = os.Link(idPath, addrPath); err != nil {
		log.Println("link error", err)
	}
}

func (fs *FileSessionStore) Get(id []byte) (s *Session) {
	return fs.get(path.Join(fs.Root, hex.EncodeToString(id)), true)
}

func (fs *FileSessionStore) GetByAddr(addr string) *Session {
	return fs.get(path.Join(fs.Root, addr), true)
}

func (fs *FileSessionStore) get(path string, checkTTL bool) (s *Session) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error", err)
		return
	}

	d := hexSession{}
	err = json.NewDecoder(f).Decode(&d)
	if err != nil {
		fmt.Println("decode error", err)
		return
	}

	s = &Session{Addr: d.Addr}

	s.ID, err = hex.DecodeString(d.ID)
	if err != nil {
		fmt.Println("decode id error", err)
		return
	}

	if checkTTL && !d.ExpireAt.IsZero() && d.ExpireAt.Before(time.Now()) {
		fs.Del(s.ID)
		return nil
	}

	s.Secret, err = hex.DecodeString(d.Secret)
	if err != nil {
		fmt.Println("decode secret error", err)
		return
	}

	return
}

func (fs *FileSessionStore) Del(id []byte) {
	sid := hex.EncodeToString(id)
	s := fs.get(path.Join(fs.Root, sid), false)
	if s == nil {
		return
	}

	os.Remove(path.Join(fs.Root, sid))
	os.Remove(path.Join(fs.Root, s.Addr))
}