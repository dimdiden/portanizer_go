package server

import (
	"encoding/json"
	"net/http"
)

type Info struct {
	Message string
}

func renderJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	switch data := data.(type) {
	case string:
		info := &Info{Message: data}
		res, err := json.Marshal(info)
		if err != nil {
			renderJSON(w, "Can not marshal output", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	default:
		res, err := json.Marshal(data)
		if err != nil {
			renderJSON(w, "Can not marshal output", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
