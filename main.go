package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	
	"onion-network/pkg/client"
	"onion-network/pkg/directory"
	"onion-network/pkg/node"
)

func main() {
	var mode = flag.String("mode", "node", "Mode: node, client, or directory")
	var port = flag.Int("port", 8080, "Port to listen on")
	var nodeType = flag.String("type", "relay", "Node type: guard, relay, or exit")
	flag.Parse()

	switch *mode {
	case "node":
		var nodeTypeEnum node.NodeType
		switch *nodeType {
		case "guard":
			nodeTypeEnum = node.Guard
		case "relay":
			nodeTypeEnum = node.Relay
		case "exit":
			nodeTypeEnum = node.Exit
		default:
			fmt.Println("Invalid node type. Use: guard, relay, or exit")
			os.Exit(1)
		}
		
		n, err := node.NewNode(nodeTypeEnum, "localhost", *port)
		if err != nil {
			log.Fatal("Failed to create node:", err)
		}
		
		fmt.Printf("Starting %s node %s on port %d\n", *nodeType, n.ID, *port)
		fmt.Printf("Node IP: %s\n", n.GetVirtualIP())
		if err := n.Start(); err != nil {
			log.Fatal("Failed to start node:", err)
		}
		
	case "client":
		onionClient := client.NewOnionClient("http://172.191.95.78:9000")
		fmt.Println("Starting onion client")
		if err := onionClient.Start(); err != nil {
			log.Fatal("Failed to start client:", err)
		}
		
	case "directory":
		ds := directory.NewDirectoryServer(*port)
		fmt.Printf("Starting directory server on port %d\n", *port)
		if err := ds.Start(); err != nil {
			log.Fatal("Failed to start directory server:", err)
		}
		
	default:
		fmt.Println("Invalid mode. Use: node, client, or directory")
		os.Exit(1)
	}
}