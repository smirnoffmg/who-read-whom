package domain

type Work struct {
	id       uint64
	title    string
	authorID uint64
}

func NewWork(id uint64, title string, authorID uint64) *Work {
	return &Work{
		id:       id,
		title:    title,
		authorID: authorID,
	}
}

func (w *Work) ID() uint64 {
	return w.id
}

func (w *Work) Title() string {
	return w.title
}

func (w *Work) AuthorID() uint64 {
	return w.authorID
}
