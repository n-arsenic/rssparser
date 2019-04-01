package crawler

/*
import (
	//	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	//	"github.com/lib/pq"
	"github.com/pkg/errors"
	"reflect"
	"rssnews/entity"
	"rssnews/services"
	"strings"
	"time"
)

type (
	Config struct {
		Create_limit  string
		Start_limit   string
		Success_limit string
	}
	Service struct {
		Config
		StTimeLim time.Time
		CrTimeLim time.Time
		SuTimeLim time.Time
	}
)

func (crawlService *Service) ReadMany() ([]entity.Channel, error) {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var (
		chanl    *entity.Channel = new(entity.Channel)
		channels []entity.Channel
		chanVal  reflect.Value = reflect.ValueOf(chanl).Elem()
		_err     error
	)
	fmt.Println(crawlService.CrTimeLim)
	//lock query
	rows, _err := sq.Select("*").
		From("channels").
		Where(sq.Eq{"status": chanl.GetNewStatus()}, sq.LtOrEq{"created_at": crawlService.CrTimeLim}).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()

	nn, val, er := sq.Select("*").
		From("channels").
		Where(sq.Or{
			sq.And{
				sq.Eq{"status": chanl.GetNewStatus()},
				sq.LtOrEq{"created_at": crawlService.CrTimeLim}},
			sq.And{
				sq.Eq{"status": chanl.GetSuccessStatus()},
				sq.LtOrEq{"parsed_at": crawlService.SuTimeLim}},
			sq.And{
				sq.Eq{"status": chanl.GetWaitStatus()},
				sq.LtOrEq{"start_parse": crawlService.StTimeLim}}}).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).ToSql()
	fmt.Println(nn, val, er)

	if _err == nil {
		columns, _ := rows.ColumnTypes() //rows.Columns()
		pointers := make([]interface{}, len(columns))

		for i, column := range columns {
			fieldVal := chanVal.FieldByName(strings.Title(column.Name()))
			if fieldVal.IsValid() {
				pointers[i] = fieldVal.Addr().Interface()
			} else {
				_err = errors.New("Structure field doesn't match table field")
				fmt.Println("field not valid")
			}
		}
		for rows.Next() {

			_err = rows.Scan(pointers...)
			if _err == nil {
				//	fmt.Println(chanl)
				channels = append(channels, *chanl)
			}

		}
	}
	return channels, _err
}

func (crawlService *Service) setTimeLimits() error {
	var past time.Time
	var _err error
	past, _err = getPastDate(crawlService.Create_limit)
	if _err == nil {
		crawlService.CrTimeLim = past
	}
	past, _err = getPastDate(crawlService.Success_limit)
	if _err == nil {
		crawlService.SuTimeLim = past
	}
	past, _err = getPastDate(crawlService.Start_limit)
	if _err == nil {
		crawlService.StTimeLim = past
	}

	return _err
}

func getPastDate(durString string) (time.Time, error) {
	var past time.Time
	dur, err := time.ParseDuration(durString)
	if err == nil {
		now := time.Now()
		past = now.Add(-dur)
	}

	return past, err
}

func New(config Config) *Service {
	service := &Service{Config: config}
	_err := service.setTimeLimits()
	if _err != nil {
		panic("Wrong time limits format")
	}
	return service
}*/
