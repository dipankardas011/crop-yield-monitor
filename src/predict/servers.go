package cropyieldmonitorpredict

import (
	"context"
	"log"

	servers "github.com/dipankardas011/crop-yield-monitor/src/predict/gen/servers"
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

// Predict implements predict.
func (s *serverssrvc) Predict(ctx context.Context, p *servers.ReqPrediction) (res *servers.Recommendations, err error) {
	res = &servers.Recommendations{}
	s.logger.Print("servers.predict")
	return
}

// GetHealth implements get health.
func (s *serverssrvc) GetHealth(ctx context.Context) (res *servers.Health, err error) {
	res = &servers.Health{}
	s.logger.Print("servers.get health")
	return
}
