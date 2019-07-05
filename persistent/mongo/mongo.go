package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FindCallback func(*mgo.Cursor, error) error

	Mongo interface {
		FindOneWithContext(context.Context, string, interface{}, interface{}, ...*options.FindOneOptions) error
		FindOne(string, interface{}, interface{}, ...*options.FindOneOptions) error

		FindWithContext(context.Context, string, interface{}, FindCallback, ...*options.FindOptions) error
		Find(string, interface{}, FindCallback, ...*options.FindOptions) error

		FindOneAndDeleteWithContext(context.Context, string, interface{}, ...*options.FindOneAndDeleteOptions) error
		FindOneAndDelete(string, interface{}, ...*options.FindOneAndDeleteOptions) error

		FindOneAndUpdateWithContext(context.Context, string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error
		FindOneAndUpdate(string, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error

		InsertWithContext(context.Context, string, interface{}, ...*options.InsertOneOptions) (*primitive.ObjectID, error)
		Insert(string, interface{}, ...*options.InsertOneOptions) (*primitive.ObjectID, error)

		InsertManyWithContext(context.Context, string, []interface{}, ...*options.InsertManyOptions) ([]primitive.ObjectID, error)
		InsertMany(string, []interface{}, ...*options.InsertManyOptions) ([]primitive.ObjectID, error)

		UpdateWithContext(context.Context, string, interface{}, interface{}, ...*options.UpdateOptions) error
		Update(string, interface{}, interface{}, ...*options.UpdateOptions) error

		UpdateManyWithContext(context.Context, string, interface{}, interface{}, ...*options.UpdateOptions) error
		UpdateMany(string, interface{}, interface{}, ...*options.UpdateOptions) error

		DeleteManyWithContext(context.Context, string, interface{}, ...*options.DeleteOptions) error
		DeleteMany(string, interface{}, ...*options.DeleteOptions) error

		DeleteWithContext(context.Context, string, interface{}, ...*options.DeleteOptions) error
		Delete(string, interface{}, ...*options.DeleteOptions) error

		// - DDL
		Indexes(string) mongo.IndexView
	}

	implementation struct {
		client   *mgo.Client
		database *mgo.Database
		logger   logs.Logger
	}
)

func New(ctx context.Context, uri, name string, logger logs.Logger) (Mongo, error) {
	if uri == "" {
		return nil, errors.New("uri is required!")
	}

	if name == "" {
		return nil, errors.New("database name is required!")
	}

	if logger == nil {
		return nil, errors.New("logger is required!")
	}

	client, err := mgo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo!")
	}

	database := client.Database(name)

	return &implementation{client, database, logger}, nil
}

func (i *implementation) FindOneWithContext(ctx context.Context, collection string, filter, object interface{}, options ...*options.FindOneOptions) error {
	sr := i.database.Collection(collection).FindOne(ctx, filter, options...)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOne failed!")
	}

	if err := sr.Decode(&object); err != nil {
		return errors.Wrap(err, "FindOne decode failed!")
	}

	return nil
}

func (i *implementation) FindOne(collection string, filter interface{}, object interface{}, options ...*options.FindOneOptions) error {
	return i.FindOneWithContext(context.Background(), collection, filter, object, options...)
}

func (i *implementation) FindWithContext(ctx context.Context, collection string, filter interface{}, callback FindCallback, options ...*options.FindOptions) error {
	cursor, err := i.database.Collection(collection).Find(context.Background(), filter, options...)

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			i.logger.Errorf("failed to close cursor %s", err)
		}
	}()

	if err != nil {
		return callback(nil, err)
	} else {
		return callback(cursor, nil)
	}
}

func (i *implementation) Find(collection string, filter interface{}, callback FindCallback, options ...*options.FindOptions) error {
	return i.FindWithContext(context.Background(), collection, filter, callback, options...)
}

func (i *implementation) FindOneAndDeleteWithContext(ctx context.Context, collection string, filter interface{}, options ...*options.FindOneAndDeleteOptions) error {
	sr := i.database.Collection(collection).FindOneAndDelete(ctx, filter, options...)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOneAndDeleteWithContext failed!")
	}

	return nil
}

func (i *implementation) FindOneAndDelete(collection string, filter interface{}, options ...*options.FindOneAndDeleteOptions) error {
	return i.FindOneAndDeleteWithContext(context.Background(), collection, filter, options...)
}

