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
		FindOneWithContext(context.Context, string, interface{}, interface{}) error
		FindOne(string, interface{}, interface{}) error

		FindWithContext(context.Context, string, interface{}, FindCallback, ...*options.FindOptions) error
		Find(string, interface{}, FindCallback, ...*options.FindOptions) error

		FindOneAndDeleteWithContext(context.Context, string, interface{}) error
		FindOneAndDelete(string, interface{}) error

		FindOneAndUpdateWithContext(context.Context, string, interface{}, interface{}) error
		FindOneAndUpdate(string, interface{}, interface{}) error

		InsertWithContext(context.Context, string, interface{}) (*primitive.ObjectID, error)
		Insert(string, interface{}) (*primitive.ObjectID, error)

		InsertManyWithContext(context.Context, string, []interface{}) ([]primitive.ObjectID, error)
		InsertMany(string, []interface{}) ([]primitive.ObjectID, error)

		UpdateWithContext(context.Context, string, interface{}, interface{}) error
		Update(string, interface{}, interface{}) error

		UpdateManyWithContext(context.Context, string, interface{}, interface{}) error
		UpdateMany(string, interface{}, interface{}) error

		DeleteManyWithContext(context.Context, string, interface{}) error
		DeleteMany(string, interface{}) error

		DeleteWithContext(context.Context, string, interface{}) error
		Delete(string, interface{}) error

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

func (i *implementation) FindOneWithContext(ctx context.Context, collection string, filter, object interface{}) error {
	sr := i.database.Collection(collection).FindOne(ctx, filter)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOne failed!")
	}

	if err := sr.Decode(&object); err != nil {
		return errors.Wrap(err, "FindOne decode failed!")
	}

	return nil
}

func (i *implementation) FindOne(collection string, filter interface{}, object interface{}) error {
	return i.FindOneWithContext(context.Background(), collection, filter, object)
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

func (i *implementation) FindOneAndDeleteWithContext(ctx context.Context, collection string, filter interface{}) error {
	sr := i.database.Collection(collection).FindOneAndDelete(ctx, filter)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOneAndDeleteWithContext failed!")
	}

	return nil
}

func (i *implementation) FindOneAndDelete(collection string, filter interface{}) error {
	return i.FindOneAndDeleteWithContext(context.Background(), collection, filter)
}

func (i *implementation) FindOneAndUpdateWithContext(ctx context.Context, collection string, filter, object interface{}) error {
	sr := i.database.Collection(collection).FindOneAndUpdate(ctx, filter, object)

	if err := sr.Err(); err != nil {
		return errors.Wrap(err, "FindOneAndUpdateWithContext failed!")
	}

	if err := sr.Decode(&object); err != nil {
		return errors.Wrap(err, "FindOneAndUpdate decode failed!")
	}

	return nil
}

func (i *implementation) FindOneAndUpdate(collection string, filter, object interface{}) error {
	return i.FindOneAndUpdateWithContext(context.Background(), collection, filter, object)
}

func (i *implementation) InsertWithContext(ctx context.Context, collection string, object interface{}) (*primitive.ObjectID, error) {
	ir, err := i.database.Collection(collection).InsertOne(ctx, object)

	if err != nil {
		return nil, errors.Wrap(err, "InsertOneWithContext failed!")
	}

	id, ok := ir.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, errors.New("InsertWithContext failed to cast ObjectID")
	}

	return &id, nil
}

func (i *implementation) Insert(collection string, object interface{}) (*primitive.ObjectID, error) {
	return i.InsertWithContext(context.Background(), collection, object)
}

func (i *implementation) InsertManyWithContext(ctx context.Context, collection string, documents []interface{}) ([]primitive.ObjectID, error) {
	ir, err := i.database.Collection(collection).InsertMany(ctx, documents)

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

func (i *implementation) InsertMany(collection string, documents []interface{}) ([]primitive.ObjectID, error) {
	return i.InsertManyWithContext(context.Background(), collection, documents)
}

func (i *implementation) UpdateWithContext(ctx context.Context, collection string, filter, object interface{}) error {
	if _, err := i.database.Collection(collection).UpdateOne(ctx, filter, object); err != nil {
		return errors.Wrap(err, "UpdateWithContext failed!")
	}

	return nil
}

func (i *implementation) Update(collection string, filter, object interface{}) error {
	return i.UpdateWithContext(context.Background(), collection, filter, object)
}

func (i *implementation) UpdateManyWithContext(ctx context.Context, collection string, filter, object interface{}) error {
	if _, err := i.database.Collection(collection).UpdateMany(ctx, filter, object); err != nil {
		return errors.Wrap(err, "UpdateManyWithContext failed!")
	}

	return nil
}

func (i *implementation) UpdateMany(collection string, filter, object interface{}) error {
	return i.UpdateManyWithContext(context.Background(), collection, filter, object)
}

func (i *implementation) DeleteManyWithContext(ctx context.Context, collection string, filter interface{}) error {
	if _, err := i.database.Collection(collection).DeleteMany(ctx, filter); err != nil {
		return errors.Wrap(err, "DeleteManyWithContext failed!")
	}

	return nil
}

func (i *implementation) DeleteMany(collection string, filter interface{}) error {
	return i.DeleteManyWithContext(context.Background(), collection, filter)
}

func (i *implementation) DeleteWithContext(ctx context.Context, collection string, filter interface{}) error {
	if _, err := i.database.Collection(collection).DeleteOne(ctx, filter); err != nil {
		return errors.Wrap(err, "DeleteWithContext failed!")
	}

	return nil
}

func (i *implementation) Delete(collection string, filter interface{}) error {
	return i.DeleteWithContext(context.Background(), collection, filter)
}

func (i *implementation) Indexes(collection string) mongo.IndexView {
	return i.database.Collection(collection).Indexes()
}
