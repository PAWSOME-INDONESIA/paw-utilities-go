package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/persistent"
)

func New(uri string) (persistent.ORM, error) {
	db, err := gorm.Open("mysql", uri)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open mysql connection!")
	}

	return &persistent.Impl{Database: db}, nil
}
