# Blueprint do Projeto: Shipwright (O seu "Mini-Portainer" em Go)

Este documento serve como a especificação arquitetural e o guia de desenvolvimento passo a passo para a criação de um gerenciador visual e leve de containers Docker Engine, inspirado no Portainer, escrito inteiramente em **Go**. 

O projeto foi desenhado sob os princípios de **Clean Architecture**, foco absoluto em **eficiência de recursos (CPU/RAM)**, portabilidade (multiplataforma) e **modularidade**, permitindo que um agente de código guie o desenvolvedor no aprendizado prático da linguagem Go e do ecossistema Docker.

---

## 1. Visão Geral e Requisitos de Design

### Objetivos Principais
1. **Eficiência Bruta:** Ocupar o mínimo de memória RAM possível (alvo < 25MB em idle) e inicialização instantânea, superando o Docker Desktop e abordagens baseadas em Electron/Webview tradicionais.
2. **Arquitetura Escalável:** Utilizar Clean Architecture adaptada para a herança idiomática do Go (interfaces implícitas, composição sobre herança).
3. **Multiplataforma Nativa:** O core da aplicação deve compilar perfeitamente para Linux, macOS e Windows.
4. **Interface Gráfica Leve (Web Integrada):** Como o Portainer, o app rodará um servidor HTTP ultra-leve em Go e servirá uma SPA (Single Page Application) minimalista ou usará SSR (Server-Side Rendering) com componentes puros. O binário final do Go irá embutir toda a interface usando a diretiva `//go:embed`, gerando um **único arquivo executável independente**.

---

## 2. Visão Geral da Arquitetura (Clean Architecture em Go)

A estrutura segue a separação de conceitos onde o núcleo do negócio (Docker) não conhece os detalhes de entrega (HTTP/Web).

```
shipwright/
├── cmd/
│   └── server/
│       └── main.go         # Ponto de entrada (Inicialização e Injeção de Dependências)
├── internal/
│   ├── domain/             # Camada de Domínio (Entidades e Contratos de Interfaces)
│   │   ├── container.go
│   │   └── stats.go
│   ├── usecase/            # Camada de Regras de Negócio (Casos de Uso)
│   │   ├── list_containers.go
│   │   ├── monitor_stats.go
│   │   └── manage_lifecycle.go
│   ├── infrastructure/     # Camada de Detalhes de Baixo Nível (Frameworks e Drivers)
│   │   ├── docker/         # Implementação do SDK oficial do Docker
│   │   │   └── client.go
│   │   └── http/           # Servidor Web (Rotas, Handlers, Middleware)
│   │       └── router.go
└── ui/                     # Interface Visual (HTML/CSS/JS puros ou Minimal SPA)
    ├── embed.go            # Arquivo de cola para usar //go:embed
    └── dist/
```

### Detalhe das Camadas:
- **Domain:** Define as estruturas puras (ex: `type Container struct`) e as interfaces que as outras camadas devem satisfazer (Inversão de Dependência). Não importa se usamos o Docker SDK ou mockamos os dados, o domínio permanece idêntico.
- **Usecase:** Orquestra os fluxos de dados. Ex: Busca a lista de containers do repositório, filtra ou formata de acordo com a regra de negócio.
- **Infrastructure/Docker:** Onde a mágica do Docker acontece. Traduz as chamadas de método das nossas interfaces para as chamadas reais do SDK oficial do Docker (`github.com/docker/docker/client`).
- **Infrastructure/HTTP:** Expõe os Casos de Uso via endpoints REST estruturados e gerencia o streaming de dados via Server-Sent Events (SSE) ou WebSockets para manter o consumo leve.

---

## 3. Fases do Roadmap de Desenvolvimento (Iterativo e Incremental)

O Agente de Código deve guiar o desenvolvedor executando e explicando rigorosamente uma fase por vez.

### Fase 1: Fundação do Domínio e SDK Docker (O Essencial)
- Configuração do projeto (`go mod init`).
- Criação das entidades de Domínio para `Container`.
- Implementação da infraestrutura do Docker conectando no Socket local (`/var/run/docker.sock` no Linux/macOS ou `//./pipe/docker_engine` no Windows).
- **Meta:** Conseguir listar todos os containers (Equivalente ao `docker ps -a`) através de um teste unitário ou comando de CLI simples.

### Fase 2: O Servidor HTTP e Embutimento da UI
- Implementação de um roteador HTTP performático usando a biblioteca nativa do Go `net/http`.
- Configuração do mecanismo de `//go:embed` para que os arquivos da pasta `ui/dist/` sejam compilados diretamente para dentro do binário final.
- **Meta:** Acessar `http://localhost:8080` no navegador e ver uma página HTML estática vinda de dentro do binário Go executável.

### Fase 3: Conectando as Pontas (Aba A, D e E)
- Expor os endpoints REST:
  - `GET /api/containers` (Listagem completa)
  - `POST /api/containers/:id/start` (Iniciar)
  - `POST /api/containers/:id/stop` (Pausar)
  - `DELETE /api/containers/prune` (Limpeza do sistema)
- Criar a interface visual mínima em HTML/TailwindCSS (via CDN ou embutido) com botões funcionais para disparar essas ações via `fetch`.

### Fase 4: Streaming em Tempo Real (Aba B e C - Logs e Stats)
- Para evitar sobrecarga por requisições constantes (polling), implementar **Server-Sent Events (SSE)** em Go para fazer o streaming dos dados do Docker diretamente para o browser.
- Endpoint de Streaming de Logs: `GET /api/containers/:id/logs`
- Endpoint de Streaming de Stats (Uso de CPU/RAM): `GET /api/containers/:id/stats`
- **Meta:** Exibir no frontend um gráfico minimalista de CPU/RAM e o terminal de logs rolando em tempo real sem travar o navegador e sem gastar CPU em excesso.

---

## 4. Instruções para o Agente de Código (Diretrizes de Interação)

Ao trabalhar com este projeto, o Agente de Código deve seguir as seguintes premissas:
1. **Didática Orientada a Código:** Antes de despejar um arquivo inteiro, explique o porquê de o Go organizar o código em pacotes (`packages`) e como o sistema de visibilidade (Letra Maiúscula vs. Minúscula) funciona.
2. **Garantia de Performance:** Evite o uso excessivo de reflexão (`reflect`) ou dependências gigantescas. Dê preferência absoluta à biblioteca padrão do Go (`stdlib`) e use apenas o SDK oficial do Docker como pacote externo inicial.
3. **Evitar Concorrência Bloqueante:** Ao implementar o streaming de Logs e Stats (Fase 4), guie o desenvolvedor no uso correto de **Goroutines** e **Channels** para gerenciar o ciclo de vida das conexões abertas de forma segura (tratando o fechamento de contextos com `context.Context`).
4. **Validação Passo a Passo:** Cada arquivo criado deve ter seu funcionamento testado através de um comando `go run` ou de um arquivo de teste simples (`*_test.go`), garantindo que o desenvolvedor aprenda a testar código em Go desde o primeiro dia.

---

Comece solicitando ao desenvolvedor para criar o diretório do projeto e iniciar o módulo Go. Guie-o pela **Fase 1**.
