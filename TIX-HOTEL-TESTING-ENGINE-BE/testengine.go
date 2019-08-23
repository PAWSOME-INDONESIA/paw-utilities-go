package main

import (
	"encoding/json"
	goflag "flag"
	"io/ioutil"
	"os"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/logging"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var filename = flag.String("file", "test-default.json", "filename for testing (.json)")

func prepare() {
	logging.Log()

}

func main() {

	prepare()

	var (
		fileContent structs.FileContent
	)

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	log.Info("Read file " + *filename + " from storages\n")
	raw, err := ioutil.ReadFile("storages/" + *filename)
	if err != nil {
		log.Warning("cannot find file " + *filename + " in storages " + err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &fileContent)
	if err != nil {
		log.Warning("json file is invalid, " + err.Error())
		os.Exit(1)
	}

	for _, v := range fileContent.Testing {
		log.Info("============================================================================================================")
		logic.MainTest(v)
		log.Info("============================================================================================================\n")
	}
}
