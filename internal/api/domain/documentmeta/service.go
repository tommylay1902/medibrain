package documentmeta

type DocumentMetaService struct {
	repo *DocumentMetaRepo
}

func NewService(repo *DocumentMetaRepo) *DocumentMetaService {
	return &DocumentMetaService{repo: repo}
}

func (s *DocumentMetaService) List() (DocumentMetaList, error) {
	result, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DocumentMetaService) Create(dm *DocumentMeta) error {
	err := s.repo.Create(dm)
	if err != nil {
		return err
	}
	return nil
}
