# TIX-HOTEL-TESTING-ENGINE-BE

## Run
```
go run testengine.go --file={test-case}.json
```
**Note** : if you not specify --file then it will use test case from storage/test-default.json

## Test Case file
you can create test case file (.json) and put it in storage folder. this is the file case format :
```
{
	"testing": [{
		"command": "CommandName",
		"data": {},
		"header": {},
		"expected": {}
	}]
}
```
Explanation : 
1. Command : Command is the case you want to test, it will hit API **based** on the command name. This is command available in this test engine : 
  - autocomplete : will test /tix-hotel-search/autocomplete [POST]
  - search : will test /tix-hotel-search/search [POST]
  - search SEO : will test /tix-hotel-search/search/seo [POST]
  - room : will test/tix-hotel-search/room [POST]
  - prebook : will test/tix-hotel-cart/hotel/prebook [POST]
  - book : will test /tix-hotel-cart/hotel/book [POST]
  - default : will test any URL that you want to test, but please put URL in expected
        ```
        "expected": {
          "url": ""
        }
        ```
2. data : it is Request Body of API that you want to test (based on Command)
3. header : it is Request Header of API that you want to test (based on Command)
4. expected : it is Response Body of API that you want to expect.

### command autocomplete, search, and room
because command autocomplete, search, and room will response array then you can expect 1 or more data in that response array. For example this is response of search
```
"data": {
    "contents": [
      {
        "id": "menara-peninsula-hotel-108001534490266880",
        "name": "Menara Peninsula Hotel",
        "address": null,
        "postalCode": null,
        "starRating": 4
        "country": {
            "id": "indonesia",
            "name": "Indonesia"
        }
        (other field)
      {
        "id": "hotel-borobudur-jakarta-108001534490282441",
        "name": "Hotel Borobudur Jakarta",
        "address": null,
        "postalCode": null,
        "starRating": 5,
        "country": {
            "id": "indonesia",
            "name": "Indonesia"
        }
        (other field)
       }
   ]
}
```

you can make expected something like this
```
"expected": {
    "responseData": [
        {
            "id": "menara-peninsula-hotel-108001534490266880",
            "name": "Menara Peninsula Hotel",
            "address": null,
            "postalCode": null,
            "starRating": 4,
        }
    ]
}
```

this will check if "responseData" you put match with the response API (you don't have to put all field, only field that you want to check). this example test will return PASS because it's match with response.

### command default
it is different for command default, you can use this command **ONLY** if the API response is array in "data".
for example you want to test /tix-hotel-search/b2b/search and this is the response :
```
"data": [
    {
      "id": "15105d18-add5-4d60-aa39-258b3e556839",
      "name": "All Seasons Jakarta Thamrin Hotel",
      "postalCode": null,
      "starRating": 3,
      "paymentOptions": null,
      (other field)
     },
     (other data)
]
```
then you can make "expected" something like this
```
"expected": {
    "url": "/tix-hotel-search/b2b/search",
    "method": "POST",
    "responseData": [
        {
            "id": "15105d18-add5-4d60-aa39-258b3e556839",
            "name": "All Seasons Jakarta Thamrin Hotel",
            "postalCode": null
        }
    ]
}
```

### Command prebook
this command will use data from **command room***, if room return 10 room then it will do prebook 10 room, but you can specify it in "expected" something like this
```
"expected": {
    "vendor": [],
    "paymentMethod": [],
    "totalRequestHit": 3
}
```
this will do prebook only 3, also you can expect response to check based on "vendor" and "paymentMethod".

### Command book
this command will use data from **command prebook***
