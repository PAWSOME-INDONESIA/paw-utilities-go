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
	BulkTemplate   = `{ "%s" : { "_index": "%s", "_type": "%s", "_id": "%s" } }`
)

func constructBulkBody(action, index, _type string, ids []string, request interface{}, upsert bool) (string, error) {
	var response strings.Builder
	var err error
	datas := []string{}

	if action != CREATE && action != UPDATE && action != DELETE {
		Log.Error(`Action must be between "create", "update" or "delete"`)
		return "", errors.Wrap(err, "Action must be between \"create\", \"update\" or \"delete\"")
	}

	if action != DELETE {
		val := reflect.ValueOf(request)
		if val.Kind() != reflect.Slice {
			Log.Error("Request must be a list")
			return "", errors.Wrap(err, "Request must be a list")
		}

		encoded, err := json.Marshal(request)

		if err != nil {
			Log.Errorf("Error parsing: %+v", request)
			return "", errors.Wrapf(err, "Error parsing: %+v", request)
		}

		replacer := strings.NewReplacer("[", "", "]", "")
		datas = strings.SplitAfter(replacer.Replace(string(encoded)), "},")
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
	body, err := constructBulkBody(action, index, _type, ids, request, upsert)

	if err != nil {
		Log.Errorf("Error constructBulkBody with Request : %+v", request)
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
			Log.Errorf("[Elastic Search] failed to close response body %s", err)
		}
	}()

	if err != nil {
		Log.Errorf("[Elastic Search] Error getting response: %s", err)
		return errors.Wrap(err, "[Elastic Search] Error getting response")
	}

	if res.IsError() {
		Log.Errorf("[Elastic Search] [%s] Error %s document Ids=%+v", res.Status(), action, ids)
		return errors.Wrapf(err, "[Elastic Search] [%s] Error %s document Ids=%+v", res.Status(), action, ids)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		Log.Errorf("Error parsing the response body: %s", err)
		return errors.Wrapf(err, "[Elastic Search] Error parsing the response body: %s", err)
	}

	return nil
}

func (e *ElasticSearch) do(action string, res *esapi.Response, err error, data interface{}) error {
	if err != nil {
		Log.Errorf("[Elastic Search] Error getting response: %+v", err)
		return errors.Wrap(err, "[Elastic Search] Error getting response")
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			Log.Errorf("[Elastic Search] failed to close response body %s", err)
		}
	}()

	if res.IsError() {
		Log.Errorf("[Elastic Search] [%s] Error %s", res.Status(), action)
		return errors.Wrapf(err, "[Elastic Search] [%s] Error %s", res.Status(), action)
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		Log.Errorf("Error parsing the response body: %+v", err)
		return errors.Wrap(err, "[Elastic Search] Error parsing the response body")
	}

	return nil
}
