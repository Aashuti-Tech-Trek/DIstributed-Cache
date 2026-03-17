/**
 * ============================================================================
 * DISTRIBUTED CACHE SERVER (Windows Version)
 * ============================================================================
 */

#include <iostream>
#include <thread>
#include <vector>
#include <string>
#include <cstring>
#include <winsock2.h>
#include <ws2tcpip.h>
#pragma comment(lib, "ws2_32.lib")
#include "cache.h"

#define PORT 8080
#define MAX_CONNECTIONS 5
#define BUFFER_SIZE 1024

using namespace std;

/**
 * Handle individual client connection
 */
void handle_client(SOCKET client_socket, ThreadSafeCache& cache, int client_id) {
    char buffer[BUFFER_SIZE];
    bool running = true;

    cout << "[CLIENT] Connected (Socket: " << client_socket << ", ID: " << client_id << ")" << endl;

    while (running) {
        memset(buffer, 0, BUFFER_SIZE);

        int bytes_read = recv(client_socket, buffer, BUFFER_SIZE - 1, 0);
        if (bytes_read <= 0) {
            cout << "[CLIENT] Disconnected (Socket: " << client_socket << ", ID: " << client_id << ")" << endl;
            break;
        }

        string command(buffer);
        command.erase(command.find_last_not_of("\r\n") + 1);

        cout << "[REQUEST] " << client_id << " > " << command << endl;

        string response;
        if (command.empty()) {
            response = "ERROR: Empty command";
        } else {
            vector<string> tokens;
            size_t pos = 0;
            string token;
            while ((pos = command.find(" ")) != string::npos) {
                token = command.substr(0, pos);
                if (!token.empty()) tokens.push_back(token);
                command.erase(0, pos + 1);
            }
            if (!command.empty()) tokens.push_back(command);

            if (tokens.empty()) {
                response = "ERROR: Invalid command";
            } else {
                string cmd = tokens[0];

                if (cmd == "SET" && tokens.size() == 3) {
                    cache.put(tokens[1], tokens[2]);
                    response = "OK";
                } else if (cmd == "GET" && tokens.size() == 2) {
                    response = cache.get(tokens[1]);
                } else if (cmd == "INFO" && tokens.size() == 1) {
                    response = "Cache Server v1.0 | Capacity: 5 | Port: " + to_string(PORT);
                } else if (cmd == "HELP" && tokens.size() == 1) {
                    response = "Commands: SET <key> <value>, GET <key>, INFO, HELP, QUIT";
                } else if (cmd == "QUIT" && tokens.size() == 1) {
                    response = "Goodbye!";
                    running = false;
                } else {
                    response = "ERROR: Unknown command or invalid syntax. Type HELP for commands.";
                }
            }
        }

        response += "\n";
        send(client_socket, response.c_str(), (int)response.length(), 0);
    }

    closesocket(client_socket);
}

/**
 * Main server function
 */
int main() {
    cout << "\n";
    cout << "╔═════════════════════════════════════════════════════════════╗" << endl;
    cout << "║                 🚀 DISTRIBUTED CACHE SERVER                 ║" << endl;
    cout << "╚═════════════════════════════════════════════════════════════╝" << endl;
    cout << endl;

    // Initialize Winsock
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        cerr << "ERROR: WSAStartup failed" << endl;
        return 1;
    }

    ThreadSafeCache cache(5);

    // Create server socket
    SOCKET server_fd = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (server_fd == INVALID_SOCKET) {
        cerr << "ERROR: Socket creation failed" << endl;
        WSACleanup();
        return 1;
    }

    // Configure server address
    sockaddr_in address;
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    // Bind socket to port
    if (bind(server_fd, (sockaddr*)&address, sizeof(address)) == SOCKET_ERROR) {
        cerr << "ERROR: Bind failed" << endl;
        closesocket(server_fd);
        WSACleanup();
        return 1;
    }

    // Listen for connections
    if (listen(server_fd, MAX_CONNECTIONS) == SOCKET_ERROR) {
        cerr << "ERROR: Listen failed" << endl;
        closesocket(server_fd);
        WSACleanup();
        return 1;
    }

    cout << "✓ Server started successfully" << endl;
    cout << "✓ Listening on port " << PORT << endl;
    cout << "✓ Cache capacity: 5 items" << endl;
    cout << "✓ Max connections: " << MAX_CONNECTIONS << endl;
    cout << endl;
    cout << "Accepting connections... (Ctrl+C to stop)" << endl;
    cout << endl;

    vector<thread> client_threads;
    int client_counter = 1;

    while (true) {
        sockaddr_in client_addr;
        int client_len = sizeof(client_addr);
        SOCKET client_socket = accept(server_fd, (sockaddr*)&client_addr, &client_len);

        if (client_socket == INVALID_SOCKET) {
            cerr << "ERROR: Accept failed" << endl;
            continue;
        }

        client_threads.push_back(thread(handle_client, client_socket, ref(cache), client_counter++));

        for (auto it = client_threads.begin(); it != client_threads.end(); ) {
            if (it->joinable()) {
                it->join();
                it = client_threads.erase(it);
            } else {
                ++it;
            }
        }
    }

    closesocket(server_fd);
    WSACleanup();
    return 0;
}