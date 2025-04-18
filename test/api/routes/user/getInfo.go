package user

import "fmt"

type GetInfoInput struct {
	UserId int `json:"user_id"`
}

type GetInfoOutput interface {
	string
}

func GetInfoTest(invoke func(GetInfoInput)) {
	invoke(GetInfoInput{UserId: 1})
	invoke(GetInfoInput{UserId: 5})
}

func GetInfo(input GetInfoInput) (any, error) {
	return fmt.Sprint("name User[", input.UserId, "] - Andry"), nil
}
