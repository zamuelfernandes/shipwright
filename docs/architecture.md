# Architecture Guide

AnchorDock is designed under **Clean Architecture** principles to separate core business rules from presentation layers, data delivery mechanisms, and external frameworks.

---

## Directory Structure

```
anchordock/
├── cmd/
│   └── server/
│       └── main.go         # Wires dependencies, loads .env, launches server & GUI window
├── docs/                   # Documentation folder
├── internal/
│   ├── domain/             # Domain Layer (Core business models & abstract repository interfaces)
│   │   ├── container.go
│   │   ├── image.go
│   │   ├── exec.go
│   │   ├── project.go
│   │   └── stats.go
│   ├── usecase/            # Usecase Layer (Application business logic & coordinates actions)
│   │   ├── list_containers.go
│   │   ├── list_images.go
│   │   ├── start_project.go
│   │   ├── stop_project.go
│   │   ├── exec_container.go
│   │   ├── stream_logs.go
│   │   └── stream_stats.go
│   └── infrastructure/     # Infrastructure Layer (External tools adapters)
│       ├── docker/         # Docker engine socket API adapter (Moby SDK)
│       │   └── client.go
│       └── http/           # HTTP Routing, SSE servers, WebSockets and Handlers (SRP)
│           ├── router.go
│           ├── container_handler.go
│           ├── image_handler.go
│           └── project_handler.go
└── ui/                     # Web Application UI assets
    ├── embed.go            # Compiles static files into Go binary via //go:embed
    └── dist/
        └── index.html      # Single-page visual console dashboard (HTML/CSS/JS/Xterm.js)
```

---

## Design Principles

### Clean Architecture Layers
1.  **Domain Layer (Core)**: Completely isolated. It knows nothing about databases, network protocol routers, or the Docker SDK. It only defines structs (`Container`, `Image`) and repository interfaces.
2.  **Usecase Layer**: Implements orchestrating actions. Each usecase does one specific job. It depends only on domain contracts, allowing the infrastructure implementation to change without affecting domain rules.
3.  **Infrastructure Layer**: Connects external frameworks. This is where the Docker client connects to `/var/run/docker.sock` and HTTP/WebSocket endpoints are exposed to the client.

### SOLID Implementation
*   **Single Responsibility Principle (SRP)**: Each HTTP endpoint is separated into dedicated handlers (`container_handler.go`, `image_handler.go`, `project_handler.go`) instead of a single massive controller.
*   **Interface Segregation Principle (ISP)**: The monolithic container repository has been split into 5 small, client-specific contracts:
    *   `ContainerRepository` (lifecycle list, start, stop, prune)
    *   `ImageRepository` (listing images)
    *   `ProjectRepository` (compose projects batch operations)
    *   `TelemetryRepository` (SSE logs and stats telemetry streams)
    *   `ExecRepository` (interactive WebSocket terminal shell execution)
*   **Dependency Inversion Principle (DIP)**: Usecases do not instantiate database or Docker clients directly. They accept abstract repositories via constructor injections. The concrete `DockerClient` adapter implements these interfaces implicitly in Go.

---

## Web UI & GUI Window Integration

*   **HTML/CSS/JS Web Interface**: Pure vanilla stack without heavy UI node framework requirements. Handled in a single-page layout using CSS Variables (for dark mode) and native DOM manipulation.
*   **Embedded Assets**: Go compiles the entire folder `ui/dist` inside the final binary using `//go:embed`. When the user runs the binary, the HTTP server serves it natively from memory.
*   **Native GUI Integration**: AnchorDock launches a background HTTP server and opens a Webview window on runtime. Closing the OS window calls `server.Shutdown` which frees ports and cleanly exits the application.
