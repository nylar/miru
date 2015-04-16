package miru

import "github.com/BurntSushi/toml"

// DefaultConfig is used if no config.toml file is found, sets the config to
// acceptable defaults.
var DefaultConfig = `
[database]
host = "localhost:28015"
name = "miru"

[tables]
index = "indexes"
document = "documents"

[api]
port = "8036"
`

// Config holds configuration information regarding the database and the port in
// which to serve on.
type Config struct {
	Database database
	Tables   tables
	Api      api
}

type database struct {
	Host string
	Name string
}

type tables struct {
	Index    string
	Document string
}

type api struct {
	Port string
}

// LoadConfig loads configuration data into the Config struct.
func LoadConfig(data string) (*Config, error) {
	var conf Config
	if _, err := toml.Decode(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
