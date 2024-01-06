# Go Reverse Proxy

This project implements a basic reverse proxy in Go, capable of routing traffic based on hostname to different backend services. It supports HTTP and HTTPS traffic and can be configured via a JSON file.

## Features

- Hostname-based routing.
- HTTP and HTTPS support.
- Configurable via a JSON file.

## Prerequisites

- Go (version 1.13 or later).
- SSL/TLS certificates (for HTTPS support).

## Installation

To install the reverse proxy, use the following command:

```bash
go install github.com/the-zucc/reverse-proxy@latest
```

This will install the proxy application in your Go bin directory.

## Configuration

Create a `config.json` file in the same directory where you will run the application. This file should contain the routing configuration for the proxy. Here's an example:

```json
{
  "routes": [
    {
      "host": "example.com",
      "target": "http://localhost:8080"
    },
    {
      "host": "another-example.com",
      "target": "http://localhost:9090"
    }
  ]
}
```

## Running the Application

After installation and creation of your configuration, you can run the application with the following command:

```bash
reverse-proxy [-https] [-cert path/to/cert.pem] [-key path/to/key.pem] [-config path/to/config.json]
```

- `-https`: Enable HTTPS support (optional).
- `-cert`: Path to the SSL certificate file (required for HTTPS).
- `-key`: Path to the SSL certificate key file (required for HTTPS).
- `-config`: Path to the configuration file (default is `config.json` in the current directory).

## Running with Docker (Optional)

If you prefer to run the application using Docker, follow these steps:

1. Build the Docker image:

   ```bash
   docker build -t go-reverse-proxy .
   ```

2. Run the Docker container:

   ```bash
   docker run -d -p 80:80 -p 443:443 go-reverse-proxy
   ```

## Contributing

Contributions to this project are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

This project was inspired by the need for a simple, configurable reverse proxy written in Go.
