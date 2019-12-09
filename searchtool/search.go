package searchtool

import (
	"context"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/util"
)

type SearchOption struct {
	Sort          []string
	ExcludedField []string
}

type SearchTool interface {
	util.Ping

	IndexExist(string) error
	IndexExistWithContext(context.Context, string) error
	CreateIndex(string, string, string) error
	CreateIndexWithContext(context.Context, string, string, string) error
	DeleteIndex(string) error
	DeleteIndexes([]string) error
	DeleteIndexesWithContext(context.Context, []string) error

	CreateDocument(string, string, string, interface{}) error
	CreateDocumentWithContext(context.Context, string, string, string, interface{}) error
	UpdateDocument(string, string, string, interface{}) error
	UpdateDocumentWithContext(context.Context, string, string, string, interface{}) error
	DeleteDocument(string, string, string) error
	DeleteDocumentWithContext(context.Context, string, string, string) error

	BulkCreateDocument(string, string, []string, interface{}) error
	BulkCreateDocumentWithContext(context.Context, string, string, []string, interface{}) error
	BulkUpdateDocument(string, string, []string, interface{}, interface{}, bool) error
	BulkUpdateDocumentWithContext(context.Context, string, string, []string, interface{}, interface{}, bool) error
	BulkDeleteDocument(string, string, []string) error
	BulkDeleteDocumentWithContext(context.Context, string, string, []string) error

	FindById(string, string, string, interface{}) error
	FindByIdWithContext(context.Context, string, string, string, interface{}) error

	Search(string, string, string, interface{}, ...SearchOption) error
	SearchWithContext(context.Context, string, string, string, interface{}, ...SearchOption) error
	SearchWithCustomQuery(context.Context, string, string, string, interface{}) error
}
