//go:generate go run generate.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"vpub/config"
	"vpub/model"
	"vpub/storage"
	"vpub/web"
)

func processor(wg *sync.WaitGroup, c chan int, replies int, s *storage.Storage, board_id int64) {
	for i := range c {
		fmt.Println("Processing topic", i)
		topicId, err := s.CreateTopic(model.Post{
			User: model.User{
				Id: 1,
			},
			Subject:  "test",
			Content:  "test",
			IsSticky: false,
			IsLocked: false,
			BoardId:  board_id,
		})
		if err != nil {
			log.Panic(err, i)
			return
		}
		for j := 0; j < replies; j++ {
			_, err := s.CreatePost(model.Post{
				User:    model.User{Id: 1},
				Subject: "test",
				Content: "test",
				TopicId: topicId,
				BoardId: board_id,
			})
			if err != nil {
				log.Panic(err, j)
				return
			}
		}
		wg.Done()
	}
}

func seedTestData(topics, replies int, s *storage.Storage) {
	forum_id, err := s.CreateForum(model.Forum{
		Name:     "test",
		Position: 0,
	})
	if err != nil {
		log.Panic(err)
		return
	}
	board_id, err := s.CreateBoard(model.Board{
		Name:        "test",
		Description: "test",
		Position:    0,
		Forum:       model.Forum{Id: forum_id},
	})
	if err != nil {
		log.Panic(err)
		return
	}

	c := make(chan int, topics)

	var wg sync.WaitGroup
	for i := 0; i < topics; i++ {
		c <- i
		wg.Add(1)
	}
	close(c)
	go processor(&wg, c, replies, s, board_id)
	wg.Wait()
	fmt.Println("finished seeding")
}

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
	//seedTestData(1000, 100, data)
	//if _, err := data.ForumById(1); err != nil {
	//	fId, _ := data.CreateForum(model.Forum{
	//		Name: "Test Area",
	//	})
	//	data.CreateBoard(model.Board{
	//		Name:        "Testing",
	//		Description: "Main testing area",
	//		Forum:       model.Forum{Id: fId},
	//	})
	//	data.CreateBoard(model.Board{
	//		Name:        "Secondary Testing",
	//		Description: "Not the main, but still a testing area!",
	//		Forum:       model.Forum{Id: fId},
	//	})
	//}

	// Let's now attempt to insert n messages at the exact same time, and see if the count increased the way it should have.
	// At the start, topic count is 0 on the board.
	// If two transactions are executed concurrently, then the result would be 1.
	// That's because they'd do 0+1 each.
	//for i := 0; i < 100; i++ {
	//	go func() {
	//		fmt.Println(i)
	//		_, err := data.CreateTopic(model.Post{
	//			User:    model.User{Id: 1},
	//			Subject: "test",
	//			Content: "test",
	//			BoardId: bId1,
	//		})
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	//}
	log.Fatal(
		web.Serve(cfg, data),
	)
}
