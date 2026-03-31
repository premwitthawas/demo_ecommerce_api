package search_config

type SearchAppConfig struct {
	Address               string
	Name                  string
	Mode                  string
	KafkaAddress          string
	ElasticsearchAddress  string
	ElasticsearchUsername string
	ElasticsearchPassword string
	OtelURL               string
}

// type SearchDBConfig struct {
// 	DatabaseURL string
// }

type SearchConfig struct {
	App *SearchAppConfig
	// DB  *ProductDBConfig
}
