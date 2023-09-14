// Code generated by goa v3.12.4, DO NOT EDIT.
//
// servers HTTP server encoders and decoders
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package server

import (
	"context"
	"io"
	"net/http"

	servers "github.com/dipankardas011/crop-yield-monitor/src/authentication/gen/servers"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeLoginResponse returns an encoder for responses returned by the servers
// login endpoint.
func EncodeLoginResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*servers.Response)
		enc := encoder(ctx, w)
		body := NewLoginResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeLoginRequest returns a decoder for requests sent to the servers login
// endpoint.
func DecodeLoginRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body LoginRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateLoginRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewLoginRequest(&body)

		return payload, nil
	}
}

// EncodeSignupResponse returns an encoder for responses returned by the
// servers signup endpoint.
func EncodeSignupResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*servers.Response)
		enc := encoder(ctx, w)
		body := NewSignupResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeSignupRequest returns a decoder for requests sent to the servers
// signup endpoint.
func DecodeSignupRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body SignupRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		payload := NewSignupSignUp(&body)

		return payload, nil
	}
}

// EncodeGetHealthResponse returns an encoder for responses returned by the
// servers get health endpoint.
func EncodeGetHealthResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*servers.Health)
		enc := encoder(ctx, w)
		body := NewGetHealthResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}