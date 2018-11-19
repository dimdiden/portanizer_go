package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dimdiden/portanizer_go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	atokenExp = 1
	rtokenExp = 4
)

// Server is the Rest API Server
type Server struct {
	asecret []byte
	rsecret []byte
	// auth    *authHandler
	User   *userHandler
	Post   *postHandler
	Tag    *tagHandler
	router *mux.Router

	logout io.Writer
}

// New will construct a Server and apply all of the necessary routes
func New(as, rs []byte, logout io.Writer, pr portanizer.PostRepo, tr portanizer.TagRepo, ur portanizer.UserRepo) *Server {
	server := Server{
		Post:   &postHandler{repo: pr},
		Tag:    &tagHandler{repo: tr},
		User:   newUserHandler(as, rs, ur),
		router: mux.NewRouter(),
		logout: logout,
	}
	server.postroutes()
	server.tagroutes()
	server.userroutes()

	return &server
}

// Run starts the server on the provided port
func (s *Server) Run(port string) error {
	fmt.Fprintf(s.logout, "[[> listening on %v port...\n", port)
	err := http.ListenAndServe(":"+port, s)
	return err
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := handlers.LoggingHandler(s.logout, s.router)
	handler.ServeHTTP(w, r)
}

func (s *Server) postroutes() {
	s.router.Handle("/posts",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.GetList))).Methods("GET")
	s.router.Handle("/posts",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.Create))).Methods("POST")

	s.router.Handle("/posts/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.Get))).Methods("GET")
	s.router.Handle("/posts/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.Update))).Methods("PATCH")
	s.router.Handle("/posts/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.Delete))).Methods("DELETE")
	s.router.Handle("/posts/{id}/tags",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Post.PutTags))).Methods("PUT")
}

func (s *Server) tagroutes() {
	s.router.Handle("/tags",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Tag.GetList))).Methods("GET")
	s.router.Handle("/tags",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Tag.Create))).Methods("POST")

	s.router.Handle("/tags/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Tag.Get))).Methods("GET")
	s.router.Handle("/tags/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Tag.Update))).Methods("PATCH")
	s.router.Handle("/tags/{id}",
		s.User.jwtMiddleware.Handler(http.HandlerFunc(s.Tag.Delete))).Methods("DELETE")
}

func (s *Server) userroutes() {
	s.router.Handle("/users", s.User.SignUp()).Methods("POST")
	s.router.Handle("/users/signin", s.User.SignIn()).Methods("POST")
	s.router.Handle("/users/refresh",
		s.User.jwtMiddleware.Handler(s.User.Refresh())).Methods("POST")
	s.router.Handle("/users/signout",
		s.User.jwtMiddleware.Handler(s.User.SignOut())).Methods("POST")
}
