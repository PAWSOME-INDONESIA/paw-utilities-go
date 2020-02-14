package dto

type (
	IntegratorB2bRequestDto struct {
		Vendor                   string                         `json:"vendor"`
		Code                     *string                        `json:"code"`
		Message                  *string                        `json:"message"`
		Error                    *string                        `json:"error"`
		MandatoryRequest         MandatoryRequestDto            `json:"mandatoryRequest"`
		HotelAvailabilityRequest HotelAvailabilityB2bRequestDto `json:"hotelAvailabilityRequest"`
	}

	IntegratorB2bResponseDto struct {
		Vendor                    string                          `json:"vendor"`
		Code                      *string                         `json:"code"`
		Message                   *string                         `json:"message"`
		Error                     *string                         `json:"error"`
		MandatoryRequest          MandatoryRequestDto             `json:"mandatoryRequest"`
		HotelAvailabilityRequest  HotelAvailabilityB2bRequestDto  `json:"hotelAvailabilityRequest"`
		HotelAvailabilityResponse HotelAvailabilityB2bResponseDto `json:"hotelAvailabilityResponse"`
	}

	HotelAvailabilityB2bRequestDto struct {
		HotelID         []string          `json:"hotelId" validate:"required"`
		MapVendorCoreID map[string]string `json:"mapVendorCoreID" validate:"required"`
		StartDate       string            `json:"startDate" validate:"required"`
		EndDate         string            `json:"endDate" validate:"required"`
		NumberOfNight   int               `json:"numberOfNights" validate:"required"`
		NumberOfRooms   int               `json:"numberOfRoom" validate:"required"`
		NumberOfAdult   int               `json:"numberOfAdults" validate:"required"`
		NumberOfChild   int               `json:"numberOfChild"`
		ChildAges       string            `json:"childAges"`
		PackageRate     int               `json:"packageRate"`
	}

	HotelAvailabilityB2bResponseDto []AvailabilityB2bResponseDto

	AvailabilityB2bResponseDto struct {
		HotelID   string               `json:"hotelId"`
		HotelRoom []RoomB2bResponseDto `json:"hotelRoom"`
	}

	RoomB2bResponseDto struct {
		RateCode                   string                            `json:"rateCode"`
		RoomID                     string                             `json:"roomId"`
		RateKey                    string                               `json:"rateKey"`
		CurrentAllotment           int                                  `json:"currentAllotment"`
		RateOccupancyPerRoom       int                                  `json:"rateOccupancyPerRoom"`
		SupplierType               string                               `json:"supplierType"`
		CancellationPolicy         string                               `json:"cancellationPolicy"`
		CancellationPolicies       []HotelB2bRoomCancellationPolicies   `json:"cancellationPolicies"`
		CancellationPoliciesV2     []HotelB2bRoomCancellationPoliciesV2 `json:"cancellationPoliciesV2"`
		CancellationPolicyInfo     []HotelB2bRoomCancellationPolicyInfo `json:"cancellationPolicyInfo"`
		CancellationDetails        []HotelB2bRoomCancellationDetails    `json:"cancellationDetails"`
		RateInfo                   HotelB2bRoomRateInfo                 `json:"rateInfo"`
		FreeBreakfast              bool                                 `json:"freeBreakfast"`
		BreakfastPax               interface{}                          `json:"breakfastPax"`
		FreeBreakfastDesc          string                               `json:"freeBreakfastDesc"`
		SoldOut                    bool                                 `json:"soldOut"`
		CheckInInstructions        string                               `json:"checkInInstructions"`
		SpecialCheckInInstructions interface{}                          `json:"specialCheckInInstructions"`
		PaymentOption              string                               `json:"paymentOption,omitempty"`
		CrossSellRate              bool                                 `json:"crossSellRate"`
	}

	HotelB2bRoomCancellationPolicies struct {
		DaysBefore int     `json:"daysBefore"`
		Amount     float64 `json:"amount"`
	}

	HotelB2bRoomCancellationPoliciesV2 struct {
		Time   string  `json:"time"`
		Amount float64 `json:"amount"`
	}

	HotelB2bRoomCancellationPolicyInfo struct {
		VersionID           int      `json:"versionId"`
		CancelTime          string   `json:"cancelTime"`
		StartWindowHours    int      `json:"startWindowHours"`
		NightCount          int      `json:"nightCount"`
		Amount              *float64 `json:"amount"`
		Percent             float64  `json:"percent"`
		CurrencyCode        string   `json:"currencyCode"`
		TimeZoneDescription string   `json:"timeZoneDescription"`
	}

	HotelB2bRoomCancellationDetails struct {
		Days                 int     `json:"days"`
		ChargeAmount         float64 `json:"chargeAmount"`
		RefundCustomerAmount float64 `json:"refundCustomerAmount"`
		RefundHotelAmount    float64 `json:"refundHotelAmount"`
	}

	HotelB2bRoomRateInfo struct {
		Currency     string                           `json:"currency"`
		Refundable   bool                             `json:"refundable"`
		Price        HotelB2bRoomRateInfoPrice        `json:"price"`
		TixPoint     interface{}                      `json:"tixPoint"`
		PriceSummary HotelB2bRoomRateInfoPriceSummary `json:"priceSummary"`
	}

	HotelB2bRoomRateInfoPrice struct {
		BaseRateWithTax      float64 `json:"baseRateWithTax"`
		RateWithTax          float64 `json:"rateWithTax"`
		TotalBaseRateWithTax float64 `json:"totalBaseRateWithTax"`
		TotalRateWithTax     float64 `json:"totalRateWithTax"`
	}

	HotelB2bRoomRateInfoPriceSummary struct {
		Total                 float64                           `json:"total"`
		TotalWithoutTax       float64                              `json:"totalWithoutTax"`
		TaxAndOtherFee        float64                              `json:"taxAndOtherFee"`
		TotalCompulsory       float64                              `json:"totalCompulsory"`
		TotalSellingRateAddOn float64                              `json:"totalSellingRateAddOn"`
		Net                   float64                              `json:"net"`
		MarkupPercentage      float64                              `json:"markupPercentage"`
		Surcharge             []HotelB2bRoomRateInfoPriceSurcharge `json:"surcharge"`
		Compulsory            interface{}                          `json:"compulsory"`
		PricePerNight         []HotelB2bRoomRateInfoPricePerNight  `json:"pricePerNight"`
		TotalObject           HotelB2bRoomRateInfoPriceTotalObject `json:"totalObject"`
		SubsidyPrice          float64                              `json:"subsidyPrice"`
		MarkupID              []string                             `json:"markupId"`
		SubsidyID             []string                             `json:"subsidyId"`
		VendorIncentive       float64                              `json:"vendorIncentive"`
	}

	HotelB2bRoomRateInfoPriceSurcharge struct {
		Type       string  `json:"type"`
		Name       string  `json:"name"`
		Rate       float64 `json:"rate"`
		ChargeUnit string  `json:"chargeUnit,omitempty"`
		Code       string  `json:"code,omitempty"`
	}

	HotelB2bRoomRateInfoPricePerNight struct {
		StayingDate string  `json:"stayingDate"`
		Rate        float64 `json:"rate"`
	}

	HotelB2bRoomRateInfoPriceTotalObject struct {
		Label string  `json:"label"`
		Value float64 `json:"value"`
	}
)
