package note

type NoteService struct {
	repo *NoteRepo
}

func NewNoteService(repo *NoteRepo) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (ns *NoteService) List() (NoteList, error) {
	notes, err := ns.repo.List()
	if err != nil {
		return nil, err
	}

	return notes, nil
}
