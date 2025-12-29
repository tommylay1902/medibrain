package seaweedclient

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			// Only return the close error if no other error occurred
			err = fmt.Errorf("closing response body: %w", closeErr)
		}
	}()

	var result AssignResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}
