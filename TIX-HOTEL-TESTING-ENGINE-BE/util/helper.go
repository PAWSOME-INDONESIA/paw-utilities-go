package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/TIX-HOTEL-TESTING-ENGINE-BE/util/structs"

	log "github.com/sirupsen/logrus"
)

const (
	BASE_URL = "http://54.169.65.147:7060/tix-hotel-autocomplete"
)

// Call Rest Api
func CallRest(method string, bodyContent, headerContent interface{}, url string) (result structs.ResponseDefault) {

	param, _ := json.Marshal(bodyContent)

	var jsonStr = []byte(param)

	log.Info("Send Req to :", url)
	log.Info("Request body :", string(jsonStr))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-Forwarded-Host", fmt.Sprintf("%s",  GetOutboundIP()))
	for k, v := range headerContent.(map[string]interface{}) {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	result = structs.ResponseDefault{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}

	return result
}

// Get Status Code
func GetStatusCode(resp structs.ResponseDefault) int {
	return resp.StatusCode
}

// Get Response Body
func GetResponseBody(resp structs.ResponseDefault) string {
	return resp.Body
}

// CheckSimilarStruct ...
func CheckSimilarStruct(x interface{}, y interface{}) bool {
	xvalue := reflect.ValueOf(x)
	// xtype := reflect.TypeOf(x)

	yvalue := reflect.ValueOf(y)

	for i, n := 0, xvalue.NumField(); i < n; i++ {
		// log.Info(xtype.Field(i).Name, xvalue.Field(i))

		switch typee := xvalue.Field(i).Kind(); typee {
		case reflect.Struct:
			{
				if !checkStructOrSlice(xvalue.Field(i), yvalue.Field(i), reflect.Struct) {
					return false
				}
				// if !reflect.DeepEqual(xvalue.Field(i).Interface(), reflect.Zero(xvalue.Field(i).Type()).Interface()) {
				// 	if !reflect.DeepEqual(xvalue.Field(i).Interface(), yvalue.Field(i).Interface()) {
				// 		return false
				// 	}
				// }
			}
		case reflect.Slice:
			{
				// if !checkStructOrSlice(xvalue.Field(i), yvalue.Field(i), reflect.Slice) {
				// 	return false
				// }
				if !reflect.DeepEqual(xvalue.Field(i).Interface(), reflect.Zero(xvalue.Field(i).Type()).Interface()) {
					if !reflect.DeepEqual(xvalue.Field(i).Interface(), yvalue.Field(i).Interface()) {
						return false
					}
				}
			}
		default:
			{
				if !checkValueIgnoreNil(xvalue.Field(i), yvalue.Field(i)) {
					return false
				}
			}
		}

	}
	return true
}

func checkStructOrSlice(xvalue reflect.Value, yvalue reflect.Value, opt reflect.Kind) bool {
	if opt == reflect.Struct {
		for i, n := 0, xvalue.NumField(); i < n; i++ {
			if !checkStructOrSlice(xvalue.Field(i), yvalue.Field(i), xvalue.Field(i).Kind()) {
				return false
			}
		}
	} else {
		if !checkValueIgnoreNil(xvalue, yvalue) {
			return false
		}
	}
	return true
}

func checkValueIgnoreNil(xvalue reflect.Value, yvalue reflect.Value) bool {
	if !reflect.DeepEqual(xvalue.Interface(), reflect.Zero(xvalue.Type()).Interface()) {
		// log.Info(xvalue, yvalue, true)
		if !reflect.DeepEqual(xvalue.Interface(), yvalue.Interface()) {
			return false
		}
	}

	return true
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	//fmt.Println(localAddr.)

	return localAddr.IP
}

func StringContaintsInSlice(str string, list []string) bool {
	for _, v := range list {

		if strings.Contains(strings.ToLower(v), strings.ToLower(str)) {
			return true
		}
	}
	return false
}
