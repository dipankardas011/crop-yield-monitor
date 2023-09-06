// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers HTTP client types
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package client

import (
	servers "github.com/dipankardas011/crop-yield-monitor/src/authentication/gen/servers"
)

// LoginRequestBody is the type of the "servers" service "login" endpoint HTTP
// request body.
type LoginRequestBody struct {
	// Username
	Username string `form:"username" json:"username" xml:"username"`
	// Password
	Password string `form:"password" json:"password" xml:"password"`
}

// SignupRequestBody is the type of the "servers" service "signup" endpoint
// HTTP request body.
type SignupRequestBody struct {
	// first name
	First *string `form:"first,omitempty" json:"first,omitempty" xml:"first,omitempty"`
	// last name
	Last *string `form:"last,omitempty" json:"last,omitempty" xml:"last,omitempty"`
	// password
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	// emailid
	Emailid *string `form:"emailid,omitempty" json:"emailid,omitempty" xml:"emailid,omitempty"`
}

// LoginResponseBody is the type of the "servers" service "login" endpoint HTTP
// response body.
type LoginResponseBody struct {
	// operation successful?
	OK *bool `form:"ok,omitempty" json:"ok,omitempty" xml:"ok,omitempty"`
	// error reason
	Error *string `form:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
	// unique user identification
	UUID *string `form:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
}

// SignupResponseBody is the type of the "servers" service "signup" endpoint
// HTTP response body.
type SignupResponseBody struct {
	// operation successful?
	OK *bool `form:"ok,omitempty" json:"ok,omitempty" xml:"ok,omitempty"`
	// error reason
	Error *string `form:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
	// unique user identification
	UUID *string `form:"uuid,omitempty" json:"uuid,omitempty" xml:"uuid,omitempty"`
}

// GetHealthResponseBody is the type of the "servers" service "get health"
// endpoint HTTP response body.
type GetHealthResponseBody struct {
	// message
	Msg *string `form:"msg,omitempty" json:"msg,omitempty" xml:"msg,omitempty"`
}

// NewLoginRequestBody builds the HTTP request body from the payload of the
// "login" endpoint of the "servers" service.
func NewLoginRequestBody(p *servers.Request) *LoginRequestBody {
	body := &LoginRequestBody{
		Username: p.Username,
		Password: p.Password,
	}
	return body
}

// NewSignupRequestBody builds the HTTP request body from the payload of the
// "signup" endpoint of the "servers" service.
func NewSignupRequestBody(p *servers.SignUp) *SignupRequestBody {
	body := &SignupRequestBody{
		First:    p.First,
		Last:     p.Last,
		Password: p.Password,
		Emailid:  p.Emailid,
	}
	return body
}

// NewLoginResponseOK builds a "servers" service "login" endpoint result from a
// HTTP "OK" response.
func NewLoginResponseOK(body *LoginResponseBody) *servers.Response {
	v := &servers.Response{
		OK:    body.OK,
		Error: body.Error,
		UUID:  body.UUID,
	}

	return v
}

// NewSignupResponseOK builds a "servers" service "signup" endpoint result from
// a HTTP "OK" response.
func NewSignupResponseOK(body *SignupResponseBody) *servers.Response {
	v := &servers.Response{
		OK:    body.OK,
		Error: body.Error,
		UUID:  body.UUID,
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
