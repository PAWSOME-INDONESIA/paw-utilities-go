package constant

const (
	// SearchTypeRegion ...
	SearchTypeRegion = "region"
	// SearchTypeCity ...
	SearchTypeCity = "city"
	// SearchTypeArea ...
	SearchTypeArea = "area"
	// SearchTypeCoordinate ...
	SearchTypeCoordinate = "coordinate"
	// SearchTypePoi ...
	SearchTypePoi = "poi"
	// SearchTypeHotelChain ...
	SearchTypeHotelChain = "hotelChain"
)

const (
	SortTypePopularity = "popularity"
)

const (
	// BookPending ...
	BookPending = "PENDING"
	// BookSuccess ...
	BookSuccess = "SUCCESS"
)

var (
	// Vendor ...
	Vendor = make(map[string]string)

	// paymentMethod ...
	PaymentMethod = make(map[int]string)

	// SearchType ...
	SearchType = make(map[string]string)
)

func init() {
	Vendor["AGODA"] = "#FB387E"
	Vendor["EXPEDIA"] = "#FEDD00"
	Vendor["HOTELBEDS"] = "#008000"
	Vendor["TIKET"] = "#0064D2"

	PaymentMethod[1] = "pay_at_hotel"
	PaymentMethod[2] = "pay_now"
	PaymentMethod[3] = "pay_at_hotel_without_cc"

	SearchType[SearchTypeRegion] = "regionId"
	SearchType[SearchTypeCity] = "cityId"
}
