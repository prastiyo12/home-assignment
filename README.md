# File Transfer Demo

## Overview
This project demonstrates a pragmatic file transfer mechanism for **offline-first POS (Point of Sale)** systems.

Each on-premise POS client maintains a **persistent outbound connection** to a central server.  
This allows the server to request data **on demand** without polling and without relying on client IP addresses.

The design reflects real-world retail environments where POS machines are deployed behind NATs or firewalls.

---

## Architecture

```
[ POS Client ]
     |
     |  WebSocket (persistent, outbound)
     v
[ Server ]
     |
     |  Command: DOWNLOAD
     v
[ POS Client ]
     |
     |  HTTP streaming upload
     v
[ Server ]
```

### Key Characteristics
- No polling
- No inbound connection to client
- NAT / firewall friendly
- Scales to many outlets
- Operationally simple

---

## Components

- **/server**
  - Cloud backend
  - Manages connected clients
  - Sends commands
  - Receives streamed files

- **/client**
  - Lightweight POS agent
  - Maintains persistent connection
  - Streams files on request
  - Designed to run as a compiled binary

---

## Authentication

A **shared API token** is used to authenticate:
- WebSocket connections
- File upload requests

This keeps the system secure while remaining simple and pragmatic.

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
CLIENT_ID=store-001
API_TOKEN=change-me
FILE_PATH=C:\temp\file_to_download.txt
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

For production or POS deployment, the client should be compiled into a binary:

```bash
go build -o pos-agent.exe
```

---

## Triggering a File Transfer

From the **server machine**:

```bash
curl http://localhost:8080/trigger
```

This simulates an administrative or system-initiated request to retrieve data from connected POS clients.

---

## Notes

- Files are streamed directly to disk to avoid high memory usage
- The client does not expose any inbound ports
- The system is suitable for offline-first retail environments
- Designed to be easy to operate and maintain by small engineering teams

---

## Disclaimer

This project is a **technical demonstration** intended for evaluation purposes.
It focuses on clarity, correctness, and pragmatic system design rather than production hardening.
