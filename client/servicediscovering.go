package client

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

// ServiceDiscovery ...
type ServiceDiscovery struct {
	ConsulAPI    *consulapi.Client
	ConsulClient *consul.Client
	Logger       log.Logger
}

// NewServiceDiscovery ...
func NewServiceDiscovery(logger log.Logger, consulAddress string) (ServiceDiscovery, error) {
	var err error
	sd := ServiceDiscovery{
		Logger: logger,
	}

	err = sd.makeConsulClient(consulAddress)
	if err != nil {
		sd.Logger.Log("method", "NewServiceDiscovery", "input", consulAddress, "err", err)
	}
	return sd, err
}

func (sd *ServiceDiscovery) makeConsulClient(consulAddress string) error {
	var err error
	config := &consulapi.Config{
		Address: consulAddress,
	}

	sd.ConsulAPI, err = consulapi.NewClient(config)

	if err == nil {
		client := consul.NewClient(sd.ConsulAPI)
		sd.ConsulClient = &client
	}

	return err
}

// ConsulInstance creates kit consul instancer which is used to find specific service
// For each service a new instance is required
func (sd *ServiceDiscovery) ConsulInstance(serviceName string, tags []string, passingOnly bool) *consul.Instancer {
	instancer := consul.NewInstancer(*sd.ConsulClient, sd.Logger, serviceName, tags, passingOnly)
	return instancer
}
