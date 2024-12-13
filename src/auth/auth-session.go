package auth

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("auth-key"))

func CreateSession(w http.ResponseWriter, r *http.Request, userID int64, role string) error {
	session, _ := store.Get(r, "auth-session")
	session.Values["userId"] = userID
	session.Values["role"] = role

	session.Options.MaxAge = 30 * 60
	return session.Save(r, w)
}

func Check(r *http.Request) bool {
	session, _ := store.Get(r, "auth-session")
	_, ok := session.Values["userId"].(int64)
	return ok
}

func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "auth-session")
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "auth-session")
	session.Options.MaxAge = -1
	return session.Save(r, w)
}