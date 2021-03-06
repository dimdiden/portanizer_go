package server

import (
	"encoding/json"
	"net/http"

	"github.com/dimdiden/portanizer_go"
	"github.com/gorilla/mux"
)

type tagHandler struct {
	repo portanizer.TagRepo
}

func (h *tagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tag, err := h.repo.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *tagHandler) GetList(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.GetList()
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, tags, http.StatusOK)
}

func (h *tagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tmp portanizer.Tag

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
	// Create Tag
	tag, err := h.repo.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *tagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var tmp portanizer.Tag
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
	// Create Tag
	tag, err := h.repo.Update(id, tmp)
	if err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *tagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.repo.Delete(id); err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "Tag "+id+" has been deleted successfully", http.StatusOK)
}
