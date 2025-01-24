# GoProxy

GoProxy is a lightweight proxy server built in Go, designed to block incoming requests based on IP address, headers, or query parameters. It provides an easy-to-use CLI to manage blocking rules and acts as a secure gateway to your backend services.

## Features
- Block requests by IP address.
- Block requests based on headers or query parameters (extendable).
- Simple CLI to manage block rules dynamically.
- Transparent proxying to your backend.

---

## Getting Started

### Prerequisites
- [Go](https://go.dev/) installed on your system (version 1.18 or higher recommended).

---

### Running the Proxy Server

To start the proxy server, run the following command:

```sh
go run ./cmd/server <BACKEND_URL>
```
<BACKEND_URL>: The URL of the backend service to which requests will be proxied.

Example
```sh
go run ./cmd/server https://www.youtube.com/
```
By default, the proxy listens on http://localhost:8080.

### Blocking an IP Address
To block an IP address from accessing the proxy, use the following command:

```sh
go run ./cmd/cli block <IP_ADDRESS>
```
Example:
```sh
go run ./cmd/cli block ::1
```
The specified IP will no longer be able to send requests through the proxy.


