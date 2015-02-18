package config

import "github.com/BurntSushi/toml"

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

func LoadConfig(data string) (*Config, error) {
	var conf Config
	if _, err := toml.Decode(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
