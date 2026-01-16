package stirling

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
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

func (sc *StirlingClient) GetMetaData(pdfBytes []byte, header *multipart.FileHeader) (*documentmeta.DocumentMeta, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", header.Filename)
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(pdfBytes))
	if err != nil {
		return nil, fmt.Errorf("write file error: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analysis/document-properties", sc.BaseURL),
		body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("not expected status code")
	}

	var dm documentmeta.DocumentMeta
	err = json.Unmarshal(respBody, &dm)
	if err != nil {
		return nil, err
	}

	return &dm, nil
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
