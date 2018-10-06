package server

import (
	"net/http"
	"os"

	"github.com/dimdiden/portanizer_go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	post   *postHandler
	tag    *tagHandler
	router *mux.Router

	logOn bool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	switch s.logOn {
	case true:
		handler = handlers.LoggingHandler(os.Stdout, s.router)
	default:
		handler = s.router
	}
	handler.ServeHTTP(w, r)
}

// NewServer will construct a Server and apply all of the necessary routes
func New(pr portanizer.PostRepo, tr portanizer.TagRepo) *Server {
	server := Server{
		post:   &postHandler{repo: pr},
		tag:    &tagHandler{repo: tr},
		router: mux.NewRouter(),
	}
	server.postroutes()
	server.tagroutes()

	return &server
}

func (s *Server) LogEnable() {
	s.logOn = true
}

func (s *Server) postroutes() {
	s.router.HandleFunc("/posts", s.post.GetList).Methods("GET")
	s.router.HandleFunc("/posts", s.post.Create).Methods("POST")

	s.router.HandleFunc("/posts/{id}", s.post.Get).Methods("GET")
	s.router.HandleFunc("/posts/{id}", s.post.Update).Methods("PATCH")
	s.router.HandleFunc("/posts/{id}", s.post.Delete).Methods("DELETE")

	s.router.HandleFunc("/posts/{id}/tags", s.post.PutTags).Methods("PUT")
}

func (s *Server) tagroutes() {
	s.router.HandleFunc("/tags", s.tag.GetList).Methods("GET")
	s.router.HandleFunc("/tags", s.tag.Create).Methods("POST")

	s.router.HandleFunc("/tags/{id}", s.tag.Get).Methods("GET")
	s.router.HandleFunc("/tags/{id}", s.tag.Update).Methods("PATCH")
	s.router.HandleFunc("/tags/{id}", s.tag.Delete).Methods("DELETE")
}
