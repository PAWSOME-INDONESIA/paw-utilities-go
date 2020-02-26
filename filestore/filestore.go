package filestore

import "context"

type (
	File struct {
		Mode        Mode
		Name, Path  string
		Content     []byte `swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
		ContentType string
	}

	Writer interface {
		Open(ctx context.Context, path string, mode Mode) (*File, error)
		Write(ctx context.Context, file *File, maxAge int) error
		Delete(ctx context.Context, file *File) error
		Close(ctx context.Context, file *File) error
		Exist(ctx context.Context, path string) (bool, error)
	}

	Mode uint
)

const (
	NEW    Mode = 1
	APPEND Mode = 2
)
