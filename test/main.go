package main

import (
	"github.com/goaperture/goaperture/test/api"
)

func main() {
	token := "123"

	api.Serve(3000, &token)

	// server := api.NewServer()
	// server.Middleware(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	// if err := server.Run(3000, &token); err != nil {
	// 	panic(err)
	// }
}
