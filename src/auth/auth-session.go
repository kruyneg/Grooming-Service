package auth

import (
	"github.com/gorilla/sessions"
	"github.com/boj/redistore"
	"net/http"
)


var store *redistore.RediStore

func InitSessionStore() error {
	var err error
	store, err = redistore.NewRediStoreWithDB(
		10,
		"tcp",
		"localhost:6379",
		"",
		"",
		"0",
		[]byte("auth-key"),
	)
	if err != nil {
		return err
	}
	store.SetMaxAge(30 * 60) // 30 минут
	store.SetKeyPrefix("session_")
	return nil
}
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