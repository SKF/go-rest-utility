package specification

type Config struct {
	swaggerUIDirectory string
	indexFilePath      string
	openAPIPath        string
	swaggerInit        *string
}

type Option func(*Config)

// WithSwaggerUIDirectory Sets another base path for the swagger endpoints
func WithSwaggerUIDirectory(swaggerUIDirectory string) Option {
	return func(config *Config) {
		config.swaggerUIDirectory = swaggerUIDirectory
	}
}

// WithIndexFilePath Sets another path to where files for swagger ui is located
func WithIndexFilePath(indexFilePath string) Option {
	return func(config *Config) {
		config.indexFilePath = indexFilePath
	}
}

// WithOpenAPIPath Sets another name for index file than index.html
func WithOpenAPIPath(openAPIPath string) Option {
	return func(config *Config) {
		config.openAPIPath = openAPIPath
	}
}

// WithSwaggerInitJSPath Sets a path for swagger-initializer.js. This is required for using swagger >= 4.0.0
func WithSwaggerInitJSPath(swaggerInitJSPath string) Option {
	return func(config *Config) {
		config.swaggerInit = &swaggerInitJSPath
	}
}
