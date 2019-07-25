package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
)

func (e *ElasticSearch) Search(index, _type, query string, data interface{}) error {
	return e.SearchWithContext(context.Background(), index, _type, query, data)
}

func (e *ElasticSearch) SearchWithContext(ctx context.Context, index, _type, query string, data interface{}) error {
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

	var jsons []string
	variable := r.Hits.Hits.([]interface{})
	for _, value := range variable {
		obj := value.(map[string]interface{})["_source"]

		jsonString, err := json.Marshal(obj)
		if err != nil {
			return errors.Wrap(err, "failed to marshal document")
		}
		jsons = append(jsons, string(jsonString))
	}

	if err = json.Unmarshal([]byte("[" + strings.Join(jsons, ",") + "]"), &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal document")
	}
	return nil
}