func (i *implementation) FindOneAndUpdateWithContext(ctx context.Context, collection string, filter, object interface{}, options ...*options.FindOneAndUpdateOptions) error {
	sr := i.database.Collection(collection).FindOneAndUpdate(ctx, filter, object, options...)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOneAndUpdateWithContext failed!")
	}

	if err := sr.Decode(&object); err != nil {
		return errors.Wrap(err, "FindOneAndUpdate decode failed!")
	}

	return nil
}

func (i *implementation) FindOneAndUpdate(collection string, filter, object interface{}, options ...*options.FindOneAndUpdateOptions) error {
	return i.FindOneAndUpdateWithContext(context.Background(), collection, filter, object, options...)
}

func (i *implementation) InsertWithContext(ctx context.Context, collection string, object interface{}, options ...*options.InsertOneOptions) (*primitive.ObjectID, error) {
	ir, err := i.database.Collection(collection).InsertOne(ctx, object, options...)

	if err != nil {
		return nil, errors.Wrap(err, "InsertOneWithContext failed!")
	}

	id, ok := ir.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, errors.New("InsertWithContext failed to cast ObjectID")
	}

	return &id, nil
}

func (i *implementation) Insert(collection string, object interface{}, options ...*options.InsertOneOptions) (*primitive.ObjectID, error) {
	return i.InsertWithContext(context.Background(), collection, object, options...)
}

func (i *implementation) InsertManyWithContext(ctx context.Context, collection string, documents []interface{}, options ...*options.InsertManyOptions) ([]primitive.ObjectID, error) {
	ir, err := i.database.Collection(collection).InsertMany(ctx, documents, options...)

	if err != nil {
		return nil, errors.Wrap(err, "InsertManyWithContext failed!")
	}

	ids := make([]primitive.ObjectID, 0)

	for _, id := range ir.InsertedIDs {
		i, ok := id.(primitive.ObjectID)

		if !ok {
			err = errors.Errorf("InsertWithContext failed to cast ObjectID %s", i)
			break
		}

		ids = append(ids, i)
	}

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (i *implementation) InsertMany(collection string, documents []interface{}, options ...*options.InsertManyOptions) ([]primitive.ObjectID, error) {
	return i.InsertManyWithContext(context.Background(), collection, documents, options...)
}

func (i *implementation) UpdateWithContext(ctx context.Context, collection string, filter, object interface{}, options ...*options.UpdateOptions) error {
	if _, err := i.database.Collection(collection).UpdateOne(ctx, filter, object, options...); err != nil {
		return errors.Wrap(err, "UpdateWithContext failed!")
	}

	return nil
}

func (i *implementation) Update(collection string, filter, object interface{}, options ...*options.UpdateOptions) error {
	return i.UpdateWithContext(context.Background(), collection, filter, object, options...)
}

func (i *implementation) UpdateManyWithContext(ctx context.Context, collection string, filter, object interface{}, options ...*options.UpdateOptions) error {
	if _, err := i.database.Collection(collection).UpdateMany(ctx, filter, object, options...); err != nil {
		return errors.Wrap(err, "UpdateManyWithContext failed!")
	}

	return nil
}

func (i *implementation) UpdateMany(collection string, filter, object interface{}, options ...*options.UpdateOptions) error {
	return i.UpdateManyWithContext(context.Background(), collection, filter, object, options...)
}

func (i *implementation) DeleteManyWithContext(ctx context.Context, collection string, filter interface{}, options ...*options.DeleteOptions) error {
	if _, err := i.database.Collection(collection).DeleteMany(ctx, filter, options...); err != nil {
		return errors.Wrap(err, "DeleteManyWithContext failed!")
	}

	return nil
}

func (i *implementation) DeleteMany(collection string, filter interface{}, options ...*options.DeleteOptions) error {
	return i.DeleteManyWithContext(context.Background(), collection, filter, options...)
}

func (i *implementation) DeleteWithContext(ctx context.Context, collection string, filter interface{}, options ...*options.DeleteOptions) error {
	if _, err := i.database.Collection(collection).DeleteOne(ctx, filter, options...); err != nil {
		return errors.Wrap(err, "DeleteWithContext failed!")
	}

	return nil
}

func (i *implementation) Delete(collection string, filter interface{}, options ...*options.DeleteOptions) error {
	return i.DeleteWithContext(context.Background(), collection, filter, options...)
}

func (i *implementation) Indexes(collection string) mongo.IndexView {
	return i.database.Collection(collection).Indexes()
}
