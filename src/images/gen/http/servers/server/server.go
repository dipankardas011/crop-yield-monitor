// Code generated by goa v3.13.0, DO NOT EDIT.
//
// servers HTTP server
//
// Command:
// $ goa gen github.com/dipankardas011/crop-yield-monitor/src/images/design

package server

import (
	"context"
	"net/http"

	servers "github.com/dipankardas011/crop-yield-monitor/src/images/gen/servers"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the servers service endpoint HTTP handlers.
type Server struct {
	Mounts              []*MountPoint
	Upload              http.Handler
	Fetch               http.Handler
	GetHealth           http.Handler
	GenHTTPOpenapi3JSON http.Handler
	Swaggerui           http.Handler
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the servers service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *servers.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
	fileSystemGenHTTPOpenapi3JSON http.FileSystem,
	fileSystemSwaggerui http.FileSystem,
) *Server {
	if fileSystemGenHTTPOpenapi3JSON == nil {
		fileSystemGenHTTPOpenapi3JSON = http.Dir(".")
	}
	if fileSystemSwaggerui == nil {
		fileSystemSwaggerui = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"Upload", "POST", "/upload"},
			{"Fetch", "GET", "/fetch"},
			{"GetHealth", "GET", "/healthz"},
			{"./gen/http/openapi3.json", "GET", "/openapi3.json"},
			{"./swaggerui", "GET", "/swaggerui"},
		},
		Upload:              NewUploadHandler(e.Upload, mux, decoder, encoder, errhandler, formatter),
		Fetch:               NewFetchHandler(e.Fetch, mux, decoder, encoder, errhandler, formatter),
		GetHealth:           NewGetHealthHandler(e.GetHealth, mux, decoder, encoder, errhandler, formatter),
		GenHTTPOpenapi3JSON: http.FileServer(fileSystemGenHTTPOpenapi3JSON),
		Swaggerui:           http.FileServer(fileSystemSwaggerui),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "servers" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Upload = m(s.Upload)
	s.Fetch = m(s.Fetch)
	s.GetHealth = m(s.GetHealth)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return servers.MethodNames[:] }

// Mount configures the mux to serve the servers endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountUploadHandler(mux, h.Upload)
	MountFetchHandler(mux, h.Fetch)
	MountGetHealthHandler(mux, h.GetHealth)
	MountGenHTTPOpenapi3JSON(mux, goahttp.Replace("", "/./gen/http/openapi3.json", h.GenHTTPOpenapi3JSON))
	MountSwaggerui(mux, goahttp.Replace("/swaggerui", "/./swaggerui", h.Swaggerui))
}

// Mount configures the mux to serve the servers endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountUploadHandler configures the mux to serve the "servers" service
// "upload" endpoint.
func MountUploadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/upload", f)
}

// NewUploadHandler creates a HTTP handler which loads the HTTP request and
// calls the "servers" service "upload" endpoint.
func NewUploadHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUploadRequest(mux, decoder)
		encodeResponse = EncodeUploadResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "upload")
		ctx = context.WithValue(ctx, goa.ServiceKey, "servers")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountFetchHandler configures the mux to serve the "servers" service "fetch"
// endpoint.
func MountFetchHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/fetch", f)
}

// NewFetchHandler creates a HTTP handler which loads the HTTP request and
// calls the "servers" service "fetch" endpoint.
func NewFetchHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeFetchRequest(mux, decoder)
		encodeResponse = EncodeFetchResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "fetch")
		ctx = context.WithValue(ctx, goa.ServiceKey, "servers")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGetHealthHandler configures the mux to serve the "servers" service "get
// health" endpoint.
func MountGetHealthHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/healthz", f)
}

// NewGetHealthHandler creates a HTTP handler which loads the HTTP request and
// calls the "servers" service "get health" endpoint.
func NewGetHealthHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeGetHealthResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get health")
		ctx = context.WithValue(ctx, goa.ServiceKey, "servers")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGenHTTPOpenapi3JSON configures the mux to serve GET request made to
// "/openapi3.json".
func MountGenHTTPOpenapi3JSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/openapi3.json", h.ServeHTTP)
}

// MountSwaggerui configures the mux to serve GET request made to "/swaggerui".
func MountSwaggerui(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/swaggerui/", h.ServeHTTP)
	mux.Handle("GET", "/swaggerui/*path", h.ServeHTTP)
}
