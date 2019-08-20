package logic

import (
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/autocomplete"
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/book"
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/prebook"
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/room"
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/search"
	"TIX-HOTEL-TESTING-ENGINE-BE/models/logic/default"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/constant"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"os"

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
