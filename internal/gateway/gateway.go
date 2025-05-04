package gateway

import (
	"github.com/jasonlovesdoggo/velo/internal/log"
	"net/http"
)

func Start(port string) error {
	// TODO: Hook into gRPC handlers or REST endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Info("Gateway listening", "port", port)
	return http.ListenAndServe(":"+port, nil)
}
