package http

import (
	"encoding/json"
	"net/http"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	postStore app.PostStore
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	post, err := h.postStore.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *PostHandler) GetList(w http.ResponseWriter, r *http.Request) {
	// var posts []*app.Post
	posts, err := h.postStore.GetList()
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, posts, http.StatusOK)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tmp app.Post
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	post, err := h.postStore.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var tmp app.Post
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	post, err := h.postStore.Update(id, tmp)
	if err != nil {
		if err == app.ErrNotFound {
			renderJSON(w, err.Error(), http.StatusNotFound)
			return
		}
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *PostHandler) PutTags(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var tagids []string
	if err := json.NewDecoder(r.Body).Decode(&tagids); err != nil {
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	post, err := h.postStore.PutTags(id, tagids)
	if err != nil {
		if err == app.ErrNotFound {
			renderJSON(w, err.Error(), http.StatusNotFound)
			return
		}
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.postStore.Delete(id)
	if err != nil {
		if err == app.ErrNotFound {
			renderJSON(w, err.Error(), http.StatusNotFound)
			return
		}
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "Tag "+id+" has been deleted successfully", http.StatusOK)
}
