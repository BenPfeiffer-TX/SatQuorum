# SatQuorum

## Description
SatQuorum is a distributed consensus system that runs multiple satnode containers orchestrated by headcli. Headcli manages the lifecycle of these nodes, starting/stopping them and facilitating message passing for quorum-based operations.

## Scope
- Overview and architecture
- Core implementation details
- Testing and validation
- Deployment instructions

## Prerequisites

### Docker
Docker must be installed and running on your system. You can verify this with:
```bash
docker ps
```

### DOCKER_HOST Configuration (Linux - Docker Desktop)
If you're using **Docker Desktop on Linux**, the Docker socket is typically located at `/home/lars/.docker/desktop/docker.sock`. Headcli will automatically detect this path, but if it fails to find a valid socket, you'll need to set the `DOCKER_HOST` environment variable:

```bash
export DOCKER_HOST=unix:///home/lars/.docker/desktop/docker.sock
```

**To make this persistent across reboots**, add the export command to your shell configuration file (`~/.bashrc`, `~/.zshrc`, or similar):

```bash
echo 'export DOCKER_HOST=unix:///home/lars/.docker/desktop/docker.sock' >> ~/.bashrc
source ~/.bashrc
```

**Alternative: Use Docker Contexts** (recommended for multi-environment setups)
```bash
# Create a custom context for Docker Desktop
docker context create desktop --docker "host=unix:///home/lars/.docker/desktop/docker.sock"
docker context use desktop
```

## Quick Start

### 1. Build the satnode image
```bash
cd /path/to/SatQuorum
docker build -t satnode:latest .
```

### 2. Run headcli
```bash
./headcli
```

Headcli will automatically detect and use your Docker socket if `DOCKER_HOST` is not explicitly set. It searches for common socket locations:
- `/home/lars/.docker/desktop/docker.sock` (Docker Desktop)
- `/var/run/docker.sock` (standard Linux location)
- `/run/user/1000/docker.sock` (user-specific Docker)

### 3. Interactive Commands
Once running, headcli provides an interactive shell:
```bash
> list
NODE ID       PORT       STATUS
satnode-0     54321      running
satnode-1     54322      running
...

> help
Available Commands:
  list    - List all running satnodes with their ports and status
  help    - Show this help message
  exit    - Exit and stop all nodes

> exit
```

## Architecture Overview

### Components

**Headcli (Master Controller)**
- Orchestrates the lifecycle of satnode containers
- Manages container creation, starting, stopping, and cleanup
- Provides interactive command-line interface for monitoring
- Sends test messages to verify node communication

**SatNode (Worker Container)**
- Runs as a Docker container on port 8080
- Implements consensus protocol logic
- Accepts JSON messages with id, payload, and timestamp fields
- Participates in quorum-based decision making

### Parallelization
Headcli uses Go goroutines to parallelize operations:
- **Container spawning**: Up to 50 concurrent container creation attempts (configurable via `MaxConcurrentSpawning` constant)
- **Message verification**: All satnodes receive test messages simultaneously

This design significantly reduces startup time compared to sequential execution.

## Configuration

### Adjusting Concurrency Limit
To modify the maximum number of concurrent container spawns, edit `internal/types/node.go`:
```go
const MaxConcurrentSpawning = 50 // Change this value as needed
```

### Environment Variables
- `DOCKER_HOST`: Docker socket path (auto-detected if not set)
- `PORT`: Satnode listening port (default: 8080, set inside container)

## Testing and Validation

### Running Tests
```bash
go test ./...
```

### Manual Verification
1. Start headcli
2. Run `list` command to see all nodes
3. Verify nodes are running with status "running"
4. Check that test messages were sent during initialization phase

## Troubleshooting

### Docker Client Initialization Fails
- Ensure Docker daemon is running: `sudo systemctl start docker`
- Check DOCKER_HOST environment variable is correct
- Verify your user has permissions to access the Docker socket: `sudo usermod -aG docker $USER`

### Failed to Start Any Nodes
- Verify satnode image exists: `docker images | grep satnode`
- Rebuild the image if needed: `docker build -t satnode:latest .`
- Check available system resources (memory, disk space)

### Permission Denied Errors
```bash
# Add your user to docker group
sudo usermod -aG docker $USER
# Log out and back in for changes to take effect
```

## License
[Add license information]
