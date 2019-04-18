package channel

import (
	sq "github.com/Masterminds/squirrel"
	"rssparser/internal/pkg/entity"
	"rssparser/internal/pkg/services"
)

type (
	ContReadManyResponse struct{}
	ContReadManyRequest  struct{}
	ContentService       struct{}
)

func (chanContService *ContentService) Create(item entity.ChannelContent) error {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	_, err := sq.
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
			item.Channel_id,
			item.Link,
			item.Title,
			item.Author,
			item.Category,
			item.Description,
			item.Pub_date,
		).RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	//fmt.Printf("%v\n", res)

	return err
}

func (channelService *ContentService) ReadMany(rq *ContReadManyRequest) *ContReadManyResponse {
	return nil
}
