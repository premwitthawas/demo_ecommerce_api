package config

type Config interface {
	IsProduction() bool
	GetAppAddress() string
	GetAppName() string
	GetMode() string
}
