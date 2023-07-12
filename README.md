# DirectAdmin Exporter

This is a simple exporter for monitoring DirectAdmin using Prometheus. It retrieves metrics from the DirectAdmin API and exposes them in a format that can be scraped by Prometheus.

![Continous Integration status](https://github.com/piotr-ku/directadmin-exporter/actions/workflows/integration.yml/badge.svg?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/piotr-ku/directadmin-exporter)](https://goreportcard.com/report/github.com/piotr-ku/directadmin-exporter) ![coverage](https://raw.githubusercontent.com/piotr-ku/directadmin-exporter/badges/.badges/main/coverage.svg)

## Installation

To install and set up the DirectAdmin Exporter, follow these steps:

1. Clone the repository:

   ```shell
   git clone https://github.com/piotr-ku/directadmin-exporter.git
   ```

2. Install the required dependencies. Make sure you have [Go](https://golang.org/doc/install) installed.

3. Build the exporter:

   ```shell
   go build -o directadmin-exporter
   ```

4. Run the exporter:

   ```shell
   ./directadmin-exporter
   ```

## Usage

Before running the DirectAdmin Exporter, make sure you have the following prerequisites:

- Go (version 1.20 or later)
- DirectAdmin API credentials

To run the application, use the following command:

```shell
./directadmin-exporter --port <port-number> --ip <ip-address> --config <config-file-path>
```

Replace the placeholders with the appropriate values:

- `<port-number>`: Port number for the HTTP server (default: 8080)
- `<ip-address>`: IP address for the HTTP server (default: 127.0.0.1)
- `<config-file-path>`: Path to the configuration file

## Configuration

The DirectAdmin Exporter requires an API configuration to connect to the DirectAdmin server. Create an environment file in the following format:

```
DIRECTADMIN_HOSTNAME=<directadmin-hostname>
DIRECTADMIN_USERNAME=<directadmin-username>
DIRECTADMIN_TOKEN=<directadmin-token>
DIRECTADMIN_PORT=<directadmin-port>
DIRECTADMIN_PROTOCOL=<directadmin-protocol>
```

- `<directadmin-hostname>`: The hostname of the DirectAdmin server
- `<directadmin-username>`: The username for the DirectAdmin API.
- `<directadmin-token>`: The token or password for the DirectAdmin API.
- `<directadmin-port>`: The port number on which the DirectAdmin server is running.
- `<directadmin-protocol>`: The protocol to use for communication with the DirectAdmin server (`http` or `https`).

When running the application, provide the path to the environment file using the `--config` flag:

```shell
./directadmin-exporter --config <config-file-path>
```

## Metrics

The DirectAdmin Exporter collects various metrics exposed by the DirectAdmin server. These metrics are scraped periodically and made available for Prometheus to scrape.

The metrics endpoint is available at `/metrics` on the HTTP server.

## Contributing

Contributions to the DirectAdmin Exporter project are welcome. To contribute, please follow these guidelines:

1. Fork the repository and create a new branch.
2. Make your changes and ensure that the code is properly tested.
3. Submit a pull request describing your changes and the problem they solve.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
