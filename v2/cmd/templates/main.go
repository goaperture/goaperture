package main

import (
	"fmt"
)

func main() {
	fmt.Printf("SERVER Run on http://localhost:%d\n", api.Server.Port)
	api.Server.Routes(&routes.Routes).WebSockets(&ws.Routes).Run()
}
