package user

import (
	"time"
)

type (
	ServiceInterface interface {
		Create(rq *CreateRequest) *CreateResponse
		ReadOne(rq *ReadOneRequest) *ReadOneResponse
	}

	CreateRequest struct {
		Login    string
		Password string
	}

	CreateResponse struct {
		Id          int
		Name        string
		Err_message string
	}

	ReadOneRequest struct {
		Login    string
		Password string
	}

	ReadOneResponse struct {
		User        LoggedUser
		Err_message string
	}

	LoggedUser struct {
		Id        int
		Name      string
		CreatedAt time.Time
	}
)
