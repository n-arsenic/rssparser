package main

import (
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"rssparser/internal/pkg/config"
	"rssparser/internal/pkg/crawler"
	pbf "rssparser/internal/pkg/protonotify"
	"rssparser/internal/pkg/services/scheduler"
	"time"
)

type (
	Server struct {
		Ch chan scheduler.Service
	}
)

var SRVlocal = new(Server)
var conf *config.Config = config.New()

func (srv *Server) InsertNotify(ctx context.Context, req *pbf.Request) (*pbf.Response, error) {

	fmt.Printf("get notify %v\n", req)

	schedulerService := scheduler.Service{}
	schedulerService.Channel_id = int(req.Id)
	schedulerService.Rss_url = req.Url
	schedulerService.Start = pq.NullTime{Time: time.Now(), Valid: true}
	schedulerService.SetWorkStatus()

	err := schedulerService.Create()
	//for new schedule only
	if err == nil {
		crawl := crawler.Crawler{Config: conf}
		go crawl.Work(SRVlocal.Ch)
		SRVlocal.Ch <- schedulerService
	}
	return &pbf.Response{Received: true}, nil
}

func main() {
	fmt.Println("asap worker is ready...")

	lis, err := net.Listen("tcp", conf.GRPC_HOST)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	SRVlocal.Ch = make(chan scheduler.Service)

	pbf.RegisterAsapWorkerServer(srv, SRVlocal)
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		//[TODO] kill goroutines && close channel
		log.Fatalf("Failed to serve: %v", err)
	}
}
