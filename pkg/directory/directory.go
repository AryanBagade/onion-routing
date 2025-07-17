package directory

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type NodeInfo struct {
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	Address   string           `json:"address"`
	Port      int              `json:"port"`
	PublicKey *rsa.PublicKey   `json:"public_key"`
	LastSeen  time.Time        `json:"last_seen"`
}

type DirectoryServer struct {
	Port  int
	Nodes map[string]*NodeInfo
	mutex sync.RWMutex
}

func NewDirectoryServer(port int) *DirectoryServer {
	return &DirectoryServer{
		Port:  port,
		Nodes: make(map[string]*NodeInfo),
	}
}

func (ds *DirectoryServer) Start() error {
	http.HandleFunc("/register", ds.handleRegister)
	http.HandleFunc("/nodes", ds.handleGetNodes)
	http.HandleFunc("/nodes/", ds.handleGetNodesByType)
	
	fmt.Printf("Directory server listening on port %d\n", ds.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", ds.Port), nil)
}

func (ds *DirectoryServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var node NodeInfo
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	node.LastSeen = time.Now()

	ds.mutex.Lock()
	ds.Nodes[node.ID] = &node
	ds.mutex.Unlock()

	fmt.Printf("Registered %s node: %s at %s:%d\n", node.Type, node.ID, node.Address, node.Port)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func (ds *DirectoryServer) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ds.mutex.RLock()
	nodes := make([]*NodeInfo, 0, len(ds.Nodes))
	for _, node := range ds.Nodes {
		if time.Since(node.LastSeen) < 5*time.Minute {
			nodes = append(nodes, node)
		}
	}
	ds.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (ds *DirectoryServer) handleGetNodesByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeType := r.URL.Path[len("/nodes/"):]
	if nodeType == "" {
		http.Error(w, "Node type required", http.StatusBadRequest)
		return
	}

	ds.mutex.RLock()
	var nodes []*NodeInfo
	for _, node := range ds.Nodes {
		if node.Type == nodeType && time.Since(node.LastSeen) < 5*time.Minute {
			nodes = append(nodes, node)
		}
	}
	ds.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}