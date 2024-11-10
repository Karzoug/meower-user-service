package health

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-user-service/internal/delivery/http/response"
	"github.com/Karzoug/meower-user-service/pkg/buildinfo"
	"github.com/Karzoug/meower-user-service/pkg/healthcheck"
)

type Checker interface {
	HealthCheck(ctx context.Context) healthcheck.Response
}

func RoutesFunc(hch Checker, logger zerolog.Logger) func(mux *http.ServeMux) {
	logger = logger.With().
		Str("component", "http server: health handlers").
		Logger()

	hdl := handlers{
		healthChecker: hch,
		logger:        logger,
	}

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /readiness", hdl.Readiness)
		mux.HandleFunc("GET /liveness", hdl.Liveness)
	}
}

type handlers struct {
	healthChecker Checker
	logger        zerolog.Logger
}

// Readiness checks if the server is ready to start accepting traffic.
func (h *handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	res := h.healthChecker.HealthCheck(r.Context())

	statusCode := http.StatusOK
	if res.Status == healthcheck.Fail {
		statusCode = http.StatusServiceUnavailable
	}

	if err := response.JSON(w, statusCode, res); err != nil {
		h.logger.Error().
			Err(err).
			Msg("couldn't write response")
	}
}

// Liveness returns status info if the service is alive and build info.
// If the app is deployed to a Kubernetes cluster, it will also return pod, node,
// and namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h *handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	res := h.healthChecker.HealthCheck(r.Context())

	statusCode := http.StatusOK
	if res.Status == healthcheck.Fail {
		statusCode = http.StatusServiceUnavailable
	}

	data := struct {
		HealthCheck healthcheck.Response `json:"healthCheck"`
		BuildInfo   buildinfo.BuildInfo  `json:"buildInfo"`
		Host        string               `json:"host,omitempty"`
		Name        string               `json:"name,omitempty"`
		PodIP       string               `json:"podIp,omitempty"`
		Node        string               `json:"node,omitempty"`
		Namespace   string               `json:"namespace,omitempty"`
		GOMAXPROCS  string               `json:"gomaxprocs,omitempty"`
	}{
		HealthCheck: res,
		BuildInfo:   buildinfo.Get(),
		Host:        host,
		Name:        os.Getenv("KUBERNETES_NAME"),
		PodIP:       os.Getenv("KUBERNETES_POD_IP"),
		Node:        os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:   os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS:  os.Getenv("GOMAXPROCS"),
	}

	if err := response.JSON(w, statusCode, data); err != nil {
		h.logger.Error().
			Err(err).
			Msg("couldn't write response")
	}
}
