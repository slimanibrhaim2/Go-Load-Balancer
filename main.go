package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	URL         string `json:"url"`
	Healthy     bool   `json:"healthy"`
	Connections int    // Track active connections
}

type Config struct {
	Port                int      `json:"port"`
	Servers             []Server `json:"servers"`
	HealthCheckInterval int      `json:"healthCheckInterval"`
}

var (
	config     Config
	serverPool []Server
	mu         sync.Mutex
)

func loadConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	serverPool = config.Servers
	log.Printf("Loaded configuration: %+v\n", config)
}

func getLeastConnectionsServer() *Server {
	mu.Lock()
	defer mu.Unlock()

	minConnections := int(^uint(0) >> 1) // maximum int value
	var candidates []*Server

	log.Println("Checking for healthy servers...")
	for i := range serverPool {
		server := &serverPool[i]
		if server.Healthy {
			log.Printf("Server %s is healthy (Connections: %d)", server.URL, server.Connections)
			if server.Connections < minConnections {
				minConnections = server.Connections
				candidates = []*Server{server}
			} else if server.Connections == minConnections {
				candidates = append(candidates, server)
			}
		} else {
			log.Printf("Server %s is unhealthy", server.URL)
		}
	}

	if len(candidates) == 0 {
		log.Println("No healthy servers available")
		return nil
	}

	// Randomly select one server among those with the least connections
	selected := candidates[rand.Intn(len(candidates))]
	selected.Connections++ // Increment connection count
	log.Printf("Selected server: %s (Connections: %d)", selected.URL, selected.Connections)
	return selected
}

func healthCheck() {
	for {
		mu.Lock()
		log.Println("Performing health checks...")
		for i, server := range serverPool {
			resp, err := http.Get(server.URL + "/health")
			if err != nil || resp.StatusCode != http.StatusOK {
				serverPool[i].Healthy = false
				log.Printf("Server %s is unhealthy: %v", server.URL, err)
			} else {
				serverPool[i].Healthy = true
				log.Printf("Server %s is healthy", server.URL)
			}
		}
		mu.Unlock()
		time.Sleep(time.Duration(config.HealthCheckInterval) * time.Second)
	}
}

// handleRequest selects a server and uses an external HTML template to display its URL.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	server := getLeastConnectionsServer()
	if server == nil {
		http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
		return
	}

	// Prepare data for the template.
	data := struct {
		ServerURL string
	}{
		ServerURL: server.URL,
	}

	// Since we're just displaying the selection, immediately decrement the connection count.
	mu.Lock()
	server.Connections--
	mu.Unlock()

	// Parse and execute the external HTML template (dashboard.html).
	tmpl, err := template.ParseFiles("dashboard.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func main() {
	// Seed the random number generator for tie-breaking.
	rand.Seed(time.Now().UnixNano())

	loadConfig()
	go healthCheck()

	http.HandleFunc("/", handleRequest)
	log.Printf("Load balancer started on port %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
