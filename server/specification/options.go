package specification

type Config struct {
	swaggerUIDirectory string
	indexFilePath      string
	openAPIPath        string
}

type Option func(*Config)

// Sets another base path for the swagger endpoints
func WithSwaggerUIDirectory(swaggerUIDirectory string) Option {
	return func(config *Config) {
		config.swaggerUIDirectory = swaggerUIDirectory
	}
}

// Sets another path to where files for swagger ui is located
func WithIndexFilePath(indexFilePath string) Option {
	return func(config *Config) {
		config.indexFilePath = indexFilePath
	}
}

// Sets another name for index file than index.html
func WithOpenAPIPath(openAPIPath string) Option {
	return func(config *Config) {
		config.openAPIPath = openAPIPath
	}
}
