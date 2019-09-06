package logic

import (
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/searchseo"
	"os"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/autocomplete"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/book"
	_default "github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/default"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/prebook"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/room"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/search"

	log "github.com/sirupsen/logrus"
)

// Command ...
type Command interface {
	Test(contentList structs.ContentList)
}

// CommandTest ...
var CommandTest = make(map[string]Command)

func init() {
	CommandTest[constant.CommandSearch] = new(search.CommandSearch)
	CommandTest[constant.CommandAutocomplete] = new(autocomplete.CommandAutocomplete)
	CommandTest[constant.CommandRoom] = new(room.CommandRoom)
	CommandTest[constant.CommandPrebook] = new(prebook.CommandPrebook)
	CommandTest[constant.CommandBook] = new(book.CommandBook)
	CommandTest[constant.CommandDefault] = new(_default.CommandDefault)
	CommandTest[constant.CommandSearchSEO] = new(searchseo.CommandSearchSEO)
}

// MainTest : Main logic of test
func MainTest(contentList structs.ContentList) {

	if _, ok := CommandTest[contentList.Command]; !ok {
		log.Warning("Command not found : " + contentList.Command)
		os.Exit(1)
		return
	}

	log.Info("Command executed : " + contentList.Command)

	CommandTest[contentList.Command].Test(contentList)

}
