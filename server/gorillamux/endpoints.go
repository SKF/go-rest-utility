package gorillamux

import (
	"net/http"

	authentication_mw "github.com/SKF/go-enlight-middleware/authentication"
	authorization_mw "github.com/SKF/go-enlight-middleware/authorization"
	http_middleware "github.com/SKF/go-utility/v2/http-middleware"
	"github.com/SKF/go-utility/v2/log"
	"github.com/gorilla/mux"
)

type Endpoint struct {
	Name           string
	PathTpl        string
	authentication *authentication_mw.Middleware
	authorization  *authorization_mw.Middleware
	router         *mux.Router
	*mux.Route
}

func NewEndpoint(name string, router *mux.Router) *Endpoint {
	return &Endpoint{
		Name:   name,
		Route:  router.NewRoute().Name(name),
		router: router,
	}
}

func (e *Endpoint) WithAuthorization(authorization *authorization_mw.Middleware) *Endpoint {
	e.authorization = authorization
	return e
}

func (e *Endpoint) WithAuthentication(authentication *authentication_mw.Middleware) *Endpoint {
	e.authentication = authentication
	return e
}

func (e *Endpoint) IgnoreInAuthentication(authentication *authentication_mw.Middleware) *Endpoint {
	e.authentication = authentication
	e.authentication.IgnoreRoute(e.Route)
	return e
}

func (e *Endpoint) Path(tpl string) *Endpoint {
	e.PathTpl = tpl
	e.Route.Path(tpl)

	err := e.Route.GetError()
	if err != nil {
		log.Errorf("Failed to create path", err)
	}

	return e
}

func (e *Endpoint) Methods(methods ...string) *Endpoint {
	e.Route.Methods(methods...)
	return e
}

func (e *Endpoint) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Endpoint {
	e.Route.HandlerFunc(f)
	return e
}

func (e *Endpoint) WithCors(settings []string) *Endpoint {
	methods, err := e.Route.GetMethods()
	if err != nil {
		log.Error(err)
	}

	corsRoute := e.router.NewRoute()
	corsRoute.Methods(http.MethodOptions)
	corsRoute.HandlerFunc(http_middleware.Options(
		methods,
		settings,
	))
	corsRoute.Path(e.PathTpl)
	e.authentication.IgnoreRoute(corsRoute)

	return e
}

func (e *Endpoint) WithAuthPolicy(policy authorization_mw.Policy) *Endpoint {
	e.authorization.SetPolicy(e.Route, policy)
	return e
}
