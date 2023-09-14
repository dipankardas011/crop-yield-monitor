// Code generated by goa v3.13.0, DO NOT EDIT.
//
// HTTP request path constructors for the servers service.
//
// Command:
// $ goa gen github.com/dipankardas011/crop-yield-monitor/src/images/design

package server

// UploadServersPath returns the URL path to the servers service upload HTTP endpoint.
func UploadServersPath() string {
	return "/upload"
}

// FetchServersPath returns the URL path to the servers service fetch HTTP endpoint.
func FetchServersPath() string {
	return "/fetch"
}

// GetHealthServersPath returns the URL path to the servers service get health HTTP endpoint.
func GetHealthServersPath() string {
	return "/healthz"
}
