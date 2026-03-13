package config

type AuthAppConfig struct {
	Address string
	Name    string
	Mode    string
}

type AuthConfig struct {
	App *AuthAppConfig
}
