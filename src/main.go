package main

import (
	"dog-service/config"
	"dog-service/logging"
	"dog-service/server/middleware"
	"dog-service/routes"
	"dog-service/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	conf := config.Load()
	fmt.Println(conf)

	log := logging.Setup(conf.Env, conf.LogPath)
	log.Info("Go, go, Gophers!")

	db, err := storage.New(conf.DBConnection)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		return
	}
	log.Info("Connected to database")
	defer db.Close()

	router := mux.NewRouter()
	router.Use(middleware.LogRequest(log))

	routes.SetupRoutes(
		router,
		conf.TemplatesPath,
		conf.StaticPath,
		log,
		&db)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router))
	fmt.Println("End of Program")
}
