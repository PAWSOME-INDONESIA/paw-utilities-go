package example

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/persistent/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	VendorShutdown struct {
		Vendor        string                   `json:"vendor"`
		IsActive      int                      `json:"isActive"`
	}
)

func main(){
	log, _ := logs.DefaultLog()
	db, _ := mongo.New(context.Background(), "mongodb://localhost:27017", "hotel_search", log)

	filter := bson.M{"storeId": "TIKETCOM"}
	opt := options.Find()
	opt.SetLimit(10)

	var result []VendorShutdown
	if err := db.FindWithContext(context.Background(), "vendor_shutdown_schedule", filter, func(cursor *mgo.Cursor, e error) error {
		if e != nil {
			return errors.WithStack(e)
		}
		for cursor.Next(context.Background()) {
			var data VendorShutdown
			if err := cursor.Decode(&data); err != nil {
				log.Error(err)
				return errors.Wrap(err, "failed to decode data")
			}
			result = append(result, data)
		}
		return nil
	}, opt); err != nil {
		log.Error(err)
	}
	log.Info(result)
}
