package elastic

import (
	"context"
	"github.com/pkg/errors"

	"github.com/elastic/go-elasticsearch/esapi"
)

func (e *ElasticSearch) FindById(index, _type, id string, data *interface{}) error {
	return e.FindByIdWithContext(context.Background(), index, _type, id, data)
}

func (e *ElasticSearch) FindByIdWithContext(ctx context.Context, index, _type, id string, data *interface{}) error {
	req := esapi.GetRequest{
		Index:        index,
		DocumentType: _type,
		DocumentID:   id,
	}

	res, err := req.Do(ctx, e.Client)

	var r GetResponse

	if err := e.do("Find Document", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to find elastic document with id %s", id)
	}

	*data = r.Source
	return nil
}
