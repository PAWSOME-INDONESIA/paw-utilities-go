package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/searchtool"
	"time"
)

const (
	DefaultShards             = 5
	DefaultReplica            = 1
	DefaultMaxIdleConnnection = 10
)

type (
	Option struct {
		Host                []string
		MaxIdleConnsPerHost int
		Log                 logs.Logger
		Shards              int
		Replica             int
		MaxPoolSize         int
		MaxBatchSize        int
	}

	GetResponse struct {
		Index           string      `json:"_index"`
		DocumentType    string      `json:"_type"`
		DocumentId      string      `json:"_id"`
		DocumentVersion int         `json:"_version"`
		Found           bool        `json:"found"`
		Source          interface{} `json:"_source"`
	}

	ElasticSearch struct {
		Option *Option
		Client *elasticsearch.Client
	}
)

//easyjson:json
type SearchResponse struct {
	Took    int         `json:"took"`
	TimeOut bool        `json:"time_out"`
	Shards  interface{} `json:"_shards"`
	Hits    SearchHits  `json:"hits"`
}

type SearchHits struct {
	Total int64       `json:"total"`
	Hits  interface{} `json:"hits"`
}

func getOption(option *Option) {
	if option.Log == nil {
		option.Log, _ = logs.DefaultLog()
	}

	if option.MaxIdleConnsPerHost == 0 {
		option.MaxIdleConnsPerHost = DefaultMaxIdleConnnection
	}

	if option.Shards == 0 {
		option.Shards = DefaultShards
	}

	if option.Replica == 0 {
		option.Replica = DefaultReplica
	}

	if option.MaxBatchSize == 0 {
		option.MaxBatchSize = 10
	}

	if option.MaxPoolSize == 0 {
		option.MaxPoolSize = 5
	}
}

func New(option *Option) (searchtool.SearchTool, error) {
	getOption(option)

	es := ElasticSearch{
		Option: option,
	}

	config := elasticsearch.Config{
		Addresses: option.Host,
	}

	client, err := elasticsearch.NewClient(config)

	if err != nil {
		return nil, errors.Wrap(err, "[Elastic Search] Error Create Client")
	}

	if _, err := client.Info(); err != nil {
		return nil, errors.Wrap(err, "[Elastic Search] Error Get Info")
	}

	es.Client = client

	return &es, nil
}

func (e *ElasticSearch) Ping() error {
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	req := esapi.PingRequest{
		ErrorTrace: true,
		Human:      true,
	}
	res, err := req.Do(ctx, e.Client)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.StatusCode == 200 {
		return nil
	}
	return errors.New("ping fail")
}
