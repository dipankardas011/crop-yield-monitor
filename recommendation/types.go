package main

type RecommendGet struct {
	Uuid        string `json:"uuid"`
	AccessToken string `json:"token"`
	// Crops       []string `json:"current_crops"`
}

type Recommendations struct {
	Crops []string
}

type Response struct {
	Error           string `json:"error"`
	Stdout          string `json:"stdout"`
	Recommendations any
}
