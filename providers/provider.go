package providers

import (
	"context"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"

	"github.com/rancher/auth/providers/local"
)

//Providers map
var providers map[string]PrincipalProvider
var providerOrderList []string

func init() {
	providerOrderList = []string{"local"}
	providers = make(map[string]PrincipalProvider)
}

//PrincipalProvider interfacse defines what methods an identity provider should implement
type PrincipalProvider interface {
	GetName() string
	AuthenticateUser(jsonInput v3.LoginInput) (v3.Principal, []v3.Principal, int, error)
	SearchPrincipals(name string, myToken v3.Token) ([]v3.Principal, int, error)
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

func AuthenticateUser(jsonInput v3.LoginInput) (v3.Principal, []v3.Principal, int, error) {
	var groupPrincipals []v3.Principal
	var userPrincipal v3.Principal
	var status int
	var err error

	for _, providerName := range providerOrderList {
		switch providerName {
		case "local":
			userPrincipal, groupPrincipals, status, err = providers[providerName].AuthenticateUser(jsonInput)
		}
	}
	return userPrincipal, groupPrincipals, status, err
}

func SearchPrincipals(name string, myToken v3.Token) ([]v3.Principal, int, error) {
	principals := make([]v3.Principal, 0)
	var status int
	var err error

	for _, providerName := range providerOrderList {
		switch providerName {
		case "local":
			localprincipals, status, err := providers[providerName].SearchPrincipals(name, myToken)
			if err != nil {
				return localprincipals, status, err
			}
			principals = append(principals, localprincipals...)
		}
	}
	return principals, status, err
}
