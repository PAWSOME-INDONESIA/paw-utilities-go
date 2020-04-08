package example

import (
	"encoding/json"
	"fmt"
	cache2 "github.com/PAWSOME-INDONESIA/paw-utilities-go/cache/redis"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sync"
	"syreclabs.com/go/faker"
	"time"
)

const ITERATION = 100000

type Hotel struct {
	Address      string   `json:"address"`
	AliasName    []string `json:"aliasName"`
	Areas        []string `json:"areas"`
	CheckInTime  string   `json:"checkInTime"`
	CheckOutTime string   `json:"checkOutTime"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
	CreatedBy    string   `json:"createdBy"`
	CreatedDate  int      `json:"createdDate"`
	Facilities   []string `json:"facilities"`
	FacilityTags []string `json:"facilityTags"`
	GeoLocation  struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"geoLocation"`
	HotelBrand struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"hotelBrand"`
	HotelChain struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"hotelChain"`
	HotelID             string `json:"hotelId"`
	ID                  string `json:"id"`
	IsActive            int    `json:"isActive"`
	IsDeleted           int    `json:"isDeleted"`
	LangAlternativeName struct {
		ID string `json:"id"`
	} `json:"langAlternativeName"`
	LangAreas []struct {
		Alias string `json:"alias"`
		En    string `json:"en"`
		ID    string `json:"id"`
	} `json:"langAreas"`
	LangCity struct {
		Alias string `json:"alias"`
		En    string `json:"en"`
		ID    string `json:"id"`
	} `json:"langCity"`
	LangCountry struct {
		Cn string `json:"cn"`
		En string `json:"en"`
		ID string `json:"id"`
	} `json:"langCountry"`
	LangDescription struct {
		En string `json:"en"`
		ID string `json:"id"`
	} `json:"langDescription"`
	LangHotelPolicy  struct{}    `json:"langHotelPolicy"`
	LangPropertyType interface{} `json:"langPropertyType"`
	LangRegion       struct {
		Alias string `json:"alias"`
		En    string `json:"en"`
		ID    string `json:"id"`
	} `json:"langRegion"`
	MainVideo struct {
		Thumbnail string `json:"thumbnail"`
		URL       string `json:"url"`
	} `json:"mainVideo"`
	Name           string        `json:"name"`
	PaymentOptions []interface{} `json:"paymentOptions"`
	PostalCode     string        `json:"postalCode"`
	PropertyType   string        `json:"propertyType"`
	PublicID       string        `json:"publicId"`
	Region         string        `json:"region"`
	ReviewSummary  struct {
		Tiket struct {
			Count          int    `json:"count"`
			Impression     string `json:"impression"`
			Rating         int    `json:"rating"`
			RatingImageURL string `json:"ratingImageUrl"`
			URL            string `json:"url"`
		} `json:"tiket"`
		Tripadvisor struct {
			Count          int    `json:"count"`
			Impression     string `json:"impression"`
			Rating         int    `json:"rating"`
			RatingImageURL string `json:"ratingImageUrl"`
			URL            string `json:"url"`
		} `json:"tripadvisor"`
	} `json:"reviewSummary"`
	StarRating  int    `json:"starRating"`
	StoreID     string `json:"storeId"`
	TimeZone    string `json:"timeZone"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedDate int    `json:"updatedDate"`
	Vendors     []struct {
		IsActive         int    `json:"isActive"`
		Name             string `json:"name"`
		VendorExternalID string `json:"vendorExternalId"`
	} `json:"vendors"`
	Version     int `json:"version"`
	WeightValue int `json:"weightValue"`
}

func (h *Hotel) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

func (h *Hotel) UnmarshalBinary(msg []byte) error {
	return json.Unmarshal(msg, h)
}

type HotelSlice struct {
	hotel []*Hotel
	mu    sync.Mutex
}

func NewHotelSlice() *HotelSlice {
	return &HotelSlice{make([]*Hotel, 0), sync.Mutex{}}
}

func (hs *HotelSlice) Append(other *Hotel) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.hotel = append(hs.hotel, other)
}

func (hs *HotelSlice) Len() int {
	return len(hs.hotel)
}

func ArrayChunk(s []Hotel, size int) [][]Hotel {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]Hotel
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

func main() {
	cpuProf, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}

	defer cpuProf.Close()

	if err := pprof.StartCPUProfile(cpuProf); err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()

	cacher, err := cache2.New(&cache2.Option{
		Address:     "localhost:6379",
		Password:    "",
		DB:          0,
		PoolSize:    100,
		PoolTimeout: 10 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	rusuh := generateOtherHotel(ITERATION)

	xxx := ArrayChunk(rusuh, 500)

	hs := NewHotelSlice()
	var wg sync.WaitGroup

	now := time.Now()
	log.Printf("start get caching %s", now)

	wg.Add(len(xxx))
	for _, e := range xxx {
		go func() {
			for _, r := range e {
				hotel := Hotel{}
				if err := cacher.Get(r.ID, &hotel); err != nil {
					log.Print(err)
				}
				hs.Append(&hotel)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	log.Printf("finish caching %s", time.Now().Sub(now).String())
	log.Printf("len of hotel slice %d", hs.Len())

	//var wg sync.WaitGroup
	//
	//hs := NewHotelSlice()
	//now := time.Now()
	//log.Printf("start get caching %s", now)
	//wg.Add(len(rusuh))
	//
	//for _, r := range rusuh {
	//	go func(id string) {
	//		hotel := Hotel{}
	//		if err := cacher.Get(id, &hotel); err != nil {
	//			log.Print(err)
	//		}
	//		hs.Append(&hotel)
	//		wg.Done()
	//	}(r.ID)
	//
	//}
	//wg.Wait()
	//log.Printf("finish get caching %s", time.Now().Sub(now).String())
	//log.Printf("len of hotel slice %d", hs.Len())

	//now := time.Now()
	//log.Printf("start caching %s\n", now)
	//for _, r := range rusuh {
	//	if err := cacher.Set(r.ID, &r); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//log.Printf("finish caching %s", time.Now().Sub(now).String())
}

func generateOtherHotel(iteration int) []Hotel {
	hotels := make([]Hotel, 0)
	for i := 0; i < iteration; i++ {
		ht := Hotel{
			ID: fmt.Sprintf("%d", i+1),
		}
		hotels = append(hotels, ht)
	}
	return hotels
}

func generateHotel(iteration int) []Hotel {
	hotels := make([]Hotel, 0)

	for i := 0; i < iteration; i++ {
		hotel := Hotel{
			Address:      faker.Address().String(),
			AliasName:    []string{faker.Address().BuildingNumber(), faker.Address().StreetName()},
			Areas:        []string{faker.Address().BuildingNumber()},
			CheckInTime:  faker.Address().StreetName(),
			CheckOutTime: faker.Address().StreetName(),
			City:         faker.Address().City(),
			Country:      faker.Address().Country(),
			CreatedBy:    "",
			CreatedDate:  0,
			Facilities:   []string{faker.Commerce().ProductName(), faker.Commerce().ProductName()},
			FacilityTags: []string{faker.Commerce().ProductName(), faker.Commerce().ProductName()},
			GeoLocation: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{},
			HotelBrand: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{},
			HotelChain: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{},
			HotelID:   "",
			ID:        fmt.Sprintf("%d", i+1),
			IsActive:  0,
			IsDeleted: 0,
			LangAlternativeName: struct {
				ID string `json:"id"`
			}{},
			LangAreas: nil,
			LangCity: struct {
				Alias string `json:"alias"`
				En    string `json:"en"`
				ID    string `json:"id"`
			}{},
			LangCountry: struct {
				Cn string `json:"cn"`
				En string `json:"en"`
				ID string `json:"id"`
			}{},
			LangDescription: struct {
				En string `json:"en"`
				ID string `json:"id"`
			}{},
			LangHotelPolicy:  struct{}{},
			LangPropertyType: nil,
			LangRegion: struct {
				Alias string `json:"alias"`
				En    string `json:"en"`
				ID    string `json:"id"`
			}{},
			MainVideo: struct {
				Thumbnail string `json:"thumbnail"`
				URL       string `json:"url"`
			}{},
			Name:           "",
			PaymentOptions: nil,
			PostalCode:     "",
			PropertyType:   "",
			PublicID:       "",
			Region:         "",
			ReviewSummary: struct {
				Tiket struct {
					Count          int    `json:"count"`
					Impression     string `json:"impression"`
					Rating         int    `json:"rating"`
					RatingImageURL string `json:"ratingImageUrl"`
					URL            string `json:"url"`
				} `json:"tiket"`
				Tripadvisor struct {
					Count          int    `json:"count"`
					Impression     string `json:"impression"`
					Rating         int    `json:"rating"`
					RatingImageURL string `json:"ratingImageUrl"`
					URL            string `json:"url"`
				} `json:"tripadvisor"`
			}{},
			StarRating:  0,
			StoreID:     "",
			TimeZone:    "",
			UpdatedBy:   "",
			UpdatedDate: 0,
			Vendors:     nil,
			Version:     0,
			WeightValue: faker.RandomInt(0, 1000),
		}

		hotels = append(hotels, hotel)
	}

	return hotels
}
