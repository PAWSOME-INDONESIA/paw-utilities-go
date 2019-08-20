package search

import (
	"TIX-HOTEL-TESTING-ENGINE-BE/util"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/constant"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	// CommandSearch ...
	CommandSearch struct{}

	// Data ...
	Data struct {
		Adult       int    `json:"adult"`
		Night       int    `json:"night"`
		Page        int    `json:"page"`
		Priority    string `json:"priorityRankingType"`
		Room        int    `json:"room"`
		SearchType  string `json:"searchType"`
		SearchValue string `json:"searchValue"`
		Sort        string `json:"sort"`
		StartDate   string `json:"startDate"`
	}

	// Expected ...
	Expected struct {
		ResponseData interface{} `json"responseData"`
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
		PaymentOption  interface{} `json:"paymentOption"`
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
		err                  error
		data                 Data
		hotelSearchResp      HotelSearchResponse
		expected             Expected
		expectedResponse     []ContentlList
		checkMatchSearchType = true
		checkAvailRoom       = true
		logMatchSearchType   = make([]string, 0)
		logAvailRoom         = make([]string, 0)
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

	if len(hotelSearchResp.Data.ContentList) == 0 {
		log.Warning("2. Check hotel list must > 0", constant.SuccessMessage[false])
		log.Warning("\nResponse : ", hotelSearchResp)
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

	// Test Database Level
	// TestDB(data.SearchValue)
}
