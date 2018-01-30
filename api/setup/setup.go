package setup

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/store/subtype"
	"github.com/rancher/norman/types"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"

	"github.com/rancher/auth/authconfig"
	"github.com/rancher/auth/tokens"
)

var (
	crdVersions = []*types.APIVersion{
		&managementSchema.Version,
	}
)

func Schemas(ctx context.Context, management *config.ManagementContext, schemas *types.Schemas) error {
	Token(schemas)
	GithubConfig(schemas)
	//subtyping adventures:
	//AuthConfig(management.UnversionedClient, schemas)

	crdStore, err := crd.NewCRDStoreFromConfig(management.RESTConfig)
	if err != nil {
		return err
	}

	var crdSchemas []*types.Schema
	for _, version := range crdVersions {
		for _, schema := range schemas.SchemasForVersion(*version) {
			crdSchemas = append(crdSchemas, schema)
		}
	}

	return crdStore.AddSchemas(ctx, crdSchemas...)
}

func AuthConfig(k8sClient rest.Interface, schemas *types.Schemas) {
	schema := schemas.Schema(&managementSchema.Version, client.AuthConfigType)
	schema.Store = authconfig.NewAuthConfigStore(k8sClient, schema)

	/*TODO: If kind is githubConfig we need these resourceActions, but for ADConfig the inputs to these actions are different

		schema.ResourceActions = map[string]types.Action{
		"configureTest": {
			Input:  "githubConfigTestInput",
			Output: "githubConfig",
		},
		"testAndApply": {
			Input:  "githubConfigApplyInput",
			Output: "githubConfig",
		},
	}
	*/

	for _, subSchema := range schemas.Schemas() {

		if subSchema.BaseType == "AuthConfig" && subSchema.ID == client.GithubConfigType {
			logrus.Infof("Found subtype %v", subSchema.ID)
			logrus.Infof("Assigning schema.Store %v", schema.Store)

			subSchema.Store = subtype.NewSubTypeStore(subSchema.ID, schema.Store)
			subSchema.Formatter = authconfig.GithubConfigFormatter
			subSchema.ActionHandler = authconfig.GithubConfigActionHandler

			schema.ResourceActions = subSchema.ResourceActions
			schema.ActionHandler = subSchema.ActionHandler
			schema.Formatter = subSchema.Formatter

		}
	}

}

func GithubConfig(schemas *types.Schemas) {
	schema := schemas.Schema(&managementSchema.Version, client.GithubConfigType)
	schema.Formatter = authconfig.GithubConfigFormatter
	schema.ActionHandler = authconfig.GithubConfigActionHandler

}

func Token(schemas *types.Schemas) {
	schema := schemas.Schema(&managementSchema.Version, client.TokenType)
	schema.CollectionActions = map[string]types.Action{
		"login": {
			Input:  "loginInput",
			Output: "token",
		},
		"logout": {},
	}
	schema.ActionHandler = tokens.TokenActionHandler
	schema.ListHandler = tokens.TokenListHandler
	schema.CreateHandler = tokens.TokenCreateHandler
	schema.DeleteHandler = tokens.TokenDeleteHandler
}
