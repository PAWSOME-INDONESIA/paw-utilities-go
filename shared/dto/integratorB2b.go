package dto

type (
	IntegratorB2bRequestDto struct {
		Vendor                   string                         `json:"vendor"`
		Code                     *string                        `json:"code"`
		Message                  *string                        `json:"message"`
		Error                    *string                        `json:"error"`
		MandatoryRequest         MandatoryRequestDto            `json:"mandatoryRequest"`
		HotelAvailabilityRequest HotelB2bAvailabilityRequestDto `json:"hotelAvailabilityRequest"`
	}

	IntegratorB2bResponseDto struct {
		Vendor                    string                          `json:"vendor"`
		Code                      *string                         `json:"code"`
		Message                   *string                         `json:"message"`
		Error                     *string                         `json:"error"`
		MandatoryRequest          MandatoryRequestDto             `json:"mandatoryRequest"`
		HotelAvailabilityRequest  HotelB2bAvailabilityRequestDto  `json:"hotelAvailabilityRequest"`
		HotelAvailabilityResponse HotelB2bAvailabilityResponseDto `json:"hotelAvailabilityResponse"`
	}

	HotelB2bAvailabilityRequestDto struct {
		HotelID       []string `json:"hotelId" validate:"required"`
		StartDate     string   `json:"startDate" validate:"required"`
		EndDate       string   `json:"endDate" validate:"required"`
		NumberOfNight int      `json:"numberOfNights" validate:"required"`
		NumberOfRooms int      `json:"numberOfRoom" validate:"required"`
		NumberOfAdult int      `json:"numberOfAdults" validate:"required"`
		NumberOfChild int      `json:"numberOfChild"`
		ChildAges     string   `json:"childAges"`
		PackageRate   int      `json:"packageRate"`
		RoomID        string   `json:"roomId"`
		RateCode      string   `json:"rateCode"`
	}

	HotelB2bAvailabilityResponseDto []AvailabilityResponseDto
)
