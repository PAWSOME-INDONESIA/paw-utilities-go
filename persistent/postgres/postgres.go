package postgres

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/persistent"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func New(uri string) (persistent.ORM, error) {
	db, err := gorm.Open("postgres", uri)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open postgres connection!")
	}

	return &persistent.Impl{Database: db}, nil
}
