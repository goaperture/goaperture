package connections

type GetTypesInput struct {
}

type GetTypesOutput interface {
	string
}

func GetTypesTest(invoke func(GetTypesInput)) {
	invoke(GetTypesInput{})
}

func GetTypes(input GetTypesInput) (any, error) {
	return "hello world", nil
}
