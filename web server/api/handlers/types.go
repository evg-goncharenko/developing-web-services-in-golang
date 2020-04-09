package handlers

import "sync"

type Response struct {
	Message string `json:"message"`
}

type User struct {
	ID       int
	Login    string
	Password string
}

type Handler struct {
	Sessions map[string]*User
	Users    map[string]*User
	Mu       *sync.Mutex
}
