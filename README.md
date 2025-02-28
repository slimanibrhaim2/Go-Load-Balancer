Go Load Balancer

A simple but powerful load balancer implementation in Go that distributes incoming HTTP requests across multiple backend servers while monitoring their health status.

Features

Dynamic backend server health checking

JSON-based configuration

Connection tracking

Multiple backend server support

Health check endpoints

Simple and extensible design

Project Structure

go-load-balancer/
├── config.json
├── main.go
├── README.md
├── server-1/
│   └── server-1.go
├── server-2/
│   └── server-2.go
└── server-3/
    └── server-3.go

Prerequisites

Go 1.16 or higher

Basic understanding of HTTP servers

Configuration

The load balancer can be configured using the config.json file:

{
    "port": 8080,
    "servers": [
      {
        "url": "http://localhost:8001",
        "healthy": false
      },
      {
        "url": "http://localhost:8002",
        "healthy": false
      },
      {
        "url": "http://localhost:8003",
        "healthy": false
      }
    ],
    "healthCheckInterval": 10
}

Configuration Details

port: The port on which the load balancer will listen.

servers: List of backend servers with their URLs.

healthCheckInterval: Interval (in seconds) for health checks.

Getting Started

Clone the repository:

git clone https://github.com/slimanibrhaim2/go-load-balancer.git
cd go-load-balancer

Start the backend servers:

# Terminal 1
cd server-1
go run server-1.go

# Terminal 2
cd ../server-2
go run server-2.go

# Terminal 3
cd ../server-3
go run server-3.go

Start the load balancer:

go run main.go

Test the load balancer:

curl http://localhost:8080

How It Works

The load balancer reads the configuration from config.json on startup.

It periodically performs health checks on all backend servers.

Incoming requests are distributed among healthy backend servers.

Each server's connection count is tracked to ensure even distribution.

Health Checking

The load balancer performs regular health checks on all backend servers. A server is considered healthy if it responds to the /health endpoint with a 200 OK status.

Contributing

Fork the repository.

Create your feature branch: git checkout -b feature/amazing-feature

Commit your changes: git commit -m 'Add some amazing feature'

Push to the branch: git push origin feature/amazing-feature

Open a Pull Request.
