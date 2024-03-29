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
	port := flag.String("port", "", "Port to run the server on")

	flag.Parse() // Parse the flags

	// Determine the default port based on the HTTPS flag
	if *port == "" {
		if *enableHTTPS {
			*port = "443"
		} else {
			*port = "80"
		}
	}

	// Load the proxy configuration from a JSON file
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up the reverse proxy based on the loaded configuration
	setupProxies(config)

	// Start the server
	addr := ":" + *port
	// Start the server in HTTPS mode if enabled, otherwise start in HTTP mode
	if *enableHTTPS {
		log.Println("Starting HTTPS server")
		log.Fatal(http.ListenAndServeTLS(addr, *certFile, *keyFile, nil))
	} else {
		log.Println("Starting HTTP server")
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}

// setupProxies configures the reverse proxy handler based on the given configuration.
func setupProxies(config *Config) {
	proxies := make(map[string]*httputil.ReverseProxy)

	for _, route := range config.Routes {
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			log.Fatalf("Error parsing target URL for host %s: %v", route.Host, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxies[route.Host] = proxy
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy, ok := proxies[r.Host]
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, r)
	})
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
