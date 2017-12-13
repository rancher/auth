package identities

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	"github.com/rancher/auth/util"
)

func (server *identityAPIServer) listIdentities(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("rAuthnSessionToken")
	if err != nil {
		log.Info("Failed to get token cookie: %v", err)
		util.ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	log.Infof("token cookie: %v %v", cookie.Name, cookie.Value)

	//getIdentities
	identities, status, err := server.getIdentities(cookie.Value)
	if err != nil {
		log.Errorf("listIdentities failed with error: %v", err)
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

func (server *identityAPIServer) searchIdentities(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("rAuthnSessionToken")
	if err != nil {
		log.Info("Failed to get token cookie: %v", err)
		util.ReturnHTTPError(w, r, http.StatusUnauthorized, "Invalid token cookie")
		return
	}

	log.Infof("token cookie: %v %v", cookie.Name, cookie.Value)

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("searchIdentities failed with error: %v", err)
		util.ReturnHTTPError(w, r, http.StatusBadRequest, fmt.Sprintf("Error reading input json data: %v", err))
		return
	}

	var jsonInput map[string]string
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &jsonInput)
		if err != nil {
			log.Errorf("searchIdentities: Error unmarshalling json request body: %v", err)
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
		log.Errorf("searchIdentities failed with error: %v", err)
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
