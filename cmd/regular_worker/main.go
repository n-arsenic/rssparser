package main

import (
	"fmt"
	"rssparser/internal/pkg/config"
	"rssparser/internal/pkg/crawler"
	"rssparser/internal/pkg/services/scheduler"
	"runtime"
	"sync"
)

var conf *config.Config = config.New()

func main() {
	fmt.Println("regular worker is ready...")
	var (
		wg     sync.WaitGroup
		top    int = 0
		bottom int = 0
	)
	schedul := new(scheduler.Service)
	tasks, err := schedul.ReadMany(conf.WORK_LIMIT)
	taskLen := len(tasks)

	if taskLen == 0 {
		if err != nil {
			fmt.Println("Load tasks failed", err)
		} else {
			fmt.Println("There is no tasks")
		}
		return
	}
	requestChan := make(chan scheduler.Service)
	crawl := crawler.Crawler{Config: conf}

	if taskLen > conf.MAX_ROUTINES {
		top = conf.MAX_ROUTINES
	} else {
		top = taskLen
	}

	for taskLen > bottom {
		taskSl := tasks[bottom:top]
		for _, task := range taskSl {
			wg.Add(1)
			//	go func() {
			crawl.Work(requestChan)
			wg.Done()
			//	}()
			requestChan <- task
		}
		bottom = top
		reserv := conf.MAX_ROUTINES - runtime.NumGoroutine()
		if (taskLen - bottom) < reserv {
			top += taskLen - bottom
		} else {
			top += reserv
		}
	}
	wg.Wait()
	fmt.Println("Main flow complete")
}
