package auth

import (
	"net/http"
)

func GetId(r *http.Request) int64 {
	session, _ := getSession(r)
	res, _ := session.Values["userId"].(int64)
	return res
}

func GetRole(r * http.Request) string {
	session, _ := getSession(r)
	res, _ := session.Values["role"].(string)
	return res
}
