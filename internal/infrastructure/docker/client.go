package docker

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/moby/moby/client"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// DockerClient é a implementação do nosso repositório de domínio conectada diretamente com a API do Docker.
// Em Flutter/Dart, isso seria a classe 'DockerRepositoryImpl' que implementa (implements) 'ContainerRepository'.
// No Go, não usamos a palavra chave 'implements'. Como o DockerClient possui o método
// 'ListContainers(ctx context.Context) ([]domain.Container, error)' com a assinatura EXATA,
// ele satisfaz implicitamente a interface ContainerRepository.
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient é o construtor da nossa infraestrutura. Ele inicia o cliente oficial do Docker.
//
// Variáveis de ambiente:
// O método 'client.FromEnv' faz com que o SDK procure onde o Docker está rodando no host.
// No Ubuntu, o padrão é o socket Unix 'unix:///var/run/docker.sock'.
// 'client.WithAPIVersionNegotiation()' faz o SDK negociar automaticamente a versão da API HTTP
// suportada pelo seu Docker instalado para evitar erros de incompatibilidade de versão.
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

// ListContainers busca os containers do host usando o SDK do Docker.
// Retorna uma lista formatada da nossa struct de domínio e um erro (se houver).
func (d *DockerClient) ListContainers(ctx context.Context) ([]domain.Container, error) {
	// 'ContainerList' é o equivalente do SDK para o comando 'docker ps'.
	// No SDK do Moby v29+, usamos 'client.ContainerListOptions' e ele retorna um 'ContainerListResult'.
	result, err := d.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var domainContainers []domain.Container
	// Mapeamos os items retornados no resultado
	for _, c := range result.Items {
		// Mapeamento (Data Mapping):
		// O SDK do Docker retorna structs muito ricas e cheias de metadados internos da API.
		// Nós mapeamos isso para a nossa entidade simplificada de Domínio.
		// Em Flutter, isso é similar a um mapper do Data Layer (DTO -> Entity) para isolar
		// o resto do aplicativo de mudanças nas propriedades cruas da API externa.
		domainContainers = append(domainContainers, domain.Container{
			ID:     c.ID,
			Names:  c.Names,
			Image:  c.Image,
			State:  string(c.State),
			Status: c.Status,
		})
	}
	return domainContainers, nil
}

// StartContainer inicia um container com o ID fornecido.
// No SDK do Moby, precisamos passar um 'client.ContainerStartOptions{}' vazio.
func (d *DockerClient) StartContainer(ctx context.Context, id string) error {
	_, err := d.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

// StopContainer para um container em execução com o ID fornecido.
// No SDK do Moby, passamos um 'client.ContainerStopOptions{}' vazio (que usa o tempo limite padrão do host).
func (d *DockerClient) StopContainer(ctx context.Context, id string) error {
	_, err := d.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

// PruneContainers remove todos os containers inativos/parados do sistema.
// É o equivalente programático do comando 'docker container prune -f'.
func (d *DockerClient) PruneContainers(ctx context.Context) error {
	_, err := d.cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	return err
}

// StreamLogs lê continuamente o fluxo de logs do Docker e os envia via channel.
func (d *DockerClient) StreamLogs(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error) {
	// Chamamos ContainerLogs no SDK do Moby
	result, err := d.cli.ContainerLogs(ctx, id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	})
	if err != nil {
		errChan <- err
		return
	}
	defer result.Close()

	reader := bufio.NewReader(result)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Tentamos ler o cabeçalho de multiplexação de 8 bytes do Docker
			header, err := reader.Peek(8)
			if err != nil {
				// Fim de arquivo (EOF) ou conexão encerrada
				if err != io.EOF && ctx.Err() == nil {
					errChan <- err
				}
				return
			}

			// O formato multiplexado do Docker sempre inicia com Stream Type (1 ou 2)
			// seguido de 3 bytes de padding zerados.
			isMultiplexed := (header[0] == 1 || header[0] == 2) &&
				header[1] == 0 && header[2] == 0 && header[3] == 0

			if isMultiplexed {
				// Descartamos os 8 bytes do cabeçalho que acabamos de validar
				reader.Discard(8)
				
				// Os bytes 4-7 em BigEndian definem o tamanho da mensagem (uint32)
				size := binary.BigEndian.Uint32(header[4:8])
				payload := make([]byte, size)
				
				_, err = io.ReadFull(reader, payload)
				if err != nil {
					if err != io.EOF && ctx.Err() == nil {
						errChan <- err
					}
					return
				}
				select {
				case <-ctx.Done():
					return
				case logsChan <- string(payload):
				}
			} else {
				// Se não for multiplexado (ex: container alocou TTY), lemos como texto puro linha por linha
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					select {
					case <-ctx.Done():
						return
					case logsChan <- scanner.Text():
					}
				}
				if err := scanner.Err(); err != nil && ctx.Err() == nil {
					errChan <- err
				}
				return
			}
		}
	}
}

