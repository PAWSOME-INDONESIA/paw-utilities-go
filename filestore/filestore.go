package filestore

type (
	File struct {
		Mode        Mode
		Name, Path  string
		Content     []byte `swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
		ContentType string
	}

	Writer interface {
		Open(path string, mode Mode) (*File, error)
		Write(file *File) error
		Delete(file *File) error
		Close(file *File) error
	}

	Mode uint
)

const (
	NEW    Mode = 1
	APPEND Mode = 2
)
