package note

type NoteService struct {
	repo *NoteRepo
}

func NewNoteService(repo *NoteRepo) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (ns *NoteService) List()
