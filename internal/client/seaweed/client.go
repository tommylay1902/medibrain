package seaweedclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/util"
)

func NewClient() *SeaWeedClient {
	return &SeaWeedClient{
		MasterURL: "http://localhost:9333",
		VolumeURL: "http://localhost:9000",
		Client:    &http.Client{},
	}
}

func (swc *SeaWeedClient) Assign() (*AssignResponse, error) {
	url := fmt.Sprintf("%s/dir/assign", swc.MasterURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error trying to get fid")
		return nil, err
	}

	var result AssignResponse

	err = util.Bind(&result, resp)
	if err != nil {
		fmt.Println("error binding Assign seaweed response to AssignRespones struct")
		return nil, err
	}
	return &result, nil
}

func (swc *SeaWeedClient) StoreFile(publicURL string, fid string, pdfBytes []byte, header *multipart.FileHeader) error {
	url := fmt.Sprintf("http://%s/%s", publicURL, fid)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return fmt.Errorf("create form file error: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(pdfBytes))
	if err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("close writer error: %v", err)
	}

	writer.Boundary()

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}

	res, err := swc.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return errors.New("not expected status code from StoreFile swc")
	}

	return nil
}

func (swc *SeaWeedClient) Delete(publicURL string, fid string) error {
	url := fmt.Sprintf("http://%s/%s", publicURL, fid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	res, err := swc.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 202 {
		return errors.New("not expected status code from delete swc")
	}
	return nil
}
