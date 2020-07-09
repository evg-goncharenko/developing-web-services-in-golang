package sessions

import (
	"context"
	"errors"
)

type Session struct {
	ID    string
	Login string
}

type sessKey string

var (
	ErrAuth            = errors.New("No session found")
	SessionKey sessKey = "sessionKey"
)

func NewSession(Login string, token string) *Session {
	return &Session{
		ID:    token,
		Login: Login,
	}
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)

	if !ok || sess == nil {
		return nil, ErrAuth
	}
	return sess, nil
}
