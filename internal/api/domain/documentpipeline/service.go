package documentpipeline

import (
	"io"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/api/domain/documentmeta"
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

func (dps *DocumentPipelineService) UploadDocumentPipeline(req *http.Request) (*documentmeta.DocumentMeta, error) {
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

	dm, err := dps.stirlingClient.GetMetaData(pdfBytes, header)
	if err != nil {
		return nil, err
	}

	assignRes, err := dps.seaweedClient.Assign()
	if err != nil {
		return nil, err
	}
	dm.Fid = assignRes.Fid
	err = dps.seaweedClient.StoreFile(assignRes.PublicURL, assignRes.Fid, pdfBytes, header)
	if err != nil {
		return nil, err
	}
	return dm, nil
}
