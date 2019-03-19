package crawler

import (
	"fmt"
	//"log"
	//"net/http"
	"database/sql"
	"github.com/lib/pq"
	"rssnews/entity"
	service "rssnews/services/crawler"

	//	"runtime"
	"time"
)

//get channel list from DB for processing
func GetTasks() (tasks []entity.Channel) {
	//select tasks with condition from database

	for i := 0; i < 100; i++ {
		task := entity.Channel{
			Id:          i,
			Rss_url:     "http://localhost:8282/rss/crawler/test" + string(i),
			Description: sql.NullString{String: "bla bla", Valid: true},
			Pub_date:    pq.NullTime{Time: time.Now(), Valid: true},
			Parsed_at:   pq.NullTime{Time: time.Now(), Valid: true},
			Status:      "success",
			Created_at:  time.Now(),
		}
		tasks = append(tasks, task)
	}
	return
}

func Start() {

	serv := service.New(service.Config{
		Create_limit:  CREATE_LIM,
		Start_limit:   START_LIM,
		Success_limit: SUCCESS_LIM,
	})
	tasks, err := serv.ReadMany()
	fmt.Println(err)

	for _, task := range tasks {
		fmt.Println("task ID = ", task)
	}

	/*
		requestChan := make(chan chan entity.Channel)

		for i := 0; i < MAX_ROUTIN; i++ {
			worker := new(Worker)
			go worker.Execute(requestChan)

		}
		for {
			tasks := GetTasks()
			if len(tasks) > 0 {
				for _, task := range tasks {
					responseChan := <-requestChan
					responseChan <- task
					fmt.Println("NEXT")

				}
				fmt.Println("Wait for tasks")
				//close chan
				//	fmt.Println("/////num before ////// ", runtime.NumGoroutine())
				fmt.Println("num : ", runtime.NumGoroutine())

			}
			time.Sleep(10 * time.Second)
		}
	*/
}
