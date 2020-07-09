package sessions

import (
	"net/http"
	"sync"
)

type Manager struct {
	data   map[string]*Session
	mutexx *sync.RWMutex
}

func NewSessions() *Manager {
	return &Manager{
		data:   make(map[string]*Session, 10),
		mutexx: &sync.RWMutex{},
	}
}

func (sm *Manager) Check(r *http.Request) (*Session, error) {
	sessionToken := r.Header.Get("Authorization")

	if sessionToken == "" {
		return nil, ErrAuth
	}

	if len(sessionToken) > 7 {
		sessionToken = sessionToken[7:]
	}

	sm.mutexx.RLock()
	sess, ok := sm.data[sessionToken]
	sm.mutexx.RUnlock()

	if !ok {
		return nil, ErrAuth
	}
	return sess, nil
}

func (sm *Manager) Create(w http.ResponseWriter, userLogin string, token string) (*Session, error) {
	session := NewSession(userLogin, token)
	sm.mutexx.Lock()
	sm.data[session.ID] = session
	sm.mutexx.Unlock()
	return session, nil
}

func (sm *Manager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	session, err := SessionFromContext(r.Context())

	if err != nil {
		return err
	}

	sm.mutexx.Lock()
	delete(sm.data, session.ID)
	sm.mutexx.Unlock()
	return nil
}
