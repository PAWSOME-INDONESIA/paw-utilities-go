package elastic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (e *ElasticSearch) DeleteDocument(index, _type, id string) error {
	return e.DeleteDocumentWithContext(context.Background(), index, _type, id)
}

func (e *ElasticSearch) DeleteDocumentWithContext(ctx context.Context, index, _type, id string) error {
	req := esapi.DeleteRequest{
		DocumentID:   id,
		Index:        index,
		DocumentType: _type,
	}

	res, err := req.Do(ctx, e.Client)

	var r map[string]interface{}

	if err := e.do("Delete Document", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to remove elastic document with id %s", id)
	}

	return nil
}
