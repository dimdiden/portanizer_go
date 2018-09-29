package http

import (
	"encoding/json"
	"net/http"

	"github.com/dimdiden/portanizer"
	"github.com/gorilla/mux"
)

type TagHandler struct {
	tagRepo portanizer.TagRepo
}

func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tag, err := h.tagRepo.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) GetList(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagRepo.GetList()
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, tags, http.StatusOK)
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	tag, err := h.tagRepo.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	tag, err := h.tagRepo.Update(id, tmp)
	if err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.tagRepo.Delete(id); err == portanizer.ErrNotFound {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "Tag "+id+" has been deleted successfully", http.StatusOK)
}
