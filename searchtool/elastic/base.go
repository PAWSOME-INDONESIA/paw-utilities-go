package elastic

import (
	"github.com/elastic/go-elasticsearch"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/searchtool"
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
	}

	GetResponse struct {
		Index           string      `json:"_index"`
		DocumentType    string      `json:"_type"`
		DocumentId      string      `json:"_id"`
		DocumentVersion int         `json:"_version"`
		Found           bool        `json:"found"`
		Source          interface{} `json:"_source"`
	}

	SearchResponse struct {
		Took    int         `json:"took"`
		TimeOut bool        `json:"time_out"`
		Shards  interface{} `json:"_shards"`
		Hits    SearchHits  `json:"hits"`
	}

	SearchHits struct {
		Total int64       `json:"total"`
		Hits  interface{} `json:"hits"`
	}

	ElasticSearch struct {
		Option *Option
		Client *elasticsearch.Client
	}
)

var Log logs.Logger

func getOption(option *Option) {
	if option.Log == nil {
		Log, _ = logs.DefaultLog()
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

	res, err := client.Info()

	if err != nil {
		return nil, errors.Wrap(err, "[Elastic Search] Error Get Info")
	}

	Log.Infof("[Elastic Search] %+v", res)

	es.Client = client

	return &es, nil
}
