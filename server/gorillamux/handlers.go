package gorillamux

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/SKF/go-rest-utility/problems"
	handlerproblems "github.com/SKF/go-rest-utility/server/gorillamux/problems"
)

func MethodNotFoundHandler(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedMethods := []string{}

		router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error { //nolint: errcheck
			var match mux.RouteMatch

			matched := route.Match(r, &match)
			if matched || match.MatchErr == mux.ErrMethodMismatch {
				if methods, err := route.GetMethods(); err == nil {
					allowedMethods = append(allowedMethods, methods...)
				}
			}

			return nil
		})

		w.Header().Set("Allow", strings.Join(allowedMethods, ","))

		problem := handlerproblems.MethodNotAllowed(
			r.Method,
			allowedMethods...,
		)
		problems.WriteResponse(r.Context(), problem, w, r)
	}
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problems.WriteResponse(r.Context(), handlerproblems.NotFound(), w, r)
	}
}
