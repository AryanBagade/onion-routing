# 🎯 Onion Network - Quick Demo Guide

**For demonstrating your already-deployed Azure onion network**

## 🚀 Prerequisites (Already Done)
- ✅ 4 Azure VMs running
- ✅ Binaries deployed to all servers
- ✅ Ports opened (9000, 8080, 8081, 8082)
- ✅ SSH access configured

## 📊 Your Network Infrastructure

| Server | IP Address | Location | Role | Port |
|--------|------------|----------|------|------|
| Directory | `172.191.95.78` | East US | Node Registry | 9000 |
| Guard | `172.201.12.43` | West Europe | Entry Point | 8080 |
| Relay | `68.218.3.154` | Australia East | Middle Hop | 8081 |
| Exit | `172.191.84.146` | East US | Exit Point | 8082 |

---

# 🎬 Demo Instructions

## Step 1: Start Directory Server (FIRST)

```bash
ssh onion@172.191.95.78
chmod +x onion-network-real-ip
./onion-network-real-ip -mode=directory -port=9000
```

**Wait for this output:**
```
Starting directory server on port 9000
Directory server listening on port 9000
```

**✅ KEEP THIS TERMINAL OPEN**

---

## Step 2: Start Guard Node (Europe)

**Open new terminal:**
```bash
ssh onion@172.201.12.43
chmod +x onion-network-real-ip
./onion-network-real-ip -mode=node -type=guard -port=8080
```

**Wait for this output:**
```
Starting guard node node_abc123 on port 8080
Node IP: 🌍 172.201.12.43 (Azure West Europe)
Registered with directory server: 200 OK
```

**✅ KEEP THIS TERMINAL OPEN**

---

## Step 3: Start Relay Node (Australia)

**Open new terminal:**
```bash
ssh onion@68.218.3.154
chmod +x onion-network-real-ip
./onion-network-real-ip -mode=node -type=relay -port=8081
```

**Wait for this output:**
```
Starting relay node node_def456 on port 8081
Node IP: 🌍 68.218.3.154 (Azure Australia East)
Registered with directory server: 200 OK
```

**✅ KEEP THIS TERMINAL OPEN**

---

## Step 4: Start Exit Node (USA)

**Open new terminal:**
```bash
ssh onion@172.191.84.146
chmod +x onion-network-real-ip
./onion-network-real-ip -mode=node -type=exit -port=8082
```

**Wait for this output:**
```
Starting exit node node_ghi789 on port 8082
Node IP: 🌍 172.191.84.146 (Azure East US)
Registered with directory server: 200 OK
```

**✅ KEEP THIS TERMINAL OPEN**

---

## Step 5: Test the Network

**Open new local terminal:**
```bash
cd /Users/aryan/Developer/TorOnionRouting/YourNetwork/onion-network
./onion-network -mode=client
```

### Create Circuit
```
onion> create
```

**Expected output:**
```
Created circuit circuit_abc123: node_guard -> node_relay -> node_exit
Created circuit: circuit_abc123
Path: node_guard -> node_relay -> node_exit
```

### Make Anonymous Request
```
onion> request https://httpbin.org/ip
```

**Expected output:**
```
Making request to https://httpbin.org/ip via circuit circuit_abc123
🧅 Creating onion encryption layers...
🔒 Encrypting request with 3 layers...
📦 Encrypted packet size: 910 bytes
✅ Encrypted request sent successfully!
```

### Exit Client
```
onion> quit
```

---

# 🔥 Live Demo Flow

## What Audience Will See

### 1. **Guard Terminal (Europe):**
```
[GUARD] Received 910 bytes of data
[GUARD] 🔓 Successfully decrypted layer, forwarding to relay
[GUARD] ✅ Successfully forwarded to RELAY
```

### 2. **Relay Terminal (Australia):**
```
[RELAY] Received 626 bytes of data
[RELAY] 🔓 Successfully decrypted layer, forwarding to exit
[RELAY] ✅ Successfully forwarded to EXIT
```

### 3. **Exit Terminal (USA):**
```
[EXIT] Received 342 bytes of data
[EXIT] 🔓 Successfully decrypted final layer
[EXIT] 🌐 Making REAL request to: https://httpbin.org/ip
[EXIT] ✅ SUCCESS! Got 30 bytes from https://httpbin.org/ip
[EXIT] 📊 Response Status: 200 OK
[EXIT] 📋 Response Preview: {"origin": "172.191.84.146"}
```

## 🎯 The Magic Moment

**Point out to audience:**
- **Your real IP**: Hidden
- **Website sees**: `172.191.84.146` (USA exit node)
- **Traffic path**: You → Europe → Australia → USA → Internet
- **Encryption**: 3 layers of RSA, each node peels one layer

---

# 🛠️ Quick Troubleshooting

## If Connection Fails

### 1. Check All Nodes Are Running
```bash
# Quick check - all should return process info
ssh onion@172.191.95.78 "ps aux | grep onion"
ssh onion@172.201.12.43 "ps aux | grep onion"  
ssh onion@68.218.3.154 "ps aux | grep onion"
ssh onion@172.191.84.146 "ps aux | grep onion"
```

### 2. Restart in Order
```bash
# Kill all first
ssh onion@172.191.95.78 "pkill onion-network"
ssh onion@172.201.12.43 "pkill onion-network"
ssh onion@68.218.3.154 "pkill onion-network"
ssh onion@172.191.84.146 "pkill onion-network"

# Then restart following Steps 1-4 above
```

### 3. Create Fresh Circuit
```bash
./onion-network -mode=client
create  # Creates new circuit with fresh keys
request https://httpbin.org/ip
```

---

# 🎪 Demo Script for Presentation

## Opening (30 seconds)
> "I'm going to demonstrate a real onion routing network that provides complete internet anonymity. This system routes your traffic through servers in Europe, Australia, and USA, with each hop adding a layer of encryption."

## Show Infrastructure (30 seconds)
> "Here are 4 terminals connected to real servers running on Microsoft Azure across 3 continents. I'll start each node and you'll see them register with our directory server."

**Start all nodes following Steps 1-4**

## Demonstrate Anonymity (60 seconds)
> "Now I'll make a request to check my IP address. Watch how the traffic flows through all three nodes..."

**Run client test - point to each terminal as traffic flows**

## The Reveal (30 seconds)
> "The website thinks the request came from our USA server at IP 172.191.84.146, but I'm actually here in [your location]. My real IP is completely hidden through multi-layer encryption and global routing."

## Summary (30 seconds)
> "This demonstrates how onion routing provides true anonymity - not just hiding your IP, but creating an untraceable path through multiple encrypted hops across the globe."

---

# ⚡ Emergency Commands

## Quick Start All (If Already Configured)
```bash
# Run in separate terminals
ssh onion@172.191.95.78 "./onion-network-real-ip -mode=directory -port=9000" &
sleep 5
ssh onion@172.201.12.43 "./onion-network-real-ip -mode=node -type=guard -port=8080" &
ssh onion@68.218.3.154 "./onion-network-real-ip -mode=node -type=relay -port=8081" &  
ssh onion@172.191.84.146 "./onion-network-real-ip -mode=node -type=exit -port=8082" &
```

## Stop All
```bash
ssh onion@172.191.95.78 "pkill onion-network"
ssh onion@172.201.12.43 "pkill onion-network"
ssh onion@68.218.3.154 "pkill onion-network"
ssh onion@172.191.84.146 "pkill onion-network"
```

---

**🎯 Your network demonstrates real-world onion routing with actual global infrastructure and multi-layer encryption. Perfect for showing how internet anonymity actually works!** 🧅🌍