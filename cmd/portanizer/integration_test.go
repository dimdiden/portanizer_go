package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dimdiden/portanizer_go"
	"github.com/dimdiden/portanizer_go/server"

	"github.com/dimdiden/portanizer_go/gorm"
)

const (
	atokenExp = 1
	rtokenExp = 4
)

var c = &conf{
	DBdriver: "sqlite3",
	DBname:   "test",

	ASecret: []byte("ASECRET"),
	RSecret: []byte("RSECRET"),

	logout: ioutil.Discard,
}
var srv *server.Server

func TestMain(m *testing.M) {
	// TODO: change to contain test environment variables
	db, err := c.openGormDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gorm.RunMigrations(db)

	if err := populate(db, "testdata.sql"); err != nil {
		log.Fatal("cannot populate db: ", err)
	}

	srv = c.openGormServer(db)

	code := m.Run()
	сlean(c.DBname)
	os.Exit(code)
}

func TestUserHandlers(t *testing.T) {
	t.Run("SignUp", testUserSignUp)
	t.Run("SignIn", testUserSignIn)
	t.Run("Refresh", testUserRefresh)
	t.Run("RefreshWithWrongToken", testUserRefreshWithWrongToken)
}

func testUserSignUp(t *testing.T) {
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"email":"test@gmail.com","password": "test"}`)
	r := httptest.NewRequest(http.MethodPost, "/users", body)
	srv.ServeHTTP(w, r)

	equals(t, `{"Message":"user has been created"}`, http.StatusOK, w)
}

func testUserSignIn(t *testing.T) {
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"email":"test@gmail.com","password": "test"}`)
	r := httptest.NewRequest(http.MethodPost, "/users/signin", body)
	srv.ServeHTTP(w, r)

	tp := struct {
		AToken string `json:"atoken"`
		RToken string `json:"rtoken"`
	}{}

	decoder := json.NewDecoder(w.Result().Body)
	err := decoder.Decode(&tp)
	ok(t, err)

	if w.Code != http.StatusOK {
		t.Error("unexpected response code: ", w.Code)
	}

	if tp.AToken == "" || tp.RToken == "" {
		t.Error("Token is empty")
		return
	}

	err = validateToken(tp.AToken, c.ASecret)
	err = validateToken(tp.RToken, c.RSecret)
	ok(t, err)
}

func testUserRefresh(t *testing.T) {
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"rtoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDI0MTQ0MzQsInN1YiI6IjEifQ.u804NtZ7Eoz5_uAZUkU0nRI1_YpARI9XwAys5Fk6YL8"}`)
	r := httptest.NewRequest(http.MethodPost, "/users/refresh", body)

	user := &portanizer.User{ID: 1, Email: "dimdiden@gmail.com", Password: "123"}

	atoken, err := issueToken(user)
	ok(t, err)
	r.Header.Add("Authorization", "Bearer "+atoken)
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("unexpected response code: ", w.Code)
	}

	tp := struct {
		AToken string `json:"atoken"`
		RToken string `json:"rtoken"`
	}{}

	decoder := json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&tp)
	ok(t, err)

	if tp.AToken == "" || tp.RToken == "" {
		t.Error("Token is empty")
		return
	}

	err = validateToken(tp.AToken, c.ASecret)
	err = validateToken(tp.RToken, c.RSecret)
	ok(t, err)
}

func testUserRefreshWithWrongToken(t *testing.T) {
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"rtoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDI0MTQ0MzQsInN1YiI6IjEifQ.u804NtZ7Eoz5_uAZUkU0nRI1_YpARI9XwAys5Fk6YL8"}`)
	r := httptest.NewRequest(http.MethodPost, "/users/refresh", body)

	user := &portanizer.User{ID: 1, Email: "svediden@gmail.com", Password: "123"}
	atoken, err := issueToken(user)
	ok(t, err)
	r.Header.Add("Authorization", "Bearer "+atoken)
	srv.ServeHTTP(w, r)

	equals(t, `{"Message":"rtoken is invalid"}`, http.StatusUnauthorized, w)
}

func TestGetTagHandlers(t *testing.T) {
	t.Run("GetTag", testGetTag)
	t.Run("GetTags", testGetTags)
}

func testGetTag(t *testing.T) {
	user := &portanizer.User{ID: 1, Email: "dimdiden@gmail.com", Password: "123"}

	atoken, err := issueToken(user)
	ok(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
	r.Header.Add("Authorization", "Bearer "+atoken)
	srv.ServeHTTP(w, r)

	equals(t, `{"ID":1,"Name":"Tag1"}`, http.StatusOK, w)
}

func testGetTags(t *testing.T) {
	user := &portanizer.User{ID: 1, Email: "dimdiden@gmail.com", Password: "123"}

	atoken, err := issueToken(user)
	ok(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tags", nil)
	r.Header.Add("Authorization", "Bearer "+atoken)
	srv.ServeHTTP(w, r)

	equals(t, `[{"ID":1,"Name":"Tag1"},{"ID":2,"Name":"Tag2"}]`, http.StatusOK, w)
}
func TestPostTagHandlers(t *testing.T) {
	t.Run("PostTag", testPostTag)
}

func testPostTag(t *testing.T) {
	user := &portanizer.User{ID: 1, Email: "dimdiden@gmail.com", Password: "123"}

	atoken, err := issueToken(user)
	ok(t, err)

	w := httptest.NewRecorder()
	body := strings.NewReader(`{"name":"Tag3"}`)
	r := httptest.NewRequest(http.MethodPost, "/tags", body)
	r.Header.Add("Authorization", "Bearer "+atoken)
	srv.ServeHTTP(w, r)

	equals(t, `{"ID":3,"Name":"Tag3"}`, http.StatusOK, w)
}

func issueToken(user *portanizer.User) (string, error) {
	// Create the token
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * atokenExp).Unix(),
		Subject:   strconv.Itoa(int(user.ID)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret
	atoken, err := token.SignedString(c.ASecret)
	if err != nil {
		return "", err
	}
	return atoken, nil
}

func validateToken(tokenString string, secret []byte) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func populate(db *gorm.DB, testdata string) error {
	file, err := os.Open(testdata)
	if err != nil {
		return err
	}
	defer file.Close()

	query, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if err := db.Exec(string(query)).Error; err != nil {
		return err
	}
	return nil
}

func сlean(db string) error {
	err := os.Remove("test.db")
	if err != nil {
		return err
	}
	return nil
}

func ok(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func equals(t *testing.T, exp string, status int, w *httptest.ResponseRecorder) {
	if w.Code != status {
		t.Error("unexpected response code: ", w.Code)
	}
	result, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Error("can not read response body")
	}
	if exp != string(result) {
		t.Error("unexpected result: ", string(result))
	}
}
