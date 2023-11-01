CURL ?= curl
TAR ?= tar
INSTALL ?= install
SED ?= sed
ECHO ?= echo

SWAGGER_UI_VERSION ?= 5.9.1

all: docs

docs/service.yaml:
	$(INSTALL) -d docs
	cp oas.yaml $@

swagger-ui-%.tar.gz:
	$(RM) -f "$@"
	$(CURL) -sSLJO "https://github.com/swagger-api/swagger-ui/archive/v$*.tar.gz"

docs: swagger-ui-${SWAGGER_UI_VERSION}.tar.gz docs/service.yaml
# Clean any previous built files and ensure the folder exists for tar to extract into.
	$(INSTALL) -d "docs"
# Extract compiled Swagger UI from source package.
	$(TAR) --extract \
        --file="$<" \
        --gzip \
        --directory=$@ \
        --strip-components=2 \
        "swagger-ui-${SWAGGER_UI_VERSION}/dist"
# Remove unnecessary files (removes ~90% of extracted size).
	cd docs && $(RM) $@*.map $@*-es-bundle*.js $@swagger-ui.js
# Insert location of OpenAPI spec to be rendered.
	$(SED) -i 's|url: "[^"]*",|url: "{{.EndpointBaseURL}}/docs/swagger/service.yaml",|g' "docs/swagger-initializer.js"
# Add goembed FS


clean:
	$(RM) -r docs swagger-ui-*.tar.gz

