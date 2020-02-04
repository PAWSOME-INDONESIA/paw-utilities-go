package dto

type (
	RoomResponseDto struct {
		RateCode                   string                            `json:"rateCode"`
		RoomID                     string                            `json:"roomId"`
		RateKey                    string                            `json:"rateKey"`
		CurrentAllotment           int                               `json:"currentAllotment"`
		RateOccupancyPerRoom       int                               `json:"rateOccupancyPerRoom"`
		SupplierType               string                            `json:"supplierType"`
		BedTypes                   []HotelRoomBedType                `json:"bedTypes"`
		SmokingPreferences         []IDNameData                      `json:"smokingPreferences"`
		CancellationPolicy         string                            `json:"cancellationPolicy"`
		CancellationPolicies       []HotelRoomCancellationPolicies   `json:"cancellationPolicies"`
		CancellationPoliciesV2     []HotelRoomCancellationPoliciesV2 `json:"cancellationPoliciesV2"`
		CancellationPolicyInfo     []HotelRoomCancellationPolicyInfo `json:"cancellationPolicyInfo"`
		CancellationDetails        []HotelRoomCancellationDetails    `json:"cancellationDetails"`
		RoomImages                 interface{}                       `json:"roomImages"`
		PayUponArrival             []HotelRoomPayUponArrival         `json:"payUponArrival"`
		RateInfo                   HotelRoomRateInfo                 `json:"rateInfo"`
		PromoInfo                  HotelRoomPromoInfo                `json:"promoInfo"`
		FreeWifi                   bool                              `json:"freeWifi"`
		FreeWifiDesc               string                            `json:"freeWifiDesc"`
		FreeBreakfast              bool                              `json:"freeBreakfast"`
		BreakfastPax               interface{}                       `json:"breakfastPax"`
		FreeBreakfastDesc          string                            `json:"freeBreakfastDesc"`
		ValueAdds                  []string                          `json:"valueAdds"`
		SoldOut                    bool                              `json:"soldOut"`
		CheckInInstructions        string                            `json:"checkInInstructions"`
		SpecialCheckInInstructions interface{}                       `json:"specialCheckInInstructions"`
		AdditionalInfo             interface{}                       `json:"additionalInfo"`
		PaymentType                string                            `json:"paymentType,omitempty"`
		PaymentOption              string                            `json:"paymentOption,omitempty"`
	}

	HotelRoomBedType struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	}

	HotelRoomCancellationPolicies struct {
		DaysBefore int     `json:"daysBefore"`
		Amount     float64 `json:"amount"`
	}

	HotelRoomCancellationPoliciesV2 struct {
		Time   string  `json:"time"`
		Amount float64 `json:"amount"`
	}

	HotelRoomCancellationPolicyInfo struct {
		VersionID           int      `json:"versionId"`
		CancelTime          string   `json:"cancelTime"`
		StartWindowHours    int      `json:"startWindowHours"`
		NightCount          int      `json:"nightCount"`
		Amount              *float64 `json:"amount"`
		Percent             float64  `json:"percent"`
		CurrencyCode        string   `json:"currencyCode"`
		TimeZoneDescription string   `json:"timeZoneDescription"`
	}

	HotelRoomCancellationDetails struct {
		Days                 int     `json:"days"`
		ChargeAmount         float64 `json:"chargeAmount"`
		RefundCustomerAmount float64 `json:"refundCustomerAmount"`
		RefundHotelAmount    float64 `json:"refundHotelAmount"`
	}

	HotelRoomPayUponArrival struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Frequency   string  `json:"frequency"`
		Unit        string  `json:"unit"`
	}

	HotelRoomRateInfo struct {
		Currency     string                        `json:"currency"`
		Refundable   bool                          `json:"refundable"`
		Price        HotelRoomRateInfoPrice        `json:"price"`
		TixPoint     interface{}                   `json:"tixPoint"`
		PriceSummary HotelRoomRateInfoPriceSummary `json:"priceSummary"`
		//MarketingFee  float64                            `json:"vendorIncentive"`
		//PaymentOption string                             `json:"paymentType"`
	}

	HotelRoomRateInfoPrice struct {
		BaseRateWithTax      float64 `json:"baseRateWithTax"`
		RateWithTax          float64 `json:"rateWithTax"`
		TotalBaseRateWithTax float64 `json:"totalBaseRateWithTax"`
		TotalRateWithTax     float64 `json:"totalRateWithTax"`
	}

	HotelRoomRateInfoPriceSummary struct {
		Total            float64                           `json:"total"`
		TotalWithoutTax  float64                           `json:"totalWithoutTax"`
		TaxAndOtherFee   float64                           `json:"taxAndOtherFee"`
		TotalCompulsory  float64                           `json:"totalCompulsory"`
		Net              float64                           `json:"net"`
		MarkupPercentage float64                           `json:"markupPercentage"`
		Surcharge        []HotelRoomRateInfoPriceSurcharge `json:"surcharge"`
		Compulsory       interface{}                       `json:"compulsory"`
		PricePerNight    []HotelRoomRateInfoPricePerNight  `json:"pricePerNight"`
		TotalObject      HotelRoomRateInfoPriceTotalObject `json:"totalObject"`
		SubsidyPrice     float64                           `json:"subsidyPrice"`
		MarketingFee     float64                           `json:"vendorIncentive"`
		MarkupID         []string                          `json:"markupId"`
		SubsidyID        []string                          `json:"subsidyId"`
	}

	HotelRoomRateInfoPriceSurcharge struct {
		Type       string  `json:"type"`
		Name       string  `json:"name"`
		Rate       float64 `json:"rate"`
		ChargeUnit string  `json:"chargeUnit,omitempty"`
		Code       string  `json:"code,omitempty"`
	}

	HotelRoomRateInfoPricePerNight struct {
		StayingDate string  `json:"stayingDate"`
		Rate        float64 `json:"rate"`
	}

	HotelRoomRateInfoPriceTotalObject struct {
		Label string  `json:"label"`
		Value float64 `json:"value"`
	}

	HotelRoomPromoInfo struct {
		PromoText    string      `json:"promoText"`
		PackageDeals bool        `json:"packageDeals"`
		MemberDeals  bool        `json:"memberDeals"`
		TonightDeals bool        `json:"tonightDeals"`
		PromoIcon    interface{} `json:"promoIcon"`
		ExpiredDate  *int64      `json:"expiredDate"`
		PromoLabelID string      `json:"promoLabelId"`
	}

	AvailabilityResponseDto struct {
		HotelID            string            `json:"hotelId"`
		HotelAnnouncement  interface{}       `json:"hotelAnnouncement"`
		HotelRoom          []RoomResponseDto `json:"hotelRoom"`
		ItineraryID        *string           `json:"itineraryId"`
		BookingToken       *string           `json:"bookingToken"`
		BookingExpiredTime int64             `json:"bookingExpiredTime"`
	}
)
