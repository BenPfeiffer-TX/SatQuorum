package types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type NodeInfo struct {
	ID   string
	Port int
}

type SatNodeManager struct {
	Nodes  map[string]NodeInfo
	Client *client.Client
	Image  string
	Count  int
}

func NewSatNodeManager(dockerClient *client.Client, count int, image string) (*SatNodeManager, error) {
	if _, err := dockerClient.Ping(context.Background(), client.PingOptions{}); err != nil {
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	manager := &SatNodeManager{
		Nodes:  make(map[string]NodeInfo),
		Client: dockerClient,
		Image:  image,
		Count:  count,
	}

	return manager, nil
}

func (m *SatNodeManager) StartNodes() error {
	ctx := context.Background()

	for i := 0; i < m.Count; i++ {
		nodeID := fmt.Sprintf("satnode-%d", i)

		portBinding := network.PortBinding{HostPort: ""}
		portMap := network.PortMap{network.MustParsePort("8080/tcp"): []network.PortBinding{portBinding}}

		createOpts := client.ContainerCreateOptions{
			Name: nodeID,
			Config: &container.Config{
				Image: m.Image,
				Env:   []string{"PORT=8080"},
			},
			HostConfig: &container.HostConfig{
				PortBindings: portMap,
			},
		}

		resp, err := m.Client.ContainerCreate(ctx, createOpts)
		if err != nil {
			fmt.Printf("Failed to create container %s: %v\n", nodeID, err)
			continue
		}

		_, err = m.Client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
		if err != nil {
			fmt.Printf("Failed to start container %s: %v\n", nodeID, err)
			m.Client.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})
			continue
		}

		time.Sleep(1 * time.Second)

		result, err := m.Client.ContainerInspect(ctx, resp.ID, client.ContainerInspectOptions{})
		if err != nil {
			fmt.Printf("Failed to inspect container %s: %v\n", nodeID, err)
			m.Client.ContainerStop(ctx, resp.ID, client.ContainerStopOptions{})
			m.Client.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})
			continue
		}

		port := ""
		if portMap, ok := result.Container.NetworkSettings.Ports[network.MustParsePort("8080/tcp")]; ok && len(portMap) > 0 {
			port = portMap[0].HostPort
		}

		if port == "" {
			fmt.Printf("Failed to get host port for container %s\n", nodeID)
			m.Client.ContainerStop(ctx, resp.ID, client.ContainerStopOptions{})
			m.Client.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})
			continue
		}

		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err != nil {
			fmt.Printf("Failed to parse host port for container %s\n", nodeID)
			m.Client.ContainerStop(ctx, resp.ID, client.ContainerStopOptions{})
			m.Client.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})
			continue
		}

		m.Nodes[nodeID] = NodeInfo{
			ID:   nodeID,
			Port: portInt,
		}

		fmt.Printf("Started %s on port %d\n", nodeID, m.Nodes[nodeID].Port)
	}

	if len(m.Nodes) == 0 {
		return fmt.Errorf("failed to start any nodes")
	}

	return nil
}

func (m *SatNodeManager) GetNodeStatus(nodeID string) string {
	ctx := context.Background()

	container, err := m.Client.ContainerInspect(ctx, nodeID, client.ContainerInspectOptions{})
	if err != nil {
		return "unknown"
	}

	state := container.Container.State
	if state == nil {
		return "unknown"
	}

	switch state.Status {
	case "running":
		return "running"
	case "paused":
		return "paused"
	case "restarting":
		return "restarting"
	case "exited", "dead":
		return "stopped"
	default:
		return "unknown"
	}
}

func (m *SatNodeManager) ListNodes() []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(m.Nodes))

	for nodeID, info := range m.Nodes {
		status := m.GetNodeStatus(nodeID)
		result = append(result, map[string]interface{}{
			"nodeID": nodeID,
			"port":   info.Port,
			"status": status,
		})
	}

	return result
}

func (m *SatNodeManager) SendTestMessages() error {
	for _, node := range m.Nodes {
		nodeID := node.ID
		port := node.Port

		msg := map[string]string{
			"id":        "verification",
			"payload":   "test message from headcli",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		jsonData, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("%s: FAILED to marshal JSON\n", nodeID)
			continue
		}

		resp, err := http.Post(fmt.Sprintf("http://localhost:%d/", port), "application/json", bytes.NewReader(jsonData))
		if err != nil {
			fmt.Printf("%s: FAILED - %v\n", nodeID, err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("%s: OK\n", nodeID)
		} else {
			fmt.Printf("%s: FAILED - status %d\n", nodeID, resp.StatusCode)
		}
	}

	return nil
}

func (m *SatNodeManager) StopAllNodes() {
	ctx := context.Background()

	for _, node := range m.Nodes {
		nodeID := node.ID

		_, err := m.Client.ContainerStop(ctx, nodeID, client.ContainerStopOptions{})
		if err != nil {
			fmt.Printf("Failed to stop %s: %v\n", nodeID, err)
		}

		_, err = m.Client.ContainerRemove(ctx, nodeID, client.ContainerRemoveOptions{Force: true})
		if err != nil {
			fmt.Printf("Failed to remove %s: %v\n", nodeID, err)
		}
	}

	m.Nodes = make(map[string]NodeInfo)
}
