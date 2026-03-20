package config

import config "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/config"

type Config interface {
	IsProduction() bool
	GetAPPConfig() *config.ProductAppConfig
	GetDBConfig() *config.ProductDBConfig
}
