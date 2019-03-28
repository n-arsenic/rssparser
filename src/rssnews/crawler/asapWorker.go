package main

import (
	"bytes"
	"encoding/xml"
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
	//	"regexp"
	pbf "rssnews/protonotify"
	"rssnews/services/scheduler"
	"time"
)

const (
	port       = ":50051"
	MAX_MEMORY = 1024 * 1024 //max requested content length
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

		wr.XMLParser(resp.Body)

		//is it success to write content? yes - update plan_start
	}
}

func (wr *Worker) XMLParser(rbody io.Reader) {

	type (
		Rss struct {
			Channel struct {
				Title       string `xml:"title"`
				Link        string `xml:"link`
				Description string `xml:"description"`
				PubDate     string `xml:"pubDate"`
				Item        []struct {
					Title       string `xml:"title"`
					Link        string `xml:"link`
					Description string `xml:"description"`
					Author      string `xml:"author"`
					Category    string `xml:"category"`
					PubDate     string `xml:"pubDate"`
				} `xml:"item"`
			} `xml:"channel"`
		}
	)

	result := &Rss{}

	limReader := io.LimitReader(rbody, MAX_MEMORY)
	buff := bytes.NewBuffer([]byte{})
	_, ierr := io.Copy(buff, limReader)

	xdec := xml.NewDecoder(bytes.NewReader(buff.Bytes()))
	xdec.CharsetReader = identReader

	fmt.Println(ierr)

	if err := xdec.Decode(result); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result", result)
}

func identReader(encoding string, input io.Reader) (io.Reader, error) {
	enc, err := htmlindex.Get(encoding)
	encReader := enc.NewDecoder().Reader(input)
	fmt.Println("ident err", err)
	return encReader, err
}
