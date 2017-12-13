package identities

import (
	"context"
	"fmt"
	//log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

type identityAPIServer struct {
	ctx              context.Context
	client           *config.ManagementContext
	identitiesClient v3.IdentityInterface
}

func newIdentityAPIServer(ctx context.Context, mgmtCtx *config.ManagementContext) (*identityAPIServer, error) {
	if mgmtCtx == nil {
		return nil, fmt.Errorf("Failed to build tokenAPIHandler, nil ManagementContext")
	}
	apiServer := &identityAPIServer{
		ctx:              ctx,
		client:           mgmtCtx,
		identitiesClient: mgmtCtx.Management.Identities(""),
	}
	return apiServer, nil
}

func (s *identityAPIServer) getIdentities(tokenKey string) ([]v3.Identity, int, error) {
	var identities []v3.Identity

	/*token, status, err := GetToken(tokenKey)

	if err != nil {
		return identities, 401, err
	} else {
		identities = append(identities, token.UserIdentity)
		identities = append(identities, token.GroupIdentities...)

		return identities, status, nil
	}*/

	identities = append(identities, getUserIdentity())
	identities = append(identities, getGroupIdentities()...)

	return identities, 0, nil

}

func (s *identityAPIServer) findIdentities(tokenKey string, name string) ([]v3.Identity, int, error) {
	var identities []v3.Identity

	/*token, status, err := GetToken(tokenKey)

	if err != nil {
		return identities, 401, err
	} else {
		identities = append(identities, token.UserIdentity)
		identities = append(identities, token.GroupIdentities...)

		return identities, status, nil
	}*/

	identities = append(identities, getUserIdentity())
	identities = append(identities, getGroupIdentities()...)

	return identities, 0, nil

}

func getUserIdentity() v3.Identity {

	identity := v3.Identity{
		LoginName:      "dummy",
		DisplayName:    "Dummy User",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "user",
		Me:             true,
		MemberOf:       false,
	}
	identity.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap_cn=dummy,dc=tad,dc=rancher,dc=io",
	}

	return identity
}

func getGroupIdentities() []v3.Identity {

	var identities []v3.Identity

	identity1 := v3.Identity{
		DisplayName:    "Admin group",
		LoginName:      "Administrators",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       true,
	}
	identity1.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap_cn=group1,dc=tad,dc=rancher,dc=io",
	}

	identity2 := v3.Identity{
		DisplayName:    "Dev group",
		LoginName:      "Developers",
		ProfilePicture: "",
		ProfileURL:     "",
		Kind:           "group",
		Me:             false,
		MemberOf:       true,
	}
	identity2.ObjectMeta = metav1.ObjectMeta{
		Name: "ldap_cn=group2,dc=tad,dc=rancher,dc=io",
	}

	identities = append(identities, identity1)
	identities = append(identities, identity2)

	return identities
}
