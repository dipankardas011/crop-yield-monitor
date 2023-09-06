package design

import (
	. "goa.design/goa/v3/dsl" // https://pkg.go.dev/goa.design/goa/v3@v3.12.4/dsl#API
)

var _ = API("crop-yield-monitor@auth", func() {
	Title("authentication")
	Description("crop yield monitor")
	Server("addersvr", func() {
		Host("development", func() {
			URI("http://0.0.0.0:8080")
		})
	})
})

var ReqAuth = Type("Request", func() {
	Description("Type to user authenication")

	Attribute("username", String, func() {
		MinLength(1)
		Description("Username")
	})

	Attribute("password", String, func() {
		MinLength(1)
		Description("Password")
	})

	Required("username", "password")
})

var ResultAuth = Type("Response", func() {

	Description("response type")

	Attribute("ok", Boolean, "operation successful?")
	Attribute("error", String, "error reason")
	Attribute("uuid", String, "unique user identification")
})

var SignupAuth = Type("SignUp", func() {
	Attribute("first", String, "firstname")
	Attribute("last", String, "lastname")

	Attribute("password", String, "password")

	Attribute("emailid", String, "email")
})

var HealthAuth = Type("Health", func() {
	Attribute("msg", String, "message")
})

var _ = Service("servers", func() {
	Description("ksctl server handlers")

	Method("login", func() {
		Payload(ReqAuth)
		Result(ResultAuth)
		HTTP(func() {
			POST("/login")
		})
	})

	Method("signup", func() {
		Payload(SignupAuth)
		Result(ResultAuth)
		HTTP(func() {
			POST("/signup")
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
