package server

import (
	"encoding/json"
	"net/http"

	"github.com/dimdiden/portanizer"
	"github.com/gorilla/mux"
)

type postHandler struct {
	postRepo portanizer.PostRepo
}

func (h *postHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	post, err := h.postRepo.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &post, http.StatusOK)
}

func (h *postHandler) GetList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postRepo.GetList()
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

	post, err := h.postRepo.Create(tmp)
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

	post, err := h.postRepo.Update(id, tmp)
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
	post, err := h.postRepo.PutTags(id, tagids)
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

	if err := h.postRepo.Delete(id); err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "Post "+id+" has been deleted successfully", http.StatusOK)
}
