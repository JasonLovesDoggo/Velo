package gateway

import (
	"crypto/subtle"
	"fmt"
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

	http.HandleFunc("/hooks/deploy", DeployHookHandler)
	// todo: ratelimit this endpoint

	log.Info("Gateway listening", "port", port)
	return http.ListenAndServe(":"+port, nil)
}

func DeployHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := r.URL.Query().Get("velo_secret")
	projectSlug := r.URL.Query().Get("project_slug")
	if secret == "" || projectSlug == "" {
		log.Warn("Deploy hook called with missing secret or project_slug")
		http.Error(w, "Missing velo_secret or project_slug query parameter", http.StatusBadRequest)
		return
	}

	allowed, err := checkSecret(w, projectSlug, secret)
	if err != nil {
		log.Error("Failed to check secret", "projectSlug", projectSlug, "error", err)
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	} else if !allowed { // don't throw error twice
		log.Warn("Deploy hook called with invalid secret", "projectSlug", projectSlug)
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
	}

	// --- Trigger Deployment ---
	log.Info("Valid deploy hook received", "projectSlug", projectSlug)
	err = triggerDeployment(projectSlug)
	if err != nil {
		log.Error("Failed to trigger deployment", "projectSlug", projectSlug, "error", err)
		http.Error(w, "Failed to trigger deployment", http.StatusInternalServerError)
		return
	}

	// --- Respond ---
	w.WriteHeader(http.StatusAccepted) // 202 Accepted is suitable for async operations
	fmt.Fprintf(w, "Deployment triggered for project: %s\n", projectSlug)
	log.Info("Deployment successfully triggered", "projectSlug", projectSlug)
}

// Placeholder for state access - replace with actual implementation
func getExpectedSecret(projectSlug string) (string, error) {
	// TODO: Replace with actual state/config lookup
	if projectSlug == "my-test-project" {
		// Example secret - store securely in practice
		return "supersecretwebhooktoken", nil
	}
	return "", fmt.Errorf("project not found: %s", projectSlug)
}

// Placeholder for triggering deployment - replace with actual implementation
func triggerDeployment(projectSlug string) error {
	log.Info("Triggering deployment", "projectSlug", projectSlug)
	// TODO:
	// 1. Fetch ServiceDefinition from state using projectSlug
	// 2. Get orchestrator manager instance
	// 3. Call manager.DeployService(def) or manager.UpdateService(id, def)
	// Example:
	// def, err := state.GetProjectDefinition(projectSlug)
	// if err != nil { return err }
	// _, err = orchestratorManager.DeployService(def) // Or UpdateService if it exists
	// return err
	return nil // Placeholder success
}

func checkSecret(w http.ResponseWriter, projectSlug, secret string) (bool, error) {
	expectedSecret, err := getExpectedSecret(projectSlug)
	if err != nil {
		log.Error("Failed to get expected secret for project", "projectSlug", projectSlug, "error", err)
		http.Error(w, "Project not found or configuration error", http.StatusNotFound)
		return false, err
	}
	if expectedSecret == "" {
		// Handle case where project exists but has no secret configured
		log.Error("No secret configured for project", "projectSlug", projectSlug)
		http.Error(w, "Deployment hook not configured for this project", http.StatusInternalServerError) // Or 400 Bad Request
		return false, fmt.Errorf("no secret configured for project: %s", projectSlug)
	}

	// Use constant time comparison to prevent timing attacks
	secretsMatch := subtle.ConstantTimeCompare([]byte(secret), []byte(expectedSecret)) == 1

	if !secretsMatch {
		log.Warn("Deploy hook called with invalid secret", "projectSlug", projectSlug)
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return false, fmt.Errorf("invalid secret for project: %s", projectSlug)
	}
	return true, nil
}
