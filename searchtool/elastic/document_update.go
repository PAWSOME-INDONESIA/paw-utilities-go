package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

func (e *ElasticSearch) UpdateDocument(index, _type, id string, request interface{}) error {
	return e.UpdateDocumentWithContext(context.Background(), index, _type, id, request)
}

func (e *ElasticSearch) UpdateDocumentWithContext(ctx context.Context, index, _type, id string, request interface{}) error {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		e.Option.Log.Errorf("[Elastic Search] Error on Parsing JSON: %+v", err)
		return errors.Wrapf(err, "[Elastic Search] Error on Parsing JSON: %+v", request)
	}

	if id == "" {
		e.Option.Log.Error("[Elastic Search] Id must be filled")
		return errors.Wrap(err, "[Elastic Search] Id must be filled")
	}

	doc := fmt.Sprintf("{ \"doc\" : %s}", jsonRequest)

	req := esapi.UpdateRequest{
		Index:        index,
		DocumentID:   id,
		DocumentType: _type,
		Body:         strings.NewReader(string(doc)),
		Refresh:      "true",
	}

	res, err := req.Do(ctx, e.Client)

	var r map[string]interface{}

	if err := e.do("Update Document", res, err, &r); err != nil {
		return errors.Wrapf(err, "failed to update elastic document with id %s", id)
	}

	return nil
}
