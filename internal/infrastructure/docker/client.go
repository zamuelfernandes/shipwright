package docker

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/moby/moby/client"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// DockerClient implements our domain repositories connecting directly to the Docker API.
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient initializes the official Docker client.
// It leverages standard Docker environment variables (e.g. DOCKER_HOST) via client.FromEnv.
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

// ListContainers retrieves container metadata from the host.
func (d *DockerClient) ListContainers(ctx context.Context) ([]domain.Container, error) {
	result, err := d.cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var domainContainers []domain.Container
	for _, c := range result.Items {
		var composeProj, composeServ string
		if c.Labels != nil {
			composeProj = c.Labels["com.docker.compose.project"]
			composeServ = c.Labels["com.docker.compose.service"]
		}

		domainContainers = append(domainContainers, domain.Container{
			ID:          c.ID,
			Names:       c.Names,
			Image:       c.Image,
			State:       string(c.State),
			Status:      c.Status,
			ComposeProj: composeProj,
			ComposeServ: composeServ,
		})
	}
	return domainContainers, nil
}

// StartContainer starts a stopped container.
func (d *DockerClient) StartContainer(ctx context.Context, id string) error {
	_, err := d.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

// StopContainer stops a running container.
func (d *DockerClient) StopContainer(ctx context.Context, id string) error {
	_, err := d.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

// PruneContainers deletes stopped containers from the system.
func (d *DockerClient) PruneContainers(ctx context.Context) error {
	_, err := d.cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	return err
}

// StreamLogs streams container logs via a channel.
func (d *DockerClient) StreamLogs(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error) {
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
			// Read the 8-byte multiplexing header from Docker stream
			header, err := reader.Peek(8)
			if err != nil {
				if err != io.EOF && ctx.Err() == nil {
					errChan <- err
				}
				return
			}

			isMultiplexed := (header[0] == 1 || header[0] == 2) &&
				header[1] == 0 && header[2] == 0 && header[3] == 0

			if isMultiplexed {
				reader.Discard(8)
				
				// Bytes 4-7 represent payload size in BigEndian uint32
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
				// Raw stream fallback (when TTY is enabled on the container)
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

// StreamStats streams real-time CPU/RAM container telemetry.
func (d *DockerClient) StreamStats(ctx context.Context, id string, statsChan chan<- domain.ContainerStats, errChan chan<- error) {
	result, err := d.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{Stream: true})
	if err != nil {
		errChan <- err
		return
	}
	defer result.Body.Close()

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

			// CPU percentage calculations (equivalent to Docker CLI calculation)
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

			// RAM percentage calculations
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

// StartProject starts all stopped containers in a Compose project concurrently.
func (d *DockerClient) StartProject(ctx context.Context, project string) error {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return err
	}

	var targetContainers []string
	for _, c := range containers {
		if c.ComposeProj == project && c.State != "running" {
			targetContainers = append(targetContainers, c.ID)
		}
	}

	if len(targetContainers) == 0 {
		return nil
	}

	errChan := make(chan error, len(targetContainers))
	var wg sync.WaitGroup

	for _, id := range targetContainers {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := d.StartContainer(ctx, cid); err != nil {
				errChan <- err
			}
		}(id)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// StopProject stops all running containers in a Compose project concurrently.
func (d *DockerClient) StopProject(ctx context.Context, project string) error {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return err
	}

	var targetContainers []string
	for _, c := range containers {
		if c.ComposeProj == project && c.State == "running" {
			targetContainers = append(targetContainers, c.ID)
		}
	}

	if len(targetContainers) == 0 {
		return nil
	}

	errChan := make(chan error, len(targetContainers))
	var wg sync.WaitGroup

	for _, id := range targetContainers {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := d.StopContainer(ctx, cid); err != nil {
				errChan <- err
			}
		}(id)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecContainer runs an interactive session (/bin/bash or fallback /bin/sh) inside a container.
func (d *DockerClient) ExecContainer(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error {
	log.Printf("[Docker] Starting Exec session on container %s...", id)
	
	execConfig := client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          true,
		Cmd:          []string{"/bin/bash"},
	}

	execCreateResp, err := d.cli.ExecCreate(ctx, id, execConfig)
	if err != nil {
		log.Printf("[Docker] Bash not available, trying fallback to sh: %v", err)
		execConfig.Cmd = []string{"/bin/sh"}
		execCreateResp, err = d.cli.ExecCreate(ctx, id, execConfig)
		if err != nil {
			return err
		}
	}

	attachResp, err := d.cli.ExecAttach(ctx, execCreateResp.ID, client.ExecAttachOptions{
		TTY: true,
	})
	if err != nil {
		return err
	}
	defer attachResp.Close()

	doneChan := make(chan struct{})

	// Container output -> stdout
	go func() {
		defer close(doneChan)
		_, _ = io.Copy(stdout, attachResp.Reader)
	}()

	// stdin -> Container input
	go func() {
		_, _ = io.Copy(attachResp.Conn, stdin)
		_ = attachResp.CloseWrite()
	}()

	// Run command asynchronously to prevent blocking/deadlock
	go func() {
		_, _ = d.cli.ExecStart(ctx, execCreateResp.ID, client.ExecStartOptions{
			TTY: true,
		})
	}()

	<-doneChan
	log.Printf("[Docker] Exec session finished.")
	return nil
}

// ListImages retrieves all local Docker images from the host.
func (d *DockerClient) ListImages(ctx context.Context) ([]domain.Image, error) {
	result, err := d.cli.ImageList(ctx, client.ImageListOptions{All: false})
	if err != nil {
		return nil, err
	}

	var domainImages []domain.Image
	for _, img := range result.Items {
		domainImages = append(domainImages, domain.Image{
			ID:      img.ID,
			Tags:    img.RepoTags,
			Size:    img.Size,
			Created: img.Created,
		})
	}
	return domainImages, nil
}
