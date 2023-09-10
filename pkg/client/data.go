package client

type DataLatest struct {
	Data map[string]ItemLatest `json:"data"`
}

type ItemLatest struct {
	High     *int `json:"high,omitempty"`
	HighTime *int `json:"highTime,omitempty"`
	Low      *int `json:"low,omitempty"`
	LowTime  *int `json:"lowTime,omitempty"`
}

type DataAvg struct {
	Data map[string]ItemAvg `json:"data"`
}

type ItemAvg struct {
	AvgHighPrice    *int `json:"avgHighPrice,omitempty"`
	HighPriceVolume *int `json:"highPriceVolume,omitempty"`
	AvgLowPrice     *int `json:"avgLowPrice,omitempty"`
	LowPriceVolume  *int `json:"lowPriceVolume,omitempty"`
}

type ItemMapping struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    int    `json:"value"`
	Examine  string `json:"examine"`
	Members  bool   `json:"members"`
	Icon     string `json:"icon"`
	Highalch *int   `json:"highalch,omitempty"`
	Lowalch  *int   `json:"lowalch,omitempty"`
	Limit    *int   `json:"limit,omitempty"`
}
