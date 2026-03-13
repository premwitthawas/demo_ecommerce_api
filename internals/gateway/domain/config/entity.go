package config

type GatewayConfig struct {
	App      *GatewayAppConfig
	Auth     *AuthServiceConfig
	Keycloak *KeyclaokAppConfig
}

type GatewayAppConfig struct {
	Address string
	Name    string
}

type KeyclaokAppConfig struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	RealmeName   string
	Redirect     string
}

type AuthServiceConfig struct {
	Url string
}
