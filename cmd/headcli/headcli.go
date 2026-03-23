package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BenPfeiffer-TX/SatQuorum/internal/types"
	"github.com/moby/moby/client"
)

func main() {
	count := 10
	image := "placeholder"

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
		cmd := scanner.Text()

		switch cmd {
		case "exit":
			fmt.Println("Stopping all nodes...")
			manager.StopAllNodes()
			return
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		manager.StopAllNodes()
		os.Exit(1)
	}
}
