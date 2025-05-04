package main

import (
	"dog-service/auth"
	"dog-service/cache"
	"dog-service/config"
	"dog-service/logging"
	"dog-service/pubsub"
	"dog-service/routes"
	"dog-service/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	conf := config.Load()
	fmt.Println(conf)

	log := logging.Setup(conf.Env, conf.LogPath)
	log.Info("Go, go, Gophers!")

	if err := auth.InitSessionStore(); err != nil {
		panic(err)
	}
	if err := pubsub.InitPublisher("redis://:@localhost:6379/0", "logs"); err != nil {
		panic(err)
	}
	redis, err := cache.NewRedisCache("redis://:@localhost:6379/0", 30 * time.Minute)
	if err != nil {
		panic(err)
	}
	db, err := storage.New(conf.DBConnection, redis)
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
