package routes

type IndexInput struct {
}

type IndexOutput interface {
	string
}

func IndexTest(invoke func(IndexInput)) {
	invoke(IndexInput{})
}

func Index(input IndexInput) (any, error) {
	return "hello world", nil
}
