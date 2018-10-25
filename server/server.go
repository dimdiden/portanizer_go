package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dimdiden/portanizer_go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	atokenExp = 1
	rtokenExp = 4
)

// ASecret is used to issue access tokens
var ASecret []byte

// RSecret is used to issue access tokens
var RSecret []byte

// Server is the Rest API Server
type Server struct {
	auth   *authHandler
	user   *userHandler
	post   *postHandler
	tag    *tagHandler
	router *mux.Router

	logout io.Writer
}

// New will construct a Server and apply all of the necessary routes
func New(logout io.Writer, pr portanizer.PostRepo, tr portanizer.TagRepo, ur portanizer.UserRepo) *Server {
	server := Server{
		post:   &postHandler{repo: pr},
		tag:    &tagHandler{repo: tr},
		user:   &userHandler{repo: ur},
		auth:   &authHandler{repo: ur},
		router: mux.NewRouter(),
		logout: logout,
	}
	server.postroutes()
	server.tagroutes()
	server.authroutes()
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
	handler := handlers.LoggingHandler(os.Stdout, s.router)
	handler.ServeHTTP(w, r)
}

func (s *Server) postroutes() {
	s.router.Handle("/posts",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.GetList))).Methods("GET")
	s.router.Handle("/posts",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.Create))).Methods("POST")

	s.router.Handle("/posts/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.Get))).Methods("GET")
	s.router.Handle("/posts/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.Update))).Methods("PATCH")
	s.router.Handle("/posts/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.Delete))).Methods("DELETE")

	s.router.Handle("/posts/{id}/tags",
		jwtMiddleware.Handler(http.HandlerFunc(s.post.PutTags))).Methods("PUT")
}

func (s *Server) tagroutes() {
	s.router.Handle("/tags",
		jwtMiddleware.Handler(http.HandlerFunc(s.tag.GetList))).Methods("GET")
	s.router.Handle("/tags",
		jwtMiddleware.Handler(http.HandlerFunc(s.tag.Create))).Methods("POST")

	s.router.Handle("/tags/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.tag.Get))).Methods("GET")
	s.router.Handle("/tags/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.tag.Update))).Methods("PATCH")
	s.router.Handle("/tags/{id}",
		jwtMiddleware.Handler(http.HandlerFunc(s.tag.Delete))).Methods("DELETE")
}

func (s *Server) userroutes() {
	s.router.HandleFunc("/users", s.user.Register).Methods("POST")
}

func (s *Server) authroutes() {
	s.router.HandleFunc("/login", s.auth.Login).Methods("POST")
	s.router.Handle("/refresh",
		jwtMiddleware.Handler(http.HandlerFunc(s.auth.Refresh))).Methods("POST")
}
