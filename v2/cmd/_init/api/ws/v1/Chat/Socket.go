package Chat

import (
	"fmt"

	"github.com/goaperture/goaperture/v2/ws/aperture"
)

var WsPrinter = aperture.WebSocket{
	Description: "Просто Вебсокет",
	Message: func(message string, conn *aperture.Conn) {
		fmt.Println(">>", message)
		conn.Send("<answer")
	},
}
