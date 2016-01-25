package sessions

import (
	"crypto/rand"
	"encoding/base64"
	//"errors"
	"io"
	"sync"
	"time"
)

type reg struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

var registry *reg

func init() {
	registry = new(reg)
	registry.sessions = make(map[string]*Session)
}

func (this *reg) new() (string, *Session) {
	var (
		ses_id string
		exists bool
	)
	for {
		ses_id = this.create_ses_id()

		this.mu.RLock()
		_, exists = this.sessions[ses_id]
		this.mu.RUnlock()

		if !exists {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	ses := this.create(ses_id)

	return ses_id, ses
}

func (this *reg) create(ses_id string) *Session {
	ses := new(Session)
	ses.last_time = time.Now()
	ses.values = make(map[string]interface{})

	this.mu.Lock()
	this.sessions[ses_id] = ses
	this.mu.Unlock()

	return ses
}

func (this reg) get(ses_id string) (*Session, bool) {
	this.mu.RLock()
	ses, is_exists := this.sessions[ses_id]
	this.mu.RUnlock()

	return ses, is_exists
}

func (this *reg) delete(ses_id string) {
	this.mu.Lock()
	delete(this.sessions, ses_id)
	this.mu.Unlock()
}

func (this reg) create_ses_id() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (this *reg) remove_old(time_now time.Time) {
	this.mu.Lock()
	defer this.mu.Unlock()

	for key, ses := range this.sessions {
		if !ses.is_actual(time_now) {
			delete(this.sessions, key)
		}
	}
}

/*


type reg struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}



func init() {
	registry = new(reg)
	registry.sessions = make(map[string]*Session)
}

func (this reg) create_ses_id() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (this *reg) set(ses_id string, ses *Session) {
	if ses_id != "" {
		this.mu.Lock()
		this.sessions[ses_id] = ses
		this.mu.Unlock()
	}
}

func (this *reg) delete(ses_id string) {
	this.mu.Lock()
	delete(this.sessions, ses_id)
	this.mu.Unlock()
}

func (this reg) get(ses_id string) (*Session, error) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if ses, ok := this.sessions[ses_id]; ok {
		return ses, nil
	}

	return nil, errors.New("not found")
}

*/
