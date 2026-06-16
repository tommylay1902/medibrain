package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type Pydocument struct {
	BaseURL string
	Client  *http.Client
}

func NewPydocument() *Pydocument {
	return &Pydocument{
		BaseURL: "http://pydocument-server:8000/api/v1",
		Client:  &http.Client{},
	}
}

func (pd *Pydocument) GetDocumentSplits(pdfBytes []byte, header *multipart.FileHeader) (*DocumentSplitResponse, error) {
	url := fmt.Sprintf("%s/document/split", pd.BaseURL)

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
		url, body)
	if err != nil {
		slog.Error("error creating request object", slog.Any("err", err))
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := pd.Client.Do(req)
	if err != nil {
		slog.Error("error executing request", slog.Any("err", err))
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body", slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		slog.Error(fmt.Sprintf("expected status code 200, actual status code: %v", resp.StatusCode))
		return nil, errors.New("not expected status code")
	}
	var result DocumentSplitResponse

	if err := json.Unmarshal(respBody, &result); err != nil {
		slog.Error("error unmarshalling response body", slog.Any("err", err))
		return nil, err
	}
	return &result, nil
}
