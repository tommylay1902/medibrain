package stirling

import (
	"fmt"
	"io"
	"net/http"
)

type StirlingClient struct {
	BaseURL string
	Client  *http.Client
}

func NewClient() *StirlingClient {
	return &StirlingClient{
		BaseURL: "http://localhost:3000",
		Client:  &http.Client{},
	}
}

// ForwardRequest forwards the entire HTTP request to Stirling
func (sc *StirlingClient) ForwardRequest(req *http.Request) (*http.Response, error) {
	// Build Stirling URL
	stirlingURL := fmt.Sprintf("%s/api/v1/analysis/document-properties", sc.BaseURL)
	// Create a new request to forward to Stirling
	forwardReq, err := http.NewRequest(req.Method, stirlingURL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("creating forward request: %w", err)
	}

	// Copy headers (important!)
	forwardReq.Header = req.Header.Clone()

	// Set content length
	forwardReq.ContentLength = req.ContentLength

	// Send the request
	resp, err := sc.Client.Do(forwardReq)
	if err != nil {
		return nil, fmt.Errorf("forwarding request: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err reading body")
	}
	fmt.Println(string(b))
	return resp, nil
}
