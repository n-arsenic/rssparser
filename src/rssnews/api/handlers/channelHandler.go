package handlers

import (
	"encoding/json"
	//	"fmt"
	"html"
	"net/http"
	"regexp"
	"rssnews/services/channel"
	"rssnews/services/user"
	"strconv"
)

func channelCreate(service *channel.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		luser, ok := user.GetUserSession(req, store)
		if !ok {
			http.Error(respwr, "Autorize error", http.StatusForbidden)
			return

		}

		err := req.ParseForm()
		if err != nil {
			http.Error(respwr, "Form error", http.StatusInternalServerError)
			return
		}
		//[TODO] form validation inside form custom pkg => NOT NULL FIELD
		rssurl := req.Form.Get("url")

		rq := &channel.CreateRequest{
			User_id: luser.Id,
			Url:     html.EscapeString(rssurl),
		}
		createResp := service.Create(rq)
		json.NewEncoder(respwr).Encode(createResp)
	})
}

func channelList(service *channel.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		luser, ok := user.GetUserSession(req, store)
		if !ok {
			http.Error(respwr, "Autorize error", http.StatusForbidden)
			return

		}
		rq := &channel.ReadManyRequest{User_id: luser.Id}
		listResp := service.ReadMany(rq)
		json.NewEncoder(respwr).Encode(listResp)
	})
}

func channelView(service *channel.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		luser, ok := user.GetUserSession(req, store)
		if !ok {
			http.Error(respwr, "Autorize error", http.StatusForbidden)
			return

		}
		//TODO create parsing and validate function
		var url = req.URL.Path
		var validPath = regexp.MustCompile("/([a-zA-Z0-9]*)$")
		var number int

		urlparts := validPath.FindStringSubmatch(url)
		if len(urlparts) > 0 {
			number, _ = strconv.Atoi(urlparts[1])
		}

		rq := &channel.ReadOneRequest{
			User_id:    luser.Id,
			Channel_id: number,
		}
		listResp := service.ReadOne(rq)
		json.NewEncoder(respwr).Encode(listResp)
	})
}

func channelDelete(service *channel.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		//проверять базу - есть ли подписанные на этот канал, если нет - удалять
	})
}

func ChanlHandlerBind(service *channel.Service) {
	http.Handle("/rss/new/", channelCreate(service))
	http.Handle("/rss/list/", channelList(service))
	http.Handle("/rss/view/", channelView(service))

}
