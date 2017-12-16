package identities

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	
	"github.com/rancher/norman/types"
	//"github.com/rancher/norman/parse"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/auth/util"
)

func (server *identityAPIServer) listIdentities(w http.ResponseWriter, r *http.Request) ([]v3.Identity, int, error) {
	identities := make([]v3.Identity, 0)
	cookie, err := r.Cookie("rAuthnSessionToken")
	if err != nil {
		logrus.Info("Failed to get token cookie: %v", err)
		//util.ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return identities, 0, err
	}

	logrus.Infof("listIdentities: token cookie: %v %v", cookie.Name, cookie.Value)

	//getIdentities
	/*identities, status, err := server.getIdentities(cookie.Value)
	if err != nil {
		logrus.Errorf("listIdentities failed with error: %v", err)
		if status == 0 {
			status = http.StatusInternalServerError
		}
		util.ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
		return
	}*/

	return server.getIdentities(cookie.Value)
}

func (server *identityAPIServer) handleListIdentities(apiContext *types.APIContext) error {
	fmt.Println("Handler called")
	logrus.Infof("Handler called %v", apiContext)
	
	var (
		err  error
	)
	
	identities, _ , err := server.listIdentities(apiContext.Response, apiContext.Request)
	
	logrus.Debugf("---getIdentities:  %v", identities)

	if err != nil {
		return err
	}
	
	//resp := &IdentityCollection{}
	
	
	/*store := apiContext.Schema.Store
	if store == nil {
		return nil
	}

	if apiContext.ID == "" {
		opts := parse.QueryOptions(apiContext, apiContext.Schema)
		data, err = store.List(apiContext, apiContext.Schema, opts)
	}

	if err != nil {
		return err
	}*/

	//apiContext.Response.Header().Set("content-type", "application/json")
	
	apiContext.WriteResponse(http.StatusOK, identities)
	return nil

}

func (server *identityAPIServer) handleSearchIdentities(apiContext *types.APIContext) error {
	//server.searchIdentities(apiContext.Response, apiContext.Request)
	return nil
}

func (server *identityAPIServer) searchIdentities(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("rAuthnSessionToken")
	if err != nil {
		logrus.Info("Failed to get token cookie: %v", err)
		util.ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	logrus.Infof("searchIdentities: token cookie: %v %v", cookie.Name, cookie.Value)

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("searchIdentities failed with error: %v", err)
		util.ReturnHTTPError(w, r, http.StatusBadRequest, fmt.Sprintf("Error reading input json data: %v", err))
		return
	}

	var jsonInput map[string]string
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &jsonInput)
		if err != nil {
			logrus.Errorf("searchIdentities: Error unmarshalling json request body: %v", err)
			util.ReturnHTTPError(w, r, http.StatusBadRequest, fmt.Sprintf("Error reading json request body: %v", err))
			return
		}
	} else {
		util.ReturnHTTPError(w, r, http.StatusBadRequest, fmt.Sprintf("Error reading input json data: %v", err))
		return
	}

	//searchIdentities
	identities, status, err := server.findIdentities(cookie.Value, jsonInput["name"])
	if err != nil {
		logrus.Errorf("searchIdentities failed with error: %v", err)
		if status == 0 {
			status = http.StatusInternalServerError
		}
		util.ReturnHTTPError(w, r, status, fmt.Sprintf("%v", err))
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(identities)
}
