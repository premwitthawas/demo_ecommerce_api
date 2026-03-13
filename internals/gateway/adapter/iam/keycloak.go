package keycloak

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	domain "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/iam"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/config"
	"github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/iam"
	"golang.org/x/oauth2"
)

type keycloakAdapter struct {
	cfg      port.GatewayConfigAdapter
	verifier *oidc.IDTokenVerifier
}

func (k *keycloakAdapter) ValidateToken(c context.Context, token string) (*domain.Claims, error) {
	if k.verifier == nil {
		return nil, errors.New("oidc verifier not intitialzed")
	}
	idToken, err := k.verifier.Verify(c, token)
	if err != nil {
		return nil, err
	}
	var claims domain.Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func NewKeycloakAdapter(ctx context.Context, cfg port.GatewayConfigAdapter) (iam.IAMAapter, error) {
	adapter := &keycloakAdapter{cfg: cfg}
	if err := adapter.CreateVerifierToken(ctx); err != nil {
		return nil, err
	}
	return adapter, nil
}

func (k *keycloakAdapter) CreateProvider(ctx context.Context) (*oidc.Provider, *oauth2.Config, error) {
	provider, err := oidc.NewProvider(ctx, k.cfg.GetKeycloakConfig().Issuer)
	if err != nil {
		return nil, nil, err
	}
	oauthConfig := &oauth2.Config{
		ClientID:     k.cfg.GetKeycloakConfig().ClientID,
		ClientSecret: k.cfg.GetKeycloakConfig().ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  k.cfg.GetKeycloakConfig().Redirect,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return provider, oauthConfig, nil
}

func (k *keycloakAdapter) CreateVerifierToken(ctx context.Context) error {
	provider, _, err := k.CreateProvider(ctx)
	if err != nil {
		return err
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID: k.cfg.GetKeycloakConfig().ClientID,
	})
	k.verifier = verifier
	return nil
}
