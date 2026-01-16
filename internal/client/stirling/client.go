package stirling

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
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

func (sc *StirlingClient) GetMetaData2(pdfBytes []byte, header *multipart.FileHeader) (*http.Response, error) {
	// First, let's create the request exactly as cURL would
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the file part
	part, err := writer.CreateFormFile("fileInput", header.Filename)
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	// Write file content
	_, err = io.Copy(part, bytes.NewReader(pdfBytes))
	if err != nil {
		return nil, fmt.Errorf("write file error: %v", err)
	}

	// CRITICAL: Close the writer to write the final boundary
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	// For debugging, let's see what we're sending
	boundary := writer.Boundary()
	fmt.Printf("Multipart boundary: %s\n", boundary)
	fmt.Printf("Content-Type will be: multipart/form-data; boundary=%s\n", boundary)

	// Create request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analysis/document-properties", sc.BaseURL),
		body)
	if err != nil {
		return nil, err
	}

	// Set headers - this is critical
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and log response
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Response: %s\n", resp.StatusCode, string(respBody))

	return &http.Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		Body:          io.NopCloser(bytes.NewReader(respBody)),
		ContentLength: int64(len(respBody)),
	}, nil
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
