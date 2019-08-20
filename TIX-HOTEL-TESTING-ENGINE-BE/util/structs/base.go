package structs

type (
	// FileContent ...
	FileContent struct {
		Testing []ContentList `json:"testing"`
	}

	// ContentList ...
	ContentList struct {
		Command  string      `json:"command"`
		Data     interface{} `json:"data"`
		Header   interface{} `json:"header"`
		Expected interface{} `json:"expected"`
	}

	// ResponseDefault ...
	ResponseDefault struct {
		StatusCode int    `json:"status"`
		Body       string `json:"body"`
	}
)
