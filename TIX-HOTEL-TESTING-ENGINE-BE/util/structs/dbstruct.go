package structs

type (
	// HotelSearchHotel ...
	HotelSearchHotel struct {
		PublicID string `bson:"publicId"`
		Name     string `bson:"name"`
	}

	// HotelCartBook ...
	HotelCartBook struct {
		OrderID int64  `bson:"orderId"`
		Status  string `bson:"status"`
	}

	// HotelSearchHotelPriorityRanking ...
	HotelSearchHotelPriorityRanking struct {
		Type     string   `bson:"type"`
		PublicID []string `bson:"publicIds"`
	}
)
