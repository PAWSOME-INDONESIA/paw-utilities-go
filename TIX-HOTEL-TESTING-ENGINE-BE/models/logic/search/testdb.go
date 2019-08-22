package search

import (
	"TIX-HOTEL-TESTING-ENGINE-BE/models/db"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/constant"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"context"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
)

// DBHotelSearch ...
var DBHotelSearch db.Mongo

func init() {
	// DBHotelSearch = db.Connect("hotel_search")
}

// TestDB : test level DB
func TestDB(publicID string) {
	var (
		resultDB []*structs.HotelSearchHotel
	)

	log.Info("Database Test Case :")

	coll := DBHotelSearch.DB().Collection("hotel")
	cur, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Warning("error DB : ", err.Error())
	}

	for cur.Next(context.Background()) {
		var elem structs.HotelSearchHotel
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		resultDB = append(resultDB, &elem)
	}

	log.Info(resultDB, publicID)
	// Check data exist
	if len(resultDB) == 0 {
		log.Warning("1. Check hotel exist ", constant.SuccessMessage[false])
	} else {
		log.Info("1. Check hotel exist ", constant.SuccessMessage[true])
	}
}
