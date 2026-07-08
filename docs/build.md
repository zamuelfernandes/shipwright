# Build & Compiling Guide

This guide explains how to compile Shipwright into optimized standalone executable binaries and configure desktop integration.

---

## Standard Compilation

To compile Shipwright into a native executable for your current operating system and architecture:

```bash
go build -o shipwright cmd/server/main.go
```

To run it:
```bash
./shipwright
```

---

## Production Optimization (Under 10MB)

Go compiles down to a single self-contained binary. To minimize the binary size, you can strip debugging symbols and DWARF tables. This usually reduces the final file size by up to 40% (often bringing it under 10MB):

```bash
go build -ldflags="-s -w" -o shipwright cmd/server/main.go
```

---

## Linux Desktop Shortcut (.desktop)

To integrate Shipwright with your Linux application launcher (GNOME, KDE, etc.) with custom icon mapping:

1.  Create a desktop entry file in your user local applications folder:
    ```bash
    nano ~/.local/share/applications/shipwright.desktop
    ```
2.  Add the following configuration (replace `/caminho/para/seu/projeto/` with your actual absolute project path):
    ```ini
    [Desktop Entry]
    Type=Application
    Name=Shipwright
    Comment=Visual Docker Manager
    Exec=/caminho/para/seu/projeto/shipwright
    Path=/caminho/para/seu/projeto
    Icon=/caminho/para/seu/projeto/ui/dist/logo.png
    Terminal=false
    Categories=Development;Utility;
    ```
3.  **GNOME Trusted Launching**: In your desktop file manager, locate the `shipwright.desktop` file, right-click it, and select **"Allow Launching"** (or *"Permitir Executar"*) to authorize the system to double-click and run it.

---

## Cross-Compilation

Go supports building for other target operating systems and platforms out of the box using `GOOS` and `GOARCH` environment variables.

> [!NOTE]
> Because Shipwright uses CGO for compile-time bindings with `webkit2gtk` GUI window components, standard cross-compilation requires a cross-compiler toolchain installed on your machine, or disabling GUI dependencies if compiling as a pure background API server.

### Compiling as a pure HTTP Headless Server (Disabling Webview GUI)

If you only want to build a backend Docker manager accessible from remote browsers (without opening a native desktop window on run):

#### Windows Binary (Headless)
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o shipwright.exe cmd/server/main.go
```

#### Linux Binary (Headless Server)
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o shipwright cmd/server/main.go
```
