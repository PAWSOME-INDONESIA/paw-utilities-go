package elastic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

func (e *ElasticSearch) CreateDocument(index, _type, id string, request interface{}) error {
	return e.CreateDocumentWithContext(context.Background(), index, _type, id, request)
}

func (e *ElasticSearch) CreateDocumentWithContext(ctx context.Context, index, _type, id string, request interface{}) error {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		e.Option.Log.Errorf("[Elastic Search] Error on Parsing JSON: %+v\n", err)
		return errors.Wrapf(err, "[Elastic Search] Error on Parsing JSON: %+v\n", request)
	}

	if id == "" {
		e.Option.Log.Error("[Elastic Search] Id must be filled")
		return errors.Wrap(err, "[Elastic Search] Id must be filled")
	}

	req := esapi.IndexRequest{
		Index:        index,
		DocumentID:   id,
		DocumentType: _type,
		Body:         strings.NewReader(string(jsonRequest)),
		Refresh:      "true",
	}

	res, err := req.Do(ctx, e.Client)

	var r map[string]interface{}

	if err := e.do("Create Document", res, err, &r); err != nil {
		return errors.Wrap(err, "failed to create elastic document!")
	}

	return nil
}
