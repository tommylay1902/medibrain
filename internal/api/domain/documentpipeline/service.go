package documentpipeline

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
	"github.com/tommylay1902/medibrain/internal/api/util"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
)

type DocumentPipelineService struct {
	// repo           *DocumentPipelineRepo
	stirlingClient *stirling.StirlingClient
	seaweedClient  *seaweedclient.SeaWeedClient
	dms            *documentmeta.DocumentMetaService
}

func NewService(
	// repo *DocumentPipelineRepo
	seaweedClient *seaweedclient.SeaWeedClient,
	stirlingClient *stirling.StirlingClient,
	dms *documentmeta.DocumentMetaService,
) *DocumentPipelineService {
	return &DocumentPipelineService{
		// repo: repo,
		seaweedClient:  seaweedClient,
		stirlingClient: stirlingClient,
		dms:            dms,
	}
}

func (dps *DocumentPipelineService) UploadDocumentPipeline(req *http.Request) error {
	metadataResponse, err := dps.stirlingClient.GetMetaData(req)
	if err != nil {
		fmt.Println("error getting metadata")
		return err
	}

	var metadata documentmeta.DocumentMeta
	err = util.Bind(&metadata, metadataResponse)
	if err != nil {
		fmt.Println("error binding response to metadata struct")
		return err
	}

	err = dps.dms.Create(&metadata)
	if err != nil {
		fmt.Println("error creating metadata in uploaddocumentpipeline")
		return err
	}
	fmt.Println("metadata created succesfully")
	res, err := dps.seaweedClient.Assign()
	if err != nil {
		fmt.Println("error getting space to store document")
		return err
	}
	fmt.Println(res)
	// _, err = dps.stirlingClient.OCRProcessing(req)
	return nil
}

func (dps *DocumentPipelineService) UploadDocumentPipeline2(req *http.Request) (*http.Response, error) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	file, header, err := req.FormFile("fileInput")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	response, err := dps.stirlingClient.GetMetaData2(pdfBytes, header)
	if err != nil {
		return nil, err
	}

	return response, nil
}
