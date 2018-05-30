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

var ListenAndServe = http.ListenAndServe

func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tag, err := h.tagStore.GetByID(id)
	if err != nil {
		ResponseWithJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	ResponseWithJSON(w, &tag, http.StatusOK)
}

func (h *TagHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var tags []*app.Tag
	tags, err := h.tagStore.GetList()
	if err != nil {
		ResponseWithJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	ResponseWithJSON(w, tags, http.StatusOK)
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tmp app.Tag
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		ResponseWithJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Create Tag
	tag, err := h.tagStore.Create(tmp)
	if err != nil {
		ResponseWithJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	ResponseWithJSON(w, &tag, http.StatusOK)
}
