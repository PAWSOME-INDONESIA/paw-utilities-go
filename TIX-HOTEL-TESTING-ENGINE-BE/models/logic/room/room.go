package room

import (
	"TIX-HOTEL-TESTING-ENGINE-BE/util"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/constant"
	"TIX-HOTEL-TESTING-ENGINE-BE/util/structs"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	RequestRoomTemp  Data
	ResponseRoomTemp RoomSearchResponse
)

type (
	// CommandRoom ...
	CommandRoom struct{}

	// Data ...
	Data struct {
		Adult     int    `json:"adult"`
		Night     int    `json:"night"`
		Room      int    `json:"room"`
		HotelID   string `json:"hotelId"`
		StartDate string `json:"startDate"`
	}

	// Expected ...
	Expected struct {
		ResponseData interface{} `json"responseData"`
	}

	// RoomSearchResponse ...
	RoomSearchResponse struct {
		Data    Contents `json:"data"`
		Message string   `json:"message"`
	}

	// Contents ...
	Contents struct {
		ContentList []ContentlList `json:"roomList"`
	}

	// ContentlList ...
	ContentlList struct {
		RoomID        string   `json:"roomId"`
		RoomName      string   `json:"roomName"`
		MaxOccupancy  int      `json:"maxOccupancy"`
		AvailRoom     int      `json:"availableRoom"`
		RateInfo      RateInfo `json:"rateInfo"`
		RateCode      string   `json:"rateCode"`
		ItemColor     string   `json:"itemColor"`
		PaymentOption string   `json:"paymentOption"`
	}

	// RateInfo ...
	RateInfo struct {
		Price        Price        `json:"price"`
		PriceSummary PriceSummary `json:"priceSummary"`
	}

	// Price ...
	Price struct {
		RateWithTax float64 `json:"rateWithTax"`
	}

	// PriceSummary ...
	PriceSummary struct {
		Total float64 `json:"total"`
	}
)

// Test ...
func (d *CommandRoom) Test(contentList structs.ContentList) {
	var (
		err              error
		data             Data
		roomSearchResp   RoomSearchResponse
		expected         Expected
		expectedResponse []ContentlList
		checkOccupancy   = true
		checkAvailRoom   = true
		checkPrice       = true
		logMaxOccupancy  = make([]string, 0)
		logAvailRoom     = make([]string, 0)
		logPrice         = make([]string, 0)
		price            float64
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
	RequestRoomTemp = data

	url := constant.BASEURL + constant.URLCommand[contentList.Command]
	res := util.CallRest("POST", data, contentList.Header, url)
	statusCode := util.GetStatusCode(res)
	respBody := util.GetResponseBody(res)

	log.Info("Test Case :")

	// 1. Check Status code
	if statusCode != 200 {
		log.Warning("1. Status Code must be 200 ", constant.SuccessMessage[false])
		log.Warning("\nResponse : ", respBody)
		return
	}
	log.Info("1. Status Code must be 200", constant.SuccessMessage[true])

	// 2. Check room list must have data (>0)
	err = json.Unmarshal([]byte(respBody), &roomSearchResp)
	if err != nil {
		log.Warning("2. Check room list must > 0", constant.SuccessMessage[false])
		log.Warning(err.Error())
		log.Warning("\nResponse : ", respBody)
		return
	}

	if len(roomSearchResp.Data.ContentList) == 0 {
		log.Warning("2. Check room list must > 0", constant.SuccessMessage[false])
		log.Warning("\nResponse : ", respBody)
	} else {
		log.Info("2. Check room list must > 0", constant.SuccessMessage[true])
	}

	// 3. Check Response Message Success
	if strings.ToUpper(roomSearchResp.Message) != constant.SUCCESS {
		log.Warning("3. Check Response Code must success ", constant.SuccessMessage[false])
	} else {
		log.Info("3. Check Response Code must success ", constant.SuccessMessage[true])
	}

	// Check
	// 4. Check Max Occupancy must less or equal than adult ( <= req.adult)
	// 5. Check Available Room muss higher or equal than room (>= req.room)
	// 6. Check Price must be sorted

	for n, value := range roomSearchResp.Data.ContentList {
		// Check Match occupancy
		if value.MaxOccupancy > data.Adult {
			checkOccupancy = false
			logMaxOccupancy = append(logMaxOccupancy, "\nRoom ID ", value.RoomID, string(value.MaxOccupancy))
		}

		// Check avail room
		if value.AvailRoom < data.Room {
			checkAvailRoom = false
			logAvailRoom = append(logAvailRoom, "\nRoom ID ", value.RoomID, string(value.AvailRoom))
		}

		// Check price must sort
		if n == 0 {
			price = value.RateInfo.Price.RateWithTax
		} else {
			if price > value.RateInfo.Price.RateWithTax {
				checkPrice = false
				logPrice = append(logPrice, "\nRoom ID ", value.RoomID, fmt.Sprintf("%f",
					value.RateInfo.Price.RateWithTax))
			}
		}

		// Check response
		for keyexpect, valueExpect := range expectedResponse {

			if util.CheckSimilarStruct(valueExpect, value) {
				checkExpected[keyexpect] = true
			}
		}
	}

	log.Info("4. Check Max Occupancy must <= (", data.Adult, ")",
		constant.SuccessMessage[checkOccupancy])
	if !checkOccupancy {
		log.Warning("Log Message :")
		log.Warning(logMaxOccupancy)
	}

	log.Info("5. Check Available Room should >= (", data.Room, ")",
		constant.SuccessMessage[checkAvailRoom])
	if !checkAvailRoom {
		log.Warning("Log Message :")
		log.Warning(logAvailRoom)
	}

	log.Info("6. Check Room Price must be sorted",
		constant.SuccessMessage[checkPrice])
	if !checkPrice {
		log.Warning("Log Message :")
		log.Warning(logPrice)
	}

	if len(expectedResponse) > 0 {
		log.Info("7. Check Expected response:")
		for keyexpect := range expectedResponse {
			log.Info("Expected response ", (keyexpect + 1), constant.SuccessMessage[checkExpected[keyexpect]])
		}
	}

	ResponseRoomTemp = roomSearchResp

}
