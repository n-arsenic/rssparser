package channel

import (
	"rssparser/internal/pkg/entity"
	"time"
)

type (
	ServiceInterface interface {
		Create(rq *CreateRequest) *CreateResponse
		ReadMany(rq *ReadManyRequest) *ReadManyResponse
		ReadOne(rq *ReadOneRequest) *ReadOneResponse
		Delete(rq *DeleteRequest) *DeleteResponse
	}

	CreateRequest struct {
		User_id int
		Url     string
	}

	CreateResponse struct {
		Id          int
		Err_message string
	}

	ReadOneRequest struct {
		User_id    int
		Channel_id int
	}

	ReadOneResponse struct {
		RssContent  []entity.ChannelContent
		Err_message string
	}

	ReadManyRequest struct {
		User_id int
	}

	ReadManyResponse struct {
		Channels    []RssList
		Err_message string
	}

	DeleteRequest struct{}

	DeleteResponse struct{}

	RssList struct {
		Id          int
		Url         string
		Link        string
		Title       string
		Description string
		Pub_date    time.Time
	}
)
