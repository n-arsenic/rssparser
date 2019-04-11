package handlers

import (
	"encoding/gob"
	"github.com/gorilla/Sessions"
	"github.com/gorilla/securecookie"
	"rssnews/services/channel"
	"rssnews/services/user"
)

var store *sessions.CookieStore

func init() {
	sessionKey := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore(sessionKey)
	gob.Register(user.GetEmptyLoggUser())

}

func Compose() {
	userServ := user.NewUserService()
	UserHandlerBind(userServ)
	chanlServ := channel.New()
	ChanlHandlerBind(chanlServ)
}
