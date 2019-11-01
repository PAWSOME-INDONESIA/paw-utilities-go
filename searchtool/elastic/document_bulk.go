package elastic

import (
	"context"
)

func (e *ElasticSearch) BulkCreateDocument(index string, _type string, ids []string, request interface{}) error {
	return e.BulkCreateDocumentWithContext(context.Background(), index, _type, ids, request)
}

func (e *ElasticSearch) BulkCreateDocumentWithContext(ctx context.Context, index, _type string, ids []string,
	request interface{}) error {
	return e.doBulk(ctx, CREATE, index, _type, ids, request, nil, false)
}

func (e *ElasticSearch) BulkUpdateDocument(index, _type string, ids []string,
	setOnInsert interface{}, setOnUpdate interface{}, upsert bool) error {
	return e.BulkUpdateDocumentWithContext(context.Background(), index, _type, ids, setOnInsert, setOnUpdate, upsert)
}

func (e *ElasticSearch) BulkUpdateDocumentWithContext(ctx context.Context, index, _type string, ids []string,
	setOnInsert interface{}, setOnUpdate interface{}, upsert bool) error {
	return e.doBulk(ctx, UPDATE, index, _type, ids, setOnInsert, setOnUpdate, upsert)
}

func (e *ElasticSearch) BulkDeleteDocument(index, _type string, ids []string) error {
	return e.BulkDeleteDocumentWithContext(context.Background(), index, _type, ids)
}

func (e *ElasticSearch) BulkDeleteDocumentWithContext(ctx context.Context, index, _type string, ids []string) error {
	return e.doBulk(ctx, DELETE, index, _type, ids, nil, nil, false)
}
