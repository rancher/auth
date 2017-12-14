package providers

import (
	"context"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"

	"github.com/rancher/auth/providers/local"
)

//Providers map
var Providers map[string]IdentityProvider
var ProviderOrderList []string

func init() {
	ProviderOrderList = []string{"local"}
	Providers = make(map[string]IdentityProvider)
}

//IdentityProvider interfacse defines what methods an identity provider should implement
type IdentityProvider interface {
	GetName() string
	AuthenticateUser(jsonInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error)
}

func Configure(ctx context.Context, mgmtCtx *config.ManagementContext) {
	for _, name := range ProviderOrderList {
		switch name {
		case "local":
			Providers[name] = local.Configure(ctx, mgmtCtx)
		}
	}
}

func AuthenticateUser(jsonInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error) {
	var groupIdentities []v3.Identity
	var userIdentity v3.Identity
	var status int
	var err error

	for _, name := range ProviderOrderList {
		switch name {
		case "local":
			userIdentity, groupIdentities, status, err = Providers[name].AuthenticateUser(jsonInput)
		}
	}
	return userIdentity, groupIdentities, status, err
}
