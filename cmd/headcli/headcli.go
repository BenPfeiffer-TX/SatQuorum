package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BenPfeiffer-TX/SatQuorum/internal/types"
	"github.com/moby/moby/client"
)

func main() {
	count := 10
	image := "satnode:latest"

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Docker client: %v\n", err)
		os.Exit(1)
	}

	manager, err := types.NewSatNodeManager(dockerClient, count, image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create SatNodeManager: %v\n", err)
		os.Exit(1)
	}

	err = manager.StartNodes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start nodes: %v\n", err)
		manager.StopAllNodes()
		os.Exit(1)
	}

	fmt.Println("\n=== Verification Phase ===")
	manager.SendTestMessages()
	fmt.Println("\n=== Interactive Mode ===")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived %s, cleaning up...\n", sig)
		manager.StopAllNodes()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())

		switch cmd {
		case "list":
			nodes := manager.ListNodes()
			if len(nodes) == 0 {
				fmt.Println("No running satnodes")
			} else {
				fmt.Printf("%-15s %-10s %s\n", "NODE ID", "PORT", "STATUS")
				fmt.Println(strings.Repeat("-", 35))
				for _, node := range nodes {
					fmt.Printf("%-15s %-10d %s\n", node["nodeID"], node["port"], node["status"])
				}
			}
		case "help":
			fmt.Println("Available Commands:")
			fmt.Println("  list    - List all running satnodes with their ports and status")
			fmt.Println("  help    - Show this help message")
			fmt.Println("  exit    - Exit and stop all nodes")
		case "exit":
			fmt.Println("Stopping all nodes...")
			manager.StopAllNodes()
			return
		default:
			fmt.Printf("Unknown command: %s (type 'help' for available commands)\n", cmd)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		manager.StopAllNodes()
		os.Exit(1)
	}
}
