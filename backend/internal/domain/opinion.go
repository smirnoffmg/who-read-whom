package domain

type Opinion struct {
	writerID      uint64
	workID        uint64
	sentiment     bool
	quote         string
	source        string
	page          *string
	statementYear *int
}

func NewOpinion(
	writerID, workID uint64,
	sentiment bool,
	quote, source string,
	page *string,
	statementYear *int,
) *Opinion {
	return &Opinion{
		writerID:      writerID,
		workID:        workID,
		sentiment:     sentiment,
		quote:         quote,
		source:        source,
		page:          page,
		statementYear: statementYear,
	}
}

func (o *Opinion) WriterID() uint64 {
	return o.writerID
}

func (o *Opinion) WorkID() uint64 {
	return o.workID
}

func (o *Opinion) Sentiment() bool {
	return o.sentiment
}

func (o *Opinion) Quote() string {
	return o.quote
}

func (o *Opinion) Source() string {
	return o.source
}

func (o *Opinion) Page() *string {
	return o.page
}

func (o *Opinion) StatementYear() *int {
	return o.statementYear
}
