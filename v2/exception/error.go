package exception

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Status struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type CustomPanic struct {
	Message    string `json:"message"`
	Code       string `json:"code"`
	Permission string `json:"permission,omitempty"`
	Status     Status `json:"status"`
}

func Fall(message, code string, status int) {
	panic(CustomPanic{Message: message, Code: getCode(code), Status: getStatus(status)})
}

func NotAccess(permission string) {
	panic(CustomPanic{Message: "Нет доступа", Code: getCode("access denied"), Permission: permission, Status: getStatus(403)})

}

func Catch(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		if v, ok := r.(CustomPanic); ok {
			(*w).WriteHeader(v.Status.Code)
			json.NewEncoder(*w).Encode(v)
		} else {
			log.Println("!undefined panic", r)
			panic(r)
		}

	}
}

func getStatus(code int) Status {
	return Status{Code: code, Text: http.StatusText(code)}
}

func getCode(code string) string {
	return strings.ReplaceAll(strings.ToUpper(code), " ", "_")
}
