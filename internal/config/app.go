package config

import (
	"kasir-api/internal/delivery/http/route"
	"net/http"
)

type BootstrapConfig struct {
	Server *http.ServeMux
	Config *Config
}

func Bootstrap(config *BootstrapConfig) {

	routeConfig := route.RouteConfig{
		Server: config.Server,
	}

	routeConfig.Setup()
}
