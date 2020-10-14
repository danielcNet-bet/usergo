package svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"

	httpkit "github.com/go-kit/kit/transport/http"
)

// Endpoints ...
type Endpoints struct {
	Log     core.Log
	Tracer  Tracer
	Service Service

	ServerEndpoints []EndpointEntry
}

// EndpointEntry holds the information needed to build a server endpoint which client can call upon
type EndpointEntry struct {
	Method   string
	Endpoint func(ctx context.Context, request interface{}) (interface{}, error)
	Dec      httpkit.DecodeRequestFunc
	Enc      httpkit.EncodeResponseFunc
}

// NewEndpoints ...
func NewEndpoints(log core.Log, tracer Tracer, service Service) Endpoints {
	endpoints := Endpoints{
		Log:     log,
		Tracer:  tracer,
		Service: service,
	}

	endpoints.AddEndpoints()

	return endpoints
}

// AddEndpoints ...
func (endpoints *Endpoints) AddEndpoints() {
	var serverEndpoints []EndpointEntry

	userbyid := EndpointEntry{
		Endpoint: endpoints.makeUserByIDEndpoint(),
		Enc:      core.EncodeReponseToJSON,
		Dec:      core.DecodeRequestFromJSON,
		Method:   "GET",
	}

	endpoints.ServerEndpoints = append(serverEndpoints, userbyid)
}

func (endpoints *Endpoints) makeUserByIDEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		var req core.Request
		var data shared.ByIDRequestData

		decoder := core.NewDecoder()

		err = decoder.MapDecode(request, &req)
		err = decoder.MapDecode(req.Data, &data)
		req.Data = data

		user, err := endpoints.Service.GetUserByID(ctx, data.ID)
		return shared.NewUserResponse(user), err
	}
}
