package documentpipeline

import (
	"mime/multipart"

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

func (dps *DocumentPipelineService) UploadDocumentPipelineWithEdit(pdfBytes []byte, header *multipart.FileHeader, apiKey string, updateDM *documentmeta.DocumentMeta) (*documentmeta.DocumentMeta, error) {
	dmBytes, err := dps.stirlingClient.UpdateMetaData(pdfBytes, apiKey, updateDM)
	if err != nil {
		return nil, err
	}

	dm, err := dps.stirlingClient.GetMetaData(dmBytes, header, apiKey)

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
