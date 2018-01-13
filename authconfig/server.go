package authconfig

import (
	"context"
	"net/http"
	"github.com/sirupsen/logrus"

	normanapi "github.com/rancher/norman/api"
	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"

)

var crdVersions = []*types.APIVersion{
	&managementSchema.Version,
}

func New(ctx context.Context, mgmtCtx *config.ManagementContext) (http.Handler, error) {
	schemas := types.NewSchemas().
		AddSchemas(managementSchema.Schemas)

	if err := setupSchemas(ctx, mgmtCtx, schemas); err != nil {
		return nil, err
	}

	server := normanapi.NewAPIServer()

	if err := server.AddSchemas(schemas); err != nil {
		return nil, err
	}

	return server, nil
}

func setupSchemas(ctx context.Context, management *config.ManagementContext, schemas *types.Schemas) error {

	schema := schemas.Schema(&managementSchema.Version, client.GithubConfigType)
	schema.Formatter = GithubConfigFormatter
	schema.ActionHandler = CustomHandler


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

func CustomHandler(actionName string, action *types.Action, request *types.APIContext) error {
	logrus.Infof("CustomHandler called %v", request)
	return nil
}

func GithubConfigFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.Actions["configureTest"] = apiContext.URLBuilder.Action("configureTest", resource)
	resource.Actions["testAndApply"] = apiContext.URLBuilder.Action("testAndApply", resource)
}


