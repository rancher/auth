package authconfig

import (
	//"strings"

	"github.com/rancher/norman/store/proxy"
	//"github.com/rancher/norman/store/transform"
	"github.com/rancher/norman/types"
	//"github.com/rancher/norman/types/convert"
	"k8s.io/client-go/rest"
)

type Store struct {
	types.Store
}

func NewAuthConfigStore(k8sClient rest.Interface, schema *types.Schema) types.Store {
	return proxy.NewProxyStore(k8sClient,
		[]string{"apis"},
		"management.cattle.io",
		"v3",
		"AuthConfig",
		"authconfigs")
}
