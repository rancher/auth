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

func Configure(ctx context.Context, mgmtCtx *config.ManagementContext) *LProvider {
	l := &LProvider{
		users:        mgmtCtx.Management.Users(""),
		groups:       mgmtCtx.Management.Groups(""),
		groupMembers: mgmtCtx.Management.GroupMembers(""),
	}
	return l
}

//LProvider implements an IdentityProvider for local auth
type LProvider struct {
	users        v3.UserInterface
	groups       v3.GroupInterface
	groupMembers v3.GroupMemberInterface
}

//GetName returns the name of the provider
func (l *LProvider) GetName() string {
	return "local"
}

func (l *LProvider) AuthenticateUser(loginInput v3.LoginInput) (v3.Identity, []v3.Identity, int, error) {
	username := loginInput.LocalCredential.Username
	pwd := loginInput.LocalCredential.Password

	var groupIdentities []v3.Identity
	var userIdentity v3.Identity

	localUser, err := l.users.Get(username, metav1.GetOptions{})

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
	return userIdentity, groupIdentities, 0, nil
}

func (l *LProvider) getGroupIdentities(user *v3.User) ([]v3.Identity, int, error) {
	var groupIdentities []v3.Identity

	if user != nil {
		externalID := user.ExternalID
		set := labels.Set(map[string]string{"io.cattle.groupmember.field.externalID": externalID})
		groupMemberList, err := l.groupMembers.List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})

		if err != nil {
			return groupIdentities, 500, fmt.Errorf("Error listing groupmembers for user: %v selector: %v  err: %v", externalID, set.AsSelector().String(), err)
		}

		for _, gm := range groupMemberList.Items {
			//find group for this member mapping
			localGroup, err := l.groups.Get(gm.GroupName, metav1.GetOptions{})
			if err != nil {
				log.Errorf("Failed to get Group resource: %v", err)
				continue
			}
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

func (l *LProvider) SearchIdentities(searchKey string, myToken v3.Token) ([]v3.Identity, int, error) {
	var identities []v3.Identity
	//is this local token?
	if myToken.AuthProvider != l.GetName() {
		return identities, 0, nil
	}

	usersList, err := l.users.List(metav1.ListOptions{})
	if err != nil {
		return identities, 0, fmt.Errorf("SearchIdentities: Error listing users %v", err)
	}
	//form user identities

	for _, user := range usersList.Items {
		//TODO: add code to filter out Users from other providers.

		if !strings.HasPrefix(user.ObjectMeta.Name, searchKey) {
			continue
		}
		userIdentity := v3.Identity{
			DisplayName:    "",
			LoginName:      user.ObjectMeta.Name,
			ProfilePicture: "",
			ProfileURL:     "",
			Kind:           "user",
			Me:             false,
			MemberOf:       false,
		}
		userIdentity.ObjectMeta = metav1.ObjectMeta{
			Name: user.ExternalID,
		}
		if l.isThisUserMe(myToken.UserIdentity, userIdentity) {
			userIdentity.Me = true
		}
		identities = append(identities, userIdentity)
	}

	groupsList, err := l.groups.List(metav1.ListOptions{})
	if err != nil {
		return identities, 0, fmt.Errorf("SearchIdentities: Error listing groups %v", err)
	}
	//form group identities

	for _, group := range groupsList.Items {
		//TODO: add code to filter out Users from other providers.

		if !strings.HasPrefix(group.ObjectMeta.Name, searchKey) {
			continue
		}

		groupIdentity := v3.Identity{
			DisplayName:    group.ObjectMeta.Name,
			LoginName:      "",
			ProfilePicture: "",
			ProfileURL:     "",
			Kind:           "group",
			Me:             false,
			MemberOf:       false,
		}
		groupIdentity.ObjectMeta = metav1.ObjectMeta{
			Name: group.ObjectMeta.Name,
		}
		if l.isMemberOf(myToken.GroupIdentities, groupIdentity) {
			groupIdentity.MemberOf = true
		}
		identities = append(identities, groupIdentity)
	}

	return identities, 0, nil
}

func (l *LProvider) isThisUserMe(me v3.Identity, other v3.Identity) bool {

	if me.ObjectMeta.Name == other.ObjectMeta.Name && me.LoginName == other.LoginName && me.Kind == other.Kind {
		return true
	}
	return false
}

func (l *LProvider) isMemberOf(myGroups []v3.Identity, other v3.Identity) bool {

	for _, mygroup := range myGroups {
		if mygroup.ObjectMeta.Name == other.ObjectMeta.Name && mygroup.Kind == other.Kind {
			return true
		}
	}
	return false
}
