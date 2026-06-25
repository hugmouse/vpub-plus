//go:generate go run generate.go
package main

import (
	"log"
	_ "net/http/pprof"
	"vpub/config"
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

	if err := web.Serve(cfg, data); err != nil {
		log.Fatal(err)
	}
}
