package searchtool

import (
	"context"
)

type SearchTool interface {
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
	BulkUpdateDocument(string, string, []string, interface{}) error
	BulkUpdateDocumentWithContext(context.Context, string, string, []string, interface{}) error
	BulkDeleteDocument(string, string, []string) error
	BulkDeleteDocumentWithContext(context.Context, string, string, []string) error

	FindById(string, string, string, *interface{}) error
	FindByIdWithContext(context.Context, string, string, string, *interface{}) error

	Search(string, string, string, *interface{}) error
	SearchWithContext(context.Context, string, string, string, *interface{}) error
}
