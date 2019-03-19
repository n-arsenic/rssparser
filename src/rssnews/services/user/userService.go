package user

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gorilla/Sessions"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"rssnews/entity"
	"rssnews/services"
)

type Service struct{}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func GetUserSession(req *http.Request, store *sessions.CookieStore) (currentUser LoggedUser, ok bool) {
	session, sessErr := store.Get(req, "rssfeed-user")
	if sessErr != nil {
		ok = false
		return
	}
	sessValues := session.Values["user"]
	fmt.Printf("%#v\n\n\n", session.Values["user"])

	currentUser, ok = sessValues.(LoggedUser)
	return
}

//for Gob register
func GetEmptyLoggUser() LoggedUser {
	luser := LoggedUser{}
	return luser
}

func (userService *Service) Create(rq *CreateRequest) *CreateResponse {
	defer services.Postgre.Close()
	services.Postgre.Connect()
	var err string
	passhash, _ := hashPassword(rq.Password)

	user := new(entity.User)
	user.Name = rq.Login
	user.Password = passhash

	query := sq.
		Insert("users").Columns("name", "password").
		Values(user.Name, user.Password).Suffix("RETURNING \"id\"").
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar)

	_err := query.QueryRow().Scan(&user.Id)

	if _err != nil {
		_err := errors.Wrapf(_err,
			"Insert new user (username=%s) into user table is failed",
			user.Name)
		fmt.Println(_err)
		err = "Error of user registration, user doesn't created"
	}

	return &CreateResponse{
		Id:          user.Id,
		Name:        user.Name,
		Err_message: err,
	}
}

func (userService *Service) ReadOne(rq *ReadOneRequest) (*ReadOneResponse, bool) {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var isset bool
	resp := new(ReadOneResponse)

	query := sq.
		Select("id, name, password", "created_at").
		From("users").
		Where("name = ?", rq.Login).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		QueryRow()

	userEntity := entity.User{}

	err := query.Scan(
		&userEntity.Id,
		&userEntity.Name,
		&userEntity.Password,
		&userEntity.Created_at,
	)

	if err != nil {
		isset = false
		err = errors.Wrapf(err,
			"Select one user (userName=%s) from user table is failed", rq.Login)
		fmt.Println(err)
		resp.Err_message = "User does't found"
		return resp, isset
	}

	if isset = checkPasswHash(rq.Password, userEntity.Password); !isset {
		resp.Err_message = "Wrong login or password"
	} else {

		resp.User.Id = userEntity.Id
		resp.User.Name = userEntity.Name
		resp.User.CreatedAt = userEntity.Created_at
		isset = true

	}

	fmt.Println(resp)

	return resp, isset
}

func NewUserService() *Service {
	return &Service{}
}
