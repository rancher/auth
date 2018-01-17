package model

import "github.com/rancher/types/apis/management.cattle.io/v3"

//TokenCollection structure contains the token collection fields as per rancher spec
type TokenCollection struct {
	Type         string     `json:"type"`
	ResourceType string     `json:"resourceType"`
	Links        *Links     `json:"links"`
	Data         []v3.Token `json:"data"`
}

type Links struct {
	Self string `json:"self,omitempty"`
}
