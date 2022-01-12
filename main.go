//go:generate go run generate.go
package main

import (
	"log"
	"vpub/config"
	"vpub/model"
	"vpub/storage"
	"vpub/web"
)

func main() {
	cfg := config.New()
	db, err := storage.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	data := storage.New(db)
	if !data.HasAdmin() {
		data.CreateUser(model.User{Name: "admin", Password: "admin", IsAdmin: true})
	}
	log.Fatal(
		web.Serve(cfg, data),
	)
}
