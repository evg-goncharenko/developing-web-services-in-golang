package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h *Handler) HandleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{"wrong method"}
		bytes, _ := json.Marshal(&resp)
		w.Write(bytes)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("cookie err: %s", err)
		w.Header().Set("Content-Type", "application/json")
		resp := Response{"no cookie"}
		bytes, _ := json.Marshal(&resp)
		w.Write(bytes)
		return
	}

	cookie.Expires = time.Now().Add(-1)

	http.SetCookie(w, cookie)

	h.Mu.Lock()
	delete(h.Sessions, cookie.Value)
	h.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
