package main

import (
	"log"

	"microauth.io/core/internal/database"
	"microauth.io/core/internal/transport/http"
	"microauth.io/core/internal/user"
)

func main() {
	port := "8080"
	log.Println("starting server on port:", port)
	db := database.New()
	err := db.Connect("postgres", "postgres://postgres:root@localhost:5432/microauth?sslmode=disable")
	if err != nil {
		log.Println(err)
		log.Fatalln("error connecting to db")
	}
	userService := user.NewService(db)
	httpServer := http.New(userService)
	httpServer.RegisterHandlers()
	httpServer.Start(port)
}
