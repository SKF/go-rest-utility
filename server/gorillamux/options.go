package gorillamux

type Config struct {
	swaggerPath string
	docsPath    string
	indexFile   string
}

type Option func(*Config)

// Sets another base path for the swagger endpoints
func WithSwaggerPath(swaggerPath string) Option {
	return func(config *Config) {
		config.swaggerPath = swaggerPath
	}
}

// Sets another path to where files for swagger ui is located
func WithDocsPath(docsPath string) Option {
	return func(config *Config) {
		config.docsPath = docsPath
	}
}

// Sets another name for index file than index.html
func WithIndexFile(indexFile string) Option {
	return func(config *Config) {
		config.indexFile = indexFile
	}
}
