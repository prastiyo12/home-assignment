# File Transfer Demo

## Overview
This project demonstrates a **pragmatic and on-demand file transfer mechanism** between multiple on-premise clients and a central server.

Each client maintains a **persistent outbound connection** to the server, allowing the server to request files **at any time** without polling and without relying on client IP addresses.

The design is suitable for environments where clients are deployed behind NATs or firewalls and direct inbound connections are not feasible.

---

## Architecture

```
[ Client ]
     |
     |  WebSocket (persistent, outbound)
     v
[ Server ]
     |
     |  Command: DOWNLOAD
     v
[ Client ]
     |
     |  HTTP streaming upload
     v
[ Server ]
```

### Key Characteristics
- On-demand file transfer
- No polling mechanism
- No inbound connection to clients
- NAT / firewall friendly
- Scales efficiently with multiple clients
- Simple and operationally practical

---

## Components

- **/server**
  - Central backend service
  - Manages connected clients
  - Sends commands to request files
  - Receives streamed file uploads

- **/client**
  - Lightweight client agent
  - Maintains a persistent connection to the server
  - Streams files upon request
  - Designed to run as a compiled binary

---

## Authentication

A **shared API token** is used to authenticate:
- WebSocket connections
- File upload requests

This approach keeps the system secure while remaining simple and easy to operate.

---

## Configuration

All configuration is provided via **environment variables**.

### Server (.env.example)
```
API_TOKEN=change-me
SERVER_PORT=8080
```

### Client (.env.example)
```
SERVER_WS=ws://SERVER_IP:8080/ws
SERVER_UPLOAD=http://SERVER_IP:8080/upload
CLIENT_ID=client-001
API_TOKEN=change-me
FILE_PATH=C:\temp\file_to_transfer.txt
```

> `.env` files are intentionally excluded from version control.

---

## Running the Server

```bash
cd server
cp .env.example .env
go run main.go
```

---

## Running the Client

```bash
cd client
cp .env.example .env
go run main.go
```

For deployment, the client can be compiled into a standalone binary:

```bash
go build -o client-agent.exe
```

---

## Triggering a File Transfer

From the **server machine**:

```bash
curl http://localhost:8080/trigger
```

This simulates a system-initiated request to retrieve a file from any connected client.

---

## Notes

- Files are streamed directly to disk to avoid high memory usage
- Clients do not expose inbound ports
- The system is designed for distributed environments with limited network access
- Focused on clarity, correctness, and pragmatic system design

---

## Disclaimer

This project is a **technical demonstration** intended for evaluation purposes.
It prioritizes simplicity and correctness over production hardening.
