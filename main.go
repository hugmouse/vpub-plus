//go:generate go run generate.go
package main

import (
	"log"
	"pboard/config"
	"pboard/storage"
	"pboard/web"
)

func main() {
	cfg := config.New()
	db, err := storage.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	data := storage.New(db)
	//go gemini.Start(data)
	//go gopher.Start(data)
	log.Fatal(
		web.Serve(cfg, data),
	)
}
