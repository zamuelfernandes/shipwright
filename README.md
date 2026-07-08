# Shipwright ⚓

[![Go Version](https://img.shields.io/github/go-mod/go-version/zamuelfernandes/shipwright?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Docker Compatible](https://img.shields.io/badge/Docker-Compatible-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/github/license/zamuelfernandes/shipwright?style=flat-square)](LICENSE)

Shipwright is an ultra-lightweight, visual local Docker container and image manager, inspired by Portainer, written entirely in **Go** and packaged as a native Linux desktop application under **12MB**.

---

## 🛠️ Tech Stack & Technologies

*   ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) **Go (Golang)** — Standard Library & Native `http.ServeMux` router.
*   ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white) **Moby / Docker Go SDK** — Direct local host Docker UNIX socket communication.
*   ![WebSockets](https://img.shields.io/badge/WebSockets-010101?style=flat-square&logo=socket.io&logoColor=white) **WebSockets & SSE** — Real-time telemetry monitoring streams and interactive terminal shell sessions.
*   ![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=flat-square&logo=html5&logoColor=white) ![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=flat-square&logo=css3&logoColor=white) ![JS](https://img.shields.io/badge/JavaScript-F7DF1E?style=flat-square&logo=javascript&logoColor=black) **Vanilla Web Stack** — Pure HTML5/CSS/JS frontend embedded in the binary via `//go:embed`.
*   ![GTK3](https://img.shields.io/badge/GTK3-7FE719?style=flat-square&logo=gnome&logoColor=white) **WebKitGTK Webview** — Native Linux desktop window integration.

---

## 🌟 The Story & Spirit of the Project

Shipwright was born out of frustration with **Docker Desktop's heavy resource consumption** (often eating up gigabytes of RAM just to list local containers). While researching lighter alternatives, I looked at existing tools like **Portainer** (which requires running another complex container structure) and **lazydocker** (a great CLI tool, but lacking a visual browser-based GUI). 

I wanted something different: a simple, visual dashboard that was ultra-lightweight, lightning fast, and focused precisely on my daily developer workflows. This curiosity sparked a coding experiment to test the limits of rapid, high-quality development in a new stack (Go + Embedded Vanilla WebUI), fully guided and paired with an **AI coding assistant**.

---

## 🚀 Key Features / Funcionalidades Principais

*   **Docker Compose Grouping / Agrupamento por Compose**: Automatically groups containers belonging to the same project / Agrupa containers do mesmo projeto Compose.
*   **Batch Operations / Comandos em Lote**: Start and stop entire projects concurrently using goroutines / Inicia e para projetos Compose concorrentemente.
*   **Interactive Web Terminal / Terminal Web Interativo**: Direct terminal command execution using **Xterm.js** and WebSockets / Console ativo interativo direto no navegador.
*   **Inline Telemetry / Telemetria Inline**: Persistent live CPU and Memory usage monitoring / Acompanhamento contínuo de consumo de CPU/RAM.
*   **Local Images Dashboard / Lista de Imagens**: Complete overview of downloaded Docker images / Tabela de imagens locais no topo do painel.
*   **Desktop Client / Janela Nativa**: Runs as a native desktop application window under **12MB** / Abre em janela nativa no Linux.

---

## 📖 Documentation & Guides / Guias e Documentação

To keep the repository clean and structured, detailed guides have been moved to the `/docs` directory:

1.  **[Setup & Installation Guide (English)](docs/setup.md)**: Prerequisites, local `.env` configuration variables, and development execution commands.
2.  **[Build & Compiling Guide (English)](docs/build.md)**: Compile commands, production size optimization, headless server compilation, and Linux `.desktop` launcher creation.
3.  **[Architecture & SOLID Guide (English)](docs/architecture.md)**: Clean Architecture package separation details, directory layouts, and SOLID implementation.

---

## 🇧🇷 Resumo do Projeto (Português)

O Shipwright é um gerenciador visual e extremamente leve de containers e imagens Docker locais, inspirado no Portainer, escrito inteiramente em **Go** e compilado como uma janela desktop nativa Linux de apenas **12MB** de tamanho final.

Ele nasceu para resolver o alto consumo de recursos do Docker Desktop de forma offline, independente e didática. Ele separa as regras de negócio em Clean Architecture aplicando princípios de SOLID (ISP, SRP, DIP) para servir como portfólio limpo de desenvolvimento de software.

Para instruções de configuração, compilação de binários e atalhos de sistema, consulte a pasta de [Documentação (docs)](docs/setup.md).

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
