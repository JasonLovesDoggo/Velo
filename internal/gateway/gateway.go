package gateway

import (
	"log"
	"net/http"
)

func Start(port string) error {
	// TODO: Hook into gRPC handlers or REST endpoints
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Gateway listening on :" + port)
	return http.ListenAndServe(":"+port, nil)
}
