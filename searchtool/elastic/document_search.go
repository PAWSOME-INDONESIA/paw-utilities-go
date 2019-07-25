package elastic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

func (e *ElasticSearch) Search(index, _type, query string, data *interface{}) error {
	return e.SearchWithContext(context.Background(), index, _type, query, data)
}

func (e *ElasticSearch) SearchWithContext(ctx context.Context, index, _type, query string, data *interface{}) error {
	body := fmt.Sprintf(SearchTemplate, query)
	Log.Info(body)

	req := esapi.SearchRequest{
		Index:        []string{index},
		DocumentType: []string{_type},
		Body:         strings.NewReader(body),
		Pretty:       true,
	}

	res, err := req.Do(ctx, e.Client)

	var r SearchResponse

	if err := e.do("Search Document", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to search elastic document with query %s", query)
	}

	*data = r.Hits
	return nil
}
