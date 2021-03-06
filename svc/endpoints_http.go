package svc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
)

// UserHTTPEndpoints ...
type UserHTTPEndpoints struct {
	HTTPEndpoints *core.HTTPEndpoints
	Service       Service
	Log           core.Log
	Tracer        core.Tracer
}

// NewUserHTTPEndpoints ...
func NewUserHTTPEndpoints(log core.Log, tracer core.Tracer, service Service) *UserHTTPEndpoints {
	userEndpoints := UserHTTPEndpoints{
		Log:           log,
		Tracer:        tracer,
		Service:       service,
		HTTPEndpoints: &core.HTTPEndpoints{},
	}

	userEndpoints.HTTPEndpoints = userEndpoints.makeEndpoints()

	return &userEndpoints
}

func (ue UserHTTPEndpoints) makeEndpoints() *core.HTTPEndpoints {
	var endpoints core.HTTPEndpoints
	var serverEndpoints []core.HTTPEndpoint

	userbyid := core.HTTPEndpoint{
		Endpoint: makeUserByIDEndpoint(ue.Service),
		Enc:      ue.encodeUserByIDReponse,
		Dec:      ue.decodeUserByIDRequest,
		Method:   "GET",
		Path:     shared.UserByIDServerRoute,
	}

	serverEndpoints = append(serverEndpoints, userbyid)
	endpoints.ServerEndpoints = serverEndpoints
	return &endpoints
}

func makeUserByIDEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		var req core.Request
		var data shared.ByIDRequestData

		decoder := core.NewDecoder()

		err = decoder.MapDecode(request, &req)
		err = decoder.MapDecode(req.Data, &data)
		req.Data = data

		user, err := service.GetUserByID(ctx, data.ID)
		return shared.NewUserResponse(user), err
	}
}

func (ue UserHTTPEndpoints) encodeUserByIDReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ue.Log.Logger.Log("method", "EncodeReponseToJSONFunc", "error", err)
	}
	return err
}

func (ue UserHTTPEndpoints) decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		ue.Log.Logger.Log(
			"level", "error",
			"method", "DecodeRequestFromJSONFunc",
			"error", err,
		)
	}

	return req, err
}
