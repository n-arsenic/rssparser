package channel

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"rssnews/entity"
	"rssnews/services"
	//	"log"
	pbf "rssnews/protonotify"
)

const (
	host = "localhost:50051"
)

type Service struct{}

func (channelService *Service) sendNotifyEvent(ent entity.Channel) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		//log.Fatalf("Couldn't connect: %v", err)
		fmt.Printf("Couldn't connect: %v", err)
	}
	defer conn.Close()
	fmt.Println("init client")
	asapcli := pbf.NewAsapWorkerClient(conn)
	fmt.Println("SEND request to server")
	//reset connection!!!!! if not reply!
	//!!!!
	resp, err := asapcli.InsertNotify(context.Background(), &pbf.Request{Id: int32(ent.Id), Url: ent.Rss_url})
	if err != nil {
		//	log.Fatalf("could not send notify: %v", err)
		fmt.Printf("could not send notify: %v", err) //panic error
	}
	fmt.Printf("Get: %s\n", resp.Received)
}

func (channelService *Service) Create(rq *CreateRequest) *CreateResponse {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var response *CreateResponse = new(CreateResponse)
	var chanl = entity.Channel{}
	var _err error

	chanl.Rss_url = rq.Url

	//check unic rows
	selQuery := sq.Select("id").
		From("channels").
		Where("rss_url = ?", chanl.Rss_url).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	_err = selQuery.Scan(&chanl.Id)

	if _err == sql.ErrNoRows {
		query := sq.
			Insert("channels").
			Columns("rss_url").
			Values(chanl.Rss_url).
			Suffix("RETURNING \"id\", \"rss_url\"").
			RunWith(services.Postgre.Db).
			PlaceholderFormat(sq.Dollar)

		_err = query.QueryRow().Scan(&chanl.Id, &chanl.Rss_url)
	}

	if _err != nil && _err != sql.ErrNoRows {
		_err := errors.Wrapf(_err,
			"Insert new channel (url=%s) into channels table is failed",
			chanl.Rss_url)
		fmt.Println(_err)
		response.Err_message = "Error of create new channel"
		return response
	}

	_, _err = sq.
		Insert("user_channels").
		Columns("channel_id", "user_id").
		Values(
			chanl.Id,
			rq.User_id,
		).RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if _err != nil {
		_err := errors.Wrapf(_err,
			"Insert new relation user to rss channel (rss_id =%s) into user_channels table is failed",
			chanl.Id)
		fmt.Println(_err)
		response.Err_message = "Error of create relation"
		return response
	}

	response.Id = chanl.Id

	//check if channel already parsed
	_, _err = sq.Select("channel_id").
		From("scheduler").
		Where("channel_id = ?", chanl.Id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	if _err == nil {
		//send event to asap worker
		channelService.sendNotifyEvent(chanl)
	} else {
		fmt.Println("Scheduler select ", _err)
	}
	return response

}

func (channelService *Service) Update(ent *entity.Channel) {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	_, err := sq.Update("channels").
		SetMap(sq.Eq{
			"title":       ent.Title,
			"link":        ent.Link,
			"description": ent.Description,
			"pub_date":    ent.Pub_date,
		}).
		Where("id = ?", ent.Id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		fmt.Println("Update of channel is failed: ", err)
	}
}

func (channelService *Service) ReadMany(rq *ReadManyRequest) *ReadManyResponse {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var response *ReadManyResponse = &ReadManyResponse{}

	rows, _err := sq.Select("ch.id", "ch.rss_url", "ch.description", "ch.pub_date").
		From("channels AS ch").
		Join("user_channels AS uch ON uch.channel_id = ch.id").
		Where("uch.user_id = ?", rq.User_id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()

	if _err != nil {
		_err = errors.Wrapf(_err,
			"Select list of rss channels for user (id =%s) is failed",
			rq.User_id)
		fmt.Println(_err)
		response.Err_message = "Error get list of channels"
		return response
	}
	for rows.Next() {
		var chanl *entity.Channel = new(entity.Channel)

		_resErr := rows.Scan(
			&chanl.Id,
			&chanl.Rss_url,
			&chanl.Description,
			&chanl.Pub_date,
		)
		if _resErr == nil {
			response.Channels = append(response.Channels, RssList{
				Id:          chanl.Id,
				Url:         chanl.Rss_url,
				Description: chanl.Description.String,
				Pub_date:    chanl.Pub_date.Time,
			})
		} else {
			fmt.Println(_resErr)
		}
	}
	return response

}

func (channelService *Service) ReadOne(rq *ReadOneRequest) *ReadOneResponse {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var response *ReadOneResponse = &ReadOneResponse{}
	var chanl *entity.Channel = new(entity.Channel)
	var _err error

	query := sq.Select("ch.id", "ch.rss_url", "ch.description", "ch.pub_date").
		From("channels AS ch").
		Join("user_channels AS uch ON uch.channel_id = ch.id").
		Where(sq.Eq{"uch.user_id": rq.User_id, "uch.channel_id": rq.Channel_id}).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	_err = query.Scan(
		&chanl.Id,
		&chanl.Rss_url,
		&chanl.Description,
		&chanl.Pub_date,
	)
	if _err != nil {
		_err = errors.Wrapf(_err,
			"Select description rss channel for user (id =%s) is failed",
			rq.User_id)
		fmt.Println(_err)
		response.Err_message = "Error get channel info"
		return response
	}
	rows, _err := sq.Select("link", "title", "description", "pub_date").
		From("channel_content").
		Where("channel_id = ?", rq.Channel_id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()
		//[TODO] if scheduler has error - print it
	for rows.Next() {
		var content *entity.ChannelContent = new(entity.ChannelContent)

		_err = query.Scan(
			&content.Link,
			&content.Title,
			&content.Description,
			&content.Pub_date,
		)
		if _err != nil {
			_err = errors.Wrapf(_err,
				"Select content of rss channel (id=%s) for user (id =%s) is failed",
				rq.Channel_id, rq.User_id)
			fmt.Println(_err)
			response.Err_message = "Error: get content of rss channel is failed"
			return response
		}
		response.RssContent = append(response.RssContent, *content)
	}
	return response
}

func NewChanlService() *Service {
	return &Service{}
}
