package scheduler

import (
	//	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"rssnews/entity"
	"rssnews/services"
	"strings"
	//	"time"
)

type (
	Config struct {
		Create_limit  string
		Start_limit   string
		Success_limit string
	}
	Service struct {
		entity.Scheduler
		Exists bool
		//	Config
		//	StTimeLim time.Time
		//	CrTimeLim time.Time
		//	SuTimeLim time.Time
	}
)

func (sche *Service) Create() error {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	sche.Exists = false

	query := sq.
		Insert("scheduler").
		Columns("channel_id", "rss_url", "start", "status").
		Values(sche.Channel_id, sche.Rss_url, sche.Start, sche.Status).
		Suffix("RETURNING \"channel_id\", \"rss_url\", \"start\", \"status\"").
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar)

	_err := query.QueryRow().Scan(
		&sche.Channel_id,
		&sche.Rss_url,
		&sche.Start,
		&sche.Status,
	)

	if _err != nil {
		if err, ok := _err.(*pq.Error); ok {
			if strings.Compare(string(err.Code), "23505") == 0 {
				return nil
			}

		}
		fmt.Println(_err)
		return _err
	}
	sche.Exists = true
	return nil
}

func (sche *Service) Update() error {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	_, err := sq.Update("scheduler").
		SetMap(sq.Eq{
			"status":  sche.Status,
			"message": sche.Message,
			"start":   sche.Start,
			"finish":  sche.Finish,
		}).
		Where("channel_id = ?", sche.Channel_id).
		Exec()

	if err != nil {
		fmt.Println(err)
	}
	return err
}

/*
func (crawlService *Service) ReadMany() ([]entity.Channel, error) {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var (
		chanl    *entity.Channel = new(entity.Channel)
		channels []entity.Channel
		chanVal  reflect.Value = reflect.ValueOf(chanl).Elem()
		_err     error
	)

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
}
*/
