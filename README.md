# ⚓ Shipwright

O Shipwright é um gerenciador visual e extremamente leve de containers Docker Engine, inspirado no Portainer, escrito inteiramente em **Go**.

A ferramenta foi projetada com base nos princípios de **Clean Architecture**, foco absoluto em **eficiência de recursos (CPU/RAM)** e **portabilidade**. Ela integra uma interface gráfica minimalista embutida diretamente dentro de um único arquivo executável usando o recurso `//go:embed` do Go, rodando de forma 100% local e offline.

---

## 🚀 Funcionalidades da V1

- **Listagem de Containers:** Detecta todos os containers locais ativos e parados através de conexão direta com o Docker Engine Socket UNIX.
- **Controle de Ciclo de Vida:** Botões funcionais na tela para iniciar, parar e executar limpeza completa de containers inativos (Docker Prune).
- **Streaming de Logs (Tempo Real):** Conexão contínua usando Server-Sent Events (SSE) para ler logs do container sem necessidade de fazer requisições constantes (polling).
- **Telemetria ao Vivo (Stats):** Medição em tempo real de consumo de CPU (%) e uso de memória RAM (MB e %) com barras de progresso dinâmicas.
- **Binário Único Independente:** Toda a interface web (HTML, CSS e JS puros) é compilada e integrada dentro do executável em Go.

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
│   │   ├── stream_logs.go
│   │   └── stream_stats.go
│   └── infrastructure/     # Camada de Infraestrutura (Drivers e Adaptadores)
│       ├── docker/         # Cliente do SDK Oficial do Moby/Docker v29
│       │   └── client.go
│       └── http/           # Servidor e Roteador HTTP Nativo (Go 1.22+)
│           └── router.go
└── ui/                     # Interface Web (HTML/CSS/JS puros)
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
