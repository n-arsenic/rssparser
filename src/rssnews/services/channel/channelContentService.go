package channel

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"rssnews/crawler"
	//	"rssnews/entity"
	"rssnews/services"
)

type (
	ContReadManyResponse struct{}
	ContReadManyRequest  struct{}
	ContentService       struct{}
)

func (chanContService *ContentService) Create(items []crawler.RssItem, chanId int) error {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	var err error

	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		_, _err := sq.
			Insert("channel_content").
			Columns(
				"channel_id",
				"link",
				"title",
				"author",
				"category",
				"description",
				"pub_date",
			).
			Values(
				chanId,
				item.Link,
				item.Title,
				item.Author,
				item.Category,
				item.Description,
				item.PubDate,
			).RunWith(services.Postgre.Db).
			PlaceholderFormat(sq.Dollar).
			Exec()
		if _err != nil {
			fmt.Println("Insert channel content is failed: ", _err)
			err = errors.New("Insert channel content was completed with errors")
		}
		//fmt.Printf("%v\n", res)
	}

	return err
}

func (channelService *ContentService) ReadMany(rq *ContReadManyRequest) *ContReadManyResponse {
	return nil
}
