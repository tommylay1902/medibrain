package stirling

import (
	"fmt"
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

func (sc *StirlingClient) GetMetaData(req *http.Request) (*http.Response, error) {
	stirlingURL := fmt.Sprintf("%s/api/v1/analysis/document-properties", sc.BaseURL)
	return sc.forwardReq(req, stirlingURL)
}

func (sc *StirlingClient) UpdateMetaData(req *http.Request) (*http.Response, error) {
	stirlingURL := fmt.Sprintf("%s/api/v1/misc/update-metadata", sc.BaseURL)
	return sc.forwardReq(req, stirlingURL)
}

func (sc *StirlingClient) OCRProcessing(req *http.Request) (*http.Response, error) {
	stirlingURL := fmt.Sprintf("%s/api/v1/misc/ocr-pdf", sc.BaseURL)
	return sc.forwardReq(req, stirlingURL)
}

func (sc *StirlingClient) forwardReq(req *http.Request, stirlingURL string) (*http.Response, error) {
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
	return resp, nil
}
