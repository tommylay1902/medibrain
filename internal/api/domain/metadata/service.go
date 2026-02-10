package metadata

type MetadataService struct {
	repo *MetadataRepo
}

func NewService(repo *MetadataRepo) *MetadataService {
	return &MetadataService{repo: repo}
}

func (s *MetadataService) List() (MetadataList, error) {
	result, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MetadataService) Create(dm *Metadata) error {
	err := s.repo.Create(dm)
	if err != nil {
		return err
	}
	return nil
}
