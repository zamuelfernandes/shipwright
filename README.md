# ⚓ Shipwright

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Engine-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

</div>

Welcome to **Shipwright**! This project is a lightweight, visual manager for Docker containers and local images, inspired by Portainer but designed to be extremely fast and resource-efficient (< 25MB RAM), written entirely in **Go**.

---

### 🛠️ Built With

*   ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) **Go (Golang)** — Standard Library & standard `http.ServeMux` router
*   ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white) **Moby / Docker Go SDK** — Native unix socket engine communication
*   ![WebSockets](https://img.shields.io/badge/WebSockets-010101?style=flat-square&logo=socket.io&logoColor=white) **WebSockets & SSE** — Bi-directional channels for terminal emulators & real-time telemetry streaming
*   ![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=flat-square&logo=html5&logoColor=white) ![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=flat-square&logo=css3&logoColor=white) ![JS](https://img.shields.io/badge/JavaScript-F7DF1E?style=flat-square&logo=javascript&logoColor=black) **Vanilla Stack** — Zero UI framework footprint, embedded directly in Go binary using `//go:embed`

---

### 🌟 The Story & Spirit of the Project

Shipwright was born out of frustration with **Docker Desktop's heavy resource consumption** (often eating up gigabytes of RAM just to list local containers). While researching lighter alternatives, I looked at existing tools like **Portainer** (which requires running another complex container structure) and **lazydocker** (a great CLI tool, but lacking a visual browser-based GUI). 

I wanted something different: a simple, visual dashboard that was ultra-lightweight, lightning fast, and focused precisely on my daily developer workflows. This curiosity sparked a coding experiment to test the limits of rapid, high-quality development in a new stack (Go + Embedded Vanilla WebUI), fully guided and paired with an **AI coding assistant**.

It was built from scratch to be:
*   **Highly Educational:** If you are learning Go, Clean Architecture, or Docker SDK, this codebase is designed to be readable, documented, and pedagogical.
*   **100% Offline & Independent:** There are no external cloud dependencies. The entire web client (HTML/CSS/JS) is compiled directly inside a single Go binary using `//go:embed`.
*   **Open & Reusable:** Feel free to fork it, read it, break it, and use it as a boilerplate for your own Go/Docker tools!

---

## English Version

### 🚀 Key Features

*   **Docker Compose Grouping:** Automatically groups containers belonging to the same Docker Compose project based on labels.
*   **Batch Actions:** Start and stop entire Compose projects concurrently using Go's goroutines and `sync.WaitGroup`.
*   **Interactive Web Terminal:** Open a live shell session (`/bin/bash` or `/bin/sh`) inside any running container directly in your browser, powered by **Xterm.js** and WebSockets.
*   **Inline Monitoring:** The logs and terminal panel expand dynamically directly under the clicked container row in the table, preventing screen clutter.
*   **Persistent Telemetry:** Real-time CPU (%) and RAM (MB/%) consumption graphs stay visible at the top of the monitor panel, even when switching tabs between logs and terminal.
*   **Local Images Dashboard:** Lists all downloaded Docker images, showing short IDs, tags, formatted sizes, and creation dates, positioned at the top of the dashboard.
*   **Collapsible Cards:** Image lists and Compose project cards can be collapsed/minimized to save screen space.

---

### 🏛️ Architecture & Inner Workings

The codebase is structured under **Clean Architecture** principles to separate business logic from delivery mechanisms and external dependencies:

```
shipwright/
├── cmd/
│   └── server/
│       └── main.go         # Entry point (Wires dependencies and starts HTTP server)
├── internal/
│   ├── domain/             # Domain Layer (Core entities and abstract interfaces)
│   │   ├── container.go
│   │   ├── image.go        # Image Entity
│   │   └── stats.go
│   ├── usecase/            # Usecase Layer (Pure business logic)
│   │   ├── list_containers.go
│   │   ├── list_images.go
│   │   ├── manage_lifecycle.go
│   │   ├── start_project.go
│   │   ├── stop_project.go
│   │   ├── exec_container.go
│   │   ├── stream_logs.go
│   │   └── stream_stats.go
│   └── infrastructure/     # Infrastructure Layer (Adapters & External tools)
│       ├── docker/         # Official Moby/Docker SDK v29 Adapter
│       │   └── client.go
│       └── http/           # Standard HTTP Router (Go 1.22+ routes) & WebSockets
│           └── router.go
└── ui/                     # UI Web Assets (Vanilla HTML/CSS/JS & Xterm.js)
    ├── embed.go            # Embedding directives (//go:embed)
    └── dist/
        └── index.html      # Single page dashboard application
```

---

### 🛠️ How to Run & Build

#### Prerequisites
- **Go 1.22+**
- **Docker Engine** running locally.
- On Linux, make sure your user belongs to the `docker` group so the binary can read `/var/run/docker.sock` without `sudo`:
  ```bash
  sudo usermod -aG docker $USER
  ```
  *(Note: You will need to log out and log back in for this change to take effect)*.

#### Running in Development Mode
1. Copy the environment configuration template:
   ```bash
   cp .env.example .env
   ```
2. (Optional) Customize variables in `.env` (like changing the execution `PORT` or custom `DOCKER_HOST`).
3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```
4. Open your browser at:
   👉 **[http://localhost:8080](http://localhost:8080)** (or the port defined in your `.env` file).

#### Building a Production Executable
Go compiles down to a single, self-contained binary. To build a highly optimized version (removing debug symbols to get a file size under 10MB):
```bash
go build -ldflags="-s -w" -o shipwright cmd/server/main.go
```
To run it:
```bash
./shipwright
```

---

## Versão em Português

Bem-vindo ao **Shipwright**! Este projeto é um gerenciador visual e extremamente leve de containers e imagens Docker locais, inspirado no Portainer, escrito inteiramente em **Go**.

---

### 🛠️ Tecnologias Utilizadas

*   ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) **Go (Golang)** — Biblioteca Padrão e roteador nativo `http.ServeMux`
*   ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white) **Moby / SDK Go do Docker** — Comunicação direta com o socket Unix do host local
*   ![WebSockets](https://img.shields.io/badge/WebSockets-010101?style=flat-square&logo=socket.io&logoColor=white) **WebSockets & SSE** — Canais bidirecionais para terminal e fluxo SSE de métricas/logs em tempo real
*   ![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=flat-square&logo=html5&logoColor=white) ![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=flat-square&logo=css3&logoColor=white) ![JS](https://img.shields.io/badge/JavaScript-F7DF1E?style=flat-square&logo=javascript&logoColor=black) **Stack Vanilla** — UI pura em HTML5, CSS e JS, embutida no binário final via `//go:embed`

---

### 🌟 A História e o Espírito do Projeto

O Shipwright nasceu da minha **insatisfação com o alto consumo de recursos do Docker Desktop** (que frequentemente consome gigabytes de memória RAM apenas para listar containers locais). Ao pesquisar alternativas mais leves, examinei soluções existentes como o **Portainer** (que exige rodar toda uma estrutura complexa de containers dedicada) e o **lazydocker** (uma excelente ferramenta de terminal CLI, mas sem uma interface gráfica visual acessível pelo navegador).

Eu queria algo diferente: um dashboard simples, visual, extremamente leve, rápido e focado exatamente nas minhas principais necessidades diárias de desenvolvimento. Essa curiosidade deu início a um experimento de código para testar os limites do desenvolvimento ágil em uma nova stack (Go + WebUI Vanilla embutida), pareado com um **assistente de programação IA**.

Ele foi feito sob medida para ser:
*   **Altamente Didático:** Se você está aprendendo Go, Clean Architecture ou o SDK do Docker, esta base de código foi desenhada para ser legível, bem documentada e pedagógica.
*   **100% Offline e Independente:** Sem dependências externas de nuvem. Toda a interface web é compilada dentro de um único binário Go usando `//go:embed`.
*   **Aberto e Reutilizável:** Fique à vontade para fazer forks, ler, quebrar e usar como base para suas próprias ferramentas em Go e Docker!

---

### 🚀 Funcionalidades Principais

*   **Agrupamento por Docker Compose:** Agrupa visualmente os containers que pertencem ao mesmo projeto Compose através de labels.
*   **Comandos em Lote do Compose:** Inicia e para projetos Compose inteiros de forma concorrente utilizando goroutines e `sync.WaitGroup` no Go.
*   **Terminal Web Interativo:** Abre uma sessão de shell ativa (`/bin/bash` ou `/bin/sh`) dentro de qualquer container em execução diretamente no navegador usando **Xterm.js** e WebSockets.
*   **Monitoramento Inline:** O painel de logs/terminal expande-se dinamicamente logo abaixo da linha do container selecionado na tabela, evitando poluição visual.
*   **Telemetria Persistente:** Gráficos dinâmicos de CPU (%) e RAM (MB/%) no topo do painel de monitoramento, que continuam visíveis mesmo quando você muda de aba para rodar comandos no terminal.
*   **Painel de Imagens Locais:** Lista todas as imagens Docker salvas na máquina local no topo do dashboard, mostrando tags, tamanhos formatados e data de criação.
*   **Cards Minimizáveis:** Possibilidade de recolher/minimizar os painéis de Imagens e Compose para economizar espaço de tela.

---

### 🏛️ Arquitetura e Funcionamento

O projeto foi estruturado seguindo os princípios de **Clean Architecture**, mantendo as regras de negócio isoladas das ferramentas externas (como roteador HTTP e o cliente Docker):

*   **Domain (Domínio):** Contém as entidades principais (`Container`, `Image`) e a assinatura de interfaces como contratos.
*   **Usecase (Casos de Uso):** Contém a lógica de negócio pura (listar containers, iniciar projetos, conectar terminal).
*   **Infrastructure (Infraestrutura):** Adaptadores e ferramentas externas. O [client.go](file:///home/samuel/Workspace/Pessoal/Go/shipwright/internal/infrastructure/docker/client.go) faz a ponte com o SDK Oficial do Docker, e o [router.go](file:///home/samuel/Workspace/Pessoal/Go/shipwright/internal/infrastructure/http/router.go) expõe a API HTTP e conexões WebSocket.
*   **UI:** Arquivos de front-end (HTML/CSS/JS puros e Xterm.js) embutidos no binário final.

---

### 🛠️ Como Executar e Compilar

#### Pré-requisitos
- **Go 1.22+**
- **Docker Engine** ativo na máquina.
- No Linux/Ubuntu, garanta que seu usuário está no grupo do docker:
  ```bash
  sudo usermod -aG docker $USER
  ```

#### Rodando em Modo de Desenvolvimento
1. Copie o modelo de variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
2. (Opcional) Edite as variáveis no `.env` (como alterar a porta de escuta `PORT` ou customizar o `DOCKER_HOST`).
3. Rode o servidor:
   ```bash
   go run cmd/server/main.go
   ```
4. Acesse no seu navegador:
   👉 **[http://localhost:8080](http://localhost:8080)** (ou a porta configurada no seu `.env`).

#### Compilando para Produção
Para compilar um binário único e super otimizado (geralmente menor do que 10MB):
```bash
go build -ldflags="-s -w" -o shipwright cmd/server/main.go
```
Para executar o binário gerado:
```bash
./shipwright
```

---

## 📄 License / Licença

This project is licensed under the **MIT License**. / Este projeto está sob a licença **MIT**. See the [LICENSE](LICENSE) file for details.
