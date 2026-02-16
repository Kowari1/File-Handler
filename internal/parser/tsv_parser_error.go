package parser

type ParseError struct {
	Line    int
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}