// Estruturas auxiliares internas para fazer o parse dos dados de stats do Docker.
// Isso evita que precisemos importar pacotes complexos e mantemos o consumo leve.
type cpuUsage struct {
	TotalUsage uint64 `json:"total_usage"`
}

type cpuStats struct {
	CPUUsage       cpuUsage `json:"cpu_usage"`
	SystemCPUUsage uint64   `json:"system_cpu_usage"`
	OnlineCPUs     uint32   `json:"online_cpus"`
}

type memoryStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

type dockerStatsJSON struct {
	CPUStats    cpuStats    `json:"cpu_stats"`
	PreCPUStats cpuStats    `json:"precpu_stats"`
	MemoryStats memoryStats `json:"memory_stats"`
}

// StreamStats lê continuamente as estatísticas de telemetria do Docker e calcula CPU % e RAM %.
func (d *DockerClient) StreamStats(ctx context.Context, id string, statsChan chan<- domain.ContainerStats, errChan chan<- error) {
	// ContainerStats do SDK oficial. Habilitamos Stream: true.
	result, err := d.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{Stream: true})
	if err != nil {
		errChan <- err
		return
	}
	defer result.Body.Close()

	// Decodificador de fluxo JSON (JSON streaming decoder) - extremamente eficiente em CPU/RAM
	decoder := json.NewDecoder(result.Body)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var raw dockerStatsJSON
			if err := decoder.Decode(&raw); err != nil {
				if err != io.EOF && ctx.Err() == nil {
					errChan <- err
				}
				return
			}

			// Cálculo da porcentagem de CPU (Padrão do Docker CLI)
			cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage) - float64(raw.PreCPUStats.CPUUsage.TotalUsage)
			systemDelta := float64(raw.CPUStats.SystemCPUUsage) - float64(raw.PreCPUStats.SystemCPUUsage)
			onlineCPUs := raw.CPUStats.OnlineCPUs
			if onlineCPUs == 0 {
				onlineCPUs = 1
			}

			cpuPercent := 0.0
			if systemDelta > 0.0 && cpuDelta > 0.0 {
				cpuPercent = (cpuDelta / systemDelta) * float64(onlineCPUs) * 100.0
			}

			// Cálculo da porcentagem de Memória RAM
			memPercent := 0.0
			if raw.MemoryStats.Limit > 0 {
				memPercent = (float64(raw.MemoryStats.Usage) / float64(raw.MemoryStats.Limit)) * 100.0
			}

			select {
			case <-ctx.Done():
				return
			case statsChan <- domain.ContainerStats{
				CPUPercent:       cpuPercent,
				MemoryUsageBytes: raw.MemoryStats.Usage,
				MemoryLimitBytes: raw.MemoryStats.Limit,
				MemoryPercent:    memPercent,
			}:
			}
		}
	}
}
