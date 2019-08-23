package db

import (
	"context"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DBCart Mongo
)

type Mongo interface {
	Client() *mongo.Client
	DB() *mongo.Database
}

type implementation struct {
	client   *mongo.Client
	database *mongo.Database
}

func init() {
	DBCart = Connect("hotel_cart")
}

// Connect : mongo db collection
func Connect(databaseName string) Mongo {
	// Set client options

	clientOptions := options.Client().ApplyURI(constant.DBURL)
	clientOptions.SetReadPreference(readpref.Primary())
	clientOptions.SetDirect(true)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Info("failed connect mongo")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB! (" + databaseName + ")")

	database := client.Database(databaseName)

	return &implementation{client: client, database: database}
}

func (i *implementation) Client() *mongo.Client {
	return i.client
}

func (i *implementation) DB() *mongo.Database {
	return i.database
}
