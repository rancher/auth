package identities

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"

	"github.com/rancher/auth/providers"
)

type identityAPIServer struct {
	ctx              context.Context
	client           *config.ManagementContext
	identitiesClient v3.IdentityInterface
	tokensClient     v3.TokenInterface
}

func newIdentityAPIServer(ctx context.Context, mgmtCtx *config.ManagementContext) (*identityAPIServer, error) {
	if mgmtCtx == nil {
		return nil, fmt.Errorf("Failed to build tokenAPIHandler, nil ManagementContext")
	}
	providers.Configure(ctx, mgmtCtx)

	apiServer := &identityAPIServer{
		ctx:              ctx,
		client:           mgmtCtx,
		identitiesClient: mgmtCtx.Management.Identities(""),
		tokensClient:     mgmtCtx.Management.Tokens(""),
	}
	return apiServer, nil
}

func (s *identityAPIServer) getIdentities(tokenKey string) ([]v3.Identity, int, error) {
	identities := make([]v3.Identity, 0)

	logrus.Debugf("getIdentities: token cookie: %v", tokenKey)

	token, err := s.getTokenCR(tokenKey)

	if err != nil {
		return identities, 401, err
	}

	//add code to make sure token is valid
	identities = append(identities, token.UserIdentity)
	identities = append(identities, token.GroupIdentities...)

	return identities, 0, nil

}

func (s *identityAPIServer) findIdentities(tokenKey string, name string) ([]v3.Identity, int, error) {
	var identities []v3.Identity
	var status int
	logrus.Debugf("searchIdentities: token cookie: %v, name: %v", tokenKey, name)

	token, err := s.getTokenCR(tokenKey)
	if err != nil {
		return identities, 401, err
	}
	identities, status, err = providers.SearchIdentities(name, *token)

	return identities, status, err
}

func (s *identityAPIServer) getTokenCR(tokenID string) (*v3.Token, error) {
	if s.client != nil {
		storedToken, err := s.tokensClient.Get(strings.ToLower(tokenID), metav1.GetOptions{})

		if err != nil {
			logrus.Info("Failed to get token resource: %v", err)
			return nil, fmt.Errorf("Failed to retrieve auth token")
		}

		logrus.Debugf("storedToken token resource: %v", storedToken)

		return storedToken, nil
	}
	return nil, fmt.Errorf("No k8s Client configured")
}
