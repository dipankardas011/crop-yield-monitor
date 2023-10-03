package main

import "fmt"

type RecommendationStatus string

const (
	RecommendationReady     RecommendationStatus = "Ready"     // now we can use the results
	RecommendationPending   RecommendationStatus = "NotReady"  // we have triggered the ML to start it
	RecommendationScheduled RecommendationStatus = "Scheduled" // it helps the ML and the recommend server to know whether the Record corresponsding to the user is being processed or it is the first time triggering
)

type Recommendations struct {
	Crops  []string
	Status RecommendationStatus
}

type Response struct {
	Error           string          `json:"error"`
	Stdout          string          `json:"stdout"`
	Recommendations Recommendations `json:"recommendations"`
}

type AuthResponse struct {
	Stdout  string `json:"stdout"`
	Error   string `json:"error"`
	Account any
}

func (r AuthResponse) String() string {
	return fmt.Sprintf("{ Err: %s, Stdout: %s, Account: %v }", r.Error, r.Stdout, r.Account)
}
