package seaweedclient

import "net/http"

type SeaWeedClient struct {
	MasterURL string
	VolumeURL string
	Client    *http.Client
}

type AssignResponse struct {
	Fid       string `json:"fid"`
	URL       string `json:"url"`
	PublicURL string `json:"publicUrl"`
	Count     int    `json:"count"`
}
