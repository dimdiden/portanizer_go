package server

import (
	"encoding/json"
	"net/http"
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
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	// Read the request body
	if err := decoder.Decode(&user); err != nil {
		renderJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !user.IsValid() {
		renderJSON(w, portanizer.ErrEmpty.Error(), http.StatusBadRequest)
		return
	}
	if !h.repo.Exists(user) {
		renderJSON(w, "user not found", http.StatusBadRequest)
		return
	}
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)
	// Set token claims
	// claims["admin"] = true
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// Sign the token with the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		renderJSON(w, err.Error(), http.StatusInternalServerError)
	}
	// Finally, response with token
	renderJSON(w, tokenString, http.StatusOK)
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler:  OnError,
})

func OnError(w http.ResponseWriter, r *http.Request, err string) {
	renderJSON(w, err, http.StatusUnauthorized)
}

// jwtMiddleware.Options.ErrorHandler = OnError

// Examples of middlewares

// func authMiddleware(f func(http.ResponseWriter, *http.Request)) http.Handler {
// 	h := http.HandlerFunc(f)
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("\nMiddle!!!")
// 		h.ServeHTTP(w, r)
// 	})
// }

// func authMiddleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	fmt.Println("\nMiddle!!!")
// 	return f
// }
