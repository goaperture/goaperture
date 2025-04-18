package main

import (
	"github.com/goaperture/goaperture/test/api"
)

func main() {
	token := "123"
	if err := api.Serve(3000, &token); err != nil {
		panic(err)
	}
}
