package _default

import (
	"encoding/json"
	"os"
	"time"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type (
	CommandDefault struct {
	}

	Expected struct {
		Url          string      `json:"url"`
		Method       string      `json:"method"`
		ResponseData interface{} `json:"responseData"`
	}

	// BodyResponse ...
	B2bSearchresponse struct {
		Data []ContentlList `json:"data"`
	}

	// ContentlList ...
	ContentlList struct {
		Id             string      `json:"id"`
		Name           string      `json:"name"`
		PostalCode     string      `json:"postalCode"`
		StarRating     float64     `json:"starRating"`
		PaymentOptions string      `json:"paymentOptions"`
		AvailableRoom  int         `json:"availableRoom"`
		Refundable     interface{} `json:"refundable"`
		RateInfo       interface{} `json:"rateInfo"`
		Promo          interface{} `json:"promo"`
		ValueAdded     interface{} `json:"valueAdded"`
		SoldOut        bool        `json:"soldOut"`
		Location       struct {
			Coordinate struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"coordinates"`
		} `json:"location"`
	}
)

func (d *CommandDefault) Test(contentList structs.ContentList) {
	var (
		err               error
		expected          Expected
		expectedResponse  []ContentlList
		b2bSearchresponse B2bSearchresponse
	)

	dataByte, _ := json.Marshal(contentList.Data)

	dataInterface := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &dataInterface)
	if err != nil {
		log.Warning("error unmarshal data  :", err.Error())
	}

	// Change StartDate if now
	if dataInterface["startDate"] == "now" {
		dataInterface["startDate"] = time.Now().Format("2006-01-02")
	}

	expectedByte, _ := json.Marshal(contentList.Expected)
	err = json.Unmarshal(expectedByte, &expected)
	if err != nil {
		log.Warning("error unmarshal expected :", err.Error())
		return
	}

	log.Info("Test Case :")

	url := constant.BASEURL + expected.Url

	res := util.CallRest(expected.Method, dataInterface, contentList.Header, url)
	statusCode := util.GetStatusCode(res)

	//Check Status Code must == 200
	if statusCode != 200 {
		log.Error(res.Body)
		os.Exit(1)
		return
	} else {

		log.Info("1. Status Code must be 200 ", constant.SuccessMessage[true])
	}

	respBody := util.GetResponseBody(res)
	messageCode := gjson.Get(respBody, "code").String()
	arrayData := gjson.Get(respBody, "data").Array()

	//Check Message Code Is SUCCESS && Data not null add info log
	if messageCode == constant.SUCCESS && len(arrayData) == 0 {
		log.Info("Empty Data")
		os.Exit(1)
		return
	}

	log.Info("2. Data not null", constant.SuccessMessage[true])

	// Get Expected Response
	expectedRespByte, _ := json.Marshal(expected.ResponseData)
	err = json.Unmarshal(expectedRespByte, &expectedResponse)
	if err != nil {
		log.Warning("error unmarshal expected response :", err.Error())
		return
	}

	checkExpected := make(map[int]bool, 0)
	for k := range expectedResponse {
		checkExpected[k] = false
	}

	err = json.Unmarshal([]byte(respBody), &b2bSearchresponse)
	if err != nil {
		log.Warning("error unmarshal b2bSearchresponse  :", err.Error())
		return
	}

	for _, value := range b2bSearchresponse.Data {
		// Check response
		for keyexpect, valueExpect := range expectedResponse {

			if util.CheckSimilarStruct(valueExpect, value) {
				checkExpected[keyexpect] = true
			}
		}
	}

	if len(expectedResponse) > 0 {
		log.Info("3. Check Expected response:")
		for keyexpect := range expectedResponse {
			log.Info("Expected response ", (keyexpect + 1), constant.SuccessMessage[checkExpected[keyexpect]])
		}
	}

}
