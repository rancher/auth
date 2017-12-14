package authenticator

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

type tokenAuthenticator struct {
	ctx          context.Context
	client       *config.ManagementContext
	tokensClient v3.TokenInterface
}

func (a *tokenAuthenticator) Authenticate(req *http.Request) (bool, string, []string, error) {
	var user string
	var groups []string

	cookie, err := req.Cookie("rAuthnSessionToken")
	if err != nil {
		return false, user, groups, fmt.Errorf("Failed to find auth cookie")
	}
	logrus.Debugf("Authenticate: token cookie: %v %v", cookie.Name, cookie.Value)

	token, err := a.getTokenCR(cookie.Value)
	if err != nil {
		return false, user, groups, err
	}

	user = token.UserIdentity.LoginName

	for _, groupIdentity := range token.GroupIdentities {
		groups = append(groups, groupIdentity.Name)
	}

	return true, user, groups, nil
}

func (a *tokenAuthenticator) getTokenCR(tokenID string) (*v3.Token, error) {
	if a.client != nil {
		storedToken, err := a.tokensClient.Get(strings.ToLower(tokenID), metav1.GetOptions{})

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve auth token, error: %v", err)
		}

		logrus.Debugf("storedToken token resource: %v", storedToken)

		return storedToken, nil
	}
	return nil, fmt.Errorf("No k8s Client configured")
}
