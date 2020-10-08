package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
)

// Proxy ...
type Proxy struct {
	cb         core.CircuitBreaker
	limmitermw endpoint.Middleware
	router     *mux.Router
	limiter    core.RateLimiter
}

type userByIDProxyMiddleware struct {
	// Context holds the context
	Context context.Context

	// Next is a the service instance
	// We need to use Next, since it is used to satisfy the middleware pattern
	// Each middleware is responbsible for a single API, yet, due to the service interface,
	// it need to implement all the service interface APIs. To support it, we use Next to obstract the implementation
	Next interface{}

	// This is the current API which we plan to support in the service interface contract
	This endpoint.Endpoint
}

// NewProxy ..
func NewProxy(cb core.CircuitBreaker, limiter core.RateLimiter, router *mux.Router) Proxy {
	return Proxy{
		cb:      cb,
		limiter: limiter,
		router:  router,
	}
}

// UserByIDMiddleware ..
func (proxy Proxy) UserByIDMiddleware(ctx context.Context, id int) UserServiceMiddleware {
	commandName := "get_user_by_id"
	var endpointer sd.FixedEndpointer
	breakermw := proxy.cb.NewDefaultHystrixCommandMiddleware(commandName)
	limitermw := proxy.limiter.NewDefaultErrorLimitterMiddleware()
	tgt, _ := proxy.router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))
	e := httptransport.NewClient("GET", tgt, core.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()
	e = breakermw(e)
	e = limitermw(e)
	endpointer = append(endpointer, e)

	lb := core.NewLoadBalancer(endpointer)
	retry := lb.DefaultRoundRobinWithRetryEndpoint(ctx)

	return func(next UserService) UserService {
		return userByIDProxyMiddleware{Context: ctx, Next: next, This: retry}
	}
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
func (proxymw userByIDProxyMiddleware) GetUserByID(id int) core.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = proxymw.This(proxymw.Context, id); err != nil {
		// TODO: need a refactor to analyze the response
		circuitOpen = true
		statusCode = 500
	}

	return core.HTTPResponse{
		Result:        res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (proxymw userByIDProxyMiddleware) GetUserByEmail(email string) core.HTTPResponse {
	svc := proxymw.Next.(UserService)
	return svc.GetUserByEmail(email)
}