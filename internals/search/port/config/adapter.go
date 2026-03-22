package search

import config "github.com/premwitthawas/demo_ecommerce_api/internals/search/model/config"

type Config interface {
	IsProduction() bool
	GetAPPConfig() *config.SearchAppConfig
}
