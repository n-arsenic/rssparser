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
	pbf "rssnews/protonotify"
	"rssnews/services/scheduler"
	"strings"
	"time"
)

const (
	port       = ":50051"
	MAX_MEMORY = 5 * 1024
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
		resp, err := http.Get(schedule.Rss_url) //write to temporary file - bounds of memory
		if err != nil {
			fmt.Println(err)
			schedule.SetError("failed to load rss page")
			schedule.Update()
			continue

		}

		defer resp.Body.Close()

		//	html, rerr := ioutil.ReadAll(resp.Body)

		//get header charset if xml encoding does not exists
		html := io.LimitReader(resp.Body, MAX_MEMORY)
		enc, _ := htmlindex.Get("windows-1252")
		tr := enc.NewDecoder().Reader(html)
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, tr)
		fmt.Println("STRING ", err, buf.String())

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
