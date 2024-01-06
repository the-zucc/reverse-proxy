package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Config represents the structure of the configuration file.
type Config struct {
	Routes []Route `json:"routes"` // Slice of Route objects for routing configuration
}

// Route defines the structure for each routing rule.
type Route struct {
	Host   string `json:"host"`   // Host name to match for the route
	Target string `json:"target"` // Target URL to proxy the request to
}

func main() {
	// Define command-line flags for configuring the server
	enableHTTPS := flag.Bool("https", false, "Enable HTTPS support")
	certFile := flag.String("cert", "cert.pem", "Path to the SSL certificate file")
	keyFile := flag.String("key", "key.pem", "Path to the SSL certificate key file")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse() // Parse the flags

	// Load the proxy configuration from a JSON file
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up the reverse proxy based on the loaded configuration
	setupProxies(config)

	// Start the server in HTTPS mode if enabled, otherwise start in HTTP mode
	if *enableHTTPS {
		log.Println("Starting HTTPS server")
		log.Fatal(http.ListenAndServeTLS(":443", *certFile, *keyFile, nil))
	} else {
		log.Println("Starting HTTP server")
		log.Fatal(http.ListenAndServe(":80", nil))
	}
}

// setupProxies configures the reverse proxy handlers based on the given configuration.
func setupProxies(config *Config) {
	for _, route := range config.Routes {
		// Parse the target URL from the configuration
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			log.Fatalf("Error parsing target URL for host %s: %v", route.Host, err)
		}

		// Create a new reverse proxy for the target URL
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Set up a handler function for each route
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Check if the request host matches the configured route host
			if r.Host == route.Host {
				// Use the proxy to handle the request
				proxy.ServeHTTP(w, r)
			}
		})
	}
}

// loadConfig reads and parses the configuration file.
func loadConfig(filename string) (*Config, error) {
	var config Config

	// Read the file contents
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the Config struct
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
