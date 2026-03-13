package config

import "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/config"

type GatewayConfigAdapter interface {
	GetConfigAuthService() *config.AuthServiceConfig
	GetAppConfig() *config.GatewayAppConfig
	GetKeycloakConfig() *config.KeyclaokAppConfig
}
