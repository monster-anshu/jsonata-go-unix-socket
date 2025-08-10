# JSONata Go Socket Service – Documentation

## Overview

This is a lightweight socket-based service written in **Go** that allows you to evaluate **JSONata** expressions without embedding a heavy JavaScript runtime into your application.

It works by:

1. Listening on a Unix socket (or TCP port).
2. Accepting JSON requests containing:
   - A JSONata expression.
   - A JSON input object.
3. Returning the evaluated result or an error.

## 1. Installation

### 1.1 Build from source

You need Go installed for compilation (only on the build machine):

```bash
git clone https://github.com/your-repo/jsonata-go-socket.git
cd jsonata-go-socket
go build -o jsonata_server main.go
```

### 1.2 Deploy the binary

Copy the `jsonata_server` binary to your server:

```bash
cp jsonata_server /usr/local/bin/
chmod +x /usr/local/bin/jsonata_server
```

---

## 2. Running the Service

### 2.1 Quick Run (Foreground)

```bash
./jsonata_server
```

### 2.3 Run as a Systemd Service (Recommended)

Create `/etc/systemd/system/jsonata.service`:

```ini
[Unit]
Description=JSONata Socket Server
After=network.target

[Service]
ExecStart=/usr/local/bin/jsonata_server
Restart=always
RestartSec=2
User=www-data
WorkingDirectory=/usr/local/bin
StandardOutput=append:/var/log/jsonata_server.log
StandardError=append:/var/log/jsonata_server_error.log

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable jsonata.service
sudo systemctl start jsonata.service
```

Check status:

```bash
sudo systemctl status jsonata.service
```

---

## 3. Example PHP Client

A working example with PHP can be found [here.](./test.php)

## 4. Error Handling

- Invalid JSON request → `"success": false, "error": "Invalid JSON"`
- Invalid JSONata expression → `"success": false, "error": "Invalid expression"`
- Evaluation errors → `"success": false, "error": "Evaluation error"`

## 5. Notes

- Ensure the socket path is accessible to the user running PHP scripts.
- Use systemd to keep the service always running.
- Works with any language that can connect to a Unix or TCP socket.
