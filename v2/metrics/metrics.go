package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func BindHanders(server *http.ServeMux) {
	server.Handle("/metrics", promhttp.Handler())
}
