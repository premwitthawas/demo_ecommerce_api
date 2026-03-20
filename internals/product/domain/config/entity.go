package config

type ProductAppConfig struct {
	Address        string
	Name           string
	Mode           string
	KafkaAddresses string
}

type ProductDBConfig struct {
	DatabaseURL string
}

type ProductConfig struct {
	App *ProductAppConfig
	DB  *ProductDBConfig
}
