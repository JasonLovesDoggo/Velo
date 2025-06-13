package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
)

// WebServer provides a web interface for Velo
type WebServer struct {
	manager manager.Manager
	server  *http.Server
}

// NewWebServer creates a new web server
func NewWebServer(mgr manager.Manager, port string) *WebServer {
	ws := &WebServer{
		manager: mgr,
	}

	mux := http.NewServeMux()

	// Static files
	mux.HandleFunc("/static/", ws.handleStatic)

	// Pages
	mux.HandleFunc("/", ws.handleHome)
	mux.HandleFunc("/deployments", ws.handleDeployments)
	mux.HandleFunc("/deploy", ws.handleDeploy)
	mux.HandleFunc("/services", ws.handleServices)

	// API endpoints
	mux.HandleFunc("/api/deployments", ws.handleAPIDeployments)
	mux.HandleFunc("/api/deploy", ws.handleAPIDeploy)
	mux.HandleFunc("/api/services", ws.handleAPIServices)

	ws.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return ws
}

// Start starts the web server
func (ws *WebServer) Start() error {
	log.Info("Starting web server", "address", ws.server.Addr)
	return ws.server.ListenAndServe()
}

// Stop stops the web server
func (ws *WebServer) Stop() error {
	return ws.server.Close()
}

