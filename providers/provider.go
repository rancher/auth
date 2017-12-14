package providers

import (
	"context"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"

	"github.com/rancher/auth/providers/local"
)

//Providers map
var providers map[string]IdentityProvider
var providerOrderList []string

func init() {
	providerOrderList = []string{"local"}
	providers = make(map[string]IdentityProvider)
}

//IdentityProvider interfacse defines what methods an identity provider should implement
type IdentityProvider interface {
	GetName() string
	AuthenticateUser(jsonInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error)
	SearchIdentities(name string, myToken v3.Token) ([]v3.Identity, int, error)
}

func Configure(ctx context.Context, mgmtCtx *config.ManagementContext) {
	for _, providerName := range providerOrderList {
		if _, exists := providers[providerName]; !exists {
			switch providerName {
			case "local":
				providers[providerName] = local.Configure(ctx, mgmtCtx)
			}
		}
	}
}

func AuthenticateUser(jsonInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error) {
	var groupIdentities []v3.Identity
	var userIdentity v3.Identity
	var status int
	var err error

	for _, providerName := range providerOrderList {
		switch providerName {
		case "local":
			userIdentity, groupIdentities, status, err = providers[providerName].AuthenticateUser(jsonInput)
		}
	}
	return userIdentity, groupIdentities, status, err
}

func SearchIdentities(name string, myToken v3.Token) ([]v3.Identity, int, error) {
	identities := make([]v3.Identity, 0)
	var status int
	var err error

	for _, providerName := range providerOrderList {
		switch providerName {
		case "local":
			localIdentities, status, err := providers[providerName].SearchIdentities(name, myToken)
			if err != nil {
				return localIdentities, status, err
			}
			identities = append(identities, localIdentities...)
		}
	}
	return identities, status, err
}
