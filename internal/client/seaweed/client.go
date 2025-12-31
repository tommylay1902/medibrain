package seaweedclient

import (
	"fmt"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/util"
)

func NewClient() *SeaWeedClient {
	return &SeaWeedClient{
		MasterURL: "http://localhost:9333",
		VolumeURL: "http://localhost:9000",
	}
}

func (swc *SeaWeedClient) Assign() (*AssignResponse, error) {
	url := fmt.Sprintf("%s/dir/assign", swc.MasterURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error trying to get fid")
		return nil, err
	}
	var result AssignResponse

	err = util.Bind(&result, resp)
	if err != nil {
		fmt.Println("error binding Assign seaweed response to AssignRespones struct")
		return nil, err
	}
	return &result, nil
}
