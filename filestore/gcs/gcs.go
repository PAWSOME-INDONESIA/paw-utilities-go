package gcs

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/filestore"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type (
	writer struct {
		bucket string
	}
)

func (w *writer) Open(path string, mode filestore.Mode) (*filestore.File, error) {
	ctx := context.Background()
	var data []byte

	//open connection to gcs
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "OPEN - could not open connection to gcs")
	}

	//if mode is NEW, fill data with empty slice of bytes
	//else fetch data from gcs
	if mode == filestore.NEW {
		data = make([]byte, 0)
	} else {
		obj, err := client.Bucket(w.bucket).Object(path).NewReader(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "OPEN - could not initialize bucket reader")
		}

		data, err = ioutil.ReadAll(obj)
		if err != nil {
			return nil, errors.Wrap(err, "OPEN - could not read current data")
		}
	}

	// close connection to gcs
	if err := client.Close(); err != nil {
		return nil, errors.Wrap(err, "OPEN - could not close connection to gcs")
	}

	filename := md5.Sum([]byte(path))

	//fill content with data from gcs/empty slice of bytes
	return &filestore.File{
		Mode:    mode,
		Path:    path,
		Name:    fmt.Sprintf("%x", filename),
		Content: data,
	}, err
}

func (w *writer) Write(file *filestore.File) error {
	ctx := context.Background()

	//open connection to gcs
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "WRITE - could not open connection to gcs")
	}

	wc := client.Bucket(w.bucket).Object(file.Name).NewWriter(ctx)

	wc.ContentType = file.ContentType

	if _, err := wc.Write(file.Content); err != nil {
		return errors.Wrapf(err, "WRITE - unable to write data to bucket %q, file %q: %v", w.bucket, file.Name, err)
	}
	if err := wc.Close(); err != nil {
		return errors.Wrapf(err, "WRITE - unable to close bucket %q, file %q: %v", w.bucket, file.Name, err)
	}

	return nil
}

func (w *writer) Delete(file *filestore.File) error {
	// - delete logic
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "DELETE - could not open connection to gcs")
	}

	bk := client.Bucket(w.bucket)
	if err := bk.Object(file.Name).Delete(ctx); err != nil {
		return errors.Wrapf(err, "DELETE - unable to delete bucket %q, file %q", w.bucket, file.Name)
	}

	return nil
}

func (w *writer) Close(file *filestore.File) error {
	// - close logic
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "CLOSE - could not open connection to gcs")
	}
	if err := client.Close(); err != nil {
		return errors.Wrapf(err, "CLOSE - unable to close file %q", file.Name)
	}
	return nil
}

func NewWriter(bucket string) (filestore.Writer, error) {
	return &writer{bucket: bucket}, nil
}
