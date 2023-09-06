// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers client
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package servers

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "servers" service client.
type Client struct {
	LoginEndpoint     goa.Endpoint
	SignupEndpoint    goa.Endpoint
	GetHealthEndpoint goa.Endpoint
}

// NewClient initializes a "servers" service client given the endpoints.
func NewClient(login, signup, getHealth goa.Endpoint) *Client {
	return &Client{
		LoginEndpoint:     login,
		SignupEndpoint:    signup,
		GetHealthEndpoint: getHealth,
	}
}

// Login calls the "login" endpoint of the "servers" service.
func (c *Client) Login(ctx context.Context, p *Request) (res *Response, err error) {
	var ires any
	ires, err = c.LoginEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Response), nil
}

// Signup calls the "signup" endpoint of the "servers" service.
func (c *Client) Signup(ctx context.Context, p *SignUp) (res *Response, err error) {
	var ires any
	ires, err = c.SignupEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Response), nil
}

// GetHealth calls the "get health" endpoint of the "servers" service.
func (c *Client) GetHealth(ctx context.Context) (res *Health, err error) {
	var ires any
	ires, err = c.GetHealthEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*Health), nil
}
