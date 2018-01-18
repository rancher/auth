package authconfig

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"

	"github.com/rancher/auth/providers"
	"github.com/rancher/auth/providers/github"
	"github.com/rancher/auth/tokens"
)

const (
	githubDefaultHostName = "https://github.com"
	reqMaxSize            = (2 * 1 << 20) + 1
)

func GithubConfigFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.Actions["configureTest"] = apiContext.URLBuilder.Action("configureTest", resource)
	resource.Actions["testAndApply"] = apiContext.URLBuilder.Action("testAndApply", resource)
}

func GithubConfigActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
	logrus.Infof("GithubConfigActionHandler called for action %v", actionName)

	if actionName == "configureTest" {
		return GithubConfigureTest(actionName, action, request)
	} else if actionName == "testAndApply" {
		return GithubConfigTestApply(actionName, action, request)
	}

	return nil
}

func GithubConfigureTest(actionName string, action *types.Action, request *types.APIContext) error {
	var githubConfig v3.GithubConfig
	githubConfigTestInput := v3.GithubConfigTestInput{}

	if err := json.NewDecoder(request.Request.Body).Decode(&githubConfigTestInput); err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}

	githubConfig = githubConfigTestInput.GithubConfig
	redirectURL := formGithubRedirectURL(githubConfig)

	logrus.Debugf("redirecting the user to %v", redirectURL)
	http.Redirect(request.Response, request.Request, redirectURL, http.StatusFound)

	return nil
}

func formGithubRedirectURL(githubConfig v3.GithubConfig) string {
	redirect := ""
	if githubConfig.Hostname != "" {
		redirect = githubConfig.Scheme + githubConfig.Hostname
	} else {
		redirect = githubDefaultHostName
	}
	redirect = redirect + "/login/oauth/authorize?client_id=" + githubConfig.ClientID + "&scope=read:org"

	return redirect
}

func GithubConfigTestApply(actionName string, action *types.Action, request *types.APIContext) error {
	var githubConfig v3.GithubConfig
	var githubCredential v3.GithubCredential
	githubConfigApplyInput := v3.GithubConfigApplyInput{}

	if err := json.NewDecoder(request.Request.Body).Decode(&githubConfigApplyInput); err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}

	githubConfig = githubConfigApplyInput.GithubConfig
	githubCredential = githubConfigApplyInput.GithubCredential
	
	//Call provider to testLogin
	p, err := providers.GetProvider("github")
	if err != nil {
		return err
	}
	githubProvider, ok := p.(*github.GProvider)
	if !ok {
		return fmt.Errorf("No github provider")
	}
	
	userPrincipal, groupPrincipals, providerInfo, status, err := githubProvider.LoginUser(githubCredential, &githubConfig)
	if err != nil {
		if status == 0 || status == 500{
			status = http.StatusInternalServerError
			return httperror.NewAPIErrorLong(status, "ServerError", fmt.Sprintf("Failed to login to github: %v", err))
		}
		return httperror.NewAPIErrorLong(status, "",
			fmt.Sprintf("Failed to login to github: %v", err))
	}
	
	//if this works, save githubConfig CR adding enabled flag
	githubConfig.Enabled = githubConfigApplyInput.Enabled
	err = githubProvider.SaveGithubConfig(githubConfig)
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Failed to save github config: %v", err))
	}
	
	//update User with github principalID?
	
	
	//create a new token, set this token as the cookie and return 200
	token, status, err := tokens.GenerateNewLoginToken(userPrincipal, groupPrincipals, providerInfo)
	if err != nil {
		log.Errorf("Login failed with error: %v", err)
		if status == 0 || status == 500 {
			status = http.StatusInternalServerError
			return httperror.NewAPIErrorLong(status, "ServerError", fmt.Sprintf("Failed to login to github: %v", err))
		}
		return httperror.NewAPIErrorLong(status, "", fmt.Sprintf("Failed to login to github: %v", err))
	}
		
	isSecure := false
	if request.Request.URL.Scheme == "https" {
		isSecure = true
	}

	tokenCookie := &http.Cookie {
		Name:     tokens.CookieName,
		Value:    token.Name,
		Secure:   isSecure,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, tokenCookie)
	request.WriteResponse(http.StatusOK, nil)
	
	return nil
}
