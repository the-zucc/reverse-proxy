package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Config holds the configuration for the reverse proxy
type Config struct {
	Routes []Route `json:"routes"`
}

// Route defines a single route configuration
type Route struct {
	Host   string `json:"host"`
	Target string `json:"target"`
}

func main() {
	// Load the configuration
	config, err := loadConfig("proxy-config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create a map to hold the host to proxy mapping
	proxies := make(map[string]*httputil.ReverseProxy)

	// Set up a reverse proxy for each route
	for _, route := range config.Routes {
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			log.Fatalf("Error parsing target URL for host %s: %v", route.Host, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxies[route.Host] = proxy
	}

	// Custom handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		proxy, ok := proxies[r.Host]
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// loadConfig loads the configuration from a JSON file
func loadConfig(filename string) (*Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
