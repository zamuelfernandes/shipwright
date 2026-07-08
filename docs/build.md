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
2.  Add the following configuration (replace `/path/to/your/project/` with your actual absolute project path):
    ```ini
    [Desktop Entry]
    Type=Application
    Name=Shipwright
    Comment=Visual Docker Manager
    Exec=/path/to/your/project/shipwright
    Path=/path/to/your/project
    Icon=/path/to/your/project/ui/dist/logo.png
    Terminal=false
    Categories=Development;Utility;
    ```
3.  **GNOME Trusted Launching**: In your desktop file manager, locate the `shipwright.desktop` file, right-click it, and select **"Allow Launching"** (or *"Permitir Executar"*) to authorize the system to double-click and run it.

---

## Cross-Compilation & Multi-OS Desktop Build

Because Shipwright compiles with a native GUI desktop window using native system library bindings, it requires **CGO** (`CGO_ENABLED=1`) to compile the graphical layer. Below are the steps to build native desktop applications for Windows, macOS, and Linux.

---

### 🪟 Windows Desktop Build

#### Option A: Build Natively on Windows
If you are compiling on a Windows host machine (requires Git Bash, Go, and a GCC toolchain like MSYS2/MinGW installed):
```bash
go build -ldflags="-H=windowsgui -s -w" -o shipwright.exe cmd/server/main.go
```
*Note: The `-H=windowsgui` flag prevents an empty terminal command prompt (CMD) from launching alongside your window.*

#### Option B: Cross-Compile from Linux to Windows
If you are on Linux and want to build the `.exe` binary:
1.  Install the MinGW-w64 cross-compiler:
    ```bash
    sudo apt install mingw-w64
    ```
2.  Compile using CGO and the cross-compiler toolchain:
    ```bash
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H=windowsgui -s -w" -o shipwright.exe cmd/server/main.go
    ```

---

### 🍏 macOS Desktop Build

#### Build Natively on macOS
If you are compiling on a macOS host machine (requires Xcode Command Line Tools installed):
```bash
go build -ldflags="-s -w" -o shipwright cmd/server/main.go
```

---

### 🌐 Compiling as a Headless Server (Disabling Webview GUI)

If you only want to build a backend Docker manager accessible from remote web browsers, you can disable CGO and build a headless server without GUI window dependencies:

#### Windows Headless Server
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o shipwright-server.exe cmd/server/main.go
```

#### Linux Headless Server
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o shipwright-server cmd/server/main.go
```

