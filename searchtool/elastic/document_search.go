package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

type SearchData struct {
	jsons []string
	mu    sync.Mutex
}

func (sd *SearchData) append(data ...string) {
	sd.mu.Lock()
	defer sd.mu.Unlock()
	sd.jsons = append(sd.jsons, data...)
}

func (e *ElasticSearch) Search(index, _type, query string, data interface{}, sort ...string) error {
	return e.SearchWithContext(context.Background(), index, _type, query, data)
}

func (e *ElasticSearch) SearchWithContext(ctx context.Context, index, _type, query string, data interface{}, sort ...string) error {
	jsons := SearchData{jsons: make([]string, 0)}

	batchSize := int64(50)
	sortQuery := `{ "_id" : "asc" }`
	if len(sort) > 0 {
		sortQuery = strings.Join(sort, ",")
	}

	body := fmt.Sprintf(SearchTemplate, 0, batchSize, query, sortQuery)
	searchResponse, err := e.search(ctx, index, _type, body)
	if err != nil {
		return errors.Wrap(err, "failed to search document")
	}
	jsonResponse, err := e.getResponse(searchResponse)
	if err == nil {
		jsons.append(jsonResponse...)
	} else {
		return errors.WithStack(err)
	}

	totalData := searchResponse.Hits.Total
	totalPage := totalData / batchSize
	if totalData%batchSize != 0 {
		totalPage += 1
	}

	if totalPage > 1 {
		var wg sync.WaitGroup

		p, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
			page := i.(int64)
			start := (page - 1) * batchSize
			body := fmt.Sprintf(SearchTemplate, start, batchSize, query, sortQuery)
			searchResponse, err = e.search(ctx, index, _type, body)
			if err != nil {
				err = errors.WithStack(err)
			}
			jsonResponse, err = e.getResponse(searchResponse)
			if err == nil {
				jsons.append(jsonResponse...)
			} else {
				err = errors.WithStack(err)
			}
			wg.Done()
		})
		defer func() {
			_ = p.Release()
		}()

		wg.Add(int(totalPage - 1))
		for page := int64(2); page <= totalPage; page++ {
			_ = p.Invoke(page)
		}
		wg.Wait()
	}

	if err != nil {
		return errors.WithStack(err)
	}

	if err := json.Unmarshal([]byte("["+strings.Join(jsons.jsons, ",")+"]"), &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal document")
	}
	return nil
}

func (e *ElasticSearch) search(ctx context.Context, index, _type, query string) (*SearchResponse, error) {
	req := esapi.SearchRequest{
		Index:        []string{index},
		DocumentType: []string{_type},
		Body:         strings.NewReader(query),
		Pretty:       true,
	}
	res, err := req.Do(ctx, e.Client)
	var r SearchResponse
	if err := e.do("Search Document", res, err, &r); err != nil {
		return nil, errors.Wrapf(err, "failed to search elastic document with query %s", query)
	}
	return &r, nil
}

func (e *ElasticSearch) getResponse(r *SearchResponse) ([]string, error) {
	var jsons = make([]string, 0)
	variable, ok := r.Hits.Hits.([]interface{})
	if !ok {
		return nil, errors.New("Failed to get response")
	}
	for _, value := range variable {
		obj := value.(map[string]interface{})["_source"]

		jsonString, err := json.Marshal(obj)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal document")
		}
		jsons = append(jsons, string(jsonString))
	}
	return jsons, nil
}
