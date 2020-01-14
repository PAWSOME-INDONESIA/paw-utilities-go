package elastic

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

func (e *ElasticSearch) UpdateDocumentByQuery(index, _type, query string, data interface{}) error {
	return e.UpdateDocumentByQueryWithContext(context.Background(), index, _type, query, data)
}

func (e *ElasticSearch) UpdateDocumentByQueryWithContext(ctx context.Context, index, _type, query string, data interface{}) error {

	refresh := true

	req := esapi.UpdateByQueryRequest{
		Index:        []string{index},
		DocumentType: []string{_type},
		Body:         strings.NewReader(query),
		Refresh:      &refresh,
		Pretty:       true,
	}

	res, err := req.Do(ctx, e.Client)

	if err := e.do("Update Document By Query", res, err, &data); err != nil {
		return errors.Wrapf(err, "failed to update elastic document with query %s", query)
	}

	return nil
}