// Page handlers
func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Velo - Deployment Platform</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .nav { margin: 20px 0; }
        .nav a { margin-right: 20px; padding: 8px 16px; background: #007acc; color: white; text-decoration: none; border-radius: 4px; }
        .nav a:hover { background: #005a9e; }
        .status { background: #e8f5e8; padding: 15px; border-radius: 4px; border-left: 4px solid #4caf50; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Velo Deployment Platform</h1>
        <div class="status">
            <strong>Status:</strong> Manager node is running and ready to accept deployments.
        </div>
        <div class="nav">
            <a href="/deployments">View Deployments</a>
            <a href="/deploy">Deploy Service</a>
            <a href="/services">Manage Services</a>
        </div>
        <h2>Welcome to Velo</h2>
        <p>Velo is a lightweight, self-hostable deployment and operations platform built on Docker Swarm.</p>
        <p>Use the navigation above to manage your deployments and services.</p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func (ws *WebServer) handleDeployments(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Deployments - Velo</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .nav { margin: 20px 0; }
        .nav a { margin-right: 20px; padding: 8px 16px; background: #007acc; color: white; text-decoration: none; border-radius: 4px; }
        .nav a:hover { background: #005a9e; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: 600; }
        .status-running { color: #4caf50; font-weight: bold; }
        .status-failed { color: #f44336; font-weight: bold; }
        .status-pending { color: #ff9800; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Deployments</h1>
        <div class="nav">
            <a href="/">Home</a>
            <a href="/deploy">Deploy Service</a>
            <a href="/services">Manage Services</a>
        </div>
        <div id="deployments">
            <p>Loading deployments...</p>
        </div>
    </div>
    
    <script>
        async function loadDeployments() {
            try {
                const response = await fetch('/api/deployments');
                const deployments = await response.json();
                
                let html = '<table><thead><tr><th>Service Name</th><th>Image</th><th>Status</th><th>ID</th></tr></thead><tbody>';
                
                if (deployments.length === 0) {
                    html += '<tr><td colspan="4" style="text-align: center; padding: 40px; color: #666;">No deployments found. <a href="/deploy">Deploy your first service</a></td></tr>';
                } else {
                    deployments.forEach(deployment => {
                        const statusClass = 'status-' + deployment.state;
                        html += ` + "`" + `<tr>
                            <td>${deployment.service.name}</td>
                            <td>${deployment.service.image}</td>
                            <td class="${statusClass}">${deployment.state}</td>
                            <td><code>${deployment.id}</code></td>
                        </tr>` + "`" + `;
                    });
                }
                
                html += '</tbody></table>';
                document.getElementById('deployments').innerHTML = html;
            } catch (error) {
                document.getElementById('deployments').innerHTML = '<p style="color: red;">Error loading deployments: ' + error.message + '</p>';
            }
        }
        
        loadDeployments();
        setInterval(loadDeployments, 5000); // Refresh every 5 seconds
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func (ws *WebServer) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ws.handleAPIDeploy(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Deploy Service - Velo</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .nav { margin: 20px 0; }
        .nav a { margin-right: 20px; padding: 8px 16px; background: #007acc; color: white; text-decoration: none; border-radius: 4px; }
        .nav a:hover { background: #005a9e; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; font-weight: 600; color: #333; }
        input, textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
        input:focus, textarea:focus { outline: none; border-color: #007acc; box-shadow: 0 0 0 2px rgba(0, 122, 204, 0.1); }
        button { background: #007acc; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; }
        button:hover { background: #005a9e; }
        .alert { padding: 15px; margin: 20px 0; border-radius: 4px; }
        .alert-success { background: #e8f5e8; border-left: 4px solid #4caf50; color: #2e7d32; }
        .alert-error { background: #ffebee; border-left: 4px solid #f44336; color: #c62828; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Deploy Service</h1>
        <div class="nav">
            <a href="/">Home</a>
            <a href="/deployments">View Deployments</a>
            <a href="/services">Manage Services</a>
        </div>
        
        <form id="deployForm">
            <div class="form-group">
                <label for="serviceName">Service Name:</label>
                <input type="text" id="serviceName" name="serviceName" required placeholder="my-web-app">
            </div>
            
            <div class="form-group">
                <label for="image">Docker Image:</label>
                <input type="text" id="image" name="image" required placeholder="nginx:latest">
            </div>
            
            <div class="form-group">
                <label for="replicas">Replicas:</label>
                <input type="number" id="replicas" name="replicas" value="1" min="1" max="10">
            </div>
            
            <div class="form-group">
                <label for="environment">Environment Variables (one per line, KEY=VALUE):</label>
                <textarea id="environment" name="environment" rows="4" placeholder="NODE_ENV=production\nPORT=3000"></textarea>
            </div>
            
            <button type="submit">Deploy Service</button>
        </form>
        
        <div id="result"></div>
    </div>
    
    <script>
        document.getElementById('deployForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const serviceName = document.getElementById('serviceName').value;
            const image = document.getElementById('image').value;
            const replicas = parseInt(document.getElementById('replicas').value);
            const envText = document.getElementById('environment').value;
            
            // Parse environment variables
            const environment = {};
            if (envText.trim()) {
                envText.split('\n').forEach(line => {
                    const [key, ...valueParts] = line.split('=');
                    if (key && valueParts.length > 0) {
                        environment[key.trim()] = valueParts.join('=').trim();
                    }
                });
            }
            
            const deployment = {
                serviceName,
                image,
                replicas,
                environment
            };
            
            try {
                const response = await fetch('/api/deploy', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(deployment)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    document.getElementById('result').innerHTML = 
                        ` + "`" + `<div class="alert alert-success">
                            <strong>Success!</strong> Service deployed successfully.<br>
                            <strong>Deployment ID:</strong> <code>${result.deploymentId}</code><br>
                            <strong>Status:</strong> ${result.status}
                        </div>` + "`" + `;
                    document.getElementById('deployForm').reset();
                } else {
                    document.getElementById('result').innerHTML = 
                        ` + "`" + `<div class="alert alert-error">
                            <strong>Error:</strong> ${result.error || 'Deployment failed'}
                        </div>` + "`" + `;
                }
            } catch (error) {
                document.getElementById('result').innerHTML = 
                    ` + "`" + `<div class="alert alert-error">
                        <strong>Error:</strong> ${error.message}
                    </div>` + "`" + `;
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func (ws *WebServer) handleServices(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Services - Velo</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        .nav { margin: 20px 0; }
        .nav a { margin-right: 20px; padding: 8px 16px; background: #007acc; color: white; text-decoration: none; border-radius: 4px; }
        .nav a:hover { background: #005a9e; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🛠️ Services</h1>
        <div class="nav">
            <a href="/">Home</a>
            <a href="/deployments">View Deployments</a>
            <a href="/deploy">Deploy Service</a>
        </div>
        <p>Service management features will be available here.</p>
        <p>This will include service scaling, updates, and monitoring capabilities.</p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

// API handlers
func (ws *WebServer) handleAPIDeployments(w http.ResponseWriter, r *http.Request) {
	// For now, return empty list - this will be implemented when the manager interface is extended
	deployments := []config.DeploymentStatus{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

func (ws *WebServer) handleAPIDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ServiceName string            `json:"serviceName"`
		Image       string            `json:"image"`
		Replicas    int               `json:"replicas"`
		Environment map[string]string `json:"environment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.ServiceName == "" || req.Image == "" {
		http.Error(w, "Service name and image are required", http.StatusBadRequest)
		return
	}

	if req.Replicas <= 0 {
		req.Replicas = 1
	}

	// Create ServiceDefinition from request
	serviceDef := config.ServiceDefinition{
		Name:        req.ServiceName,
		Image:       req.Image,
		Environment: req.Environment,
		Replicas:    req.Replicas,
	}

	// Deploy the service using the manager
	deploymentID, err := ws.manager.DeployService(serviceDef)
	if err != nil {
		http.Error(w, fmt.Sprintf("Deployment failed: %v", err), http.StatusInternalServerError)
		return
	}
	response := struct {
		DeploymentID string `json:"deploymentId"`
		Status       string `json:"status"`
	}{
		DeploymentID: deploymentID,
		Status:       "deployed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) handleAPIServices(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement service management API
	services := []struct {
		Name     string `json:"name"`
		Image    string `json:"image"`
		Replicas int    `json:"replicas"`
		Status   string `json:"status"`
	}{
		{Name: "example-service", Image: "nginx:latest", Replicas: 2, Status: "running"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (ws *WebServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Simple static file serving
	file := filepath.Base(r.URL.Path)
	switch file {
	case "style.css":
		w.Header().Set("Content-Type", "text/css")
		// Return basic CSS if needed
	default:
		http.NotFound(w, r)
	}
}
