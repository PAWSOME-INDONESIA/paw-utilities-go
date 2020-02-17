package dto

type (
	IntegratorB2BRequestDto struct {
		Vendor                   string                         `json:"vendor"`
		Code                     *string                        `json:"code"`
		Message                  *string                        `json:"message"`
		Error                    *string                        `json:"error"`
		MandatoryRequest         MandatoryRequestDto            `json:"mandatoryRequest"`
		HotelAvailabilityRequest HotelAvailabilityB2BRequestDto `json:"hotelAvailabilityRequest"`
	}

	IntegratorB2BResponseDto struct {
		Vendor                    string                          `json:"vendor"`
		Code                      *string                         `json:"code"`
		Message                   *string                         `json:"message"`
		Error                     *string                         `json:"error"`
		MandatoryRequest          MandatoryRequestDto             `json:"mandatoryRequest"`
		HotelAvailabilityRequest  HotelAvailabilityB2BRequestDto  `json:"hotelAvailabilityRequest"`
		HotelAvailabilityResponse HotelAvailabilityB2BResponseDto `json:"hotelAvailabilityResponse"`
	}

	HotelAvailabilityB2BRequestDto struct {
		HotelIDs        []string          `json:"hotelIds"`
		MapVendorCoreID map[string]string `json:"mapVendorCoreId"`
		StartDate       string            `json:"startDate" validate:"required"`
		EndDate         string            `json:"endDate" validate:"required"`
		NumberOfNight   int               `json:"numberOfNights" validate:"required"`
		NumberOfRooms   int               `json:"numberOfRoom" validate:"required"`
		NumberOfAdult   int               `json:"numberOfAdults" validate:"required"`
		NumberOfChild   int               `json:"numberOfChild"`
		ChildAges       string            `json:"childAges"`
		PackageRate     int               `json:"packageRate"`
	}

	HotelAvailabilityB2BResponseDto []AvailabilityB2BResponseDto

	AvailabilityB2BResponseDto struct {
		HotelID   string               `json:"hotelId"`
		HotelRoom []RoomB2BResponseDto `json:"hotelRoom"`
	}

	RoomB2BResponseDto struct {
		RateCode                   string                               `json:"rateCode"`
		RoomID                     string                               `json:"roomId"`
		RateKey                    string                               `json:"rateKey"`
		CurrentAllotment           int                                  `json:"currentAllotment"`
		RateOccupancyPerRoom       int                                  `json:"rateOccupancyPerRoom"`
		SupplierType               string                               `json:"supplierType"`
		CancellationPolicy         string                               `json:"cancellationPolicy"`
		CancellationPolicies       []HotelB2BRoomCancellationPolicies   `json:"cancellationPolicies"`
		CancellationPoliciesV2     []HotelB2BRoomCancellationPoliciesV2 `json:"cancellationPoliciesV2"`
		CancellationPolicyInfo     []HotelB2BRoomCancellationPolicyInfo `json:"cancellationPolicyInfo"`
		CancellationDetails        []HotelB2BRoomCancellationDetails    `json:"cancellationDetails"`
		RateInfo                   HotelB2BRoomRateInfo                 `json:"rateInfo"`
		FreeBreakfast              bool                                 `json:"freeBreakfast"`
		BreakfastPax               interface{}                          `json:"breakfastPax"`
		FreeBreakfastDesc          string                               `json:"freeBreakfastDesc"`
		SoldOut                    bool                                 `json:"soldOut"`
		CheckInInstructions        string                               `json:"checkInInstructions"`
		SpecialCheckInInstructions interface{}                          `json:"specialCheckInInstructions"`
		PaymentOption              string                               `json:"paymentOption,omitempty"`
		CrossSellRate              bool                                 `json:"crossSellRate"`
	}

	HotelB2BRoomCancellationPolicies struct {
		DaysBefore int     `json:"daysBefore"`
		Amount     float64 `json:"amount"`
	}

	HotelB2BRoomCancellationPoliciesV2 struct {
		Time   string  `json:"time"`
		Amount float64 `json:"amount"`
	}

	HotelB2BRoomCancellationPolicyInfo struct {
		VersionID           int      `json:"versionId"`
		CancelTime          string   `json:"cancelTime"`
		StartWindowHours    int      `json:"startWindowHours"`
		NightCount          int      `json:"nightCount"`
		Amount              *float64 `json:"amount"`
		Percent             float64  `json:"percent"`
		CurrencyCode        string   `json:"currencyCode"`
		TimeZoneDescription string   `json:"timeZoneDescription"`
	}

	HotelB2BRoomCancellationDetails struct {
		Days                 int     `json:"days"`
		ChargeAmount         float64 `json:"chargeAmount"`
		RefundCustomerAmount float64 `json:"refundCustomerAmount"`
		RefundHotelAmount    float64 `json:"refundHotelAmount"`
	}

	HotelB2BRoomRateInfo struct {
		Currency     string                           `json:"currency"`
		Refundable   bool                             `json:"refundable"`
		Price        HotelB2BRoomRateInfoPrice        `json:"price"`
		TixPoint     interface{}                      `json:"tixPoint"`
		PriceSummary HotelB2BRoomRateInfoPriceSummary `json:"priceSummary"`
	}

	HotelB2BRoomRateInfoPrice struct {
		BaseRateWithTax      float64 `json:"baseRateWithTax"`
		RateWithTax          float64 `json:"rateWithTax"`
		TotalBaseRateWithTax float64 `json:"totalBaseRateWithTax"`
		TotalRateWithTax     float64 `json:"totalRateWithTax"`
	}

	HotelB2BRoomRateInfoPriceSummary struct {
		Total                 float64                              `json:"total"`
		TotalWithoutTax       float64                              `json:"totalWithoutTax"`
		TaxAndOtherFee        float64                              `json:"taxAndOtherFee"`
		TotalCompulsory       float64                              `json:"totalCompulsory"`
		TotalSellingRateAddOn float64                              `json:"totalSellingRateAddOn"`
		Net                   float64                              `json:"net"`
		MarkupPercentage      float64                              `json:"markupPercentage"`
		Surcharge             []HotelB2BRoomRateInfoPriceSurcharge `json:"surcharge"`
		Compulsory            interface{}                          `json:"compulsory"`
		PricePerNight         []HotelB2BRoomRateInfoPricePerNight  `json:"pricePerNight"`
		TotalObject           HotelB2BRoomRateInfoPriceTotalObject `json:"totalObject"`
		SubsidyPrice          float64                              `json:"subsidyPrice"`
		MarkupID              []string                             `json:"markupId"`
		SubsidyID             []string                             `json:"subsidyId"`
		VendorIncentive       float64                              `json:"vendorIncentive"`
	}

	HotelB2BRoomRateInfoPriceSurcharge struct {
		Type       string  `json:"type"`
		Name       string  `json:"name"`
		Rate       float64 `json:"rate"`
		ChargeUnit string  `json:"chargeUnit,omitempty"`
		Code       string  `json:"code,omitempty"`
	}

	HotelB2BRoomRateInfoPricePerNight struct {
		StayingDate string  `json:"stayingDate"`
		Rate        float64 `json:"rate"`
	}

	HotelB2BRoomRateInfoPriceTotalObject struct {
		Label string  `json:"label"`
		Value float64 `json:"value"`
	}
)
