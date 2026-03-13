package dotenvx

import (
	"sync"

	domain "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/config"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/config"
	pkgs_env "github.com/premwitthawas/demo_ecommerce_api/pkgs/env"
)

type gatewayConfig struct {
	*domain.GatewayConfig
}

func (g *gatewayConfig) GetKeycloakConfig() *domain.KeyclaokAppConfig {
	return g.Keycloak
}

func (g *gatewayConfig) GetAppConfig() *domain.GatewayAppConfig {
	return g.App
}

func (g *gatewayConfig) GetConfigAuthService() *domain.AuthServiceConfig {
	return g.Auth
}

var (
	once sync.Once
	cfg  *gatewayConfig
)

func NewGatewayConfig() port.GatewayConfigAdapter {
	once.Do(func() {
		cfg = &gatewayConfig{
			GatewayConfig: &domain.GatewayConfig{
				App: &domain.GatewayAppConfig{
					Address: pkgs_env.GetEnvString("APP_ADDRESS", "127.0.0.1:6001"),
					Name:    pkgs_env.GetEnvString("APP_NAME", "gateway"),
				},
				Auth: &domain.AuthServiceConfig{
					Url: pkgs_env.GetEnvString("AUTH_SVC_BASE_URL", "http://127.0.0.1:5001"),
				},
				Keycloak: &domain.KeyclaokAppConfig{
					ClientID:     pkgs_env.GetEnvString("KEYCLOAK_CLIENT_ID", "gateway-service"),
					ClientSecret: pkgs_env.GetEnvString("KEYCLOAK_CLIENT_SECRET", "DJXNLvfeBJsF9G2ZwHfdBtfo3xmQDWVu"),
					Issuer:       pkgs_env.GetEnvString("KEYCLOAK_ISSUER", "http://localhost:5000/realms/master"),
					RealmeName:   pkgs_env.GetEnvString("KEYCLOAK_REALMS", "master"),
					Redirect:     pkgs_env.GetEnvString("KEYCLOAK_REDIRECT", "http://localhost:6001"),
				},
			},
		}
	})
	return cfg
}
