package cropyieldmonitorauth

import (
	"context"
	"log"

	servers "github.com/dipankardas011/crop-yield-monitor/src/authentication/gen/servers"
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

// Login implements login.
func (s *serverssrvc) Login(ctx context.Context, p *servers.Request) (res *servers.Response, err error) {
	res = &servers.Response{}
	s.logger.Print("servers.login")
	return
}

// Signup implements signup.
func (s *serverssrvc) Signup(ctx context.Context, p *servers.SignUp) (res *servers.Response, err error) {
	res = &servers.Response{}
	s.logger.Print("servers.signup")
	return
}

// GetHealth implements get health.
func (s *serverssrvc) GetHealth(ctx context.Context) (res *servers.Health, err error) {
	res = &servers.Health{}
	s.logger.Print("servers.get health")
	return
}
