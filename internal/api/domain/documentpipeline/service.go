package documentpipeline

import (
	"fmt"
	"net/http"

	"github.com/tommylay1902/medibrain/internal/client/stirling"
)

type DocumentPipelineService struct {
	// repo           *DocumentPipelineRepo
	stirlingClient *stirling.StirlingClient
}

func NewService(
	// repo *DocumentPipelineRepo
	stirlingClient *stirling.StirlingClient,
) *DocumentPipelineService {
	return &DocumentPipelineService{
		// repo: repo,
		stirlingClient: stirlingClient,
	}
}

func (dps *DocumentPipelineService) UploadDocumentPipeline(req *http.Request) error {
	_, err := dps.stirlingClient.ForwardRequest(req)
	if err != nil {
		fmt.Println("error getting metadata")
		return err
	}
	return nil
}
