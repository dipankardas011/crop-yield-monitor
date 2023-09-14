package design

import (
	. "goa.design/goa/v3/dsl" // https://pkg.go.dev/goa.design/goa/v3@v3.12.4/dsl#API
)

var _ = API("crop-yield-monitor@predict", func() {
	Title("predictions")
	Description("crop yield monitor")
	Server("addersvr", func() {
		Host("development", func() {
			URI("http://0.0.0.0:8082")
		})
	})
})

var ReqPred = Type("ReqPrediction", func() {
	Description("request payload for prediction")

	Attribute("uuid", String, func() {
		Description("unique userid")
		Example("1")
		MinLength(1)
	})

	Required("uuid")

})

var ResultPred = Type("Recommendations", func() {

	Description("response type")

	Attribute("ok", Boolean, "operation successful?")
	Attribute("waiting", Boolean, "processing in progress")
	Attribute("error", String, "error reason")

	Attribute("recommendations", func() {
		Description("recommendations")

		Attribute("crops", ArrayOf(String), "recommended crops", func() {})

		Attribute("need_water", Boolean)

	})

	Required("ok", "error", "waiting")
})

var HealthAuth = Type("Health", func() {
	Attribute("msg", String, "message")
})

var _ = Service("servers", func() {
	Description("server handlers")

	Method("predict", func() {
		Payload(ReqPred)
		Result(ResultPred)
		HTTP(func() {
			GET("/predict")
		})
	})

	Method("get health", func() {
		Result(HealthAuth)
		HTTP(func() {
			GET("/healthz")
		})
	})

	Files("/openapi3.json", "./gen/http/openapi3.json")
	Files("/swaggerui/{*path}", "./swaggerui")
})
