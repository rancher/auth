package identities

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	
	normanapi "github.com/rancher/norman/api"
	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	//"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

/*identityAPIHandler is a wrapper over the mux router serving /v3/identities API
type identityAPIHandler struct {
	identityRouter http.Handler
}

func (h *identityAPIHandler) getRouter() http.Handler {
	return h.identityRouter
}

func (h *identityAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.getRouter().ServeHTTP(w, r)
}

/*func NewIdentityAPIHandler(ctx context.Context, mgmtCtx *config.ManagementContext) (http.Handler, error) {
	router, err := newIdentityRouter(ctx, mgmtCtx)
	if err != nil {
		return nil, err
	}
	return &identityAPIHandler{identityRouter: router}, nil
}

/*newIdentityRouter creates and configures a mux router for /v3/identities APIs
func newIdentityRouter(ctx context.Context, mgmtCtx *config.ManagementContext) (*mux.Router, error) {
	apiServer, err := newIdentityAPIServer(ctx, mgmtCtx)
	if err != nil {
		return nil, err
	}

	newSchemas := factory.Schemas(&managementSchema.Version)
	
	schemas := types.NewSchemas().
		AddSchemas(newSchemas)
		

	if err := setupSchemas(ctx, mgmtCtx, schemas, apiServer); err != nil {
		return nil, err
	}

	server := normanapi.NewAPIServer()

	if err := server.AddSchemas(schemas); err != nil {
		return nil, err
	}

	router := mux.NewRouter().StrictSlash(true)
	// Application routes
	router.Methods("GET").Path("/v3/identities").Handler(http.HandlerFunc(apiServer.listIdentities))
	router.Methods("POST").Path("/v3/identities").Queries("action", "search").Handler(http.HandlerFunc(apiServer.searchIdentities))

	router.Handle("/v3/identities", server).Methods("GET")

	return router, nil
}*/

func NewIdentityAPIHandler(ctx context.Context, mgmtCtx *config.ManagementContext) (http.Handler, error) {
	
	apiServer, err := newIdentityAPIServer(ctx, mgmtCtx)
	if err != nil {
		return nil, err
	}
		
	schemas := types.NewSchemas() //.AddSchemas(managementSchema.Schemas)

	if err := setupSchemas(ctx, mgmtCtx, schemas, apiServer); err != nil {
		return nil, err
	}

	server := normanapi.NewAPIServer()

	if err := server.AddSchemas(schemas); err != nil {
		return nil, err
	}

	return server, nil
}

var crdVersions = []*types.APIVersion{
	&managementSchema.Version,
}

func setupSchemas(ctx context.Context, management *config.ManagementContext, schemas *types.Schemas, apiServer *identityAPIServer) error {

	schemas.MustImportAndCustomize(&managementSchema.Version, v3.Identity{}, func(schema *types.Schema) {
		schema.CollectionMethods = []string{http.MethodGet}
		//schema.CollectionActions = []string{http.MethodGet}
		schema.ResourceMethods = []string{http.MethodGet}
		schema.ListHandler = apiServer.handleListIdentities
		schema.PluralName = "identities"
	})
	
	//schema := schemas.Schema(&managementSchema.Version, v3.IdentityType)
	//schema.ListHandler = apiServer.handleListIdentities
	
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

func CustomHandler(apiContext *types.APIContext) error {
	fmt.Println("Handler called")
	logrus.Infof("Handler called %v", apiContext)
	/*err := handler(apiContext)
	if err != nil {
		logrus.Errorf("Error during subscribe %v", err)
	}*/
	return nil
}
