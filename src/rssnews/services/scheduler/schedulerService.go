package scheduler

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"reflect"
	"rssnews/entity"
	"rssnews/services"
	"strings"
	"time"
)

type (
	Service struct {
		entity.Scheduler
	}
)

func (sche *Service) Create() error {
	defer services.Postgre.Close()
	services.Postgre.Connect()

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
	/*
		if _err != nil {
			if err, ok := _err.(*pq.Error); ok {
				//Duplicate key
				if strings.Compare(string(err.Code), "23505") == 0 {
					return nil
				}

			}
			fmt.Println(_err)
			return _err
		}
	*/
	return _err
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

func (sche *Service) ReadMany(tl time.Duration) ([]Service, error) {
	defer services.Postgre.Close()
	services.Postgre.Connect()

	var (
		scheduler    *entity.Scheduler = &sche.Scheduler
		schedulers   []Service
		schedulerVal reflect.Value = reflect.ValueOf(scheduler).Elem()
		now                        = time.Now()
		timeLimit                  = now.Add(-tl) //-time.Duration(tl))
		_err         error
	)

	rows, _err := sq.Update("scheduler").
		SetMap(sq.Eq{
			"status": scheduler.GetWorkStatus(),
			"start":  now,
		}).
		Where(sq.Or{
			sq.And{
				sq.Eq{"status": scheduler.GetSuccessStatus()},
				sq.LtOrEq{"plan_start": now}},
			sq.And{
				sq.Eq{"status": scheduler.GetWorkStatus()},
				sq.LtOrEq{"start": timeLimit}}}).
		Suffix("RETURNING *").
		RunWith(services.Postgre.Db).
		PlaceholderFormat(sq.Dollar).
		Query()

		/*
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
		*/
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
				schedulers = append(schedulers, *sche)
			}
		}
	}

	return schedulers, _err
}
