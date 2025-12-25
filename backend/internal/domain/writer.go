package domain

type Writer struct {
	id        uint64
	name      string
	birthYear int
	deathYear *int
	bio       *string
}

func NewWriter(id uint64, name string, birthYear int, deathYear *int, bio *string) *Writer {
	return &Writer{
		id:        id,
		name:      name,
		birthYear: birthYear,
		deathYear: deathYear,
		bio:       bio,
	}
}

func (w *Writer) ID() uint64 {
	return w.id
}

func (w *Writer) Name() string {
	return w.name
}

func (w *Writer) BirthYear() int {
	return w.birthYear
}

func (w *Writer) DeathYear() *int {
	return w.deathYear
}

func (w *Writer) Bio() *string {
	return w.bio
}
