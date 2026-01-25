package exception

import (
	"encoding/json"
	"net/http"
)

func Fall(message, code string) {
	panic(message)
}

func NotAccess(permission string) {
	panic(permission)
}

func Catch(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		json.NewEncoder(*w).Encode("error")
	}
}
