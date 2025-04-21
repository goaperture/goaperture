package main

import (
	"github.com/goaperture/goaperture/test/api"
)

func main() {
	token := "123"
	api.Serve(3000, &token)
}
