package main

import (
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	//	"regexp"
	"database/sql"
	"rssnews/crawler"
	"rssnews/entity"
	pbf "rssnews/protonotify"
	"rssnews/services/channel"
	"rssnews/services/scheduler"
	"time"
)

const (
	port = ":50051"
)

type (
	Server struct {
		Ch chan scheduler.Service
	}
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

	go Work(srvLocal.Ch)

	pbf.RegisterAsapWorkerServer(srv, srvLocal)
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		//[TODO] kill goroutines && close channel
		log.Fatalf("Failed to serve: %v", err)
	}
}

//[TODO] allocate this to another pkg
func Work(ch chan scheduler.Service) {
	fmt.Println("init work")
	for {
		//condition for close channel
		schedule := <-ch

		resp, err := http.Get(schedule.Rss_url)
		if err != nil {
			fmt.Println(err)
			schedule.SetError("failed to load rss page")
			schedule.Update()
			continue
		}

		defer resp.Body.Close()

		var result *crawler.Rss = crawler.XMLParser(resp.Body)

		var chnlEnt *entity.Channel = new(entity.Channel)
		var chnl *channel.Service = new(channel.Service)
		var chanlCont *channel.ContentService = new(channel.ContentService)

		chnlEnt.Id = schedule.Channel_id
		chnlEnt.Title = result.Channel.Title
		chnlEnt.Link = result.Channel.Link
		chnlEnt.Description = sql.NullString{
			String: result.Channel.Description,
			Valid:  true,
		}
		pubDate, _ := time.Parse(time.RFC1123Z, result.Channel.PubDate)
		chnlEnt.Pub_date = pq.NullTime{
			Time:  pubDate,
			Valid: true,
		}

		chnl.Update(chnlEnt)
		cerr := chanlCont.Create(result.Channel.Item, schedule.Channel_id)

		if cerr != nil {
			schedule.SetError(cerr.Error())
		} else {
			schedule.SetSuccessStatus()
			schedule.SetFinish()
		}
		schedule.SetPlanStart(crawler.PARSE_PERIOD)
		schedule.Update()

	}
}
