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

	if cfg.SessionKey == "your32byteslongsessionkeyhere" {
		log.Println("[Warning] You forgot to change your Session Key. Make sure to not expose this instance publicly.")
	}

	if cfg.CSRFKey == "your32byteslongcsrfkeyhere" {
		log.Println("[Warning] Remember to change the CSRF key, you're using the default one.")
	}

	db, err := storage.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	data := storage.New(db)
	adminExists, err := data.HasAdmin()
	if err != nil {
		log.Fatal(err)
	}

	if !adminExists {
		userID, err := data.CreateUser("admin", model.UserCreationRequest{
			Name:     "admin",
			Password: "admin",
			IsAdmin:  true,
		})
		if err != nil {
			log.Fatal(err)
		}

		forumID, err := data.CreateForum(model.ForumRequest{
			Name:     "Getting started",
			Position: 0,
			IsLocked: false,
		})
		if err != nil {
			log.Fatal(err)
		}

		boardID, err := data.CreateBoard(model.BoardRequest{
			Name:        "Your first board",
			Description: "This is an example board, you can check out what articles it haves.",
			Position:    0,
			IsLocked:    false,
			ForumId:     forumID,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = data.CreateTopic(userID, model.TopicRequest{
			Subject: "Change your admin account password NOW",
			Content: `Perhaps it would be a good idea to change your default admin password!
			
## Changing the password.
			
Navigate to [/admin/users](/admin/users) route. There, you should be able to find the **admin** user and set a new password to whatever you prefer.`,
			BoardId:  boardID,
			IsLocked: true,
			IsSticky: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = data.CreateTopic(userID, model.TopicRequest{
			Subject:  "What to do?",
			Content:  `Navigate to [/admin](/admin) route and see what you can poke in there. Just make sure to not delete yourself, okay? You will be locked in weird state otherwise.`,
			BoardId:  boardID,
			IsLocked: false,
			IsSticky: false,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(
		web.Serve(cfg, data),
	)
}
