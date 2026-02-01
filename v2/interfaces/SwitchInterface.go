package interfaces

import (
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
)

type Switch interface {
	Handler(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request)
}
