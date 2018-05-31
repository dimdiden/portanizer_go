package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/dimdiden/portanizer_sop/mock"
)

func TestTagHandler(t *testing.T) {

	var ts mock.TagService
	// var ps mock.Pos

	ts.GetIdFn = func(id string) (*app.Tag, error) {
		if id != "100" {
			t.Fatalf("unexpected id: %v", id)
		}
		return &app.Tag{ID: 100, Name: "Tag100"}, nil
	}

	router := NewServer(&ts).router
	srv := httptest.NewServer(router)
	defer srv.Close()

	r, err := http.NewRequest("GET", srv.URL+"/tags/100", nil)
	if err != nil {
		t.Fatal("could not create request: ", err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatal("could make request: ", err)
	}

	fmt.Println(res.StatusCode)

	// Validate mock.
	if !ts.GetIdInvoked {
		t.Fatal("expected Tag() to be invoked")
	}
}
