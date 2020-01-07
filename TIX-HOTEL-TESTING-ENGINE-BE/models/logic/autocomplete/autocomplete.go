package autocomplete

import (
	"encoding/json"
	"os"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// CommandSearch ...
type (
	CommandAutocomplete struct{}

	Expected struct {
		HotelCount    int         `json:"hotelCount"`
		StatusCode    int         `json:"statusCode"`
		ResponseData  interface{} `json:"responseData"`
		ExactResponse bool        `json:"exactResponse"`
	}

	// HotelSearchResponse ...
	AutocompleteResponse struct {
		Data []ContentlList `json:"data"`
	}

	// ContentlList ...
	ContentlList struct {
		Type       string  `json:"type"`
		PublicId   string  `json:"publicId"`
		Name       string  `json:"name"`
		Location   string  `json:"location"`
		HotelCount float64 `json:"hotelCount"`
		Country    IDName  `json:"country"`
		Region     IDName  `json:"region"`
		City       IDName  `json:"city"`
		Area       IDName  `json:"area"`
	}

	// IDName ...
	IDName struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

// Test ...
func (d *CommandAutocomplete) Test(contentList structs.ContentList) {
	var (
		err                  error
		expectedResponse     []ContentlList
		expected             Expected
		autocompleteResponse AutocompleteResponse
		check = true
	)

	dataByte, _ := json.Marshal(contentList.Expected)

	url := constant.BASEURL_AUTOCOMPLETE + constant.URLCommand[contentList.Command]

	res := util.CallRest("POST", contentList.Data, contentList.Header, url)
	statusCode := util.GetStatusCode(res)
	expectStatusCode := gjson.GetBytes(dataByte, "statusCode").Int()

	//Check Status Code must == 200
	if statusCode != int(expectStatusCode) {
		log.Error(expectStatusCode)
		log.Error(res.Body)
		os.Exit(1)
		return
	}
	respBody := util.GetResponseBody(res)
	err = json.Unmarshal([]byte(respBody), &autocompleteResponse)
	if err != nil {
		//log.Warning("2. Check hotel list must > 0", constant.SuccessMessage[false])
		log.Warning(err.Error())
		os.Exit(1)
		return
	}

	arrayData := gjson.Get(respBody, "data").Array()

	messageCode := gjson.Get(respBody, "code").String()

	//Check Message Code != SUCCESS add err log
	if messageCode != constant.SUCCESS {
		log.Error(res.Body)
		os.Exit(1)
		return
	}
	log.Info("Test Case :")
	log.Info("1. Status Code must be 200", constant.SuccessMessage[true])

	//Check Message Code Is SUCCESS && Data not null add info log
	if messageCode == constant.SUCCESS && len(arrayData) == 0 {
		log.Info("Empty Data")
		os.Exit(1)
		return
	}

	log.Info("2. Data not null", constant.SuccessMessage[true])

	//Check Expected HotelCount NOT NULL with Field Type=['AREA', 'REGION', 'COUNTRY', 'CITY']
	var resultHotelCountNotNull []string
	//var checkPublicIdAndType bool
	for _, v := range arrayData {
		typeData := gjson.Get(v.String(), "type").String()
		hotelCount := gjson.Get(v.String(), "hotelCount").Int()
		name := gjson.Get(v.String(), "name").String()
		// log.Info(hotelCount, gjson.GetBytes(dataByte, "hotelCount").Int())
		if typeData == "AREA" || typeData == "REGION" || typeData == "COUNTRY" || typeData == "CITY" {
			if hotelCount < gjson.GetBytes(dataByte, "hotelCount").Int() {
				resultHotelCountNotNull = append(resultHotelCountNotNull, name+`(type is `+typeData+`)`)
			}
		}
	}

	if len(resultHotelCountNotNull) > 0 {
		log.Info("3. if type=['AREA', 'REGION', 'COUNTRY', 'CITY'] hotelCount must >= 1 [ FAILED ]")
		for _, v := range resultHotelCountNotNull {
			log.Info(v)
		}

		os.Exit(1)
		return
	} else {
		log.Info("3. if type=['AREA', 'REGION', 'COUNTRY', 'CITY'] hotelCount must >= 1 [ SUCCESS ]")
	}

	// Get expected
	expectedByte, _ := json.Marshal(contentList.Expected)
	err = json.Unmarshal(expectedByte, &expected)
	if err != nil {
		log.Warning("error unmarshal expected :", err.Error())
		return
	}

	// Get Expected Response
	expectedRespByte, _ := json.Marshal(expected.ResponseData)
	err = json.Unmarshal(expectedRespByte, &expectedResponse)
	if err != nil {
		log.Warning("error unmarshal expected :", err.Error())
		return
	}

	checkExpected := make(map[int]bool, 0)
	for k := range expectedResponse {
		checkExpected[k] = false
	}

	for key, value := range autocompleteResponse.Data {
		// Check response
		if expected.ExactResponse {
			for keyexpect, valueExpect := range expectedResponse {
				if util.CheckSimilarStruct(valueExpect, value) && key == keyexpect {
					checkExpected[keyexpect] = true
				}
			}
		} else {
			for keyexpect, valueExpect := range expectedResponse {
				if util.CheckSimilarStruct(valueExpect, value) {
					checkExpected[keyexpect] = true
				}
			}
		}

	}

	if len(expectedResponse) > 0 {
		log.Info("4. Check Expected response:")
		for keyexpect := range expectedResponse {
			expect := checkExpected[keyexpect]
			log.Info("Expected response ", (keyexpect + 1), constant.SuccessMessage[expect])
			if !expect {
				check = false
			}
		}

		if !check {
			mar, _ := json.Marshal(autocompleteResponse.Data)
			log.Infof("response %s", string(mar))
		}
	}

}
