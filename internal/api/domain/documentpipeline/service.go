package documentpipeline

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
)

// TODO: passing in shouldnt be pointer maybe? need to do research
type DocumentPipelineService struct {
	dmRepo         *documentmeta.DocumentMetaRepo
	stirlingClient *stirling.StirlingClient
	seaweedClient  *seaweedclient.SeaWeedClient
	dms            *documentmeta.DocumentMetaService
}

func NewService(
	dmRepo *documentmeta.DocumentMetaRepo,
	seaweedClient *seaweedclient.SeaWeedClient,
	stirlingClient *stirling.StirlingClient,
	dms *documentmeta.DocumentMetaService,
) *DocumentPipelineService {
	return &DocumentPipelineService{
		dmRepo:         dmRepo,
		seaweedClient:  seaweedClient,
		stirlingClient: stirlingClient,
		dms:            dms,
	}
}

func (dps *DocumentPipelineService) UploadDocumentPipelineWithEdit(req *http.Request) (*documentmeta.DocumentMeta, error) {
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

	apiKey := req.Header.Get("X-API-KEY")
	dm, err := dps.stirlingClient.GetMetaData(pdfBytes, header, apiKey)
	if err != nil {
		return nil, err
	}

	assignRes, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}

	dm.PdfFid = assignRes.Fid
	err = dps.seaweedClient.StoreFile(assignRes.PublicURL, assignRes.Fid, pdfBytes, header)
	if err != nil {
		return nil, err
	}

	thumbnail, err := dps.stirlingClient.GenerateThumbnail(pdfBytes, apiKey)
	if err != nil {
		return nil, err
	}

	pdfAssign, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}

	dm.ThumbnailFid = pdfAssign.Fid
	err = dps.seaweedClient.StoreFile(pdfAssign.PublicURL, pdfAssign.Fid, thumbnail, header)
	if err != nil {
		return nil, err
	}

	err = dps.dmRepo.Create(dm)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func (dps *DocumentPipelineService) UploadDocumentPipeline(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*documentmeta.DocumentMeta, error) {
	dm, err := dps.stirlingClient.GetMetaData(pdfBytes, header, apiKey)
	if err != nil {
		return nil, err
	}

	assignRes, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}

	dm.PdfFid = assignRes.Fid
	err = dps.seaweedClient.StoreFile(assignRes.PublicURL, assignRes.Fid, pdfBytes, header)
	if err != nil {
		return nil, err
	}

	thumbnail, err := dps.stirlingClient.GenerateThumbnail(pdfBytes, apiKey)
	if err != nil {
		return nil, err
	}

	pdfAssign, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}

	dm.ThumbnailFid = pdfAssign.Fid
	err = dps.seaweedClient.StoreFile(pdfAssign.PublicURL, pdfAssign.Fid, thumbnail, header)
	if err != nil {
		return nil, err
	}
	err = dps.dmRepo.Create(dm)
	if err != nil {
		return nil, err
	}
	return dm, nil
}
