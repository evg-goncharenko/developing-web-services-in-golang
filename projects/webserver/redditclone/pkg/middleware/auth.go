package middleware

import (
	"context"
	"fmt"
	"net/http"

	"redditclone/pkg/sessions"
)

func Auth(sm *sessions.Manager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")
		sess, err := sm.Check(r)

		if err == nil {
			ctx := context.WithValue(r.Context(), sessions.SessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, errAuth := sessions.SessionFromContext(r.Context())

		if errAuth != nil {
			http.Error(w, `{"message":"Failed to authenticate"}`, http.StatusNonAuthoritativeInfo)
			return
		}
		next(w, r)
	})
}
