package scheduler

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"reflect"
	"rssnews/crawler"
	"rssnews/entity"
	"rssnews/services"
	"strings"
	"time"
)

//future rss parsing

type (
	Service struct {
		entity.Scheduler
		Exists bool
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
			"status":     sche.Status,
			"message":    sche.Message,
			"start":      sche.Start,
			"finish":     sche.Finish,
			"plan_start": sche.Plan_start,
		}).
		Where("channel_id = ?", sche.Channel_id).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		fmt.Println("Scheduler update failed: ", err)
	}
	return err
}

func (sche *Service) ReadMany() ([]entity.Scheduler, error) {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var (
		scheduler    *entity.Scheduler = &sche.Scheduler
		schedulers   []entity.Scheduler
		schedulerVal reflect.Value = reflect.ValueOf(scheduler).Elem()
		now                        = time.Now()
		timeLimit                  = now.Add(-time.Duration(crawler.WORK_LIMIT))
		_err         error
	)

	rows, _err := sq.Select("*").
		From("scheduler").
		Where(sq.Or{
			sq.And{
				sq.Eq{"status": scheduler.GetSuccessStatus()},
				sq.LtOrEq{"plan_start": now}},
			sq.And{
				sq.Eq{"status": scheduler.GetWorkStatus()},
				sq.LtOrEq{"start": timeLimit}}}).
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()

	if _err == nil {
		columns, _ := rows.ColumnTypes() //rows.Columns()
		pointers := make([]interface{}, len(columns))

		for i, column := range columns {
			fieldVal := schedulerVal.FieldByName(strings.Title(column.Name()))
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
				schedulers = append(schedulers, *scheduler)
			}

		}
	}
	return schedulers, _err
}
