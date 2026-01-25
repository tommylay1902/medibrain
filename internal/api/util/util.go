package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Bind[T any](save *T, response *http.Response) error {
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println("error closing response body stream")
		}
	}()

	if err := json.NewDecoder(response.Body).Decode(&save); err != nil {
		return err
	}
	return nil
}
