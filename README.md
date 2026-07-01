# ⚓ Shipwright

O Shipwright é um gerenciador visual e extremamente leve de containers Docker Engine, inspirado no Portainer, escrito inteiramente em **Go**.

A ferramenta foi projetada com base nos princípios de **Clean Architecture**, foco absoluto em **eficiência de recursos (CPU/RAM)** e **portabilidade**. Ela integra uma interface gráfica minimalista embutida diretamente dentro de um único arquivo executável usando o recurso `//go:embed` do Go, rodando de forma 100% local e offline.

---

## 🚀 Funcionalidades Principais

### Gerenciamento de Projetos e Containers (V2)
- **Agrupamento por Docker Compose:** Agrupa visualmente os containers que pertencem ao mesmo projeto do Docker Compose, facilitando o gerenciamento.
- **Comandos em Lote do Compose:** Inicia e para projetos do Compose de forma concorrente utilizando goroutines e `sync.WaitGroup` no Go (equivalente aos comandos `docker compose up` e `down` em lote).
- **Listagem e Controle de Containers Avulsos:** Detecta todos os containers locais ativos e parados através de conexão direta com o Docker Engine Socket UNIX.
- **Controle de Ciclo de Vida Individual:** Botões para iniciar, parar e executar limpeza completa de containers inativos (Docker Prune).

### Monitoramento e Interatividade em Tempo Real (V2)
- **Terminal Interativo no Navegador (Exec):** Abre um console web interativo e funcional emulando um terminal VT100 (via **Xterm.js**) utilizando WebSockets para conexões bidirecionais diretas com o shell do container (`/bin/bash` ou `/bin/sh`).
- **Streaming de Logs (Tempo Real):** Conexão contínua usando Server-Sent Events (SSE) para ler logs do container sem necessidade de fazer requisições constantes (polling).
- **Telemetria ao Vivo (Stats):** Medição em tempo real de consumo de CPU (%) e uso de memória RAM (MB e %) com barras de progresso dinâmicas.

---

## 🏛️ Arquitetura do Projeto

O projeto segue os princípios de **Clean Architecture** para Go, separando as regras de negócio de detalhes de infraestrutura (como o SDK do Docker ou o servidor HTTP):

```
shipwright/
├── cmd/
│   └── server/
│       └── main.go         # Ponto de entrada (Inicialização e Injeção de Dependências)
├── internal/
│   ├── domain/             # Camada de Domínio (Entidades e Interfaces Abstratas)
│   │   ├── container.go
│   │   └── stats.go
│   ├── usecase/            # Camada de Regras de Negócio (Casos de Uso)
│   │   ├── list_containers.go
│   │   ├── manage_lifecycle.go
│   │   ├── start_project.go
│   │   ├── stop_project.go
│   │   ├── exec_container.go
│   │   ├── stream_logs.go
│   │   └── stream_stats.go
│   └── infrastructure/     # Camada de Infraestrutura (Drivers e Adaptadores)
│       ├── docker/         # Cliente do SDK Oficial do Moby/Docker v29
│       │   └── client.go
│       └── http/           # Servidor e Roteador HTTP Nativo (Go 1.22+ e WebSockets)
│           └── router.go
└── ui/                     # Interface Web (HTML/CSS/JS puros e Xterm.js)
    ├── embed.go            # Configuração do //go:embed
    └── dist/
        └── index.html      # Dashboard do Console
```

---

## 🛠️ Como Executar Localmente

### Pré-requisitos
- **Go** (versão 1.22+ recomendada, testado com Go 1.26.0)
- **Docker Engine** ativo na máquina local.
- No Linux/Ubuntu, certifique-se de que seu usuário pertence ao grupo `docker` para que o programa possa ler o socket `/var/run/docker.sock` sem precisar de privilégios de `sudo`:
  ```bash
  sudo usermod -aG docker $USER
  ```
  *(Nota: Se você acabou de rodar o comando acima, é necessário fazer logout e login no sistema para aplicar as permissões)*.

### Executando em Modo de Desenvolvimento
Clone o repositório e execute a aplicação diretamente com o comando do Go:
```bash
# Executa e compila o servidor local
go run cmd/server/main.go
```

Acesse o console no seu navegador:
👉 **[http://localhost:8080](http://localhost:8080)**

---

## 📄 Licença

Este projeto está sob os termos da licença **MIT**. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
