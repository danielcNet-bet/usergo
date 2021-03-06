package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/client"
	"github.com/thelotter-enterprise/usergo/core"
)

func TestClientIntegration(t *testing.T) {
	serviceName := "test"
	logger := makeLogger()
	ctx := context.Background()
	id := 1

	serviceDiscoverator := makeServiceDiscovery(logger)
	c := client.NewServiceClientWithDefaults(logger, serviceDiscoverator, serviceName)

	response := c.GetUserByID(ctx, id)

	if response.Data == nil {
		t.Fail()
	}
}

func makeLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", ":8080", "caller", log.DefaultCaller)

	return logger
}

func makeServiceDiscovery(logger log.Logger) *core.ServiceDiscovery {
	consulAddress := "localhost:8500"
	sd := core.NewServiceDiscovery(logger)
	sd.WithConsul(consulAddress)
	return &sd
}
