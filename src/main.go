package main

import (
	"dog-service/auth"
	"dog-service/config"
	"dog-service/logging"
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

	db.SaveLoginPassword(0, "admin", conf.AdminPassword, auth.RoleAdmin)

	router := mux.NewRouter()

	routes.SetupRoutes(
		router,
		conf.TemplatesPath,
		conf.StaticPath,
		log,
		&db)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router))
	fmt.Println("End of Program")
}
