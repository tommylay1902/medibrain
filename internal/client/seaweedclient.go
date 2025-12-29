package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SeaWeedClient struct{}

func GetVolumeSpace() ([]byte, error) {
	resp, err := http.Get("http://localhost:9333/dir/assign")
	if err != nil {
		fmt.Println("error trying to get fid")
		return nil, err
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error trying to read body")
		return nil, err
	}

	jsonResult, err := json.Marshal(bytes)
	if err != nil {
		fmt.Println("error converting response into json")
		return nil, err
	}
	return jsonResult, nil
}
