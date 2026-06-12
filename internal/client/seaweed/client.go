// Package seaweedclient interacts with the seaweed store
package seaweedclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/util"
)

func NewClient() *SeaWeedClient {
	return &SeaWeedClient{
		MasterURL: "http://seaweedfs-master:9333",
		VolumeURL: "http://seaweedfs-volume:9000",
		Client:    &http.Client{},
	}
}

func (swc *SeaWeedClient) Assign() (*AssignResponse, error) {
	url := fmt.Sprintf("%s/dir/assign", swc.MasterURL)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("error trying to get fid", slog.Any("err", err))
		return nil, err
	}

	var result AssignResponse

	err = util.Bind(&result, resp)
	if err != nil {
		slog.Error("error binding Assign seaweed response to AssignRespones struct", slog.Any("err", err))
		return nil, err
	}
	return &result, nil
}

func (swc *SeaWeedClient) StoreFile(publicURL string, fid string, pdfBytes []byte, header *multipart.FileHeader) error {
	url := fmt.Sprintf("%s/%s", swc.VolumeURL, fid)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		slog.Error("create form file error", slog.Any("err", err))
		return fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(pdfBytes))
	if err != nil {
		slog.Error("write file error: ", slog.Any("err", err))
		return fmt.Errorf("write file error: %v", err)
	}

	err = writer.Close()
	if err != nil {
		slog.Error("close writer error:", slog.Any("err", err))
		return fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		slog.Error("error creating request object", slog.Any("err", err))
		return err
	}

	res, err := swc.Client.Do(req)
	if err != nil {
		slog.Error("Error running request object: ", slog.Any("err", err))
		return err
	}

	if res.StatusCode != 201 {
		slog.Error(fmt.Sprintf("expected status code: 201, recieved status code: %v", res.StatusCode))
		return errors.New("not expected status code from StoreFile swc")
	}

	return nil
}

func (swc *SeaWeedClient) Delete(publicURL string, fid string) error {
	url := fmt.Sprintf("%s/%s", swc.VolumeURL, fid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		slog.Error("error creating delete request object", slog.Any("err", err))
		return err
	}

	res, err := swc.Client.Do(req)
	if err != nil {
		slog.Error("error executing request object", slog.Any("err", err))
		return err
	}

	if res.StatusCode != 202 {

		slog.Error(fmt.Sprintf("expected status code: 202, recieved status code: %v", res.StatusCode))
		return errors.New("not expected status code from delete swc")
	}
	return nil
}
