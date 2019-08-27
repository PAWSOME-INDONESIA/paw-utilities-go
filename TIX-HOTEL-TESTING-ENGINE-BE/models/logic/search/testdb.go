package search

import (
	"context"
	"strings"
	"time"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/db"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBHotelSearch ...
var DBHotelSearch db.Mongo

func init() {
	DBHotelSearch = db.Connect("hotel_search")
}

// GetPriority : test level DB
func GetPriority(ID string, searchType string, typee string, datee string) (elem structs.HotelSearchHotelPriorityRanking) {
	var (
		// resultDB []*structs.HotelSearchHotelPriorityRanking

		ctx = context.Background()
	)

	if _, ok := constant.SearchType[searchType]; !ok {
		return
	}

	log.Info("Database Test Case :")

	coll := DBHotelSearch.DB().Collection(constant.ColHotelPriorityRanking)
	opts := options.FindOneOptions{
		Projection: bson.M{},
		Sort:       bson.M{"_id": 1},
	}
	startDate, _ := time.Parse("2006-01-02", datee)
	filter := bson.D{
		{"type", typee},
		{constant.SearchType[searchType], ID},
		{"isDeleted", 0},
		{"startDate", bson.M{"$lte": startDate}},
		{"endDate", bson.M{"$gte": startDate}},
		{"cityId", ""},
		{"areaId", ""},
	}
	// if searchType == constant.SearchTypeRegion {
	// 	filter = append(filter, bson.E{"cityId", ""})
	// }
	// fmt.Println(typee, ID, constant.SearchType[searchType])
	err := coll.FindOne(ctx, filter, &opts).Decode(&elem)
	// defer cur.Close(ctx)
	if err != nil {
		// log.Warning(cur)

		if strings.Contains(err.Error(), "no documents") {
			log.Warning("Error : failed to find priority for ", ID)
		} else {
			log.Warning("error DB : ", err.Error())
		}
		return
	}

	// fmt.Println(elem.PublicID)
	return

}
