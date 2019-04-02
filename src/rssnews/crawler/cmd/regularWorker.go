package main

import (
	"fmt"
	"rssnews/crawler"
	"rssnews/entity"
	"rssnews/services/scheduler"
	"runtime"
	"time"
)

type ()

func main() {
	var top int = 0
	var bottom int = 0
	schedul := new(scheduler.Service)
	tasks, err := schedul.ReadMany()
	taskLen := len(tasks)

	if taskLen == 0 {
		if err != nil {
			fmt.Println("Load tasks failed", err)
		}
		return
	}
	requestChan := make(chan entity.Scheduler)

	if taskLen > crawler.MAX_ROUTINES {
		top = crawler.MAX_ROUTINES
	} else {
		top = taskLen
	}

	for taskLen > bottom {
		taskSl := tasks[bottom:top]
		for _, task := range taskSl {
			go work(requestChan)
			requestChan <- task
		}
		bottom = top
		reserv := crawler.MAX_ROUTINES - runtime.NumGoroutine()
		if (taskLen - bottom) < reserv {
			top += taskLen - bottom
		} else {
			top += reserv
		}
	}
	time.Sleep(10 * time.Second)
}

func work(requestChan chan entity.Scheduler) {
	ent := <-requestChan
	fmt.Println(ent.Channel_id)
}
