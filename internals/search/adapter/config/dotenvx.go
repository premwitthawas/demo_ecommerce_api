package search_dotnevx

import (
	"sync"

	config "github.com/premwitthawas/demo_ecommerce_api/internals/search/model/config"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/search/port/config"
	pkgs_env "github.com/premwitthawas/demo_ecommerce_api/pkgs/env"
)

type envAdapter struct {
	*config.SearchConfig
}

func (e *envAdapter) GetAPPConfig() *config.SearchAppConfig {
	return e.App
}

// func (e *envAdapter) GetDBConfig() *domain.ProductDBConfig {
// 	return e.DB
// }

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
			SearchConfig: &config.SearchConfig{
				App: &config.SearchAppConfig{
					Address:               pkgs_env.GetEnvString("APP_ADDRESS", ":5003"),
					Name:                  pkgs_env.GetEnvString("APP_NAME", "search-service"),
					Mode:                  pkgs_env.GetEnvString("APP_MODE", "dev"),
					KafkaAddress:          pkgs_env.GetEnvString("APP_KAFKA_ADREESES", "localhost:9092"),
					ElasticsearchAddress:  pkgs_env.GetEnvString("APP_ELASTICSEARC_ADDRESS", "APP_ELASTICSEARC_ADDRESS"),
					ElasticsearchUsername: pkgs_env.GetEnvString("APP_ELASTICSEARC_USERNAME", "APP_ELASTICSEARC_USERNAME"),
					ElasticsearchPassword: pkgs_env.GetEnvString("APP_ELASTICSEARC_PASSWORD", "APP_ELASTICSEARC_PASSWORD"),
				},
			},
		}
	})
	return cfg
}
