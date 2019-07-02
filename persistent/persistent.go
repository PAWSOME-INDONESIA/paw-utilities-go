package persistent

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type (
	ORM interface {
		Close() error

		FindOne(interface{}, ...interface{}) error
		Find(interface{}, interface{}, ...interface{}) error

		Create(interface{}) error
		Update(interface{}) error
		Delete(interface{}) error
		SoftDelete(interface{}) error

		// Exec is used to execute sql Create, Update or Delete
		Exec(string, ...interface{}) error

		// RawSql is used to execute Select
		RawSqlWithObject(string, interface{}, ...interface{}) error
		RawSql(string, ...interface{}) (*sql.Rows, error)

		Begin() ORM
		Commit() error
		Rollback() error
	}

	Impl struct {
		Database *gorm.DB
	}
)

func (o *Impl) Close() error {
	if err := o.Database.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}

	return nil
}

func (o *Impl) FindOne(object interface{}, where ...interface{}) error {
	res := o.Database.First(object, where...)

	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.Wrap(err, fmt.Sprintf("could not find record %+v", where...))
		} else {
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func (o *Impl) Find(object interface{}, query interface{}, condition ...interface{}) error {
	res := o.Database.Where(query, condition...).Find(object)

	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to query %+v with condition %+v", query, condition)
	}

	return nil
}

func (o *Impl) Create(object interface{}) error {
	res := o.Database.Create(object)

	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create object %+v", object)
	}

	return nil
}

func (o *Impl) Update(object interface{}) error {
	res := o.Database.Update(object)

	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to update object %+v", object)
	}

	return nil
}

func (o *Impl) Delete(object interface{}) error {
	res := o.Database.Unscoped().Delete(object)

	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to delete object %+v", object)
	}

	return nil
}

func (o *Impl) SoftDelete(object interface{}) error {
	res := o.Database.Delete(object)

	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to soft delete object %+v", object)
	}

	return nil
}

func (o *Impl) Begin() ORM {
	return &Impl{Database: o.Database.Begin()}
}

func (o *Impl) Rollback() error {
	res := o.Database.Rollback()

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to rollback transaction!")
	}

	return nil
}

func (o *Impl) Commit() error {
	res := o.Database.Commit()

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to commit transaction!")
	}

	return nil
}

func (o *Impl) Exec(sql string, args ...interface{}) error {
	res := o.Database.Exec(sql, args...)

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to exec sql!")
	}

	return nil
}

func (o *Impl) RawSqlWithObject(sql string, object interface{}, args ...interface{}) error {
	res := o.Database.Raw(sql, args...).Scan(object)

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to query sql!")
	}

	return nil
}

func (o *Impl) RawSql(sql string, args ...interface{}) (*sql.Rows, error) {
	return o.Database.Raw(sql, args...).Rows()
}
