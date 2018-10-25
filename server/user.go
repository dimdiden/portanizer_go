package server

import (
	"encoding/json"
	"net/http"

	"github.com/dimdiden/portanizer_go"
)

type userHandler struct {
	repo portanizer.UserRepo
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var tmp portanizer.User
	// Read and decode the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tmp); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := tmp.IsValid(); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.repo.Create(tmp)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, "user has been created", http.StatusOK)
}
