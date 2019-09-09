package constant

const (
	// CommandAutocomplete ...
	CommandAutocomplete = "autocomplete"
	// CommandSearch ...
	CommandSearch = "search"
	// CommandRoom ...
	CommandRoom = "room"
	// CommandPrebook ...
	CommandPrebook = "prebook"
	// CommandBook ...
	CommandBook = "book"
	// CommandGeneral
	CommandDefault = "default"
	// CommandSearchSEO
	CommandSearchSEO = "searchseo"
)

// URLCommand ...
var URLCommand = make(map[string]string, 0)

func init() {
	URLCommand[CommandAutocomplete] = "/tix-hotel-autocomplete/autocomplete"
	URLCommand[CommandSearch] = "/tix-hotel-search/search"
	URLCommand[CommandRoom] = "/tix-hotel-search/room"
	URLCommand[CommandPrebook] = "/tix-hotel-cart/hotel/prebook"
	URLCommand[CommandBook] = "/tix-hotel-cart/hotel/book"
	URLCommand[CommandDefault] = "/tix-hotel-search/general"
	URLCommand[CommandSearchSEO] = "/tix-hotel-search/search/seo"
}
