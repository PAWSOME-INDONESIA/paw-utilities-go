package searchseo

import (
	"encoding/json"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"strings"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util"

	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

type (
	// CommandSearchSEO ...
	CommandSearchSEO struct{}

	// Data ...
	Data struct {
		Filter      Filter `json:"filter,omitempty"`
		Page        int    `json:"page"`
		SearchType  string `json:"searchType"`
		SearchValue string `json:"searchValue"`
	}

	// Filter ...
	Filter struct {
		Type string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
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
	}

	// IDName ...
	IDName struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

// Test ...
func (d *CommandSearchSEO) Test(contentList structs.ContentList) {
	var (
		err                      error
		data                     Data
		hotelSearchResp          HotelSearchResponse
		expected                 Expected
		expectedResponse         []ContentlList
		checkMatchSearchType     = true
		checkNameEmpty           = true
		logNameEmpty             = make([]string, 0)
		logMatchSearchType   	 = make([]string, 0)
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
	} else {
		log.Info("2. Check hotel list must > 0", constant.SuccessMessage[true])
	}

	var locationDetail []string
	result2 := gjson.Get(respBody, "data.searchDetail.searchLocation").Array()
	for _, v := range result2 {
		locationDetail = append(locationDetail, v.String())
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
	}

	log.Info("4. Check hotel list should match search type (", data.SearchType, ")",
		constant.SuccessMessage[checkMatchSearchType])
	if !checkMatchSearchType {
		log.Warning("Log Message :")
		log.Warning(logMatchSearchType)
	}

	if len(expectedResponse) > 0 {
		log.Info("5. Check Expected response:")
		for keyexpect := range expectedResponse {
			log.Info("Expected response ", (keyexpect + 1), constant.SuccessMessage[checkExpected[keyexpect]])
		}
	}

	log.Info("6. Check Hotel Name is not empty",
		constant.SuccessMessage[checkNameEmpty])
	if !checkNameEmpty {
		log.Warning("Log Message :")
		log.Warning(logNameEmpty)
	}
}
