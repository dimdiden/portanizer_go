package http

import (
	"net/http"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/gorilla/mux"
)

type Server struct {
	tags   *TagHandler
	router *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer will construct a Server and apply all of the necessary routes
func NewServer(ts app.TagStore) *Server {
	server := Server{
		tags: &TagHandler{
			tagStore: ts,
		},
		router: mux.NewRouter(),
	}
	server.routes()
	return &server
}

func (s *Server) routes() {
	s.router.HandleFunc("/tags", s.tags.GetList).Methods("GET")
	s.router.HandleFunc("/tags/{id}", s.tags.Get).Methods("GET")

	s.router.HandleFunc("/tags", s.tags.Create).Methods("POST")
	s.router.HandleFunc("/tags/{id}", s.tags.Update).Methods("PATCH")
	s.router.HandleFunc("/tags/{id}", s.tags.Delete).Methods("DELETE")
}
