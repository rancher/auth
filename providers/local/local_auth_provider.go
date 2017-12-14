package local

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

//Constants for github
const (
	Name = "local"
)

func Configure(ctx context.Context, mgmtCtx *config.ManagementContext) *LProvider {
	l := &LProvider{
		Users:        mgmtCtx.Management.Users(""),
		Groups:       mgmtCtx.Management.Groups(""),
		GroupMembers: mgmtCtx.Management.GroupMembers(""),
	}
	return l
}

//LProvider implements an IdentityProvider for local auth
type LProvider struct {
	Users        v3.UserInterface
	Groups       v3.GroupInterface
	GroupMembers v3.GroupMemberInterface
}

//GetName returns the name of the provider
func (l *LProvider) GetName() string {
	return Name
}

func (l *LProvider) AuthenticateUser(loginInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error) {
	username := loginInput.LocalCredential.Username
	pwd := loginInput.LocalCredential.Password

	var groupIdentities []v3.Identity
	var userIdentity v3.Identity

	localUser, err := l.Users.Get(username, metav1.GetOptions{})

	if err != nil {
		log.Info("Failed to get User resource: %v", err)
		return userIdentity, groupIdentities, 401, fmt.Errorf("Invalid Credentials")
	}

	if strings.EqualFold(pwd, localUser.Secret) {
		userIdentity = v3.Identity{
			DisplayName:    "",
			LoginName:      localUser.ObjectMeta.Name,
			ProfilePicture: "",
			ProfileURL:     "",
			Kind:           "user",
			Me:             true,
			MemberOf:       false,
		}
		userIdentity.ObjectMeta = metav1.ObjectMeta{
			Name: localUser.ExternalID,
		}

		var status int
		groupIdentities, status, err = l.getGroupIdentities(localUser)

		if err != nil {
			log.Info("Failed to get group identities for local user: %v, user: %v", err, localUser.ObjectMeta.Name)
			return userIdentity, groupIdentities, status, fmt.Errorf("Error getting group identities for local user %v", err)
		}

	} else {
		return userIdentity, groupIdentities, 401, fmt.Errorf("Invalid Credentials")
	}

	log.Info("List userIdentity resource: %v", userIdentity)
	log.Info("List groupIdentities resource: %v", groupIdentities)

	return userIdentity, groupIdentities, 0, nil
}

func (l *LProvider) getGroupIdentities(user *v3.User) ([]v3.Identity, int, error) {
	var groupIdentities []v3.Identity

	if user != nil {
		externalID := user.ExternalID
		set := labels.Set(map[string]string{"io.cattle.groupmember.field.externalID": externalID})
		groupMemberList, err := l.GroupMembers.List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})

		if err != nil {
			return groupIdentities, 500, fmt.Errorf("Error listing groupmembers for user: %v selector: %v  err: %v", externalID, set.AsSelector().String(), err)
		}

		for _, gm := range groupMemberList.Items {
			log.Info("List groupmember resource: %v", gm)

			//find group for this member mapping
			localGroup, err := l.Groups.Get(gm.GroupName, metav1.GetOptions{})
			if err != nil {
				log.Errorf("Failed to get Group resource: %v", err)
				continue
			}
			log.Info("List group resource: %v", localGroup)
			groupIdentity := v3.Identity{
				DisplayName:    localGroup.Name,
				LoginName:      "",
				ProfilePicture: "",
				ProfileURL:     "",
				Kind:           "group",
				Me:             false,
				MemberOf:       false,
			}
			groupIdentity.ObjectMeta = metav1.ObjectMeta{
				Name: localGroup.Name,
			}
			groupIdentities = append(groupIdentities, groupIdentity)
		}
	}

	return groupIdentities, 0, nil
}
