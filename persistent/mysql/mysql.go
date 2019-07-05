package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/persistent"
)

func New(uri string, option *persistent.Option) (persistent.ORM, error) {
	db, err := gorm.Open("mysql", uri)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open mysql connection!")
	}

	db.DB().SetMaxIdleConns(option.MaxIdleConnection)
	db.DB().SetMaxOpenConns(option.MaxOpenConnection)
	db.DB().SetConnMaxLifetime(option.ConnMaxLifetime)

	return &persistent.Impl{Database: db}, nil
}
