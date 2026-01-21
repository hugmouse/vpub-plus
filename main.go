//go:generate go run generate.go
package main

import (
	"log"
	"math/rand"
	_ "net/http/pprof"
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
			Description: "This is a sample board.",
			Position:    0,
			IsLocked:    false,
			ForumID:     forumID,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = data.CreateTopic(userID, model.TopicRequest{
			Subject: "Change your admin account password",
			Content: `It might be a good idea to change your default Admin password!
			
## To change the password.
			
Navigate to [/admin/users](/admin/users). Find the **admin** user and change the password.`,
			BoardID:  boardID,
			IsLocked: true,
			IsSticky: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = data.CreateTopic(userID, model.TopicRequest{
			Subject:  "What to do?",
			Content:  `Navigate to the [/admin](/admin) route and see what you can change in there. Just make sure you do not delete yourself, okay? Otherwise you will be locked in a weird state.`,
			BoardID:  boardID,
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
