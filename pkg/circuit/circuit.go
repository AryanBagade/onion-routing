package circuit

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type NodeInfo struct {
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	Address   string           `json:"address"`
	Port      int              `json:"port"`
	PublicKey *rsa.PublicKey   `json:"public_key"`
}

type Circuit struct {
	ID    string
	Nodes []NodeInfo
	Path  []string
	mutex sync.RWMutex
}

type CircuitManager struct {
	DirectoryURL string
	Circuits     map[string]*Circuit
	mutex        sync.RWMutex
}

func NewCircuitManager(directoryURL string) *CircuitManager {
	if directoryURL == "" {
		directoryURL = "http://172.191.95.78:9000"
	}
	return &CircuitManager{
		DirectoryURL: directoryURL,
		Circuits:     make(map[string]*Circuit),
	}
}

func (cm *CircuitManager) CreateCircuit() (*Circuit, error) {
	// Get available nodes from directory
	guardNodes, err := cm.getNodesByType("guard")
	if err != nil {
		return nil, fmt.Errorf("failed to get guard nodes: %v", err)
	}
	
	relayNodes, err := cm.getNodesByType("relay")
	if err != nil {
		return nil, fmt.Errorf("failed to get relay nodes: %v", err)
	}
	
	exitNodes, err := cm.getNodesByType("exit")
	if err != nil {
		return nil, fmt.Errorf("failed to get exit nodes: %v", err)
	}

	if len(guardNodes) == 0 || len(relayNodes) == 0 || len(exitNodes) == 0 {
		return nil, errors.New("insufficient nodes for circuit creation")
	}

	// Select one node of each type (simple selection for now)
	circuit := &Circuit{
		ID: generateCircuitID(),
		Nodes: []NodeInfo{
			guardNodes[0], // Guard node
			relayNodes[0], // Relay node
			exitNodes[0],  // Exit node
		},
		Path: []string{guardNodes[0].ID, relayNodes[0].ID, exitNodes[0].ID},
	}

	cm.mutex.Lock()
	cm.Circuits[circuit.ID] = circuit
	cm.mutex.Unlock()

	fmt.Printf("Created circuit %s: %s -> %s -> %s\n", 
		circuit.ID, guardNodes[0].ID, relayNodes[0].ID, exitNodes[0].ID)

	return circuit, nil
}

func (cm *CircuitManager) getNodesByType(nodeType string) ([]NodeInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/nodes/%s", cm.DirectoryURL, nodeType))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (cm *CircuitManager) GetCircuit(circuitID string) (*Circuit, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	circuit, exists := cm.Circuits[circuitID]
	return circuit, exists
}

func (cm *CircuitManager) DestroyCircuit(circuitID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	delete(cm.Circuits, circuitID)
	fmt.Printf("Destroyed circuit %s\n", circuitID)
}

func generateCircuitID() string {
	return "circuit_" + randomString(12)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)] // Simple deterministic selection for now
	}
	return string(b)
}