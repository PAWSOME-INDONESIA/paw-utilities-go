package mysql

import (
	"github.com/PAWSOME-INDONESIA/paw-utilities-go/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/PAWSOME-INDONESIA/paw-utilities-go/persistent"
)

func New(uri string, option *persistent.Option, logger logs.Logger) (persistent.ORM, error) {
	db, err := gorm.Open("mysql", uri)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open mysql connection!")
	}

	db.SetLogger(logger)
	db.LogMode(option.LogMode)

	db.DB().SetMaxIdleConns(option.MaxIdleConnection)
	db.DB().SetMaxOpenConns(option.MaxOpenConnection)
	db.DB().SetConnMaxLifetime(option.ConnMaxLifetime)

	return &persistent.Impl{Database: db, Logger: logger}, nil
}
