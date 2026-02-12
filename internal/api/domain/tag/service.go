package tag

type TagService struct {
	repo *TagRepo
}

func NewTagService(repo *TagRepo) *TagService {
	return &TagService{
		repo: repo,
	}
}

func (ts *TagService) List() (TagList, error) {
	tags, err := ts.repo.List()

	return tags, err
}
