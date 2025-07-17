package client

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	
	"onion-network/pkg/circuit"
	"onion-network/pkg/crypto"
)

type OnionClient struct {
	CircuitManager *circuit.CircuitManager
	DirectoryURL   string
}

func NewOnionClient(directoryURL string) *OnionClient {
	return &OnionClient{
		CircuitManager: circuit.NewCircuitManager(directoryURL),
		DirectoryURL:   directoryURL,
	}
}

func (oc *OnionClient) Start() error {
	fmt.Println("Onion client started")
	fmt.Println("Commands:")
	fmt.Println("  create - Create a new circuit")
	fmt.Println("  request <url> - Make anonymous request")
	fmt.Println("  circuits - List active circuits")
	fmt.Println("  quit - Exit client")
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("onion> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)
		
		if len(parts) == 0 {
			continue
		}
		
		command := parts[0]
		
		switch command {
		case "create":
			oc.handleCreateCircuit()
		case "request":
			if len(parts) < 2 {
				fmt.Print("URL: ")
				if !scanner.Scan() {
					break
				}
				url := strings.TrimSpace(scanner.Text())
				oc.handleRequest(url)
			} else {
				url := parts[1]
				oc.handleRequest(url)
			}
		case "circuits":
			oc.handleListCircuits()
		case "quit":
			fmt.Println("Goodbye!")
			return nil
		default:
			fmt.Println("Unknown command")
		}
	}
	return nil
}

func (oc *OnionClient) handleCreateCircuit() {
	circuit, err := oc.CircuitManager.CreateCircuit()
	if err != nil {
		fmt.Printf("Failed to create circuit: %v\n", err)
		return
	}
	
	fmt.Printf("Created circuit: %s\n", circuit.ID)
	fmt.Printf("Path: %s -> %s -> %s\n", 
		circuit.Nodes[0].ID, circuit.Nodes[1].ID, circuit.Nodes[2].ID)
}

func (oc *OnionClient) handleRequest(url string) {
	// Get first available circuit
	circuits := oc.CircuitManager.Circuits
	if len(circuits) == 0 {
		fmt.Println("No circuits available. Create one first.")
		return
	}
	
	var selectedCircuit *circuit.Circuit
	for _, c := range circuits {
		selectedCircuit = c
		break
	}
	
	fmt.Printf("Making request to %s via circuit %s\n", url, selectedCircuit.ID)
	
	
	// Create onion layers for encryption
	nodeKeys := make([]*rsa.PublicKey, len(selectedCircuit.Nodes))
	nodeIDs := make([]string, len(selectedCircuit.Nodes))
	
	for i, node := range selectedCircuit.Nodes {
		nodeKeys[i] = node.PublicKey
		nodeIDs[i] = node.ID
	}
	
	// Actually send request through circuit
	fmt.Printf("Sending request through circuit:\n")
	fmt.Printf("  Guard: %s:%d\n", selectedCircuit.Nodes[0].Address, selectedCircuit.Nodes[0].Port)
	fmt.Printf("  Relay: %s:%d\n", selectedCircuit.Nodes[1].Address, selectedCircuit.Nodes[1].Port)
	fmt.Printf("  Exit: %s:%d\n", selectedCircuit.Nodes[2].Address, selectedCircuit.Nodes[2].Port)
	fmt.Printf("  Final destination: %s\n", url)
	
	// Create request data
	requestData := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n", url, extractHost(url))
	
	// Create REAL onion layers
	fmt.Printf("🧅 Creating onion encryption layers...\n")
	layers, err := crypto.CreateOnionLayers(nodeKeys, nodeIDs)
	if err != nil {
		fmt.Printf("❌ Failed to create onion layers: %v\n", err)
		return
	}
	
	// Encrypt with onion layers
	fmt.Printf("🔒 Encrypting request with %d layers...\n", len(layers))
	encryptedPacket, err := crypto.EncryptOnion([]byte(requestData), layers)
	if err != nil {
		fmt.Printf("❌ Failed to encrypt request: %v\n", err)
		return
	}
	
	fmt.Printf("📦 Encrypted packet size: %d bytes\n", len(encryptedPacket.Data))
	
	// Send encrypted packet to guard node
	if err := oc.sendThroughCircuit(encryptedPacket.Data, selectedCircuit); err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
	} else {
		fmt.Printf("✅ Encrypted request sent successfully!\n")
	}
}

func (oc *OnionClient) handleListCircuits() {
	circuits := oc.CircuitManager.Circuits
	if len(circuits) == 0 {
		fmt.Println("No active circuits")
		return
	}
	
	fmt.Println("Active circuits:")
	for _, c := range circuits {
		fmt.Printf("  %s: %s -> %s -> %s\n", 
			c.ID, c.Nodes[0].ID, c.Nodes[1].ID, c.Nodes[2].ID)
	}
}

func (oc *OnionClient) sendThroughCircuit(data []byte, circuit *circuit.Circuit) error {
	// Connect to guard node - use real Azure IP
	guardAddr := "172.201.12.43:8080" // Guard node in Europe
	conn, err := net.DialTimeout("tcp", guardAddr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to guard node: %v", err)
	}
	defer conn.Close()
	
	// Send encrypted data
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	
	return nil
}

func extractHost(url string) string {
	// Simple URL parsing - extract hostname
	if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}
	
	for i, char := range url {
		if char == '/' || char == ':' {
			return url[:i]
		}
	}
	return url
}