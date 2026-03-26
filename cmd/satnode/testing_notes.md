Testing Instructions (Once satnode Docker Image is Ready)
Step 1: Create a satnode Docker image
# From the SatQuorum directory, create a Dockerfile for satnode
cat > cmd/satnode/Dockerfile << 'EOF'
FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o satnode ./cmd/satnode
EXPOSE 8080
CMD ["./satnode"]
EOF
# Build and tag the image (use "placeholder" or your preferred name)
docker build -t placeholder:latest -f cmd/satnode/Dockerfile .
Step 2: Run headcli with custom parameters
# Basic usage - start 10 satnodes using default "placeholder" image
./headcli
# Start a specific number of nodes
./headcli -count=5
# Use a different Docker image name
./headcli -image=my-satnode:v1.0
# Combine options
./headcli -count=3 -image=my-satnode:latest
Expected Output:
Started satnode-0 on port 32768
Started satnode-1 on port 32769
Started satnode-2 on port 32770
...
=== Verification Phase ===
satnode-0: OK
satnode-1: OK
satnode-2: OK
...
=== Interactive Mode ===
Available Commands:
  list    - List all running satnodes with their ports and status
  help    - Show this help message
  exit    - Exit and stop all nodes

Example usage:
> list
NODE ID           PORT       STATUS
satnode-0         32768      running
satnode-1         32769      running
...
> help
Available Commands:
  list    - List all running satnodes with their ports and status
  help    - Show this help message
  exit    - Exit and stop all nodes
> exit
Stopping all nodes...
Step 3: Verify containers are running
# Check running satnodes
docker ps --filter "name=satnode" --format "{{.Names}} - Port: {{.Ports}}"
# Inspect a specific container's port mapping
docker inspect satnode-0 --format '{{(index .NetworkSettings.Ports "8080/tcp")[0].HostPort}}'
Step 4: Test message delivery (manual)
# Send a test JSON message to a running satnode
curl -X POST http://localhost:<port>/ \
  -H "Content-Type: application/json" \
  -d '{"id":"test","payload":"hello world","timestamp":"2026-03-23T12:00:00Z"}'
# Check satnode logs
docker logs satnode-0
Step 5: Cleanup (if containers persist)
# Stop and remove all satnodes manually
docker stop $(docker ps -q --filter "name=satnode")
docker rm -f $(docker ps -aq --filter "name=satnode")
---
