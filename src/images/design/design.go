package design

import (
	. "goa.design/goa/v3/dsl" // https://pkg.go.dev/goa.design/goa/v3@v3.12.4/dsl#API
)

var _ = API("crop-yield-monitor@images", func() {
	Title("image")
	Description("crop yield monitor")
	Server("addersvr", func() {
		Host("development", func() {
			URI("http://0.0.0.0:8081")
		})
	})
})

var ReqGet = Type("ReqGet", func() {
	Description("request payload for fetching image")

	Attribute("uuid", String, func() {
		Description("unique userid")
		Example("1")
		MinLength(1)
	})

	Required("uuid")

})

var ReqUpload = Type("ReqUpload", func() {
	Description("request payload for uploading")

	Attribute("uuid", String, func() {
		Description("unique userid")
		Example("1")
		MinLength(1)
	})

	Attribute("image", Bytes, func() {
		Description("image in byte array")
	})

	Required("image", "uuid")
})

var ResultImg = Type("Response", func() {

	Description("response type")

	Attribute("ok", Boolean, "operation successful?")
	Attribute("error", String, "error reason")
	Attribute("image", Bytes, "image in []byte")

	Required("ok", "error")
})

var HealthAuth = Type("Health", func() {
	Attribute("msg", String, "message")
})

var _ = Service("servers", func() {
	Description("server handlers")

	Method("upload", func() {
		Payload(ReqUpload)
		Result(ResultImg)

		HTTP(func() {
			POST("/upload")
		})
	})

	Method("fetch", func() {
		Payload(ReqGet)
		Result(ResultImg)

		HTTP(func() {
			GET("/fetch")
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
