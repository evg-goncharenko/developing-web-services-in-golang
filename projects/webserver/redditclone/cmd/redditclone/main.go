package main

import (
	"net/http"

	"redditclone/pkg/handlers"
	"redditclone/pkg/interfaces"
	"redditclone/pkg/middleware"
	"redditclone/pkg/posts"
	"redditclone/pkg/sessions"
	"redditclone/pkg/users"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	sm := sessions.NewSessions()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	userRep := users.NewUserRep()
	postsRep := posts.NewPostRep()

	var uh interfaces.IUserHandler = &handlers.UserHandler{
		Logger:   logger,
		UserRep:  userRep,
		Sessions: sm,
	}
	var ph interfaces.IPostHandler = &handlers.PostHandler{
		Logger:   logger,
		UserRep:  userRep,
		PostRep:  postsRep,
		Sessions: sm,
	}

	ar := middleware.AuthRequired
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("template")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("template/static"))))

	restRouter := mux.NewRouter()
	restRouter.HandleFunc("/register", uh.Registration).Methods("POST")
	restRouter.HandleFunc("/login", uh.Login).Methods("POST")

	restRouter.HandleFunc("/posts/", ph.ListOfAllPosts).Methods("GET")
	restRouter.HandleFunc("/posts", ar(ph.Add)).Methods("POST")

	restRouter.HandleFunc("/posts/{CATEGORY_NAME}", ph.CategoryList).Methods("GET")
	restRouter.HandleFunc("/post/{POST_ID}", ph.Details).Methods("GET")

	restRouter.HandleFunc("/post/{POST_ID}", ar(ph.AddComment)).Methods("POST")
	restRouter.HandleFunc("/post/{POST_ID}/{COMMENT_ID}", ar(ph.DeleteComment)).Methods("DELETE")

	restRouter.HandleFunc("/post/{POST_ID}/upvote", ar(ph.RatingUp)).Methods("GET")
	restRouter.HandleFunc("/post/{POST_ID}/downvote", ar(ph.RatingDown)).Methods("GET")

	restRouter.HandleFunc("/post/{POST_ID}", ar(ph.DeletePost)).Methods("DELETE")
	restRouter.HandleFunc("/user/{USER_LOGIN}", ph.UserList).Methods("GET")

	r.PathPrefix("/api").Handler(middleware.JsonContent(http.StripPrefix("/api", restRouter)))

	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	http.ListenAndServe(addr, mux)
}
