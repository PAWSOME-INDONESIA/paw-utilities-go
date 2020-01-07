package constant

const (
	// BASEURL ...
	BASEURL = "https://lb-stag.tiket.com"
	// BASEURL_AUTOCOMPLETE ...
	BASEURL_AUTOCOMPLETE = "http://192.168.40.35:7012"
	// BASEURL_CART ...
	BASEURL_CART = "http://54.169.65.147:7130"

	// DBURL ...
	DBURL = "mongodb://dev_hotel:7SXWhM7032n0UHb@13.250.57.121:27017"
	// DBURL = "mongodb://13.250.57.121:8080"
)

var (
	// SuccessMessage ...
	SuccessMessage = make(map[bool]string)
)

func init() {
	SuccessMessage[true] = "[ SUCCESS ]"
	SuccessMessage[false] = "[ FAILED ]"
}
