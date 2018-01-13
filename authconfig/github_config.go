package authconfig

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"

	"github.com/rancher/types/apis/management.cattle.io/v3"
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
	return nil
}
