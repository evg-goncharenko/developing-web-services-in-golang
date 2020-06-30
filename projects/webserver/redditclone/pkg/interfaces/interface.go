package interfaces

import "net/http"

type IUserHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type IPostHandler interface {
	ListOfAllPosts(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)

	CategoryList(w http.ResponseWriter, r *http.Request)
	Details(w http.ResponseWriter, r *http.Request)

	AddComment(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)

	RatingUp(w http.ResponseWriter, r *http.Request)
	RatingDown(w http.ResponseWriter, r *http.Request)

	DeletePost(w http.ResponseWriter, r *http.Request)
	UserList(w http.ResponseWriter, r *http.Request)
}
