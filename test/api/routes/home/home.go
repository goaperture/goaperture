package home

type HomeInput struct {
}

type HomeOutput interface {
	string
}

func HomeTest(invoke func(HomeInput)) {
	invoke(HomeInput{})
}

func Home(input HomeInput) (any, error) {
	return "hello world", nil
}
