package search

import (
	"TIX-HOTEL-TESTING-ENGINE-BE/util"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/constant"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"encoding/json"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

type (
	// CommandSearch ...
	CommandSearch struct{}

	// Data ...
	Data struct {
		Adult       int    `json:"adult"`
		Filter      Filter `json:"filter"`
		Night       int    `json:"night"`
		Page        int    `json:"page"`
		Priority    string `json:"priorityRankingType"`
		Room        int    `json:"room"`
		SearchType  string `json:"searchType"`
		SearchValue string `json:"searchValue"`
		Sort        string `json:"sort"`
		StartDate   string `json:"startDate"`
	}

	// Filter ...
	Filter struct {
		PaymentOptions []string `json:"paymentOptions"`
	}

	// Expected ...
	Expected struct {
		ResponseData interface{} `json:"responseData"`
	}

	// HotelSearchResponse ...
	HotelSearchResponse struct {
		Data    Contents `json:"data"`
		Message string   `json:"message"`
	}

	// Contents ...
	Contents struct {
		ContentList []ContentlList `json:"contents"`
	}

	// ContentlList ...
	ContentlList struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Address    string  `json:"address"`
		PostalCode string  `json:"postalCode"`
		StarRating float64 `json:"starRating"`
		Country    IDName  `json:"country"`
		Region     IDName  `json:"region"`
		City       IDName  `json:"city"`
		Area       IDName  `json:"area"`
		Category   IDName  `json:"category"`
		Location   struct {
			Coordinate struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"coordinates"`
		} `json:"location"`
		AvailRoom      int         `json:"availableRoom"`
		Reviews        interface{} `json:"reviews"`
		MainImage      interface{} `json:"mainImage"`
		MainFacilities interface{} `json:"mainFacilities"`
		PoiDistance    float64     `json:"poiDistance"`
		ItemColor      interface{} `json:"itemColor"`
		Refundable     bool        `json:"refundable"`
		RateInfo       interface{} `json:"rateInfo"`
		Promo          interface{} `json:"promo"`
		ValueAdded     interface{} `json:"valueAdded"`
		PaymentOption  string      `json:"paymentOption"`
		Soldout        bool        `json:"soldOut"`
	}

	// IDName ...
	IDName struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

// Test ...
func (d *CommandSearch) Test(contentList structs.ContentList) {
	var (
		err                      error
		data                     Data
		hotelSearchResp          HotelSearchResponse
		expected                 Expected
		expectedResponse         []ContentlList
		checkMatchSearchType     = true
		checkAvailRoom           = true
		checkNameEmpty           = true
		checkFilterPaymentOption = true
		logMatchSearchType       = make([]string, 0)
		logAvailRoom             = make([]string, 0)
		logNameEmpty             = make([]string, 0)
		logFilterPaymentOption   = make([]string, 0)
	)

	// Get Data
	dataByte, _ := json.Marshal(contentList.Data)
	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		log.Warning("error unmarshal data :", err.Error())
		return
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

	// Change StartDate if now
	if data.StartDate == "now" {
		data.StartDate = time.Now().Format("2006-01-02")
	}

	url := constant.BASEURL + constant.URLCommand[contentList.Command]
	res := util.CallRest("POST", data, contentList.Header, url)
	statusCode := util.GetStatusCode(res)
	respBody := util.GetResponseBody(res)

	log.Info("Test Case :")

	// 1. Check Status code
	if statusCode != 200 {
		log.Warning("1. Status Code must be 200 ", constant.SuccessMessage[false])
		log.Warning("\nResponse : ", respBody)
	} else {
		log.Info("1. Status Code must be 200", constant.SuccessMessage[true])
	}

	// 2. Check hotel list must have data (>0)
	err = json.Unmarshal([]byte(respBody), &hotelSearchResp)
	if err != nil {
		log.Warning("2. Check hotel list must > 0", constant.SuccessMessage[false])
		log.Warning(err.Error())
		log.Warning("\nResponse : ", respBody)
		return
	}

	var locationDetail []string
	result2 := gjson.Get(respBody, "data.searchDetail.searchLocation").Array()
	for _, v := range result2 {
		locationDetail = append(locationDetail, v.String())
	}

	result3 := gjson.Get(respBody, "data.contents").Array()
	locationInIndonesia := util.StringContaintsInSlice("indonesia", locationDetail)
	if !locationInIndonesia && util.StringContaintsInSlice(constant.PaymentMethod[1], data.Filter.PaymentOptions) {
		if len(result3) > 0 {
			log.Info("2. Expected Response if paymentOptions is ['" + constant.PaymentMethod[1] + "'] and country other than Indonesia : List Hotels Must Null " + constant.SuccessMessage[false])
		} else {
			log.Info("2. Expected Response if paymentOptions is ['" + constant.PaymentMethod[1] + "'] and country other than Indonesia : List Hotels Must Null " + constant.SuccessMessage[true])
		}

		if len(hotelSearchResp.Data.ContentList) > 0 {
			for _, value := range hotelSearchResp.Data.ContentList {
				// check if filter pay_at_hotel, should not have data here because it's not indo
				if value.PaymentOption == constant.PaymentMethod[1] {
					logFilterPaymentOption = append(logFilterPaymentOption, "\nHotel ID ", value.ID, value.PaymentOption)
					checkFilterPaymentOption = false
				}
			}

			log.Info("3. Check Hotel list should not have payment options : ["+constant.PaymentMethod[1]+"] "+strings.Join(data.Filter.PaymentOptions, ","),
				constant.SuccessMessage[checkFilterPaymentOption])
			if !checkFilterPaymentOption {
				log.Warning("Log Message :")
				log.Warning(logFilterPaymentOption)
			}
		}

		return

	}

	if len(hotelSearchResp.Data.ContentList) == 0 {
		log.Warning("2. Check hotel list must > 0", constant.SuccessMessage[false])
		log.Warning("\nResponse : ", hotelSearchResp)
		checkAvailRoom, checkFilterPaymentOption, checkMatchSearchType, checkNameEmpty = false, false, false, false
	} else {
		log.Info("2. Check hotel list must > 0", constant.SuccessMessage[true])
	}

	// 3. Check Response Message Success
	if strings.ToUpper(hotelSearchResp.Message) != constant.SUCCESS {
		log.Warning("3. Check Response Code must success ", constant.SuccessMessage[false])
	} else {
		log.Info("3. Check Response Code must success ", constant.SuccessMessage[true])
	}

	// Check
	// 4. Check hotel list should match search type (REGION, POI, etc ...)
	// 5. Available Room must be (>= req.room)

	for _, value := range hotelSearchResp.Data.ContentList {

		// Check Match search type
		switch searchType := data.SearchType; searchType {
		case constant.SearchTypeRegion:
			{
				if value.Region.ID != data.SearchValue {
					checkMatchSearchType = false
					logMatchSearchType = append(logMatchSearchType, "\nHotel ID ", value.ID, searchType, value.Region.ID)
				}
			}
		case constant.SearchTypeCity:
			{
				if value.City.ID != data.SearchValue {
					checkMatchSearchType = false
					logMatchSearchType = append(logMatchSearchType, "\nHotel ID ", value.ID, searchType, value.City.ID)
				}

			}
		case constant.SearchTypeArea:
			{
				if value.Area.ID != data.SearchValue {
					checkMatchSearchType = false
					logMatchSearchType = append(logMatchSearchType, "\nHotel ID ", value.ID, searchType, value.Area.ID)
				}
			}
		}

		// Check avail room
		if value.AvailRoom < data.Room {
			checkAvailRoom = false
			logAvailRoom = append(logAvailRoom, "\nHotel ID ", value.ID, string(value.AvailRoom))
		}

		// Check response
		for keyexpect, valueExpect := range expectedResponse {

			if util.CheckSimilarStruct(valueExpect, value) {
				checkExpected[keyexpect] = true
			}
		}

		// Check name empty
		if &value.Name == nil || value.Name == " " || value.Name == "" {
			checkNameEmpty = false
			logNameEmpty = append(logNameEmpty, "\nHotel ID ", value.ID, string(value.Name))
		}

		// Check filter payment option
		if len(data.Filter.PaymentOptions) > 0 {
			isFound := false
			for _, valPayment := range data.Filter.PaymentOptions {
				// log.Info(valPayment, " == ", value.PaymentOption)
				if value.PaymentOption == valPayment {
					isFound = true
				}
			}
			if !isFound {
				logFilterPaymentOption = append(logFilterPaymentOption, "\nHotel ID ", value.ID, value.PaymentOption)
				checkFilterPaymentOption = false
			}
		} else {
			if !locationInIndonesia {
				// check pay_at_hotel should not in list of other indo
				if value.PaymentOption == constant.PaymentMethod[1] {
					logFilterPaymentOption = append(logFilterPaymentOption, "\nHotel ID ", value.ID, value.PaymentOption)
					checkFilterPaymentOption = false
				}
			}
		}
	}

	log.Info("4. Check hotel list should match search type (", data.SearchType, ")",
		constant.SuccessMessage[checkMatchSearchType])
	if !checkMatchSearchType {
		log.Warning("Log Message :")
		log.Warning(logMatchSearchType)
	}

	log.Info("5. Check Available Room should >= (", data.Room, ")",
		constant.SuccessMessage[checkAvailRoom])
	if !checkAvailRoom {
		log.Warning("Log Message :")
		log.Warning(logAvailRoom)
	}

	if len(expectedResponse) > 0 {
		log.Info("6. Check Expected response:")
		for keyexpect := range expectedResponse {
			log.Info("Expected response ", (keyexpect + 1), constant.SuccessMessage[checkExpected[keyexpect]])
		}
	}

	log.Info("7. Check Hotel Name is not empty",
		constant.SuccessMessage[checkNameEmpty])
	if !checkNameEmpty {
		log.Warning("Log Message :")
		log.Warning(logNameEmpty)
	}

	if len(data.Filter.PaymentOptions) > 0 {
		log.Info("8. Check Hotel should match with filter payment options : "+strings.Join(data.Filter.PaymentOptions, ","),
			constant.SuccessMessage[checkFilterPaymentOption])
		if !checkFilterPaymentOption {
			log.Warning("Log Message :")
			log.Warning(logFilterPaymentOption)
		}
	} else {
		if !locationInIndonesia {
			// check pay_at_hotel should not in list of other indo
			log.Info("8. Check Hotel should not have payment options : "+constant.PaymentMethod[1],
				constant.SuccessMessage[checkFilterPaymentOption])
			if !checkFilterPaymentOption {
				log.Warning("Log Message :")
				log.Warning(logFilterPaymentOption)
			}
		}
	}

	//interfaceData := make(map[string]interface{})
	//marshalRequestData, _ := json.Marshal(contentList.Data)
	//err = json.Unmarshal(marshalRequestData, &interfaceData)
	//if err != nil {
	//	log.Warning("error unmarshal interfaceData :", err.Error())
	//}
	//if interfaceData["startDate"] == "now" {
	//	interfaceData["startDate"] = time.Now().Format("2006-01-02")
	//}
	//
	//res2 := util.CallRest("POST", interfaceData, contentList.Header, url)
	//respBody2 := util.GetResponseBody(res2)

}