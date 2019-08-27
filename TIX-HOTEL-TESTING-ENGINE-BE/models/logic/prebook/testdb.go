package prebook

import (
	"context"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/db"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DBHotel ...
var (
	DBHotel = db.DBCart
)

// func init() {
// 	DBHotel = db.Connect("hotel_cart")
// }

// TestDB : test level DB
func TestDB(id string) {
	var (
		resultDB []*structs.HotelCartBook
		cekRow   = make(map[string]map[string]bool, 0)
	)

	log.Info("Database Test Case :")

	coll := DBHotel.DB().Collection("book")
	prim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warn(id)
		log.Fatal(err)
	}

	cur, err := coll.Find(context.Background(), bson.M{"_id": prim})
	if err != nil {
		log.Warning("error DB : ", err.Error())
	}

	for cur.Next(context.Background()) {
		var elem structs.HotelCartBook
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		resultDB = append(resultDB, &elem)

		// Cek status must pending
		// log.Info(elem)
		cekStasus := map[string]bool{"cek status": false}
		if elem.Status == constant.BookPending {
			cekStasus = map[string]bool{"cek status": true}
		}
		cekRow[id] = cekStasus
		// log.Info(cekRow)
	}

	// log.Info(resultDB, id)
	// Check data exist
	if len(resultDB) == 0 {
		log.Warning("--- Check prebook exist ", constant.SuccessMessage[false])
	} else {
		log.Info("--- Check prebook exist ", constant.SuccessMessage[true])
	}

	log.Info("Check row DB prebook :")
	for key, value := range cekRow {
		log.Info("--- [cek status] ", key, " ", constant.SuccessMessage[value["cek status"]])
	}
}
