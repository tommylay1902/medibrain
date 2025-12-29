package seaweedclient

type SeaWeedClient struct {
	MasterURL string
	VolumeURL string
}

type AssignResponse struct {
	Fid       string `json:"fid"`
	URL       string `json:"url"`
	PublicURL string `json:"publicUrl"`
	Count     int    `json:"count"`
}
