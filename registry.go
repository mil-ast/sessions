package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
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

func (this reg) create_ses_id() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (this *reg) add(ses_id string, ses *Session) {
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

func (this *reg) remove_old(time_now time.Time) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	for key, ses := range this.sessions {
		if !ses.is_actual(time_now) {
			delete(this.sessions, key)
		}
	}
}
