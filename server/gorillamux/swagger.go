package gorillamux

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	authentication_mw "github.com/SKF/go-enlight-middleware/authentication"
	"github.com/SKF/go-utility/v2/log"
	"github.com/gorilla/mux"
)

type IndexTemplate struct {
	PublicURL string
}

// Register a swagger yml to endpoints:
//  /docs/swagger
//  /docs/swagger/
//  /docs/swagger/index.html
func SetupSwaggerEndpoints(openapiFilesPath, publicURL string, router *mux.Router, authentication *authentication_mw.Middleware, opts ...Option) {
	config := &Config{
		swaggerPath: "/docs/swagger",
		docsPath:    "docs/",
		indexFile:   "index.html",
	}

	for _, opt := range opts {
		opt(config)
	}

	swaggerPath := "/docs/swagger"

	localDocsPath := path.Join(openapiFilesPath, config.docsPath)
	swaggerFilePath := path.Join(localDocsPath, config.indexFile)

	var templates = template.Must(template.ParseFiles(swaggerFilePath))

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		templateName := config.indexFile

		t := templates.Lookup(templateName)
		if err := t.Execute(w, IndexTemplate{PublicURL: publicURL}); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}

	NewEndpoint("APISpecificationIndex", router).
		IgnoreInAuthentication(authentication).
		Methods(http.MethodGet).
		Path(swaggerPath + "/index.html").
		HandlerFunc(indexHandler)

	swaggerIndexPublicPath := publicURL + "/" + swaggerPath + "/index.html"
	NewEndpoint("APISpecificationTrailingSlash", router).
		IgnoreInAuthentication(authentication).
		Methods(http.MethodGet).
		Path(swaggerPath + "/").
		Handler(http.RedirectHandler(swaggerIndexPublicPath, http.StatusMovedPermanently))

	NewEndpoint("APISpecificationNoTrailingSlash", router).
		IgnoreInAuthentication(authentication).
		Methods(http.MethodGet).
		Path(swaggerPath).
		Handler(http.RedirectHandler(swaggerIndexPublicPath, http.StatusMovedPermanently))

	swaggerFilterFS := filterFS{rootDir: localDocsPath, allowdExtenstions: []string{"yaml", "png", "html", "css", "js"}}
	NewEndpoint("APISpecification", router).
		IgnoreInAuthentication(authentication).
		PathPrefix(swaggerPath).
		Methods(http.MethodGet).
		Handler(http.StripPrefix(swaggerPath, http.FileServer(swaggerFilterFS)))
}

type filterFS struct {
	rootDir           string
	allowdExtenstions []string
}

func open(rootDir, name string, extensions []string) (http.File, error) {
	re := `[a-zA-z0-9-_]*\.(` + strings.Join(extensions, "|") + ")"

	matches, err := regexp.MatchString(re, name)
	if matches {
		return os.Open(path.Join(rootDir, name))
	}

	log.WithError(err).Warn("trying to request invalid file")

	return nil, os.ErrNotExist
}

func (fs filterFS) Open(name string) (http.File, error) {
	return open(fs.rootDir, name, fs.allowdExtenstions)
}
