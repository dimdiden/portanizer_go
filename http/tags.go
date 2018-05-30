package http

import (
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
