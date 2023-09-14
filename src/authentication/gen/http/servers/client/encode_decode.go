// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers HTTP client encoders and decoders
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	servers "github.com/dipankardas011/crop-yield-monitor/src/authentication/gen/servers"
	goahttp "goa.design/goa/v3/http"
)

// BuildLoginRequest instantiates a HTTP request object with method and path
// set to call the "servers" service "login" endpoint
func (c *Client) BuildLoginRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: LoginServersPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("servers", "login", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeLoginRequest returns an encoder for requests sent to the servers login
// server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*servers.Request)
		if !ok {
			return goahttp.ErrInvalidType("servers", "login", "*servers.Request", v)
		}
		body := NewLoginRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("servers", "login", err)
		}
		return nil
	}
}

// DecodeLoginResponse returns a decoder for responses returned by the servers
// login endpoint. restoreBody controls whether the response body should be
// restored after having been read.
func DecodeLoginResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body LoginResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("servers", "login", err)
			}
			res := NewLoginResponseOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("servers", "login", resp.StatusCode, string(body))
		}
	}
}

// BuildSignupRequest instantiates a HTTP request object with method and path
// set to call the "servers" service "signup" endpoint
func (c *Client) BuildSignupRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: SignupServersPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("servers", "signup", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeSignupRequest returns an encoder for requests sent to the servers
// signup server.
func EncodeSignupRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*servers.SignUp)
		if !ok {
			return goahttp.ErrInvalidType("servers", "signup", "*servers.SignUp", v)
		}
		body := NewSignupRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("servers", "signup", err)
		}
		return nil
	}
}

// DecodeSignupResponse returns a decoder for responses returned by the servers
// signup endpoint. restoreBody controls whether the response body should be
// restored after having been read.
func DecodeSignupResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body SignupResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("servers", "signup", err)
			}
			res := NewSignupResponseOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("servers", "signup", resp.StatusCode, string(body))
		}
	}
}

// BuildGetHealthRequest instantiates a HTTP request object with method and
// path set to call the "servers" service "get health" endpoint
func (c *Client) BuildGetHealthRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetHealthServersPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("servers", "get health", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetHealthResponse returns a decoder for responses returned by the
// servers get health endpoint. restoreBody controls whether the response body
// should be restored after having been read.
func DecodeGetHealthResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetHealthResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("servers", "get health", err)
			}
			res := NewGetHealthHealthOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("servers", "get health", resp.StatusCode, string(body))
		}
	}
}