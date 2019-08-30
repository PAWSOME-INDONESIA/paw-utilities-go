package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/pkg/errors"
)

const (
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

const (
	SearchTemplate = `{ "from" : %d, "size" : %d, "query" : %s , "sort" : [%s] }`
	BulkTemplate   = `{ "%s" : { "_index": "%s", "_type": "%s", "_id": "%s" } }`
)

func (e *ElasticSearch) constructBulkBody(action, index, _type string, ids []string, request interface{}, upsert bool) (string, error) {
	var response strings.Builder
	var err error
	datas := []string{}

	if action != CREATE && action != UPDATE && action != DELETE {
		e.Option.Log.Error(`Action must be between "create", "update" or "delete"`)
		return "", errors.Wrap(err, "Action must be between \"create\", \"update\" or \"delete\"")
	}

	if action != DELETE {
		val := reflect.ValueOf(request)
		if val.Kind() != reflect.Slice {
			e.Option.Log.Error("Request must be a list")
			return "", errors.Wrap(err, "Request must be a list")
		}

		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			encoded, err := json.Marshal(elem)
			if err != nil {
				e.Option.Log.Errorf("Error parsing: %+v", elem)
				return "", errors.Wrapf(err, "Error parsing: %+v", elem)
			} else {
				datas = append(datas, string(encoded))
			}
		}
	}

	for i, id := range ids {
		req := fmt.Sprintf(BulkTemplate, action, index, _type, id)
		response.WriteString(fmt.Sprintf("%s\n", req))

		if action == DELETE {
			continue
		}

		formatDoc := ""
		switch action {
		case CREATE:
			formatDoc = "%s\n"
		case UPDATE:
			formatDoc = "{ \"doc\" : %s, \"doc_as_upsert\" : %t } }\n"
		}

		var data = datas[i]
		if data[len(data)-1] == ',' {
			data = data[:len(data)-1]
		}

		response.WriteString(fmt.Sprintf(formatDoc, data, upsert))
	}
	return response.String(), nil
}

func (e *ElasticSearch) doBulk(ctx context.Context, action, index, _type string, ids []string, request interface{}, upsert bool) error {
	body, err := e.constructBulkBody(action, index, _type, ids, request, upsert)

	if err != nil {
		e.Option.Log.Errorf("Error constructBulkBody with Request : %+v", request)
		return errors.Wrapf(err, "Error constructBulkBody with Request : %+v", request)
	}

	req := esapi.BulkRequest{
		Index:        index,
		DocumentType: _type,
		Body:         strings.NewReader(body),
		Refresh:      "true",
	}

	res, err := req.Do(ctx, e.Client)

	defer func() {
		if err := res.Body.Close(); err != nil {
			e.Option.Log.Errorf("[Elastic Search] failed to close response body %s", err)
		}
	}()

	if err != nil {
		e.Option.Log.Errorf("[Elastic Search] Error getting response: %s", err)
		return errors.Wrap(err, "[Elastic Search] Error getting response")
	}

	if res.IsError() {
		e.Option.Log.Errorf("[Elastic Search] [%+v] Error %s document Ids=%+v", res.String(), action, ids)
		return errors.Wrapf(err, "[Elastic Search] [%+v] Error %s document Ids=%+v", res.String(), action, ids)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		e.Option.Log.Errorf("Error parsing the response body: %s", err)
		return errors.Wrapf(err, "[Elastic Search] Error parsing the response body: %s", err)
	}

	return nil
}

func (e *ElasticSearch) do(action string, res *esapi.Response, err error, data interface{}) error {
	if err != nil {
		e.Option.Log.Errorf("[Elastic Search] Error getting response: %+v", err)
		return errors.Wrap(err, "[Elastic Search] Error getting response")
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			e.Option.Log.Errorf("[Elastic Search] failed to close response body %s", err)
		}
	}()

	if res.IsError() {
		e.Option.Log.Errorf("[Elastic Search] [%+v] Error %s", res.String(), action)
		return errors.Wrapf(err, "[Elastic Search] [%+v] Error %s", res.String(), action)
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		e.Option.Log.Errorf("Error parsing the response body: %+v", err)
		return errors.Wrap(err, "[Elastic Search] Error parsing the response body")
	}

	return nil
}
