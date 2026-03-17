# 🚀 Distributed Cache System (Redis-like)

A **production-ready**, **multi-threaded** distributed in-memory cache built in C++ from scratch. This project demonstrates **system design**, **concurrency**, **networking**, and **data structures** at an intermediate-to-advanced level.

<align="center">

```
┌─────────────┐      ┌──────────────┐      ┌──────────────┐      ┌─────────────┐
│  Client 1   │      │  Client 2    │      │  Client 3    │      │  Client 4   │
└──────┬──────┘      └──────┬───────┘      └──────┬───────┘      └──────┬──────┘
       │                    │                     │                     │
       └────────────────────┼─────────────────────┼─────────────────────┘
                            │
                   ┌────────▼────────┐
                   │   TCP Server    │
                   │   (Port 8080)   │
                   └────────┬────────┘
                            │
                  ┌─────────┴─────────┐
                  │                   │
            ┌─────▼──────┐    ┌─────▼──────┐
            │ Thread 1   │    │ Thread 2   │
            └─────┬──────┘    └─────┬──────┘
                  │                │
                  └────────┬───────┘
                           │
                  ┌────────▼────────┐
                  │ ThreadSafeCache │
                  │  (Mutex Locked) │
                  └────────┬────────┘
                           │
                  ┌────────▼────────┐
                  │   LRU Cache     │
                  │ (Linked List +  │
                  │   HashMap)      │
                  └─────────────────┘
```

</align="center">

---

## 📋 Table of Contents

