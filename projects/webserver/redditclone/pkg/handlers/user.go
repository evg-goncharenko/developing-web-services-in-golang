package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"redditclone/pkg/sessions"
	"redditclone/pkg/users"

	"go.uber.org/zap"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	UserRep  *users.UserRep
	Sessions *sessions.Manager
}

const (
	CantReading         = `{"message":"Error with reading"}`
	AlreadyRegistration = `{"message":"This account is already registered"}`
	FailedAuthorization = `{"message":"Invalid username or password"}`
	UserError           = `{"message":"Error with user"}`
)

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {
	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if errRead != nil {
		http.Error(w, CantReading, http.StatusInternalServerError)
		return
	}

	_, errAuthorization := sessions.SessionFromContext(r.Context())

	if errAuthorization == nil {
		http.Error(w, AlreadyRegistration, http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{}
	json.Unmarshal(body, &data)
	login, takeLogin := data["username"]
	password, takePass := data["password"]

	if !takeLogin || !takePass {
		http.Error(w, FailedAuthorization, http.StatusBadRequest)
		return
	}

	user, err := h.UserRep.Registration(login.(string), password.(string))

	if err != nil {
		http.Error(w, FailedAuthorization, http.StatusBadRequest)
		return
	}

	h.Sessions.Create(w, user.Login, user.GetToken())
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, string(user.GetToken()))))
	return
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	_, errAuthorization := sessions.SessionFromContext(r.Context())

	if errAuthorization == nil {
		http.Error(w, AlreadyRegistration, http.StatusBadRequest)
		return
	}

	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if errRead != nil {
		http.Error(w, CantReading, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{}
	json.Unmarshal(body, &data)
	login, hasLogin := data["username"]
	password, hasPass := data["password"]

	if !hasLogin || !hasPass {
		http.Error(w, UserError, http.StatusBadRequest)
		return
	}

	fmt.Println(data)
	fmt.Println("Before authorization")
	user, err := h.UserRep.Authorize(login.(string), password.(string))

	if err != nil {
		http.Error(w, FailedAuthorization, http.StatusBadRequest)
		return
	}

	fmt.Println("After authorization")
	token := user.GetToken()
	h.Sessions.Create(w, user.Login, token)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, string(token))))
	return
}
