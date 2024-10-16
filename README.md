# Docker Stats Collector

## What This Service Does

The Docker Stats Collector is a Go-based service that periodically collects performance statistics from running Docker containers on a host system. It gathers information such as:

- CPU usage
- Memory usage
- Network I/O
- Container image

This data is stored in a SQLite database, allowing for easy querying and analysis of container performance over time. The service is designed to run continuously, providing ongoing monitoring of your Docker environment.

## Running the Service on a Server

To run this service on your Ubuntu server:

1. Ensure Docker is installed and running on your server.

2. Copy the compiled `docker-stats-collector` binary to `/usr/local/bin/` on your server.

3. Copy the `docker-stats-collector.service` file from this repository to `/etc/systemd/system/` on your server.

4. Start and enable the service:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start docker-stats-collector
   sudo systemctl enable docker-stats-collector
   ```

The service will now start automatically on system boot and restart if it crashes.

To check the status of the service:

```bash
sudo systemctl status docker-stats-collector
```

To view logs:

```bash
sudo journalctl -u docker-stats-collector
```

## Building the Service

### Dependencies

Before building, ensure you have the following installed:

- Go (version 1.17 or later)
- SQLite3 development libraries
- GCC (for CGO compilation)

On Ubuntu, you can install these with:

```bash
sudo apt-get update
sudo apt-get install golang sqlite3 libsqlite3-dev gcc
```

On macOS, if you're cross-compiling for Linux:

```bash
brew install go sqlite3 FiloSottile/musl-cross/musl-cross
```

### Build Process

1. Clone this repository:

   ```bash
   git clone https://github.com/yourusername/docker-stats-collector.git
   cd docker-stats-collector
   ```

2. Build the binary:

   - On Linux:

     ```bash
     go build -o docker-stats-collector -ldflags="-w -s" .
     ```

   - On macOS (cross-compiling for Linux):
     ```bash
     CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
     go build -o docker-stats-collector -ldflags="-w -s" .
     ```

   Alternatively, you can use the provided Makefile:

   ```bash
   make build
   ```

3. The resulting `docker-stats-collector` binary can be deployed to your Linux server.

## Configuration

The service accepts the following command-line flags:

- `-interval`: Data collection interval in seconds (default: 5)
- `-db`: Path to SQLite database file (default: "container_stats.db")
- `-log`: Path to log file (default: "docker_stats.log")

You can modify these in the `docker-stats-collector.service` file to adjust the service's behavior.
