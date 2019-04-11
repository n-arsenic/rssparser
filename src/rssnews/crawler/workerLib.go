package crawler

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"log"
	"net/http"
	"rssnews/entity"
	"rssnews/services/channel"
	"rssnews/services/scheduler"
	"time"
)

const (
	MAX_MEMORY   = 1024 * 1024 //max requested content length
	PARSE_PERIOD = 3 * time.Hour
	WORK_LIMIT   = 30 * time.Minute
	MAX_ROUTINES = 10
)

type (
	RssItem struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Author      string `xml:"author"`
		Category    string `xml:"category"`
		PubDate     string `xml:"pubDate"`
	}
	Rss struct {
		Channel struct {
			Title       string    `xml:"title"`
			Link        string    `xml:"link"`
			Description string    `xml:"description"`
			PubDate     string    `xml:"pubDate"`
			Item        []RssItem `xml:"item"`
		} `xml:"channel"`
	}
)

func XMLParser(rbody io.Reader) *Rss {
	result := &Rss{}
	limReader := io.LimitReader(rbody, MAX_MEMORY)
	buff := bytes.NewBuffer([]byte{})
	_, ierr := io.Copy(buff, limReader)
	xdec := xml.NewDecoder(bytes.NewReader(buff.Bytes()))
	xdec.CharsetReader = identReader

	fmt.Println(ierr)

	if err := xdec.Decode(result); err != nil {
		log.Fatal(err)
	}

	return result
}

func identReader(encoding string, input io.Reader) (io.Reader, error) {
	enc, err := htmlindex.Get(encoding)
	encReader := enc.NewDecoder().Reader(input)
	if err != nil {
		fmt.Println("ident err", err)
	}
	return encReader, err
}

//[TODO] allocate this to another pkg
func Work(ch chan scheduler.Service) {
	fmt.Println("init work")

	schedule := <-ch

	fmt.Println(schedule.Channel_id)

	resp, err := http.Get(schedule.Rss_url)
	if err != nil {
		fmt.Println(err)
		schedule.SetError("failed to load rss page")
		schedule.Update()
		return
	}

	defer resp.Body.Close()

	var result *Rss = XMLParser(resp.Body)

	var (
		chnl          *channel.Service        = channel.New(schedule.Channel_id)
		chanlCont     *channel.ContentService = new(channel.ContentService)
		items                                 = result.Channel.Item
		hasErrors     bool                    = false
		oldPubDate    time.Time               = chnl.Pub_date.Time
		maxPubDate    time.Time               = oldPubDate
		oldPubDateSec int64                   = oldPubDate.Unix()
	)
	//if channel is new => channel pub date is NULL
	if chnl.Pub_date.Valid == false {
		chnl.Title = sql.NullString{
			String: result.Channel.Title,
			Valid:  result.Channel.Title != "",
		}
		chnl.Link = sql.NullString{
			String: result.Channel.Link,
			Valid:  result.Channel.Link != "",
		}
		chnl.Description = sql.NullString{
			String: result.Channel.Description,
			Valid:  result.Channel.Description != "",
		}
	}

	for _, item := range items {
		var (
			itemDate, _ = time.Parse(time.RFC1123, item.PubDate)
			itemDataSec = itemDate.Unix()
		)
		fmt.Println(schedule.Channel_id, item.PubDate, itemDate)
		if itemDataSec > oldPubDateSec {
			_err := chanlCont.Create(entity.ChannelContent{
				Channel_id:  schedule.Channel_id,
				Link:        item.Link,
				Title:       item.Title,
				Author:      item.Author,
				Category:    item.Category,
				Description: item.Description,
				Pub_date:    itemDate,
			})
			if _err != nil {
				fmt.Println("Insert channel content is failed: ", _err)
				hasErrors = true
			} else if itemDataSec > maxPubDate.Unix() {
				maxPubDate = itemDate
			}
		}
	}

	chnl.Pub_date = pq.NullTime{
		Time:  maxPubDate,
		Valid: maxPubDate.Unix() > 0,
	}

	chnl.Update()

	if hasErrors {
		schedule.SetError("Insert channel content was completed with errors")
	} else {
		schedule.SetSuccessStatus()
		schedule.SetFinish()
	}
	schedule.SetPlanStart(PARSE_PERIOD)
	schedule.Update()
	fmt.Println("Routine flow complete")
}
