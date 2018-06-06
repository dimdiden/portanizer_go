package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestGetTagHandlers(t *testing.T) {
	t.Run("GetTag", testGetTag)
	t.Run("GetTagList", testGetTagList)
}

func testGetTag(t *testing.T) {

	var ts mock.TagStore

	ts.GetIdFn = func(id string) (*app.Tag, error) {
		if id != "100" {
			t.Fatalf("unexpected id: %v", id)
		}
		return &app.Tag{ID: 100, Name: "Tag100"}, nil
	}

	handler := NewTagServer(&ts).router

	w, err := sendRequest(handler, "GET", "/tags/100", nil)
	ok(t, err)
	equals(t, `{"ID":100,"Name":"Tag100"}`, w)

	if !ts.GetIdInvoked {
		t.Fatal("expected TagStore to be invoked")
	}
}

func testGetTagList(t *testing.T) {

	var ts mock.TagStore

	ts.GetListFn = func() ([]*app.Tag, error) {
		return []*app.Tag{&app.Tag{ID: 1, Name: "Tag1"}, &app.Tag{ID: 2, Name: "Tag2"}}, nil
	}

	handler := NewTagServer(&ts).router

	w, err := sendRequest(handler, "GET", "/tags", nil)
	ok(t, err)
	equals(t, `[{"ID":1,"Name":"Tag1"},{"ID":2,"Name":"Tag2"}]`, w)

	if !ts.GetListInvoked {
		t.Fatal("expected TagStore to be invoked")
	}
}

func TestCreateTagHandlers(t *testing.T) {
	t.Run("CreateTag", testCreateTag)
	t.Run("CreateWithEmptyTagName", testCreateWithEmptyTagName)
	t.Run("CreateWithUnknownTagField", testCreateWithUnknownTagField)
	t.Run("CreateWithExistingTagField", testCreateWithExistingTagField)
}

func testCreateTag(t *testing.T) {

	var ts mock.TagStore

	ts.CreateFn = func(tag app.Tag) (*app.Tag, error) {
		if tag.Name != "Tag1" {
			t.Fatalf("unexpected tag Name: %v", tag.Name)
		}
		return &app.Tag{ID: 1, Name: "Tag1"}, nil
	}

	handler := NewTagServer(&ts).router
	body := strings.NewReader(`{"Name": "Tag1"}`)

	w, err := sendRequest(handler, "POST", "/tags", body)
	ok(t, err)
	equals(t, `{"ID":1,"Name":"Tag1"}`, w)

	if !ts.CreateInvoked {
		t.Fatal("expected TagStore to be invoked")
	}
}

func testCreateWithEmptyTagName(t *testing.T) {

	var ts mock.TagStore

	ts.CreateFn = func(tag app.Tag) (*app.Tag, error) {
		if tag.Name != "" {
			t.Fatalf("unexpected tag Name: %v", tag.Name)
		}
		return nil, app.ErrEmpty
	}

	handler := NewTagServer(&ts).router
	body := strings.NewReader(`{"Name": ""}`)

	w, err := sendRequest(handler, "POST", "/tags", body)
	ok(t, err)
	equals(t, `{"Message":"Record has empty field"}`, w)

	if ts.CreateInvoked {
		t.Fatal("expected TagStore NOT to be invoked")
	}
}

func testCreateWithUnknownTagField(t *testing.T) {

	var ts mock.TagStore

	ts.CreateFn = func(tag app.Tag) (*app.Tag, error) {
		if &tag.Name != nil {
			t.Fatalf("Tag Name field should be incorrect in this test case")
		}
		return nil, nil
	}

	handler := NewTagServer(&ts).router
	body := strings.NewReader(`{"Nam": "Tag1"}`)

	w, err := sendRequest(handler, "POST", "/tags", body)
	ok(t, err)
	equals(t, `{"Message":"json: unknown field \"Nam\""}`, w)

	if ts.CreateInvoked {
		t.Fatal("expected TagStore NOT to be invoked")
	}
}

func testCreateWithExistingTagField(t *testing.T) {
	var ts mock.TagStore

	ts.CreateFn = func(tag app.Tag) (*app.Tag, error) {
		if tag.Name != "Tag1" {
			t.Fatalf("unexpected tag Name: %v", tag.Name)
		}
		return nil, app.ErrExists
	}

	handler := NewTagServer(&ts).router
	body := strings.NewReader(`{"Name": "Tag1"}`)

	w, err := sendRequest(handler, "POST", "/tags", body)
	ok(t, err)
	equals(t, `{"Message":"Record already exists in the database"}`, w)

	if !ts.CreateInvoked {
		t.Fatal("expected TagStore to be invoked")
	}
}

func TestUpdateTagHandlers(t *testing.T) {
	t.Run("UpdateTag", testUpdateTag)
}

func testUpdateTag(t *testing.T) {

	var ts mock.TagStore

	ts.UpdateFn = func(id string, tag app.Tag) (*app.Tag, error) {
		if id != "1" {
			t.Fatalf("unexpected id: %v", id)
		}
		if tag.Name != "Tag2" {
			t.Fatalf("unexpected tag Name: %v", tag.Name)
		}
		return &app.Tag{ID: 1, Name: "Tag2"}, nil
	}

	handler := NewTagServer(&ts).router
	body := strings.NewReader(`{"Name": "Tag2"}`)

	w, err := sendRequest(handler, "PATCH", "/tags/1", body)
	ok(t, err)
	equals(t, `{"ID":1,"Name":"Tag2"}`, w)

	if !ts.UpdateInvoked {
		t.Fatal("expected TagStore to be invoked")
	}
}

func sendRequest(handler http.Handler, method, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	handler.ServeHTTP(w, r)
	return w, nil
}

func ok(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func equals(t *testing.T, exp string, w *httptest.ResponseRecorder) {
	result, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatalf("can not read response body")
	}
	if exp != string(result) {
		t.Fatalf("unexpected result: %v", string(result))
	}
}
