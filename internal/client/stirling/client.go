// Package stirling interacts with the stirling api
package stirling

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
		BaseURL: "http://Stirling:8080/api/v1",
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
		slog.Error("create form file error", slog.Any("err", err))
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		slog.Error("write file error", slog.Any("err", err))
		return nil, fmt.Errorf("write file error: %v", err)
	}

	writer.WriteField("outputFormat", "txt")
	err = writer.Close()
	if err != nil {
		slog.Error("close writer error", slog.Any("err", err))
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/convert/pdf/text", sc.BaseURL),
		body)
	if err != nil {
		slog.Error("error creating request object", slog.Any("err", err))
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-KEY", apiKey)
	resp, err := sc.Client.Do(req)
	if err != nil {
		slog.Error("error executing request", slog.Any("err", err))
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body", slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		slog.Error(fmt.Sprintf("expected status code 200, actual status code: %v", resp.StatusCode))
		return nil, errors.New("not expected status code")
	}
	result := string(respBody)
	return &result, nil
}

func (sc *StirlingClient) GetMetaData(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*metadata.Metadata, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", header.Filename)
	if err != nil {
		slog.Error("create form file error", slog.Any("err", err))
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		slog.Error("write file error", slog.Any("err", err))
		return nil, fmt.Errorf("write file error: %v", err)
	}

	err = writer.Close()
	if err != nil {
		slog.Error("close writer error", slog.Any("err", err))
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/analysis/document-properties", sc.BaseURL),
		body)
	if err != nil {
		slog.Error("error creating request object", slog.Any("err", err))
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-KEY", apiKey)
	resp, err := sc.Client.Do(req)
	if err != nil {
		slog.Error("error executing request", slog.Any("err", err))
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body", slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		slog.Error(fmt.Sprintf("expected status code 200, actual status code: %v", resp.StatusCode))
		return nil, errors.New("not expected status code")
	}

	var dm metadata.Metadata
	err = json.Unmarshal(respBody, &dm)
	if err != nil {
		slog.Error("error unmarshaling response body")
		return nil, err
	}

	return &dm, nil
}

func (sc *StirlingClient) GenerateThumbnail(pdfBytes []byte, apiKey string) ([]byte, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	stirlingURL := fmt.Sprintf("%s/convert/pdf/img", sc.BaseURL)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", "thumbnail.jpeg")
	if err != nil {
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		slog.Error("write file error", slog.Any("err", err))
		return nil, fmt.Errorf("write file error: %v", err)
	}

	writer.WriteField("pageNumbers", "1")
	writer.WriteField("imageFormat", "jpeg")
	writer.WriteField("singleOrMultiple", "single")
	writer.WriteField("colorType", "greyscale")
	writer.WriteField("dpi", "300")
	err = writer.Close()
	if err != nil {
		slog.Error("close writer error", slog.Any("err", err))
		return nil, fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()
	req, err := http.NewRequest("POST",
		stirlingURL,
		body)
	if err != nil {
		slog.Error("error creating request body", slog.Any("err", err))
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-Key", apiKey)
	resp, err := sc.Client.Do(req)
	if err != nil {
		slog.Error("error execiting request", slog.Any("err", err))
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading request body", slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		slog.Error(fmt.Sprintf("expected status code 200, actual status code: %v", resp.StatusCode))
		return nil, errors.New("generating thumbnail, status not expected status code")
	}

	return respBody, nil
}

func (sc *StirlingClient) UpdateMetaData(pdfBytes []byte, apiKey string, dm *metadata.Metadata) ([]byte, error) {
	var preservedBuf bytes.Buffer

	pdfReader := bytes.NewReader(pdfBytes)
	tee := io.TeeReader(pdfReader, &preservedBuf)
	stirlingURL := fmt.Sprintf("%s/misc/update-metadata", sc.BaseURL)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileInput", "")
	if err != nil {
		slog.Error("create form file error", slog.Any("err", err))
		return nil, fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, tee)
	if err != nil {
		slog.Error("error copying from tee reader", slog.Any("err", err))
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

	writer.WriteField("keywords", dm.Keywords)

	writer.Close()

	req, err := http.NewRequest("POST",
		stirlingURL,
		body)
	if err != nil {
		slog.Error("error creating req object", slog.Any("err", err))
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-API-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error executing request", slog.Any("err", err))
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading resp body", slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		slog.Error(fmt.Sprintf("expected status code 200, actual status code: %v", resp.StatusCode))
		return nil, errors.New("not expected status code")
	}

	return respBody, nil
}
