package channel

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"rssparser/internal/pkg/entity"
	"rssparser/internal/pkg/services"
	//	"log"
	"github.com/lib/pq"
	pbf "rssparser/internal/pkg/protonotify"
	"time"
)

const (
	host = "localhost:50051"
)

type Service struct {
	entity.Channel
}

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

		if _err == nil {
			//send event to asap worker
			channelService.sendNotifyEvent(chanl)
		}
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

	//TODO it is nesessary write sheduler task here!!
	return response

}

func (channelService *Service) UpdatePubDate(date time.Time, chid int) {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	pubDate := pq.NullTime{
		Time:  date,
		Valid: true,
	}

	_, err := sq.Update("channels").
		SetMap(sq.Eq{
			"pub_date": pubDate,
		}).
		Where("id = ?", chid).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		fmt.Println("Update pub date of channel is failed: ", err)
	}
}

//update with embedded data
func (chanl *Service) Update() {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	data, err := sq.Update("channels").
		SetMap(sq.Eq{
			"title":       chanl.Title,
			"link":        chanl.Link,
			"description": chanl.Description,
			"pub_date":    chanl.Pub_date,
		}).
		Where("id = ?", chanl.Id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()
	fmt.Printf("%v\n", data)
	if err != nil {
		fmt.Println("Update of channel is failed: ", err)
	}
}

func (channelService *Service) ReadMany(rq *ReadManyRequest) *ReadManyResponse {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var response *ReadManyResponse = &ReadManyResponse{}

	rows, _err := sq.Select(
		"ch.id",
		"ch.rss_url",
		"ch.link",
		"ch.title",
		"ch.description",
		"ch.pub_date",
	).From("channels AS ch").
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
			&chanl.Link,
			&chanl.Title,
			&chanl.Description,
			&chanl.Pub_date,
		)
		if _resErr == nil {
			response.Channels = append(response.Channels, RssList{
				Id:          chanl.Id,
				Url:         chanl.Rss_url,
				Title:       chanl.Title.String,
				Link:        chanl.Link.String,
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
	rows, _err := sq.Select(
		"channel_id",
		"link",
		"title",
		"author",
		"category",
		"description",
		"pub_date",
	).From("channel_content").
		Where("channel_id = ?", rq.Channel_id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()

	//[TODO] if scheduler has error - print it
	for rows.Next() {
		var content *entity.ChannelContent = new(entity.ChannelContent)

		_err = rows.Scan(
			&content.Channel_id,
			&content.Link,
			&content.Title,
			&content.Author,
			&content.Category,
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

func New(id ...int) *Service {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	var chanl *Service = &Service{}

	if (len(id) == 0) || len(id) > 1 {
		return chanl
	}

	query := sq.Select(
		"id",
		"rss_url",
		"link",
		"title",
		"description",
		"pub_date",
		"created_at",
	).
		From("channels").
		Where("id = ?", id[0]).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	_err := query.Scan(
		&chanl.Id,
		&chanl.Rss_url,
		&chanl.Link,
		&chanl.Title,
		&chanl.Description,
		&chanl.Pub_date,
		&chanl.Created_at,
	)

	if _err != nil {
		fmt.Println("Error of creation channel with data: ", _err)
	}
	return chanl
}
