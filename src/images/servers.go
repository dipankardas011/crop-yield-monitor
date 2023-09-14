package cropyieldmonitorimages

import (
	"context"
	"log"

	servers "github.com/dipankardas011/crop-yield-monitor/src/images/gen/servers"
)

// servers service example implementation.
// The example methods log the requests and return zero values.
type serverssrvc struct {
	logger *log.Logger
}

// NewServers returns the servers service implementation.
func NewServers(logger *log.Logger) servers.Service {
	return &serverssrvc{logger}
}

// Upload implements upload.
func (s *serverssrvc) Upload(ctx context.Context, p *servers.ReqUpload) (res *servers.Response, err error) {
	res = &servers.Response{}
	s.logger.Print("servers.upload")
	return
}

// Fetch implements fetch.
func (s *serverssrvc) Fetch(ctx context.Context, p *servers.ReqGet) (res *servers.Response, err error) {
	res = &servers.Response{}
	s.logger.Print("servers.fetch")
	return
}

// GetHealth implements get health.
func (s *serverssrvc) GetHealth(ctx context.Context) (res *servers.Health, err error) {
	res = &servers.Health{}
	s.logger.Print("servers.get health")
	return
}
