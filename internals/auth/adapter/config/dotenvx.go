package dotnevx

import (
	"sync"

	domain "github.com/premwitthawas/demo_ecommerce_api/internals/auth/domain/config"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/auth/port/config"
	pkgs_env "github.com/premwitthawas/demo_ecommerce_api/pkgs/env"
)

type envAdapter struct {
	*domain.AuthConfig
}

func (e *envAdapter) GetMode() string {
	return e.App.Mode
}

func (e *envAdapter) GetAppAddress() string {
	return e.App.Address
}

func (e *envAdapter) GetAppName() string {
	return e.App.Name
}

func (e *envAdapter) IsProduction() bool {
	return e.App.Mode == "prod"
}

var (
	once sync.Once
	cfg  *envAdapter
)

func NewConfig() port.Config {
	once.Do(func() {
		cfg = &envAdapter{
			AuthConfig: &domain.AuthConfig{
				App: &domain.AuthAppConfig{
					Address: pkgs_env.GetEnvString("APP_ADDRESS", ":5001"),
					Name:    pkgs_env.GetEnvString("APP_NAME", "auth"),
					Mode:    pkgs_env.GetEnvString("APP_MODE", "dev"),
				},
			},
		}
	})
	return cfg
}
