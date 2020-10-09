package proxanne

import (
	"errors"
	"fmt"
	"net/http/httputil"
	"net/url"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	ErrNoRoutes = errors.New("no routes found")
)

type Config struct {
	Routes []struct {
		Matches string `yaml:"matches"`
		Target  string `yaml:"target"`
	} `yaml:"routes"`
}

func ParseConfig(raw []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(raw, config)
	return config, err
}

func BuildRouter(config *Config) (Router, error) {
	if len(config.Routes) == 0 {
		return nil, ErrNoRoutes
	}

	router := Router{}

	for _, route := range config.Routes {
		matchesRegexp, err := regexp.Compile(route.Matches)
		if err != nil {
			return nil, fmt.Errorf("error in route %v: matches is not a valid regex", route)
		}

		targetURL, err := url.Parse(route.Target)
		if err != nil {
			return nil, fmt.Errorf("error in route %v: target is not a valid url", route)
		}

		router = append(
			router,
			&Route{
				Matches: matchesRegexp,
				Target:  httputil.NewSingleHostReverseProxy(targetURL),
			},
		)
	}

	return router, nil
}
