package stirling

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
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

func (sc *StirlingClient) GetTextFromPdf(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*string, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", header.Filename)
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		return nil, fmt.Errorf("write file error: %v", err)
	}

	writer.WriteField("outputFormat", "txt")
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/convert/pdf/text", sc.BaseURL),
		body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-KEY", apiKey)
	resp, err := sc.Client.Do(req)
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
	result := string(respBody)
	return &result, nil
}

func (sc *StirlingClient) GetMetaData(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*metadata.DocumentMeta, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", header.Filename)
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
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
	req.Header.Add("X-API-KEY", apiKey)
	resp, err := sc.Client.Do(req)
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

	var dm metadata.DocumentMeta
	err = json.Unmarshal(respBody, &dm)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (sc *StirlingClient) GenerateThumbnail(pdfBytes []byte, apiKey string) ([]byte, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	stirlingURL := fmt.Sprintf("%s/api/v1/convert/pdf/img", sc.BaseURL)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", "thumbnail.jpeg")
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		return nil, fmt.Errorf("write file error: %v", err)
	}

	writer.WriteField("pageNumbers", "1")
	writer.WriteField("imageFormat", "jpeg")
	writer.WriteField("singleOrMultiple", "single")
	writer.WriteField("colorType", "greyscale")
	writer.WriteField("dpi", "300")
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()
	req, err := http.NewRequest("POST",
		stirlingURL,
		body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-Key", apiKey)
	resp, err := sc.Client.Do(req)
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

	return respBody, nil
}

func (sc *StirlingClient) UpdateMetaData(pdfBytes []byte, apiKey string, dm *metadata.DocumentMeta) ([]byte, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	stirlingURL := fmt.Sprintf("%s/api/v1/misc/update-metadata", sc.BaseURL)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", "")
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		return nil, fmt.Errorf("write file error: %v", err)
	}
	writer.WriteField("deletaAll", "false")
	if dm.Title != nil {
		writer.WriteField("title", *dm.Title)
	}
	if dm.Author != nil {
		writer.WriteField("author", *dm.Author)
	}
	if dm.Subject != nil {
		writer.WriteField("subject", *dm.Subject)
	}
	if dm.CreationDate != nil {
		writer.WriteField("creationDate", *dm.CreationDate)
	}
	if dm.ModificationDate != nil {
		writer.WriteField("modificationDate", *dm.ModificationDate)
	}

	// writer.Boundary()
	writer.Close()
	req, err := http.NewRequest("POST",
		stirlingURL,
		body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-Key", apiKey)
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

	return respBody, nil
}

// func (sc *StirlingClient) OCRProcessing(req *http.Request) (*http.Response, error) {
// 	stirlingURL := fmt.Sprintf("%s/api/v1/misc/ocr-pdf", sc.BaseURL)
// }
