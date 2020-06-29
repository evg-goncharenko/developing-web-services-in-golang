package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/posts"
	"redditclone/pkg/sessions"
	"redditclone/pkg/users"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	Logger   *zap.SugaredLogger
	UserRep  *users.UserRep
	PostRep  *posts.PostRep
	Sessions *sessions.Manager
}

const (
	BadPost         = `{"message":"Error with post"}`
	BadIdPost       = `{"message":"Error with id post"}`
	BadResearchPost = `{"message":"Error with finding post"}`
	CantAddPost     = `{"message":"Error with adding post"}`
	CantGetPost     = `{"message":"Error with getting post"}`
	CantDeletePost  = `{"message":"Error with deletion post"}`
	CantReadPost    = `{"message":"Error with reading post"}`

	BadComment      = `{"message":"Error with comment"}`
	CantAddComment  = `{"message":"Error with adding comment"}`
	BadIdComment    = `{"message":"Error with id comment"}`
	CantFindComment = `{"message":"Error with finding comment"}`

	BadCategory = `{"Error with category"}`

	BadUsername      = `{"message":"Error with getting username"}`
	BadUser          = `{"message":"Error with user"}`
	CantRead         = `{"message":"Error with reading"}`
	CantDeleteAuthor = `{"message":"Error with deletion author"}`
	BadVote          = `{"message":"Error with Vote"}`
)

func (ph *PostHandler) ListOfAllPosts(w http.ResponseWriter, r *http.Request) {
	post, err := ph.PostRep.GetAll()

	if err != nil {
		http.Error(w, BadPost, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(post)
	w.Write(data)
}

func (ph *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["POST_ID"], 10, 64)

	if err != nil {
		http.Error(w, BadIdPost, http.StatusBadRequest)
		return
	}

	sess, _ := sessions.SessionFromContext(r.Context())
	user, _ := ph.UserRep.Get(sess.Login)
	post, errP := ph.PostRep.Get(posts.PostIdType(id))

	if errP != nil {
		http.Error(w, BadResearchPost, http.StatusInternalServerError)
	}

	if post.Author.ID != user.ID {
		http.Error(w, CantDeletePost, http.StatusBadRequest)
	}

	ph.PostRep.DeletePost(post.ID)
	w.Write([]byte(`{"message":"success"}`))
}

func (ph *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if errRead != nil {
		http.Error(w, CantReadPost, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	idInt, err := strconv.ParseUint(vars["POST_ID"], 10, 64)

	if err != nil {
		http.Error(w, BadIdPost, http.StatusBadRequest)
		return
	}

	id := posts.PostIdType(idInt)
	commentData := make(map[string]interface{})
	json.Unmarshal(body, &commentData)
	message, hasComment := commentData["comment"]

	if _, isStr := message.(string); !hasComment || !isStr {
		http.Error(w, BadComment, http.StatusBadRequest)
		return
	}

	sess, _ := sessions.SessionFromContext(r.Context())
	u, _ := ph.UserRep.Get(sess.Login)
	comment := &posts.Comment{Author: u, Body: message.(string), Created: time.Now()}
	_, err = ph.PostRep.AddComment(id, comment)

	if err != nil {
		http.Error(w, CantAddComment, http.StatusInternalServerError)
		return
	}

	post, _ := ph.PostRep.Get(id)
	data, _ := json.Marshal(post)
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (ph *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.SessionFromContext(r.Context())
	user, _ := ph.UserRep.Get(session.Login)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["POST_ID"], 10, 64)

	if err != nil {
		http.Error(w, BadIdPost, http.StatusBadRequest)
		return
	}

	postId := posts.PostIdType(id)
	cId, err := strconv.ParseUint(vars["COMMENT_ID"], 10, 32)

	commentId := uint32(cId)

	if err != nil {
		http.Error(w, BadIdComment, http.StatusBadRequest)
		return
	}

	comment, err := ph.PostRep.GetComment(postId, commentId)

	if err != nil {
		http.Error(w, CantFindComment, http.StatusInternalServerError)
		return
	}

	if comment.Author.ID != user.ID {
		http.Error(w, CantDeleteAuthor, http.StatusBadRequest)
		return
	}

	ph.PostRep.DeleteComment(postId, commentId)
	p, _ := ph.PostRep.Get(postId)
	data, _ := json.Marshal(p)
	w.Write(data)
}

func (ph *PostHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catName, ok := vars["CATEGORY_NAME"]

	if !ok {
		http.Error(w, BadCategory, http.StatusBadRequest)
		return
	}

	category := posts.NameToType(catName)
	post, err := ph.PostRep.GetCategoryPosts(category)

	if err != nil {
		http.Error(w, CantGetPost, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(post)
	w.Write(data)
}

func (ph *PostHandler) UserList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Login, ok := vars["USER_LOGIN"]

	if !ok {
		http.Error(w, BadUsername, http.StatusBadRequest)
		return
	}

	user, okUser := ph.UserRep.Get(Login)

	if !okUser {
		http.Error(w, BadUser, http.StatusBadRequest)
		return
	}

	post, err := ph.PostRep.GetUserPosts(user.ID)

	if err != nil {
		http.Error(w, CantGetPost, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(post)
	w.Write(data)
}

func (ph *PostHandler) Add(w http.ResponseWriter, r *http.Request) {
	body, errRead := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if errRead != nil {
		http.Error(w, CantRead, http.StatusInternalServerError)
		return
	}

	sess, _ := sessions.SessionFromContext(r.Context())
	user, _ := ph.UserRep.Get(sess.Login)
	p := posts.NewPost()

	if json.Unmarshal(body, p) != nil {
		http.Error(w, BadPost, http.StatusBadRequest)
		return
	}

	p.Author = user
	p.AddVote(&posts.Vote{Author: user.ID, Vote: posts.RatingUp})
	p.Score = 1
	p.Created = time.Now()
	_, err := ph.PostRep.Add(p)

	if err != nil {
		http.Error(w, CantAddPost, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(p)
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (ph *PostHandler) Details(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["POST_ID"], 10, 64)

	if err != nil {
		http.Error(w, BadIdPost, http.StatusBadRequest)
		return
	}

	post, err := ph.PostRep.Get(posts.PostIdType(id))

	if err != nil {
		http.Error(w, BadPost, http.StatusNotFound)
		return
	}

	post.Views++
	ph.PostRep.Update(post.ID, post)
	data, _ := json.Marshal(post)
	w.Write(data)
}

func (ph *PostHandler) RatingUp(w http.ResponseWriter, r *http.Request) {
	ph.Vote(w, r, posts.RatingUp)
}

func (ph *PostHandler) RatingDown(w http.ResponseWriter, r *http.Request) {
	ph.Vote(w, r, posts.RatingDown)
}

func (ph *PostHandler) Vote(w http.ResponseWriter, r *http.Request, v posts.VoteType) {
	session, _ := sessions.SessionFromContext(r.Context())
	user, _ := ph.UserRep.Get(session.Login)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["POST_ID"], 10, 64)

	if err != nil {
		http.Error(w, BadIdPost, http.StatusBadRequest)
		return
	}

	postId := posts.PostIdType(id)
	post, err := ph.PostRep.Vote(postId, user.ID, v)

	if err != nil {
		http.Error(w, BadVote, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(post)
	w.Write(data)
}
