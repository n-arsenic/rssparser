package main

import (
	//	"encoding/xml"
	"bytes"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"golang.org/x/text/encoding/htmlindex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	pbf "rssnews/protonotify"
	"rssnews/services/scheduler"
	"time"
)

const (
	port       = ":50051"
	MAX_MEMORY = 1 * 1024 //max requested content length
)

type (
	Server struct {
		Ch chan scheduler.Service
	}

	Worker struct{}
)

func (srv *Server) InsertNotify(ctx context.Context, req *pbf.Request) (*pbf.Response, error) {

	fmt.Printf("get notify %v\n", req)

	schedulerService := scheduler.Service{}
	schedulerService.Channel_id = int(req.Id)
	schedulerService.Rss_url = req.Url
	schedulerService.Start = pq.NullTime{Time: time.Now(), Valid: true}
	schedulerService.SetWorkStatus()

	err := schedulerService.Create()
	//send message may be buffered? channel
	if schedulerService.Exists && err == nil {
		srv.Ch <- schedulerService
	}
	return &pbf.Response{Received: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	srvLocal := new(Server)
	srvLocal.Ch = make(chan scheduler.Service)
	worker := new(Worker)

	go worker.Work(srvLocal.Ch)

	pbf.RegisterAsapWorkerServer(srv, srvLocal)
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		//[TODO] kill goroutines && close channel
		log.Fatalf("Failed to serve: %v", err)
	}
}

//[TODO] allocate this to another pkg
func (wr *Worker) Work(ch chan scheduler.Service) {
	fmt.Println("init work")
	//	fmt.Printf("v\n", charmap.All)
	for {
		//condition for close channel
		schedule := <-ch

		var content string
		resp, err := http.Get(schedule.Rss_url)
		if err != nil {
			fmt.Println(err)
			schedule.SetError("failed to load rss page")
			schedule.Update()
			continue

		}

		defer resp.Body.Close()

		limReader := io.LimitReader(resp.Body, MAX_MEMORY)
		buff := bytes.NewBuffer([]byte{})
		_, ierr := io.Copy(buff, limReader)
		content = buff.String()

		fmt.Println(ierr, content)

		reg, _ := regexp.Compile(`xml[\s]+version.+encoding="([a-zA-Z0-9\-:]+)"`)
		match := reg.FindStringSubmatch(content)

		if len(match) > 1 {
			charset := match[1]
			enc, err := htmlindex.Get(charset)

			fmt.Println(err)

			bufReader := bytes.NewReader(buff.Bytes())
			encReader := enc.NewDecoder().Reader(bufReader)
			buf := &bytes.Buffer{}
			_, ioerr := io.Copy(buf, encReader)
			content = buf.String()

			fmt.Println(ioerr, content)
		}

		//is it success to write content? yes - update plan_start
	}
}

/*
func (wr *Worker) HtmlParser(content []byte) {
	var data struct {
		Channel []Channel `xml:"channel"`
	}

	if err := xml.Unmarshal([]byte(content), &data); err != nil {
		log.Fatal(err)
	}
}*/
