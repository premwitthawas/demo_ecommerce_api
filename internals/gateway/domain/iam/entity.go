package iam

import pkg_debug "github.com/premwitthawas/demo_ecommerce_api/pkgs/debug"

type Claims struct {
	// jwt.Claims
	Subject        string                    `json:"sub"`
	Email          string                    `json:"email"`
	FamilyName     string                    `json:"family_name"`
	GivenName      string                    `json:"given_name"`
	Name           string                    `json:"name"`
	RealmAccess    map[string]any            `json:"realm_access"`
	ResourceAccess map[string]map[string]any `json:"resource_access"`
	Issuer         string                    `json:"iss"`
	Audience       any                       `json:"aud"`
	Scrope         string                    `json:"scope"`
}

func NewCliamIAM(claims *Claims) (*Claims, error) {
	if claims.Subject == "" {
		return nil, ErrIAMSubjectIsEmpty
	}
	return claims, nil
}

func (e *Claims) IsRolePermisison(role string, key string) bool {
	pkg_debug.Debug(e.ResourceAccess[key])
	return true
}
