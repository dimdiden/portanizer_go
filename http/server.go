package http

import (
	"net/http"
	"os"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	post   *PostHandler
	tag    *TagHandler
	router *mux.Router

	logOn bool
}

var ListenAndServe = http.ListenAndServe

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
func NewServer(ts app.TagStore, ps app.PostStore) *Server {
	server := Server{
		tag:    &TagHandler{tagStore: ts},
		post:   &PostHandler{postStore: ps},
		router: mux.NewRouter(),
	}
	server.tagroutes()
	server.postroutes()

	return &server
}

func (s *Server) LogHttpEnable() {
	s.logOn = true
}

func (s *Server) tagroutes() {
	s.router.HandleFunc("/tags", s.tag.GetList).Methods("GET")
	s.router.HandleFunc("/tags", s.tag.Create).Methods("POST")

	s.router.HandleFunc("/tags/{id}", s.tag.Get).Methods("GET")
	s.router.HandleFunc("/tags/{id}", s.tag.Update).Methods("PATCH")
	s.router.HandleFunc("/tags/{id}", s.tag.Delete).Methods("DELETE")
}

func (s *Server) postroutes() {
	s.router.HandleFunc("/posts", s.post.GetList).Methods("GET")
	s.router.HandleFunc("/posts", s.post.Create).Methods("POST")

	s.router.HandleFunc("/posts/{id}", s.post.Get).Methods("GET")
	s.router.HandleFunc("/posts/{id}", s.post.Update).Methods("PATCH")
	s.router.HandleFunc("/posts/{id}", s.post.Delete).Methods("DELETE")

	s.router.HandleFunc("/posts/{id}/tags", s.post.PutTags).Methods("PUT")
}
