package authenticator

import (
	"context"
	"net/http"

	"github.com/rancher/types/config"
)

type Authenticator interface {
	Authenticate(req *http.Request) (authed bool, user string, groups []string, err error)
}

func NewAuthenticator(ctx context.Context, mgmtCtx *config.ManagementContext) Authenticator {
	return &tokenAuthenticator{
		ctx:          ctx,
		client:       mgmtCtx,
		tokensClient: mgmtCtx.Management.Tokens(""),
	}
}
