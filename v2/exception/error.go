package exception

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/goaperture/goaperture/v2/responce"
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

func Fall(message any, code string, status int) {
	panic(CustomPanic{Message: fmt.Sprintf("%s", message), Code: getCode(code), Status: getStatus(status)})
}

func NotAccess(permission string) {
	panic(CustomPanic{Message: "Нет доступа", Code: getCode("access denied"), Permission: permission, Status: getStatus(403)})
}

func Catch(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		if v, ok := r.(CustomPanic); ok {
			(*w).WriteHeader(v.Status.Code)
			json.NewEncoder(*w).Encode(responce.Error(v))
		} else {
			log.Println("!undefined panic", r)

			(*w).WriteHeader(500)
			json.NewEncoder(*w).Encode(responce.Error(CustomPanic{Message: fmt.Sprintf("%s", r)}))
		}
	}
}

func getStatus(code int) Status {
	return Status{Code: code, Text: http.StatusText(code)}
}

func getCode(code string) string {
	return strings.ReplaceAll(strings.ToUpper(code), " ", "_")
}
