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

type authHandler struct {
	repo portanizer.UserRepo
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user portanizer.User
	var err error
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	// Read the request body
	if err = decoder.Decode(&user); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = user.IsValid(); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.repo.Exists(&user)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	atoken, rtoken, err := issueTokens(user.ID)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.RToken = rtoken

	if err := h.repo.Refresh(&user); err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokens := struct {
		AToken string `json:"atoken"`
		RToken string `json:"rtoken"`
	}{
		atoken,
		user.RToken,
	}
	// Finally, response with tokens
	renderJSON(w, tokens, http.StatusOK)
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: appkeyFunc,
	SigningMethod:       jwt.SigningMethodHS256,
	ErrorHandler:        onError,
})

func onError(w http.ResponseWriter, r *http.Request, err string) {
	renderJSON(w, err, http.StatusUnauthorized)
}

var appkeyFunc = func(token *jwt.Token) (interface{}, error) {
	return ASecret, nil
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var user portanizer.User
	var err error
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	// Read the request body
	if err = decoder.Decode(&user); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.RToken == "" {
		renderJSON(w, "please provide rtoken", http.StatusBadRequest)
		return
	}
	// Validate refresh token and parse claims from it
	token, err := jwt.Parse(user.RToken, func(token *jwt.Token) (interface{}, error) {
		return RSecret, nil
	})
	if err != nil {
		renderJSON(w, "invalid rtoken. suspicious activity: "+err.Error(), http.StatusBadRequest)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	sub, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = uint(sub)

	if err = h.repo.Valid(&user); err != nil {
		renderJSON(w, "invalid tokens. suspicious activity: "+err.Error(), http.StatusBadRequest)
		return
	}

	atoken, rtoken, err := issueTokens(user.ID)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.RToken = rtoken

	if err := h.repo.Refresh(&user); err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := struct {
		AToken string `json:"atoken"`
		RToken string `json:"rtoken"`
	}{
		atoken,
		user.RToken,
	}
	// Finally, response with token
	renderJSON(w, tokens, http.StatusOK)
	return
}

func issueTokens(uid uint) (atoken, rtoken string, err error) {
	// Create the token
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * atokenExp).Unix(),
		Subject:   strconv.Itoa(int(uid)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret
	atoken, err = token.SignedString(ASecret)
	if err != nil {
		return "", "", err
	}
	// Change expires to generate refresh token
	claims.ExpiresAt = time.Now().Add(time.Hour * rtokenExp).Unix()
	rtoken, err = token.SignedString(RSecret)
	if err != nil {
		return "", "", err
	}
	return atoken, rtoken, nil
}
