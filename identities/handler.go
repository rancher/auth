package identities

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/rancher/types/config"
)

//identityAPIHandler is a wrapper over the mux router serving /v3/identities API
type identityAPIHandler struct {
	identityRouter http.Handler
}

func (h *identityAPIHandler) getRouter() http.Handler {
	return h.identityRouter
}

func (h *identityAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.getRouter().ServeHTTP(w, r)
}

func NewIdentityAPIHandler(ctx context.Context, mgmtCtx *config.ManagementContext) (http.Handler, error) {
	router, err := newIdentityRouter(ctx, mgmtCtx)
	if err != nil {
		return nil, err
	}
	return &identityAPIHandler{identityRouter: router}, nil
}

//newIdentityRouter creates and configures a mux router for /v3/identities APIs
func newIdentityRouter(ctx context.Context, mgmtCtx *config.ManagementContext) (*mux.Router, error) {
	apiServer, err := newIdentityAPIServer(ctx, mgmtCtx)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter().StrictSlash(true)
	// Application routes
	router.Methods("GET").Path("/v3/identities").Handler(http.HandlerFunc(apiServer.listIdentities))
	router.Methods("POST").Path("/v3/identities").Queries("action", "search").Handler(http.HandlerFunc(apiServer.searchIdentities))

	return router, nil
}
