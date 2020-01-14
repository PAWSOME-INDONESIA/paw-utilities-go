package elastic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (e *ElasticSearch) DeleteIndex(index string) error {
	return e.DeleteIndexesWithContext(context.Background(), []string{index})
}

func (e *ElasticSearch) DeleteIndexes(indexes []string) error {
	return e.DeleteIndexesWithContext(context.Background(), indexes)
}

func (e *ElasticSearch) DeleteIndexesWithContext(ctx context.Context, indexes []string) error {
	req := esapi.IndicesDeleteRequest{
		Index: indexes,
	}

	res, err := req.Do(ctx, e.Client)

	var r map[string]interface{}

	if err := e.do("Delete Index", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to delete elastic index with index %+v", indexes)
	}

	return nil
}
