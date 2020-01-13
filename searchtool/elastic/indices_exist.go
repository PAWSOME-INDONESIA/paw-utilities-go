package elastic

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

func (e *ElasticSearch) IndexExist(index string) error {
	return e.IndexExistWithContext(context.Background(), index)
}

func (e *ElasticSearch) IndexExistWithContext(ctx context.Context, index string) error {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	res, err := req.Do(ctx, e.Client)

	if err != nil {
		e.Option.Log.Infof("[Elastic Search] Error getting response: %s", err)
		return errors.Wrap(err, "[Elastic Search] Error getting response")
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			e.Option.Log.Errorf("[ELATIC SEARCH] failed to close response body")
		}
	}()

	if res.IsError() {
		return nil
	}

	return errors.New(fmt.Sprintf("elastic index %s already exists", index))
}