1. [Quick Start](#-quick-start)
2. [Project Structure](#-project-structure)
3. [Component Breakdown](#-component-breakdown)
4. [Architecture & Design](#-architecture--design)
5. [Commands](#-commands)
6. [Usage Guide](#-usage-guide)
7. [Implementation Details](#-implementation-details)
8. [Interview Explanation](#-interview-explanation)
9. [Performance Metrics](#-performance-metrics)
10. [Future Enhancements](#-future-enhancements)

---

## ⚡ Quick Start

### Prerequisites
- **g++ compiler** (C++11 or higher)
- **Linux/macOS/WSL** (for socket programming)
- **Make** (optional, but recommended)

### Build & Run

**Terminal 1 - Start Server:**
```bash
make
make run_server
```

**Terminal 2 - Start Client:**
```bash
make run_client
```

Or manually:
```bash
# Compile
g++ src/server.cpp src/cache.cpp src/lru.cpp -o server -pthread
g++ client/client.cpp -o client

# Run
./server    # Terminal 1
./client    # Terminal 2
```

### Quick Test
```bash
# In client terminal
SET user john
GET user
SET age 25
GET age
INFO
QUIT
```

---

## 📁 Project Structure

```
distributed-cache/
│
├── src/
│   ├── lru.h         ←  LRU cache interface
│   ├── lru.cpp       ←  LRU cache implementation
│   ├── cache.h       ←  Thread-safe wrapper interface
│   ├── cache.cpp     ←  Thread-safe wrapper implementation
│   └── server.cpp    ←  TCP server (main)
│
├── client/
│   └── client.cpp    ←  Interactive client for testing
│
├── Makefile          ←  Build automation
└── README.md         ←  This file
```

---

## 🧠 Component Breakdown

### 1️⃣ **lru.h / lru.cpp** - LRU Cache Core

**What it does:**
- Implements **Least Recently Used (LRU)** eviction policy
- Uses **Doubly Linked List** + **Hash Map** for O(1) operations
- Automatically removes least-used items when full

**Key Data Structures:**
```cpp
list<pair<string, string>> cacheList;  // Doubly linked list
unordered_map<string, iterator> cacheMap;  // Fast lookup
```

**Why this design?**
- **GET**: O(1) - Direct map lookup
- **PUT**: O(1) - Insert at front, remove from back
- **EVICTION**: O(1) - Pop back element

**Real-world analogy:**
Imagine a small desk with only 5 notebooks. When you get a new notebook and the desk is full, you remove the one you haven't used recently. The list keeps track of "most recent" at the front.

---

### 2️⃣ **cache.h / cache.cpp** - Thread-Safe Wrapper

**What it does:**
- Wraps LRU cache with **mutex lock**
- Prevents race conditions in multithreaded environment
- Ensures atomic operations

**Key Concept - Race Condition (Without Locks):**
```
Thread 1: Check if 'a' is LRU
Thread 2: Check if 'a' is LRU  ← Both see same thing
Thread 1: Delete 'a' (evict)
Thread 2: Access 'a'           ← CRASH! Item was deleted
```

**Solution - With Locks:**
```
Thread 1: LOCK → Check 'a' → UNLOCK
Thread 2: LOCK (waits...) 
         → Check 'a' → UNLOCK
```

**Implementation:**
```cpp
lock_guard<mutex> lock(mtx);  // RAII: lock acquired
// Operation happens here
// lock released automatically in destructor
```

---

### 3️⃣ **server.cpp** - TCP Server

**What it does:**
- Listens on **Port 8080**
- Accepts multiple **concurrent clients**
- Spawns **separate thread** for each client
- Parses and executes commands
- Sends responses back

**Processing Flow:**
```
Client A sends "SET user john"
         ↓
Server receives and spawns Thread A
         ↓
Thread A: Parse command → Acquire cache lock
         ↓
Thread A: Execute SET → Release lock
         ↓
Thread A: Send "OK" response
         ↓
Thread A: Wait for next command
```

**Multi-threading Benefit:**
- 100 clients? No problem! 🚀
- Each gets own thread, non-blocking
- Cache access synchronized by mutex

---

### 4️⃣ **client.cpp** - Interactive CLI Client

**What it does:**
- Interactive protocol client
- Sends commands to server
- Displays responses
- User-friendly interface

**Example Session:**
```
$ ./client
✓ Connected to server at 127.0.0.1:8080
Command: SET username alice
Response: OK
Command: GET username
Response: alice
Command: QUIT
```

---

## 🏗️ Architecture & Design

### System Architecture

```
┌─────────────────────────────────────────────────┐
│         APPLICATION LAYER (Client)              │
│  - Interactive CLI interface                    │
│  - Send commands (SET, GET, INFO)               │
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│      NETWORK LAYER (TCP Sockets)                │
│  - Port 8080                                    │
│  - Bidirectional communication                  │
│  - Protocol: TEXT (simple strings)              │
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│    CONCURRENCY LAYER (Multithreading)           │
│  - Thread per client                            │
│  - Independent execution                       │
│  - Mutex-protected shared cache                 │
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│      SYNCHRONIZATION LAYER (Mutex)              │
│  - Ensures thread safety                        │
│  - Prevents race conditions                     │
│  - Atomic cache operations                      │
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│      MEMORY LAYER (Storage)                     │
│  - ThreadSafeCache wrapper                      │
│  - LRU Cache logic                              │
│  - Doubly Linked List + Hash Map                │
└─────────────────────────────────────────────────┘
```

### Data Flow for SET Command

```
Client: "SET username alice"  →
                              →  Server receives
                              →  Spawn thread
                              →  Parse: cmd=SET, key=username, value=alice
                              →  Acquire mutex
                              →  cache.put("username", "alice")
                              →  Release mutex
                              →  Send "OK"
Client receives "OK"  ←
```

### Data Flow for GET Command

```
Client: "GET username"  →
                        →  Server receives
                        →  Spawn thread
                        →  Parse: cmd=GET, key=username
                        →  Acquire mutex
                        →  result = cache.get("username")
                        →  Release mutex
                        →  Send result
Client receives "alice"  ←
```

---

## 📡 Commands

### Supported Commands

| Command | Syntax | Example | Response |
|---------|--------|---------|----------|
| **SET** | `SET <key> <value>` | `SET user john` | `OK` |
| **GET** | `GET <key>` | `GET user` | `john` or `NULL` |
| **INFO** | `INFO` | `INFO` | Server info |
| **HELP** | `HELP` | `HELP` | Display help |
| **QUIT** | `QUIT` | `QUIT` | Exit client |

### Command Examples

```bash
# Store data
SET name alice
SET age 25
SET city seattle

# Retrieve data
GET name           # Returns: alice
GET age            # Returns: 25
GET email          # Returns: NULL (doesn't exist)

# Server info
INFO               # Returns: Cache Server v1.0 | Capacity: 5 | Port: 8080
```

---

## 📖 Usage Guide

### Scenario 1: Basic Operations

**Terminal 1:**
```bash
$ make run_server
🚀 DISTRIBUTED CACHE SERVER STARTED
Available Connections: 5
Cache Capacity: 5
Accepting connections... (Ctrl+C to stop)
```

**Terminal 2:**
```bash
$ ./client
✓ Connected to server at 127.0.0.1:8080
Type HELP for commands

Command: SET username alice
Response: OK

Command: SET counter 1
Response: OK

Command: GET counter
Response: 1

Command: QUIT
```

### Scenario 2: LRU Eviction

```bash
# Capacity is 5 items
Command: SET a 1
Response: OK
Command: SET b 2
Response: OK
Command: SET c 3
Response: OK
Command: SET d 4
Response: OK
Command: SET e 5
Response: OK         # Cache is now FULL [e, d, c, b, a]

Command: SET f 6
Response: OK         # This triggers eviction!
                     # 'a' removed (least recently used)
                     # Cache: [f, e, d, c, b]

Command: GET a
Response: NULL       # 'a' was evicted!
```

### Scenario 3: Multiple Clients

**Terminal 1 - Server:**
```bash
$ ./server
[CLIENT] Connected (Socket: 5)
[REQUEST] 5 > SET user1 alice
[CLIENT] Connected (Socket: 6)
[REQUEST] 6 > SET user2 bob
[REQUEST] 5 > GET user1
```

**Terminal 2 - Client 1:**
```bash
$ ./client
Command: SET user1 alice
Response: OK
Command: GET user1
Response: alice
```

**Terminal 3 - Client 2:**
```bash
$ ./client
Command: SET user2 bob
Response: OK
Command: GET user2
Response: bob
```

**Note:** Both clients execute concurrently without interfering!

---

## 🔍 Implementation Details

### LRU Cache Algorithm

#### GET Operation:

```
function GET(key):
    if key not in map:
        return "NULL"
    
    node = map[key]              # O(1) lookup
    value = node.value
    
    list.remove(node)            # O(1) removal
    list.push_front(key, value)  # O(1) addition
    map[key] = list.front()      # O(1) update
    
    return value
```

**Time: O(1), Space: O(1)**

#### PUT Operation:

```
function PUT(key, value):
    if key in map:
        list.remove(map[key])    # O(1) - Update existing
    
    else if list.size() == capacity:
        last_node = list.back()
        map.remove(last_node.key)    # O(1) - Evict LRU
        list.pop_back()
    
    list.push_front(key, value)  # O(1) - Add new
    map[key] = list.front()      # O(1) - Update map
```

**Time: O(1), Space: O(1) amortized**

### Thread Safety

**Problem:** Multiple threads accessing cache simultaneously

```cpp
// WITHOUT mutex: DANGER!
Thread 1: if (cache not full) ...
Thread 2: if (cache not full) ...  ← Both threads enter
Thread 1: insert item
Thread 2: insert item              ← DATA CORRUPTION!
```

**Solution:** Wrap with mutex

```cpp
// WITH mutex: SAFE!
lock_guard<mutex> lock(mtx);  // Only one thread enters
// Critical section - protected
```

### Socket Programming

**Creating a server socket:**

```cpp
int server_fd = socket(AF_INET, SOCK_STREAM, 0);
                       ^^^^^^  ^^^^^^^^^^^
                       IPv4    TCP protocol

struct sockaddr_in address;
address.sin_family = AF_INET;
address.sin_addr.s_addr = INADDR_ANY;  // Listen on all interfaces
address.sin_port = htons(8080);        // Port 8080

bind(server_fd, (struct sockaddr*)&address, sizeof(address));
listen(server_fd, 5);  // Queue up to 5 pending connections

int new_socket = accept(server_fd, ...);  // Accept incoming connection
```

---

## 💬 Interview Explanation

### The Elevator Pitch (30 seconds)

> "I built a multi-threaded distributed cache system in C++ using TCP sockets. It implements an LRU eviction policy using a doubly linked list and hash map for O(1) operations. To handle concurrency, I used mutex locks ensuring thread safety across multiple client requests. The system follows a client-server architecture and supports real-time command processing."

### Deep Dive (2-3 minutes)

**Interviewer:** "Tell me about your Distributed Cache project."

**You:** 
> "Sure! It's a Redis-like cache built from scratch in C++. Here's the architecture:
> 
> **Core Components:**
> 1. **LRU Cache** - Uses a doubly linked list and hash map. The list maintains Most-Recently-Used items at the front. When we GET/PUT, we move items to the front to mark them as recently used. When capacity is reached, we evict from the back. All O(1)!
> 
> 2. **Thread-Safe Wrapper** - Wraps the LRU cache with a mutex. Why? Multiple clients connect simultaneously, spawning independent threads. Without synchronization, we'd have race conditions. The mutex ensures only one thread accesses the cache at a time, preventing data corruption.
> 
> 3. **TCP Server** - Listens on port 8080, accepts connections, spawns a thread per client. Each thread independently reads commands, parses them, accesses the threadsafe cache, and sends responses back. So 100 concurrent clients? No problem!
> 
> 4. **Protocol** - Simple text-based: 'SET key value' and 'GET key'. Easy to implement, easy to test.
> 
> **Key Design Decisions:**
> - Doubly Linked List + HashMap instead of just a HashMap (O(1) eviction required)
> - Thread pool pattern naturally emerges (server spawns thread per client)
> - Mutex for synchronization (fine-grained locking would be a future enhancement)
> 
> **Complexity:**
> - GET/PUT: O(1) time, O(n) space where n = capacity
> - Scale: Each thread ~1MB, so ~1000 concurrent clients possible on modern hardware"

### Technical Questions You Might Get

**Q: Why Doubly Linked List + HashMap and not just HashMap?**

A: Good question! A HashMap alone can't evict in O(1) time. When we need to evict, we need to find the LRU item efficiently. With a linked list, the LRU item is always at the back—O(1) access. The map gives us O(1) key lookup, and the list gives us O(1) eviction. Together: O(1) GET, PUT, and eviction.

**Q: What happens if two threads access the same key simultaneously?**

A: The mutex ensures mutual exclusion. Only one thread can hold the lock at a time, so operations on the cache are serialized. Thread A acquires lock, completes operation, releases lock. Thread B acquires lock, completes operation. No race conditions.

**Q: What's a race condition example?**

A: Sure! Without locks:
- Thread 1: Check if 'key' is in cache (yes)
- Thread 2: Check if 'key' is in cache (yes)
- Thread 1: Delete 'key' (eviction)
- Thread 2: Access 'key' → CRASH, dangling pointer!

With locks, operations are atomic and non-interruptible.

**Q: How would you scale this to multiple nodes?**

A: Great follow-up! I'd implement:
1. **Consistent Hashing** - Each node owns a hash range, clients route requests
2. **Replication** - Each node replicates to 2-3 others (fault tolerance)
3. **Peer-to-peer communication** - Nodes gossip to stay in sync

This turns it into a true distributed cache.

---

## 📊 Performance Metrics

### Single-threaded Baseline

```
Operation      | Time       | Space
───────────────┼────────────┼──────────
GET hit        | O(1)       | O(1)
GET miss       | O(1)       | O(1)
PUT new        | O(1)       | O(1)
PUT update     | O(1)       | O(1)
PUT eviction   | O(1)       | O(1)

Total Space    | O(n)       | n = capacity
```

### Multi-threaded Overhead

```
Mutex lock/unlock per operation: ~10-100 nanoseconds (modern CPU)
Context switch: ~1-10 microseconds (OS dependent)

Low overhead for reasonable request rate (<1M requests/sec)
```

### Scalability

```
Concurrent Clients | Memory (approx)
──────────────────┼─────────────────
10                | 10 MB
100               | 100 MB
1000              | 1 GB
```

Each thread ~1MB stack on Linux. Cache data negligible for small values.

---

## 🚀 Future Enhancements

### Phase 2: Distributed Capabilities

- [ ] **Consistent Hashing** - Distribute cache across nodes
- [ ] **Replication** - Fault tolerance
- [ ] **Peer Discovery** - Automatic node detection
- [ ] **Data Rebalancing** - Handle node joins/leaves

### Phase 3: Persistence

- [ ] **Disk Snapshots** - Save cache to disk
- [ ] **Write-Ahead Logging** - Durability
- [ ] **Database Integration** - SQL backend fallback

### Phase 4: Performance

- [ ] **Benchmarking** - Latency/throughput metrics
- [ ] **Connection Pooling** - Optimize client connections
- [ ] **Compression** - Reduce network overhead
- [ ] **Read-Write Locks** - Fine-grained concurrency

### Phase 5: DevOps

- [ ] **Docker Containerization** - Easy deployment
- [ ] **Kubernetes Integration** - Cloud-native scaling
- [ ] **Prometheus Metrics** - Monitoring
- [ ] **Admin Dashboard** - Web UI for management

### Phase 6: Advanced Features

- [ ] **Pub/Sub Messaging** - Event system
- [ ] **Transactions** - Multi-command ACID
- [ ] **TTL (Time-To-Live)** - Auto-expiry
- [ ] **Lua Scripting** - Server-side logic

---

## 🧪 Testing

### Manual Testing

```bash
# Start server
./server

# In another terminal, test manually
./client
Command: SET test_key test_value
Command: GET test_key
Command: SET a 1
Command: SET b 2
...
```

### Automated Testing (Makefile)

```bash
make test
```

### Advanced Testing Ideas

```bash
# Stress test (1000 sequential operations)
for i in {1..1000}; do echo "SET key$i value$i"; sleep 0.001; done | ./client

# Concurrency test (multiple clients)
./client &
./client &
./client &
...

# Network test (external machine)
# Edit client.cpp: SERVER_IP = "192.168.1.100"
```

---

## 🔧 Troubleshooting

### "Address already in use" Error

**Problem:** Port 8080 is already bound

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process (on Linux/macOS)
kill -9 <PID>
```

### Compilation Error: "undefined reference to pthread"

**Problem:** Missing `-pthread` flag

**Solution:** Use the Makefile or:
```bash
g++ src/server.cpp src/cache.cpp src/lru.cpp -o server -pthread
```

### Client Cannot Connect

**Problem:** Server not running

**Solution:**
```bash
# Terminal 1
./server

# Terminal 2 (only after server starts)
./client
```

---

## 📚 Learning Resources

### Data Structures
- Understanding Doubly Linked Lists: Time complexity O(1) operations
- Hash Maps: O(1) average lookup and insertion
- LRU implementation patterns

### Networking
- TCP/IP socket programming
- Client-server model
- Network protocols

### Concurrency
- Mutex and synchronization primitives
- Race conditions and thread safety
- RAII pattern in C++

### System Design
- Caching strategies and eviction policies
- Scalability patterns
- Fault tolerance (replication, redundancy)

---

## 📝 License

This project is open source and available as-is for educational and professional use.

---

## 🎯 Summary

| Aspect | Details |
|--------|---------|
| **Language** | C++11 |
| **Architecture** | Client-Server (TCP) |
| **Concurrency** | Multi-threaded (std::thread) |
| **Synchronization** | Mutex locks |
| **Data Structure** | Doubly Linked List + HashMap |
| **Eviction Policy** | LRU (O(1)) |
| **Capacity** | Configurable (default: 5) |
| **Port** | 8080 |
| **Protocol** | Text-based commands |
| **Time Complexity** | GET/PUT: O(1) |
| **Space Complexity** | O(capacity) |

---

## ✨ Credits

Built as a portfolio project to demonstrate:
✔️ System Design  
✔️ Backend Engineering  
✔️ Concurrency & Multithreading  
✔️ Networking (TCP/IP)  
✔️ Data Structures & Algorithms  
✔️ Clean Code & Documentation  

Perfect for FAANG interviews and senior engineering roles! 🚀

---

**Happy Caching!** If you have questions or improvements, feel free to contribute! 💪
