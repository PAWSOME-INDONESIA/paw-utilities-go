package elastic

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

const ConstantSetting = `{
	"settings":{
		"number_of_shards": %d,
		"number_of_replicas": %d
	},
	"mappings": {
		"%s" : %s
	}
}`

func constructIndexMapping(shards, replica int, _type, mapping string) string {
	if mapping == "" {
		mapping = "{}"
	}
	return fmt.Sprintf(ConstantSetting, shards, replica, _type, mapping)
}

func (e *ElasticSearch) CreateIndex(index, _type, mapping string) error {
	return e.CreateIndexWithContext(context.Background(), index, _type, mapping)
}

func (e *ElasticSearch) CreateIndexWithContext(ctx context.Context, index, _type, mapping string) error {
	err := e.IndexExistWithContext(ctx, index)
	if err != nil {
		e.Option.Log.Info(err)
		return errors.Wrap(err, "Error Check Index Exist")
	}

	e.Option.Log.Info("[ELASTIC SEARCH] Create Index")

	encoded := constructIndexMapping(e.Option.Shards, e.Option.Replica, _type, mapping)

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(encoded),
	}

	res, err := req.Do(ctx, e.Client)

	var r map[string]interface{}

	if err := e.do("Create Index", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to create elastic index with index %s", index)
	}

	return nil
}
