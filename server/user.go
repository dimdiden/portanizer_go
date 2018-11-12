package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dimdiden/portanizer_go"
)

type userHandler struct {
	asecret       []byte
	rsecret       []byte
	jwtMiddleware *jwtmiddleware.JWTMiddleware
	userRepo      portanizer.UserRepo
}

type tokenPair struct {
	AToken string `json:"atoken"`
	RToken string `json:"rtoken"`
}

func newUserHandler(as, rs []byte, ur portanizer.UserRepo) *userHandler {
	h := &userHandler{asecret: as, rsecret: rs, userRepo: ur}

	h.jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: h.appkeyFunc,
		SigningMethod:       jwt.SigningMethodHS256,
		ErrorHandler:        onError,
	})
	return h
}

func (h *userHandler) appkeyFunc(token *jwt.Token) (interface{}, error) {
	return h.asecret, nil
}

func onError(w http.ResponseWriter, r *http.Request, err string) {
	renderJSON(w, err, http.StatusUnauthorized)
}

func (h *userHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err := h.userRepo.Create(tmp)
		if err != nil {
			renderJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderJSON(w, "user has been created", http.StatusOK) // <= incorrect format in response
	}
}

func (h *userHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tmpUser portanizer.User
		// Read the request body
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		// Read the request body
		if err := decoder.Decode(&tmpUser); err != nil {
			renderJSON(w, "could not decode body: "+err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: do something with this validation
		if err := tmpUser.IsValid(); err != nil {
			renderJSON(w, err, http.StatusBadRequest)
			return
		}

		user, err := h.userRepo.GetByCreds(tmpUser.Email, tmpUser.Password)
		if err != nil {
			renderJSON(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tp, err := h.issueTokens(user)
		if err != nil {
			renderJSON(w, "could not issue tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		renderJSON(w, tp, http.StatusOK)
	}
}

func (h *userHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)

		user, err := h.userRepo.GetByID(uid)
		if err != nil {
			renderJSON(w, err.Error(), http.StatusNotFound)
			return
		}

		var tmpTpair tokenPair
		// Read the request body
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		// Read the request body
		if err := decoder.Decode(&tmpTpair); err != nil {
			renderJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		if tmpTpair.RToken == "" {
			renderJSON(w, "no rtoken in the request body", http.StatusBadRequest)
			return
		}

		if user.RToken != tmpTpair.RToken {
			renderJSON(w, "rtoken is invalid", http.StatusUnauthorized)
			h.userRepo.EmptyRToken(uid)
			return
		}

		tp, err := h.issueTokens(user)
		if err != nil {
			renderJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderJSON(w, tp, http.StatusOK)
	}
}

func (h *userHandler) SignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		err := h.userRepo.EmptyRToken(uid)
		if err != nil {
			renderJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		renderJSON(w, "rtoken has been removed", http.StatusOK)
	}
}

func (h *userHandler) issueTokens(user *portanizer.User) (*tokenPair, error) {
	// Create the token
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * atokenExp).Unix(),
		Subject:   strconv.Itoa(int(user.ID)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret
	atoken, err := token.SignedString(h.asecret)
	if err != nil {
		return nil, err
	}
	// Change expires to generate refresh token
	claims.ExpiresAt = time.Now().Add(time.Hour * rtokenExp).Unix()
	rtoken, err := token.SignedString(h.rsecret)
	if err != nil {
		return nil, err
	}

	user.RToken = rtoken
	if err := h.userRepo.Refresh(user); err != nil {
		return nil, err
	}

	return &tokenPair{atoken, rtoken}, nil
}

func getUserID(r *http.Request) string {
	token := r.Context().Value("user").(*jwt.Token)
	m := token.Claims.(jwt.MapClaims)
	return m["sub"].(string)
}
