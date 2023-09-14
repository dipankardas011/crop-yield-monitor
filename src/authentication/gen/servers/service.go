// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers service
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package servers

import (
	"context"
)

// ksctl server handlers
type Service interface {
	// Login implements login.
	Login(context.Context, *Request) (res *Response, err error)
	// Signup implements signup.
	Signup(context.Context, *SignUp) (res *Response, err error)
	// GetHealth implements get health.
	GetHealth(context.Context) (res *Health, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "servers"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"login", "signup", "get health"}

// Health is the result type of the servers service get health method.
type Health struct {
	// message
	Msg *string
}

// Request is the payload type of the servers service login method.
type Request struct {
	// Username
	Username string
	// Password
	Password string
}

// Response is the result type of the servers service login method.
type Response struct {
	// operation successful?
	OK *bool
	// error reason
	Error *string
	// unique user identification
	UUID *string
}

// SignUp is the payload type of the servers service signup method.
type SignUp struct {
	// first name
	First *string
	// last name
	Last *string
	// password
	Password *string
	// emailid
	Emailid *string
}