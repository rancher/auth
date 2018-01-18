package setup

import (
	"context"

	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"

	"github.com/rancher/auth/authconfig"
)

var (
	crdVersions = []*types.APIVersion{
		&managementSchema.Version,
	}
)

func Schemas(ctx context.Context, management *config.ManagementContext, schemas *types.Schemas) error {

	GithubConfig(schemas)

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

func GithubConfig(schemas *types.Schemas) {
	schema := schemas.Schema(&managementSchema.Version, client.GithubConfigType)
	schema.Formatter = authconfig.GithubConfigFormatter
	schema.ActionHandler = authconfig.GithubConfigActionHandler
}
