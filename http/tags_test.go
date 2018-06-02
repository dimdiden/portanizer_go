package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/dimdiden/portanizer_sop/mock"
	"github.com/gorilla/mux"
)

func NewTagServer(ts app.TagStore) *Server {
	server := Server{
		tag:    &TagHandler{tagStore: ts},
		router: mux.NewRouter(),
	}
	server.tagroutes()
	return &server
}

func TestTagHandler(t *testing.T) {

	var ts mock.TagStore

	ts.GetIdFn = func(id string) (*app.Tag, error) {
		if id != "100" {
			t.Fatalf("unexpected id: %v", id)
		}
		return &app.Tag{ID: 100, Name: "Tag100"}, nil
	}

	router := NewTagServer(&ts).router

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/tags/100", nil)

	router.ServeHTTP(w, r)

	if !ts.GetIdInvoked {
		t.Fatal("expected Tag() to be invoked")
	}
}
