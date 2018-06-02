package http

import (
	"encoding/json"
	"net/http"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/gorilla/mux"
)

type TagHandler struct {
	tagStore app.TagStore
}

func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tag, err := h.tagStore.GetByID(id)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) GetList(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagStore.GetList()
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, tags, http.StatusOK)
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tmp app.Tag
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Create Tag
	tag, err := h.tagStore.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var tmp app.Tag
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		renderJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Create Tag
	tag, err := h.tagStore.Update(id, tmp)
	if err != nil {
		if err == app.ErrNotFound {
			renderJSON(w, err.Error(), http.StatusNotFound)
			return
		}
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.tagStore.Delete(id)
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
