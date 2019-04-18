package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"rssparser/internal/pkg/services/user"
)

func userCreate(service *user.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			respwr.WriteHeader(http.StatusInternalServerError)
			respwr.Write([]byte("form error"))
			return
		}

		login := req.Form.Get("login")
		password := req.Form.Get("password")
		//passrepeat := req.Form.Get("passrepeat")
		//if password != passrepeat
		//return error message without header error

		rq := &user.CreateRequest{
			Login:    html.EscapeString(login),
			Password: html.EscapeString(password),
		}
		//return json
		createResp := service.Create(rq)
		json.NewEncoder(respwr).Encode(createResp)
	})
}

func userLogin(service *user.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		//hook for authentication problem with restarted server
		req.Header.Del("Cookie")

		session, sessErr := store.Get(req, "rssfeed-user")
		if sessErr != nil {
			http.Error(respwr, "Session error", http.StatusInternalServerError)
			fmt.Println(sessErr.Error())
			return
		}

		err := req.ParseForm()
		if err != nil {
			respwr.WriteHeader(http.StatusInternalServerError)
			respwr.Write([]byte("Form error"))
			fmt.Println(err.Error())
			return
		}
		login := req.Form.Get("login")
		password := req.Form.Get("password")
		rq := &user.ReadOneRequest{
			Login:    html.EscapeString(login),
			Password: html.EscapeString(password),
		}
		readOneResp, issetUser := service.ReadOne(rq)
		if issetUser {
			session.Values["user"] = readOneResp.User
			if err := session.Save(req, respwr); err != nil {
				http.Error(respwr, "Failed Autorization", http.StatusForbidden)
				readOneResp.Err_message += " User session does not create"
				fmt.Println(err.Error())
				return
			}
		}
		json.NewEncoder(respwr).Encode(readOneResp)

	})
}

//может быть не нужен сервис
func userLogout(service *user.Service) http.Handler {
	return http.HandlerFunc(func(respwr http.ResponseWriter, req *http.Request) {
		session, sessErr := store.Get(req, "rssfeed-user")
		if sessErr != nil {
			http.Error(respwr, "Session creation error", http.StatusInternalServerError)
			fmt.Println(sessErr.Error())
			return
		}

		session.Values["user"] = nil
		session.Options.MaxAge = -1

		sessErr = session.Save(req, respwr)

		if sessErr != nil {
			fmt.Println(sessErr.Error())
			return
		}

	})
}

func UserHandlerBind(service *user.Service) {
	http.Handle("/register/", userCreate(service))
	http.Handle("/login/", userLogin(service))
	http.Handle("/logout/", userLogout(service))
}
