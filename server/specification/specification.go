package specification

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"text/template"

	"go.opencensus.io/trace"

	"github.com/SKF/go-rest-utility/problems"
)

const (
	swaggerUIDirectory = "swagger-ui"
	indexFilePath      = "swagger-ui/index.html"
	openAPIPath        = "/openapi.yaml"
)

type SwaggerFS struct {
	FS fs.FS

	SwaggerURLPrefix      string
	APIEndpointsURLPrefix string

	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
}

type swaggerFSHandler struct {
	*SwaggerFS
	*Config
	pathsToCompile []string
	initOnce       sync.Once
	compiledFiles  map[string]CompiledFile
}

type CompiledFile struct {
	Contents *bytes.Buffer
	Info     fs.FileInfo
}

func (swaggerfs *SwaggerFS) Handler(opts ...Option) http.Handler {
	config := &Config{
		swaggerUIDirectory: swaggerUIDirectory,
		indexFilePath:      indexFilePath,
		openAPIPath:        openAPIPath,
	}

	for _, opt := range opts {
		opt(config)
	}

	pathsToCompile := []string{
		config.indexFilePath,
		config.openAPIPath,
	}

	if swaggerfs.FS == nil {
		panic("no swagger filesystem set")
	}

	if swaggerfs.NotFoundHandler == nil {
		panic("no NotFoundHandler set")
	}

	if swaggerfs.MethodNotAllowedHandler == nil {
		panic("no MethodNotAllowedHandler set")
	}

	return &swaggerFSHandler{
		SwaggerFS:      swaggerfs,
		Config:         config,
		pathsToCompile: pathsToCompile,
	}
}

func (swagger *swaggerFSHandler) init(r *http.Request) func() {
	return func() {
		baseURL := *BaseURLFromRequest(r)

		endpointBaseURL := baseURL
		endpointBaseURL.Path = path.Clean(swagger.APIEndpointsURLPrefix)

		yamlURL := baseURL
		yamlURL.Path = path.Clean(swagger.SwaggerURLPrefix) + swagger.openAPIPath

		data := struct {
			EndpointBaseURL *url.URL
			YamlURL         *url.URL
		}{
			EndpointBaseURL: &endpointBaseURL,
			YamlURL:         &yamlURL,
		}

		for _, path := range swagger.pathsToCompile {
			if err := swagger.Compile(path, data); err != nil {
				panic(err)
			}
		}
	}
}

func (swagger *swaggerFSHandler) Open(name string) (io.ReadSeekCloser, fs.FileInfo, error) {
	if compiled, found := swagger.compiledFiles[name]; found {
		reader := bytes.NewReader(compiled.Contents.Bytes())
		return NopSeekCloser(reader), compiled.Info, nil
	}

	f, err := http.FS(swagger.FS).Open(name)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open static file: %w", err)
	}

	d, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to stat static file: %w", err)
	}

	return f, d, nil
}

func (swagger *swaggerFSHandler) Compile(name string, data interface{}) error {
	pattern := strings.TrimPrefix(name, "/")

	t, err := template.ParseFS(swagger.FS, pattern)
	if err != nil {
		return fmt.Errorf("unable to parse template for compilation: %w", err)
	}

	f, err := swagger.FS.Open(pattern)
	if err != nil {
		return fmt.Errorf("unable to open file for compilation: %w", err)
	}

	defer f.Close()

	compiledFile := CompiledFile{
		Contents: new(bytes.Buffer),
		Info:     nil,
	}

	if compiledFile.Info, err = f.Stat(); err != nil {
		return fmt.Errorf("unable to stat template file: %w", err)
	}

	if err := t.Execute(compiledFile.Contents, data); err != nil {
		return fmt.Errorf("unable to process template: %w", err)
	}

	if swagger.compiledFiles == nil {
		swagger.compiledFiles = make(map[string]CompiledFile)
	}

	swagger.compiledFiles[name] = compiledFile

	return nil
}

func (swagger *swaggerFSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "Server/ServeSwagger")
	defer span.End()

	swagger.initOnce.Do(swagger.init(r))

	if r.Method != http.MethodGet {
		if swagger.MethodNotAllowedHandler != nil {
			swagger.MethodNotAllowedHandler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "405 only GET is allowed", http.StatusMethodNotAllowed)

		return
	}

	// Extract the desired file by removing the URLPrefix
	// 	/docs/swagger/openapi.yaml => /openapi.yaml
	name := strings.TrimPrefix(
		path.Clean(r.URL.Path),
		path.Clean(swagger.SwaggerURLPrefix),
	)

	if name == "" {
		name = swagger.indexFilePath
	} else if name != swagger.openAPIPath {
		name = swagger.swaggerUIDirectory + name
	}

	f, d, err := swagger.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if swagger.NotFoundHandler != nil {
				swagger.NotFoundHandler.ServeHTTP(w, r)
				return
			}

			http.Error(w, "404 page not found", http.StatusNotFound)

			return
		}

		problems.WriteResponse(ctx, err, w, r)

		return
	}

	defer f.Close()

	// http.ServeContent will attempt to do this from the systems MIME database.
	// We need to use our own here to avoid issues where the CI system does not
	// have an up to date database. If the extension is unknown it will fallback to
	// using the system database.
	if ctype := ContentTypeFromFilename(d.Name()); ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
