package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SKF/go-utility/v2/log"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"

	"github.com/SKF/go-rest-utility/server/gorillamux"
	"github.com/SKF/go-rest-utility/server/specification"
)

var exampleFS = os.DirFS(".")

func main() {
	router := setupRouter()

	httpServer := &http.Server{
		Addr:           ":8080",
		Handler:        &ochttp.Handler{Handler: router},
		ReadTimeout:    30 * time.Second, //nolint:gomnd
		WriteTimeout:   30 * time.Second, //nolint:gomnd
		MaxHeaderBytes: 1 << 20,          //nolint:gomnd
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serverClosed := make(chan struct{}, 1)

	go func() {
		<-sigs
		log.Info("Will try to exit gracefully")
		close(serverClosed)
	}()

	go func() {
		log.WithField("port", httpServer.Addr).Info("Will start to listen and serve")

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server ListenAndServe")
			sigs <- syscall.SIGTERM
		}
	}()

	<-serverClosed
	log.Info("Exiting")
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/").Subrouter()
	docs := router.PathPrefix("/docs").Subrouter()

	err := SetupSwaggerEndpoints(docs, api)
	if err != nil {
		log.Fatal(err)
	}

	return router
}

func SetupSwaggerEndpoints(docsEndpoints, apiEndpoints *mux.Router) error {
	route := docsEndpoints.NewRoute().PathPrefix("/swagger/")

	swaggerURL, err := route.URLPath()
	if err != nil {
		return err
	}

	endpointURL, err := apiEndpoints.Path("/").BuildOnly().URLPath()
	if err != nil {
		return err
	}

	swaggerFS := &specification.SwaggerFS{
		FS: exampleFS,

		SwaggerURLPrefix:      swaggerURL.Path,
		APIEndpointsURLPrefix: endpointURL.Path,

		NotFoundHandler:         gorillamux.NotFoundHandler(),
		MethodNotAllowedHandler: gorillamux.MethodNotFoundHandler(apiEndpoints),
	}

	opts := []specification.Option{
		specification.WithIndexFilePath("docs/index.html"),
		specification.WithSwaggerUIDirectory("docs"),
		specification.WithOpenAPIPath("oas.yaml"),
		specification.WithSwaggerInitJSPath("docs/swagger-initializer.js"),
	}
	route.Handler(swaggerFS.Handler(opts...))

	docsEndpoints.NewRoute().
		Path("/swagger").
		Handler(http.RedirectHandler(swaggerURL.Path, http.StatusMovedPermanently))

	return nil
}
