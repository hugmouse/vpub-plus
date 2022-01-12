//go:generate go run generate.go
package main

import (
	"log"
	"math/rand"
	"time"
	"vpub/config"
	"vpub/model"
	"vpub/storage"
	"vpub/web"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cfg := config.New()
	db, err := storage.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	data := storage.New(db)
	if !data.HasAdmin() {
		if _, err := data.CreateUser(model.User{Name: "admin", Password: "admin", IsAdmin: true}, "admin"); err != nil {
			log.Fatal(err)
		}
	}
	log.Fatal(
		web.Serve(cfg, data),
	)
}
