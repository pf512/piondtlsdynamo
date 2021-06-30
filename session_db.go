package piondtlsdynamo

import (
	"encoding/hex"
	sessiondb "github.com/pf512/piondtlsdynamo/db"
	"log"
	"os"
	"path"
	"time"
)

// DbSessionStore is a simple db based SessionStore.
// You need set a root path to store the session data.
// And you can set an optional TTL to avoid long time session.
//
// DbSessionStore only clean session while fetching.  If you
// want clean more aggressively, you could call the Clean() func.
type DbSessionStore struct {
	// Root store the session dir root path.
	Root string
	// TTL store the session store time duration.
	TTL time.Duration
}

/*
type hexSession struct {
	ID       string    `json:"id"`
	Secret   string    `json:"secret"`
	Addr     string    `json:"addr"`
	ExpireAt time.Time `json:"expire_at"`
}
*/


func (db *DbSessionStore) Set(s *Session, isClient bool) {
	d := hexSession{
		ID:     hex.EncodeToString(s.ID),
		Secret: hex.EncodeToString(s.Secret),
		Addr:   s.Addr,
	}

	if db.TTL > 0 {
		d.ExpireAt = time.Now().Add(db.TTL)
	}

	/*
	idPath := path.Join(db.Root, hex.EncodeToString(s.ID))
	f, err := os.OpenFile(idPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Println("open file error", err)
		return
	}
	*/

	sessiondb.StoreSession(d.ID, d.Secret, d.Addr, d.ExpireAt)

	/*
	if err = json.NewEncoder(f).Encode(d); err != nil {
		log.Println("encode error", err)
		return
	}

	 */

	if !isClient {
		return
	}

	/*
	addrPath := path.Join(db.Root, s.Addr)
	if err = os.Link(idPath, addrPath); err != nil {
		log.Println("link error", err)
	}
	 */
}

func (db *DbSessionStore) Get(id []byte) (s *Session) {
	return db.get(path.Join(hex.EncodeToString(id)), true)
}

func (db *DbSessionStore) GetByAddr(addr string) *Session {
	// todo not used by server
	return db.get(path.Join(db.Root, addr), true)
}

func (db *DbSessionStore) get(path string, checkTTL bool) (s *Session) {

	/*
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Println("open file error", err)
		return
	}
	 */

	dbsession, err := sessiondb.RetrieveSession(path)

	/*
	err = json.NewDecoder(f).Decode(&d)
	if err != nil {
		log.Println("decode error", err)
		return
	}
	 */

	s = &Session{Addr: dbsession.Address}

	s.ID, err = hex.DecodeString(dbsession.ID)
	if err != nil {
		log.Println("decode id error", err)
		return
	}

	if checkTTL && !dbsession.Expiration.IsZero() && dbsession.Expiration.Before(time.Now()) {
		db.Del(s.ID)
		return nil
	}

	s.Secret, err = hex.DecodeString(dbsession.Secret)
	if err != nil {
		log.Println("decode secret error", err)
		return
	}

	return
}

func (db *DbSessionStore) Del(id []byte) {
	sid := hex.EncodeToString(id)
	s := db.get(path.Join(db.Root, sid), false)
	if s == nil {
		return
	}

	//os.Remove(path.Join(db.Root, sid))
	sessiondb.DeleteSessionID(sid)
	//os.Remove(path.Join(db.Root, s.Addr))
}

func (db *DbSessionStore) Clean() error {
	files, err := os.ReadDir(db.Root)
	if err != nil {
		return err
	}

	for _, f := range files {
		db.get(path.Join(db.Root, f.Name()), true)
	}

	return nil
}
