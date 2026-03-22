package product_dotnevx

import (
	"sync"

	proudct_config "github.com/premwitthawas/demo_ecommerce_api/internals/product/model/config"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port/config"
	pkgs_env "github.com/premwitthawas/demo_ecommerce_api/pkgs/env"
)

type envAdapter struct {
	*proudct_config.ProductConfig
}

func (e *envAdapter) GetAPPConfig() *proudct_config.ProductAppConfig {
	return e.App
}

func (e *envAdapter) GetDBConfig() *proudct_config.ProductDBConfig {
	return e.DB
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
			ProductConfig: &proudct_config.ProductConfig{
				App: &proudct_config.ProductAppConfig{
					Address:        pkgs_env.GetEnvString("APP_ADDRESS", ":5002"),
					Name:           pkgs_env.GetEnvString("APP_NAME", "product"),
					Mode:           pkgs_env.GetEnvString("APP_MODE", "dev"),
					KafkaAddresses: pkgs_env.GetEnvString("APP_KAFKA_ADREESES", "localhost:9092"),
				},
				DB: &proudct_config.ProductDBConfig{
					DatabaseURL: pkgs_env.GetEnvString("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:6432/product?sslmode=disable"),
				},
			},
		}
	})
	return cfg
}
