package example

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PAWSOME-INDONESIA/paw-utilities-go/logs"
	"github.com/PAWSOME-INDONESIA/paw-utilities-go/searchtool/elastic"
)

func esBulk() {
	type Tweet struct {
		ID      string `json:"id"`
		User    string `json:"user"`
		Message string `json:"message"`
	}

	var (
		Users    = []string{"robin", "ardo", "surya", "daniel", "galih", "jannes", "kevin", "jastian", "gabriel"}
		Messages = []string{"Hi!", "How r u?", "Im good", "Thank you", "okay"}
	)

	log, _ := logs.DefaultLog()

	option := elastic.Option{
		Log:  log,
		Host: []string{"http://localhost:9200"},
	}

	es, err := elastic.New(&option)
	if err != nil {
		panic(err)
	}

	esIndex := "tweet"
	esType := "tweet"
	var ids = make([]string, 0)
	var tweets = make([]Tweet, 0)

	rand.Seed(time.Now().UTC().UnixNano())
	for i := 1; i <= 60000; i++ {
		u := Users[rand.Intn(9)]
		ids = append(ids, fmt.Sprintf("%d", i))
		tweets = append(tweets, Tweet{
			ID:      fmt.Sprintf("%d", i),
			User:    fmt.Sprintf("%s", u),
			Message: fmt.Sprintf("message from %s: %s", u, Messages[rand.Intn(5)]),
		})
	}

	if err := es.BulkUpdateDocument(esIndex, esType, ids, tweets, true); err != nil {
		log.Error(err)
	}
}
