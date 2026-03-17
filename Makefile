# ============================================================================
# MAKEFILE - Build Configuration for Distributed Cache
# ============================================================================
# 
# Purpose:
# --------
# Automates compilation of all source files
# 
# Targets:
# --------
# make              - Compile server and client
# make server       - Compile server only
# make client       - Compile client only
# make clean        - Remove compiled files
# make run_server   - Compile and run server
# make run_client   - Compile and run client
# make test         - Run automated tests
# 
# Compiler Flags:
# ---------------
# -std=c++11        - Use C++11 standard
# -pthread          - Enable multithreading
# -Wall             - Show all warnings
# -O2               - Optimization level 2
# 
# ============================================================================

# Compiler
CXX = g++

# Compiler flags
CXXFLAGS = -std=c++11 -Wall -O2 -pthread

# Source files
SRC_SERVER = src/server.cpp src/cache.cpp src/lru.cpp
SRC_CLIENT = client/client.cpp

# Output executables
BINDIR = bin
EXEC_SERVER = $(BINDIR)/server
EXEC_CLIENT = $(BINDIR)/client

# Default target: Build all
all: $(EXEC_SERVER) $(EXEC_CLIENT)

# ============================================================================
# SERVER BUILD
# ============================================================================
$(EXEC_SERVER): $(SRC_SERVER)
	@mkdir -p $(BINDIR)
	@echo "🔨 Compiling server..."
	$(CXX) $(CXXFLAGS) $(SRC_SERVER) -o $(EXEC_SERVER)
	@echo "✓ Server compiled: $(EXEC_SERVER)"

# ============================================================================
# CLIENT BUILD
# ============================================================================
$(EXEC_CLIENT): $(SRC_CLIENT)
	@mkdir -p $(BINDIR)
	@echo "🔨 Compiling client..."
	$(CXX) $(CXXFLAGS) $(SRC_CLIENT) -o $(EXEC_CLIENT)
	@echo "✓ Client compiled: $(EXEC_CLIENT)"

# ============================================================================
# CONVENIENCE TARGETS
# ============================================================================

# Compile server only
server: $(EXEC_SERVER)

# Compile client only
client: $(EXEC_CLIENT)

# Run server
run_server: $(EXEC_SERVER)
	@echo "\n🚀 Starting server..."
	./$(EXEC_SERVER)

# Run client (requires separate terminal)
run_client: $(EXEC_CLIENT)
	@echo "\n💻 Starting client..."
	./$(EXEC_CLIENT)

# Clean compiled files
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(EXEC_SERVER) $(EXEC_CLIENT)
	@echo "✓ Clean complete"

# Rebuild everything
rebuild: clean all
	@echo "✓ Rebuild complete"

# ============================================================================
# TEST TARGET
# ============================================================================

# Simple automated test
test: all
	@echo "\n════════════════════════════════════════════════════════════"
	@echo "  AUTOMATED TESTS"
	@echo "════════════════════════════════════════════════════════════\n"
	@echo "Test 1: Starting server in background..."
	./$(EXEC_SERVER) &
	@sleep 2
	@echo "\nTest 2: Running client tests..."
	@(echo "SET key1 value1"; echo "GET key1"; echo "SET key2 value2"; echo "SET key3 value3"; echo "GET key1"; echo "GET key3"; echo "QUIT") | ./$(EXEC_CLIENT)
	@echo "\n✓ Tests completed"

# ============================================================================
# HELP TARGET
# ============================================================================

help:
	@echo "╔═════════════════════════════════════════════════════════════╗"
	@echo "║            DISTRIBUTED CACHE - MAKEFILE HELP               ║"
	@echo "╚═════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available targets:"
	@echo "─────────────────────────────────────────────────────────────"
	@echo "  make              - Compile both server and client"
	@echo "  make server       - Compile server only"
	@echo "  make client       - Compile client only"
	@echo "  make run_server   - Build and run server"
	@echo "  make run_client   - Build and run client"
	@echo "  make rebuild      - Clean and rebuild"
	@echo "  make clean        - Remove compiled files"
	@echo "  make test         - Run automated tests"
	@echo "  make help         - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "─────────────────────────────────────────────────────────────"
	@echo "  Terminal 1: make run_server"
	@echo "  Terminal 2: make run_client"
	@echo ""

# Mark these as non-file targets
.PHONY: all server client run_server run_client clean rebuild test help
