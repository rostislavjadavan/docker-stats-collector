# Docker Stats Collector

The Docker Stats Collector is a Go-based service that periodically collects performance statistics from running Docker containers on a host system. It gathers information such as:

- CPU usage
- Memory usage
- Network I/O
- Container image

This data is stored in a SQLite database, allowing for easy querying and analysis of container performance over time. The service is designed to run continuously, providing ongoing monitoring of your Docker environment.
