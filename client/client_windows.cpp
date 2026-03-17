/**
 * ============================================================================
 * DISTRIBUTED CACHE CLIENT (Windows Version)
 * ============================================================================
 */

#include <iostream>
#include <string>
#include <cstring>
#include <winsock2.h>
#include <ws2tcpip.h>
#pragma comment(lib, "ws2_32.lib")

#define SERVER_IP "127.0.0.1"
#define SERVER_PORT 8080
#define BUFFER_SIZE 1024

using namespace std;

/**
 * Main client function
 */
int main() {
    cout << "\n";
    cout << "╔═════════════════════════════════════════════════════════════╗" << endl;
    cout << "║                 💻 DISTRIBUTED CACHE CLIENT                 ║" << endl;
    cout << "╚═════════════════════════════════════════════════════════════╝" << endl;
    cout << endl;

    // Initialize Winsock
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        cerr << "ERROR: WSAStartup failed" << endl;
        return 1;
    }

    // Create client socket
    SOCKET client_socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (client_socket == INVALID_SOCKET) {
        cerr << "ERROR: Socket creation failed" << endl;
        WSACleanup();
        return 1;
    }

    // Configure server address
    sockaddr_in server_addr;
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(SERVER_PORT);
    server_addr.sin_addr.s_addr = inet_addr(SERVER_IP);

    // Connect to server
    if (connect(client_socket, (sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
        cerr << "ERROR: Connection failed. Make sure server is running on " << SERVER_IP << ":" << SERVER_PORT << endl;
        closesocket(client_socket);
        WSACleanup();
        return 1;
    }

    cout << "✓ Connected to server at " << SERVER_IP << ":" << SERVER_PORT << endl;
    cout << "Type HELP for commands or QUIT to exit" << endl;
    cout << endl;

    char buffer[BUFFER_SIZE];
    bool running = true;

    while (running) {
        cout << "Command: ";
        cout.flush();

        string command;
        getline(cin, command);

        if (command.empty()) {
            continue;
        }

        command += "\n";
        send(client_socket, command.c_str(), (int)command.length(), 0);

        memset(buffer, 0, BUFFER_SIZE);
        int bytes_read = recv(client_socket, buffer, BUFFER_SIZE - 1, 0);

        if (bytes_read <= 0) {
            cout << "ERROR: Connection lost" << endl;
            running = false;
            break;
        }

        string response(buffer);
        if (!response.empty() && response.back() == '\n') {
            response.pop_back();
        }

        cout << "Response: " << response << endl;

        if (command.find("QUIT") == 0 || command.find("quit") == 0) {
            running = false;
        }

        cout << endl;
    }

    closesocket(client_socket);
    WSACleanup();
    cout << "✓ Disconnected from server" << endl;

    return 0;
}