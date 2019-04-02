package crawler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"log"
	"time"
)

const (
	MAX_MEMORY   = 1024 * 1024 //max requested content length
	PARSE_PERIOD = 3 * time.Hour
	WORK_LIMIT   = 30 * time.Minute
)

type (
	RssItem struct {
		Title       string `xml:"title"`
		Link        string `xml:"link`
		Description string `xml:"description"`
		Author      string `xml:"author"`
		Category    string `xml:"category"`
		PubDate     string `xml:"pubDate"`
	}
	Rss struct {
		Channel struct {
			Title       string    `xml:"title"`
			Link        string    `xml:"link`
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

/*
type (
	Worker struct {
		DB_Object entity.Channel
		tlabel    string
	}
)

func (exec *Worker) SetTask(channel entity.Channel) *Worker {
	exec.DB_Object = channel //.(entity.Channel)
	return exec
}

//load rss content from rss channel
func (exec *Worker) load() *Worker {
	//load
	//	fmt.Println("Load rss content ID = ", exec.DB_Object.Id)
	time.Sleep(1 * time.Second)

	return exec
}

//parse rss page
func (exec *Worker) parsing() *Worker {
	//	fmt.Println("Parsing rss content ID = ", exec.DB_Object.Id)
	return exec

}

//save parsed rss content to database
func (exec *Worker) save() *Worker {
	//ВНИМАТЕЛЬНО проверять статус и лимиты времени - могут быть дубликаты
	fmt.Println("Save rss content ID = ", exec.DB_Object.Id)
	return exec
}

func (w *Worker) Execute(requestChan chan chan entity.Channel) {
	responseChan := make(chan entity.Channel)

	for {
		requestChan <- responseChan
		response := <-responseChan
		w.DB_Object = response
		w.load().parsing().save()
		//	close(responseChan)
	}
}
*/
