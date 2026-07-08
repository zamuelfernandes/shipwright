# Installation & Setup Guide

This guide describes how to configure and run Shipwright on your local machine.

---

## Prerequisites

Before running or developing Shipwright, ensure your host machine has the following tools installed:

1.  **Go (Golang)**: Version `1.22` or higher.
2.  **Docker Engine**: Active and running on your local machine.
3.  **WebKitGTK Development Libraries** (Required for GUI window desktop client):
    *   **Ubuntu/Debian**: `sudo apt install libwebkit2gtk-4.1-dev build-essential`
    *   **Fedora**: `sudo dnf install gtk3-devel webkit2gtk4.1-devel`
    *   **Arch Linux**: `sudo pacman -S webkit2gtk-4.1`

> [!IMPORTANT]
> **Docker Permissions**: Ensure your user belongs to the `docker` system group so the application can access the Docker UNIX socket (`/var/run/docker.sock`) without root/sudo privilege:
> ```bash
> sudo usermod -aG docker $USER
> newgrp docker
> ```

---

## Local Configuration (.env)

Shipwright supports custom parameters loaded from a local environment file. 

1.  Copy the example template file:
    ```bash
    cp .env.example .env
    ```
2.  Open `.env` in your editor and configure the variables:
    *   `PORT`: Port of the web server (Defaults to `8080`).
    *   `APP_ENV`: Application environment (`development` / `production`).
    *   `DOCKER_HOST`: The endpoint to reach your Docker Engine. Defaults to your local OS socket or named pipe.

---

## Running in Development Mode

To run Shipwright without compiling a persistent executable binary (useful during codebase editing):

```bash
go run cmd/server/main.go
```

This starts the background HTTP server on `http://localhost:8080` (or your configured port) and opens a native desktop window showing the interface.
