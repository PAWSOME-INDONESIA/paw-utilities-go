package example

import (
	"fmt"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/searchtool/elastic"
	"math/rand"
	"sync"
	"time"
)

func main() {
	type Tweet struct {
		ID      string `json:"id"`
		User    string `json:"user"`
		Message string `json:"message"`
	}

	var (
		Users    = []string{"robin", "ardo", "surya", "daniel", "galih", "jannes", "kevin", "jastian", "gabriel"}
	)

	log, _ := logs.DefaultLog()
	option := elastic.Option{
		Log:  log,
		Host: []string{"http://localhost:9200"},
		MaxBatchSize: 100,
		MaxPoolSize: 10,
	}

	es, err := elastic.New(&option)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	esIndex := "tweet"
	esType := "tweet"

	totalQuery := 8
	var wg sync.WaitGroup
	wg.Add(totalQuery)
	for i := 0; i < totalQuery ; i++ {
		u := Users[rand.Intn(9)]
		query := `{
			"term" : {
				"user": "%s"
			}
		}`
		q := fmt.Sprintf(query, u)
		sort := []string{`{"user":"asc"}`,`{"id":"asc"}`}

		go func () {
			start := time.Now()
			var tweets []Tweet
			if err := es.Search(esIndex, esType, q, &tweets, sort...); err != nil {
				log.Error(err)
				return
			}
			end := time.Now()
			log.Infof("%d - %+v", len(tweets), end.Sub(start))
			wg.Done()
		}()
	}
	wg.Wait()
}
