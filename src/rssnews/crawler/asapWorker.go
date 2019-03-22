package main

import (
	"log"
	"net"

	//	"database/sql"
	"fmt"
	//	sq "github.com/Masterminds/squirrel"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pbf "rssnews/protonotify"
	"time"
)

const (
	port = ":50051"
)

type Server struct {
	Ch chan *pbf.Request
}

func (srv *Server) InsertNotify(ctx context.Context, req *pbf.Request) (*pbf.Response, error) {
	fmt.Printf("get notify %v\n", req)

	//time.Sleep(time.Second * 5)
	/*
		//write to scheduler if not exists
		//- надо ли делать селект или обработать ошибку не уникальности?

		defer services.Postgre.Close()
		services.Postgre.Connect()
		var cid int
		query := sq.
			Insert("scheduler").
			Columns("channel_id", "rss_url", "start", "status").
			Values(req.Id, req.Url, time.Now(), "work").
			Suffix("RETURNING \"id\"").
			RunWith(services.Postgre.Db).
			PlaceholderFormat(sq.Dollar)

		_err = query.QueryRow().Scan(&cid)

		fmt.Printf("%#v", _err)
	*/

	//send message to buffered channel
	srv.Ch <- req
	return &pbf.Response{Received: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	srvLocal := new(Server)
	srvLocal.Ch = make(chan *pbf.Request)

	pbf.RegisterAsapWorkerServer(srv, srvLocal)
	reflection.Register(srv)
	go work(srvLocal.Ch)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func work(ch chan *pbf.Request) {
	fmt.Println("init work")
	for {
		//contition for close channel
		rq := <-ch
		time.Sleep(time.Second * 10)
		fmt.Println("I am!!!!!", rq)
	}
}

/*
	fmt.Println("async")
	ch := make(chan *pbf.Request)
	go work(ch)
	for {
		fmt.Println("1")

		if conn, err := lis.Accept(); err == nil {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
				os.Exit(1)
			}
			pdata := new(pbf.Request)
			err = proto.Unmarshal(buf.Bytes(), pdata)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
				os.Exit(1)
			}
			ch <- pdata
		}
	}


*/
