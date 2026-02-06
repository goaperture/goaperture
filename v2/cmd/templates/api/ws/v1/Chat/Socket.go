package Chat

import (
	"fmt"

	"github.com/goaperture/goaperture/v2/ws/aperture"
)

var WsPrinter = aperture.WebSocket{
	Description: "Просто Вебсокет",
	Message: func(message any, conn *aperture.Conn) {
		fmt.Println(">>", message)
		conn.SendText("<< answer")
	},
}
