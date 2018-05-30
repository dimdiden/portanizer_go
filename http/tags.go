package http

import (
	"encoding/json"
	"fmt"
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

func ResponseWithJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	switch data := data.(type) {
	case string:
		str := fmt.Sprintf("{\"Message\": \"%s\"}", data)
		res := []byte(str)
		w.Write(res)
	default:
		res, err := json.Marshal(data)
		if err != nil {
			ResponseWithJSON(w, "Can not marshal output", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
