package server

import (
	"encoding/json"
	"net/http"

	portanizer "github.com/dimdiden/portanizer_go"
	"github.com/gorilla/mux"
)

type postHandler struct {
	repo portanizer.PostRepo
}

func (h *postHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	post, err := h.repo.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *postHandler) GetList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repo.GetList()
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, posts, http.StatusOK)
}

func (h *postHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tmp portanizer.Post
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	// Read the request body
	if err := decoder.Decode(&tmp); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tmp.IsValid() {
		renderJSON(w, portanizer.ErrEmpty.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.repo.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *postHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var tmp portanizer.Post
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	// Read the request body
	if err := decoder.Decode(&tmp); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tmp.IsValid() {
		renderJSON(w, portanizer.ErrEmpty.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.repo.Update(id, tmp)
	if err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *postHandler) PutTags(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var tagids []string
	if err := json.NewDecoder(r.Body).Decode(&tagids); err != nil { // <= need to add check for non-numeral values
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	post, err := h.repo.PutTags(id, tagids)
	if err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *postHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.repo.Delete(id); err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "Post "+id+" has been deleted successfully", http.StatusOK)
}
