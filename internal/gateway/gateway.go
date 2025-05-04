package gateway

import (
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/pkg/core"
	"net/http"
)

// an HTTP gateway for external access to the Velo system. TODO :)

func Start(port string) error {
	// TODO: Hook into gRPC handlers or REST endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Velo version: " + core.Version))
	})

	log.Info("Gateway listening", "port", port)
	return http.ListenAndServe(":"+port, nil)
}
