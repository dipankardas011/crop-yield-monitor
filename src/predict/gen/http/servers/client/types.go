// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers HTTP client types
//
// Command:
// $ goa gen github.com/dipankardas011/crop-yield-monitor/src/predict/design

package client

import (
	servers "github.com/dipankardas011/crop-yield-monitor/src/predict/gen/servers"
	goa "goa.design/goa/v3/pkg"
)

// PredictRequestBody is the type of the "servers" service "predict" endpoint
// HTTP request body.
type PredictRequestBody struct {
	// unique userid
	UUID string `form:"uuid" json:"uuid" xml:"uuid"`
}

// PredictResponseBody is the type of the "servers" service "predict" endpoint
// HTTP response body.
type PredictResponseBody struct {
	// operation successful?
	OK *bool `form:"ok,omitempty" json:"ok,omitempty" xml:"ok,omitempty"`
	// processing in progress
	Waiting *bool `form:"waiting,omitempty" json:"waiting,omitempty" xml:"waiting,omitempty"`
	// error reason
	Error *string `form:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
	// recommendations
	Recommendations *struct {
		// recommended crops
		Crops     []string `form:"crops" json:"crops" xml:"crops"`
		NeedWater *bool    `form:"need_water" json:"need_water" xml:"need_water"`
	} `form:"recommendations,omitempty" json:"recommendations,omitempty" xml:"recommendations,omitempty"`
}

// GetHealthResponseBody is the type of the "servers" service "get health"
// endpoint HTTP response body.
type GetHealthResponseBody struct {
	// message
	Msg *string `form:"msg,omitempty" json:"msg,omitempty" xml:"msg,omitempty"`
}

// NewPredictRequestBody builds the HTTP request body from the payload of the
// "predict" endpoint of the "servers" service.
func NewPredictRequestBody(p *servers.ReqPrediction) *PredictRequestBody {
	body := &PredictRequestBody{
		UUID: p.UUID,
	}
	return body
}

// NewPredictRecommendationsOK builds a "servers" service "predict" endpoint
// result from a HTTP "OK" response.
func NewPredictRecommendationsOK(body *PredictResponseBody) *servers.Recommendations {
	v := &servers.Recommendations{
		OK:      *body.OK,
		Waiting: *body.Waiting,
		Error:   *body.Error,
	}
	if body.Recommendations != nil {
		v.Recommendations = &struct {
			// recommended crops
			Crops     []string
			NeedWater *bool
		}{
			NeedWater: body.Recommendations.NeedWater,
		}
		if body.Recommendations.Crops != nil {
			v.Recommendations.Crops = make([]string, len(body.Recommendations.Crops))
			for i, val := range body.Recommendations.Crops {
				v.Recommendations.Crops[i] = val
			}
		}
	}

	return v
}

// NewGetHealthHealthOK builds a "servers" service "get health" endpoint result
// from a HTTP "OK" response.
func NewGetHealthHealthOK(body *GetHealthResponseBody) *servers.Health {
	v := &servers.Health{
		Msg: body.Msg,
	}

	return v
}

// ValidatePredictResponseBody runs the validations defined on
// PredictResponseBody
func ValidatePredictResponseBody(body *PredictResponseBody) (err error) {
	if body.OK == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ok", "body"))
	}
	if body.Error == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("error", "body"))
	}
	if body.Waiting == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("waiting", "body"))
	}
	return
}
