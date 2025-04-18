package filter

type GetFilterInfoInput struct {
}

type GetFilterInfoOutput interface {
	string
}

func GetFilterInfoTest(invoke func(GetFilterInfoInput)) {
	invoke(GetFilterInfoInput{})
}

func GetFilterInfo(input GetFilterInfoInput) (any, error) {
	return "hello world", nil
}
