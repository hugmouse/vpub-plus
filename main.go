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
	adminExists, err := data.HasAdmin()
	if !adminExists {
		if _, err := data.CreateUser("admin", model.UserCreationRequest{
			Name:     "admin",
			Password: "admin",
			IsAdmin:  true,
		}); err != nil {
			log.Fatal(err)
		}
	}
	log.Fatal(
		web.Serve(cfg, data),
	)
}
