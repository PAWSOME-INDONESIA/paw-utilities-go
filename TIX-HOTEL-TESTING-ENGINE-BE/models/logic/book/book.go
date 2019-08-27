package book

import (
	"encoding/json"
	"strings"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/constant"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/models/logic/prebook"

	log "github.com/sirupsen/logrus"
)

type (
	// CommandBook ...
	CommandBook struct{}

	// Data ...
	Data struct {
		CartID        string      `json:"cartId"`
		ContactPerson interface{} `json:"contactPerson"`
		Traveller     interface{} `json:"traveller"`
	}

	// BookResponse ...
	BookResponse struct {
		Code    string       `json:"code"`
		Message string       `json:"Message"`
		Data    BookDataResp `json:"data"`
	}
	// BookDataResp ...
	BookDataResp struct {
		OrderID int64 `json:"orderId"`
	}
)

// Test ...
func (d *CommandBook) Test(contentList structs.ContentList) {
	var (
		err         error
		data        Data
		prebookResp = prebook.ResponsePreBook
	)

	// Get Data
	dataByte, _ := json.Marshal(contentList.Data)
	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		log.Warning("error unmarshal data :", err.Error())
		return
	}

	if data.CartID != "" {
		prebookResp.Data = nil
		prebookResp.Data = append(prebookResp.Data, prebook.ListPreBookResponse{
			CartId: data.CartID,
		})
	}

	log.Info("Test case :")
	log.Info("try book ", len(prebookResp.Data), " cartId")

	for _, value := range prebookResp.Data {
		var (
			bookResp BookResponse
		)
		data.CartID = value.CartId

		url := constant.BASEURL_CART + constant.URLCommand[contentList.Command]
		res := util.CallRest("POST", data, contentList.Header, url)
		statusCode := util.GetStatusCode(res)
		respBody := util.GetResponseBody(res)
		// log.Info(contentList.Header)
		// log.Info(respBody)
		log.Info("Test Case :")

		// 1. Check Status code
		if statusCode != 200 {
			log.Warning("1. Status Code must be 200 ", constant.SuccessMessage[false])
			log.Warning("\nResponse : ", respBody)
		} else {
			log.Info("1. Status Code must be 200 ", constant.SuccessMessage[true])
		}

		// 2. Check Response code
		err = json.Unmarshal([]byte(respBody), &bookResp)
		if err != nil {
			log.Warning("2. Check Response Code must success ", constant.SuccessMessage[false])
			log.Warning(err.Error())
			return
		}

		if strings.ToUpper(bookResp.Message) != constant.SUCCESS {
			log.Info("2. Check Response Code must success ", constant.SuccessMessage[false])
		} else {
			log.Info("2. Check Response Code must success ", constant.SuccessMessage[true])
		}

		// 3. Check Order ID is not zero
		if bookResp.Data.OrderID == 0 {
			log.Warning("3. Check Order ID (", (bookResp.Data.OrderID), ") is not zero", constant.SuccessMessage[false])
		} else {
			log.Info("3. Check Order ID (", (bookResp.Data.OrderID), ") is not zero ", constant.SuccessMessage[true])
		}

		// Test DB Level
		TestDB(data.CartID)
		log.Info("\n")

	}

}
