package dto

import (
	"github.com/labstack/echo/v4"
)

type (
	AppContext struct {
		echo.Context
		MandatoryRequestDto
	}

	ErrorCode struct {
		Code                 string
		Message              string
		WrappedError         error
		FrontEndErrorMessage string
	}

	BaseResponseDto struct {
		Code       string      `json:"code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		Errors     []string    `json:"errors"`
		ServerTime int64       `json:"serverTime"`
	}

	MandatoryRequestDto struct {
		StoreID           string `json:"storeId" validate:"required"`
		ChannelID         string `json:"channelId" validate:"required"`
		RequestID         string `json:"requestId" validate:"required"`
		ServiceID         string `json:"serviceId" validate:"required"`
		Username          string `json:"username" validate:"required"`
		Language          string `json:"lang,omitempty"`
		Login             int    `json:"login,omitempty"`
		CustomerUserAgent string `json:"customerUserAgent,omitempty"`
		CustomerIPAddress string `json:"customerIpAddress,omitempty"`
		CustomerSessionId string `json:"customerSessionId,omitempty"`
		Currency          string `json:"currency,omitempty"`
		ResellerID        string `json:"resellerId,omitempty"`
		PreviousOrderID   string `json:"prevOrderId,omitempty"`
	}
)

func (e *ErrorCode) Error() string {
	return e.Message
}
