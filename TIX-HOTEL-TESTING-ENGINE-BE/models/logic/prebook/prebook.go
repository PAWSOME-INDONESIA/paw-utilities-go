package prebook

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/room"

	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var (
	ResponsePreBook ResponseData
)

type (
	CommandPrebook struct{}

	RequestRoom struct {
		Adult     int    `json:"adult"`
		HotelId   string `json:"hotelId"`
		Night     int    `json:"night"`
		RateCode  string `json:"rateCode"`
		Room      int    `json:"room"`
		StartDate string `json:"startDate"`
		RoomId    string `json:"roomId"`
	}

	Expected struct {
		TotalRequestHit int      `json:"totalRequestHit"`
		Vendor          []string `json:"vendor"`
		PaymentMethod   []int    `json:"paymentMethod"`
	}

	ResponseData struct {
		Data []ListPreBookResponse `json:"data"`
	}

	ListPreBookResponse struct {
		CartId string `json:"cartId"`
	}
)

func (d *CommandPrebook) Test(contentList structs.ContentList) {

	//Convert interface to struct
	var (
		expected         Expected
		results          ResponseData
		responseRoomTemp []room.ContentlList
		totalPrebook     int
	)
	dataE, _ := json.Marshal(contentList.Expected)
	json.Unmarshal(dataE, &expected)

	url := constant.BASEURL_CART + constant.URLCommand[contentList.Command]

	// Filter responseRoomTemp
	for _, value := range room.ResponseRoomTemp.Data.ContentList {
		isOk := true
		foundVendor := true
		foundPayment := true
		if len(expected.Vendor) > 0 {
			foundVendor = false
			for _, val := range expected.Vendor {
				// log.Info(value.ItemColor, "==", constant.Vendor[strings.ToUpper(val)])
				if value.ItemColor == constant.Vendor[strings.ToUpper(val)] {
					foundVendor = true
					break
				}
			}
		}
		isOk = foundVendor

		if isOk {
			if len(expected.PaymentMethod) > 0 {
				foundPayment = false
				for _, val := range expected.PaymentMethod {
					if value.PaymentOption == constant.PaymentMethod[val] {
						foundPayment = true
						break
					}
				}
			}
			isOk = foundPayment
		}

		if isOk {
			responseRoomTemp = append(responseRoomTemp, value)
		}
	}

	// log.Info(len(responseRoomTemp), expected.TotalRequestHit)
	totalPrebook = len(responseRoomTemp)
	if expected.TotalRequestHit > 0 && expected.TotalRequestHit < len(responseRoomTemp) {
		totalPrebook = expected.TotalRequestHit
		responseRoomTemp = responseRoomTemp[:totalPrebook]
	}

	log.Info("Test Case :")
	log.Info("1.try prebook " + strconv.Itoa(totalPrebook) + " rooms")
	if len(expected.Vendor) > 0 {
		log.Info("Vendor : " + strings.Join(expected.Vendor, ""))
	}
	if len(expected.PaymentMethod) > 0 {
		log.Info("Payment Method : " + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(expected.PaymentMethod)), ""), "[]"))
	}

	for _, v := range responseRoomTemp {

		requestRoom := RequestRoom{
			Adult:     room.RequestRoomTemp.Adult,
			HotelId:   room.RequestRoomTemp.HotelID,
			Night:     room.RequestRoomTemp.Night,
			RateCode:  v.RateCode,
			Room:      room.RequestRoomTemp.Room,
			StartDate: room.RequestRoomTemp.StartDate,
			RoomId:    v.RoomID,
		}

		res := util.CallRest("POST", &requestRoom, contentList.Header, url)
		statusCode := util.GetStatusCode(res)
		respBody := util.GetResponseBody(res)
		if statusCode == 200 {
			log.Info("--- Status Code " + constant.SuccessMessage[true])

			responseData := gjson.Get(respBody, "data")
			cartId := gjson.Get(responseData.String(), "cartId").String()
			if cartId != "" {

				results.Data = append(results.Data, ListPreBookResponse{CartId: cartId})
				log.Info("--- Response Data " + constant.SuccessMessage[true])
				TestDB(cartId)
			}

		}

	}

	ResponsePreBook = results

}
