// Package checkapi maintains the web based api for system access.
package checkapi

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type api struct {
	build string
	log   *logger.Logger
}

func newAPI(build string, log *logger.Logger) *api {
	return &api{
		build: build,

		log: log,
	}
}

func (api *api) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, data, statusCode)
}

func (api *api) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      api.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// This handler provides a free timer loop.

	return web.Respond(ctx, w, data, http.StatusOK)

}
