package tokens

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/rancher/auth/providers"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

const (
	defaultSecret          = "secret"
	defaultTokenTTL        = "57600000"
	defaultRefreshTokenTTL = "7200000"
)

type tokenAPIServer struct {
	ctx          context.Context
	client       *config.ManagementContext
	tokensClient v3.TokenInterface
}

func newTokenAPIServer(ctx context.Context, mgmtCtx *config.ManagementContext) (*tokenAPIServer, error) {
	if mgmtCtx == nil {
		return nil, fmt.Errorf("Failed to build tokenAPIHandler, nil ManagementContext")
	}
	providers.Configure(ctx, mgmtCtx)

	apiServer := &tokenAPIServer{
		ctx:          ctx,
		client:       mgmtCtx,
		tokensClient: mgmtCtx.Management.Tokens(""),
	}

	return apiServer, nil
}

//createLoginToken will authenticate with provider and creates a token CR
func (s *tokenAPIServer) createLoginToken(jsonInput v3.LoginInput) (v3.Token, int, error) {

	logrus.Info("Create Token Invoked %v", jsonInput)

	// Authenticate User
	userIdentity, groupIdentities, status, err := providers.AuthenticateUser(jsonInput)

	if status == 0 && err == nil {
		logrus.Info("User Authenticated")
		if s.client != nil {
			key, err := generateKey()
			if err != nil {
				logrus.Info("Failed to generate token key: %v", err)
				return v3.Token{}, 0, fmt.Errorf("Failed to generate token key")
			}

			//check that there is no token with this key?

			ttl := jsonInput.TTLMillis
			refreshTTL := jsonInput.IdentityRefreshTTLMillis
			if ttl == "" {
				ttl = defaultTokenTTL               //16 hrs
				refreshTTL = defaultRefreshTokenTTL //2 hrs
			}

			k8sToken := &v3.Token{
				TokenID:                  key,
				UserIdentity:             userIdentity,
				GroupIdentities:          groupIdentities,
				IsDerived:                false,
				TTLMillis:                ttl,
				IdentityRefreshTTLMillis: refreshTTL,
				User:         userIdentity.LoginName,
				ExternalID:   userIdentity.ObjectMeta.Name,
				AuthProvider: getAuthProviderName(userIdentity.ObjectMeta.Name),
			}
			rToken, err := s.createK8sTokenCR(k8sToken)
			return rToken, 0, err
		}
		logrus.Info("Client nil %v", s.client)
		return v3.Token{}, 500, fmt.Errorf("No k8s Client configured")
	}

	return v3.Token{}, status, err
}

//CreateDerivedToken will create a jwt token for the authenticated user
func (s *tokenAPIServer) createDerivedToken(jsonInput v3.Token, tokenID string) (v3.Token, int, error) {

	logrus.Info("Create Derived Token Invoked")

	token, err := s.getK8sTokenCR(tokenID)

	if err != nil {
		return v3.Token{}, 401, err
	}

	if s.client != nil {
		key, err := generateKey()
		if err != nil {
			logrus.Info("Failed to generate token key: %v", err)
			return v3.Token{}, 0, fmt.Errorf("Failed to generate token key")
		}

		ttl := jsonInput.TTLMillis
		refreshTTL := jsonInput.IdentityRefreshTTLMillis
		if ttl == "" {
			ttl = defaultTokenTTL               //16 hrs
			refreshTTL = defaultRefreshTokenTTL //2 hrs
		}

		k8sToken := &v3.Token{
			TokenID:                  key,
			UserIdentity:             token.UserIdentity,
			GroupIdentities:          token.GroupIdentities,
			IsDerived:                true,
			TTLMillis:                ttl,
			IdentityRefreshTTLMillis: refreshTTL,
			User:         token.User,
			ExternalID:   token.ExternalID,
			AuthProvider: token.AuthProvider,
		}
		rToken, err := s.createK8sTokenCR(k8sToken)
		return rToken, 0, err

	}
	logrus.Info("Client nil %v", s.client)
	return v3.Token{}, 500, fmt.Errorf("No k8s Client configured")

}

func (s *tokenAPIServer) createK8sTokenCR(k8sToken *v3.Token) (v3.Token, error) {
	if s.client != nil {

		labels := make(map[string]string)
		labels["io.cattle.token.field.externalID"] = k8sToken.ExternalID

		k8sToken.APIVersion = "management.cattle.io/v3"
		k8sToken.Kind = "Token"
		k8sToken.ObjectMeta = metav1.ObjectMeta{
			Name:   strings.ToLower(k8sToken.TokenID),
			Labels: labels,
		}
		createdToken, err := s.tokensClient.Create(k8sToken)

		if err != nil {
			logrus.Info("Failed to create token resource: %v", err)
			return v3.Token{}, err
		}
		logrus.Info("Created Token %v", createdToken)
		return *createdToken, nil
	}

	return v3.Token{}, fmt.Errorf("No k8s Client configured")
}

func (s *tokenAPIServer) getK8sTokenCR(tokenID string) (*v3.Token, error) {
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

//GetTokens will list all tokens of the authenticated user - login and derived
func (s *tokenAPIServer) getTokens(tokenID string) ([]v3.Token, int, error) {
	logrus.Info("GET Token Invoked")
	tokens := make([]v3.Token, 0)

	if s.client != nil {
		storedToken, err := s.tokensClient.Get(strings.ToLower(tokenID), metav1.GetOptions{})

		if err != nil {
			logrus.Info("Failed to get token resource: %v", err)
			return tokens, 401, fmt.Errorf("Failed to retrieve auth token")
		}
		logrus.Debugf("storedToken token resource: %v", storedToken)
		externalID := storedToken.ExternalID
		set := labels.Set(map[string]string{"io.cattle.token.field.externalID": externalID})
		tokenList, err := s.tokensClient.List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})
		if err != nil {
			return tokens, 0, fmt.Errorf("Error getting tokens for user: %v selector: %v  err: %v", externalID, set.AsSelector().String(), err)
		}

		for _, t := range tokenList.Items {
			tokens = append(tokens, t)
		}
		return tokens, 0, nil

	}
	logrus.Info("Client nil %v", s.client)
	return tokens, 500, fmt.Errorf("No k8s Client configured")
}

func (s *tokenAPIServer) deleteToken(tokenKey string) (int, error) {
	logrus.Info("DELETE Token Invoked")

	if s.client != nil {
		err := s.tokensClient.Delete(strings.ToLower(tokenKey), &metav1.DeleteOptions{})

		if err != nil {
			return 500, fmt.Errorf("Failed to delete token")
		}
		logrus.Info("Deleted Token")
		return 0, nil

	}
	logrus.Info("Client nil %v", s.client)
	return 500, fmt.Errorf("No k8s Client configured")
}
