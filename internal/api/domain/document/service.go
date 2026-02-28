package document

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/tommylay1902/medibrain/internal/api/domain/metadata"
	"github.com/tommylay1902/medibrain/internal/client/rag"
	seaweedclient "github.com/tommylay1902/medibrain/internal/client/seaweed"
	"github.com/tommylay1902/medibrain/internal/client/stirling"
)

// TODO: passing in shouldnt be pointer maybe? need to do research
type DocumentPipelineService struct {
	dmRepo         *metadata.MetadataRepo
	stirlingClient *stirling.StirlingClient
	seaweedClient  *seaweedclient.SeaWeedClient
	dms            *metadata.MetadataService
	ragClient      *rag.Rag
}

func NewService(
	dmRepo *metadata.MetadataRepo,
	seaweedClient *seaweedclient.SeaWeedClient,
	stirlingClient *stirling.StirlingClient,
	dms *metadata.MetadataService,
	ragClient *rag.Rag,
) *DocumentPipelineService {
	return &DocumentPipelineService{
		dmRepo:         dmRepo,
		seaweedClient:  seaweedClient,
		stirlingClient: stirlingClient,
		dms:            dms,
		ragClient:      ragClient,
	}
}

func (dps *DocumentPipelineService) UploadDocumentPipelineWithEdit(pdfBytes []byte, header *multipart.FileHeader, apiKey string, updateDM *metadata.Metadata) (*metadata.Metadata, error) {
	dmBytes, err := dps.stirlingClient.UpdateMetaData(pdfBytes, apiKey, updateDM)
	if err != nil {
		return nil, err
	}

	dm, err := dps.stirlingClient.GetMetaData(dmBytes, header, apiKey)
	if err != nil {
		return nil, err
	}

	assignRes, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}

	dm.PdfFid = assignRes.Fid
	err = dps.seaweedClient.StoreFile(assignRes.PublicURL, assignRes.Fid, dmBytes, header)
	if err != nil {
		fmt.Println("error storing file")
		return nil, err
	}

	thumbnail, err := dps.stirlingClient.GenerateThumbnail(dmBytes, apiKey)
	if err != nil {
		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{err}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, err
	}

	pdfAssign, err := dps.seaweedClient.Assign()
	if err != nil {
		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{err}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, err
	}

	dm.ThumbnailFid = pdfAssign.Fid
	err = dps.seaweedClient.StoreFile(pdfAssign.PublicURL, pdfAssign.Fid, thumbnail, header)
	if err != nil {
		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{err}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, err
	}

	err = dps.dmRepo.Create(dm)
	if err != nil {
		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid, dm.ThumbnailFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{err}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, err
	}

	return dm, nil
}

func (dps *DocumentPipelineService) UploadDocumentPipeline(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*metadata.Metadata, error) {
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
		assignErr := fmt.Errorf("assign failed: %w", err)

		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{assignErr}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, assignErr
	}

	dm.ThumbnailFid = pdfAssign.Fid
	err = dps.seaweedClient.StoreFile(pdfAssign.PublicURL, pdfAssign.Fid, thumbnail, header)
	if err != nil {
		assignErr := fmt.Errorf("assign failed: %w", err)

		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{assignErr}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, assignErr
	}

	err = dps.dmRepo.Create(dm)
	if err != nil {
		assignErr := fmt.Errorf("assign failed: %w", err)
		cleanupErrors := dps.cleanupResources(assignRes.PublicURL, dm.PdfFid, dm.ThumbnailFid)
		if len(cleanupErrors) > 0 {
			allErrors := append([]error{assignErr}, cleanupErrors...)
			return nil, errors.Join(allErrors...)
		}
		return nil, assignErr
	}

	return dm, nil
}

func (dps *DocumentPipelineService) UploadChunks(pdfBytes []byte, header *multipart.FileHeader, apiKey string) (*metadata.Metadata, error) {
	fmt.Println(string(pdfBytes))
	return nil, nil
}

func (dps *DocumentPipelineService) cleanupResources(publicURL string, fids ...string) []error {
	var errors []error
	for _, fid := range fids {
		if fid != "" && publicURL != "" {
			if err := dps.seaweedClient.Delete(publicURL, fid); err != nil {
				errors = append(errors, fmt.Errorf("failed to cleanup %s: %w", fid, err))
			}
		}
	}
	return errors
}
