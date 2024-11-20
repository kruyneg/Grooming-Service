package main

import (
	"dog-service/config"
	"dog-service/logging"
	"dog-service/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	conf := config.Load()
	fmt.Println(conf)
	
	log := logging.Setup(conf.Env, conf.LogPath)
	log.Info("Go, go, Gophers!")

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	router.Use(middleware.LogRequest(log))

	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router)
}
